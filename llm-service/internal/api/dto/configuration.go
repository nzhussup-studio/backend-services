package dto

type ConfigurationRequest struct {
	Model                    *string `json:"model,omitempty"`
	SystemPromptEN           *string `json:"system_prompt_en,omitempty"`
	SystemPromptDE           *string `json:"system_prompt_de,omitempty"`
	SystemPromptKZ           *string `json:"system_prompt_kz,omitempty"`
	EnableParallelGeneration *bool   `json:"enable_parallel_generation,omitempty"`
}

type ConfigurationResponse struct {
	Status                   int    `json:"status"`
	Model                    string `json:"model"`
	SystemPromptEN           string `json:"system_prompt_en"`
	SystemPromptDE           string `json:"system_prompt_de"`
	SystemPromptKZ           string `json:"system_prompt_kz"`
	EnableParallelGeneration bool   `json:"enable_parallel_generation"`
}
