package service

import (
	"fmt"
	"image-service/internal/cache"
	custom_errors "image-service/internal/errors"
	"image-service/internal/model"
	"image-service/internal/repository"
	"image-service/internal/utils"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/go-playground/validator/v10"
)

type ImageService struct {
	storage  *repository.Storage
	redis    cache.RedisClientInterface
	validate *validator.Validate
	cfg      ImageConfig
}

func (s *ImageService) UploadImage(albumID string, files []*multipart.FileHeader) ([]*model.Image, error) {
	var (
		savedImages []*model.Image
		mu          sync.Mutex
		wg          sync.WaitGroup
		errChan     = make(chan error, len(files))
	)

	// Pre-validate image types
	for _, file := range files {
		contentType := file.Header.Get("Content-Type")
		imageType := model.ImageType(contentType)
		if _, ok := model.AllowedTypes[imageType]; !ok {
			return nil, custom_errors.NewError(custom_errors.ErrBadRequest, fmt.Sprintf("invalid image type: %s. Only JPEG, PNG, and HEIC are allowed", contentType))
		}
	}

	for _, file := range files {
		wg.Add(1)
		go func(file *multipart.FileHeader) {
			defer wg.Done()

			if file.Size > 0 && file.Size > s.cfg.MaxUploadBytes {
				errChan <- custom_errors.NewError(custom_errors.ErrBadRequest, fmt.Sprintf("image file too large: max %d MB", s.cfg.MaxUploadBytes/1024/1024))
				return
			}

			fileData, err := file.Open()
			if err != nil {
				errChan <- custom_errors.NewError(custom_errors.ErrInternalServer, "failed to open image file")
				return
			}
			defer fileData.Close()

			image := &model.Image{
				Type: model.ImageType(file.Header.Get("Content-Type")),
			}

			limitedReader := io.LimitReader(fileData, s.cfg.MaxUploadBytes+1)
			data, err := io.ReadAll(limitedReader)
			if err != nil {
				errChan <- custom_errors.NewError(custom_errors.ErrInternalServer, "failed to read image file")
				return
			}
			if int64(len(data)) > s.cfg.MaxUploadBytes {
				errChan <- custom_errors.NewError(custom_errors.ErrBadRequest, fmt.Sprintf("image file too large: max %d MB", s.cfg.MaxUploadBytes/1024/1024))
				return
			}

			extension := strings.ToLower(strings.Split(string(image.Type), "/")[1])

			compressedData, finalExt, err := utils.CompressImage(data, extension, s.cfg.ResizeWidth, s.cfg.JPEGQuality)
			if err != nil {
				errChan <- err
				return
			}

			image.Type = mimeTypeFromExtension(finalExt)
			image.Data = compressedData

			savedImage, err := s.storage.Image.Upload(albumID, image)
			if err != nil {
				errChan <- err
				return
			}

			cacheKey := fmt.Sprintf("image:%s:%s", albumID, savedImage.ID)
			cacheValue := fmt.Sprintf("%s/%s/%s", s.storage.Path, albumID, savedImage.ID)
			s.redis.Set(cacheKey, cacheValue)

			mu.Lock()
			savedImages = append(savedImages, savedImage)
			mu.Unlock()
		}(file)
	}

	wg.Wait()
	close(errChan)

	if len(errChan) > 0 {
		return nil, <-errChan
	}

	s.redis.Del(fmt.Sprintf("album_%s", albumID))

	return savedImages, nil
}

func (s *ImageService) RenameImage(albumID string, imageID string, newName string) (*model.Image, error) {
	if err := s.validate.Var(newName, "required,alphanum"); err != nil {
		return nil, custom_errors.NewError(custom_errors.ErrBadRequest, "invalid image name")
	}

	if newName != filepath.Base(newName) || strings.Contains(newName, string(os.PathSeparator)) || strings.Contains(newName, "..") {
		return nil, custom_errors.NewError(custom_errors.ErrBadRequest, "image name must not include path components or directory traversal")
	}

	extension := filepath.Ext(imageID)
	newName = fmt.Sprintf("%s%s", newName, extension)

	image, err := s.storage.Image.Rename(albumID, imageID, newName)
	if err != nil {
		return nil, err
	}

	// CACHE EVICTION
	cacheKey := fmt.Sprintf("image:%s:%s", albumID, imageID)
	s.redis.Del(cacheKey)
	s.redis.Del(fmt.Sprintf("album_%s", albumID))

	return image, nil
}

func (s *ImageService) DeleteImage(albumID string, imageID string) error {
	err := s.storage.Image.Delete(albumID, imageID)
	if err != nil {
		return err
	}

	// CACHE EVICTION
	cacheKey := fmt.Sprintf("image:%s:%s", albumID, imageID)
	s.redis.Del(cacheKey)
	s.redis.Del(fmt.Sprintf("album_%s", albumID))

	return nil
}

func (s *ImageService) ServeImage(albumID string, imageID string) (string, error) {
	// CACHE CHECK
	cacheKey := fmt.Sprintf("image:%s:%s", albumID, imageID)
	var cachedImagePath string
	err := s.redis.Get(cacheKey, &cachedImagePath)
	if err == nil {
		return cachedImagePath, nil
	}

	imagePath := filepath.Join(s.storage.Path, albumID, imageID)
	if _, err := os.Stat(imagePath); os.IsNotExist(err) {
		return "", custom_errors.NewError(custom_errors.ErrNotFound, fmt.Sprintf("image %s not found in album %s", imageID, albumID))
	}

	// CACHING
	s.redis.Set(cacheKey, imagePath)
	return imagePath, nil
}

func mimeTypeFromExtension(ext string) model.ImageType {
	switch strings.ToLower(ext) {
	case "jpeg", "jpg":
		return model.JPEG
	case "png":
		return model.PNG
	case "heic":
		return model.HEIC
	default:
		return model.ImageType(fmt.Sprintf("image/%s", ext))
	}
}
