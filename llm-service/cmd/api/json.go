package main

import (
	"github.com/gin-gonic/gin"
)

type jsonResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type errorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"` // omit if empty
}

func constructJSONResponse(ctx *gin.Context, status int, message string) {
	resp := jsonResponse{
		Status:  status,
		Message: message,
	}
	ctx.JSON(status, resp)
}

func constructErrorResponse(ctx *gin.Context, status int, errorMessage string) {
	resp := errorResponse{
		Status:  status,
		Message: errorMessage,
	}

	if status >= 500 {
		resp.Error = "Internal Server Error"
	}

	ctx.JSON(status, resp)
}
