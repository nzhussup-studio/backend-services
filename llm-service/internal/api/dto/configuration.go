package dto

type ConfigurationRequest struct {
	Model                    *string `json:"model,omitempty"`
	SystemPromptEN           *string `json:"system_prompt_en,omitempty"`
	SystemPromptDE           *string `json:"system_prompt_de,omitempty"`
	SystemPromptKK           *string `json:"system_prompt_kk,omitempty"`
	EnableParallelGeneration *bool   `json:"enable_parallel_generation,omitempty"`
}

type ConfigurationResponse struct {
	Status                   int    `json:"status"`
	Model                    string `json:"model"`
	SystemPromptEN           string `json:"system_prompt_en"`
	SystemPromptDE           string `json:"system_prompt_de"`
	SystemPromptKK           string `json:"system_prompt_kk"`
	EnableParallelGeneration bool   `json:"enable_parallel_generation"`
}
