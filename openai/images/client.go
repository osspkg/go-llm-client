// Package images implements OpenAI image operations.
package images

import (
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"strconv"

	"go.osspkg.com/llm-client/openai/internal/request"
	"go.osspkg.com/llm-client/pkg/stream"
	"go.osspkg.com/llm-client/pkg/transport"
)

// Client calls image endpoints.
type Client struct{ transport *transport.Client }

// New creates an image client.
func New(client *transport.Client) *Client { return &Client{transport: client} }

// Generate creates images.
func (client *Client) Generate(ctx context.Context, input Request) (Response, error) {
	var output Response
	err := request.JSON(ctx, client.transport, request.Post, "/images/generations", "images.generate", input, &output)
	return output, err
}

// GenerateStream creates images and returns their typed SSE events.
func (client *Client) GenerateStream(ctx context.Context, input Request) (stream.Iterator[StreamEvent], error) {
	input.Stream = true
	body, err := request.Stream(ctx, client.transport, request.Post, "/images/generations", "images.generate_stream", input, "text/event-stream")
	if err != nil {
		return nil, err
	}
	return stream.NewSSE(body, request.Decode[StreamEvent], 0), nil
}

// Edit edits an image using the provider's multipart request shape.
func (client *Client) Edit(ctx context.Context, input EditInput) (Response, error) {
	var output Response
	body, contentType, err := editBody(input)
	if err != nil {
		return output, err
	}
	data, err := request.Reader(ctx, client.transport, request.Post, "/images/edits", "images.edit", body, contentType)
	if err != nil {
		return output, err
	}
	return request.Decode[Response](data)
}

// EditStream edits images and returns their typed SSE events.
func (client *Client) EditStream(ctx context.Context, input EditInput) (stream.Iterator[StreamEvent], error) {
	input.Stream = true
	body, contentType, err := editBody(input)
	if err != nil {
		return nil, err
	}
	streamBody, err := request.StreamReader(ctx, client.transport, request.Post, "/images/edits", "images.edit_stream", body, contentType, "text/event-stream")
	if err != nil {
		return nil, err
	}
	return stream.NewSSE(streamBody, request.Decode[StreamEvent], 0), nil
}

// Variation creates a variation of an image.
func (client *Client) Variation(ctx context.Context, input VariationInput) (Response, error) {
	var output Response
	body, contentType, err := variationBody(input)
	if err != nil {
		return output, err
	}
	data, err := request.Reader(ctx, client.transport, request.Post, "/images/variations", "images.variation", body, contentType)
	if err != nil {
		return output, err
	}
	return request.Decode[Response](data)
}

func editBody(input EditInput) (io.ReadCloser, string, error) {
	if input.Filename == "" || input.Image == nil || input.Prompt == "" {
		return nil, "", errors.New("filename, image, and prompt are required")
	}
	return multipartBody(input.Filename, input.Image, input.Mask, func(writer *multipart.Writer) error {
		if err := writer.WriteField("prompt", input.Prompt); err != nil {
			return err
		}
		if input.Model != "" {
			if err := writer.WriteField("model", input.Model); err != nil {
				return err
			}
		}
		if input.N > 0 {
			if err := writer.WriteField("n", strconv.Itoa(input.N)); err != nil {
				return err
			}
		}
		if input.Size != "" {
			if err := writer.WriteField("size", input.Size); err != nil {
				return err
			}
		}
		if input.Stream {
			return writer.WriteField("stream", "true")
		}
		return nil
	}, "image", "mask")
}

func variationBody(input VariationInput) (io.ReadCloser, string, error) {
	if input.Filename == "" || input.Image == nil {
		return nil, "", errors.New("filename and image are required")
	}
	return multipartBody(input.Filename, input.Image, nil, func(writer *multipart.Writer) error {
		if input.Model != "" {
			if err := writer.WriteField("model", input.Model); err != nil {
				return err
			}
		}
		if input.N > 0 {
			if err := writer.WriteField("n", strconv.Itoa(input.N)); err != nil {
				return err
			}
		}
		if input.Size != "" {
			return writer.WriteField("size", input.Size)
		}
		return nil
	}, "image", "")
}

func multipartBody(filename string, image io.Reader, mask io.Reader, fields func(*multipart.Writer) error, imageField, maskField string) (io.ReadCloser, string, error) { //nolint:revive // multipart field parameters describe the wire contract.
	reader, pipeWriter := io.Pipe()
	multipartWriter := multipart.NewWriter(pipeWriter)
	go func() {
		err := fields(multipartWriter)
		if err == nil {
			var part io.Writer
			part, err = multipartWriter.CreateFormFile(imageField, filename)
			if err == nil {
				_, err = io.Copy(part, image)
			}
		}
		if err == nil && mask != nil {
			var part io.Writer
			part, err = multipartWriter.CreateFormFile(maskField, "mask.png")
			if err == nil {
				_, err = io.Copy(part, mask)
			}
		}
		if err == nil {
			err = multipartWriter.Close()
		}
		if err != nil {
			_ = pipeWriter.CloseWithError(fmt.Errorf("write image multipart: %w", err))
			return
		}
		_ = pipeWriter.Close()
	}()
	return reader, multipartWriter.FormDataContentType(), nil
}
