package middleware

import (
	"net/http"

	"handsoft/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RequireSuperAdmin permite acceso solo si el usuario tiene un rol con IsSuperAdmin=true.
// Usa los nombres de roles que ya vienen en el JWT (CtxRolesKey).
func RequireSuperAdmin(db *gorm.DB) gin.HandlerFunc {
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

		var count int64
		if err := db.Model(&models.Role{}).
			Where("name IN ?", roleNames).
			Where("is_super_admin = ?", true).
			Count(&count).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "error checking super admin"})
			return
		}

		if count == 0 {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}

		c.Next()
	}
}
