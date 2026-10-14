package middleware

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func RequireAdmin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if strings.TrimSpace(ctx.GetHeader("X-User-ID")) == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		for _, role := range decodeTrustedRoles(ctx.GetHeader("X-User-Roles")) {
			if role == "ROLE_ADMIN" {
				ctx.Next()
				return
			}
		}

		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
	}
}

func decodeTrustedRoles(value string) []string {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil
	}

	encodings := []*base64.Encoding{
		base64.StdEncoding,
		base64.RawStdEncoding,
		base64.URLEncoding,
		base64.RawURLEncoding,
	}
	for _, encoding := range encodings {
		decoded, err := encoding.DecodeString(value)
		if err != nil {
			continue
		}
		var roles []string
		if json.Unmarshal(decoded, &roles) == nil {
			return roles
		}
	}

	if !strings.ContainsAny(value, "[]{}\"") {
		parts := strings.Split(value, ",")
		roles := make([]string, 0, len(parts))
		for _, part := range parts {
			if role := strings.TrimSpace(part); role != "" {
				roles = append(roles, role)
			}
		}
		return roles
	}

	return nil
}
