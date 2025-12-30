package services

import (
	"bytes"
	"image"
	"image/color"
	"image/jpeg"
	"testing"
)

func TestImageProcessor_ProcessUpload(t *testing.T) {
	tests := []struct {
		name          string
		imageWidth    int
		imageHeight   int
		expectResize  bool
		expectedError bool
	}{
		{
			name:          "Small image - no resize needed",
			imageWidth:    800,
			imageHeight:   600,
			expectResize:  false,
			expectedError: false,
		},
		{
			name:          "Large image - should resize",
			imageWidth:    3000,
			imageHeight:   2000,
			expectResize:  true,
			expectedError: false,
		},
		{
			name:          "Square image",
			imageWidth:    1500,
			imageHeight:   1500,
			expectResize:  false,
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a test image
			img := createTestImage(tt.imageWidth, tt.imageHeight)
			buf := new(bytes.Buffer)
			err := jpeg.Encode(buf, img, &jpeg.Options{Quality: 90})
			if err != nil {
				t.Fatalf("Failed to encode test image: %v", err)
			}

			// Process the image
			processor := NewImageProcessor()
			fullImage, thumbnail, err := processor.ProcessUpload(bytes.NewReader(buf.Bytes()), "test.jpg")

			if tt.expectedError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			// Verify full image
			if len(fullImage) == 0 {
				t.Error("Full image is empty")
			}

			// Verify thumbnail
			if len(thumbnail) == 0 {
				t.Error("Thumbnail is empty")
			}

			// Decode and verify dimensions
			fullImg, _, err := image.Decode(bytes.NewReader(fullImage))
			if err != nil {
				t.Fatalf("Failed to decode full image: %v", err)
			}

			bounds := fullImg.Bounds()
			if tt.expectResize {
				// Check that image was resized
				if bounds.Dx() > MaxImageWidth || bounds.Dy() > MaxImageHeight {
					t.Errorf("Image not resized properly: got %dx%d, max %dx%d",
						bounds.Dx(), bounds.Dy(), MaxImageWidth, MaxImageHeight)
				}
			}

			// Verify thumbnail dimensions
			thumbImg, _, err := image.Decode(bytes.NewReader(thumbnail))
			if err != nil {
				t.Fatalf("Failed to decode thumbnail: %v", err)
			}

			thumbBounds := thumbImg.Bounds()
			if thumbBounds.Dx() > ThumbnailSize || thumbBounds.Dy() > ThumbnailSize {
				t.Errorf("Thumbnail too large: got %dx%d, max %dx%d",
					thumbBounds.Dx(), thumbBounds.Dy(), ThumbnailSize, ThumbnailSize)
			}
		})
	}
}

func TestImageProcessor_ResizeImage(t *testing.T) {
	processor := NewImageProcessor()

	tests := []struct {
		name        string
		inputWidth  int
		inputHeight int
		maxWidth    int
		maxHeight   int
		expectWidth int
	}{
		{
			name:        "No resize needed",
			inputWidth:  800,
			inputHeight: 600,
			maxWidth:    1000,
			maxHeight:   1000,
			expectWidth: 800,
		},
		{
			name:        "Width exceeds max",
			inputWidth:  2000,
			inputHeight: 1000,
			maxWidth:    1000,
			maxHeight:   1000,
			expectWidth: 1000,
		},
		{
			name:        "Height exceeds max",
			inputWidth:  1000,
			inputHeight: 2000,
			maxWidth:    1000,
			maxHeight:   1000,
			expectWidth: 500,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			img := createTestImage(tt.inputWidth, tt.inputHeight)
			resized := processor.resizeImage(img, tt.maxWidth, tt.maxHeight)

			bounds := resized.Bounds()
			if bounds.Dx() != tt.expectWidth {
				t.Errorf("Expected width %d, got %d", tt.expectWidth, bounds.Dx())
			}
		})
	}
}

func TestImageProcessor_CompressImage(t *testing.T) {
	processor := NewImageProcessor()
	img := createTestImage(1000, 1000)

	compressed, err := processor.compressImage(img, "jpeg")
	if err != nil {
		t.Fatalf("Failed to compress image: %v", err)
	}

	if len(compressed) == 0 {
		t.Error("Compressed image is empty")
	}

	// Verify it's a valid JPEG
	_, _, err = image.Decode(bytes.NewReader(compressed))
	if err != nil {
		t.Errorf("Compressed image is not valid: %v", err)
	}
}

// Helper function to create a test image
func createTestImage(width, height int) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// Fill with a gradient pattern
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r := uint8((x * 255) / width)
			g := uint8((y * 255) / height)
			b := uint8(128)
			img.Set(x, y, color.RGBA{R: r, G: g, B: b, A: 255})
		}
	}

	return img
}
