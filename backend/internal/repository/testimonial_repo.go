package repository

import (
	"mygo-immigration/backend/internal/model"

	"gorm.io/gorm"
)

type TestimonialRepo struct {
	db *gorm.DB
}

func (r *TestimonialRepo) FindByID(id uint64) (*model.Testimonial, error) {
	var t model.Testimonial
	err := r.db.First(&t, id).Error
	if err != nil {
		return nil, err
	}
	return &t, nil
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

func (r *TestimonialRepo) FindByIDs(ids []uint64) ([]model.Testimonial, error) {
	var items []model.Testimonial
	err := r.db.Where("id IN ?", ids).Find(&items).Error
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

// DeleteByProjectID soft-deletes all testimonials belonging to a project.
func (r *TestimonialRepo) DeleteByProjectID(projectID uint64) error {
	return r.db.Where("project_id = ?", projectID).Delete(&model.Testimonial{}).Error
}

func (r *TestimonialRepo) HardDelete(id uint64) error {
	return r.db.Unscoped().Delete(&model.Testimonial{}, id).Error
}

// FindAllAvatarURLs returns non-empty avatar_url values referencing /uploads/ (unscoped).
func (r *TestimonialRepo) FindAllAvatarURLs() ([]string, error) {
	var urls []string
	err := r.db.Unscoped().Model(&model.Testimonial{}).
		Where("avatar_url LIKE ?", "%/uploads/%").
		Pluck("avatar_url", &urls).Error
	return urls, err
}
