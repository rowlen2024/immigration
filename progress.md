# Progress Log — MyGo 移民

**关联计划:** task_plan.md

---

## Session 2 — 2026-05-20 — 接口规范化

### 分析阶段
- 通读 `router.go` 全部路由定义（~150 行）
- 通读 13 个 handler 文件，识别响应格式
- 扫描前台 8 个页面 + 后台 13 个页面的所有 API 调用
- 交叉比对前后端，找出 5 类问题

### 实施阶段 — 后端
- `handler/faq_handler.go`: `ListFAQs` → 非分页 `Success`
- `handler/project_handler.go`: `AdminListProjects` `?all=true` → `Success`
- `handler/page_handler.go`: `AdminListPages` `?all=true` → `Success`
- `handler/case_handler.go`: `AdminListCases` 新增 `?all=true` 支持
- `handler/lawyer_handler.go`: 默认分页 + `?all=true` 全量
- `repository/lawyer_repo.go`: 新增 `FindPaginated`
- `service/lawyer_svc.go`: 新增 `ListPaginated`
- `handler/testimonial_handler.go`: 新增 `AdminListTestimonials`（支持 `?project_id=`）
- `handler/home_handler.go`: 新增 `GetAdminHomeConfig`，与公共 `GetHomeConfig` 分离
- `router/router.go`: 新增 `GET /admin/testimonials`、`GET /admin/site-config`，绑定分离的 handler

### 实施阶段 — 前端
- `homepage.vue`: 3 处公共接口 → admin 接口 + 响应类型修正
- `settings.vue`: `GET /site-config` → `GET /admin/site-config`
- `faqs.vue` / `cases.vue` / `navigation.vue` / `projects.vue` / `lawyers.vue`: `?all=true` 响应类型修正

### 验证
- `go build ./...` 后端编译通过
- `npx nuxi typecheck` 前端无新增错误（现有错误均为 Nuxt auto-import 预存问题）

### 产出物
| 文件 | 操作 | 说明 |
|------|------|------|
| CLAUDE.md | 更新 | 新增"API 接口规范"章节（含速查表） |
| findings.md | 更新 | 记录 5 类接口问题 |
| 后端 8 文件 | 修改 | handler/repo/service/router |
| 前端 6 文件 | 修改 | admin 页面 API 调用修正 |

---

## Session 1 — 2026-05-17

### 分析阶段
- 完成前台代码全面探索（14 个页面 + 组件）
- 识别 11 个 UI Bug（3 严重 + 5 中等 + 3 轻微）
- 运行 UI/UX Pro Max 设计系统搜索，获取 Trust & Authority 风格推荐
- 输出《前端UI优化方案.md》含三个可选方案
- 用户选择方案 B，输出《前端UI设计稿.md》含 12 节详细设计稿
- 生成 `design-preview.html` 效果预览文件供用户审核

### 实施阶段 — Phase 1
- 重写 `variables.css`：统一 16 组设计令牌
- 更新 `global.css`：金色渐变按钮 + section-header 更新
- 更新 `nuxt.config.ts`：引入 Noto Serif SC + Noto Sans SC
- 验证：types 生成通过

### 实施阶段 — Phase 2
- 创建 `CaseCard.vue` 统一组件
- 更新 `cases.vue`、`index.vue`、`projects/[slug].vue` 使用 CaseCard
- 更新 `Header.vue`：Emoji → SVG chevron ×4 处
- 验证：前端构建通过（exit 0）

### 实施阶段 — Phase 3
- FAQ 筛选按钮改为 pill 样式
- 联系表单 focus 金色边框 + 阴影
- 首页信任条毛玻璃 + 金色分隔线
- 首页项目卡片金色渐变叠加 + pill 标签
- 验证：前端 + 后端构建均通过

### 实施阶段 — Phase 4
- FAQ 手风琴 SVG 十字图标 + 平滑展开动画
- ContactSidebar 移动端底部固定条
- 提取 TestimonialCarousel 组件（消除 130+ 行重复）
- 评价卡片重新设计（去掉引号，暖色渐变 + 金色点缀 + 验证标签）
- 验证：前后端构建均通过（exit 0）

### 清理
- 确认 0 个 `&#9662;` Emoji 残留
- 确认 0 处 `tm-quote` 旧样式残留

### 产出物
| 文件 | 类型 |
|------|------|
| `前端UI优化方案.md` | 文档 |
| `前端UI设计稿.md` | 文档 |
| `design-preview.html` | 预览 |
| `testimonial-preview.html` | 预览 |
