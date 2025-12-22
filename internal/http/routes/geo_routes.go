package routes

import (
	"handsoft/internal/http/handlers"

	"github.com/gin-gonic/gin"
)

func RegisterGeoRoutes(api *gin.RouterGroup, deps Deps) {
	h := &handlers.GeoHandler{DB: deps.DB}

	geo := api.Group("/geo")
	{
		geo.GET("/regions", h.Regions)
		geo.GET("/regions/:regionId/cities", h.CitiesByRegion)
		geo.GET("/cities/:cityId/communes", h.CommunesByCity)
		geo.GET("/communes/:communeId", h.CommuneDetail)
		geo.GET("/communes", h.SearchCommunes) // ?search=...
	}
}
