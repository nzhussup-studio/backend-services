package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestIdentityFromHeadersReadsTrustedIdentity(t *testing.T) {
	gin.SetMode(gin.TestMode)
	request := httptest.NewRequest(http.MethodGet, "/", nil)
	request.Header.Set("X-User-ID", "user-123")
	request.Header.Set("X-User-Name", "admin")
	context, _ := gin.CreateTestContext(httptest.NewRecorder())
	context.Request = request

	identity, err := IdentityFromHeaders(context)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if identity.Subject != "user-123" || identity.Username != "admin" {
		t.Fatalf("unexpected identity: %#v", identity)
	}
}

func TestIdentityFromHeadersRejectsMissingSubject(t *testing.T) {
	gin.SetMode(gin.TestMode)
	context, _ := gin.CreateTestContext(httptest.NewRecorder())
	context.Request = httptest.NewRequest(http.MethodGet, "/", nil)

	if _, err := IdentityFromHeaders(context); err == nil {
		t.Fatal("expected missing identity to be rejected")
	}
}
