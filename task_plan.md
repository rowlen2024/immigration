# Task Plan — MyGo 移民

**更新时间:** 2026-05-20

---

## Phase 5: 接口规范化 ✅ (2026-05-20 完成)

**目标:** 规范前后台 API 使用，确保前后台接口严格隔离、响应格式统一。

### 完成项

- [x] **后台接口补充**: 新增 `GET /admin/cases?all=true`、`GET /admin/testimonials`、`GET /admin/site-config`
- [x] **响应格式规范化**: `?all=true` 返回 `Success` 非分页；分页返回 `SuccessPaginated`
- [x] **公共接口修正**: `GET /faqs` 改为非分页；`GetHomeConfig` 分离公共/admin handler
- [x] **Lawyer 分页**: 新增 repo/service 分页方法，handler 默认分页 + `?all=true`
- [x] **前端修正**: 6 个 admin 页面的 API 调用改为正确 admin 端点 + 响应类型修正
- [x] **规范文档**: CLAUDE.md 新增"API 接口规范"章节（含模板代码 + 接口速查表）

### 修改文件（14 个）

| 层 | 文件 | 说明 |
|----|------|------|
| handler | `faq_handler.go` | ListFAQs → 非分页 |
| handler | `project_handler.go` | AdminListProjects ?all=true → Success |
| handler | `page_handler.go` | AdminListPages ?all=true → Success |
| handler | `case_handler.go` | AdminListCases 新增 ?all=true |
| handler | `lawyer_handler.go` | 分页 + ?all=true |
| handler | `testimonial_handler.go` | 新增 AdminListTestimonials |
| handler | `home_handler.go` | 新增 GetAdminHomeConfig |
| repo | `lawyer_repo.go` | 新增 FindPaginated |
| service | `lawyer_svc.go` | 新增 ListPaginated |
| router | `router.go` | 新增路由 + handler 绑定 |
| frontend | `homepage.vue` | 公共接口 → admin 接口 |
| frontend | `settings.vue` | /site-config → /admin/site-config |
| frontend | `faqs/cases/navigation/projects/lawyers.vue` | ?all=true 响应类型修正 |

---

## Phase 1–4: 前台 UI 优化 ✅ (2026-05-17 完成)
  - [x] 统一 CSS 变量（合并双套体系，更新色值）
  - [x] 引入 Noto Serif SC + Noto Sans SC Web 字体
  - [x] 添加 Type Scale、完整阴影层级、圆角系统、间距系统
  - [x] 按钮改为金色渐变白字 + 悬浮上浮效果

- [x] **Phase 2: 组件提取与优化**
  - [x] 提取 CaseCard 组件（消除 index/cases/projects 三处重复）
  - [x] 提取 TestimonialCarousel 组件（消除 index/projects 两处重复）
  - [x] Header Emoji `▼` → SVG chevron 图标
  - [x] FAQAccordion Emoji `+/−` → SVG 十字图标

- [x] **Phase 3: 核心页面视觉升级**
  - [x] 首页信任条毛玻璃效果 + 金色分隔线
  - [x] 首页项目卡片图片金色渐变叠加 + pill 标签
  - [x] FAQ 筛选按钮 pill 样式
  - [x] 联系表单 focus 金色边框 + 阴影
  - [x] 评价卡片重新设计（暖色渐变 + 金色竖线 + 验证标签）

- [x] **Phase 4: 动效打磨与验证**
  - [x] FAQ 手风琴平滑展开/收起动画
  - [x] 全局过渡时间统一（design tokens）
  - [x] 移动端适配（ContactSidebar 底部条、断点系统）
  - [x] 前端构建验证通过
  - [x] 后端构建验证通过

## 已修改文件清单

| 文件 | 操作 | 说明 |
|------|------|------|
| `frontend/assets/css/variables.css` | 重写 | 统一设计令牌，新增 type scale / 阴影 / 间距 / 动画变量 |
| `frontend/assets/css/global.css` | 更新 | 金色渐变按钮、section-header 更新 |
| `frontend/nuxt.config.ts` | 更新 | 引入 Google Fonts |
| `frontend/components/CaseCard.vue` | **新建** | 统一案例卡片组件 |
| `frontend/components/TestimonialCarousel.vue` | **新建** | 统一评价轮播组件 |
| `frontend/pages/cases.vue` | 更新 | 使用 CaseCard 组件 |
| `frontend/pages/index.vue` | 更新 | CaseCard + TestimonialCarousel + 卡片/信任条样式 |
| `frontend/pages/projects/[slug].vue` | 更新 | CaseCard + TestimonialCarousel |
| `frontend/pages/faq.vue` | 更新 | Pill 筛选按钮 |
| `frontend/pages/contact.vue` | 更新 | 表单焦点金色效果 |
| `frontend/components/global/Header.vue` | 更新 | Emoji → SVG 图标 |
| `frontend/components/project/FAQAccordion.vue` | 更新 | SVG 图标 + 平滑动画 |

## 延迟项（后续可做）

- Hero 动态光晕动效（需显著 JS 改动）
- 暗色模式支持（需 CSS 双主题）
- i18n 国际化（需提取硬编码文本）
