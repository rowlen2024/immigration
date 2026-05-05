package service

import (
	"fmt"

	"mygo-immigration/backend/internal/model"
	"mygo-immigration/backend/internal/repository"
)

// SearchResults holds the combined results of a search across FAQs and Pages.
type SearchResults struct {
	FAQs  []model.FAQ  `json:"faqs"`
	Pages []model.Page `json:"pages"`
}

// SearchService performs full-text search across content types.
type SearchService struct {
	faqRepo  repository.FAQRepository
	pageRepo repository.PageRepository
}

// Search searches FAQs and Pages for the given keyword.
func (s *SearchService) Search(keyword string) (*SearchResults, error) {
	faqs, err := s.faqRepo.Search(keyword)
	if err != nil {
		return nil, fmt.Errorf("failed to search faqs: %w", err)
	}

	pages, err := s.pageRepo.Search(keyword)
	if err != nil {
		return nil, fmt.Errorf("failed to search pages: %w", err)
	}

	return &SearchResults{
		FAQs:  faqs,
		Pages: pages,
	}, nil
}
