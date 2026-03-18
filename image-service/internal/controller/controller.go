package controller

import (
	"image-service/internal/config/messaging"
	"image-service/internal/service"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	service         *service.Service
	AlbumController interface {
		Get(*gin.Context)
		GetPreview(*gin.Context)
		Create(*gin.Context)
		Update(*gin.Context)
		Delete(*gin.Context)
	}
	ImageController interface {
		Upload(*gin.Context)
		Delete(*gin.Context)
		Serve(*gin.Context)
		Rename(*gin.Context)
	}
	CacheController interface {
		ClearCache(*gin.Context)
	}
}

func NewController(service *service.Service, producer *messaging.Producer) *Controller {
	return &Controller{
		service:         service,
		AlbumController: &AlbumController{service: service},
		ImageController: &ImageController{service: service, producer: producer},
		CacheController: &CacheController{service: service},
	}
}
