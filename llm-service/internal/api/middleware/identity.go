package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// RequireGatewayIdentity requires the identity header injected by Envoy after
// it has authenticated and authorized the request.
func RequireGatewayIdentity() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if strings.TrimSpace(ctx.GetHeader("X-User-ID")) == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		ctx.Next()
	}
}
