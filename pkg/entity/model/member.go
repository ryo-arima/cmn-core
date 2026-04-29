package model

import "time"

// PgMembers is the GORM/PostgreSQL model for the members table.
type PgMembers struct {
	ID        uint       `gorm:"primaryKey,autoIncrement"`
	UUID      string
	GroupUUID string
	UserUUID  string
	Role      string
	CreatedAt *time.Time
	UpdatedAt *time.Time
	DeletedAt *time.Time
}
