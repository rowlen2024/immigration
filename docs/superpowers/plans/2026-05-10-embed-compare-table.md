# Embed Compare Table in Project Detail — Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Replace the "查看对比详情" link on the project detail page with an inline N-way comparison table, and generalize the backend compare API from 2-project to N-project.

**Architecture:** Backend `CompareRow` gains `Values []string` instead of hardcoded `A`/`B`. `CompareRows` iterates all projects dynamically. Frontend detail page fetches compare API and renders a horizontally-scrollable N+1 column table with sticky first column. `compare.vue` adapted to use `values[0]`/`values[1]`.

**Tech Stack:** Go (Gin + GORM), Nuxt 3 / Vue 3, Element Plus (unused in this change)

---

### Task 1: Backend — Change CompareRow struct to N-way

**Files:**
- Modify: `backend/internal/service/project_svc.go:71-75`

- [ ] **Step 1: Replace CompareRow struct**

```go
// CompareRow represents a single comparison row.
type CompareRow struct {
	Label  string   `json:"label"`
	Values []string `json:"values"`
}
```

- [ ] **Step 2: Verify the change compiles (will fail until CompareRows is updated in Task 2)**

Run: `cd backend && go build ./... 2>&1`
Expected: compilation errors in `CompareRows` referencing `.A` and `.B` — expected, will fix in Task 2.

---

### Task 2: Backend — Rewrite CompareRows for N-way

**Files:**
- Modify: `backend/internal/service/project_svc.go:89-122`

- [ ] **Step 1: Replace the CompareRows function body**

```go
// CompareRows returns formatted comparison rows for N projects (2–5).
func (s *ProjectService) CompareRows(slugs []string) (*CompareResult, error) {
	projects, err := s.Compare(slugs)
	if err != nil {
		return nil, err
	}
	if len(projects) < 2 {
		return nil, errors.New("需要至少两个项目进行对比")
	}

	projInfo := make([]CompareProject, len(projects))
	for i, p := range projects {
		projInfo[i] = CompareProject{Title: p.Name, Slug: p.Slug}
	}

	rows := []CompareRow{
		{Label: "投资金额", Values: pluck(projects, func(p model.Project) string { return p.InvestmentAmount })},
		{Label: "办理周期", Values: pluck(projects, func(p model.Project) string { return p.ProcessingPeriod })},
		{Label: "适合人群", Values: pluck(projects, func(p model.Project) string { return p.TargetCrowd })},
		{Label: "申请条件", Values: pluck(projects, func(p model.Project) string { return joinRequirements(p.Requirements) })},
		{Label: "费用总计", Values: pluck(projects, func(p model.Project) string { return p.CostsTotal })},
		{Label: "流程步骤", Values: pluck(projects, func(p model.Project) string { return fmt.Sprintf("%d 个阶段", len(p.TimelinePhases)) })},
	}

	return &CompareResult{Projects: projInfo, Rows: rows}, nil
}

func pluck(projects []model.Project, fn func(model.Project) string) []string {
	values := make([]string, len(projects))
	for i, p := range projects {
		values[i] = fn(p)
	}
	return values
}
```

- [ ] **Step 2: Build to verify no compilation errors**

Run: `cd backend && go build ./... 2>&1`
Expected: clean build, no errors.

---

### Task 3: Backend — Update test for N-way comparison

**Files:**
- Modify: `backend/internal/handler/project_handler_test.go:227-241`

- [ ] **Step 1: Repurpose the TooMany test — 6 slugs still returns 500, but 3 slugs should succeed**

Replace the existing `TestProjectHandler_CompareProjects_TooMany` (lines 227-241):

```go
func TestProjectHandler_CompareProjects_TooMany(t *testing.T) {
	mockRepo := &handlerMockProjectRepo{}
	projectSvc := service.NewProjectService(mockRepo)
	h := &Handler{svc: &service.Service{Project: projectSvc}}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = makeGetRequest("/api/v1/projects/compare?slugs=a,b,c,d,e,f")

	h.CompareProjects(c)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status 500 for too many slugs, got %d", w.Code)
	}
}
```

This test stays as-is because the `Compare` function (line 60-62) still limits to 5 slugs. But the mockRepo's `findBySlugs` is nil, so `Compare` will check `len(slugs) > 5` (6 > 5 → error) and return before calling `FindBySlugs`. The test remains valid.

- [ ] **Step 2: Add a test for 3-way comparison**

Add after the TooMany test:

```go
func TestProjectHandler_CompareProjects_ThreeWay(t *testing.T) {
	mockRepo := &handlerMockProjectRepo{
		findBySlugs: func(slugs []string) ([]model.Project, error) {
			return []model.Project{
				{ID: 1, Name: "Project A", Slug: "a", InvestmentAmount: "100万", ProcessingPeriod: "12月", TargetCrowd: "投资者"},
				{ID: 2, Name: "Project B", Slug: "b", InvestmentAmount: "200万", ProcessingPeriod: "24月", TargetCrowd: "企业家"},
				{ID: 3, Name: "Project C", Slug: "c", InvestmentAmount: "300万", ProcessingPeriod: "36月", TargetCrowd: "高净值"},
			}, nil
		},
	}
	projectSvc := service.NewProjectService(mockRepo)
	h := &Handler{svc: &service.Service{Project: projectSvc}}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = makeGetRequest("/api/v1/projects/compare?slugs=a,b,c")

	h.CompareProjects(c)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d, body: %s", w.Code, w.Body.String())
	}

	var body map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}
	data := body["data"].(map[string]interface{})
	projects := data["projects"].([]interface{})
	if len(projects) != 3 {
		t.Errorf("expected 3 projects in response, got %d", len(projects))
	}
	rows := data["rows"].([]interface{})
	firstRow := rows[0].(map[string]interface{})
	values := firstRow["values"].([]interface{})
	if len(values) != 3 {
		t.Errorf("expected 3 values per row, got %d", len(values))
	}
}
```

Note: add `"encoding/json"` to imports if not already present (check existing imports — `json` is already used in the test file).

- [ ] **Step 3: Run tests**

Run: `cd backend && go test ./internal/handler/... -run "Compare" -v 2>&1`
Expected: all Compare tests pass (Success, NoSlugs, TooMany, ThreeWay).

---

### Task 4: Backend — Full test suite

**Files:** None new.

- [ ] **Step 1: Run all backend tests**

Run: `cd backend && go test ./... -count=1 2>&1`
Expected: all tests pass.

- [ ] **Step 2: Commit backend changes**

```bash
git add backend/internal/service/project_svc.go backend/internal/handler/project_handler_test.go
git commit -m "feat: generalize CompareRows to support N-way project comparison"
```

---

### Task 5: Frontend — Adapt compare.vue to new Values API

**Files:**
- Modify: `frontend/pages/compare.vue:53-57,73-86,103-106,139-141`

- [ ] **Step 1: Update ComparisonData interface and table rendering**

The `ComparisonData` interface (lines 103-106) and `getColClass` (lines 139-141) stay the same shape since compare.vue only does 2 projects. But `row.a` / `row.b` references must become `row.values[0]` / `row.values[1]`.

In `<template>` lines 53-57, change:
```vue
              <tr v-for="row in comparison.rows" :key="row.label">
                <td class="row-label">{{ row.label }}</td>
                <td :class="getColClass(row.a, row.b, 'a')">{{ row.a }}</td>
                <td :class="getColClass(row.a, row.b, 'b')">{{ row.b }}</td>
              </tr>
```
to:
```vue
              <tr v-for="row in comparison.rows" :key="row.label">
                <td class="row-label">{{ row.label }}</td>
                <td class="col-a">{{ row.values[0] }}</td>
                <td class="col-b">{{ row.values[1] }}</td>
              </tr>
```

Remove `getColClass` function (lines 139-141) since it's no longer used:
```typescript
// DELETE these lines:
const getColClass = (_valueA: string, _valueB: string, col: 'a' | 'b') => {
  if (col === 'a') return 'col-a';
  return 'col-b';
};
```

- [ ] **Step 2: Verify typecheck**

Run: `cd frontend && npx nuxi typecheck 2>&1 | grep -i "compare"` 
Expected: no compare-related errors (pre-existing project-wide auto-import errors are unrelated).

---

### Task 6: Frontend — Embed compare table in project detail page

**Files:**
- Modify: `frontend/pages/projects/[slug].vue:77-84,96-169`

- [ ] **Step 1: Add compare API fetch composable**

In `<script setup>`, after the existing `useFetch` for project data and the `project` computed, add:

```typescript
// Compare data fetch
interface CompareRowData {
  label: string;
  values: string[];
}

interface CompareTableData {
  projects: Array<{ title: string; slug: string }>;
  rows: CompareRowData[];
}

const compareSlugs = computed(() => {
  const cfg = project.value.compare_config;
  if (!cfg || !cfg.compare_with || cfg.compare_with.length < 2) return '';
  return cfg.compare_with.join(',');
});

const {
  data: compareRaw,
  pending: comparePending,
  error: compareErrorRaw,
} = useFetch<{ data: CompareTableData }>(
  () => `/api/v1/projects/compare?slugs=${compareSlugs.value}`,
);

const compareData = computed<CompareTableData | null>(() => {
  const raw = compareRaw.value as any;
  if (raw?.data?.rows) return raw.data;
  if (raw?.rows) return raw as CompareTableData;
  return null;
});

const compareError = computed(() =>
  compareErrorRaw.value ? '加载对比数据失败' : null
);
```

- [ ] **Step 2: Replace the link button with inline table**

Replace lines 77-84 (current compare section):
```vue
        <section v-if="project.compare_config && project.compare_config.compare_with.length >= 2" class="detail-section">
          <h2 class="detail-section-title">项目对比</h2>
          <div class="compare-link-wrap">
            <NuxtLink :to="`/compare?slugs=${project.compare_config.compare_with.join(',')}`" class="btn-primary">
              查看对比详情
            </NuxtLink>
          </div>
        </section>
```

With:
```vue
        <section v-if="project.compare_config && project.compare_config.compare_with.length >= 2" class="detail-section">
          <h2 class="detail-section-title">项目对比</h2>
          <div v-if="comparePending" class="compare-loading">加载对比数据...</div>
          <div v-else-if="compareError" class="compare-error">{{ compareError }}</div>
          <div v-else-if="compareData" class="compare-table-wrap">
            <table class="compare-table">
              <thead>
                <tr>
                  <th class="compare-label-col">对比项目</th>
                  <th v-for="(proj, i) in compareData.projects" :key="i" :class="`compare-col-${i}`">
                    {{ proj.title }}
                  </th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="row in compareData.rows" :key="row.label">
                  <td class="compare-label">{{ row.label }}</td>
                  <td v-for="(val, j) in row.values" :key="j" :class="`compare-col-${j}`">{{ val }}</td>
                </tr>
              </tbody>
            </table>
          </div>
        </section>
```

- [ ] **Step 3: Add styles to `<style scoped>` section**

Add after the existing `.compare-link-wrap` style block (or replace it):

```css
.compare-table-wrap {
  overflow-x: auto;
  -webkit-overflow-scrolling: touch;
}

.compare-table {
  width: 100%;
  border-collapse: collapse;
  font-size: 14px;
  min-width: max-content;
}

.compare-table thead {
  background-color: var(--primary);
  color: var(--bg-white);
}

.compare-table thead th {
  padding: 14px 16px;
  font-weight: 600;
  text-align: left;
  white-space: nowrap;
}

.compare-table thead th:first-child {
  border-radius: var(--radius-md) 0 0 0;
}

.compare-table thead th:last-child {
  border-radius: 0 var(--radius-md) 0 0;
}

.compare-label-col {
  position: sticky;
  left: 0;
  z-index: 2;
  background-color: var(--primary);
  min-width: 120px;
}

.compare-table td {
  padding: 12px 16px;
  border-bottom: 1px solid var(--border-color);
  color: var(--text-secondary);
  line-height: 1.6;
  min-width: 180px;
}

.compare-table tbody tr:nth-child(even) {
  background-color: var(--bg-light);
}

.compare-label {
  font-weight: 600;
  color: var(--text-primary);
  white-space: nowrap;
  position: sticky;
  left: 0;
  z-index: 1;
  min-width: 120px;
}

.compare-table tbody tr:nth-child(even) .compare-label {
  background-color: var(--bg-light);
}

.compare-table tbody tr:nth-child(odd) .compare-label {
  background-color: var(--bg-white);
}

.compare-loading,
.compare-error {
  text-align: center;
  padding: 24px;
  color: var(--text-light);
  font-size: 14px;
}

.compare-error {
  color: #c62828;
}
```

NOTE: Also remove the old `.compare-link-wrap` style block (line 373-375) since it's no longer used.

- [ ] **Step 4: Verify typecheck**

Run: `cd frontend && npx nuxi typecheck 2>&1 | grep -E "\[slug\]|compare"` 
Expected: no new errors from the modified files.

---

### Task 7: Frontend — Commit & verify

**Files:** None new.

- [ ] **Step 1: Commit frontend changes**

```bash
git add frontend/pages/compare.vue frontend/pages/projects/\[slug\].vue
git commit -m "feat: embed N-way compare table in project detail page"
```

- [ ] **Step 2: Final typecheck**

Run: `cd frontend && npx nuxi typecheck 2>&1`
Expected: only pre-existing auto-import errors; no new errors.

- [ ] **Step 3: Final backend test**

Run: `cd backend && go test ./... -count=1 2>&1`
Expected: all tests pass.
