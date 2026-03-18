package cache

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRedisClient mocks the RedisAPI interface
type MockRedisClient struct {
	mock.Mock
}

func (m *MockRedisClient) Set(ctx context.Context, key string, value any, expiration time.Duration) *redis.StatusCmd {
	args := m.Called(ctx, key, value, expiration)
	return args.Get(0).(*redis.StatusCmd)
}

func (m *MockRedisClient) Get(ctx context.Context, key string) *redis.StringCmd {
	args := m.Called(ctx, key)
	return args.Get(0).(*redis.StringCmd)
}

func (m *MockRedisClient) Del(ctx context.Context, keys ...string) *redis.IntCmd {
	args := m.Called(ctx, keys)
	return args.Get(0).(*redis.IntCmd)
}

func (m *MockRedisClient) Ping(ctx context.Context) *redis.StatusCmd {
	args := m.Called(ctx)
	return args.Get(0).(*redis.StatusCmd)
}

// Helper functions to create redis.Cmd mocks
func newStatusCmd(err error) *redis.StatusCmd {
	cmd := redis.NewStatusCmd(context.Background())
	cmd.SetErr(err)
	return cmd
}

func newStringCmd(val string, err error) *redis.StringCmd {
	cmd := redis.NewStringCmd(context.Background())
	cmd.SetVal(val)
	cmd.SetErr(err)
	return cmd
}

func newIntCmd(val int64, err error) *redis.IntCmd {
	cmd := redis.NewIntCmd(context.Background())
	cmd.SetVal(val)
	cmd.SetErr(err)
	return cmd
}

func newPingStatusCmd(val string, err error) *redis.StatusCmd {
	cmd := redis.NewStatusCmd(context.Background())
	cmd.SetVal(val)
	cmd.SetErr(err)
	return cmd
}

func TestRedisClient_SetGetDel(t *testing.T) {
	mockRedis := new(MockRedisClient)
	client := &RedisClient{
		Client:   mockRedis,
		Duration: time.Minute,
	}

	type testCase struct {
		run         func() error
		setupMock   func()
		expectedErr error
	}

	key := "testkey"
	value := map[string]string{"foo": "bar"}
	valueJSON, _ := json.Marshal(value)

	testCases := map[string]testCase{
		"set_success": {
			setupMock: func() {
				mockRedis.On("Set", mock.Anything, key, valueJSON, client.Duration).
					Return(newStatusCmd(nil)).Once()
			},
			run: func() error {
				return client.Set(key, value)
			},
			expectedErr: nil,
		},
		"set_marshal_failure": {
			setupMock: func() {}, // no mock needed
			run: func() error {
				// Trying to set a non-marshalable value (channel)
				return client.Set(key, make(chan int))
			},
			expectedErr: errors.New("failed to marshal value for cache"),
		},
		"get_success": {
			setupMock: func() {
				mockRedis.On("Get", mock.Anything, key).
					Return(newStringCmd(string(valueJSON), nil)).Once()
			},
			run: func() error {
				var dest map[string]string
				return client.Get(key, &dest)
			},
			expectedErr: nil,
		},
		"get_not_found": {
			setupMock: func() {
				mockRedis.On("Get", mock.Anything, key).
					Return(newStringCmd("", redis.Nil)).Once()
			},
			run: func() error {
				var dest map[string]string
				return client.Get(key, &dest)
			},
			expectedErr: ErrNotFound,
		},
		"get_unmarshal_failure": {
			setupMock: func() {
				// Return invalid JSON
				mockRedis.On("Get", mock.Anything, key).
					Return(newStringCmd("invalid_json", nil)).Once()
			},
			run: func() error {
				var dest map[string]string
				return client.Get(key, &dest)
			},
			expectedErr: errors.New("failed to unmarshal value from cache"),
		},
		"del_success": {
			setupMock: func() {
				mockRedis.On("Del", mock.Anything, []string{key}).
					Return(newIntCmd(1, nil)).Once()
			},
			run: func() error {
				return client.Del(key)
			},
			expectedErr: nil,
		},
		"del_failure": {
			setupMock: func() {
				mockRedis.On("Del", mock.Anything, []string{key}).
					Return(newIntCmd(0, errors.New("some del error"))).Once()
			},
			run: func() error {
				return client.Del(key)
			},
			expectedErr: errors.New("failed to delete cache key"),
		},
		"ping_success": {
			setupMock: func() {
				mockRedis.On("Ping", mock.Anything).
					Return(newPingStatusCmd("PONG", nil)).Once()
			},
			run: func() error {
				return client.Ping()
			},
			expectedErr: nil,
		},
		"ping_failure": {
			setupMock: func() {
				mockRedis.On("Ping", mock.Anything).
					Return(newPingStatusCmd("", errors.New("ping error"))).Once()
			},
			run: func() error {
				return client.Ping()
			},
			expectedErr: errors.New("failed to ping Redis"),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tc.setupMock()
			err := tc.run()
			if tc.expectedErr != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}
			mockRedis.AssertExpectations(t)
			mockRedis.ExpectedCalls = nil // reset for next test
		})
	}
}
