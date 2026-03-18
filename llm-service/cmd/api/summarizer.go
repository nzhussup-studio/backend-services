package main

import (
	"llm-service/internal/summarizer"
	"llm-service/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// handleGetSummarizer godoc
// @Summary      Generate professional profile summary
// @Description  Retrieves structured personal data (e.g., work experience, education), generates a professional summary using a large language model (LLM), and returns it in the requested language.
// @Tags         summarizer
// @Produce      json
// @Param        lang  query     string  false  "Language code for the summary output. Supported values: 'en' (English), 'kz' (Kazakh), 'de' (German). Defaults to 'en'."
// @Success      200   {string}  string  "Generated professional summary"
// @Failure      500   {object}  map[string]string  "Internal server error with error details"
// @Router       /v1/llm/summarize [get]
func (app *app) handleGetSummarizer(ctx *gin.Context) {
	lang := ctx.DefaultQuery("lang", "en")
	if !utils.IsValidLanguage(lang) {
		lang = "en" // Default to English if the provided language is invalid
	}

	s := summarizer.NewSummarizer(
		app.config.summarizerConfig.API_KEY,
		app.config.summarizerConfig.API_URL,
		[]string{
			app.config.servicesConfig.WORK_EXPERIENCE_URL,
			app.config.servicesConfig.EDUCATION_URL,
			app.config.servicesConfig.PROJECTS_URL,
			app.config.servicesConfig.SKILLS_URL,
			app.config.servicesConfig.CERTIFICATES_URL,
		},
		http.DefaultClient,
		app.redis,
		lang,
	)

	summary, err := s.Summarize(ctx.Request.Context())
	if err != nil {
		constructErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	if summary == "" {
		constructErrorResponse(ctx, http.StatusInternalServerError, summarizer.ErrNoSummaryReceived.Error())
		return
	}

	constructJSONResponse(ctx, http.StatusOK, summary)
}
