// Package files provides typed files API operations.
package files

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"net/url"
	"strconv"
	"strings"

	"go.osspkg.com/llm-client/anthropic/internal/request"
	llmerrors "go.osspkg.com/llm-client/pkg/errors"
	"go.osspkg.com/llm-client/pkg/transport"
)

// Client describes the Client API value.
type Client struct{ transport *transport.Client }

// New creates a client on shared transport.
func New(client *transport.Client) *Client { return &Client{transport: client} }

// List performs the List operation.
func (client *Client) List(ctx context.Context, limit int, after, before string) (ListResponse, error) {
	query := url.Values{}
	if limit > 0 {
		query.Set("limit", strconv.Itoa(limit))
	}
	if after != "" {
		query.Set("after_id", after)
	}
	if before != "" {
		query.Set("before_id", before)
	}
	endpoint := "/files"
	if encoded := query.Encode(); encoded != "" {
		endpoint += "?" + encoded
	}
	var output ListResponse
	err := request.JSON(ctx, client.transport, request.Get, endpoint, "files.list", nil, &output)
	return output, err
}

// Get performs the Get operation.
func (client *Client) Get(ctx context.Context, id string) (File, error) {
	var output File
	err := request.JSON(ctx, client.transport, request.Get, "/files/"+url.PathEscape(id), "files.get", nil, &output)
	return output, err
}

// Delete performs the Delete operation.
func (client *Client) Delete(ctx context.Context, id string) (DeleteResponse, error) {
	var output DeleteResponse
	err := request.JSON(ctx, client.transport, request.Delete, "/files/"+url.PathEscape(id), "files.delete", nil, &output)
	return output, err
}

// Content performs the Content operation.
func (client *Client) Content(ctx context.Context, id string) (io.ReadCloser, error) {
	return client.transport.Stream(ctx, http.MethodGet, "/files/"+url.PathEscape(id)+"/content", "files.content", nil, "", "application/octet-stream")
}

// Upload performs the Upload operation.
func (client *Client) Upload(ctx context.Context, filename, mediaType string, body io.Reader) (File, error) {
	if strings.TrimSpace(filename) == "" || strings.ContainsAny(filename, "\r\n\"") || strings.ContainsAny(mediaType, "\r\n") || body == nil {
		return File{}, fmt.Errorf("%w: invalid multipart upload", llmerrors.ErrInvalidRequest)
	}
	reader, writer := io.Pipe()
	multipartWriter := multipart.NewWriter(writer)
	go func() {
		defer func() { _ = writer.Close() }()
		header := make(textproto.MIMEHeader)
		header.Set("Content-Disposition", `form-data; name="file"; filename="`+filename+`"`)
		if mediaType != "" {
			header.Set("Content-Type", mediaType)
		}
		part, err := multipartWriter.CreatePart(header)
		if err == nil {
			_, err = io.Copy(part, body)
		}
		if err == nil {
			err = multipartWriter.WriteField("purpose", "user_data")
		}
		if err == nil {
			err = multipartWriter.Close()
		}
		if err != nil {
			_ = writer.CloseWithError(err)
		}
	}()
	data, err := client.transport.RequestReader(ctx, http.MethodPost, "/files", "files.upload", reader, multipartWriter.FormDataContentType())
	if err != nil {
		return File{}, err
	}
	return request.Decode[File](data)
}
