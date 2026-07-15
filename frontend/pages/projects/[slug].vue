<template>
  <div class="project-detail">
    <div class="container">
      <ProjectBreadcrumb :label="project.title" />
    </div>

    <section class="detail-hero" :style="heroFallbackStyle">
      <ResponsiveImage
        v-if="project.cover_image"
        :src="project.cover_image"
        alt=""
        variant="lg"
        :variants="project.cover_image_variants"
        loading="eager"
        fetchpriority="high"
        sizes="100vw"
        class="detail-hero-bg"
      />
      <div class="detail-hero-overlay"></div>
      <div class="container">
        <h1 class="detail-title">{{ project.title }}</h1>
        <p class="detail-summary">{{ project.summary }}</p>
      </div>
    </section>

    <nav v-if="tabs.length > 0" class="tab-nav" ref="tabNavRef">
      <div class="container">
        <div class="tab-nav-scroll">
          <button
            v-for="tab in tabs"
            :key="tab.id"
            class="tab-btn"
            :class="{ active: activeTab === tab.id }"
            @click="scrollToSection(tab.id)"
          >{{ tab.label }}</button>
        </div>
      </div>
    </nav>

    <div class="container">
      <div v-if="pending" class="page-skeleton-wrapper"><PageSkeleton variant="detail" /></div>
      <div v-else-if="error" class="error-state">{{ error }}</div>
      <template v-else>
        <section id="overview" class="detail-section">
          <h2 class="detail-section-title">项目概览</h2>
          <ProjectQuickFacts :facts="quickFacts" />
          <div v-if="project.description" class="detail-content" style="margin-top: 24px;">
            <p>{{ project.description }}</p>
          </div>
        </section>

        <section id="requirements" v-if="requirements.length > 0" class="detail-section">
          <h2 class="detail-section-title">申请条件</h2>
          <ProjectRequirementsChecklist :items="requirements" />
        </section>

        <section id="cost" v-if="costTable.length > 0" class="detail-section">
          <h2 class="detail-section-title">费用明细</h2>
          <ProjectCostTable :rows="costTable" />
        </section>

        <section id="timeline" v-if="timelinePhases.length > 0" class="detail-section">
          <h2 class="detail-section-title">申请流程</h2>
          <ProjectTimeline :phases="timelinePhases" />
        </section>

        <section id="advantages" v-if="advantages.length > 0" class="detail-section">
          <h2 class="detail-section-title">项目优势</h2>
          <div class="advantages-grid">
            <div v-for="(adv, index) in advantages" :key="index" class="advantage-card">
              <div class="advantage-icon">
                <span v-if="getIconByName(adv.icon)" v-html="getIconSvg(adv.icon, 22, '#C8963E')" class="advantage-svg"></span>
                <span v-else class="advantage-svg-fallback">
                  <span v-html="getIconSvg('star', 22, '#C8963E')"></span>
                </span>
              </div>
              <h3 class="advantage-title">{{ adv.title }}</h3>
              <p class="advantage-desc">{{ adv.description }}</p>
            </div>
          </div>
        </section>

        <section id="cases" v-if="project.cases.length > 0" class="detail-section">
          <h2 class="detail-section-title">成功案例</h2>
          <div class="case-grid">
            <CaseCard
              v-for="c in project.cases"
              :key="c.name"
              :slug="c.slug"
              :name="c.name"
              :country="c.country"
              :summary="c.content ? stripHtml(c.content) : ''"
              :image="c.photo"
              :image-variants="c.photo_variants"
              :meta-text="`${c.country} | ${c.amount} | ${c.period}`"
            />
          </div>
        </section>

        <section id="testimonials" v-if="project.testimonials.length > 0" class="detail-section">
          <h2 class="detail-section-title">客户评价</h2>
          <TestimonialCarousel :testimonials="project.testimonials" />
        </section>

        <section id="news" v-if="project.news.length > 0" class="detail-section">
          <h2 class="detail-section-title">最新资讯</h2>
          <div class="news-list">
            <NuxtLink v-for="n in project.news" :key="n.id" :to="`/pages/${n.slug}`" class="news-item">
              <ResponsiveImage v-if="n.cover" :src="n.cover" :alt="n.title" variant="thumb" :variants="n.cover_variants" class="news-cover" loading="lazy" />
              <div class="news-body">
                <h4 class="news-title">{{ n.title }}</h4>
                <span v-if="n.date" class="news-date">{{ new Date(n.date).toLocaleDateString('zh-CN') }}</span>
              </div>
            </NuxtLink>
          </div>
        </section>

        <section id="compare" v-if="project.compare_config && project.compare_config.compare_with.length >= 2" class="detail-section">
          <h2 class="detail-section-title">项目对比</h2>
          <div v-if="comparePending" class="compare-loading">加载对比数据...</div>
          <div v-else-if="compareError" class="compare-error">{{ compareError }}</div>
          <template v-else-if="compareData">
            <!-- 桌面端：传统表格 -->
            <div class="compare-table-wrap" :class="{ 'hidden-mobile': true }">
              <table class="compare-table">
                <thead>
                  <tr>
                    <th class="compare-label-col">对比项目</th>
                    <th v-for="(proj, i) in compareData.projects" :key="i" class="compare-col-header">
                      {{ proj.title }}
                    </th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="row in compareData.rows" :key="row.label">
                    <td class="compare-label">{{ row.label }}</td>
                    <td v-for="(val, j) in row.values" :key="j" class="compare-value">
                      <template v-if="row.items?.[j]?.length">
                        <div class="compare-requirements-grid">
                          <span v-for="(item, k) in row.items[j]" :key="k">{{ item }}</span>
                        </div>
                      </template>
                      <template v-else>{{ val }}</template>
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>
            <!-- 移动端：下拉选择单属性对比 -->
            <div
              class="compare-switcher hidden-desktop"
              ref="compareCardRef"
              @touchstart="onCompareTouchStart"
              @touchend="onCompareTouchEnd"
            >
              <!-- 属性下拉选择 -->
              <div class="cmp-switch-select-wrap">
                <select
                  v-model="compareIdx"
                  class="cmp-switch-select"
                >
                  <option v-for="(crow, ci) in compareData.rows" :key="crow.label" :value="ci">
                    {{ crow.label }}{{ hasDiff(crow) ? ' · 有差异' : '' }}
                  </option>
                </select>
                <svg class="cmp-switch-select-arrow" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
                  <polyline points="6 9 12 15 18 9"/>
                </svg>
              </div>
              <!-- 当前属性内容 -->
              <div class="cmp-switch-body">
                <div class="cmp-switch-cols" :class="{ 'cols-2': compareData.projects.length === 2, 'cols-3': compareData.projects.length >= 3 }">
                  <div v-for="(proj, idx) in compareData.projects" :key="idx" class="cmp-switch-col">
                    <span class="cmp-switch-proj">{{ proj.title }}</span>
                    <template v-if="currentRow.items?.[idx]?.length">
                      <ul class="cmp-switch-items">
                        <li v-for="(item, k) in currentRow.items[idx]" :key="k">{{ item }}</li>
                      </ul>
                    </template>
                    <span v-else class="cmp-switch-val">{{ currentRow.values[idx] }}</span>
                  </div>
                </div>
              </div>
            </div>
          </template>
        </section>

        <section id="faqs" v-if="faqs.length > 0" class="detail-section">
          <h2 class="detail-section-title">常见问题</h2>
          <ProjectFaqListSection
            :items="paginatedFaqs"
            :page="faqPage"
            :per-page="faqPerPage"
            :total="faqs.length"
            @page-change="changeFaqPage"
          />
        </section>

        <ConsultCTA
          :title="`对${project.title}感兴趣？`"
          description="立即联系我们，专业顾问为您一对一解答"
        />
      </template>
    </div>
  </div>
</template>

<script setup lang="ts">
import { getIconByName, getIconSvg } from '~/composables/lucideIcons';
import type { ImageVariantInfo } from '~/utils/image'
import { stripHtml } from '~/utils/html'
import { buildServiceJsonLd, buildFAQPageJsonLd, toJsonLdConfig, toJsonLdScripts } from '~/utils/jsonld'

const route = useRoute();
const slug = computed(() => route.params.slug as string);
const projectDataKey = computed(() => `public:project:${slug.value}`);

interface ApiRequirement { label: string; is_required: boolean; }
interface ApiCostItem { name: string; amount: string; note: string; }
interface ApiTimelinePhase { phase_number: number; title: string; description: string; duration: string; }
interface ApiFAQ { question: string; answer: string; }
interface ApiCase { slug: string; name: string; country_from: string; investment_amount: string; processing_period: string; content: string; photo_url: string; photo_variants?: Record<string, ImageVariantInfo>; }
interface ApiTestimonial { id: number; avatar_url: string; avatar_variants?: Record<string, ImageVariantInfo>; nickname: string; rating: number; content: string; }
interface ApiNewsPage { id: number; title: string; slug: string; cover_image: string; cover_image_variants?: Record<string, ImageVariantInfo>; created_at: string; }
interface ApiCompareConfig { compare_with: string[]; compare_fields: string[]; }

interface ApiAdvantage { icon: string; icon_type: string; title: string; description: string; }

interface ApiProject {
  name: string;
  tagline: string;
  country: string;
  cover_image: string;
  cover_image_variants?: Record<string, ImageVariantInfo>;
  investment_amount: string;
  processing_period: string;
  target_crowd: string;
  overview_title: string;
  overview_text: string;
  cta_text: string;
  hero_title: string;
  hero_desc: string;
  hero_gradient: string;
  requirements: ApiRequirement[];
  cost_items: ApiCostItem[];
  timeline_phases: ApiTimelinePhase[];
  faqs: ApiFAQ[];
  cases: ApiCase[];
  testimonials: ApiTestimonial[];
  news: ApiNewsPage[];
  compare_config: ApiCompareConfig | null;
  advantages: ApiAdvantage[];
}

usePublicDataFreshness(() => [projectDataKey.value]);

onMounted(setupSectionObserver);
onUnmounted(cleanupSectionObserver);

const { data, pending, error, refresh } = await useFetch<{ data: ApiProject }>(
  () => `/api/v1/projects/${slug.value}`,
  { key: projectDataKey },
);

if (error.value?.statusCode === 404) {
  throw createError({
    statusCode: 404,
    statusMessage: '项目不存在',
    fatal: true,
  });
}

const project = computed(() => {
  const p = data.value?.data;
  return {
    title: p?.name || slug.value,
    summary: p?.tagline || '',
    description: p?.overview_text || '',
    cover_image: p?.cover_image || '',
    cover_image_variants: p?.cover_image_variants,
    investment_amount: p?.investment_amount || '',
    processing_period: p?.processing_period || '',
    target_crowd: p?.target_crowd || '',
    requirements: (p?.requirements || []).map((r) => ({ text: r.label, met: r.is_required })),
    cost_table: (p?.cost_items || []).map((c) => ({ item: c.name, amount: c.amount, note: c.note })),
    timeline: (p?.timeline_phases || []).map((t) => ({
      phase: `第${t.phase_number}步`,
      title: t.title,
      description: t.description,
      period: t.duration,
    })),
    faqs: (p?.faqs || []).map((f) => ({ question: f.question, answer: f.answer })),
    cases: (p?.cases || []).map((c) => ({
      slug: c.slug,
      name: c.name,
      country: c.country_from,
      amount: c.investment_amount,
      period: c.processing_period,
      content: c.content,
      photo: c.photo_url,
      photo_variants: c.photo_variants,
    })),
    testimonials: (p?.testimonials || []).map((t) => ({
      id: t.id,
      avatar_url: t.avatar_url,
      avatar_variants: t.avatar_variants,
      nickname: t.nickname,
      rating: t.rating,
      content: t.content,
    })),
    news: (p?.news || []).map((n) => ({
      id: n.id,
      title: n.title,
      slug: n.slug,
      cover: n.cover_image,
      cover_variants: n.cover_image_variants,
      date: n.created_at,
    })),
    compare_config: p?.compare_config || null,
    advantages: (p?.advantages || []).map((a) => ({
      icon: a.icon,
      icon_type: a.icon_type,
      title: a.title,
      description: a.description,
    })),
  };
});

// Compare data fetch
interface CompareRowData {
  label: string;
  values: string[];
  items?: string[][];
}

interface CompareTableData {
  projects: Array<{ title: string; slug: string }>;
  rows: CompareRowData[];
}

const compareSlugs = computed(() => {
  const cfg = project.value.compare_config;
  if (!cfg || !cfg.compare_with || cfg.compare_with.length < 2) return '';
  return cfg.compare_with.join(',');
});

const compareQuery = computed(() => {
  const slugs = compareSlugs.value;
  if (!slugs) return {};
  const params: Record<string, string> = { slugs };
  const fields = project.value.compare_config?.compare_fields;
  if (fields && fields.length > 0) {
    params.fields = fields.join(',');
  }
  return params;
});

const {
  data: compareRaw,
  pending: comparePending,
  error: compareErrorRaw,
  refresh: refreshCompare,
} = useFetch<{ data: CompareTableData }>('/api/v1/projects/compare', {
  query: compareQuery,
  immediate: false,
  watch: false,
});

watch(compareSlugs, (slugs) => {
  if (slugs) refreshCompare();
}, { immediate: true });

const compareData = computed<CompareTableData | null>(() => {
  const raw = compareRaw.value as any;
  if (raw?.data?.rows) return raw.data;
  if (raw?.rows) return raw as CompareTableData;
  return null;
});

const compareError = computed(() =>
  compareErrorRaw.value ? '加载对比数据失败' : null
);

// 移动端 — 当前对比属性索引 + 触摸滑动
const compareIdx = ref(0);
const compareCardRef = ref<HTMLElement | null>(null);
let compareTouchStartX = 0;

function goCompareIdx(idx: number) {
  compareIdx.value = Math.max(0, Math.min(idx, (compareData.value?.rows.length || 1) - 1));
}

function onCompareTouchStart(e: TouchEvent) {
  compareTouchStartX = e.touches[0].clientX;
}

function onCompareTouchEnd(e: TouchEvent) {
  const diff = compareTouchStartX - e.changedTouches[0].clientX;
  if (Math.abs(diff) > 50) {
    goCompareIdx(compareIdx.value + (diff > 0 ? 1 : -1));
  }
}

const currentRow = computed(() => {
  return compareData.value?.rows[compareIdx.value] ?? { label: '', values: [], items: [] };
});

function hasDiff(row: { values: string[]; items?: string[][] }): boolean {
  const items = row.items;
  if (items && items[0]) {
    const base = items[0];
    for (let i = 1; i < items.length; i++) {
      if (!items[i] || base.length !== items[i].length) return true;
      if (items[i].some((v, j) => v !== base[j])) return true;
    }
    return false;
  }
  const vals = row.values;
  for (let i = 1; i < vals.length; i++) {
    if (vals[i] !== vals[0]) return true;
  }
  return false;
}

useSeo({
  title: project.value.title || '项目详情',
  description: project.value.summary || '',
  breadcrumbLabel: project.value.title,
});

// Structured data for search engines (Baidu AI, Google Knowledge Graph)
const { siteConfig } = useMygoSiteConfig()

useHead(() => {
  const p = project.value
  if (!p?.title) return {}

  const base = siteConfig.value?.canonical_base || ''
  const pageUrl = base ? base + route.path : undefined

  const rated = (p.testimonials || []).filter((t: any) => t.rating > 0)
  const avgRating = rated.length > 0
    ? Number((rated.reduce((s: number, t: any) => s + t.rating, 0) / rated.length).toFixed(1))
    : null

  return {
    script: toJsonLdScripts(
      buildServiceJsonLd({
        name: p.title,
        description: p.description || p.summary || '',
        category: '移民服务',
        url: pageUrl,
        image: p.cover_image,
        investmentAmount: p.investment_amount,
        avgRating,
        reviewCount: rated.length,
      }, toJsonLdConfig(siteConfig.value)),
      buildFAQPageJsonLd(p.faqs || []),
    ),
  }
})

const heroFallbackStyle = computed(() => {
  return project.value.cover_image
    ? {}
    : { background: 'linear-gradient(135deg, #1a3a5c, #2d5a8e)' };
});

const quickFacts = computed(() => [
  { label: '投资金额', value: project.value.investment_amount },
  { label: '办理周期', value: project.value.processing_period },
  { label: '适合人群', value: project.value.target_crowd },
]);

const requirements = computed(() => project.value.requirements);
const costTable = computed(() => project.value.cost_table);
const timelinePhases = computed(() => project.value.timeline);
const faqs = computed(() => project.value.faqs);
const advantages = computed(() => project.value.advantages);

const faqPage = ref(1);
const faqPerPage = 10;
const paginatedFaqs = computed(() => {
  const start = (faqPage.value - 1) * faqPerPage;
  return faqs.value.slice(start, start + faqPerPage);
});
const changeFaqPage = (p: number) => {
  faqPage.value = p;
  const el = document.getElementById('faqs');
  if (el) window.scrollTo({ top: el.getBoundingClientRect().top + window.scrollY - 80, behavior: 'smooth' });
};

interface TabItem {
  id: string;
  label: string;
}

const tabs = computed<TabItem[]>(() => {
  const items: TabItem[] = [];
  items.push({ id: 'overview', label: '项目概览' });
  if (requirements.value.length > 0) items.push({ id: 'requirements', label: '申请条件' });
  if (costTable.value.length > 0) items.push({ id: 'cost', label: '费用明细' });
  if (timelinePhases.value.length > 0) items.push({ id: 'timeline', label: '申请流程' });
  if (advantages.value.length > 0) items.push({ id: 'advantages', label: '项目优势' });
  if (project.value.cases.length > 0) items.push({ id: 'cases', label: '成功案例' });
  if (project.value.testimonials.length > 0) items.push({ id: 'testimonials', label: '客户评价' });
  if (project.value.news.length > 0) items.push({ id: 'news', label: '最新资讯' });
  if (project.value.compare_config && project.value.compare_config.compare_with.length >= 2) items.push({ id: 'compare', label: '项目对比' });
  if (faqs.value.length > 0) items.push({ id: 'faqs', label: '常见问题' });
  return items;
});

const activeTab = ref('overview');
const tabNavRef = ref<HTMLElement | null>(null);

function scrollToSection(id: string) {
  const el = document.getElementById(id);
  if (!el) return;
  const navH = tabNavRef.value?.offsetHeight || 0;
  const top = el.getBoundingClientRect().top + window.scrollY - navH - 64;
  window.scrollTo({ top, behavior: 'smooth' });
}

let observer: IntersectionObserver | null = null;

function setupSectionObserver() {
  const headerH = (document.querySelector('.site-header') as HTMLElement)?.offsetHeight || 64;
  observer = new IntersectionObserver(
    (entries) => {
      for (const entry of entries) {
        if (entry.isIntersecting) {
          activeTab.value = entry.target.id;
        }
      }
    },
    { rootMargin: `-${(tabNavRef.value?.offsetHeight || 0) + headerH}px 0px -60% 0px` },
  );
  for (const tab of tabs.value) {
    const el = document.getElementById(tab.id);
    if (el) observer.observe(el);
  }

  // 滚动指示器阴影
  const scrollEl = tabNavRef.value?.querySelector('.tab-nav-scroll');
  if (scrollEl) {
    const updateShadows = () => {
      const hasLeft = scrollEl.scrollLeft > 4;
      const hasRight = scrollEl.scrollLeft + scrollEl.clientWidth < scrollEl.scrollWidth - 4;
      tabNavRef.value?.classList.toggle('shadow-left', hasLeft);
      tabNavRef.value?.classList.toggle('shadow-right', hasRight);
      if (tabNavRef.value) {
        tabNavRef.value.classList.toggle('shadow-right', hasRight);
      }
    };
    updateShadows();
    scrollEl.addEventListener('scroll', updateShadows, { passive: true });
    window.addEventListener('resize', updateShadows, { passive: true });
    (scrollEl as any)._shadowCleanup = () => {
      scrollEl.removeEventListener('scroll', updateShadows);
      window.removeEventListener('resize', updateShadows);
    };
  }

  // 动态追踪 header 高度，消除 tab-nav 与 header 之间的间隔
  const headerEl = document.querySelector('.site-header') as HTMLElement | null;
  const syncTabTop = () => {
    if (tabNavRef.value && headerEl) {
      tabNavRef.value.style.top = headerEl.offsetHeight + 'px';
    }
  };
  syncTabTop();
  window.addEventListener('scroll', syncTabTop, { passive: true });
  window.addEventListener('resize', syncTabTop, { passive: true });
  (tabNavRef.value as any)._topCleanup = () => {
    window.removeEventListener('scroll', syncTabTop);
    window.removeEventListener('resize', syncTabTop);
  };
}

// 切换 tab 时滚动按钮到可见区
watch(activeTab, (id) => {
  nextTick(() => {
    const btn = tabNavRef.value?.querySelector('.tab-btn.active');
    if (btn) {
      btn.scrollIntoView({ behavior: 'smooth', block: 'nearest', inline: 'center' });
    }
  });
});

function cleanupSectionObserver() {
  observer?.disconnect();
  observer = null;
  const scrollEl = tabNavRef.value?.querySelector('.tab-nav-scroll');
  if (scrollEl && (scrollEl as any)._shadowCleanup) {
    (scrollEl as any)._shadowCleanup();
  }
  if (tabNavRef.value && (tabNavRef.value as any)._topCleanup) {
    (tabNavRef.value as any)._topCleanup();
  }
}
</script>

<style scoped>
.tab-nav {
  position: sticky;
  top: var(--header-height);
  z-index: 50;
  background: var(--bg-white);
  border-bottom: 1px solid var(--border-color);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);
}

.tab-nav::before,
.tab-nav::after {
  content: '';
  position: absolute;
  top: 0;
  bottom: 0;
  width: 32px;
  pointer-events: none;
  z-index: 1;
  opacity: 0;
  transition: opacity 0.2s;
}

.tab-nav::before {
  left: 0;
  background: linear-gradient(to left, transparent, var(--bg-white));
}

.tab-nav::after {
  right: 0;
  background: linear-gradient(to right, transparent, var(--bg-white));
}

.tab-nav.shadow-left::before {
  opacity: 1;
}

.tab-nav.shadow-right::after {
  opacity: 1;
}

.tab-nav-scroll {
  display: flex;
  gap: 0;
  overflow-x: auto;
  -webkit-overflow-scrolling: touch;
  scrollbar-width: none;
}

.tab-nav-scroll::-webkit-scrollbar {
  display: none;
}

.tab-btn {
  flex-shrink: 0;
  padding: 14px 20px;
  font-size: 15px;
  font-family: var(--font-sans);
  font-weight: 500;
  color: var(--text-secondary);
  background: none;
  border: none;
  border-bottom: 2px solid transparent;
  cursor: pointer;
  transition: color 0.2s, border-color 0.2s;
  white-space: nowrap;
}

.tab-btn:hover {
  color: var(--primary);
}

.tab-btn:active {
  background: rgba(0, 0, 0, 0.04);
}

.tab-btn.active {
  color: var(--primary);
  border-bottom-color: var(--primary);
  font-weight: 600;
}

.detail-hero {
  position: relative;
  min-height: 400px;
  display: flex;
  align-items: center;
  background-size: cover;
  background-position: center;
  padding: 80px 0;
  color: var(--bg-white);
  margin-bottom: 0;
}

.detail-hero-bg {
  position: absolute;
  inset: 0;
  width: 100%;
  height: 100%;
  object-fit: cover;
  object-position: center;
  z-index: 0;
}

.detail-hero-overlay {
  position: absolute;
  inset: 0;
  background: linear-gradient(135deg, rgba(15, 36, 64, 0.85), rgba(26, 58, 92, 0.7));
  z-index: 1;
}

.detail-hero .container {
  position: relative;
  z-index: 2;
}

.detail-title {
  font-size: 42px;
  font-weight: 800;
  margin-bottom: 12px;
}

.detail-summary {
  font-size: 18px;
  opacity: 0.9;
  max-width: 600px;
  line-height: 1.6;
}

.detail-section {
  padding: 56px 0;
  border-bottom: 1px solid var(--border-color);
}

.detail-section:last-of-type {
  border-bottom: none;
}

.detail-section-title {
  font-size: 28px;
  font-weight: 700;
  color: var(--text-primary);
  margin-bottom: 28px;
  position: relative;
  padding-bottom: 12px;
}

.detail-section-title::after {
  content: '';
  position: absolute;
  bottom: 0;
  left: 0;
  width: 60px;
  height: 3px;
  background-color: var(--accent);
  border-radius: 2px;
}

.detail-content {
  font-size: 16px;
  color: var(--text-secondary);
  line-height: 1.9;
}

.loading-state,
.error-state {
  text-align: center;
  padding: 80px 20px;
  color: var(--text-light);
  font-size: 16px;
}

.error-state {
  color: #c62828;
}

.case-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
  gap: 24px;
}

@media (max-width: 767px) {
  .case-grid {
    grid-template-columns: 1fr;
    gap: 16px;
  }
}

.news-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.news-item {
  display: flex;
  gap: 16px;
  padding: 12px;
  border: 1px solid var(--border-color);
  border-radius: var(--radius-md);
  text-decoration: none;
  color: inherit;
  transition: box-shadow 0.2s;
}

.news-item:hover {
  box-shadow: var(--shadow-sm);
}

.news-item:active {
  background: var(--bg-gray);
  transform: scale(0.985);
}

@media (max-width: 767px) {
  .news-cover {
    width: 100px;
    height: 68px;
  }

  .news-title {
    font-size: 14px;
  }
}

.news-cover {
  width: 120px;
  height: 80px;
  aspect-ratio: 3 / 2;
  object-fit: cover;
  border-radius: 4px;
  flex-shrink: 0;
}

.news-body {
  display: flex;
  flex-direction: column;
  justify-content: center;
}

.news-title {
  font-size: 16px;
  font-weight: 600;
  margin-bottom: 6px;
}

.news-date {
  font-size: 13px;
  color: var(--text-light);
}

.compare-table-wrap {
  position: relative;
  overflow-x: auto;
  -webkit-overflow-scrolling: touch;
  background:
    linear-gradient(to right, var(--bg-white) 30%, transparent),
    linear-gradient(to right, transparent, var(--bg-white) 70%) 100% 0,
    linear-gradient(to right, rgba(0,0,0,0.08), transparent),
    linear-gradient(to left, rgba(0,0,0,0.08), transparent);
  background-repeat: no-repeat;
  background-size: 40px 100%, 40px 100%, 14px 100%, 14px 100%;
  background-attachment: local, local, scroll, scroll;
}

/* 覆盖全局金色滚动条为低调灰色 */
.compare-table-wrap::-webkit-scrollbar {
  height: 6px;
}
.compare-table-wrap::-webkit-scrollbar-track {
  background: var(--bg-light);
  border-radius: 3px;
}
.compare-table-wrap::-webkit-scrollbar-thumb {
  background: #c4c8cf;
  border-radius: 3px;
}

.compare-table {
  width: 100%;
  border-collapse: collapse;
  font-size: 14px;
  table-layout: fixed;
}

.compare-table thead {
  background-color: var(--primary);
  color: var(--bg-white);
}

.compare-table thead th {
  padding: 14px 16px;
  font-weight: 600;
  text-align: left;
}

.compare-table thead th:first-child {
  border-radius: var(--radius-md) 0 0 0;
}

.compare-table thead th:last-child {
  border-radius: 0 var(--radius-md) 0 0;
}

.compare-label-col {
  position: sticky;
  left: 0;
  z-index: 2;
  background-color: var(--primary);
  width: 120px;
}

.compare-col-header {
  /* table-layout:fixed 下自动等分剩余宽度 */
  word-break: break-word;
  overflow-wrap: break-word;
  hyphens: auto;
}

.compare-table td {
  padding: 12px 16px;
  border-bottom: 1px solid var(--border-color);
  color: var(--text-secondary);
  line-height: 1.6;
  word-break: break-word;
  overflow-wrap: break-word;
}

.compare-table tbody tr:nth-child(even) {
  background-color: var(--bg-light);
}

.compare-label {
  font-weight: 600;
  color: var(--text-primary);
  position: sticky;
  left: 0;
  z-index: 1;
  width: 120px;
}

.compare-table tbody tr:nth-child(even) .compare-label {
  background-color: var(--bg-light);
}

.compare-table tbody tr:nth-child(odd) .compare-label {
  background-color: var(--bg-white);
}

.compare-value {
  /* table-layout:fixed 下自动等分 */
  word-break: break-word;
  overflow-wrap: break-word;
}

.compare-requirements-grid {
  display: grid;
  grid-template-columns: 1fr;
  gap: 4px;
}

.compare-requirements-grid span {
  word-break: break-word;
  overflow-wrap: break-word;
  hyphens: auto;
  line-height: 1.5;
}

.compare-loading,
.compare-error {
  text-align: center;
  padding: 24px;
  color: var(--text-light);
  font-size: 14px;
}

.compare-error {
  color: #c62828;
}

/* ══════════════════════════════════════════════
   桌面/移动端 显示切换
   ══════════════════════════════════════════════ */
@media (min-width: 768px) {
  .hidden-desktop {
    display: none;
  }
}

@media (max-width: 767px) {
  .hidden-mobile {
    display: none;
  }
}

/* ══════════════════════════════════════════════
   移动端 — 下拉选择单属性对比
   ══════════════════════════════════════════════ */
@media (max-width: 767px) {
.compare-switcher {
  display: flex;
  flex-direction: column;
  gap: 0;
}

/* ── 下拉选择器 ── */
.cmp-switch-select-wrap {
  position: relative;
  margin-bottom: 16px;
}

.cmp-switch-select {
  width: 100%;
  padding: 12px 40px 12px 14px;
  font-size: 15px;
  font-weight: 600;
  font-family: var(--font-sans);
  color: var(--color-primary);
  background: var(--bg-white);
  border: 1px solid var(--border-color);
  border-radius: var(--radius-md);
  appearance: none;
  -webkit-appearance: none;
  cursor: pointer;
  min-height: 48px;
  line-height: 1.4;
}

.cmp-switch-select:focus-visible {
  outline: 2px solid var(--accent);
  outline-offset: -2px;
  border-color: var(--accent);
}

.cmp-switch-select-arrow {
  position: absolute;
  right: 14px;
  top: 50%;
  transform: translateY(-50%);
  color: var(--accent);
  pointer-events: none;
}

/* ── 内容区 ── */
.cmp-switch-body {
  background: var(--bg-white);
  border: 1px solid var(--border-color);
  border-radius: var(--radius-lg);
  padding: 20px 16px 16px;
}

.cmp-switch-cols {
  display: flex;
}

.cmp-switch-cols.cols-2 .cmp-switch-col {
  flex: 1;
  min-width: 0;
}

.cmp-switch-cols.cols-3 .cmp-switch-col {
  flex: 1;
  min-width: 0;
}

.cmp-switch-col {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.cmp-switch-col + .cmp-switch-col {
  border-left: 1px solid var(--border-color);
  padding-left: 16px;
  margin-left: 0;
}

.cmp-switch-col:first-child {
  padding-right: 16px;
}

/* 项目名 */
.cmp-switch-proj {
  font-size: 12px;
  font-weight: 600;
  color: var(--accent-dark);
  letter-spacing: 0.5px;
  text-transform: uppercase;
  word-break: break-word;
}

/* 普通值 */
.cmp-switch-val {
  font-size: 14px;
  color: var(--text-primary);
  line-height: 1.6;
  word-break: break-word;
}

/* 申请条件列表 */
.cmp-switch-items {
  list-style: none;
  padding: 0;
  margin: 0;
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.cmp-switch-items li {
  font-size: 13px;
  color: var(--text-secondary);
  line-height: 1.5;
  padding-left: 14px;
  position: relative;
}

.cmp-switch-items li::before {
  content: '';
  position: absolute;
  left: 0;
  top: 8px;
  width: 5px;
  height: 5px;
  border-radius: 50%;
  background: var(--accent);
}
} /* /@media mobile compare-switcher */

@media (max-width: 767px) {
  .detail-hero {
    min-height: 280px;
    padding: 48px 0;
  }

  .detail-title {
    font-size: 30px;
  }

  .detail-summary {
    font-size: 16px;
  }

  .detail-section {
    padding: 36px 0;
  }

  .detail-section-title {
    font-size: 24px;
  }
}

@media (max-width: 374px) {
  .detail-hero {
    min-height: 220px;
    padding: 36px 0;
  }

  .detail-title {
    font-size: 26px;
  }

  .detail-summary {
    font-size: 14px;
  }
}

.advantages-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 24px;
}

@media (max-width: 1023px) {
  .advantages-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (max-width: 767px) {
  .advantages-grid {
    grid-template-columns: repeat(2, 1fr);
    gap: 12px;
  }
}

@media (max-width: 374px) {
  .advantages-grid {
    grid-template-columns: 1fr;
  }
}

.advantage-card {
  text-align: center;
  padding: 24px 16px;
  border-radius: 12px;
  background: var(--bg-card, #fafafa);
  transition: transform 0.2s, box-shadow 0.2s;
}

.advantage-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.08);
}

.advantage-card:active {
  transform: translateY(0);
  box-shadow: none;
  background: var(--bg-gray);
}

@media (max-width: 767px) {
  .advantage-card {
    padding: 16px 12px;
  }

  .advantage-icon {
    width: 40px;
    height: 40px;
    margin-bottom: 8px;
  }

  .advantage-title {
    font-size: 14px;
  }

  .advantage-desc {
    font-size: 12px;
  }
}

.advantage-icon {
  margin-bottom: 12px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 48px;
  height: 48px;
  border-radius: 12px;
  background: rgba(200, 150, 62, 0.1);
}

.advantage-svg,
.advantage-svg-fallback {
  display: inline-flex;
  align-items: center;
  justify-content: center;
}

.advantage-title {
  font-size: 16px;
  font-weight: 600;
  color: var(--color-text);
  margin: 0 0 8px;
}

.advantage-desc {
  font-size: 13px;
  color: var(--color-text-muted, #666);
  line-height: 1.6;
  margin: 0;
}
</style>
