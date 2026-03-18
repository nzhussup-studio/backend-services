package repository

import (
	"image-service/internal/model"
	"sync"
)

type Storage struct {
	Path        string
	ApiBasePath string
	Album       interface {
		Create(*model.AlbumPreview) (*model.AlbumPreview, error)
		GetPreview() ([]*model.AlbumPreview, error)
		Get(string) (*model.Album, error)
		Delete(string) error
		Update(string, *model.AlbumPreview) (*model.AlbumPreview, error)
	}
	Image interface {
		Upload(string, *model.Image) (*model.Image, error)
		Rename(string, string, string) (*model.Image, error)
		Delete(string, string) error
	}
}

var albumLocks sync.Map

var getAlbumLock = func(id string) *sync.Mutex {
	lock, _ := albumLocks.LoadOrStore(id, &sync.Mutex{})
	return lock.(*sync.Mutex)
}

func NewStorage(path string, apiBasePath string) *Storage {
	return &Storage{
		Path:        path,
		ApiBasePath: apiBasePath,
		Album:       &AlbumRepository{Path: path, ApiBasePath: apiBasePath},
		Image:       &ImageRepository{Path: path, ApiBasePath: apiBasePath},
	}
}
