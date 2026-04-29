package model

import "time"

// PgResource is the GORM/PostgreSQL model for the resources table.
type PgResource struct {
	ID          uint       `gorm:"primaryKey;autoIncrement"`
	UUID        string     `gorm:"uniqueIndex;not null"`
	Name        string     `gorm:"not null"`
	Description string
	OwnerGroup  string     // IDP group ID of the owning group
	CreatedBy   string     // IDP user ID of the creator
	UpdatedBy   string     // IDP user ID of the last updater
	DeletedBy   string     // IDP user ID of the deleter
	CreatedAt   *time.Time
	UpdatedAt   *time.Time
	DeletedAt   *time.Time `gorm:"index"`
}

// PgResourceGroupRole records which role a group has on a resource.
// role: "viewer" | "editor" | "owner"
type PgResourceGroupRole struct {
	ID           uint       `gorm:"primaryKey;autoIncrement"`
	ResourceUUID string     `gorm:"not null;index"`
	GroupID      string     `gorm:"not null;index"` // IDP group ID
	Role         string     `gorm:"not null"`      // viewer | editor | owner
	CreatedAt    *time.Time
	UpdatedAt    *time.Time
}
