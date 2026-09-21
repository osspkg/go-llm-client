/*
 *  Copyright (c) 2026 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

package auth

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode"
)

const (
	// DefaultCodexIssuer is the OpenAI issuer used by Codex device login.
	DefaultCodexIssuer = "https://auth.openai.com"
	// DefaultCodexClientID is the public OAuth client ID registered for Codex.
	DefaultCodexClientID = "app_EMoamEEZ73f0CkXaXp7hrann"

	defaultCodexPollTimeout           = 15 * time.Minute
	defaultCodexPollInterval          = 5 * time.Second
	defaultCodexHTTPTimeout           = 30 * time.Second
	defaultCodexResponseLimit   int64 = 1 << 20
	defaultCodexRefreshSkew           = 5 * time.Minute
	defaultCodexRefreshCooldown       = time.Minute
	jwtParts                          = 3
)

var (
	// ErrInvalidCodexConfig indicates invalid device-auth configuration.
	ErrInvalidCodexConfig = errors.New("auth: invalid codex config")
	// ErrInvalidCodexCredentials indicates that the supplied token set is unusable.
	ErrInvalidCodexCredentials = errors.New("auth: invalid codex credentials")
	// ErrInvalidCodexDeviceCode indicates an incomplete device code.
	ErrInvalidCodexDeviceCode = errors.New("auth: invalid codex device code")
	// ErrCodexAccessTokenExpired indicates that an access token expired without a refresh token.
	ErrCodexAccessTokenExpired = errors.New("auth: codex access token expired")
	// ErrCodexResponseTooLarge indicates that an auth endpoint exceeded the response limit.
	ErrCodexResponseTooLarge = errors.New("auth: codex auth response too large")
	// ErrCodexServerResponse indicates a non-successful auth-server response.
	ErrCodexServerResponse = errors.New("auth: codex auth server response")
)

// CodexHTTPError describes a non-successful response from a Codex auth endpoint.
// Its error does not retain the response body because auth responses can contain
// authorization codes or tokens.
type CodexHTTPError struct {
	Operation  string
	StatusCode int
}

// Error implements error.
func (err *CodexHTTPError) Error() string {
	if err == nil {
		return "auth: codex auth server response"
	}
	return fmt.Sprintf("auth: codex %s: http status %d", err.Operation, err.StatusCode)
}

// Unwrap exposes the common Codex auth-server response category.
func (err *CodexHTTPError) Unwrap() error { return ErrCodexServerResponse }

// CodexOption configures a Codex device-authenticator.
type CodexOption func(*CodexDeviceAuth) error

// CodexDeviceAuth performs the OpenAI Codex device-code login flow.
//
// The authenticator is immutable after construction and safe for concurrent use.
// It does not persist credentials; callers decide how and where to store them.
type CodexDeviceAuth struct {
	issuer      string
	clientID    string
	httpClient  *http.Client
	pollTimeout time.Duration
}

// CodexDeviceCode contains the URL and one-time code shown to the user.
// The code is short-lived and must be treated as sensitive. It is intended to
// be obtained from RequestDeviceCode or Login, not constructed by callers.
type CodexDeviceCode struct {
	VerificationURL string
	UserCode        string
	ExpiresAt       time.Time

	deviceAuthID string
	interval     time.Duration
}

// CodexTokens contains credentials returned by the Codex OAuth token endpoint.
// Treat every field as a secret and never log or include it in diagnostics.
type CodexTokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	IDToken      string `json:"id_token"`
}

// CodexSession is a concurrent-safe, refreshable Codex credential session.
// Use HeaderProvider to connect it to openai.WithAuthProvider or another client
// accepting auth.HeaderProvider.
type CodexSession struct {
	provider *codexTokenProvider
}

// NewCodexDeviceAuth creates an authenticator for the Codex device-code flow.
func NewCodexDeviceAuth(options ...CodexOption) (*CodexDeviceAuth, error) {
	authenticator := &CodexDeviceAuth{
		issuer:      DefaultCodexIssuer,
		clientID:    DefaultCodexClientID,
		httpClient:  &http.Client{Timeout: defaultCodexHTTPTimeout},
		pollTimeout: defaultCodexPollTimeout,
	}
	for _, option := range options {
		if option == nil {
			continue
		}
		if err := option(authenticator); err != nil {
			return nil, err
		}
	}
	return authenticator, nil
}

// WithCodexIssuer changes the auth issuer. Use the default HTTPS issuer in
// production; HTTP is useful for a controlled local test server.
func WithCodexIssuer(issuer string) CodexOption {
	return func(authenticator *CodexDeviceAuth) error {
		normalized, err := normalizeCodexIssuer(issuer)
		if err != nil {
			return err
		}
		authenticator.issuer = normalized
		return nil
	}
}

// WithCodexClientID changes the public OAuth client ID used by the flow.
func WithCodexClientID(clientID string) CodexOption {
	return func(authenticator *CodexDeviceAuth) error {
		clientID = strings.TrimSpace(clientID)
		if clientID == "" || strings.IndexFunc(clientID, unicode.IsSpace) >= 0 {
			return fmt.Errorf("%w: client id", ErrInvalidCodexConfig)
		}
		authenticator.clientID = clientID
		return nil
	}
}

// WithCodexHTTPClient supplies the HTTP client used for auth endpoints.
func WithCodexHTTPClient(httpClient *http.Client) CodexOption {
	return func(authenticator *CodexDeviceAuth) error {
		if httpClient == nil {
			return fmt.Errorf("%w: nil http client", ErrInvalidCodexConfig)
		}
		authenticator.httpClient = httpClient
		return nil
	}
}

// WithCodexPollTimeout sets the maximum time spent polling after a device code
// is issued. The caller's context can impose a shorter deadline.
func WithCodexPollTimeout(timeout time.Duration) CodexOption {
	return func(authenticator *CodexDeviceAuth) error {
		if timeout <= 0 {
			return fmt.Errorf("%w: poll timeout", ErrInvalidCodexConfig)
		}
		authenticator.pollTimeout = timeout
		return nil
	}
}

// RequestDeviceCode starts a Codex device login and returns the URL and code to
// show to the user. It does not wait for browser authorization.
func (authenticator *CodexDeviceAuth) RequestDeviceCode(ctx context.Context) (CodexDeviceCode, error) {
	if err := authenticator.validate(); err != nil {
		return CodexDeviceCode{}, err
	}
	if err := validateCodexContext(ctx); err != nil {
		return CodexDeviceCode{}, err
	}

	body, err := marshalCodexJSON(codexUserCodeRequest{ClientID: authenticator.clientID})
	if err != nil {
		return CodexDeviceCode{}, fmt.Errorf("marshal codex device request: %w", err)
	}
	body, err = authenticator.post(ctx, "request device code", "/api/accounts/deviceauth/usercode", "application/json", body)
	if err != nil {
		return CodexDeviceCode{}, err
	}

	var response codexUserCodeResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return CodexDeviceCode{}, fmt.Errorf("decode codex device response: %w", err)
	}
	userCode := response.UserCode
	if strings.TrimSpace(userCode) == "" {
		userCode = response.UserCodeAlt
	}
	if strings.TrimSpace(response.DeviceAuthID) == "" || strings.TrimSpace(userCode) == "" {
		return CodexDeviceCode{}, fmt.Errorf("%w: response is missing device auth id or user code", ErrInvalidCodexDeviceCode)
	}
	interval := response.Interval.duration()
	if interval <= 0 {
		interval = defaultCodexPollInterval
	}

	return CodexDeviceCode{
		VerificationURL: authenticator.issuer + "/codex/device",
		UserCode:        strings.TrimSpace(userCode),
		ExpiresAt:       time.Now().Add(authenticator.pollTimeout),
		deviceAuthID:    strings.TrimSpace(response.DeviceAuthID),
		interval:        interval,
	}, nil
}

// Login runs the complete device-code flow. showCode is called once with the
// browser URL and one-time code before polling begins; the returned session
// refreshes access tokens when their expiry is known.
func (authenticator *CodexDeviceAuth) Login(ctx context.Context, showCode func(CodexDeviceCode) error) (*CodexSession, error) {
	if showCode == nil {
		return nil, fmt.Errorf("%w: nil code callback", ErrInvalidCodexConfig)
	}
	deviceCode, err := authenticator.RequestDeviceCode(ctx)
	if err != nil {
		return nil, err
	}
	if err := showCode(deviceCode); err != nil {
		return nil, fmt.Errorf("show codex device code: %w", err)
	}
	return authenticator.Complete(ctx, deviceCode)
}

// Complete waits for browser authorization of a previously requested device
// code and exchanges the returned authorization code for Codex tokens.
func (authenticator *CodexDeviceAuth) Complete(ctx context.Context, deviceCode CodexDeviceCode) (*CodexSession, error) {
	if err := authenticator.validate(); err != nil {
		return nil, err
	}
	if err := validateCodexContext(ctx); err != nil {
		return nil, err
	}
	if strings.TrimSpace(deviceCode.deviceAuthID) == "" || strings.TrimSpace(deviceCode.UserCode) == "" {
		return nil, fmt.Errorf("%w: response is missing device auth id or user code", ErrInvalidCodexDeviceCode)
	}

	pollContext, cancel := context.WithTimeout(ctx, authenticator.pollTimeout)
	defer cancel()

	codeResponse, err := authenticator.pollDeviceCode(pollContext, deviceCode)
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(codeResponse.AuthorizationCode) == "" || strings.TrimSpace(codeResponse.CodeVerifier) == "" {
		return nil, fmt.Errorf("%w: response is missing authorization code or code verifier", ErrInvalidCodexCredentials)
	}

	tokens, err := authenticator.exchangeAuthorizationCode(pollContext, codeResponse)
	if err != nil {
		return nil, err
	}
	return authenticator.newSession(tokens), nil
}

// NewSession creates a refreshable session from credentials loaded by the
// caller, for example from a protected credential store.
func (authenticator *CodexDeviceAuth) NewSession(tokens CodexTokens) (*CodexSession, error) {
	if err := authenticator.validate(); err != nil {
		return nil, err
	}
	if strings.TrimSpace(tokens.AccessToken) == "" {
		return nil, fmt.Errorf("%w: missing access token", ErrInvalidCodexCredentials)
	}
	return authenticator.newSessionFromTokens(tokens), nil
}

// HeaderProvider returns the per-request provider for this session.
func (session *CodexSession) HeaderProvider() HeaderProvider {
	if session == nil || session.provider == nil {
		return nil
	}
	return session.provider.headers
}

// Tokens returns the current credentials, including any refreshed token values.
func (session *CodexSession) Tokens() CodexTokens {
	if session == nil || session.provider == nil {
		return CodexTokens{}
	}
	return session.provider.tokensSnapshot()
}

type codexUserCodeRequest struct {
	ClientID string `json:"client_id"`
}

type codexUserCodeResponse struct {
	DeviceAuthID string        `json:"device_auth_id"`
	UserCode     string        `json:"user_code"`
	UserCodeAlt  string        `json:"usercode"`
	Interval     codexDuration `json:"interval"`
}

type codexDeviceTokenRequest struct {
	DeviceAuthID string `json:"device_auth_id"`
	UserCode     string `json:"user_code"`
}

type codexDeviceTokenResponse struct {
	AuthorizationCode string `json:"authorization_code"`
	CodeChallenge     string `json:"code_challenge"`
	CodeVerifier      string `json:"code_verifier"`
}

type codexTokenResponse struct {
	AccessToken  string        `json:"access_token"`
	RefreshToken string        `json:"refresh_token"`
	IDToken      string        `json:"id_token"`
	ExpiresIn    codexDuration `json:"expires_in"`
}

type codexDuration time.Duration

// UnmarshalJSON decodes a duration represented in seconds as a number or string.
func (duration *codexDuration) UnmarshalJSON(data []byte) error {
	raw := strings.TrimSpace(string(data))
	if raw == "" || raw == "null" {
		return nil
	}
	if raw[0] == '"' {
		var value string
		if err := json.Unmarshal(data, &value); err != nil {
			return err
		}
		raw = strings.TrimSpace(value)
	}
	seconds, err := strconv.ParseInt(raw, 10, 64)
	if err != nil {
		return fmt.Errorf("parse seconds: %w", err)
	}
	maxSeconds := int64((time.Duration(1<<63 - 1)) / time.Second)
	if seconds < 0 || seconds > maxSeconds {
		return errors.New("seconds out of range")
	}
	*duration = codexDuration(time.Duration(seconds) * time.Second)
	return nil
}

func (duration *codexDuration) duration() time.Duration {
	if duration == nil {
		return 0
	}
	return time.Duration(*duration)
}

func (authenticator *CodexDeviceAuth) pollDeviceCode(ctx context.Context, deviceCode CodexDeviceCode) (codexDeviceTokenResponse, error) {
	interval := deviceCode.interval
	if interval <= 0 {
		interval = defaultCodexPollInterval
	}
	requestBody, err := marshalCodexJSON(codexDeviceTokenRequest{
		DeviceAuthID: deviceCode.deviceAuthID,
		UserCode:     deviceCode.UserCode,
	})
	if err != nil {
		return codexDeviceTokenResponse{}, fmt.Errorf("marshal codex device poll: %w", err)
	}

	for {
		body, postErr := authenticator.post(ctx, "poll device code", "/api/accounts/deviceauth/token", "application/json", requestBody)
		if postErr == nil {
			var response codexDeviceTokenResponse
			if err := json.Unmarshal(body, &response); err != nil {
				return codexDeviceTokenResponse{}, fmt.Errorf("decode codex device token response: %w", err)
			}
			return response, nil
		}

		var httpErr *CodexHTTPError
		if !errors.As(postErr, &httpErr) || (httpErr.StatusCode != http.StatusForbidden && httpErr.StatusCode != http.StatusNotFound) {
			return codexDeviceTokenResponse{}, postErr
		}
		if err := waitCodex(ctx, interval); err != nil {
			return codexDeviceTokenResponse{}, fmt.Errorf("poll codex device code: %w", err)
		}
	}
}

func (authenticator *CodexDeviceAuth) exchangeAuthorizationCode(ctx context.Context, response codexDeviceTokenResponse) (codexTokenResponse, error) {
	form := url.Values{
		"grant_type":    {"authorization_code"},
		"code":          {strings.TrimSpace(response.AuthorizationCode)},
		"redirect_uri":  {authenticator.issuer + "/deviceauth/callback"},
		"client_id":     {authenticator.clientID},
		"code_verifier": {strings.TrimSpace(response.CodeVerifier)},
	}
	body, err := authenticator.post(ctx, "exchange authorization code", "/oauth/token", "application/x-www-form-urlencoded", []byte(form.Encode()))
	if err != nil {
		return codexTokenResponse{}, err
	}

	var tokens codexTokenResponse
	if err := json.Unmarshal(body, &tokens); err != nil {
		return codexTokenResponse{}, fmt.Errorf("decode codex token response: %w", err)
	}
	if strings.TrimSpace(tokens.AccessToken) == "" {
		return codexTokenResponse{}, fmt.Errorf("%w: token response is missing access token", ErrInvalidCodexCredentials)
	}
	return tokens, nil
}

func (authenticator *CodexDeviceAuth) refresh(ctx context.Context, refreshToken string) (codexTokenResponse, error) {
	form := url.Values{
		"grant_type":    {"refresh_token"},
		"refresh_token": {refreshToken},
		"client_id":     {authenticator.clientID},
	}
	body, err := authenticator.post(ctx, "refresh access token", "/oauth/token", "application/x-www-form-urlencoded", []byte(form.Encode()))
	if err != nil {
		return codexTokenResponse{}, err
	}

	var tokens codexTokenResponse
	if err := json.Unmarshal(body, &tokens); err != nil {
		return codexTokenResponse{}, fmt.Errorf("decode codex refresh response: %w", err)
	}
	if strings.TrimSpace(tokens.AccessToken) == "" {
		return codexTokenResponse{}, fmt.Errorf("%w: refresh response is missing access token", ErrInvalidCodexCredentials)
	}
	return tokens, nil
}

func (authenticator *CodexDeviceAuth) newSession(tokens codexTokenResponse) *CodexSession {
	return &CodexSession{provider: newCodexTokenProvider(authenticator, CodexTokens{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		IDToken:      tokens.IDToken,
	}, tokens.ExpiresIn.duration())}
}

func (authenticator *CodexDeviceAuth) newSessionFromTokens(tokens CodexTokens) *CodexSession {
	return &CodexSession{provider: newCodexTokenProvider(authenticator, tokens, 0)}
}

func newCodexTokenProvider(authenticator *CodexDeviceAuth, tokens CodexTokens, expiresIn time.Duration) *codexTokenProvider {
	tokens.AccessToken = strings.TrimSpace(tokens.AccessToken)
	tokens.RefreshToken = strings.TrimSpace(tokens.RefreshToken)
	tokens.IDToken = strings.TrimSpace(tokens.IDToken)
	expiresAt := time.Time{}
	if expiresIn > 0 {
		expiresAt = time.Now().Add(expiresIn)
	}
	if expiresAt.IsZero() {
		expiresAt = jwtExpiry(tokens.AccessToken)
	}
	gate := make(chan struct{}, 1)
	gate <- struct{}{}
	return &codexTokenProvider{
		authenticator: authenticator,
		tokens:        tokens,
		expiresAt:     expiresAt,
		refreshGate:   gate,
	}
}

type codexTokenProvider struct {
	authenticator *CodexDeviceAuth
	mu            sync.Mutex
	tokens        CodexTokens
	expiresAt     time.Time
	refreshGate   chan struct{}
	refreshErr    error
	refreshAt     time.Time
}

func (provider *codexTokenProvider) headers(ctx context.Context, _ RequestMeta) (http.Header, error) {
	if err := validateCodexContext(ctx); err != nil {
		return nil, err
	}

	for {
		provider.mu.Lock()
		tokens := provider.tokens
		expiresAt := provider.expiresAt
		refreshErr := provider.refreshErr
		refreshAt := provider.refreshAt
		provider.mu.Unlock()

		if !needsCodexRefresh(expiresAt) {
			return http.Header{"Authorization": []string{"Bearer " + tokens.AccessToken}}, nil
		}
		if tokens.RefreshToken == "" {
			return nil, ErrCodexAccessTokenExpired
		}
		if refreshErr != nil && time.Since(refreshAt) < defaultCodexRefreshCooldown {
			return nil, refreshErr
		}

		select {
		case <-provider.refreshGate:
		case <-ctx.Done():
			return nil, ctx.Err()
		}

		provider.mu.Lock()
		if !needsCodexRefresh(provider.expiresAt) {
			provider.mu.Unlock()
			provider.refreshGate <- struct{}{}
			continue
		}
		refreshToken := provider.tokens.RefreshToken
		provider.mu.Unlock()

		refreshed, err := provider.authenticator.refresh(ctx, refreshToken)
		if err != nil {
			provider.mu.Lock()
			if !errors.Is(err, context.Canceled) && !errors.Is(err, context.DeadlineExceeded) {
				provider.refreshErr = err
				provider.refreshAt = time.Now()
			}
			provider.mu.Unlock()
			provider.refreshGate <- struct{}{}
			return nil, fmt.Errorf("refresh codex token: %w", err)
		}

		provider.mu.Lock()
		provider.tokens.AccessToken = strings.TrimSpace(refreshed.AccessToken)
		if strings.TrimSpace(refreshed.RefreshToken) != "" {
			provider.tokens.RefreshToken = strings.TrimSpace(refreshed.RefreshToken)
		}
		if strings.TrimSpace(refreshed.IDToken) != "" {
			provider.tokens.IDToken = strings.TrimSpace(refreshed.IDToken)
		}
		expiresIn := refreshed.ExpiresIn.duration()
		if expiresIn > 0 {
			provider.expiresAt = time.Now().Add(expiresIn)
		} else {
			provider.expiresAt = jwtExpiry(provider.tokens.AccessToken)
		}
		provider.refreshErr = nil
		provider.refreshAt = time.Time{}
		provider.mu.Unlock()
		provider.refreshGate <- struct{}{}
	}
}

func (provider *codexTokenProvider) tokensSnapshot() CodexTokens {
	provider.mu.Lock()
	defer provider.mu.Unlock()
	return provider.tokens
}

func needsCodexRefresh(expiresAt time.Time) bool {
	return !expiresAt.IsZero() && !time.Now().Before(expiresAt.Add(-defaultCodexRefreshSkew))
}

func jwtExpiry(token string) time.Time {
	parts := strings.Split(token, ".")
	if len(parts) != jwtParts {
		return time.Time{}
	}
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		payload, err = base64.URLEncoding.DecodeString(parts[1])
	}
	if err != nil {
		return time.Time{}
	}
	var claims struct {
		ExpiresAt int64 `json:"exp"`
	}
	if err := json.Unmarshal(payload, &claims); err != nil || claims.ExpiresAt <= 0 {
		return time.Time{}
	}
	return time.Unix(claims.ExpiresAt, 0)
}

func (authenticator *CodexDeviceAuth) post(ctx context.Context, operation, endpoint, contentType string, body []byte) ([]byte, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, authenticator.issuer+endpoint, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("create codex %s request: %w", operation, err)
	}
	request.Header.Set("Content-Type", contentType)
	request.Header.Set("Accept", "application/json")
	response, err := authenticator.httpClient.Do(request)
	if err != nil {
		if response != nil && response.Body != nil {
			_ = response.Body.Close()
		}
		return nil, fmt.Errorf("execute codex %s request: %w", operation, err)
	}
	if response == nil {
		return nil, fmt.Errorf("execute codex %s request: empty response", operation)
	}
	data, readErr := readCodexResponse(response.Body)
	if readErr != nil {
		return nil, fmt.Errorf("read codex %s response: %w", operation, readErr)
	}
	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return nil, &CodexHTTPError{Operation: operation, StatusCode: response.StatusCode}
	}
	return data, nil
}

func readCodexResponse(body io.ReadCloser) ([]byte, error) {
	if body == nil {
		return nil, errors.New("nil response body")
	}
	data, readErr := io.ReadAll(io.LimitReader(body, defaultCodexResponseLimit+1))
	closeErr := body.Close()
	if readErr != nil {
		return nil, readErr
	}
	if closeErr != nil {
		return nil, closeErr
	}
	if int64(len(data)) > defaultCodexResponseLimit {
		return nil, ErrCodexResponseTooLarge
	}
	return data, nil
}

func marshalCodexJSON(value any) ([]byte, error) { return json.Marshal(value) }

func waitCodex(ctx context.Context, interval time.Duration) error {
	timer := time.NewTimer(interval)
	defer timer.Stop()
	select {
	case <-timer.C:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func normalizeCodexIssuer(issuer string) (string, error) {
	issuer = strings.TrimRight(strings.TrimSpace(issuer), "/")
	parsed, err := url.Parse(issuer)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" || parsed.User != nil || parsed.RawQuery != "" || parsed.Fragment != "" || (parsed.Scheme != "http" && parsed.Scheme != "https") {
		return "", fmt.Errorf("%w: issuer", ErrInvalidCodexConfig)
	}
	return issuer, nil
}

func validateCodexContext(ctx context.Context) error {
	if ctx == nil {
		return fmt.Errorf("%w: nil context", ErrInvalidCodexConfig)
	}
	return nil
}

func (authenticator *CodexDeviceAuth) validate() error {
	if authenticator == nil || authenticator.issuer == "" || authenticator.clientID == "" || authenticator.httpClient == nil || authenticator.pollTimeout <= 0 {
		return fmt.Errorf("%w: incomplete authenticator", ErrInvalidCodexConfig)
	}
	return nil
}
