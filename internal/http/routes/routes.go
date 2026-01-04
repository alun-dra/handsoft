package routes

import (
	"os"
	"time"

	"handsoft/internal/http/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Deps = dependencias compartidas por TODAS las rutas
type Deps struct {
	DB        *gorm.DB
	JWTSecret string
	Issuer    string
	AccessTTL time.Duration
}

// Register es el punto único de entrada de rutas
func Register(r *gin.Engine, deps Deps) {

	// ============================
	// API KEY (obligatoria)
	// ============================
	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		panic("API_KEY no está definido en variables de entorno")
	}

	api := r.Group("/api")

	// 🔐 Middleware API-Key global para /api
	api.Use(middleware.RequireAPIKey(middleware.APIKeyConfig{
		ValidKeys: map[string]bool{
			apiKey: true,
		},
		HeaderName: "X-API-Key",
	}))

	// ============================
	// RUTAS
	// ============================

	// Auth (login, refresh, etc.)
	RegisterAuthRoutes(api, deps)

	// Usuarios (JWT requerido internamente)
	RegisterUserRoutes(api, deps)

	// Inventario
	RegisterInventoryRoutes(api, deps)

	// Geografía (regiones, comunas, etc.)
	RegisterGeoRoutes(api, deps)

	// Admin (roles, permisos, etc.)
	RegisterAdminRoutes(api, deps)
}
