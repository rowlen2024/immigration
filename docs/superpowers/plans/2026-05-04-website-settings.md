# Website Settings Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add a "网站设置" admin page with 22 site-wide config fields, replace hardcoded values in Header/Footer/useSeo with dynamic config, and auto-inject SEO/GEO/analytics tags.

**Architecture:** Reuse the existing `home_configs` key-value table with a new `site` key. Backend adds a `GetSiteConfig`/`UpdateSiteConfig` service + a public GET and admin PUT route. Frontend adds a new admin settings page with grouped forms and tooltips, a `useSiteConfig` composable, and wires the config into Header, Footer, useSeo, and app-wide head injections.

**Tech Stack:** Go + Gin + GORM (backend), Nuxt 3 + Element Plus + Pinia (frontend)

---

## File Structure

| Action | File | Purpose |
|--------|------|---------|
| Modify | `backend/internal/service/home_config_svc.go` | Add `GetSiteConfig`/`UpdateSiteConfig` |
| Modify | `backend/internal/handler/home_handler.go` | Add `GetSiteConfig`/`UpdateSiteConfig` handlers |
| Modify | `backend/internal/router/router.go` | Register `GET /site-config` + `PUT /admin/site-config` |
| Create | `frontend/composables/useSiteConfig.ts` | Fetch + reactive cache of site config |
| Create | `frontend/pages/admin/settings.vue` | Admin settings page with grouped forms + tooltips |
| Modify | `frontend/layouts/admin.vue` | Add "网站设置" sidebar link |
| Modify | `frontend/components/global/Header.vue` | Read site_name + site_logo from config |
| Modify | `frontend/components/global/Footer.vue` | Read contact/copyright from config |
| Modify | `frontend/composables/useSeo.ts` | Read site_name for title suffix |
| Modify | `frontend/app.vue` | Inject favicon, analytics, JSON-LD, custom code |

---

### Task 1: Backend — Add SiteConfig type and service methods

**Files:**
- Modify: `backend/internal/service/home_config_svc.go`

- [ ] **Step 1: Add SiteConfig struct and service methods**

Add after the existing `HomeConfigData` struct:

```go
// SiteConfig holds all site-wide settings.
type SiteConfig struct {
	SiteName                string   `json:"site_name"`
	SiteLogo                string   `json:"site_logo"`
	SiteFavicon             string   `json:"site_favicon"`
	SeoTitle                string   `json:"seo_title"`
	SeoDescription          string   `json:"seo_description"`
	SeoKeywords             string   `json:"seo_keywords"`
	OgImage                 string   `json:"og_image"`
	CanonicalBase           string   `json:"canonical_base"`
	OrganizationName        string   `json:"organization_name"`
	OrganizationDescription string   `json:"organization_description"`
	OrganizationLogo        string   `json:"organization_logo"`
	OrganizationURL         string   `json:"organization_url"`
	SameAs                  []string `json:"same_as"`
	ContactPhone            string   `json:"contact_phone"`
	ContactEmail            string   `json:"contact_email"`
	ContactAddress          string   `json:"contact_address"`
	ContactWechat           string   `json:"contact_wechat"`
	GATrackingID            string   `json:"ga_tracking_id"`
	BaiduTongjiID           string   `json:"baidu_tongji_id"`
	CustomHeadCode          string   `json:"custom_head_code"`
	CustomBodyCode          string   `json:"custom_body_code"`
	CopyrightText           string   `json:"copyright_text"`
	ICPNumber               string   `json:"icp_number"`
}

// DefaultSiteConfig returns sensible zero-value defaults for optional fields.
func DefaultSiteConfig() *SiteConfig {
	return &SiteConfig{
		SiteName:   "MyGo移民",
		SeoTitle:   "{site_name} | 专业投资移民服务",
		CopyrightText: "© {year} {site_name}. All rights reserved.",
	}
}
```

Add two methods to `HomeConfigService`:

```go
// GetSiteConfig returns the site configuration, falling back to defaults.
func (s *HomeConfigService) GetSiteConfig() (*SiteConfig, error) {
	cfg, err := s.repo.FindByKey("site")
	if err != nil {
		return DefaultSiteConfig(), nil
	}

	var data SiteConfig
	if err := json.Unmarshal(cfg.ConfigValue, &data); err != nil {
		return nil, fmt.Errorf("failed to parse site config: %w", err)
	}
	return &data, nil
}

// UpdateSiteConfig replaces the entire site configuration.
func (s *HomeConfigService) UpdateSiteConfig(input *SiteConfig) error {
	raw, err := json.Marshal(input)
	if err != nil {
		return fmt.Errorf("failed to marshal site config: %w", err)
	}

	cfg := &model.HomeConfig{
		ConfigKey:   "site",
		ConfigValue: raw,
	}
	// Try update first, create if not found
	if err := s.repo.Update(cfg); err != nil {
		if err := s.repo.Create(cfg); err != nil {
			return fmt.Errorf("failed to save site config: %w", err)
		}
	}
	return nil
}
```

Import `"encoding/json"` and `"mygo-immigration/backend/internal/model"` are already imported. Verify `"fmt"` is imported.

- [ ] **Step 2: Verify compilation**

Run: `cd backend && go build ./...`
Expected: builds without errors.

- [ ] **Step 3: Commit**

```bash
git add backend/internal/service/home_config_svc.go
git commit -m "feat: add SiteConfig type and service methods for website settings"
```

---

### Task 2: Backend — Add site-config handlers

**Files:**
- Modify: `backend/internal/handler/home_handler.go`

- [ ] **Step 1: Add GetSiteConfig and UpdateSiteConfig handlers**

Append after the existing `UpdateHomeConfig` method:

```go
func (h *Handler) GetSiteConfig(c *gin.Context) {
	cfg, err := h.svc.HomeConfig.GetSiteConfig()
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Error(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.Success(cfg))
}

func (h *Handler) UpdateSiteConfig(c *gin.Context) {
	var cfg service.SiteConfig
	if err := c.ShouldBindJSON(&cfg); err != nil {
		c.JSON(http.StatusBadRequest, dto.Error(400, "invalid request body"))
		return
	}

	if err := h.svc.HomeConfig.UpdateSiteConfig(&cfg); err != nil {
		c.JSON(http.StatusInternalServerError, dto.Error(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.Success(nil))
}
```

- [ ] **Step 2: Verify compilation**

Run: `cd backend && go build ./...`
Expected: builds without errors.

- [ ] **Step 3: Commit**

```bash
git add backend/internal/handler/home_handler.go
git commit -m "feat: add GetSiteConfig and UpdateSiteConfig HTTP handlers"
```

---

### Task 3: Backend — Register routes

**Files:**
- Modify: `backend/internal/router/router.go`

- [ ] **Step 1: Register the two new routes**

Add in the public routes section (after `api.GET("/navigation", ...)`):

```go
api.GET("/site-config", h.GetSiteConfig)
```

Add in the admin routes section (after the `home-config` admin block):

```go
admin.PUT("/site-config", middleware.RBAC("content:write"), h.UpdateSiteConfig)
```

- [ ] **Step 2: Verify compilation and run tests**

Run: `cd backend && go build ./...`
Expected: builds without errors.

Run: `cd backend && go test ./... -v -count=1`
Expected: existing tests pass.

- [ ] **Step 3: Commit**

```bash
git add backend/internal/router/router.go
git commit -m "feat: register GET /site-config and PUT /admin/site-config routes"
```

---

### Task 4: Frontend — Create useSiteConfig composable

**Files:**
- Create: `frontend/composables/useSiteConfig.ts`

- [ ] **Step 1: Create the composable**

```typescript
interface SiteConfig {
  site_name: string;
  site_logo: string;
  site_favicon: string;
  seo_title: string;
  seo_description: string;
  seo_keywords: string;
  og_image: string;
  canonical_base: string;
  organization_name: string;
  organization_description: string;
  organization_logo: string;
  organization_url: string;
  same_as: string[];
  contact_phone: string;
  contact_email: string;
  contact_address: string;
  contact_wechat: string;
  ga_tracking_id: string;
  baidu_tongji_id: string;
  custom_head_code: string;
  custom_body_code: string;
  copyright_text: string;
  icp_number: string;
}

const siteConfig = ref<SiteConfig | null>(null);
let fetchPromise: Promise<void> | null = null;

export const useSiteConfig = () => {
  const fetch = async () => {
    if (siteConfig.value) return;
    if (fetchPromise) {
      await fetchPromise;
      return;
    }

    fetchPromise = (async () => {
      try {
        const api = useApi();
        siteConfig.value = await api<SiteConfig>('/site-config');
      } catch {
        // Silently use defaults — non-critical
        siteConfig.value = null;
      }
    })();

    await fetchPromise;
    fetchPromise = null;
  };

  return {
    siteConfig: readonly(siteConfig),
    fetch,
  };
};
```

- [ ] **Step 2: Verify TypeScript check**

Run: `cd frontend && npx nuxi prepare && npx nuxi typecheck`
Expected: no errors.

- [ ] **Step 3: Commit**

```bash
git add frontend/composables/useSiteConfig.ts
git commit -m "feat: add useSiteConfig composable with fetch and cache"
```

---

### Task 5: Frontend — Create /admin/settings page

**Files:**
- Create: `frontend/pages/admin/settings.vue`

- [ ] **Step 1: Create the settings page with grouped forms and tooltips**

```vue
<template>
  <div class="page-wrapper">
    <div class="page-header">
      <h2 class="page-title">网站设置</h2>
    </div>

    <div v-loading="loading" class="settings-body">
      <el-card v-for="group in groups" :key="group.key" class="settings-card">
        <template #header>
          <h3 class="card-title">{{ group.label }}</h3>
        </template>

        <el-form label-width="180px" label-position="right">
          <el-form-item
            v-for="field in group.fields"
            :key="field.key"
            :label="field.label"
          >
            <div class="field-wrap">
              <!-- same_as: dynamic array input -->
              <div v-if="field.key === 'same_as'" class="array-input">
                <div
                  v-for="(item, idx) in form.same_as"
                  :key="idx"
                  class="array-row"
                >
                  <el-input v-model="form.same_as[idx]" placeholder="https://" />
                  <el-button
                    type="danger"
                    :icon="Delete"
                    circle
                    size="small"
                    @click="removeSameAs(idx)"
                  />
                </div>
                <el-button type="primary" text @click="addSameAs">
                  + 添加链接
                </el-button>
              </div>

              <!-- custom_head_code / custom_body_code: textarea -->
              <el-input
                v-else-if="field.textarea"
                v-model="form[field.key]"
                type="textarea"
                :rows="field.rows || 6"
                class="monospace-input"
              />

              <!-- default text input -->
              <el-input v-else v-model="form[field.key]" :placeholder="field.placeholder" />

              <el-tooltip
                :content="field.tip"
                placement="right"
                effect="dark"
                raw-content
              >
                <span class="tip-icon">?</span>
              </el-tooltip>
            </div>
          </el-form-item>
        </el-form>
      </el-card>

      <div class="save-bar">
        <el-button type="primary" size="large" :loading="saving" @click="save">
          保存设置
        </el-button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { Delete } from '@element-plus/icons-vue';

definePageMeta({ layout: 'admin', middleware: ['auth'] });

interface SiteConfig {
  [key: string]: any;
  site_name: string;
  site_logo: string;
  site_favicon: string;
  seo_title: string;
  seo_description: string;
  seo_keywords: string;
  og_image: string;
  canonical_base: string;
  organization_name: string;
  organization_description: string;
  organization_logo: string;
  organization_url: string;
  same_as: string[];
  contact_phone: string;
  contact_email: string;
  contact_address: string;
  contact_wechat: string;
  ga_tracking_id: string;
  baidu_tongji_id: string;
  custom_head_code: string;
  custom_body_code: string;
  copyright_text: string;
  icp_number: string;
}

const defaultForm = (): SiteConfig => ({
  site_name: 'MyGo移民',
  site_logo: '',
  site_favicon: '',
  seo_title: '{site_name} | 专业投资移民服务',
  seo_description: '',
  seo_keywords: '',
  og_image: '',
  canonical_base: '',
  organization_name: '',
  organization_description: '',
  organization_logo: '',
  organization_url: '',
  same_as: [],
  contact_phone: '',
  contact_email: '',
  contact_address: '',
  contact_wechat: '',
  ga_tracking_id: '',
  baidu_tongji_id: '',
  custom_head_code: '',
  custom_body_code: '',
  copyright_text: '© {year} {site_name}. All rights reserved.',
  icp_number: '',
});

const tips: Record<string, string> = {
  site_name: '网站在导航栏、浏览器标题栏等位置显示的名称',
  site_logo: '网站主 Logo 图片地址，显示在页头导航栏',
  site_favicon: '浏览器标签页和书签栏显示的小图标，建议 32×32px',
  seo_title: '搜索引擎结果页显示的标题模板。<br/><code>{site_name}</code> 会被替换为网站名称',
  seo_description: '搜索引擎结果页显示的页面描述，建议 120-160 字，各内页无独立描述时使用此值',
  seo_keywords: '页面关键词，用逗号分隔。现代搜索引擎权重已降低，但仍建议填写',
  og_image: '在微信、Facebook 等社交平台分享链接时显示的默认预览图',
  canonical_base: '网站首选访问域名，用于规范搜索引擎索引（如 https://www.example.com）',
  organization_name: '向 Google 知识图谱和 AI 搜索声明的企业法律实体全称',
  organization_description: '用于生成 AI 搜索摘要和知识面板的企业介绍，建议包含成立时间、核心业务和规模',
  organization_logo: 'Google 知识面板展示的 Logo，建议 112×112px 以上高清 PNG',
  organization_url: '声明机构的官方网址，用于多平台交叉验证网站真实性',
  same_as: '与官网关联的社交媒体链接（LinkedIn、公众号、小红书等）。AI 搜索引擎通过双向链接验证实体可信度',
  contact_phone: '网站底部和结构化数据中展示的客服电话',
  contact_email: '网站底部和结构化数据中展示的客服邮箱',
  contact_address: '公司办公地址，显示在网站底部',
  contact_wechat: '企业微信号，显示在网站底部联系方式区域',
  ga_tracking_id: 'Google Analytics 4 衡量 ID（格式：G-XXXXXXXX），用于网站流量分析',
  baidu_tongji_id: '百度统计站点 ID，用于中国国内流量分析',
  custom_head_code: '插入到每个页面 &lt;head&gt; 标签内的自定义代码（meta 验证标签、第三方脚本等）',
  custom_body_code: '插入到每个页面 &lt;/body&gt; 闭合标签前的自定义代码',
  copyright_text: '页脚版权声明。<code>{year}</code> 动态替换为当前年份，<code>{site_name}</code> 替换为网站名称',
  icp_number: 'ICP 备案号（中国大陆运营网站必需），如 沪ICP备XXXXXXXX号',
};

const form = ref<SiteConfig>(defaultForm());
const loading = ref(true);
const saving = ref(false);

interface FieldDef {
  key: string;
  label: string;
  placeholder?: string;
  textarea?: boolean;
  rows?: number;
  tip: string;
}

interface GroupDef {
  key: string;
  label: string;
  fields: FieldDef[];
}

const groups: GroupDef[] = [
  {
    key: 'basic', label: '基础信息',
    fields: [
      { key: 'site_name', label: '网站名称', placeholder: 'MyGo移民', tip: tips.site_name },
      { key: 'site_logo', label: '网站 Logo URL', placeholder: '/images/logo.png', tip: tips.site_logo },
      { key: 'site_favicon', label: 'Favicon URL', placeholder: '/favicon.ico', tip: tips.site_favicon },
    ],
  },
  {
    key: 'seo', label: 'SEO 设置',
    fields: [
      { key: 'seo_title', label: '标题模板', tip: tips.seo_title },
      { key: 'seo_description', label: 'Meta 描述', tip: tips.seo_description },
      { key: 'seo_keywords', label: '关键词', tip: tips.seo_keywords },
      { key: 'og_image', label: 'OG 分享图', tip: tips.og_image },
      { key: 'canonical_base', label: '首选域名', placeholder: 'https://www.example.com', tip: tips.canonical_base },
    ],
  },
  {
    key: 'geo', label: 'GEO / 结构化数据',
    fields: [
      { key: 'organization_name', label: '机构名称', tip: tips.organization_name },
      { key: 'organization_description', label: '机构描述', tip: tips.organization_description },
      { key: 'organization_logo', label: '机构 Logo', tip: tips.organization_logo },
      { key: 'organization_url', label: '官网 URL', placeholder: 'https://www.example.com', tip: tips.organization_url },
      { key: 'same_as', label: '社交媒体链接', tip: tips.same_as },
    ],
  },
  {
    key: 'contact', label: '联系方式',
    fields: [
      { key: 'contact_phone', label: '客服电话', placeholder: '400-xxx-xxxx', tip: tips.contact_phone },
      { key: 'contact_email', label: '客服邮箱', placeholder: 'info@example.com', tip: tips.contact_email },
      { key: 'contact_address', label: '公司地址', tip: tips.contact_address },
      { key: 'contact_wechat', label: '企业微信', tip: tips.contact_wechat },
    ],
  },
  {
    key: 'third_party', label: '第三方代码',
    fields: [
      { key: 'ga_tracking_id', label: 'Google Analytics ID', placeholder: 'G-XXXXXXXX', tip: tips.ga_tracking_id },
      { key: 'baidu_tongji_id', label: '百度统计 ID', tip: tips.baidu_tongji_id },
      { key: 'custom_head_code', label: '自定义 Head', textarea: true, rows: 6, tip: tips.custom_head_code },
      { key: 'custom_body_code', label: '自定义 Body', textarea: true, rows: 6, tip: tips.custom_body_code },
    ],
  },
  {
    key: 'footer', label: '页脚设置',
    fields: [
      { key: 'copyright_text', label: '版权声明', tip: tips.copyright_text },
      { key: 'icp_number', label: 'ICP 备案号', placeholder: '沪ICP备XXXXXXXX号', tip: tips.icp_number },
    ],
  },
];

const load = async () => {
  loading.value = true;
  try {
    const api = useApi();
    const data = await api<SiteConfig>('/site-config');
    if (data) {
      form.value = { ...defaultForm(), ...data };
    }
  } finally {
    loading.value = false;
  }
};

const save = async () => {
  saving.value = true;
  try {
    const api = useApi();
    await api('/admin/site-config', {
      method: 'PUT',
      body: JSON.parse(JSON.stringify(form.value)),
    });
    ElMessage.success('设置已保存');
  } catch {
    ElMessage.error('保存失败');
  } finally {
    saving.value = false;
  }
};

const addSameAs = () => form.value.same_as.push('');
const removeSameAs = (idx: number) => form.value.same_as.splice(idx, 1);

onMounted(load);
</script>

<style scoped>
.settings-body {
  max-width: 800px;
}

.settings-card {
  margin-bottom: 20px;
}

.card-title {
  margin: 0;
  font-size: 16px;
  font-weight: 600;
}

.field-wrap {
  display: flex;
  align-items: center;
  gap: 8px;
  width: 100%;
}

.field-wrap > .el-input,
.field-wrap > .el-input-number {
  flex: 1;
}

.tip-icon {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 18px;
  height: 18px;
  border-radius: 50%;
  background: #c0c4cc;
  color: #fff;
  font-size: 12px;
  font-weight: 700;
  cursor: help;
  flex-shrink: 0;
}

.tip-icon:hover {
  background: #909399;
}

.monospace-input :deep(textarea) {
  font-family: 'Courier New', Courier, monospace;
  font-size: 13px;
}

.array-input {
  flex: 1;
}

.array-row {
  display: flex;
  gap: 8px;
  margin-bottom: 8px;
  align-items: center;
}

.save-bar {
  padding: 20px 0;
  text-align: center;
}
</style>
```

- [ ] **Step 2: Verify TypeScript check**

Run: `cd frontend && npx nuxi typecheck`
Expected: no type errors.

- [ ] **Step 3: Commit**

```bash
git add frontend/pages/admin/settings.vue
git commit -m "feat: add /admin/settings page with grouped forms and tooltips"
```

---

### Task 6: Frontend — Add sidebar link in admin layout

**Files:**
- Modify: `frontend/layouts/admin.vue`

- [ ] **Step 1: Add "网站设置" nav item**

Add after the "媒体库" link (before `</nav>`) in the template:

```vue
<NuxtLink to="/admin/settings" class="nav-item" active-class="active">
  <span>网站设置</span>
</NuxtLink>
```

- [ ] **Step 2: Verify TypeScript check**

Run: `cd frontend && npx nuxi typecheck`
Expected: no type errors.

- [ ] **Step 3: Commit**

```bash
git add frontend/layouts/admin.vue
git commit -m "feat: add 网站设置 link to admin sidebar"
```

---

### Task 7: Frontend — Wire Header to use site config

**Files:**
- Modify: `frontend/components/global/Header.vue`

- [ ] **Step 1: Make logo dynamic**

In `<script setup>`, add:

```typescript
const { siteConfig, fetch: fetchSiteConfig } = useSiteConfig();
```

After `onMounted(() => { fetchNav(); });`, add fetch call:

```typescript
onMounted(() => {
  fetchNav();
  fetchSiteConfig();
});
```

Replace the hardcoded logo text:

```vue
<NuxtLink to="/" class="header-logo">
  <img v-if="siteConfig?.site_logo" :src="siteConfig.site_logo" :alt="siteConfig?.site_name || 'MyGo移民'" class="logo-img" />
  <span v-else class="logo-text">{{ siteConfig?.site_name || 'MyGo移民' }}</span>
</NuxtLink>
```

Add style for the logo image:

```css
.logo-img {
  height: 36px;
  width: auto;
}
```

- [ ] **Step 2: Verify TypeScript check**

Run: `cd frontend && npx nuxi typecheck`
Expected: no type errors.

- [ ] **Step 3: Commit**

```bash
git add frontend/components/global/Header.vue
git commit -m "feat: wire Header logo to site config"
```

---

### Task 8: Frontend — Wire Footer to use site config

**Files:**
- Modify: `frontend/components/global/Footer.vue`

- [ ] **Step 1: Make contact info and copyright dynamic**

In `<script setup>`, add:

```typescript
const { siteConfig, fetch: fetchSiteConfig } = useSiteConfig();

onMounted(() => { fetchSiteConfig(); });
```

Replace the hardcoded contact items:

```vue
<li class="contact-item">
  <span class="contact-label">电话：</span>
  <a v-if="siteConfig?.contact_phone" :href="`tel:${siteConfig.contact_phone}`" class="footer-link">{{ siteConfig.contact_phone }}</a>
  <span v-else class="footer-link">400-xxx-xxxx</span>
</li>
<li class="contact-item">
  <span class="contact-label">邮箱：</span>
  <a v-if="siteConfig?.contact_email" :href="`mailto:${siteConfig.contact_email}`" class="footer-link">{{ siteConfig.contact_email }}</a>
  <span v-else class="footer-link">info@mygo-immigration.com</span>
</li>
<li class="contact-item">
  <span class="contact-label">地址：</span>
  <span>{{ siteConfig?.contact_address || '上海市浦东新区陆家嘴金融中心' }}</span>
</li>
<li class="contact-item">
  <span class="contact-label">微信：</span>
  <span>{{ siteConfig?.contact_wechat || 'MyGo_Immigration' }}</span>
</li>
```

Replace the copyright line:

```vue
<p class="copyright">
  {{ copyrightText }}
</p>
```

Add the computed in `<script setup>`:

```typescript
const copyrightText = computed(() => {
  const template = siteConfig.value?.copyright_text || '© {year} {site_name}. All rights reserved.';
  return template
    .replace('{year}', String(new Date().getFullYear()))
    .replace('{site_name}', siteConfig.value?.site_name || 'MyGo移民');
});
```

- [ ] **Step 2: Verify TypeScript check**

Run: `cd frontend && npx nuxi typecheck`
Expected: no type errors.

- [ ] **Step 3: Commit**

```bash
git add frontend/components/global/Footer.vue
git commit -m "feat: wire Footer contact/copyright to site config"
```

---

### Task 9: Frontend — Wire useSeo to use site config

**Files:**
- Modify: `frontend/composables/useSeo.ts`

- [ ] **Step 1: Use site_name for title suffix**

```typescript
interface SeoOptions {
  title?: string;
  description?: string;
  jsonLd?: Record<string, unknown>;
}

export const useSeo = (options: SeoOptions) => {
  const { siteConfig } = useSiteConfig();
  const siteName = computed(() => siteConfig.value?.site_name || 'MyGo移民');

  const fullTitle = computed(() => {
    if (options.title) {
      return `${options.title} | ${siteName.value}`;
    }
    const seoTitle = siteConfig.value?.seo_title || '{site_name} | 专业投资移民服务';
    return seoTitle.replace('{site_name}', siteName.value);
  });

  const head = computed(() => {
    const h: Record<string, unknown> = {
      title: fullTitle.value,
      meta: [] as Record<string, unknown>[],
    };

    if (siteConfig.value?.og_image || options.description) {
      const metas: Record<string, unknown>[] = [];

      if (options.description || siteConfig.value?.seo_description) {
        const desc = options.description || siteConfig.value?.seo_description || '';
        metas.push(
          { name: 'description', content: desc },
          { property: 'og:title', content: fullTitle.value },
          { property: 'og:description', content: desc },
        );
      }

      if (siteConfig.value?.og_image) {
        metas.push({ property: 'og:image', content: siteConfig.value.og_image });
      }

      if (siteConfig.value?.canonical_base) {
        metas.push({ rel: 'canonical', href: siteConfig.value.canonical_base });
      }

      (h.meta as Record<string, unknown>[]).push(...metas);
    }

    if (options.jsonLd) {
      h.script = [
        {
          type: 'application/ld+json',
          innerHTML: JSON.stringify(options.jsonLd),
        },
      ];
    }

    return h;
  });

  useHead(head);
};
```

- [ ] **Step 2: Verify TypeScript check**

Run: `cd frontend && npx nuxi typecheck`
Expected: no type errors.

- [ ] **Step 3: Commit**

```bash
git add frontend/composables/useSeo.ts
git commit -m "feat: wire useSeo title suffix and og/canonical to site config"
```

---

### Task 10: Frontend — Inject global head tags in app.vue

**Files:**
- Modify: `frontend/app.vue`

- [ ] **Step 1: Read current app.vue**

First, read the current `frontend/app.vue` to see the existing structure.

- [ ] **Step 2: Add global injections**

Add in `<script setup>`:

```typescript
const { siteConfig, fetch: fetchSiteConfig } = useSiteConfig();

onMounted(() => { fetchSiteConfig(); });

// Inject favicon
useHead(() => {
  const links: Record<string, unknown>[] = [];
  const scripts: Record<string, unknown>[] = [];

  if (siteConfig.value?.site_favicon) {
    links.push({ rel: 'icon', type: 'image/x-icon', href: siteConfig.value.site_favicon });
  }

  // Inject GA
  if (siteConfig.value?.ga_tracking_id) {
    scripts.push({
      async: true,
      src: `https://www.googletagmanager.com/gtag/js?id=${siteConfig.value.ga_tracking_id}`,
    });
    scripts.push({
      innerHTML: `window.dataLayer=window.dataLayer||[];function gtag(){dataLayer.push(arguments);}gtag('js',new Date());gtag('config','${siteConfig.value.ga_tracking_id}');`,
    });
  }

  // Inject Baidu Tongji
  if (siteConfig.value?.baidu_tongji_id) {
    scripts.push({
      innerHTML: `var _hmt=_hmt||[];(function(){var hm=document.createElement("script");hm.src="https://hm.baidu.com/hm.js?${siteConfig.value.baidu_tongji_id}";var s=document.getElementsByTagName("script")[0];s.parentNode.insertBefore(hm,s);})();`,
    });
  }

  // Inject custom head code (raw HTML — treated as text node in head)
  if (siteConfig.value?.custom_head_code) {
    scripts.push({
      innerHTML: siteConfig.value.custom_head_code,
    });
  }

  return { link: links, script: scripts };
});

// Inject Organization JSON-LD
useHead(() => {
  const org = siteConfig.value;
  if (!org?.organization_name) return {};

  const jsonLd: Record<string, unknown> = {
    '@context': 'https://schema.org',
    '@type': 'Organization',
    name: org.organization_name,
    url: org.organization_url || '',
  };

  if (org.organization_description) {
    jsonLd.description = org.organization_description;
  }
  if (org.organization_logo) {
    jsonLd.logo = org.organization_logo;
  }
  if (org.same_as && org.same_as.length > 0) {
    jsonLd.sameAs = org.same_as.filter((s) => s.trim() !== '');
  }
  if (org.contact_phone || org.contact_email) {
    jsonLd.contactPoint = {
      '@type': 'ContactPoint',
      telephone: org.contact_phone || '',
      email: org.contact_email || '',
      contactType: 'customer service',
    };
  }
  if (org.contact_address) {
    jsonLd.address = {
      '@type': 'PostalAddress',
      streetAddress: org.contact_address,
    };
  }

  return {
    script: [
      {
        type: 'application/ld+json',
        innerHTML: JSON.stringify(jsonLd),
      },
    ],
  };
});

// Inject custom body-end code via a client-only script
if (import.meta.client && siteConfig.value?.custom_body_code) {
  watch(
    () => siteConfig.value?.custom_body_code,
    (code) => {
      if (code) {
        const el = document.createElement('div');
        el.innerHTML = code;
        document.body.appendChild(el);
      }
    },
    { immediate: true },
  );
}
```

- [ ] **Step 3: Verify TypeScript check**

Run: `cd frontend && npx nuxi typecheck`
Expected: no type errors.

- [ ] **Step 4: Commit**

```bash
git add frontend/app.vue
git commit -m "feat: inject favicon, analytics, JSON-LD, and custom code from site config"
```

---

### Task 11: Integration test — Manual verification

**Files:**
- None (verification only)

- [ ] **Step 1: Launch services**

Run: `make up` (MySQL)
Run: `make dev-backend` (backend on :8080)
Run: `make dev-frontend` (frontend on :3000)

- [ ] **Step 2: Check API**

```bash
# Public endpoint
curl http://localhost:8080/api/v1/site-config | jq .

# Update endpoint (requires token)
curl -X PUT http://localhost:8080/api/v1/admin/site-config \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{"site_name":"Test Site","seo_description":"Test desc"}'
```

Expected: GET returns default config; PUT saves and returns 200.

- [ ] **Step 3: Check admin page**

Open `http://localhost:3000/admin/settings`
Expected:
- 6 grouped cards rendered
- Each field has a `?` tooltip icon
- Hovering shows the tip text
- Fill in values, click "保存设置" → success message
- Refresh page → values persist

- [ ] **Step 4: Check consumers**

- Header shows configured site name/logo
- Footer shows configured contact info and copyright
- View page source → `og:*` meta tags present
- View page source → JSON-LD Organization schema present
- If GA ID set → GA script tag appears

- [ ] **Step 5: Commit (if any fixes needed)**

```bash
git add -A
git commit -m "fix: integration fixes for site settings feature"
```

---

## Backend Test (optional but recommended)

### Task 12: Backend — Handler test for site-config

**Files:**
- Create: `backend/internal/handler/site_config_handler_test.go`

- [ ] **Step 1: Write handler tests**

```go
package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"mygo-immigration/backend/internal/dto"
)

func TestGetSiteConfig_ReturnsDefaults(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/api/v1/site-config", nil)

	h := setupHandler()
	h.GetSiteConfig(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var resp dto.Response
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp.Code != 200 {
		t.Fatalf("expected code 200, got %d", resp.Code)
	}
}
```

Note: `setupHandler()` needs a real DB connection — follow the existing test pattern from `home_handler_test.go` or `auth_handler_test.go`. If those tests require MySQL, write this test to match the same pattern.

- [ ] **Step 2: Run the test**

Run: `cd backend && go test ./internal/handler -run TestGetSiteConfig -v`
Expected: PASS

- [ ] **Step 3: Commit**

```bash
git add backend/internal/handler/site_config_handler_test.go
git commit -m "test: add handler test for GET /site-config"
```
