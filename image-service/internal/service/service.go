package service

import (
	"image-service/internal/auth"
	"image-service/internal/cache"
	appconfig "image-service/internal/config"
	"image-service/internal/model"
	"image-service/internal/repository"
	"mime/multipart"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type Service struct {
	storage      *repository.Storage
	validate     *validator.Validate
	imageConfig  appconfig.ImageConfig
	AlbumService interface {
		GetAlbumsPreview(string) ([]*model.AlbumPreview, error)
		GetAlbum(*gin.Context, string) (*model.Album, error)
		CreateAlbum(*model.AlbumPreview) (*model.AlbumPreview, error)
		UpdateAlbum(string, *model.AlbumPreview) (*model.AlbumPreview, error)
		DeleteAlbum(string) error
	}
	ImageService interface {
		UploadImage(string, []*multipart.FileHeader) ([]*model.Image, error)
		RenameImage(string, string, string) (*model.Image, error)
		DeleteImage(string, string) error
		ServeImage(string, string) (string, error)
	}
	CacheService interface {
		ClearCache() error
	}
}

type ImageConfig = appconfig.ImageConfig

func NewService(storage *repository.Storage, redis *cache.RedisClient, securityConfig *auth.AuthConfig, validate *validator.Validate, imageCfg ImageConfig) *Service {
	return &Service{
		storage:      storage,
		validate:     validate,
		imageConfig:  imageCfg,
		AlbumService: &AlbumService{storage: storage, redis: redis, securityConfig: securityConfig, validate: validate},
		ImageService: &ImageService{storage: storage, redis: redis, validate: validate, cfg: imageCfg},
		CacheService: &CacheService{redis: redis},
	}
}
