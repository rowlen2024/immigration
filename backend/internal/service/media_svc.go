package service

import (
	"errors"
	"fmt"

	"mygo-immigration/backend/internal/model"
	"mygo-immigration/backend/internal/repository"
)

// MediaService handles business logic for media file metadata.
type MediaService struct {
	repo *repository.MediaRepo
}

// List returns paginated media entries.
func (s *MediaService) List(page, perPage int, search string) ([]model.Media, int64, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 10
	}

	media, err := s.repo.FindAll(search)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list media: %w", err)
	}

	total := int64(len(media))
	start := (page - 1) * perPage
	if start >= len(media) {
		return []model.Media{}, total, nil
	}
	end := start + perPage
	if end > len(media) {
		end = len(media)
	}
	return media[start:end], total, nil
}

// Upload saves media metadata (placeholder - file upload handled at handler level).
func (s *MediaService) Upload(media *model.Media) (*model.Media, error) {
	if media == nil {
		return nil, errors.New("media is nil")
	}
	if media.Filename == "" {
		return nil, errors.New("media filename is required")
	}
	if media.URL == "" {
		return nil, errors.New("media url is required")
	}
	if err := s.repo.Create(media); err != nil {
		return nil, fmt.Errorf("failed to save media metadata: %w", err)
	}
	return media, nil
}

// Delete removes a media entry by ID.
func (s *MediaService) Delete(id uint64) error {
	if id == 0 {
		return errors.New("media id is required")
	}
	if err := s.repo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete media: %w", err)
	}
	return nil
}
