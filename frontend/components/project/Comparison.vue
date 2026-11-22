<template>
  <div class="comparison-page">
    <div class="container">
      <ProjectBreadcrumb :items="breadcrumbs" />
      <h1 class="page-title">项目对比</h1>
      <p class="page-subtitle">对照投资、周期和申请条件，缩小您的选择范围。</p>

      <div class="comparison-selectors">
        <PublicProjectSelect
          v-model="selectedA"
          label="项目 A"
          :options="options"
          :loading="optionsPending"
          :error="optionsErrorMessage"
          @retry="refreshOptions()"
        />
        <PublicProjectSelect
          v-model="selectedB"
          label="项目 B"
          :options="options"
          :loading="optionsPending"
          :error="optionsErrorMessage"
          @retry="refreshOptions()"
        />
      </div>

      <div v-if="canCompare" class="comparison-actions">
        <button type="button" class="compare-button" @click="setPair(selectedB, selectedA)">交换两个项目</button>
        <button type="button" class="compare-button compare-button--dark" @click="copyLink">复制对比链接</button>
      </div>
      <p v-if="copyMessage" class="copy-message" role="status">{{ copyMessage }}</p>
      <label v-if="copyFailed" class="copy-field">
        对比链接
        <input type="url" :value="shareUrl" readonly @focus="selectLinkText" />
      </label>

      <p v-if="sameProject" class="comparison-state" role="status">请选择两个不同的项目进行对比。</p>
      <p v-else-if="!canCompare" class="comparison-state">请选择两个项目进行对比。</p>
      <div v-else-if="pending" class="comparison-loading" role="status" aria-label="正在加载对比结果">
        <PageSkeleton variant="detail" aria-hidden="true" />
      </div>
      <div v-else-if="error || !comparison" class="comparison-state comparison-state--error">
        <p role="status">暂时无法加载对比结果。请确认两个项目仍可查看，或重试。</p>
        <button type="button" class="compare-button" @click="refreshComparison()">重试对比</button>
      </div>
      <section v-else aria-label="对比结果">
        <div class="comparison-table-wrap">
          <table class="comparison-table">
            <thead>
              <tr>
                <th scope="col">对比项目</th>
                <th v-for="project in comparison.projects" :key="project.slug" scope="col">{{ project.title }}</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="row in comparison.rows" :key="row.label">
                <th scope="row">{{ row.label }}</th>
                <td v-for="cell in row.cells" :key="cell.project.slug">
                  <ul v-if="cell.items?.length" class="comparison-list" role="list">
                    <li v-for="(item, itemIndex) in cell.items" :key="itemIndex">{{ item }}</li>
                  </ul>
                  <span v-else>{{ cell.value || '暂未提供' }}</span>
                </td>
              </tr>
              <tr>
                <th scope="row">项目详情</th>
                <td v-for="project in comparison.projects" :key="project.slug">
                  <NuxtLink :to="projectPath(project.slug)" class="comparison-detail-link" :aria-label="`查看${project.title}详情`">
                    查看项目详情 <span aria-hidden="true">→</span>
                  </NuxtLink>
                </td>
              </tr>
            </tbody>
          </table>
        </div>

        <div class="comparison-groups">
          <section v-for="row in comparison.rows" :key="row.label">
            <h2>{{ row.label }}</h2>
            <dl class="comparison-pair">
              <div v-for="cell in row.cells" :key="cell.project.slug">
                <dt>{{ cell.project.title }}</dt>
                <dd>
                  <ul v-if="cell.items?.length" class="comparison-list" role="list">
                    <li v-for="(item, itemIndex) in cell.items" :key="itemIndex">{{ item }}</li>
                  </ul>
                  <span v-else>{{ cell.value || '暂未提供' }}</span>
                </dd>
              </div>
            </dl>
          </section>
          <section>
            <h2>项目详情</h2>
            <div class="comparison-pair">
              <NuxtLink v-for="project in comparison.projects" :key="project.slug" :to="projectPath(project.slug)" class="comparison-detail-link">
                {{ project.title }} <span aria-hidden="true">→</span>
              </NuxtLink>
            </div>
          </section>
        </div>
      </section>
    </div>

    <div v-if="comparison && !pending && !error" class="compare-consult-bar">
      <div class="container">
        <NuxtLink :to="{ path: '/contact', query: { project: selectedA, compare: selectedB } }" class="compare-button compare-button--gold">
          咨询这两个项目 <span aria-hidden="true">→</span>
        </NuxtLink>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { HOME_BREADCRUMB } from '~/utils/breadcrumb'

interface ProjectOption {
  slug: string
  name: string
}

interface ProjectOptionsResponse {
  data: ProjectOption[] | null
  pagination?: { total: number }
}

interface ComparisonData {
  projects: Array<{ slug: string; title: string }>
  rows: Array<{ label: string; values: string[]; items?: (string[] | null)[] }>
}

const route = useRoute()
const router = useRouter()
const isLegacyLink = computed(() => Boolean(route.params.a && route.params.b))
const readSlug = (value: unknown) => typeof value === 'string' && !value.includes(',') ? value : ''
const selectedA = computed<string>({
  get: () => readSlug(isLegacyLink.value ? route.params.a : route.query.a),
  set: (slug: string) => { void setPair(slug, selectedB.value) },
})
const selectedB = computed<string>({
  get: () => readSlug(isLegacyLink.value ? route.params.b : route.query.b),
  set: (slug: string) => { void setPair(selectedA.value, slug) },
})
const sameProject = computed(() => Boolean(selectedA.value && selectedA.value === selectedB.value))
const canCompare = computed(() => Boolean(selectedA.value && selectedB.value && !sameProject.value))
const requestSlugs = computed(() => canCompare.value ? [selectedA.value, selectedB.value].sort() : [])
const compareKey = computed(() => `public:compare:${requestSlugs.value.join(':')}`)
const optionsKey = 'public:project-options:comparison'
const initialCompareKey = canCompare.value && (import.meta.server || useNuxtApp().isHydrating)
  ? compareKey.value
  : ''

usePublicDataFreshness(() => [
  { versionKey: 'public:projects:list', dataKey: optionsKey },
  // 手动选择会直接请求新结果，只对首屏 SSR 数据做版本复查。
  ...(compareKey.value === initialCompareKey ? [initialCompareKey] : []),
])

const { data: projectOptions, pending: optionsPending, error: optionsError, refresh: refreshOptions } =
  useAsyncData(optionsKey, async (_app, { signal }) => {
    const projects: ProjectOption[] = []
    const perPage = 500
    for (let page = 1; ; page++) {
      const { data, pagination } = await $fetch<ProjectOptionsResponse>('/api/v1/projects/options', {
        query: { page, per_page: perPage },
        signal,
      })
      if (!data?.length) return projects
      projects.push(...data)
      if (!pagination || page * perPage >= pagination.total) return projects
    }
  })

const { data: result, pending, error, refresh: refreshComparison } =
  useAsyncData<{ data: ComparisonData | null }>(compareKey, (_app, { signal }) => {
    if (!requestSlugs.value.length) return Promise.resolve({ data: null })
    return $fetch('/api/v1/projects/compare', {
      query: { slugs: requestSlugs.value.join(',') },
      signal,
      retry: 0,
    })
  }, { dedupe: 'defer' })

const options = computed(() => {
  const choices = (projectOptions.value ?? []).map(({ slug, name }) => ({ value: slug, label: name }))
  if (!choices.length) return choices
  return [{ value: '', label: '请选择项目' }, ...choices]
})
const optionsErrorMessage = computed(() => optionsError.value ? '项目选项加载失败。' : undefined)
const priorityRows = ['投资金额', '办理周期', '适合人群', '申请条件']
const comparison = computed(() => {
  const data = result.value?.data
  if (!canCompare.value || !data) return null
  const indexes = [selectedA.value, selectedB.value].map(slug => data.projects.findIndex(project => project.slug === slug))
  if (indexes.some(index => index < 0)) return null

  return {
    projects: indexes.map(index => data.projects[index]!),
    rows: [
      ...data.rows.filter(row => priorityRows.includes(row.label)).sort((a, b) => priorityRows.indexOf(a.label) - priorityRows.indexOf(b.label)),
      ...data.rows.filter(row => !priorityRows.includes(row.label)),
    ].map(row => ({
      label: row.label,
      cells: indexes.map(index => ({
        project: data.projects[index]!,
        value: row.values[index],
        items: row.items?.[index],
      })),
    })),
  }
})

const breadcrumbs = computed(() => {
  const base = [HOME_BREADCRUMB, { label: '项目对比', href: '/compare' }]
  const projects = comparison.value?.projects
  if (!isLegacyLink.value || !projects) return base
  return [...base, { label: projects.map(project => project.title).join(' 与 '), href: route.path }]
})
useSeo({ title: '项目对比', breadcrumbs })

const projectPath = (slug: string) => `/projects/${encodeURIComponent(slug)}`
const copyMessage = ref('')
const copyFailed = ref(false)
const shareUrl = ref('')

function setPair(a: string, b: string) {
  return router.replace({ path: '/compare', query: { ...route.query, a: a || undefined, b: b || undefined } })
}

watch([selectedA, selectedB], () => {
  copyMessage.value = ''
  copyFailed.value = false
})

async function copyLink() {
  const href = router.resolve({ path: '/compare', query: { a: selectedA.value, b: selectedB.value } }).href
  shareUrl.value = new URL(href, window.location.origin).href
  try {
    await navigator.clipboard.writeText(shareUrl.value)
    copyMessage.value = '已复制对比链接。'
    copyFailed.value = false
  } catch {
    copyMessage.value = '请手动复制对比链接。'
    copyFailed.value = true
  }
}

function selectLinkText(event: FocusEvent) {
  if (event.target instanceof HTMLInputElement) event.target.select()
}
</script>

<style scoped>
.comparison-page {
  padding-bottom: 28px;
}

.page-title {
  font-family: var(--font-serif);
  font-size: clamp(28px, 4vw, 42px);
  line-height: 1.25;
  margin-bottom: 8px;
}

.page-subtitle {
  color: var(--color-text-secondary);
  margin-bottom: 24px;
}

.comparison-selectors {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 16px;
}

.comparison-actions {
  display: flex;
  justify-content: flex-end;
  flex-wrap: wrap;
  gap: 8px;
  margin-top: 12px;
}

.compare-button {
  display: inline-flex;
  justify-content: center;
  align-items: center;
  gap: 8px;
  min-height: 44px;
  padding: 8px 14px;
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  color: var(--color-text);
  background: var(--bg-white);
  font: 600 14px/1.4 var(--font-sans);
  cursor: pointer;
}

.compare-button:hover {
  border-color: var(--color-primary);
}

.compare-button--dark {
  color: var(--bg-white);
  background: var(--color-primary);
  border-color: var(--color-primary);
}

.compare-button--dark:hover {
  background: var(--color-primary-hover);
}

.compare-button--gold {
  background: var(--gradient-gold);
  border-color: transparent;
}

.compare-button--gold:hover {
  background: var(--gradient-gold-hover);
}

.copy-message {
  margin-top: 6px;
  color: var(--color-text-secondary);
  font-size: 14px;
  text-align: right;
}

.copy-field {
  display: block;
  margin-top: 8px;
  font-size: 14px;
}

.copy-field input {
  display: block;
  width: 100%;
  min-height: 44px;
  padding: 8px 12px;
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  color: var(--color-text);
  font: inherit;
}

.comparison-state,
.comparison-loading,
.comparison-table-wrap {
  margin-top: 16px;
}

.comparison-state {
  padding: 28px 16px;
  border: 1px solid var(--color-border);
  border-radius: var(--radius-lg);
  background: var(--color-bg-light);
  color: var(--color-text-secondary);
  text-align: center;
}

.comparison-state--error .compare-button {
  margin-top: 12px;
}

.comparison-table-wrap {
  border: 1px solid var(--color-border);
  border-radius: var(--radius-lg);
  overflow: hidden;
}

.comparison-table {
  width: 100%;
  table-layout: fixed;
  border-collapse: collapse;
  text-align: left;
  font-size: 15px;
}

.comparison-table thead {
  background: #edf2f7;
}

.comparison-table thead th:first-child {
  width: 18%;
}

.comparison-table :is(th, td) {
  padding: 13px 16px;
  border-right: 1px solid var(--color-border);
  border-bottom: 1px solid var(--color-border);
  vertical-align: top;
  overflow-wrap: anywhere;
}

.comparison-table tr > :last-child {
  border-right: 0;
}

.comparison-table tbody tr:last-child > * {
  border-bottom: 0;
}

.comparison-table tbody tr:nth-child(even) {
  background: var(--color-bg-light);
}

.comparison-list li + li {
  margin-top: 4px;
}

.comparison-detail-link {
  color: var(--color-primary-light);
  font-weight: 700;
  text-decoration: underline;
  text-underline-offset: 3px;
}

.comparison-groups {
  display: none;
}

.compare-consult-bar {
  margin-top: 20px;
  text-align: right;
}

@media (max-width: 767px) {
  .comparison-page {
    padding-bottom: 0;
  }

  .comparison-selectors {
    grid-template-columns: minmax(0, 1fr);
    gap: 12px;
  }

  .comparison-table-wrap {
    display: none;
  }

  .comparison-groups {
    display: grid;
    gap: 16px;
    margin-top: 18px;
  }

  .comparison-groups h2 {
    margin-bottom: 6px;
    font-size: 17px;
  }

  .comparison-pair {
    display: grid;
    grid-template-columns: repeat(2, minmax(0, 1fr));
    border: 1px solid var(--color-border);
    border-radius: var(--radius-md);
  }

  .comparison-pair > * {
    padding: 10px;
    overflow-wrap: anywhere;
    font-size: 14px;
  }

  .comparison-pair > :nth-child(2) {
    border-inline-start: 1px solid var(--color-border);
  }

  .comparison-pair > div {
    display: grid;
    grid-row: span 2;
    grid-template-rows: subgrid;
  }

  .comparison-pair dt {
    margin-bottom: 5px;
    color: var(--color-text-secondary);
    font-size: 13px;
  }

  .compare-consult-bar {
    position: sticky;
    z-index: 20;
    bottom: 0;
    min-height: calc(66px + env(safe-area-inset-bottom, 0px));
    margin-top: 16px;
    border-top: 1px solid var(--color-border);
    padding: 10px 0 calc(10px + env(safe-area-inset-bottom, 0px));
    background: var(--bg-white);
  }

  .compare-consult-bar .compare-button {
    width: 100%;
  }
}
</style>
