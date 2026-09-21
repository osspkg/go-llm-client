// Package transport provides bounded HTTP execution shared by provider clients.
package transport

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"go.osspkg.com/llm-client/pkg/auth"
	llmerrors "go.osspkg.com/llm-client/pkg/errors"
)

const (
	defaultMaxResponseBody = 32 << 20
	defaultMaxRequestBody  = 64 << 20
	defaultRequestTimeout  = 2 * time.Minute
	defaultHeaderTimeout   = 30 * time.Second
	defaultDialTimeout     = 10 * time.Second
	defaultKeepAlive       = 30 * time.Second
	defaultIdleConnTimeout = 90 * time.Second
	maxErrorMessage        = 4 << 10
)

// RetryPolicy controls bounded retries for safe requests.
type RetryPolicy struct {
	MaxAttempts int
	Backoff     time.Duration
}

// Option configures a Client.
type Option func(*Client) error

// OperationGate validates operation metadata before network execution.
type OperationGate func(string) error

// Client executes provider requests against a configured base URL.
type Client struct {
	baseURL         *url.URL
	httpClient      *http.Client
	auth            auth.HeaderProvider
	defaultHeaders  http.Header
	maxResponseBody int64
	maxRequestBody  int64
	requestTimeout  time.Duration
	streamTimeout   time.Duration
	retry           RetryPolicy
	gate            OperationGate
}

// New creates a bounded HTTP client with stdlib defaults.
func New(baseURL string, options ...Option) (*Client, error) {
	parsed, err := parseBaseURL(baseURL)
	if err != nil {
		return nil, err
	}
	client := &Client{
		baseURL:         parsed,
		maxResponseBody: defaultMaxResponseBody,
		maxRequestBody:  defaultMaxRequestBody,
		requestTimeout:  defaultRequestTimeout,
		defaultHeaders:  make(http.Header),
		retry:           RetryPolicy{MaxAttempts: 1},
		httpClient:      &http.Client{Transport: defaultTransport()},
	}
	for _, option := range options {
		if option == nil {
			continue
		}
		if err := option(client); err != nil {
			return nil, err
		}
	}
	return client, nil
}

// WithHTTPClient replaces the internal HTTP client, primarily for tests.
func WithHTTPClient(httpClient *http.Client) Option {
	return func(client *Client) error {
		if httpClient == nil {
			return fmt.Errorf("%w: nil http client", llmerrors.ErrInvalidConfig)
		}
		client.httpClient = httpClient
		return nil
	}
}

// WithBaseURL replaces the endpoint URL after construction.
func WithBaseURL(baseURL string) Option {
	return func(client *Client) error {
		parsed, err := parseBaseURL(baseURL)
		if err != nil {
			return err
		}
		client.baseURL = parsed
		return nil
	}
}

func parseBaseURL(value string) (*url.URL, error) {
	parsed, err := url.Parse(value)
	if err != nil || parsed.Host == "" || (parsed.Scheme != "http" && parsed.Scheme != "https") || parsed.User != nil {
		return nil, fmt.Errorf("%w: base url", llmerrors.ErrInvalidConfig)
	}
	return parsed, nil
}

// WithAuthProvider configures per-request authentication headers.
func WithAuthProvider(provider auth.HeaderProvider) Option {
	return func(client *Client) error {
		client.auth = provider
		return nil
	}
}

// WithHeader adds a defensive copy of a default header.
func WithHeader(name, value string) Option {
	return func(client *Client) error {
		if strings.TrimSpace(name) == "" {
			return fmt.Errorf("%w: empty header name", llmerrors.ErrInvalidConfig)
		}
		client.defaultHeaders.Add(name, value)
		return nil
	}
}

// WithMaxResponseBody sets the maximum buffered response size.
func WithMaxResponseBody(limit int64) Option {
	return func(client *Client) error {
		if limit <= 0 {
			return fmt.Errorf("%w: max response body", llmerrors.ErrInvalidConfig)
		}
		client.maxResponseBody = limit
		return nil
	}
}

// WithMaxRequestBody sets the maximum bytes accepted from a request reader.
func WithMaxRequestBody(limit int64) Option {
	return func(client *Client) error {
		if limit <= 0 {
			return fmt.Errorf("%w: max request body", llmerrors.ErrInvalidConfig)
		}
		client.maxRequestBody = limit
		return nil
	}
}

// WithRequestTimeout sets the default timeout for non-streaming requests.
func WithRequestTimeout(timeout time.Duration) Option {
	return func(client *Client) error {
		if timeout <= 0 {
			return fmt.Errorf("%w: request timeout", llmerrors.ErrInvalidConfig)
		}
		client.requestTimeout = timeout
		return nil
	}
}

// WithStreamTimeout sets an optional total timeout for streaming requests.
func WithStreamTimeout(timeout time.Duration) Option {
	return func(client *Client) error {
		if timeout < 0 {
			return fmt.Errorf("%w: stream timeout", llmerrors.ErrInvalidConfig)
		}
		client.streamTimeout = timeout
		return nil
	}
}

// WithRetryPolicy enables bounded retries for idempotent HTTP methods.
func WithRetryPolicy(policy RetryPolicy) Option {
	return func(client *Client) error {
		if policy.MaxAttempts < 1 || policy.MaxAttempts > 5 || policy.Backoff < 0 {
			return fmt.Errorf("%w: retry policy", llmerrors.ErrInvalidConfig)
		}
		client.retry = policy
		return nil
	}
}

// WithOperationGate configures a provider-specific capability gate without
// making the shared transport depend on a provider package.
func WithOperationGate(gate OperationGate) Option {
	return func(client *Client) error {
		client.gate = gate
		return nil
	}
}

// Request executes a request and buffers its bounded response body.
func (client *Client) Request(ctx context.Context, method, endpoint, operation string, body []byte, contentType string) ([]byte, error) { //nolint:revive // positional arguments mirror the low-level HTTP contract.
	if ctx == nil {
		return nil, llmerrors.ErrInvalidRequest
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if client.gate != nil {
		if err := client.gate(operation); err != nil {
			return nil, err
		}
	}
	requestContext, cancel := client.withTimeout(ctx, false)
	defer cancel()
	attempts := client.retry.MaxAttempts
	if attempts == 0 || !retryableMethod(method) {
		attempts = 1
	}
	for attempt := 1; attempt <= attempts; attempt++ {
		response, err := client.do(requestContext, method, endpoint, operation, bytes.NewReader(body), contentType, "") //nolint:bodyclose // readResponse owns and closes the body.
		if err != nil {
			return nil, err
		}
		data, readErr := client.readResponse(response)
		if readErr != nil {
			return nil, readErr
		}
		if response.StatusCode >= http.StatusOK && response.StatusCode < http.StatusMultipleChoices {
			return data, nil
		}
		if attempt < attempts && response.StatusCode >= 500 {
			if err := sleep(requestContext, client.retry.Backoff, attempt); err != nil {
				return nil, err
			}
			continue
		}
		return nil, newHTTPError(response, data)
	}
	return nil, llmerrors.ErrProtocol
}

// RequestReader executes one non-retried request from a streaming body.
// The reader is consumed at most up to the configured request body limit.
func (client *Client) RequestReader(ctx context.Context, method, endpoint, operation string, body io.Reader, contentType string) ([]byte, error) { //nolint:revive // positional arguments mirror the low-level HTTP contract.
	if ctx == nil {
		return nil, llmerrors.ErrInvalidRequest
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if client.gate != nil {
		if err := client.gate(operation); err != nil {
			return nil, err
		}
	}
	requestContext, cancel := client.withTimeout(ctx, false)
	defer cancel()
	response, err := client.do(requestContext, method, endpoint, operation, body, contentType, "") //nolint:bodyclose // readResponse owns and closes the response.
	if err != nil {
		return nil, err
	}
	data, readErr := client.readResponse(response)
	if readErr != nil {
		return nil, readErr
	}
	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return nil, newHTTPError(response, data)
	}
	return data, nil
}

// Stream opens a response body for a streaming request. The caller owns and must close it.
func (client *Client) Stream(ctx context.Context, method, endpoint, operation string, body []byte, contentType, accept string) (io.ReadCloser, error) { //nolint:revive // positional arguments mirror the low-level streaming contract.
	return client.streamReader(ctx, method, endpoint, operation, bytes.NewReader(body), contentType, accept)
}

// StreamReader opens a streaming response for a one-shot reader request body.
// The reader is consumed at most once and the caller owns the returned body.
func (client *Client) StreamReader(ctx context.Context, method, endpoint, operation string, body io.Reader, contentType, accept string) (io.ReadCloser, error) { //nolint:revive // positional arguments mirror the low-level streaming contract.
	return client.streamReader(ctx, method, endpoint, operation, body, contentType, accept)
}

func (client *Client) streamReader(ctx context.Context, method, endpoint, operation string, body io.Reader, contentType, accept string) (io.ReadCloser, error) { //nolint:revive // internal helper keeps low-level transport metadata together.
	if ctx == nil {
		return nil, llmerrors.ErrInvalidRequest
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if client.gate != nil {
		if err := client.gate(operation); err != nil {
			return nil, err
		}
	}
	requestContext, cancel := client.withTimeout(ctx, true)
	response, err := client.do(requestContext, method, endpoint, operation, body, contentType, accept) //nolint:bodyclose // the returned stream owns the body.
	if err != nil {
		cancel()
		return nil, err
	}
	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		data, readErr := client.readResponse(response)
		cancel()
		if readErr != nil {
			return nil, readErr
		}
		return nil, newHTTPError(response, data)
	}
	return &cancelingBody{ReadCloser: response.Body, cancel: cancel}, nil
}

// HTTPClient returns the configured HTTP client for transports that need to upgrade a connection.
func (client *Client) HTTPClient() *http.Client { return client.httpClient }

// WebSocketURL resolves an HTTP base URL as a WebSocket URL.
func (client *Client) WebSocketURL(endpoint string) (string, error) {
	target, err := client.resolve(endpoint)
	if err != nil {
		return "", err
	}
	switch target.Scheme {
	case "http":
		target.Scheme = "ws"
	case "https":
		target.Scheme = "wss"
	default:
		return "", fmt.Errorf("%w: websocket scheme", llmerrors.ErrInvalidConfig)
	}
	return target.String(), nil
}

// WebSocketHeaders returns default and per-request authentication headers.
func (client *Client) WebSocketHeaders(ctx context.Context, operation string) (http.Header, error) {
	if client.gate != nil {
		if err := client.gate(operation); err != nil {
			return nil, err
		}
	}
	headers := client.defaultHeaders.Clone()
	if client.auth == nil {
		return headers, nil
	}
	authHeaders, err := client.auth(ctx, auth.RequestMeta{
		Method:    http.MethodGet,
		Domain:    client.baseURL.Hostname(),
		Operation: operation,
		WebSocket: true,
	})
	if err != nil {
		return nil, fmt.Errorf("authenticate %s: %w", operation, err)
	}
	for key, values := range authHeaders {
		for _, value := range values {
			headers.Add(key, value)
		}
	}
	return headers, nil
}

func (client *Client) do(ctx context.Context, method, endpoint, operation string, body io.Reader, contentType, accept string) (*http.Response, error) { //nolint:revive // internal transport contract keeps request metadata together.
	target, err := client.resolve(endpoint)
	if err != nil {
		return nil, err
	}
	requestBody := io.Reader(bytes.NewReader(nil))
	if body != nil {
		requestBody = &limitedReader{reader: body, remaining: client.maxRequestBody}
	}
	request, err := http.NewRequestWithContext(ctx, method, target.String(), requestBody)
	if err != nil {
		return nil, fmt.Errorf("create %s request: %w", operation, err)
	}
	request.Header = client.defaultHeaders.Clone()
	if contentType != "" {
		request.Header.Set("Content-Type", contentType)
	}
	if accept != "" {
		request.Header.Set("Accept", accept)
	}
	if client.auth != nil {
		headers, authErr := client.auth(ctx, auth.RequestMeta{
			Method:    method,
			URL:       target.String(),
			Domain:    target.Hostname(),
			Operation: operation,
			Streaming: accept != "",
		})
		if authErr != nil {
			_ = request.Body.Close()
			return nil, fmt.Errorf("authenticate %s: %w", operation, authErr)
		}
		for key, values := range headers {
			for _, value := range values {
				request.Header.Add(key, value)
			}
		}
	}
	response, err := client.httpClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("execute %s request: %w", operation, err)
	}
	return response, nil
}

type limitedReader struct {
	reader    io.Reader
	remaining int64
}

func (reader *limitedReader) Read(buffer []byte) (int, error) {
	if reader.remaining == 0 {
		var probe [1]byte
		n, err := reader.reader.Read(probe[:])
		if n > 0 {
			return 0, llmerrors.ErrBodyTooLarge
		}
		return 0, err
	}
	if int64(len(buffer)) > reader.remaining {
		buffer = buffer[:reader.remaining]
	}
	n, err := reader.reader.Read(buffer)
	reader.remaining -= int64(n)
	return n, err
}

// Close closes the wrapped reader when it owns a close operation.
func (reader *limitedReader) Close() error {
	closer, ok := reader.reader.(io.Closer)
	if !ok {
		return nil
	}
	return closer.Close()
}

func (client *Client) resolve(endpoint string) (*url.URL, error) {
	parsed, err := url.Parse(endpoint)
	if err != nil || parsed.IsAbs() || parsed.Host != "" {
		return nil, fmt.Errorf("%w: endpoint url", llmerrors.ErrInvalidRequest)
	}
	resolved := *client.baseURL
	basePath := strings.TrimRight(client.baseURL.EscapedPath(), "/")
	endpointPath := strings.TrimLeft(parsed.EscapedPath(), "/")
	joinedPath := path.Join(basePath, endpointPath)
	decodedPath, err := url.PathUnescape(joinedPath)
	if err != nil {
		return nil, fmt.Errorf("%w: endpoint path", llmerrors.ErrInvalidRequest)
	}
	resolved.Path = decodedPath
	resolved.RawPath = joinedPath
	resolved.RawQuery = parsed.RawQuery
	resolved.Fragment = ""
	return &resolved, nil
}

func (client *Client) readResponse(response *http.Response) ([]byte, error) {
	defer func() { _ = response.Body.Close() }()
	data, err := io.ReadAll(io.LimitReader(response.Body, client.maxResponseBody+1))
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}
	if int64(len(data)) > client.maxResponseBody {
		return nil, llmerrors.ErrBodyTooLarge
	}
	return data, nil
}

func (client *Client) withTimeout(ctx context.Context, streaming bool) (context.Context, context.CancelFunc) {
	timeout := client.requestTimeout
	if streaming {
		timeout = client.streamTimeout
	}
	if timeout <= 0 {
		return ctx, func() {}
	}
	if _, ok := ctx.Deadline(); ok {
		return ctx, func() {}
	}
	return context.WithTimeout(ctx, timeout)
}

func defaultTransport() http.RoundTripper {
	transport, ok := http.DefaultTransport.(*http.Transport)
	if !ok {
		return http.DefaultTransport
	}
	copyTransport := transport.Clone()
	copyTransport.DialContext = (&net.Dialer{Timeout: defaultDialTimeout, KeepAlive: defaultKeepAlive}).DialContext
	copyTransport.TLSHandshakeTimeout = defaultHeaderTimeout
	copyTransport.ResponseHeaderTimeout = defaultHeaderTimeout
	copyTransport.IdleConnTimeout = defaultIdleConnTimeout
	copyTransport.TLSClientConfig = &tls.Config{MinVersion: tls.VersionTLS12}
	return copyTransport
}

func retryableMethod(method string) bool {
	switch method {
	case http.MethodGet, http.MethodHead, http.MethodOptions, http.MethodDelete:
		return true
	default:
		return false
	}
}

func sleep(ctx context.Context, backoff time.Duration, attempt int) error {
	if backoff <= 0 {
		return nil
	}
	timer := time.NewTimer(backoff * time.Duration(attempt))
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}

func newHTTPError(response *http.Response, data []byte) error {
	output := &llmerrors.HTTPError{
		StatusCode: response.StatusCode,
		RequestID:  response.Header.Get("x-request-id"),
		Message:    http.StatusText(response.StatusCode),
	}
	providerError := decodeProviderError(data)
	output.ProviderCode = providerError.code
	output.ProviderType = providerError.providerType
	output.Param = providerError.param
	if providerError.message != "" {
		output.Message = providerError.message
	}
	if len(output.Message) > maxErrorMessage {
		output.Message = output.Message[:maxErrorMessage]
	}
	return output
}

type providerErrorFields struct {
	code         string
	providerType string
	param        string
	message      string
}

func decodeProviderError(data []byte) providerErrorFields {
	var payload struct {
		Error   json.RawMessage `json:"error"`
		Message string          `json:"message"`
	}
	if err := json.Unmarshal(data, &payload); err != nil {
		return providerErrorFields{}
	}
	var providerError struct {
		Code    string `json:"code"`
		Type    string `json:"type"`
		Param   string `json:"param"`
		Message string `json:"message"`
	}
	if err := json.Unmarshal(payload.Error, &providerError); err == nil && providerError.Message != "" {
		return providerErrorFields{code: providerError.Code, providerType: providerError.Type, param: providerError.Param, message: providerError.Message}
	}
	var providerMessage string
	if err := json.Unmarshal(payload.Error, &providerMessage); err == nil && providerMessage != "" {
		return providerErrorFields{message: providerMessage}
	}
	return providerErrorFields{message: payload.Message}
}

type cancelingBody struct {
	io.ReadCloser
	cancel context.CancelFunc
}

// Close closes the response body and cancels the stream context.
func (body *cancelingBody) Close() error {
	err := body.ReadCloser.Close()
	body.cancel()
	return err
}
