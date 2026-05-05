# 首页配置管理 — Implementation Plan

> **For agentic workers:** Use superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 在管理后台新增"首页配置"页面，统一管理首页轮播 Slide、项目展示区、优势卡片。

**Architecture:** 后端更新 `HomeConfigService` 的 HeroSlide 结构体（对齐种子数据字段）并新增 `ProjectShowcaseConfig`。前端新建 `/admin/homepage` 页（三张 el-card 布局），侧边栏新增菜单项。完全复用现有 API 路由，权限沿用 `admin:read` / `content:write`。

**Tech Stack:** Go + Gin + GORM（后端）, Nuxt 3 + Element Plus + TypeScript（前端）

---

### Task 1: Backend — 更新 service 层结构体和 Get()

**Files:**
- Modify: `backend/internal/service/home_config_svc.go`

- [ ] **Step 1: 更新 HeroSlide 结构体**

当前代码约第 17-22 行，替换 HeroSlide：

```go
// HeroSlide represents a slide in the hero section of the homepage.
type HeroSlide struct {
	Title       string `json:"title"`
	Desc        string `json:"desc"`
	ProjectSlug string `json:"project_slug"`
	Gradient    string `json:"gradient"`
	Image       string `json:"image"`
	Link        string `json:"link"`
}
```

变化：`Subtitle` → `Desc`，`ImageURL` → `Image`，新增 `ProjectSlug`、`Gradient`。与数据库种子数据字段对齐。

- [ ] **Step 2: 新增 ProjectShowcaseConfig 结构体**

在 AdvantageItem 结构体之后（约第 29 行后）插入：

```go
// ProjectShowcaseConfig holds the project showcase section settings.
type ProjectShowcaseConfig struct {
	SectionTitle    string   `json:"section_title"`
	SectionSubtitle string   `json:"section_subtitle"`
	FeaturedSlugs   []string `json:"featured_slugs"`
}
```

- [ ] **Step 3: 更新 HomeConfigData 结构体**

当前代码约第 32-35 行，替换为：

```go
// HomeConfigData holds the parsed homepage configuration data.
type HomeConfigData struct {
	HeroSlides      []HeroSlide           `json:"hero_slides"`
	AdvantageItems  []AdvantageItem       `json:"advantage_items"`
	ProjectShowcase *ProjectShowcaseConfig `json:"project_showcase"`
}
```

使用指针类型 `*ProjectShowcaseConfig`，当数据库无此 key 时返回 `null`，前端据此判断是否需要默认值。

- [ ] **Step 4: 更新 Get() 方法，新增解析 project_showcase**

当前 `Get()` 方法约第 38-49 行。在解析 `advantage_items` 之后（约第 46 行后），新增：

```go
if projCfg, err := s.repo.FindByKey("project_showcase"); err == nil {
	var psc ProjectShowcaseConfig
	if err := json.Unmarshal(projCfg.ConfigValue, &psc); err == nil {
		data.ProjectShowcase = &psc
	}
}
```

完整的 `Get()` 方法变为：

```go
func (s *HomeConfigService) Get() (*HomeConfigData, error) {
	data := &HomeConfigData{}

	if heroCfg, err := s.repo.FindByKey("hero_slides"); err == nil {
		json.Unmarshal(heroCfg.ConfigValue, &data.HeroSlides)
	}
	if advCfg, err := s.repo.FindByKey("advantage_items"); err == nil {
		json.Unmarshal(advCfg.ConfigValue, &data.AdvantageItems)
	}
	if projCfg, err := s.repo.FindByKey("project_showcase"); err == nil {
		var psc ProjectShowcaseConfig
		if err := json.Unmarshal(projCfg.ConfigValue, &psc); err == nil {
			data.ProjectShowcase = &psc
		}
	}

	return data, nil
}
```

- [ ] **Step 5: 运行测试验证**

```bash
cd backend && go build ./...
```

预期：编译成功，无报错。

- [ ] **Step 6: Commit**

```bash
cd backend && git add internal/service/home_config_svc.go && git commit -m "feat: update HeroSlide fields and add ProjectShowcaseConfig"
```

---

### Task 2: Frontend — 新建 admin/homepage.vue 页面

**Files:**
- Create: `frontend/pages/admin/homepage.vue`

这是本次改动最大的文件——首页配置管理页面。分步构建。

- [ ] **Step 1: 创建基础骨架（script + template 空壳 + 数据加载）**

```vue
<template>
  <div class="homepage-config" v-loading="loading">
    <div class="page-header">
      <h2 class="page-title">首页配置</h2>
    </div>
    <div class="config-body">
      <!-- Hero Slides Card -->
      <el-card class="config-card">
        <template #header><h3 class="card-title">轮播管理</h3></template>
        <!-- 稍后填充 -->
      </el-card>

      <!-- Project Showcase Card -->
      <el-card class="config-card">
        <template #header><h3 class="card-title">项目展示区</h3></template>
        <!-- 稍后填充 -->
      </el-card>

      <!-- Advantage Items Card -->
      <el-card class="config-card">
        <template #header><h3 class="card-title">优势管理</h3></template>
        <!-- 稍后填充 -->
      </el-card>
    </div>
  </div>
</template>

<script setup lang="ts">
definePageMeta({ layout: 'admin', middleware: ['auth'] });

interface HeroSlide {
  title: string;
  desc: string;
  project_slug: string;
  gradient: string;
  image: string;
  link: string;
}

interface AdvantageItem {
  icon: string;
  title: string;
  description: string;
}

interface ProjectShowcase {
  section_title: string;
  section_subtitle: string;
  featured_slugs: string[];
}

const heroSlides = ref<HeroSlide[]>([]);
const advantageItems = ref<AdvantageItem[]>([]);
const projectShowcase = ref<ProjectShowcase>({
  section_title: '',
  section_subtitle: '',
  featured_slugs: [],
});

const allProjects = ref<Array<{ slug: string; title: string }>>([]);
const loading = ref(true);

const load = async () => {
  loading.value = true;
  try {
    const api = useApi();
    const [config, projects] = await Promise.all([
      api<{
        hero_slides: HeroSlide[];
        advantage_items: AdvantageItem[];
        project_showcase: ProjectShowcase | null;
      }>('/admin/home-config'),
      api<Array<{ slug: string; title: string }>>('/projects'),
    ]);

    if (config) {
      heroSlides.value = config.hero_slides || [];
      advantageItems.value = config.advantage_items || [];
      if (config.project_showcase) {
        projectShowcase.value = config.project_showcase;
      }
    }
    
    // Deduplicate projects by slug
    if (projects) {
      const seen = new Set<string>();
      allProjects.value = (projects as Array<{ slug: string; title: string }>).filter((p) => {
        if (seen.has(p.slug)) return false;
        seen.add(p.slug);
        return true;
      });
    }
  } finally {
    loading.value = false;
  }
};

onMounted(load);
</script>

<style scoped>
.page-header { margin-bottom: 24px; }
.page-title { font-size: 22px; font-weight: 600; }
.config-body { max-width: 900px; }
.config-card { margin-bottom: 20px; }
.card-title { margin: 0; font-size: 16px; font-weight: 600; }
</style>
```

- [ ] **Step 2: 实现轮播管理 Card 的主体逻辑**

在 script setup 中新增轮播相关方法、dialog 状态：

```ts
// --- Hero Slides ---
const slideDialogVisible = ref(false);
const slideForm = ref<HeroSlide>(blankSlide());
const slideEditIndex = ref(-1);
const slideSaving = ref(false);

function blankSlide(): HeroSlide {
  return { title: '', desc: '', project_slug: '', gradient: '', image: '', link: '' };
}

function openAddSlide() {
  slideEditIndex.value = -1;
  slideForm.value = blankSlide();
  slideDialogVisible.value = true;
}

function openEditSlide(index: number) {
  slideEditIndex.value = index;
  slideForm.value = { ...heroSlides.value[index] };
  slideDialogVisible.value = true;
}

function removeSlide(index: number) {
  heroSlides.value.splice(index, 1);
}

function moveSlide(index: number, direction: -1 | 1) {
  const target = index + direction;
  if (target < 0 || target >= heroSlides.value.length) return;
  const items = [...heroSlides.value];
  [items[index], items[target]] = [items[target], items[index]];
  heroSlides.value = items;
}

async function saveSlide() {
  if (!slideForm.value.title.trim()) {
    ElMessage.warning('请填写标题');
    return;
  }
  if (slideEditIndex.value >= 0) {
    heroSlides.value[slideEditIndex.value] = { ...slideForm.value };
  } else {
    heroSlides.value.push({ ...slideForm.value });
  }
  slideDialogVisible.value = false;
  await saveSlides();
}

async function saveSlides() {
  slideSaving.value = true;
  try {
    const api = useApi();
    await api('/admin/home-config', {
      method: 'PUT',
      body: { hero_slides: heroSlides.value },
    });
    ElMessage.success('轮播已保存');
  } catch {
    ElMessage.error('保存失败');
  } finally {
    slideSaving.value = false;
  }
}
```

更新 template 中轮播 Card：

```html
<el-card class="config-card">
  <template #header>
    <div class="card-header">
      <h3 class="card-title">轮播管理</h3>
      <el-button type="primary" size="small" @click="openAddSlide">新增 Slide</el-button>
    </div>
  </template>
  <div v-if="heroSlides.length === 0" class="empty-hint">暂无轮播，点击"新增 Slide"添加。</div>
  <div v-else class="slide-list">
    <div v-for="(slide, i) in heroSlides" :key="i" class="slide-item">
      <span class="slide-label">{{ slide.title || '(无标题)' }}</span>
      <div class="slide-actions">
        <el-button size="small" :disabled="i === 0" @click="moveSlide(i, -1)">↑</el-button>
        <el-button size="small" :disabled="i === heroSlides.length - 1" @click="moveSlide(i, 1)">↓</el-button>
        <el-button size="small" @click="openEditSlide(i)">编辑</el-button>
        <el-button size="small" type="danger" @click="removeSlide(i)">删除</el-button>
      </div>
    </div>
  </div>
</el-card>
```

- [ ] **Step 3: 实现轮播编辑弹窗**

```html
<el-dialog
  v-model="slideDialogVisible"
  :title="slideEditIndex >= 0 ? '编辑 Slide' : '新增 Slide'"
  width="560px"
  destroy-on-close
>
  <el-form label-width="100px">
    <el-form-item label="标题" required>
      <el-input v-model="slideForm.title" placeholder="主标题" />
    </el-form-item>
    <el-form-item label="描述">
      <el-input v-model="slideForm.desc" placeholder="描述文案" />
    </el-form-item>
    <el-form-item label="关联项目">
      <el-select v-model="slideForm.project_slug" placeholder="(可选)" clearable>
        <el-option v-for="p in allProjects" :key="p.slug" :label="p.title" :value="p.slug" />
      </el-select>
    </el-form-item>
    <el-form-item label="背景渐变色">
      <el-input v-model="slideForm.gradient" placeholder="linear-gradient(135deg, #1a3a5c, #2d5a8e)" />
    </el-form-item>
    <el-form-item label="背景图 URL">
      <el-input v-model="slideForm.image" placeholder="图片地址" />
    </el-form-item>
    <el-form-item label="跳转链接">
      <el-input v-model="slideForm.link" placeholder="点击跳转链接(可选)" />
    </el-form-item>
  </el-form>
  <template #footer>
    <el-button @click="slideDialogVisible = false">取消</el-button>
    <el-button type="primary" @click="saveSlide">确定</el-button>
  </template>
</el-dialog>
```

- [ ] **Step 4: 实现项目展示区 Card**

script 中新增：

```ts
// --- Project Showcase ---
const showcaseSaving = ref(false);

// featured project ordering
function moveFeatured(index: number, direction: -1 | 1) {
  const target = index + direction;
  if (target < 0 || target >= projectShowcase.value.featured_slugs.length) return;
  const slugs = [...projectShowcase.value.featured_slugs];
  [slugs[index], slugs[target]] = [slugs[target], slugs[index]];
  projectShowcase.value.featured_slugs = slugs;
}

function removeFeatured(index: number) {
  projectShowcase.value.featured_slugs.splice(index, 1);
}

const availableProjects = computed(() => {
  const featured = new Set(projectShowcase.value.featured_slugs);
  return allProjects.value.filter((p) => !featured.has(p.slug));
});

function addFeatured(slug: string) {
  if (!projectShowcase.value.featured_slugs.includes(slug)) {
    projectShowcase.value.featured_slugs.push(slug);
  }
}

async function saveShowcase() {
  showcaseSaving.value = true;
  try {
    const api = useApi();
    await api('/admin/home-config', {
      method: 'PUT',
      body: { project_showcase: projectShowcase.value },
    });
    ElMessage.success('项目展示区已保存');
  } catch {
    ElMessage.error('保存失败');
  } finally {
    showcaseSaving.value = false;
  }
}
```

template 中项目展示区 Card：

```html
<el-card class="config-card">
  <template #header><h3 class="card-title">项目展示区</h3></template>
  <el-form label-width="100px">
    <el-form-item label="区域标题">
      <el-input v-model="projectShowcase.section_title" placeholder="精选移民项目" />
    </el-form-item>
    <el-form-item label="区域副标题">
      <el-input v-model="projectShowcase.section_subtitle" placeholder="为您量身定制的最佳移民方案" />
    </el-form-item>
    <el-form-item label="精选项目">
      <div class="featured-area">
        <div v-if="projectShowcase.featured_slugs.length === 0" class="empty-hint">
          未选择精选项目，首页将展示全部项目。
        </div>
        <div v-for="(slug, i) in projectShowcase.featured_slugs" :key="slug" class="featured-row">
          <span class="featured-name">{{ allProjects.find((p) => p.slug === slug)?.title || slug }}</span>
          <div class="featured-actions">
            <el-button size="small" :disabled="i === 0" @click="moveFeatured(i, -1)">↑</el-button>
            <el-button size="small" :disabled="i === projectShowcase.featured_slugs.length - 1" @click="moveFeatured(i, 1)">↓</el-button>
            <el-button size="small" type="danger" @click="removeFeatured(i)">移除</el-button>
          </div>
        </div>
        <el-select
          v-if="availableProjects.length > 0"
          value=""
          placeholder="添加项目..."
          clearable
          @change="(val: string) => { if (val) addFeatured(val); }"
          class="add-project-select"
        >
          <el-option v-for="p in availableProjects" :key="p.slug" :label="p.title" :value="p.slug" />
        </el-select>
      </div>
    </el-form-item>
  </el-form>
  <div class="card-footer">
    <el-button type="primary" :loading="showcaseSaving" @click="saveShowcase">保存</el-button>
  </div>
</el-card>
```

- [ ] **Step 5: 实现优势管理 Card**

script 中新增：

```ts
// --- Advantage Items ---
const advDialogVisible = ref(false);
const advForm = ref<AdvantageItem>({ icon: '', title: '', description: '' });
const advEditIndex = ref(-1);
const advSaving = ref(false);

function openAddAdv() {
  advEditIndex.value = -1;
  advForm.value = { icon: '', title: '', description: '' };
  advDialogVisible.value = true;
}

function openEditAdv(index: number) {
  advEditIndex.value = index;
  advForm.value = { ...advantageItems.value[index] };
  advDialogVisible.value = true;
}

function removeAdv(index: number) {
  advantageItems.value.splice(index, 1);
}

function moveAdv(index: number, direction: -1 | 1) {
  const target = index + direction;
  if (target < 0 || target >= advantageItems.value.length) return;
  const items = [...advantageItems.value];
  [items[index], items[target]] = [items[target], items[index]];
  advantageItems.value = items;
}

async function saveAdv() {
  if (!advForm.value.title.trim()) {
    ElMessage.warning('请填写标题');
    return;
  }
  if (advEditIndex.value >= 0) {
    advantageItems.value[advEditIndex.value] = { ...advForm.value };
  } else {
    advantageItems.value.push({ ...advForm.value });
  }
  advDialogVisible.value = false;
  await saveAdvantages();
}

async function saveAdvantages() {
  advSaving.value = true;
  try {
    const api = useApi();
    await api('/admin/home-config', {
      method: 'PUT',
      body: { advantage_items: advantageItems.value },
    });
    ElMessage.success('优势项已保存');
  } catch {
    ElMessage.error('保存失败');
  } finally {
    advSaving.value = false;
  }
}
```

template 中优势管理 Card：

```html
<el-card class="config-card">
  <template #header>
    <div class="card-header">
      <h3 class="card-title">优势管理</h3>
      <el-button type="primary" size="small" @click="openAddAdv">新增优势项</el-button>
    </div>
  </template>
  <div v-if="advantageItems.length === 0" class="empty-hint">暂无优势项，点击"新增优势项"添加。</div>
  <div v-else class="adv-list">
    <div v-for="(item, i) in advantageItems" :key="i" class="adv-item">
      <span class="adv-icon">{{ item.icon }}</span>
      <div class="adv-info">
        <strong>{{ item.title }}</strong>
        <span class="adv-desc">{{ item.description }}</span>
      </div>
      <div class="adv-actions">
        <el-button size="small" :disabled="i === 0" @click="moveAdv(i, -1)">↑</el-button>
        <el-button size="small" :disabled="i === advantageItems.length - 1" @click="moveAdv(i, 1)">↓</el-button>
        <el-button size="small" @click="openEditAdv(i)">编辑</el-button>
        <el-button size="small" type="danger" @click="removeAdv(i)">删除</el-button>
      </div>
    </div>
  </div>
</el-card>
```

优势编辑弹窗：

```html
<el-dialog
  v-model="advDialogVisible"
  :title="advEditIndex >= 0 ? '编辑优势项' : '新增优势项'"
  width="500px"
  destroy-on-close
>
  <el-form label-width="80px">
    <el-form-item label="图标">
      <el-input v-model="advForm.icon" placeholder="emoji 或图标名" />
    </el-form-item>
    <el-form-item label="标题" required>
      <el-input v-model="advForm.title" placeholder="标题" />
    </el-form-item>
    <el-form-item label="描述">
      <el-input v-model="advForm.description" type="textarea" :rows="3" placeholder="描述文案" />
    </el-form-item>
  </el-form>
  <template #footer>
    <el-button @click="advDialogVisible = false">取消</el-button>
    <el-button type="primary" @click="saveAdv">确定</el-button>
  </template>
</el-dialog>
```

- [ ] **Step 6: 补充完整样式**

在 `<style scoped>` 末尾追加：

```css
.card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}
.empty-hint {
  color: #909399;
  font-size: 14px;
  padding: 16px 0;
  text-align: center;
}
/* Slide list */
.slide-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 10px 0;
  border-bottom: 1px solid #ebeef5;
}
.slide-item:last-child { border-bottom: none; }
.slide-label { font-size: 14px; color: #303133; }
.slide-actions { display: flex; gap: 6px; flex-shrink: 0; }
/* Advantage list */
.adv-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 10px 0;
  border-bottom: 1px solid #ebeef5;
}
.adv-item:last-child { border-bottom: none; }
.adv-icon { font-size: 24px; width: 36px; text-align: center; flex-shrink: 0; }
.adv-info { flex: 1; display: flex; flex-direction: column; gap: 2px; }
.adv-info strong { font-size: 14px; }
.adv-desc { font-size: 12px; color: #909399; }
.adv-actions { display: flex; gap: 6px; flex-shrink: 0; }
/* Featured area */
.featured-area { width: 100%; }
.featured-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 6px 0;
  border-bottom: 1px solid #ebeef5;
}
.featured-row:last-child { border-bottom: none; }
.featured-name { font-size: 14px; }
.featured-actions { display: flex; gap: 6px; }
.add-project-select { margin-top: 8px; width: 100%; }
/* Card footer */
.card-footer { text-align: center; padding-top: 16px; }
```

- [ ] **Step 7: 验证编译**

```bash
cd frontend && npx nuxi typecheck
```

预期：类型检查通过。

- [ ] **Step 8: Commit**

```bash
git add frontend/pages/admin/homepage.vue
git commit -m "feat: add homepage config admin page"
```

---

### Task 3: Frontend — 侧边栏新增菜单项

**Files:**
- Modify: `frontend/layouts/admin.vue`

- [ ] **Step 1: 在"页面管理"和"导航管理"之间插入菜单项**

找到当前 sidebar-nav 中这两行（约第 14-17 行）：

```html
<NuxtLink to="/admin/pages" class="nav-item" active-class="active">
  <span>页面管理</span>
</NuxtLink>
<NuxtLink to="/admin/navigation" class="nav-item" active-class="active">
  <span>导航管理</span>
</NuxtLink>
```

替换为：

```html
<NuxtLink to="/admin/pages" class="nav-item" active-class="active">
  <span>页面管理</span>
</NuxtLink>
<NuxtLink to="/admin/homepage" class="nav-item" active-class="active">
  <span>首页配置</span>
</NuxtLink>
<NuxtLink to="/admin/navigation" class="nav-item" active-class="active">
  <span>导航管理</span>
</NuxtLink>
```

- [ ] **Step 2: Commit**

```bash
git add frontend/layouts/admin.vue && git commit -m "feat: add homepage config to admin sidebar"
```

---

### Task 4: Frontend — 仪表盘快捷入口新增

**Files:**
- Modify: `frontend/pages/admin/index.vue`

- [ ] **Step 1: 在快捷入口网格中新增"首页配置"**

找到 links-grid（约第 33-48 行），在合适位置插入：

```html
<NuxtLink to="/admin/homepage" class="quick-link">
  <span>首页配置</span>
</NuxtLink>
```

建议放在"页面管理"和"导航管理"之间，与侧边栏顺序一致。当前代码约第 33 行，在"页面管理"那行之后插入。

- [ ] **Step 2: Commit**

```bash
git add frontend/pages/admin/index.vue && git commit -m "feat: add homepage config quick link to dashboard"
```

---

### Task 5: Frontend — 首页 index.vue 适配新字段

**Files:**
- Modify: `frontend/pages/index.vue`

- [ ] **Step 1: 更新 hero_slides 字段映射**

当前代码约第 163-170 行：

```ts
if (data && Array.isArray(data.hero_slides)) {
  heroSlides.value = (data.hero_slides as Array<Record<string, string>>).map((s) => ({
    title: s.title || '',
    subtitle: s.subtitle || '',
    image: s.image_url || '',
  }));
}
```

替换为：

```ts
if (data && Array.isArray(data.hero_slides)) {
  heroSlides.value = (data.hero_slides as Array<Record<string, string>>).map((s) => ({
    title: s.title || '',
    subtitle: s.desc || '',
    image: s.image || '',
    link: s.link || '',
    gradient: s.gradient || '',
  }));
}
```

同时更新 `HeroSlide` 接口（约第 111-115 行），新增 `link` 和 `gradient` 字段：

```ts
interface HeroSlide {
  title: string;
  subtitle: string;
  image: string;
  link?: string;
  gradient?: string;
}
```

- [ ] **Step 2: 更新项目展示区，读取 project_showcase 配置**

当前项目展示区（约第 142-260 行）使用的是硬编码默认值。需新增读取 `project_showcase` 配置。

在 `homeConfig` 解包后新增项目展示区配置的读取。找到当前 `const projectCards = computed<...>(() => {` 约第 182 行，在其上方新增：

```ts
// 从 homeConfig 中读取项目展示区配置
const showcaseConfig = computed(() => {
  if (homeConfig.value) {
    const config = homeConfig.value as unknown as Record<string, unknown>;
    const data = config.data as Record<string, unknown> | undefined;
    if (data && data.project_showcase) {
      return data.project_showcase as {
        section_title?: string;
        section_subtitle?: string;
        featured_slugs?: string[];
      };
    }
  }
  return null;
});

// 展示用标题和副标题
const sectionTitle = computed(() => showcaseConfig.value?.section_title || '精选移民项目');
const sectionSubtitle = computed(() => showcaseConfig.value?.section_subtitle || '为您量身定制的最佳移民方案');
```

然后让 template 中的 section 标题和副标题使用这两个 computed。改第 45-46 行：

```html
<h2 class="section-title">{{ sectionTitle }}</h2>
<p class="section-subtitle">{{ sectionSubtitle }}</p>
```

在 `projectCards` computed 的 API 分支中（约第 188-196 行），返回数据前按 `featured_slugs` 过滤排序：

```ts
if (apiProjects && apiProjects.length > 0) {
  const featured = showcaseConfig.value?.featured_slugs;
  let items = apiProjects.map((p) => ({
    slug: p.slug,
    title: p.title,
    description: p.description,
    image: p.cover_image || '',
    features: [] as string[],
    link: `/projects/${p.slug}`,
  }));

  if (featured && featured.length > 0) {
    // Sort: featured projects first, in specified order
    const orderMap = new Map(featured.map((s: string, i: number) => [s, i]));
    items.sort((a, b) => {
      const ai = orderMap.get(a.slug);
      const bi = orderMap.get(b.slug);
      if (ai !== undefined && bi !== undefined) return ai - bi;
      if (ai !== undefined) return -1;
      if (bi !== undefined) return 1;
      return 0;
    });
  }

  return items;
}
```

- [ ] **Step 3: 优势区字段保持不变**

当前优势区代码不需要改动（icon/title/description 映射已匹配）。

- [ ] **Step 4: 验证编译**

```bash
cd frontend && npx nuxi typecheck
```

预期：类型检查通过。

- [ ] **Step 5: Commit**

```bash
git add frontend/pages/index.vue && git commit -m "feat: update homepage to read new hero_slides fields and project_showcase config"
```

---

### Task 6: 端到端验证

- [ ] **Step 1: 启动后端**

```bash
cd backend && go run ./cmd/server/main.go
```

预期：后端启动在 8080 端口。

- [ ] **Step 2: 启动前端**

```bash
cd frontend && npm run dev
```

预期：前端启动在 3000 端口。

- [ ] **Step 3: 功能验证**

1. 登录管理后台（admin / admin123 或实际账号）
2. 侧边栏出现"首页配置"菜单项，点击进入
3. 轮播管理：新增/编辑/删除/排序 slide，保存
4. 项目展示区：修改标题/副标题，选择精选项目，排序，保存
5. 优势管理：新增/编辑/删除/排序，保存
6. 访问前台首页：验证轮播/项目展示/优势区正确渲染
7. 数据库验证：检查 `home_configs` 表中数据正确持久化

- [ ] **Step 4: 运行后端测试**

```bash
cd backend && go test ./... -v
```

预期：全部测试通过。

- [ ] **Step 5: 运行前端类型检查**

```bash
cd frontend && npx nuxi typecheck
```

预期：类型检查通过。
