package request

// CreateResource is the request body for creating a new resource.
type CreateResource struct {
	Name        string `json:"name"        binding:"required"`
	Description string `json:"description"`
}

// UpdateResource is the request body for updating an existing resource.
type UpdateResource struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// SetResourceGroupRole adds or updates the role a group has on a resource.
type SetResourceGroupRole struct {
	GroupUUID string `json:"group_uuid" binding:"required"`
	Role      string `json:"role"       binding:"required,oneof=viewer editor owner"`
}

// DeleteResourceGroupRole specifies the group to remove from a resource's group-role list.
type DeleteResourceGroupRole struct {
	GroupUUID string `json:"group_uuid" binding:"required"`
}
