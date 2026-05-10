# Project Management Enhancement Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Enhance project management with success cases tab, news tab, compare config tab, conditional frontend rendering, and scrollbar fix.

**Architecture:** Follows existing layered pattern (Model → Repo → Service → Handler → Router). New models CompareConfig/ProjectNews. New repo/service/handler files per existing naming convention. Frontend admin projects.vue gains 3 tabs. Frontend project detail page gains conditional sections.

**Tech Stack:** Go + Gin + GORM (backend), Nuxt 3 + Element Plus (frontend), MySQL 8.0

---

## Task 1: Database Migration

**Files:**
- Create: `database/migrations/000019_add_compare_config_and_project_news.up.sql`
- Create: `database/migrations/000019_add_compare_config_and_project_news.down.sql`

- [ ] **Step 1: Create up migration**

```sql
CREATE TABLE IF NOT EXISTS `compare_configs` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `project_id` BIGINT UNSIGNED NOT NULL,
  `compare_with` JSON NOT NULL,
  `compare_fields` JSON NOT NULL,
  `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_project_id` (`project_id`),
  CONSTRAINT `fk_cc_project` FOREIGN KEY (`project_id`) REFERENCES `projects` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `project_news` (
  `project_id` BIGINT UNSIGNED NOT NULL,
  `page_id` BIGINT UNSIGNED NOT NULL,
  `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`project_id`, `page_id`),
  CONSTRAINT `fk_pn_project` FOREIGN KEY (`project_id`) REFERENCES `projects` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_pn_page` FOREIGN KEY (`page_id`) REFERENCES `pages` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
```

- [ ] **Step 2: Create down migration**

```sql
DROP TABLE IF EXISTS `project_news`;
DROP TABLE IF EXISTS `compare_configs`;
```

- [ ] **Step 3: Commit**

```bash
git add database/migrations/000019_add_compare_config_and_project_news.up.sql database/migrations/000019_add_compare_config_and_project_news.down.sql
git commit -m "feat: add compare_configs and project_news tables"
```

---

## Task 2: New Models + Project Model Update

**Files:**
- Create: `backend/internal/model/compare_config.go`
- Create: `backend/internal/model/project_news.go`
- Modify: `backend/internal/model/project.go`

- [ ] **Step 1: Create CompareConfig model**

```go
package model

import (
	"time"

	"gorm.io/datatypes"
)

type CompareConfig struct {
	ID            uint64         `gorm:"primaryKey;autoIncrement" json:"id"`
	ProjectID     uint64         `gorm:"uniqueIndex;not null" json:"project_id"`
	CompareWith   datatypes.JSON `gorm:"type:json;not null" json:"compare_with"`
	CompareFields datatypes.JSON `gorm:"type:json;not null" json:"compare_fields"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
}

func (CompareConfig) TableName() string { return "compare_configs" }
```

- [ ] **Step 2: Create ProjectNews model**

```go
package model

import "time"

type ProjectNews struct {
	ProjectID uint64    `gorm:"primaryKey" json:"project_id"`
	PageID    uint64    `gorm:"primaryKey" json:"page_id"`
	CreatedAt time.Time `json:"created_at"`
	Page      *Page     `gorm:"foreignKey:PageID" json:"page,omitempty"`
}

func (ProjectNews) TableName() string { return "project_news" }
```

- [ ] **Step 3: Update Project model — add new associations**

In `backend/internal/model/project.go`, add to the Project struct:
```go
Cases         []Case         `gorm:"foreignKey:ProjectID" json:"cases,omitempty"`
News          []Page         `gorm:"many2many:project_news;" json:"news,omitempty"`
CompareConfig *CompareConfig `gorm:"foreignKey:ProjectID" json:"compare_config,omitempty"`
```

Note: `Cases` already exists at line 43. `News` and `CompareConfig` are new.

- [ ] **Step 4: Add `gorm.io/datatypes` dependency**

```bash
cd backend && go get gorm.io/datatypes
```

- [ ] **Step 5: Verify compilation**

```bash
cd backend && go build ./...
```
Expected: builds successfully.

- [ ] **Step 6: Commit**

```bash
git add backend/internal/model/compare_config.go backend/internal/model/project_news.go backend/internal/model/project.go backend/go.mod backend/go.sum
git commit -m "feat: add CompareConfig and ProjectNews models"
```

---

## Task 3: Compare Fields Config

**Files:**
- Create: `backend/internal/config/compare_fields.go`

- [ ] **Step 1: Create compare fields constants**

```go
package config

type CompareField struct {
	Key   string `json:"key"`
	Label string `json:"label"`
	From  string `json:"from"`
}

var CompareFields = []CompareField{
	{Key: "investment_amount", Label: "投资金额", From: "project"},
	{Key: "processing_period", Label: "办理周期", From: "project"},
	{Key: "target_crowd", Label: "适合人群", From: "project"},
	{Key: "country", Label: "国家", From: "project"},
	{Key: "costs_total", Label: "费用总计", From: "project"},
	{Key: "requirements_count", Label: "申请条件数", From: "requirements"},
	{Key: "timeline_steps", Label: "流程步骤数", From: "timeline"},
	{Key: "overview_text", Label: "项目介绍", From: "project"},
	{Key: "tagline", Label: "标语", From: "project"},
	{Key: "policy_title", Label: "政策标题", From: "project"},
}
```

- [ ] **Step 2: Verify compilation**

```bash
cd backend && go build ./...
```
Expected: builds successfully.

- [ ] **Step 3: Commit**

```bash
git add backend/internal/config/compare_fields.go
git commit -m "feat: add compare fields configuration"
```

---

## Task 4: CompareConfig Repository

**Files:**
- Create: `backend/internal/repository/compare_config_repo.go`
- Modify: `backend/internal/repository/interfaces.go`
- Modify: `backend/internal/repository/repository.go`

- [ ] **Step 1: Create CompareConfigRepo**

```go
package repository

import (
	"mygo-immigration/backend/internal/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type CompareConfigRepo struct {
	db *gorm.DB
}

func (r *CompareConfigRepo) FindByProjectID(projectID uint64) (*model.CompareConfig, error) {
	var cfg model.CompareConfig
	err := r.db.Where("project_id = ?", projectID).First(&cfg).Error
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

func (r *CompareConfigRepo) Upsert(cfg *model.CompareConfig) error {
	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "project_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"compare_with", "compare_fields", "updated_at"}),
	}).Create(cfg).Error
}

func (r *CompareConfigRepo) DeleteByProjectID(projectID uint64) error {
	return r.db.Where("project_id = ?", projectID).Delete(&model.CompareConfig{}).Error
}
```

- [ ] **Step 2: Add CompareConfigRepository interface**

In `backend/internal/repository/interfaces.go`, append:
```go
// CompareConfigRepository defines the interface for compare config data access.
type CompareConfigRepository interface {
	FindByProjectID(projectID uint64) (*model.CompareConfig, error)
	Upsert(cfg *model.CompareConfig) error
	DeleteByProjectID(projectID uint64) error
}
```

- [ ] **Step 3: Wire CompareConfigRepo into Repository struct**

In `backend/internal/repository/repository.go`, add to struct:
```go
CompareConfig *CompareConfigRepo
```

In `New()`:
```go
CompareConfig: &CompareConfigRepo{db: db},
```

- [ ] **Step 4: Verify compilation**

```bash
cd backend && go build ./...
```
Expected: builds successfully.

- [ ] **Step 5: Commit**

```bash
git add backend/internal/repository/compare_config_repo.go backend/internal/repository/interfaces.go backend/internal/repository/repository.go
git commit -m "feat: add CompareConfig repository"
```

---

## Task 5: Case Repository — Hard Delete

**Files:**
- Modify: `backend/internal/repository/case_repo.go`

- [ ] **Step 1: Add HardDelete method**

```go
func (r *CaseRepo) HardDelete(id uint64) error {
	return r.db.Unscoped().Delete(&model.Case{}, id).Error
}
```

- [ ] **Step 2: Verify compilation**

```bash
cd backend && go build ./...
```
Expected: builds successfully.

- [ ] **Step 3: Commit**

```bash
git add backend/internal/repository/case_repo.go
git commit -m "feat: add HardDelete to CaseRepo"
```

---

## Task 6: CompareConfig Service

**Files:**
- Create: `backend/internal/service/compare_config_svc.go`
- Modify: `backend/internal/service/service.go`

- [ ] **Step 1: Create CompareConfigService**

```go
package service

import (
	"fmt"

	"mygo-immigration/backend/internal/model"
	"mygo-immigration/backend/internal/repository"
)

type CompareConfigService struct {
	repo repository.CompareConfigRepository
}

func (s *CompareConfigService) GetByProjectID(projectID uint64) (*model.CompareConfig, error) {
	cfg, err := s.repo.FindByProjectID(projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to get compare config: %w", err)
	}
	return cfg, nil
}

func (s *CompareConfigService) Save(cfg *model.CompareConfig) (*model.CompareConfig, error) {
	if err := s.repo.Upsert(cfg); err != nil {
		return nil, fmt.Errorf("failed to save compare config: %w", err)
	}
	return cfg, nil
}

func (s *CompareConfigService) DeleteByProjectID(projectID uint64) error {
	if err := s.repo.DeleteByProjectID(projectID); err != nil {
		return fmt.Errorf("failed to delete compare config: %w", err)
	}
	return nil
}
```

- [ ] **Step 2: Wire into Service struct**

In `backend/internal/service/service.go`, add to struct:
```go
CompareConfig *CompareConfigService
```

In `New()`:
```go
CompareConfig: &CompareConfigService{repo: repo.CompareConfig},
```

- [ ] **Step 3: Verify compilation**

```bash
cd backend && go build ./...
```
Expected: builds successfully.

- [ ] **Step 4: Commit**

```bash
git add backend/internal/service/compare_config_svc.go backend/internal/service/service.go
git commit -m "feat: add CompareConfig service"
```

---

## Task 7: Case Service — Hard Delete

**Files:**
- Modify: `backend/internal/service/case_svc.go`

- [ ] **Step 1: Add HardDelete method**

```go
// HardDelete permanently removes a case study by ID.
func (s *CaseService) HardDelete(id uint64) error {
	if id == 0 {
		return errors.New("case id is required")
	}
	if err := s.repo.HardDelete(id); err != nil {
		return fmt.Errorf("failed to hard delete case: %w", err)
	}
	return nil
}
```

Note: The `CaseService` needs access to `HardDelete`. The current `CaseService` uses `repository.CaseRepository` interface which doesn't have `HardDelete`. Add `HardDelete` to the `CaseRepository` interface in `interfaces.go`:

```go
type CaseRepository interface {
	FindByProjectID(projectID uint64) ([]model.Case, error)
	FindAll(search string) ([]model.Case, error)
	Create(c *model.Case) error
	Update(c *model.Case) error
	Delete(id uint64) error
	HardDelete(id uint64) error
}
```

- [ ] **Step 2: Verify compilation**

```bash
cd backend && go build ./...
```
Expected: builds successfully.

- [ ] **Step 3: Commit**

```bash
git add backend/internal/service/case_svc.go backend/internal/repository/interfaces.go
git commit -m "feat: add HardDelete to CaseService and CaseRepository interface"
```

---

## Task 8: Backend Handlers — Cases, News, CompareConfig

**Files:**
- Modify: `backend/internal/handler/case_handler.go`
- Create: `backend/internal/handler/compare_config_handler.go`
- Modify: `backend/internal/router/router.go`
- Modify: `backend/internal/handler/project_handler.go`

- [ ] **Step 1: Add project-scoped case handlers to case_handler.go**

Append to `backend/internal/handler/case_handler.go`:

```go
// ListProjectCases returns cases belonging to a project.
func (h *Handler) ListProjectCases(c *gin.Context) {
	projectID, err := parseIDParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "invalid project id"))
		return
	}

	cases, err := h.svc.Case.ListByProject(projectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Error(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.Success(cases))
}

// CreateProjectCase creates a case bound to a project.
func (h *Handler) CreateProjectCase(c *gin.Context) {
	projectID, err := parseIDParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "invalid project id"))
		return
	}

	var caseModel model.Case
	if err := c.ShouldBindJSON(&caseModel); err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "invalid request"))
		return
	}
	caseModel.ProjectID = &projectID

	created, err := h.svc.Case.Create(&caseModel)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Error(500, err.Error()))
		return
	}

	c.JSON(http.StatusCreated, dto.Success(created))
}

// DeleteProjectCase hard-deletes a case.
func (h *Handler) DeleteProjectCase(c *gin.Context) {
	caseID, err := parseIDParam(c, "cid")
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "invalid case id"))
		return
	}

	if err := h.svc.Case.HardDelete(caseID); err != nil {
		c.JSON(http.StatusInternalServerError, dto.Error(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.Success(nil))
}
```

Also add `ListByProject` to `CaseService`:

```go
func (s *CaseService) ListByProject(projectID uint64) ([]model.Case, error) {
	cases, err := s.repo.FindByProjectID(projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to list cases by project: %w", err)
	}
	return cases, nil
}
```

- [ ] **Step 2: Create compare_config_handler.go**

```go
package handler

import (
	"net/http"

	"mygo-immigration/backend/internal/config"
	"mygo-immigration/backend/internal/dto"
	"mygo-immigration/backend/internal/model"

	"github.com/gin-gonic/gin"
)

// GetCompareConfig returns the compare config for a project.
func (h *Handler) GetCompareConfig(c *gin.Context) {
	projectID, err := parseIDParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "invalid project id"))
		return
	}

	cfg, err := h.svc.CompareConfig.GetByProjectID(projectID)
	if err != nil {
		c.JSON(http.StatusOK, dto.Success(nil)) // return null if not found
		return
	}
	c.JSON(http.StatusOK, dto.Success(cfg))
}

type saveCompareConfigRequest struct {
	ProjectID     uint64   `json:"project_id"`
	CompareWith   []string `json:"compare_with"`
	CompareFields []string `json:"compare_fields"`
}

// SaveCompareConfig upserts the compare config for a project.
func (h *Handler) SaveCompareConfig(c *gin.Context) {
	projectID, err := parseIDParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "invalid project id"))
		return
	}

	var req saveCompareConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "invalid request"))
		return
	}

	// Validate: need at least 2 projects
	if len(req.CompareWith) < 2 {
		// Delete config if only 1 project
		_ = h.svc.CompareConfig.DeleteByProjectID(projectID)
		c.JSON(http.StatusOK, dto.Success(nil))
		return
	}

	cfg := &model.CompareConfig{
		ProjectID:     projectID,
		CompareWith:   toJSON(req.CompareWith),
		CompareFields: toJSON(req.CompareFields),
	}

	saved, err := h.svc.CompareConfig.Save(cfg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Error(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.Success(saved))
}

// ListCompareFields returns available compare field definitions.
func (h *Handler) ListCompareFields(c *gin.Context) {
	c.JSON(http.StatusOK, dto.Success(config.CompareFields))
}
```

Note: Need a `toJSON` helper for converting `[]string` to `datatypes.JSON`. Add this to the handler package:

In a new file `backend/internal/handler/helpers.go` (or add to existing handler files):
```go
package handler

import (
	"encoding/json"
	"gorm.io/datatypes"
)

func toJSON(v interface{}) datatypes.JSON {
	b, _ := json.Marshal(v)
	return datatypes.JSON(b)
}
```

Actually, use `encoding/json` to marshal and store as `datatypes.JSON`:

```go
import "encoding/json"
import "gorm.io/datatypes"

func toJSON(v interface{}) datatypes.JSON {
	b, _ := json.Marshal(v)
	return datatypes.JSON(b)
}
```

- [ ] **Step 3: Add project news handlers to a new file or project_handler.go**

Add to `backend/internal/handler/project_handler.go`:

```go
type addNewsRequest struct {
	PageIDs []uint64 `json:"page_ids"`
}

// ListProjectNews returns news pages linked to a project.
func (h *Handler) ListProjectNews(c *gin.Context) {
	projectID, err := parseIDParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "invalid project id"))
		return
	}

	news, err := h.svc.Project.ListNews(projectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Error(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.Success(news))
}

// AddProjectNews links news pages to a project.
func (h *Handler) AddProjectNews(c *gin.Context) {
	projectID, err := parseIDParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "invalid project id"))
		return
	}

	var req addNewsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "invalid request"))
		return
	}

	if err := h.svc.Project.AddNews(projectID, req.PageIDs); err != nil {
		c.JSON(http.StatusInternalServerError, dto.Error(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.Success(nil))
}

// RemoveProjectNews unlinks a news page from a project.
func (h *Handler) RemoveProjectNews(c *gin.Context) {
	projectID, err := parseIDParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "invalid project id"))
		return
	}

	pageID, err := parseIDParam(c, "page_id")
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "invalid page id"))
		return
	}

	if err := h.svc.Project.RemoveNews(projectID, pageID); err != nil {
		c.JSON(http.StatusInternalServerError, dto.Error(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.Success(nil))
}
```

Need to add `ProjectNewsRepo` and service methods. Add to `ProjectRepo`:

```go
func (r *ProjectRepo) FindNews(projectID uint64) ([]model.Page, error) {
	var news []model.Page
	err := r.db.Joins("JOIN project_news ON project_news.page_id = pages.id").
		Where("project_news.project_id = ?", projectID).
		Order("project_news.created_at DESC").
		Find(&news).Error
	return news, err
}

func (r *ProjectRepo) AddNews(projectID uint64, pageIDs []uint64) error {
	for _, pageID := range pageIDs {
		err := r.db.Exec("INSERT IGNORE INTO project_news (project_id, page_id) VALUES (?, ?)", projectID, pageID).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *ProjectRepo) RemoveNews(projectID, pageID uint64) error {
	return r.db.Exec("DELETE FROM project_news WHERE project_id = ? AND page_id = ?", projectID, pageID).Error
}
```

Add to `ProjectService`:

```go
func (s *ProjectService) ListNews(projectID uint64) ([]model.Page, error) {
	return s.repo.FindNews(projectID)
}

func (s *ProjectService) AddNews(projectID uint64, pageIDs []uint64) error {
	return s.repo.AddNews(projectID, pageIDs)
}

func (s *ProjectService) RemoveNews(projectID, pageID uint64) error {
	return s.repo.RemoveNews(projectID, pageID)
}
```

Update `ProjectRepository` interface:
```go
FindNews(projectID uint64) ([]model.Page, error)
AddNews(projectID uint64, pageIDs []uint64) error
RemoveNews(projectID, pageID uint64) error
```

- [ ] **Step 4: Update router.go — add new routes**

In the admin projects group, add:
```go
// Cases sub-resources
projects.GET("/:id/cases", middleware.RBAC("admin:read"), h.ListProjectCases)
projects.POST("/:id/cases", middleware.RBAC("projects:write"), h.CreateProjectCase)
projects.PUT("/:id/cases/:cid", middleware.RBAC("projects:write"), h.UpdateCase)
projects.DELETE("/:id/cases/:cid", middleware.RBAC("projects:write"), h.DeleteProjectCase)

// News sub-resources
projects.GET("/:id/news", middleware.RBAC("admin:read"), h.ListProjectNews)
projects.POST("/:id/news", middleware.RBAC("projects:write"), h.AddProjectNews)
projects.DELETE("/:id/news/:page_id", middleware.RBAC("projects:write"), h.RemoveProjectNews)

// Compare config
projects.GET("/:id/compare-config", middleware.RBAC("admin:read"), h.GetCompareConfig)
projects.PUT("/:id/compare-config", middleware.RBAC("projects:write"), h.SaveCompareConfig)
```

And add the compare fields endpoint outside the projects group:
```go
admin.GET("/compare-fields", middleware.RBAC("admin:read"), h.ListCompareFields)
```

- [ ] **Step 5: Update GetProject to Preload News and CompareConfig**

In `backend/internal/repository/project_repo.go`, update `FindBySlug`:
```go
func (r *ProjectRepo) FindBySlug(slug string) (*model.Project, error) {
	var project model.Project
	err := r.db.
		Preload("Requirements").
		Preload("CostItems").
		Preload("TimelinePhases").
		Preload("Milestones").
		Preload("FAQs").
		Preload("Cases").
		Preload("News").
		Preload("CompareConfig").
		Where("slug = ?", slug).
		First(&project).Error
	if err != nil {
		return nil, err
	}
	return &project, nil
}
```

- [ ] **Step 6: Verify compilation and run tests**

```bash
cd backend && go build ./... && go test ./... -v 2>&1 | tail -20
```
Expected: builds, all tests pass.

- [ ] **Step 7: Commit**

```bash
git add backend/
git commit -m "feat: add project cases, news, and compare config handlers and routes"
```

---

## Task 9: Frontend Admin — Scrollbar Fix + 3 New Tabs

**Files:**
- Modify: `frontend/pages/admin/projects.vue`

This is the largest single change. Make these modifications:

### 9a: Scrollbar Fix

- [ ] **Step 1: Add fixed="right" to sub-tab table action columns**

In the requirements, costItems, and timelinePhases tab tables, add `fixed="right"` to the action column (the last `<el-table-column label="操作">`):

```html
<el-table-column label="操作" width="120" fixed="right">
```

### 9b: Success Cases Tab

- [ ] **Step 2: Add case types and state**

Add to script section:
```ts
interface CaseItem {
  id: number;
  name: string;
  country_from: string;
  investment_amount: string;
  processing_period: string;
  description: string;
  photo_url: string;
  sort_order: number;
}
```

Add to SubState and subTypeLabels:
```ts
type SubType = 'requirement' | 'costItem' | 'timelinePhase' | 'caseItem';
const subTypeLabels: Record<SubType, string> = {
  requirement: '申请条件',
  costItem: '费用明细',
  timelinePhase: '申请流程',
  caseItem: '成功案例',
};
```

Add `cases` to subData:
```ts
const subData = reactive<SubState>({
  requirements: [],
  costItems: [],
  timelinePhases: [],
  cases: [],
});
```

In loadSubData, add:
```ts
const cases = await api<CaseItem[]>(`/admin/projects/${editingId.value}/cases`);
subData.cases = cases ?? [];
```

Add case form defaults:
```ts
case 'caseItem':
  return { name: '', country_from: '', investment_amount: '', processing_period: '', description: '', photo_url: '', sort_order: 0 };
```

- [ ] **Step 3: Add case tab pane and dialog template**

Add after the timeline tab-pane:
```html
<el-tab-pane label="成功案例" name="cases">
  <div style="margin-bottom: 12px">
    <el-button type="primary" size="small" @click="openSubDialog('caseItem')">添加案例</el-button>
  </div>
  <el-table :data="subData.cases" border size="small" max-height="360">
    <el-table-column prop="name" label="名称" min-width="120" />
    <el-table-column prop="country_from" label="来源国" width="80" />
    <el-table-column prop="investment_amount" label="投资金额" width="100" />
    <el-table-column prop="processing_period" label="处理周期" width="90" />
    <el-table-column prop="sort_order" label="排序" width="60" />
    <el-table-column label="操作" width="120" fixed="right">
      <template #default="{ row: r }">
        <div class="table-actions">
          <button class="action-btn" @click="openSubDialog('caseItem', r)">编辑</button>
          <el-popconfirm title="确定删除？" @confirm="deleteSubItem('caseItem', r.id)">
            <template #reference>
              <button class="action-btn danger">删除</button>
            </template>
          </el-popconfirm>
        </div>
      </template>
    </el-table-column>
  </el-table>
</el-tab-pane>
```

Add to the dialog template (after the timelinePhase template block):
```html
<template v-else-if="subType === 'caseItem'">
  <el-form-item label="名称" prop="name">
    <el-input v-model="subForm.name" />
  </el-form-item>
  <el-form-item label="来源国" prop="country_from">
    <el-input v-model="subForm.country_from" />
  </el-form-item>
  <el-form-item label="投资金额" prop="investment_amount">
    <el-input v-model="subForm.investment_amount" />
  </el-form-item>
  <el-form-item label="处理周期" prop="processing_period">
    <el-input v-model="subForm.processing_period" />
  </el-form-item>
  <el-form-item label="描述" prop="description">
    <el-input v-model="subForm.description" type="textarea" :rows="3" />
  </el-form-item>
  <el-form-item label="照片URL" prop="photo_url">
    <el-input v-model="subForm.photo_url" />
  </el-form-item>
  <el-form-item label="排序" prop="sort_order">
    <el-input-number v-model="subForm.sort_order" :min="0" />
  </el-form-item>
</template>
```

Update handleSubSave to handle `caseItem`:
```ts
else if (subType.value === 'caseItem') endpoint += '/cases';
```

Update deleteSubItem:
```ts
else if (type === 'caseItem') endpoint += `/cases/${id}`;
```

### 9c: News Tab

- [ ] **Step 4: Add news types and state**

```ts
interface NewsItem {
  id: number;
  title: string;
  slug: string;
  status: string;
  created_at: string;
}

interface NewsLink {
  page_id: number;
  page?: NewsItem;
  created_at: string;
}
```

Add reactive state:
```ts
const subNews = ref<NewsLink[]>([]);
const newsDialogVisible = ref(false);
const newsSearchQuery = ref('');
const newsOptions = ref<NewsItem[]>([]);
const newsSelected = ref<number[]>([]);
```

Add load functions:
```ts
const loadNews = async () => {
  if (!editingId.value) return;
  try {
    const api = useApi();
    const data = await api<NewsLink[]>(`/admin/projects/${editingId.value}/news`);
    subNews.value = data ?? [];
  } catch { subNews.value = []; }
};

const openNewsDialog = async () => {
  try {
    const api = useApi();
    const data = await api<{ items: NewsItem[] }>('/admin/pages?page_type=news&status=published&all=true');
    newsOptions.value = data.items ?? [];
  } catch { newsOptions.value = []; }
  newsSelected.value = [];
  newsDialogVisible.value = true;
};

const addNewsLinks = async () => {
  if (!editingId.value || newsSelected.value.length === 0) return;
  try {
    const api = useApi();
    await api(`/admin/projects/${editingId.value}/news`, {
      method: 'POST',
      body: { page_ids: newsSelected.value },
    });
    newsDialogVisible.value = false;
    loadNews();
  } catch { ElMessage.error('添加资讯失败'); }
};

const removeNewsLink = async (pageId: number) => {
  if (!editingId.value) return;
  try {
    const api = useApi();
    await api(`/admin/projects/${editingId.value}/news/${pageId}`, { method: 'DELETE' });
    loadNews();
  } catch { ElMessage.error('解除关联失败'); }
};
```

Add news tab pane:
```html
<el-tab-pane label="资讯" name="news">
  <div style="margin-bottom: 12px">
    <el-button type="primary" size="small" @click="openNewsDialog">添加资讯</el-button>
  </div>
  <el-table :data="subNews" border size="small" max-height="360">
    <el-table-column label="标题" min-width="180">
      <template #default="{ row }">
        <span>{{ row.page?.title || '(已删除)' }}</span>
      </template>
    </el-table-column>
    <el-table-column label="状态" width="80">
      <template #default="{ row }">
        <span :class="['status-pill', row.page?.status === 'published' ? 'published' : 'draft']">
          {{ row.page?.status === 'published' ? '已发布' : '草稿' }}
        </span>
      </template>
    </el-table-column>
    <el-table-column label="操作" width="80" fixed="right">
      <template #default="{ row }">
        <div class="table-actions">
          <el-popconfirm title="确定解除关联？" @confirm="removeNewsLink(row.page_id)">
            <template #reference>
              <button class="action-btn danger">移除</button>
            </template>
          </el-popconfirm>
        </div>
      </template>
    </el-table-column>
  </el-table>
</el-tab-pane>
```

Add news dialog after the existing sub-entity dialog:
```html
<el-dialog v-model="newsDialogVisible" title="添加资讯" width="500px" destroy-on-close>
  <el-select v-model="newsSelected" multiple filterable placeholder="搜索新闻页面..." style="width: 100%">
    <el-option v-for="n in newsOptions" :key="n.id" :label="n.title" :value="n.id" />
  </el-select>
  <template #footer>
    <el-button @click="newsDialogVisible = false">取消</el-button>
    <el-button type="primary" :disabled="newsSelected.length === 0" @click="addNewsLinks">确认添加</el-button>
  </template>
</el-dialog>
```

Update `onDialogOpened` to also load news:
```ts
const onDialogOpened = () => {
  if (editingId.value) {
    loadSubData();
    loadNews();
  }
};
```

### 9d: Compare Config Tab

- [ ] **Step 5: Add compare config state and logic**

```ts
interface CompareField {
  key: string;
  label: string;
  from: string;
}

interface CompareConfig {
  id?: number;
  project_id: number;
  compare_with: string[];
  compare_fields: string[];
}

interface ProjectOption {
  slug: string;
  name: string;
}
```

Add state:
```ts
const compareFields = ref<CompareField[]>([]);
const compareConfig = reactive<CompareConfig>({ project_id: 0, compare_with: [], compare_fields: [] });
const projectOptions = ref<ProjectOption[]>([]);
```

Add load/save functions:
```ts
const loadCompareConfig = async () => {
  if (!editingId.value) return;
  try {
    const api = useApi();
    const [fields, cfg] = await Promise.all([
      api<CompareField[]>('/admin/compare-fields'),
      api<CompareConfig>(`/admin/projects/${editingId.value}/compare-config`),
    ]);
    compareFields.value = fields ?? [];
    if (cfg) {
      compareConfig.compare_with = cfg.compare_with ?? [];
      compareConfig.compare_fields = cfg.compare_fields ?? [];
    }
  } catch {}
};

const loadProjectOptions = async () => {
  try {
    const api = useApi();
    const data = await api<{ items: ProjectOption[] }>('/admin/projects?all=true');
    projectOptions.value = data.items ?? [];
  } catch { projectOptions.value = []; }
};

const saveCompareConfig = async () => {
  if (!editingId.value) return;
  try {
    const api = useApi();
    await api(`/admin/projects/${editingId.value}/compare-config`, {
      method: 'PUT',
      body: {
        project_id: Number(editingId.value),
        compare_with: compareConfig.compare_with,
        compare_fields: compareConfig.compare_fields,
      },
    });
    ElMessage.success('已保存');
  } catch { ElMessage.error('保存失败'); }
};
```

Add compare config tab pane:
```html
<el-tab-pane label="项目对比" name="compare">
  <el-form label-position="top">
    <el-form-item label="对比项目">
      <el-select
        v-model="compareConfig.compare_with"
        multiple
        filterable
        placeholder="选择对比项目（至少选 2 个）"
        style="width: 100%"
      >
        <el-option
          v-for="p in projectOptions"
          :key="p.slug"
          :label="p.name"
          :value="p.slug"
        />
      </el-select>
      <div style="font-size: 12px; color: var(--el-text-color-secondary); margin-top: 4px">
        当前项目默认参与对比，请至少追加 1 个其他项目
      </div>
    </el-form-item>
    <el-form-item label="对比属性">
      <el-checkbox-group v-model="compareConfig.compare_fields">
        <el-checkbox v-for="f in compareFields" :key="f.key" :value="f.key" :label="f.key">
          {{ f.label }}
        </el-checkbox>
      </el-checkbox-group>
    </el-form-item>
    <el-button type="primary" @click="saveCompareConfig">保存对比配置</el-button>
  </el-form>
</el-tab-pane>
```

Update `onDialogOpened` to also load compare config:
```ts
const onDialogOpened = () => {
  if (editingId.value) {
    loadSubData();
    loadNews();
    loadCompareConfig();
    loadProjectOptions();
  }
};
```

- [ ] **Step 6: Verify typecheck**

```bash
cd frontend && npx nuxi typecheck
```
Expected: no new type errors.

- [ ] **Step 7: Commit**

```bash
git add frontend/pages/admin/projects.vue
git commit -m "feat: add cases, news, compare tabs and fix scrollbar in project admin"
```

---

## Task 10: Frontend Project Detail — Conditional Sections

**Files:**
- Modify: `frontend/pages/projects/[slug].vue`

- [ ] **Step 1: Update ApiProject interface**

Replace the ApiProject interface:
```ts
interface ApiRequirement { label: string; is_required: boolean; }
interface ApiCostItem { name: string; amount: string; note: string; }
interface ApiTimelinePhase { phase_number: number; title: string; description: string; duration: string; }
interface ApiFAQ { question: string; answer: string; }
interface ApiCase { name: string; country_from: string; investment_amount: string; processing_period: string; description: string; photo_url: string; }
interface ApiNewsPage { id: number; title: string; slug: string; cover_image: string; created_at: string; }
interface ApiCompareConfig { compare_with: string[]; compare_fields: string[]; }

interface ApiProject {
  name: string;
  tagline: string;
  country: string;
  cover_image: string;
  investment_amount: string;
  processing_period: string;
  target_crowd: string;
  overview_title: string;
  overview_text: string;
  cta_text: string;
  hero_title: string;
  hero_desc: string;
  hero_gradient: string;
  requirements: ApiRequirement[];
  cost_items: ApiCostItem[];
  timeline_phases: ApiTimelinePhase[];
  faqs: ApiFAQ[];
  cases: ApiCase[];
  news: ApiNewsPage[];
  compare_config: ApiCompareConfig | null;
}
```

- [ ] **Step 2: Update project computed to include new data**

Add to the project computed:
```ts
cases: (p?.cases || []).map((c) => ({
  name: c.name,
  country: c.country_from,
  amount: c.investment_amount,
  period: c.processing_period,
  description: c.description,
  photo: c.photo_url,
})),
news: (p?.news || []).map((n) => ({
  id: n.id,
  title: n.title,
  slug: n.slug,
  cover: n.cover_image,
  date: n.created_at,
})),
compare_config: p?.compare_config || null,
```

- [ ] **Step 3: Add success cases section**

After the FAQ section:
```html
<section v-if="project.cases.length > 0" class="detail-section">
  <h2 class="detail-section-title">成功案例</h2>
  <div class="case-grid">
    <div v-for="c in project.cases" :key="c.name" class="case-card">
      <img v-if="c.photo" :src="c.photo" :alt="c.name" class="case-photo" />
      <div class="case-body">
        <h4 class="case-name">{{ c.name }}</h4>
        <p class="case-meta">{{ c.country }} | {{ c.amount }} | {{ c.period }}</p>
        <p v-if="c.description" class="case-desc">{{ c.description }}</p>
      </div>
    </div>
  </div>
</section>
```

- [ ] **Step 4: Add news section**

```html
<section v-if="project.news.length > 0" class="detail-section">
  <h2 class="detail-section-title">最新资讯</h2>
  <div class="news-list">
    <NuxtLink
      v-for="n in project.news"
      :key="n.id"
      :to="`/pages/${n.slug}`"
      class="news-item"
    >
      <img v-if="n.cover" :src="n.cover" :alt="n.title" class="news-cover" />
      <div class="news-body">
        <h4 class="news-title">{{ n.title }}</h4>
        <span v-if="n.date" class="news-date">{{ new Date(n.date).toLocaleDateString('zh-CN') }}</span>
      </div>
    </NuxtLink>
  </div>
</section>
```

- [ ] **Step 5: Add compare section**

```html
<section v-if="project.compare_config && project.compare_config.compare_with.length >= 2" class="detail-section">
  <h2 class="detail-section-title">项目对比</h2>
  <div class="compare-link">
    <NuxtLink :to="`/compare?slugs=${project.compare_config.compare_with.join(',')}`" class="btn-primary">
      查看对比详情
    </NuxtLink>
  </div>
</section>
```

- [ ] **Step 6: Add styles for new sections**

Append to `<style scoped>`:
```css
.case-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
  gap: 24px;
}

.case-card {
  background: var(--bg-white);
  border: 1px solid var(--border-color);
  border-radius: var(--radius-md);
  overflow: hidden;
}

.case-photo {
  width: 100%;
  height: 200px;
  object-fit: cover;
}

.case-body {
  padding: 16px;
}

.case-name {
  font-size: 18px;
  font-weight: 600;
  margin-bottom: 6px;
}

.case-meta {
  font-size: 13px;
  color: var(--text-light);
  margin-bottom: 8px;
}

.case-desc {
  font-size: 14px;
  color: var(--text-secondary);
  line-height: 1.6;
}

.news-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.news-item {
  display: flex;
  gap: 16px;
  padding: 12px;
  border: 1px solid var(--border-color);
  border-radius: var(--radius-md);
  text-decoration: none;
  color: inherit;
  transition: box-shadow 0.2s;
}

.news-item:hover {
  box-shadow: var(--shadow-sm);
}

.news-cover {
  width: 120px;
  height: 80px;
  object-fit: cover;
  border-radius: 4px;
  flex-shrink: 0;
}

.news-body {
  display: flex;
  flex-direction: column;
  justify-content: center;
}

.news-title {
  font-size: 16px;
  font-weight: 600;
  margin-bottom: 6px;
}

.news-date {
  font-size: 13px;
  color: var(--text-light);
}

.compare-link {
  text-align: center;
}

@media (max-width: 767px) {
  .news-cover {
    width: 80px;
    height: 60px;
  }
}
```

- [ ] **Step 7: Verify typecheck**

```bash
cd frontend && npx nuxi typecheck
```
Expected: no new type errors.

- [ ] **Step 8: Commit**

```bash
git add frontend/pages/projects/[slug].vue
git commit -m "feat: add conditional cases, news, compare sections to project detail"
```

---

## Task 11: End-to-End Verification

**Files:** None (verification only)

- [ ] **Step 1: Run all backend tests**

```bash
cd backend && go test ./... -v 2>&1 | tail -30
```
Expected: all tests pass.

- [ ] **Step 2: Verify backend build**

```bash
cd backend && go build ./...
```
Expected: builds successfully.

- [ ] **Step 3: Verify frontend typecheck**

```bash
cd frontend && npx nuxi typecheck
```
Expected: no new type errors.

- [ ] **Step 4: Check git log**

```bash
git log --oneline -11
```
Expected: all 10+implementation commits in sequence.
