package repository_test

import (
	"encoding/json"
	"image-service/internal/model"
	"image-service/internal/repository"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAlbumRepository_Update(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "album_repo_test")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir) // Clean up

	repo := &repository.AlbumRepository{Path: tempDir}
	albumID := "test-album"

	albumPath := filepath.Join(tempDir, albumID)
	err = os.Mkdir(albumPath, 0755)
	assert.NoError(t, err)

	metaFilePath := filepath.Join(albumPath, "meta.json")
	initialMeta := map[string]interface{}{
		"Title":           "Old Title",
		"Description":     "Old Description",
		"Date":            "Old Date",
		"ImageCount":      1,
		"Type":            "Gallery",
		"PreviewImageURL": "old-preview.jpg",
	}
	initialMetaBytes, _ := json.MarshalIndent(initialMeta, "", "  ")
	err = os.WriteFile(metaFilePath, initialMetaBytes, 0644)
	assert.NoError(t, err)

	newMeta := &model.AlbumPreview{
		Title:           "New Title",
		Desc:            "New Description",
		Date:            "2025-04-28",
		Type:            "Gallery",
		PreviewImageURL: "new-preview.jpg",
	}

	updated, err := repo.Update(albumID, newMeta)
	assert.NoError(t, err)
	assert.Equal(t, "New Title", updated.Title)
	assert.Equal(t, albumID, updated.ID)

	content, err := os.ReadFile(metaFilePath)
	assert.NoError(t, err)

	var updatedMeta map[string]interface{}
	err = json.Unmarshal(content, &updatedMeta)
	assert.NoError(t, err)

	assert.Equal(t, "New Title", updatedMeta["Title"])
	assert.Equal(t, "New Description", updatedMeta["Description"])
	assert.Equal(t, "new-preview.jpg", updatedMeta["PreviewImageURL"])
}

func TestAlbumRepository_Update_Concurrent(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "album_repo_test_concurrent")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	repo := &repository.AlbumRepository{Path: tempDir}
	albumID := "test-album"

	albumPath := filepath.Join(tempDir, albumID)
	err = os.Mkdir(albumPath, 0755)
	assert.NoError(t, err)

	metaFilePath := filepath.Join(albumPath, "meta.json")

	initialMeta := map[string]interface{}{
		"Title":           "Initial Title",
		"Description":     "Initial Description",
		"Date":            "2025-04-28",
		"ImageCount":      0,
		"Type":            "Gallery",
		"PreviewImageURL": "initial-preview.jpg",
	}
	initialMetaBytes, _ := json.MarshalIndent(initialMeta, "", "  ")
	err = os.WriteFile(metaFilePath, initialMetaBytes, 0644)
	assert.NoError(t, err)

	var wg sync.WaitGroup
	updateCount := 1000

	for i := 0; i < updateCount; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			newMeta := &model.AlbumPreview{
				Title:           "Title " + time.Now().Format(time.RFC3339Nano),
				Desc:            "Desc " + time.Now().Format(time.RFC3339Nano),
				Date:            "2025-04-28",
				Type:            "Gallery",
				PreviewImageURL: "preview-" + time.Now().Format("150405.000") + ".jpg",
			}

			_, err := repo.Update(albumID, newMeta)
			assert.NoError(t, err)
		}(i)
	}

	wg.Wait()

	content, err := os.ReadFile(metaFilePath)
	assert.NoError(t, err)

	var finalMeta map[string]any
	err = json.Unmarshal(content, &finalMeta)
	assert.NoError(t, err)

	t.Logf("Final meta.json content:\n%v", finalMeta)
}
