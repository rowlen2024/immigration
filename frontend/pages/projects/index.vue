<template>
  <div class="projects-page">
    <div class="container">
      <ProjectBreadcrumb label="项目列表" />

      <h1 class="page-title">移民项目</h1>
      <p class="page-subtitle">探索最适合您的投资移民项目</p>

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

      <div v-if="loadingMore" class="page-skeleton-wrapper" style="margin-top: 32px;">
        <PageSkeleton variant="cards" :count="3" />
      </div>

      <div v-if="allLoaded && items.length > 0" class="end-of-list">已加载全部项目</div>

      <div v-if="!initialLoading && items.length === 0 && !loadError" class="empty-state">
        暂无移民项目
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
const loadError = ref<string | null>(null)
const sentinel = ref<HTMLElement | null>(null)
let observer: IntersectionObserver | null = null

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

const { pending, error: fetchError } = await useFetch(
  () => `/api/v1/projects?page=1&per_page=${PER_PAGE}`,
  {
    onResponse({ response }) {
      const body = response._data as any
      if (body?.data) {
        items.value = (body.data as ApiProject[]).map(mapProject)
        totalCount.value = body.pagination?.total ?? 0
      }
    },
  },
)

const initialLoading = computed(() => pending.value && items.value.length === 0)
const allLoaded = computed(() => items.value.length >= totalCount.value && totalCount.value > 0)
const computedError = computed(() => (fetchError.value ? '加载失败，请刷新重试' : null))
watchEffect(() => { loadError.value = computedError.value })

async function loadMore() {
  if (loadingMore.value || allLoaded.value) return
  loadingMore.value = true
  const nextPage = page.value + 1
  try {
    const raw = await $fetch(`/api/v1/projects?page=${nextPage}&per_page=${PER_PAGE}`)
    const body = raw as any
    const newItems = (body.data as ApiProject[]).map(mapProject)
    items.value.push(...newItems)
    totalCount.value = body.pagination?.total ?? totalCount.value
    page.value = nextPage
  } catch {
    // silent retry on next scroll
  } finally {
    loadingMore.value = false
  }
}

onMounted(() => {
  $fetch(`/api/v1/projects?page=1&per_page=${PER_PAGE}`)
    .then((v: any) => {
      items.value = (v.data as ApiProject[]).map(mapProject)
      totalCount.value = v.pagination?.total ?? 0
    })
    .catch(() => {})

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
  margin-bottom: 40px;
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
  .projects-grid {
    grid-template-columns: 1fr;
  }

  .page-title {
    font-size: 28px;
  }
}
</style>
