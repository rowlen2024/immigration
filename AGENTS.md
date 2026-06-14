# AGENTS.md

This file provides guidance to Codex (Codex.ai/code) when working with code in this repository.

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

## CodeGraph 代码知识图谱（常驻）

项目已安装 [CodeGraph](https://github.com/colbymchenry/codegraph)（v0.9.9），提供语义级代码搜索和调用链追踪。MCP Server 已在 `.mcp.json` 配置，Codex 会话启动时自动加载。

**索引状态**: 187 文件 · 2,709 节点 · 6,056 边（Go + Vue + TypeScript）

**MCP 工具速查**（前缀 `mcp__codegraph__`）：

| 工具 | 用途 | 典型场景 |
|------|------|---------|
| `codegraph_search` | 按名称搜索符号 | 快速定位函数/类型定义 |
| `codegraph_context` | 一站式获取相关符号+关系+源码 | 理解某个任务的代码上下文 |
| `codegraph_trace` | 追踪两个符号之间的调用路径 | 理解数据流/调用链 |
| `codegraph_impact` | 分析修改某符号的影响范围 | 重构前评估影响面 |
| `codegraph_callers` | 查找调用者 | 修改前确认所有引用点 |
| `codegraph_callees` | 查找被调用者 | 理解函数内部依赖 |
| `codegraph_explore` | 获取一组相关符号及关系图 | 快速了解一个模块 |
| `codegraph_node` | 获取单个符号详情（可选源码） | 查看具体实现 |
| `codegraph_files` | 获取索引中的文件结构 | 替代文件系统扫描 |
| `codegraph_status` | 检查索引健康和统计 | 确认索引是否最新 |

**使用原则**：
- 重构公共接口时优先用 `codegraph_impact` + `codegraph_callers`，结果比 grep 更准确（可追踪 interface→impl 动态派发）
- 理解陌生调用链时优先用 `codegraph_trace`，比手动 grep 效率高 10 倍+
- 代码改动后索引不会自动更新，必要时跑 `codegraph sync` 增量同步
- 如果工具返回结果不理想，回退到 Grep/Glob

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

### 列表接口重构规范（2026-05-30 起执行）

**已重构模块**: case（案例）、lead（咨询）、testimonial（客户评价）、lawyer（律师）、page（页面）、faq（常见问题）、project（项目）。其余模块（media 等）后续按此规范逐一重构。

**核心规则**:

1. **三层各保留一个 List 方法**: Handler → `ListXxx(c *gin.Context)`，Service → `List(req dto.XxxListRequest) ([]model.Xxx, int64, error)`，Repository → `FindAll(filter XxxFilter) ([]model.Xxx, int64, error)`

2. **请求结构体**: 每个资源定义独立的 `XxxListRequest`（放在 `dto/` 下），嵌入通用分页字段。字段名按数据库字段命名（如 `country_from`、`name`），不用泛化的 `search`/`q`。

   ```go
   type CaseListRequest struct {
       Page        int     `form:"page" binding:"omitempty,min=1"`
       PerPage     int     `form:"per_page" binding:"omitempty,min=1,max=100"`
       ProjectID   *uint64 `form:"project_id"`
       CountryFrom string  `form:"country_from"`
       Name        string  `form:"name"` // LIKE 模糊匹配
   }
   ```

3. **分页规则**: 传了 `page` + `per_page` → SQL 分页 + `dto.SuccessPaginated`；都没传 → 全量 + `dto.Success`。无论分不分页，所有筛选参数均生效。

4. **废弃 `?all=true`**: 不传分页参数即全量，不再需要 `?all=true` 开关。同步修改前端调用方。

5. **Repository Filter**: 每个资源定义 `XxxFilter` 结构体（放在 `repository/`），Repository 的 `FindAll` 方法内部根据 filter 字段动态构建 GORM 查询，Page>0 时加 LIMIT/OFFSET。

6. **公共 PaginationRequest**: `dto.PaginationRequest` 只保留 `Page` + `PerPage`，不含任何业务字段（已移除 `Status`/`Q`）。

7. **子资源路由保留**: 如 `/admin/projects/:id/cases`，保留独立路由，内部构造 `XxxListRequest` 调用同一个 `List` 方法。

8. **检索字段用 LIKE**: 数据库字段 `name` 等文本检索保持 `LIKE '%xxx%'` 模糊匹配。

**重构检查清单**（每个模块）:
- [ ] 新建 `dto/xxx_request.go` — `XxxListRequest`
- [ ] 更新 `repository/interfaces.go` — 接口改为单个 `FindAll(XxxFilter)`
- [ ] 更新 `repository/xxx_repo.go` — 删除多个 Find 方法，新增 `XxxFilter` + `FindAll`
- [ ] 更新 `service/xxx_svc.go` — 删除多个 List 方法，新增单个 `List`
- [ ] 更新 `handler/xxx_handler.go` — List handler 改为 `ShouldBindQuery` + 统一调用
- [ ] 更新 `handler/lead_handler.go` 等引用 `PaginationRequest.Status` 的地方（如适用）
- [ ] 如果有 `project_svc.go` 调用 `xxxRepo.FindByProjectID`，改为 `FindAll(XxxFilter{ProjectID: &id})`
- [ ] 前端：`search` → 实际字段名，去掉 `?all=true`
- [ ] `go build ./...` + `go test ./...` 通过
- [ ] 更新 mock 和测试用例

### 前后台接口隔离

| 端 | 前缀 | 鉴权 |
|----|------|------|
| 前台 | `/api/v1/<resource>` | 无 |
| 后台 | `/api/v1/admin/<resource>` | JWT + RBAC |

禁止交叉调用，禁止共用 handler。

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
