package model

import (
	"time"

	"gorm.io/gorm"
)

type Project struct {
	ID                 uint64         `gorm:"primaryKey;autoIncrement" json:"id"`
	Slug               string         `gorm:"uniqueIndex;size:64;not null" json:"slug"`
	Name               string         `gorm:"size:128;not null" json:"name"`
	Country            string         `gorm:"size:64;not null;default:''" json:"country"`
	FlagEmoji          string         `gorm:"size:16;not null;default:''" json:"flag_emoji"`
	Tagline            string         `gorm:"size:255;not null;default:''" json:"tagline"`
	InvestmentAmount   string         `gorm:"size:64;not null;default:''" json:"investment_amount"`
	InvestmentValue    *float64       `gorm:"type:decimal(12,2)" json:"investment_value"`
	InvestmentCurrency string         `gorm:"size:3;not null;default:'USD'" json:"investment_currency"`
	ProcessingPeriod   string         `gorm:"size:64;not null;default:''" json:"processing_period"`
	TargetCrowd        string         `gorm:"size:255;not null;default:''" json:"target_crowd"`
	OverviewTitle      string         `gorm:"size:255;not null;default:''" json:"overview_title"`
	OverviewText       string         `gorm:"type:text" json:"overview_text"`
	PolicyTitle        string         `gorm:"size:255;not null;default:''" json:"policy_title"`
	PolicyText         string         `gorm:"type:text" json:"policy_text"`
	CostsTotal         string         `gorm:"size:128;not null;default:''" json:"costs_total"`
	CostsNote          string         `gorm:"size:255;not null;default:''" json:"costs_note"`
	CtaText            string         `gorm:"size:128;not null;default:'立即咨询'" json:"cta_text"`
	HeroTitle          string         `gorm:"size:255;not null;default:''" json:"hero_title"`
	HeroDesc           string         `gorm:"size:512;not null;default:''" json:"hero_desc"`
	HeroGradient       string         `gorm:"size:255;not null;default:''" json:"hero_gradient"`
	CoverImage         string         `gorm:"size:512;not null;default:''" json:"cover_image"`
	SortOrder          int            `gorm:"not null;default:0" json:"sort_order"`
	Status             int8           `gorm:"not null;default:1" json:"status"`
	DeletedAt          gorm.DeletedAt `gorm:"index" json:"-"`
	CreatedAt          time.Time      `json:"created_at"`
	UpdatedAt          time.Time      `json:"updated_at"`

	Requirements   []Requirement   `gorm:"foreignKey:ProjectID" json:"requirements,omitempty"`
	CostItems      []CostItem      `gorm:"foreignKey:ProjectID" json:"cost_items,omitempty"`
	TimelinePhases []TimelinePhase `gorm:"foreignKey:ProjectID" json:"timeline_phases,omitempty"`
	Milestones     []Milestone     `gorm:"foreignKey:ProjectID" json:"milestones,omitempty"`
	Advantages     []ProjectAdvantage `gorm:"foreignKey:ProjectID" json:"advantages,omitempty"`
	FAQs           []FAQ           `gorm:"foreignKey:ProjectID" json:"faqs,omitempty"`
	Cases          []Case          `gorm:"foreignKey:ProjectID" json:"cases,omitempty"`
	Testimonials   []Testimonial   `gorm:"foreignKey:ProjectID" json:"testimonials,omitempty"`
	News           []Page         `gorm:"many2many:project_news;" json:"news,omitempty"`
	CompareConfig  *CompareConfig `gorm:"foreignKey:ProjectID" json:"compare_config,omitempty"`
}

func (Project) TableName() string { return "projects" }
