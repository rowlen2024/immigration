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

    <div class="container">
      <div v-if="pending" class="loading-state">加载中...</div>
      <div v-else-if="error" class="error-state">{{ error }}</div>
      <template v-else>
        <section class="detail-section">
          <h2 class="detail-section-title">项目概览</h2>
          <ProjectQuickFacts :facts="quickFacts" />
        </section>

        <section v-if="project.description" class="detail-section">
          <h2 class="detail-section-title">项目介绍</h2>
          <div class="detail-content">
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

        <section v-if="faqs.length > 0" class="detail-section">
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
const route = useRoute();
const slug = route.params.slug as string;

interface ApiRequirement { label: string; is_required: boolean; }
interface ApiCostItem { name: string; amount: string; note: string; }
interface ApiTimelinePhase { phase_number: number; title: string; description: string; duration: string; }
interface ApiFAQ { question: string; answer: string; }

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
  };
});

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
</script>

<style scoped>
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
</style>
