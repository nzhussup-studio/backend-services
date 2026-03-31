package service

import (
	"testing"

	"image-service/internal/auth"
	"image-service/internal/cache"
	appconfig "image-service/internal/config"
	"image-service/internal/repository"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

func TestNewService(t *testing.T) {
	mockStorage := &repository.Storage{}
	mockRedis := &cache.RedisClient{}
	mockSecurity := &auth.AuthConfig{}
	validate := validator.New()
	imageCfg := appconfig.ImageConfig{}

	svc := NewService(mockStorage, mockRedis, mockSecurity, validate, imageCfg)

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
