package controller

import (
	"image-service/internal/model"
	"mime/multipart"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
)

// IMAGE SERVICE

type MockImageService struct {
	mock.Mock
}

func (m *MockImageService) UploadImage(albumID string, files []*multipart.FileHeader) ([]*model.Image, error) {
	args := m.Called(albumID, files)
	return args.Get(0).([]*model.Image), args.Error(1)
}

func (m *MockImageService) DeleteImage(albumID, imageID string) error {
	args := m.Called(albumID, imageID)
	return args.Error(0)
}

func (m *MockImageService) ServeImage(albumID, imageID string) (string, error) {
	args := m.Called(albumID, imageID)
	return args.String(0), args.Error(1)
}

func (m *MockImageService) RenameImage(albumID, imageID, newName string) (*model.Image, error) {
	args := m.Called(albumID, imageID, newName)
	if img, ok := args.Get(0).(*model.Image); ok {
		return img, args.Error(1)
	}
	return nil, args.Error(1)
}

// PRDUCER

type MockProducer struct {
	mock.Mock
}

func (m *MockProducer) SendMessage(message []byte) error {
	args := m.Called(message)
	return args.Error(0)
}

func (m *MockProducer) Close() {
	m.Called()
}

// CACHE SERVICE

type MockCacheService struct {
	mock.Mock
}

func (m *MockCacheService) ClearCache() error {
	args := m.Called()
	return args.Error(0)
}

// ALBUM SERVICE

type MockAlbumService struct {
	mock.Mock
}

func (m *MockAlbumService) GetAlbum(c *gin.Context, id string) (*model.Album, error) {
	args := m.Called(c, id)
	if album, ok := args.Get(0).(*model.Album); ok {
		return album, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockAlbumService) GetAlbumsPreview(typeQuery string) ([]*model.AlbumPreview, error) {
	args := m.Called(typeQuery)
	if albums, ok := args.Get(0).([]*model.AlbumPreview); ok {
		return albums, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockAlbumService) CreateAlbum(album *model.AlbumPreview) (*model.AlbumPreview, error) {
	args := m.Called(album)
	if created, ok := args.Get(0).(*model.AlbumPreview); ok {
		return created, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockAlbumService) DeleteAlbum(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockAlbumService) UpdateAlbum(id string, album *model.AlbumPreview) (*model.AlbumPreview, error) {
	args := m.Called(id, album)
	if updated, ok := args.Get(0).(*model.AlbumPreview); ok {
		return updated, args.Error(1)
	}
	return nil, args.Error(1)
}
