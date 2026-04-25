package response

import "time"

// Resource is the representation of a resource returned to the client.
type Resource struct {
	ID          uint       `json:"id"`
	UUID        string     `json:"uuid"`
	Name        string     `json:"name"`
	Description string     `json:"description,omitempty"`
	CreatedBy   string     `json:"created_by"`
	UpdatedBy   string     `json:"updated_by,omitempty"`
	DeletedBy   string     `json:"deleted_by,omitempty"`
	CreatedAt   *time.Time `json:"created_at,omitempty"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
}

// Resources wraps a list of resources.
type Resources struct {
	Code      string     `json:"code"`
	Message   string     `json:"message"`
	Resources []Resource `json:"resources,omitempty"`
}

// SingleResource wraps a single resource.
type SingleResource struct {
	Code     string    `json:"code"`
	Message  string    `json:"message"`
	Resource *Resource `json:"resource,omitempty"`
}

// ResourceGroupRole is the representation of a group-role entry on a resource.
type ResourceGroupRole struct {
	ResourceUUID string `json:"resource_uuid"`
	GroupUUID    string `json:"group_uuid"`
	Role         string `json:"role"`
}

// ResourceGroupRoles wraps a list of group-role entries.
type ResourceGroupRoles struct {
	Code    string              `json:"code"`
	Message string              `json:"message"`
	Groups  []ResourceGroupRole `json:"groups,omitempty"`
}
