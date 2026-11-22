<template>
  <div class="projects-page">
    <ProjectBrowseHero
      :breadcrumbs="breadcrumbs"
      title="移民项目"
      description="按地区、名称或项目类型，缩小您的选择范围。"
      :images="heroImages"
    />

    <div class="project-tools">
      <div class="container">
        <ProjectCategoryDirectory label="按项目类型浏览">
          <template #leading>
            <form ref="searchForm" class="search-form" role="search" aria-label="查找项目" @submit.prevent="searchProjects">
              <div class="search-field">
                <label for="project-country">国家或地区</label>
                <input
                  v-if="countryOptionsError"
                  id="project-country"
                  v-model="draftFilters.country"
                  type="text"
                  placeholder="输入国家或地区"
                  aria-describedby="country-options-error"
                />
                <select
                  v-else
                  id="project-country"
                  v-model="draftFilters.country"
                  class="public-select"
                  :disabled="countryOptionsPending && !countries.length"
                >
                  <option value="">全部地区</option>
                  <option v-for="country in countries" :key="country" :value="country">{{ country }}</option>
                </select>
                <span v-if="countryOptionsError" id="country-options-error" class="search-help">
                  地区列表加载失败，可以手动输入。
                </span>
              </div>
              <div class="search-field">
                <label for="project-name">项目名称</label>
                <input
                  id="project-name"
                  v-model="draftFilters.name"
                  type="search"
                  placeholder="例如 EB-5 或黄金签证"
                />
              </div>
              <button type="submit" class="search-button" :disabled="searching">
                {{ searching ? '查找中…' : '查找项目' }}
              </button>
            </form>
          </template>
        </ProjectCategoryDirectory>
      </div>
    </div>

    <section
      ref="results"
      class="container project-results"
      aria-labelledby="project-results-heading"
      :aria-busy="initialLoading || loadingMore"
    >
      <header class="results-heading">
        <h2 id="project-results-heading" tabindex="-1">{{ hasAppliedFilters ? '查找结果' : '全部项目' }}</h2>
        <span v-if="hasLoaded" class="result-count" role="status">{{ totalCount }} 个项目</span>
        <span v-if="items.length" class="compare-hint">选择 2 个项目进行对比</span>
      </header>

      <div v-if="hasAppliedFilters" class="applied-filters">
        <span v-if="appliedFilters.country" class="applied-filter">{{ appliedFilters.country }}</span>
        <span v-if="appliedFilters.name" class="applied-filter">{{ appliedFilters.name }}</span>
        <button type="button" class="clear-button" @click="resetSearch">清除条件</button>
      </div>

      <PageSkeleton v-if="initialLoading" variant="cards" :count="6" />

      <div v-else-if="loadError && !hasLoaded" class="empty-state">
        <p role="status">{{ loadError }}</p>
        <button type="button" class="secondary-button" @click="fetchFirstPage">重试</button>
      </div>

      <div v-else-if="items.length" ref="grid" class="projects-grid">
        <article
          v-for="(project, index) in items"
          :key="project.slug"
          class="project-result"
        >
          <ProjectCard v-bind="project" :image-variant="index" class="result-card" />
          <label class="compare-choice">
            <input
              :id="`compare-choice-${project.slug}`"
              type="checkbox"
              :checked="isSelected(project.slug)"
              :disabled="isFull && !isSelected(project.slug)"
              :aria-label="`加入对比：${project.title}`"
              @change="toggle({ slug: project.slug, title: project.title })"
            />
            {{ isSelected(project.slug) ? '已加入对比' : '加入对比' }}
          </label>
        </article>
      </div>

      <div v-else class="empty-state">
        <p>{{ hasAppliedFilters ? '没有找到符合条件的项目' : '暂无移民项目' }}</p>
        <p v-if="hasAppliedFilters">可以更换地区或项目名称再试。</p>
      </div>

      <aside v-if="selection.length" class="compare-tray" aria-label="项目对比">
        <strong role="status">已选 {{ selection.length }} 个项目</strong>
        <ul class="compare-selection" aria-label="已选项目">
          <li v-for="project in selection" :key="project.slug">
            <button
              type="button"
              class="compare-chip"
              :aria-label="`移除 ${project.title}`"
              @click="removeComparisonChoice(project.slug)"
            >
              {{ project.title }} <span aria-hidden="true">×</span>
            </button>
          </li>
        </ul>
        <NuxtLink v-if="isFull" :to="compareTo" class="search-button compare-link">
          去对比 <span aria-hidden="true">→</span>
        </NuxtLink>
        <span v-else class="compare-hint">再选 1 个项目</span>
      </aside>

      <div v-if="hasMore" class="more-projects">
        <p v-if="loadError" role="status">{{ loadError }}</p>
        <button type="button" class="secondary-button" :disabled="loadingMore" @click="loadMore">
          {{ loadingMore ? '加载中…' : loadError ? '重新加载' : '查看更多项目' }}
          <span v-if="!loadingMore && !loadError" aria-hidden="true">↓</span>
        </button>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import { HOME_BREADCRUMB } from '~/utils/breadcrumb'
import type { ImageVariantInfo } from '~/utils/image'
import { projectRegionLabel } from '~/utils/project-browse'
import { toProjectCard, type ProjectCardItem, type ProjectCardSource } from '~/utils/project-card'

definePageMeta({
  middleware: (to, from) => {
    if (to.path === from.path || from.path.startsWith('/projects/')) return
    useState('project-overview:previous-browse').value = null
  },
})

const breadcrumbs = [
  HOME_BREADCRUMB,
  { label: '移民项目', href: '/projects' },
]
const description = '按国家或地区、项目名称与移民方式浏览北极星移民的项目，了解全球身份及居留方案。'

useSeo({ title: '移民项目', description, breadcrumbs })

interface CatalogProject extends ProjectCardSource {
  country?: string
}

interface ProjectsResponse {
  data: CatalogProject[] | null
  pagination?: {
    total: number
  }
}

interface ProjectFilters {
  country: string
  name: string
}

interface OverviewSnapshot {
  filters: ProjectFilters
  heroImages: { src: string; variants?: Record<string, ImageVariantInfo> }[]
  results: {
    items: ProjectCardItem[]
    page: number
    totalCount: number
  } | null
}

const PER_PAGE = 6
const REGION_BATCH_SIZE = 500
const versionKey = 'public:projects:list'
const countryOptionsKey = 'public:projects:list:countries'
const previousBrowse = useState<OverviewSnapshot | null>('project-overview:previous-browse', () => null)
const restored = previousBrowse.value
const draftFilters = reactive<ProjectFilters>(restored ? { ...restored.filters } : { country: '', name: '' })
const appliedFilters = reactive<ProjectFilters>({ ...draftFilters })
const items = ref<ProjectCardItem[]>(restored?.results?.items ?? [])
const page = ref(restored?.results?.page ?? 1)
const totalCount = ref(restored?.results?.totalCount ?? 0)
const hasLoaded = ref(Boolean(restored?.results))
const searching = ref(false)
const loadingMore = ref(false)
const interactionError = ref('')
const results = ref<HTMLElement | null>(null)
const grid = ref<HTMLElement | null>(null)
const searchForm = ref<HTMLFormElement | null>(null)
const hasAppliedFilters = computed(() => Boolean(appliedFilters.country || appliedFilters.name))
const { selection, isSelected, isFull, toggle, remove, compareTo } = useProjectComparison()
let requestToken = 0
let requestController: AbortController | undefined

function buildProjectsUrl(targetPage: number) {
  const params = new URLSearchParams({
    status: '1',
    page: String(targetPage),
    per_page: String(PER_PAGE),
  })
  if (appliedFilters.country) {
    params.set('country', appliedFilters.country === '中国香港' ? '香港' : appliedFilters.country)
  }
  if (appliedFilters.name) params.set('name', appliedFilters.name)
  return '/api/v1/projects?' + params.toString()
}

const initialUrl = buildProjectsUrl(1)
const initialKey = hasAppliedFilters.value
  ? versionKey + ':page1:' + encodeURIComponent(initialUrl)
  : versionKey + ':page1'

usePublicDataFreshness([
  { versionKey, dataKey: initialKey },
  { versionKey, dataKey: countryOptionsKey },
  { versionKey: 'public:project-regions', dataKey: countryOptionsKey },
])

onScopeDispose(() => {
  requestToken++
  requestController?.abort()
})

const { data: countryOptions, pending: countryOptionsPending, error: countryOptionsError } = useLazyAsyncData(
  countryOptionsKey,
  async (_nuxtApp, { signal }) => {
    let labels = await fetchProjectRegionNames(signal).catch((error) => {
      if (signal.aborted) throw error
      return null
    })
    if (!labels) {
      const { data } = await $fetch<{ data: { country?: string }[] | null }>('/api/v1/projects?status=1', { signal })
      labels = (data ?? []).map(project => projectRegionLabel(project.country)).filter(Boolean)
    }
    return [...new Set(labels)].sort((first, second) => first.localeCompare(second, 'zh-CN'))
  },
  { server: false },
)

async function fetchProjectRegionNames(signal: AbortSignal) {
  const names = new Set<string>()
  for (let optionPage = 1; ; optionPage++) {
    const options = await $fetch<{ data: { id: number }[]; pagination: { total: number } }>(
      '/api/v1/projects/options',
      { query: { page: optionPage, per_page: REGION_BATCH_SIZE }, signal },
    )
    if (!options.pagination || !Array.isArray(options.data)) return null
    if (!options.data.length) return options.pagination.total === 0 ? [] : null

    const ids = options.data.map(project => project.id)
    const { data } = await $fetch<{ data: { project_id: number; region: { name: string } | null }[] }>(
      '/api/v1/project-regions',
      { query: { ids: ids.join(',') }, signal },
    )
    const labels = new Map(data.map(row => [row.project_id, projectRegionLabel(row.region?.name)]))
    for (const id of ids) {
      const label = labels.get(id)
      // 存量项目可能只有 country 字段，需要回退到旧列表查询。
      if (!label) return null
      names.add(label)
    }
    if (optionPage * REGION_BATCH_SIZE >= options.pagination.total) return [...names]
  }
}

const countries = computed(() => {
  const choices = countryOptions.value ?? []
  const selected = draftFilters.country
  return selected && !choices.includes(selected) ? [...choices, selected] : choices
})

const { data: firstPageResponse, pending, error: fetchError } = await useFetch<ProjectsResponse>(
  initialUrl,
  { key: initialKey, immediate: !restored?.results },
)

const heroImages = computed(() => {
  if (restored) return restored.heroImages
  const projects = firstPageResponse.value?.data ?? []
  const project = projects.find(item => item.cover_image && item.slug === 'cies')
    ?? projects.find(item => item.cover_image && projectRegionLabel(item.country) === '中国香港')
    ?? projects.find(item => item.cover_image)
  return project?.cover_image
    ? [{ src: project.cover_image, variants: project.cover_image_variants }]
    : []
})

onBeforeRouteLeave((to) => {
  // 进入项目详情或类目后返回时，恢复结果高度和筛选条件。
  previousBrowse.value = to.path.startsWith('/projects/')
    ? {
        filters: { ...appliedFilters },
        heroImages: heroImages.value,
        results: hasLoaded.value ? { items: items.value, page: page.value, totalCount: totalCount.value } : null,
      }
    : null
})

const initialLoading = computed(() => searching.value || (!hasLoaded.value && requestToken === 0 && pending.value))
const hasMore = computed(() => hasLoaded.value && page.value * PER_PAGE < totalCount.value)
const loadError = computed(() => interactionError.value
  || (requestToken === 0 && fetchError.value && !hasLoaded.value ? '项目加载失败，请重试。' : ''))

function replaceResults(body: ProjectsResponse) {
  items.value = (body.data ?? []).map(toProjectCard)
  totalCount.value = body.pagination?.total ?? items.value.length
  page.value = 1
  hasLoaded.value = true
}

watch(firstPageResponse, (body) => {
  if (body && requestToken === 0) replaceResults(body)
}, { immediate: !restored?.results })

async function fetchFirstPage() {
  const token = ++requestToken
  requestController?.abort()
  requestController = new AbortController()
  searching.value = true
  loadingMore.value = false
  interactionError.value = ''
  page.value = 1
  items.value = []
  totalCount.value = 0
  hasLoaded.value = false

  try {
    const body = await $fetch<ProjectsResponse>(buildProjectsUrl(1), { signal: requestController.signal })
    if (token === requestToken) replaceResults(body)
  } catch {
    if (token === requestToken) interactionError.value = '项目加载失败，请重试。'
  } finally {
    if (token === requestToken) searching.value = false
  }
}

function searchProjects() {
  appliedFilters.country = projectRegionLabel(draftFilters.country)
  appliedFilters.name = draftFilters.name.trim()
  void fetchFirstPage()
  results.value?.scrollIntoView({ block: 'nearest' })
}

function resetSearch() {
  Object.assign(draftFilters, { country: '', name: '' })
  Object.assign(appliedFilters, draftFilters)
  void fetchFirstPage()
  searchForm.value?.scrollIntoView({ block: 'nearest' })
  searchForm.value?.querySelector<HTMLElement>('#project-country:not(:disabled), #project-name')
    ?.focus({ preventScroll: true })
}

async function removeComparisonChoice(slug: string) {
  remove(slug)
  await nextTick()
  const choice = document.getElementById(`compare-choice-${slug}`)
    ?? results.value?.querySelector<HTMLElement>('.compare-chip')
    ?? document.getElementById('project-results-heading')
  choice?.focus()
}

async function loadMore() {
  if (loadingMore.value || searching.value || !hasMore.value) return
  const token = requestToken
  const nextPage = page.value + 1
  const previousLength = items.value.length
  requestController = new AbortController()
  loadingMore.value = true
  interactionError.value = ''

  try {
    const body = await $fetch<ProjectsResponse>(buildProjectsUrl(nextPage), { signal: requestController.signal })
    if (token !== requestToken) return
    const existing = new Set(items.value.map(project => project.slug))
    const nextItems = (body.data ?? []).filter(project => !existing.has(project.slug))
    items.value.push(...nextItems.map(toProjectCard))
    totalCount.value = body.pagination?.total ?? items.value.length
    page.value = nextPage
    if (nextItems.length) {
      await nextTick()
      const firstNewCard = grid.value?.children[previousLength]?.querySelector<HTMLAnchorElement>('a')
      firstNewCard?.focus({ preventScroll: true })
    }
  } catch {
    if (token === requestToken) interactionError.value = '更多项目加载失败，请重试。'
  } finally {
    if (token === requestToken) loadingMore.value = false
  }
}
</script>

<style scoped>
.project-tools {
  padding-block: 16px;
  border-bottom: 1px solid var(--color-border);
  background: var(--color-bg-light);
}

.search-form {
  display: grid;
  grid-template-columns: clamp(110px, 20vw, 236px) minmax(0, 1fr) auto;
  align-items: end;
  gap: 8px;
  min-width: 0;
  scroll-margin-top: calc(var(--header-height) + 16px);
}

.search-field {
  position: relative;
  display: grid;
  gap: 4px;
  min-width: 0;
}

.search-field label {
  color: var(--color-text-secondary);
  font-size: var(--text-xs);
  font-weight: 600;
}

.search-field input {
  width: 100%;
  min-width: 0;
  min-height: 44px;
  padding: 8px 10px;
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  background: var(--color-bg-surface);
  color: var(--color-text);
  font: inherit;
}

.search-help {
  color: var(--color-text-secondary);
  font-size: var(--text-xs);
}

.search-button,
.secondary-button,
.clear-button {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  min-height: 44px;
  padding: 8px 14px;
  border-radius: var(--radius-sm);
  font: inherit;
  font-weight: 600;
  cursor: pointer;
}

.search-button,
.secondary-button {
  border: 1px solid var(--color-border);
  white-space: nowrap;
}

.search-button {
  border-color: var(--color-primary);
  background: var(--color-primary);
  color: #fff;
}

.search-button:hover:not(:disabled) {
  background: var(--color-primary-hover);
}

.search-button:disabled,
.secondary-button:disabled {
  opacity: 0.75;
  cursor: wait;
}

.secondary-button {
  background: var(--color-bg-surface);
  color: var(--color-primary);
}

.secondary-button:hover:not(:disabled) {
  border-color: var(--color-accent);
  background: var(--color-accent-soft);
}

.project-results {
  padding-block: clamp(24px, 4vw, 42px) clamp(48px, 7vw, 88px);
  scroll-margin-top: calc(var(--header-height) + 16px);
}

.results-heading,
.applied-filters {
  display: flex;
  align-items: baseline;
  flex-wrap: wrap;
  gap: 4px 12px;
  margin-bottom: 16px;
}

.results-heading h2 {
  font-size: var(--text-xl);
}

.result-count,
.compare-hint {
  color: var(--color-text-secondary);
  font-size: var(--text-sm);
}

.results-heading .compare-hint {
  margin-inline-start: auto;
}

.applied-filters {
  align-items: center;
}

.applied-filter {
  padding: 4px 10px;
  border-radius: var(--radius-sm);
  background: var(--color-accent-soft);
  overflow-wrap: anywhere;
}

.clear-button {
  border: 0;
  background: transparent;
  color: var(--color-text-secondary);
  text-decoration: underline;
  text-underline-offset: 3px;
}

.projects-grid {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: clamp(14px, 2vw, 24px);
}

.project-result {
  display: flex;
  min-width: 0;
  flex-direction: column;
  border: 1px solid var(--color-border);
  border-radius: var(--radius-lg);
  background: var(--color-bg-surface);
  box-shadow: var(--shadow-sm);
}

.project-result:hover,
.project-result:focus-within {
  border-color: var(--color-accent-dark);
}

.result-card,
.result-card:hover {
  flex: 1;
  height: auto;
  border: 0;
  border-radius: var(--radius-lg) var(--radius-lg) 0 0;
  box-shadow: none;
  transform: none;
}

.compare-choice {
  display: flex;
  align-items: center;
  align-self: stretch;
  gap: 8px;
  min-height: 44px;
  margin-inline: 20px;
  border-top: 1px solid var(--color-border);
  color: var(--color-text);
  font-size: var(--text-sm);
  font-weight: 600;
  cursor: pointer;
}

.compare-choice:has(input:disabled) {
  color: var(--color-text-muted);
  cursor: not-allowed;
}

.compare-choice input {
  width: 16px;
  height: 16px;
  accent-color: var(--color-primary);
  cursor: inherit;
}

.compare-tray {
  position: sticky;
  bottom: var(--public-contact-bar-height);
  z-index: 5;
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 8px 12px;
  margin-top: 20px;
  padding: 12px clamp(12px, 2vw, 20px);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  background: var(--color-bg-light);
}

.compare-selection {
  display: flex;
  flex: 1 1 16rem;
  flex-wrap: wrap;
  gap: 6px;
  min-width: 0;
}

.compare-selection li {
  min-width: 0;
  max-width: 100%;
}

.compare-chip {
  display: flex;
  align-items: center;
  gap: 8px;
  max-width: 100%;
  min-height: 44px;
  padding: 4px 12px;
  border: 1px solid var(--color-border);
  border-radius: 999px;
  background: var(--color-bg-surface);
  color: var(--color-text);
  font: inherit;
  text-align: start;
  overflow-wrap: anywhere;
  cursor: pointer;
}

.compare-chip:hover {
  border-color: var(--color-accent-dark);
  background: var(--color-accent-soft);
}

.compare-link {
  margin-inline-start: auto;
}

.empty-state,
.more-projects {
  display: grid;
  justify-items: center;
  gap: 10px;
  text-align: center;
  color: var(--color-text-secondary);
}

.empty-state {
  padding: 24px;
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
}

.more-projects {
  margin-top: 24px;
}

@media (max-width: 1023px) {
  .projects-grid {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
}

@media (max-width: 600px) {
  .search-form {
    width: 100%;
    grid-template-columns: minmax(0, 1fr) auto;
  }

  .search-field:first-child {
    grid-column: 1 / -1;
  }

  .search-field input {
    font-size: 16px;
  }

  .projects-grid {
    grid-template-columns: minmax(0, 1fr);
  }

  .results-heading .compare-hint {
    flex-basis: 100%;
    margin-inline-start: 0;
  }
}
</style>
