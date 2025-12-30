package services

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"path/filepath"

	"github.com/nomdb/backend/internal/logger"
	"golang.org/x/image/draw"
)

const (
	// MaxImageWidth is the maximum width for full-size images
	MaxImageWidth = 1920
	// MaxImageHeight is the maximum height for full-size images
	MaxImageHeight = 1920
	// ThumbnailSize is the size for thumbnail images
	ThumbnailSize = 200
	// JPEGQuality is the quality setting for JPEG compression
	JPEGQuality = 85
)

// ImageProcessor handles image processing operations
type ImageProcessor struct{}

// NewImageProcessor creates a new image processor
func NewImageProcessor() *ImageProcessor {
	return &ImageProcessor{}
}

// ProcessUpload processes an uploaded image: resize, compress, and generate thumbnail
func (ip *ImageProcessor) ProcessUpload(file io.Reader, filename string) (fullImage []byte, thumbnail []byte, err error) {
	// Decode the image
	img, format, err := image.Decode(file)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to decode image: %w", err)
	}

	logger.Debug("Processing image: format=%s, size=%dx%d", format, img.Bounds().Dx(), img.Bounds().Dy())

	// Resize full image if needed
	resizedImg := ip.resizeImage(img, MaxImageWidth, MaxImageHeight)

	// Compress full image
	fullImage, err = ip.compressImage(resizedImg, format)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to compress image: %w", err)
	}

	// Generate thumbnail
	thumbnailImg := ip.resizeImage(img, ThumbnailSize, ThumbnailSize)
	thumbnail, err = ip.compressImage(thumbnailImg, format)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to compress thumbnail: %w", err)
	}

	logger.Debug("Image processed: full=%d bytes, thumbnail=%d bytes", len(fullImage), len(thumbnail))

	return fullImage, thumbnail, nil
}

// resizeImage resizes an image to fit within maxWidth and maxHeight while maintaining aspect ratio
func (ip *ImageProcessor) resizeImage(img image.Image, maxWidth, maxHeight int) image.Image {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	// Check if resizing is needed
	if width <= maxWidth && height <= maxHeight {
		return img
	}

	// Calculate new dimensions maintaining aspect ratio
	ratio := float64(width) / float64(height)
	newWidth := maxWidth
	newHeight := int(float64(newWidth) / ratio)

	if newHeight > maxHeight {
		newHeight = maxHeight
		newWidth = int(float64(newHeight) * ratio)
	}

	logger.Debug("Resizing image from %dx%d to %dx%d", width, height, newWidth, newHeight)

	// Create new image
	dst := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))

	// Use high-quality scaling (bilinear interpolation)
	draw.BiLinear.Scale(dst, dst.Bounds(), img, img.Bounds(), draw.Over, nil)

	return dst
}

// compressImage compresses an image to JPEG or PNG format
func (ip *ImageProcessor) compressImage(img image.Image, format string) ([]byte, error) {
	var buf bytes.Buffer

	switch format {
	case "jpeg", "jpg":
		err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: JPEGQuality})
		if err != nil {
			return nil, err
		}
	case "png":
		// PNG doesn't have quality settings, but we can use default compression
		encoder := png.Encoder{CompressionLevel: png.DefaultCompression}
		err := encoder.Encode(&buf, img)
		if err != nil {
			return nil, err
		}
	default:
		// Default to JPEG for unknown formats
		err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: JPEGQuality})
		if err != nil {
			return nil, err
		}
	}

	return buf.Bytes(), nil
}

// SaveImage saves image bytes to a file
func (ip *ImageProcessor) SaveImage(data []byte, filepath string) error {
	// Ensure directory exists
	dir := filepath[:len(filepath)-len(filepath[len(filepath)-1:])]
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Write file
	if err := os.WriteFile(filepath, data, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// GetThumbnailPath generates the thumbnail path from the original path
func GetThumbnailPath(originalPath string) string {
	dir := filepath.Dir(originalPath)
	filename := filepath.Base(originalPath)
	ext := filepath.Ext(filename)
	nameWithoutExt := filename[:len(filename)-len(ext)]

	return filepath.Join(dir, "thumbnails", nameWithoutExt+"_thumb"+ext)
}
