package model

type ErrorResponse struct {
	Error string `json:"error" example:"Unauthorized"`
}

type SuccessResponse struct {
	Message string `json:"message,omitempty" example:"OK"`
}
