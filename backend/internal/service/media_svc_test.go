package service

import (
	"errors"
	"testing"

	"mygo-immigration/backend/internal/dto"
	"mygo-immigration/backend/internal/model"
	"mygo-immigration/backend/internal/repository"
)

type mockMediaRepo struct {
	findAllFn func(filter repository.MediaFilter) ([]model.Media, int64, error)
}

func (m *mockMediaRepo) FindAll(filter repository.MediaFilter) ([]model.Media, int64, error) {
	if m.findAllFn != nil {
		return m.findAllFn(filter)
	}
	return nil, 0, nil
}

func (m *mockMediaRepo) FindByID(id uint64) (*model.Media, error) { return nil, nil }
func (m *mockMediaRepo) Create(media *model.Media) error          { return nil }
func (m *mockMediaRepo) Delete(id uint64) error                   { return nil }
func (m *mockMediaRepo) DeleteByIDPermanently(id uint64) error    { return nil }

func TestMedia_List_ForwardsFilter(t *testing.T) {
	repo := &mockMediaRepo{
		findAllFn: func(filter repository.MediaFilter) ([]model.Media, int64, error) {
			if filter.Search != "hero" {
				t.Errorf("expected Search='hero', got %q", filter.Search)
			}
			if filter.Page != 2 {
				t.Errorf("expected Page=2, got %d", filter.Page)
			}
			if filter.PerPage != 12 {
				t.Errorf("expected PerPage=12, got %d", filter.PerPage)
			}
			return []model.Media{{ID: 2, Filename: "hero.jpg"}}, 25, nil
		},
	}
	svc := &MediaService{repo: repo}

	items, total, err := svc.List(dto.MediaListRequest{
		PaginationRequest: dto.PaginationRequest{Page: 2, PerPage: 12},
		Search:            "hero",
	})
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if total != 25 {
		t.Errorf("expected total 25, got %d", total)
	}
	if len(items) != 1 || items[0].ID != 2 {
		t.Fatalf("expected media item ID 2, got %#v", items)
	}
}

func TestMedia_List_Error(t *testing.T) {
	repo := &mockMediaRepo{
		findAllFn: func(filter repository.MediaFilter) ([]model.Media, int64, error) {
			return nil, 0, errors.New("db error")
		},
	}
	svc := &MediaService{repo: repo}

	_, _, err := svc.List(dto.MediaListRequest{})
	if err == nil {
		t.Fatal("expected error from repo")
	}
}
