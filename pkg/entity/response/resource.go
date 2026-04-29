package response

import "time"

// RrResource is the representation of a resource returned to the client.
type RrResource struct {
	ID          uint       `json:"id"`
	UUID        string     `json:"uuid"`
	Name        string     `json:"name"`
	Description string     `json:"description,omitempty"`
	OwnerGroup  string     `json:"owner_group,omitempty"`
	CreatedBy   string     `json:"created_by"`
	UpdatedBy   string     `json:"updated_by,omitempty"`
	DeletedBy   string     `json:"deleted_by,omitempty"`
	CreatedAt   *time.Time `json:"created_at,omitempty"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
}

// RrResources wraps a list of resources.
type RrResources struct {
	Code      string       `json:"code"`
	Message   string       `json:"message"`
	Resources []RrResource `json:"resources,omitempty"`
}

// RrSingleResource wraps a single resource.
type RrSingleResource struct {
	Code     string      `json:"code"`
	Message  string      `json:"message"`
	Resource *RrResource `json:"resource,omitempty"`
}

// RrResourceGroupRole is the representation of a group-role entry on a resource.
type RrResourceGroupRole struct {
	ResourceUUID string `json:"resource_uuid"`
	GroupID      string `json:"group_id"`
	Role         string `json:"role"`
}

// RrResourceGroupRoles wraps a list of group-role entries.
type RrResourceGroupRoles struct {
	Code    string                `json:"code"`
	Message string                `json:"message"`
	Groups  []RrResourceGroupRole `json:"groups,omitempty"`
}
