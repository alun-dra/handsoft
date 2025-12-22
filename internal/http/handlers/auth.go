package handlers

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"handsoft/internal/auth"
	"handsoft/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuthHandler struct {
	DB        *gorm.DB
	JWTConfig auth.JWTConfig
}

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Username string `json:"username" binding:"required,min=3,max=32"`
	Password string `json:"password" binding:"required,min=8,max=72"`

	FullName string `json:"full_name"`
	Phone    string `json:"phone"`

	// Dirección (obligatoria)
	CommuneID    uint   `json:"commune_id" binding:"required"`
	Street       string `json:"street" binding:"required"`
	StreetNumber string `json:"street_number" binding:"required"`

	IsCondominium           bool   `json:"is_condominium"`
	CondominiumHouseNumber  string `json:"condominium_house_number"`
	BuildingNumber          string `json:"building_number"`
	ApartmentNumber         string `json:"apartment_number"`
	Extra                   string `json:"extra"`
}

type LoginRequest struct {
	Login    string `json:"login" binding:"required"` // email o username
	Password string `json:"password" binding:"required"`
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Normalización
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))
	req.Username = strings.TrimSpace(req.Username)
	req.Phone = strings.TrimSpace(req.Phone)
	req.FullName = strings.TrimSpace(req.FullName)

	req.Street = strings.TrimSpace(req.Street)
	req.StreetNumber = strings.TrimSpace(req.StreetNumber)
	req.CondominiumHouseNumber = strings.TrimSpace(req.CondominiumHouseNumber)
	req.BuildingNumber = strings.TrimSpace(req.BuildingNumber)
	req.ApartmentNumber = strings.TrimSpace(req.ApartmentNumber)
	req.Extra = strings.TrimSpace(req.Extra)

	// Chequear duplicados
	var count int64
	h.DB.Model(&models.User{}).
		Where("email = ? OR username = ?", req.Email, req.Username).
		Count(&count)
	if count > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "email o username ya existe"})
		return
	}

	// Validar que la comuna exista (y por ende ciudad/región/país se infiere)
	var commune models.Commune
	if err := h.DB.First(&commune, req.CommuneID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "commune_id no existe"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no se pudo validar la comuna"})
		return
	}

	hash, err := auth.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no se pudo procesar la contraseña"})
		return
	}

	// Rol por defecto (user)
	var role models.Role
	if err := h.DB.Where("name = ?", "user").First(&role).Error; err != nil {
		role = models.Role{Name: "user", Description: "Rol base"}
		if err := h.DB.Create(&role).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "no se pudo crear rol base"})
			return
		}
	}

	// Reutilizar dirección si ya existe EXACTAMENTE igual
	var addr models.Address
	addrQuery := h.DB.Where(`
		commune_id = ? AND
		street = ? AND
		street_number = ? AND
		is_condominium = ? AND
		COALESCE(condominium_house_number,'') = ? AND
		COALESCE(building_number,'') = ? AND
		COALESCE(apartment_number,'') = ? AND
		COALESCE(extra,'') = ?
	`,
		req.CommuneID,
		req.Street,
		req.StreetNumber,
		req.IsCondominium,
		req.CondominiumHouseNumber,
		req.BuildingNumber,
		req.ApartmentNumber,
		req.Extra,
	)

	if err := addrQuery.First(&addr).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// No existe -> crear
			addr = models.Address{
				CommuneID:               req.CommuneID,
				Street:                  req.Street,
				StreetNumber:            req.StreetNumber,
				IsCondominium:           req.IsCondominium,
				CondominiumHouseNumber:  req.CondominiumHouseNumber,
				BuildingNumber:          req.BuildingNumber,
				ApartmentNumber:         req.ApartmentNumber,
				Extra:                   req.Extra,
			}
			if err := h.DB.Create(&addr).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "no se pudo crear dirección"})
				return
			}
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "no se pudo consultar dirección"})
			return
		}
	}

	// Crear usuario con dirección
	u := models.User{
		Email:        req.Email,
		Username:     req.Username,
		PasswordHash: hash,
		IsActive:     true,

		AddressID: &addr.ID,

		Contacts: models.Contact{
			FullName: req.FullName,
		},

		Roles: []models.Role{role},
	}

	// Si viene teléfono, lo guardamos como phone principal
	if req.Phone != "" {
		u.Phones = []models.UserPhone{
			{
				Label:  "mobile",
				Number: req.Phone,
				IsMain: true,
			},
		}
	}

	if err := h.DB.Create(&u).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no se pudo crear usuario"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":       u.ID,
		"email":    u.Email,
		"username": u.Username,
		"created":  u.CreatedAt.Format(time.RFC3339),
		"address_id": addr.ID,
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	login := strings.ToLower(strings.TrimSpace(req.Login))

	var u models.User
	err := h.DB.Preload("Roles").
		Where("email = ? OR username = ?", login, req.Login).
		First(&u).Error
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "credenciales inválidas"})
		return
	}

	if !u.IsActive {
		c.JSON(http.StatusForbidden, gin.H{"error": "usuario desactivado"})
		return
	}

	if !auth.CheckPassword(u.PasswordHash, req.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "credenciales inválidas"})
		return
	}

	roles := make([]string, 0, len(u.Roles))
	for _, r := range u.Roles {
		roles = append(roles, r.Name)
	}

	token, err := auth.SignAccessToken(h.JWTConfig, u.ID, roles)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no se pudo generar token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token": token,
		"token_type":   "Bearer",
		"expires_in":   int(h.JWTConfig.AccessTTL.Seconds()),
	})
}
