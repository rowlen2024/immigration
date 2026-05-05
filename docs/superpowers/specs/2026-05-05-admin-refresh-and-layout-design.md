# Admin List Pages: Refresh Button & Layout Fix

## Summary

Add a refresh button to all 8 admin list pages, and fix the horizontal layout of users and cases table pages so they fill the container width consistently with other list pages.

## Refresh Button

### Pages with toolbar (4 pages)

Refresh button placed at the right end of the toolbar row, pushed via `margin-left: auto`.

| Page | Toolbar content |
|------|----------------|
| `projects.vue` | search input + status filter + **refresh** |
| `cases.vue` | search input + **refresh** |
| `faqs.vue` | search input + project filter + **refresh** |
| `pages.vue` | search input + status filter + **refresh** |

### Pages without toolbar (4 pages)

Refresh button placed in the page header row, to the right of the title, before the action button.

| Page | Header right area |
|------|------------------|
| `users.vue` | **refresh** + "新建用户" |
| `leads.vue` | **refresh** (no action button) |
| `media.vue` | **refresh** + upload button |
| `navigation.vue` | **refresh** + "新建导航" (viewer role hides both via v-if) |

### Button implementation

```html
<el-button :icon="Refresh" circle @click="loadList" :loading="loading" />
```

- Icon: `Refresh` from `@element-plus/icons-vue` (already used in the project)
- `circle` variant for compact icon-only appearance
- Bound to existing `loading` ref — shows spin animation during data fetch
- Calls the same `loadList()` / `loadTree()` function already used by pagination and CRUD operations

## Layout Fix

### Problem

`users.vue` and `cases.vue` use fixed `width` on ALL table columns. No column stretches to fill remaining space, leaving a blank gap on the right side of the table.

### Fix

Change the primary content column from `width` to `min-width`:

| Page | Column | Before | After |
|------|--------|--------|-------|
| `users.vue` | 用户名 (`username`) | `width="140"` | `min-width="140"` |
| `users.vue` | 显示名称 (`display_name`) | `width="140"` | `min-width="140"` |
| `cases.vue` | 姓名 (`name`) | `width="120"` | `min-width="140"` |
| `cases.vue` | 来源国家 (`country_from`) | `width="120"` | `min-width="120"` |

This matches the pattern already used in `projects.vue`, `faqs.vue`, and `pages.vue` where at least one column uses `min-width` to absorb extra space. Role, status, sort, and action columns keep fixed `width` since their content is narrow and predictable.

## Files changed

1. `frontend/pages/admin/users.vue` — refresh button + layout fix
2. `frontend/pages/admin/cases.vue` — refresh button + layout fix
3. `frontend/pages/admin/projects.vue` — refresh button
4. `frontend/pages/admin/faqs.vue` — refresh button
5. `frontend/pages/admin/pages.vue` — refresh button
6. `frontend/pages/admin/leads.vue` — refresh button
7. `frontend/pages/admin/media.vue` — refresh button
8. `frontend/pages/admin/navigation.vue` — refresh button
