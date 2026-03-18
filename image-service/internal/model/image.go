package model

type ImageType string

var ExtensionsMap = map[string]ImageType{
	".jpeg": JPEG,
	".jpg":  JPG,
	".png":  PNG,
	".heic": HEIC,
}

const (
	JPEG ImageType = "image/jpeg"
	JPG  ImageType = "image/jpg"
	PNG  ImageType = "image/png"
	HEIC ImageType = "image/heic"
)

var AllowedTypes = map[ImageType]bool{
	JPEG: true,
	JPG:  true,
	PNG:  true,
	HEIC: true,
}

type Image struct {
	ID   string    `json:"id"`
	Type ImageType `json:"type"`
	Data []byte    `json:"data"`
	URL  string    `json:"url"`
}
