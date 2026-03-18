package repository

import (
	"encoding/json"
	"fmt"
	custom_errors "image-service/internal/errors"
	"image-service/internal/model"
	"image-service/internal/utils"
	"os"
	"path/filepath"
)

type AlbumRepository struct {
	Path        string
	ApiBasePath string
}

func (a *AlbumRepository) Create(m *model.AlbumPreview) (*model.AlbumPreview, error) {
	albumID := utils.GenerateUUID()
	albumPath := filepath.Join(a.Path, albumID)

	if _, err := os.Stat(albumPath); err == nil {
		return nil, custom_errors.NewError(custom_errors.ErrConflict, "album already exists")
	}

	lock := getAlbumLock(albumID)
	lock.Lock()
	defer lock.Unlock()

	err := os.MkdirAll(albumPath, os.ModePerm)
	if err != nil {
		return nil, custom_errors.NewError(custom_errors.ErrInternalServer, "failed to create album")
	}

	metaFilePath := filepath.Join(albumPath, "meta.json")
	metaFile, err := os.Create(metaFilePath)
	if err != nil {
		return nil, custom_errors.NewError(custom_errors.ErrInternalServer, "failed to create metadata file")
	}
	defer metaFile.Close()

	albumMetadata := map[string]any{
		"Title":           m.Title,
		"Description":     m.Desc,
		"Date":            m.Date,
		"ImageCount":      0,
		"Type":            m.Type,
		"PreviewImageURL": m.PreviewImageURL,
	}

	encoder := json.NewEncoder(metaFile)
	err = encoder.Encode(albumMetadata)
	if err != nil {
		return nil, custom_errors.NewError(custom_errors.ErrInternalServer, "failed to write metadata to JSON")
	}

	m.ID = albumID
	return m, nil
}

func (a *AlbumRepository) Get(id string) (*model.Album, error) {
	albumPath := filepath.Join(a.Path, id)
	if _, err := os.Stat(albumPath); os.IsNotExist(err) {
		return nil, custom_errors.NewError(custom_errors.ErrNotFound, "album not found")
	}

	lock := getAlbumLock(id)
	lock.Lock()
	defer lock.Unlock()

	metaFilePath := filepath.Join(albumPath, "meta.json")
	albumMetadata, err := utils.LoadMetaData(metaFilePath)
	if err != nil {
		return nil, err
	}

	album := &model.Album{
		ID:    id,
		Title: albumMetadata["Title"].(string),
		Desc:  albumMetadata["Description"].(string),
		Date:  albumMetadata["Date"].(string),
		Type:  model.AlbumType(albumMetadata["Type"].(string)),
	}

	files, err := os.ReadDir(albumPath)
	if err != nil {
		return nil, custom_errors.NewError(custom_errors.ErrInternalServer, "failed to read album directory")
	}

	var images []*model.Image
	for _, file := range files {
		if file.IsDir() || !utils.IsImageFile(file.Name()) {
			continue
		}

		imageURL := fmt.Sprintf("%s/%s/%s", a.ApiBasePath, id, file.Name())

		image := &model.Image{
			ID:   file.Name(),
			Type: utils.DetermineImageType(file.Name()),
			Data: nil,
			URL:  imageURL,
		}
		images = append(images, image)
	}

	album.Images = images
	return album, nil
}

func (a *AlbumRepository) GetPreview() ([]*model.AlbumPreview, error) {
	dirs, err := os.ReadDir(a.Path)
	if err != nil {
		return nil, custom_errors.NewError(custom_errors.ErrInternalServer, "failed to read storage path")
	}

	var albumPreviews []*model.AlbumPreview

	for _, dir := range dirs {
		if dir.IsDir() {
			// scoped anonymous function to manage lock/unlock per album
			albumPreview, err := func(dirName string) (*model.AlbumPreview, error) {
				lock := getAlbumLock(dirName)
				lock.Lock()
				defer lock.Unlock()

				metaFilePath := filepath.Join(a.Path, dirName, "meta.json")
				albumMetadata, err := utils.LoadMetaData(metaFilePath)
				if err != nil {
					return nil, err
				}

				imageCount, err := utils.GetImageCount(albumMetadata)
				if err != nil {
					return nil, err
				}

				return &model.AlbumPreview{
					ID:              dirName,
					Title:           albumMetadata["Title"].(string),
					Desc:            albumMetadata["Description"].(string),
					Date:            albumMetadata["Date"].(string),
					ImageCount:      imageCount,
					Type:            model.AlbumType(albumMetadata["Type"].(string)),
					PreviewImageURL: albumMetadata["PreviewImageURL"].(string),
				}, nil
			}(dir.Name())

			if err != nil {
				return nil, err
			}

			albumPreviews = append(albumPreviews, albumPreview)
		}
	}
	return albumPreviews, nil
}

func (a *AlbumRepository) Delete(id string) error {
	albumPath := filepath.Join(a.Path, id)
	if _, err := os.Stat(albumPath); os.IsNotExist(err) {
		return custom_errors.NewError(custom_errors.ErrNotFound, "album not found")
	}

	lock := getAlbumLock(id)
	lock.Lock()
	defer lock.Unlock()

	err := os.RemoveAll(albumPath)
	if err != nil {
		return custom_errors.NewError(custom_errors.ErrInternalServer, "failed to delete album directory")
	}

	return nil
}

func (a *AlbumRepository) Update(id string, m *model.AlbumPreview) (*model.AlbumPreview, error) {
	albumPath := filepath.Join(a.Path, id)
	if _, err := os.Stat(albumPath); os.IsNotExist(err) {
		return nil, custom_errors.NewError(custom_errors.ErrNotFound, "album not found")
	}

	lock := getAlbumLock(id)
	lock.Lock()
	defer lock.Unlock()

	metaFilePath := filepath.Join(albumPath, "meta.json")
	albumMetadata, err := utils.LoadMetaData(metaFilePath)
	if err != nil {
		return nil, err
	}

	albumMetadata["Title"] = m.Title
	albumMetadata["Description"] = m.Desc
	albumMetadata["Date"] = m.Date
	albumMetadata["Type"] = m.Type
	albumMetadata["PreviewImageURL"] = m.PreviewImageURL

	metaFile, err := os.Create(metaFilePath)
	if err != nil {
		return nil, custom_errors.NewError(custom_errors.ErrInternalServer, "failed to open metadata file for writing")
	}
	defer metaFile.Close()

	encoder := json.NewEncoder(metaFile)
	err = encoder.Encode(albumMetadata)
	if err != nil {
		return nil, custom_errors.NewError(custom_errors.ErrInternalServer, "failed to write updated metadata to JSON")
	}

	m.ID = id
	return m, nil
}
