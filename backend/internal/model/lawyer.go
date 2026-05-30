package model

import (
	"encoding/json"
	"strings"
	"time"

	"gorm.io/gorm"
)

type Lawyer struct {
	ID        uint64         `gorm:"primaryKey;autoIncrement" json:"id"`
	PhotoURL  string         `gorm:"size:512;not null;default:''" json:"photo_url"`
	Name      string         `gorm:"size:64;not null;default:''" json:"name"`
	Title     string         `gorm:"size:128;not null;default:''" json:"title"`
	Tags      string         `gorm:"size:512;not null;default:''" json:"-"`
	SortOrder int            `gorm:"not null;default:0" json:"sort_order"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

func (Lawyer) TableName() string { return "lawyers" }

func (l Lawyer) MarshalJSON() ([]byte, error) {
	type Alias Lawyer
	tags := []string{}
	if l.Tags != "" {
		for _, t := range strings.Split(l.Tags, ",") {
			trimmed := strings.TrimSpace(t)
			if trimmed != "" {
				tags = append(tags, trimmed)
			}
		}
	}
	return json.Marshal(&struct {
		Tags []string `json:"tags"`
		*Alias
	}{
		Tags:  tags,
		Alias: (*Alias)(&l),
	})
}

func (l *Lawyer) UnmarshalJSON(data []byte) error {
	type Alias Lawyer
	aux := &struct {
		Tags []string `json:"tags"`
		*Alias
	}{
		Alias: (*Alias)(l),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	l.Tags = strings.Join(aux.Tags, ",")
	return nil
}
