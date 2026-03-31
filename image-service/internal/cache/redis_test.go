package cache

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestCheckHealth(t *testing.T) {
	mock := &MockRedisAPI{
		PingFunc: func(ctx context.Context) error { return nil },
	}
	r := &RedisClient{Client: mock}
	r.CheckHealth()

	mock.PingFunc = func(ctx context.Context) error { return errors.New("ping error") }
	r.CheckHealth()
}

func TestSet(t *testing.T) {
	mock := &MockRedisAPI{
		SetFunc: func(ctx context.Context, key string, value any, expiration time.Duration) error {
			assert.Equal(t, "key1", key)
			b, err := json.Marshal(value)
			assert.NoError(t, err)
			assert.NotEmpty(t, b)
			return nil
		},
	}
	r := &RedisClient{Client: mock, Duration: time.Second * 10}
	r.Set("key1", map[string]string{"foo": "bar"})
}

func TestSet_MarshalError(t *testing.T) {
	mock := &MockRedisAPI{
		SetFunc: func(ctx context.Context, key string, value any, expiration time.Duration) error {
			return nil
		},
	}
	r := &RedisClient{Client: mock, Duration: time.Second * 10}
	r.Set("key", make(chan int))
}

func TestGet_Success(t *testing.T) {
	expected := map[string]string{"foo": "bar"}
	jsonData, _ := json.Marshal(expected)

	mock := &MockRedisAPI{
		GetFunc: func(ctx context.Context, key string) (string, error) {
			assert.Equal(t, "key1", key)
			return string(jsonData), nil
		},
	}
	r := &RedisClient{Client: mock}
	var dest map[string]string
	err := r.Get("key1", &dest)
	assert.NoError(t, err)
	assert.Equal(t, expected, dest)
}

func TestGet_KeyNotFound(t *testing.T) {
	mock := &MockRedisAPI{
		GetFunc: func(ctx context.Context, key string) (string, error) {
			return "", redis.Nil
		},
	}
	r := &RedisClient{Client: mock}
	var dest map[string]string
	err := r.Get("missing", &dest)
	assert.Error(t, err)
}

func TestGet_UnmarshalError(t *testing.T) {
	mock := &MockRedisAPI{
		GetFunc: func(ctx context.Context, key string) (string, error) {
			return "invalid json", nil
		},
	}
	r := &RedisClient{Client: mock}
	var dest map[string]string
	err := r.Get("key", &dest)
	assert.Error(t, err)
}

func TestDel_Success(t *testing.T) {
	called := false
	mock := &MockRedisAPI{
		DelFunc: func(ctx context.Context, keys ...string) error {
			called = true
			assert.Equal(t, []string{"key1"}, keys)
			return nil
		},
	}
	r := &RedisClient{Client: mock}
	r.Del("key1")
	assert.True(t, called)
}

func TestDel_FailureFlushAllSuccess(t *testing.T) {
	delCalled := false
	flushCalled := false

	mock := &MockRedisAPI{
		DelFunc: func(ctx context.Context, keys ...string) error {
			delCalled = true
			return errors.New("del error")
		},
		FlushAllFunc: func(ctx context.Context) error {
			flushCalled = true
			return nil
		},
	}
	r := &RedisClient{Client: mock}
	r.Del("key1")
	assert.True(t, delCalled)
	assert.True(t, flushCalled)
}

func TestDel_FailureFlushAllFailure(t *testing.T) {
	delCalled := false
	flushCalled := false

	mock := &MockRedisAPI{
		DelFunc: func(ctx context.Context, keys ...string) error {
			delCalled = true
			return errors.New("del error")
		},
		FlushAllFunc: func(ctx context.Context) error {
			flushCalled = true
			return errors.New("flush error")
		},
	}
	r := &RedisClient{Client: mock}
	r.Del("key1")
	assert.True(t, delCalled)
	assert.True(t, flushCalled)
}

func TestFlushAll_Success(t *testing.T) {
	mock := &MockRedisAPI{
		FlushAllFunc: func(ctx context.Context) error {
			return nil
		},
	}
	r := &RedisClient{Client: mock}
	err := r.FlushAll()
	assert.NoError(t, err)
}

func TestFlushAll_Failure(t *testing.T) {
	mock := &MockRedisAPI{
		FlushAllFunc: func(ctx context.Context) error {
			return errors.New("flush error")
		},
	}
	r := &RedisClient{Client: mock}
	err := r.FlushAll()
	assert.NoError(t, err)
}
