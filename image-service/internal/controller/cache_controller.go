package controller

import (
	custom_errors "image-service/internal/errors"
	"image-service/internal/json"
	"image-service/internal/service"

	"github.com/gin-gonic/gin"
)

type CacheController struct {
	service *service.Service
}

func NewCacheController(service *service.Service) *CacheController {
	return &CacheController{service: service}
}

// ClearCache godoc
// @Summary Clear the image cache
// @Description This endpoint clears the server-side cache for images.
// @Tags Cache
// @Accept  json
// @Produce  json
// @Success 200 {object} model.SuccessResponse "Cache cleared successfully"
// @Failure 500 {object} model.ErrorResponse "Internal Server Error"
// @Router /v1/album/cache [delete]
// @Security ApiKeyAuth
func (ctrl *CacheController) ClearCache(c *gin.Context) {

	err := ctrl.service.CacheService.ClearCache()
	if err != nil {
		custom_errors.MapErrors(c, err)
		return
	}

	json.ConstructJsonResponseSuccess(c, nil, "Cache cleared successfully")
}
