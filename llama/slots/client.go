// Package slots provides typed slots API operations.
package slots

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"go.osspkg.com/llm-client/llama/internal/request"
	llmerrors "go.osspkg.com/llm-client/pkg/errors"
	"go.osspkg.com/llm-client/pkg/transport"
)

// Client describes the Client API value.
type Client struct{ transport *transport.Client }

// New creates a client on shared transport.
func New(c *transport.Client) *Client { return &Client{c} }

// List performs the List operation.
func (c *Client) List(ctx context.Context) (ListResponse, error) {
	var o ListResponse
	err := request.JSON(ctx, c.transport, http.MethodGet, "/slots", "slots.list", nil, &o)
	return o, err
}

// Action performs the Action operation.
func (c *Client) Action(ctx context.Context, id int, action, filename string) (ActionResponse, error) {
	if id < 0 {
		return ActionResponse{}, fmt.Errorf("%w: slot id", llmerrors.ErrInvalidRequest)
	}
	if action != ActionSave && action != ActionRestore && action != ActionErase {
		return ActionResponse{}, fmt.Errorf("%w: unknown slot action", llmerrors.ErrInvalidRequest)
	}
	q := url.Values{"action": []string{action}}
	if filename != "" {
		q.Set("filename", filename)
	}
	var o ActionResponse
	err := request.JSON(ctx, c.transport, http.MethodPost, "/slots/"+strconv.Itoa(id)+"?"+q.Encode(), "slots."+action, nil, &o)
	return o, err
}

// Save stores a slot's state in filename.
func (c *Client) Save(ctx context.Context, id int, filename string) (ActionResponse, error) {
	return c.Action(ctx, id, ActionSave, filename)
}

// Restore loads a slot's state from filename.
func (c *Client) Restore(ctx context.Context, id int, filename string) (ActionResponse, error) {
	return c.Action(ctx, id, ActionRestore, filename)
}

// Erase removes a slot's cached state.
func (c *Client) Erase(ctx context.Context, id int) (ActionResponse, error) {
	return c.Action(ctx, id, ActionErase, "")
}
