package model

import "github.com/go-playground/validator/v10"

type AlbumType string

const (
	Private    AlbumType = "private"
	SemiPublic AlbumType = "semi-public"
	Public     AlbumType = "public"
)

func (t AlbumType) IsValid() bool {
	switch t {
	case Private, SemiPublic, Public:
		return true
	}
	return false
}

func ValidateAlbumType(fl validator.FieldLevel) bool {
	val, ok := fl.Field().Interface().(AlbumType)
	if !ok {
		return false
	}
	return val.IsValid()
}

type Album struct {
	ID     string    `json:"id"`
	Title  string    `json:"title" validate:"required"`
	Desc   string    `json:"desc"`
	Date   string    `json:"date"`
	Type   AlbumType `json:"type" validate:"required"`
	Images []*Image  `json:"images"`
}

type AlbumPreview struct {
	ID              string    `json:"id"`
	Title           string    `json:"title" validate:"required"`
	Desc            string    `json:"desc"`
	Date            string    `json:"date"`
	ImageCount      int       `json:"image_count"`
	Type            AlbumType `json:"type" validate:"required,albumtype"`
	PreviewImageURL string    `json:"preview_image"`
}
