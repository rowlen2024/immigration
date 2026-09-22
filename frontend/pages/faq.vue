<template>
  <div class="faq-page">
    <div class="container">
      <ProjectBreadcrumb :items="breadcrumbs" />

      <h1 class="page-title">北极星移民投资移民常见问题</h1>
      <p class="page-subtitle">北极星移民为您解答投资移民的申请条件、费用、办理周期与项目选择等常见问题。</p>

      <!-- 项目筛选 -->
      <div class="faq-project-picker">
        <ProjectPicker
          :projects="projectFilters"
          :model-value="activeFilter"
          subtitle="按项目查看对应常见问题"
          @update:model-value="changeFilter"
        />
      </div>

      <!-- FAQ List -->
      <ProjectFaqListSection
        :items="items"
        :loading="pending"
        :error="faqError"
        empty-text="暂无该分类的常见问题"
        :page="page"
        :per-page="perPage"
        :total="totalItems"
        @page-change="changePage"
      />

      <!-- CTA -->
      <ConsultCTA
        title="没有找到您的问题？"
        description="联系我们的专业顾问，获取一对一解答"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { HOME_BREADCRUMB } from '~/utils/breadcrumb'

/** 常见问题页面的稳定面包屑路径。 */
const breadcrumbs = [
  HOME_BREADCRUMB,
  { label: '常见问题', href: '/faq' },
]

useSeo({
  title: '常见问题',
  description: '北极星移民常见问题解答，涵盖美国EB-5、香港投资移民、巴拿马购房移民等投资移民相关问题。',
  breadcrumbs,
});

const page = ref(1);
const perPage = 10;
const activeFilter = ref('all');
const faqListKey = computed(() => `public:faqs:list:${activeFilter.value}:${page.value}`)

usePublicDataFreshness(() => [
  { versionKey: 'public:faqs:list', dataKey: faqListKey.value },
  'public:faqs:projects',
])

interface FaqItem {
  id: number;
  question: string;
  answer: string;
  project_id: number | null;
  project_name: string;
  project_slug: string;
  is_global: boolean;
  sort_order: number;
}

// 获取筛选按钮所用的项目列表
const { data: projectsRaw } = await useFetch('/api/v1/faqs/projects', {
  key: 'public:faqs:projects',
  onResponseError() { /* keep filters empty if API fails */ },
})

const projectFilters = computed(() => {
  const list = (projectsRaw.value as any)?.data ?? []
  return list.map((p: any) => ({
    id: p.id,
    slug: p.slug,
    name: p.name,
  }))
})

/** FAQ 项目筛选项。 */
type ProjectFilter = {
  id: number;
  slug: string;
  name: string;
}

// 用 getter 函数使 useFetch 在 page/activeFilter 变化时自动重新请求
const { data: faqRaw, pending, error: fetchError, refresh } = await useFetch(
  () => {
    const params = new URLSearchParams({
      page: String(page.value),
      per_page: String(perPage),
    })
    if (activeFilter.value !== 'all') {
      const filter = projectFilters.value.find((f: any) => f.slug === activeFilter.value)
      if (filter) params.set('project_id', String(filter.id))
    }
    return `/api/v1/faqs?${params.toString()}`
  },
  {
    key: faqListKey,
    onResponseError() {
      // error handled via computed
    },
  }
)

const items = computed<FaqItem[]>(() => {
  const raw = faqRaw.value as any
  return raw?.data ?? []
})

const totalItems = computed(() => {
  const raw = faqRaw.value as any
  return raw?.pagination?.total ?? 0
})

const faqError = computed(() => {
  return fetchError.value ? '加载常见问题失败，请稍后重试' : null
})

const changeFilter = (slug: string) => {
  activeFilter.value = slug;
  page.value = 1;
};

const changePage = (p: number) => {
  page.value = p;
  const prefersReduced = window.matchMedia('(prefers-reduced-motion: reduce)').matches;
  window.scrollTo({ top: 0, behavior: prefersReduced ? 'instant' : 'smooth' });
};

import { buildFAQPageJsonLd, toJsonLdScripts } from '~/utils/jsonld'

// FAQPage structured data
useHead(() => {
  return {
    script: toJsonLdScripts(buildFAQPageJsonLd(items.value)),
  }
})

</script>

<style scoped>
.page-title {
  font-size: var(--text-3xl);
  font-weight: 700;
  color: var(--color-text);
  margin-bottom: 12px;
}

.page-subtitle {
  font-size: var(--text-base);
  color: var(--color-text-secondary);
  margin-bottom: 32px;
}

.faq-project-picker {
  margin-bottom: 40px;
}

/* styles moved to ProjectFaqListSection / ConsultCTA components */

@media (max-width: 767px) {
  .page-title {
    font-size: 28px;
  }

  .page-subtitle {
    font-size: 14px;
    margin-bottom: 24px;
  }

  .faq-project-picker {
    margin-bottom: 28px;
  }
}
</style>
