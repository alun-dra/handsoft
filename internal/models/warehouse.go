package models

import "time"

type SpaceType string

const (
	SpaceTypeOpenArea SpaceType = "open_area"
	SpaceTypeBuilding SpaceType = "building"
)

type Space struct {
	ID        uint      `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time

	Name        string    `gorm:"not null"`
	Type        SpaceType `gorm:"type:varchar(20);not null"` // open_area | building
	Description string

	// Si es building: tendrá Floors
	Floors []SpaceFloor `gorm:"constraint:OnDelete:CASCADE;"`
	// Si es open_area: tendrá una sola Warehouse (la "principal")
	Warehouses []Warehouse `gorm:"constraint:OnDelete:CASCADE;"`
}

type SpaceFloor struct {
	ID        uint      `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time

	SpaceID uint `gorm:"index;not null"`
	Number  int  `gorm:"not null"` // 1,2,3...

	Warehouses []Warehouse `gorm:"constraint:OnDelete:CASCADE;"`
}

type Warehouse struct {
	ID        uint      `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time

	// Pertenece a Space. Si es building, además tendrá FloorID.
	SpaceID uint `gorm:"index;not null"`
	FloorID *uint `gorm:"index"` // null para open_area

	Name string `gorm:"not null"`

	// m² configurable por bodega
	AreaM2 float64 `gorm:"not null;default:0"`

	// Pallets a ras de piso
	PalletsFloor int `gorm:"not null;default:0"`

	// ¿Tiene racks?
	HasRacks bool `gorm:"not null;default:false"`

	// Racks (si HasRacks=true)
	Racks []WarehouseRack `gorm:"constraint:OnDelete:CASCADE;"`
}

type WarehouseRack struct {
	ID        uint      `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time

	WarehouseID uint `gorm:"index;not null"`

	// Nomenclatura libre: "Rack 1A", "RX-01", etc.
	Label string `gorm:"not null"`

	Levels          int     `gorm:"not null;default:0"` // pisos del rack
	PalletsPerLevel int     `gorm:"not null;default:0"` // pallets por nivel
	LengthM         float64 `gorm:"not null;default:0"` // opcional, metros del rack
}
