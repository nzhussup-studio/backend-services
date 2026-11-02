package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RequireTrustedIdentityMiddleware protects an entire route group. The edge
// gateway is responsible for validating the token and admin role before
// injecting the trusted identity headers.
func RequireTrustedIdentityMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		identity, err := IdentityFromHeaders(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}
		c.Set("subject", identity.Subject)
		c.Set("username", identity.Username)
		c.Next()
	}
}
