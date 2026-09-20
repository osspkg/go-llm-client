// Package stateful implements OpenAI assistant, thread, run, vector-store,
// container, and evaluation resource operations.
package stateful

import (
	"context"
	"net/http"
	"net/url"

	"go.osspkg.com/llm-client/openai/internal/request"
	"go.osspkg.com/llm-client/pkg/transport"
)

// Client calls stateful OpenAI domains through explicit operation methods.
type Client struct{ transport *transport.Client }

// New creates a stateful client.
func New(client *transport.Client) *Client { return &Client{transport: client} }

// CreateAssistant creates an assistant.
func (client *Client) CreateAssistant(ctx context.Context, input Request) (Resource, error) {
	return client.create(ctx, "/assistants", "assistants.create", input)
}

// GetAssistant returns an assistant.
func (client *Client) GetAssistant(ctx context.Context, id string) (Resource, error) {
	return client.get(ctx, "/assistants/"+url.PathEscape(id), "assistants.get")
}

// ListAssistants lists assistants.
func (client *Client) ListAssistants(ctx context.Context) (ListResponse, error) {
	return client.list(ctx, "/assistants", "assistants.list")
}

// DeleteAssistant deletes an assistant.
func (client *Client) DeleteAssistant(ctx context.Context, id string) (DeleteResponse, error) {
	return client.delete(ctx, "/assistants/"+url.PathEscape(id), "assistants.delete")
}

// CreateThread creates a thread.
func (client *Client) CreateThread(ctx context.Context, input Request) (Resource, error) {
	return client.create(ctx, "/threads", "threads.create", input)
}

// GetThread returns a thread.
func (client *Client) GetThread(ctx context.Context, id string) (Resource, error) {
	return client.get(ctx, "/threads/"+url.PathEscape(id), "threads.get")
}

// DeleteThread deletes a thread.
func (client *Client) DeleteThread(ctx context.Context, id string) (DeleteResponse, error) {
	return client.delete(ctx, "/threads/"+url.PathEscape(id), "threads.delete")
}

// CreateRun creates a run in a thread.
func (client *Client) CreateRun(ctx context.Context, threadID string, input Request) (Resource, error) {
	return client.create(ctx, "/threads/"+url.PathEscape(threadID)+"/runs", "runs.create", input)
}

// GetRun returns a run.
func (client *Client) GetRun(ctx context.Context, threadID, runID string) (Resource, error) {
	return client.get(ctx, "/threads/"+url.PathEscape(threadID)+"/runs/"+url.PathEscape(runID), "runs.get")
}

// CancelRun cancels a run.
func (client *Client) CancelRun(ctx context.Context, threadID, runID string) (Resource, error) {
	return client.post(ctx, "/threads/"+url.PathEscape(threadID)+"/runs/"+url.PathEscape(runID)+"/cancel", "runs.cancel")
}

// CreateVectorStore creates a vector store.
func (client *Client) CreateVectorStore(ctx context.Context, input Request) (Resource, error) {
	return client.create(ctx, "/vector_stores", "vector_stores.create", input)
}

// ListVectorStores lists vector stores.
func (client *Client) ListVectorStores(ctx context.Context) (ListResponse, error) {
	return client.list(ctx, "/vector_stores", "vector_stores.list")
}

// CreateContainer creates a container resource.
func (client *Client) CreateContainer(ctx context.Context, input Request) (Resource, error) {
	return client.create(ctx, "/containers", "containers.create", input)
}

// CreateEval creates an evaluation resource.
func (client *Client) CreateEval(ctx context.Context, input Request) (Resource, error) {
	return client.create(ctx, "/evals", "evals.create", input)
}

func (client *Client) create(ctx context.Context, endpoint, operation string, input Request) (Resource, error) {
	var output Resource
	err := request.JSON(ctx, client.transport, request.Post, endpoint, operation, input, &output)
	return output, err
}

func (client *Client) get(ctx context.Context, endpoint, operation string) (Resource, error) {
	var output Resource
	err := request.JSON(ctx, client.transport, request.Get, endpoint, operation, nil, &output)
	return output, err
}

func (client *Client) post(ctx context.Context, endpoint, operation string) (Resource, error) {
	var output Resource
	err := request.JSON(ctx, client.transport, http.MethodPost, endpoint, operation, nil, &output)
	return output, err
}

func (client *Client) list(ctx context.Context, endpoint, operation string) (ListResponse, error) {
	var output ListResponse
	err := request.JSON(ctx, client.transport, request.Get, endpoint, operation, nil, &output)
	return output, err
}

func (client *Client) delete(ctx context.Context, endpoint, operation string) (DeleteResponse, error) {
	var output DeleteResponse
	err := request.JSON(ctx, client.transport, http.MethodDelete, endpoint, operation, nil, &output)
	return output, err
}
