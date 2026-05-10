package repository

import (
	"mygo-immigration/backend/internal/model"
	"time"

	"gorm.io/gorm"
)

type ProjectRepo struct {
	db *gorm.DB
}

func (r *ProjectRepo) FindBySlug(slug string) (*model.Project, error) {
	var project model.Project
	err := r.db.
		Preload("Requirements").
		Preload("CostItems").
		Preload("TimelinePhases").
		Preload("Milestones").
		Preload("FAQs").
		Preload("Cases").
		Where("slug = ?", slug).
		First(&project).Error
	if err != nil {
		return nil, err
	}
	return &project, nil
}

func (r *ProjectRepo) FindAll(page, perPage int, search, status string) ([]model.Project, int64, error) {
	var projects []model.Project
	var total int64

	q := r.db.Model(&model.Project{})
	if search != "" {
		q = q.Where("name LIKE ?", "%"+search+"%")
	}
	if status != "" {
		q = q.Where("status = ?", status)
	}

	q.Count(&total)

	offset := (page - 1) * perPage
	err := q.Session(&gorm.Session{}).Order("sort_order asc").Offset(offset).Limit(perPage).Find(&projects).Error
	if err != nil {
		return nil, 0, err
	}
	return projects, total, nil
}

func (r *ProjectRepo) FindBySlugs(slugs []string) ([]model.Project, error) {
	var projects []model.Project
	err := r.db.
		Where("slug IN ?", slugs).
		Preload("Requirements").
		Preload("CostItems").
		Preload("TimelinePhases").
		Preload("Milestones").
		Preload("FAQs").
		Preload("Cases").
		Find(&projects).Error
	if err != nil {
		return nil, err
	}
	return projects, nil
}

func (r *ProjectRepo) Create(project *model.Project) error {
	return r.db.Create(project).Error
}

func (r *ProjectRepo) Update(project *model.Project) error {
	return r.db.Omit("Requirements", "CostItems", "TimelinePhases", "Milestones", "FAQs", "Cases").Save(project).Error
}

func (r *ProjectRepo) Delete(id uint64) error {
	return r.db.Delete(&model.Project{}, id).Error
}

func (r *ProjectRepo) Count() (int64, error) {
	var c int64
	err := r.db.Model(&model.Project{}).Count(&c).Error
	return c, err
}

func (r *ProjectRepo) CountByRange(start, end time.Time) (int64, error) {
	var c int64
	err := r.db.Model(&model.Project{}).Where("created_at >= ? AND created_at < ?", start, end).Count(&c).Error
	return c, err
}
