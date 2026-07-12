package repository

import (
	"encoding/json"
	"time"

	"mygo-immigration/backend/internal/model"
)

// UserFilter holds optional filters for user queries.
type UserFilter struct {
	Role     string
	Status   *int8
	Username string
	Page     int
	PerPage  int
}

// UserRepository defines the interface for user data access.
type UserRepository interface {
	FindByUsername(username string) (*model.User, error)
	FindAll(filter UserFilter) ([]model.User, int64, error)
	Create(user *model.User) error
	Update(user *model.User) error
	FindByID(id uint64) (*model.User, error)
	PatchUpdate(id uint64, updates map[string]interface{}) error
	Delete(id uint64) error
}

// RoleFilter holds optional filters for role queries.
type RoleFilter struct {
	Status *int8
}

// PermissionOverrideInput represents one user-level permission override.
type PermissionOverrideInput struct {
	PermissionCode string
	Effect         string
}

// RBACRepository defines the interface for role and permission data access.
type RBACRepository interface {
	FindPermissions() ([]model.Permission, error)
	FindPermissionsByCodes(codes []string) ([]model.Permission, error)
	FindRoles(filter RoleFilter) ([]model.Role, error)
	FindRoleByID(id uint64) (*model.Role, error)
	FindRoleByCode(code string) (*model.Role, error)
	CreateRole(role *model.Role) error
	UpdateRole(role *model.Role) error
	DeleteRole(id uint64) error
	FindRolePermissionCodes(roleID uint64) ([]string, error)
	ReplaceRolePermissions(roleID uint64, permissionCodes []string) error
	FindUserPermissionOverrides(userID uint64) ([]model.UserPermissionOverride, error)
	ReplaceUserPermissionOverrides(userID uint64, overrides []PermissionOverrideInput) error
	FindEffectivePermissionCodes(userID uint64) ([]string, error)
}

// ProjectFilter holds optional query conditions for project list queries.
type ProjectFilter struct {
	Name    string
	Country string
	Status  string
	Page    int
	PerPage int
}

type ProjectOptionRow struct {
	ID   uint64
	Slug string
	Name string
}

// ProjectRepository defines the interface for project data access.
type ProjectRepository interface {
	FindByID(id uint64) (*model.Project, error)
	FindBySlug(slug string) (*model.Project, error)
	FindAll(filter ProjectFilter) ([]model.Project, int64, error)
	FindOptions(filter ProjectFilter) ([]ProjectOptionRow, int64, error)
	FindBySlugs(slugs []string) ([]model.Project, error)
	FindBySlugsLight(slugs []string) ([]model.Project, error)
	FindByIDsLight(ids []uint64) ([]model.Project, error)
	Create(project *model.Project) error
	Update(project *model.Project) error
	Delete(id uint64) error
	FindNews(projectID uint64) ([]model.Page, error)
	AddNews(projectID uint64, pageIDs []uint64) error
	RemoveNews(projectID, pageID uint64) error
	DeleteNewsByProjectID(projectID uint64) error
	FindAllCoverImages() ([]string, error)
	Count() (int64, error)
	CountByRange(start, end time.Time) (int64, error)
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
	FindByID(id uint64) (*model.FAQ, error)
	FindAll(params FAQQueryParams) ([]FAQWithProject, int64, error)
	FindDistinctProjects() ([]model.Project, error)
	Create(faq *model.FAQ) error
	Update(faq *model.FAQ) error
	Delete(id uint64) error
	DeleteByProjectID(projectID uint64) error
	Search(keyword string) ([]model.FAQ, error)
}

// PageFilter holds optional query conditions for page list queries.
type PageFilter struct {
	PageType string
	Title    string
	Status   string
	Page     int
	PerPage  int
}

type PageOptionRow struct {
	ID    uint64
	Slug  string
	Title string
}

// PageRepository defines the interface for page data access.
type PageRepository interface {
	FindByID(id uint64) (*model.Page, error)
	FindBySlug(slug string) (*model.Page, error)
	FindAll(filter PageFilter) ([]model.Page, int64, error)
	FindOptions(filter PageFilter) ([]PageOptionRow, int64, error)
	FindBySlugPublished(slug string) (*model.Page, error)
	FindProjectsByPageID(pageID uint64) ([]model.PageProject, error)
	FindRelatedBySlug(slug string, limit int) ([]model.Page, error)
	Create(page *model.Page) error
	Update(page *model.Page) error
	Delete(id uint64) error
	Search(keyword string) ([]model.Page, error)
	FindAllCoverImages() ([]string, error)
	FindAllContents() ([]string, error)
	Count() (int64, error)
	CountByRange(start, end time.Time) (int64, error)
}

// LeadFilter holds optional filters for lead queries.
type LeadFilter struct {
	Status            string
	Name              string
	Email             string
	InterestedProject string
	Page              int
	PerPage           int
}

// LeadRepository defines the interface for lead data access.
type LeadRepository interface {
	FindAll(filter LeadFilter) ([]model.Lead, int64, error)
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
	FindByParentID(parentID uint64) ([]model.Navigation, error)
	CountByProjectID(projectID uint64) (int64, error)
	CountByPageID(pageID uint64) (int64, error)
}

// CaseFilter holds optional filters for case queries.
type CaseFilter struct {
	ProjectID   *uint64
	CountryFrom string
	Name        string
	Page        int
	PerPage     int
}

type CaseOptionRow struct {
	ID   uint64
	Name string
}

// CaseRepository defines the interface for case data access.
type CaseRepository interface {
	FindByID(id uint64) (*model.Case, error)
	FindByIDs(ids []uint64) ([]model.Case, error)
	FindAll(filter CaseFilter) ([]model.Case, int64, error)
	FindOptions(filter CaseFilter) ([]CaseOptionRow, int64, error)
	FindBySlug(slug string) (*model.Case, error)
	Create(c *model.Case) error
	Update(c *model.Case) error
	Delete(id uint64) error
	DeleteByProjectID(projectID uint64) error
	FindAllPhotoURLs() ([]string, error)
	FindAllContents() ([]string, error)
	Count() (int64, error)
	CountByRange(start, end time.Time) (int64, error)
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
	DeleteByProjectID(projectID uint64) error
}

// CostItemRepository defines the interface for cost item data access.
type CostItemRepository interface {
	FindByProjectID(projectID uint64) ([]model.CostItem, error)
	Create(costItem *model.CostItem) error
	Update(costItem *model.CostItem) error
	Delete(id uint64) error
	DeleteByProjectID(projectID uint64) error
}

// TimelinePhaseRepository defines the interface for timeline phase data access.
type TimelinePhaseRepository interface {
	FindByProjectID(projectID uint64) ([]model.TimelinePhase, error)
	Create(phase *model.TimelinePhase) error
	Update(phase *model.TimelinePhase) error
	Delete(id uint64) error
	DeleteByProjectID(projectID uint64) error
}

// ProjectAdvantageRepository defines the interface for project advantage data access.
type ProjectAdvantageRepository interface {
	FindByProjectID(projectID uint64) ([]model.ProjectAdvantage, error)
	Create(adv *model.ProjectAdvantage) error
	Update(adv *model.ProjectAdvantage) error
	Delete(id uint64) error
	DeleteByProjectID(projectID uint64) error
}

// MilestoneRepository defines the interface for milestone data access.
type MilestoneRepository interface {
	DeleteByProjectID(projectID uint64) error
}

// TestimonialFilter holds optional filters for testimonial queries.
type TestimonialFilter struct {
	ProjectID *uint64
	Nickname  string
	Rating    *uint8
	Page      int
	PerPage   int
}

type TestimonialOptionRow struct {
	ID       uint64
	Nickname string
}

// TestimonialRepository defines the interface for testimonial data access.
type TestimonialRepository interface {
	FindByID(id uint64) (*model.Testimonial, error)
	FindByIDs(ids []uint64) ([]model.Testimonial, error)
	FindAll(filter TestimonialFilter) ([]model.Testimonial, int64, error)
	FindOptions(filter TestimonialFilter) ([]TestimonialOptionRow, int64, error)
	Create(t *model.Testimonial) error
	Update(t *model.Testimonial) error
	Delete(id uint64) error
	DeleteByProjectID(projectID uint64) error
	FindAllAvatarURLs() ([]string, error)
}

// HomeConfigRepository defines the interface for home config data access.
type HomeConfigRepository interface {
	FindByKey(key string) (*model.HomeConfig, error)
	FindAll() ([]model.HomeConfig, error)
	Create(cfg *model.HomeConfig) error
	Update(cfg *model.HomeConfig) error
	FindAllConfigValues() ([]json.RawMessage, error)
}

// MediaRepository defines the interface for media data access.
type MediaFilter struct {
	Search  string
	Page    int
	PerPage int
}

type MediaRepository interface {
	FindAll(filter MediaFilter) ([]model.Media, int64, error)
	FindByID(id uint64) (*model.Media, error)
	Create(media *model.Media) error
	Delete(id uint64) error
	DeleteByIDPermanently(id uint64) error
}

// LawyerFilter holds optional query conditions for lawyer list queries.
type LawyerFilter struct {
	Name    string
	Page    int
	PerPage int
}

// LawyerRepository defines the interface for lawyer data access.
type LawyerRepository interface {
	FindAll(filter LawyerFilter) ([]model.Lawyer, int64, error)
	FindByID(id uint64) (*model.Lawyer, error)
	Create(item *model.Lawyer) error
	Update(item *model.Lawyer) error
	Delete(id uint64) error
	FindAllPhotoURLs() ([]string, error)
	Count() (int64, error)
	CountByRange(start, end time.Time) (int64, error)
}
