package middleware

import (
	"net/http"
	"strings"

	"handsoft/internal/auth"

	"github.com/gin-gonic/gin"
)

const (
	CtxUserIDKey = "userID"
	CtxRolesKey  = "roles"
)

func AuthJWT(cfg auth.JWTConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		h := c.GetHeader("Authorization")
		if h == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "faltó Authorization header"})
			return
		}

		parts := strings.SplitN(h, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "formato Authorization inválido (Bearer token)"})
			return
		}

		tokenStr := strings.TrimSpace(parts[1])
		if tokenStr == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token vacío"})
			return
		}

		claims, err := auth.VerifyAccessToken(cfg, tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token inválido o expirado"})
			return
		}

		// Guardamos datos útiles para handlers
		c.Set(CtxUserIDKey, claims.UserID)
		c.Set(CtxRolesKey, claims.Roles)

		c.Next()
	}
}
