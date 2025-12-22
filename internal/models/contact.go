package models

import "time"

type Contact struct {
	ID        uint      `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time

	UserID uint `gorm:"uniqueIndex;not null"`

	FullName string
	// EmailAlt, AddressLine, etc. (si quieres)
}
