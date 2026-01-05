package handlers

import (
	"handsoft/internal/models"

	"gorm.io/gorm"
)

func isSuperAdminByRoles(db *gorm.DB, roleNames []string) (bool, error) {
	var count int64
	err := db.Model(&models.Role{}).
		Where("name IN ?", roleNames).
		Where("is_super_admin = ?", true).
		Count(&count).Error
	return count > 0, err
}

func permissionsByRoles(db *gorm.DB, roleNames []string) ([]string, error) {
	// Trae permisos por join role_permissions
	var perms []models.Permission
	err := db.Model(&models.Permission{}).
		Select("DISTINCT permissions.id, permissions.code").
		Joins("JOIN role_permissions rp ON rp.permission_id = permissions.id").
		Joins("JOIN roles r ON r.id = rp.role_id").
		Where("r.name IN ?", roleNames).
		Order("permissions.code asc").
		Find(&perms).Error
	if err != nil {
		return nil, err
	}

	out := make([]string, 0, len(perms))
	for _, p := range perms {
		out = append(out, p.Code)
	}
	return out, nil
}
