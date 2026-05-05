# FAQ 项目归属 & 筛选 设计

**日期**: 2026-05-05  
**状态**: 待实现

## 动机

1. FAQ 管理列表无法显示所属项目（API 不返回 project name）
2. 管理后台搜索参数传到后端但后端不处理
3. 公开 FAQ 页筛选按钮写死，非动态项目
4. `is_global` 概念语义不清

## 核心逻辑规则

- `project_id` 决定 FAQ 出现在哪个项目详情页（非空 = 归属项目，空 = 不归属任何项目）
- `is_global` 纯粹控制 FAQ 是否出现在公开 FAQ 总页（`/faq`）
- 一个 FAQ 可以同时 `project_id` 非空 + `is_global=true`（既在项目详情页展示，也在 FAQ 总页展示）
- 两者互不干扰

---

## 一、数据 & API

### 新增 DTO

`backend/internal/dto/faq_response.go`（新建）：

```go
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

### Repository 层

```go
type FAQQueryParams struct {
    ProjectID *uint64
    IsGlobal  *bool
    Search    string   // 模糊匹配 question OR answer
    Page      int
    PerPage   int
}

// FAQRepository 接口更新
type FAQRepository interface {
    FindAll(params FAQQueryParams) ([]model.FAQ, int64, error)  // 改：支持过滤+分页
    Create(faq *model.FAQ) error
    Update(faq *model.FAQ) error
    Delete(id uint64) error
    // 移除 FindByProjectID, FindGlobal, Search（被 FindAll 统一替代）
}
```

`FindAll` 使用 GORM `Joins("LEFT JOIN projects ON projects.id = faqs.project_id")` 一次查询带出 project name/slug。

### Service 层

| 方法 | 变更 |
|------|------|
| `List(projectID *uint64, isGlobal *bool)` | 改签名，接收可选筛选，返回 `[]dto.FAQResponse` |
| `AdminList(projectID *uint64, search string, page, perPage int)` | 改签名，接收项目筛选+搜索，返回 `[]dto.FAQResponse, int64` |

### 路由 & Handler

| 端点 | 变更 |
|------|------|
| `GET /api/v1/faqs` | 新增可选参数 `project_id`, `is_global`；返回 `FAQResponse[]`（含 project_name/slug）；默认不过滤 |
| `GET /api/v1/admin/faqs` | 新增可选参数 `project_id`；修复 `search` 参数使其实际生效；返回 `FAQResponse[]` |

---

## 二、前端变更

### 公开 FAQ 页 (`pages/faq.vue`)

- 移除所有硬编码 mock 回退数据
- 调用 `GET /api/v1/faqs`，API 失败则显示错误空状态
- 从返回结果中提取 `project_slug` + `project_name` 去重，动态生成筛选按钮
- "全部"按钮显示所有 FAQ；各项目按钮按 `project_slug` 内存过滤
- 全局 FAQ（`is_global=true`）始终出现在"全部"视图中

### 项目详情页 (`pages/projects/[slug].vue`)

- 新增 FAQ 区块：调用 `GET /api/v1/faqs?project_id=<当前项目ID>`
- API 失败/无数据时静默隐藏该区块
- 展示风格与项目详情页其他区块一致

### 管理 FAQ 页 (`pages/admin/faqs.vue`)

- 添加项目下拉筛选器（与搜索框同行，可选、可清空）
- "所属项目"列改为显示 `row.project_name`（去掉无用的 `row.project`）
- 后端修复后搜索框 debounce 实际生效
- 创建/编辑表单项目选择器保留不变

---

## 三、错误处理

| 场景 | 处理 |
|------|------|
| 公开 FAQ API 失败 | 显示 "加载常见问题失败，请稍后重试" |
| 项目详情 FAQ API 失败 | 静默隐藏 FAQ 区块 |
| 管理后台 API 失败 | 维持现有 ElMessage.error |

## 四、测试

- **Repository**: `TestFAQRepo_FindAll_QueryParams` — 覆盖各参数组合
- **Service**: `TestFAQService_List_WithProjectFilter`, `TestFAQService_AdminList_WithSearch`
- **Handler**: `TestListFAQs_ResponseContainsProjectName`
- **前端**: `npx nuxi typecheck`

## 五、涉及文件

| 文件 | 改动 |
|------|------|
| `backend/internal/dto/faq_response.go` | 新建 |
| `backend/internal/repository/interfaces.go` | FAQRepository 接口更新 |
| `backend/internal/repository/faq_repo.go` | 重写 FindAll，新增 QueryParams |
| `backend/internal/service/faq_svc.go` | List/AdminList 签名改，填充 project 信息 |
| `backend/internal/service/faq_svc_test.go` | 更新测试匹配新签名 |
| `backend/internal/handler/faq_handler.go` | handler 读取新参数，返回 FAQResponse |
| `backend/internal/handler/faq_handler_test.go` | 更新测试 |
| `frontend/pages/faq.vue` | 动态筛选按钮，移除 mock |
| `frontend/pages/projects/[slug].vue` | 新增 FAQ 区块 |
| `frontend/pages/admin/faqs.vue` | 搜索修复 + 项目筛选器 + project_name 列 |
