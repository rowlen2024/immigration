package dto

import "mygo-immigration/backend/internal/model"

type NavigationRequest struct {
	model.Navigation
	CategoryID *uint64 `json:"category_id"`
}

type AdminNavigationCategory struct {
	ID     uint64 `json:"id"`
	Name   string `json:"name"`
	Slug   string `json:"slug"`
	Status string `json:"status"`
}

type AdminNavigationItem struct {
	model.Navigation
	CategoryID *uint64                  `json:"category_id,omitempty"`
	Category   *AdminNavigationCategory `json:"category,omitempty"`
	Children   []AdminNavigationItem    `json:"children"`
}

func NewAdminNavigationItem(nav model.Navigation) AdminNavigationItem {
	item := AdminNavigationItem{
		Navigation: nav,
		CategoryID: nav.CategoryID,
		Children:   NewAdminNavigationItems(nav.Children),
	}
	if category := nav.Category; category != nil {
		item.Category = &AdminNavigationCategory{
			ID:     category.ID,
			Name:   category.Name,
			Slug:   category.Slug,
			Status: category.Status,
		}
	}
	return item
}

func NewAdminNavigationItems(items []model.Navigation) []AdminNavigationItem {
	if items == nil {
		return nil
	}
	result := make([]AdminNavigationItem, len(items))
	for i, nav := range items {
		result[i] = NewAdminNavigationItem(nav)
	}
	return result
}
