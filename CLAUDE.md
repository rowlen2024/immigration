# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

# MyGo 移民

投资移民服务全栈 Web 应用 — Go (Gin + GORM) + Nuxt 3 SPA + MySQL 8.0。

## 常用命令

```bash
make up              # 启动 MySQL (docker-compose, 端口 3307)
make dev-backend     # Go 后端 (go run, 端口 8080)
make dev-frontend    # Nuxt 前端 (npm run dev, 端口 3000)
make test            # 后端测试 (go test ./...)
make test-cover      # 后端测试覆盖率
make build           # 前端构建
make install         # 安装所有依赖

# 精细命令
cd backend && go build ./...                              # 编译检查
cd backend && go test ./internal/service -run TestAuth     # 单个测试
cd frontend && npx nuxi typecheck                         # 前端类型检查
```

## 架构

```
Handler → Service → Repository → Model (GORM)
```

- 19 个 GORM 模型，16 个 repository 实现文件，16 个 handler 文件
- **构造函数注入**: `service.New()` 使用 `New*Service()` 构造函数注入依赖，交叉依赖（如 homeConfigSvc）在构造函数后单独赋值
- **级联删除**: `ProjectService.Delete()` 拆分为 `preDeleteCleanup` / `cascadeDeleteProjectResources` / `postDeleteHomeConfigCleanup` 三个辅助方法，无跨表事务
- **共享工具**: `repository/helpers.go` 提供 `CountByModel[T]` / `CountByModelRange[T]` / `PluckUploadsByColumn[T]` 泛型辅助；`service/sanitizer.go` 共享 HTML 净化策略
- **所有分页均为 SQL 分页**: `PageRepo.FindAllPaginated` 已接入 `PageService.AdminList`
- `.env` 查找顺序: cwd → `../.env`；不覆盖已有环境变量；自定义解析器（`bufio.Scanner`）

## 关键约束与陷阱

### API 信封

后端统一响应 `{code, message, data}`。分页接口 `data` 内含 `{items, total, page, perPage}`。

- **useApi()** (后台): 自动附加 Bearer token + 解包信封 → 直接拿到 `data`
- **useFetch / $fetch** (前台): 需手动取 `.data`，`$fetch` 不解包

### SPA 缓存陷阱

Nuxt SPA 模式对 `useFetch` 结果做 payload 缓存，客户端导航可能返回过期数据。新增 composable 必须在 `onMounted` 中用 `$fetch` 强制刷新。参考 `composables/useNavigation.ts`、`composables/useSiteConfig.ts`。**关键:** `$fetch` 不解包信封，必须手动 `.data`。

### 骨架屏

公众页面必须接 `PageSkeleton`（5 种 variant: hero/cards/list/detail/content）。`useFetch` 的 `pending` 自动管理加载态，无需手动维护。

### 路由注册顺序

字面量路由必须在参数化路由之前注册。`/projects/compare` 在 `/projects/:slug` 之前，否则 `:slug` 捕获 "compare"。所有 `/:slug` 类路由放在最后。

### 前后台接口隔离

| 端 | 前缀 | 鉴权 |
|----|------|------|
| 前台 | `/api/v1/<resource>` | 无 |
| 后台 | `/api/v1/admin/<resource>` | JWT + RBAC |

禁止交叉调用，禁止共用 handler。所有 admin 列表接口须同时支持分页和 `?all=true`。

### RBAC 角色

`admin`(全部) / `editor`(读+内容写+leads读) / `viewer`(只读)

### 导航约束

最大深度 3 层，防循环引用，防自引用。project→`/projects/<slug>`，page→`/pages/<slug>`。

### 对比系统

`/compare` (下拉选择) 和 `/compare/[a]-vs-[b]` (URL 直达) 两种方式，修改时需同步两处。项目详情页在 compare_config 存在时内嵌对比表。

### 其他

- 保留 `gorm.DeletedAt` 字段，但所有 `Delete()` 操作统一使用 `Unscoped()` 硬删除
- JWT Claims 类型 (`model.JWTClaims`) 由 middleware 和 service 共享
- 前端无测试，类型检查靠 `npx nuxi typecheck`
- Handler 测试用 `httptest.NewRecorder` + `gin.CreateTestContext`，连接真实 MySQL
- JWT 密钥从 `.env` 读取，生产务必更换

### 待办

- **RBAC 动态化**: 角色-权限映射硬编码在 `middleware/rbac.go`，待后续改为数据库配置

## 工作流规则

1. **中文沟通**
2. **禁止自动 git**: 改代码不自动 add/commit/push，除非明确要求
3. **变更前影响检查**: 修改公共接口前 Grep 所有引用点，列出影响报告等待确认
4. **防止回退**: 新功能不得破坏已有功能
5. **编译验证**: 后端改动后必须 `go build ./...` 通过
