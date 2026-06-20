package service

import (
	"encoding/json"
	"fmt"
	"mygo-immigration/backend/internal/logging"

	"mygo-immigration/backend/internal/model"
	"mygo-immigration/backend/internal/repository"
)

// HomeConfigService handles business logic for the homepage configuration.
type HomeConfigService struct {
	repo            repository.HomeConfigRepository
	projectRepo     repository.ProjectRepository
	caseRepo        repository.CaseRepository
	testimonialRepo repository.TestimonialRepository
	versionRepo     *repository.PublicVersionRepo
}

func (s *HomeConfigService) RegisterPublicVersions(reg *PublicVersionRegistry) {
	reg.Register("public:site-config", func(string) (repository.PublicVersion, error) {
		return tableVersion(s.versionRepo, "home_configs", "config_key = ?", "site")
	})
	reg.Register("public:home", s.publicHomeVersion)
}

func (s *HomeConfigService) publicHomeVersion(string) (repository.PublicVersion, error) {
	versions := make([]repository.PublicVersion, 0, 4)
	homeCfgVersion, err := tableVersion(
		s.versionRepo,
		"home_configs",
		"config_key IN ?",
		[]string{"hero_slides", "advantage_items", "advantage_section", "project_showcase", "case_showcase", "testimonial_showcase", "hero_trust"},
	)
	if err != nil {
		return repository.PublicVersion{}, err
	}
	versions = append(versions, homeCfgVersion)

	if slugs := s.featuredProjectSlugs(); len(slugs) > 0 {
		v, err := s.versionRepo.VersionFromQuery(
			"SELECT MAX(updated_at) AS updated_at, COUNT(*) AS count FROM projects WHERE deleted_at IS NULL AND slug IN ?",
			slugs,
		)
		if err != nil {
			return repository.PublicVersion{}, err
		}
		versions = append(versions, v)
	}
	if ids := s.featuredCaseIDs(); len(ids) > 0 {
		v, err := s.versionRepo.VersionFromQuery(
			"SELECT MAX(updated_at) AS updated_at, COUNT(*) AS count FROM cases WHERE deleted_at IS NULL AND id IN ?",
			ids,
		)
		if err != nil {
			return repository.PublicVersion{}, err
		}
		versions = append(versions, v)
	}
	if ids := s.featuredTestimonialIDs(); len(ids) > 0 {
		v, err := s.versionRepo.VersionFromQuery(
			"SELECT MAX(updated_at) AS updated_at, COUNT(*) AS count FROM testimonials WHERE deleted_at IS NULL AND id IN ?",
			ids,
		)
		if err != nil {
			return repository.PublicVersion{}, err
		}
		versions = append(versions, v)
	}
	return repository.MergePublicVersions(versions...), nil
}

func (s *HomeConfigService) featuredProjectSlugs() []string {
	cfg, err := s.repo.FindByKey("project_showcase")
	if err != nil {
		return nil
	}
	var data ProjectShowcaseConfig
	if json.Unmarshal(cfg.ConfigValue, &data) != nil {
		return nil
	}
	return data.FeaturedSlugs
}

func (s *HomeConfigService) featuredCaseIDs() []uint64 {
	cfg, err := s.repo.FindByKey("case_showcase")
	if err != nil {
		return nil
	}
	var data CaseShowcaseConfig
	if json.Unmarshal(cfg.ConfigValue, &data) != nil {
		return nil
	}
	return data.FeaturedCaseIDs
}

func (s *HomeConfigService) featuredTestimonialIDs() []uint64 {
	cfg, err := s.repo.FindByKey("testimonial_showcase")
	if err != nil {
		return nil
	}
	var data TestimonialShowcaseConfig
	if json.Unmarshal(cfg.ConfigValue, &data) != nil {
		return nil
	}
	return data.FeaturedTestimonialIDs
}

// HeroSlide represents a slide in the hero section of the homepage.
type HeroSlide struct {
	Title       string `json:"title"`
	Desc        string `json:"desc"`
	ProjectSlug string `json:"project_slug"`
	Gradient    string `json:"gradient"`
	Image       string `json:"image"`
	Link        string `json:"link"`

	ImageVariants map[string]model.ImageVariantInfo `json:"image_variants,omitempty"`
}

// TrustItem represents a single trust stat in the hero section.
type TrustItem struct {
	Number string `json:"number"`
	Label  string `json:"label"`
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

	ImageVariants map[string]model.ImageVariantInfo `json:"image_variants,omitempty"`
}

// FeaturedProject holds lightweight project data embedded in home-config.
type FeaturedProject struct {
	Name         string `json:"name"`
	Slug         string `json:"slug"`
	Tagline      string `json:"tagline"`
	CoverImage   string `json:"cover_image"`
	OverviewText string `json:"overview_text"`

	CoverImageVariants map[string]model.ImageVariantInfo `json:"cover_image_variants,omitempty"`
}

// FeaturedCase holds lightweight case data embedded in home-config.
type FeaturedCase struct {
	ID          uint64 `json:"id"`
	Slug        string `json:"slug"`
	Name        string `json:"name"`
	CountryFrom string `json:"country_from"`
	PhotoURL    string `json:"photo_url"`
	Content     string `json:"content"`
	ProjectName string `json:"project_name"`

	PhotoVariants map[string]model.ImageVariantInfo `json:"photo_variants,omitempty"`
}

// FeaturedTestimonial holds lightweight testimonial data embedded in home-config.
type FeaturedTestimonial struct {
	ID        uint64 `json:"id"`
	Nickname  string `json:"nickname"`
	AvatarURL string `json:"avatar_url"`
	Rating    uint8  `json:"rating"`
	Content   string `json:"content"`

	AvatarVariants map[string]model.ImageVariantInfo `json:"avatar_variants,omitempty"`
}

// ProjectShowcaseConfig holds the project showcase section settings.
type ProjectShowcaseConfig struct {
	SectionTitle     string            `json:"section_title"`
	SectionSubtitle  string            `json:"section_subtitle"`
	FeaturedSlugs    []string          `json:"featured_slugs"`
	FeaturedProjects []FeaturedProject `json:"featured_projects"`
}

// CaseShowcaseConfig holds the case showcase section settings.
type CaseShowcaseConfig struct {
	SectionTitle    string         `json:"section_title"`
	SectionSubtitle string         `json:"section_subtitle"`
	FeaturedCaseIDs []uint64       `json:"featured_case_ids"`
	FeaturedCases   []FeaturedCase `json:"featured_cases"`
}

// TestimonialShowcaseConfig holds the testimonial showcase section settings.
type TestimonialShowcaseConfig struct {
	SectionTitle           string                `json:"section_title"`
	SectionSubtitle        string                `json:"section_subtitle"`
	FeaturedTestimonialIDs []uint64              `json:"featured_testimonial_ids"`
	FeaturedTestimonials   []FeaturedTestimonial `json:"featured_testimonials"`
}

// HomeConfigData holds the parsed homepage configuration data.
type HomeConfigData struct {
	HeroSlides          []HeroSlide                `json:"hero_slides"`
	AdvantageItems      []AdvantageItem            `json:"advantage_items"`
	AdvantageSection    *AdvantageSectionConfig    `json:"advantage_section"`
	ProjectShowcase     *ProjectShowcaseConfig     `json:"project_showcase"`
	CaseShowcase        *CaseShowcaseConfig        `json:"case_showcase"`
	TestimonialShowcase *TestimonialShowcaseConfig `json:"testimonial_showcase"`
	TrustItems          []TrustItem                `json:"hero_trust"`
}

// Get returns the homepage configuration with parsed sections and embedded featured items.
func (s *HomeConfigService) Get() (*HomeConfigData, error) {
	data := &HomeConfigData{}

	if heroCfg, err := s.repo.FindByKey("hero_slides"); err == nil {
		if err := json.Unmarshal(heroCfg.ConfigValue, &data.HeroSlides); err != nil {
			logging.Logger.Warn("home_config: failed to unmarshal hero_slides", "error", err)
		}
	}
	for i := range data.HeroSlides {
		data.HeroSlides[i].ImageVariants = ResolveImageVariants(data.HeroSlides[i].Image, UploadContextHomepageSlide)
	}
	if advCfg, err := s.repo.FindByKey("advantage_items"); err == nil {
		if err := json.Unmarshal(advCfg.ConfigValue, &data.AdvantageItems); err != nil {
			logging.Logger.Warn("home_config: failed to unmarshal advantage_items", "error", err)
		}
	}
	if advSecCfg, err := s.repo.FindByKey("advantage_section"); err == nil {
		var asc AdvantageSectionConfig
		if err := json.Unmarshal(advSecCfg.ConfigValue, &asc); err == nil {
			asc.ImageVariants = ResolveImageVariants(asc.Image, UploadContextGeneral)
			data.AdvantageSection = &asc
		}
	}
	if projCfg, err := s.repo.FindByKey("project_showcase"); err == nil {
		var psc ProjectShowcaseConfig
		if err := json.Unmarshal(projCfg.ConfigValue, &psc); err == nil {
			s.loadFeaturedProjects(&psc)
			data.ProjectShowcase = &psc
		} else {
			logging.Logger.Warn("home_config: failed to unmarshal project_showcase", "error", err)
		}
	} else {
		logging.Logger.Warn("home_config: project_showcase config not found", "error", err)
	}
	if caseCfg, err := s.repo.FindByKey("case_showcase"); err == nil {
		var csc CaseShowcaseConfig
		if err := json.Unmarshal(caseCfg.ConfigValue, &csc); err == nil {
			s.loadFeaturedCases(&csc)
			data.CaseShowcase = &csc
		} else {
			logging.Logger.Warn("home_config: failed to unmarshal case_showcase", "error", err)
		}
	} else {
		logging.Logger.Warn("home_config: case_showcase config not found", "error", err)
	}
	if testimonialCfg, err := s.repo.FindByKey("testimonial_showcase"); err == nil {
		var tsc TestimonialShowcaseConfig
		if err := json.Unmarshal(testimonialCfg.ConfigValue, &tsc); err == nil {
			s.loadFeaturedTestimonials(&tsc)
			data.TestimonialShowcase = &tsc
		} else {
			logging.Logger.Warn("home_config: failed to unmarshal testimonial_showcase", "error", err)
		}
	} else {
		logging.Logger.Warn("home_config: testimonial_showcase config not found", "error", err)
	}
	if trustCfg, err := s.repo.FindByKey("hero_trust"); err == nil {
		if err := json.Unmarshal(trustCfg.ConfigValue, &data.TrustItems); err != nil {
			logging.Logger.Warn("home_config: failed to unmarshal hero_trust", "error", err)
		}
	}

	return data, nil
}

func (s *HomeConfigService) loadFeaturedProjects(psc *ProjectShowcaseConfig) {
	if len(psc.FeaturedSlugs) == 0 || s.projectRepo == nil {
		return
	}
	projects, err := s.projectRepo.FindBySlugsLight(psc.FeaturedSlugs)
	if err != nil {
		logging.Logger.Warn("home_config: failed to load featured projects", "error", err)
		return
	}
	// Self-healing: remove stale slugs
	if len(projects) < len(psc.FeaturedSlugs) {
		foundSlugs := make(map[string]bool, len(projects))
		for _, p := range projects {
			foundSlugs[p.Slug] = true
		}
		cleanSlugs := make([]string, 0, len(projects))
		for _, slug := range psc.FeaturedSlugs {
			if foundSlugs[slug] {
				cleanSlugs = append(cleanSlugs, slug)
			}
		}
		logging.Logger.Warn("home_config: detected stale featured project refs, auto-healing",
			"before", len(psc.FeaturedSlugs), "after", len(cleanSlugs))
		psc.FeaturedSlugs = cleanSlugs
		_ = s.persistFromJSONArrayStr("project_showcase", "featured_slugs", cleanSlugs)
	}
	items := make([]FeaturedProject, 0, len(projects))
	for _, p := range projects {
		items = append(items, FeaturedProject{
			Name:               p.Name,
			Slug:               p.Slug,
			Tagline:            p.Tagline,
			CoverImage:         p.CoverImage,
			CoverImageVariants: ResolveImageVariants(p.CoverImage, UploadContextProject),
			OverviewText:       p.OverviewText,
		})
	}
	psc.FeaturedProjects = items
}

func (s *HomeConfigService) loadFeaturedCases(csc *CaseShowcaseConfig) {
	if len(csc.FeaturedCaseIDs) == 0 || s.caseRepo == nil {
		return
	}
	cases, err := s.caseRepo.FindByIDs(csc.FeaturedCaseIDs)
	if err != nil {
		logging.Logger.Warn("home_config: failed to load featured cases", "error", err)
		return
	}
	// Self-healing: remove stale IDs (soft-deleted or hard-deleted records not returned by GORM)
	if len(cases) < len(csc.FeaturedCaseIDs) {
		foundIDs := make(map[uint64]bool, len(cases))
		for _, c := range cases {
			foundIDs[c.ID] = true
		}
		cleanIDs := make([]uint64, 0, len(cases))
		for _, id := range csc.FeaturedCaseIDs {
			if foundIDs[id] {
				cleanIDs = append(cleanIDs, id)
			}
		}
		logging.Logger.Warn("home_config: detected stale featured case refs, auto-healing",
			"before", len(csc.FeaturedCaseIDs), "after", len(cleanIDs))
		csc.FeaturedCaseIDs = cleanIDs
		_ = s.persistFromJSONArray("case_showcase", "featured_case_ids", cleanIDs)
	}
	// Preserve configured order
	orderMap := make(map[uint64]int)
	for i, id := range csc.FeaturedCaseIDs {
		orderMap[id] = i
	}
	items := make([]FeaturedCase, len(cases))
	for i, c := range cases {
		projectName := ""
		if c.Project != nil {
			projectName = c.Project.Name
		}
		items[i] = FeaturedCase{
			ID:            c.ID,
			Slug:          c.Slug,
			Name:          c.Name,
			CountryFrom:   c.CountryFrom,
			PhotoURL:      c.PhotoURL,
			PhotoVariants: ResolveImageVariants(c.PhotoURL, UploadContextCase),
			Content:       c.Content,
			ProjectName:   projectName,
		}
	}
	// Sort by configured order
	sortByOrder(items, func(item FeaturedCase) int {
		if idx, ok := orderMap[item.ID]; ok {
			return idx
		}
		return len(orderMap)
	})
	csc.FeaturedCases = items
}

func (s *HomeConfigService) loadFeaturedTestimonials(tsc *TestimonialShowcaseConfig) {
	if len(tsc.FeaturedTestimonialIDs) == 0 || s.testimonialRepo == nil {
		return
	}
	testimonials, err := s.testimonialRepo.FindByIDs(tsc.FeaturedTestimonialIDs)
	if err != nil {
		logging.Logger.Warn("home_config: failed to load featured testimonials", "error", err)
		return
	}
	// Self-healing: remove stale IDs
	if len(testimonials) < len(tsc.FeaturedTestimonialIDs) {
		foundIDs := make(map[uint64]bool, len(testimonials))
		for _, t := range testimonials {
			foundIDs[t.ID] = true
		}
		cleanIDs := make([]uint64, 0, len(testimonials))
		for _, id := range tsc.FeaturedTestimonialIDs {
			if foundIDs[id] {
				cleanIDs = append(cleanIDs, id)
			}
		}
		logging.Logger.Warn("home_config: detected stale featured testimonial refs, auto-healing",
			"before", len(tsc.FeaturedTestimonialIDs), "after", len(cleanIDs))
		tsc.FeaturedTestimonialIDs = cleanIDs
		_ = s.persistFromJSONArray("testimonial_showcase", "featured_testimonial_ids", cleanIDs)
	}
	orderMap := make(map[uint64]int)
	for i, id := range tsc.FeaturedTestimonialIDs {
		orderMap[id] = i
	}
	items := make([]FeaturedTestimonial, len(testimonials))
	for i, t := range testimonials {
		items[i] = FeaturedTestimonial{
			ID:             t.ID,
			Nickname:       t.Nickname,
			AvatarURL:      t.AvatarURL,
			AvatarVariants: ResolveImageVariants(t.AvatarURL, UploadContextTestimonial),
			Rating:         t.Rating,
			Content:        t.Content,
		}
	}
	sortByOrder(items, func(item FeaturedTestimonial) int {
		if idx, ok := orderMap[item.ID]; ok {
			return idx
		}
		return len(orderMap)
	})
	tsc.FeaturedTestimonials = items
}

// removeFromJSONArray removes a uint64 value from a JSON array field within a home_config row.
func (s *HomeConfigService) removeFromJSONArray(configKey, fieldName string, targetID uint64) error {
	cfg, err := s.repo.FindByKey(configKey)
	if err != nil {
		return nil // row doesn't exist, nothing to clean
	}
	var m map[string]json.RawMessage
	if err := json.Unmarshal(cfg.ConfigValue, &m); err != nil {
		return fmt.Errorf("home_config: failed to unmarshal %s: %w", configKey, err)
	}
	rawIDs, ok := m[fieldName]
	if !ok {
		return nil
	}
	var ids []uint64
	if err := json.Unmarshal(rawIDs, &ids); err != nil {
		return fmt.Errorf("home_config: failed to unmarshal %s.%s: %w", configKey, fieldName, err)
	}
	newIDs := make([]uint64, 0, len(ids))
	for _, id := range ids {
		if id != targetID {
			newIDs = append(newIDs, id)
		}
	}
	if len(newIDs) == len(ids) {
		return nil // target not found
	}
	newRaw, err := json.Marshal(newIDs)
	if err != nil {
		return fmt.Errorf("home_config: failed to marshal %s.%s: %w", configKey, fieldName, err)
	}
	m[fieldName] = newRaw
	updated, err := json.Marshal(m)
	if err != nil {
		return fmt.Errorf("home_config: failed to marshal %s: %w", configKey, err)
	}
	cfg.ConfigValue = updated
	return s.repo.Update(cfg)
}

// removeFromJSONArrayStr removes a string value from a JSON string array field.
func (s *HomeConfigService) removeFromJSONArrayStr(configKey, fieldName string, target string) error {
	cfg, err := s.repo.FindByKey(configKey)
	if err != nil {
		return nil
	}
	var m map[string]json.RawMessage
	if err := json.Unmarshal(cfg.ConfigValue, &m); err != nil {
		return fmt.Errorf("home_config: failed to unmarshal %s: %w", configKey, err)
	}
	rawSlugs, ok := m[fieldName]
	if !ok {
		return nil
	}
	var slugs []string
	if err := json.Unmarshal(rawSlugs, &slugs); err != nil {
		return fmt.Errorf("home_config: failed to unmarshal %s.%s: %w", configKey, fieldName, err)
	}
	newSlugs := make([]string, 0, len(slugs))
	for _, s := range slugs {
		if s != target {
			newSlugs = append(newSlugs, s)
		}
	}
	if len(newSlugs) == len(slugs) {
		return nil
	}
	newRaw, err := json.Marshal(newSlugs)
	if err != nil {
		return fmt.Errorf("home_config: failed to marshal %s.%s: %w", configKey, fieldName, err)
	}
	m[fieldName] = newRaw
	updated, err := json.Marshal(m)
	if err != nil {
		return fmt.Errorf("home_config: failed to marshal %s: %w", configKey, err)
	}
	cfg.ConfigValue = updated
	return s.repo.Update(cfg)
}

// RemoveFeaturedCaseID removes a case ID from case_showcase.featured_case_ids.
func (s *HomeConfigService) RemoveFeaturedCaseID(caseID uint64) error {
	return s.removeFromJSONArray("case_showcase", "featured_case_ids", caseID)
}

// RemoveFeaturedProjectSlug removes a project slug from project_showcase.featured_slugs.
func (s *HomeConfigService) RemoveFeaturedProjectSlug(slug string) error {
	return s.removeFromJSONArrayStr("project_showcase", "featured_slugs", slug)
}

// RemoveFeaturedTestimonialID removes a testimonial ID from testimonial_showcase.featured_testimonial_ids.
func (s *HomeConfigService) RemoveFeaturedTestimonialID(testimonialID uint64) error {
	return s.removeFromJSONArray("testimonial_showcase", "featured_testimonial_ids", testimonialID)
}

// CleanupFeaturedRefs removes stale case and testimonial IDs from home_config featured lists.
// Reads once, conditionally writes once per section — fixed 2-4 SQL regardless of stale count.
func (s *HomeConfigService) CleanupFeaturedRefs(caseIDs, testimonialIDs []uint64) error {
	if len(caseIDs) > 0 {
		staleCaseSet := make(map[uint64]bool, len(caseIDs))
		for _, id := range caseIDs {
			staleCaseSet[id] = true
		}
		cfg, err := s.repo.FindByKey("case_showcase")
		if err == nil {
			var m map[string]json.RawMessage
			if json.Unmarshal(cfg.ConfigValue, &m) == nil {
				if raw, ok := m["featured_case_ids"]; ok {
					var ids []uint64
					if json.Unmarshal(raw, &ids) == nil {
						cleaned := make([]uint64, 0, len(ids))
						for _, id := range ids {
							if !staleCaseSet[id] {
								cleaned = append(cleaned, id)
							}
						}
						if len(cleaned) < len(ids) {
							newRaw, _ := json.Marshal(cleaned)
							m["featured_case_ids"] = newRaw
							updated, _ := json.Marshal(m)
							cfg.ConfigValue = updated
							_ = s.repo.Update(cfg)
						}
					}
				}
			}
		}
	}
	if len(testimonialIDs) > 0 {
		staleTestimonialSet := make(map[uint64]bool, len(testimonialIDs))
		for _, id := range testimonialIDs {
			staleTestimonialSet[id] = true
		}
		cfg, err := s.repo.FindByKey("testimonial_showcase")
		if err == nil {
			var m map[string]json.RawMessage
			if json.Unmarshal(cfg.ConfigValue, &m) == nil {
				if raw, ok := m["featured_testimonial_ids"]; ok {
					var ids []uint64
					if json.Unmarshal(raw, &ids) == nil {
						cleaned := make([]uint64, 0, len(ids))
						for _, id := range ids {
							if !staleTestimonialSet[id] {
								cleaned = append(cleaned, id)
							}
						}
						if len(cleaned) < len(ids) {
							newRaw, _ := json.Marshal(cleaned)
							m["featured_testimonial_ids"] = newRaw
							updated, _ := json.Marshal(m)
							cfg.ConfigValue = updated
							_ = s.repo.Update(cfg)
						}
					}
				}
			}
		}
	}
	return nil
}

func sortByOrder[T any](items []T, keyFn func(T) int) {
	// insertion sort — items are small (typically < 10)
	for i := 1; i < len(items); i++ {
		j := i
		for j > 0 && keyFn(items[j]) < keyFn(items[j-1]) {
			items[j], items[j-1] = items[j-1], items[j]
			j--
		}
	}
}

// persistFromJSONArray writes a uint64 array back to a home_config JSON field.
func (s *HomeConfigService) persistFromJSONArray(configKey, fieldName string, ids []uint64) error {
	cfg, err := s.repo.FindByKey(configKey)
	if err != nil {
		return err
	}
	var m map[string]json.RawMessage
	if err := json.Unmarshal(cfg.ConfigValue, &m); err != nil {
		return err
	}
	raw, err := json.Marshal(ids)
	if err != nil {
		return err
	}
	m[fieldName] = raw
	updated, err := json.Marshal(m)
	if err != nil {
		return err
	}
	cfg.ConfigValue = updated
	return s.repo.Update(cfg)
}

// persistFromJSONArrayStr writes a string array back to a home_config JSON field.
func (s *HomeConfigService) persistFromJSONArrayStr(configKey, fieldName string, items []string) error {
	cfg, err := s.repo.FindByKey(configKey)
	if err != nil {
		return err
	}
	var m map[string]json.RawMessage
	if err := json.Unmarshal(cfg.ConfigValue, &m); err != nil {
		return err
	}
	raw, err := json.Marshal(items)
	if err != nil {
		return err
	}
	m[fieldName] = raw
	updated, err := json.Marshal(m)
	if err != nil {
		return err
	}
	cfg.ConfigValue = updated
	return s.repo.Update(cfg)
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
