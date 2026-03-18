package controller

import (
	"errors"
	"image-service/internal/service"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCacheController_ClearCache_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	mockService := new(MockCacheService)
	mockService.On("ClearCache").Return(nil)

	svc := &service.Service{CacheService: mockService}
	controller := NewCacheController(svc)

	router.DELETE("/v1/album/cache", controller.ClearCache)

	req, _ := http.NewRequest(http.MethodDelete, "/v1/album/cache", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	mockService.AssertExpectations(t)
}

func TestCacheController_ClearCache_Failure(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	mockService := new(MockCacheService)
	mockService.On("ClearCache").Return(errors.New("mock cache error"))

	svc := &service.Service{CacheService: mockService}
	controller := NewCacheController(svc)

	router.DELETE("/v1/album/cache", controller.ClearCache)

	req, _ := http.NewRequest(http.MethodDelete, "/v1/album/cache", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusInternalServerError, resp.Code)
	mockService.AssertExpectations(t)
}
