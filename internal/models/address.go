package models

import "time"

type Address struct {
	ID        uint      `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time

	// Clave: comuna define ciudad/región/país por relaciones
	CommuneID uint    `gorm:"index;not null"`
	Commune   Commune `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`

	Street       string `gorm:"not null"`
	StreetNumber string `gorm:"not null"` // "123", "S/N", "12-A"

	// Condominio (casas)
	IsCondominium          bool
	CondominiumHouseNumber string

	// Departamentos
	BuildingNumber  string
	ApartmentNumber string

	Extra string

	Users []User
}
