package dto

// UpdateCaseRequest 案例更新请求（slug 由后端自动生成，不可修改）
type UpdateCaseRequest struct {
	Name             string   `json:"name"`
	CountryFrom      string   `json:"country_from"`
	ProjectID        *uint64  `json:"project_id"`
	InvestmentAmount string   `json:"investment_amount"`
	InvestmentValue  *float64 `json:"investment_value"`
	ProcessingPeriod string   `json:"processing_period"`
	Content          string   `json:"content"`
	PhotoURL         string   `json:"photo_url"`
	SortOrder        int      `json:"sort_order"`
}

// UpdatePageRequest 页面更新请求
type UpdatePageRequest struct {
	ProjectID       *uint64 `json:"project_id"`
	Title           string  `json:"title"`
	Slug            string  `json:"slug"`
	Content         string  `json:"content"`
	CoverImage      string  `json:"cover_image"`
	MetaTitle       string  `json:"meta_title"`
	MetaDescription string  `json:"meta_description"`
	Template        string  `json:"template"`
	PageType        string  `json:"page_type"`
	Status          string  `json:"status"`
	SortOrder       int     `json:"sort_order"`
}

// UpdateProjectRequest 项目更新请求
type UpdateProjectRequest struct {
	Slug               string   `json:"slug"`
	Name               string   `json:"name"`
	Country            string   `json:"country"`
	FlagEmoji          string   `json:"flag_emoji"`
	Tagline            string   `json:"tagline"`
	InvestmentAmount   string   `json:"investment_amount"`
	InvestmentValue    *float64 `json:"investment_value"`
	InvestmentCurrency string   `json:"investment_currency"`
	ProcessingPeriod   string   `json:"processing_period"`
	TargetCrowd        string   `json:"target_crowd"`
	OverviewTitle      string   `json:"overview_title"`
	OverviewText       string   `json:"overview_text"`
	PolicyTitle        string   `json:"policy_title"`
	PolicyText         string   `json:"policy_text"`
	CostsTotal         string   `json:"costs_total"`
	CostsNote          string   `json:"costs_note"`
	CtaText            string   `json:"cta_text"`
	HeroTitle          string   `json:"hero_title"`
	HeroDesc           string   `json:"hero_desc"`
	HeroGradient       string   `json:"hero_gradient"`
	CoverImage         string   `json:"cover_image"`
	SortOrder          int      `json:"sort_order"`
	Status             int8     `json:"status"`
}

// UpdateFAQRequest FAQ 更新请求
type UpdateFAQRequest struct {
	ProjectID *uint64 `json:"project_id"`
	Question  string  `json:"question"`
	Answer    string  `json:"answer"`
	IsGlobal  bool    `json:"is_global"`
	SortOrder int     `json:"sort_order"`
}

// UpdateTestimonialRequest 客户评价更新请求
type UpdateTestimonialRequest struct {
	AvatarURL string `json:"avatar_url"`
	Nickname  string `json:"nickname"`
	Rating    uint8  `json:"rating"`
	Content   string `json:"content"`
	SortOrder int    `json:"sort_order"`
}

// UpdateUserRequest 用户更新请求
type UpdateUserRequest struct {
	DisplayName string `json:"display_name"`
	Role        string `json:"role"`
	Status      *int8  `json:"status"`
	Password    string `json:"password,omitempty"`
}
