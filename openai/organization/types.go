package organization

//go:generate easyjson -all types.go

// User describes an organization user.
type User struct {
	ID        string `json:"id"`
	Name      string `json:"name,omitempty"`
	Email     string `json:"email,omitempty"`
	Role      string `json:"role,omitempty"`
	CreatedAt int64  `json:"created_at,omitempty"`
}

// Project describes an organization project.
type Project struct {
	ID        string `json:"id"`
	Object    string `json:"object"`
	Name      string `json:"name,omitempty"`
	Status    string `json:"status,omitempty"`
	CreatedAt int64  `json:"created_at,omitempty"`
}

// CreateProjectRequest creates a project.
type CreateProjectRequest struct {
	Name string `json:"name"`
}

// UpdateProjectRequest changes a project name.
type UpdateProjectRequest struct {
	Name string `json:"name"`
}

// UpdateUserRequest changes an organization user's role.
type UpdateUserRequest struct {
	Role string `json:"role"`
}

// ListUsersResponse lists users.
type ListUsersResponse struct {
	Object  string `json:"object"`
	Data    []User `json:"data"`
	HasMore bool   `json:"has_more"`
}

// ListProjectsResponse lists projects.
type ListProjectsResponse struct {
	Object  string    `json:"object"`
	Data    []Project `json:"data"`
	HasMore bool      `json:"has_more"`
}

// DeleteResponse confirms deletion.
type DeleteResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Deleted bool   `json:"deleted"`
}

// APIKey describes a project API key.
type APIKey struct {
	ID        string `json:"id"`
	Name      string `json:"name,omitempty"`
	CreatedAt int64  `json:"created_at,omitempty"`
	Redacted  string `json:"redacted_value,omitempty"`
}

// ListAPIKeysResponse lists project API keys.
type ListAPIKeysResponse struct {
	Object string   `json:"object"`
	Data   []APIKey `json:"data"`
}

// Invite describes an organization invitation.
type Invite struct {
	ID        string   `json:"id"`
	Email     string   `json:"email,omitempty"`
	Role      string   `json:"role,omitempty"`
	Projects  []string `json:"projects,omitempty"`
	Status    string   `json:"status,omitempty"`
	CreatedAt int64    `json:"created_at,omitempty"`
}

// CreateInviteRequest creates an organization invitation.
type CreateInviteRequest struct {
	Email    string   `json:"email"`
	Role     string   `json:"role"`
	Projects []string `json:"projects,omitempty"`
}

// ListInvitesResponse lists organization invitations.
type ListInvitesResponse struct {
	Object  string   `json:"object"`
	Data    []Invite `json:"data"`
	HasMore bool     `json:"has_more"`
}
