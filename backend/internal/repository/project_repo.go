package repository

import (
	"mygo-immigration/backend/internal/model"
	"time"

	"gorm.io/gorm"
)

type ProjectRepo struct {
	db *gorm.DB
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
	err := r.db.
		Preload("Requirements", func(db *gorm.DB) *gorm.DB { return db.Order("sort_order asc") }).
		Preload("CostItems", func(db *gorm.DB) *gorm.DB { return db.Order("sort_order asc") }).
		Preload("TimelinePhases", func(db *gorm.DB) *gorm.DB { return db.Order("sort_order asc") }).
		Preload("Milestones", func(db *gorm.DB) *gorm.DB { return db.Order("sort_order asc") }).
		Preload("FAQs", func(db *gorm.DB) *gorm.DB { return db.Order("sort_order asc") }).
		Preload("Cases", func(db *gorm.DB) *gorm.DB { return db.Order("sort_order asc") }).
		Preload("News").
		Preload("CompareConfig").
		Preload("Advantages", func(db *gorm.DB) *gorm.DB { return db.Order("sort_order asc") }).
		Preload("Testimonials", func(db *gorm.DB) *gorm.DB { return db.Order("sort_order asc") }).
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

// FindAllWithoutPagination returns all projects matching filters, no pagination.
func (r *ProjectRepo) FindAllWithoutPagination(search, status string) ([]model.Project, error) {
	var projects []model.Project
	q := r.db.Model(&model.Project{})
	if search != "" {
		q = q.Where("name LIKE ?", "%"+search+"%")
	}
	if status != "" {
		q = q.Where("status = ?", status)
	}
	err := q.Order("sort_order asc").Find(&projects).Error
	return projects, err
}

func (r *ProjectRepo) FindBySlugs(slugs []string) ([]model.Project, error) {
	var projects []model.Project
	err := r.db.
		Where("slug IN ?", slugs).
		Preload("Requirements", func(db *gorm.DB) *gorm.DB { return db.Order("sort_order asc") }).
		Preload("CostItems", func(db *gorm.DB) *gorm.DB { return db.Order("sort_order asc") }).
		Preload("TimelinePhases", func(db *gorm.DB) *gorm.DB { return db.Order("sort_order asc") }).
		Preload("Milestones", func(db *gorm.DB) *gorm.DB { return db.Order("sort_order asc") }).
		Preload("FAQs", func(db *gorm.DB) *gorm.DB { return db.Order("sort_order asc") }).
		Preload("Cases", func(db *gorm.DB) *gorm.DB { return db.Order("sort_order asc") }).
		Preload("Advantages", func(db *gorm.DB) *gorm.DB { return db.Order("sort_order asc") }).
		Find(&projects).Error
	if err != nil {
		return nil, err
	}
	return projects, nil
}

// FindBySlugsLight returns projects by slugs without heavy preloads.
func (r *ProjectRepo) FindBySlugsLight(slugs []string) ([]model.Project, error) {
	var projects []model.Project
	err := r.db.Where("slug IN ?", slugs).Find(&projects).Error
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

// DeleteNewsByProjectID hard-deletes all project_news rows for a project.
func (r *ProjectRepo) DeleteNewsByProjectID(projectID uint64) error {
	return r.db.Where("project_id = ?", projectID).Delete(&model.ProjectNews{}).Error
}

// FindAllCoverImages returns non-empty cover_image values referencing /uploads/ (unscoped).
func (r *ProjectRepo) FindAllCoverImages() ([]string, error) {
	var urls []string
	err := r.db.Unscoped().Model(&model.Project{}).
		Where("cover_image LIKE ?", "%/uploads/%").
		Pluck("cover_image", &urls).Error
	return urls, err
}

// RemoveNews unlinks a news page from a project.
func (r *ProjectRepo) RemoveNews(projectID, pageID uint64) error {
	return r.db.Exec(
		"DELETE FROM project_news WHERE project_id = ? AND page_id = ?",
		projectID, pageID,
	).Error
}
