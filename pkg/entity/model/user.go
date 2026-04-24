package model

import "time"

type Users struct {
	ID        uint       `gorm:"primaryKey,autoIncrement"`
	UUID      string
	Email     string
	Password  string
	Name      string
	CreatedAt *time.Time
	UpdatedAt *time.Time
	DeletedAt *time.Time
}
