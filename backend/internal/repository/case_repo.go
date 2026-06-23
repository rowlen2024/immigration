package repository

import (
	"mygo-immigration/backend/internal/model"
	"time"

	"gorm.io/gorm"
)

type CaseRepo struct {
	db *gorm.DB
}

func (r *CaseRepo) FindByID(id uint64) (*model.Case, error) {
	var c model.Case
	err := r.db.Preload("Project").First(&c, id).Error
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *CaseRepo) FindBySlug(slug string) (*model.Case, error) {
	var c model.Case
	err := r.db.Preload("Project").Where("slug = ?", slug).First(&c).Error
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *CaseRepo) FindByIDs(ids []uint64) ([]model.Case, error) {
	var cases []model.Case
	err := r.db.Preload("Project").Where("id IN ?", ids).Find(&cases).Error
	if err != nil {
		return nil, err
	}
	return cases, nil
}

func (r *CaseRepo) FindAll(filter CaseFilter) ([]model.Case, int64, error) {
	var cases []model.Case
	var total int64

	q := r.db.Model(&model.Case{}).Preload("Project")
	if filter.ProjectID != nil {
		q = q.Where("project_id = ?", *filter.ProjectID)
	}
	if filter.CountryFrom != "" {
		q = q.Where("country_from = ?", filter.CountryFrom)
	}
	if filter.Name != "" {
		q = q.Where("name LIKE ?", "%"+filter.Name+"%")
	}

	if err := q.Session(&gorm.Session{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	q = q.Order("sort_order asc, id asc")
	if filter.Page > 0 && filter.PerPage > 0 {
		offset := (filter.Page - 1) * filter.PerPage
		q = q.Offset(offset).Limit(filter.PerPage)
	}

	if err := q.Find(&cases).Error; err != nil {
		return nil, 0, err
	}
	return cases, total, nil
}

func (r *CaseRepo) FindOptions(filter CaseFilter) ([]CaseOptionRow, int64, error) {
	var cases []CaseOptionRow
	var total int64

	q := r.db.Model(&model.Case{})
	if filter.ProjectID != nil {
		q = q.Where("project_id = ?", *filter.ProjectID)
	}
	if filter.CountryFrom != "" {
		q = q.Where("country_from = ?", filter.CountryFrom)
	}
	if filter.Name != "" {
		q = q.Where("name LIKE ?", "%"+filter.Name+"%")
	}

	if err := q.Session(&gorm.Session{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	q = q.Select("id", "name").Order("sort_order asc, id asc")
	if filter.Page > 0 && filter.PerPage > 0 {
		offset := (filter.Page - 1) * filter.PerPage
		q = q.Offset(offset).Limit(filter.PerPage)
	}

	if err := q.Find(&cases).Error; err != nil {
		return nil, 0, err
	}
	return cases, total, nil
}

func (r *CaseRepo) Create(c *model.Case) error {
	return r.db.Create(c).Error
}

func (r *CaseRepo) Update(c *model.Case) error {
	return r.db.Omit("created_at").Save(c).Error
}

func (r *CaseRepo) Delete(id uint64) error {
	return r.db.Unscoped().Delete(&model.Case{}, id).Error
}

// DeleteByProjectID soft-deletes all cases belonging to a project.
func (r *CaseRepo) DeleteByProjectID(projectID uint64) error {
	return r.db.Unscoped().Where("project_id = ?", projectID).Delete(&model.Case{}).Error
}

func (r *CaseRepo) Count() (int64, error) {
	return CountByModel[model.Case](r.db)
}

// FindAllPhotoURLs returns non-empty photo_url values referencing /uploads/ (unscoped).
func (r *CaseRepo) FindAllPhotoURLs() ([]string, error) {
	return PluckUploadsByColumn[model.Case](r.db, "photo_url")
}

// FindAllContents returns content values that contain /uploads/ references (unscoped).
func (r *CaseRepo) FindAllContents() ([]string, error) {
	return PluckUploadsByColumn[model.Case](r.db, "content")
}

func (r *CaseRepo) CountByRange(start, end time.Time) (int64, error) {
	var c int64
	err := r.db.Model(&model.Case{}).Where("created_at >= ? AND created_at < ?", start, end).Count(&c).Error
	return c, err
}
