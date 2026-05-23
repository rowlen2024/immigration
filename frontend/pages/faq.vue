<template>
  <div class="faq-page">
    <div class="container">
      <ProjectBreadcrumb />

      <h1 class="page-title">常见问题</h1>
      <p class="page-subtitle">关于投资移民的常见问题解答</p>

      <!-- Filter Buttons -->
      <div class="faq-filters">
        <button
          class="filter-btn"
          :class="{ active: activeFilter === 'all' }"
          @click="changeFilter('all')"
        >
          全部
        </button>
        <button
          v-for="filter in projectFilters"
          :key="filter.slug"
          class="filter-btn"
          :class="{ active: activeFilter === filter.slug }"
          @click="changeFilter(filter.slug)"
        >
          {{ filter.label }}
        </button>
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
useSeo({
  title: '常见问题',
  description: '北极星移民常见问题解答，涵盖美国EB-5、香港投资移民、巴拿马购房移民等投资移民相关问题。',
});

const page = ref(1);
const perPage = 10;
const activeFilter = ref('all');

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
const { data: projectsRaw, refresh: refreshProjects } = await useFetch('/api/v1/faqs/projects', {
  onResponseError() { /* keep filters empty if API fails */ },
})

const projectFilters = computed(() => {
  const list = (projectsRaw.value as any)?.data ?? []
  return list.map((p: any) => ({ id: p.id, slug: p.slug, label: p.name }))
})

// 用 getter 函数使 useFetch 在 page/activeFilter 变化时自动重新请求
const { data: faqRaw, pending, error: fetchError, refresh } = await useFetch(
  () => {
    const params = new URLSearchParams({
      page: String(page.value),
      per_page: String(perPage),
    })
    if (activeFilter.value !== 'all') {
      const filter = projectFilters.value.find(f => f.slug === activeFilter.value)
      if (filter) params.set('project_id', String(filter.id))
    }
    return `/api/v1/faqs?${params.toString()}`
  },
  {
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

// FAQPage structured data
useHead(() => {
  if (items.value.length === 0) return {};
  return {
    script: [
      {
        type: 'application/ld+json',
        innerHTML: JSON.stringify({
          '@context': 'https://schema.org',
          '@type': 'FAQPage',
          mainEntity: items.value.map((faq) => ({
            '@type': 'Question',
            name: faq.question,
            acceptedAnswer: {
              '@type': 'Answer',
              text: faq.answer,
            },
          })),
        }),
      },
    ],
  };
});

// 客户端刷新确保数据最新
onMounted(() => {
  $fetch('/api/v1/faqs?page=1&per_page=10').then(v => { faqRaw.value = v }).catch(() => {})
  $fetch('/api/v1/faqs/projects').then(v => { projectsRaw.value = v }).catch(() => {})
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

.faq-filters {
  display: flex;
  gap: 12px;
  margin-bottom: 40px;
  flex-wrap: nowrap;
  overflow-x: auto;
  -webkit-overflow-scrolling: touch;
  touch-action: pan-x;
  scrollbar-width: none;
  padding-bottom: 4px;
  position: relative;
}

.faq-filters::-webkit-scrollbar {
  display: none;
}

.filter-btn {
  padding: 10px 20px;
  min-height: 44px;
  font-size: 13px;
  font-weight: 500;
  font-family: var(--font-sans);
  background-color: var(--bg-white);
  color: var(--color-text-secondary);
  border: 1.5px solid var(--color-border);
  border-radius: var(--radius-full);
  cursor: pointer;
  transition: all var(--duration-fast) var(--ease-out);
  white-space: nowrap;
  flex-shrink: 0;
  display: inline-flex;
  align-items: center;
}

.filter-btn:hover {
  border-color: var(--color-accent);
  color: var(--color-accent);
}

.filter-btn.active {
  background-color: var(--color-accent);
  color: #fff;
  border-color: var(--color-accent);
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

  .faq-filters {
    gap: 8px;
    margin-bottom: 28px;
  }

  .filter-btn {
    padding: 8px 16px;
    min-height: 44px;
    font-size: 12px;
  }

}
</style>
