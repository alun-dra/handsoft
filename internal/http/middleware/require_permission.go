package middleware

import (
	"net/http"
	"strings"

	"handsoft/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RequirePermission valida que el usuario (según roles en el JWT) tenga un permiso.
// - Bypass total si alguno de sus roles tiene IsSuperAdmin=true.
// - Soporta comodín "modulo:*" además del permiso exacto.
func RequirePermission(db *gorm.DB, permissionCode string) gin.HandlerFunc {
	return func(c *gin.Context) {

		rolesAny, ok := c.Get(CtxRolesKey)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "no roles in context"})
			return
		}

		roleNames, ok := rolesAny.([]string)
		if !ok || len(roleNames) == 0 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "roles invalid or empty"})
			return
		}

		// 1) Bypass si es super admin (por flag real en DB)
		isSuper, err := hasSuperAdminRole(db, roleNames)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "error checking super admin"})
			return
		}
		if isSuper {
			c.Next()
			return
		}

		// 2) Permiso exacto
		okPerm, err := roleHasPermission(db, roleNames, permissionCode)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "error checking permission"})
			return
		}
		if okPerm {
			c.Next()
			return
		}

		// 3) Soporte comodín "modulo:*"
		// Ej: si pides "dashboard:export_excel", acepta "dashboard:*"
		if moduleWildcard(permissionCode) != "" {
			okWild, err := roleHasPermission(db, roleNames, moduleWildcard(permissionCode))
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "error checking wildcard permission"})
				return
			}
			if okWild {
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"error":      "forbidden",
			"permission": permissionCode,
		})
	}
}

func moduleWildcard(code string) string {
	parts := strings.SplitN(code, ":", 2)
	if len(parts) != 2 || parts[0] == "" {
		return ""
	}
	return parts[0] + ":*"
}

func hasSuperAdminRole(db *gorm.DB, roleNames []string) (bool, error) {
	var count int64
	err := db.Model(&models.Role{}).
		Where("name IN ?", roleNames).
		Where("is_super_admin = ?", true).
		Count(&count).Error
	return count > 0, err
}

func roleHasPermission(db *gorm.DB, roleNames []string, permissionCode string) (bool, error) {
	var count int64
	err := db.Model(&models.Permission{}).
		Joins("JOIN role_permissions rp ON rp.permission_id = permissions.id").
		Joins("JOIN roles r ON r.id = rp.role_id").
		Where("r.name IN ?", roleNames).
		Where("permissions.code = ?", permissionCode).
		Count(&count).Error
	return count > 0, err
}
