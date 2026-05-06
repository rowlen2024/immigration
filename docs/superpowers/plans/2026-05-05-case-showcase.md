# 首页成功案例展示区 实现计划

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 在首页"精选移民项目"下方新增"成功案例展示区"，管理员可在后台配置标题/副标题/精选案例。

**Architecture:** 完全复用 `project_showcase` 模式 — 后端新增 `case_showcase` config key（`CaseShowcaseConfig` 结构体），前端 admin 新增 tab，首页新增 section。卡片样式复用 `/cases` 页面，去掉"成功获批"徽章。

**Tech Stack:** Go (Gin + GORM), Nuxt 3 (Vue 3 + Element Plus)

---

### Task 1: Backend — 新增 CaseShowcaseConfig 和 HomeConfigData 字段

**Files:**
- Modify: `immigration/backend/internal/service/home_config_svc.go`

- [ ] **Step 1: 新增 CaseShowcaseConfig 结构体**

在 `ProjectShowcaseConfig` 之后添加：

```go
// CaseShowcaseConfig holds the case showcase section settings.
type CaseShowcaseConfig struct {
	SectionTitle    string   `json:"section_title"`
	SectionSubtitle string   `json:"section_subtitle"`
	FeaturedCaseIDs []uint64 `json:"featured_case_ids"`
}
```

- [ ] **Step 2: HomeConfigData 新增 CaseShowcase 字段**

在 `HomeConfigData` 结构体的 `ProjectShowcase` 字段后添加：

```go
CaseShowcase *CaseShowcaseConfig `json:"case_showcase"`
```

- [ ] **Step 3: Get() 方法新增 case_showcase key 解析**

在 `Get()` 方法中，`project_showcase` 解析代码块之后添加：

```go
if caseCfg, err := s.repo.FindByKey("case_showcase"); err == nil {
	var csc CaseShowcaseConfig
	if err := json.Unmarshal(caseCfg.ConfigValue, &csc); err == nil {
		data.CaseShowcase = &csc
	}
}
```

- [ ] **Step 4: 编译验证**

```bash
cd immigration/backend && go build ./...
```

Expected: 编译成功，无错误。

---

### Task 2: Frontend Admin — 新增"案例展示区"tab

**Files:**
- Modify: `immigration/frontend/pages/admin/homepage.vue`

- [ ] **Step 1: 新增 CaseShowcase 接口和响应式状态**

在 `<script setup>` 中 `interface ProjectShowcase` 后面添加接口：

```typescript
interface CaseShowcase {
  section_title: string;
  section_subtitle: string;
  featured_case_ids: number[];
}
```

在 `const projectShowcase` 后面添加响应式状态和保存 loading：

```typescript
const caseShowcase = ref<CaseShowcase>({
  section_title: '',
  section_subtitle: '',
  featured_case_ids: [],
});

interface CaseOption {
  id: number;
  name: string;
}

const allCases = ref<CaseOption[]>([]);
const caseSaving = ref(false);
```

- [ ] **Step 2: load() 中加载案例数据和 case_showcase 配置**

在 `load()` 函数中，将 `Promise.all` 改为同时请求三个 API。修改 `load` 函数中的 API 调用部分：

在现有 `Promise.all` 调用处，将 `[config, projects]` 改为 `[config, projects, casesData]`，并新增 cases 请求和解析：

```typescript
const api = useApi();
const [config, projects, casesData] = await Promise.all([
  api<{
    hero_slides: HeroSlide[];
    advantage_items: AdvantageItem[];
    advantage_section: { section_title: string; section_subtitle: string; image: string } | null;
    project_showcase: ProjectShowcase | null;
    case_showcase: CaseShowcase | null;
  }>('/admin/home-config'),
  api<{ items: ProjectOption[] }>('/projects'),
  api<{ items: CaseOption[] } | CaseOption[]>('/cases'),
]);

// ... existing config parsing ...

if (config.case_showcase) {
  caseShowcase.value = config.case_showcase;
}

// Parse cases
if (casesData) {
  const items = Array.isArray(casesData) ? casesData : (casesData as { items: CaseOption[] }).items;
  if (items) {
    allCases.value = items;
  }
}
```

- [ ] **Step 3: 新增辅助方法和 computed**

在 `availableProjects` computed 后面添加：

```typescript
const availableCases = computed(() => {
  const featured = new Set(caseShowcase.value.featured_case_ids);
  return allCases.value.filter((c) => !featured.has(c.id));
});

function getCaseTitle(id: number): string {
  return allCases.value.find((c) => c.id === id)?.name || String(id);
}

function moveCaseFeatured(index: number, direction: -1 | 1) {
  const target = index + direction;
  if (target < 0 || target >= caseShowcase.value.featured_case_ids.length) return;
  const ids = [...caseShowcase.value.featured_case_ids];
  [ids[index], ids[target]] = [ids[target], ids[index]];
  caseShowcase.value.featured_case_ids = ids;
}

function removeCaseFeatured(index: number) {
  caseShowcase.value.featured_case_ids.splice(index, 1);
}

function addCaseFeatured(id: number) {
  if (!caseShowcase.value.featured_case_ids.includes(id)) {
    caseShowcase.value.featured_case_ids.push(id);
  }
}

async function saveCaseShowcase() {
  caseSaving.value = true;
  try {
    const api = useApi();
    await api('/admin/home-config', {
      method: 'PUT',
      body: { case_showcase: caseShowcase.value },
    });
    ElMessage.success('案例展示区已保存');
  } catch {
    ElMessage.error('保存失败');
  } finally {
    caseSaving.value = false;
  }
}
```

- [ ] **Step 4: 新增"案例展示区"tab 模板**

在 `<el-tab-pane label="优势管理" name="advantages">` 结束标签 `</el-tab-pane>` 之后、`</el-tabs>` 之前添加：

```html
<el-tab-pane label="案例展示区" name="cases">
  <el-card class="config-card">
    <template #header><h3 class="admin-card-title">案例展示区</h3></template>
    <el-form label-width="100px">
      <el-form-item label="区域标题">
        <el-input v-model="caseShowcase.section_title" placeholder="成功案例" />
      </el-form-item>
      <el-form-item label="区域副标题">
        <el-input v-model="caseShowcase.section_subtitle" placeholder="数百家庭的信赖之选" />
      </el-form-item>
      <el-form-item label="精选案例">
        <div class="featured-area">
          <div v-if="caseShowcase.featured_case_ids.length === 0" class="admin-empty-hint">
            未选择精选案例，首页将展示全部案例。
          </div>
          <div v-else class="config-list">
            <div v-for="(id, i) in caseShowcase.featured_case_ids" :key="id" class="config-item">
              <span class="config-item-name">{{ getCaseTitle(id) }}</span>
              <div class="config-item-actions">
                <button class="action-btn" :disabled="i === 0" @click="moveCaseFeatured(i, -1)">↑</button>
                <button class="action-btn" :disabled="i === caseShowcase.featured_case_ids.length - 1" @click="moveCaseFeatured(i, 1)">↓</button>
                <button class="action-btn danger" @click="removeCaseFeatured(i)">移除</button>
              </div>
            </div>
          </div>
          <el-select
            v-if="availableCases.length > 0"
            value=""
            placeholder="添加案例..."
            clearable
            @change="(val: number) => { if (val) addCaseFeatured(val) }"
            class="add-project-select"
          >
            <el-option
              v-for="c in availableCases"
              :key="c.id"
              :label="c.name"
              :value="c.id"
            />
          </el-select>
        </div>
      </el-form-item>
    </el-form>
    <div class="card-footer">
      <el-button type="primary" :loading="caseSaving" @click="saveCaseShowcase">保存</el-button>
    </div>
  </el-card>
</el-tab-pane>
```

- [ ] **Step 5: 类型检查**

```bash
cd immigration/frontend && npx nuxi typecheck 2>&1 | tail -20
```

Expected: 无新增类型错误。

---

### Task 3: Frontend Index — 新增案例展示 section

**Files:**
- Modify: `immigration/frontend/pages/index.vue`

- [ ] **Step 1: 扩展 homeConfig 类型映射**

在 `<script setup>` 中，`showcaseConfig` computed 之后，`advantageSection` computed 之前，新增 `caseShowcaseConfig` computed：

```typescript
const caseShowcaseConfig = computed(() => {
  if (homeConfig.value) {
    const config = homeConfig.value as unknown as Record<string, unknown>;
    const data = config.data as Record<string, unknown> | undefined;
    if (data && data.case_showcase) {
      return data.case_showcase as {
        section_title?: string;
        section_subtitle?: string;
        featured_case_ids?: number[];
      };
    }
  }
  return null;
});

const caseTitle = computed(() => caseShowcaseConfig.value?.section_title || '成功案例');
const caseSubtitle = computed(() => caseShowcaseConfig.value?.section_subtitle || '');
```

- [ ] **Step 2: 新增案例数据获取和 featuredCases computed**

在 `advantageSection` computed 之前（即 Step 1 代码之后），添加案例数据类型和请求：

```typescript
interface CaseItem {
  id: number;
  name: string;
  country_from: string;
  photo_url: string;
  description: string;
  project?: { name: string };
}

const { data: casesData, pending: pendingCases, error: errorCasesRaw } = await useFetch<{
  data?: CaseItem[];
}>('/api/v1/cases', {
  onResponseError() {},
});

const featuredCases = computed<CaseItem[]>(() => {
  const apiData = casesData.value as { data?: CaseItem[] } | null;
  const all = apiData?.data ?? [];
  if (all.length === 0) return [];

  const featured = caseShowcaseConfig.value?.featured_case_ids;
  if (featured && featured.length > 0) {
    const orderMap = new Map(featured.map((id: number, i: number) => [id, i]));
    return all
      .filter((c) => orderMap.has(c.id))
      .sort((a, b) => {
        const ai = orderMap.get(a.id);
        const bi = orderMap.get(b.id);
        if (ai !== undefined && bi !== undefined) return ai - bi;
        return 0;
      });
  }

  return all;
});
```

更新 `pending` computed，添加 cases：

```typescript
const pending = computed(() => ({
  projects: pendingProjects.value,
  cases: pendingCases.value,
}));
```

更新 `error` computed，添加 cases：

```typescript
const error = computed(() => ({
  projects: errorProjectsRaw.value ? '加载失败，请刷新重试' : null,
  cases: errorCasesRaw.value ? '加载失败，请刷新重试' : null,
}));
```

- [ ] **Step 3: 新增案例展示 section 模板**

在 `<!-- Advantages Section -->` 的 `</section>` 之后，`<!-- CTA Banner -->` 之前，插入：

```html
<!-- Cases Section -->
<section class="section cases-section">
  <div class="container">
    <div class="section-header">
      <h2>{{ caseTitle }}</h2>
      <p v-if="caseSubtitle">{{ caseSubtitle }}</p>
    </div>

    <div v-if="pending.cases" class="loading-state">
      <div class="skeleton" style="height:360px;"></div>
    </div>
    <div v-else-if="error.cases" class="error-state">
      <div class="error-card">
        <span v-html="getIconSvg('alert-circle', 24, '#c0392b')"></span>
        <p>加载失败，请刷新重试</p>
      </div>
    </div>
    <div v-else class="cases-grid">
      <div v-for="item in featuredCases" :key="item.id" class="case-card reveal">
        <div class="case-image">
          <img
            :src="item.photo_url || ''"
            :alt="item.name"
            loading="lazy"
          />
        </div>
        <div class="case-body">
          <div class="case-meta">
            <span class="case-country">{{ item.country_from }}</span>
            <span v-if="item.project?.name" class="case-project">{{ item.project.name }}</span>
          </div>
          <h3 class="case-name">{{ item.name }}</h3>
          <p class="case-desc">{{ item.description }}</p>
        </div>
      </div>
    </div>

    <div v-if="!pending.cases && featuredCases.length === 0" class="empty-state">
      暂无成功案例
    </div>
  </div>
</section>
```

- [ ] **Step 4: 新增案例展示区 CSS**

在 `<style scoped>` 中，CTA Section 的 CSS 之前，添加：

```css
/* Cases Section */
.cases-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 32px;
  margin-bottom: 48px;
}

.case-card {
  background-color: var(--bg-white);
  border-radius: var(--radius-lg);
  overflow: hidden;
  box-shadow: var(--shadow-sm);
  transition: box-shadow 0.3s ease, transform 0.3s ease;
}

.case-card:hover {
  box-shadow: var(--shadow-md);
  transform: translateY(-4px);
}

.case-image {
  height: 200px;
  overflow: hidden;
}

.case-image img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: transform 0.3s ease;
}

.case-card:hover .case-image img {
  transform: scale(1.05);
}

.case-body {
  padding: 24px;
}

.case-meta {
  display: flex;
  gap: 8px;
  margin-bottom: 12px;
}

.case-country,
.case-project {
  font-size: 12px;
  font-weight: 600;
  padding: 2px 10px;
  border-radius: var(--radius-sm);
}

.case-country {
  background-color: rgba(26, 58, 92, 0.1);
  color: var(--primary);
}

.case-project {
  background-color: rgba(200, 150, 62, 0.1);
  color: var(--accent-dark);
}

.case-name {
  font-size: 18px;
  font-weight: 700;
  color: var(--text-primary);
  margin-bottom: 8px;
}

.case-desc {
  font-size: 14px;
  color: var(--text-secondary);
  line-height: 1.7;
  margin-bottom: 0;
}

.empty-state {
  text-align: center;
  padding: 60px 20px;
  color: var(--text-light);
  font-size: 15px;
}
```

在响应式媒体查询中添加 cases 的断点。在 `@media (max-width: 1023px)` 块内添加：

```css
.cases-grid {
  grid-template-columns: repeat(2, 1fr);
}
```

在 `@media (max-width: 767px)` 块内添加：

```css
.cases-grid {
  grid-template-columns: 1fr;
}
```

注意：现有的 `@media (max-width: 1023px)` 和 `@media (max-width: 767px)` 已经存在，需要在对应花括号内追加 `.cases-grid` 规则，不要创建重复的媒体查询块。

- [ ] **Step 5: 类型检查**

```bash
cd immigration/frontend && npx nuxi typecheck 2>&1 | tail -20
```

Expected: 无新增类型错误。

---

### Task 4: 验证

- [ ] **Step 1: 后端测试**

```bash
cd immigration/backend && go test ./...
```

Expected: 全部测试通过。

- [ ] **Step 2: 前端类型检查**

```bash
cd immigration/frontend && npx nuxi typecheck 2>&1 | tail -20
```

Expected: 无类型错误。

- [ ] **Step 3: 手动验证清单**

| 验证项 | 操作 |
|--------|------|
| 管理后台 | 打开 `/admin/homepage`，确认第 4 个 tab"案例展示区"可见 |
| 编辑标题/副标题 | 填写标题和副标题，点击保存，刷新页面确认持久化 |
| 选择精选案例 | 从下拉选择案例，排序/移除，保存后确认 |
| 首页渲染 | 打开首页，确认案例区在项目区和优势区之间 |
| 首页卡片 | 确认卡片有照片/来源国/项目/姓名/描述，无"成功获批"徽章 |
| 副标题为空 | 副标题清空保存后，首页不显示副标题行 |
