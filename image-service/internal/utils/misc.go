package utils

import (
	"image-service/internal/model"
	"path/filepath"
	"strings"
)

func IsImageFile(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	return ext == ".jpg" || ext == ".jpeg" || ext == ".png"
}

func DetermineImageType(filename string) model.ImageType {
	ext := strings.ToLower(filepath.Ext(filename))
	if ext == ".png" {
		return model.PNG
	} else if ext == ".jpg" {
		return model.JPG
	} else {
		return model.JPEG
	}
}
