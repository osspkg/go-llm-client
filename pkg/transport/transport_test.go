package transport_test

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"go.osspkg.com/llm-client/pkg/auth"
	llmerrors "go.osspkg.com/llm-client/pkg/errors"
	"go.osspkg.com/llm-client/pkg/transport"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (fn roundTripFunc) RoundTrip(request *http.Request) (*http.Response, error) { return fn(request) }

func TestRequestUsesBaseURLAndPerRequestAuth(t *testing.T) {
	client, err := transport.New("https://example.test/v1",
		transport.WithHTTPClient(&http.Client{Transport: roundTripFunc(func(request *http.Request) (*http.Response, error) {
			if request.URL.String() != "https://example.test/v1/models" {
				t.Errorf("url = %s", request.URL)
			}
			if request.Header.Get("Authorization") != "Bearer secret" {
				t.Errorf("authorization header = %q", request.Header.Get("Authorization"))
			}
			return &http.Response{StatusCode: http.StatusOK, Header: make(http.Header), Body: io.NopCloser(strings.NewReader(`{"ok":true}`))}, nil
		})}),
		transport.WithAuthProvider(auth.StaticBearer("secret")),
	)
	if err != nil {
		t.Fatal(err)
	}
	data, err := client.Request(context.Background(), http.MethodGet, "/models", "models.list", nil, "")
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != `{"ok":true}` {
		t.Fatalf("body = %s", data)
	}
}

func TestRequestPreservesEscapedPathSegments(t *testing.T) {
	client, err := transport.New("https://example.test/v1", transport.WithHTTPClient(&http.Client{Transport: roundTripFunc(func(request *http.Request) (*http.Response, error) {
		if request.URL.EscapedPath() != "/v1/models/model%2Fwith%2Fslashes" {
			t.Errorf("escaped path = %s", request.URL.EscapedPath())
		}
		return &http.Response{StatusCode: http.StatusOK, Header: make(http.Header), Body: io.NopCloser(strings.NewReader("{}"))}, nil
	})}))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := client.Request(context.Background(), http.MethodGet, "/models/model%2Fwith%2Fslashes", "models.get", nil, ""); err != nil {
		t.Fatal(err)
	}
}

func TestRequestReaderIsBounded(t *testing.T) {
	client, err := transport.New("https://example.test", transport.WithMaxRequestBody(2), transport.WithHTTPClient(&http.Client{Transport: roundTripFunc(func(request *http.Request) (*http.Response, error) {
		_, err := io.ReadAll(request.Body)
		if !errors.Is(err, llmerrors.ErrBodyTooLarge) {
			t.Errorf("request body error = %v", err)
		}
		return nil, err
	})}))
	if err != nil {
		t.Fatal(err)
	}
	_, err = client.RequestReader(context.Background(), http.MethodPost, "/upload", "upload", strings.NewReader("123"), "text/plain")
	if !errors.Is(err, llmerrors.ErrBodyTooLarge) {
		t.Fatalf("error = %v", err)
	}
}

func TestRequestMapsStatusAndBodyLimit(t *testing.T) {
	client, err := transport.New("https://example.test", transport.WithHTTPClient(&http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: http.StatusTooManyRequests, Header: http.Header{"X-Request-Id": []string{"request-id"}}, Body: io.NopCloser(strings.NewReader(`{"error":{"code":"rate_limit","type":"throttle","message":"slow down"}}`))}, nil
	})}))
	if err != nil {
		t.Fatal(err)
	}
	_, err = client.Request(context.Background(), http.MethodGet, "/models", "models.list", nil, "")
	var httpErr *llmerrors.HTTPError
	if !errors.As(err, &httpErr) || httpErr.StatusCode != http.StatusTooManyRequests || httpErr.RequestID != "request-id" || httpErr.ProviderCode != "rate_limit" {
		t.Fatalf("error = %v", err)
	}

	client, err = transport.New("https://example.test", transport.WithMaxResponseBody(2), transport.WithHTTPClient(&http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: http.StatusOK, Header: make(http.Header), Body: io.NopCloser(strings.NewReader("123"))}, nil
	})}))
	if err != nil {
		t.Fatal(err)
	}
	_, err = client.Request(context.Background(), http.MethodGet, "/models", "models.list", nil, "")
	if !errors.Is(err, llmerrors.ErrBodyTooLarge) {
		t.Fatalf("error = %v", err)
	}
}
