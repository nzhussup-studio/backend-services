package controller

import (
	"net/http"

	"account-service/internal/model"

	"github.com/gin-gonic/gin"
)

// HealthCheckHandler godoc
// @Summary Health check endpoint
// @Description Returns 200 OK if the service is up
// @Tags Health
// @Produce json
// @Success 200 {object} model.SuccessResponse "Service is healthy"
// @Router /health [get]
func HealthCheckHandler(c *gin.Context) {
	c.JSON(http.StatusOK, model.SuccessResponse{Message: "OK"})
}
