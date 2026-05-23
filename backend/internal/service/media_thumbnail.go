package service

import (
	"fmt"
	"image"
	"image/jpeg"
	"os"
	"path/filepath"
	"strings"

	"github.com/disintegration/imaging"
	_ "golang.org/x/image/webp"
)

type variantSpec struct {
	name string
	maxW int
	maxH int
	mode string // "fit" or "thumb"
}

var variantSpecs = []variantSpec{
	{name: "thumb", maxW: 200, maxH: 200, mode: "thumb"},
	{name: "sm", maxW: 400, maxH: 300, mode: "fit"},
	{name: "md", maxW: 800, maxH: 450, mode: "fit"},
	{name: "lg", maxW: 1920, maxH: 800, mode: "fit"},
}

const compressThreshold = 5 * 1024 * 1024 // 5MB

// GenerateVariants creates JPEG thumbnails from the source image at savePath.
// baseName is the filename without extension (e.g. "1779509379388292500").
// Returns a map of variant name -> relative URL path.
func (s *MediaService) GenerateVariants(savePath, baseName string) (map[string]string, error) {
	img, err := imaging.Open(savePath, imaging.AutoOrientation(true))
	if err != nil {
		return nil, fmt.Errorf("GenerateVariants: open: %w", err)
	}

	variants := make(map[string]string)
	for _, spec := range variantSpecs {
		var processed image.Image
		if spec.mode == "thumb" {
			processed = imaging.Thumbnail(img, spec.maxW, spec.maxH, imaging.Lanczos)
		} else {
			processed = imaging.Fit(img, spec.maxW, spec.maxH, imaging.Lanczos)
		}

		filename := baseName + "_" + spec.name + ".jpg"
		variantPath := filepath.Join("./uploads", filename)

		f, err := os.Create(variantPath)
		if err != nil {
			return nil, fmt.Errorf("GenerateVariants: create %s: %w", spec.name, err)
		}
		if err := jpeg.Encode(f, processed, &jpeg.Options{Quality: 80}); err != nil {
			f.Close()
			return nil, fmt.Errorf("GenerateVariants: encode %s: %w", spec.name, err)
		}
		f.Close()

		variants[spec.name] = "/uploads/" + filename
	}
	return variants, nil
}

// CompressIfLarge re-encodes images over 5MB at lower quality to reduce file size.
// Returns the file path to use (may differ from savePath if format changed) and the new file size.
func (s *MediaService) CompressIfLarge(savePath string, fileSize int64) (string, int64, error) {
	if fileSize <= compressThreshold {
		return savePath, fileSize, nil
	}

	img, err := imaging.Open(savePath, imaging.AutoOrientation(true))
	if err != nil {
		return "", 0, fmt.Errorf("CompressIfLarge: open: %w", err)
	}

	ext := strings.ToLower(filepath.Ext(savePath))
	basePath := savePath[:len(savePath)-len(ext)]
	newPath := basePath + ".jpg"

	f, err := os.Create(newPath)
	if err != nil {
		return "", 0, fmt.Errorf("CompressIfLarge: create: %w", err)
	}
	defer f.Close()

	if err := jpeg.Encode(f, img, &jpeg.Options{Quality: 80}); err != nil {
		return "", 0, fmt.Errorf("CompressIfLarge: encode: %w", err)
	}

	if newPath != savePath {
		os.Remove(savePath)
	}

	info, err := os.Stat(newPath)
	if err != nil {
		return newPath, 0, nil
	}

	return newPath, info.Size(), nil
}
