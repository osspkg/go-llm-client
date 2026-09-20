// Package responses implements the OpenAI Responses bounded context.
package responses

import (
	"context"
	"net/http"
	"net/url"

	"go.osspkg.com/llm-client/openai/internal/request"
	"go.osspkg.com/llm-client/pkg/stream"
	"go.osspkg.com/llm-client/pkg/transport"
)

// Client calls Responses endpoints.
type Client struct{ transport *transport.Client }

// New creates a Responses client on shared transport.
func New(client *transport.Client) *Client { return &Client{transport: client} }

// Create creates a non-streaming response.
func (client *Client) Create(ctx context.Context, input Request) (Response, error) {
	var output Response
	err := request.JSON(ctx, client.transport, request.Post, "/responses", "responses.create", input, &output)
	return output, err
}

// CreateStream creates a typed SSE response stream.
func (client *Client) CreateStream(ctx context.Context, input Request) (stream.Iterator[StreamEvent], error) {
	input.Stream = true
	body, err := request.Stream(ctx, client.transport, http.MethodPost, "/responses", "responses.create_stream", input, "text/event-stream")
	if err != nil {
		return nil, err
	}
	return stream.NewSSE(body, request.Decode[StreamEvent], 0), nil
}

// Get returns a stored response.
func (client *Client) Get(ctx context.Context, id string) (Response, error) {
	var output Response
	err := request.JSON(ctx, client.transport, request.Get, "/responses/"+url.PathEscape(id), "responses.get", nil, &output)
	return output, err
}

// Delete removes a stored response.
func (client *Client) Delete(ctx context.Context, id string) (DeleteResponse, error) {
	var output DeleteResponse
	err := request.JSON(ctx, client.transport, request.Delete, "/responses/"+url.PathEscape(id), "responses.delete", nil, &output)
	return output, err
}

// Cancel cancels an in-progress response.
func (client *Client) Cancel(ctx context.Context, id string) (Response, error) {
	var output Response
	err := request.JSON(ctx, client.transport, request.Post, "/responses/"+url.PathEscape(id)+"/cancel", "responses.cancel", nil, &output)
	return output, err
}

// InputItems returns the cursor page of items supplied to a response.
func (client *Client) InputItems(ctx context.Context, id string) (InputItemsResponse, error) {
	var output InputItemsResponse
	err := request.JSON(ctx, client.transport, request.Get, "/responses/"+url.PathEscape(id)+"/input_items", "responses.input_items.list", nil, &output)
	return output, err
}
