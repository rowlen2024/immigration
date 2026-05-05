package model

import (
	"time"

	"gorm.io/gorm"
)

type Lead struct {
	ID                uint64         `gorm:"primaryKey;autoIncrement" json:"id"`
	Name              string         `gorm:"size:128;not null;default:''" json:"name"`
	Phone             string         `gorm:"size:32;not null;default:''" json:"phone"`
	Email             string         `gorm:"size:128;not null;default:''" json:"email"`
	InterestedProject string         `gorm:"size:64;not null;default:'';index" json:"interested_project"`
	Message           string         `gorm:"type:text" json:"message"`
	Status            string         `gorm:"type:enum('new','contacted','qualified','closed');not null;default:'new';index" json:"status"`
	Notes             string         `gorm:"type:text" json:"notes"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"-"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
}

func (Lead) TableName() string { return "leads" }
