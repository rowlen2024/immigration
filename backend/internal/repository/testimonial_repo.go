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

func (r *TestimonialRepo) FindByIDs(ids []uint64) ([]model.Testimonial, error) {
	var items []model.Testimonial
	err := r.db.Where("id IN ?", ids).Find(&items).Error
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (r *TestimonialRepo) FindAll(filter TestimonialFilter) ([]model.Testimonial, int64, error) {
	var items []model.Testimonial
	var total int64

	q := r.db.Model(&model.Testimonial{})
	if filter.ProjectID != nil {
		q = q.Where("project_id = ?", *filter.ProjectID)
	}
	if filter.Nickname != "" {
		q = q.Where("nickname LIKE ?", "%"+filter.Nickname+"%")
	}
	if filter.Rating != nil {
		q = q.Where("rating = ?", *filter.Rating)
	}

	if err := q.Session(&gorm.Session{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	q = q.Order("sort_order asc, id desc")
	if filter.Page > 0 && filter.PerPage > 0 {
		offset := (filter.Page - 1) * filter.PerPage
		q = q.Offset(offset).Limit(filter.PerPage)
	}

	if err := q.Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

func (r *TestimonialRepo) FindOptions(filter TestimonialFilter) ([]TestimonialOptionRow, int64, error) {
	var items []TestimonialOptionRow
	var total int64

	q := r.db.Model(&model.Testimonial{})
	if filter.ProjectID != nil {
		q = q.Where("project_id = ?", *filter.ProjectID)
	}
	if filter.Nickname != "" {
		q = q.Where("nickname LIKE ?", "%"+filter.Nickname+"%")
	}
	if filter.Rating != nil {
		q = q.Where("rating = ?", *filter.Rating)
	}

	if err := q.Session(&gorm.Session{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	q = q.Select("id", "nickname").Order("sort_order asc, id desc")
	if filter.Page > 0 && filter.PerPage > 0 {
		offset := (filter.Page - 1) * filter.PerPage
		q = q.Offset(offset).Limit(filter.PerPage)
	}

	if err := q.Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

func (r *TestimonialRepo) Create(t *model.Testimonial) error {
	return r.db.Create(t).Error
}

func (r *TestimonialRepo) Update(t *model.Testimonial) error {
	return r.db.Omit("created_at").Save(t).Error
}

func (r *TestimonialRepo) Delete(id uint64) error {
	return r.db.Unscoped().Delete(&model.Testimonial{}, id).Error
}

// DeleteByProjectID deletes all testimonials belonging to a project.
func (r *TestimonialRepo) DeleteByProjectID(projectID uint64) error {
	return r.db.Unscoped().Where("project_id = ?", projectID).Delete(&model.Testimonial{}).Error
}

// FindAllAvatarURLs returns non-empty avatar_url values referencing /uploads/ (unscoped).
func (r *TestimonialRepo) FindAllAvatarURLs() ([]string, error) {
	var urls []string
	err := r.db.Unscoped().Model(&model.Testimonial{}).
		Where("avatar_url LIKE ?", "%/uploads/%").
		Pluck("avatar_url", &urls).Error
	return urls, err
}
