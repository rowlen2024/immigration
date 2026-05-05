# 首页成功案例展示区 — 设计文档

**日期**: 2026-05-05
**目标**: 在首页"精选移民项目"下方新增"成功案例展示区"，复用项目展示区的管理模式

---

## 范围

| 模块 | 管理方式 | 说明 |
|------|---------|------|
| 案例展示区 | 文案 + 精选案例 | section_title/section_subtitle/featured_case_ids |

- 副标题为空时不显示
- 案例卡片复用 `/cases` 页面样式，去掉"成功获批"徽章
- 位置：首页"精选移民项目"和"优势区"之间
- 权限沿用：查看 `admin:read`，编辑 `content:write`

---

## 数据模型

### CaseShowcaseConfig（新增，对标 ProjectShowcaseConfig）

```go
type CaseShowcaseConfig struct {
    SectionTitle    string  `json:"section_title"`
    SectionSubtitle string  `json:"section_subtitle"`
    FeaturedCaseIDs []uint64 `json:"featured_case_ids"`
}
```

- `featured_case_ids` 为空数组或 null 表示展示全部案例（按 sort_order），非空则按列表顺序过滤和排序
- 存储为 `home_configs` 表的 `config_key = "case_showcase"`

### HomeConfigData 新增字段

```go
type HomeConfigData struct {
    // ... 现有字段 ...
    CaseShowcase *CaseShowcaseConfig `json:"case_showcase"`
}
```

---

## 后端改动

| 文件 | 改动 |
|------|------|
| `backend/internal/service/home_config_svc.go` | 新增 `CaseShowcaseConfig` 结构体；`HomeConfigData` 新增 `CaseShowcase` 字段；`Get()` 新增解析 `case_showcase` key |

Handler、Repository、Router 无需改动 — 完全复用现有 `PUT /api/v1/admin/home-config` 的 `Update()` 方法，前端传 `{"case_showcase": {...}}` 即可。

---

## 前端改动

### 管理后台 `admin/homepage.vue`

- `<el-tabs>` 新增第 4 个 tab：`<el-tab-pane label="案例展示区" name="cases">`
- 表单结构对标"项目展示区"：
  - `el-input`：区域标题（placeholder: "成功案例"）
  - `el-input`：区域副标题（placeholder: "数百家庭的信赖之选"）
  - 精选案例列表：从案例下拉选择（调用 `/api/v1/cases`），支持排序/移除
- 独立保存按钮，PUT `{"case_showcase": {...}}`
- 类型定义新增 `CaseShowcase` interface

**案例下拉数据来源**：加载时并行请求 `/admin/home-config` 和 `/cases`（admin 接口返回全部案例）获取可选案例列表。已有 `case_handler.go` 提供 `GET /admin/cases` 路由，返回全部案例。

### 首页 `index.vue`

- 新增 `<section class="section cases-section">` 插入在 `projects-section` 和 `advantages-section` 之间
- `<script>` 中新增：
  - `caseShowcaseConfig` computed：从 homeConfig 读取 `case_showcase`
  - `caseTitle` / `caseSubtitle` computed：回退默认值
  - `featuredCases` computed：根据 `featured_case_ids` 过滤/排序案例
  - 加载案例数据：复用现有 `/api/v1/cases` 接口（已包含 project 关联数据）
- 卡片 HTML/CSS 完全复用 `cases.vue` 的 `.case-card` 样式，仅去掉 `.case-result` 部分
- 副标题为空时 `v-if="caseSubtitle"` 不渲染

卡片模板骨架：

```html
<div class="cases-grid">
  <div v-for="item in featuredCases" :key="item.id" class="case-card">
    <div class="case-image">
      <img :src="item.photo_url" :alt="item.name" loading="lazy" />
    </div>
    <div class="case-body">
      <div class="case-meta">
        <span class="case-country">{{ item.country_from }}</span>
        <span class="case-project">{{ item.project?.name }}</span>
      </div>
      <h3 class="case-name">{{ item.name }}</h3>
      <p class="case-desc">{{ item.description }}</p>
    </div>
  </div>
</div>
```

---

## 改动清单汇总

| 层 | 文件 | 改动 |
|------|------|------|
| 后端 | `backend/internal/service/home_config_svc.go` | 新增 `CaseShowcaseConfig`；`HomeConfigData` 加字段；`Get()` 解析新 key |
| 前端 | `frontend/pages/admin/homepage.vue` | 新增"案例展示区"tab + 表单逻辑 |
| 前端 | `frontend/pages/index.vue` | 新增案例展示 section + 数据逻辑 |

---

## 测试要点

| 场景 | 验证点 |
|------|--------|
| 管理员配置 | 案例展示区 tab 可编辑标题/副标题/精选案例，保存后持久化 |
| 案例选择 | 下拉列出全部案例，选中后显示在列表，可排序/移除 |
| 空精选列表 | `featured_case_ids` 为空时首页展示全部案例 |
| 副标题为空 | 首页不显示副标题行 |
| 首页渲染 | 案例区位于项目区和优势区之间，3 列网格，卡片样式与 /cases 页面一致 |
| 首页卡片 | 卡片无"成功获批"徽章，有照片/来源国/项目/姓名/描述 |
| 权限 | viewer 只读，editor/admin 可编辑保存 |
