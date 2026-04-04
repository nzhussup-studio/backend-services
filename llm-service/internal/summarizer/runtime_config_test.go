package summarizer

import (
	"strings"
	"testing"
)

func strPtr(v string) *string { return &v }
func boolPtr(v bool) *bool    { return &v }

func restoreRuntimeConfig(t *testing.T, cfg RuntimeConfig) {
	t.Helper()
	err := UpdateRuntimeConfig(RuntimeConfigUpdate{
		Model:                    strPtr(cfg.Model),
		SystemPromptEN:           strPtr(cfg.SystemPromptEN),
		SystemPromptDE:           strPtr(cfg.SystemPromptDE),
		SystemPromptKZ:           strPtr(cfg.SystemPromptKZ),
		EnableParallelGeneration: boolPtr(cfg.EnableParallelGeneration),
	})
	if err != nil {
		t.Fatalf("failed to restore runtime config: %v", err)
	}
}

func TestRuntimeConfig_UpdateSuccess(t *testing.T) {
	original := GetRuntimeConfig()
	defer restoreRuntimeConfig(t, original)

	model := "openai/gpt-4.1-mini"
	promptEN := "EN prompt override"
	promptDE := "DE prompt override"
	promptKZ := "KZ prompt override"
	parallel := false

	err := UpdateRuntimeConfig(RuntimeConfigUpdate{
		Model:                    &model,
		SystemPromptEN:           &promptEN,
		SystemPromptDE:           &promptDE,
		SystemPromptKZ:           &promptKZ,
		EnableParallelGeneration: &parallel,
	})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	cfg := GetRuntimeConfig()
	if cfg.Model != model {
		t.Fatalf("expected model %q, got %q", model, cfg.Model)
	}
	if cfg.SystemPromptEN != promptEN || cfg.SystemPromptDE != promptDE || cfg.SystemPromptKZ != promptKZ {
		t.Fatalf("prompt overrides were not applied")
	}
	if cfg.EnableParallelGeneration != parallel {
		t.Fatalf("expected parallel=%v, got %v", parallel, cfg.EnableParallelGeneration)
	}
}

func TestRuntimeConfig_UpdateValidation(t *testing.T) {
	original := GetRuntimeConfig()
	defer restoreRuntimeConfig(t, original)

	tests := []struct {
		name   string
		update RuntimeConfigUpdate
	}{
		{
			name:   "empty update",
			update: RuntimeConfigUpdate{},
		},
		{
			name: "empty model",
			update: RuntimeConfigUpdate{
				Model: strPtr("   "),
			},
		},
		{
			name: "empty english prompt",
			update: RuntimeConfigUpdate{
				SystemPromptEN: strPtr(""),
			},
		},
		{
			name: "too long english prompt",
			update: RuntimeConfigUpdate{
				SystemPromptEN: strPtr(strings.Repeat("a", maxPromptChars+1)),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := UpdateRuntimeConfig(tt.update); err == nil {
				t.Fatalf("expected validation error, got nil")
			}
		})
	}
}
