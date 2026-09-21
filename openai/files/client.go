// Package files implements OpenAI file operations.
package files

import (
	"context"
	"errors"
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
	data, requestErr := request.Reader(ctx, client.transport, request.Post, "/files", "files.upload", body, contentType)
	bodyErr := body.Close()
	if requestErr != nil {
		return output, requestErr
	}
	if bodyErr != nil {
		return output, bodyErr
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

func uploadBody(input UploadInput) (*transport.MultipartBody, string, error) {
	if input.Filename == "" || input.Purpose == "" || input.File == nil {
		return nil, "", errors.New("filename, purpose, and file are required")
	}
	body, contentType := transport.NewMultipartBody(func(multipartWriter *multipart.Writer) error {
		if err := multipartWriter.WriteField("purpose", input.Purpose); err != nil {
			return err
		}
		part, err := multipartWriter.CreateFormFile("file", input.Filename)
		if err != nil {
			return err
		}
		_, err = io.Copy(part, input.File)
		return err
	})
	return body, contentType, nil
}
