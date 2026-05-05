# Admin UI Redesign Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Redesign all 12 admin pages with Modern Minimal style — slate color palette, dark sidebar with Lucide icons, borderless tables, drawer forms.

**Architecture:** Three-phase incremental rewrite. Phase 1 establishes new CSS tokens and layout shell. Phase 2 applies the unified list+drawer template to all 6 CRUD pages. Phase 3 rebuilds dashboard, homepage config, and polishes remaining special pages.

**Tech Stack:** Nuxt 3 (SPA), Element Plus, existing `lucideIcons.ts` composable, no new dependencies.

---

## File Map

| File | Action | Phase |
|------|--------|-------|
| `frontend/assets/css/variables.css` | Rewrite color tokens | 1 |
| `frontend/assets/css/admin.css` | Rewrite shared admin styles | 1 |
| `frontend/assets/css/global.css` | Minor updates | 1 |
| `frontend/layouts/admin.vue` | Full sidebar+topbar rewrite | 1 |
| `frontend/pages/admin/projects.vue` | List+drawer rewrite | 2 |
| `frontend/pages/admin/pages.vue` | List+drawer rewrite | 2 |
| `frontend/pages/admin/faqs.vue` | List+drawer rewrite | 2 |
| `frontend/pages/admin/cases.vue` | List+drawer rewrite | 2 |
| `frontend/pages/admin/users.vue` | List+drawer rewrite | 2 |
| `frontend/pages/admin/leads.vue` | List+drawer rewrite | 2 |
| `frontend/pages/admin/index.vue` | Dashboard rewrite | 3 |
| `frontend/pages/admin/homepage.vue` | Tab layout rewrite | 3 |
| `frontend/pages/admin/navigation.vue` | Icon button polish | 3 |
| `frontend/pages/admin/media.vue` | Grid polish | 3 |
| `frontend/pages/admin/settings.vue` | Layout polish | 3 |
| `frontend/pages/admin/login.vue` | Visual upgrade | 3 |

---

## Phase 1: Core Framework

### Task 1: Rewrite CSS Variables

**Files:**
- Modify: `frontend/assets/css/variables.css`

- [ ] **Step 1: Replace all color and design tokens**

Replace the entire content of `variables.css`:

```css
:root {
  /* Brand */
  --color-primary: #0f172a;
  --color-primary-hover: #1e293b;
  --color-accent: #e2a83e;
  --color-accent-hover: #c9952e;
  --color-accent-soft: #fef9c3;

  /* Surfaces */
  --color-bg-app: #f8fafc;
  --color-bg-surface: #ffffff;
  --color-bg-sidebar: #0f172a;
  --color-bg-sidebar-hover: rgba(255, 255, 255, 0.08);
  --color-bg-sidebar-active: rgba(255, 255, 255, 0.16);

  /* Text */
  --color-text: #0f172a;
  --color-text-secondary: #64748b;
  --color-text-muted: #94a3b8;
  --color-text-sidebar: rgba(255, 255, 255, 0.65);
  --color-text-sidebar-active: #ffffff;

  /* Borders */
  --color-border: #e2e8f0;
  --color-border-light: #f1f5f9;
  --color-border-sidebar: rgba(255, 255, 255, 0.08);

  /* Semantic */
  --color-success: #16a34a;
  --color-success-soft: #dcfce7;
  --color-warning: #d97706;
  --color-warning-soft: #fef3c7;
  --color-danger: #dc2626;
  --color-danger-soft: #fef2f2;
  --color-info: #3b82f6;
  --color-info-soft: #dbeafe;

  /* Radii */
  --radius-sm: 6px;
  --radius-md: 8px;
  --radius-lg: 12px;

  /* Shadows */
  --shadow-sm: 0 1px 2px rgba(0, 0, 0, 0.04);
  --shadow-md: 0 4px 12px rgba(0, 0, 0, 0.06);
  --shadow-lg: 0 8px 24px rgba(0, 0, 0, 0.08);

  /* Typography */
  --font-sans: 'PingFang SC', 'Microsoft YaHei', 'Helvetica Neue', sans-serif;
  --font-mono: 'SF Mono', 'Fira Code', 'Courier New', monospace;

  /* Layout */
  --max-width: 1200px;
  --header-height: 72px;
  --sidebar-width: 220px;
  --sidebar-collapsed-width: 56px;
  --topbar-height: 56px;

  /* Legacy aliases — keep old var names working during transition */
  --primary: var(--color-primary);
  --primary-dark: #0a0f1a;
  --primary-light: #1e293b;
  --accent: var(--color-accent);
  --accent-dark: var(--color-accent-hover);
  --accent-light: #ecc663;
  --text-primary: var(--color-text);
  --text-secondary: var(--color-text-secondary);
  --text-light: var(--color-text-muted);
  --bg-white: var(--color-bg-surface);
  --bg-light: var(--color-bg-app);
  --bg-gray: #f0f1f3;
  --border-color: var(--color-border);
}
```

- [ ] **Step 2: Verify file was written correctly**

Run: `head -5 frontend/assets/css/variables.css`
Expected: `:root {`

---

### Task 2: Rewrite Global CSS

**Files:**
- Modify: `frontend/assets/css/global.css`

- [ ] **Step 1: Update body background and button styles**

Replace the content of `global.css`:

```css
*,
*::before,
*::after {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

html {
  font-size: 16px;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
}

body {
  font-family: var(--font-sans);
  color: var(--color-text);
  background-color: var(--color-bg-surface);
  line-height: 1.6;
}

a {
  color: inherit;
  text-decoration: none;
}

ul,
ol {
  list-style: none;
}

img {
  max-width: 100%;
  height: auto;
  display: block;
}

.container {
  max-width: var(--max-width);
  margin: 0 auto;
  padding: 0 20px;
}

.btn-primary {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 10px 24px;
  background-color: var(--color-primary);
  color: #fff;
  border: none;
  border-radius: var(--radius-md);
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: background-color 0.2s ease, transform 0.15s ease;
}

.btn-primary:hover {
  background-color: var(--color-primary-hover);
}

.btn-outline {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 10px 24px;
  background-color: transparent;
  color: var(--color-primary);
  border: 1.5px solid var(--color-border);
  border-radius: var(--radius-md);
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
}

.btn-outline:hover {
  border-color: var(--color-primary);
  background-color: var(--color-bg-app);
}
```

---

### Task 3: Rewrite Shared Admin CSS

**Files:**
- Modify: `frontend/assets/css/admin.css`

- [ ] **Step 1: Replace with new shared admin styles**

Replace the entire content of `admin.css`:

```css
/* ===== Shared Admin Page Styles ===== */

/* Page structure */
.admin-page-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 24px;
}

.admin-page-title {
  font-size: 22px;
  font-weight: 600;
  color: var(--color-text);
}

.admin-section-title {
  font-size: 16px;
  font-weight: 600;
  color: var(--color-text);
  margin-bottom: 16px;
}

/* Toolbar (search + filter + actions) */
.admin-toolbar {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 16px;
}

.admin-toolbar .admin-search-input {
  width: 300px;
}

.admin-toolbar .admin-filter-select {
  width: 140px;
}

/* Pagination */
.admin-pagination-wrap {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}

/* Table wrapper */
.admin-table-wrap {
  background: var(--color-bg-surface);
  border-radius: var(--radius-md);
  border: 1px solid var(--color-border);
  overflow: hidden;
}

/* Stats grid */
.admin-stat-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 16px;
  margin-bottom: 24px;
}

@media (max-width: 1024px) {
  .admin-stat-grid { grid-template-columns: repeat(3, 1fr); }
}
@media (max-width: 767px) {
  .admin-stat-grid { grid-template-columns: repeat(2, 1fr); }
}

.admin-stat-card {
  background: var(--color-bg-surface);
  padding: 20px;
  border-radius: var(--radius-md);
  border: 1px solid var(--color-border-light);
  transition: transform 0.15s ease, box-shadow 0.15s ease;
}

.admin-stat-card:hover {
  transform: translateY(-2px);
  box-shadow: var(--shadow-md);
}

.admin-stat-card .stat-icon {
  width: 36px;
  height: 36px;
  border-radius: var(--radius-sm);
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 12px;
}

.admin-stat-label {
  font-size: 13px;
  color: var(--color-text-secondary);
  margin-bottom: 4px;
}

.admin-stat-value {
  font-size: 28px;
  font-weight: 700;
  color: var(--color-text);
  line-height: 1.2;
  margin-bottom: 4px;
}

.admin-stat-trend {
  font-size: 12px;
}

.admin-stat-trend.up { color: var(--color-success); }
.admin-stat-trend.down { color: var(--color-danger); }
.admin-stat-trend.neutral { color: var(--color-text-muted); }

/* Quick links grid */
.admin-links-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 12px;
}

@media (max-width: 1024px) {
  .admin-links-grid { grid-template-columns: repeat(2, 1fr); }
}

.admin-quick-link {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 16px;
  background: var(--color-bg-surface);
  border-radius: var(--radius-md);
  border: 1px solid var(--color-border-light);
  color: var(--color-text);
  font-size: 14px;
  font-weight: 500;
  transition: all 0.2s ease;
}

.admin-quick-link:hover {
  border-color: var(--color-primary);
  box-shadow: var(--shadow-sm);
}

/* Empty state */
.admin-empty-state {
  text-align: center;
  padding: 60px 20px;
}

.admin-empty-state .empty-icon {
  width: 48px;
  height: 48px;
  margin: 0 auto 16px;
  color: var(--color-text-muted);
}

.admin-empty-state .empty-title {
  font-size: 15px;
  font-weight: 500;
  color: var(--color-text-secondary);
  margin-bottom: 8px;
}

.admin-empty-state .empty-desc {
  font-size: 13px;
  color: var(--color-text-muted);
  margin-bottom: 20px;
}

/* Loading skeleton */
.admin-loading-state {
  padding: 40px 0;
}

/* Detail view (leads detail) */
.admin-detail {
  padding: 0;
}

.admin-detail-row {
  display: flex;
  padding: 10px 0;
  border-bottom: 1px solid var(--color-border-light);
}

.admin-detail-row:last-child {
  border-bottom: none;
}

.admin-detail-label {
  width: 100px;
  color: var(--color-text-secondary);
  font-size: 13px;
  flex-shrink: 0;
}

.admin-detail-value {
  color: var(--color-text);
  font-size: 14px;
}

/* Card styles */
.admin-card {
  background: var(--color-bg-surface);
  border-radius: var(--radius-md);
  border: 1px solid var(--color-border-light);
  padding: 24px;
  margin-bottom: 16px;
}

.admin-card-title {
  margin: 0;
  font-size: 16px;
  font-weight: 600;
  color: var(--color-text);
}

/* Settings page */
.admin-settings-card {
  margin-bottom: 16px;
}

.admin-field-wrap {
  display: flex;
  align-items: center;
  gap: 8px;
  width: 100%;
}

.admin-field-wrap > .el-input,
.admin-field-wrap > .el-input-number,
.admin-field-wrap > .el-select {
  flex: 1;
}

.admin-tip-icon {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 18px;
  height: 18px;
  border-radius: 50%;
  background: var(--color-text-muted);
  color: #fff;
  font-size: 11px;
  font-weight: 700;
  cursor: help;
  flex-shrink: 0;
  transition: background-color 0.15s ease;
}

.admin-tip-icon:hover {
  background: var(--color-text-secondary);
}

.admin-save-bar {
  padding: 16px 0;
  text-align: center;
}

/* Media grid */
.admin-media-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 16px;
}

@media (max-width: 1024px) {
  .admin-media-grid { grid-template-columns: repeat(3, 1fr); }
}
@media (max-width: 767px) {
  .admin-media-grid { grid-template-columns: repeat(2, 1fr); }
}

.admin-media-card {
  background: var(--color-bg-surface);
  border-radius: var(--radius-md);
  overflow: hidden;
  border: 1px solid var(--color-border-light);
  transition: box-shadow 0.2s ease;
  cursor: pointer;
}

.admin-media-card:hover {
  box-shadow: var(--shadow-md);
}

.admin-media-preview {
  width: 100%;
  height: 160px;
  background: var(--color-bg-app);
  display: flex;
  align-items: center;
  justify-content: center;
  overflow: hidden;
}

.admin-media-preview img {
  max-width: 100%;
  max-height: 100%;
  object-fit: contain;
}

.admin-media-info {
  padding: 10px 12px;
}

.admin-media-filename {
  font-size: 13px;
  color: var(--color-text);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  margin-bottom: 4px;
}

.admin-media-url {
  font-size: 11px;
  color: var(--color-text-muted);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  margin-bottom: 4px;
  font-family: var(--font-mono);
}

.admin-media-size {
  font-size: 12px;
  color: var(--color-text-secondary);
  margin-bottom: 8px;
}

/* Table overrides for borderless style */
.admin-table-wrap .el-table {
  --el-table-border-color: var(--color-border-light);
  --el-table-header-bg-color: var(--color-bg-app);
  --el-table-row-hover-bg-color: var(--color-bg-app);
  --el-table-header-text-color: var(--color-text-secondary);
  --el-table-text-color: var(--color-text);
}

.admin-table-wrap .el-table th.el-table__cell {
  font-size: 12px;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  border-bottom: 2px solid var(--color-border);
  height: 44px;
}

.admin-table-wrap .el-table td.el-table__cell {
  height: 52px;
  border-bottom: 1px solid var(--color-border-light);
}

/* Action button group in table rows */
.table-actions {
  display: flex;
  align-items: center;
  gap: 4px;
}

.table-actions .action-btn {
  padding: 4px 8px;
  font-size: 13px;
  color: var(--color-text-secondary);
  background: none;
  border: none;
  cursor: pointer;
  border-radius: var(--radius-sm);
  transition: all 0.15s ease;
}

.table-actions .action-btn:hover {
  color: var(--color-primary);
  background: var(--color-bg-app);
}

.table-actions .action-btn.danger:hover {
  color: var(--color-danger);
  background: var(--color-danger-soft);
}

/* Status pill tags */
.status-pill {
  display: inline-flex;
  align-items: center;
  padding: 2px 10px;
  border-radius: 100px;
  font-size: 12px;
  font-weight: 500;
}

.status-pill.published,
.status-pill.success {
  background: var(--color-success-soft);
  color: var(--color-success);
}

.status-pill.draft,
.status-pill.info {
  background: var(--color-info-soft);
  color: var(--color-info);
}

.status-pill.warning {
  background: var(--color-warning-soft);
  color: var(--color-warning);
}

.status-pill.danger {
  background: var(--color-danger-soft);
  color: var(--color-danger);
}

/* Homepage config card */
.config-card {
  margin-bottom: 0;
}

.config-card .el-card__header {
  padding: 16px 20px;
  border-bottom: 1px solid var(--color-border-light);
}

.config-card .el-card__body {
  padding: 20px;
}

.config-list {
  display: flex;
  flex-direction: column;
}

.config-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 0;
  border-bottom: 1px solid var(--color-border-light);
}

.config-item:last-child {
  border-bottom: none;
}

.config-item-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 2px;
  min-width: 0;
}

.config-item-info strong {
  font-size: 14px;
}

.config-item-desc {
  font-size: 12px;
  color: var(--color-text-muted);
}

.config-item-name {
  font-size: 14px;
  flex: 1;
}

.config-item-actions {
  display: flex;
  gap: 4px;
  flex-shrink: 0;
}

.config-list-actions {
  padding-top: 12px;
}

.card-footer {
  text-align: center;
  padding-top: 16px;
  border-top: 1px solid var(--color-border-light);
  margin-top: 16px;
}

/* Slide thumbnail */
.slide-thumb {
  width: 120px;
  height: 68px;
  object-fit: cover;
  border-radius: var(--radius-sm);
  flex-shrink: 0;
  border: 1px solid var(--color-border-light);
}

.slide-label {
  font-size: 14px;
  color: var(--color-text-muted);
  flex: 1;
}

/* Advantage icon preview */
.adv-icon-preview {
  width: 36px;
  height: 36px;
  background: var(--color-bg-app);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  border: 1px solid var(--color-border-light);
}

.adv-icon-svg {
  width: 18px;
  height: 18px;
  color: var(--color-accent);
}

.adv-icon-emoji {
  font-size: 18px;
}

/* Navigation tree */
.tree-node {
  display: flex;
  align-items: center;
  gap: 12px;
  flex: 1;
  padding: 4px 0;
}

.tree-node-label {
  font-weight: 500;
  min-width: 120px;
}

.tree-node-link {
  color: var(--color-text-muted);
  font-size: 13px;
  font-family: var(--font-mono);
  flex: 1;
}

.tree-node-actions {
  display: flex;
  gap: 2px;
  margin-left: auto;
}

/* Filter bar */
.admin-filter-bar {
  margin-bottom: 16px;
}

/* Section form */
.section-form {
  padding-bottom: 8px;
  margin-bottom: 8px;
  border-bottom: 1px solid var(--color-border-light);
}

/* Add project select */
.add-project-select {
  margin-top: 8px;
  width: 100%;
}

/* Link preview */
.link-preview {
  margin-top: 6px;
  font-size: 13px;
  color: var(--color-text-secondary);
}

.link-preview code {
  background: var(--color-bg-app);
  padding: 2px 6px;
  border-radius: 3px;
  font-family: var(--font-mono);
}

/* Monospace textarea */
.monospace-input :deep(textarea) {
  font-family: var(--font-mono);
  font-size: 13px;
}

/* Array input (settings) */
.array-input {
  flex: 1;
}

.array-row {
  display: flex;
  gap: 8px;
  margin-bottom: 8px;
  align-items: center;
}

/* Element Plus theme overrides */
:root {
  --el-color-primary: var(--color-primary);
  --el-color-primary-light-3: #3b4258;
  --el-color-primary-light-5: #6b7280;
  --el-color-primary-light-7: #9ca3af;
  --el-color-primary-light-8: #cbd5e1;
  --el-color-primary-light-9: #e2e8f0;
  --el-color-success: var(--color-success);
  --el-color-warning: var(--color-warning);
  --el-color-danger: var(--color-danger);
  --el-color-info: var(--color-info);
  --el-border-radius-base: var(--radius-sm);
  --el-font-family: var(--font-sans);
}
```

---

### Task 4: Rewrite Admin Layout (Sidebar + Topbar)

**Files:**
- Modify: `frontend/layouts/admin.vue`

- [ ] **Step 1: Replace template with new sidebar + topbar**

Replace the `<template>` block:

```vue
<template>
  <div class="admin-layout">
    <div
      class="sidebar-overlay"
      :class="{ 'overlay-visible': mobileOpen }"
      @click="mobileOpen = false"
    ></div>

    <aside
      class="admin-sidebar"
      :class="{ collapsed: sidebarCollapsed, 'mobile-open': mobileOpen }"
    >
      <!-- Logo -->
      <div class="sidebar-header">
        <NuxtLink to="/admin" class="sidebar-logo" @click="closeMobile">
          <div class="logo-icon">
            <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <path d="m3 9 9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z"/>
              <polyline points="9 22 9 12 15 12 15 22"/>
            </svg>
          </div>
          <span class="logo-text">MyGo 管理后台</span>
        </NuxtLink>
      </div>

      <!-- Navigation -->
      <nav class="sidebar-nav">
        <NuxtLink
          v-for="item in navItems"
          :key="item.to"
          :to="item.to"
          class="nav-item"
          :class="{ active: isActive(item.to) }"
          @click="closeMobile"
        >
          <span class="nav-icon" v-html="item.icon"></span>
          <span class="nav-label">{{ item.label }}</span>
        </NuxtLink>
      </nav>

      <!-- User footer -->
      <div class="sidebar-footer">
        <div class="user-info">
          <div class="user-avatar">{{ userInitial }}</div>
          <div class="user-meta">
            <div class="user-name">{{ userName }}</div>
            <div class="user-role">{{ userRoleLabel }}</div>
          </div>
        </div>
        <button class="sidebar-collapse-btn" @click="toggleSidebar" :title="sidebarCollapsed ? '展开侧边栏' : '折叠侧边栏'">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <polyline v-if="sidebarCollapsed" points="9 18 15 12 9 6"/>
            <polyline v-else points="15 18 9 12 15 6"/>
          </svg>
        </button>
      </div>
    </aside>

    <!-- Main area -->
    <div class="admin-main">
      <header class="admin-topbar">
        <button class="menu-toggle" @click="toggleSidebar">
          <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <line x1="4" x2="20" y1="6" y2="6"/>
            <line x1="4" x2="20" y1="12" y2="12"/>
            <line x1="4" x2="20" y1="18" y2="18"/>
          </svg>
        </button>

        <div class="topbar-right">
          <NuxtLink to="/" target="_blank" class="topbar-link">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M18 13v6a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V8a2 2 0 0 1 2-2h6"/><polyline points="15 3 21 3 21 9"/><line x1="10" x2="21" y1="14" y2="3"/></svg>
            <span>访问网站</span>
          </NuxtLink>
          <button class="topbar-link logout-btn" @click="handleLogout">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"/><polyline points="16 17 21 12 16 7"/><line x1="21" x2="9" y1="12" y2="12"/></svg>
            <span>退出登录</span>
          </button>
        </div>
      </header>

      <main class="admin-content">
        <slot />
      </main>
    </div>
  </div>
</template>
```

- [ ] **Step 2: Replace script block**

Replace the `<script setup>` block:

```typescript
<script setup lang="ts">
import { getIconSvg } from '~/composables/lucideIcons';

const route = useRoute();

const sidebarCollapsed = ref(false);
const mobileOpen = ref(false);

interface NavItem {
  to: string;
  label: string;
  icon: string;
}

const navItems: NavItem[] = [
  { to: '/admin', label: '控制台', icon: getIconSvg('bar-chart', 20) },
  { to: '/admin/projects', label: '项目管理', icon: getIconSvg('folder', 20) },
  { to: '/admin/pages', label: '页面管理', icon: getIconSvg('file-text', 20) },
  { to: '/admin/homepage', label: '首页配置', icon: getIconSvg('home', 20) },
  { to: '/admin/navigation', label: '导航管理', icon: getIconSvg('compass', 20) },
  { to: '/admin/faqs', label: 'FAQ 管理', icon: getIconSvg('help-circle', 20) },
  { to: '/admin/cases', label: '案例管理', icon: getIconSvg('users', 20) },
  { to: '/admin/leads', label: '咨询管理', icon: getIconSvg('message-circle', 20) },
  { to: '/admin/users', label: '用户管理', icon: getIconSvg('shield', 20) },
  { to: '/admin/media', label: '媒体库', icon: getIconSvg('download', 20) },
  { to: '/admin/settings', label: '网站设置', icon: getIconSvg('settings', 20) },
];

const isActive = (to: string) => {
  if (to === '/admin') return route.path === '/admin';
  return route.path.startsWith(to);
};

function toggleSidebar() {
  if (window.innerWidth < 768) {
    mobileOpen.value = !mobileOpen.value;
  } else {
    sidebarCollapsed.value = !sidebarCollapsed.value;
  }
}

function closeMobile() {
  mobileOpen.value = false;
}

const handleLogout = () => {
  const { logout } = useAuth();
  logout();
};

const { user } = useAuth();

const userName = computed(() => (user.value as any)?.display_name || (user.value as any)?.username || '管理员');
const userRoleLabel = computed(() => {
  const role = (user.value as any)?.role;
  if (role === 'admin') return '管理员';
  if (role === 'editor') return '编辑者';
  return '只读用户';
});
const userInitial = computed(() => userName.value.charAt(0).toUpperCase());
</script>
```

- [ ] **Step 3: Replace style block**

Replace the `<style scoped>` block:

```css
<style scoped>
.admin-layout {
  display: flex;
  min-height: 100vh;
  background: var(--color-bg-app);
}

/* Sidebar overlay (mobile) */
.sidebar-overlay {
  display: none;
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.4);
  z-index: 150;
}

/* Sidebar */
.admin-sidebar {
  width: var(--sidebar-width);
  background: var(--color-bg-sidebar);
  color: var(--color-text-sidebar);
  display: flex;
  flex-direction: column;
  flex-shrink: 0;
  transition: width 0.25s ease;
  overflow: hidden;
  position: sticky;
  top: 0;
  height: 100vh;
}

.sidebar-header {
  padding: 20px 16px;
  border-bottom: 1px solid var(--color-border-sidebar);
}

.sidebar-logo {
  display: flex;
  align-items: center;
  gap: 10px;
  color: #fff;
  text-decoration: none;
}

.logo-icon {
  width: 32px;
  height: 32px;
  background: var(--color-accent);
  border-radius: var(--radius-sm);
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  flex-shrink: 0;
}

.logo-text {
  font-size: 15px;
  font-weight: 700;
  white-space: nowrap;
  transition: opacity 0.2s ease;
}

/* Nav */
.sidebar-nav {
  flex: 1;
  padding: 12px 8px;
  overflow-y: auto;
  overflow-x: hidden;
}

.nav-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px 12px;
  margin-bottom: 2px;
  color: var(--color-text-sidebar);
  border-radius: var(--radius-sm);
  font-size: 14px;
  white-space: nowrap;
  border-left: 3px solid transparent;
  transition: all 0.15s ease;
}

.nav-item:hover {
  color: #fff;
  background: var(--color-bg-sidebar-hover);
}

.nav-item.active {
  color: #fff;
  background: var(--color-bg-sidebar-active);
  border-left-color: var(--color-accent);
}

.nav-icon {
  width: 20px;
  height: 20px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.nav-icon :deep(svg) {
  width: 20px;
  height: 20px;
}

.nav-label {
  transition: opacity 0.2s ease;
}

/* Sidebar footer */
.sidebar-footer {
  padding: 12px;
  border-top: 1px solid var(--color-border-sidebar);
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 10px;
  min-width: 0;
}

.user-avatar {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  background: #334155;
  color: var(--color-accent);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 14px;
  font-weight: 600;
  flex-shrink: 0;
}

.user-meta {
  min-width: 0;
}

.user-name {
  font-size: 13px;
  font-weight: 500;
  color: #fff;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.user-role {
  font-size: 11px;
  color: var(--color-text-sidebar);
}

.sidebar-collapse-btn {
  width: 28px;
  height: 28px;
  border: none;
  background: none;
  color: var(--color-text-sidebar);
  cursor: pointer;
  border-radius: var(--radius-sm);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  transition: all 0.15s ease;
}

.sidebar-collapse-btn:hover {
  color: #fff;
  background: var(--color-bg-sidebar-hover);
}

/* Collapsed sidebar */
.admin-sidebar.collapsed {
  width: var(--sidebar-collapsed-width);
}

.admin-sidebar.collapsed .logo-text,
.admin-sidebar.collapsed .nav-label,
.admin-sidebar.collapsed .user-meta,
.admin-sidebar.collapsed .user-avatar {
  opacity: 0;
  width: 0;
  overflow: hidden;
  padding: 0;
  margin: 0;
}

.admin-sidebar.collapsed .sidebar-logo {
  justify-content: center;
}

.admin-sidebar.collapsed .nav-item {
  justify-content: center;
  padding: 10px;
  border-left: none;
}

.admin-sidebar.collapsed .nav-item.active {
  border-left: none;
  border-radius: var(--radius-sm);
}

.admin-sidebar.collapsed .sidebar-footer {
  justify-content: center;
}

.admin-sidebar.collapsed .sidebar-header {
  padding: 16px 8px;
  text-align: center;
}

/* Main area */
.admin-main {
  flex: 1;
  display: flex;
  flex-direction: column;
  min-width: 0;
}

/* Topbar */
.admin-topbar {
  height: var(--topbar-height);
  background: var(--color-bg-surface);
  border-bottom: 1px solid var(--color-border);
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 20px;
  flex-shrink: 0;
  position: sticky;
  top: 0;
  z-index: 50;
}

.menu-toggle {
  background: none;
  border: none;
  cursor: pointer;
  padding: 6px;
  color: var(--color-text-secondary);
  border-radius: var(--radius-sm);
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.15s ease;
}

.menu-toggle:hover {
  color: var(--color-text);
  background: var(--color-bg-app);
}

.topbar-right {
  display: flex;
  align-items: center;
  gap: 4px;
}

.topbar-link {
  font-size: 13px;
  color: var(--color-text-secondary);
  background: none;
  border: none;
  cursor: pointer;
  padding: 8px 12px;
  border-radius: var(--radius-sm);
  display: flex;
  align-items: center;
  gap: 6px;
  transition: all 0.15s ease;
}

.topbar-link:hover {
  color: var(--color-text);
  background: var(--color-bg-app);
}

.logout-btn:hover {
  color: var(--color-danger);
}

/* Content */
.admin-content {
  flex: 1;
  padding: 24px;
  overflow-y: auto;
}

@media (max-width: 767px) {
  .admin-content {
    padding: 16px;
  }
}

/* Mobile */
@media (max-width: 767px) {
  .sidebar-overlay {
    display: block;
    opacity: 0;
    pointer-events: none;
    transition: opacity 0.25s ease;
  }

  .sidebar-overlay.overlay-visible {
    opacity: 1;
    pointer-events: auto;
  }

  .admin-sidebar {
    position: fixed;
    top: 0;
    left: 0;
    height: 100vh;
    z-index: 200;
    transform: translateX(-100%);
    transition: transform 0.25s ease;
  }

  .admin-sidebar.mobile-open {
    transform: translateX(0);
  }

  .admin-sidebar.collapsed {
    width: var(--sidebar-width);
  }
}
</style>
```

- [ ] **Step 4: Verify the layout compiles**

Run: `cd frontend && npx nuxi prepare`
Expected: No errors.

---

## Phase 2: List Page Templates

### Task 5: Rewrite Projects Page (projects.vue)

**Files:**
- Modify: `frontend/pages/admin/projects.vue`

- [ ] **Step 1: Write the complete new projects.vue**

Replace entire file with:

```vue
<template>
  <div>
    <div class="admin-page-header">
      <h2 class="admin-page-title">项目管理</h2>
      <el-button type="primary" @click="openCreate">新建项目</el-button>
    </div>

    <!-- Toolbar -->
    <div class="admin-toolbar">
      <el-input
        v-model="searchQuery"
        placeholder="搜索项目名称..."
        :prefix-icon="Search"
        clearable
        class="admin-search-input"
        @input="onSearch"
      />
      <el-select
        v-model="statusFilter"
        placeholder="状态筛选"
        clearable
        class="admin-filter-select"
        @change="loadList"
      >
        <el-option label="全部" value="" />
        <el-option label="已发布" value="1" />
        <el-option label="草稿" value="0" />
      </el-select>
    </div>

    <!-- Table -->
    <div class="admin-table-wrap">
      <el-table :data="list" v-loading="loading">
        <el-table-column prop="name" label="项目名称" min-width="180">
          <template #default="{ row }">
            <div>
              <div class="row-title">{{ row.name }}</div>
              <div class="row-meta">{{ row.country }}</div>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="investment_amount" label="投资金额" width="130" />
        <el-table-column prop="processing_period" label="办理周期" width="120" />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <span :class="['status-pill', row.status === 1 ? 'published' : 'draft']">
              {{ row.status === 1 ? '已发布' : '草稿' }}
            </span>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="140" fixed="right">
          <template #default="{ row }">
            <div class="table-actions">
              <button class="action-btn" @click="openEdit(row)">编辑</button>
              <el-popconfirm
                title="确定删除该项目？"
                confirm-button-text="删除"
                cancel-button-text="取消"
                @confirm="handleDelete(row.id)"
              >
                <template #reference>
                  <button class="action-btn danger">删除</button>
                </template>
              </el-popconfirm>
            </div>
          </template>
        </el-table-column>
      </el-table>
    </div>

    <!-- Empty state -->
    <div v-if="!loading && list.length === 0" class="admin-empty-state">
      <div class="empty-icon" v-html="getIconSvg('folder', 48)"></div>
      <div class="empty-title">暂无项目</div>
      <div class="empty-desc">点击上方按钮创建第一个项目</div>
      <el-button type="primary" @click="openCreate">新建项目</el-button>
    </div>

    <div class="admin-pagination-wrap" v-if="total > pageSize">
      <el-pagination
        v-model:current-page="page"
        :page-size="pageSize"
        :total="total"
        layout="total, prev, pager, next"
        @current-change="loadList"
      />
    </div>

    <!-- Drawer -->
    <el-drawer
      v-model="drawerVisible"
      :title="editingId ? '编辑项目' : '新建项目'"
      size="560px"
      destroy-on-close
    >
      <el-tabs v-model="activeTab">
        <el-tab-pane label="基本信息" name="basic">
          <el-form ref="formRef" :model="form" :rules="rules" label-position="top">
            <el-row :gutter="12">
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
            <el-row :gutter="12">
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
            <el-row :gutter="12">
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
            <el-row :gutter="12">
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
              <ImageInput v-model="form.cover_image" placeholder="图片 URL 或上传" />
            </el-form-item>
            <el-form-item label="Hero 描述" prop="hero_desc">
              <el-input v-model="form.hero_desc" type="textarea" :rows="2" />
            </el-form-item>
            <el-form-item label="Hero 渐变" prop="hero_gradient">
              <el-input v-model="form.hero_gradient" />
            </el-form-item>
            <el-row :gutter="12">
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

        <el-tab-pane label="申请条件" name="requirements">
          <div style="margin-bottom: 12px">
            <el-button type="primary" size="small" @click="openSubDialog('requirement')">添加条件</el-button>
          </div>
          <el-table :data="subData.requirements" border size="small">
            <el-table-column prop="label" label="条件描述" min-width="180" />
            <el-table-column label="必需" width="70">
              <template #default="{ row: r }">
                <span :class="['status-pill', r.is_required ? 'published' : 'draft']">
                  {{ r.is_required ? '必需' : '可选' }}
                </span>
              </template>
            </el-table-column>
            <el-table-column prop="sort_order" label="排序" width="60" />
            <el-table-column label="操作" width="120">
              <template #default="{ row: r }">
                <div class="table-actions">
                  <button class="action-btn" @click="openSubDialog('requirement', r)">编辑</button>
                  <el-popconfirm title="确定删除？" @confirm="deleteSubItem('requirement', r.id)">
                    <template #reference>
                      <button class="action-btn danger">删除</button>
                    </template>
                  </el-popconfirm>
                </div>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>

        <el-tab-pane label="费用明细" name="costItems">
          <div style="margin-bottom: 12px">
            <el-button type="primary" size="small" @click="openSubDialog('costItem')">添加费用</el-button>
          </div>
          <el-table :data="subData.costItems" border size="small">
            <el-table-column prop="name" label="费用名称" min-width="140" />
            <el-table-column prop="amount" label="金额" width="100" />
            <el-table-column prop="note" label="说明" min-width="160" />
            <el-table-column prop="sort_order" label="排序" width="60" />
            <el-table-column label="操作" width="120">
              <template #default="{ row: r }">
                <div class="table-actions">
                  <button class="action-btn" @click="openSubDialog('costItem', r)">编辑</button>
                  <el-popconfirm title="确定删除？" @confirm="deleteSubItem('costItem', r.id)">
                    <template #reference>
                      <button class="action-btn danger">删除</button>
                    </template>
                  </el-popconfirm>
                </div>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>

        <el-tab-pane label="申请流程" name="timelinePhases">
          <div style="margin-bottom: 12px">
            <el-button type="primary" size="small" @click="openSubDialog('timelinePhase')">添加步骤</el-button>
          </div>
          <el-table :data="subData.timelinePhases" border size="small">
            <el-table-column prop="phase_number" label="步骤号" width="70" />
            <el-table-column prop="title" label="标题" min-width="130" />
            <el-table-column prop="description" label="描述" min-width="160" />
            <el-table-column prop="duration" label="周期" width="90" />
            <el-table-column prop="sort_order" label="排序" width="60" />
            <el-table-column label="操作" width="120">
              <template #default="{ row: r }">
                <div class="table-actions">
                  <button class="action-btn" @click="openSubDialog('timelinePhase', r)">编辑</button>
                  <el-popconfirm title="确定删除？" @confirm="deleteSubItem('timelinePhase', r.id)">
                    <template #reference>
                      <button class="action-btn danger">删除</button>
                    </template>
                  </el-popconfirm>
                </div>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>
      </el-tabs>

      <template #footer>
        <el-button @click="drawerVisible = false">取消</el-button>
        <el-button type="primary" :loading="saving" @click="handleSave">保存</el-button>
      </template>
    </el-drawer>

    <!-- Sub-entity dialog -->
    <el-dialog
      v-model="subDialogVisible"
      :title="subDialogTitle"
      width="500px"
      destroy-on-close
    >
      <el-form ref="subFormRef" :model="subForm" label-position="top">
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
  </div>
</template>
```

- [ ] **Step 2: Write the script block**

The script logic remains identical to the existing `projects.vue` script block (lines 298-571 of the current file). Only two additions at the top of `<script setup>`:

```typescript
<script setup lang="ts">
import { Search } from '@element-plus/icons-vue';
import { ElMessage, type FormInstance, type FormRules } from 'element-plus';
import ImageInput from '~/components/admin/ImageInput.vue';
import { getIconSvg } from '~/composables/lucideIcons';

// ... (all existing logic unchanged)

// Replace dialogVisible → drawerVisible
const drawerVisible = ref(false);

// Add search/filter state
const searchQuery = ref('');
const statusFilter = ref('');

let searchTimer: ReturnType<typeof setTimeout>;
const onSearch = () => {
  clearTimeout(searchTimer);
  searchTimer = setTimeout(() => {
    page.value = 1;
    loadList();
  }, 300);
};

// Update loadList to use search/status
const loadList = async () => {
  loading.value = true;
  try {
    const api = useApi();
    let url = `/admin/projects?page=${page.value}&per_page=${pageSize.value}`;
    if (searchQuery.value) url += `&search=${encodeURIComponent(searchQuery.value)}`;
    if (statusFilter.value) url += `&status=${statusFilter.value}`;
    const data = await api<{ items: Project[]; total: number }>(url);
    list.value = data.items ?? [];
    total.value = data.total ?? 0;
  } catch {
    list.value = [];
    ElMessage.error('加载项目列表失败');
  } finally {
    loading.value = false;
  }
};

// Replace dialogVisible with drawerVisible everywhere
// (Use Edit to replace: dialogVisible → drawerVisible)
</script>
```

- [ ] **Step 3: Add scoped styles**

```css
<style scoped>
.row-title {
  font-size: 14px;
  font-weight: 500;
  color: var(--color-text);
}

.row-meta {
  font-size: 12px;
  color: var(--color-text-muted);
  margin-top: 2px;
}
</style>
```

- [ ] **Step 4: Verify the page compiles**

Run: `cd frontend && npx nuxi typecheck 2>&1 | head -20`
Expected: No new type errors.

---

### Task 6: Rewrite Pages Page (pages.vue)

**Files:**
- Modify: `frontend/pages/admin/pages.vue`

This follows the exact same template as Task 5 (projects.vue). Key differences:

- [ ] **Step 1: Apply the unified list template**

Replace template with search + filter + borderless table + drawer, matching the projects.vue pattern but with pages-specific fields:

```vue
<template>
  <div>
    <div class="admin-page-header">
      <h2 class="admin-page-title">页面管理</h2>
      <el-button type="primary" @click="openCreate">新建页面</el-button>
    </div>

    <div class="admin-toolbar">
      <el-input
        v-model="searchQuery"
        placeholder="搜索页面标题..."
        :prefix-icon="Search"
        clearable
        class="admin-search-input"
        @input="onSearch"
      />
      <el-select v-model="statusFilter" placeholder="状态筛选" clearable class="admin-filter-select" @change="loadList">
        <el-option label="全部" value="" />
        <el-option label="已发布" value="published" />
        <el-option label="草稿" value="draft" />
      </el-select>
    </div>

    <div class="admin-table-wrap">
      <el-table :data="list" v-loading="loading">
        <el-table-column prop="title" label="标题" min-width="180">
          <template #default="{ row }">
            <div>
              <div class="row-title">{{ row.title }}</div>
              <div class="row-meta">/{{ row.slug }}</div>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="template" label="模板" width="100" />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <span :class="['status-pill', row.status === 'published' ? 'published' : 'draft']">
              {{ row.status === 'published' ? '已发布' : '草稿' }}
            </span>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="140" fixed="right">
          <template #default="{ row }">
            <div class="table-actions">
              <button class="action-btn" @click="openEdit(row)">编辑</button>
              <el-popconfirm title="确定删除该页面？" confirm-button-text="删除" cancel-button-text="取消" @confirm="handleDelete(row.id)">
                <template #reference>
                  <button class="action-btn danger">删除</button>
                </template>
              </el-popconfirm>
            </div>
          </template>
        </el-table-column>
      </el-table>
    </div>

    <div v-if="!loading && list.length === 0" class="admin-empty-state">
      <div class="empty-icon" v-html="getIconSvg('file-text', 48)"></div>
      <div class="empty-title">暂无页面</div>
      <div class="empty-desc">点击上方按钮创建第一个页面</div>
      <el-button type="primary" @click="openCreate">新建页面</el-button>
    </div>

    <div class="admin-pagination-wrap" v-if="total > pageSize">
      <el-pagination v-model:current-page="page" :page-size="pageSize" :total="total" layout="total, prev, pager, next" @current-change="loadList" />
    </div>

    <el-drawer v-model="drawerVisible" :title="editingId ? '编辑页面' : '新建页面'" size="560px" destroy-on-close>
      <el-form ref="formRef" :model="form" :rules="rules" label-position="top">
        <el-form-item label="标题" prop="title">
          <el-input v-model="form.title" />
        </el-form-item>
        <el-form-item label="标识(slug)" prop="slug">
          <el-input v-model="form.slug" />
        </el-form-item>
        <el-form-item label="内容" prop="content">
          <el-input v-model="form.content" type="textarea" :rows="8" />
        </el-form-item>
        <el-row :gutter="12">
          <el-col :span="12">
            <el-form-item label="SEO 标题" prop="meta_title">
              <el-input v-model="form.meta_title" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="SEO 描述" prop="meta_description">
              <el-input v-model="form.meta_description" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="12">
          <el-col :span="12">
            <el-form-item label="模板" prop="template">
              <el-select v-model="form.template">
                <el-option label="默认" value="default" />
                <el-option label="全宽" value="fullwidth" />
                <el-option label="落地页" value="landing" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="状态" prop="status">
              <el-select v-model="form.status">
                <el-option label="草稿" value="draft" />
                <el-option label="已发布" value="published" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>
      <template #footer>
        <el-button @click="drawerVisible = false">取消</el-button>
        <el-button type="primary" :loading="saving" @click="handleSave">保存</el-button>
      </template>
    </el-drawer>
  </div>
</template>
```

- [ ] **Step 2: Add script updates**

Add imports and search logic, replace `dialogVisible` with `drawerVisible` in the existing script block.

---

### Task 7: Rewrite FAQs Page (faqs.vue)

**Files:**
- Modify: `frontend/pages/admin/faqs.vue`

Same unified list+drawer template. Table columns: 问题 / 所属项目 / 全局 / 排序 / 操作. Drawer form fields: 问题, 回答, 所属项目(select), 全局(switch), 排序.

---

### Task 8: Rewrite Cases Page (cases.vue)

**Files:**
- Modify: `frontend/pages/admin/cases.vue`

Same template. Table columns: 姓名 / 来源国家 / 项目 / 排序 / 操作. Drawer form: 姓名, 来源国家, 所属项目(select), 描述(textarea), 排序.

---

### Task 9: Rewrite Users Page (users.vue)

**Files:**
- Modify: `frontend/pages/admin/users.vue`

Same template. Table columns: 用户名 / 显示名称 / 角色(pill) / 状态(pill) / 操作. Keep the existing role edit dialog and toggle-status logic; wrap in drawer for create.

---

### Task 10: Rewrite Leads Page (leads.vue)

**Files:**
- Modify: `frontend/pages/admin/leads.vue`

Same template. Table columns: 姓名 / 电话 / 邮箱 / 感兴趣项目 / 状态(pill) / 创建时间. Click row opens drawer with detail + status update form (reuse existing detail dialog logic but in drawer).

---

## Phase 3: Special Pages

### Task 11: Rewrite Dashboard (index.vue)

**Files:**
- Modify: `frontend/pages/admin/index.vue`

- [ ] **Step 1: Replace template with new dashboard design**

```vue
<template>
  <div>
    <h2 class="admin-page-title" style="margin-bottom: 24px">控制台</h2>

    <!-- Stats -->
    <div class="admin-stat-grid">
      <div class="admin-stat-card">
        <div class="stat-icon" style="background: var(--color-info-soft); color: var(--color-info);" v-html="getIconSvg('folder', 18)"></div>
        <div class="admin-stat-label">项目总数</div>
        <div class="admin-stat-value">{{ stats.totalProjects ?? '-' }}</div>
        <div class="admin-stat-trend neutral">—</div>
      </div>
      <div class="admin-stat-card">
        <div class="stat-icon" style="background: var(--color-warning-soft); color: var(--color-warning);" v-html="getIconSvg('message-circle', 18)"></div>
        <div class="admin-stat-label">咨询总数</div>
        <div class="admin-stat-value">{{ stats.totalLeads ?? '-' }}</div>
        <div class="admin-stat-trend neutral">—</div>
      </div>
      <div class="admin-stat-card">
        <div class="stat-icon" style="background: var(--color-success-soft); color: var(--color-success);" v-html="getIconSvg('file-text', 18)"></div>
        <div class="admin-stat-label">页面总数</div>
        <div class="admin-stat-value">{{ stats.totalPages ?? '-' }}</div>
        <div class="admin-stat-trend neutral">—</div>
      </div>
      <div class="admin-stat-card">
        <div class="stat-icon" style="background: #ede9fe; color: #7c3aed;" v-html="getIconSvg('users', 18)"></div>
        <div class="admin-stat-label">案例总数</div>
        <div class="admin-stat-value">{{ stats.totalCases ?? '-' }}</div>
        <div class="admin-stat-trend neutral">—</div>
      </div>
    </div>

    <!-- Content row -->
    <div style="display: grid; grid-template-columns: 2fr 1fr; gap: 16px;">
      <!-- Recent leads -->
      <div class="admin-card">
        <h3 class="admin-section-title">最近咨询</h3>
        <div v-if="recentLeads.length === 0" style="text-align: center; padding: 40px 0; color: var(--color-text-muted); font-size: 14px;">
          暂无咨询记录
        </div>
        <div v-else>
          <div
            v-for="lead in recentLeads"
            :key="lead.id"
            style="display: flex; justify-content: space-between; align-items: center; padding: 12px 0; border-bottom: 1px solid var(--color-border-light);"
          >
            <div>
              <div style="font-size: 14px; font-weight: 500;">{{ lead.name }}</div>
              <div style="font-size: 12px; color: var(--color-text-secondary);">{{ lead.interested_project || '未指定项目' }}</div>
            </div>
            <span style="font-size: 12px; color: var(--color-text-muted);">{{ formatTime(lead.created_at) }}</span>
          </div>
        </div>
      </div>

      <!-- Quick links -->
      <div class="admin-card">
        <h3 class="admin-section-title">快捷操作</h3>
        <div style="display: flex; flex-direction: column; gap: 8px;">
          <NuxtLink to="/admin/projects" class="admin-quick-link">
            <span v-html="getIconSvg('plus', 16)"></span>
            <span>新建项目</span>
          </NuxtLink>
          <NuxtLink to="/admin/pages" class="admin-quick-link">
            <span v-html="getIconSvg('plus', 16)"></span>
            <span>新建页面</span>
          </NuxtLink>
          <NuxtLink to="/admin/faqs" class="admin-quick-link">
            <span v-html="getIconSvg('plus', 16)"></span>
            <span>新建 FAQ</span>
          </NuxtLink>
          <NuxtLink to="/admin/leads" class="admin-quick-link">
            <span v-html="getIconSvg('message-circle', 16)"></span>
            <span>查看咨询</span>
          </NuxtLink>
        </div>
      </div>
    </div>
  </div>
</template>
```

- [ ] **Step 2: Add script with data loading**

```typescript
<script setup lang="ts">
definePageMeta({ layout: 'admin', middleware: 'auth' });
import { getIconSvg } from '~/composables/lucideIcons';

interface DashboardStats {
  totalProjects: number;
  totalPages: number;
  totalLeads: number;
  totalCases: number;
}

interface RecentLead {
  id: string;
  name: string;
  interested_project: string;
  created_at: string;
}

const stats = reactive<DashboardStats>({
  totalProjects: 0,
  totalPages: 0,
  totalLeads: 0,
  totalCases: 0,
});

const recentLeads = ref<RecentLead[]>([]);

const loadStats = async () => {
  try {
    const api = useApi();
    const data = await api<DashboardStats>('/admin/dashboard/stats');
    Object.assign(stats, data);
  } catch {
    // silent
  }
};

const loadRecentLeads = async () => {
  try {
    const api = useApi();
    const data = await api<{ items: RecentLead[] }>('/admin/leads?page=1&per_page=5');
    recentLeads.value = data?.items ?? [];
  } catch {
    // silent
  }
};

const formatTime = (ts: string) => {
  if (!ts) return '';
  const d = new Date(ts);
  const now = new Date();
  const diff = now.getTime() - d.getTime();
  const mins = Math.floor(diff / 60000);
  if (mins < 60) return `${mins}分钟前`;
  const hours = Math.floor(mins / 60);
  if (hours < 24) return `${hours}小时前`;
  const days = Math.floor(hours / 24);
  return `${days}天前`;
};

onMounted(() => {
  loadStats();
  loadRecentLeads();
});
</script>
```

---

### Task 12: Rewrite Homepage Config (homepage.vue)

**Files:**
- Modify: `frontend/pages/admin/homepage.vue`

- [ ] **Step 1: Replace template with tab layout**

Change the main structure from stacked cards to tabs. Wrap the three card sections inside `<el-tabs>`. Keep all existing slide/advantage/showcase logic identical, only change the layout wrapper.

```vue
<template>
  <div class="homepage-config" v-loading="loading">
    <div class="admin-page-header">
      <h2 class="admin-page-title">首页配置</h2>
    </div>

    <el-tabs v-model="activeConfigTab" type="border-card">
      <el-tab-pane label="轮播管理" name="slides">
        <!-- ... existing hero slides card content ... -->
      </el-tab-pane>
      <el-tab-pane label="项目展示区" name="showcase">
        <!-- ... existing project showcase card content ... -->
      </el-tab-pane>
      <el-tab-pane label="优势管理" name="advantages">
        <!-- ... existing advantage items card content ... -->
      </el-tab-pane>
    </el-tabs>

    <!-- ... existing dialogs unchanged ... -->
  </div>
</template>
```

- [ ] **Step 2: Add activeConfigTab ref**

Add `const activeConfigTab = ref('slides');` to the script. Everything else stays the same.

---

### Task 13: Polish Navigation Page (navigation.vue)

**Files:**
- Modify: `frontend/pages/admin/navigation.vue`

- [ ] **Step 1: Replace action buttons with icon buttons**

In the tree node template, replace the three text buttons (`添加子级`, `编辑`, `删除`) with icon-only small buttons:

```vue
<span v-if="!isViewer" class="tree-node-actions">
  <el-button size="small" :icon="Plus" circle @click.stop="openCreate(data.id)" title="添加子级" />
  <el-button size="small" :icon="Edit" circle @click.stop="openEdit(data)" title="编辑" />
  <el-popconfirm title="确定删除该导航项？" @confirm="handleDelete(data.id)">
    <template #reference>
      <el-button size="small" :icon="Delete" circle type="danger" title="删除" />
    </template>
  </el-popconfirm>
</span>
```

- [ ] **Step 2: Add icon imports**

Add `import { Plus, Edit, Delete } from '@element-plus/icons-vue';` to the script block.

- [ ] **Step 3: Update the link preview to use drawer instead of dialog**

Replace `el-dialog` with `el-drawer` (size="500px"), same form content.

---

### Task 14: Polish Media Page (media.vue)

**Files:**
- Modify: `frontend/pages/admin/media.vue`

- [ ] **Step 1: Add image preview on click**

Wrap media card images with a click handler that opens a preview:

```vue
<div class="admin-media-card" @click="previewImage = item.url">
  <!-- existing content -->
</div>
```

Add a preview overlay component:

```vue
<div v-if="previewImage" class="image-preview-overlay" @click="previewImage = ''">
  <img :src="previewImage" class="image-preview-img" @click.stop />
  <button class="image-preview-close" @click="previewImage = ''">
    <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="white" stroke-width="2"><path d="M18 6 6 18"/><path d="m6 6 12 12"/></svg>
  </button>
</div>
```

Add `const previewImage = ref('');` and the preview overlay CSS:

```css
.image-preview-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.85);
  z-index: 1000;
  display: flex;
  align-items: center;
  justify-content: center;
}
.image-preview-img {
  max-width: 90vw;
  max-height: 90vh;
  object-fit: contain;
}
.image-preview-close {
  position: absolute;
  top: 20px;
  right: 20px;
  background: none;
  border: none;
  cursor: pointer;
  padding: 8px;
  border-radius: 50%;
  transition: background 0.15s;
}
.image-preview-close:hover {
  background: rgba(255,255,255,0.1);
}
```

- [ ] **Step 2: Improve upload button styling**

Add icon to upload button:

```vue
<el-button type="primary" :icon="Upload">上传图片</el-button>
```

Import: `import { Upload } from '@element-plus/icons-vue';`

---

### Task 15: Polish Settings Page (settings.vue)

**Files:**
- Modify: `frontend/pages/admin/settings.vue`

- [ ] **Step 1: Adjust form label width**

In both `<el-form>` instances on the page, set `label-width="160px"`.

- [ ] **Step 2: Adjust card spacing**

Ensure all `el-card` instances have `class="admin-settings-card"` (they already do — verify CSS handles the 16px margin-bottom from admin.css).

---

### Task 16: Upgrade Login Page (login.vue)

**Files:**
- Modify: `frontend/pages/admin/login.vue`

- [ ] **Step 1: Add logo area to login card**

Add before the title:

```vue
<div class="login-logo">
  <svg width="40" height="40" viewBox="0 0 24 24" fill="none" stroke="var(--color-accent)" stroke-width="2">
    <path d="m3 9 9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z"/>
    <polyline points="9 22 9 12 15 12 15 22"/>
  </svg>
</div>
```

- [ ] **Step 2: Update style to use new color tokens**

Replace the style block:

```css
<style scoped>
.login-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #0f172a 0%, #1e293b 100%);
}

.login-card {
  width: 400px;
  padding: 40px;
  background: var(--color-bg-surface);
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-lg);
}

.login-logo {
  text-align: center;
  margin-bottom: 8px;
}

.login-title {
  font-size: 22px;
  font-weight: 700;
  text-align: center;
  margin-bottom: 32px;
  color: var(--color-text);
}

.login-btn {
  width: 100%;
  height: 42px;
  font-size: 15px;
  border-radius: var(--radius-md);
}

.error-msg {
  color: var(--color-danger);
  text-align: center;
  font-size: 14px;
  margin-top: 12px;
}
</style>
```

---

### Task 17: Final Verification

- [ ] **Step 1: Type check**

Run: `cd frontend && npx nuxi typecheck`

- [ ] **Step 2: Start dev server and visually verify all pages**

Run: `cd frontend && npm run dev`

Check each page at:
- `http://localhost:3000/admin` — Dashboard
- `http://localhost:3000/admin/projects` — Projects list + drawer
- `http://localhost:3000/admin/pages` — Pages list + drawer
- `http://localhost:3000/admin/faqs` — FAQs list + drawer
- `http://localhost:3000/admin/cases` — Cases list + drawer
- `http://localhost:3000/admin/leads` — Leads list + drawer
- `http://localhost:3000/admin/users` — Users list + drawer
- `http://localhost:3000/admin/homepage` — Homepage config tabs
- `http://localhost:3000/admin/navigation` — Navigation tree
- `http://localhost:3000/admin/media` — Media grid
- `http://localhost:3000/admin/settings` — Settings form
- `http://localhost:3000/admin/login` — Login page

Verify:
- Sidebar icons render correctly
- Sidebar collapse works (desktop)
- Mobile hamburger menu works
- Table hover effects
- Drawer opens/closes smoothly
- Status pills display correct colors
- Empty states show when no data
- Login page gradient renders

- [ ] **Step 3: Commit**

```bash
git add frontend/assets/css/ frontend/layouts/ frontend/pages/admin/
git commit -m "feat: admin UI redesign — Modern Minimal style with slate palette, Lucide sidebar icons, borderless tables, drawer forms

Co-Authored-By: Claude Opus 4.7 <noreply@anthropic.com>"
```
