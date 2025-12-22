package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"handsoft/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type GeoHandler struct {
	DB *gorm.DB
}

func (h *GeoHandler) Regions(c *gin.Context) {
	var regions []models.Region
	if err := h.DB.Preload("Country").
		Joins("Country").
		Where("countries.code = ?", "CL").
		Order("regions.name ASC").
		Find(&regions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no se pudieron cargar regiones"})
		return
	}

	out := make([]gin.H, 0, len(regions))
	for _, r := range regions {
		out = append(out, gin.H{"id": r.ID, "name": r.Name, "code": r.Code})
	}
	c.JSON(http.StatusOK, out)
}

func (h *GeoHandler) CitiesByRegion(c *gin.Context) {
	regionID, err := strconv.Atoi(c.Param("regionId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "regionId inválido"})
		return
	}

	var cities []models.City
	if err := h.DB.Where("region_id = ?", regionID).
		Order("name ASC").
		Find(&cities).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no se pudieron cargar ciudades"})
		return
	}

	out := make([]gin.H, 0, len(cities))
	for _, ct := range cities {
		out = append(out, gin.H{"id": ct.ID, "name": ct.Name})
	}
	c.JSON(http.StatusOK, out)
}

func (h *GeoHandler) CommunesByCity(c *gin.Context) {
	cityID, err := strconv.Atoi(c.Param("cityId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cityId inválido"})
		return
	}

	var communes []models.Commune
	if err := h.DB.Where("city_id = ?", cityID).
		Order("name ASC").
		Find(&communes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no se pudieron cargar comunas"})
		return
	}

	out := make([]gin.H, 0, len(communes))
	for _, co := range communes {
		out = append(out, gin.H{"id": co.ID, "name": co.Name})
	}
	c.JSON(http.StatusOK, out)
}

func (h *GeoHandler) CommuneDetail(c *gin.Context) {
	communeID, err := strconv.Atoi(c.Param("communeId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "communeId inválido"})
		return
	}

	var co models.Commune
	if err := h.DB.
		Preload("City").
		Preload("City.Region").
		Preload("City.Region.Country").
		First(&co, communeID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "comuna no encontrada"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":   co.ID,
		"name": co.Name,
		"city": gin.H{"id": co.City.ID, "name": co.City.Name},
		"region": gin.H{"id": co.City.Region.ID, "name": co.City.Region.Name, "code": co.City.Region.Code},
		"country": gin.H{"id": co.City.Region.Country.ID, "name": co.City.Region.Country.Name, "code": co.City.Region.Country.Code},
	})
}

func (h *GeoHandler) SearchCommunes(c *gin.Context) {
	q := strings.TrimSpace(c.Query("search"))
	if q == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "search requerido"})
		return
	}

	var communes []models.Commune
	if err := h.DB.
		Where("LOWER(name) LIKE LOWER(?)", "%"+q+"%").
		Order("name ASC").
		Limit(30).
		Find(&communes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no se pudo buscar comunas"})
		return
	}

	out := make([]gin.H, 0, len(communes))
	for _, co := range communes {
		out = append(out, gin.H{"id": co.ID, "name": co.Name})
	}
	c.JSON(http.StatusOK, out)
}
