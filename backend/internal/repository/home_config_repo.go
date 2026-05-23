package repository

import (
	"encoding/json"
	"mygo-immigration/backend/internal/model"

	"gorm.io/gorm"
)

type HomeConfigRepo struct {
	db *gorm.DB
}

func (r *HomeConfigRepo) FindByKey(key string) (*model.HomeConfig, error) {
	var cfg model.HomeConfig
	err := r.db.Where("config_key = ?", key).First(&cfg).Error
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

func (r *HomeConfigRepo) FindAll() ([]model.HomeConfig, error) {
	var configs []model.HomeConfig
	err := r.db.Find(&configs).Error
	if err != nil {
		return nil, err
	}
	return configs, nil
}

func (r *HomeConfigRepo) Create(cfg *model.HomeConfig) error {
	return r.db.Create(cfg).Error
}

func (r *HomeConfigRepo) Update(cfg *model.HomeConfig) error {
	return r.db.Omit("created_at").Save(cfg).Error
}

// FindAllConfigValues returns all config_value JSON from home_configs table.
func (r *HomeConfigRepo) FindAllConfigValues() ([]json.RawMessage, error) {
	var values []json.RawMessage
	err := r.db.Model(&model.HomeConfig{}).Pluck("config_value", &values).Error
	return values, err
}
