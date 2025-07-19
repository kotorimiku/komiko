package utils

import (
	"bytes"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"math"
	"path/filepath"
	"strings"

	"github.com/chai2010/webp"
	"github.com/gen2brain/avif"
	"golang.org/x/image/draw"
)

func ImgToPng(imageReader io.Reader) ([]byte, error) {
	img, format, err := image.Decode(imageReader)
	if err != nil {
		return nil, err
	}

	if format == "png" {
		data, err := io.ReadAll(imageReader)
		if err != nil {
			return nil, err
		}
		return data, nil
	}

	var buf bytes.Buffer

	err = png.Encode(&buf, img)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func ImgToJpg(imageReader io.Reader) ([]byte, error) {
	img, format, err := image.Decode(imageReader)
	if err != nil {
		return nil, err
	}

	if format == "jpeg" {
		data, err := io.ReadAll(imageReader)
		if err != nil {
			return nil, err
		}
		return data, nil
	}

	var buf bytes.Buffer

	err = jpeg.Encode(&buf, img, nil)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func ImgToAvif(imageReader io.Reader) ([]byte, error) {
	img, format, err := image.Decode(imageReader)
	if err != nil {
		return nil, err
	}

	if format == "avif" {
		data, err := io.ReadAll(imageReader)
		if err != nil {
			return nil, err
		}
		return data, nil
	}

	var buf bytes.Buffer

	err = avif.Encode(&buf, img, avif.Options{Quality: 0, Speed: 10, ChromaSubsampling: image.YCbCrSubsampleRatio420})
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func ImgToWebp(imageReader io.Reader) ([]byte, error) {
	img, format, err := image.Decode(imageReader)
	if err != nil {
		return nil, err
	}

	if format == "webp" {
		data, err := io.ReadAll(imageReader)
		if err != nil {
			return nil, err
		}
		return data, nil
	}

	var buf bytes.Buffer

	err = webp.Encode(&buf, img, nil)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func IsImageFile(name string) bool {
	ext := strings.ToLower(filepath.Ext(name))
	return ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".webp" || ext == ".avif"
}

func ImageMiMeType(data []byte) string {
	// PNG
	if bytes.HasPrefix(data, []byte{0x89, 0x50, 0x4E, 0x47}) {
		return "image/png"
	}

	// JPEG
	if bytes.HasPrefix(data, []byte{0xFF, 0xD8, 0xFF}) {
		return "image/jpeg"
	}

	// GIF
	if bytes.HasPrefix(data, []byte{0x47, 0x49, 0x46, 0x38}) {
		return "image/gif"
	}

	// WebP
	if bytes.HasPrefix(data, []byte{'R', 'I', 'F', 'F'}) && bytes.HasPrefix(data[8:], []byte{'W', 'E', 'B', 'P'}) {
		return "image/webp"
	}

	if len(data) <= 4 {
		return "text/plain; charset=utf-8"
	}

	// AVIF
	if bytes.HasPrefix(data[4:], []byte("ftypavif")) {
		return "image/avif"
	}

	return "image/jpeg"
}

func IsImage(data []byte) bool {
	// PNG
	if bytes.HasPrefix(data, []byte{0x89, 0x50, 0x4E, 0x47}) {
		return true
	}

	// JPEG
	if bytes.HasPrefix(data, []byte{0xFF, 0xD8, 0xFF}) {
		return true
	}

	// GIF
	if bytes.HasPrefix(data, []byte{0x47, 0x49, 0x46, 0x38}) {
		return true
	}

	// WebP
	if bytes.HasPrefix(data, []byte{'R', 'I', 'F', 'F'}) && bytes.HasPrefix(data[8:], []byte{'W', 'E', 'B', 'P'}) {
		return true
	}

	// AVIF
	if bytes.HasPrefix(data[4:], []byte("ftypavif")) {
		return true
	}

	return false
}

func GetImageSize(data io.Reader) (width, height int, err error) {
	img, _, err := image.DecodeConfig(data)
	if err != nil {
		return 0, 0, err
	}
	return img.Width, img.Height, nil
}

func ResizeImageToMinDimensionWebp(imageReader io.Reader, minWidth, minHeight int) ([]byte, error) {
	img, _, err := image.Decode(imageReader)
	if err != nil {
		return nil, err
	}

	bounds := img.Bounds()
	originalWidth := bounds.Dx()
	originalHeight := bounds.Dy()

	if originalWidth <= minWidth || originalHeight <= minHeight {
		var buf bytes.Buffer
		err = webp.Encode(&buf, img, nil)
		if err != nil {
			return nil, err
		}
		return buf.Bytes(), nil
	}

	widthRatio := float64(minWidth) / float64(originalWidth)
	heightRatio := float64(minHeight) / float64(originalHeight)

	ratio := math.Max(widthRatio, heightRatio)

	targetWidth := int(math.Round(float64(originalWidth) * ratio))
	targetHeight := int(math.Round(float64(originalHeight) * ratio))

	dst := image.NewRGBA(image.Rect(0, 0, targetWidth, targetHeight))

	draw.CatmullRom.Scale(dst, dst.Bounds(), img, img.Bounds(), draw.Over, nil)

	var buf bytes.Buffer
	err = webp.Encode(&buf, dst, &webp.Options{Quality: 90})
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
