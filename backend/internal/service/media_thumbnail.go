package service

import (
	"fmt"
	"image"
	"image/jpeg"
	"os"
	"path"
	"path/filepath"
	"strings"

	"mygo-immigration/backend/internal/dto"

	"github.com/disintegration/imaging"
	_ "golang.org/x/image/webp"
)

// 上传 context 常量 — 业务模块与变体规格的绑定关系
const (
	UploadContextLawyer        = "lawyer"
	UploadContextProject       = "project"
	UploadContextCase          = "case"
	UploadContextTestimonial   = "testimonial"
	UploadContextHomepageSlide = "homepage-slide"
	UploadContextPageCover     = "page-cover"
	UploadContextQRCode        = "qr-code"
	UploadContextOGImage       = "og-image"
	UploadContextFavicon       = "favicon"
	UploadContextGeneral       = "general"
)

// VariantSpec describes a single image variant to generate.
type VariantSpec struct {
	Name string
	MaxW int
	MaxH int
	Mode string // "fit", "fill", or "thumb"
}

// GetContextSpecs returns the variant specs for a given context.
func GetContextSpecs(context string) ([]VariantSpec, bool) {
	raw, ok := contextSpecs[context]
	if !ok {
		return nil, false
	}
	specs := make([]VariantSpec, len(raw))
	for i, s := range raw {
		specs[i] = VariantSpec{Name: s.Name, MaxW: s.MaxW, MaxH: s.MaxH, Mode: s.Mode}
	}
	return specs, true
}

// contextSpecs maps upload contexts to their variant specifications.
// "general" is the default for backward compatibility.
var contextSpecs = map[string][]VariantSpec{
	"lawyer": {
		{Name: "thumb", MaxW: 150, MaxH: 200, Mode: "fill"},
		{Name: "sm", MaxW: 300, MaxH: 400, Mode: "fill"},
		{Name: "md", MaxW: 600, MaxH: 800, Mode: "fill"},
	},
	"project": {
		{Name: "thumb", MaxW: 200, MaxH: 113, Mode: "fill"},
		{Name: "sm", MaxW: 400, MaxH: 225, Mode: "fill"},
		{Name: "md", MaxW: 800, MaxH: 450, Mode: "fill"},
		{Name: "lg", MaxW: 1920, MaxH: 800, Mode: "fill"},
	},
	"case": {
		{Name: "thumb", MaxW: 200, MaxH: 113, Mode: "fill"},
		{Name: "sm", MaxW: 400, MaxH: 225, Mode: "fill"},
		{Name: "md", MaxW: 800, MaxH: 450, Mode: "fill"},
	},
	"testimonial": {
		{Name: "thumb", MaxW: 200, MaxH: 200, Mode: "fill"},
		{Name: "sm", MaxW: 400, MaxH: 400, Mode: "fill"},
		{Name: "md", MaxW: 800, MaxH: 800, Mode: "fill"},
	},
	"homepage-slide": {
		{Name: "thumb", MaxW: 240, MaxH: 100, Mode: "fill"},
		{Name: "sm", MaxW: 480, MaxH: 200, Mode: "fill"},
		{Name: "md", MaxW: 960, MaxH: 400, Mode: "fill"},
		{Name: "lg", MaxW: 1920, MaxH: 800, Mode: "fill"},
	},
	"page-cover": {
		{Name: "thumb", MaxW: 200, MaxH: 133, Mode: "fill"},
		{Name: "sm", MaxW: 360, MaxH: 240, Mode: "fill"},
		{Name: "md", MaxW: 720, MaxH: 480, Mode: "fill"},
	},
	"qr-code": {
		{Name: "thumb", MaxW: 150, MaxH: 150, Mode: "fit"},
		{Name: "sm", MaxW: 300, MaxH: 300, Mode: "fit"},
		{Name: "md", MaxW: 500, MaxH: 500, Mode: "fit"},
	},
	"og-image": {
		{Name: "sm", MaxW: 600, MaxH: 315, Mode: "fill"},
		{Name: "md", MaxW: 1200, MaxH: 630, Mode: "fill"},
	},
	"favicon": {
		{Name: "thumb", MaxW: 32, MaxH: 32, Mode: "fill"},
	},
	"general": {
		{Name: "thumb", MaxW: 200, MaxH: 200, Mode: "thumb"},
		{Name: "sm", MaxW: 400, MaxH: 300, Mode: "fit"},
		{Name: "md", MaxW: 800, MaxH: 450, Mode: "fit"},
		{Name: "lg", MaxW: 1920, MaxH: 800, Mode: "fit"},
	},
}

var validContexts = func() map[string]bool {
	m := make(map[string]bool, len(contextSpecs))
	for k := range contextSpecs {
		m[k] = true
	}
	return m
}()

const compressThreshold = 5 * 1024 * 1024 // 5MB

// GenerateVariants creates JPEG variants using the "general" context (backward compatible).
func (s *MediaService) GenerateVariants(savePath, baseName string) (map[string]string, error) {
	return GenerateVariantsForContext(savePath, baseName, "general")
}

// GenerateVariantsForContext creates JPEG variants according to the context's spec.
func GenerateVariantsForContext(savePath, baseName, context string) (map[string]string, error) {
	specs, ok := contextSpecs[context]
	if !ok {
		specs = contextSpecs["general"]
	}

	img, err := imaging.Open(savePath, imaging.AutoOrientation(true))
	if err != nil {
		return nil, fmt.Errorf("GenerateVariantsForContext: open: %w", err)
	}

	variants := make(map[string]string)
	for _, spec := range specs {
		var processed image.Image
		switch spec.Mode {
		case "thumb":
			processed = imaging.Thumbnail(img, spec.MaxW, spec.MaxH, imaging.Lanczos)
		case "fill":
			processed = imaging.Fill(img, spec.MaxW, spec.MaxH, imaging.Center, imaging.Lanczos)
		default:
			processed = imaging.Fit(img, spec.MaxW, spec.MaxH, imaging.Lanczos)
		}

		filename := baseName + "_" + spec.Name + ".jpg"
		variantPath := filepath.Join("./uploads", filename)

		f, err := os.Create(variantPath)
		if err != nil {
			return nil, fmt.Errorf("GenerateVariantsForContext: create %s: %w", spec.Name, err)
		}
		if err := jpeg.Encode(f, processed, &jpeg.Options{Quality: 80}); err != nil {
			f.Close()
			return nil, fmt.Errorf("GenerateVariantsForContext: encode %s: %w", spec.Name, err)
		}
		f.Close()

		variants[spec.Name] = "/uploads/" + filename
	}
	return variants, nil
}

// IsValidContext returns whether the context string is a known upload context.
func IsValidContext(context string) bool {
	if context == "" {
		return true // empty means general
	}
	return validContexts[context]
}

// CompressIfLarge re-encodes images over 5MB at lower quality to reduce file size.
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

// ResolveImageVariants 根据原图 URL 和 context 计算出所有可用变体的 URL 和宽度。
// 纯计算函数，不访问数据库，不检查文件是否存在。
func ResolveImageVariants(baseURL string, context string) map[string]dto.ImageVariantInfo {
	specs, ok := contextSpecs[context]
	if !ok {
		specs = contextSpecs[UploadContextGeneral]
	}

	base := strings.TrimSuffix(baseURL, path.Ext(baseURL))
	dir := path.Dir(baseURL)
	name := path.Base(base)

	result := make(map[string]dto.ImageVariantInfo, len(specs))
	for _, spec := range specs {
		result[spec.Name] = dto.ImageVariantInfo{
			URL:   path.Join(dir, name+"_"+spec.Name+".jpg"),
			Width: spec.MaxW,
		}
	}
	return result
}
