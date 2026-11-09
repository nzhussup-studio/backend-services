package dto

type GenerationRequest struct {
	Language string `json:"language" binding:"required" example:"en"`
}

type SummaryResponse struct {
	Status      string `json:"status" example:"generating"`
	Language    string `json:"language" example:"en"`
	Summary     string `json:"summary,omitempty"`
	PollAfterMS int    `json:"poll_after_ms,omitempty" example:"1500"`
}
