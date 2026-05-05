# Project Detail Admin & Frontend-Backend Alignment

## Problem

1. Admin project editor only edits top-level Project fields. Sub-entities (Requirements, CostItems, TimelinePhases) cannot be managed through the admin UI.
2. Dedicated project pages (`usa/eb5.vue`, `hongkong/cies.vue`, `panama/property.vue`) use hardcoded fallback data because they don't unwrap the API response envelope and expect wrong field names.
3. Compare API returns raw `Project[]` but frontend expects `{ projects, rows }` format — hardcoded fallback data is always used.
4. Old `/usa/eb5`, `/hongkong/cies`, `/panama/property` URLs need redirects to `/projects/:slug`.

## Solution

Unify all project detail pages under `projects/[slug].vue`, add sub-entity CRUD to admin, and fix the compare API data format.

### 1. Backend: Sub-entity CRUD Endpoints (NEW)

Add admin-only endpoints for managing requirements, cost items, and timeline phases per project:

```
GET    /admin/projects/:id/requirements        # list
POST   /admin/projects/:id/requirements        # create
PUT    /admin/projects/:id/requirements/:rid    # update
DELETE /admin/projects/:id/requirements/:rid    # delete

GET    /admin/projects/:id/cost-items
POST   /admin/projects/:id/cost-items
PUT    /admin/projects/:id/cost-items/:cid
DELETE /admin/projects/:id/cost-items/:cid

GET    /admin/projects/:id/timeline-phases
POST   /admin/projects/:id/timeline-phases
PUT    /admin/projects/:id/timeline-phases/:tid
DELETE /admin/projects/:id/timeline-phases/:tid
```

GET uses `admin:read`, POST/PUT/DELETE use `projects:write`.

**Also:** Change `ProjectRepo.Update` from `db.Save(project)` to `db.Omit("Requirements", "CostItems", "TimelinePhases", "Milestones", "FAQs", "Cases").Save(project)` so the main project update doesn't touch associations (now managed via separate endpoints).

**Files to create/modify:**
- `backend/internal/repository/interfaces.go` — add SubEntityRepository interface
- `backend/internal/repository/sub_repo.go` (new) — CRUD for requirements/cost_items/timeline_phases
- `backend/internal/service/sub_svc.go` (new) — business logic
- `backend/internal/handler/sub_handler.go` (new) — HTTP handlers
- `backend/internal/router/router.go` — register new routes

All require `projects:write` RBAC. Return standard `{ code, message, data }` envelope.

### 2. Backend: Fix Compare API

**Current:** `GET /api/v1/projects/compare?slugs=a,b` returns `{ code, data: [Project, Project] }`

**New:** Transform to `{ code, data: { projects: [...], rows: [...] } }` where rows contain label-based comparisons across these dimensions:

| Label | Source (project A vs B) |
|-------|-------------------------|
| 投资金额 | `investment_amount` |
| 办理周期 | `processing_period` |
| 适合人群 | `target_crowd` |
| 投资金额（数值） | `investment_value` + `investment_currency` |
| 申请条件 | concatenated `requirements[].label` |
| 费用总计 | `costs_total` |
| 流程步骤数 | count of `timeline_phases` |

**Files to modify:**
- `backend/internal/service/project_svc.go` — add CompareRows method
- `backend/internal/handler/project_handler.go` — update CompareProjects handler

### 3. Frontend: Admin Project Editor Overhaul

Convert the single-form dialog in `admin/projects.vue` to a tabbed dialog with four panels:

**Tab 1: 基本信息** — Keep existing form fields (slug, name, country, investment_amount, etc.)

**Tab 2: 申请条件** — Table with columns (label, is_required, sort_order) + Add/Edit/Delete buttons. Inline editing via a small dialog or row editing.

**Tab 3: 费用明细** — Table with columns (name, amount, note, sort_order) + CRUD.

**Tab 4: 申请流程** — Table with columns (phase_number, title, description, duration, sort_order) + CRUD.

On edit dialog open, fetch sub-entities via new admin list endpoints. On save/create of sub-entities, call the new CRUD endpoints.

### 4. Frontend: Delete Dedicated Project Pages

Delete these files:
- `frontend/pages/usa/eb5.vue`
- `frontend/pages/hongkong/cies.vue`
- `frontend/pages/panama/property.vue`

### 5. Frontend: Old URL Redirects

Add redirects in `nuxt.config.ts`:
- `/usa/eb5` → `/projects/eb5`
- `/hongkong/cies` → `/projects/cies`
- `/panama/property` → `/projects/panama`

### 6. Frontend: Fix Compare Page

- Remove `getDefaultComparison()` hardcoded fallback
- Update to use the new API response format
- The `compare.vue` structure already matches the target format, just remove fallback and fix data access

### 7. Frontend: Remove Unused Milestones/Policy Sections

No change needed — `projects/[slug].vue` already correctly renders what the API returns.

## Files Changed

| File | Action |
|------|--------|
| `backend/internal/repository/interfaces.go` | Add SubRepository interfaces |
| `backend/internal/repository/sub_repo.go` | NEW |
| `backend/internal/service/sub_svc.go` | NEW |
| `backend/internal/handler/sub_handler.go` | NEW |
| `backend/internal/handler/project_handler.go` | Update CompareProjects |
| `backend/internal/service/project_svc.go` | Add compare transform logic |
| `backend/internal/router/router.go` | Register sub-entity routes |
| `frontend/pages/admin/projects.vue` | Major: tabbed editor + sub-entity CRUD |
| `frontend/pages/usa/eb5.vue` | DELETE |
| `frontend/pages/hongkong/cies.vue` | DELETE |
| `frontend/pages/panama/property.vue` | DELETE |
| `frontend/pages/compare.vue` | Remove hardcoded fallback, fix API data access |
| `frontend/nuxt.config.ts` | Add redirects |

## Not Changing

- Models (already complete)
- Database migrations and seed data
- `projects/[slug].vue` (data flow is correct)
- `frontend/components/project/*` (components are correct)
- FAQ admin page (already has independent management)
