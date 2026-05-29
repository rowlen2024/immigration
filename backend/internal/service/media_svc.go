package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"mygo-immigration/backend/internal/model"
	"mygo-immigration/backend/internal/repository"
)

// MediaService handles business logic for media file metadata.
type MediaService struct {
	repo             repository.MediaRepository
	projectRepo      repository.ProjectRepository
	caseRepo         repository.CaseRepository
	pageRepo         repository.PageRepository
	lawyerRepo       repository.LawyerRepository
	testimonialRepo  repository.TestimonialRepository
	homeConfigRepo   repository.HomeConfigRepository
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

// Delete removes a media entry by ID, including its physical files and variants.
func (s *MediaService) Delete(id uint64) error {
	if id == 0 {
		return errors.New("media id is required")
	}

	m, err := s.repo.FindByID(id)
	if err != nil {
		return fmt.Errorf("media not found: %w", err)
	}

	filesToDelete := make(map[string]bool)
	if m.Filename != "" {
		filesToDelete[m.Filename] = true
	}
	if m.Variants != nil {
		var v map[string]string
		if err := json.Unmarshal(m.Variants, &v); err == nil {
			for _, url := range v {
				name := filepath.Base(url)
				if name != "" {
					filesToDelete[name] = true
				}
			}
		}
	}
	for name := range filesToDelete {
		filePath := filepath.Join("./uploads", name)
		if err := os.Remove(filePath); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("failed to delete file %s: %w", name, err)
		}
	}

	if err := s.repo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete media: %w", err)
	}
	return nil
}

var imgSrcRegex = regexp.MustCompile(`<img[^>]*\s+src\s*=\s*["']([^"']*)["'][^>]*>`)

var uploadURLRegex = regexp.MustCompile(`(?:https?://[^/"'\s,}]+)?(/uploads/[^"'\s,}]+)`)

// FindUnused returns non-deleted media records whose URL is not referenced anywhere.
func (s *MediaService) FindUnused() ([]model.Media, error) {
	mediaList, err := s.repo.FindAll("")
	if err != nil {
		return nil, fmt.Errorf("failed to list media: %w", err)
	}
	if len(mediaList) == 0 {
		return nil, nil
	}

	referenced := make(map[string]bool)
	s.collectColumnURLs(referenced)
	s.collectRichTextURLs(referenced)
	s.collectHomeConfigURLs(referenced)

	var unused []model.Media
	for _, m := range mediaList {
		url := m.URL
		if !referenced[url] {
			unused = append(unused, m)
		}
	}
	return unused, nil
}

func (s *MediaService) collectColumnURLs(refs map[string]bool) {
	addURL := func(u string) {
		u = strings.TrimSpace(u)
		if u != "" {
			refs[u] = true
		}
	}
	if urls, err := s.projectRepo.FindAllCoverImages(); err == nil {
		for _, u := range urls {
			addURL(u)
		}
	}
	if urls, err := s.caseRepo.FindAllPhotoURLs(); err == nil {
		for _, u := range urls {
			addURL(u)
		}
	}
	if urls, err := s.pageRepo.FindAllCoverImages(); err == nil {
		for _, u := range urls {
			addURL(u)
		}
	}
	if urls, err := s.lawyerRepo.FindAllPhotoURLs(); err == nil {
		for _, u := range urls {
			addURL(u)
		}
	}
	if urls, err := s.testimonialRepo.FindAllAvatarURLs(); err == nil {
		for _, u := range urls {
			addURL(u)
		}
	}
}

func (s *MediaService) collectRichTextURLs(refs map[string]bool) {
	addMatches := func(content string) {
		matches := imgSrcRegex.FindAllStringSubmatch(content, -1)
		for _, m := range matches {
			if len(m) > 1 {
				u := strings.TrimSpace(m[1])
				if u != "" {
					refs[u] = true
				}
			}
		}
	}
	if contents, err := s.caseRepo.FindAllContents(); err == nil {
		for _, c := range contents {
			addMatches(c)
		}
	}
	if contents, err := s.pageRepo.FindAllContents(); err == nil {
		for _, c := range contents {
			addMatches(c)
		}
	}
}

func (s *MediaService) collectHomeConfigURLs(refs map[string]bool) {
	values, err := s.homeConfigRepo.FindAllConfigValues()
	if err != nil {
		return
	}
	for _, raw := range values {
		s.extractUploadURLs(string(raw), refs)
	}
}

func (s *MediaService) extractUploadURLs(content string, refs map[string]bool) {
	matches := uploadURLRegex.FindAllStringSubmatch(content, -1)
	for _, m := range matches {
		if len(m) > 1 && m[1] != "" {
			refs[m[1]] = true
		}
	}
}

// CleanupUnused deletes the specified media DB records and their physical files.
// Returns the number successfully deleted and a list of failures.
func (s *MediaService) CleanupUnused(ids []uint64) (int, []string, error) {
	if len(ids) == 0 {
		return 0, nil, nil
	}

	var deleted int
	var failed []string

	for _, id := range ids {
		m, err := s.repo.FindByID(id)
		if err != nil {
			failed = append(failed, fmt.Sprintf("id=%d: not found", id))
			continue
		}

		// Delete physical files: original + variants
		filesToDelete := make(map[string]bool)
		if m.Filename != "" {
			filesToDelete[m.Filename] = true
		}
		if m.Variants != nil {
			var v map[string]string
			if err := json.Unmarshal(m.Variants, &v); err == nil {
				for _, url := range v {
					name := filepath.Base(url)
					if name != "" {
						filesToDelete[name] = true
					}
				}
			}
		}
		var delErr error
		for name := range filesToDelete {
			filePath := filepath.Join("./uploads", name)
			if err := os.Remove(filePath); err != nil && !os.IsNotExist(err) {
				delErr = err
				break
			}
		}
		if delErr != nil && !os.IsNotExist(delErr) {
			failed = append(failed, fmt.Sprintf("%s: %v", m.Filename, delErr))
			continue
		}

		if err := s.repo.DeleteByIDPermanently(id); err != nil {
			failed = append(failed, fmt.Sprintf("%s: %v", m.Filename, err))
			continue
		}
		deleted++
	}

	return deleted, failed, nil
}
