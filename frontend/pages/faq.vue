<template>
  <div class="faq-page">
    <div class="container">
      <ProjectBreadcrumb />

      <h1 class="page-title">常见问题</h1>
      <p class="page-subtitle">关于投资移民的常见问题解答</p>

      <!-- Filter Buttons -->
      <div ref="projectPickerRef" class="project-picker">
        <div class="faq-filters">
          <button
            type="button"
            class="filter-btn"
            :class="{ active: activeFilter === 'all' }"
            @click="changeFilter('all')"
          >
            全部
          </button>
          <button
            v-for="filter in visibleProjectFilters"
            :key="filter.slug"
            type="button"
            class="filter-btn"
            :class="{ active: activeFilter === filter.slug }"
            @click="changeFilter(filter.slug)"
          >
            {{ filter.label }}
          </button>
          <button
            v-if="selectedHiddenFilter"
            type="button"
            class="filter-btn active selected-hidden-filter"
            @click="changeFilter(selectedHiddenFilter.slug)"
          >
            {{ selectedHiddenFilter.label }}
          </button>
          <button
            v-if="hiddenProjectFilters.length > 0"
            type="button"
            class="filter-btn more-projects-btn"
            :class="{ open: showProjectPicker }"
            :aria-expanded="showProjectPicker"
            aria-controls="faq-project-picker-panel"
            @click="toggleProjectPicker"
          >
            更多项目
            <span class="more-projects-count">{{ hiddenProjectFilters.length }}</span>
          </button>
        </div>

        <div
          v-if="showProjectPicker"
          id="faq-project-picker-panel"
          class="project-picker-panel"
        >
          <div class="project-picker-handle" aria-hidden="true"></div>
          <div class="project-picker-header">
            <div>
              <h2 class="project-picker-title">选择项目</h2>
              <p class="project-picker-subtitle">按项目查看对应常见问题</p>
            </div>
            <button
              type="button"
              class="project-picker-close"
              aria-label="关闭项目筛选"
              @click="closeProjectPicker"
            >
              ×
            </button>
          </div>
          <input
            v-model.trim="projectSearch"
            class="project-picker-search"
            type="search"
            placeholder="搜索项目名称"
          >
          <div v-if="filteredHiddenProjectFilters.length > 0" class="project-picker-grid">
            <button
              v-for="filter in filteredHiddenProjectFilters"
              :key="filter.slug"
              type="button"
              class="project-picker-option"
              :class="{ active: activeFilter === filter.slug }"
              @click="changeFilter(filter.slug)"
            >
              <span>{{ filter.label }}</span>
              <span v-if="activeFilter === filter.slug" class="project-picker-check">✓</span>
            </button>
          </div>
          <div v-else class="project-picker-empty">未找到相关项目</div>
        </div>
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
const showProjectPicker = ref(false);
const projectSearch = ref('');
const projectPickerRef = ref<HTMLElement | null>(null);
const visibleProjectLimit = 5;
const faqListKey = computed(() => `public:faqs:list:${activeFilter.value}:${page.value}`)

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

const toggleProjectPicker = () => {
  showProjectPicker.value = !showProjectPicker.value;
};

const closeProjectPicker = () => {
  showProjectPicker.value = false;
  projectSearch.value = '';
};

const handleProjectPickerClickOutside = (event: MouseEvent) => {
  if (!showProjectPicker.value) return;
  const target = event.target as Node;
  if (projectPickerRef.value && !projectPickerRef.value.contains(target)) {
    closeProjectPicker();
  }
};

const handleProjectPickerKeydown = (event: KeyboardEvent) => {
  if (event.key === 'Escape') {
    closeProjectPicker();
  }
};

onMounted(() => {
  document.addEventListener('click', handleProjectPickerClickOutside)
  document.addEventListener('keydown', handleProjectPickerKeydown)
})

onBeforeUnmount(() => {
  document.removeEventListener('click', handleProjectPickerClickOutside)
  document.removeEventListener('keydown', handleProjectPickerKeydown)
})

// 获取筛选按钮所用的项目列表
const { data: projectsRaw, refresh: refreshProjects } = await useFetch('/api/v1/faqs/projects', {
  key: 'public:faqs:projects',
  onResponseError() { /* keep filters empty if API fails */ },
})

const projectFilters = computed(() => {
  const list = (projectsRaw.value as any)?.data ?? []
  return list.map((p: any) => ({ id: p.id, slug: p.slug, label: p.name }))
})

type ProjectFilter = {
  id: number;
  slug: string;
  label: string;
}

const visibleProjectFilters = computed<ProjectFilter[]>(() => {
  return projectFilters.value.slice(0, visibleProjectLimit)
})

const hiddenProjectFilters = computed<ProjectFilter[]>(() => {
  return projectFilters.value.slice(visibleProjectLimit)
})

const selectedHiddenFilter = computed<ProjectFilter | null>(() => {
  if (activeFilter.value === 'all') return null
  const inVisible = visibleProjectFilters.value.some((filter) => filter.slug === activeFilter.value)
  if (inVisible) return null
  return hiddenProjectFilters.value.find((filter) => filter.slug === activeFilter.value) ?? null
})

const filteredHiddenProjectFilters = computed<ProjectFilter[]>(() => {
  const keyword = projectSearch.value.trim().toLowerCase()
  if (!keyword) return hiddenProjectFilters.value
  return hiddenProjectFilters.value.filter((filter) => {
    return filter.label.toLowerCase().includes(keyword) || filter.slug.toLowerCase().includes(keyword)
  })
})

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
usePublicDataFreshness(() => [
  { versionKey: 'public:faqs:list', dataKey: faqListKey.value },
  'public:faqs:projects',
])

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
  closeProjectPicker();
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

.project-picker {
  position: relative;
  margin-bottom: 40px;
}

.faq-filters {
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
  align-items: center;
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
  gap: 8px;
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

.selected-hidden-filter {
  box-shadow: 0 4px 14px rgba(191, 154, 96, 0.22);
}

.more-projects-btn {
  background-color: var(--bg-light);
  color: var(--color-text);
}

.more-projects-btn.open,
.more-projects-btn:hover {
  border-color: var(--color-accent);
  color: var(--color-accent);
}

.more-projects-count {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-width: 22px;
  height: 22px;
  padding: 0 7px;
  border-radius: var(--radius-full);
  background-color: rgba(191, 154, 96, 0.12);
  color: var(--color-accent);
  font-size: 12px;
  line-height: 1;
}

.project-picker-panel {
  position: absolute;
  top: calc(100% + 12px);
  left: 0;
  z-index: 20;
  width: min(720px, 100%);
  padding: 18px;
  background-color: var(--bg-white);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-lg);
  box-shadow: 0 24px 60px rgba(32, 32, 32, 0.14);
}

.project-picker-handle {
  display: none;
}

.project-picker-header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 16px;
  margin-bottom: 14px;
}

.project-picker-title {
  margin: 0 0 4px;
  font-size: 16px;
  font-weight: 700;
  color: var(--color-text);
}

.project-picker-subtitle {
  margin: 0;
  font-size: 13px;
  color: var(--color-text-secondary);
}

.project-picker-close {
  width: 32px;
  height: 32px;
  border: 0;
  border-radius: var(--radius-full);
  background-color: var(--bg-light);
  color: var(--color-text-secondary);
  cursor: pointer;
  font-size: 20px;
  line-height: 1;
}

.project-picker-close:hover {
  color: var(--color-accent);
}

.project-picker-search {
  width: 100%;
  height: 42px;
  margin-bottom: 14px;
  padding: 0 14px;
  border: 1.5px solid var(--color-border);
  border-radius: var(--radius-md);
  background-color: var(--bg-white);
  color: var(--color-text);
  font-size: 14px;
  font-family: var(--font-sans);
  outline: none;
}

.project-picker-search:focus {
  border-color: var(--color-accent);
  box-shadow: 0 0 0 3px rgba(191, 154, 96, 0.12);
}

.project-picker-grid {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 8px;
  max-height: 280px;
  overflow-y: auto;
  padding-right: 2px;
}

.project-picker-option {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 10px;
  min-height: 44px;
  padding: 10px 12px;
  border: 1.5px solid transparent;
  border-radius: var(--radius-md);
  background-color: var(--bg-white);
  color: var(--color-text);
  cursor: pointer;
  font-size: 13px;
  font-family: var(--font-sans);
  text-align: left;
  transition: all var(--duration-fast) var(--ease-out);
}

.project-picker-option:hover {
  background-color: var(--bg-light);
  border-color: var(--color-border);
}

.project-picker-option.active {
  background-color: rgba(191, 154, 96, 0.12);
  border-color: rgba(191, 154, 96, 0.45);
  color: var(--color-accent);
  font-weight: 600;
}

.project-picker-check {
  flex-shrink: 0;
}

.project-picker-empty {
  min-height: 88px;
  display: flex;
  align-items: center;
  justify-content: center;
  border: 1px dashed var(--color-border);
  border-radius: var(--radius-md);
  color: var(--color-text-secondary);
  font-size: 14px;
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
  }

  .filter-btn {
    padding: 8px 16px;
    min-height: 44px;
    font-size: 12px;
  }

  .project-picker {
    margin-bottom: 28px;
  }

  .project-picker-panel {
    position: fixed;
    top: auto;
    left: 0;
    right: 0;
    bottom: 0;
    width: 100%;
    max-height: 78vh;
    padding: 12px 16px 18px;
    border-right: 0;
    border-bottom: 0;
    border-left: 0;
    border-radius: 24px 24px 0 0;
    box-shadow: 0 -18px 45px rgba(32, 32, 32, 0.18);
  }

  .project-picker-handle {
    display: block;
    width: 44px;
    height: 4px;
    margin: 0 auto 14px;
    border-radius: var(--radius-full);
    background-color: var(--color-border);
  }

  .project-picker-grid {
    grid-template-columns: 1fr;
    max-height: calc(78vh - 164px);
  }

  .project-picker-title {
    font-size: 15px;
  }

  .project-picker-subtitle {
    font-size: 12px;
  }

}
</style>
