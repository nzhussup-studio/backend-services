package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClearCache(t *testing.T) {
	mockRedis := new(MockRedisClient)

	cacheService := &CacheService{
		redis: mockRedis,
	}

	mockRedis.On("FlushAll").Return(nil)

	err := cacheService.ClearCache()

	assert.NoError(t, err)
	mockRedis.AssertExpectations(t)
}
