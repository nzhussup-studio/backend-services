package auth

import (
	"encoding/base64"
	"encoding/json"
	"strings"

	custom_errors "image-service/internal/errors"

	"github.com/gin-gonic/gin"
)

type Identity struct {
	Subject  string
	Username string
	Roles    []string
}

func IdentityFromHeaders(c *gin.Context) (Identity, error) {
	subject := strings.TrimSpace(c.GetHeader("X-User-ID"))
	if subject == "" {
		return Identity{}, custom_errors.NewError(custom_errors.ErrUnauthorized, "Trusted identity is missing")
	}

	roles, err := decodeRoles(c.GetHeader("X-User-Roles"))
	if err != nil {
		return Identity{}, custom_errors.NewError(custom_errors.ErrUnauthorized, "Trusted roles are invalid")
	}

	return Identity{
		Subject:  subject,
		Username: firstNonEmpty(strings.TrimSpace(c.GetHeader("X-User-Name")), subject),
		Roles:    roles,
	}, nil
}

func decodeRoles(value string) ([]string, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil, nil
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
		if err := json.Unmarshal(decoded, &roles); err == nil {
			return roles, nil
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
		return roles, nil
	}

	return nil, custom_errors.NewError(custom_errors.ErrUnauthorized, "Invalid roles header")
}

func isAdmin(roles []string) bool {
	for _, role := range roles {
		if role == "ROLE_ADMIN" {
			return true
		}
	}
	return false
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return "unknown"
}

func shouldAuthenticate(c *gin.Context, rules []AuthRule) bool {
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

var CheckIsAdmin = func(c *gin.Context, _ *AuthConfig) error {
	identity, err := IdentityFromHeaders(c)
	if err != nil {
		return err
	}
	if isAdmin(identity.Roles) {
		return nil
	}
	return custom_errors.NewError(custom_errors.ErrForbidden, "User is not an admin")
}
