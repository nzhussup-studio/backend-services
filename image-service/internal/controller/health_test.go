package controller

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestHealthCheckHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/v1/album/health", nil)

	HealthCheckHandler(c)

	assert.Equal(t, http.StatusOK, w.Code)
	expected := `{"status":200}` + "\n"
	assert.JSONEq(t, expected, w.Body.String())
}
