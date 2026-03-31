package utils

import (
	"bytes"
	"encoding/json"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"strings"

	"github.com/adrium/goheif"
	"github.com/disintegration/imaging"
	"github.com/nfnt/resize"
	"github.com/rwcarlsen/goexif/exif"

	custom_errors "image-service/internal/errors"
)

// CompressImage decodes, normalises (applies EXIF orientation), resizes, and re-encodes
// the provided image bytes. It returns the compressed bytes and the final extension
// actually used for encoding ("jpeg" for both jpg/jpeg/heic inputs, "png" otherwise).
func CompressImage(data []byte, extension string, resizeWidth uint, jpegQuality int) ([]byte, string, error) {
	extension = normalizeExtension(extension)

	img, err := decodeImage(data, extension)
	if err != nil {
		return nil, "", err
	}

	if resizeWidth == 0 {
		resizeWidth = 800
	}
	if jpegQuality == 0 {
		jpegQuality = 80
	}

	// HEIC is always transcoded to JPEG for compatibility.
	targetExt := extension
	if extension == "heic" {
		targetExt = "jpeg"
	}

	compressedImage := resize.Resize(resizeWidth, 0, img, resize.Lanczos3)

	var buf bytes.Buffer

	switch targetExt {
	case "jpeg", "jpg":
		err = jpeg.Encode(&buf, compressedImage, &jpeg.Options{Quality: jpegQuality})
	case "png":
		err = png.Encode(&buf, compressedImage)
	default:
		return nil, "", custom_errors.NewError(custom_errors.ErrBadRequest, "unsupported image format. only jpeg, jpg, and png are allowed")
	}

	if err != nil {
		return nil, "", custom_errors.NewError(custom_errors.ErrInternalServer, "failed to encode image")
	}

	return buf.Bytes(), targetExt, nil
}

func decodeImage(data []byte, extension string) (image.Image, error) {
	extension = normalizeExtension(extension)

	if extension == "heic" {
		return processHEIC(data)
	}

	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, custom_errors.NewError(custom_errors.ErrBadRequest, "failed to decode image")
	}

	orientation := extractOrientation(bytes.NewReader(data))
	if orientation > 1 {
		img = applyOrientation(img, orientation)
	}

	return img, nil
}

func processHEIC(data []byte) (image.Image, error) {
	reader := bytes.NewReader(data)

	img, err := goheif.Decode(reader)
	if err != nil {
		return nil, custom_errors.NewError(custom_errors.ErrBadRequest, "failed to decode HEIC image")
	}

	// Reset cursor before extracting EXIF; goheif.Decode consumes the reader.
	if _, err := reader.Seek(0, io.SeekStart); err != nil {
		return nil, custom_errors.NewError(custom_errors.ErrInternalServer, "failed to rewind HEIC reader")
	}

	exifBytes, err := goheif.ExtractExif(reader)
	if err != nil {
		// Orientation is helpful but non-critical; keep the decoded image.
		return img, nil
	}
	if len(exifBytes) == 0 {
		return img, nil
	}

	orientation := extractOrientation(bytes.NewReader(exifBytes))
	if orientation > 1 {
		img = applyOrientation(img, orientation)
	}

	return img, nil
}

func extractOrientation(r io.Reader) int {
	x, err := exif.Decode(r)
	if err != nil {
		return 1
	}

	tag, err := x.Get(exif.Orientation)
	if err != nil {
		return 1
	}

	val, err := tag.Int(0)
	if err != nil {
		return 1
	}

	return val
}

func normalizeExtension(ext string) string {
	ext = strings.ToLower(ext)
	ext = strings.TrimPrefix(ext, ".")
	return ext
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
