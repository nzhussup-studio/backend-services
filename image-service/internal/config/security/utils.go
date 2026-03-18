package security

import (
	"bytes"
	"encoding/json"
	custom_errors "image-service/internal/errors"
	"io"
	"net/http"
	"slices"
	"strings"

	"github.com/gin-gonic/gin"
)

type ValidationRequest struct {
	Token string `json:"token"`
}

type ValidationResponse struct {
	Username string   `json:"username"`
	Roles    []string `json:"roles"`
}

func Validate(c *gin.Context, config *AuthConfig, token string) (ValidationResponse, error) {

	reqBody, _ := json.Marshal(ValidationRequest{Token: token})

	var resp *http.Response
	resp, err := http.Post(config.AuthServiceURL+config.ValidationURL, "application/json", bytes.NewBuffer(reqBody))
	if err != nil || resp.StatusCode == http.StatusInternalServerError {
		return ValidationResponse{}, custom_errors.NewError(custom_errors.ErrInternalServer, "Failed to validate token")
	}
	if resp.StatusCode != http.StatusOK {
		return ValidationResponse{}, custom_errors.NewError(custom_errors.ErrUnauthorized, "Invalid token")
	}

	body, _ := io.ReadAll(resp.Body)
	var validationResponse ValidationResponse
	json.Unmarshal(body, &validationResponse)

	return validationResponse, nil
}

func GetToken(c *gin.Context) (string, error) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		return "", custom_errors.NewError(custom_errors.ErrUnauthorized, "Authorization header missing or invalid")
	}
	token := strings.TrimPrefix(authHeader, "Bearer ")
	return token, nil
}

func isAdmin(roles []string) bool {
	return slices.Contains(roles, "ROLE_ADMIN")
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

var CheckIsAdmin = func(c *gin.Context, config *AuthConfig) error {
	token, err := GetToken(c)
	if err != nil {
		return err
	}
	validationResponse, err := Validate(c, config, token)
	if err != nil {
		return err
	}
	if isAdmin(validationResponse.Roles) {
		return nil
	}
	return custom_errors.NewError(custom_errors.ErrForbidden, "User is not an admin")
}
