# FAQ 项目归属 & 筛选 实现计划

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** FAQ 管理列表显示所属项目名，管理端搜索+项目筛选可用，公开 FAQ 页动态筛选按钮替代 mock 数据。

**Architecture:** 后端新增 `FAQResponse` DTO，Repository 新增 `FAQWithProject` 结构体 + `FindAll` 返回 JOIN 结果，Service 映射为 DTO，Handler 读取新 query 参数。前端移除 mock，动态生成筛选按钮，管理端增加项目下拉筛选器。

**Tech Stack:** Go + Gin + GORM / Nuxt 3 + Element Plus + TypeScript

**公开 FAQ 端点不提供搜索参数**（FAQ 量小，前端内存过滤即可）。`Search` 方法保留在接口中（为 search_svc 所用）。

---

### Task 1: 创建 FAQResponse DTO 与 FAQWithProject 结构体

**Files:**
- Create: `backend/internal/dto/faq_response.go`
- Modify: `backend/internal/repository/faq_repo.go` (add FAQWithProject at top)

- [ ] **Step 1: 创建 DTO 文件**

```go
// backend/internal/dto/faq_response.go
package dto

// FAQResponse is the API response struct for FAQ entries, including project info.
type FAQResponse struct {
	ID          uint64  `json:"id"`
	Question    string  `json:"question"`
	Answer      string  `json:"answer"`
	ProjectID   *uint64 `json:"project_id"`
	ProjectName string  `json:"project_name"`
	ProjectSlug string  `json:"project_slug"`
	IsGlobal    bool    `json:"is_global"`
	SortOrder   int     `json:"sort_order"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}
```

- [ ] **Step 2: 添加 FAQWithProject 到 faq_repo.go 顶部（in `repository` package）**

在 `package repository` 和 imports 之后、`type FAQRepo struct` 之前插入：

```go
// FAQWithProject holds an FAQ row with joined project columns.
type FAQWithProject struct {
	model.FAQ
	ProjectName string `gorm:"column:project_name"`
	ProjectSlug string `gorm:"column:project_slug"`
}
```

- [ ] **Step 3: 验证编译**

```bash
cd backend && go build ./...
```
Expected: builds.

---

### Task 2: 更新 FAQRepository 接口

**Files:**
- Modify: `backend/internal/repository/interfaces.go:27-36`

- [ ] **Step 1: 替换 FAQRepository 接口定义**

在 `interfaces.go` 中，将第 27-36 行的 `FAQRepository interface` 块替换为：

```go
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
```

注意：`FAQQueryParams` 放在 `FAQRepository` 之前，在 interface 块外面的文件级作用域。

- [ ] **Step 2: 验证编译（预期失败 — faq_repo.go 未适配新签名）**

```bash
cd backend && go build ./... 2>&1
```
Expected: `faq_repo.go` compilation error. Proceed to Task 3.

---

### Task 3: 重写 FAQRepo.FindAll（使用 LEFT JOIN）

**Files:**
- Modify: `backend/internal/repository/faq_repo.go`

- [ ] **Step 1: 替换 faq_repo.go 全部内容**

```go
package repository

import (
	"mygo-immigration/backend/internal/model"

	"gorm.io/gorm"
)

// FAQWithProject holds an FAQ row with joined project columns.
type FAQWithProject struct {
	model.FAQ
	ProjectName string `gorm:"column:project_name"`
	ProjectSlug string `gorm:"column:project_slug"`
}

type FAQRepo struct {
	db *gorm.DB
}

func (r *FAQRepo) FindAll(params FAQQueryParams) ([]FAQWithProject, int64, error) {
	var results []FAQWithProject
	var total int64

	q := r.db.Model(&model.FAQ{}).
		Select("faqs.*, projects.name AS project_name, projects.slug AS project_slug").
		Joins("LEFT JOIN projects ON projects.id = faqs.project_id AND projects.deleted_at IS NULL")

	if params.ProjectID != nil {
		q = q.Where("faqs.project_id = ?", *params.ProjectID)
	}
	if params.IsGlobal != nil {
		q = q.Where("faqs.is_global = ?", *params.IsGlobal)
	}
	if params.Search != "" {
		like := "%" + params.Search + "%"
		q = q.Where("faqs.question LIKE ? OR faqs.answer LIKE ?", like, like)
	}

	// Count total matching rows (without LIMIT/OFFSET).
	if err := q.Session(&gorm.Session{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if params.Page < 1 {
		params.Page = 1
	}
	if params.PerPage < 1 || params.PerPage > 100 {
		params.PerPage = 10
	}
	offset := (params.Page - 1) * params.PerPage

	err := q.
		Order("faqs.sort_order asc").
		Offset(offset).
		Limit(params.PerPage).
		Find(&results).Error
	if err != nil {
		return nil, 0, err
	}
	return results, total, nil
}

func (r *FAQRepo) Create(faq *model.FAQ) error {
	return r.db.Create(faq).Error
}

func (r *FAQRepo) Update(faq *model.FAQ) error {
	return r.db.Omit("created_at").Save(faq).Error
}

func (r *FAQRepo) Delete(id uint64) error {
	return r.db.Delete(&model.FAQ{}, id).Error
}

func (r *FAQRepo) Search(keyword string) ([]model.FAQ, error) {
	var faqs []model.FAQ
	err := r.db.
		Where("question LIKE ? OR answer LIKE ?", "%"+keyword+"%", "%"+keyword+"%").
		Order("sort_order asc").
		Find(&faqs).Error
	if err != nil {
		return nil, err
	}
	return faqs, nil
}
```

- [ ] **Step 2: 验证编译**

```bash
cd backend && go build ./...
```
Expected: builds.

---

### Task 4: 更新 FAQService（新签名 + 映射 DTO）

**Files:**
- Modify: `backend/internal/service/faq_svc.go`

- [ ] **Step 1: 替换 faq_svc.go 全部内容**

```go
package service

import (
	"errors"
	"fmt"
	"time"

	"mygo-immigration/backend/internal/dto"
	"mygo-immigration/backend/internal/model"
	"mygo-immigration/backend/internal/repository"
)

// FAQService handles business logic for FAQ entries.
type FAQService struct {
	repo repository.FAQRepository
}

// NewFAQService creates a new FAQService.
func NewFAQService(repo repository.FAQRepository) *FAQService {
	return &FAQService{repo: repo}
}

// List returns FAQs, optionally filtered by project or global flag.
func (s *FAQService) List(projectID *uint64, isGlobal *bool) ([]dto.FAQResponse, error) {
	results, _, err := s.repo.FindAll(repository.FAQQueryParams{
		ProjectID: projectID,
		IsGlobal:  isGlobal,
		Page:      1,
		PerPage:   1000,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list faqs: %w", err)
	}
	return toFAQResponses(results), nil
}

// AdminList returns paginated FAQs with optional project filter and search.
func (s *FAQService) AdminList(projectID *uint64, search string, page, perPage int) ([]dto.FAQResponse, int64, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 10
	}

	results, total, err := s.repo.FindAll(repository.FAQQueryParams{
		ProjectID: projectID,
		Search:    search,
		Page:      page,
		PerPage:   perPage,
	})
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list faqs: %w", err)
	}
	return toFAQResponses(results), total, nil
}

// toFAQResponses converts FAQWithProject rows to DTO responses.
func toFAQResponses(rows []repository.FAQWithProject) []dto.FAQResponse {
	result := make([]dto.FAQResponse, len(rows))
	for i, r := range rows {
		result[i] = dto.FAQResponse{
			ID:          r.ID,
			Question:    r.Question,
			Answer:      r.Answer,
			ProjectID:   r.ProjectID,
			ProjectName: r.ProjectName,
			ProjectSlug: r.ProjectSlug,
			IsGlobal:    r.IsGlobal,
			SortOrder:   r.SortOrder,
			CreatedAt:   r.CreatedAt.Format(time.RFC3339),
			UpdatedAt:   r.UpdatedAt.Format(time.RFC3339),
		}
	}
	return result
}

// Create creates a new FAQ entry.
func (s *FAQService) Create(faq *model.FAQ) (*model.FAQ, error) {
	if faq == nil {
		return nil, errors.New("faq is nil")
	}
	if faq.Question == "" {
		return nil, errors.New("faq question is required")
	}
	if faq.Answer == "" {
		return nil, errors.New("faq answer is required")
	}
	if err := s.repo.Create(faq); err != nil {
		return nil, fmt.Errorf("failed to create faq: %w", err)
	}
	return faq, nil
}

// Update updates an existing FAQ entry.
func (s *FAQService) Update(id uint64, faq *model.FAQ) (*model.FAQ, error) {
	if faq == nil {
		return nil, errors.New("faq is nil")
	}
	if id == 0 {
		return nil, errors.New("faq id is required")
	}
	faq.ID = id
	if err := s.repo.Update(faq); err != nil {
		return nil, fmt.Errorf("failed to update faq: %w", err)
	}
	return faq, nil
}

// Delete removes an FAQ entry by ID.
func (s *FAQService) Delete(id uint64) error {
	if id == 0 {
		return errors.New("faq id is required")
	}
	if err := s.repo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete faq: %w", err)
	}
	return nil
}
```

- [ ] **Step 2: 验证编译**

```bash
cd backend && go build ./...
```
Expected: builds.

---

### Task 5: 更新 FAQHandler（读取新 query 参数，返回 FAQResponse）

**Files:**
- Modify: `backend/internal/handler/faq_handler.go`

- [ ] **Step 1: 替换 faq_handler.go 全部内容**

```go
package handler

import (
	"net/http"
	"strconv"

	"mygo-immigration/backend/internal/dto"
	"mygo-immigration/backend/internal/model"

	"github.com/gin-gonic/gin"
)

func (h *Handler) ListFAQs(c *gin.Context) {
	var projectID *uint64
	if v := c.Query("project_id"); v != "" {
		id, err := strconv.ParseUint(v, 10, 64)
		if err == nil {
			projectID = &id
		}
	}

	var isGlobal *bool
	if v := c.Query("is_global"); v != "" {
		b, err := strconv.ParseBool(v)
		if err == nil {
			isGlobal = &b
		}
	}

	faqs, err := h.svc.FAQ.List(projectID, isGlobal)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Error(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.Success(faqs))
}

func (h *Handler) AdminListFAQs(c *gin.Context) {
	page, perPage := parsePagination(c)

	var projectID *uint64
	if v := c.Query("project_id"); v != "" {
		id, err := strconv.ParseUint(v, 10, 64)
		if err == nil {
			projectID = &id
		}
	}

	search := c.Query("search")

	faqs, total, err := h.svc.FAQ.AdminList(projectID, search, page, perPage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Error(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessPaginated(faqs, page, perPage, total))
}

func (h *Handler) CreateFAQ(c *gin.Context) {
	var faq model.FAQ
	if err := c.ShouldBindJSON(&faq); err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "invalid request"))
		return
	}

	created, err := h.svc.FAQ.Create(&faq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Error(500, err.Error()))
		return
	}

	c.JSON(http.StatusCreated, dto.Success(created))
}

func (h *Handler) UpdateFAQ(c *gin.Context) {
	id, err := parseIDParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "invalid faq id"))
		return
	}

	var faq model.FAQ
	if err := c.ShouldBindJSON(&faq); err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "invalid request"))
		return
	}

	updated, err := h.svc.FAQ.Update(id, &faq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Error(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.Success(updated))
}

func (h *Handler) DeleteFAQ(c *gin.Context) {
	id, err := parseIDParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "invalid faq id"))
		return
	}

	if err := h.svc.FAQ.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, dto.Error(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.Success(nil))
}
```

- [ ] **Step 2: 验证编译**

```bash
cd backend && go build ./...
```
Expected: builds.

---

### Task 6: 更新测试文件中的 mock（匹配新接口）

**Files:**
- Modify: `backend/internal/service/faq_svc_test.go`
- Modify: `backend/internal/handler/faq_handler_test.go`
- Modify: `backend/internal/service/search_svc_test.go`

- [ ] **Step 1: 更新 service/faq_svc_test.go — 替换 mock 和测试**

替换文件全部内容：

```go
package service

import (
	"errors"
	"testing"

	"mygo-immigration/backend/internal/dto"
	"mygo-immigration/backend/internal/model"
	"mygo-immigration/backend/internal/repository"
)

// mockFAQRepo implements repository.FAQRepository for testing.
type mockFAQRepo struct {
	findAllFn func(params repository.FAQQueryParams) ([]repository.FAQWithProject, int64, error)
	createFn  func(faq *model.FAQ) error
	updateFn  func(faq *model.FAQ) error
	deleteFn  func(id uint64) error
	searchFn  func(keyword string) ([]model.FAQ, error)
}

func (m *mockFAQRepo) FindAll(params repository.FAQQueryParams) ([]repository.FAQWithProject, int64, error) {
	if m.findAllFn != nil {
		return m.findAllFn(params)
	}
	return nil, 0, nil
}
func (m *mockFAQRepo) Create(faq *model.FAQ) error {
	if m.createFn != nil {
		return m.createFn(faq)
	}
	return nil
}
func (m *mockFAQRepo) Update(faq *model.FAQ) error {
	if m.updateFn != nil {
		return m.updateFn(faq)
	}
	return nil
}
func (m *mockFAQRepo) Delete(id uint64) error {
	if m.deleteFn != nil {
		return m.deleteFn(id)
	}
	return nil
}
func (m *mockFAQRepo) Search(keyword string) ([]model.FAQ, error) {
	if m.searchFn != nil {
		return m.searchFn(keyword)
	}
	return nil, nil
}

func makeFAQRow(id uint64, q, a string) repository.FAQWithProject {
	return repository.FAQWithProject{
		FAQ: model.FAQ{ID: id, Question: q, Answer: a},
	}
}

func TestFAQ_List(t *testing.T) {
	sampleRows := []repository.FAQWithProject{
		makeFAQRow(1, "Q1", "A1"),
		makeFAQRow(2, "Q2", "A2"),
	}

	repo := &mockFAQRepo{
		findAllFn: func(params repository.FAQQueryParams) ([]repository.FAQWithProject, int64, error) {
			return sampleRows, 2, nil
		},
	}

	svc := NewFAQService(repo)

	faqs, err := svc.List(nil, nil)
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if len(faqs) != 2 {
		t.Errorf("expected 2 faqs, got %d", len(faqs))
	}
	// Verify DTO fields
	if faqs[0].Question != "Q1" {
		t.Errorf("expected Q1, got %s", faqs[0].Question)
	}
}

func TestFAQ_List_WithProjectFilter(t *testing.T) {
	pid := uint64(5)
	repo := &mockFAQRepo{
		findAllFn: func(params repository.FAQQueryParams) ([]repository.FAQWithProject, int64, error) {
			if params.ProjectID == nil || *params.ProjectID != 5 {
				t.Errorf("expected ProjectID=5, got %v", params.ProjectID)
			}
			return []repository.FAQWithProject{}, 0, nil
		},
	}
	svc := NewFAQService(repo)
	_, err := svc.List(&pid, nil)
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
}

func TestFAQ_List_Error(t *testing.T) {
	repo := &mockFAQRepo{
		findAllFn: func(params repository.FAQQueryParams) ([]repository.FAQWithProject, int64, error) {
			return nil, 0, errors.New("db error")
		},
	}

	svc := NewFAQService(repo)

	_, err := svc.List(nil, nil)
	if err == nil {
		t.Fatal("expected error from repo")
	}
}

func TestFAQ_List_Empty(t *testing.T) {
	repo := &mockFAQRepo{
		findAllFn: func(params repository.FAQQueryParams) ([]repository.FAQWithProject, int64, error) {
			return []repository.FAQWithProject{}, 0, nil
		},
	}

	svc := NewFAQService(repo)

	faqs, err := svc.List(nil, nil)
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if len(faqs) != 0 {
		t.Errorf("expected 0 faqs, got %d", len(faqs))
	}
}

func TestFAQ_Create_Success(t *testing.T) {
	created := false
	repo := &mockFAQRepo{
		createFn: func(faq *model.FAQ) error {
			created = true
			faq.ID = 100
			return nil
		},
	}

	svc := NewFAQService(repo)

	faq, err := svc.Create(&model.FAQ{Question: "Test Q?", Answer: "Test A."})
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if !created {
		t.Error("expected Create to be called on repo")
	}
	if faq.ID != 100 {
		t.Errorf("expected ID 100, got %d", faq.ID)
	}
}

func TestFAQ_Create_NilFAQ(t *testing.T) {
	repo := &mockFAQRepo{}
	svc := NewFAQService(repo)

	_, err := svc.Create(nil)
	if err == nil {
		t.Fatal("expected error for nil faq")
	}
}

func TestFAQ_Create_MissingQuestion(t *testing.T) {
	repo := &mockFAQRepo{}
	svc := NewFAQService(repo)

	_, err := svc.Create(&model.FAQ{Answer: "Some answer"})
	if err == nil {
		t.Fatal("expected error for missing question")
	}
}

func TestFAQ_Create_MissingAnswer(t *testing.T) {
	repo := &mockFAQRepo{}
	svc := NewFAQService(repo)

	_, err := svc.Create(&model.FAQ{Question: "Some question?"})
	if err == nil {
		t.Fatal("expected error for missing answer")
	}
}

func TestFAQ_Create_SpecialCharacters(t *testing.T) {
	var savedFAQ *model.FAQ
	repo := &mockFAQRepo{
		createFn: func(faq *model.FAQ) error {
			savedFAQ = faq
			faq.ID = 1
			return nil
		},
	}

	svc := NewFAQService(repo)

	question := "What about <script>alert('xss')</script> & special chars?"
	answer := "Answer with <b>bold</b> and &amp; encoding"
	faq, err := svc.Create(&model.FAQ{Question: question, Answer: answer})
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if savedFAQ == nil {
		t.Fatal("expected Create to be called on repo")
	}
	if savedFAQ.Question != question {
		t.Errorf("expected question to be saved as-is")
	}
	if faq.ID != 1 {
		t.Errorf("expected ID 1, got %d", faq.ID)
	}
}

func TestFAQ_Create_RepoError(t *testing.T) {
	repo := &mockFAQRepo{
		createFn: func(faq *model.FAQ) error {
			return errors.New("db error")
		},
	}

	svc := NewFAQService(repo)

	_, err := svc.Create(&model.FAQ{Question: "Q", Answer: "A"})
	if err == nil {
		t.Fatal("expected error from repo")
	}
}

func TestFAQ_Update_Success(t *testing.T) {
	updated := false
	repo := &mockFAQRepo{
		updateFn: func(faq *model.FAQ) error {
			updated = true
			return nil
		},
	}

	svc := NewFAQService(repo)

	faq, err := svc.Update(1, &model.FAQ{Question: "Updated Q", Answer: "Updated A"})
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if !updated {
		t.Error("expected Update to be called on repo")
	}
	if faq.ID != 1 {
		t.Errorf("expected ID 1, got %d", faq.ID)
	}
}

func TestFAQ_Update_NilFAQ(t *testing.T) {
	repo := &mockFAQRepo{}
	svc := NewFAQService(repo)

	_, err := svc.Update(1, nil)
	if err == nil {
		t.Fatal("expected error for nil faq in update")
	}
}

func TestFAQ_Update_ZeroID(t *testing.T) {
	repo := &mockFAQRepo{}
	svc := NewFAQService(repo)

	_, err := svc.Update(0, &model.FAQ{Question: "Q", Answer: "A"})
	if err == nil {
		t.Fatal("expected error for zero id")
	}
}

func TestFAQ_Update_RepoError(t *testing.T) {
	repo := &mockFAQRepo{
		updateFn: func(faq *model.FAQ) error {
			return errors.New("db error")
		},
	}

	svc := NewFAQService(repo)

	_, err := svc.Update(1, &model.FAQ{Question: "Q", Answer: "A"})
	if err == nil {
		t.Fatal("expected error from repo")
	}
}

func TestFAQ_Delete_Success(t *testing.T) {
	deleted := false
	repo := &mockFAQRepo{
		deleteFn: func(id uint64) error {
			deleted = true
			return nil
		},
	}

	svc := NewFAQService(repo)

	err := svc.Delete(1)
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if !deleted {
		t.Error("expected Delete to be called on repo")
	}
}

func TestFAQ_Delete_ZeroID(t *testing.T) {
	repo := &mockFAQRepo{}
	svc := NewFAQService(repo)

	err := svc.Delete(0)
	if err == nil {
		t.Fatal("expected error for zero id")
	}
}

func TestFAQ_Delete_RepoError(t *testing.T) {
	repo := &mockFAQRepo{
		deleteFn: func(id uint64) error {
			return errors.New("db error")
		},
	}

	svc := NewFAQService(repo)

	err := svc.Delete(1)
	if err == nil {
		t.Fatal("expected error from repo")
	}
}

func TestFAQ_AdminList_Success(t *testing.T) {
	sampleRows := []repository.FAQWithProject{
		makeFAQRow(1, "Q1", "A1"),
		makeFAQRow(2, "Q2", "A2"),
		makeFAQRow(3, "Q3", "A3"),
		makeFAQRow(4, "Q4", "A4"),
	}

	repo := &mockFAQRepo{
		findAllFn: func(params repository.FAQQueryParams) ([]repository.FAQWithProject, int64, error) {
			// Slice to return based on pagination
			start := (params.Page - 1) * params.PerPage
			end := start + params.PerPage
			if end > len(sampleRows) {
				end = len(sampleRows)
			}
			if start >= len(sampleRows) {
				return []repository.FAQWithProject{}, int64(len(sampleRows)), nil
			}
			return sampleRows[start:end], int64(len(sampleRows)), nil
		},
	}

	svc := NewFAQService(repo)

	faqs, total, err := svc.AdminList(nil, "", 1, 2)
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if total != 4 {
		t.Errorf("expected total 4, got %d", total)
	}
	if len(faqs) != 2 {
		t.Errorf("expected 2 faqs on page 1 with perPage=2, got %d", len(faqs))
	}
}

func TestFAQ_AdminList_WithSearch(t *testing.T) {
	repo := &mockFAQRepo{
		findAllFn: func(params repository.FAQQueryParams) ([]repository.FAQWithProject, int64, error) {
			if params.Search != "apply" {
				t.Errorf("expected search='apply', got '%s'", params.Search)
			}
			return []repository.FAQWithProject{}, 0, nil
		},
	}

	svc := NewFAQService(repo)
	_, _, err := svc.AdminList(nil, "apply", 1, 10)
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
}

func TestFAQ_AdminList_WithProjectFilter(t *testing.T) {
	pid := uint64(3)
	repo := &mockFAQRepo{
		findAllFn: func(params repository.FAQQueryParams) ([]repository.FAQWithProject, int64, error) {
			if params.ProjectID == nil || *params.ProjectID != 3 {
				t.Errorf("expected ProjectID=3, got %v", params.ProjectID)
			}
			return []repository.FAQWithProject{}, 0, nil
		},
	}

	svc := NewFAQService(repo)
	_, _, err := svc.AdminList(&pid, "", 1, 10)
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
}

func TestFAQ_AdminList_DefaultPagination(t *testing.T) {
	repo := &mockFAQRepo{
		findAllFn: func(params repository.FAQQueryParams) ([]repository.FAQWithProject, int64, error) {
			if params.Page != 1 {
				t.Errorf("expected default page 1, got %d", params.Page)
			}
			if params.PerPage != 10 {
				t.Errorf("expected default perPage 10, got %d", params.PerPage)
			}
			return []repository.FAQWithProject{}, 0, nil
		},
	}

	svc := NewFAQService(repo)
	_, _, err := svc.AdminList(nil, "", 0, 0)
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
}

func TestFAQ_AdminList_RepoError(t *testing.T) {
	repo := &mockFAQRepo{
		findAllFn: func(params repository.FAQQueryParams) ([]repository.FAQWithProject, int64, error) {
			return nil, 0, errors.New("db error")
		},
	}

	svc := NewFAQService(repo)

	_, _, err := svc.AdminList(nil, "", 1, 10)
	if err == nil {
		t.Fatal("expected error from repo")
	}
}

func TestFAQ_toFAQResponses_WithProject(t *testing.T) {
	pid := uint64(10)
	rows := []repository.FAQWithProject{
		{
			FAQ:         model.FAQ{ID: 1, Question: "Q1", Answer: "A1", ProjectID: &pid, IsGlobal: true, SortOrder: 0},
			ProjectName: "Test Project",
			ProjectSlug: "test-project",
		},
		{
			FAQ:         model.FAQ{ID: 2, Question: "Q2", Answer: "A2", ProjectID: nil, IsGlobal: false, SortOrder: 1},
			ProjectName: "",
			ProjectSlug: "",
		},
	}

	responses := toFAQResponses(rows)
	if len(responses) != 2 {
		t.Fatalf("expected 2 responses, got %d", len(responses))
	}
	if responses[0].ProjectName != "Test Project" {
		t.Errorf("expected ProjectName='Test Project', got '%s'", responses[0].ProjectName)
	}
	if responses[0].ProjectSlug != "test-project" {
		t.Errorf("expected ProjectSlug='test-project', got '%s'", responses[0].ProjectSlug)
	}
	if responses[0].ProjectID == nil || *responses[0].ProjectID != 10 {
		t.Error("expected ProjectID=10")
	}
	if responses[1].ProjectName != "" {
		t.Errorf("expected empty ProjectName, got '%s'", responses[1].ProjectName)
	}
	if _, ok := interface{}(responses[0]).(dto.FAQResponse); !ok {
		t.Error("expected dto.FAQResponse type")
	}
}
```

- [ ] **Step 2: 更新 handler/faq_handler_test.go — 替换 mock 和测试**

替换文件全部内容：

```go
package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"mygo-immigration/backend/internal/dto"
	"mygo-immigration/backend/internal/model"
	"mygo-immigration/backend/internal/repository"
	"mygo-immigration/backend/internal/service"

	"github.com/gin-gonic/gin"
)

// handlerMockFAQRepo implements repository.FAQRepository.
type handlerMockFAQRepo struct {
	findAll func(params repository.FAQQueryParams) ([]repository.FAQWithProject, int64, error)
	create  func(faq *model.FAQ) error
	update  func(faq *model.FAQ) error
	delete  func(id uint64) error
}

func (m *handlerMockFAQRepo) FindAll(params repository.FAQQueryParams) ([]repository.FAQWithProject, int64, error) {
	if m.findAll != nil {
		return m.findAll(params)
	}
	return nil, 0, nil
}
func (m *handlerMockFAQRepo) Create(faq *model.FAQ) error {
	if m.create != nil {
		return m.create(faq)
	}
	return nil
}
func (m *handlerMockFAQRepo) Update(faq *model.FAQ) error {
	if m.update != nil {
		return m.update(faq)
	}
	return nil
}
func (m *handlerMockFAQRepo) Delete(id uint64) error {
	if m.delete != nil {
		return m.delete(id)
	}
	return nil
}
func (m *handlerMockFAQRepo) Search(keyword string) ([]model.FAQ, error) { return nil, nil }

func TestFAQHandler_ListFAQs(t *testing.T) {
	mockRepo := &handlerMockFAQRepo{
		findAll: func(params repository.FAQQueryParams) ([]repository.FAQWithProject, int64, error) {
			return []repository.FAQWithProject{
				{FAQ: model.FAQ{ID: 1, Question: "How to apply?", Answer: "Fill the form."}},
				{FAQ: model.FAQ{ID: 2, Question: "Processing time?", Answer: "2-3 months."}},
			}, 2, nil
		},
	}

	faqSvc := service.NewFAQService(mockRepo)
	h := &Handler{svc: &service.Service{FAQ: faqSvc}}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = makeGetRequest("/api/v1/faqs")

	h.ListFAQs(c)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d, body: %s", w.Code, w.Body.String())
	}
}

func TestFAQHandler_ListFAQs_WithProjectID(t *testing.T) {
	mockRepo := &handlerMockFAQRepo{
		findAll: func(params repository.FAQQueryParams) ([]repository.FAQWithProject, int64, error) {
			if params.ProjectID == nil || *params.ProjectID != 5 {
				t.Errorf("expected ProjectID=5 in params, got %v", params.ProjectID)
			}
			return []repository.FAQWithProject{}, 0, nil
		},
	}

	faqSvc := service.NewFAQService(mockRepo)
	h := &Handler{svc: &service.Service{FAQ: faqSvc}}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = makeGetRequest("/api/v1/faqs?project_id=5")

	h.ListFAQs(c)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

func TestFAQHandler_ListFAQs_ResponseContainsProjectName(t *testing.T) {
	pid := uint64(10)
	mockRepo := &handlerMockFAQRepo{
		findAll: func(params repository.FAQQueryParams) ([]repository.FAQWithProject, int64, error) {
			return []repository.FAQWithProject{
				{
					FAQ:         model.FAQ{ID: 1, Question: "Q1", Answer: "A1", ProjectID: &pid, IsGlobal: false},
					ProjectName: "EB-5 投资移民",
					ProjectSlug: "eb5",
				},
			}, 1, nil
		},
	}

	faqSvc := service.NewFAQService(mockRepo)
	h := &Handler{svc: &service.Service{FAQ: faqSvc}}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = makeGetRequest("/api/v1/faqs")

	h.ListFAQs(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	var resp dto.Response
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}
	// Data is []interface{} from JSON unmarshal — check first item
	arr, ok := resp.Data.([]interface{})
	if !ok || len(arr) == 0 {
		t.Fatal("expected data array with items")
	}
	item := arr[0].(map[string]interface{})
	if item["project_name"] != "EB-5 投资移民" {
		t.Errorf("expected project_name='EB-5 投资移民', got '%v'", item["project_name"])
	}
	if item["project_slug"] != "eb5" {
		t.Errorf("expected project_slug='eb5', got '%v'", item["project_slug"])
	}
}

func TestFAQHandler_ListFAQs_Empty(t *testing.T) {
	mockRepo := &handlerMockFAQRepo{
		findAll: func(params repository.FAQQueryParams) ([]repository.FAQWithProject, int64, error) {
			return []repository.FAQWithProject{}, 0, nil
		},
	}

	faqSvc := service.NewFAQService(mockRepo)
	h := &Handler{svc: &service.Service{FAQ: faqSvc}}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = makeGetRequest("/api/v1/faqs")

	h.ListFAQs(c)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200 for empty list, got %d", w.Code)
	}
}

func TestFAQHandler_ListFAQs_ServiceError(t *testing.T) {
	mockRepo := &handlerMockFAQRepo{
		findAll: func(params repository.FAQQueryParams) ([]repository.FAQWithProject, int64, error) {
			return nil, 0, errors.New("db error")
		},
	}

	faqSvc := service.NewFAQService(mockRepo)
	h := &Handler{svc: &service.Service{FAQ: faqSvc}}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = makeGetRequest("/api/v1/faqs")

	h.ListFAQs(c)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status 500 for service error, got %d", w.Code)
	}
}

func TestFAQHandler_AdminListFAQs_Success(t *testing.T) {
	mockRepo := &handlerMockFAQRepo{
		findAll: func(params repository.FAQQueryParams) ([]repository.FAQWithProject, int64, error) {
			return []repository.FAQWithProject{
				{FAQ: model.FAQ{ID: 1, Question: "Q1", Answer: "A1"}},
				{FAQ: model.FAQ{ID: 2, Question: "Q2", Answer: "A2"}},
			}, 2, nil
		},
	}

	faqSvc := service.NewFAQService(mockRepo)
	h := &Handler{svc: &service.Service{FAQ: faqSvc}}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = makeGetRequest("/api/v1/admin/faqs?page=1&per_page=10")

	h.AdminListFAQs(c)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d, body: %s", w.Code, w.Body.String())
	}
}

func TestFAQHandler_AdminListFAQs_WithSearch(t *testing.T) {
	mockRepo := &handlerMockFAQRepo{
		findAll: func(params repository.FAQQueryParams) ([]repository.FAQWithProject, int64, error) {
			if params.Search != "apply" {
				t.Errorf("expected search='apply', got '%s'", params.Search)
			}
			return []repository.FAQWithProject{}, 0, nil
		},
	}

	faqSvc := service.NewFAQService(mockRepo)
	h := &Handler{svc: &service.Service{FAQ: faqSvc}}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = makeGetRequest("/api/v1/admin/faqs?page=1&per_page=10&search=apply")

	h.AdminListFAQs(c)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

func TestFAQHandler_AdminListFAQs_ServiceError(t *testing.T) {
	mockRepo := &handlerMockFAQRepo{
		findAll: func(params repository.FAQQueryParams) ([]repository.FAQWithProject, int64, error) {
			return nil, 0, errors.New("db error")
		},
	}

	faqSvc := service.NewFAQService(mockRepo)
	h := &Handler{svc: &service.Service{FAQ: faqSvc}}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = makeGetRequest("/api/v1/admin/faqs?page=1&per_page=10")

	h.AdminListFAQs(c)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status 500, got %d", w.Code)
	}
}

func TestFAQHandler_CreateFAQ_Success(t *testing.T) {
	mockRepo := &handlerMockFAQRepo{
		create: func(faq *model.FAQ) error {
			faq.ID = 100
			return nil
		},
	}

	faqSvc := service.NewFAQService(mockRepo)
	h := &Handler{svc: &service.Service{FAQ: faqSvc}}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = makePostRequest("/api/v1/admin/faqs", model.FAQ{
		Question: "New question?",
		Answer:   "New answer.",
	})

	h.CreateFAQ(c)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status 201, got %d, body: %s", w.Code, w.Body.String())
	}
}

func TestFAQHandler_CreateFAQ_InvalidJSON(t *testing.T) {
	h := &Handler{}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/admin/faqs", nil)
	c.Request.Header.Set("Content-Type", "application/json")

	h.CreateFAQ(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400 for invalid JSON, got %d", w.Code)
	}
}

func TestFAQHandler_CreateFAQ_MissingFields(t *testing.T) {
	faqSvc := service.NewFAQService(&handlerMockFAQRepo{})
	h := &Handler{svc: &service.Service{FAQ: faqSvc}}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = makePostRequest("/api/v1/admin/faqs", model.FAQ{})

	h.CreateFAQ(c)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status 500 for missing fields, got %d", w.Code)
	}
}

func TestFAQHandler_UpdateFAQ_Success(t *testing.T) {
	mockRepo := &handlerMockFAQRepo{
		update: func(faq *model.FAQ) error {
			return nil
		},
	}

	faqSvc := service.NewFAQService(mockRepo)
	h := &Handler{svc: &service.Service{FAQ: faqSvc}}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = makePostRequest("/api/v1/admin/faqs/5", model.FAQ{
		Question: "Updated Q",
		Answer:   "Updated A",
	})
	c.Params = gin.Params{{Key: "id", Value: "5"}}

	h.UpdateFAQ(c)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d, body: %s", w.Code, w.Body.String())
	}
}

func TestFAQHandler_UpdateFAQ_InvalidID(t *testing.T) {
	h := &Handler{}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = makePostRequest("/api/v1/admin/faqs/abc", model.FAQ{})
	c.Params = gin.Params{{Key: "id", Value: "abc"}}

	h.UpdateFAQ(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400 for invalid id, got %d", w.Code)
	}
}

func TestFAQHandler_DeleteFAQ_Success(t *testing.T) {
	mockRepo := &handlerMockFAQRepo{
		delete: func(id uint64) error {
			return nil
		},
	}

	faqSvc := service.NewFAQService(mockRepo)
	h := &Handler{svc: &service.Service{FAQ: faqSvc}}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = makeGetRequest("/api/v1/admin/faqs/5")
	c.Params = gin.Params{{Key: "id", Value: "5"}}

	h.DeleteFAQ(c)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d, body: %s", w.Code, w.Body.String())
	}
}

func TestFAQHandler_DeleteFAQ_InvalidID(t *testing.T) {
	h := &Handler{}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = makeGetRequest("/api/v1/admin/faqs/abc")
	c.Params = gin.Params{{Key: "id", Value: "abc"}}

	h.DeleteFAQ(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400 for invalid id, got %d", w.Code)
	}
}
```

- [ ] **Step 3: 更新 service/search_svc_test.go — 更新 mock 实现新接口**

替换 `searchMockFAQRepo` 结构体及其方法（约第 10-27 行）为：

```go
// searchMockFAQRepo implements repository.FAQRepository for search service tests.
type searchMockFAQRepo struct {
	searchFn func(keyword string) ([]model.FAQ, error)
}

func (m *searchMockFAQRepo) FindAll(params repository.FAQQueryParams) ([]repository.FAQWithProject, int64, error) {
	return nil, 0, nil
}
func (m *searchMockFAQRepo) Create(faq *model.FAQ) error { return nil }
func (m *searchMockFAQRepo) Update(faq *model.FAQ) error { return nil }
func (m *searchMockFAQRepo) Delete(id uint64) error      { return nil }
func (m *searchMockFAQRepo) Search(keyword string) ([]model.FAQ, error) {
	if m.searchFn != nil {
		return m.searchFn(keyword)
	}
	return nil, nil
}
```

注意：在文件顶部 import 块中加入 `"mygo-immigration/backend/internal/repository"`（如果尚未引入）。

- [ ] **Step 4: 运行后端全部测试**

```bash
cd backend && go test ./... -v 2>&1
```
Expected: all tests pass.

---

### Task 7: 更新前端管理 FAQ 页 — 项目筛选 + project_name 列

**Files:**
- Modify: `frontend/pages/admin/faqs.vue`

- [ ] **Step 1: 更新 script 部分 — 修改类型定义、添加筛选器、修复请求逻辑**

在 `<script setup>` 中，修改以下部分：

```typescript
// 修改 Faq 接口（第 102-110 行附近）：
interface Faq {
  id: string;
  question: string;
  answer: string;
  project_id: string | null;
  project_name: string;     // 改：原为 project: string
  is_global: boolean;
  sort_order: number;
}

// 在 searchQuery ref 后面新增项目筛选器（第 124 行 searchQuery 之后）：
const projectFilter = ref<string | null>(null);

// 修改 loadList（第 160 行附近）：
const loadList = async () => {
  loading.value = true;
  try {
    const api = useApi();
    let url = `/admin/faqs?page=${page.value}&per_page=${pageSize.value}`;
    if (searchQuery.value) url += `&search=${encodeURIComponent(searchQuery.value)}`;
    if (projectFilter.value) url += `&project_id=${projectFilter.value}`;
    const data = await api<{ items: Faq[]; total: number }>(url);
    list.value = data.items ?? [];
    total.value = data.total ?? 0;
  } catch {
    list.value = [];
    ElMessage.error('加载FAQ列表失败');
  } finally {
    loading.value = false;
  }
};

// 修改 onSearch（第 127 行附近）—— 筛选器变化时也要重置 page：
const onFilterChange = () => {
  page.value = 1;
  loadList();
};
```

- [ ] **Step 2: 更新 template — 添加项目下拉筛选器，修复表格列**

在搜索框的 `<div class="admin-toolbar">` 内，搜索框后面添加项目筛选器：

```html
<div class="admin-toolbar">
  <div class="admin-toolbar-row">
    <el-input
      v-model="searchQuery"
      placeholder="搜索问题..."
      :prefix-icon="Search"
      clearable
      class="admin-search-input"
      @input="onSearch"
    />
    <el-select
      v-model="projectFilter"
      placeholder="按项目筛选"
      clearable
      class="admin-project-filter"
      @change="onFilterChange"
    >
      <el-option
        v-for="p in projects"
        :key="p.id"
        :label="p.name"
        :value="String(p.id)"
      />
    </el-select>
  </div>
</div>
```

修改表格"所属项目"列（第 26 行）：

```html
<el-table-column prop="project_name" label="所属项目" width="160" />
```

- [ ] **Step 3: 添加样式**

在 `<style scoped>` 块中添加：

```css
.admin-toolbar-row {
  display: flex;
  gap: 12px;
  align-items: center;
}
.admin-project-filter {
  width: 200px;
  flex-shrink: 0;
}
```

- [ ] **Step 4: 验证前端类型检查**

```bash
cd frontend && npx nuxi typecheck 2>&1
```
Expected: no new type errors.

---

### Task 8: 更新前端公开 FAQ 页 — 动态筛选按钮，移除 mock

**Files:**
- Modify: `frontend/pages/faq.vue`

- [ ] **Step 1: 替换 script 部分**

```typescript
const breadcrumbs = [
  { label: '首页', link: '/' },
  { label: '常见问题' },
];

useSeo({
  title: '常见问题',
  description: 'MyGo移民常见问题解答，涵盖美国EB-5、香港投资移民、巴拿马购房移民等投资移民相关问题。',
  breadcrumbs,
});

const activeFilter = ref('all');

interface FaqItem {
  id: number;
  question: string;
  answer: string;
  project_id: number | null;
  project_name: string;
  project_slug: string;
  is_global: boolean;
  sort_order: number;
}

const { data, pending, error } = await useFetch<{ data: FaqItem[] }>('/api/v1/faqs');

const allFaqs = computed<FaqItem[]>(() => {
  return (data.value as { data?: FaqItem[] })?.data ?? [];
});

// Dynamically extract project filters from actual FAQ data.
const projectFilters = computed(() => {
  const seen = new Set<string>();
  const result: { slug: string; label: string }[] = [];
  for (const faq of allFaqs.value) {
    if (faq.project_slug && !seen.has(faq.project_slug)) {
      seen.add(faq.project_slug);
      result.push({ slug: faq.project_slug, label: faq.project_name || faq.project_slug });
    }
  }
  return result;
});

const filteredFaqs = computed(() => {
  if (activeFilter.value === 'all') {
    return allFaqs.value;
  }
  return allFaqs.value.filter((faq) => faq.project_slug === activeFilter.value);
});

// FAQPage structured data
useHead(() => {
  if (allFaqs.value.length === 0) return {};
  return {
    script: [
      {
        type: 'application/ld+json',
        innerHTML: JSON.stringify({
          '@context': 'https://schema.org',
          '@type': 'FAQPage',
          mainEntity: allFaqs.value.map((faq) => ({
            '@type': 'Question',
            name: faq.question,
            acceptedAnswer: {
              '@type': 'Answer',
              text: faq.answer,
            },
          })),
        }),
      },
    ],
  };
});
```

- [ ] **Step 2: 更新 template 中的筛选按钮循环**

将第 18-26 行替换为：

```html
<div class="faq-filters">
  <button
    class="filter-btn"
    :class="{ active: activeFilter === 'all' }"
    @click="activeFilter = 'all'"
  >
    全部
  </button>
  <button
    v-for="filter in projectFilters"
    :key="filter.slug"
    class="filter-btn"
    :class="{ active: activeFilter === filter.slug }"
    @click="activeFilter = filter.slug"
  >
    {{ filter.label }}
  </button>
</div>
```

- [ ] **Step 3: 修改错误状态文案**

将第 31 行错误状态从 `{{ error }}` 改为：

```html
<div v-else-if="error" class="error-state">加载常见问题失败，请稍后重试</div>
```

- [ ] **Step 4: 验证前端类型检查**

```bash
cd frontend && npx nuxi typecheck 2>&1
```
Expected: no new type errors.

---

### Task 9: 端到端验证

- [ ] **Step 1: 启动后端并手动测试 API**

```bash
# 确保 MySQL 运行中
cd backend && go test ./internal/service -run TestFAQ -v
```

- [ ] **Step 2: 全部后端测试通过**

```bash
cd backend && go test ./... -v 2>&1
```
Expected: all tests pass.

---

## 变更文件总览

| 文件 | 操作 |
|------|------|
| `backend/internal/dto/faq_response.go` | 新建 |
| `backend/internal/repository/interfaces.go` | 修改 FAQRepository + FAQQueryParams |
| `backend/internal/repository/faq_repo.go` | 重写 FindAll，新增 FAQWithProject |
| `backend/internal/service/faq_svc.go` | 修改 List/AdminList 签名，新增 toFAQResponses |
| `backend/internal/service/faq_svc_test.go` | 更新 mock + 测试 |
| `backend/internal/handler/faq_handler.go` | 读取新 query 参数 |
| `backend/internal/handler/faq_handler_test.go` | 更新 mock + 测试 |
| `backend/internal/service/search_svc_test.go` | 更新 mock（新接口） |
| `frontend/pages/faq.vue` | 动态筛选按钮，移除 mock |
| `frontend/pages/admin/faqs.vue` | 项目筛选器 + project_name 列 |

**无需改动：**
- `frontend/pages/projects/[slug].vue` — 项目详情的 FAQ 区块已存在，数据通过项目 API 的 Preload("FAQs") 返回
- `backend/internal/router/router.go` — 现有路由已够用
- 数据库 — 无 schema 变更
