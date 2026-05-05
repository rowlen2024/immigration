# Project Detail Admin & Frontend-Backend Alignment — Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add sub-entity CRUD (requirements/cost_items/timeline_phases) to admin project editor, unify project detail pages under `[slug].vue`, fix compare API, and remove hardcoded fallback data.

**Architecture:** Backend follows existing layered pattern (Handler→Service→Repository→Model). Three new repos/services/handlers for sub-entities. Frontend admin dialog converts to tabs. Dedicated project pages deleted, redirects added.

**Tech Stack:** Go/Gin/GORM backend, Nuxt 3/Element Plus frontend.

---

### Task 1: Backend — Sub-entity Repository

**Files:**
- Create: `backend/internal/repository/project_sub_repo.go`
- Modify: `backend/internal/repository/interfaces.go`
- Modify: `backend/internal/repository/repository.go`

- [ ] **Step 1: Add sub-entity repository interfaces**

In `backend/internal/repository/interfaces.go`, append after the `CaseRepository` interface:

```go
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
```

- [ ] **Step 2: Create sub-entity repository implementations**

Create `backend/internal/repository/project_sub_repo.go`:

```go
package repository

import (
	"mygo-immigration/backend/internal/model"

	"gorm.io/gorm"
)

type RequirementRepo struct {
	db *gorm.DB
}

func (r *RequirementRepo) FindByProjectID(projectID uint64) ([]model.Requirement, error) {
	var items []model.Requirement
	err := r.db.Where("project_id = ?", projectID).Order("sort_order asc").Find(&items).Error
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (r *RequirementRepo) Create(item *model.Requirement) error {
	return r.db.Create(item).Error
}

func (r *RequirementRepo) Update(item *model.Requirement) error {
	return r.db.Save(item).Error
}

func (r *RequirementRepo) Delete(id uint64) error {
	return r.db.Delete(&model.Requirement{}, id).Error
}

type CostItemRepo struct {
	db *gorm.DB
}

func (r *CostItemRepo) FindByProjectID(projectID uint64) ([]model.CostItem, error) {
	var items []model.CostItem
	err := r.db.Where("project_id = ?", projectID).Order("sort_order asc").Find(&items).Error
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (r *CostItemRepo) Create(item *model.CostItem) error {
	return r.db.Create(item).Error
}

func (r *CostItemRepo) Update(item *model.CostItem) error {
	return r.db.Save(item).Error
}

func (r *CostItemRepo) Delete(id uint64) error {
	return r.db.Delete(&model.CostItem{}, id).Error
}

type TimelinePhaseRepo struct {
	db *gorm.DB
}

func (r *TimelinePhaseRepo) FindByProjectID(projectID uint64) ([]model.TimelinePhase, error) {
	var items []model.TimelinePhase
	err := r.db.Where("project_id = ?", projectID).Order("sort_order asc").Find(&items).Error
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (r *TimelinePhaseRepo) Create(item *model.TimelinePhase) error {
	return r.db.Create(item).Error
}

func (r *TimelinePhaseRepo) Update(item *model.TimelinePhase) error {
	return r.db.Save(item).Error
}

func (r *TimelinePhaseRepo) Delete(id uint64) error {
	return r.db.Delete(&model.TimelinePhase{}, id).Error
}
```

- [ ] **Step 3: Wire into Repository struct**

In `backend/internal/repository/repository.go`, add fields to the struct:

```go
type Repository struct {
	Project        *ProjectRepo
	User           *UserRepo
	FAQ            *FAQRepo
	Case           *CaseRepo
	Page           *PageRepo
	Lead           *LeadRepo
	HomeConfig     *HomeConfigRepo
	Media          *MediaRepo
	Nav            *NavRepo
	Requirement    *RequirementRepo
	CostItem       *CostItemRepo
	TimelinePhase  *TimelinePhaseRepo
}
```

And in `New()`:

```go
func New(db *gorm.DB) *Repository {
	return &Repository{
		Project:       &ProjectRepo{db: db},
		User:          &UserRepo{db: db},
		FAQ:           &FAQRepo{db: db},
		Case:          &CaseRepo{db: db},
		Page:          &PageRepo{db: db},
		Lead:          &LeadRepo{db: db},
		HomeConfig:    &HomeConfigRepo{db: db},
		Media:         &MediaRepo{db: db},
		Nav:           &NavRepo{db: db},
		Requirement:   &RequirementRepo{db: db},
		CostItem:      &CostItemRepo{db: db},
		TimelinePhase: &TimelinePhaseRepo{db: db},
	}
}
```

- [ ] **Step 4: Run Go build to verify**

```bash
cd backend && go build ./...
```

Expected: build succeeds.

---

### Task 2: Backend — Sub-entity Service

**Files:**
- Create: `backend/internal/service/project_sub_svc.go`
- Modify: `backend/internal/service/service.go`

- [ ] **Step 1: Create sub-entity services**

Create `backend/internal/service/project_sub_svc.go`:

```go
package service

import (
	"errors"
	"fmt"

	"mygo-immigration/backend/internal/model"
	"mygo-immigration/backend/internal/repository"
)

type RequirementService struct {
	repo repository.RequirementRepository
}

func NewRequirementService(repo repository.RequirementRepository) *RequirementService {
	return &RequirementService{repo: repo}
}

func (s *RequirementService) List(projectID uint64) ([]model.Requirement, error) {
	items, err := s.repo.FindByProjectID(projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to list requirements: %w", err)
	}
	return items, nil
}

func (s *RequirementService) Create(projectID uint64, item *model.Requirement) (*model.Requirement, error) {
	if item == nil {
		return nil, errors.New("requirement is nil")
	}
	if item.Label == "" {
		return nil, errors.New("label is required")
	}
	item.ProjectID = projectID
	if err := s.repo.Create(item); err != nil {
		return nil, fmt.Errorf("failed to create requirement: %w", err)
	}
	return item, nil
}

func (s *RequirementService) Update(projectID uint64, id uint64, item *model.Requirement) (*model.Requirement, error) {
	if item == nil {
		return nil, errors.New("requirement is nil")
	}
	if id == 0 {
		return nil, errors.New("id is required")
	}
	item.ID = id
	item.ProjectID = projectID
	if err := s.repo.Update(item); err != nil {
		return nil, fmt.Errorf("failed to update requirement: %w", err)
	}
	return item, nil
}

func (s *RequirementService) Delete(projectID uint64, id uint64) error {
	if id == 0 {
		return errors.New("id is required")
	}
	_ = projectID
	if err := s.repo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete requirement: %w", err)
	}
	return nil
}

type CostItemService struct {
	repo repository.CostItemRepository
}

func NewCostItemService(repo repository.CostItemRepository) *CostItemService {
	return &CostItemService{repo: repo}
}

func (s *CostItemService) List(projectID uint64) ([]model.CostItem, error) {
	items, err := s.repo.FindByProjectID(projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to list cost items: %w", err)
	}
	return items, nil
}

func (s *CostItemService) Create(projectID uint64, item *model.CostItem) (*model.CostItem, error) {
	if item == nil {
		return nil, errors.New("cost item is nil")
	}
	if item.Name == "" {
		return nil, errors.New("name is required")
	}
	item.ProjectID = projectID
	if err := s.repo.Create(item); err != nil {
		return nil, fmt.Errorf("failed to create cost item: %w", err)
	}
	return item, nil
}

func (s *CostItemService) Update(projectID uint64, id uint64, item *model.CostItem) (*model.CostItem, error) {
	if item == nil {
		return nil, errors.New("cost item is nil")
	}
	if id == 0 {
		return nil, errors.New("id is required")
	}
	item.ID = id
	item.ProjectID = projectID
	if err := s.repo.Update(item); err != nil {
		return nil, fmt.Errorf("failed to update cost item: %w", err)
	}
	return item, nil
}

func (s *CostItemService) Delete(projectID uint64, id uint64) error {
	if id == 0 {
		return errors.New("id is required")
	}
	_ = projectID
	if err := s.repo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete cost item: %w", err)
	}
	return nil
}

type TimelinePhaseService struct {
	repo repository.TimelinePhaseRepository
}

func NewTimelinePhaseService(repo repository.TimelinePhaseRepository) *TimelinePhaseService {
	return &TimelinePhaseService{repo: repo}
}

func (s *TimelinePhaseService) List(projectID uint64) ([]model.TimelinePhase, error) {
	items, err := s.repo.FindByProjectID(projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to list timeline phases: %w", err)
	}
	return items, nil
}

func (s *TimelinePhaseService) Create(projectID uint64, item *model.TimelinePhase) (*model.TimelinePhase, error) {
	if item == nil {
		return nil, errors.New("timeline phase is nil")
	}
	if item.Title == "" {
		return nil, errors.New("title is required")
	}
	item.ProjectID = projectID
	if err := s.repo.Create(item); err != nil {
		return nil, fmt.Errorf("failed to create timeline phase: %w", err)
	}
	return item, nil
}

func (s *TimelinePhaseService) Update(projectID uint64, id uint64, item *model.TimelinePhase) (*model.TimelinePhase, error) {
	if item == nil {
		return nil, errors.New("timeline phase is nil")
	}
	if id == 0 {
		return nil, errors.New("id is required")
	}
	item.ID = id
	item.ProjectID = projectID
	if err := s.repo.Update(item); err != nil {
		return nil, fmt.Errorf("failed to update timeline phase: %w", err)
	}
	return item, nil
}

func (s *TimelinePhaseService) Delete(projectID uint64, id uint64) error {
	if id == 0 {
		return errors.New("id is required")
	}
	_ = projectID
	if err := s.repo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete timeline phase: %w", err)
	}
	return nil
}
```

- [ ] **Step 2: Wire into Service struct**

In `backend/internal/service/service.go`, add fields to the struct:

```go
type Service struct {
	Project        *ProjectService
	Auth           *AuthService
	User           *UserService
	FAQ            *FAQService
	Page           *PageService
	Case           *CaseService
	Lead           *LeadService
	HomeConfig     *HomeConfigService
	Media          *MediaService
	Nav            *NavService
	Search         *SearchService
	Requirement    *RequirementService
	CostItem       *CostItemService
	TimelinePhase  *TimelinePhaseService
}
```

And in `New()`:

```go
func New(repo *repository.Repository, cfg *config.Config) *Service {
	return &Service{
		Project:       &ProjectService{repo: repo.Project, navRepo: repo.Nav},
		Auth:          &AuthService{repo: repo.User, cfg: cfg},
		User:          &UserService{repo: repo.User},
		FAQ:           &FAQService{repo: repo.FAQ},
		Page:          &PageService{repo: repo.Page, navRepo: repo.Nav},
		Case:          &CaseService{repo: repo.Case},
		Lead:          &LeadService{repo: repo.Lead},
		HomeConfig:    &HomeConfigService{repo: repo.HomeConfig},
		Media:         &MediaService{repo: repo.Media},
		Nav:           &NavService{repo: repo.Nav, projectRepo: repo.Project, pageRepo: repo.Page},
		Search:        &SearchService{faqRepo: repo.FAQ, pageRepo: repo.Page},
		Requirement:   &RequirementService{repo: repo.Requirement},
		CostItem:      &CostItemService{repo: repo.CostItem},
		TimelinePhase: &TimelinePhaseService{repo: repo.TimelinePhase},
	}
}
```

- [ ] **Step 3: Run Go build to verify**

```bash
cd backend && go build ./...
```

Expected: build succeeds.

---

### Task 3: Backend — Sub-entity Handler & Routes

**Files:**
- Create: `backend/internal/handler/project_sub_handler.go`
- Modify: `backend/internal/router/router.go`

- [ ] **Step 1: Create sub-entity handlers**

Create `backend/internal/handler/project_sub_handler.go`:

```go
package handler

import (
	"net/http"

	"mygo-immigration/backend/internal/dto"
	"mygo-immigration/backend/internal/model"

	"github.com/gin-gonic/gin"
)

// Requirements

func (h *Handler) ListRequirements(c *gin.Context) {
	projectID, err := parseIDParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "invalid project id"))
		return
	}
	items, err := h.svc.Requirement.List(projectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Error(500, err.Error()))
		return
	}
	c.JSON(http.StatusOK, dto.Success(items))
}

func (h *Handler) CreateRequirement(c *gin.Context) {
	projectID, err := parseIDParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "invalid project id"))
		return
	}
	var item model.Requirement
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "invalid request"))
		return
	}
	created, err := h.svc.Requirement.Create(projectID, &item)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Error(500, err.Error()))
		return
	}
	c.JSON(http.StatusCreated, dto.Success(created))
}

func (h *Handler) UpdateRequirement(c *gin.Context) {
	projectID, err := parseIDParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "invalid project id"))
		return
	}
	rid, err := parseIDParam(c, "rid")
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "invalid requirement id"))
		return
	}
	var item model.Requirement
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "invalid request"))
		return
	}
	updated, err := h.svc.Requirement.Update(projectID, rid, &item)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Error(500, err.Error()))
		return
	}
	c.JSON(http.StatusOK, dto.Success(updated))
}

func (h *Handler) DeleteRequirement(c *gin.Context) {
	projectID, err := parseIDParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "invalid project id"))
		return
	}
	rid, err := parseIDParam(c, "rid")
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "invalid requirement id"))
		return
	}
	if err := h.svc.Requirement.Delete(projectID, rid); err != nil {
		c.JSON(http.StatusInternalServerError, dto.Error(500, err.Error()))
		return
	}
	c.JSON(http.StatusOK, dto.Success(nil))
}

// Cost Items

func (h *Handler) ListCostItems(c *gin.Context) {
	projectID, err := parseIDParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "invalid project id"))
		return
	}
	items, err := h.svc.CostItem.List(projectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Error(500, err.Error()))
		return
	}
	c.JSON(http.StatusOK, dto.Success(items))
}

func (h *Handler) CreateCostItem(c *gin.Context) {
	projectID, err := parseIDParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "invalid project id"))
		return
	}
	var item model.CostItem
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "invalid request"))
		return
	}
	created, err := h.svc.CostItem.Create(projectID, &item)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Error(500, err.Error()))
		return
	}
	c.JSON(http.StatusCreated, dto.Success(created))
}

func (h *Handler) UpdateCostItem(c *gin.Context) {
	projectID, err := parseIDParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "invalid project id"))
		return
	}
	cid, err := parseIDParam(c, "cid")
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "invalid cost item id"))
		return
	}
	var item model.CostItem
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "invalid request"))
		return
	}
	updated, err := h.svc.CostItem.Update(projectID, cid, &item)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Error(500, err.Error()))
		return
	}
	c.JSON(http.StatusOK, dto.Success(updated))
}

func (h *Handler) DeleteCostItem(c *gin.Context) {
	projectID, err := parseIDParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "invalid project id"))
		return
	}
	cid, err := parseIDParam(c, "cid")
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "invalid cost item id"))
		return
	}
	if err := h.svc.CostItem.Delete(projectID, cid); err != nil {
		c.JSON(http.StatusInternalServerError, dto.Error(500, err.Error()))
		return
	}
	c.JSON(http.StatusOK, dto.Success(nil))
}

// Timeline Phases

func (h *Handler) ListTimelinePhases(c *gin.Context) {
	projectID, err := parseIDParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "invalid project id"))
		return
	}
	items, err := h.svc.TimelinePhase.List(projectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Error(500, err.Error()))
		return
	}
	c.JSON(http.StatusOK, dto.Success(items))
}

func (h *Handler) CreateTimelinePhase(c *gin.Context) {
	projectID, err := parseIDParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "invalid project id"))
		return
	}
	var item model.TimelinePhase
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "invalid request"))
		return
	}
	created, err := h.svc.TimelinePhase.Create(projectID, &item)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Error(500, err.Error()))
		return
	}
	c.JSON(http.StatusCreated, dto.Success(created))
}

func (h *Handler) UpdateTimelinePhase(c *gin.Context) {
	projectID, err := parseIDParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "invalid project id"))
		return
	}
	tid, err := parseIDParam(c, "tid")
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "invalid timeline phase id"))
		return
	}
	var item model.TimelinePhase
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "invalid request"))
		return
	}
	updated, err := h.svc.TimelinePhase.Update(projectID, tid, &item)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Error(500, err.Error()))
		return
	}
	c.JSON(http.StatusOK, dto.Success(updated))
}

func (h *Handler) DeleteTimelinePhase(c *gin.Context) {
	projectID, err := parseIDParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "invalid project id"))
		return
	}
	tid, err := parseIDParam(c, "tid")
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "invalid timeline phase id"))
		return
	}
	if err := h.svc.TimelinePhase.Delete(projectID, tid); err != nil {
		c.JSON(http.StatusInternalServerError, dto.Error(500, err.Error()))
		return
	}
	c.JSON(http.StatusOK, dto.Success(nil))
}
```

- [ ] **Step 2: Register routes**

In `backend/internal/router/router.go`, inside the `admin := api.Group("/admin")` block, after the projects CRUD group, add:

```go
				// Sub-entity routes nested under projects
				projects.GET("/:id/requirements", middleware.RBAC("admin:read"), h.ListRequirements)
				projects.POST("/:id/requirements", middleware.RBAC("projects:write"), h.CreateRequirement)
				projects.PUT("/:id/requirements/:rid", middleware.RBAC("projects:write"), h.UpdateRequirement)
				projects.DELETE("/:id/requirements/:rid", middleware.RBAC("projects:write"), h.DeleteRequirement)

				projects.GET("/:id/cost-items", middleware.RBAC("admin:read"), h.ListCostItems)
				projects.POST("/:id/cost-items", middleware.RBAC("projects:write"), h.CreateCostItem)
				projects.PUT("/:id/cost-items/:cid", middleware.RBAC("projects:write"), h.UpdateCostItem)
				projects.DELETE("/:id/cost-items/:cid", middleware.RBAC("projects:write"), h.DeleteCostItem)

				projects.GET("/:id/timeline-phases", middleware.RBAC("admin:read"), h.ListTimelinePhases)
				projects.POST("/:id/timeline-phases", middleware.RBAC("projects:write"), h.CreateTimelinePhase)
				projects.PUT("/:id/timeline-phases/:tid", middleware.RBAC("projects:write"), h.UpdateTimelinePhase)
				projects.DELETE("/:id/timeline-phases/:tid", middleware.RBAC("projects:write"), h.DeleteTimelinePhase)
```

These go inside the `projects` route group (after the existing `projects.DELETE("/:id", ...)` line).

- [ ] **Step 3: Run Go build to verify**

```bash
cd backend && go build ./...
```

Expected: build succeeds.

---

### Task 4: Backend — Fix Project Update (Don't Touch Associations)

**File:** Modify: `backend/internal/repository/project_repo.go`

- [ ] **Step 1: Change Update to omit associations**

In `backend/internal/repository/project_repo.go`, change the `Update` method from:

```go
func (r *ProjectRepo) Update(project *model.Project) error {
	return r.db.Save(project).Error
}
```

To:

```go
func (r *ProjectRepo) Update(project *model.Project) error {
	return r.db.Omit("Requirements", "CostItems", "TimelinePhases", "Milestones", "FAQs", "Cases").Save(project).Error
}
```

- [ ] **Step 2: Run Go build**

```bash
cd backend && go build ./...
```

Expected: build succeeds.

---

### Task 5: Backend — Fix Compare API

**Files:**
- Modify: `backend/internal/service/project_svc.go`
- Modify: `backend/internal/handler/project_handler.go`

- [ ] **Step 1: Add CompareRows to project service**

In `backend/internal/service/project_svc.go`, after the existing `Compare` method, add:

```go
// CompareRow represents a single comparison row.
type CompareRow struct {
	Label string `json:"label"`
	A     string `json:"a"`
	B     string `json:"b"`
}

// CompareResult holds the full comparison output.
type CompareResult struct {
	Projects []CompareProject `json:"projects"`
	Rows     []CompareRow     `json:"rows"`
}

// CompareProject holds minimal project info for the comparison header.
type CompareProject struct {
	Title string `json:"title"`
	Slug  string `json:"slug"`
}

// CompareRows returns formatted comparison rows for two projects.
func (s *ProjectService) CompareRows(slugs []string) (*CompareResult, error) {
	projects, err := s.Compare(slugs)
	if err != nil {
		return nil, err
	}
	if len(projects) < 2 {
		return nil, errors.New("需要两个项目进行对比")
	}

	a, b := projects[0], projects[1]

	projInfo := []CompareProject{
		{Title: a.Name, Slug: a.Slug},
		{Title: b.Name, Slug: b.Slug},
	}

	reqsA := joinRequirements(a.Requirements)
	reqsB := joinRequirements(b.Requirements)

	phaseCountA := fmt.Sprintf("%d 个阶段", len(a.TimelinePhases))
	phaseCountB := fmt.Sprintf("%d 个阶段", len(b.TimelinePhases))

	rows := []CompareRow{
		{Label: "投资金额", A: a.InvestmentAmount, B: b.InvestmentAmount},
		{Label: "办理周期", A: a.ProcessingPeriod, B: b.ProcessingPeriod},
		{Label: "适合人群", A: a.TargetCrowd, B: b.TargetCrowd},
		{Label: "申请条件", A: reqsA, B: reqsB},
		{Label: "费用总计", A: a.CostsTotal, B: b.CostsTotal},
		{Label: "流程步骤", A: phaseCountA, B: phaseCountB},
	}

	return &CompareResult{Projects: projInfo, Rows: rows}, nil
}

func joinRequirements(reqs []model.Requirement) string {
	if len(reqs) == 0 {
		return ""
	}
	labels := make([]string, len(reqs))
	for i, r := range reqs {
		prefix := ""
		if r.IsRequired {
			prefix = "✓ "
		} else {
			prefix = "○ "
		}
		labels[i] = prefix + r.Label
	}
	return strings.Join(labels, "；")
}
```

Add `"strings"` and `"errors"` to the imports if not already present (errors is already there, verify strings).

- [ ] **Step 2: Update CompareProjects handler**

In `backend/internal/handler/project_handler.go`, change the `CompareProjects` method from:

```go
func (h *Handler) CompareProjects(c *gin.Context) {
	slugsParam := c.Query("slugs")
	if slugsParam == "" {
		c.JSON(http.StatusBadRequest, dto.Error(400, "slugs query param is required"))
		return
	}

	slugs := strings.Split(slugsParam, ",")

	projects, err := h.svc.Project.Compare(slugs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Error(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.Success(projects))
}
```

To:

```go
func (h *Handler) CompareProjects(c *gin.Context) {
	slugsParam := c.Query("slugs")
	if slugsParam == "" {
		c.JSON(http.StatusBadRequest, dto.Error(400, "slugs query param is required"))
		return
	}

	slugs := strings.Split(slugsParam, ",")

	result, err := h.svc.Project.CompareRows(slugs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Error(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.Success(result))
}
```

- [ ] **Step 3: Run Go build**

```bash
cd backend && go build ./...
```

Expected: build succeeds.

---

### Task 6: Frontend — Delete Dedicated Pages & Add Redirects

**Files:**
- Delete: `frontend/pages/usa/eb5.vue`
- Delete: `frontend/pages/hongkong/cies.vue`
- Delete: `frontend/pages/panama/property.vue`
- Modify: `frontend/nuxt.config.ts`

- [ ] **Step 1: Delete the three dedicated project pages**

```bash
rm frontend/pages/usa/eb5.vue
rm frontend/pages/hongkong/cies.vue
rm frontend/pages/panama/property.vue
```

- [ ] **Step 2: Check if parent directories are now empty and remove them**

```bash
rmdir frontend/pages/usa 2>/dev/null; rmdir frontend/pages/hongkong 2>/dev/null; rmdir frontend/pages/panama 2>/dev/null; echo "done"
```

- [ ] **Step 3: Add route redirects to nuxt.config.ts**

In `frontend/nuxt.config.ts`, add a `routeRules` section inside the `nitro: {}` block. Change:

```ts
  nitro: {
    devProxy: {
      '/api': 'http://localhost:8080',
      '/uploads': 'http://localhost:8080',
    },
  },
```

To:

```ts
  nitro: {
    devProxy: {
      '/api': 'http://localhost:8080',
      '/uploads': 'http://localhost:8080',
    },
    routeRules: {
      '/usa/eb5': { redirect: '/projects/eb5' },
      '/hongkong/cies': { redirect: '/projects/cies' },
      '/panama/property': { redirect: '/projects/panama' },
    },
  },
```

- [ ] **Step 4: Run Nuxt typecheck**

```bash
cd frontend && npx nuxi typecheck
```

Expected: no new errors.

---

### Task 7: Frontend — Admin Project Editor with Tabs

**File:** Modify: `frontend/pages/admin/projects.vue`

- [ ] **Step 1: Replace the dialog template section**

Replace the entire `<el-dialog>` block (lines 46-168) with the tabbed version below.

Find the `<el-dialog>` opening tag and replace everything from `<el-dialog` through `</el-dialog>` with:

```vue
    <el-dialog
      v-model="dialogVisible"
      :title="editingId ? '编辑项目' : '新建项目'"
      width="900px"
      destroy-on-close
      @opened="onDialogOpened"
    >
      <el-tabs v-model="activeTab" type="border-card">
        <!-- Tab 1: 基本信息 -->
        <el-tab-pane label="基本信息" name="basic">
          <el-form ref="formRef" :model="form" :rules="rules" label-width="120px">
            <el-row :gutter="16">
              <el-col :span="12">
                <el-form-item label="标识(slug)" prop="slug">
                  <el-input v-model="form.slug" />
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item label="项目名称" prop="name">
                  <el-input v-model="form.name" />
                </el-form-item>
              </el-col>
            </el-row>
            <el-row :gutter="16">
              <el-col :span="12">
                <el-form-item label="国家" prop="country">
                  <el-input v-model="form.country" />
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item label="国旗图标" prop="flag_emoji">
                  <el-input v-model="form.flag_emoji" />
                </el-form-item>
              </el-col>
            </el-row>
            <el-form-item label="标语" prop="tagline">
              <el-input v-model="form.tagline" />
            </el-form-item>
            <el-row :gutter="16">
              <el-col :span="12">
                <el-form-item label="投资金额" prop="investment_amount">
                  <el-input v-model="form.investment_amount" />
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item label="投资价值" prop="investment_value">
                  <el-input v-model="form.investment_value" />
                </el-form-item>
              </el-col>
            </el-row>
            <el-form-item label="办理周期" prop="processing_period">
              <el-input v-model="form.processing_period" />
            </el-form-item>
            <el-form-item label="目标人群" prop="target_crowd">
              <el-input v-model="form.target_crowd" />
            </el-form-item>
            <el-form-item label="概览标题" prop="overview_title">
              <el-input v-model="form.overview_title" />
            </el-form-item>
            <el-form-item label="概览内容" prop="overview_text">
              <el-input v-model="form.overview_text" type="textarea" :rows="3" />
            </el-form-item>
            <el-form-item label="政策标题" prop="policy_title">
              <el-input v-model="form.policy_title" />
            </el-form-item>
            <el-form-item label="政策内容" prop="policy_text">
              <el-input v-model="form.policy_text" type="textarea" :rows="3" />
            </el-form-item>
            <el-row :gutter="16">
              <el-col :span="12">
                <el-form-item label="费用总计" prop="costs_total">
                  <el-input v-model="form.costs_total" />
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item label="费用备注" prop="costs_note">
                  <el-input v-model="form.costs_note" />
                </el-form-item>
              </el-col>
            </el-row>
            <el-form-item label="CTA 文字" prop="cta_text">
              <el-input v-model="form.cta_text" />
            </el-form-item>
            <el-form-item label="Hero 标题" prop="hero_title">
              <el-input v-model="form.hero_title" />
            </el-form-item>
            <el-form-item label="封面图片">
              <div style="display: flex; gap: 8px; width: 100%">
                <el-input v-model="form.cover_image" placeholder="图片 URL 或上传" style="flex: 1" />
                <el-upload
                  :action="uploadUrl"
                  :headers="uploadHeaders"
                  accept=".jpg,.jpeg,.png,.webp"
                  :show-file-list="false"
                  :on-success="handleCoverUploadSuccess"
                >
                  <el-button>上传</el-button>
                </el-upload>
              </div>
            </el-form-item>
            <el-form-item label="Hero 描述" prop="hero_desc">
              <el-input v-model="form.hero_desc" type="textarea" :rows="2" />
            </el-form-item>
            <el-form-item label="Hero 渐变" prop="hero_gradient">
              <el-input v-model="form.hero_gradient" />
            </el-form-item>
            <el-row :gutter="16">
              <el-col :span="12">
                <el-form-item label="排序" prop="sort_order">
                  <el-input-number v-model="form.sort_order" :min="0" />
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item label="状态" prop="status">
                  <el-select v-model="form.status">
                    <el-option label="草稿" :value="0" />
                    <el-option label="已发布" :value="1" />
                  </el-select>
                </el-form-item>
              </el-col>
            </el-row>
          </el-form>
        </el-tab-pane>

        <!-- Tab 2: 申请条件 -->
        <el-tab-pane label="申请条件" name="requirements">
          <div style="margin-bottom: 12px">
            <el-button type="primary" size="small" @click="openSubDialog('requirement')">添加条件</el-button>
          </div>
          <el-table :data="subData.requirements" border stripe>
            <el-table-column prop="label" label="条件描述" min-width="200" />
            <el-table-column label="是否必需" width="100">
              <template #default="{ row }">
                <el-tag :type="row.is_required ? 'success' : 'info'">
                  {{ row.is_required ? '必需' : '可选' }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="sort_order" label="排序" width="80" />
            <el-table-column label="操作" width="150">
              <template #default="{ row }">
                <el-button size="small" @click="openSubDialog('requirement', row)">编辑</el-button>
                <el-popconfirm title="确定删除？" @confirm="deleteSubItem('requirement', row.id)">
                  <template #reference>
                    <el-button size="small" type="danger">删除</el-button>
                  </template>
                </el-popconfirm>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>

        <!-- Tab 3: 费用明细 -->
        <el-tab-pane label="费用明细" name="costItems">
          <div style="margin-bottom: 12px">
            <el-button type="primary" size="small" @click="openSubDialog('costItem')">添加费用</el-button>
          </div>
          <el-table :data="subData.costItems" border stripe>
            <el-table-column prop="name" label="费用名称" min-width="150" />
            <el-table-column prop="amount" label="金额" width="120" />
            <el-table-column prop="note" label="说明" min-width="180" />
            <el-table-column prop="sort_order" label="排序" width="80" />
            <el-table-column label="操作" width="150">
              <template #default="{ row }">
                <el-button size="small" @click="openSubDialog('costItem', row)">编辑</el-button>
                <el-popconfirm title="确定删除？" @confirm="deleteSubItem('costItem', row.id)">
                  <template #reference>
                    <el-button size="small" type="danger">删除</el-button>
                  </template>
                </el-popconfirm>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>

        <!-- Tab 4: 申请流程 -->
        <el-tab-pane label="申请流程" name="timelinePhases">
          <div style="margin-bottom: 12px">
            <el-button type="primary" size="small" @click="openSubDialog('timelinePhase')">添加步骤</el-button>
          </div>
          <el-table :data="subData.timelinePhases" border stripe>
            <el-table-column prop="phase_number" label="步骤号" width="80" />
            <el-table-column prop="title" label="标题" min-width="150" />
            <el-table-column prop="description" label="描述" min-width="200" />
            <el-table-column prop="duration" label="周期" width="100" />
            <el-table-column prop="sort_order" label="排序" width="80" />
            <el-table-column label="操作" width="150">
              <template #default="{ row }">
                <el-button size="small" @click="openSubDialog('timelinePhase', row)">编辑</el-button>
                <el-popconfirm title="确定删除？" @confirm="deleteSubItem('timelinePhase', row.id)">
                  <template #reference>
                    <el-button size="small" type="danger">删除</el-button>
                  </template>
                </el-popconfirm>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>
      </el-tabs>

      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="saving" @click="handleSave">保存</el-button>
      </template>
    </el-dialog>

    <!-- Sub-entity edit dialog -->
    <el-dialog
      v-model="subDialogVisible"
      :title="subDialogTitle"
      width="500px"
      destroy-on-close
    >
      <el-form ref="subFormRef" :model="subForm" label-width="100px">
        <template v-if="subType === 'requirement'">
          <el-form-item label="条件描述" prop="label">
            <el-input v-model="subForm.label" />
          </el-form-item>
          <el-form-item label="是否必需" prop="is_required">
            <el-switch v-model="subForm.is_required" />
          </el-form-item>
          <el-form-item label="排序" prop="sort_order">
            <el-input-number v-model="subForm.sort_order" :min="0" />
          </el-form-item>
        </template>
        <template v-else-if="subType === 'costItem'">
          <el-form-item label="费用名称" prop="name">
            <el-input v-model="subForm.name" />
          </el-form-item>
          <el-form-item label="金额" prop="amount">
            <el-input v-model="subForm.amount" />
          </el-form-item>
          <el-form-item label="说明" prop="note">
            <el-input v-model="subForm.note" />
          </el-form-item>
          <el-form-item label="排序" prop="sort_order">
            <el-input-number v-model="subForm.sort_order" :min="0" />
          </el-form-item>
        </template>
        <template v-else-if="subType === 'timelinePhase'">
          <el-form-item label="步骤号" prop="phase_number">
            <el-input-number v-model="subForm.phase_number" :min="1" />
          </el-form-item>
          <el-form-item label="标题" prop="title">
            <el-input v-model="subForm.title" />
          </el-form-item>
          <el-form-item label="描述" prop="description">
            <el-input v-model="subForm.description" type="textarea" :rows="3" />
          </el-form-item>
          <el-form-item label="周期" prop="duration">
            <el-input v-model="subForm.duration" />
          </el-form-item>
          <el-form-item label="排序" prop="sort_order">
            <el-input-number v-model="subForm.sort_order" :min="0" />
          </el-form-item>
        </template>
      </el-form>
      <template #footer>
        <el-button @click="subDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="subSaving" @click="handleSubSave">保存</el-button>
      </template>
    </el-dialog>
```

- [ ] **Step 2: Replace the script section**

Replace the entire `<script setup lang="ts">` block (lines 172-327) with:

```ts
<script setup lang="ts">
import { ElMessage, type FormInstance, type FormRules } from 'element-plus';

definePageMeta({ layout: 'admin', middleware: 'auth' });

interface Project {
  id: string;
  slug: string;
  name: string;
  country: string;
  flag_emoji: string;
  tagline: string;
  investment_amount: string;
  investment_value: string;
  processing_period: string;
  target_crowd: string;
  overview_title: string;
  overview_text: string;
  policy_title: string;
  policy_text: string;
  costs_total: string;
  costs_note: string;
  cta_text: string;
  hero_title: string;
  hero_desc: string;
  hero_gradient: string;
  cover_image: string;
  sort_order: number;
  status: number;
}

interface SubEntity {
  id?: number;
  sort_order: number;
  [key: string]: any;
}

const list = ref<Project[]>([]);
const loading = ref(false);
const saving = ref(false);
const page = ref(1);
const pageSize = ref(10);
const total = ref(0);

const dialogVisible = ref(false);
const editingId = ref<string | null>(null);
const formRef = ref<FormInstance>();
const activeTab = ref('basic');

const defaultForm = (): Partial<Project> => ({
  slug: '',
  name: '',
  country: '',
  flag_emoji: '',
  tagline: '',
  investment_amount: '',
  investment_value: '',
  processing_period: '',
  target_crowd: '',
  overview_title: '',
  overview_text: '',
  policy_title: '',
  policy_text: '',
  costs_total: '',
  costs_note: '',
  cta_text: '',
  hero_title: '',
  hero_desc: '',
  hero_gradient: '',
  cover_image: '',
  sort_order: 0,
  status: 0,
});

const form = reactive<Partial<Project>>(defaultForm());

const rules: FormRules = {
  slug: [{ required: true, message: '请输入标识', trigger: 'blur' }],
  name: [{ required: true, message: '请输入项目名称', trigger: 'blur' }],
  country: [{ required: true, message: '请输入国家', trigger: 'blur' }],
};

type SubType = 'requirement' | 'costItem' | 'timelinePhase';
const subTypeLabels: Record<SubType, string> = {
  requirement: '申请条件',
  costItem: '费用明细',
  timelinePhase: '申请流程',
};

interface SubState {
  requirements: any[];
  costItems: any[];
  timelinePhases: any[];
}

const subData = reactive<SubState>({
  requirements: [],
  costItems: [],
  timelinePhases: [],
});

const subDialogVisible = ref(false);
const subSaving = ref(false);
const subType = ref<SubType>('requirement');
const subEditingId = ref<number | null>(null);
const subFormRef = ref<FormInstance>();
const subForm = reactive<Record<string, any>>({});
const subDialogTitle = computed(() => {
  const prefix = subEditingId.value ? '编辑' : '新增';
  return `${prefix}${subTypeLabels[subType.value]}`;
});

const defaultSubForm = (type: SubType): Record<string, any> => {
  switch (type) {
    case 'requirement':
      return { label: '', is_required: true, sort_order: 0 };
    case 'costItem':
      return { name: '', amount: '', note: '', sort_order: 0 };
    case 'timelinePhase':
      return { phase_number: 1, title: '', description: '', duration: '', sort_order: 0 };
  }
};

const loadList = async () => {
  loading.value = true;
  try {
    const api = useApi();
    const data = await api<{ items: Project[]; total: number }>(
      `/admin/projects?page=${page.value}&per_page=${pageSize.value}`
    );
    list.value = data.items ?? [];
    total.value = data.total ?? 0;
  } catch {
    list.value = [];
    ElMessage.error('加载项目列表失败');
  } finally {
    loading.value = false;
  }
};

const openCreate = () => {
  editingId.value = null;
  Object.assign(form, defaultForm());
  resetSubData();
  activeTab.value = 'basic';
  dialogVisible.value = true;
};

const openEdit = (row: Project) => {
  editingId.value = row.id;
  Object.assign(form, row);
  resetSubData();
  activeTab.value = 'basic';
  dialogVisible.value = true;
};

const onDialogOpened = () => {
  if (editingId.value) {
    loadSubData();
  }
};

const loadSubData = async () => {
  if (!editingId.value) return;
  const api = useApi();
  try {
    const [reqs, costs, phases] = await Promise.all([
      api<any[]>(`/admin/projects/${editingId.value}/requirements`),
      api<any[]>(`/admin/projects/${editingId.value}/cost-items`),
      api<any[]>(`/admin/projects/${editingId.value}/timeline-phases`),
    ]);
    subData.requirements = reqs ?? [];
    subData.costItems = costs ?? [];
    subData.timelinePhases = phases ?? [];
  } catch {
    // sub-data load is best-effort
  }
};

const resetSubData = () => {
  subData.requirements = [];
  subData.costItems = [];
  subData.timelinePhases = [];
};

const handleSave = async () => {
  const valid = await formRef.value?.validate().catch(() => false);
  if (!valid) return;

  saving.value = true;
  try {
    const api = useApi();
    if (editingId.value) {
      await api(`/admin/projects/${editingId.value}`, {
        method: 'PUT',
        body: form,
      });
    } else {
      const created = await api<{ id: string }>('/admin/projects', { method: 'POST', body: form });
      if (created?.id) {
        editingId.value = String(created.id);
      }
    }
    dialogVisible.value = false;
    loadList();
  } catch {
    ElMessage.error(editingId.value ? '更新项目失败' : '创建项目失败');
  } finally {
    saving.value = false;
  }
};

const handleDelete = async (id: string) => {
  try {
    const api = useApi();
    await api(`/admin/projects/${id}`, { method: 'DELETE' });
    loadList();
  } catch {
    ElMessage.error('删除项目失败');
  }
};

// Sub-entity CRUD
const openSubDialog = (type: SubType, row?: any) => {
  subType.value = type;
  subEditingId.value = row?.id ?? null;
  Object.assign(subForm, row ? { ...row } : defaultSubForm(type));
  if (type === 'requirement') {
    subForm.is_required = subForm.is_required === true || subForm.is_required === 1;
  }
  subDialogVisible.value = true;
};

const handleSubSave = async () => {
  if (!editingId.value) {
    ElMessage.warning('请先保存项目基本信息');
    return;
  }
  subSaving.value = true;
  try {
    const api = useApi();
    let endpoint = `/admin/projects/${editingId.value}`;
    if (subType.value === 'requirement') endpoint += '/requirements';
    else if (subType.value === 'costItem') endpoint += '/cost-items';
    else endpoint += '/timeline-phases';

    if (subEditingId.value) {
      endpoint += `/${subEditingId.value}`;
      await api(endpoint, { method: 'PUT', body: subForm });
    } else {
      await api(endpoint, { method: 'POST', body: subForm });
    }
    subDialogVisible.value = false;
    loadSubData();
  } catch {
    ElMessage.error('保存失败');
  } finally {
    subSaving.value = false;
  }
};

const deleteSubItem = async (type: SubType, id: number) => {
  if (!editingId.value) return;
  try {
    const api = useApi();
    let endpoint = `/admin/projects/${editingId.value}`;
    if (type === 'requirement') endpoint += `/requirements/${id}`;
    else if (type === 'costItem') endpoint += `/cost-items/${id}`;
    else endpoint += `/timeline-phases/${id}`;
    await api(endpoint, { method: 'DELETE' });
    loadSubData();
  } catch {
    ElMessage.error('删除失败');
  }
};

const uploadUrl = '/api/v1/admin/media/upload';

const uploadHeaders = computed(() => {
  const token = import.meta.client ? localStorage.getItem('token') : null;
  return token ? { Authorization: `Bearer ${token}` } : {};
});

const handleCoverUploadSuccess = (res: any) => {
  if (res?.data?.url) {
    form.cover_image = res.data.url;
  } else if (res?.url) {
    form.cover_image = res.url;
  }
};

onMounted(() => {
  loadList();
});
</script>
```

- [ ] **Step 3: Run Nuxt typecheck**

```bash
cd frontend && npx nuxi typecheck
```

Expected: no new errors.

---

### Task 8: Frontend — Fix Compare Page

**File:** Modify: `frontend/pages/compare.vue`

- [ ] **Step 1: Remove hardcoded fallback and fix data access**

In `frontend/pages/compare.vue`, remove the `watchEffect` block (lines 121-128) and the `getDefaultComparison` function (lines 130-181). Then update the data fetching to use the new API format.

Replace everything from `const onSelect` (line 110) through `const getColClass` (line 188) with:

```ts
const onSelect = () => {
  if (selectedA.value && selectedB.value) {
    if (selectedA.value === selectedB.value) {
      selectedB.value = '';
      return;
    }
    refreshComparison();
  }
};

const getColClass = (_valueA: string, _valueB: string, col: 'a' | 'b') => {
  if (col === 'a') return 'col-a';
  return 'col-b';
};
```

Then update the `comparison` computed (line 105-107) to properly access data from the API envelope. Change:

```ts
const comparison = computed(() => comparisonRaw.value || null);
```

To:

```ts
const comparison = computed(() => {
  const raw = comparisonRaw.value as any;
  if (raw?.rows) return raw;        // new API format
  if (raw?.data?.rows) return raw.data; // envelope format
  return null;
});
```

- [ ] **Step 2: Also remove the comparisonRaw watchEffect that references getDefaultComparison**

Delete lines 121-128:
```ts
watchEffect(() => {
  if (!comparisonRaw.value && selectedA.value && selectedB.value) {
    const defaults = getDefaultComparison(selectedA.value, selectedB.value);
    if (defaults) {
      comparisonRaw.value = defaults as unknown as typeof comparisonRaw.value;
    }
  }
});
```

- [ ] **Step 3: Remove unused `watchEffect` import if no longer needed**

Check if `watchEffect` is used elsewhere in the file. If it was only used for the removed block, remove it from the import on line 60. The current import is from `vue` — verify it only contains what's still used. (In the original file, `watchEffect` is the only additional import from `vue` beyond the defaults; remove it from the import.)

Change:
```ts
import { watchEffect, computed, ref, onMounted, nextTick } from 'vue';
```

To:
```ts
import { computed, ref, onMounted, nextTick } from 'vue';
```

(If `watchEffect` was the only use of that import style; verify by checking the file.)

- [ ] **Step 4: Run Nuxt typecheck**

```bash
cd frontend && npx nuxi typecheck
```

Expected: no new errors.

---

### Task 9: Final Verification

- [ ] **Step 1: Run Go tests**

```bash
cd backend && go test ./...
```

Expected: all tests pass.

- [ ] **Step 2: Check for any remaining references to deleted pages**

```bash
cd frontend && grep -r "usa/eb5\|hongkong/cies\|panama/property" --include="*.vue" --include="*.ts" .
```

Expected: only the redirects in `nuxt.config.ts` should appear.

- [ ] **Step 3: Verify Go build one final time**

```bash
cd backend && go build ./...
```

Expected: build succeeds.
