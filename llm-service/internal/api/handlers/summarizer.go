package handlers

import (
	"llm-service/configs"
	"llm-service/internal/api/dto"
	"llm-service/internal/summarizer"
	"llm-service/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Cache interface {
	Set(key string, value any) error
	Get(key string, dest any) error
	Del(key string) error
	Ping() error
}

type Handler struct {
	config *configs.Config
	redis  Cache
}

func New(config *configs.Config, redis Cache) *Handler {
	return &Handler{
		config: config,
		redis:  redis,
	}
}

// handleGetSummarizer godoc
// @Summary      Generate professional profile summary
// @Description  Retrieves structured personal data (e.g., work experience, education), generates a professional summary using a large language model (LLM), and returns it in the requested language.
// @Tags         summarizer
// @Produce      json
// @Param        lang  query     string  false  "Language code for the summary output. Supported values: 'en' (English), 'kz' (Kazakh), 'de' (German). Defaults to 'en'."
// @Success      200   {object}  dto.APIResponse  "Generated professional summary"
// @Failure      400   {object}  dto.APIResponse  "Invalid query parameters"
// @Failure      500   {object}  dto.APIResponse  "Internal server error with error details"
// @Router       /v1/llm/summarize [get]
func (h *Handler) GetSummarizer(ctx *gin.Context) {
	var query dto.SummarizeQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.APIResponse{
			Status:  http.StatusBadRequest,
			Message: "invalid query parameters",
			Error:   "Bad Request",
		})
		return
	}

	lang := query.Lang
	if lang == "" {
		lang = "en"
	}
	if !utils.IsValidLanguage(lang) {
		lang = "en" // Default to English if the provided language is invalid
	}

	svc := summarizer.NewSummarizer(
		h.config.Summarizer.APIKey,
		h.config.Summarizer.APIURL,
		[]string{
			h.config.Services.WorkExperienceURL,
			h.config.Services.EducationURL,
			h.config.Services.ProjectsURL,
			h.config.Services.SkillsURL,
			h.config.Services.CertificatesURL,
		},
		http.DefaultClient,
		h.redis,
		lang,
	)

	summary, err := svc.Summarize(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.APIResponse{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
			Error:   "Internal Server Error",
		})
		return
	}

	if summary == "" {
		ctx.JSON(http.StatusInternalServerError, dto.APIResponse{
			Status:  http.StatusInternalServerError,
			Message: summarizer.ErrNoSummaryReceived.Error(),
			Error:   "Internal Server Error",
		})
		return
	}

	ctx.JSON(http.StatusOK, dto.APIResponse{
		Status:  http.StatusOK,
		Message: summary,
	})
}
