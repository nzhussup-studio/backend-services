package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
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

func AuthMiddlewareWithClient(authServiceURL string, client *http.Client) gin.HandlerFunc {
	if client == nil {
		client = http.DefaultClient
	}

	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.Next()
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")

		payload, err := json.Marshal(ValidationRequest{Token: token})
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal request"})
			return
		}

		req, err := http.NewRequest("POST", authServiceURL+"/auth/validate", bytes.NewBuffer(payload))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
			return
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token validation failed"})
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Access denied", "details": string(body)})
			return
		}

		var result ValidationResponse
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode validation response"})
			return
		}

		c.Request.Header.Set("X-User-Username", result.Username)
		if len(result.Roles) > 0 {
			c.Request.Header.Set("X-User-Roles", strings.Join(result.Roles, ","))
		} else {
			c.Request.Header.Set("X-User-Roles", "")
		}

		c.Next()
	}
}
