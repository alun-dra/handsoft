package handlers

import (
	"net/http"
	"strconv"

	"handsoft/internal/models"

	"github.com/gin-gonic/gin"
)

func (h *AdminHandler) ListPermissions(c *gin.Context) {
	var perms []models.Permission
	if err := h.DB.Order("id asc").Find(&perms).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db_error"})
		return
	}
	c.JSON(http.StatusOK, perms)
}

func (h *AdminHandler) GetRolePermissions(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_id"})
		return
	}

	var role models.Role
	if err := h.DB.Preload("Permissions").First(&role, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "role_not_found"})
		return
	}

	c.JSON(http.StatusOK, role.Permissions)
}

type setRolePermsReq struct {
	PermissionCodes []string `json:"permission_codes" binding:"required"`
}

func (h *AdminHandler) SetRolePermissions(c *gin.Context) {
	roleID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_id"})
		return
	}

	var role models.Role
	if err := h.DB.First(&role, roleID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "role_not_found"})
		return
	}

	var req setRolePermsReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_body"})
		return
	}

	var perms []models.Permission
	if len(req.PermissionCodes) > 0 {
		if err := h.DB.Where("code IN ?", req.PermissionCodes).Find(&perms).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "db_error"})
			return
		}
	}

	if err := h.DB.Model(&role).Association("Permissions").Replace(perms); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot_set_permissions"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"role_id":  role.ID,
		"assigned": len(perms),
	})
}
