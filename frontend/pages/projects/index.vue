<template>
  <div class="projects-page">
    <div class="container">
      <ProjectBreadcrumb label="项目列表" />

      <h1 class="page-title">移民项目</h1>
      <p class="page-subtitle">探索最适合您的投资移民项目</p>

      <form class="search-form" @submit.prevent="searchProjects">
        <div class="search-field">
          <label for="project-country">国家</label>
          <input
            id="project-country"
            v-model="draftFilters.country"
            type="text"
            placeholder="请输入国家"
          />
        </div>
        <div class="search-field">
          <label for="project-name">项目名称</label>
          <input
            id="project-name"
            v-model="draftFilters.name"
            type="text"
            placeholder="请输入项目名称"
          />
        </div>
        <div class="search-actions">
          <button type="submit" class="btn-primary search-button" :disabled="searching">
            {{ searching ? '搜索中...' : '搜索' }}
          </button>
          <button type="button" class="reset-button" :disabled="searching" @click="resetSearch">
            重置
          </button>
        </div>
      </form>

      <div v-if="initialLoading" class="page-skeleton-wrapper">
        <PageSkeleton variant="cards" :count="6" />
      </div>

      <div v-else-if="loadError && items.length === 0" class="error-state">{{ loadError }}</div>

      <div v-else class="projects-grid">
        <ProjectCard
          v-for="(project, idx) in items"
          :key="project.slug"
          :slug="project.slug"
          :title="project.title"
          :description="project.description"
          :image="project.image"
          :features="project.features"
          :link="project.link"
          :image-variant="idx"
          :image-variants="project.imageVariants"
        />
      </div>

      <div v-if="(loadError || retryingLoadMore) && items.length > 0" class="inline-error">
        <span>{{ loadError || '正在重新加载...' }}</span>
        <button type="button" class="retry-button" :disabled="loadingMore" @click="retryLoadMore">
          {{ loadingMore ? '加载中...' : '重新加载' }}
        </button>
      </div>

      <div v-if="loadingMore" class="page-skeleton-wrapper" style="margin-top: 32px;">
        <PageSkeleton variant="cards" :count="3" />
      </div>

      <div v-if="allLoaded && items.length > 0" class="end-of-list">已加载全部项目</div>

      <div v-if="!initialLoading && items.length === 0 && !loadError" class="empty-state">
        {{ hasAppliedFilters ? '未找到符合条件的项目' : '暂无移民项目' }}
      </div>

      <div ref="sentinel" class="scroll-sentinel"></div>
    </div>
  </div>
</template>

<script setup lang="ts">
useSeo({ title: '移民项目', description: '探索最适合您的投资移民项目' })

import type { ImageVariantInfo } from '~/utils/image'

interface ApiProject {
  slug: string
  name: string
  tagline: string
  overview_text: string
  cover_image: string
  cover_image_variants?: Record<string, ImageVariantInfo>
  investment_amount: string
  processing_period: string
  target_crowd: string
}

interface ProjectsResponse {
  data: ApiProject[]
  pagination?: {
    page: number
    per_page: number
    total: number
  }
}

interface ProjectItem {
  slug: string
  title: string
  description: string
  image: string
  imageVariants?: Record<string, ImageVariantInfo>
  features: string[]
  link: string
}

const PER_PAGE = 12
const page = ref(1)
const items = ref<ProjectItem[]>([])
const totalCount = ref(0)
const loadingMore = ref(false)
const retryingLoadMore = ref(false)
const searching = ref(false)
const interactionError = ref<string | null>(null)
const hasLoaded = ref(false)
const sentinel = ref<HTMLElement | null>(null)
const draftFilters = reactive({ country: '', name: '' })
const appliedFilters = reactive({ country: '', name: '' })
let observer: IntersectionObserver | null = null
let requestToken = 0

function mapProject(api: ApiProject): ProjectItem {
  const features: string[] = []
  if (api.investment_amount) features.push(`投资金额：${api.investment_amount}`)
  if (api.processing_period) features.push(`办理周期：${api.processing_period}`)
  if (api.target_crowd) features.push(`适合人群：${api.target_crowd}`)
  return {
    slug: api.slug,
    title: api.name,
    description: api.tagline || api.overview_text || '',
    image: api.cover_image || '',
    imageVariants: api.cover_image_variants,
    features,
    link: `/projects/${api.slug}`,
  }
}

function buildProjectsUrl(targetPage: number) {
  const params = new URLSearchParams({
    page: String(targetPage),
    per_page: String(PER_PAGE),
  })
  if (appliedFilters.country) params.set('country', appliedFilters.country)
  if (appliedFilters.name) params.set('name', appliedFilters.name)
  return `/api/v1/projects?${params.toString()}`
}

usePublicDataFreshness([{ versionKey: 'public:projects:list', dataKey: 'public:projects:list:page1' }])

onMounted(() => {
  observer = new IntersectionObserver(
    (entries) => {
      if (entries[0]?.isIntersecting) loadMore()
    },
    { rootMargin: '200px' },
  )
  if (sentinel.value) observer.observe(sentinel.value)
})

onUnmounted(() => {
  observer?.disconnect()
})

const { data: firstPageResponse, pending, error: fetchError } = await useFetch<ProjectsResponse>(
  () => `/api/v1/projects?page=1&per_page=${PER_PAGE}`,
  {
    key: 'public:projects:list:page1',
  },
)

watch(firstPageResponse, (body) => {
  if (!body || requestToken !== 0) return
  items.value = body.data.map(mapProject)
  totalCount.value = body.pagination?.total ?? 0
  page.value = 1
  hasLoaded.value = true
}, { immediate: true })

const initialLoading = computed(() => ((requestToken === 0 && pending.value) || searching.value) && items.value.length === 0)
const allLoaded = computed(() => hasLoaded.value && items.value.length >= totalCount.value)
const hasAppliedFilters = computed(() => Boolean(appliedFilters.country || appliedFilters.name))
const loadError = computed(() => interactionError.value || (requestToken === 0 && fetchError.value ? '加载失败，请刷新重试' : null))

async function fetchFirstPage() {
  const token = ++requestToken
  searching.value = true
  loadingMore.value = false
  interactionError.value = null
  page.value = 1
  items.value = []
  totalCount.value = 0
  hasLoaded.value = false

  try {
    const raw = await $fetch(buildProjectsUrl(1))
    if (token !== requestToken) return
    const body = raw as any
    items.value = (body.data as ApiProject[]).map(mapProject)
    totalCount.value = body.pagination?.total ?? 0
    hasLoaded.value = true
  } catch {
    if (token === requestToken) interactionError.value = '加载失败，请稍后重试'
  } finally {
    if (token === requestToken) searching.value = false
  }
}

function searchProjects() {
  appliedFilters.country = draftFilters.country.trim()
  appliedFilters.name = draftFilters.name.trim()
  fetchFirstPage()
}

function resetSearch() {
  draftFilters.country = ''
  draftFilters.name = ''
  appliedFilters.country = ''
  appliedFilters.name = ''
  fetchFirstPage()
}

async function retryLoadMore() {
  retryingLoadMore.value = true
  await loadMore()
  retryingLoadMore.value = false
}

async function loadMore() {
  if (loadingMore.value || searching.value || !hasLoaded.value || allLoaded.value) return
  const token = requestToken
  interactionError.value = null
  loadingMore.value = true
  const nextPage = page.value + 1
  try {
    const raw = await $fetch(buildProjectsUrl(nextPage))
    if (token !== requestToken) return
    const body = raw as any
    const newItems = (body.data as ApiProject[]).map(mapProject)
    items.value.push(...newItems)
    totalCount.value = body.pagination?.total ?? totalCount.value
    page.value = nextPage
  } catch {
    if (token === requestToken) interactionError.value = '加载更多失败，请稍后重试'
  } finally {
    if (token === requestToken) loadingMore.value = false
  }
}

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
  color: var(--color-text-muted);
  margin-bottom: 24px;
}

.search-form {
  display: flex;
  align-items: flex-end;
  gap: 16px;
  margin-bottom: 40px;
  padding: 24px;
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  background: var(--color-bg-surface);
}

.search-field {
  display: flex;
  flex: 1;
  flex-direction: column;
  gap: 8px;
}

.search-field label {
  color: var(--color-text);
  font-size: 14px;
  font-weight: 600;
}

.search-field input {
  width: 100%;
  min-height: 44px;
  padding: 10px 14px;
  border: 1px solid var(--color-border);
  border-radius: var(--radius-sm);
  background: var(--color-bg-surface);
  color: var(--color-text);
  font: inherit;
  transition: border-color 0.2s, box-shadow 0.2s;
}

.search-field input:focus {
  border-color: var(--color-accent);
  box-shadow: var(--shadow-focus);
  outline: none;
}

.search-actions {
  display: flex;
  gap: 12px;
}

.search-button,
.reset-button {
  min-height: 44px;
  padding: 10px 24px;
  cursor: pointer;
  font: inherit;
  font-weight: 600;
}

.reset-button {
  border: 1px solid var(--color-border);
  border-radius: var(--radius-sm);
  background: var(--color-bg-surface);
  color: var(--color-text);
}

.search-button:disabled,
.reset-button:disabled {
  cursor: not-allowed;
  opacity: 0.6;
}

.search-button:disabled:hover {
  box-shadow: var(--shadow-gold);
  transform: none;
}

.projects-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 24px;
}

.error-state,
.empty-state {
  text-align: center;
  padding: 80px 20px;
  color: var(--color-text-muted);
  font-size: 16px;
}

.error-state {
  color: var(--color-danger);
}

.inline-error {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 12px;
  padding: 24px 20px 0;
  color: var(--color-danger);
  text-align: center;
}

.retry-button {
  min-height: 44px;
  padding: 8px 18px;
  border: 1px solid var(--color-danger);
  border-radius: var(--radius-sm);
  background: var(--color-bg-surface);
  color: var(--color-danger);
  cursor: pointer;
  font: inherit;
  font-weight: 600;
}

.retry-button:disabled {
  cursor: not-allowed;
  opacity: 0.6;
}

.end-of-list {
  text-align: center;
  padding: 40px 20px;
  color: var(--text-light);
  font-size: 14px;
}

.scroll-sentinel {
  height: 1px;
}

@media (max-width: 1023px) {
  .projects-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (max-width: 767px) {
  .search-form,
  .search-actions {
    align-items: stretch;
    flex-direction: column;
  }

  .search-form {
    padding: 20px;
  }

  .projects-grid {
    grid-template-columns: 1fr;
  }

  .page-title {
    font-size: 28px;
  }
}
</style>
