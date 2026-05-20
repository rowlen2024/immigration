# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

# MyGo 移民 (MyGo Immigration)

投资移民服务全栈 Web 应用 — Go 后端 + Nuxt 3 前端 + MySQL 数据库。

## 项目结构

```
immigration/
├── backend/                  # Go 后端 (Gin + GORM)
│   ├── cmd/server/main.go    # 入口
│   └── internal/
│       ├── config/           # 环境变量加载 (.env)
│       ├── database/         # MySQL 连接 & AutoMigrate
│       ├── model/            # GORM 模型 (14个)
│       ├── repository/       # 数据访问层 (interfaces + 实现)
│       ├── service/          # 业务逻辑层 (含 bluemonday 消毒)
│       ├── handler/          # HTTP handler (路由处理)
│       ├── middleware/        # CORS, JWT Auth, RBAC, RateLimit
│       ├── router/           # 路由注册 (Gin)
│       └── dto/              # Request/Response 结构体
├── frontend/                 # Nuxt 3 前端 (SPA 模式)
│   ├── pages/                # 页面组件 (文件路由)
│   │   ├── [...slug].vue     # 动态 CMS 页面 (catch-all, 渲染 DB 中的 page)
│   │   └── admin/            # 管理后台 (auth 中间件保护)
│   ├── components/           # RichEditor (Tiptap), Header, Footer, admin/*, project/*
│   ├── composables/          # useApi, useAuth, useSeo, useSiteConfig
│   ├── layouts/              # default.vue (公众), admin.vue (分组侧边栏)
│   └── middleware/auth.ts    # 客户端路由鉴权
├── database/migrations/      # MySQL 迁移 (13组, 000001-000013)
├── docker-compose.yml        # MySQL 8.0 (端口 3307)
├── Makefile                  # 开发命令
└── .env                      # 环境变量
```

## 技术栈

| 层 | 技术 |
|---|------|
| 后端框架 | **Gin** (github.com/gin-gonic/gin) |
| ORM | **GORM** (gorm.io/gorm) + MySQL driver |
| 认证 | **JWT** (golang-jwt/jwt/v5), bcrypt |
| 前端框架 | **Nuxt 3** (SPA, ssr: false) |
| UI 库 | **Element Plus** |
| 状态管理 | **Pinia** |
| 数据库 | **MySQL 8.0** (Docker, utf8mb4) |
| HTML 消毒 | **bluemonday** (用户输入) |

## 后端架构

分层架构：`Handler → Service → Repository → Model (GORM)`

- **Model**: GORM 结构体，对应数据库表。共 14 个模型：`User`, `Project`(+`Requirement`/`CostItem`/`TimelinePhase`/`Milestone`), `FAQ`, `Case`, `Page`, `Lead`, `Media`, `Navigation`, `HomeConfig`, `OperationLog`
- **Repository**: 接口定义在 `interfaces.go`，实现按模型拆分文件
- **Service**: 业务逻辑，通过 `Service` 结构体聚合
- **Handler**: 通过 `Handler` 结构体注入 Service，处理 HTTP 请求
- **DTO**: `request.go` / `response.go` — 统一响应格式 `{code, message, data}`
- **组合模式**: `service.go` 和 `repository.go` 分别定义 `Service` 和 `Repository` 组合结构体，聚合所有子模块。`service.New()` 使用直接结构体字面量注入依赖，而非各子模块的 `New*Service()` 构造函数（那些构造函数缺少跨模块依赖注入）
- **无事务支持**: 服务层无法传递 GORM 事务，所有 repository 方法使用自己的 `db *gorm.DB`，多步操作无法保证原子性
- **分页注意**: `CaseService.AdminList()` 和 `PageService.AdminList()` 是内存分页（先查全表再切片），`ProjectService.List()` 是 SQL 分页

## API 路由 (所有前缀 `/api/v1`)

### 公开路由
| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/health` | 健康检查 |
| GET | `/projects` | 项目列表 |
| GET | `/projects/:slug` | 项目详情 (含关联数据) |
| GET | `/projects/compare` | 项目对比 (?slugs=a,b) |
| GET | `/faqs` | FAQ 列表 |
| GET | `/pages` | 页面列表 |
| GET | `/pages/*slug` | 页面详情 (动态路由) |
| GET | `/cases` | 案例列表 |
| GET | `/home-config` | 首页配置 |
| GET | `/navigation` | 导航树 (嵌套 children) |
| GET | `/search` | 全文搜索 |
| POST | `/leads` | 提交咨询 (限流 10/min) |
| POST | `/auth/login` | 登录 (限流 5/min) |
| POST | `/auth/refresh` | 刷新 token |

### Admin 路由 (需 Bearer Token)
- `admin:read` 权限：dashboard/stats, 查看项目/FAQ/页面/案例/首页配置
- `content:write` 权限：CRUD FAQ/页面/案例，上传媒体，更新首页配置
- `projects:write` 权限：CRUD 项目
- `leads:read` 权限：查看/更新咨询
- `admin:write` 权限：用户管理
- 导航 CRUD 使用 `admin:read`（查看）和 `content:write`（增改删）

三种角色 RBAC：`admin`（全部权限）、`editor`（读+内容写+leads读）、`viewer`（只读）

**关键：路由注册顺序** — `/projects/compare` 必须在 `/projects/:slug` 之前注册，否则 `:slug` 会捕获 "compare"。同理 `/:slug` 类路由应在所有字面量路由之后注册。

**Admin 子资源路由** 统一使用 `/admin/projects/:id/<资源名>` 模式，包含 requirements、cost-items、timeline-phases、advantages、cases、news、compare-config，各有独立的 CRUD 端点。

## 前端架构

- **SSR 关闭** (SPA)，开发端口 3000，API 代理到 `localhost:8080`
- **useApi()** composable 封装了 `$fetch`，自动：附加 Bearer token、解包 `{data}` 信封、处理 401 跳转登录页
- **useAuth()** composable 管理 token 和用户状态（localStorage 持久化）
- **布局**：默认布局 `default.vue`（Header + Footer），后台布局 `admin.vue`（侧边栏 + 顶栏）
  - 后台侧边栏使用分组结构：2 个可展开组（内容管理 4 子项、系统设置 2 子项）+ 5 个独立项（控制台/项目管理/FAQ/案例/咨询）
  - 侧边栏支持折叠（桌面端 `sidebarCollapsed`）和移动端抽屉（`mobileOpen` + overlay）
  - 当前路由匹配的子组自动展开
- **认证中间件**：`auth.ts` 检查 `isAuthenticated`，未登录重定向到 `/admin/login`
- **动态页面**：`[...slug].vue` 匹配所有 CMS 自定义页面路径
- **富文本编辑器**：`components/RichEditor.vue` — 基于 Tiptap v2，支持表格/图片/视频/代码块/文字样式/源码切换，图片上传复用 `/api/v1/admin/media/upload`
  - 后端对应 bluemonday 自定义策略（`page_svc.go`），除标准元素外额外允许 table/iframe/video/span style
  - **bluemonday 策略重复**: `page_svc.go` 和 `case_svc.go` 有相同的 sanitizer 策略定义，修改时需同步两处
- **API 调用两种路径**:
  - `useApi()` 自动解包 `{data}` 信封并附加 token，用于所有 admin 页面
  - 公众页面（`projects/[slug].vue`、`[...slug].vue`、`compare.vue`）使用 `useFetch` / `$fetch` 直接调用，**需手动处理信封**：通过 `transform` 回调或访问 `.data` 属性
- **对比系统**: `/compare`（下拉选择器）和 `/compare/[a]-vs-[b]`（URL 直达，含硬编码降级数据）两种方式并存，修改对比逻辑时需同步两处
- **项目详情的页内对比**: `projects/[slug].vue` 在 compare_config 存在时单独请求 `/api/v1/projects/compare` 内嵌展示对比表
- **`?all=true` 参数**: 管理端的项目列表支持此参数以跳过默认过滤，用于下拉选择等场景

## 数据库

- Docker Compose 启动 `mygo-mysql`（端口 3307），MySQL 首次初始化时自动按序执行 `database/migrations/` 下的 SQL 脚本
- 13 组迁移脚本 (000001–000013)：覆盖 users, projects/requirements/cost_items/timeline_phases/milestones, faqs, cases, pages, leads, media/home_configs, seed_data, navigations
- 启动时 GORM AutoMigrate 同步所有模型（代码中的 `database.AutoMigrate()` 包含全部 14 个模型，其中 operation_log 仅由 AutoMigrate 创建）

## 配置与启动

- 自定义 `.env` 解析器（`config/config.go`），不依赖第三方库，`bufio.Scanner` 逐行读取
- **查找顺序**: 先在 cwd 找 `.env`，找不到则查 `../.env`（适配从子目录运行的情况）
- **不覆盖已有环境变量**: `os.Getenv(key) == ""` 时才设置，系统环境变量优先
- 默认值: 端口 8080, DB `mygo:mygo123@/mygo_immigration`, JWT secret `change-me-in-production`, CORS origin `http://localhost:3000`
- **启动顺序**: `config.Load()` → `database.InitMySQL()` → `database.RunMigrations()` → `database.AutoMigrate()`
- **迁移系统**: SQL 文件按文件名排序执行，有自定义分号分隔器（尊重字符串字面量内的分号）
- `CompareFields` 硬编码在 `compare_fields.go`，共 10 个对比维度
- Nuxt 配置中有**双重代理**: `nitro.devProxy`（SSR 构建时）和 `vite.server.proxy`（开发 HMR 时），都代理 `/api` 和 `/uploads` 到 `localhost:8080`

## 导航约束

- 最大深度 3 层，创建/更新时服务层校验
- 防循环引用：不能将节点设为自身的后代
- 防自引用：父节点不能是自己
- 链接自动生成：project 类型 → `/projects/<slug>`，page 类型 → `/pages/<slug>`（`nav_svc.go` 的 `fillLink()` 方法）

## 常用命令

```bash
make up              # 启动 MySQL (docker-compose)
make dev-backend     # Go 后端开发 (go run)
make dev-frontend    # Nuxt 前端开发 (npm run dev)
make test            # 后端测试 (go test ./...)
make test-cover      # 后端测试覆盖率
make build           # 前端构建
make install         # 安装所有依赖

# 其他
cd backend && go test ./... -v                        # 运行全部后端测试
cd backend && go test ./internal/service -run TestAuth   # 运行单个测试
cd frontend && npx nuxi prepare                       # 生成前端类型 (.nuxt/types)
cd frontend && npx nuxi typecheck                     # 前端类型检查
```

## 测试

- Go 测试文件与源码同目录，命名 `*_test.go`：handler 和 service 各有 6-8 个测试文件
- Handler 测试使用 `httptest.NewRecorder` + `gin.CreateTestContext` 模拟 HTTP 请求
- Service 测试直接调用 service 方法，使用真实 MySQL 数据库（测试隔离依赖数据库状态）
- 前端暂无单元测试，类型检查通过 `npx nuxi typecheck` 验证

## 开发注意事项

- JWT 密钥从 `.env` 读取，生产环境务必更换
- 密码使用 bcrypt 哈希存储
- 用户输入（页面内容等）经过 bluemonday 消毒
- 所有 GORM 模型使用软删除 (`gorm.DeletedAt`)
- **API 信封处理**：后端响应 `{code, message, data}`，`useApi()` 自动解包返回 `data`；若用 `useFetch` / `$fetch` 直接调用，需手动提取 `.data`
  - 分页接口返回 `{items, total, page, perPage}`（在 `data` 内）
- 新建 `.go` / `.vue` 文件时的命名规范：Go 用 snake_case 文件名，Vue 用 PascalCase 组件名或 kebab-case 路由页名

## 工作流规则

1. **中文沟通**: 所有对话使用中文。
2. **禁止自动 Git**: 改完代码后不执行 git add/commit/push，除非用户明确要求。
3. **变更前影响检查**: 修改任何公共接口（函数签名、API 路由、DTO 结构体、Vue composable/props/emits）时，必须先用 Grep 搜索所有引用点，列出影响报告并等待用户确认后，才能开始改代码。
4. **防止回退**: 新增功能时不得破坏已有功能。改动完成后必须确认相关已有功能不受影响。
5. **编译验证**: 每次后端改动完成后运行 `go build ./...` 验证编译通过。

## API 接口规范 (2026-05-20 制定)

**所有新增/修改接口必须严格遵循此规范。**

### 1. 前后台接口隔离

| 端 | 路由前缀 | 鉴权 | 说明 |
|----|---------|------|------|
| 前台（公开） | `/api/v1/<resource>` | 无需 auth | 仅供前台页面使用 |
| 后台（管理） | `/api/v1/admin/<resource>` | 需 JWT + RBAC | 仅供后台页面使用 |

**禁止事项:**
- 前台页面不得调用 `/admin/` 前缀接口
- 后台页面不得调用前台接口（无 `/admin/` 前缀）
- admin 和 public 不得共用同一个 handler（即使逻辑相同，也需分离为两个方法）

### 2. 响应格式规范

**分页列表** → `PaginatedResponse` → `{code, data, message, pagination}`
```go
c.JSON(http.StatusOK, dto.SuccessPaginated(items, page, perPage, total))
// pagination: { page, per_page, total }
```

**全量列表（?all=true）** → `Response` → `{code, data, message}`
```go
c.JSON(http.StatusOK, dto.Success(items))
// data 为原始数组，不包装 pagination
```

**单条查询** → `Response` → `{code, data, message}`

**创建/更新/删除** → `Response` → `{code, data, message}`

### 3. ?all=true 模式

- 所有 admin 列表接口必须同时支持分页和 `?all=true` 全量查询
- `?all=true` 返回 `dto.Success(items)`（非分页），前端 `useApi()` 解包为原始数组
- 不带 `?all=true` 时默认分页，返回 `dto.SuccessPaginated`，前端 `useApi()` 解包为 `{items, total, page, perPage}`

Handler 模板：
```go
func (h *Handler) AdminListXXX(c *gin.Context) {
    if c.Query("all") == "true" {
        items, _, err := h.svc.XXX.List(1, 1000, ...)
        // ...
        c.JSON(http.StatusOK, dto.Success(items))
        return
    }
    page, perPage := parsePagination(c)
    items, total, err := h.svc.XXX.List(page, perPage, ...)
    c.JSON(http.StatusOK, dto.SuccessPaginated(items, page, perPage, total))
}
```

### 4. 前端 API 调用规则

| 端 | 调用方式 | 解包行为 |
|----|---------|---------|
| 前台页面 | `useFetch` / `$fetch` 直接调用 | 手动访问 `response.data` |
| 后台页面 | `useApi()` composable | 自动解包：paginated→`{items,total}`，非分页→原始数据 |

**后台页面 `?all=true` 调用模板**：
```typescript
const api = useApi()
const data = await api<ItemType[]>('/admin/resource?all=true')
// data 为 ItemType[]，直接使用，无需 .items
```

### 5. 路由注册顺序

- 字面量路由在参数化路由之前注册（如 `/projects/compare` 在 `/projects/:slug` 之前）
- 公开路由和 admin 路由分属不同 `gin.RouterGroup`

### 6. 现有接口状态速查

| 接口 | 分页? | ?all=true? | Handler 方法 |
|------|-------|-------------|-------------|
| `GET /projects` | ✓ | — | `ListProjects` |
| `GET /faqs` | ✗ (全量) | — | `ListFAQs` |
| `GET /cases` | ✗ (全量) | — | `ListCases` |
| `GET /testimonials` | ✗ (全量) | — | `ListAllTestimonials` |
| `GET /lawyers` | ✗ (全量) | — | `ListLawyers` |
| `GET /home-config` | ✗ | — | `GetHomeConfig` (公共) |
| `GET /site-config` | ✗ | — | `GetSiteConfig` (公共) |
| `GET /navigation` | ✗ (树) | — | `GetNavigation` |
| `GET /admin/projects` | ✓ | ✓ | `AdminListProjects` |
| `GET /admin/faqs` | ✓ | ✗ | `AdminListFAQs` |
| `GET /admin/pages` | ✓ | ✓ | `AdminListPages` |
| `GET /admin/cases` | ✓ | ✓ | `AdminListCases` |
| `GET /admin/lawyers` | ✓ | ✓ | `AdminListLawyers` |
| `GET /admin/testimonials` | ✗ (全量) | — | `AdminListTestimonials` |
| `GET /admin/leads` | ✓ | ✗ | `AdminListLeads` |
| `GET /admin/users` | ✓ | ✗ | `AdminListUsers` |
| `GET /admin/media` | ✓ | ✗ | `ListMedia` |
| `GET /admin/navigation` | ✗ (树) | — | `AdminListNavigationTree` |
| `GET /admin/home-config` | ✗ | — | `GetAdminHomeConfig` |
| `GET /admin/site-config` | ✗ | — | `GetSiteConfig` |
| `GET /admin/compare-fields` | ✗ | — | `ListCompareFields` |
| `GET /admin/dashboard/stats` | ✗ | — | `DashboardStats` |

**新增接口时，将新条目追加到此表。**
