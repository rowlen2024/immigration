package repository

import "mygo-immigration/backend/internal/model"

// UserRepository defines the interface for user data access.
type UserRepository interface {
	FindByUsername(username string) (*model.User, error)
	FindAll() ([]model.User, error)
	FindAllPaginated(page, perPage int) ([]model.User, int64, error)
	Create(user *model.User) error
	Update(user *model.User) error
	FindByID(id uint64) (*model.User, error)
	PatchUpdate(id uint64, updates map[string]interface{}) error
	Delete(id uint64) error
}

// ProjectRepository defines the interface for project data access.
type ProjectRepository interface {
	FindBySlug(slug string) (*model.Project, error)
	FindAll(page, perPage int, search, status string) ([]model.Project, int64, error)
	FindBySlugs(slugs []string) ([]model.Project, error)
	Create(project *model.Project) error
	Update(project *model.Project) error
	Delete(id uint64) error
}

// FAQQueryParams holds optional filters for FAQ queries.
type FAQQueryParams struct {
	ProjectID *uint64
	IsGlobal  *bool
	Search    string
	Page      int
	PerPage   int
}

// FAQRepository defines the interface for FAQ data access.
type FAQRepository interface {
	FindAll(params FAQQueryParams) ([]FAQWithProject, int64, error)
	Create(faq *model.FAQ) error
	Update(faq *model.FAQ) error
	Delete(id uint64) error
	Search(keyword string) ([]model.FAQ, error)
}

// PageRepository defines the interface for page data access.
type PageRepository interface {
	FindBySlug(slug string) (*model.Page, error)
	FindAll(pageType, search, status string) ([]model.Page, error)
	FindAllPublished() ([]model.Page, error)
	FindBySlugPublished(slug string) (*model.Page, error)
	FindByProjectID(projectID uint64) ([]model.Page, error)
	Create(page *model.Page) error
	Update(page *model.Page) error
	Delete(id uint64) error
	Search(keyword string) ([]model.Page, error)
}

// LeadRepository defines the interface for lead data access.
type LeadRepository interface {
	FindAll(page, perPage int, status string) ([]model.Lead, int64, error)
	Create(lead *model.Lead) error
	UpdateStatus(id uint64, status string, notes string) error
	Delete(id uint64) error
}

// NavigationRepository defines the interface for navigation data access.
type NavigationRepository interface {
	FindAll() ([]model.Navigation, error)
	FindAllActive() ([]model.Navigation, error)
	FindAllActiveByPosition(position string) ([]model.Navigation, error)
	FindByID(id uint64) (*model.Navigation, error)
	Create(nav *model.Navigation) error
	Update(nav *model.Navigation) error
	Delete(id uint64) error
	HasChildren(parentID uint64) (bool, error)
	CountByProjectID(projectID uint64) (int64, error)
	CountByPageID(pageID uint64) (int64, error)
}

// CaseRepository defines the interface for case data access.
type CaseRepository interface {
	FindByProjectID(projectID uint64) ([]model.Case, error)
	FindAll(search string) ([]model.Case, error)
	Create(c *model.Case) error
	Update(c *model.Case) error
	Delete(id uint64) error
	HardDelete(id uint64) error
}

// CompareConfigRepository defines the interface for compare config data access.
type CompareConfigRepository interface {
	FindByProjectID(projectID uint64) (*model.CompareConfig, error)
	Upsert(cfg *model.CompareConfig) error
	DeleteByProjectID(projectID uint64) error
}

// RequirementRepository defines the interface for requirement data access.
type RequirementRepository interface {
	FindByProjectID(projectID uint64) ([]model.Requirement, error)
	Create(requirement *model.Requirement) error
	Update(requirement *model.Requirement) error
	Delete(id uint64) error
}

// CostItemRepository defines the interface for cost item data access.
type CostItemRepository interface {
	FindByProjectID(projectID uint64) ([]model.CostItem, error)
	Create(costItem *model.CostItem) error
	Update(costItem *model.CostItem) error
	Delete(id uint64) error
}

// TimelinePhaseRepository defines the interface for timeline phase data access.
type TimelinePhaseRepository interface {
	FindByProjectID(projectID uint64) ([]model.TimelinePhase, error)
	Create(phase *model.TimelinePhase) error
	Update(phase *model.TimelinePhase) error
	Delete(id uint64) error
}
