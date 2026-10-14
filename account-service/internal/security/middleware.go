package security

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		subject := strings.TrimSpace(c.GetHeader("X-User-ID"))
		if subject == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		c.Set("subject", subject)
		c.Set("username", strings.TrimSpace(c.GetHeader("X-User-Name")))
		c.Set("email", strings.TrimSpace(c.GetHeader("X-User-Email")))
		c.Next()
	}
}
