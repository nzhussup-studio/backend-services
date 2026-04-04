package summarizer

import (
	"strings"
	"testing"

	"llm-service/internal/model"

	"github.com/stretchr/testify/require"
)

func TestPrompts_getPromptBase(t *testing.T) {
	original := GetRuntimeConfig()
	defer restoreRuntimeConfig(t, original)

	customEN := "Custom EN Prompt"
	customDE := "Custom DE Prompt"
	customKZ := "Custom KZ Prompt"

	require.NoError(t, UpdateRuntimeConfig(RuntimeConfigUpdate{
		SystemPromptEN: &customEN,
		SystemPromptDE: &customDE,
		SystemPromptKZ: &customKZ,
	}))

	pd := &model.PersonalData{}

	t.Run("selects german prompt for de", func(t *testing.T) {
		s := &Summarizer{lang: "de"}
		systemPrompt, userPrompt, err := s.getPromptBase(pd)

		require.NoError(t, err)
		require.Equal(t, customDE, systemPrompt)
		require.Contains(t, userPrompt, "structured profile data")
	})

	t.Run("falls back to english for unsupported language", func(t *testing.T) {
		s := &Summarizer{lang: "fr"}
		systemPrompt, userPrompt, err := s.getPromptBase(pd)

		require.NoError(t, err)
		require.Equal(t, customEN, systemPrompt)
		require.Contains(t, userPrompt, `"work_experience": null`)
	})

	t.Run("includes personal data json in user prompt", func(t *testing.T) {
		company := "OpenAI"
		s := &Summarizer{lang: "en"}
		systemPrompt, userPrompt, err := s.getPromptBase(&model.PersonalData{
			WorkExperience: []*model.WorkExperience{
				{Company: company},
			},
		})

		require.NoError(t, err)
		require.Equal(t, customEN, systemPrompt)
		require.True(t, strings.Contains(userPrompt, company))
	})
}
