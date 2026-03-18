package repository

import (
	"errors"
	custom_errors "image-service/internal/errors"
	"image-service/internal/model"
	"image-service/internal/utils"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testLock = &sync.Mutex{}

func init() {
	getAlbumLock = func(albumID string) *sync.Mutex {
		return testLock
	}
}

func setupTestAlbum(t *testing.T, basePath, albumID string) string {
	albumPath := filepath.Join(basePath, albumID)
	err := os.MkdirAll(albumPath, os.ModePerm)
	assert.NoError(t, err)

	metaDataPath := filepath.Join(albumPath, "meta.json")
	initialMeta := `{"ImageCount": 0}`
	err = os.WriteFile(metaDataPath, []byte(initialMeta), 0644)
	assert.NoError(t, err)

	return albumPath
}

func TestUpload_Success(t *testing.T) {
	basePath := t.TempDir()
	albumID := "album1"
	setupTestAlbum(t, basePath, albumID)

	repo := &ImageRepository{
		Path:        basePath,
		ApiBasePath: "/api/images",
	}

	img := &model.Image{
		Type: "image/jpeg",
		Data: []byte("test image data"),
	}

	uploadedImage, err := repo.Upload(albumID, img)
	assert.NoError(t, err)
	assert.NotNil(t, uploadedImage)
	assert.NotEmpty(t, uploadedImage.ID)
	assert.Nil(t, uploadedImage.Data)
	assert.Contains(t, uploadedImage.URL, albumID)
	assert.Contains(t, uploadedImage.URL, uploadedImage.ID)

	// Check file exists
	imagePath := filepath.Join(basePath, albumID, uploadedImage.ID)
	_, err = os.Stat(imagePath)
	assert.NoError(t, err)

	// Check meta.json increment
	metaData, err := utils.LoadMetaData(filepath.Join(basePath, albumID, "meta.json"))
	assert.NoError(t, err)
	count, err := utils.GetImageCount(metaData)
	assert.NoError(t, err)
	assert.Equal(t, 1, count)
}

func TestUpload_AlbumNotFound(t *testing.T) {
	basePath := t.TempDir()
	repo := &ImageRepository{
		Path:        basePath,
		ApiBasePath: "/api/images",
	}

	img := &model.Image{
		Type: "image/jpeg",
		Data: []byte("test image data"),
	}

	_, err := repo.Upload("nonexistent-album", img)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, custom_errors.ErrNotFound))
}

func TestDelete_Success(t *testing.T) {
	basePath := t.TempDir()
	albumID := "album1"
	setupTestAlbum(t, basePath, albumID)

	// Prepare an image file to delete
	imageID := "testimg.jpg"
	imagePath := filepath.Join(basePath, albumID, imageID)
	err := os.WriteFile(imagePath, []byte("image data"), 0644)
	assert.NoError(t, err)

	// Update meta to 1
	metaDataPath := filepath.Join(basePath, albumID, "meta.json")
	metaData, err := utils.LoadMetaData(metaDataPath)
	assert.NoError(t, err)
	err = utils.IncrementImageCountMeta(metaDataPath, metaData, 1)
	assert.NoError(t, err)

	repo := &ImageRepository{
		Path:        basePath,
		ApiBasePath: "/api/images",
	}

	err = repo.Delete(albumID, imageID)
	assert.NoError(t, err)

	// Confirm file deleted
	_, err = os.Stat(imagePath)
	assert.True(t, os.IsNotExist(err))

	// Confirm meta decremented
	metaData, err = utils.LoadMetaData(metaDataPath)
	assert.NoError(t, err)
	count, err := utils.GetImageCount(metaData)
	assert.NoError(t, err)
	assert.Equal(t, 0, count)
}

func TestDelete_ImageNotFound(t *testing.T) {
	basePath := t.TempDir()
	albumID := "album1"
	setupTestAlbum(t, basePath, albumID)

	repo := &ImageRepository{
		Path:        basePath,
		ApiBasePath: "/api/images",
	}

	err := repo.Delete(albumID, "nonexistent.jpg")
	assert.Error(t, err)
	assert.True(t, errors.Is(err, custom_errors.ErrNotFound))
}

func TestRename(t *testing.T) {
	basePath := t.TempDir()
	albumID := "album1"
	albumPath := setupTestAlbum(t, basePath, albumID)

	repo := &ImageRepository{
		Path:        basePath,
		ApiBasePath: "/api/images",
	}

	// Create initial image
	originalName := "test.jpg"
	originalPath := filepath.Join(albumPath, originalName)
	err := os.WriteFile(originalPath, []byte("image data"), 0644)
	assert.NoError(t, err)

	// Create a conflicting image file
	conflictName := "conflict.jpg"
	conflictPath := filepath.Join(albumPath, conflictName)
	err = os.WriteFile(conflictPath, []byte("conflict data"), 0644)
	assert.NoError(t, err)

	tests := []struct {
		name         string
		albumID      string
		imageID      string
		newName      string
		expectError  bool
		expectedBase error
	}{
		{
			name:         "successful rename",
			albumID:      albumID,
			imageID:      originalName,
			newName:      "renamed.jpg",
			expectError:  false,
			expectedBase: nil,
		},
		{
			name:         "image not found",
			albumID:      albumID,
			imageID:      "missing.jpg",
			newName:      "renamed.jpg",
			expectError:  true,
			expectedBase: custom_errors.ErrNotFound,
		},
		{
			name:         "new name already exists",
			albumID:      albumID,
			imageID:      originalName,
			newName:      conflictName,
			expectError:  true,
			expectedBase: custom_errors.ErrConflict,
		},
		{
			name:         "case-only rename not allowed",
			albumID:      albumID,
			imageID:      originalName,
			newName:      "TEST.jpg",
			expectError:  true,
			expectedBase: custom_errors.ErrConflict,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Recreate original file if it was renamed in previous test
			if tt.imageID == originalName {
				_, err := os.Stat(originalPath)
				if os.IsNotExist(err) {
					_ = os.WriteFile(originalPath, []byte("image data"), 0644)
				}
			}

			img, err := repo.Rename(tt.albumID, tt.imageID, tt.newName)

			if tt.expectError {
				assert.Error(t, err)
				assert.True(t, errors.Is(err, tt.expectedBase), "expected base error: %v, got: %v", tt.expectedBase, err)
				assert.Nil(t, img)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, img)
				assert.Equal(t, tt.newName, img.ID)

				// Confirm file physically exists
				newPath := filepath.Join(albumPath, tt.newName)
				_, statErr := os.Stat(newPath)
				assert.NoError(t, statErr)
			}
		})
	}
}
