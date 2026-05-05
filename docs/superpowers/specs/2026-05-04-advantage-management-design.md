# 优势管理增强 — 设计文档

**日期**: 2026-05-04  
**状态**: 已确认

## 概述

增强首页优势管理功能：新增区域标题/副标题配置、优势项图标支持 SVG（Lucide）选择、统一管理后台三卡片 UI 风格。

## 一、数据模型

### 后端结构体变更 (`backend/internal/service/home_config_svc.go`)

**AdvantageItem** — 新增 `IconType` 字段：
```go
type AdvantageItem struct {
    Icon        string `json:"icon"`        // Lucide 图标名，如 "search"
    IconType    string `json:"icon_type"`   // 固定值 "lucide"，旧数据缺失时 fallback 为 emoji
    Title       string `json:"title"`
    Description string `json:"description"`
}
```

**AdvantageSectionConfig** — 新增结构体：
```go
type AdvantageSectionConfig struct {
    SectionTitle    string          `json:"section_title"`
    SectionSubtitle string          `json:"section_subtitle"`
    Items           []AdvantageItem `json:"items"`
}
```

### 存储

沿用 `home_configs` 表（key-value JSON），拆分为两条 config_key：
- `advantage_section` — 存 `{"section_title":"...", "section_subtitle":"..."}` 
- `advantage_items` — 存 `[{"icon":"search", "icon_type":"lucide", "title":"...", "description":"..."}, ...]`

**向后兼容**：前端渲染时若 `icon_type !== 'lucide'`（旧数据），降级为文本渲染 emoji。

### HomeConfigData 更新

```go
type HomeConfigData struct {
    HeroSlides         []HeroSlide              `json:"hero_slides"`
    AdvantageItems     []AdvantageItem          `json:"advantage_items"`
    AdvantageSection   *AdvantageSectionConfig  `json:"advantage_section"`  // 新增
    ProjectShowcase    *ProjectShowcaseConfig   `json:"project_showcase"`
}
```

## 二、前端图标系统

### 图标数据 (`frontend/composables/lucideIcons.ts`)

- 100 个精选 Lucide 图标，每个存储 SVG path 数据
- 分类：全部、商务、金融、法律、教育、通信、通用
- 总大小约 20KB（gzip ~5KB），不引入外部依赖

### 图标组件 (`frontend/components/admin/IconPicker.vue`)

- 弹窗内 8 列图标网格
- 分类 tabs 筛选 + 搜索框按名称过滤
- 选中高亮（金色边框 + 浅金背景）
- 底部显示当前选中图标名

### 首页渲染 (`frontend/pages/index.vue`)

- 图标渲染：`<SvgIcon :name="adv.icon" :size="28" />` 替换 `{{ adv.icon }}`
- 样式对标设计稿：56px 圆形容器 + `rgba(26,58,92,0.08)` 背景 + `#c8963e` 金色
- 区域标题/副标题从 API (`advantage_section`) 读取，替换硬编码文案
- Fallback：`icon_type !== 'lucide'` 时渲染纯文本 emoji

## 三、管理后台

### 优势管理卡片 (`frontend/pages/admin/homepage.vue`)

新增：
- 区域标题输入框
- 区域副标题输入框
- 图标选择器替换纯文本输入框（`<IconPicker v-model="advForm.icon" />`）
- 列表中图标以圆形 SVG 预览（36px 容器）
- 保存时同时写入 `advantage_section` + `advantage_items`

### 三卡片统一 UI 风格

轮播管理、项目展示区、优势管理采用一致规范：
- **卡片结构**：灰色表头（仅标题，无操作按钮）+ 白色内容区
- **区域设置**：标题/副标题字段位于列表上方
- **列表项**：统一边框圆角、左内容右操作布局
- **新增按钮**：列表右下角，右对齐
- **保存按钮**：底部居中，上方分隔线区隔

轮播管理卡片的"新增 Slide"按钮从表头移至内容区右下角。

## 四、改动文件清单

| 层 | 文件 | 操作 |
|----|------|------|
| 后端 service | `backend/internal/service/home_config_svc.go` | 改 — 新增结构体，更新 Get() |
| 前端首页 | `frontend/pages/index.vue` | 改 — 标题/副标题 API 化，图标 SVG 化 |
| 前端管理 | `frontend/pages/admin/homepage.vue` | 改 — 三卡片 UI 统一，优势管理增强 |
| 前端新增 | `frontend/components/admin/IconPicker.vue` | 新建 — 图标选择器组件 |
| 前端新增 | `frontend/composables/lucideIcons.ts` | 新建 — 100 个 SVG 图标数据 |
