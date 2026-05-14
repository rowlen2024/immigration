package service

import (
	"encoding/json"
	"fmt"

	"mygo-immigration/backend/internal/model"
	"mygo-immigration/backend/internal/repository"
)

// HomeConfigService handles business logic for the homepage configuration.
type HomeConfigService struct {
	repo *repository.HomeConfigRepo
}

// HeroSlide represents a slide in the hero section of the homepage.
type HeroSlide struct {
	Title       string `json:"title"`
	Desc        string `json:"desc"`
	ProjectSlug string `json:"project_slug"`
	Gradient    string `json:"gradient"`
	Image       string `json:"image"`
	Link        string `json:"link"`
}

// AdvantageItem represents an advantage card on the homepage.
type AdvantageItem struct {
	Icon        string `json:"icon"`
	IconType    string `json:"icon_type"` // "lucide" for svg icons, empty for legacy emoji
	Title       string `json:"title"`
	Description string `json:"description"`
}

// AdvantageSectionConfig holds the advantage section title and subtitle.
type AdvantageSectionConfig struct {
	SectionTitle    string `json:"section_title"`
	SectionSubtitle string `json:"section_subtitle"`
	Image           string `json:"image"`
}

// ProjectShowcaseConfig holds the project showcase section settings.
type ProjectShowcaseConfig struct {
	SectionTitle    string   `json:"section_title"`
	SectionSubtitle string   `json:"section_subtitle"`
	FeaturedSlugs   []string `json:"featured_slugs"`
}

// CaseShowcaseConfig holds the case showcase section settings.
type CaseShowcaseConfig struct {
	SectionTitle    string   `json:"section_title"`
	SectionSubtitle string   `json:"section_subtitle"`
	FeaturedCaseIDs []uint64 `json:"featured_case_ids"`
}

// HomeConfigData holds the parsed homepage configuration data.
type HomeConfigData struct {
	HeroSlides      []HeroSlide              `json:"hero_slides"`
	AdvantageItems  []AdvantageItem          `json:"advantage_items"`
	AdvantageSection  *AdvantageSectionConfig  `json:"advantage_section"`
	ProjectShowcase *ProjectShowcaseConfig   `json:"project_showcase"`
	CaseShowcase    *CaseShowcaseConfig      `json:"case_showcase"`
}

// Get returns the homepage configuration with parsed hero_slides and advantage_items.
func (s *HomeConfigService) Get() (*HomeConfigData, error) {
	data := &HomeConfigData{}

	if heroCfg, err := s.repo.FindByKey("hero_slides"); err == nil {
		json.Unmarshal(heroCfg.ConfigValue, &data.HeroSlides)
	}
	if advCfg, err := s.repo.FindByKey("advantage_items"); err == nil {
		json.Unmarshal(advCfg.ConfigValue, &data.AdvantageItems)
	}
	if advSecCfg, err := s.repo.FindByKey("advantage_section"); err == nil {
		var asc AdvantageSectionConfig
		if err := json.Unmarshal(advSecCfg.ConfigValue, &asc); err == nil {
			data.AdvantageSection = &asc
		}
	}
	if projCfg, err := s.repo.FindByKey("project_showcase"); err == nil {
		var psc ProjectShowcaseConfig
		if err := json.Unmarshal(projCfg.ConfigValue, &psc); err == nil {
			data.ProjectShowcase = &psc
		}
	}
	if caseCfg, err := s.repo.FindByKey("case_showcase"); err == nil {
		var csc CaseShowcaseConfig
		if err := json.Unmarshal(caseCfg.ConfigValue, &csc); err == nil {
			data.CaseShowcase = &csc
		}
	}

	return data, nil
}

// SiteConfig holds all site-wide settings.
type SiteConfig struct {
	SiteName                string   `json:"site_name"`
	SiteLogo                string   `json:"site_logo"`
	SiteFavicon             string   `json:"site_favicon"`
	SeoTitle                string   `json:"seo_title"`
	SeoDescription          string   `json:"seo_description"`
	SeoKeywords             string   `json:"seo_keywords"`
	OgImage                 string   `json:"og_image"`
	CanonicalBase           string   `json:"canonical_base"`
	OrganizationName        string   `json:"organization_name"`
	OrganizationDescription string   `json:"organization_description"`
	OrganizationLogo        string   `json:"organization_logo"`
	OrganizationURL         string   `json:"organization_url"`
	SameAs                  []string `json:"same_as"`
	ContactPhone            string   `json:"contact_phone"`
	ContactPhone2           string   `json:"contact_phone_2"`
	ContactEmail            string   `json:"contact_email"`
	ContactAddress          string   `json:"contact_address"`
	ContactWechat           string   `json:"contact_wechat"`
	ContactWechatMP         string   `json:"contact_wechat_mp"`
	ContactWechatVideo      string   `json:"contact_wechat_video"`
	GATrackingID            string   `json:"ga_tracking_id"`
	BaiduTongjiID           string   `json:"baidu_tongji_id"`
	CustomHeadCode          string   `json:"custom_head_code"`
	CustomBodyCode          string   `json:"custom_body_code"`
	CopyrightText           string   `json:"copyright_text"`
	ICPNumber               string   `json:"icp_number"`
	FooterTagline           string   `json:"footer_tagline"`
}

// DefaultSiteConfig returns sensible zero-value defaults.
func DefaultSiteConfig() *SiteConfig {
	return &SiteConfig{
		SiteName:      "北极星移民",
		SeoTitle:      "{site_name} | 专业投资移民服务",
		CopyrightText: "© {year} {site_name}. All rights reserved.",
	}
}

// GetSiteConfig returns the site configuration, falling back to defaults.
func (s *HomeConfigService) GetSiteConfig() (*SiteConfig, error) {
	cfg, err := s.repo.FindByKey("site")
	if err != nil {
		return DefaultSiteConfig(), nil
	}

	var data SiteConfig
	if err := json.Unmarshal(cfg.ConfigValue, &data); err != nil {
		return nil, fmt.Errorf("failed to parse site config: %w", err)
	}
	return &data, nil
}

// UpdateSiteConfig replaces the entire site configuration.
func (s *HomeConfigService) UpdateSiteConfig(input *SiteConfig) error {
	raw, err := json.Marshal(input)
	if err != nil {
		return fmt.Errorf("failed to marshal site config: %w", err)
	}

	existing, err := s.repo.FindByKey("site")
	if err != nil {
		cfg := &model.HomeConfig{
			ConfigKey:   "site",
			ConfigValue: raw,
		}
		return s.repo.Create(cfg)
	}

	existing.ConfigValue = raw
	return s.repo.Update(existing)
}

// Update saves one or more homepage configuration entries.
func (s *HomeConfigService) Update(configs map[string]json.RawMessage) error {
	for key, rawValue := range configs {
		existing, err := s.repo.FindByKey(key)
		if err != nil {
			cfg := &model.HomeConfig{
				ConfigKey:   key,
				ConfigValue: rawValue,
			}
			if err := s.repo.Create(cfg); err != nil {
				return fmt.Errorf("failed to save home config key %s: %w", key, err)
			}
		} else {
			existing.ConfigValue = rawValue
			if err := s.repo.Update(existing); err != nil {
				return fmt.Errorf("failed to save home config key %s: %w", key, err)
			}
		}
	}
	return nil
}
