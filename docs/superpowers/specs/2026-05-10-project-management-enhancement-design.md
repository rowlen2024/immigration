# 项目管理系统增强设计

**日期**: 2026-05-10
**状态**: 待实施

## 概述

对项目管理后台和前台项目详情页进行六项增强：滚动条修复、成功案例 Tab、资讯 Tab、项目对比配置、前台条件展示、API 数据扩展。

## 变更清单

| # | 变更 | 层级 |
|---|------|------|
| 1 | 修复申请流程 tab 的编辑/删除按钮被滚动条隐藏 | Frontend Admin |
| 2 | 新增成功案例 Tab（1:N 关系，案例生命周期依附项目） | Backend + Frontend |
| 3 | 新增资讯 Tab（M:N 关联 page_type=news 页面） | Backend + Frontend |
| 4 | 新增项目对比配置 Tab（选择对比项目 + 对比属性） | Backend + Frontend |
| 5 | 前台项目详情页条件展示（有内容则展示） | Frontend |
| 6 | API 扩展 Preload cases/news/compare_config | Backend |

---

## 1. 滚动条修复

### 问题
项目编辑 drawer（width=560px）中，申请流程 tab 的表格列总宽超过 drawer 宽度，操作列被横向滚动条遮挡。

### 修复
给所有子 tab 表格的**操作列**添加 `fixed="right"` 属性，使编辑/删除按钮始终可见。同时可适当调整 table `max-height` 避免竖向过度滚动。

---

## 2. 成功案例 Tab

### 关系
Project ↔ Case：1:N（现有模型已支持，Case 有 `ProjectID *uint64`）

### 后台操作

**Tab 内容**：表格展示当前项目已绑定的案例（名称、来源国、投资金额、处理周期、排序、操作按钮）

**新增案例**：
- 弹窗表单，复用 Case 模型字段：name, country_from, investment_amount, processing_period, description, photo_url, sort_order
- project_id 自动填充为当前项目 ID（必选项）

**编辑案例**：弹窗编辑，只能编辑属于当前项目的案例

**删除案例**：直接硬删除案例记录（`DELETE FROM cases WHERE id = ?`）

**选择已有案例**：也可通过下拉选择已有的未绑定案例（`project_id IS NULL`），绑定到当前项目

### 后端 API

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/admin/projects/:id/cases` | 查询已绑定案例 |
| POST | `/admin/projects/:id/cases` | 新增案例（绑定到项目） |
| PUT | `/admin/projects/:id/cases/:cid` | 编辑案例 |
| DELETE | `/admin/projects/:id/cases/:cid` | 删除案例（硬删除） |
| GET | `/admin/unbound-cases` | 查询未绑定的案例（供选择） |

Handler 层新增函数：`ListProjectCases`, `CreateCase`, `UpdateCase`, `DeleteCase`, `ListUnboundCases`

Service 层复用一个 `CaseService`（现有），扩展支持 `FindUnbound()` 查询。

删除逻辑从软删除改为硬删除（`r.db.Unscoped().Delete(&model.Case{}, id)` 或直接 `r.db.Delete(&model.Case{}, id)` 因为 GORM 默认用 DeletedAt 软删除，需要强制硬删）。

---

## 3. 资讯 Tab

### 关系
Project ↔ News（page_type='news' 的 Page）：M:N，中间表 `project_news`

### 数据模型

新建中间表：

```sql
CREATE TABLE project_news (
  project_id BIGINT UNSIGNED NOT NULL,
  page_id    BIGINT UNSIGNED NOT NULL,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (project_id, page_id),
  FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE,
  FOREIGN KEY (page_id) REFERENCES pages(id) ON DELETE CASCADE
);
```

新建 Go model：

```go
type ProjectNews struct {
    ProjectID uint64    `gorm:"primaryKey" json:"project_id"`
    PageID    uint64    `gorm:"primaryKey" json:"page_id"`
    CreatedAt time.Time `json:"created_at"`
    Page      *Page     `gorm:"foreignKey:PageID" json:"page,omitempty"`
}
```

Project model 新增关联：

```go
News []Page `gorm:"many2many:project_news;" json:"news,omitempty"`
```

### 后台操作

**Tab 内容**：表格展示当前项目已关联的新闻页面（标题、状态、关联时间、操作按钮）

**添加关联**：下拉搜索/多选已发布的新闻页面。数据源：`GET /api/v1/admin/pages?page_type=news&status=published&all=true`

**解除关联**：从 `project_news` 删除记录，不删除页面本身

### 后端 API

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/admin/projects/:id/news` | 查询已关联新闻 |
| POST | `/admin/projects/:id/news` | 添加关联（body: `{page_ids: [1,2,3]}`） |
| DELETE | `/admin/projects/:id/news/:page_id` | 解除关联 |

### 删除保护

解除资讯关联时无需额外检查，直接删除即可。删除新闻页面本身时（page 管理），CASCADE 自动清理 project_news 记录。

---

## 4. 项目对比配置 Tab

### 数据模型

```sql
CREATE TABLE compare_configs (
  id             BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
  project_id     BIGINT UNSIGNED NOT NULL UNIQUE,
  compare_with   JSON NOT NULL,    -- ["slug-1", "slug-2"]
  compare_fields JSON NOT NULL,    -- ["investment_amount", "country", ...]
  created_at     DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated_at     DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE
);
```

Go model：

```go
type CompareConfig struct {
    ID            uint64         `gorm:"primaryKey;autoIncrement" json:"id"`
    ProjectID     uint64         `gorm:"uniqueIndex;not null" json:"project_id"`
    CompareWith   datatypes.JSON `gorm:"type:json;not null" json:"compare_with"`
    CompareFields datatypes.JSON `gorm:"type:json;not null" json:"compare_fields"`
    CreatedAt     time.Time      `json:"created_at"`
    UpdatedAt     time.Time      `json:"updated_at"`
}
```

`compare_with`: JSON 数组，存储参与对比的项目 slug 列表。当前项目默认在列表首位（前端保证）。
`compare_fields`: JSON 数组，管理员勾选的对比属性标识列表。

### 可对比属性定义

后端常量（`config/compare_fields.go` 或直接在 model 中定义）：

```go
type CompareField struct {
    Key   string `json:"key"`
    Label string `json:"label"`
    From  string `json:"from"` // "project" | "requirements" | "timeline" | "cost_items"
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

### 后台操作

**Tab 内容**：
- 上半部分：选择对比项目。展示当前项目（默认选中，不可移除）+ 多选其他已发布项目
- 下半部分：勾选对比属性。展示所有 `CompareFields`，管理员勾选要参与对比的维度

**保存**：单行配置，按 project_id upsert

### 后端 API

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/admin/projects/:id/compare-config` | 获取对比配置 |
| PUT | `/admin/projects/:id/compare-config` | 保存对比配置（upsert） |
| GET | `/admin/compare-fields` | 获取可对比属性列表（常量） |

### 前台展示

项目详情页，如果该项目的 compare_config 存在且 `compare_with` 非空、`compare_fields` 非空，则渲染对比表格：

行 = compare_fields（属性），列 = 对比项目（当前 + 已选）。
对比数据通过请求 `GET /api/v1/projects/compare?slugs=a,b,c` 获取（现有 API）。

---

## 5. 前台条件展示

### 原则
所有模块统一使用 `v-if` 条件渲染：**后端返回数据为空（空数组 / 空字符串 / null）则不渲染 section**。不再展示任何默认/占位数据。

### 模块清单

| # | 模块 | 展示条件 | 数据来源 |
|---|------|---------|---------|
| 1 | 项目概览 | 始终展示 | project |
| 2 | 项目介绍 | `overview_text != ''` | `project.overview_text` |
| 3 | 申请条件 | `requirements.length > 0` | `project.requirements[]` |
| 4 | 费用明细 | `cost_items.length > 0` | `project.cost_items[]` |
| 5 | 申请流程 | `timeline_phases.length > 0` | `project.timeline_phases[]` |
| 6 | 常见问题 | `faqs.length > 0` | `project.faqs[]` |
| 7 | 成功案例 | `cases.length > 0` | `project.cases[]` |
| 8 | 最新资讯 | `news.length > 0` | `project.news[]` |
| 9 | 项目对比 | compare_config 存在且非空 | `project.compare_config` |

### 此逻辑对项目详情页的影响

现有 `pages/projects/[slug].vue` 已对 modules 2-6 使用 `v-if`。需要：
- 扩展 ApiProject 接口
- 新增 modules 7-9 的渲染组件
- 移除所有默认数据逻辑

---

## 6. 后端 API 数据扩展

`GET /api/v1/projects/:slug`（`GetProject` handler → `ProjectService.GetBySlug`）：

当前 Preload：Requirements, CostItems, TimelinePhases, Milestones, FAQs, Cases
扩展 Preload：
- `Cases`（已有）
- `News`（新增，通过 `many2many:project_news` 关联）
- `CompareConfig`（新增，通过 `HasOne` 关联）

Project model 新增：

```go
News          []Page         `gorm:"many2many:project_news;" json:"news,omitempty"`
CompareConfig *CompareConfig `gorm:"foreignKey:ProjectID" json:"compare_config,omitempty"`
```

Repository `FindBySlug` 添加 Preload：

```go
Preload("News").
Preload("CompareConfig")
```

---

## 数据库迁移

**000019_add_compare_config_and_project_news.up.sql**：

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

---

## 测试要点

| # | 测试场景 | 预期结果 |
|---|---------|----------|
| 1 | 项目编辑 drawer 申请流程 tab | 操作按钮不被滚动条遮挡 |
| 2 | 成功案例 Tab — 新增案例 | 案例创建成功，绑定到项目 |
| 3 | 成功案例 Tab — 删除案例 | 案例被硬删除 |
| 4 | 成功案例 Tab — 选择已有案例 | 未绑定案例列表正确过滤 |
| 5 | 资讯 Tab — 添加关联 | project_news 记录创建 |
| 6 | 资讯 Tab — 解除关联 | project_news 记录删除 |
| 7 | 项目对比 — 保存配置 | compare_configs upsert |
| 8 | 项目对比 — 前台展示 | 有配置时展示对比表格 |
| 9 | 项目详情页 — 无案例时 | 成功案例 section 不渲染 |
| 10 | 项目详情页 — 无资讯时 | 最新资讯 section 不渲染 |
| 11 | 项目详情页 — 无对比配置时 | 项目对比 section 不渲染 |
| 12 | 项目详情页 — 无 FAQ 时 | 常见问题 section 不渲染 |

---

## 文件变更清单

| 文件 | 操作 | 说明 |
|------|------|------|
| `backend/internal/model/project.go` | 修改 | 新增 News, CompareConfig 关联 |
| `backend/internal/model/compare_config.go` | 新建 | CompareConfig model |
| `backend/internal/model/project_news.go` | 新建 | ProjectNews model |
| `backend/internal/config/compare_fields.go` | 新建 | 可对比属性常量 |
| `backend/internal/repository/case_repo.go` | 修改 | 新增 FindUnbound, 硬删除 |
| `backend/internal/repository/compare_config_repo.go` | 新建 | CompareConfig CRUD |
| `backend/internal/repository/interfaces.go` | 修改 | 新增接口 |
| `backend/internal/service/case_svc.go` | 修改 | 扩展方法 |
| `backend/internal/service/compare_config_svc.go` | 新建 | 对比配置业务逻辑 |
| `backend/internal/handler/case_handler.go` | 修改 | 新增项目案例 CRUD handler |
| `backend/internal/handler/compare_config_handler.go` | 新建 | 对比配置 handler |
| `backend/internal/handler/project_handler.go` | 修改 | GetProject 扩展 Preload |
| `backend/internal/router/router.go` | 修改 | 新增路由 |
| `database/migrations/000019_add_compare_config_and_project_news.up.sql` | 新建 | 建表 |
| `database/migrations/000019_add_compare_config_and_project_news.down.sql` | 新建 | 回滚 |
| `frontend/pages/admin/projects.vue` | 修改 | 新增 3 个 tab + 滚动条修复 |
| `frontend/pages/projects/[slug].vue` | 修改 | 新增 section + 条件展示 |
