package repository

import (
	"mygo-immigration/backend/internal/model"
	"time"

	"gorm.io/gorm"
)

type CaseRepo struct {
	db *gorm.DB
}

func (r *CaseRepo) FindByProjectID(projectID uint64) ([]model.Case, error) {
	var cases []model.Case
	err := r.db.
		Where("project_id = ?", projectID).
		Order("sort_order asc").
		Find(&cases).Error
	if err != nil {
		return nil, err
	}
	return cases, nil
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

func (r *CaseRepo) FindAll(search string) ([]model.Case, error) {
	var cases []model.Case
	q := r.db.Preload("Project").Order("sort_order asc")
	if search != "" {
		q = q.Where("name LIKE ?", "%"+search+"%")
	}
	err := q.Find(&cases).Error
	if err != nil {
		return nil, err
	}
	return cases, nil
}

func (r *CaseRepo) Create(c *model.Case) error {
	return r.db.Create(c).Error
}

func (r *CaseRepo) Update(c *model.Case) error {
	return r.db.Omit("created_at").Save(c).Error
}

func (r *CaseRepo) Delete(id uint64) error {
	return r.db.Delete(&model.Case{}, id).Error
}

// HardDelete permanently removes a case record (bypasses soft delete).
func (r *CaseRepo) HardDelete(id uint64) error {
	return r.db.Unscoped().Delete(&model.Case{}, id).Error
}

func (r *CaseRepo) Count() (int64, error) {
	var c int64
	err := r.db.Model(&model.Case{}).Count(&c).Error
	return c, err
}

func (r *CaseRepo) CountByRange(start, end time.Time) (int64, error) {
	var c int64
	err := r.db.Model(&model.Case{}).Where("created_at >= ? AND created_at < ?", start, end).Count(&c).Error
	return c, err
}
