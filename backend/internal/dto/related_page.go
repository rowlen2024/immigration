package dto

import (
	"time"

	"mygo-immigration/backend/internal/model"
)

type RelatedPage struct {
	ID                 uint64                            `json:"id"`
	Title              string                            `json:"title"`
	Slug               string                            `json:"slug"`
	CoverImage         string                            `json:"cover_image"`
	CoverImageVariants map[string]model.ImageVariantInfo `json:"cover_image_variants,omitempty"`
	CreatedAt          time.Time                         `json:"created_at"`
}
