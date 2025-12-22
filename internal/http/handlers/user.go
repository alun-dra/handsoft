package handlers

import (
	"net/http"

	"handsoft/internal/http/middleware"
	"handsoft/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserHandler struct {
	DB *gorm.DB
}

func (h *UserHandler) Me(c *gin.Context) {
	userIDAny, ok := c.Get(middleware.CtxUserIDKey)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no autenticado"})
		return
	}

	userID, ok := userIDAny.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "userID inválido"})
		return
	}

	var u models.User

	// Cargamos relaciones relevantes
	// - Contacts: datos personales
	// - Phones: teléfonos
	// - Roles: roles del usuario
	// - Address -> Commune -> City -> Region -> Country: ubicación oficial
	q := h.DB.
		Preload("Contacts").
		Preload("Phones").
		Preload("Roles").
		Preload("Address").
		Preload("Address.Commune").
		Preload("Address.Commune.City").
		Preload("Address.Commune.City.Region").
		Preload("Address.Commune.City.Region.Country")

	if err := q.First(&u, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "usuario no encontrado"})
		return
	}

	// Roles
	roles := make([]string, 0, len(u.Roles))
	for _, r := range u.Roles {
		roles = append(roles, r.Name)
	}

	// Phones
	phones := make([]gin.H, 0, len(u.Phones))
	for _, p := range u.Phones {
		phones = append(phones, gin.H{
			"label":   p.Label,
			"number":  p.Number,
			"is_main": p.IsMain,
		})
	}

	// Address + Location (desde Address)
	var address any = nil
	location := gin.H{
		"country": nil,
		"region":  nil,
		"city":    nil,
		"commune": nil,
	}

	if u.Address != nil {
		// Address completo
		address = gin.H{
			"id":                      u.Address.ID,
			"street":                  u.Address.Street,
			"street_number":           u.Address.StreetNumber,
			"is_condominium":          u.Address.IsCondominium,
			"condominium_house_number": u.Address.CondominiumHouseNumber,
			"building_number":         u.Address.BuildingNumber,
			"apartment_number":        u.Address.ApartmentNumber,
			"extra":                   u.Address.Extra,
			"commune_id":              u.Address.CommuneID,
		}

		// Location inferida desde la comuna de Address
		co := u.Address.Commune
		location["commune"] = gin.H{
			"id":   co.ID,
			"name": co.Name,
		}

		ct := co.City
		if ct.ID != 0 {
			location["city"] = gin.H{
				"id":   ct.ID,
				"name": ct.Name,
			}
		}

		rg := ct.Region
		if rg.ID != 0 {
			location["region"] = gin.H{
				"id":   rg.ID,
				"name": rg.Name,
				"code": rg.Code,
			}
		}

		cn := rg.Country
		if cn.ID != 0 {
			location["country"] = gin.H{
				"id":   cn.ID,
				"name": cn.Name,
				"code": cn.Code,
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"id":       u.ID,
		"email":    u.Email,
		"username": u.Username,
		"isActive": u.IsActive,
		"contact": gin.H{
			"full_name": u.Contacts.FullName,
		},
		"phones":   phones,
		"address":  address,
		"location": location,
		"roles":    roles,
	})
}
