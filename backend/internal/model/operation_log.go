package model

import (
	"encoding/json"
	"time"
)

type OperationLog struct {
	ID         uint64          `gorm:"primaryKey;autoIncrement" json:"id"`
	OperatorID *uint64         `gorm:"index" json:"operator_id"`
	Action     string          `gorm:"size:64;not null;index" json:"action"`
	Target     string          `gorm:"size:128;not null;default:''" json:"target"`
	TargetID   *uint64         `json:"target_id"`
	IP         string          `gorm:"size:45;not null;default:''" json:"ip"`
	Details    json.RawMessage `gorm:"type:json" json:"details"`
	CreatedAt  time.Time       `gorm:"index" json:"created_at"`
}

func (OperationLog) TableName() string { return "operation_logs" }
