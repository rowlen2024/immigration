# Admin Refresh Button & Layout Fix — Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add a refresh button to all 8 admin list pages and fix horizontal table layout on users/cases pages.

**Architecture:** Inline `<el-button>` per page — no shared component. Each page already has `loading` ref and `loadList()`/`loadTree()` function; the refresh button binds to both. Layout fix changes two columns per page from fixed `width` to `min-width`.

**Tech Stack:** Vue 3 + Nuxt 3 + Element Plus + TypeScript

---

### Task 1: users.vue — refresh button + layout fix

**Files:**
- Modify: `frontend/pages/admin/users.vue`

- [ ] **Step 1: Add Refresh icon import**

In the `<script setup>` section, add `Refresh` to the Element Plus icons import (currently only imports from `element-plus`, needs a new import line):

```typescript
import { ElMessage, type FormInstance, type FormRules } from 'element-plus';
import { Refresh } from '@element-plus/icons-vue';
import { getIconSvg } from '~/composables/lucideIcons';
```

Insert `import { Refresh } from '@element-plus/icons-vue';` between the existing element-plus import and the lucideIcons import.

- [ ] **Step 2: Add refresh button in header (no toolbar)**

Replace the header section (lines 2-6):

```html
<div class="admin-page-header">
  <h2 class="admin-page-title">用户管理</h2>
  <el-button v-if="isAdmin" type="primary" @click="openCreate">新建用户</el-button>
</div>
```

With:

```html
<div class="admin-page-header">
  <h2 class="admin-page-title">用户管理</h2>
  <div style="display:flex;align-items:center;gap:8px;">
    <el-button :icon="Refresh" circle @click="loadList" :loading="loading" />
    <el-button v-if="isAdmin" type="primary" @click="openCreate">新建用户</el-button>
  </div>
</div>
```

- [ ] **Step 3: Fix layout — username and display_name columns**

Change line 10 (`username` column):
```
<el-table-column prop="username" label="用户名" width="140">
```
To:
```
<el-table-column prop="username" label="用户名" min-width="140">
```

Change line 15 (`display_name` column):
```
<el-table-column prop="display_name" label="显示名称" width="140" />
```
To:
```
<el-table-column prop="display_name" label="显示名称" min-width="140" />
```

- [ ] **Step 4: Typecheck**

Run: `cd frontend && npx nuxi typecheck`

---

### Task 2: cases.vue — refresh button + layout fix

**Files:**
- Modify: `frontend/pages/admin/cases.vue`

- [ ] **Step 1: Add Refresh to icon import**

Change line 106:
```typescript
import { Search } from '@element-plus/icons-vue';
```
To:
```typescript
import { Search, Refresh } from '@element-plus/icons-vue';
```

- [ ] **Step 2: Add refresh button in toolbar**

In the toolbar (lines 8-17), add the refresh button after the search input:

```html
<div class="admin-toolbar">
  <el-input
    v-model="searchQuery"
    placeholder="搜索姓名..."
    :prefix-icon="Search"
    clearable
    class="admin-search-input"
    @input="onSearch"
  />
  <el-button :icon="Refresh" circle @click="loadList" :loading="loading" style="margin-left:auto;" />
</div>
```

Key change: add `style="margin-left:auto;"` on the refresh button to push it to the right.

- [ ] **Step 3: Fix layout — name and country_from columns**

Change line 21 (`name` column):
```
<el-table-column prop="name" label="姓名" width="120">
```
To:
```
<el-table-column prop="name" label="姓名" min-width="140">
```

Change line 26 (`country_from` column):
```
<el-table-column prop="country_from" label="来源国家" width="120" />
```
To:
```
<el-table-column prop="country_from" label="来源国家" min-width="120" />
```

- [ ] **Step 4: Typecheck**

Run: `cd frontend && npx nuxi typecheck`

---

### Task 3: projects.vue — refresh button

**Files:**
- Modify: `frontend/pages/admin/projects.vue`

- [ ] **Step 1: Add Refresh to icon import**

Change line 347:
```typescript
import { Search } from '@element-plus/icons-vue';
```
To:
```typescript
import { Search, Refresh } from '@element-plus/icons-vue';
```

- [ ] **Step 2: Add refresh button in toolbar**

Replace the toolbar section (lines 9-29):

```html
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
  <el-button :icon="Refresh" circle @click="loadList" :loading="loading" style="margin-left:auto;" />
</div>
```

The only addition is the `<el-button>` line at the end of the toolbar div.

- [ ] **Step 3: Typecheck**

Run: `cd frontend && npx nuxi typecheck`

---

### Task 4: faqs.vue — refresh button

**Files:**
- Modify: `frontend/pages/admin/faqs.vue`

- [ ] **Step 1: Add Refresh to icon import**

Change line 112:
```typescript
import { Search } from '@element-plus/icons-vue';
```
To:
```typescript
import { Search, Refresh } from '@element-plus/icons-vue';
```

- [ ] **Step 2: Add refresh button in toolbar**

Replace the toolbar section (lines 8-33):

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
    <el-button :icon="Refresh" circle @click="loadList" :loading="loading" style="margin-left:auto;" />
  </div>
</div>
```

Add the `<el-button>` line before the closing `</div>` of `admin-toolbar-row`.

- [ ] **Step 3: Typecheck**

Run: `cd frontend && npx nuxi typecheck`

---

### Task 5: pages.vue — refresh button

**Files:**
- Modify: `frontend/pages/admin/pages.vue`

- [ ] **Step 1: Add Refresh to icon import**

Change line 120:
```typescript
import { Search } from '@element-plus/icons-vue';
```
To:
```typescript
import { Search, Refresh } from '@element-plus/icons-vue';
```

- [ ] **Step 2: Add refresh button in toolbar**

Replace the toolbar section (lines 8-22):

```html
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
  <el-button :icon="Refresh" circle @click="loadList" :loading="loading" style="margin-left:auto;" />
</div>
```

Add the `<el-button>` line at the end of the toolbar div.

- [ ] **Step 3: Typecheck**

Run: `cd frontend && npx nuxi typecheck`

---

### Task 6: leads.vue — refresh button

**Files:**
- Modify: `frontend/pages/admin/leads.vue`

- [ ] **Step 1: Add Refresh icon import**

Change line 100:
```typescript
import { ElMessage } from 'element-plus';
```
Add after it:
```typescript
import { Refresh } from '@element-plus/icons-vue';
```

- [ ] **Step 2: Add refresh button in header**

Replace the header section (lines 2-5):

```html
<div class="admin-page-header">
  <h2 class="admin-page-title">咨询管理</h2>
  <el-button :icon="Refresh" circle @click="loadList" :loading="loading" />
</div>
```

- [ ] **Step 3: Typecheck**

Run: `cd frontend && npx nuxi typecheck`

---

### Task 7: media.vue — refresh button

**Files:**
- Modify: `frontend/pages/admin/media.vue`

- [ ] **Step 1: Add Refresh icon import**

Change line 72:
```typescript
import { ElMessage } from 'element-plus';
```
Add after it:
```typescript
import { Refresh } from '@element-plus/icons-vue';
```

- [ ] **Step 2: Add refresh button in header**

Replace the header section (lines 2-18):

```html
<div class="admin-page-header">
  <h2 class="admin-page-title">媒体库</h2>
  <div style="display:flex;align-items:center;gap:8px;">
    <el-button :icon="Refresh" circle @click="loadList" :loading="loading" />
    <el-upload
      :action="uploadUrl"
      :headers="uploadHeaders"
      accept=".jpg,.jpeg,.png,.webp"
      :show-file-list="false"
      :on-success="handleUploadSuccess"
      :on-error="handleUploadError"
    >
      <el-button type="primary">
        <span v-html="getIconSvg('upload', 16)" style="margin-right:6px;vertical-align:middle"></span>
        上传图片
      </el-button>
    </el-upload>
  </div>
</div>
```

Wrap the existing upload button and the new refresh button in a flex container.

- [ ] **Step 3: Typecheck**

Run: `cd frontend && npx nuxi typecheck`

---

### Task 8: navigation.vue — refresh button

**Files:**
- Modify: `frontend/pages/admin/navigation.vue`

- [ ] **Step 1: Add Refresh icon import**

Change line 166:
```typescript
import { ElMessage, type FormInstance, type FormRules } from 'element-plus';
```
Add after it:
```typescript
import { Refresh } from '@element-plus/icons-vue';
```

- [ ] **Step 2: Add refresh button in header**

Replace the header section (lines 2-6):

```html
<div class="admin-page-header">
  <h2 class="admin-page-title">导航管理</h2>
  <div v-if="!isViewer" style="display:flex;align-items:center;gap:8px;">
    <el-button :icon="Refresh" circle @click="loadTree" :loading="loading" />
    <el-button type="primary" @click="openCreate()">新建导航</el-button>
  </div>
</div>
```

Note: `loadTree` instead of `loadList` (navigation uses a tree, not a table). The `v-if="!isViewer"` wraps both buttons so viewers see neither.

- [ ] **Step 3: Typecheck**

Run: `cd frontend && npx nuxi typecheck`

---

### Task 9: Final verification

- [ ] **Step 1: Full typecheck**

Run: `cd frontend && npx nuxi typecheck`
Expected: No type errors

- [ ] **Step 2: Verify all 8 files have Refresh imported**

Run: `cd frontend && grep -rn "Refresh" pages/admin/`
Expected: 8 files with `Refresh` in import or usage

- [ ] **Step 3: Verify no width-only columns remain on users/cases primary columns**

Run: `cd frontend && grep -n "width=" pages/admin/users.vue pages/admin/cases.vue`
Expected: Only role/status/sort/action columns have `width`; username/display_name/name/country_from use `min-width`
