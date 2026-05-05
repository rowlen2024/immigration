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

## 数据库

- Docker Compose 启动 `mygo-mysql`（端口 3307），MySQL 首次初始化时自动按序执行 `database/migrations/` 下的 SQL 脚本
- 13 组迁移脚本 (000001–000013)：覆盖 users, projects/requirements/cost_items/timeline_phases/milestones, faqs, cases, pages, leads, media/home_configs, seed_data, navigations
- 启动时 GORM AutoMigrate 同步所有模型（代码中的 `database.AutoMigrate()` 包含全部 14 个模型，其中 operation_log 仅由 AutoMigrate 创建）

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
- **导航链接生成**：page 类型导航的链接为 `/<slug>`（非 `/pages/<slug>`），在 `nav_svc.go` 的 `fillLink()`/`fillLinks()` 中生成
- 新建 `.go` / `.vue` 文件时的命名规范：Go 用 snake_case 文件名，Vue 用 PascalCase 组件名或 kebab-case 路由页名
