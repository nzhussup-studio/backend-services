package utils

import (
	"encoding/json"
	"errors"
	"image"
	"os"
	"path/filepath"
	"testing"

	custom_errors "image-service/internal/errors"

	"github.com/stretchr/testify/assert"
)

func TestGetImageCount_Success(t *testing.T) {
	meta := map[string]any{"ImageCount": 5.0}
	count, err := GetImageCount(meta)
	assert.NoError(t, err)
	assert.Equal(t, 5, count)
}

func TestGetImageCount_Failure(t *testing.T) {
	meta := map[string]any{"ImageCount": "not a float"}
	count, err := GetImageCount(meta)
	assert.Error(t, err)
	assert.Equal(t, 0, count)
	assert.True(t, errors.Is(err, custom_errors.ErrInternalServer))
}

func TestLoadMetaData_Success(t *testing.T) {
	tmpDir := t.TempDir()
	metaPath := filepath.Join(tmpDir, "meta.json")
	content := `{"ImageCount": 10}`
	err := os.WriteFile(metaPath, []byte(content), 0644)
	assert.NoError(t, err)

	meta, err := LoadMetaData(metaPath)
	assert.NoError(t, err)
	assert.Equal(t, float64(10), meta["ImageCount"])
}

func TestLoadMetaData_FileNotFound(t *testing.T) {
	_, err := LoadMetaData("nonexistent.json")
	assert.Error(t, err)
	assert.True(t, errors.Is(err, custom_errors.ErrInternalServer))
}

func TestLoadMetaData_InvalidJson(t *testing.T) {
	tmpDir := t.TempDir()
	metaPath := filepath.Join(tmpDir, "meta.json")
	err := os.WriteFile(metaPath, []byte("{invalid json}"), 0644)
	assert.NoError(t, err)

	_, err = LoadMetaData(metaPath)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, custom_errors.ErrInternalServer))
}

func TestIncrementImageCountMeta_Success(t *testing.T) {
	tmpDir := t.TempDir()
	metaPath := filepath.Join(tmpDir, "meta.json")
	meta := map[string]any{"ImageCount": 2.0}
	content, _ := json.Marshal(meta)
	err := os.WriteFile(metaPath, content, 0644)
	assert.NoError(t, err)

	err = IncrementImageCountMeta(metaPath, meta, 3)
	assert.NoError(t, err)

	updatedMeta, err := LoadMetaData(metaPath)
	assert.NoError(t, err)
	assert.Equal(t, float64(5), updatedMeta["ImageCount"])
}

func TestIncrementImageCountMeta_InvalidType(t *testing.T) {
	tmpDir := t.TempDir()
	metaPath := filepath.Join(tmpDir, "meta.json")
	meta := map[string]any{"ImageCount": "invalid"}
	content, _ := json.Marshal(meta)
	err := os.WriteFile(metaPath, content, 0644)
	assert.NoError(t, err)

	err = IncrementImageCountMeta(metaPath, meta, 1)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, custom_errors.ErrInternalServer))
}

func TestDecrementImageCountMeta_Success(t *testing.T) {
	tmpDir := t.TempDir()
	metaPath := filepath.Join(tmpDir, "meta.json")
	meta := map[string]any{"ImageCount": 5.0}
	content, _ := json.Marshal(meta)
	err := os.WriteFile(metaPath, content, 0644)
	assert.NoError(t, err)

	err = DecrementImageCountMeta(metaPath, meta, 3)
	assert.NoError(t, err)

	updatedMeta, err := LoadMetaData(metaPath)
	assert.NoError(t, err)
	assert.Equal(t, float64(2), updatedMeta["ImageCount"])
}

func TestDecrementImageCountMeta_ClampsToZero(t *testing.T) {
	tmpDir := t.TempDir()
	metaPath := filepath.Join(tmpDir, "meta.json")
	meta := map[string]any{"ImageCount": 1.0}
	content, _ := json.Marshal(meta)
	err := os.WriteFile(metaPath, content, 0644)
	assert.NoError(t, err)

	err = DecrementImageCountMeta(metaPath, meta, 5)
	assert.NoError(t, err)

	updatedMeta, err := LoadMetaData(metaPath)
	assert.NoError(t, err)
	assert.Equal(t, float64(0), updatedMeta["ImageCount"])
}

func TestDecrementImageCountMeta_InvalidType(t *testing.T) {
	tmpDir := t.TempDir()
	metaPath := filepath.Join(tmpDir, "meta.json")
	meta := map[string]any{"ImageCount": "invalid"}
	content, _ := json.Marshal(meta)
	err := os.WriteFile(metaPath, content, 0644)
	assert.NoError(t, err)

	err = DecrementImageCountMeta(metaPath, meta, 1)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, custom_errors.ErrInternalServer))
}

func TestApplyOrientation(t *testing.T) {
	rect := image.Rect(0, 0, 100, 100)
	img := image.NewRGBA(rect)

	// Just test that it returns an image (cannot easily test rotated image here)
	for _, orientation := range []int{1, 2, 3, 4, 5, 6, 7, 8, 999} {
		out := applyOrientation(img, orientation)
		assert.NotNil(t, out)
	}
}
