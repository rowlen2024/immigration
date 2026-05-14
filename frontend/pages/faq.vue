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
      <div v-if="pending" class="loading-state">加载中...</div>
      <div v-else-if="error" class="error-state">加载常见问题失败，请稍后重试</div>
      <div v-else class="faq-list">
        <ProjectFAQAccordion :items="items" />
      </div>

      <div v-if="!pending && items.length === 0" class="empty-state">
        暂无该分类的常见问题
      </div>

      <Pagination
        :page="page"
        :per-page="perPage"
        :total="totalItems"
        @change="changePage"
      />

      <!-- CTA -->
      <section class="faq-cta">
        <h3>没有找到您的问题？</h3>
        <p>联系我们的专业顾问，获取一对一解答</p>
        <NuxtLink to="/contact" class="btn-primary">免费咨询</NuxtLink>
      </section>
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

const items = ref<FaqItem[]>([]);
const totalItems = ref(0);
const pending = ref(true);
const error = ref<string | null>(null);

const fetchFaqs = async () => {
  pending.value = true;
  error.value = null;
  try {
    const params = new URLSearchParams({
      page: String(page.value),
      per_page: String(perPage),
    });
    if (activeFilter.value !== 'all') {
      const filter = projectFilters.value.find(f => f.slug === activeFilter.value);
      if (filter) {
        params.set('project_id', String(filter.id));
      }
    }
    const res = await $fetch<any>(`/api/v1/faqs?${params.toString()}`);
    items.value = res?.data ?? [];
    totalItems.value = res?.pagination?.total ?? 0;
  } catch {
    error.value = '加载常见问题失败，请稍后重试';
  } finally {
    pending.value = false;
  }
};

// Fetch project list for filter buttons.
const projectFilters = ref<{ id: number; slug: string; label: string }[]>([]);
const fetchProjects = async () => {
  try {
    const res = await $fetch<any>('/api/v1/projects?per_page=100');
    const list = res?.data ?? [];
    projectFilters.value = list.map((p: any) => ({ id: p.id, slug: p.slug, label: p.name }));
  } catch { /* keep filters empty if projects API fails */ }
};

const changeFilter = (slug: string) => {
  activeFilter.value = slug;
  page.value = 1;
};

const changePage = (p: number) => {
  page.value = p;
  window.scrollTo({ top: 0, behavior: 'smooth' });
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

watch([activeFilter, page], () => {
  fetchFaqs();
});

onMounted(() => {
  fetchFaqs();
  fetchProjects();
});
</script>

<style scoped>
.page-title {
  font-size: 36px;
  font-weight: 700;
  color: var(--text-primary);
  margin-bottom: 12px;
}

.page-subtitle {
  font-size: 16px;
  color: var(--text-light);
  margin-bottom: 32px;
}

.faq-filters {
  display: flex;
  gap: 12px;
  margin-bottom: 40px;
  flex-wrap: wrap;
}

.filter-btn {
  padding: 10px 24px;
  font-size: 14px;
  font-weight: 600;
  font-family: var(--font-sans);
  background-color: var(--bg-light);
  color: var(--text-secondary);
  border: 2px solid var(--border-color);
  border-radius: var(--radius-md);
  cursor: pointer;
  transition: all 0.3s ease;
}

.filter-btn:hover {
  border-color: var(--accent);
  color: var(--accent-dark);
}

.filter-btn.active {
  background-color: var(--primary);
  color: var(--bg-white);
  border-color: var(--primary);
}

.faq-list {
  margin-bottom: 48px;
}

.faq-cta {
  text-align: center;
  padding: 48px;
  background-color: var(--bg-light);
  border-radius: var(--radius-lg);
  margin-bottom: 48px;
}

.faq-cta h3 {
  font-size: 24px;
  font-weight: 700;
  color: var(--text-primary);
  margin-bottom: 8px;
}

.faq-cta p {
  font-size: 15px;
  color: var(--text-secondary);
  margin-bottom: 24px;
}

.loading-state,
.error-state,
.empty-state {
  text-align: center;
  padding: 60px 20px;
  color: var(--text-light);
  font-size: 16px;
}

.error-state {
  color: #c62828;
}

@media (max-width: 767px) {
  .page-title {
    font-size: 28px;
  }
}
</style>
