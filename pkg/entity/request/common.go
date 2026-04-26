package request

import "time"

type Common struct {
	ID        uint       `json:"id"`
	UUID      string     `json:"uuid"`
	Name      string     `json:"name"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

// User represents a request to create or update a local DB user.
type User struct {
	UUID     string `json:"uuid"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

// Group represents a request to create or update a local DB group.
type Group struct {
	UUID string `json:"uuid"`
	Name string `json:"name"`
}

// Member represents a request to create or update a local DB group membership.
type Member struct {
	UUID      string `json:"uuid"`
	UserUUID  string `json:"user_uuid"`
	GroupUUID string `json:"group_uuid"`
}

// Login holds credentials for a login request.
type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

