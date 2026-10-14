package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// TrustedIdentityMiddleware trusts identity only after the edge gateway has validated the
// token and replaced the X-User-* headers. Network policy must prevent callers
// from reaching this service without passing through that gateway.
func TrustedIdentityMiddleware(protectedReads *ProtectedReadRules) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == http.MethodGet && !RequiresIdentity(c, protectedReads.Rules) {
			c.Next()
			return
		}

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
