package repository

import (
	"mygo-immigration/backend/internal/model"
	"time"

	"gorm.io/gorm"
)

type LawyerRepo struct {
	db *gorm.DB
}

func (r *LawyerRepo) FindAll(filter LawyerFilter) ([]model.Lawyer, int64, error) {
	var items []model.Lawyer
	var total int64

	q := r.db.Model(&model.Lawyer{})
	if filter.Name != "" {
		q = q.Where("name LIKE ?", "%"+filter.Name+"%")
	}

	if err := q.Session(&gorm.Session{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	q = q.Order("sort_order ASC, id ASC")
	if filter.Page > 0 && filter.PerPage > 0 {
		offset := (filter.Page - 1) * filter.PerPage
		q = q.Offset(offset).Limit(filter.PerPage)
	}

	if err := q.Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

func (r *LawyerRepo) FindByID(id uint64) (*model.Lawyer, error) {
	var item model.Lawyer
	err := r.db.First(&item, id).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *LawyerRepo) Create(item *model.Lawyer) error {
	return r.db.Create(item).Error
}

func (r *LawyerRepo) Update(item *model.Lawyer) error {
	return r.db.Omit("created_at").Save(item).Error
}

func (r *LawyerRepo) Delete(id uint64) error {
	return r.db.Unscoped().Delete(&model.Lawyer{}, id).Error
}

func (r *LawyerRepo) FindAllPhotoURLs() ([]string, error) {
	return PluckUploadsByColumn[model.Lawyer](r.db, "photo_url")
}

func (r *LawyerRepo) Count() (int64, error) {
	return CountByModel[model.Lawyer](r.db)
}

func (r *LawyerRepo) CountByRange(start, end time.Time) (int64, error) {
	var c int64
	err := r.db.Model(&model.Lawyer{}).Where("created_at >= ? AND created_at < ?", start, end).Count(&c).Error
	return c, err
}
