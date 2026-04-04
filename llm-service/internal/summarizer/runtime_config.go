package summarizer

import (
	"errors"
	"fmt"
	"llm-service/internal/cache"
	"strings"
	"sync"
	"time"
)

const (
	defaultModelID        = "nvidia/nemotron-3-super-120b-a12b:free"
	maxPromptChars        = 8000
	runtimeConfigCacheKey = "llm_runtime_config"
)

type RuntimeConfig struct {
	Model                    string
	SystemPromptEN           string
	SystemPromptDE           string
	SystemPromptKZ           string
	EnableParallelGeneration bool
}

type RuntimeConfigUpdate struct {
	Model                    *string
	SystemPromptEN           *string
	SystemPromptDE           *string
	SystemPromptKZ           *string
	EnableParallelGeneration *bool
}

var runtimeConfig = struct {
	mu  sync.RWMutex
	cfg RuntimeConfig
}{
	cfg: RuntimeConfig{
		Model:                    defaultModelID,
		SystemPromptEN:           SYSTEM_PROMPT_EN,
		SystemPromptDE:           SYSTEM_PROMPT_DE,
		SystemPromptKZ:           SYSTEM_PROMPT_KZ,
		EnableParallelGeneration: true,
	},
}

var runtimeConfigStore cache.Cacher

type runtimeConfigTTLStore interface {
	SetWithTTL(key string, value any, ttl time.Duration) error
}

func GetRuntimeConfig() RuntimeConfig {
	runtimeConfig.mu.RLock()
	defer runtimeConfig.mu.RUnlock()
	return runtimeConfig.cfg
}

func InitRuntimeConfigStore(store cache.Cacher) error {
	if store == nil {
		return nil
	}
	runtimeConfigStore = store

	var persisted RuntimeConfig
	err := store.Get(runtimeConfigCacheKey, &persisted)
	if err != nil {
		if errors.Is(err, cache.ErrNotFound) {
			return nil
		}
		return fmt.Errorf("failed to load runtime config from redis: %w", err)
	}

	// Empty persisted payload should not override runtime defaults.
	if strings.TrimSpace(persisted.Model) == "" &&
		strings.TrimSpace(persisted.SystemPromptEN) == "" &&
		strings.TrimSpace(persisted.SystemPromptDE) == "" &&
		strings.TrimSpace(persisted.SystemPromptKZ) == "" {
		return nil
	}

	runtimeConfig.mu.Lock()
	defer runtimeConfig.mu.Unlock()

	cfg := runtimeConfig.cfg
	if strings.TrimSpace(persisted.Model) != "" {
		cfg.Model = strings.TrimSpace(persisted.Model)
	}
	if strings.TrimSpace(persisted.SystemPromptEN) != "" {
		cfg.SystemPromptEN = strings.TrimSpace(persisted.SystemPromptEN)
	}
	if strings.TrimSpace(persisted.SystemPromptDE) != "" {
		cfg.SystemPromptDE = strings.TrimSpace(persisted.SystemPromptDE)
	}
	if strings.TrimSpace(persisted.SystemPromptKZ) != "" {
		cfg.SystemPromptKZ = strings.TrimSpace(persisted.SystemPromptKZ)
	}
	cfg.EnableParallelGeneration = persisted.EnableParallelGeneration

	runtimeConfig.cfg = cfg
	return nil
}

func UpdateRuntimeConfig(update RuntimeConfigUpdate) error {
	if update.Model == nil &&
		update.SystemPromptEN == nil &&
		update.SystemPromptDE == nil &&
		update.SystemPromptKZ == nil &&
		update.EnableParallelGeneration == nil {
		return fmt.Errorf("at least one configuration field must be provided")
	}

	var modelValue string
	if update.Model != nil {
		modelValue = strings.TrimSpace(*update.Model)
		if modelValue == "" {
			return fmt.Errorf("model cannot be empty")
		}
	}

	var promptEN string
	if update.SystemPromptEN != nil {
		promptEN = strings.TrimSpace(*update.SystemPromptEN)
		if promptEN == "" {
			return fmt.Errorf("system_prompt_en cannot be empty")
		}
		if len(promptEN) > maxPromptChars {
			return fmt.Errorf("system_prompt_en exceeds %d characters", maxPromptChars)
		}
	}

	var promptDE string
	if update.SystemPromptDE != nil {
		promptDE = strings.TrimSpace(*update.SystemPromptDE)
		if promptDE == "" {
			return fmt.Errorf("system_prompt_de cannot be empty")
		}
		if len(promptDE) > maxPromptChars {
			return fmt.Errorf("system_prompt_de exceeds %d characters", maxPromptChars)
		}
	}

	var promptKZ string
	if update.SystemPromptKZ != nil {
		promptKZ = strings.TrimSpace(*update.SystemPromptKZ)
		if promptKZ == "" {
			return fmt.Errorf("system_prompt_kz cannot be empty")
		}
		if len(promptKZ) > maxPromptChars {
			return fmt.Errorf("system_prompt_kz exceeds %d characters", maxPromptChars)
		}
	}

	runtimeConfig.mu.Lock()

	if update.Model != nil {
		runtimeConfig.cfg.Model = modelValue
	}

	if update.SystemPromptEN != nil {
		runtimeConfig.cfg.SystemPromptEN = promptEN
	}

	if update.SystemPromptDE != nil {
		runtimeConfig.cfg.SystemPromptDE = promptDE
	}

	if update.SystemPromptKZ != nil {
		runtimeConfig.cfg.SystemPromptKZ = promptKZ
	}

	if update.EnableParallelGeneration != nil {
		runtimeConfig.cfg.EnableParallelGeneration = *update.EnableParallelGeneration
	}

	updated := runtimeConfig.cfg
	runtimeConfig.mu.Unlock()

	if err := persistRuntimeConfig(updated); err != nil {
		return err
	}
	return nil
}

func persistRuntimeConfig(cfg RuntimeConfig) error {
	if runtimeConfigStore == nil {
		return nil
	}

	if storeWithTTL, ok := runtimeConfigStore.(runtimeConfigTTLStore); ok {
		if err := storeWithTTL.SetWithTTL(runtimeConfigCacheKey, cfg, 0); err != nil {
			return fmt.Errorf("failed to persist runtime config: %w", err)
		}
		return nil
	}

	if err := runtimeConfigStore.Set(runtimeConfigCacheKey, cfg); err != nil {
		return fmt.Errorf("failed to persist runtime config: %w", err)
	}
	return nil
}
