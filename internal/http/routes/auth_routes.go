package routes

import (
	"handsoft/internal/auth"
	"handsoft/internal/http/handlers"

	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(api *gin.RouterGroup, deps Deps) {
	authHandler := &handlers.AuthHandler{
		DB: deps.DB,
		JWTConfig: auth.JWTConfig{
			Secret:    deps.JWTSecret,
			Issuer:    deps.Issuer,
			AccessTTL: deps.AccessTTL,
		},
	}

	authRoutes := api.Group("/auth")
	{
		authRoutes.POST("/register", authHandler.Register)
		authRoutes.POST("/login", authHandler.Login)
	}
}
