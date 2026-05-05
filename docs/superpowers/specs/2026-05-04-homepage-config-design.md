# 首页配置管理 — 设计文档

**日期**: 2026-05-04
**目标**: 在管理后台新增"首页配置"页面，统一管理首页轮播、项目展示区、优势区

---

## 范围

| 模块 | 管理方式 | 说明 |
|------|---------|------|
| 轮播 Slide | 增删改排序 | title/desc/project_slug/gradient/image/link |
| 项目展示区 | 文案 + 精选项目 | section_title/section_subtitle/featured_slugs |
| 优势卡片 | 增删改排序 | icon/title/description |
| CTA 区 | 不管 | 保持硬编码 |

权限：查看 `admin:read`，编辑 `content:write`（沿用现有路由权限）。

---

## 数据模型

三组数据各存为 `home_configs` 表的一个 config_key。

### HeroSlide（更新）

```go
type HeroSlide struct {
    Title       string `json:"title"`
    Desc        string `json:"desc"`
    ProjectSlug string `json:"project_slug"`
    Gradient    string `json:"gradient"`
    Image       string `json:"image"`
    Link        string `json:"link"`
}
```

比现有多出 `Desc`、`ProjectSlug`、`Gradient` 字段，与种子数据对齐；新增 `Link`。`ImageURL` 重命名为 `Image` 与种子数据一致。

### AdvantageItem（不变）

```go
type AdvantageItem struct {
    Icon        string `json:"icon"`
    Title       string `json:"title"`
    Description string `json:"description"`
}
```

### ProjectShowcaseConfig（新增）

```go
type ProjectShowcaseConfig struct {
    SectionTitle    string   `json:"section_title"`
    SectionSubtitle string   `json:"section_subtitle"`
    FeaturedSlugs   []string `json:"featured_slugs"`
}
```

`FeaturedSlugs` 为空数组表示展示全部项目（按默认顺序），非空则按列表顺序过滤和排序。

### HomeConfigData（更新）

```go
type HomeConfigData struct {
    HeroSlides      []HeroSlide            `json:"hero_slides"`
    AdvantageItems  []AdvantageItem        `json:"advantage_items"`
    ProjectShowcase ProjectShowcaseConfig  `json:"project_showcase"`
}
```

---

## API

完全复用现有路由，无需新增。

| 方法 | 路由 | 权限 | 说明 |
|------|------|------|------|
| GET | `/api/v1/home-config` | 无 | 公开获取首页配置 |
| GET | `/api/v1/admin/home-config` | `admin:read` | 管理端获取 |
| PUT | `/api/v1/admin/home-config` | `content:write` | 管理端保存 |

后端改动仅限 service 层：结构体字段更新、`Get()` 方法新增 `project_showcase` key 的解析。Handler/repository 无需改动。

---

## 前端

### 路由与菜单

- 新增 Nuxt 路由页 `/admin/homepage`
- 侧边栏 `admin.vue` 新增"首页配置"菜单项（active-class 模式），放在"页面管理"和"导航管理"之间
- 仪表盘 `index.vue` 快捷入口也补上

### 页面结构（方案一：单页多卡片）

三张 `<el-card>`，各自独立加载/保存：

**轮播管理 Card**
- Slide 列表，每项显示缩略标题，右侧操作按钮（上下排序/编辑/删除）
- "新增 Slide" 按钮
- 编辑/新增用 `<el-dialog>` 弹窗表单（title 必填校验）
- 排序用上下箭头交换位置

**项目展示区 Card**
- `el-input`：区域标题、副标题
- 精选项目：从项目列表（GET /projects）多选 + 上下箭头排序
- 单独保存按钮

**优势管理 Card**
- 优势项列表，每项显示 icon/标题/描述，右侧操作按钮
- 新增/编辑用 `<el-dialog>` 弹窗或行内表单
- 排序用上下箭头

### 数据流

```
加载：GET /admin/home-config → 解包三组数据 → 分别注入各 card 的响应式状态

保存：PUT /admin/home-config { "hero_slides": [...] } → 成功 toast → 失败 toast + 回滚状态
```

三张 card 独立保存。保存只发当前 card 对应的一组数据，不全量覆盖。

### 首页联动

`index.vue` 调整：
- hero_slides：字段映射 `image_url` → `image`，新增读取 `link`
- 项目展示区：读取 `project_showcase` 中的 section 标题/副标题，按 `featured_slugs` 过滤排序项目
- 优势区：字段映射不变

---

## 改动清单

### 后端

| 文件 | 改动 |
|------|------|
| `backend/internal/service/home_config_svc.go` | HeroSlide 结构体更新字段；新增 ProjectShowcaseConfig 结构体；HomeConfigData 新增 ProjectShowcase 字段；Get() 新增解析 project_showcase key |

### 前端

| 文件 | 改动 |
|------|------|
| `frontend/pages/admin/homepage.vue` | **新建** — 首页配置管理页面 |
| `frontend/layouts/admin.vue` | 侧边栏新增"首页配置"菜单项 |
| `frontend/pages/admin/index.vue` | 快捷入口新增"首页配置" |
| `frontend/pages/index.vue` | hero_slides 字段映射更新；项目展示区读取 project_showcase 配置 |

---

## 测试要点

| 场景 | 验证点 |
|------|--------|
| 新增 slide | 弹窗必填校验，保存后列表更新 |
| 编辑 slide | 弹窗回显，修改后保存 |
| 删除 slide | 确认后移除，至少保留 1 张 |
| 排序 slide | 上下箭头交换位置，保存后持久化 |
| 项目展示区 | 标题/副标题编辑，精选项目多选排序 |
| 优势项 CRUD | 新增/编辑/删除/排序 |
| 权限 | viewer 只读，editor/admin 可写 |
| 首页联动 | 管理端保存后，首页正确渲染新数据 |
