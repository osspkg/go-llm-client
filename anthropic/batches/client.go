// Package batches provides typed batches API operations.
package batches

import (
	"context"
	"net/http"
	"net/url"
	"strconv"

	"go.osspkg.com/llm-client/anthropic/internal/request"
	"go.osspkg.com/llm-client/pkg/stream"
	"go.osspkg.com/llm-client/pkg/transport"
)

// Client describes the Client API value.
type Client struct{ transport *transport.Client }

// New creates a client on shared transport.
func New(client *transport.Client) *Client { return &Client{transport: client} }

// Create performs the Create operation.
func (client *Client) Create(ctx context.Context, input Request) (Batch, error) {
	var output Batch
	err := request.JSON(ctx, client.transport, http.MethodPost, "/messages/batches", "batches.create", input, &output)
	return output, err
}

// List performs the List operation.
func (client *Client) List(ctx context.Context, limit int, after, before string) (ListResponse, error) {
	q := url.Values{}
	if limit > 0 {
		q.Set("limit", strconv.Itoa(limit))
	}
	if after != "" {
		q.Set("after_id", after)
	}
	if before != "" {
		q.Set("before_id", before)
	}
	e := "/messages/batches"
	if s := q.Encode(); s != "" {
		e += "?" + s
	}
	var o ListResponse
	err := request.JSON(ctx, client.transport, http.MethodGet, e, "batches.list", nil, &o)
	return o, err
}

// Get performs the Get operation.
func (client *Client) Get(ctx context.Context, id string) (Batch, error) {
	var o Batch
	err := request.JSON(ctx, client.transport, http.MethodGet, "/messages/batches/"+url.PathEscape(id), "batches.get", nil, &o)
	return o, err
}

// Cancel performs the Cancel operation.
func (client *Client) Cancel(ctx context.Context, id string) (Batch, error) {
	var o Batch
	err := request.JSON(ctx, client.transport, http.MethodPost, "/messages/batches/"+url.PathEscape(id)+"/cancel", "batches.cancel", nil, &o)
	return o, err
}

// Results opens the bounded JSONL result stream for a completed batch.
func (client *Client) Results(ctx context.Context, id string) (stream.Iterator[Result], error) {
	body, err := request.Stream(ctx, client.transport, http.MethodGet, "/messages/batches/"+url.PathEscape(id)+"/results", "batches.results", nil, "application/jsonl")
	if err != nil {
		return nil, err
	}
	return stream.NewNDJSON(body, request.Decode[Result], 0), nil
}
