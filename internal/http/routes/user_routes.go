package routes

import (
	"handsoft/internal/auth"
	"handsoft/internal/http/handlers"
	"handsoft/internal/http/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(api *gin.RouterGroup, deps Deps) {
	jwtCfg := auth.JWTConfig{
		Secret:    deps.JWTSecret,
		Issuer:    deps.Issuer,
		AccessTTL: deps.AccessTTL,
	}

	userH := &handlers.UserHandler{DB: deps.DB}

	users := api.Group("/users")
	users.Use(middleware.AuthJWT(jwtCfg))
	{
		users.GET("/me", userH.Me)
	}
}
