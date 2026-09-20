package files_test

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"go.osspkg.com/llm-client/openai/files"
	"go.osspkg.com/llm-client/pkg/transport"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (fn roundTripFunc) RoundTrip(request *http.Request) (*http.Response, error) { return fn(request) }

func TestUploadStreamsBoundedMultipartBody(t *testing.T) {
	client, err := transport.New("https://example.test/v1", transport.WithHTTPClient(&http.Client{Transport: roundTripFunc(func(request *http.Request) (*http.Response, error) {
		if request.Header.Get("Content-Type") == "" || !strings.HasPrefix(request.Header.Get("Content-Type"), "multipart/form-data;") {
			t.Errorf("content type = %q", request.Header.Get("Content-Type"))
		}
		body, readErr := io.ReadAll(request.Body)
		if readErr != nil {
			t.Errorf("read multipart body: %v", readErr)
		}
		if !strings.Contains(string(body), "purpose") || !strings.Contains(string(body), "file-content") {
			t.Errorf("multipart body = %q", body)
		}
		return &http.Response{StatusCode: http.StatusOK, Header: make(http.Header), Body: io.NopCloser(strings.NewReader(`{"id":"file-1","filename":"test.json"}`))}, nil
	})}))
	if err != nil {
		t.Fatal(err)
	}
	output, err := files.New(client).Upload(context.Background(), files.UploadInput{
		Filename: "test.json",
		Purpose:  "fine-tune",
		File:     strings.NewReader("file-content"),
	})
	if err != nil {
		t.Fatal(err)
	}
	if output.ID != "file-1" {
		t.Fatalf("output = %#v", output)
	}
}
