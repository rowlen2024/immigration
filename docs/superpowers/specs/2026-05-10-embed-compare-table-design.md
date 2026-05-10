# Embed Compare Table in Project Detail Page

**Date**: 2026-05-10
**Status**: approved

## Goal

Replace the "查看对比详情" link button on the project detail page with an inline N-way comparison table. Support 2-5 projects with horizontal scroll.

## Current State

- `projects/[slug].vue` line 77-84: shows `<NuxtLink>` button when `compare_config.compare_with.length >= 2`
- `compare.vue`: standalone page with 2 dropdowns, hardcoded 3-column table (label | A | B)
- Backend `CompareRows`: hardcodes `a, b := projects[0], projects[1]`, returns `CompareRow{A, B}` fixed to 2 projects
- Backend `Compare` function: limits to max 5 slugs (already supports N)

## Changes

### 1. Backend: N-way CompareRow

**File**: `backend/internal/service/project_svc.go`

- `CompareRow.A`, `CompareRow.B` → `CompareRow.Values []string`
- `CompareRows`: iterate all projects instead of `a, b := projects[0], projects[1]`
  - Build `projInfo` dynamically for all projects
  - For each row, collect values from all projects into `Values` slice
- Remove effective 2-project-only restriction (already capped at 5 by `Compare`)

**File**: `backend/internal/handler/project_handler_test.go`

- `TestProjectHandler_CompareProjects_TooMany` (5+ slugs → 500): update expected result to success (3 projects is fine), or rename/repurpose to test 6+ slugs → 500

### 2. Frontend: Embed table in project detail page

**File**: `frontend/pages/projects/[slug].vue`

- Replace the `<NuxtLink>` button (lines 77-84) with a full comparison table
- Fetch `/api/v1/projects/compare?slugs=...` when `compare_config` is present
- Table: header row (对比项目 + project names), body rows (label + values)
- First column sticky left, `overflow-x: auto` on wrapper
- Each data column `min-width: 180px` to prevent squeeze
- Handle loading / empty / error states

### 3. Frontend: Adapt compare.vue to new API

**File**: `frontend/pages/compare.vue`

- `row.a` → `row.values[0]`, `row.b` → `row.values[1]`
- Table rendering already uses `v-for` for data columns, scope stays at 2 projects for the selector page

## Non-Goals

- `compare/[a]-vs-[b].vue` hardcoded detail page: not in scope (uses mock data)
- Removing `compare.vue` selector page: it remains as standalone entry point

## Test Plan

- Backend: `go test ./internal/... -count=1` — verify CompareRows with 2, 3, and 5 slugs
- Frontend: `npx nuxi typecheck` — verify no new type errors
- Manual: open a project detail page with `compare_config`, verify table renders with correct data and horizontal scroll works for 3+ projects
