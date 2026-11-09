package summarizer

import (
	"encoding/json"
	"fmt"
	"sync"
	"testing"
	"time"

	"llm-service/internal/cache"

	"github.com/stretchr/testify/require"
)

type mockCache struct {
	mu    sync.Mutex
	store map[string][]byte
}

func newMockCache() *mockCache                       { return &mockCache{store: map[string][]byte{}} }
func (m *mockCache) Set(key string, value any) error { return m.SetWithTTL(key, value, 0) }
func (m *mockCache) SetWithTTL(key string, value any, _ time.Duration) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	b, err := json.Marshal(value)
	if err == nil {
		m.store[key] = b
	}
	return err
}
func (m *mockCache) SetNX(key string, value any, _ time.Duration) (bool, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.store[key]; ok {
		return false, nil
	}
	b, err := json.Marshal(value)
	if err == nil {
		m.store[key] = b
	}
	return true, err
}
func (m *mockCache) Get(key string, dest any) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	b, ok := m.store[key]
	if !ok {
		return fmt.Errorf("%w: missing", cache.ErrNotFound)
	}
	return json.Unmarshal(b, dest)
}
func (m *mockCache) Del(key string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.store, key)
	return nil
}

func TestGetSummaryState(t *testing.T) {
	store := newMockCache()
	s := NewSummarizer("", "", nil, nil, store)

	missing, err := s.Get("en")
	require.NoError(t, err)
	require.Equal(t, "missing", missing.Status)

	require.NoError(t, store.Set("llm:generation:en", Generation{JobID: "job", Status: "generating"}))
	running, err := s.Get("en")
	require.NoError(t, err)
	require.Equal(t, "generating", running.Status)
	require.Equal(t, 1500, running.PollAfterMS)

	require.NoError(t, store.Set("llm:summary:en", Summary{Text: "Ready"}))
	completed, err := s.Get("en")
	require.NoError(t, err)
	require.Equal(t, "completed", completed.Status)
	require.Equal(t, "Ready", completed.Summary)
}

func TestHashValueIsStable(t *testing.T) {
	first, err := hashValue(map[string]string{"name": "Nurzhanat"})
	require.NoError(t, err)
	second, err := hashValue(map[string]string{"name": "Nurzhanat"})
	require.NoError(t, err)
	require.Equal(t, first, second)
}
