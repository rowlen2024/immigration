package service

import (
	"encoding/json"
	"errors"
	"testing"

	"mygo-immigration/backend/internal/logging"
	"mygo-immigration/backend/internal/model"
)

type mockHomeConfigRepo struct {
	configs     map[string]json.RawMessage
	createCalls int
	updateCalls int
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
	m.createCalls++
	m.configs[cfg.ConfigKey] = append(json.RawMessage(nil), cfg.ConfigValue...)
	return nil
}

func (m *mockHomeConfigRepo) Update(cfg *model.HomeConfig) error {
	m.updateCalls++
	m.configs[cfg.ConfigKey] = append(json.RawMessage(nil), cfg.ConfigValue...)
	return nil
}

func (m *mockHomeConfigRepo) FindAllConfigValues() ([]json.RawMessage, error) {
	return nil, nil
}

func TestHomeConfigServiceGetAdminSkipsFeaturedDetails(t *testing.T) {
	svc := &HomeConfigService{
		repo: &mockHomeConfigRepo{
			configs: map[string]json.RawMessage{
				"project_showcase":     json.RawMessage(`{"featured_project_ids":[21]}`),
				"case_showcase":        json.RawMessage(`{"featured_case_ids":[12]}`),
				"testimonial_showcase": json.RawMessage(`{"featured_testimonial_ids":[34]}`),
			},
		},
	}

	data, err := svc.GetAdmin()
	if err != nil {
		t.Fatalf("GetAdmin returned error: %v", err)
	}

	raw, err := json.Marshal(data.ProjectShowcase)
	if err != nil {
		t.Fatalf("marshal project showcase: %v", err)
	}
	var projectShowcase map[string]json.RawMessage
	if err := json.Unmarshal(raw, &projectShowcase); err != nil {
		t.Fatalf("unmarshal project showcase: %v", err)
	}
	var projectIDs []uint64
	if err := json.Unmarshal(projectShowcase["featured_project_ids"], &projectIDs); err != nil {
		t.Fatalf("expected featured_project_ids in admin response: %v", err)
	}
	if len(projectIDs) != 1 || projectIDs[0] != 21 {
		t.Fatalf("expected featured project ids to be preserved, got %#v", projectIDs)
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

func TestHomeConfigServiceGetPreservesFeaturedProjectIDOrder(t *testing.T) {
	repo := &mockHomeConfigRepo{configs: map[string]json.RawMessage{
		"project_showcase":     json.RawMessage(`{"featured_project_ids":[3,1,2]}`),
		"case_showcase":        json.RawMessage(`{}`),
		"testimonial_showcase": json.RawMessage(`{}`),
	}}
	svc := &HomeConfigService{
		repo: repo,
		projectRepo: &mockProjectRepo{findByIDsLightFn: func(ids []uint64) ([]model.Project, error) {
			return []model.Project{
				{ID: 1, Slug: "one", Name: "项目一"},
				{ID: 2, Slug: "two", Name: "项目二"},
				{ID: 3, Slug: "three", Name: "项目三"},
			}, nil
		}},
	}

	data, err := svc.Get()
	if err != nil {
		t.Fatalf("Get returned error: %v", err)
	}
	got := data.ProjectShowcase.FeaturedProjects
	if len(got) != 3 || got[0].Slug != "three" || got[1].Slug != "one" || got[2].Slug != "two" {
		t.Fatalf("expected configured project order [three one two], got %#v", got)
	}
}

func TestHomeConfigServiceGetUsesLatestSlugForFeaturedProjectID(t *testing.T) {
	repo := &mockHomeConfigRepo{configs: map[string]json.RawMessage{
		"project_showcase":     json.RawMessage(`{"featured_project_ids":[7]}`),
		"case_showcase":        json.RawMessage(`{}`),
		"testimonial_showcase": json.RawMessage(`{}`),
	}}
	currentSlug := "old-slug"
	svc := &HomeConfigService{
		repo: repo,
		projectRepo: &mockProjectRepo{findByIDsLightFn: func(ids []uint64) ([]model.Project, error) {
			return []model.Project{{ID: 7, Slug: currentSlug}}, nil
		}},
	}

	before, err := svc.Get()
	if err != nil {
		t.Fatalf("Get returned error before slug update: %v", err)
	}
	if got := before.ProjectShowcase.FeaturedProjects[0].Slug; got != "old-slug" {
		t.Fatalf("expected old slug before update, got %q", got)
	}

	currentSlug = "new-slug"
	after, err := svc.Get()
	if err != nil {
		t.Fatalf("Get returned error after slug update: %v", err)
	}
	if got := after.ProjectShowcase.FeaturedProjects[0].Slug; got != "new-slug" {
		t.Fatalf("expected latest slug after update, got %q", got)
	}

	var stored struct {
		FeaturedProjectIDs []uint64 `json:"featured_project_ids"`
	}
	if err := json.Unmarshal(repo.configs["project_showcase"], &stored); err != nil {
		t.Fatalf("unmarshal stored config: %v", err)
	}
	if len(stored.FeaturedProjectIDs) != 1 || stored.FeaturedProjectIDs[0] != 7 {
		t.Fatalf("expected featured project id 7 to remain, got %#v", stored.FeaturedProjectIDs)
	}
}

func TestHomeConfigServiceGetRemovesStaleFeaturedProjectIDs(t *testing.T) {
	logging.Init("error")
	repo := &mockHomeConfigRepo{configs: map[string]json.RawMessage{
		"project_showcase":     json.RawMessage(`{"section_title":"精选项目","featured_project_ids":[3,99,1]}`),
		"case_showcase":        json.RawMessage(`{}`),
		"testimonial_showcase": json.RawMessage(`{}`),
	}}
	svc := &HomeConfigService{
		repo: repo,
		projectRepo: &mockProjectRepo{findByIDsLightFn: func(ids []uint64) ([]model.Project, error) {
			return []model.Project{{ID: 1, Slug: "one"}, {ID: 3, Slug: "three"}}, nil
		}},
	}

	data, err := svc.Get()
	if err != nil {
		t.Fatalf("Get returned error: %v", err)
	}
	got := data.ProjectShowcase.FeaturedProjects
	if len(got) != 2 || got[0].Slug != "three" || got[1].Slug != "one" {
		t.Fatalf("expected stale id to be removed without changing order, got %#v", got)
	}

	var stored struct {
		SectionTitle       string   `json:"section_title"`
		FeaturedProjectIDs []uint64 `json:"featured_project_ids"`
	}
	if err := json.Unmarshal(repo.configs["project_showcase"], &stored); err != nil {
		t.Fatalf("unmarshal stored config: %v", err)
	}
	if stored.SectionTitle != "精选项目" {
		t.Fatalf("expected unrelated fields to be preserved, got %q", stored.SectionTitle)
	}
	if len(stored.FeaturedProjectIDs) != 2 || stored.FeaturedProjectIDs[0] != 3 || stored.FeaturedProjectIDs[1] != 1 {
		t.Fatalf("expected cleaned ids [3 1], got %#v", stored.FeaturedProjectIDs)
	}
}

func TestHomeConfigServiceUpdateRejectsDuplicateFeaturedProjectIDs(t *testing.T) {
	repo := &mockHomeConfigRepo{configs: map[string]json.RawMessage{}}
	svc := &HomeConfigService{repo: repo, projectRepo: &mockProjectRepo{}}

	err := svc.Update(map[string]json.RawMessage{
		"project_showcase": json.RawMessage(`{"featured_project_ids":[1,1]}`),
	})
	if err == nil {
		t.Fatal("expected duplicate featured project ids to be rejected")
	}
	if repo.createCalls != 0 || repo.updateCalls != 0 {
		t.Fatal("expected invalid config not to be persisted")
	}
}

func TestHomeConfigServiceUpdateRequiresFeaturedProjectIDs(t *testing.T) {
	for _, raw := range []json.RawMessage{
		json.RawMessage(`{"section_title":"精选项目"}`),
		json.RawMessage(`{"featured_project_ids":null}`),
	} {
		repo := &mockHomeConfigRepo{configs: map[string]json.RawMessage{}}
		svc := &HomeConfigService{repo: repo, projectRepo: &mockProjectRepo{}}

		err := svc.Update(map[string]json.RawMessage{"project_showcase": raw})
		if err == nil {
			t.Fatalf("expected featured_project_ids to be required for %s", raw)
		}
		if repo.createCalls != 0 || repo.updateCalls != 0 {
			t.Fatal("expected invalid config not to be persisted")
		}
	}
}

func TestHomeConfigServiceUpdateRejectsMissingFeaturedProjectID(t *testing.T) {
	repo := &mockHomeConfigRepo{configs: map[string]json.RawMessage{}}
	svc := &HomeConfigService{
		repo: repo,
		projectRepo: &mockProjectRepo{findByIDsLightFn: func(ids []uint64) ([]model.Project, error) {
			return []model.Project{{ID: 1}}, nil
		}},
	}

	err := svc.Update(map[string]json.RawMessage{
		"project_showcase": json.RawMessage(`{"featured_project_ids":[1,2]}`),
	})
	if err == nil {
		t.Fatal("expected missing featured project id to be rejected")
	}
	if repo.createCalls != 0 || repo.updateCalls != 0 {
		t.Fatal("expected invalid config not to be persisted")
	}
}

func TestHomeConfigServiceUpdatePersistsOnlyProjectShowcaseConfig(t *testing.T) {
	repo := &mockHomeConfigRepo{configs: map[string]json.RawMessage{
		"project_showcase": json.RawMessage(`{"featured_project_ids":[]}`),
	}}
	svc := &HomeConfigService{
		repo: repo,
		projectRepo: &mockProjectRepo{findByIDsLightFn: func(ids []uint64) ([]model.Project, error) {
			return []model.Project{{ID: 2}, {ID: 1}}, nil
		}},
	}

	err := svc.Update(map[string]json.RawMessage{
		"project_showcase": json.RawMessage(`{"section_title":"精选项目","section_subtitle":"副标题","featured_project_ids":[2,1],"featured_projects":[{"slug":"stale"}],"featured_slugs":["legacy"]}`),
	})
	if err != nil {
		t.Fatalf("Update returned error: %v", err)
	}

	var stored map[string]json.RawMessage
	if err := json.Unmarshal(repo.configs["project_showcase"], &stored); err != nil {
		t.Fatalf("unmarshal stored config: %v", err)
	}
	if _, ok := stored["featured_projects"]; ok {
		t.Fatal("expected derived featured_projects not to be persisted")
	}
	if _, ok := stored["featured_slugs"]; ok {
		t.Fatal("expected legacy featured_slugs not to be persisted")
	}
	var ids []uint64
	if err := json.Unmarshal(stored["featured_project_ids"], &ids); err != nil {
		t.Fatalf("unmarshal featured_project_ids: %v", err)
	}
	if len(ids) != 2 || ids[0] != 2 || ids[1] != 1 {
		t.Fatalf("expected ids [2 1], got %#v", ids)
	}
}

func TestHomeConfigServiceRemoveFeaturedProjectID(t *testing.T) {
	repo := &mockHomeConfigRepo{configs: map[string]json.RawMessage{
		"project_showcase": json.RawMessage(`{"section_title":"精选项目","featured_project_ids":[3,2,1]}`),
	}}
	svc := &HomeConfigService{repo: repo}

	if err := svc.RemoveFeaturedProjectID(2); err != nil {
		t.Fatalf("RemoveFeaturedProjectID returned error: %v", err)
	}

	var stored struct {
		SectionTitle       string   `json:"section_title"`
		FeaturedProjectIDs []uint64 `json:"featured_project_ids"`
	}
	if err := json.Unmarshal(repo.configs["project_showcase"], &stored); err != nil {
		t.Fatalf("unmarshal stored config: %v", err)
	}
	if stored.SectionTitle != "精选项目" {
		t.Fatalf("expected unrelated fields to be preserved, got %q", stored.SectionTitle)
	}
	if len(stored.FeaturedProjectIDs) != 2 || stored.FeaturedProjectIDs[0] != 3 || stored.FeaturedProjectIDs[1] != 1 {
		t.Fatalf("expected ids [3 1], got %#v", stored.FeaturedProjectIDs)
	}
}
