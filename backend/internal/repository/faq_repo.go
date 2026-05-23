package repository

import (
	"mygo-immigration/backend/internal/model"

	"gorm.io/gorm"
)

// FAQWithProject holds an FAQ row with joined project columns.
type FAQWithProject struct {
	model.FAQ
	ProjectName string `gorm:"column:project_name"`
	ProjectSlug string `gorm:"column:project_slug"`
}

type FAQRepo struct {
	db *gorm.DB
}

func (r *FAQRepo) FindByID(id uint64) (*model.FAQ, error) {
	var faq model.FAQ
	err := r.db.First(&faq, id).Error
	if err != nil {
		return nil, err
	}
	return &faq, nil
}

func (r *FAQRepo) FindAll(params FAQQueryParams) ([]FAQWithProject, int64, error) {
	var results []FAQWithProject
	var total int64

	q := r.db.Model(&model.FAQ{}).
		Select("faqs.*, projects.name AS project_name, projects.slug AS project_slug").
		Joins("LEFT JOIN projects ON projects.id = faqs.project_id AND projects.deleted_at IS NULL")

	if params.ProjectID != nil {
		q = q.Where("faqs.project_id = ?", *params.ProjectID)
	}
	if params.IsGlobal != nil {
		q = q.Where("faqs.is_global = ?", *params.IsGlobal)
	}
	if params.Search != "" {
		like := "%" + params.Search + "%"
		q = q.Where("faqs.question LIKE ? OR faqs.answer LIKE ?", like, like)
	}

	// Count total matching rows (without LIMIT/OFFSET).
	if err := q.Session(&gorm.Session{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if params.Page < 1 {
		params.Page = 1
	}
	if params.PerPage < 1 || params.PerPage > 100 {
		params.PerPage = 10
	}
	offset := (params.Page - 1) * params.PerPage

	err := q.
		Order("faqs.sort_order asc").
		Offset(offset).
		Limit(params.PerPage).
		Find(&results).Error
	if err != nil {
		return nil, 0, err
	}
	return results, total, nil
}

func (r *FAQRepo) Create(faq *model.FAQ) error {
	return r.db.Create(faq).Error
}

func (r *FAQRepo) Update(faq *model.FAQ) error {
	return r.db.Omit("created_at").Save(faq).Error
}

func (r *FAQRepo) Delete(id uint64) error {
	return r.db.Delete(&model.FAQ{}, id).Error
}

// DeleteByProjectID soft-deletes all FAQs belonging to a project.
func (r *FAQRepo) DeleteByProjectID(projectID uint64) error {
	return r.db.Where("project_id = ?", projectID).Delete(&model.FAQ{}).Error
}

// FindDistinctProjects returns projects that have at least one FAQ.
func (r *FAQRepo) FindDistinctProjects() ([]model.Project, error) {
	var projects []model.Project
	err := r.db.
		Distinct("projects.*").
		Joins("INNER JOIN faqs ON faqs.project_id = projects.id").
		Where("projects.deleted_at IS NULL").
		Order("projects.sort_order asc").
		Find(&projects).Error
	if err != nil {
		return nil, err
	}
	return projects, nil
}

func (r *FAQRepo) Search(keyword string) ([]model.FAQ, error) {
	var faqs []model.FAQ
	err := r.db.
		Where("question LIKE ? OR answer LIKE ?", "%"+keyword+"%", "%"+keyword+"%").
		Order("sort_order asc").
		Find(&faqs).Error
	if err != nil {
		return nil, err
	}
	return faqs, nil
}
