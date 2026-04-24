package model

import "time"

type Members struct {
	ID        uint       `gorm:"primaryKey,autoIncrement"`
	UUID      string
	GroupUUID string
	UserUUID  string
	Role      string
	CreatedAt *time.Time
	UpdatedAt *time.Time
	DeletedAt *time.Time
}
