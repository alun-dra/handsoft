package handlers

import (
	"net/http"
	"strconv"

	"handsoft/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type WarehouseModule struct {
	DB *gorm.DB
}

type createSpaceReq struct {
	Name        string `json:"name" binding:"required"`
	Type        string `json:"type" binding:"required"` // open_area | building
	Description string `json:"description"`

	// Si open_area, creamos la bodega principal de inmediato:
	OpenAreaWarehouse struct {
		Name         string  `json:"name" binding:"required"`
		AreaM2       float64 `json:"area_m2"`
		PalletsFloor int     `json:"pallets_floor"`
		HasRacks     bool    `json:"has_racks"`
		Racks        []struct {
			Label           string  `json:"label" binding:"required"`
			Levels          int     `json:"levels"`
			PalletsPerLevel int     `json:"pallets_per_level"`
			LengthM         float64 `json:"length_m"`
		} `json:"racks"`
	} `json:"open_area_warehouse"`
}

func (h *WarehouseModule) CreateSpace(c *gin.Context) {
	var req createSpaceReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_body"})
		return
	}

	st := models.SpaceType(req.Type)
	if st != models.SpaceTypeOpenArea && st != models.SpaceTypeBuilding {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_space_type"})
		return
	}

	err := h.DB.Transaction(func(tx *gorm.DB) error {
		space := models.Space{
			Name:        req.Name,
			Type:        st,
			Description: req.Description,
		}
		if err := tx.Create(&space).Error; err != nil {
			return err
		}

		// open_area => crear 1 bodega principal
		if st == models.SpaceTypeOpenArea {
			wreq := req.OpenAreaWarehouse
			wh := models.Warehouse{
				SpaceID:      space.ID,
				Name:         wreq.Name,
				AreaM2:       wreq.AreaM2,
				PalletsFloor: wreq.PalletsFloor,
				HasRacks:     wreq.HasRacks,
			}
			if err := tx.Create(&wh).Error; err != nil {
				return err
			}

			if wreq.HasRacks {
				for _, r := range wreq.Racks {
					rack := models.WarehouseRack{
						WarehouseID:    wh.ID,
						Label:          r.Label,
						Levels:         r.Levels,
						PalletsPerLevel:r.PalletsPerLevel,
						LengthM:        r.LengthM,
					}
					if err := tx.Create(&rack).Error; err != nil {
						return err
					}
				}
			}
		}

		c.JSON(http.StatusCreated, gin.H{"id": space.ID})
		return nil
	})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cannot_create_space"})
		return
	}
}

func (h *WarehouseModule) ListSpaces(c *gin.Context) {
	var spaces []models.Space
	if err := h.DB.Order("id asc").Find(&spaces).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db_error"})
		return
	}
	c.JSON(http.StatusOK, spaces)
}

func (h *WarehouseModule) GetSpace(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var space models.Space
	if err := h.DB.
		Preload("Floors.Warehouses.Racks").
		Preload("Warehouses.Racks").
		First(&space, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "space_not_found"})
		return
	}
	c.JSON(http.StatusOK, space)
}

type createFloorReq struct {
	Number int `json:"number" binding:"required"`
}

func (h *WarehouseModule) CreateFloor(c *gin.Context) {
	spaceID, _ := strconv.Atoi(c.Param("id"))

	var space models.Space
	if err := h.DB.First(&space, spaceID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "space_not_found"})
		return
	}
	if space.Type != models.SpaceTypeBuilding {
		c.JSON(http.StatusBadRequest, gin.H{"error": "space_is_not_building"})
		return
	}

	var req createFloorReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_body"})
		return
	}

	floor := models.SpaceFloor{SpaceID: uint(spaceID), Number: req.Number}
	if err := h.DB.Create(&floor).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cannot_create_floor"})
		return
	}

	c.JSON(http.StatusCreated, floor)
}

type createWarehouseReq struct {
	Name         string  `json:"name" binding:"required"`
	AreaM2       float64 `json:"area_m2"`
	PalletsFloor int     `json:"pallets_floor"`
	HasRacks     bool    `json:"has_racks"`
}

func (h *WarehouseModule) CreateWarehouseInFloor(c *gin.Context) {
	floorID, _ := strconv.Atoi(c.Param("floorId"))

	var floor models.SpaceFloor
	if err := h.DB.First(&floor, floorID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "floor_not_found"})
		return
	}

	var req createWarehouseReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_body"})
		return
	}

	w := models.Warehouse{
		SpaceID:      floor.SpaceID,
		FloorID:      &floor.ID,
		Name:         req.Name,
		AreaM2:       req.AreaM2,
		PalletsFloor: req.PalletsFloor,
		HasRacks:     req.HasRacks,
	}
	if err := h.DB.Create(&w).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cannot_create_warehouse"})
		return
	}

	c.JSON(http.StatusCreated, w)
}

type updateWarehouseConfigReq struct {
	AreaM2       *float64 `json:"area_m2"`
	PalletsFloor *int     `json:"pallets_floor"`
	HasRacks     *bool    `json:"has_racks"`

	// si HasRacks=true, mandas racks completos (replace)
	Racks []struct {
		Label           string  `json:"label" binding:"required"`
		Levels          int     `json:"levels"`
		PalletsPerLevel int     `json:"pallets_per_level"`
		LengthM         float64 `json:"length_m"`
	} `json:"racks"`
}

func (h *WarehouseModule) UpdateWarehouseConfig(c *gin.Context) {
	warehouseID, _ := strconv.Atoi(c.Param("id"))

	var w models.Warehouse
	if err := h.DB.Preload("Racks").First(&w, warehouseID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "warehouse_not_found"})
		return
	}

	var req updateWarehouseConfigReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_body"})
		return
	}

	err := h.DB.Transaction(func(tx *gorm.DB) error {
		if req.AreaM2 != nil {
			w.AreaM2 = *req.AreaM2
		}
		if req.PalletsFloor != nil {
			w.PalletsFloor = *req.PalletsFloor
		}
		if req.HasRacks != nil {
			w.HasRacks = *req.HasRacks
		}

		if err := tx.Save(&w).Error; err != nil {
			return err
		}

		// Si HasRacks = true, hacemos replace de racks (porque la nomenclatura es libre)
		if w.HasRacks {
			// borrar racks existentes
			if err := tx.Where("warehouse_id = ?", w.ID).Delete(&models.WarehouseRack{}).Error; err != nil {
				return err
			}
			for _, r := range req.Racks {
				rack := models.WarehouseRack{
					WarehouseID:     w.ID,
					Label:           r.Label,
					Levels:          r.Levels,
					PalletsPerLevel: r.PalletsPerLevel,
					LengthM:         r.LengthM,
				}
				if err := tx.Create(&rack).Error; err != nil {
					return err
				}
			}
		}

		return nil
	})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cannot_update_config"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ok": true})
}

func (h *WarehouseModule) GetWarehouse(c *gin.Context) {
	warehouseID, _ := strconv.Atoi(c.Param("id"))

	var w models.Warehouse
	if err := h.DB.Preload("Racks").First(&w, warehouseID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "warehouse_not_found"})
		return
	}

	// capacidad calculada
	totalRack := 0
	for _, r := range w.Racks {
		totalRack += r.Levels * r.PalletsPerLevel
	}
	total := w.PalletsFloor + totalRack

	c.JSON(http.StatusOK, gin.H{
		"warehouse": w,
		"capacity": gin.H{
			"pallets_floor": w.PalletsFloor,
			"pallets_racks": totalRack,
			"pallets_total": total,
		},
	})
}
