package repository

import (
	"mygo-immigration/backend/internal/model"

	"gorm.io/gorm"
)

type MediaRepo struct {
	db *gorm.DB
}

func (r *MediaRepo) FindAll(filter MediaFilter) ([]model.Media, int64, error) {
	var media []model.Media
	var total int64

	q := r.db.Model(&model.Media{})
	if filter.Search != "" {
		like := "%" + filter.Search + "%"
		q = q.Where("filename LIKE ? OR original_name LIKE ?", like, like)
	}

	if err := q.Session(&gorm.Session{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	q = q.Order("created_at desc, id desc")
	if filter.Page > 0 && filter.PerPage > 0 {
		offset := (filter.Page - 1) * filter.PerPage
		q = q.Offset(offset).Limit(filter.PerPage)
	}

	if err := q.Find(&media).Error; err != nil {
		return nil, 0, err
	}
	return media, total, nil
}

func (r *MediaRepo) FindByID(id uint64) (*model.Media, error) {
	var m model.Media
	err := r.db.First(&m, id).Error
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *MediaRepo) Create(media *model.Media) error {
	return r.db.Create(media).Error
}

func (r *MediaRepo) Delete(id uint64) error {
	return r.db.Unscoped().Delete(&model.Media{}, id).Error
}

// DeleteByIDPermanently hard-deletes a media record (bypasses soft delete).
func (r *MediaRepo) DeleteByIDPermanently(id uint64) error {
	return r.db.Unscoped().Delete(&model.Media{}, id).Error
}
