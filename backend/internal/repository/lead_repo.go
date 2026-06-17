package repository

import (
	"mygo-immigration/backend/internal/model"
	"time"

	"gorm.io/gorm"
)

type LeadRepo struct {
	db *gorm.DB
}

func (r *LeadRepo) FindAll(filter LeadFilter) ([]model.Lead, int64, error) {
	var leads []model.Lead
	var total int64

	q := r.db.Model(&model.Lead{})
	if filter.Status != "" {
		q = q.Where("status = ?", filter.Status)
	}
	if filter.Name != "" {
		q = q.Where("name LIKE ?", "%"+filter.Name+"%")
	}
	if filter.Email != "" {
		q = q.Where("email = ?", filter.Email)
	}
	if filter.InterestedProject != "" {
		q = q.Where("interested_project = ?", filter.InterestedProject)
	}

	if err := q.Session(&gorm.Session{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	q = q.Order("created_at desc, id desc")
	if filter.Page > 0 && filter.PerPage > 0 {
		offset := (filter.Page - 1) * filter.PerPage
		q = q.Offset(offset).Limit(filter.PerPage)
	}

	if err := q.Find(&leads).Error; err != nil {
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
