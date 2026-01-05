package routes

import (
	"handsoft/internal/auth"
	"handsoft/internal/http/handlers"
	"handsoft/internal/http/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterWarehouseRoutes(api *gin.RouterGroup, deps Deps) {
	jwtCfg := auth.JWTConfig{
		Secret:    deps.JWTSecret,
		Issuer:    deps.Issuer,
		AccessTTL: deps.AccessTTL,
	}

	h := &handlers.WarehouseModule{DB: deps.DB}

	wh := api.Group("/warehouse")
	wh.Use(middleware.AuthJWT(jwtCfg))
	{
		// Espacios
		wh.POST("/spaces", middleware.RequirePermission(deps.DB, "warehouse:create"), h.CreateSpace)
		wh.GET("/spaces", middleware.RequirePermission(deps.DB, "warehouse:read"), h.ListSpaces)
		wh.GET("/spaces/:id", middleware.RequirePermission(deps.DB, "warehouse:read"), h.GetSpace)

		// Pisos (solo building)
		wh.POST("/spaces/:id/floors", middleware.RequirePermission(deps.DB, "warehouse:create"), h.CreateFloor)

		// Bodegas
		wh.POST("/floors/:floorId/warehouses", middleware.RequirePermission(deps.DB, "warehouse:create"), h.CreateWarehouseInFloor)
		wh.PUT("/warehouses/:id/config", middleware.RequirePermission(deps.DB, "warehouse:update"), h.UpdateWarehouseConfig)

		wh.GET("/warehouses/:id", middleware.RequirePermission(deps.DB, "warehouse:read"), h.GetWarehouse)
	}
}
