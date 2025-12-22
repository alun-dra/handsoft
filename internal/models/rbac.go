package models

import "time"

type Role struct {
	ID        uint      `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time

	Name        string `gorm:"uniqueIndex;not null"`
	Description string

	IsSuperAdmin bool `gorm:"default:false"`

	Permissions []Permission `gorm:"many2many:role_permissions;"`
}

type Permission struct {
	ID        uint      `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time

	Code        string `gorm:"uniqueIndex;not null"`
	Description string
}
