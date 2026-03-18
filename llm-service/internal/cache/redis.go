package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	ErrNotFound = fmt.Errorf("cache: key not found")
	ErrCache    = fmt.Errorf("cache: internal error")
)

type Cacher interface {
	Set(key string, value any) error
	Get(key string, dest any) error
	Del(key string) error
}

type RedisAPI interface {
	Set(ctx context.Context, key string, value any, expiration time.Duration) *redis.StatusCmd
	Get(ctx context.Context, key string) *redis.StringCmd
	Del(ctx context.Context, keys ...string) *redis.IntCmd
	Ping(ctx context.Context) *redis.StatusCmd
}

type RedisClient struct {
	Client   RedisAPI
	Duration time.Duration
}

const wrapper = "%w: %v"

func NewRedisClient(addr, password string, db int, duration time.Duration) *RedisClient {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	return &RedisClient{
		Client:   client,
		Duration: duration,
	}
}

func (r *RedisClient) Set(key string, value any) error {
	ctx := context.Background()

	jsonData, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf(wrapper, err, "failed to marshal value for cache")
	}

	err = r.Client.Set(ctx, key, jsonData, r.Duration).Err()
	if err != nil {
		return fmt.Errorf(wrapper, err, "failed to set value in cache")
	}

	return nil
}

func (r *RedisClient) Get(key string, dest any) error {
	ctx := context.Background()

	val, err := r.Client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return fmt.Errorf(wrapper, ErrNotFound, "key does not exist in cache")
		}
		return fmt.Errorf(wrapper, err, "failed to get value from cache")
	}

	err = json.Unmarshal([]byte(val), dest)
	if err != nil {
		return fmt.Errorf(wrapper, err, "failed to unmarshal value from cache")
	}

	return nil
}

func (r *RedisClient) Del(key string) error {
	ctx := context.Background()

	err := r.Client.Del(ctx, key).Err()
	if err != nil {
		return fmt.Errorf(wrapper, err, "failed to delete cache key")
	}
	return nil
}

func (r *RedisClient) Ping() error {
	ctx := context.Background()

	if err := r.Client.Ping(ctx).Err(); err != nil {
		return fmt.Errorf(wrapper, err, "failed to ping Redis")
	}
	return nil
}
