package models

import "time"

type UserPhone struct {
	ID        uint      `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time

	UserID uint   `gorm:"index;not null"`
	Label  string `gorm:"not null"` // mobile, work, home
	Number string `gorm:"not null"`
	IsMain bool   `gorm:"default:false"`
}
