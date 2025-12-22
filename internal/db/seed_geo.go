package db

import (
	"errors"

	"handsoft/internal/models"
	"gorm.io/gorm"
)

func SeedChile(gdb *gorm.DB) error {
	var cl models.Country
	err := gdb.Where("code = ?", "CL").First(&cl).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			cl = models.Country{Name: "Chile", Code: "CL"}
			return gdb.Create(&cl).Error
		}
		return err
	}
	return nil
}
