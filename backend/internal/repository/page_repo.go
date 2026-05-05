package repository

import (
	"mygo-immigration/backend/internal/model"
	"time"

	"gorm.io/gorm"
)

type PageRepo struct {
	db *gorm.DB
}

func (r *PageRepo) FindBySlug(slug string) (*model.Page, error) {
	var page model.Page
	err := r.db.Where("slug = ?", slug).First(&page).Error
	if err != nil {
		return nil, err
	}
	return &page, nil
}

func (r *PageRepo) FindAll() ([]model.Page, error) {
	var pages []model.Page
	err := r.db.
		Order("sort_order asc").
		Find(&pages).Error
	if err != nil {
		return nil, err
	}
	return pages, nil
}

func (r *PageRepo) FindByProjectID(projectID uint64) ([]model.Page, error) {
	var pages []model.Page
	err := r.db.
		Where("project_id = ?", projectID).
		Order("sort_order asc").
		Find(&pages).Error
	if err != nil {
		return nil, err
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
	return r.db.Delete(&model.Page{}, id).Error
}

func (r *PageRepo) Count() (int64, error) {
	var c int64
	err := r.db.Model(&model.Page{}).Count(&c).Error
	return c, err
}

func (r *PageRepo) CountByRange(start, end time.Time) (int64, error) {
	var c int64
	err := r.db.Model(&model.Page{}).Where("created_at >= ? AND created_at < ?", start, end).Count(&c).Error
	return c, err
}

func (r *PageRepo) Search(keyword string) ([]model.Page, error) {
	var pages []model.Page
	err := r.db.
		Where("title LIKE ? OR content LIKE ?", "%"+keyword+"%", "%"+keyword+"%").
		Order("sort_order asc").
		Find(&pages).Error
	if err != nil {
		return nil, err
	}
	return pages, nil
}
