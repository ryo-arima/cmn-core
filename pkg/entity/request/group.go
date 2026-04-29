package request

// RrCreateGroup is the request body for creating a new group in the IdP.
type RrCreateGroup struct {
	Name string `json:"name" binding:"required"`
}

// RrUpdateGroup is the request body for updating an existing group in the IdP.
type RrUpdateGroup struct {
	Name string `json:"name" binding:"required"`
}
