/*
 *  Copyright (c) 2026 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

// Package runs implements the OpenAI Runs bounded context.
package runs

import (
	"context"

	"go.osspkg.com/llm-client/openai/stateful"
	"go.osspkg.com/llm-client/pkg/transport"
)

// Request is a run request payload.
type Request = stateful.Request

// Resource is a run resource.
type Resource = stateful.Resource

// ListResponse is a cursor page of run resources.
type ListResponse = stateful.ListResponse

// Client calls run endpoints.
type Client struct{ stateful *stateful.Client }

// New creates a runs client.
func New(client *transport.Client) *Client { return &Client{stateful: stateful.New(client)} }

// Create creates a run in a thread.
func (client *Client) Create(ctx context.Context, threadID string, input Request) (Resource, error) {
	return client.stateful.CreateRun(ctx, threadID, input)
}

// Get returns a run.
func (client *Client) Get(ctx context.Context, threadID, runID string) (Resource, error) {
	return client.stateful.GetRun(ctx, threadID, runID)
}

// Cancel cancels a run.
func (client *Client) Cancel(ctx context.Context, threadID, runID string) (Resource, error) {
	return client.stateful.CancelRun(ctx, threadID, runID)
}

// List returns runs in a thread.
func (client *Client) List(ctx context.Context, threadID string) (ListResponse, error) {
	return client.stateful.ListRuns(ctx, threadID)
}

// Update updates a run in a thread.
func (client *Client) Update(ctx context.Context, threadID, runID string, input Request) (Resource, error) {
	return client.stateful.UpdateRun(ctx, threadID, runID, input)
}

// ListSteps returns steps for a run.
func (client *Client) ListSteps(ctx context.Context, threadID, runID string) (ListResponse, error) {
	return client.stateful.ListRunSteps(ctx, threadID, runID)
}

// GetStep returns a single run step.
func (client *Client) GetStep(ctx context.Context, threadID, runID, stepID string) (Resource, error) {
	return client.stateful.GetRunStep(ctx, threadID, runID, stepID)
}

// SubmitToolOutputs submits tool outputs to a run.
func (client *Client) SubmitToolOutputs(ctx context.Context, threadID, runID string, input Request) (Resource, error) {
	return client.stateful.SubmitToolOutputs(ctx, threadID, runID, input)
}
