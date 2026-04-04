package handlers

import (
	"llm-service/internal/api/dto"
	"net/http"

	"github.com/gin-gonic/gin"
)

// handleGetHealth godoc
// @Summary      Health check endpoint
// @Description  Checks the connectivity and health of dependent services, particularly Redis.
// @Tags         health
// @Produce      json
// @Success      200  {object}  dto.APIResponse  "Status OK"
// @Failure      500  {object}  dto.APIResponse  "Redis connection failed"
// @Router       /v1/llm/health [get]
func (h *Handler) GetHealth(ctx *gin.Context) {

	if err := h.redis.Ping(); err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.APIResponse{
			Status:  http.StatusInternalServerError,
			Message: "Redis connection failed",
			Error:   "Internal Server Error",
		})
		return
	}

	ctx.JSON(http.StatusOK, dto.APIResponse{
		Status:  http.StatusOK,
		Message: "ok",
	})
}
