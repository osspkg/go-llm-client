// Package files implements OpenAI file operations.
package files

import (
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"

	"go.osspkg.com/llm-client/openai/internal/request"
	"go.osspkg.com/llm-client/pkg/transport"
)

// Client calls file endpoints.
type Client struct{ transport *transport.Client }

// New creates a file client.
func New(client *transport.Client) *Client { return &Client{transport: client} }

// UploadInput describes a file upload.
type UploadInput struct {
	Filename string
	Purpose  string
	File     io.Reader
}

// List returns uploaded files.
func (client *Client) List(ctx context.Context) (ListResponse, error) {
	var output ListResponse
	err := request.JSON(ctx, client.transport, request.Get, "/files", "files.list", nil, &output)
	return output, err
}

// Upload uploads a file.
func (client *Client) Upload(ctx context.Context, input UploadInput) (File, error) {
	var output File
	body, contentType, err := uploadBody(input)
	if err != nil {
		return output, err
	}
	data, err := request.Reader(ctx, client.transport, request.Post, "/files", "files.upload", body, contentType)
	if err != nil {
		return output, err
	}
	return request.Decode[File](data)
}

// Get returns file metadata.
func (client *Client) Get(ctx context.Context, id string) (File, error) {
	var output File
	err := request.JSON(ctx, client.transport, request.Get, "/files/"+url.PathEscape(id), "files.get", nil, &output)
	return output, err
}

// Delete deletes a file.
func (client *Client) Delete(ctx context.Context, id string) (DeleteResponse, error) {
	var output DeleteResponse
	err := request.JSON(ctx, client.transport, http.MethodDelete, "/files/"+url.PathEscape(id), "files.delete", nil, &output)
	return output, err
}

// Content downloads file content with the configured response limit.
func (client *Client) Content(ctx context.Context, id string) ([]byte, error) {
	return request.Bytes(ctx, client.transport, request.Get, "/files/"+url.PathEscape(id)+"/content", "files.content", nil, "")
}

func uploadBody(input UploadInput) (io.ReadCloser, string, error) {
	if input.Filename == "" || input.Purpose == "" || input.File == nil {
		return nil, "", errors.New("filename, purpose, and file are required")
	}
	reader, writer := io.Pipe()
	multipartWriter := multipart.NewWriter(writer)
	go func() {
		err := multipartWriter.WriteField("purpose", input.Purpose)
		if err == nil {
			var part io.Writer
			part, err = multipartWriter.CreateFormFile("file", input.Filename)
			if err == nil {
				_, err = io.Copy(part, input.File)
			}
		}
		if err == nil {
			err = multipartWriter.Close()
		}
		if err != nil {
			_ = writer.CloseWithError(fmt.Errorf("write file multipart: %w", err))
			return
		}
		_ = writer.Close()
	}()
	return reader, writerBoundary(multipartWriter), nil
}

func writerBoundary(writer *multipart.Writer) string { return writer.FormDataContentType() }
