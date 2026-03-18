package utils

import (
	"image-service/internal/model"
	"testing"
)

func TestIsImageFile(t *testing.T) {
	tests := []struct {
		filename string
		expected bool
	}{
		{"image.jpg", true},
		{"image.jpeg", true},
		{"image.png", true},
		{"image.gif", false},
		{"document.pdf", false},
	}

	for _, test := range tests {
		result := IsImageFile(test.filename)
		if result != test.expected {
			t.Errorf("IsImageFile(%s) = %v; want %v", test.filename, result, test.expected)
		}
	}
}

func TestDetermineImageType(t *testing.T) {
	tests := []struct {
		filename string
		expected model.ImageType
	}{
		{"image.jpg", model.JPG},
		{"image.jpeg", model.JPEG},
		{"image.png", model.PNG},
		{"image.gif", model.JPEG},
		{"document.pdf", model.JPEG},
	}

	for _, test := range tests {
		result := DetermineImageType(test.filename)
		if result != test.expected {
			t.Errorf("DetermineImageType(%s) = %v; want %v", test.filename, result, test.expected)
		}
	}
}
