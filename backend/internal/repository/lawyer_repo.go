package repository

import (
	"mygo-immigration/backend/internal/model"

	"gorm.io/gorm"
)

type LawyerRepo struct {
	db *gorm.DB
}

func (r *LawyerRepo) FindAll() ([]model.Lawyer, error) {
	var items []model.Lawyer
	err := r.db.Order("sort_order ASC, id ASC").Find(&items).Error
	return items, err
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
	return r.db.Delete(&model.Lawyer{}, id).Error
}
