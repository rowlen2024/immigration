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
		Preload("News").
		Preload("CompareConfig").
		Preload("Advantages").
		Preload("Testimonials").
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

	// Clone for count to avoid mutation
	countQ := q.Session(&gorm.Session{})
	if err := countQ.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * perPage
	err := q.Order("sort_order asc").Offset(offset).Limit(perPage).Find(&projects).Error
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
		Preload("Advantages").
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
	return r.db.Omit("Requirements", "CostItems", "TimelinePhases", "Milestones", "FAQs", "Cases", "Advantages").Save(project).Error
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

// FindNews returns news pages linked to a project via project_news.
func (r *ProjectRepo) FindNews(projectID uint64) ([]model.Page, error) {
	var news []model.Page
	err := r.db.
		Joins("JOIN project_news ON project_news.page_id = pages.id").
		Where("project_news.project_id = ?", projectID).
		Where("pages.deleted_at IS NULL").
		Order("project_news.created_at DESC").
		Find(&news).Error
	return news, err
}

// AddNews links news pages to a project.
func (r *ProjectRepo) AddNews(projectID uint64, pageIDs []uint64) error {
	for _, pageID := range pageIDs {
		err := r.db.Exec(
			"INSERT IGNORE INTO project_news (project_id, page_id) VALUES (?, ?)",
			projectID, pageID,
		).Error
		if err != nil {
			return err
		}
	}
	return nil
}

// RemoveNews unlinks a news page from a project.
func (r *ProjectRepo) RemoveNews(projectID, pageID uint64) error {
	return r.db.Exec(
		"DELETE FROM project_news WHERE project_id = ? AND page_id = ?",
		projectID, pageID,
	).Error
}
