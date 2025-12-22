package handlers

import (
	"net/http"
	"strings"

	"handsoft/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AdminHandler struct {
	DB *gorm.DB
}

type UpdateUserRolesRequest struct {
	Roles []string `json:"roles" binding:"required,min=1"`
}

// PUT /api/admin/users/:id/roles
func (h *AdminHandler) UpdateUserRoles(c *gin.Context) {
	var req UpdateUserRolesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// usuario objetivo
	var user models.User
	if err := h.DB.Preload("Roles").First(&user, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "usuario no encontrado"})
		return
	}

	// normalizar + dedup
	seen := map[string]bool{}
	want := make([]string, 0, len(req.Roles))
	for _, r := range req.Roles {
		r = strings.TrimSpace(strings.ToLower(r))
		if r == "" || seen[r] {
			continue
		}
		seen[r] = true
		want = append(want, r)
	}
	if len(want) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "roles inválidos"})
		return
	}

	// cargar roles desde DB
	var roles []models.Role
	if err := h.DB.Where("name IN ?", want).Find(&roles).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no se pudieron cargar roles"})
		return
	}
	if len(roles) != len(want) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "uno o más roles no existen"})
		return
	}

	// Reemplazar roles
	if err := h.DB.Model(&user).Association("Roles").Replace(roles); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no se pudieron actualizar roles"})
		return
	}

	out := make([]string, 0, len(roles))
	for _, r := range roles {
		out = append(out, r.Name)
	}

	c.JSON(http.StatusOK, gin.H{
		"user_id": user.ID,
		"roles":   out,
	})
}
