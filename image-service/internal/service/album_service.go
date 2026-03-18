package service

import (
	"fmt"
	"image-service/internal/config/cache"
	"image-service/internal/config/security"
	custom_errors "image-service/internal/errors"
	"image-service/internal/model"
	"image-service/internal/repository"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type AlbumService struct {
	storage        *repository.Storage
	redis          cache.RedisClientInterface
	securityConfig *security.AuthConfig
	validate       *validator.Validate
}

func (s *AlbumService) GetAlbum(c *gin.Context, id string) (*model.Album, error) {
	var cachedAlbum *model.Album
	err := s.redis.Get(fmt.Sprintf("album_%s", id), &cachedAlbum)
	if err == nil {
		if err = allowedToReturn(c, cachedAlbum, s.securityConfig); err != nil {
			return nil, err
		}
		return cachedAlbum, nil
	}

	album, err := s.storage.Album.Get(id)
	if err != nil {
		return nil, err
	}
	if err = allowedToReturn(c, album, s.securityConfig); err != nil {
		return nil, err
	}

	s.redis.Set(fmt.Sprintf("album_%s", id), album)
	return album, nil
}

func (s *AlbumService) GetAlbumsPreview(typeQuery string) ([]*model.AlbumPreview, error) {

	albums, err := s.storage.Album.GetPreview()
	if err != nil {
		return nil, err
	}

	switch typeQuery {
	case string(model.Private):
		albums = filterAlbums(albums, func(a *model.AlbumPreview) bool {
			return a.Type == model.Private
		})
	case string(model.SemiPublic):
		albums = filterAlbums(albums, func(a *model.AlbumPreview) bool {
			return a.Type == model.SemiPublic
		})
	case string(model.Public):
		albums = filterAlbums(albums, func(a *model.AlbumPreview) bool {
			return a.Type == model.Public
		})
	default:
		break
	}
	return albums, nil
}

func (s *AlbumService) CreateAlbum(album *model.AlbumPreview) (*model.AlbumPreview, error) {
	album.Title = strings.TrimSpace(album.Title)
	album.Desc = strings.TrimSpace(album.Desc)

	if album.Type == "" {
		album.Type = model.Private
	}

	if err := s.validate.Struct(album); err != nil {
		return nil, custom_errors.NewError(custom_errors.ErrBadRequest, err.Error())
	}

	createdAlbum, err := s.storage.Album.Create(album)
	if err != nil {
		return nil, err
	}

	return createdAlbum, nil
}

func (s *AlbumService) DeleteAlbum(id string) error {

	s.redis.Del(fmt.Sprintf("album_%s", id))

	err := s.storage.Album.Delete(id)
	if err != nil {
		return err
	}

	return nil
}

func (s *AlbumService) UpdateAlbum(id string, album *model.AlbumPreview) (*model.AlbumPreview, error) {
	album.Title = strings.TrimSpace(album.Title)
	album.Desc = strings.TrimSpace(album.Desc)

	if err := s.validate.Struct(album); err != nil {
		return nil, custom_errors.NewError(custom_errors.ErrBadRequest, err.Error())
	}

	updatedAlbum, err := s.storage.Album.Update(id, album)
	if err != nil {
		return nil, err
	}

	s.redis.Del(fmt.Sprintf("album_%s", id))

	return updatedAlbum, nil
}

func filterAlbums(albums []*model.AlbumPreview, predicate func(*model.AlbumPreview) bool) []*model.AlbumPreview {
	var filtered []*model.AlbumPreview
	for _, album := range albums {
		if predicate(album) {
			filtered = append(filtered, album)
		}
	}
	return filtered
}

func allowedToReturn(c *gin.Context, album *model.Album, securityConfig *security.AuthConfig) error {
	if album.Type == model.Private {
		if err := security.CheckIsAdmin(c, securityConfig); err != nil {
			return err
		}
	}
	return nil
}
