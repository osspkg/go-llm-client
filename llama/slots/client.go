// Package slots provides typed slots API operations.
package slots

import (
	"context"
	"net/http"
	"net/url"
	"strconv"

	"go.osspkg.com/llm-client/llama/internal/request"
	llmerrors "go.osspkg.com/llm-client/pkg/errors"
	"go.osspkg.com/llm-client/pkg/transport"
)

// Client calls native slot operations.
type Client struct{ transport *transport.Client }

// New creates a client on shared transport.
func New(client *transport.Client) *Client { return &Client{transport: client} }

// List returns the current state of all server slots.
func (client *Client) List(ctx context.Context) (ListResponse, error) {
	var output ListResponse
	err := request.JSON(ctx, client.transport, http.MethodGet, "/slots", "slots.list", nil, &output)
	return output, err
}

func (client *Client) action(ctx context.Context, id int, action, filename string) (ActionResponse, error) {
	if id < 0 {
		return ActionResponse{}, &llmerrors.ValidationError{Field: "id", Msg: "must be non-negative"}
	}
	if action != ActionSave && action != ActionRestore && action != ActionErase {
		return ActionResponse{}, &llmerrors.ValidationError{Field: "action", Msg: "must be save, restore, or erase"}
	}
	q := url.Values{"action": []string{action}}
	if filename != "" {
		q.Set("filename", filename)
	}
	var output ActionResponse
	err := request.JSON(ctx, client.transport, http.MethodPost, "/slots/"+strconv.Itoa(id)+"?"+q.Encode(), "slots."+action, nil, &output)
	return output, err
}

// Save stores a slot's state in filename.
func (client *Client) Save(ctx context.Context, id int, filename string) (ActionResponse, error) {
	return client.action(ctx, id, ActionSave, filename)
}

// Restore loads a slot's state from filename.
func (client *Client) Restore(ctx context.Context, id int, filename string) (ActionResponse, error) {
	return client.action(ctx, id, ActionRestore, filename)
}

// Erase removes a slot's cached state.
func (client *Client) Erase(ctx context.Context, id int) (ActionResponse, error) {
	return client.action(ctx, id, ActionErase, "")
}
