package routes

import (
	"handsoft/internal/auth"
	"handsoft/internal/http/handlers"
	"handsoft/internal/http/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterAdminRoutes(api *gin.RouterGroup, deps Deps) {
	h := &handlers.AdminHandler{DB: deps.DB}

	admin := api.Group("/admin")
	admin.Use(middleware.AuthJWT(auth.JWTConfig{
		Secret:    deps.JWTSecret,
		Issuer:    deps.Issuer,
		AccessTTL: deps.AccessTTL,
	}))
	admin.Use(middleware.RequireRole("super_admin"))
	{
		admin.PUT("/users/:id/roles", h.UpdateUserRoles)
	}
}
