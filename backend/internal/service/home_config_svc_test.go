package service

import (
	"encoding/json"
	"errors"
	"testing"

	"mygo-immigration/backend/internal/model"
)

type mockHomeConfigRepo struct {
	configs map[string]json.RawMessage
}

func (m *mockHomeConfigRepo) FindByKey(key string) (*model.HomeConfig, error) {
	if value, ok := m.configs[key]; ok {
		return &model.HomeConfig{ConfigKey: key, ConfigValue: value}, nil
	}
	return nil, errors.New("not found")
}

func (m *mockHomeConfigRepo) FindAll() ([]model.HomeConfig, error) {
	return nil, nil
}

func (m *mockHomeConfigRepo) Create(cfg *model.HomeConfig) error {
	return nil
}

func (m *mockHomeConfigRepo) Update(cfg *model.HomeConfig) error {
	return nil
}

func (m *mockHomeConfigRepo) FindAllConfigValues() ([]json.RawMessage, error) {
	return nil, nil
}

func TestHomeConfigServiceGetAdminSkipsFeaturedDetails(t *testing.T) {
	svc := &HomeConfigService{
		repo: &mockHomeConfigRepo{
			configs: map[string]json.RawMessage{
				"project_showcase":     json.RawMessage(`{"featured_slugs":["canada-suv"]}`),
				"case_showcase":        json.RawMessage(`{"featured_case_ids":[12]}`),
				"testimonial_showcase": json.RawMessage(`{"featured_testimonial_ids":[34]}`),
			},
		},
	}

	data, err := svc.GetAdmin()
	if err != nil {
		t.Fatalf("GetAdmin returned error: %v", err)
	}

	if got := data.ProjectShowcase.FeaturedSlugs; len(got) != 1 || got[0] != "canada-suv" {
		t.Fatalf("expected featured slugs to be preserved, got %#v", got)
	}
	if got := data.CaseShowcase.FeaturedCaseIDs; len(got) != 1 || got[0] != 12 {
		t.Fatalf("expected featured case ids to be preserved, got %#v", got)
	}
	if got := data.TestimonialShowcase.FeaturedTestimonialIDs; len(got) != 1 || got[0] != 34 {
		t.Fatalf("expected featured testimonial ids to be preserved, got %#v", got)
	}

	if len(data.ProjectShowcase.FeaturedProjects) != 0 {
		t.Fatalf("expected admin config to skip featured project details, got %#v", data.ProjectShowcase.FeaturedProjects)
	}
	if len(data.CaseShowcase.FeaturedCases) != 0 {
		t.Fatalf("expected admin config to skip featured case details, got %#v", data.CaseShowcase.FeaturedCases)
	}
	if len(data.TestimonialShowcase.FeaturedTestimonials) != 0 {
		t.Fatalf("expected admin config to skip featured testimonial details, got %#v", data.TestimonialShowcase.FeaturedTestimonials)
	}
}
