package repository

import (
	custom_errors "image-service/internal/errors"
	"image-service/internal/model"
	"image-service/internal/utils"
	"os"
	"path/filepath"
	"strings"
)

type ImageRepository struct {
	Path        string
	ApiBasePath string
}

func (i *ImageRepository) Upload(albumID string, image *model.Image) (*model.Image, error) {
	images, err := i.UploadBatch(albumID, []*model.Image{image})
	if err != nil {
		return nil, err
	}
	return images[0], nil
}

// UploadBatch persists a complete upload under one album lock and updates the
// album image count once, rather than rewriting meta.json for every image.
func (i *ImageRepository) UploadBatch(albumID string, images []*model.Image) ([]*model.Image, error) {
	if len(images) == 0 {
		return []*model.Image{}, nil
	}

	albumPath := filepath.Join(i.Path, albumID)
	metaDataPath := filepath.Join(albumPath, "meta.json")

	if _, err := os.Stat(albumPath); os.IsNotExist(err) {
		return nil, custom_errors.NewError(custom_errors.ErrNotFound, "album not found")
	}

	lock := getAlbumLock(albumID)
	lock.Lock()
	defer lock.Unlock()

	if err := os.MkdirAll(albumPath, os.ModePerm); err != nil {
		return nil, custom_errors.NewError(custom_errors.ErrInternalServer, "failed to create image directory")
	}

	metaData, err := utils.LoadMetaData(metaDataPath)
	if err != nil {
		return nil, err
	}

	type pendingImage struct {
		image *model.Image
		id    string
	}
	pending := make([]pendingImage, len(images))
	writtenPaths := make([]string, 0, len(images))
	rollback := func() {
		for _, path := range writtenPaths {
			_ = os.Remove(path)
		}
	}

	for index, image := range images {
		extension := "." + strings.Split(string(image.Type), "/")[1]
		var imagePath string
		for {
			imagePath = filepath.Join(albumPath, utils.GenerateUUID()+extension)
			if _, err := os.Stat(imagePath); os.IsNotExist(err) {
				break
			}
		}
		if err := os.WriteFile(imagePath, image.Data, os.ModePerm); err != nil {
			rollback()
			return nil, custom_errors.NewError(custom_errors.ErrInternalServer, "failed to save image data")
		}
		writtenPaths = append(writtenPaths, imagePath)
		pending[index] = pendingImage{image: image, id: filepath.Base(imagePath)}
	}

	if err := utils.IncrementImageCountMeta(metaDataPath, metaData, len(images)); err != nil {
		rollback()
		return nil, err
	}

	saved := make([]*model.Image, len(pending))
	for index, item := range pending {
		item.image.ID = item.id
		item.image.Data = nil
		item.image.URL = filepath.Join(i.ApiBasePath, albumID, "images", item.id)
		saved[index] = item.image
	}
	return saved, nil
}

func (i *ImageRepository) Rename(albumID string, imageID string, newName string) (*model.Image, error) {
	imagePath := filepath.Join(i.Path, albumID, imageID)
	newImagePath := filepath.Join(i.Path, albumID, newName)

	lock := getAlbumLock(albumID)
	lock.Lock()
	defer lock.Unlock()

	if _, err := os.Stat(imagePath); os.IsNotExist(err) {
		return nil, custom_errors.NewError(custom_errors.ErrNotFound, "image not found")
	}

	if _, err := os.Stat(newImagePath); err == nil {
		return nil, custom_errors.NewError(custom_errors.ErrConflict, "image with the new name already exists")
	} else if !os.IsNotExist(err) {
		return nil, custom_errors.NewError(custom_errors.ErrInternalServer, "failed to validate new image name")
	}

	if strings.EqualFold(imageID, newName) && imageID != newName {
		return nil, custom_errors.NewError(custom_errors.ErrConflict, "renaming with case-difference only is not allowed")
	}

	err := os.Rename(imagePath, newImagePath)
	if err != nil {
		return nil, custom_errors.NewError(custom_errors.ErrInternalServer, "failed to rename image")
	}

	return &model.Image{
		ID:   newName,
		URL:  filepath.Join(i.ApiBasePath, albumID, "images", newName),
		Type: model.ExtensionsMap[(filepath.Ext(newName))],
	}, nil
}

func (i *ImageRepository) Delete(albumID string, imageID string) error {
	imagePath := filepath.Join(i.Path, albumID, imageID)
	metaDataPath := filepath.Join(i.Path, albumID, "meta.json")

	lock := getAlbumLock(albumID)
	lock.Lock()
	defer lock.Unlock()

	if _, err := os.Stat(imagePath); os.IsNotExist(err) {
		return custom_errors.NewError(custom_errors.ErrNotFound, "image not found")
	}

	err := os.Remove(imagePath)
	if err != nil {
		return custom_errors.NewError(custom_errors.ErrInternalServer, "failed to delete image")
	}

	metaData, err := utils.LoadMetaData(metaDataPath)
	if err != nil {
		return err
	}

	err = utils.DecrementImageCountMeta(metaDataPath, metaData, 1)
	if err != nil {
		return err
	}

	return nil
}
