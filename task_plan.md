# Task Plan — MyGo 移民前台 UI 优化

**创建时间:** 2026-05-17  
**方案:** 方案 B — 品牌视觉升级  
**状态:** 已完成

## 目标

在不改变现有功能的前提下，优化官网前台 UI，提升美观度、视觉层次和推广转化效果。

## 阶段进度

- [x] **Phase 1: 设计令牌基础设施**
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
