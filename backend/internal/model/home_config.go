package model

import (
	"encoding/json"
	"time"
)

type HomeConfig struct {
	ID          uint64          `gorm:"primaryKey;autoIncrement" json:"id"`
	ConfigKey   string          `gorm:"uniqueIndex;size:64;not null" json:"config_key"`
	ConfigValue json.RawMessage `gorm:"type:json;not null" json:"config_value"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

func (HomeConfig) TableName() string { return "home_configs" }
