package routes

import (
	"time"

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
	api := r.Group("/api")

	RegisterAuthRoutes(api, deps)
	RegisterUserRoutes(api, deps)
	RegisterInventoryRoutes(api, deps)
	RegisterGeoRoutes(api, deps)
	RegisterAdminRoutes(api, deps)
}
