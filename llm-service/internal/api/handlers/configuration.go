package handlers

import (
	"llm-service/internal/api/dto"
	"llm-service/internal/summarizer"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetConfiguration godoc
// @Summary      Get current LLM configuration
// @Description  Returns the active model and language prompts.
// @Tags         configuration
// @Produce      json
// @Success      200  {object}  dto.ConfigurationResponse  "Current configuration"
// @Router       /v1/admin/llm/configuration [get]
func (h *Handler) GetConfiguration(ctx *gin.Context) {
	cfg := summarizer.GetRuntimeConfig()
	ctx.JSON(http.StatusOK, dto.ConfigurationResponse{
		Status:         http.StatusOK,
		Model:          cfg.Model,
		SystemPromptEN: cfg.SystemPromptEN,
		SystemPromptDE: cfg.SystemPromptDE,
		SystemPromptKK: cfg.SystemPromptKK,
	})
}

// PutConfiguration godoc
// @Summary      Update LLM configuration
// @Description  Updates the active model and language prompts.
// @Tags         configuration
// @Accept       json
// @Produce      json
// @Param        request  body      dto.ConfigurationRequest  true  "Configuration overrides"
// @Success      200      {object}  dto.ConfigurationResponse "Updated configuration"
// @Failure      400      {object}  dto.APIResponse           "Invalid configuration payload"
// @Router       /v1/admin/llm/configuration [put]
func (h *Handler) PutConfiguration(ctx *gin.Context) {
	var req dto.ConfigurationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.APIResponse{
			Status:  http.StatusBadRequest,
			Message: "invalid request payload",
			Error:   "Bad Request",
		})
		return
	}

	update := summarizer.RuntimeConfigUpdate{
		Model:          req.Model,
		SystemPromptEN: req.SystemPromptEN,
		SystemPromptDE: req.SystemPromptDE,
		SystemPromptKK: req.SystemPromptKK,
	}

	if err := summarizer.UpdateRuntimeConfig(update); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.APIResponse{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
			Error:   "Bad Request",
		})
		return
	}

	cfg := summarizer.GetRuntimeConfig()
	ctx.JSON(http.StatusOK, dto.ConfigurationResponse{
		Status:         http.StatusOK,
		Model:          cfg.Model,
		SystemPromptEN: cfg.SystemPromptEN,
		SystemPromptDE: cfg.SystemPromptDE,
		SystemPromptKK: cfg.SystemPromptKK,
	})
}
