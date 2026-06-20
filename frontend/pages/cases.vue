<template>
  <div class="cases-page">
    <div class="container">
      <ProjectBreadcrumb />

      <h1 class="page-title">成功案例</h1>
      <p class="page-subtitle">每一位客户的成功获批，都是我们最大的骄傲</p>

      <div v-if="initialLoading" class="page-skeleton-wrapper"><PageSkeleton variant="cards" :count="6" /></div>

      <div v-else-if="loadError && items.length === 0" class="error-state">{{ loadError }}</div>

      <div v-else class="cases-grid">
        <CaseCard
          v-for="item in items"
          :key="item.id"
          :slug="item.slug"
          :name="item.name"
          :country="item.country"
          :project="item.project || undefined"
          :summary="item.summary"
          :image="item.image"
          :image-variants="item.imageVariants"
          :show-result="true"
          result-text="成功获批"
        />
      </div>

      <div v-if="loadingMore" class="page-skeleton-wrapper" style="margin-top: 32px;">
        <PageSkeleton variant="cards" :count="3" />
      </div>

      <div v-if="allLoaded && items.length > 0" class="end-of-list">已加载全部案例</div>

      <div v-if="!initialLoading && items.length === 0 && !loadError" class="empty-state">
        暂无成功案例展示
      </div>

      <div ref="sentinel" class="scroll-sentinel"></div>
    </div>
  </div>
</template>

<script setup lang="ts">
useSeo({ title: '成功案例' })

import type { ImageVariantInfo } from '~/utils/image'

interface ApiCaseItem {
  id: number
  slug: string
  name: string
  country_from: string
  photo_url: string
  photo_variants?: Record<string, ImageVariantInfo>
  content: string
  project?: { name: string }
}

interface CaseItem {
  id: string
  slug: string
  name: string
  country: string
  project: string
  summary: string
  image: string
  imageVariants?: Record<string, ImageVariantInfo>
}

import { stripHtml } from '~/utils/html'

const PER_PAGE = 12
const page = ref(1)
const items = ref<CaseItem[]>([])
const totalCount = ref(0)
const loadingMore = ref(false)
const loadError = ref<string | null>(null)
const sentinel = ref<HTMLElement | null>(null)
let observer: IntersectionObserver | null = null

function mapCase(api: ApiCaseItem): CaseItem {
  return {
    id: String(api.id),
    slug: api.slug,
    name: api.name,
    country: api.country_from,
    project: api.project?.name ?? '',
    summary: stripHtml(api.content),
    image: api.photo_url,
    imageVariants: api.photo_variants,
  }
}

const { pending, error: fetchError } = await useFetch(
  () => `/api/v1/cases?page=1&per_page=${PER_PAGE}`,
  {
    key: 'public:cases:list:page1',
    onResponse({ response }) {
      const body = response._data as any
      if (body?.data) {
        items.value = (body.data as ApiCaseItem[]).map(mapCase)
        totalCount.value = body.pagination?.total ?? 0
      }
    },
  },
)
usePublicDataFreshness([{ versionKey: 'public:cases:list', dataKey: 'public:cases:list:page1' }])

const initialLoading = computed(() => pending.value && items.value.length === 0)
const allLoaded = computed(() => items.value.length >= totalCount.value && totalCount.value > 0)
const computedError = computed(() => (fetchError.value ? '加载失败，请刷新重试' : null))
watchEffect(() => { loadError.value = computedError.value })

async function loadMore() {
  if (loadingMore.value || allLoaded.value) return
  loadingMore.value = true
  const nextPage = page.value + 1
  try {
    const raw = await $fetch(`/api/v1/cases?page=${nextPage}&per_page=${PER_PAGE}`)
    const body = raw as any
    const newItems = (body.data as ApiCaseItem[]).map(mapCase)
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

.cases-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 32px;
  margin-bottom: 0;
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
  .cases-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (max-width: 767px) {
  .cases-grid {
    grid-template-columns: 1fr;
  }

  .page-title {
    font-size: 28px;
  }
}
</style>
