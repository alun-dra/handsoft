package models

import "time"

type User struct {
	ID        uint      `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time

	Email        string `gorm:"uniqueIndex;not null"`
	Username     string `gorm:"uniqueIndex;not null"`
	PasswordHash string `gorm:"not null"`
	IsActive     bool   `gorm:"default:true"`

	// Ubicación administrativa (opcional, puede convivir con Address)
	CommuneID *uint
	Commune   *Commune `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`

	// Dirección principal (muchos usuarios pueden compartir la misma Address)
	AddressID *uint
	Address   *Address `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`

	// Contacto y teléfonos
	Contacts Contact     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Phones   []UserPhone `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	// Roles
	Roles []Role `gorm:"many2many:user_roles;"`
}
