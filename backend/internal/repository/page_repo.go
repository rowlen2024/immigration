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

	q = q.Order("sort_order asc, id asc")
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
		Order("sort_order asc, id asc").
		Find(&pages).Error
	if err != nil {
		return nil, err
	}
	return pages, nil
}
