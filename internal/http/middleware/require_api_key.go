package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type APIKeyConfig struct {
	
	ValidKeys map[string]bool


	HeaderName string
}

func RequireAPIKey(cfg APIKeyConfig) gin.HandlerFunc {
	headerName := cfg.HeaderName
	if headerName == "" {
		headerName = "X-API-Key"
	}

	return func(c *gin.Context) {
		key := strings.TrimSpace(c.GetHeader(headerName))
		if key == "" || !cfg.ValidKeys[key] {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "missing_or_invalid_api_key",
			})
			return
		}

		c.Next()
	}
}
