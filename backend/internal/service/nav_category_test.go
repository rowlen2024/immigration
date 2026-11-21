package service

import (
	"encoding/json"
	"strings"
	"testing"

	"mygo-immigration/backend/internal/dto"
	"mygo-immigration/backend/internal/model"
	"mygo-immigration/backend/internal/repository"
)

type navCategoryFixture struct {
	category model.ProjectCategory
}

func (f *navCategoryFixture) FindByID(uint64) (*model.ProjectCategory, error) {
	return &f.category, nil
}

func (f *navCategoryFixture) FindByIDs([]uint64) ([]model.ProjectCategory, error) {
	return []model.ProjectCategory{f.category}, nil
}

type categoryNavFixture struct {
	repository.NavigationRepository
	items []model.Navigation
}

func (f categoryNavFixture) FindAllActiveByPosition(string) ([]model.Navigation, error) {
	return append([]model.Navigation(nil), f.items...), nil
}

func TestCategoryNavResolvesURLWithoutReplacingShortcuts(t *testing.T) {
	parentID, categoryID := uint64(8), uint64(9)
	shortcut := "/projects/eb5"
	categories := &navCategoryFixture{category: model.ProjectCategory{
		ID: categoryID, Name: "热门项目", Slug: "hot", Status: model.ProjectCategoryPublished,
	}}
	svc := &NavService{
		categoryRepo: categories,
		repo: categoryNavFixture{items: []model.Navigation{
			{ID: parentID, LinkType: "category", CategoryID: &categoryID},
			{ID: 10, ParentID: &parentID, LinkType: "custom", Link: &shortcut},
			{ID: 11, LinkType: "category", CategoryID: &categoryID},
		}},
	}

	assertTree := func(wantLink string, wantRoots int) model.Navigation {
		t.Helper()
		tree, err := svc.GetTree("header")
		if err != nil || len(tree) != wantRoots || len(tree[0].Children) != 1 {
			t.Fatalf("invalid navigation tree: %+v (%v)", tree, err)
		}
		parent, child := tree[0], tree[0].Children[0]
		if child.ID != 10 || child.Link == nil || *child.Link != shortcut {
			t.Fatalf("existing shortcut did not survive the category navigation: %+v", child)
		}
		gotLink := ""
		if parent.Link != nil {
			gotLink = *parent.Link
		}
		if gotLink != wantLink {
			t.Fatalf("wrong category link, expected %q: %+v", wantLink, parent)
		}
		return parent
	}
	parent := assertTree("/projects/category/hot", 2)
	publicJSON, err := json.Marshal(parent)
	if err != nil || strings.Contains(string(publicJSON), `"category_id"`) || strings.Contains(string(publicJSON), `"category":`) {
		t.Fatalf("unexpected public navigation payload: %s (%v)", publicJSON, err)
	}
	admin := dto.NewAdminNavigationItem(parent)
	if admin.CategoryID == nil || *admin.CategoryID != categoryID || admin.Category == nil || admin.Category.Name != "热门项目" {
		t.Fatalf("admin navigation missing selected category: %+v", admin)
	}

	categories.category.Status = model.ProjectCategoryDraft
	assertTree("", 1)
}
