package service

import (
	"image-service/internal/model"

	"github.com/stretchr/testify/mock"
)

// ALBUM REPOSITORY

type MockAlbumRepo struct {
	mock.Mock
}

func (m *MockAlbumRepo) Create(album *model.AlbumPreview) (*model.AlbumPreview, error) {
	return m.Called(album).Get(0).(*model.AlbumPreview), m.Called(album).Error(1)
}

func (m *MockAlbumRepo) Get(id string) (*model.Album, error) {
	return m.Called(id).Get(0).(*model.Album), m.Called(id).Error(1)
}

func (m *MockAlbumRepo) GetPreview() ([]*model.AlbumPreview, error) {
	return m.Called().Get(0).([]*model.AlbumPreview), m.Called().Error(1)
}

func (m *MockAlbumRepo) Delete(id string) error {
	return m.Called(id).Error(0)
}

func (m *MockAlbumRepo) Update(id string, album *model.AlbumPreview) (*model.AlbumPreview, error) {
	return m.Called(id, album).Get(0).(*model.AlbumPreview), m.Called(id, album).Error(1)
}

// IMAGE REPOSITORY

type MockImageRepo struct {
	mock.Mock
}

func (m *MockImageRepo) Upload(albumID string, image *model.Image) (*model.Image, error) {
	args := m.Called(albumID, image)
	var img *model.Image
	if res := args.Get(0); res != nil {
		img = res.(*model.Image)
	}
	return img, args.Error(1)
}

func (m *MockImageRepo) Delete(albumID string, imageID string) error {
	args := m.Called(albumID, imageID)
	return args.Error(0)
}

func (m *MockImageRepo) Rename(albumID string, imageID string, newName string) (*model.Image, error) {
	args := m.Called(albumID, imageID, newName)
	var img *model.Image
	if res := args.Get(0); res != nil {
		img = res.(*model.Image)
	}
	return img, args.Error(1)
}

// REDIS CLIENT

type MockRedisClient struct {
	mock.Mock
}

func (m *MockRedisClient) Set(key string, value any) {
	m.Called(key, value)
}

func (m *MockRedisClient) Get(key string, dest any) error {
	args := m.Called(key, dest)
	return args.Error(0)
}

func (m *MockRedisClient) Del(key string) {
	m.Called(key)
}

func (m *MockRedisClient) FlushAll() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockRedisClient) CheckHealth() {
	m.Called()
}
