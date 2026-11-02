package main

import (
	"image-service/internal/cache"
	appconfig "image-service/internal/config"
	"image-service/internal/controller"
	"image-service/internal/model"
	"image-service/internal/repository"
	"image-service/internal/service"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
	validate.RegisterValidation("albumtype", model.ValidateAlbumType)
}

type app struct {
	config     *appconfig.Config
	Controller *controller.Controller
	Service    *service.Service
	Storage    *repository.Storage
	Redis      *cache.RedisClient
}

func newApp(config *appconfig.Config) *app {
	redisClient := cache.NewRedisClient(
		config.Redis.Addr,
		config.Redis.Password,
		config.Redis.DB,
		config.Redis.TTL,
	)
	storage := repository.NewStorage(config.Storage.Path, config.API.BasePath)
	service := service.NewService(storage, redisClient, validate, config.Image)
	controller := controller.NewController(service)

	return &app{
		config:     config,
		Controller: controller,
		Service:    service,
		Storage:    storage,
		Redis:      redisClient,
	}
}

func (a *app) Run() {
	router := a.GetRouter()
	router.Run(a.config.Server.Addr)
}
