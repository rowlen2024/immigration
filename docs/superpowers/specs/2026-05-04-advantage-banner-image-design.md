# 优势区域 Banner 图 — 设计说明

**日期**: 2026-05-04
**范围**: 首页配置 → 优势管理 → 新增区域图片（选填）

---

## 概述

在优势区域的 `advantage_section` 配置中新增可选 `image` 字段，用于在优势卡片 grid 下方展示一张 banner 图。该图片由管理员在后台"优势管理"tab 中上传/填写，前台条件渲染。

## 数据结构变更

### AdvantageSectionConfig

```go
type AdvantageSectionConfig struct {
    SectionTitle    string `json:"section_title"`
    SectionSubtitle string `json:"section_subtitle"`
    Image           string `json:"image"` // 新增，选填
}
```

存储不变：序列化到 `home_configs` 表 `config_key = "advantage_section"` 行。

### API

| 接口 | 方法 | 变更 |
|------|------|------|
| `/api/v1/home-config` | GET | `advantage_section` 对象返回 `image` 字段 |
| `/api/v1/admin/home-config` | PUT | `advantage_section.image` 可写入 |

现有接口自动支持，无需新增路由或 handler。

## 改动文件

### 1. `backend/internal/service/home_config_svc.go`

`AdvantageSectionConfig` struct 加 `Image string` 字段。

### 2. `frontend/pages/admin/homepage.vue`

"优势管理" tab 表单中，在"区域标题"和"区域副标题"下方新增 `ImageInput` 组件：

- label: "区域图片"
- 选填，使用已有 `ImageInput` 组件
- 绑定到 `form.advantage_section.image`

### 3. `frontend/pages/index.vue`

在 `advantages-grid` div 下方条件渲染 banner 图：

```html
<img
  v-if="homeConfig?.advantage_section?.image"
  :src="homeConfig.advantage_section.image"
  alt="优势区域图"
  class="advantage-banner"
/>
```

CSS：图片宽度撑满容器，响应式。
