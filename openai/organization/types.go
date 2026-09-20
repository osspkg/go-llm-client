package organization

//go:generate easyjson -all types.go

// User describes an organization user.
type User struct {
	// ID identifies the resource, event, or request.
	ID string `json:"id"`
	// Name identifies the function, model, project, or resource by name.
	Name string `json:"name,omitempty"`
	// Email contains the invite or user email address.
	Email string `json:"email,omitempty"`
	// Role identifies the organization or message role.
	Role string `json:"role,omitempty"`
	// CreatedAt contains the Unix creation timestamp.
	CreatedAt int64 `json:"created_at,omitempty"`
}

// Project describes an organization project.
type Project struct {
	// ID identifies the resource, event, or request.
	ID string `json:"id"`
	// Object identifies the provider resource or collection type.
	Object string `json:"object"`
	// Name identifies the function, model, project, or resource by name.
	Name string `json:"name,omitempty"`
	// Status reports the current provider processing or lifecycle state.
	Status string `json:"status,omitempty"`
	// CreatedAt contains the Unix creation timestamp.
	CreatedAt int64 `json:"created_at,omitempty"`
}

// CreateProjectRequest creates a project.
type CreateProjectRequest struct {
	// Name identifies the function, model, project, or resource by name.
	Name string `json:"name"`
}

// UpdateProjectRequest changes a project name.
type UpdateProjectRequest struct {
	// Name identifies the function, model, project, or resource by name.
	Name string `json:"name"`
}

// UpdateUserRequest changes an organization user's role.
type UpdateUserRequest struct {
	// Role identifies the organization or message role.
	Role string `json:"role"`
}

// ListUsersResponse lists users.
type ListUsersResponse struct {
	// Object identifies the provider resource or collection type.
	Object string `json:"object"`
	// Data contains the returned records.
	Data []User `json:"data"`
	// HasMore reports whether another page is available.
	HasMore bool `json:"has_more"`
}

// ListProjectsResponse lists projects.
type ListProjectsResponse struct {
	// Object identifies the provider resource or collection type.
	Object string `json:"object"`
	// Data contains the returned records.
	Data []Project `json:"data"`
	// HasMore reports whether another page is available.
	HasMore bool `json:"has_more"`
}

// DeleteResponse confirms deletion.
type DeleteResponse struct {
	// ID identifies the resource, event, or request.
	ID string `json:"id"`
	// Object identifies the provider resource or collection type.
	Object string `json:"object"`
	// Deleted confirms that the provider deleted the resource.
	Deleted bool `json:"deleted"`
}

// APIKey describes a project API key.
type APIKey struct {
	// ID identifies the resource, event, or request.
	ID string `json:"id"`
	// Name identifies the function, model, project, or resource by name.
	Name string `json:"name,omitempty"`
	// CreatedAt contains the Unix creation timestamp.
	CreatedAt int64 `json:"created_at,omitempty"`
	// Redacted contains the provider-redacted API key value.
	Redacted string `json:"redacted_value,omitempty"`
}

// ListAPIKeysResponse lists project API keys.
type ListAPIKeysResponse struct {
	// Object identifies the provider resource or collection type.
	Object string `json:"object"`
	// Data contains the returned records.
	Data []APIKey `json:"data"`
}

// Invite describes an organization invitation.
type Invite struct {
	// ID identifies the resource, event, or request.
	ID string `json:"id"`
	// Email contains the invite or user email address.
	Email string `json:"email,omitempty"`
	// Role identifies the organization or message role.
	Role string `json:"role,omitempty"`
	// Projects lists projects associated with an invitation.
	Projects []string `json:"projects,omitempty"`
	// Status reports the current provider processing or lifecycle state.
	Status string `json:"status,omitempty"`
	// CreatedAt contains the Unix creation timestamp.
	CreatedAt int64 `json:"created_at,omitempty"`
}

// CreateInviteRequest creates an organization invitation.
type CreateInviteRequest struct {
	// Email contains the invite or user email address.
	Email string `json:"email"`
	// Role identifies the organization or message role.
	Role string `json:"role"`
	// Projects lists projects associated with an invitation.
	Projects []string `json:"projects,omitempty"`
}

// ListInvitesResponse lists organization invitations.
type ListInvitesResponse struct {
	// Object identifies the provider resource or collection type.
	Object string `json:"object"`
	// Data contains the returned records.
	Data []Invite `json:"data"`
	// HasMore reports whether another page is available.
	HasMore bool `json:"has_more"`
}
