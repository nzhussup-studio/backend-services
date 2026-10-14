package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware trusts identity only after the edge gateway has validated the
// token and replaced the X-User-* headers. Network policy must prevent callers
// from reaching this service without passing through that gateway.
func AuthMiddleware(config *AuthConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == http.MethodGet && !shouldAuthenticate(c, config.Rules) {
			c.Next()
			return
		}

		identity, err := IdentityFromHeaders(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		if !isAdmin(identity.Roles) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
			c.Abort()
			return
		}

		c.Set("subject", identity.Subject)
		c.Set("username", identity.Username)
		c.Set("roles", identity.Roles)
		c.Next()
	}
}
