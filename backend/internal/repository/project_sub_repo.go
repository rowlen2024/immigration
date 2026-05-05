package repository

import (
	"mygo-immigration/backend/internal/model"

	"gorm.io/gorm"
)

type RequirementRepo struct {
	db *gorm.DB
}

func (r *RequirementRepo) FindByProjectID(projectID uint64) ([]model.Requirement, error) {
	var items []model.Requirement
	err := r.db.Where("project_id = ?", projectID).Order("sort_order asc").Find(&items).Error
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (r *RequirementRepo) Create(item *model.Requirement) error {
	return r.db.Create(item).Error
}

func (r *RequirementRepo) Update(item *model.Requirement) error {
	return r.db.Omit("created_at").Save(item).Error
}

func (r *RequirementRepo) Delete(id uint64) error {
	return r.db.Delete(&model.Requirement{}, id).Error
}

type CostItemRepo struct {
	db *gorm.DB
}

func (r *CostItemRepo) FindByProjectID(projectID uint64) ([]model.CostItem, error) {
	var items []model.CostItem
	err := r.db.Where("project_id = ?", projectID).Order("sort_order asc").Find(&items).Error
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (r *CostItemRepo) Create(item *model.CostItem) error {
	return r.db.Create(item).Error
}

func (r *CostItemRepo) Update(item *model.CostItem) error {
	return r.db.Omit("created_at").Save(item).Error
}

func (r *CostItemRepo) Delete(id uint64) error {
	return r.db.Delete(&model.CostItem{}, id).Error
}

type TimelinePhaseRepo struct {
	db *gorm.DB
}

func (r *TimelinePhaseRepo) FindByProjectID(projectID uint64) ([]model.TimelinePhase, error) {
	var items []model.TimelinePhase
	err := r.db.Where("project_id = ?", projectID).Order("sort_order asc").Find(&items).Error
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (r *TimelinePhaseRepo) Create(item *model.TimelinePhase) error {
	return r.db.Create(item).Error
}

func (r *TimelinePhaseRepo) Update(item *model.TimelinePhase) error {
	return r.db.Omit("created_at").Save(item).Error
}

func (r *TimelinePhaseRepo) Delete(id uint64) error {
	return r.db.Delete(&model.TimelinePhase{}, id).Error
}
