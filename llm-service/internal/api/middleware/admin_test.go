package middleware

import (
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestRequireAdmin(t *testing.T) {
	gin.SetMode(gin.TestMode)
	tests := []struct {
		name   string
		userID string
		roles  string
		want   int
	}{
		{name: "missing identity", want: http.StatusUnauthorized},
		{name: "normal user", userID: "user-1", roles: encodeRoles(`["ROLE_USER"]`), want: http.StatusForbidden},
		{name: "admin", userID: "admin-1", roles: encodeRoles(`["ROLE_USER","ROLE_ADMIN"]`), want: http.StatusNoContent},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			router.GET("/configuration", RequireAdmin(), func(ctx *gin.Context) { ctx.Status(http.StatusNoContent) })
			req := httptest.NewRequest(http.MethodGet, "/configuration", nil)
			req.Header.Set("X-User-ID", tt.userID)
			req.Header.Set("X-User-Roles", tt.roles)
			response := httptest.NewRecorder()
			router.ServeHTTP(response, req)
			if response.Code != tt.want {
				t.Fatalf("expected %d, got %d", tt.want, response.Code)
			}
		})
	}
}

func encodeRoles(value string) string {
	return base64.StdEncoding.EncodeToString([]byte(value))
}
