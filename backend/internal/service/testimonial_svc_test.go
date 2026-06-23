package service

import (
	"testing"

	"mygo-immigration/backend/internal/dto"
	"mygo-immigration/backend/internal/model"
	"mygo-immigration/backend/internal/repository"
)

type mockTestimonialRepo struct {
	findOptionsFn func(filter repository.TestimonialFilter) ([]repository.TestimonialOptionRow, int64, error)
}

func (m *mockTestimonialRepo) FindByID(id uint64) (*model.Testimonial, error) {
	return nil, nil
}

func (m *mockTestimonialRepo) FindByIDs(ids []uint64) ([]model.Testimonial, error) {
	return nil, nil
}

func (m *mockTestimonialRepo) FindAll(filter repository.TestimonialFilter) ([]model.Testimonial, int64, error) {
	return nil, 0, nil
}

func (m *mockTestimonialRepo) FindOptions(filter repository.TestimonialFilter) ([]repository.TestimonialOptionRow, int64, error) {
	if m.findOptionsFn != nil {
		return m.findOptionsFn(filter)
	}
	return nil, 0, nil
}

func (m *mockTestimonialRepo) Create(t *model.Testimonial) error {
	return nil
}

func (m *mockTestimonialRepo) Update(t *model.Testimonial) error {
	return nil
}

func (m *mockTestimonialRepo) Delete(id uint64) error {
	return nil
}

func (m *mockTestimonialRepo) DeleteByProjectID(projectID uint64) error {
	return nil
}

func (m *mockTestimonialRepo) FindAllAvatarURLs() ([]string, error) {
	return nil, nil
}

func TestTestimonialService_Options(t *testing.T) {
	repo := &mockTestimonialRepo{
		findOptionsFn: func(filter repository.TestimonialFilter) ([]repository.TestimonialOptionRow, int64, error) {
			if filter.Nickname != "Li" {
				t.Fatalf("expected nickname filter Li, got %q", filter.Nickname)
			}
			if filter.Page != 3 || filter.PerPage != 30 {
				t.Fatalf("expected pagination 3/30, got %d/%d", filter.Page, filter.PerPage)
			}
			return []repository.TestimonialOptionRow{{ID: 11, Nickname: "Li"}}, 1, nil
		},
	}
	svc := &TestimonialService{repo: repo}

	items, total, err := svc.Options(dto.TestimonialOptionRequest{
		OptionPaginationRequest: dto.OptionPaginationRequest{Page: 3, PerPage: 30},
		Nickname:                "Li",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if total != 1 || len(items) != 1 || items[0].Nickname != "Li" {
		t.Fatalf("unexpected options result: total=%d items=%v", total, items)
	}
}
