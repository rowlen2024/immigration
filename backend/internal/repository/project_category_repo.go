package repository

import (
	"time"

	"mygo-immigration/backend/internal/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ProjectCategoryRepo struct {
	db *gorm.DB
}

type ProjectCategoryFilter struct {
	Name    string
	Status  string
	Page    int
	PerPage int
}

type ProjectCategoryAdminRow struct {
	model.ProjectCategory `gorm:"embedded"`
	ProjectCount          int64 `json:"project_count"`
}

type ProjectCategoryPublicRow struct {
	ID           uint64
	Name         string
	Slug         string
	ProjectCount int64
	UpdatedAt    time.Time
}

type ProjectCategoryProjectRow struct {
	ID        uint64 `json:"id"`
	Name      string `json:"name"`
	Slug      string `json:"slug"`
	Status    int8   `json:"status"`
	SortOrder int32  `json:"sort_order"`
}

type ProjectCategoryAssignment struct {
	ProjectID uint64
	ID        uint64
	Name      string
	Slug      string
	Status    string
}

func (r *ProjectCategoryRepo) FindByID(id uint64) (*model.ProjectCategory, error) {
	var category model.ProjectCategory
	if err := r.db.First(&category, id).Error; err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *ProjectCategoryRepo) FindByIDs(ids []uint64) ([]model.ProjectCategory, error) {
	categories := make([]model.ProjectCategory, 0)
	if len(ids) == 0 {
		return categories, nil
	}
	err := r.db.Where("id IN ?", ids).Order("id").Find(&categories).Error
	return categories, err
}

func (r *ProjectCategoryRepo) FindPublishedBySlug(slug string) (*model.ProjectCategory, error) {
	var category model.ProjectCategory
	err := r.db.Where("slug = ? AND status = ?", slug, model.ProjectCategoryPublished).First(&category).Error
	if err != nil {
		return nil, err
	}
	if category.Slug != slug {
		return nil, gorm.ErrRecordNotFound
	}
	return &category, nil
}

func (r *ProjectCategoryRepo) categorySearch(filter ProjectCategoryFilter) *gorm.DB {
	query := r.db.Model(&model.ProjectCategory{})
	if filter.Name != "" {
		query = query.Where("name LIKE ?", "%"+filter.Name+"%")
	}
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	return query
}

func (r *ProjectCategoryRepo) FindAll(filter ProjectCategoryFilter) ([]ProjectCategoryAdminRow, int64, error) {
	var total int64
	query := r.categorySearch(filter)
	if err := query.Session(&gorm.Session{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	rows := make([]ProjectCategoryAdminRow, 0)
	err := query.Select(`project_categories.*,
		(SELECT COUNT(*) FROM project_category_memberships m
		 JOIN projects p ON p.id = m.project_id AND p.deleted_at IS NULL
		 WHERE m.category_id = project_categories.id) AS project_count`).
		Order("id asc").Offset((filter.Page - 1) * filter.PerPage).Limit(filter.PerPage).Scan(&rows).Error
	return rows, total, err
}

func (r *ProjectCategoryRepo) FindOptions(filter ProjectCategoryFilter) ([]model.ProjectCategory, int64, error) {
	var total int64
	query := r.categorySearch(filter)
	if err := query.Session(&gorm.Session{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	rows := make([]model.ProjectCategory, 0)
	err := query.Select("id", "name", "slug", "status").Order("id asc").
		Offset((filter.Page - 1) * filter.PerPage).Limit(filter.PerPage).Find(&rows).Error
	return rows, total, err
}

func (r *ProjectCategoryRepo) FindPublicSummaries() ([]ProjectCategoryPublicRow, error) {
	rows := make([]ProjectCategoryPublicRow, 0)
	err := r.db.Raw(`SELECT c.id, c.name, c.slug, COUNT(*) AS project_count,
		CAST(GREATEST(c.updated_at, COALESCE(MAX(p.updated_at), c.updated_at)) AS DATETIME(6)) AS updated_at
		FROM project_categories c
		JOIN project_category_memberships m ON m.category_id = c.id
		JOIN projects p ON p.id = m.project_id AND p.status = 1 AND p.deleted_at IS NULL
		WHERE c.status = ?
		GROUP BY c.id, c.name, c.slug, c.updated_at
		ORDER BY c.id`, model.ProjectCategoryPublished).Scan(&rows).Error
	return rows, err
}

func (r *ProjectCategoryRepo) FindPublishedProjects(categoryID uint64) ([]model.Project, error) {
	projects := make([]model.Project, 0)
	err := r.db.Model(&model.Project{}).
		Select("projects.id", "projects.name", "projects.slug", "projects.country", "projects.flag_emoji", "projects.tagline",
			"projects.overview_text", "projects.investment_amount", "projects.processing_period", "projects.target_crowd", "projects.cover_image", "projects.updated_at").
		Joins("JOIN project_category_memberships m ON m.project_id = projects.id").
		Where("m.category_id = ? AND projects.status = 1", categoryID).
		Order("m.sort_order ASC, m.project_id ASC").Find(&projects).Error
	return projects, err
}

func (r *ProjectCategoryRepo) FindProjects(categoryID uint64, page, perPage int) ([]ProjectCategoryProjectRow, int64, error) {
	query := r.db.Model(&model.Project{}).
		Joins("JOIN project_category_memberships m ON m.project_id = projects.id").Where("m.category_id = ?", categoryID)
	var total int64
	if err := query.Session(&gorm.Session{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	rows := make([]ProjectCategoryProjectRow, 0)
	err := query.Select("projects.id, projects.name, projects.slug, projects.status, m.sort_order").
		Order("m.sort_order ASC, m.project_id ASC").Offset((page - 1) * perPage).Limit(perPage).Scan(&rows).Error
	return rows, total, err
}

func (r *ProjectCategoryRepo) AssignmentsByProjectIDs(ids []uint64) ([]ProjectCategoryAssignment, error) {
	rows := make([]ProjectCategoryAssignment, 0)
	if len(ids) == 0 {
		return rows, nil
	}
	err := r.db.Table("project_category_memberships AS m").
		Select("m.project_id, c.id, c.name, c.slug, c.status").Joins("JOIN project_categories c ON c.id = m.category_id").
		Where("m.project_id IN ?", ids).Order("c.id").Scan(&rows).Error
	return rows, err
}

func (r *ProjectCategoryRepo) Create(category *model.ProjectCategory) error {
	return r.db.Create(category).Error
}

func (r *ProjectCategoryRepo) Update(id uint64, mutate func(*model.ProjectCategory) error) (*model.ProjectCategory, error) {
	var category model.ProjectCategory
	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&category, id).Error; err != nil {
			return err
		}
		if err := mutate(&category); err != nil {
			return err
		}
		return tx.Omit("created_at").Save(&category).Error
	})
	return &category, err
}

func (r *ProjectCategoryRepo) Delete(id uint64) error {
	result := r.db.Delete(&model.ProjectCategory{}, id)
	if result.Error == nil && result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return result.Error
}

func (r *ProjectCategoryRepo) UpdateProjectOrder(categoryID, projectID uint64, sortOrder int32) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var membership model.ProjectCategoryMembership
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("category_id = ? AND project_id = ?", categoryID, projectID).First(&membership).Error; err != nil {
			return err
		}
		if membership.SortOrder == sortOrder {
			return nil
		}
		if err := tx.Model(&membership).Update("sort_order", sortOrder).Error; err != nil {
			return err
		}
		return touchProjectCategories(tx, []uint64{categoryID})
	})
}

func touchProjectCategories(tx *gorm.DB, categoryIDs []uint64) error {
	if len(categoryIDs) == 0 {
		return nil
	}
	return tx.Model(&model.ProjectCategory{}).Where("id IN ?", categoryIDs).
		UpdateColumn("updated_at", gorm.Expr("CURRENT_TIMESTAMP(6)")).Error
}
