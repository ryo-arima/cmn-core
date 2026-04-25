package model

import "time"

// Resource is the core domain entity managed by cmn-core.
type Resource struct {
	ID          uint       `gorm:"primaryKey;autoIncrement"`
	UUID        string     `gorm:"uniqueIndex;not null"`
	Name        string     `gorm:"not null"`
	Description string
	CreatedBy   string     // user UUID of the creator
	UpdatedBy   string     // user UUID of the last updater
	DeletedBy   string     // user UUID of the deleter
	CreatedAt   *time.Time
	UpdatedAt   *time.Time
	DeletedAt   *time.Time `gorm:"index"`
}

// ResourceGroupRole records which role a group has on a resource.
// role: "viewer" | "editor" | "owner"
type ResourceGroupRole struct {
	ID           uint       `gorm:"primaryKey;autoIncrement"`
	ResourceUUID string     `gorm:"not null;index"`
	GroupUUID    string     `gorm:"not null;index"`
	Role         string     `gorm:"not null"` // viewer | editor | owner
	CreatedAt    *time.Time
	UpdatedAt    *time.Time
}
