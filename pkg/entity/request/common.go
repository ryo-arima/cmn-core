package request

import "time"

type RrCommon struct {
	ID        uint       `json:"id"`
	UUID      string     `json:"uuid"`
	Name      string     `json:"name"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

// RrUser represents a request to create or update a local DB user.
type RrUser struct {
	UUID     string `json:"uuid"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

// RrGroup represents a request to create or update a local DB group.
type RrGroup struct {
	UUID string `json:"uuid"`
	Name string `json:"name"`
}

// RrMember represents a request to create or update a local DB group membership.
type RrMember struct {
	UUID      string `json:"uuid"`
	UserUUID  string `json:"user_uuid"`
	GroupUUID string `json:"group_uuid"`
}

// RrLogin holds credentials for a login request.
type RrLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

