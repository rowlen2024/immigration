package repository

import (
	"mygo-immigration/backend/internal/model"

	"gorm.io/gorm"
)

type TestimonialRepo struct {
	db *gorm.DB
}

func (r *TestimonialRepo) FindByProjectID(projectID uint64) ([]model.Testimonial, error) {
	var items []model.Testimonial
	err := r.db.
		Where("project_id = ?", projectID).
		Order("sort_order asc").
		Find(&items).Error
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (r *TestimonialRepo) FindAll() ([]model.Testimonial, error) {
	var items []model.Testimonial
	err := r.db.Order("sort_order asc").Find(&items).Error
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (r *TestimonialRepo) Create(t *model.Testimonial) error {
	return r.db.Create(t).Error
}

func (r *TestimonialRepo) Update(t *model.Testimonial) error {
	return r.db.Omit("created_at").Save(t).Error
}

func (r *TestimonialRepo) HardDelete(id uint64) error {
	return r.db.Unscoped().Delete(&model.Testimonial{}, id).Error
}
