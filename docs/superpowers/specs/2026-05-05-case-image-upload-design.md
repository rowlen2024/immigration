# 案例管理图片上传及字段补全

## 背景

公开成功案例页面（`/cases`）图片无法显示，根因是后台案例管理缺少图片上传功能。此外后台表单仅覆盖了部分 Case 字段，需要补全。

## 后端 Model

`model/case.go` 已包含所有字段（含 `photo_url`），无需修改。

## 后端 Repo

`repository/case_repo.go` — `FindAll()` 查询添加 `Preload("Project")`，案例列表返回时携带 `project` 名称。

## 前端 Admin 页面 (`pages/admin/cases.vue`)

### 表单新增 4 个字段

| 字段 | 组件 | 对应模型字段 | 说明 |
|------|------|------------|------|
| 封面图片 | `ImageInput` | `photo_url` | 复用已有组件，支持 URL 输入/上传/媒体库选择 |
| 投资金额 | `el-input` | `investment_amount` | 文本型，如"80万美元" |
| 投资数额 | `el-input-number` | `investment_value` | 数字型 |
| 办理周期 | `el-input` | `processing_period` | 文本型，如"28个月" |

### 表格列补充

在现有列基础上增加 `photo_url`（缩略图预览）列，便于查看。

### TypeScript 接口

`CaseItem` 接口补充缺失字段，`defaultForm()` 和 `openEdit()` 也同步更新。

## 前端 Public 页面 (`pages/cases.vue`)

### 修正字段映射

- `item.image` → 从 API 的 `item.photo_url` 映射
- `item.country` → 从 API 的 `item.country_from` 映射
- `item.project` → 从 API 的 `item.project`（Preload 后的关联对象）

### 移除硬编码 fallback

删除 70-125 行的硬编码示例数据，API 返回空数组时直接展示空状态。

### "成功获批"标签

保持固定文本，不新增数据库字段。

## 未改动的文件

- `handler/case_handler.go` — JSON 绑定自动映射新字段
- `service/case_svc.go` — 透传无业务逻辑变更
- 路由和权限 — 无需调整
- 数据库和迁移 — 字段已存在无需新增

## 图片展示规格

公开案例页 `.case-image` 为 3 列网格中的卡片图（高度 200px，`object-fit: cover`），建议上传 800×600 左右的横向图片，控制文件大小。
