package model

import "time"

type Groups struct {
	ID        uint       `gorm:"primaryKey,autoIncrement"`
	UUID      string
	Name      string
	CreatedAt *time.Time
	UpdatedAt *time.Time
	DeletedAt *time.Time
}
