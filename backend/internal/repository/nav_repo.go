package repository

import (
	"mygo-immigration/backend/internal/model"

	"gorm.io/gorm"
)

type NavRepo struct {
	db *gorm.DB
}

func (r *NavRepo) FindAll() ([]model.Navigation, error) {
	var items []model.Navigation
	err := r.db.Where("deleted_at IS NULL").Order("sort_order asc, id asc").Find(&items).Error
	return items, err
}

func (r *NavRepo) FindAllActive() ([]model.Navigation, error) {
	var items []model.Navigation
	err := r.db.Where("status = 1 AND deleted_at IS NULL").Order("sort_order asc, id asc").Find(&items).Error
	return items, err
}

func (r *NavRepo) FindAllActiveByPosition(position string) ([]model.Navigation, error) {
	var items []model.Navigation
	err := r.db.Where("status = 1 AND deleted_at IS NULL AND display_position IN ?", []string{position, "both"}).
		Order("sort_order asc, id asc").
		Find(&items).Error
	return items, err
}

func (r *NavRepo) FindByID(id uint64) (*model.Navigation, error) {
	var item model.Navigation
	err := r.db.Where("id = ? AND deleted_at IS NULL", id).First(&item).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *NavRepo) Create(nav *model.Navigation) error {
	return r.db.Create(nav).Error
}

func (r *NavRepo) Update(nav *model.Navigation) error {
	return r.db.Omit("created_at").Save(nav).Error
}

func (r *NavRepo) Delete(id uint64) error {
	return r.db.Delete(&model.Navigation{}, id).Error
}

func (r *NavRepo) HasChildren(parentID uint64) (bool, error) {
	var count int64
	err := r.db.Model(&model.Navigation{}).
		Where("parent_id = ? AND deleted_at IS NULL", parentID).
		Count(&count).Error
	return count > 0, err
}

func (r *NavRepo) CountByProjectID(projectID uint64) (int64, error) {
	var count int64
	err := r.db.Model(&model.Navigation{}).
		Where("project_id = ? AND deleted_at IS NULL", projectID).
		Count(&count).Error
	return count, err
}

func (r *NavRepo) CountByPageID(pageID uint64) (int64, error) {
	var count int64
	err := r.db.Model(&model.Navigation{}).
		Where("page_id = ? AND deleted_at IS NULL", pageID).
		Count(&count).Error
	return count, err
}
