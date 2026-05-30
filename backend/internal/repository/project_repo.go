package repository

import (
	"mygo-immigration/backend/internal/model"
	"time"

	"gorm.io/gorm"
)

type ProjectRepo struct {
	db *gorm.DB
}

// preloadScope defines which associations to preload in a Project query.
type preloadScope int

const (
	preloadCompare preloadScope = iota + 1 // compare page: 7 common associations
	preloadDetail                         // detail page: all 10 associations
)

// withAssociations applies ordered, sorted Preloads appropriate to the scope.
// Centralizing Preloads here ensures new associations are added once, not per query method.
func (r *ProjectRepo) withAssociations(db *gorm.DB, scope preloadScope) *gorm.DB {
	sorted := func(db *gorm.DB) *gorm.DB { return db.Order("sort_order asc") }

	db = db.
		Preload("Requirements", sorted).
		Preload("CostItems", sorted).
		Preload("TimelinePhases", sorted).
		Preload("Milestones", sorted).
		Preload("FAQs", sorted).
		Preload("Cases", sorted).
		Preload("Advantages", sorted)

	if scope == preloadDetail {
		db = db.
			Preload("News").
			Preload("CompareConfig").
			Preload("Testimonials", sorted)
	}
	return db
}

func (r *ProjectRepo) FindByID(id uint64) (*model.Project, error) {
	var project model.Project
	err := r.db.First(&project, id).Error
	if err != nil {
		return nil, err
	}
	return &project, nil
}

func (r *ProjectRepo) FindBySlug(slug string) (*model.Project, error) {
	var project model.Project
	err := r.withAssociations(r.db, preloadDetail).
		Where("slug = ?", slug).
		First(&project).Error
	if err != nil {
		return nil, err
	}
	return &project, nil
}

func (r *ProjectRepo) FindAll(filter ProjectFilter) ([]model.Project, int64, error) {
	var projects []model.Project
	var total int64

	q := r.db.Model(&model.Project{})
	if filter.Name != "" {
		q = q.Where("name LIKE ?", "%"+filter.Name+"%")
	}
	if filter.Status != "" {
		q = q.Where("status = ?", filter.Status)
	}

	if err := q.Session(&gorm.Session{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	q = q.Order("sort_order asc")
	if filter.Page > 0 && filter.PerPage > 0 {
		offset := (filter.Page - 1) * filter.PerPage
		q = q.Offset(offset).Limit(filter.PerPage)
	}

	if err := q.Find(&projects).Error; err != nil {
		return nil, 0, err
	}
	return projects, total, nil
}

func (r *ProjectRepo) FindBySlugs(slugs []string) ([]model.Project, error) {
	var projects []model.Project
	err := r.withAssociations(r.db.Where("slug IN ?", slugs), preloadCompare).
		Find(&projects).Error
	if err != nil {
		return nil, err
	}
	return projects, nil
}

// FindBySlugsLight returns projects by slugs without heavy preloads.
func (r *ProjectRepo) FindBySlugsLight(slugs []string) ([]model.Project, error) {
	var projects []model.Project
	err := r.db.Where("slug IN ?", slugs).Order("sort_order asc").Find(&projects).Error
	if err != nil {
		return nil, err
	}
	return projects, nil
}

func (r *ProjectRepo) Create(project *model.Project) error {
	return r.db.Create(project).Error
}

func (r *ProjectRepo) Update(project *model.Project) error {
	return r.db.Omit("Requirements", "CostItems", "TimelinePhases", "Milestones", "FAQs", "Cases", "Advantages", "created_at").Save(project).Error
}

func (r *ProjectRepo) Delete(id uint64) error {
	return r.db.Unscoped().Delete(&model.Project{}, id).Error
}

func (r *ProjectRepo) Count() (int64, error) {
	return CountByModel[model.Project](r.db)
}

func (r *ProjectRepo) CountByRange(start, end time.Time) (int64, error) {
	return CountByModelRange[model.Project](r.db, start, end)
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

// DeleteNewsByProjectID hard-deletes all project_news rows for a project.
func (r *ProjectRepo) DeleteNewsByProjectID(projectID uint64) error {
	return r.db.Unscoped().Where("project_id = ?", projectID).Delete(&model.ProjectNews{}).Error
}

// FindAllCoverImages returns non-empty cover_image values referencing /uploads/ (unscoped).
func (r *ProjectRepo) FindAllCoverImages() ([]string, error) {
	return PluckUploadsByColumn[model.Project](r.db, "cover_image")
}

// RemoveNews unlinks a news page from a project.
func (r *ProjectRepo) RemoveNews(projectID, pageID uint64) error {
	return r.db.Exec(
		"DELETE FROM project_news WHERE project_id = ? AND page_id = ?",
		projectID, pageID,
	).Error
}
