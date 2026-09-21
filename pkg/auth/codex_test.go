/*
 *  Copyright (c) 2026 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

package auth

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"
	"time"
)

func TestUnit_CodexDeviceAuth_Login(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(codexLoginHandler(t))
	defer server.Close()

	authenticator, err := NewCodexDeviceAuth(
		WithCodexIssuer(server.URL),
		WithCodexPollTimeout(time.Second),
	)
	if err != nil {
		t.Fatal(err)
	}

	var shown CodexDeviceCode
	session, err := authenticator.Login(t.Context(), func(code CodexDeviceCode) error {
		shown = code
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if shown.VerificationURL != server.URL+"/codex/device" || shown.UserCode != "ABCD-EFGH" {
		t.Fatalf("shown code = %#v", shown)
	}
	if got := session.Tokens(); got != (CodexTokens{AccessToken: "access", RefreshToken: "refresh", IDToken: "id"}) {
		t.Fatalf("tokens = %#v", got)
	}

	headers, err := session.HeaderProvider()(t.Context(), RequestMeta{Method: http.MethodGet})
	if err != nil {
		t.Fatal(err)
	}
	if got := headers.Get("Authorization"); got != "Bearer access" {
		t.Fatalf("authorization = %q", got)
	}
}

func codexLoginHandler(t *testing.T) http.Handler {
	t.Helper()
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		switch request.URL.Path {
		case "/api/accounts/deviceauth/usercode":
			handleCodexUserCode(t, writer, request)
		case "/api/accounts/deviceauth/token":
			handleCodexDeviceToken(t, writer, request)
		case "/oauth/token":
			handleCodexTokenExchange(t, writer, request)
		default:
			http.NotFound(writer, request)
		}
	})
}

func handleCodexUserCode(t *testing.T, writer http.ResponseWriter, request *http.Request) {
	t.Helper()
	var input codexUserCodeRequest
	if err := json.NewDecoder(request.Body).Decode(&input); err != nil {
		t.Errorf("decode device request: %v", err)
	}
	if input.ClientID != DefaultCodexClientID {
		t.Errorf("client id = %q, want %q", input.ClientID, DefaultCodexClientID)
	}
	writeCodexJSON(t, writer, `{"device_auth_id":"device-1","usercode":"ABCD-EFGH","interval":"1"}`)
}

func handleCodexDeviceToken(t *testing.T, writer http.ResponseWriter, request *http.Request) {
	t.Helper()
	var input codexDeviceTokenRequest
	if err := json.NewDecoder(request.Body).Decode(&input); err != nil {
		t.Errorf("decode poll request: %v", err)
	}
	if input.DeviceAuthID != "device-1" || input.UserCode != "ABCD-EFGH" {
		t.Errorf("poll request = %#v", input)
	}
	writeCodexJSON(t, writer, `{"authorization_code":"auth-code","code_verifier":"verifier"}`)
}

func handleCodexTokenExchange(t *testing.T, writer http.ResponseWriter, request *http.Request) {
	t.Helper()
	body, err := io.ReadAll(request.Body)
	if err != nil {
		t.Errorf("read token request: %v", err)
		return
	}
	values, err := url.ParseQuery(string(body))
	if err != nil {
		t.Errorf("parse token request: %v", err)
	}
	for key, want := range map[string]string{
		"grant_type":    "authorization_code",
		"code":          "auth-code",
		"redirect_uri":  "http://" + request.Host + "/deviceauth/callback",
		"client_id":     DefaultCodexClientID,
		"code_verifier": "verifier",
	} {
		if got := values.Get(key); got != want {
			t.Errorf("token form %s = %q, want %q", key, got, want)
		}
	}
	writeCodexJSON(t, writer, `{"access_token":"access","refresh_token":"refresh","id_token":"id","expires_in":3600}`)
}

func TestUnit_CodexDeviceAuth_PollCancellation(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		switch request.URL.Path {
		case "/api/accounts/deviceauth/usercode":
			writeCodexJSON(t, writer, `{"device_auth_id":"device-1","user_code":"ABCD-EFGH","interval":60}`)
		case "/api/accounts/deviceauth/token":
			cancel()
			writer.WriteHeader(http.StatusNotFound)
		default:
			http.NotFound(writer, request)
		}
	}))
	defer server.Close()

	authenticator, err := NewCodexDeviceAuth(WithCodexIssuer(server.URL))
	if err != nil {
		t.Fatal(err)
	}
	deviceCode, err := authenticator.RequestDeviceCode(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	_, err = authenticator.Complete(ctx, deviceCode)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("error = %v, want context cancellation", err)
	}
}

func TestUnit_CodexSession_RefreshesExpiringToken(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.URL.Path != "/oauth/token" {
			http.NotFound(writer, request)
			return
		}
		body, err := io.ReadAll(request.Body)
		if err != nil {
			t.Errorf("read refresh request: %v", err)
			return
		}
		values, err := url.ParseQuery(string(body))
		if err != nil {
			t.Errorf("parse refresh request: %v", err)
		}
		if values.Get("grant_type") != "refresh_token" || values.Get("refresh_token") != "refresh" {
			t.Errorf("refresh form = %s", body)
		}
		writeCodexJSON(t, writer, `{"access_token":"new-access","expires_in":"3600"}`)
	}))
	defer server.Close()

	authenticator, err := NewCodexDeviceAuth(WithCodexIssuer(server.URL))
	if err != nil {
		t.Fatal(err)
	}
	session, err := authenticator.NewSession(CodexTokens{
		AccessToken:  fakeJWT(time.Now().Add(time.Minute)),
		RefreshToken: "refresh",
	})
	if err != nil {
		t.Fatal(err)
	}
	headers, err := session.HeaderProvider()(t.Context(), RequestMeta{})
	if err != nil {
		t.Fatal(err)
	}
	if got := headers.Get("Authorization"); got != "Bearer new-access" {
		t.Fatalf("authorization = %q", got)
	}
	if got := session.Tokens().RefreshToken; got != "refresh" {
		t.Fatalf("refresh token = %q, want preserved token", got)
	}
}

func TestUnit_CodexDeviceAuth_HTTPError(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.NotFoundHandler())
	defer server.Close()
	authenticator, err := NewCodexDeviceAuth(WithCodexIssuer(server.URL))
	if err != nil {
		t.Fatal(err)
	}
	_, err = authenticator.RequestDeviceCode(t.Context())
	var httpErr *CodexHTTPError
	if !errors.As(err, &httpErr) || httpErr.StatusCode != http.StatusNotFound {
		t.Fatalf("error = %v, want typed 404 error", err)
	}
	if !errors.Is(err, ErrCodexServerResponse) {
		t.Fatalf("error = %v, want server response category", err)
	}
}

func writeCodexJSON(t *testing.T, writer http.ResponseWriter, body string) {
	t.Helper()
	writer.Header().Set("Content-Type", "application/json")
	if _, err := writer.Write([]byte(body)); err != nil {
		t.Errorf("write response: %v", err)
	}
}

func fakeJWT(expiry time.Time) string {
	header := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"none","typ":"JWT"}`))
	payload := base64.RawURLEncoding.EncodeToString(fmt.Appendf(nil, `{"exp":%s}`, strconv.FormatInt(expiry.Unix(), 10)))
	return header + "." + payload + ".signature"
}
