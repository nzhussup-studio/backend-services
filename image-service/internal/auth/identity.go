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

// RequiresIdentity reports whether this request is one of the protected reads.
// Other protected routes are guarded by the middleware unconditionally.
func RequiresIdentity(c *gin.Context, rules []ProtectedReadRule) bool {
	requestPath := c.Request.URL.Path
	for _, rule := range rules {
		if rule.Path != requestPath {
			continue
		}
		match := true
		for key, val := range rule.QueryParams {
			if c.Query(key) != val {
				match = false
				break
			}
		}
		if match {
			return true
		}
	}
	return false
}

// RequireIdentity verifies a gateway-injected identity for a request that Envoy
// has already authorized.
func RequireIdentity(c *gin.Context) error {
	_, err := IdentityFromHeaders(c)
	return err
}
