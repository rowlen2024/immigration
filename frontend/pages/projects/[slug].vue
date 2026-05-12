<template>
  <div class="project-detail">
    <div class="container">
      <ProjectBreadcrumb :label="project.title" />
    </div>

    <section class="detail-hero" :style="heroStyle">
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
      <div v-if="pending" class="loading-state">加载中...</div>
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
            <NuxtLink v-for="c in project.cases" :key="c.name" :to="'/case/' + c.slug" class="case-card">
              <img v-if="c.photo" :src="c.photo" :alt="c.name" class="case-photo" />
              <div class="case-body">
                <h4 class="case-name">{{ c.name }}</h4>
                <p class="case-meta">{{ c.country }} | {{ c.amount }} | {{ c.period }}</p>
                <p v-if="c.content" class="case-desc">{{ stripHtml(c.content) }}</p>
              </div>
            </NuxtLink>
          </div>
        </section>

        <section id="news" v-if="project.news.length > 0" class="detail-section">
          <h2 class="detail-section-title">最新资讯</h2>
          <div class="news-list">
            <NuxtLink v-for="n in project.news" :key="n.id" :to="`/pages/${n.slug}`" class="news-item">
              <img v-if="n.cover" :src="n.cover" :alt="n.title" class="news-cover" />
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
          <div v-else-if="compareData" class="compare-table-wrap">
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
        </section>

        <section id="faqs" v-if="faqs.length > 0" class="detail-section">
          <h2 class="detail-section-title">常见问题</h2>
          <ProjectFAQAccordion :items="faqs" />
        </section>

        <section class="detail-cta">
          <h3>对{{ project.title }}感兴趣？</h3>
          <p>立即联系我们，专业顾问为您一对一解答</p>
          <NuxtLink to="/contact" class="btn-primary">免费咨询</NuxtLink>
        </section>
      </template>
    </div>
  </div>
</template>

<script setup lang="ts">
import { getIconByName, getIconSvg } from '~/composables/lucideIcons';

function stripHtml(html: string): string {
  if (!html) return '';
  return html.replace(/<[^>]+>/g, '').replace(/&nbsp;/g, ' ').slice(0, 80);
}

const route = useRoute();
const slug = route.params.slug as string;

interface ApiRequirement { label: string; is_required: boolean; }
interface ApiCostItem { name: string; amount: string; note: string; }
interface ApiTimelinePhase { phase_number: number; title: string; description: string; duration: string; }
interface ApiFAQ { question: string; answer: string; }
interface ApiCase { slug: string; name: string; country_from: string; investment_amount: string; processing_period: string; content: string; photo_url: string; }
interface ApiNewsPage { id: number; title: string; slug: string; cover_image: string; created_at: string; }
interface ApiCompareConfig { compare_with: string[]; compare_fields: string[]; }

interface ApiAdvantage { icon: string; icon_type: string; title: string; description: string; }

interface ApiProject {
  name: string;
  tagline: string;
  country: string;
  cover_image: string;
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
  news: ApiNewsPage[];
  compare_config: ApiCompareConfig | null;
  advantages: ApiAdvantage[];
}

const { data, pending, error } = await useFetch<{ data: ApiProject }>(`/api/v1/projects/${slug}`);

const project = computed(() => {
  const p = data.value?.data;
  return {
    title: p?.name || slug,
    summary: p?.tagline || '',
    description: p?.overview_text || '',
    cover_image: p?.cover_image || '',
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
    })),
    news: (p?.news || []).map((n) => ({
      id: n.id,
      title: n.title,
      slug: n.slug,
      cover: n.cover_image,
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
} = useFetch<{ data: CompareTableData }>('/api/v1/projects/compare', {
  query: compareQuery,
});

const compareData = computed<CompareTableData | null>(() => {
  const raw = compareRaw.value as any;
  if (raw?.data?.rows) return raw.data;
  if (raw?.rows) return raw as CompareTableData;
  return null;
});

const compareError = computed(() =>
  compareErrorRaw.value ? '加载对比数据失败' : null
);

useSeo({
  title: project.value.title || '项目详情',
  description: project.value.summary || '',
  breadcrumbLabel: project.value.title,
});

const heroStyle = computed(() => {
  const img = project.value.cover_image;
  return img
    ? { backgroundImage: `linear-gradient(135deg, rgba(15, 36, 64, 0.85), rgba(26, 58, 92, 0.7)), url(${img})` }
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

onMounted(() => {
  observer = new IntersectionObserver(
    (entries) => {
      for (const entry of entries) {
        if (entry.isIntersecting) {
          activeTab.value = entry.target.id;
        }
      }
    },
    { rootMargin: `-${(tabNavRef.value?.offsetHeight || 0) + 64}px 0px -60% 0px` },
  );
  for (const tab of tabs.value) {
    const el = document.getElementById(tab.id);
    if (el) observer.observe(el);
  }
});

onUnmounted(() => {
  observer?.disconnect();
  observer = null;
});
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

.tab-btn.active {
  color: var(--primary);
  border-bottom-color: var(--primary);
  font-weight: 600;
}

.detail-hero {
  background-size: cover;
  background-position: center;
  padding: 80px 0;
  color: var(--bg-white);
  margin-bottom: 0;
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

.detail-cta {
  text-align: center;
  padding: 48px 0;
  background-color: var(--bg-light);
  border-radius: var(--radius-lg);
  margin: 40px 0;
}

.detail-cta h3 {
  font-size: 24px;
  font-weight: 700;
  color: var(--text-primary);
  margin-bottom: 8px;
}

.detail-cta p {
  font-size: 15px;
  color: var(--text-secondary);
  margin-bottom: 24px;
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

.case-card {
  background: var(--bg-white);
  border: 1px solid var(--border-color);
  border-radius: var(--radius-md);
  overflow: hidden;
}

.case-photo {
  width: 100%;
  height: 200px;
  object-fit: cover;
}

.case-body {
  padding: 16px;
}

.case-name {
  font-size: 18px;
  font-weight: 600;
  margin-bottom: 6px;
}

.case-meta {
  font-size: 13px;
  color: var(--text-light);
  margin-bottom: 8px;
}

.case-desc {
  font-size: 14px;
  color: var(--text-secondary);
  line-height: 1.6;
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

.news-cover {
  width: 120px;
  height: 80px;
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
  overflow-x: auto;
  -webkit-overflow-scrolling: touch;
}

.compare-table {
  width: 100%;
  border-collapse: collapse;
  font-size: 14px;
  min-width: max-content;
}

.compare-table thead {
  background-color: var(--primary);
  color: var(--bg-white);
}

.compare-table thead th {
  padding: 14px 16px;
  font-weight: 600;
  text-align: left;
  white-space: nowrap;
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
  min-width: 120px;
}

.compare-col-header {
  min-width: 180px;
}

.compare-table td {
  padding: 12px 16px;
  border-bottom: 1px solid var(--border-color);
  color: var(--text-secondary);
  line-height: 1.6;
}

.compare-table tbody tr:nth-child(even) {
  background-color: var(--bg-light);
}

.compare-label {
  font-weight: 600;
  color: var(--text-primary);
  white-space: nowrap;
  position: sticky;
  left: 0;
  z-index: 1;
  min-width: 120px;
}

.compare-table tbody tr:nth-child(even) .compare-label {
  background-color: var(--bg-light);
}

.compare-table tbody tr:nth-child(odd) .compare-label {
  background-color: var(--bg-white);
}

.compare-value {
  min-width: 180px;
}

.compare-requirements-grid {
  display: grid;
  grid-template-columns: 1fr;
  gap: 4px;
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

@media (max-width: 767px) {
  .detail-hero {
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

.advantages-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 24px;
}

@media (max-width: 992px) {
  .advantages-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (max-width: 576px) {
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
