package models

import "time"

type Country struct {
	ID        uint      `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time

	Name    string   `gorm:"uniqueIndex;not null"`
	Code    string   `gorm:"uniqueIndex;not null"` 
	Regions []Region
}

type Region struct {
	ID        uint      `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time

	CountryID uint    `gorm:"index;not null"`
	Country   Country `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`

	Name   string `gorm:"not null"`
	Code   string `gorm:"index"`
	Cities []City
}

type City struct {
	ID        uint      `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time

	RegionID uint   `gorm:"index;not null"`
	Region   Region `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`

	Name     string    `gorm:"not null"`
	Communes []Commune
}

type Commune struct {
	ID        uint      `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time

	CityID uint `gorm:"index;not null"`
	City   City `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`

	Name  string `gorm:"not null"`
	Users []User
}
