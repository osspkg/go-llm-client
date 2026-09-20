// Package fine_tuning implements OpenAI fine-tuning lifecycle operations.
package fine_tuning

import (
	"context"
	"net/http"
	"net/url"

	"go.osspkg.com/llm-client/openai/internal/request"
	"go.osspkg.com/llm-client/pkg/pagination"
	"go.osspkg.com/llm-client/pkg/stream"
	"go.osspkg.com/llm-client/pkg/transport"
)

// Client calls fine-tuning endpoints.
type Client struct{ transport *transport.Client }

// New creates a fine-tuning client on shared transport.
func New(client *transport.Client) *Client { return &Client{transport: client} }

// Create starts a fine-tuning job.
func (client *Client) Create(ctx context.Context, input CreateRequest) (Job, error) {
	var output Job
	err := request.JSON(ctx, client.transport, request.Post, "/fine_tuning/jobs", "fine_tuning.create", input, &output)
	return output, err
}

// List returns fine-tuning jobs.
func (client *Client) List(ctx context.Context) (ListResponse, error) {
	return client.ListWithParams(ctx, pagination.Params{})
}

// ListWithParams returns fine-tuning jobs using cursor parameters.
func (client *Client) ListWithParams(ctx context.Context, params pagination.Params) (ListResponse, error) {
	var output ListResponse
	err := request.JSON(ctx, client.transport, request.Get, withQuery("/fine_tuning/jobs", params.Values()), "fine_tuning.list", nil, &output)
	return output, err
}

// Get returns a fine-tuning job.
func (client *Client) Get(ctx context.Context, id string) (Job, error) {
	var output Job
	err := request.JSON(ctx, client.transport, request.Get, "/fine_tuning/jobs/"+url.PathEscape(id), "fine_tuning.get", nil, &output)
	return output, err
}

// Cancel cancels a fine-tuning job.
func (client *Client) Cancel(ctx context.Context, id string) (Job, error) {
	var output Job
	err := request.JSON(ctx, client.transport, http.MethodPost, "/fine_tuning/jobs/"+url.PathEscape(id)+"/cancel", "fine_tuning.cancel", nil, &output)
	return output, err
}

// ListEvents opens a bounded SSE event stream for a fine-tuning job.
func (client *Client) ListEvents(ctx context.Context, id string) (stream.Iterator[Event], error) {
	body, err := request.Stream(ctx, client.transport, request.Get, "/fine_tuning/jobs/"+url.PathEscape(id)+"/events", "fine_tuning.events", nil, "text/event-stream")
	if err != nil {
		return nil, err
	}
	return stream.NewSSE(body, request.Decode[Event], 0), nil
}

func withQuery(endpoint string, values url.Values) string {
	if len(values) == 0 {
		return endpoint
	}
	return endpoint + "?" + values.Encode()
}
