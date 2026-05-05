# 图片上传功能统一改造设计

## 背景

当前系统中多个管理页面需要填写图片地址，但大部分只有纯文本输入框，缺乏上传功能。需要统一改造，支持上传 + 媒体库选择。

## 现状分析

| 页面 | 字段 | 当前状态 |
|------|------|---------|
| `admin/homepage.vue` | 轮播背景图 | 纯文本输入 |
| `admin/projects.vue` | 封面图片 | 已有内联上传，但实现与页面耦合 |
| `admin/settings.vue` | Logo / Favicon / OG图 / 机构Logo / 机构URL | 纯文本输入（5 个字段） |
| `admin/media.vue` | 媒体库管理 | 独立页面，无法在其他表单中引用；文件大小显示 NaN（字段名不匹配），缺少 URL 显示 |

后端已有 `POST /api/v1/admin/media/upload` 和 `GET /api/v1/admin/media`，但搜索支持需确认。

## 交互方案

采用**上传 + 浏览结合的方案**。每个图片字段旁提供两个按钮：

- **上传**：直接选文件上传，自动填入 URL
- **浏览**：打开媒体库对话框，浏览/搜索已有图片并选择

## 新增组件

### ImageInput.vue

可复用的图片输入组件，替代 `<el-input>` 纯文本输入。

Props:
- `modelValue: string` — v-model 绑定的 URL
- `placeholder?: string` — 输入框占位文字

交互：
- 输入框可手动填写 URL
- 「上传」按钮用 `el-upload` 上传到 `/api/v1/admin/media/upload`，成功后自动填入 URL
- 「浏览」按钮打开 `MediaPicker` 对话框
- 有 URL 时显示图片预览缩略图

### MediaPicker.vue

媒体库选择对话框（基于 `el-dialog`）。

左侧网格（3-4 列）：
- 调用 `GET /api/v1/admin/media?page=&perPage=&search=` 获取媒体列表
- 分页加载

右侧详情面板：
- 选中图片的大图预览
- 文件名、尺寸、大小信息

顶部工具栏：
- 搜索框：按文件名模糊搜索
- 上传按钮：可在此对话框直接上传新图片

底部操作：
- 取消：关闭对话框，不改变当前值
- 确认选择：将选中图片的 URL 填入输入框

## 页面修改

### admin/homepage.vue

轮播幻灯片的「背景图 URL」字段从 `el-input` 替换为 `ImageInput`。

### admin/projects.vue

封面图片字段从现有的内联上传实现替换为 `ImageInput`，消除重复代码。

### admin/settings.vue

以下字段从 `el-input` 替换为 `ImageInput`：

- `site_logo` — 网站 Logo
- `site_favicon` — Favicon
- `og_image` — OG 分享图
- `organization_logo` — 机构 Logo
- `organization_url` — 机构官网 URL（非图片字段，保持原样不变）

共 4 个图片字段需要改造（organization_url 不涉及图片）。

### admin/media.vue 修复

修复两个问题：

1. **文件大小 NaN**：前端接口字段 `size` 与后端返回的 `size_bytes` 不匹配，改为 `size_bytes`
2. **缺少文件路径**：在媒体卡片信息区增加 URL 显示，方便复制使用

## 后端改动

- `GET /api/v1/admin/media` 确认/补充 `search` 查询参数支持（按 `original_name` 或 `filename` 模糊搜索）

## 测试计划

- 手动测试所有页面的图片上传、浏览选择、手动输入三种方式
- 确认 media 表数据正确
- 确认已有图片显示预览缩略图
