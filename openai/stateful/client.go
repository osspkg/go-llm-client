/*
 *  Copyright (c) 2026 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

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

// UpdateAssistant updates an assistant resource.
func (client *Client) UpdateAssistant(ctx context.Context, id string, input Request) (Resource, error) {
	return client.create(ctx, "/assistants/"+url.PathEscape(id), "assistants.update", input)
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

// UpdateThread updates a thread resource.
func (client *Client) UpdateThread(ctx context.Context, id string, input Request) (Resource, error) {
	return client.create(ctx, "/threads/"+url.PathEscape(id), "threads.update", input)
}

// CreateThreadAndRun creates a thread and starts its initial run.
func (client *Client) CreateThreadAndRun(ctx context.Context, input Request) (Resource, error) {
	return client.create(ctx, "/threads/runs", "threads.runs.create", input)
}

// ListMessages returns messages in a thread.
func (client *Client) ListMessages(ctx context.Context, threadID string) (ListResponse, error) {
	return client.list(ctx, "/threads/"+url.PathEscape(threadID)+"/messages", "threads.messages.list")
}

// CreateMessage adds a message to a thread.
func (client *Client) CreateMessage(ctx context.Context, threadID string, input Request) (Resource, error) {
	return client.create(ctx, "/threads/"+url.PathEscape(threadID)+"/messages", "threads.messages.create", input)
}

// GetMessage returns a message in a thread.
func (client *Client) GetMessage(ctx context.Context, threadID, messageID string) (Resource, error) {
	return client.get(ctx, "/threads/"+url.PathEscape(threadID)+"/messages/"+url.PathEscape(messageID), "threads.messages.get")
}

// UpdateMessage updates a message in a thread.
func (client *Client) UpdateMessage(ctx context.Context, threadID, messageID string, input Request) (Resource, error) {
	return client.create(ctx, "/threads/"+url.PathEscape(threadID)+"/messages/"+url.PathEscape(messageID), "threads.messages.update", input)
}

// DeleteMessage removes a message from a thread.
func (client *Client) DeleteMessage(ctx context.Context, threadID, messageID string) (DeleteResponse, error) {
	return client.delete(ctx, "/threads/"+url.PathEscape(threadID)+"/messages/"+url.PathEscape(messageID), "threads.messages.delete")
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

// ListRuns returns runs in a thread.
func (client *Client) ListRuns(ctx context.Context, threadID string) (ListResponse, error) {
	return client.list(ctx, "/threads/"+url.PathEscape(threadID)+"/runs", "runs.list")
}

// UpdateRun updates a run in a thread.
func (client *Client) UpdateRun(ctx context.Context, threadID, runID string, input Request) (Resource, error) {
	return client.create(ctx, "/threads/"+url.PathEscape(threadID)+"/runs/"+url.PathEscape(runID), "runs.update", input)
}

// ListRunSteps returns steps for a run.
func (client *Client) ListRunSteps(ctx context.Context, threadID, runID string) (ListResponse, error) {
	return client.list(ctx, "/threads/"+url.PathEscape(threadID)+"/runs/"+url.PathEscape(runID)+"/steps", "runs.steps.list")
}

// GetRunStep returns a single run step.
func (client *Client) GetRunStep(ctx context.Context, threadID, runID, stepID string) (Resource, error) {
	return client.get(ctx, "/threads/"+url.PathEscape(threadID)+"/runs/"+url.PathEscape(runID)+"/steps/"+url.PathEscape(stepID), "runs.steps.get")
}

// SubmitToolOutputs submits tool outputs to a run.
func (client *Client) SubmitToolOutputs(ctx context.Context, threadID, runID string, input Request) (Resource, error) {
	return client.create(ctx, "/threads/"+url.PathEscape(threadID)+"/runs/"+url.PathEscape(runID)+"/submit_tool_outputs", "runs.tool_outputs.submit", input)
}

// CreateVectorStore creates a vector store.
func (client *Client) CreateVectorStore(ctx context.Context, input Request) (Resource, error) {
	return client.create(ctx, "/vector_stores", "vector_stores.create", input)
}

// ListVectorStores lists vector stores.
func (client *Client) ListVectorStores(ctx context.Context) (ListResponse, error) {
	return client.list(ctx, "/vector_stores", "vector_stores.list")
}

// GetVectorStore returns a vector store.
func (client *Client) GetVectorStore(ctx context.Context, id string) (Resource, error) {
	return client.get(ctx, "/vector_stores/"+url.PathEscape(id), "vector_stores.get")
}

// UpdateVectorStore updates a vector store.
func (client *Client) UpdateVectorStore(ctx context.Context, id string, input Request) (Resource, error) {
	return client.create(ctx, "/vector_stores/"+url.PathEscape(id), "vector_stores.update", input)
}

// DeleteVectorStore deletes a vector store.
func (client *Client) DeleteVectorStore(ctx context.Context, id string) (DeleteResponse, error) {
	return client.delete(ctx, "/vector_stores/"+url.PathEscape(id), "vector_stores.delete")
}

// CreateContainer creates a container resource.
func (client *Client) CreateContainer(ctx context.Context, input Request) (Resource, error) {
	return client.create(ctx, "/containers", "containers.create", input)
}

// ListContainers returns containers visible to the current project.
func (client *Client) ListContainers(ctx context.Context) (ListResponse, error) {
	return client.list(ctx, "/containers", "containers.list")
}

// GetContainer returns a container resource.
func (client *Client) GetContainer(ctx context.Context, id string) (Resource, error) {
	return client.get(ctx, "/containers/"+url.PathEscape(id), "containers.get")
}

// DeleteContainer deletes a container resource.
func (client *Client) DeleteContainer(ctx context.Context, id string) (DeleteResponse, error) {
	return client.delete(ctx, "/containers/"+url.PathEscape(id), "containers.delete")
}

// CreateEval creates an evaluation resource.
func (client *Client) CreateEval(ctx context.Context, input Request) (Resource, error) {
	return client.create(ctx, "/evals", "evals.create", input)
}

// ListEvals returns evaluations visible to the current project.
func (client *Client) ListEvals(ctx context.Context) (ListResponse, error) {
	return client.list(ctx, "/evals", "evals.list")
}

// GetEval returns an evaluation resource.
func (client *Client) GetEval(ctx context.Context, id string) (Resource, error) {
	return client.get(ctx, "/evals/"+url.PathEscape(id), "evals.get")
}

// UpdateEval updates an evaluation resource.
func (client *Client) UpdateEval(ctx context.Context, id string, input Request) (Resource, error) {
	return client.create(ctx, "/evals/"+url.PathEscape(id), "evals.update", input)
}

// DeleteEval deletes an evaluation resource.
func (client *Client) DeleteEval(ctx context.Context, id string) (DeleteResponse, error) {
	return client.delete(ctx, "/evals/"+url.PathEscape(id), "evals.delete")
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
