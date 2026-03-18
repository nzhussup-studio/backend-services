package service

import (
	"testing"

	"image-service/internal/config/cache"
	"image-service/internal/config/security"
	"image-service/internal/repository"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

func TestNewService(t *testing.T) {
	mockStorage := &repository.Storage{}
	mockRedis := &cache.RedisClient{}
	mockSecurity := &security.AuthConfig{}
	validate := validator.New()

	svc := NewService(mockStorage, mockRedis, mockSecurity, validate)

	assert.NotNil(t, svc)
	assert.Equal(t, mockStorage, svc.storage)
	assert.Equal(t, validate, svc.validate)

	assert.NotNil(t, svc.AlbumService)
	assert.NotNil(t, svc.ImageService)
	assert.NotNil(t, svc.CacheService)

	_, ok := svc.AlbumService.(*AlbumService)
	assert.True(t, ok, "AlbumService should be of type *AlbumService")

	_, ok = svc.ImageService.(*ImageService)
	assert.True(t, ok, "ImageService should be of type *ImageService")

	_, ok = svc.CacheService.(*CacheService)
	assert.True(t, ok, "CacheService should be of type *CacheService")
}
