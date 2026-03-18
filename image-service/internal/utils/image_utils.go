package utils

import (
	"bytes"
	"encoding/json"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"sync"

	"github.com/adrium/goheif"
	"github.com/disintegration/imaging"
	"github.com/nfnt/resize"
	"github.com/rwcarlsen/goexif/exif"

	custom_errors "image-service/internal/errors"
)

var MetadataMutex sync.Mutex

func CompressImage(reader *bytes.Reader, extension string) ([]byte, error) {
	var img image.Image
	var err error
	if extension == "heic" {
		var err error
		img, err = ProcessHEIC(reader)
		if err != nil {
			return nil, err
		}
		extension = "jpeg"
	} else {
		img, _, err = image.Decode(reader)
		if err != nil {
			return nil, custom_errors.NewError(custom_errors.ErrBadRequest, "failed to decode image")
		}
	}

	compressedImage := resize.Resize(800, 0, img, resize.Lanczos3)

	var buf bytes.Buffer

	switch extension {
	case "jpeg", "jpg":
		err = jpeg.Encode(&buf, compressedImage, &jpeg.Options{Quality: 80})
	case "png":
		err = png.Encode(&buf, compressedImage)
	default:
		return nil, custom_errors.NewError(custom_errors.ErrBadRequest, "unsupported image format. only jpeg, jpg, and png are allowed")
	}

	if err != nil {
		return nil, custom_errors.NewError(custom_errors.ErrInternalServer, "failed to encode image")
	}

	return buf.Bytes(), nil
}

func ProcessHEIC(reader *bytes.Reader) (image.Image, error) {
	img, err := goheif.Decode(reader)
	if err != nil {
		return nil, custom_errors.NewError(custom_errors.ErrBadRequest, "failed to decode HEIC image")
	}

	exifBytes, err := goheif.ExtractExif(reader)
	if err != nil {
		return nil, custom_errors.NewError(custom_errors.ErrBadRequest, "failed to extract EXIF data from HEIC image")
	}
	if len(exifBytes) == 0 {
		return img, nil
	}

	x, err := exif.Decode(bytes.NewReader(exifBytes))
	if err != nil {
		return nil, custom_errors.NewError(custom_errors.ErrBadRequest, "failed to decode EXIF data")
	}

	orientation, err := x.Get(exif.Orientation)
	if err != nil {
		return img, nil // No orientation tag found, fallback to unrotated image
	}
	orientationValue, err := orientation.Int(0)
	if err != nil {
		return img, nil // Fallback to unrotated image
	}

	return applyOrientation(img, orientationValue), nil
}

func applyOrientation(img image.Image, orientation int) image.Image {
	switch orientation {
	case 2:
		return imaging.FlipH(img)
	case 3:
		return imaging.Rotate180(img)
	case 4:
		return imaging.FlipV(img)
	case 5:
		return imaging.Transpose(img)
	case 6:
		return imaging.Rotate270(img) // 90 CW
	case 7:
		return imaging.Transverse(img)
	case 8:
		return imaging.Rotate90(img) // 270 CW
	default:
		return img
	}
}

func GetImageCount(metaData map[string]any) (int, error) {
	imageCountFloat, ok := metaData["ImageCount"].(float64)
	if !ok {
		return 0, custom_errors.NewError(custom_errors.ErrInternalServer, "invalid type for ImageCount")
	}
	return int(imageCountFloat), nil
}

func LoadMetaData(metaDataPath string) (map[string]any, error) {
	content, err := os.ReadFile(metaDataPath)
	if err != nil {
		return nil, custom_errors.NewError(custom_errors.ErrInternalServer, "failed to read metadata file")
	}

	var metaData map[string]any
	err = json.Unmarshal(content, &metaData)
	if err != nil {
		return nil, custom_errors.NewError(custom_errors.ErrInternalServer, "failed to parse metadata JSON")
	}

	return metaData, nil
}

func IncrementImageCountMeta(metaDataPath string, metaData map[string]any, n int) error {
	countFloat, ok := metaData["ImageCount"].(float64)
	if !ok {
		return custom_errors.NewError(custom_errors.ErrInternalServer, "invalid type for ImageCount")
	}

	count := int(countFloat)
	count += n
	metaData["ImageCount"] = count

	newContent, err := json.MarshalIndent(metaData, "", "  ")
	if err != nil {
		return custom_errors.NewError(custom_errors.ErrInternalServer, "failed to marshal updated metadata")
	}

	err = os.WriteFile(metaDataPath, newContent, os.ModePerm)
	if err != nil {
		return custom_errors.NewError(custom_errors.ErrInternalServer, "failed to write updated metadata")
	}

	return nil
}

func DecrementImageCountMeta(metaDataPath string, metaData map[string]any, n int) error {
	countFloat, ok := metaData["ImageCount"].(float64)
	if !ok {
		return custom_errors.NewError(custom_errors.ErrInternalServer, "invalid type for ImageCount")
	}

	count := int(countFloat)
	count -= n
	if count < 0 {
		count = 0
	}
	metaData["ImageCount"] = count

	newContent, err := json.MarshalIndent(metaData, "", "  ")
	if err != nil {
		return custom_errors.NewError(custom_errors.ErrInternalServer, "failed to marshal updated metadata")
	}

	err = os.WriteFile(metaDataPath, newContent, os.ModePerm)
	if err != nil {
		return custom_errors.NewError(custom_errors.ErrInternalServer, "failed to write updated metadata")
	}

	return nil
}
