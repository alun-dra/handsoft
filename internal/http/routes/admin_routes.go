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

	admin := api.Group("/admin")
	admin.Use(
		middleware.AuthJWT(jwtCfg),
		middleware.RequireSuperAdmin(deps.DB),
	)

	// Roles CRUD
	admin.GET("/roles", handlers.ListRoles(deps))
	admin.POST("/roles", handlers.CreateRole(deps))
	admin.PUT("/roles/:id", handlers.UpdateRole(deps))
	admin.DELETE("/roles/:id", handlers.DeleteRole(deps))

	// Permissions
	admin.GET("/permissions", handlers.ListPermissions(deps))

	// Asignar permisos a un rol (replace)
	admin.PUT("/roles/:id/permissions", handlers.SetRolePermissions(deps))
	admin.GET("/roles/:id/permissions", handlers.GetRolePermissions(deps))
}
