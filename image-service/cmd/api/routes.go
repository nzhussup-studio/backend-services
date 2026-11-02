package main

import (
	"image-service/internal/auth"
	"image-service/internal/controller"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func (a *app) GetRouter() *gin.Engine {
	r := gin.Default()
	r.MaxMultipartMemory = 8 << 20
	r.GET("/healthz", controller.HealthCheckHandler)

	v1 := r.Group(a.config.API.BasePath)
	v1.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	v1.GET("", a.Controller.AlbumController.GetPublicPreview)
	v1.GET("/:id", a.Controller.AlbumController.GetPublic)
	v1.GET("/:id/images/:imageID", a.Controller.ImageController.ServePublic)

	admin := r.Group("/v1/admin/albums")
	admin.Use(auth.RequireTrustedIdentityMiddleware())
	admin.DELETE("/cache", a.Controller.CacheController.ClearCache)
	admin.GET("", a.Controller.AlbumController.GetPreview)
	admin.GET("/:id", a.Controller.AlbumController.Get)
	admin.POST("", a.Controller.AlbumController.Create)
	admin.PUT("/:id", a.Controller.AlbumController.Update)
	admin.DELETE("/:id", a.Controller.AlbumController.Delete)
	admin.POST("/:id/images", a.Controller.ImageController.Upload)
	admin.GET("/:id/uploads/:jobID", a.Controller.ImageController.UploadStatus)
	admin.DELETE("/:id/images/:imageID", a.Controller.ImageController.Delete)
	admin.PATCH("/:id/images/:imageID", a.Controller.ImageController.Rename)
	admin.GET("/:id/images/:imageID", a.Controller.ImageController.Serve)

	return r
}
