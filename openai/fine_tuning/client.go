// Package finetuning implements OpenAI fine-tuning lifecycle operations.
package finetuning

import (
	"context"
	"net/http"
	"net/url"

	"go.osspkg.com/llm-client/openai/internal/request"
	"go.osspkg.com/llm-client/pkg/pagination"
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

// ListEvents returns the cursor page of fine-tuning job events.
func (client *Client) ListEvents(ctx context.Context, id string) (ListEventsResponse, error) {
	return client.ListEventsWithParams(ctx, id, pagination.Params{})
}

// ListEventsWithParams returns fine-tuning job events using cursor parameters.
func (client *Client) ListEventsWithParams(ctx context.Context, id string, params pagination.Params) (ListEventsResponse, error) {
	var output ListEventsResponse
	err := request.JSON(ctx, client.transport, request.Get, withQuery("/fine_tuning/jobs/"+url.PathEscape(id)+"/events", params.Values()), "fine_tuning.events.list", nil, &output)
	return output, err
}

// ListCheckpoints returns the cursor page of checkpoints for a fine-tuning job.
func (client *Client) ListCheckpoints(ctx context.Context, id string) (ListCheckpointsResponse, error) {
	return client.ListCheckpointsWithParams(ctx, id, pagination.Params{})
}

// ListCheckpointsWithParams returns fine-tuning checkpoints using cursor parameters.
func (client *Client) ListCheckpointsWithParams(ctx context.Context, id string, params pagination.Params) (ListCheckpointsResponse, error) {
	var output ListCheckpointsResponse
	err := request.JSON(ctx, client.transport, request.Get, withQuery("/fine_tuning/jobs/"+url.PathEscape(id)+"/checkpoints", params.Values()), "fine_tuning.checkpoints.list", nil, &output)
	return output, err
}

// Pause pauses a queued or running fine-tuning job.
func (client *Client) Pause(ctx context.Context, id string) (Job, error) {
	return client.transition(ctx, id, "pause")
}

// Resume resumes a paused fine-tuning job.
func (client *Client) Resume(ctx context.Context, id string) (Job, error) {
	return client.transition(ctx, id, "resume")
}

// ListCheckpointPermissions returns permissions for a fine-tuned checkpoint.
func (client *Client) ListCheckpointPermissions(ctx context.Context, checkpoint string) (ListCheckpointPermissionsResponse, error) {
	var output ListCheckpointPermissionsResponse
	err := request.JSON(ctx, client.transport, request.Get, "/fine_tuning/checkpoints/"+url.PathEscape(checkpoint)+"/permissions", "fine_tuning.checkpoint_permissions.list", nil, &output)
	return output, err
}

// CreateCheckpointPermissions grants projects access to a fine-tuned checkpoint.
func (client *Client) CreateCheckpointPermissions(ctx context.Context, checkpoint string, input CreateCheckpointPermissionRequest) (ListCheckpointPermissionsResponse, error) {
	var output ListCheckpointPermissionsResponse
	err := request.JSON(ctx, client.transport, request.Post, "/fine_tuning/checkpoints/"+url.PathEscape(checkpoint)+"/permissions", "fine_tuning.checkpoint_permissions.create", input, &output)
	return output, err
}

// DeleteCheckpointPermission revokes a project's checkpoint permission.
func (client *Client) DeleteCheckpointPermission(ctx context.Context, checkpoint, permissionID string) (DeleteCheckpointPermissionResponse, error) {
	var output DeleteCheckpointPermissionResponse
	err := request.JSON(ctx, client.transport, request.Delete, "/fine_tuning/checkpoints/"+url.PathEscape(checkpoint)+"/permissions/"+url.PathEscape(permissionID), "fine_tuning.checkpoint_permissions.delete", nil, &output)
	return output, err
}

// RunGrader executes an upstream-defined alpha grader.
func (client *Client) RunGrader(ctx context.Context, input RunGraderRequest) (RunGraderResponse, error) {
	var output RunGraderResponse
	err := request.JSON(ctx, client.transport, request.Post, "/fine_tuning/alpha/graders/run", "fine_tuning.graders.run", input, &output)
	return output, err
}

// ValidateGrader validates an upstream-defined alpha grader.
func (client *Client) ValidateGrader(ctx context.Context, input ValidateGraderRequest) (ValidateGraderResponse, error) {
	var output ValidateGraderResponse
	err := request.JSON(ctx, client.transport, request.Post, "/fine_tuning/alpha/graders/validate", "fine_tuning.graders.validate", input, &output)
	return output, err
}

func (client *Client) transition(ctx context.Context, id, transition string) (Job, error) {
	var output Job
	err := request.JSON(ctx, client.transport, http.MethodPost, "/fine_tuning/jobs/"+url.PathEscape(id)+"/"+transition, "fine_tuning."+transition, nil, &output)
	return output, err
}

func withQuery(endpoint string, values url.Values) string {
	if len(values) == 0 {
		return endpoint
	}
	return endpoint + "?" + values.Encode()
}
