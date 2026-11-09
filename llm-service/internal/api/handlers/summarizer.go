package handlers

import (
	"llm-service/configs"
	"llm-service/internal/api/dto"
	"llm-service/internal/cache"
	"llm-service/internal/summarizer"
	"llm-service/internal/utils"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Cache interface {
	cache.Cacher
	Ping() error
}

type Handler struct {
	config *configs.Config
	redis  Cache
}

func New(config *configs.Config, redis Cache) *Handler {
	if err := summarizer.InitRuntimeConfigStore(redis); err != nil {
		slog.Error("failed to initialize runtime config from redis", slog.String("error", err.Error()))
	}
	return &Handler{config: config, redis: redis}
}

// StartSummaryGeneration godoc
// @Summary      Start or join summary generation
// @Description  Returns a current cached summary or starts missing and stale language summaries in parallel.
// @Tags         summarizer
// @Accept       json
// @Produce      json
// @Param        request  body  dto.GenerationRequest  true  "Requested language"
// @Success      200  {object}  dto.SummaryResponse  "Current summary"
// @Success      202  {object}  dto.SummaryResponse  "Generation in progress"
// @Failure      400  {object}  dto.APIResponse      "Invalid request"
// @Failure      500  {object}  dto.APIResponse      "Internal server error"
// @Router       /v1/llm/summaries/generations [post]
func (h *Handler) StartSummaryGeneration(ctx *gin.Context) {
	var request dto.GenerationRequest
	if err := ctx.ShouldBindJSON(&request); err != nil || !utils.IsValidLanguage(request.Language) {
		ctx.JSON(http.StatusBadRequest, dto.APIResponse{Status: http.StatusBadRequest, Message: "language must be one of en, kk, or de", Error: "Bad Request"})
		return
	}
	result, err := h.summarizer().Start(ctx.Request.Context(), request.Language)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.APIResponse{Status: http.StatusInternalServerError, Message: err.Error(), Error: "Internal Server Error"})
		return
	}
	writeSummaryResult(ctx, result)
}

// GetSummary godoc
// @Summary      Get summary or generation status
// @Tags         summarizer
// @Produce      json
// @Param        language  path  string  true  "Language" Enums(en,kk,de)
// @Success      200  {object}  dto.SummaryResponse  "Generated summary"
// @Success      202  {object}  dto.SummaryResponse  "Generation in progress"
// @Failure      400  {object}  dto.APIResponse      "Invalid language"
// @Failure      404  {object}  dto.SummaryResponse  "No summary or active generation"
// @Failure      500  {object}  dto.APIResponse      "Internal server error"
// @Failure      503  {object}  dto.SummaryResponse  "Generation failed"
// @Router       /v1/llm/summaries/{language} [get]
func (h *Handler) GetSummary(ctx *gin.Context) {
	language := ctx.Param("language")
	if !utils.IsValidLanguage(language) {
		ctx.JSON(http.StatusBadRequest, dto.APIResponse{Status: http.StatusBadRequest, Message: "language must be one of en, kk, or de", Error: "Bad Request"})
		return
	}
	result, err := h.summarizer().Get(language)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.APIResponse{Status: http.StatusInternalServerError, Message: err.Error(), Error: "Internal Server Error"})
		return
	}
	writeSummaryResult(ctx, result)
}

func (h *Handler) summarizer() *summarizer.Summarizer {
	return summarizer.NewSummarizer(h.config.Summarizer.APIKey, h.config.Summarizer.APIURL, []string{
		h.config.Services.WorkExperienceURL,
		h.config.Services.EducationURL,
		h.config.Services.ProjectsURL,
		h.config.Services.SkillsURL,
		h.config.Services.CertificatesURL,
	}, nil, h.redis)
}

func writeSummaryResult(ctx *gin.Context, result summarizer.Result) {
	response := dto.SummaryResponse{Status: result.Status, Language: result.Language, Summary: result.Summary, PollAfterMS: result.PollAfterMS}
	switch result.Status {
	case "completed":
		ctx.JSON(http.StatusOK, response)
	case "generating":
		ctx.JSON(http.StatusAccepted, response)
	case "missing":
		ctx.JSON(http.StatusNotFound, response)
	default:
		ctx.JSON(http.StatusServiceUnavailable, response)
	}
}
