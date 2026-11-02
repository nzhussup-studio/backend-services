package auth

import (
	"strings"

	custom_errors "image-service/internal/errors"

	"github.com/gin-gonic/gin"
)

type Identity struct {
	Subject  string
	Username string
}

func IdentityFromHeaders(c *gin.Context) (Identity, error) {
	subject := strings.TrimSpace(c.GetHeader("X-User-ID"))
	if subject == "" {
		return Identity{}, custom_errors.NewError(custom_errors.ErrUnauthorized, "Trusted identity is missing")
	}

	return Identity{
		Subject:  subject,
		Username: firstNonEmpty(strings.TrimSpace(c.GetHeader("X-User-Name")), subject),
	}, nil
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return "unknown"
}
