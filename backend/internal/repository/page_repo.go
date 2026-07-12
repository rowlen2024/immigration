package repository

import (
	"mygo-immigration/backend/internal/model"
	"time"

	"gorm.io/gorm"
)

type PageRepo struct {
	db *gorm.DB
}

func (r *PageRepo) FindByID(id uint64) (*model.Page, error) {
	var page model.Page
	err := r.db.First(&page, id).Error
	if err != nil {
		return nil, err
	}
	return &page, nil
}

func (r *PageRepo) FindBySlug(slug string) (*model.Page, error) {
	var page model.Page
	err := r.db.Where("slug = ?", slug).First(&page).Error
	if err != nil {
		return nil, err
	}
	return &page, nil
}

func (r *PageRepo) FindAll(filter PageFilter) ([]model.Page, int64, error) {
	var pages []model.Page
	var total int64

	q := r.db.Model(&model.Page{})
	if filter.PageType != "" {
		q = q.Where("page_type = ?", filter.PageType)
	}
	if filter.Title != "" {
		q = q.Where("title LIKE ?", "%"+filter.Title+"%")
	}
	if filter.Status != "" {
		q = q.Where("status = ?", filter.Status)
	}

	if err := q.Session(&gorm.Session{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	q = q.Order("is_pinned desc, sort_order asc, id desc")
	if filter.Page > 0 && filter.PerPage > 0 {
		offset := (filter.Page - 1) * filter.PerPage
		q = q.Offset(offset).Limit(filter.PerPage)
	}

	if err := q.Find(&pages).Error; err != nil {
		return nil, 0, err
	}
	return pages, total, nil
}

func (r *PageRepo) FindOptions(filter PageFilter) ([]PageOptionRow, int64, error) {
	var pages []PageOptionRow
	var total int64

	q := r.db.Model(&model.Page{})
	if filter.PageType != "" {
		q = q.Where("page_type = ?", filter.PageType)
	}
	if filter.Title != "" {
		q = q.Where("title LIKE ?", "%"+filter.Title+"%")
	}
	if filter.Status != "" {
		q = q.Where("status = ?", filter.Status)
	}

	if err := q.Session(&gorm.Session{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	q = q.Select("id", "slug", "title").Order("is_pinned desc, sort_order asc, id desc")
	if filter.Page > 0 && filter.PerPage > 0 {
		offset := (filter.Page - 1) * filter.PerPage
		q = q.Offset(offset).Limit(filter.PerPage)
	}

	if err := q.Find(&pages).Error; err != nil {
		return nil, 0, err
	}
	return pages, total, nil
}

func (r *PageRepo) FindBySlugPublished(slug string) (*model.Page, error) {
	var page model.Page
	err := r.db.Where("slug = ? AND status = ?", slug, "published").First(&page).Error
	if err != nil {
		return nil, err
	}
	return &page, nil
}

func (r *PageRepo) FindProjectsByPageID(pageID uint64) ([]model.PageProject, error) {
	var projects []model.PageProject
	err := r.db.Model(&model.Project{}).
		Select("projects.id, projects.name, projects.slug").
		Joins("JOIN project_news ON project_news.project_id = projects.id").
		Where("project_news.page_id = ?", pageID).
		Order("projects.id asc").
		Find(&projects).Error
	if err != nil {
		return nil, err
	}
	return projects, nil
}

func (r *PageRepo) FindRelatedBySlug(slug string, limit int) ([]model.Page, error) {
	var current model.Page
	if err := r.db.Where("slug = ? AND status = ? AND page_type = ?", slug, "published", "news").First(&current).Error; err != nil {
		return nil, err
	}
	if limit <= 0 {
		return []model.Page{}, nil
	}

	projectIDs := r.db.Model(&model.ProjectNews{}).
		Select("project_id").
		Where("page_id = ?", current.ID)

	pages := make([]model.Page, 0, limit)
	err := r.db.Model(&model.Page{}).
		Distinct("pages.*").
		Joins("JOIN project_news ON project_news.page_id = pages.id").
		Where("project_news.project_id IN (?)", projectIDs).
		Where("pages.id <> ? AND pages.status = ? AND pages.page_type = ?", current.ID, "published", "news").
		Order("pages.is_pinned desc, pages.created_at desc, pages.id desc").
		Limit(limit).
		Find(&pages).Error
	if err != nil {
		return nil, err
	}

	excludedIDs := []uint64{current.ID}
	for _, page := range pages {
		excludedIDs = append(excludedIDs, page.ID)
	}

	if len(pages) < limit && len(current.Tags) > 0 {
		var tagged []model.Page
		err = r.db.Model(&model.Page{}).
			Joins("JOIN JSON_TABLE(pages.tags, '$[*]' COLUMNS(tag VARCHAR(255) PATH '$')) AS page_tags").
			Where("pages.id NOT IN ?", excludedIDs).
			Where("pages.status = ? AND pages.page_type = ?", "published", "news").
			Where("BINARY page_tags.tag IN ?", current.Tags).
			Group("pages.id").
			Order("COUNT(DISTINCT page_tags.tag) desc, pages.is_pinned desc, pages.created_at desc, pages.id desc").
			Limit(limit - len(pages)).
			Find(&tagged).Error
		if err != nil {
			return nil, err
		}
		pages = append(pages, tagged...)
		for _, page := range tagged {
			excludedIDs = append(excludedIDs, page.ID)
		}
	}

	if len(pages) < limit {
		var fallback []model.Page
		err = r.db.Model(&model.Page{}).
			Where("id NOT IN ?", excludedIDs).
			Where("status = ? AND page_type = ?", "published", "news").
			Order("is_pinned desc, created_at desc, id desc").
			Limit(limit - len(pages)).
			Find(&fallback).Error
		if err != nil {
			return nil, err
		}
		pages = append(pages, fallback...)
	}
	return pages, nil
}

func (r *PageRepo) Create(page *model.Page) error {
	return r.db.Create(page).Error
}

func (r *PageRepo) Update(page *model.Page) error {
	return r.db.Omit("created_at").Save(page).Error
}

func (r *PageRepo) Delete(id uint64) error {
	return r.db.Unscoped().Delete(&model.Page{}, id).Error
}

func (r *PageRepo) Count() (int64, error) {
	return CountByModel[model.Page](r.db)
}

func (r *PageRepo) CountByRange(start, end time.Time) (int64, error) {
	return CountByModelRange[model.Page](r.db, start, end)
}

func (r *PageRepo) FindAllCoverImages() ([]string, error) {
	return PluckUploadsByColumn[model.Page](r.db, "cover_image")
}

func (r *PageRepo) FindAllContents() ([]string, error) {
	return PluckUploadsByColumn[model.Page](r.db, "content")
}

func (r *PageRepo) Search(keyword string) ([]model.Page, error) {
	var pages []model.Page
	err := r.db.
		Where("title LIKE ? OR content LIKE ?", "%"+keyword+"%", "%"+keyword+"%").
		Order("sort_order asc, id desc").
		Find(&pages).Error
	if err != nil {
		return nil, err
	}
	return pages, nil
}
