package repository

import (
	"mygo-immigration/backend/internal/model"
	"time"

	"gorm.io/gorm"
)

type LeadRepo struct {
	db *gorm.DB
}

func (r *LeadRepo) FindAll(page, perPage int, status string) ([]model.Lead, int64, error) {
	var leads []model.Lead
	var total int64

	q := r.db.Model(&model.Lead{})
	if status != "" {
		q = q.Where("status = ?", status)
	}

	countQ := q.Session(&gorm.Session{})
	if err := countQ.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * perPage
	err := q.Order("created_at desc").Offset(offset).Limit(perPage).Find(&leads).Error
	if err != nil {
		return nil, 0, err
	}
	return leads, total, nil
}

func (r *LeadRepo) Create(lead *model.Lead) error {
	return r.db.Create(lead).Error
}

func (r *LeadRepo) UpdateStatus(id uint64, status string, notes string) error {
	return r.db.Model(&model.Lead{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status": status,
		"notes":  notes,
	}).Error
}

func (r *LeadRepo) Delete(id uint64) error {
	return r.db.Unscoped().Delete(&model.Lead{}, id).Error
}

func (r *LeadRepo) Count() (int64, error) {
	return CountByModel[model.Lead](r.db)
}

func (r *LeadRepo) CountByStatus(status string) (int64, error) {
	var c int64
	err := r.db.Model(&model.Lead{}).Where("status = ?", status).Count(&c).Error
	return c, err
}

func (r *LeadRepo) CountByRange(start, end time.Time) (int64, error) {
	return CountByModelRange[model.Lead](r.db, start, end)
}
