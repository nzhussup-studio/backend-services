package service

import (
	"errors"
	"image-service/internal/config/security"
	"image-service/internal/model"
	"image-service/internal/repository"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateAlbum_Success(t *testing.T) {
	mockAlbumRepo := new(MockAlbumRepo)
	mockRedis := new(MockRedisClient)

	storage := &repository.Storage{Album: mockAlbumRepo}
	validate := validator.New()
	validate.RegisterValidation("albumtype", model.ValidateAlbumType)
	service := &AlbumService{
		storage:  storage,
		redis:    mockRedis,
		validate: validate,
	}

	album := &model.AlbumPreview{Title: "Test Album", Desc: "Desc", Type: model.Private}

	mockAlbumRepo.On("Create", album).Return(album, nil)

	created, err := service.CreateAlbum(album)

	assert.NoError(t, err)
	assert.Equal(t, album, created)
	mockAlbumRepo.AssertExpectations(t)
}

func TestGetAlbum_FromCache(t *testing.T) {
	mockAlbumRepo := new(MockAlbumRepo)
	mockRedis := new(MockRedisClient)

	storage := &repository.Storage{Album: mockAlbumRepo}
	validate := validator.New()
	validate.RegisterValidation("albumtype", model.ValidateAlbumType)
	service := &AlbumService{
		storage:  storage,
		redis:    mockRedis,
		validate: validate,
	}

	albumID := "album1"
	album := &model.Album{ID: albumID, Type: model.Public}

	mockRedis.On("Get", "album_"+albumID, mock.Anything).Run(func(args mock.Arguments) {
		dest := args.Get(1).(**model.Album)
		*dest = album
	}).Return(nil)

	result, err := service.GetAlbum(&gin.Context{}, albumID)

	assert.NoError(t, err)
	assert.Equal(t, album, result)
	mockRedis.AssertExpectations(t)
}

func TestDeleteAlbum_Success(t *testing.T) {
	mockAlbumRepo := new(MockAlbumRepo)
	mockRedis := new(MockRedisClient)

	storage := &repository.Storage{Album: mockAlbumRepo}
	service := &AlbumService{
		storage: storage,
		redis:   mockRedis,
	}

	albumID := "album1"

	mockRedis.On("Del", "album_"+albumID).Return(nil)
	mockAlbumRepo.On("Delete", albumID).Return(nil)

	err := service.DeleteAlbum(albumID)

	assert.NoError(t, err)
	mockRedis.AssertExpectations(t)
	mockAlbumRepo.AssertExpectations(t)
}

func setupAlbumServiceWithMock(album *model.Album) (*AlbumService, *MockAlbumRepo, *MockRedisClient) {
	mockAlbumRepo := new(MockAlbumRepo)
	mockRedis := new(MockRedisClient)

	storage := &repository.Storage{Album: mockAlbumRepo}
	validate := validator.New()
	validate.RegisterValidation("albumtype", model.ValidateAlbumType)

	service := &AlbumService{
		storage:  storage,
		redis:    mockRedis,
		validate: validate,
	}

	mockRedis.On("Get", "album_"+album.ID, mock.Anything).Run(func(args mock.Arguments) {
		dest := args.Get(1).(**model.Album)
		*dest = album
	}).Return(nil)

	return service, mockAlbumRepo, mockRedis
}

func TestGetAlbum_Private_Admin(t *testing.T) {
	album := &model.Album{ID: "album1", Type: model.Private}
	service, _, mockRedis := setupAlbumServiceWithMock(album)

	original := security.CheckIsAdmin
	security.CheckIsAdmin = func(ctx *gin.Context, config *security.AuthConfig) error {
		return nil
	}
	defer func() { security.CheckIsAdmin = original }()

	result, err := service.GetAlbum(&gin.Context{}, album.ID)

	assert.NoError(t, err)
	assert.Equal(t, album, result)
	mockRedis.AssertExpectations(t)
}

func TestGetAlbum_Private_NotAdmin(t *testing.T) {
	album := &model.Album{ID: "album2", Type: model.Private}
	service, _, mockRedis := setupAlbumServiceWithMock(album)

	original := security.CheckIsAdmin
	security.CheckIsAdmin = func(ctx *gin.Context, config *security.AuthConfig) error {
		return errors.New("not an admin")
	}
	defer func() { security.CheckIsAdmin = original }()

	result, err := service.GetAlbum(&gin.Context{}, album.ID)

	assert.Error(t, err)
	assert.Nil(t, result)
	mockRedis.AssertExpectations(t)
}
