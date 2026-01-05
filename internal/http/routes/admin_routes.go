package routes

import (
	"handsoft/internal/auth"
	"handsoft/internal/http/handlers"
	"handsoft/internal/http/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterAdminRoutes(api *gin.RouterGroup, deps Deps) {
	jwtCfg := auth.JWTConfig{
		Secret:    deps.JWTSecret,
		Issuer:    deps.Issuer,
		AccessTTL: deps.AccessTTL,
	}

	adminH := &handlers.AdminHandler{DB: deps.DB}

	admin := api.Group("/admin")
	admin.Use(
		middleware.AuthJWT(jwtCfg),
		middleware.RequireSuperAdmin(deps.DB),
	)

	// Roles CRUD
	admin.GET("/roles", adminH.ListRoles)
	admin.POST("/roles", adminH.CreateRole)
	admin.PUT("/roles/:id", adminH.UpdateRole)
	admin.DELETE("/roles/:id", adminH.DeleteRole)

	// Permissions
	admin.GET("/permissions", adminH.ListPermissions)

	// Asignar permisos a un rol (replace)
	admin.PUT("/roles/:id/permissions", adminH.SetRolePermissions)
	admin.GET("/roles/:id/permissions", adminH.GetRolePermissions)
}
