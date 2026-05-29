package repository

import (
	"mygo-immigration/backend/internal/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type CompareConfigRepo struct {
	db *gorm.DB
}

func (r *CompareConfigRepo) FindByProjectID(projectID uint64) (*model.CompareConfig, error) {
	var cfg model.CompareConfig
	err := r.db.Where("project_id = ?", projectID).First(&cfg).Error
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

func (r *CompareConfigRepo) Upsert(cfg *model.CompareConfig) error {
	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "project_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"compare_with", "compare_fields", "updated_at"}),
	}).Create(cfg).Error
}

func (r *CompareConfigRepo) DeleteByProjectID(projectID uint64) error {
	return r.db.Unscoped().Where("project_id = ?", projectID).Delete(&model.CompareConfig{}).Error
}
