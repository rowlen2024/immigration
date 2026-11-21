package service

import (
	"errors"
	"testing"

	"mygo-immigration/backend/internal/dto"
	"mygo-immigration/backend/internal/model"
)

func TestProjectCategorySlugRemainsFrozenAfterUnpublishing(t *testing.T) {
	category := model.ProjectCategory{Status: model.ProjectCategoryDraft}
	req := dto.SaveProjectCategoryRequest{Name: "购房移民", Slug: "property"}
	if err := applyProjectCategoryRequest(&category, req); err != nil {
		t.Fatal(err)
	}
	req.Slug = "real-estate"
	req.Status = model.ProjectCategoryPublished
	if err := applyProjectCategoryRequest(&category, req); err != nil {
		t.Fatal(err)
	}
	if category.PublishedAt == nil {
		t.Fatal("first publication time was not saved")
	}
	firstPublishedAt := *category.PublishedAt
	req.Status = model.ProjectCategoryDraft
	if err := applyProjectCategoryRequest(&category, req); err != nil {
		t.Fatal(err)
	}
	req.Slug = "renamed"
	if err := applyProjectCategoryRequest(&category, req); !errors.Is(err, ErrProjectCategorySlugFrozen) {
		t.Fatalf("rename after unpublishing: %v, want ErrProjectCategorySlugFrozen", err)
	}
	if category.Slug != "real-estate" || !category.PublishedAt.Equal(firstPublishedAt) {
		t.Fatalf("rejected update changed published slug or time: %+v", category)
	}
}
