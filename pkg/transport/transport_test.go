package transport_test

import (
	"context"
	"errors"
	"io"
	"mime/multipart"
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
	var metadata auth.RequestMeta
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
		transport.WithAuthProvider(func(ctx context.Context, value auth.RequestMeta) (http.Header, error) {
			metadata = value
			return auth.StaticBearer("secret")(ctx, value)
		}),
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
	if metadata.Domain != "example.test" {
		t.Errorf("auth domain = %q, want example.test", metadata.Domain)
	}
	if metadata.URL != "https://example.test/v1/models" {
		t.Errorf("auth url = %q", metadata.URL)
	}
}

func TestWebSocketHeadersIncludeDestinationDomain(t *testing.T) {
	var metadata auth.RequestMeta
	client, err := transport.New("https://realtime.example.test/v1", transport.WithAuthProvider(func(_ context.Context, value auth.RequestMeta) (http.Header, error) {
		metadata = value
		return make(http.Header), nil
	}))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := client.WebSocketHeaders(t.Context(), "realtime.connect"); err != nil {
		t.Fatal(err)
	}
	if metadata.Domain != "realtime.example.test" {
		t.Errorf("auth domain = %q, want realtime.example.test", metadata.Domain)
	}
	if !metadata.WebSocket {
		t.Error("websocket metadata flag is false")
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

func TestBaseURLRejectsCredentialsAndNonHTTPSchemes(t *testing.T) {
	for _, value := range []string{"https://user:secret@example.test", "ftp://example.test"} {
		t.Run(value, func(t *testing.T) {
			_, err := transport.New(value)
			if !errors.Is(err, llmerrors.ErrInvalidConfig) {
				t.Fatalf("error = %v", err)
			}
		})
	}
}

func TestMultipartBodyWaitsForWriterAndReportsWriterError(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		body, contentType := transport.NewMultipartBody(func(writer *multipart.Writer) error {
			return writer.WriteField("name", "value")
		})
		if contentType == "" {
			t.Fatal("content type is empty")
		}
		data, err := io.ReadAll(body)
		if err != nil {
			t.Fatal(err)
		}
		if err := body.Close(); err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(string(data), "name") || !strings.Contains(string(data), "value") {
			t.Fatalf("multipart body = %q", data)
		}
	})

	t.Run("writer error", func(t *testing.T) {
		wantErr := errors.New("source failed")
		body, _ := transport.NewMultipartBody(func(*multipart.Writer) error { return wantErr })
		_, readErr := io.ReadAll(body)
		if !errors.Is(readErr, wantErr) {
			t.Fatalf("read error = %v, want %v", readErr, wantErr)
		}
		if closeErr := body.Close(); !errors.Is(closeErr, wantErr) {
			t.Fatalf("close error = %v, want %v", closeErr, wantErr)
		}
	})
}
