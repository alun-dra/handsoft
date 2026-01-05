package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"handsoft/internal/models"
	"handsoft/internal/http/routes"

	"github.com/gin-gonic/gin"
)

type createRoleReq struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	IsSuperAdmin bool  `json:"is_super_admin"`
}

func ListRoles(deps routes.Deps) gin.HandlerFunc {
	return func(c *gin.Context) {
		var roles []models.Role
		if err := deps.DB.Order("id asc").Find(&roles).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "db_error"})
			return
		}
		c.JSON(http.StatusOK, roles)
	}
}

func CreateRole(deps routes.Deps) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req createRoleReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_body"})
			return
		}

		req.Name = strings.TrimSpace(req.Name)
		if req.Name == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "name_required"})
			return
		}

		role := models.Role{
			Name:         req.Name,
			Description:  req.Description,
			IsSuperAdmin: req.IsSuperAdmin,
		}

		if err := deps.DB.Create(&role).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "cannot_create_role"})
			return
		}

		c.JSON(http.StatusCreated, role)
	}
}

type updateRoleReq struct {
	Name         *string `json:"name"`
	Description  *string `json:"description"`
	IsSuperAdmin *bool   `json:"is_super_admin"`
}

func UpdateRole(deps routes.Deps) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_id"})
			return
		}

		var role models.Role
		if err := deps.DB.First(&role, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "role_not_found"})
			return
		}

		var req updateRoleReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_body"})
			return
		}

		if req.Name != nil {
			name := strings.TrimSpace(*req.Name)
			if name == "" {
				c.JSON(http.StatusBadRequest, gin.H{"error": "name_empty"})
				return
			}
			role.Name = name
		}
		if req.Description != nil {
			role.Description = *req.Description
		}
		if req.IsSuperAdmin != nil {
			role.IsSuperAdmin = *req.IsSuperAdmin
		}

		if err := deps.DB.Save(&role).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "cannot_update_role"})
			return
		}

		c.JSON(http.StatusOK, role)
	}
}

func DeleteRole(deps routes.Deps) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_id"})
			return
		}

		// opcional: proteger super_admin para no borrarlo accidentalmente
		var role models.Role
		if err := deps.DB.First(&role, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "role_not_found"})
			return
		}
		if role.IsSuperAdmin {
			c.JSON(http.StatusBadRequest, gin.H{"error": "cannot_delete_super_admin_role"})
			return
		}

		if err := deps.DB.Delete(&models.Role{}, id).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot_delete_role"})
			return
		}

		c.Status(http.StatusNoContent)
	}
}
