package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// MockRedisAPI implements RedisAPI interface with customizable funcs
type MockRedisAPI struct {
	PingFunc     func(ctx context.Context) error
	SetFunc      func(ctx context.Context, key string, value any, expiration time.Duration) error
	GetFunc      func(ctx context.Context, key string) (string, error)
	DelFunc      func(ctx context.Context, keys ...string) error
	FlushAllFunc func(ctx context.Context) error
}

func (m *MockRedisAPI) Ping(ctx context.Context) *redis.StatusCmd {
	cmd := redis.NewStatusCmd(ctx)
	if m.PingFunc != nil {
		cmd.SetErr(m.PingFunc(ctx))
	}
	return cmd
}

func (m *MockRedisAPI) Set(ctx context.Context, key string, value any, expiration time.Duration) *redis.StatusCmd {
	cmd := redis.NewStatusCmd(ctx)
	if m.SetFunc != nil {
		cmd.SetErr(m.SetFunc(ctx, key, value, expiration))
	}
	return cmd
}

func (m *MockRedisAPI) Get(ctx context.Context, key string) *redis.StringCmd {
	cmd := redis.NewStringCmd(ctx)
	if m.GetFunc != nil {
		val, err := m.GetFunc(ctx, key)
		if err != nil {
			cmd.SetErr(err)
		} else {
			cmd.SetVal(val)
		}
	}
	return cmd
}

func (m *MockRedisAPI) Del(ctx context.Context, keys ...string) *redis.IntCmd {
	cmd := redis.NewIntCmd(ctx)
	if m.DelFunc != nil {
		cmd.SetErr(m.DelFunc(ctx, keys...))
	}
	return cmd
}

func (m *MockRedisAPI) FlushAll(ctx context.Context) *redis.StatusCmd {
	cmd := redis.NewStatusCmd(ctx)
	if m.FlushAllFunc != nil {
		cmd.SetErr(m.FlushAllFunc(ctx))
	}
	return cmd
}
