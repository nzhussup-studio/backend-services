package repository

import (
	"testing"
)

func TestNewStorage(t *testing.T) {
	path := "/tmp/images"
	apiBasePath := "/api/images"

	storage := NewStorage(path, apiBasePath)

	if storage.Path != path {
		t.Errorf("expected Path to be %q, got %q", path, storage.Path)
	}

	if storage.ApiBasePath != apiBasePath {
		t.Errorf("expected ApiBasePath to be %q, got %q", apiBasePath, storage.ApiBasePath)
	}

	if storage.Album == nil {
		t.Error("expected Album repository to be initialized, got nil")
	}

	if storage.Image == nil {
		t.Error("expected Image repository to be initialized, got nil")
	}
}
