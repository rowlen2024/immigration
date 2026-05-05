# Admin Sidebar Menu Reorganization

## Summary

Restructure the admin sidebar from 11 flat menu items into 7 grouped entries with collapsible sub-menus for groups that have children.

## Menu Structure

| Top-level | Children |
|-----------|----------|
| 控制台 | *(独立)* |
| 内容管理 | 页面管理, 首页配置, 导航管理, 媒体库 |
| 项目管理 | *(独立)* |
| FAQ 管理 | *(独立)* |
| 案例管理 | *(独立)* |
| 咨询管理 | *(独立)* |
| 系统设置 | 用户管理, 网站设置 |

## Changes

**Single file**: `frontend/layouts/admin.vue`

- Replace flat `navItems` array with nested `NavGroup` structure (groups with optional children)
- Template: render groups with `children` as expandable sections, standalone items as direct `<NuxtLink>`
- Auto-expand group when child route is active
- Collapsed sidebar: groups show only icon, sub-items hidden
- Add expand/collapse chevron arrow with rotation transition
- Sub-items indented 12px, inherit existing `nav-item` styles
