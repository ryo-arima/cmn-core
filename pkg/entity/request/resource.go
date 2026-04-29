package request

// RrCreateResource is the request body for creating a new resource.
type RrCreateResource struct {
	Name        string `json:"name"        binding:"required"`
	Description string `json:"description"`
}

// RrUpdateResource is the request body for updating an existing resource.
type RrUpdateResource struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// RrSetResourceGroupRole adds or updates the role a group has on a resource.
type RrSetResourceGroupRole struct {
	GroupUUID string `json:"group_uuid" binding:"required"`
	Role      string `json:"role"       binding:"required,oneof=viewer editor owner"`
}

// RrDeleteResourceGroupRole specifies the group to remove from a resource's group-role list.
type RrDeleteResourceGroupRole struct {
	GroupUUID string `json:"group_uuid" binding:"required"`
}

// LoResourceQueryFilter holds optional filter conditions for resource queries (internal use only).
type LoResourceQueryFilter struct {
	CreatedBy  string
	GroupUUIDs []string // user's group UUIDs for membership-based access
}
