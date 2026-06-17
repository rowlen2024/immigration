package repository

import (
	"mygo-immigration/backend/internal/model"

	"gorm.io/gorm"
)

type MediaRepo struct {
	db *gorm.DB
}

func (r *MediaRepo) FindAll(search string) ([]model.Media, error) {
	var media []model.Media
	q := r.db.Order("created_at desc, id desc")
	if search != "" {
		like := "%" + search + "%"
		q = q.Where("filename LIKE ? OR original_name LIKE ?", like, like)
	}
	err := q.Find(&media).Error
	if err != nil {
		return nil, err
	}
	return media, nil
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
