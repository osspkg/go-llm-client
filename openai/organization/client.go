// Package organization implements OpenAI organization and project operations.
package organization

import (
	"context"
	"net/http"
	"net/url"

	"go.osspkg.com/llm-client/openai/internal/request"
	"go.osspkg.com/llm-client/pkg/transport"
)

// Client calls organization endpoints.
type Client struct{ transport *transport.Client }

// New creates an organization client.
func New(client *transport.Client) *Client { return &Client{transport: client} }

// ListUsers lists organization users.
func (client *Client) ListUsers(ctx context.Context) (ListUsersResponse, error) {
	var output ListUsersResponse
	err := request.JSON(ctx, client.transport, request.Get, "/organization/users", "organization.users.list", nil, &output)
	return output, err
}

// GetUser returns an organization user.
func (client *Client) GetUser(ctx context.Context, id string) (User, error) {
	var output User
	err := request.JSON(ctx, client.transport, request.Get, "/organization/users/"+url.PathEscape(id), "organization.users.get", nil, &output)
	return output, err
}

// UpdateUser changes an organization user's role.
func (client *Client) UpdateUser(ctx context.Context, id string, input UpdateUserRequest) (User, error) {
	var output User
	err := request.JSON(ctx, client.transport, http.MethodPost, "/organization/users/"+url.PathEscape(id), "organization.users.update", input, &output)
	return output, err
}

// DeleteUser removes an organization user.
func (client *Client) DeleteUser(ctx context.Context, id string) (DeleteResponse, error) {
	var output DeleteResponse
	err := request.JSON(ctx, client.transport, http.MethodDelete, "/organization/users/"+url.PathEscape(id), "organization.users.delete", nil, &output)
	return output, err
}

// ListProjects lists organization projects.
func (client *Client) ListProjects(ctx context.Context) (ListProjectsResponse, error) {
	var output ListProjectsResponse
	err := request.JSON(ctx, client.transport, request.Get, "/organization/projects", "organization.projects.list", nil, &output)
	return output, err
}

// CreateProject creates a project.
func (client *Client) CreateProject(ctx context.Context, name string) (Project, error) {
	input := CreateProjectRequest{Name: name}
	var output Project
	err := request.JSON(ctx, client.transport, request.Post, "/organization/projects", "organization.projects.create", input, &output)
	return output, err
}

// GetProject returns a project.
func (client *Client) GetProject(ctx context.Context, id string) (Project, error) {
	var output Project
	err := request.JSON(ctx, client.transport, request.Get, "/organization/projects/"+url.PathEscape(id), "organization.projects.get", nil, &output)
	return output, err
}

// UpdateProject changes a project name.
func (client *Client) UpdateProject(ctx context.Context, id string, input UpdateProjectRequest) (Project, error) {
	var output Project
	err := request.JSON(ctx, client.transport, http.MethodPost, "/organization/projects/"+url.PathEscape(id), "organization.projects.update", input, &output)
	return output, err
}

// DeleteProject deletes a project.
func (client *Client) DeleteProject(ctx context.Context, id string) (DeleteResponse, error) {
	var output DeleteResponse
	err := request.JSON(ctx, client.transport, http.MethodDelete, "/organization/projects/"+url.PathEscape(id), "organization.projects.delete", nil, &output)
	return output, err
}

// ListProjectAPIKeys lists API keys for a project.
func (client *Client) ListProjectAPIKeys(ctx context.Context, projectID string) (ListAPIKeysResponse, error) {
	var output ListAPIKeysResponse
	err := request.JSON(ctx, client.transport, request.Get, "/organization/projects/"+url.PathEscape(projectID)+"/api_keys", "organization.project_api_keys.list", nil, &output)
	return output, err
}

// DeleteProjectAPIKey deletes a project API key.
func (client *Client) DeleteProjectAPIKey(ctx context.Context, projectID, keyID string) (DeleteResponse, error) {
	var output DeleteResponse
	err := request.JSON(ctx, client.transport, http.MethodDelete, "/organization/projects/"+url.PathEscape(projectID)+"/api_keys/"+url.PathEscape(keyID), "organization.project_api_keys.delete", nil, &output)
	return output, err
}

// ListInvites lists organization invitations.
func (client *Client) ListInvites(ctx context.Context) (ListInvitesResponse, error) {
	var output ListInvitesResponse
	err := request.JSON(ctx, client.transport, request.Get, "/organization/invites", "organization.invites.list", nil, &output)
	return output, err
}

// CreateInvite creates an organization invitation.
func (client *Client) CreateInvite(ctx context.Context, input CreateInviteRequest) (Invite, error) {
	var output Invite
	err := request.JSON(ctx, client.transport, request.Post, "/organization/invites", "organization.invites.create", input, &output)
	return output, err
}

// GetInvite returns an organization invitation.
func (client *Client) GetInvite(ctx context.Context, id string) (Invite, error) {
	var output Invite
	err := request.JSON(ctx, client.transport, request.Get, "/organization/invites/"+url.PathEscape(id), "organization.invites.get", nil, &output)
	return output, err
}

// DeleteInvite deletes an organization invitation.
func (client *Client) DeleteInvite(ctx context.Context, id string) (DeleteResponse, error) {
	var output DeleteResponse
	err := request.JSON(ctx, client.transport, http.MethodDelete, "/organization/invites/"+url.PathEscape(id), "organization.invites.delete", nil, &output)
	return output, err
}
