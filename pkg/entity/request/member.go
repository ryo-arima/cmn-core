package request

// RrAddGroupMember is the request body for adding a user to a group.
type RrAddGroupMember struct {
	UserID string `json:"user_id" binding:"required"`
	Role   string `json:"role"    binding:"required,oneof=owner editor viewer"`
}

// RrRemoveGroupMember is the request body for removing a user from a group.
type RrRemoveGroupMember struct {
	UserID string `json:"user_id" binding:"required"`
}
