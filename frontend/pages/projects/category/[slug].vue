<template>
  <div class="category-page">
    <div v-if="pending && !category" class="container category-content">
      <PageSkeleton variant="cards" :count="6" />
    </div>

    <template v-else-if="category">
      <ProjectBrowseHero
        :breadcrumbs="breadcrumbs"
        :title="category.name"
        :description="categoryIntro"
        :images="heroImages"
      >
        <div class="category-actions">
          <a v-if="projects.length" href="#category-projects" class="category-jump">
            查看本类项目
            <span aria-hidden="true">↓</span>
          </a>
          <NuxtLink to="/contact" class="category-hero-contact">咨询顾问</NuxtLink>
        </div>
      </ProjectBrowseHero>

      <div class="container category-content">
        <section id="category-projects" class="category-results" aria-labelledby="category-results-title" tabindex="-1">
          <ProjectCategoryDirectory label="浏览其他项目类型" :current-slug="slug">
            <template #leading>
              <div class="category-list-tools">
                <div class="category-list-title">
                  <h2 id="category-results-title">本类项目</h2>
                  <span role="status">{{ filteredProjects.length }} 个项目</span>
                </div>
                <label v-if="showRegionFilter" class="category-region">
                  <span>国家或地区</span>
                  <select v-model="selectedRegion" class="public-select" name="category-region">
                    <option value="">全部地区</option>
                    <option v-for="region in regions" :key="region" :value="region">{{ region }}</option>
                  </select>
                </label>
              </div>
            </template>
          </ProjectCategoryDirectory>

          <template v-if="visibleProjects.length">
            <div ref="projectGrid" class="category-projects">
              <ProjectCard
                v-for="(project, index) in visibleProjects"
                :key="project.slug"
                v-bind="project"
                :image-variant="index"
              />
            </div>
            <div v-if="visibleProjects.length < filteredProjects.length" class="category-more">
              <button type="button" class="category-directory-link" @click="showMore">查看更多项目</button>
            </div>
          </template>
          <div v-else class="category-state">
            <p>该类目暂无已发布项目。</p>
            <NuxtLink to="/projects" class="category-directory-link">查看全部移民项目</NuxtLink>
          </div>
        </section>

        <section v-if="guideQuestions.length && projects.length" class="category-guide" aria-labelledby="category-guide-title">
          <h2 id="category-guide-title">选择{{ category.name }}前先了解</h2>
          <div>
            <details v-for="question in guideQuestions" :key="question.title">
              <summary>{{ question.title }}</summary>
              <p>{{ question.answer }}</p>
            </details>
          </div>
        </section>

        <aside class="category-consultation">
          <div>
            <h2>想缩小选择范围？</h2>
            <p>和顾问聊聊你的目的地与家庭规划。</p>
          </div>
          <NuxtLink to="/contact" class="category-consultation-link">咨询顾问 <span aria-hidden="true">→</span></NuxtLink>
        </aside>
      </div>
    </template>

    <div v-else class="container category-content">
      <div class="category-state" role="status">
        <p>类目信息加载失败。</p>
        <button type="button" class="category-directory-link" @click="refresh()">重试</button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { HOME_BREADCRUMB } from '~/utils/breadcrumb'
import type { ImageVariantInfo } from '~/utils/image'
import { projectRegionLabel } from '~/utils/project-browse'
import { toProjectCard, type ProjectCardSource } from '~/utils/project-card'

interface CategoryProject extends ProjectCardSource {
  country?: string
}

interface ProjectCategory {
  name: string
  slug: string
  intro: string
  seo_title: string
  seo_description: string
  projects: CategoryProject[]
}

interface CategoryGuide {
  intro: string
  questions?: { title: string; answer: string }[]
}

const guideByCategory = new Map<string, CategoryGuide>([
  ['investment', {
    intro: '通过投资规划海外身份。先按目的地找到相关方案，再查看各项目的投资方式、申请条件与办理流程。',
    questions: [
      {
        title: '选投资项目时，先看哪些信息？',
        answer: '可以先确认目标地区、可投入预算和家庭成员安排，再逐项查看投资方式、申请条件和办理流程。',
      },
      {
        title: '同一地区有多个项目，怎么选？',
        answer: '先按地区缩小范围，再进入项目详情，对照各项目的申请条件、投资方式和费用构成。',
      },
    ],
  }],
  ['property', {
    intro: '有海外置业计划，可以先了解购房与身份申请的关系，再对照各项目的房产及申请条件。',
    questions: [
      {
        title: '看购房项目时，先确认什么？',
        answer: '先确认计划置业的地区，再到对应项目详情核对房产要求、申请条件和持有要求。',
      },
    ],
  }],
  ['entrepreneur', {
    intro: '从计划经营的国家或地区出发，了解各方案的申请条件和经营要求。',
  }],
  ['citizenship-by-donation', {
    intro: '浏览捐款入籍项目，逐项查看申请条件、办理流程与家庭成员相关要求。',
  }],
  ['talent-and-study', {
    intro: '从职业背景、专业经历或升学计划出发，了解各地的人才及留学方案。',
  }],
  ['long-stay-visas', {
    intro: '从工作或生活的目的地出发，查看各地的长期签证与居留选择。',
  }],
  ['cross-border-services', {
    intro: '浏览面向跨境生活和经营的服务，查看各项服务的办理条件和所需材料。',
  }],
  ['hot', {
    intro: '浏览北极星移民精选项目，了解不同目的地的身份规划选择。',
  }],
])

const PAGE_SIZE = 6
const route = useRoute()
const slug = computed(() => String(route.params.slug))
const categoryKey = computed(() => `public:project-categories:detail:${slug.value}`)
const browseState = useState<Record<string, { region: string; limit: number }>>('project-category-browse', () => ({}))
usePublicDataFreshness(() => [categoryKey.value])

const { data, error, pending, refresh } = await useFetch<{ data: ProjectCategory }>(
  () => `/api/v1/project-categories/${encodeURIComponent(slug.value)}`,
  { key: categoryKey },
)

function categoryNotFound() {
  return createError({ statusCode: 404, statusMessage: 'Project category not found', message: '项目类目不存在', fatal: true })
}

if (error.value?.statusCode === 404) throw categoryNotFound()
if (import.meta.server && error.value) {
  const event = useRequestEvent()
  if (event) setResponseStatus(event, error.value.statusCode ?? 502)
}

watch(error, (issue) => {
  if (issue?.statusCode === 404) showError(categoryNotFound())
})

const category = computed(() => {
  const result = data.value?.data
  return result?.slug === slug.value ? result : undefined
})
const guide = computed(() => guideByCategory.get(slug.value))
const guideQuestions = computed(() => guide.value?.questions ?? [])
const categoryIntro = computed(() => category.value?.intro.trim() || guide.value?.intro
  || '浏览本类方案，了解各项目的特点、申请条件与办理流程。')
const projects = computed(() => category.value?.projects ?? [])
const regions = computed(() => Array.from(new Set(projects.value.map(project => projectRegionLabel(project.country))))
  .filter(Boolean)
  .sort((first, second) => first.localeCompare(second, 'zh-CN')))
const showRegionFilter = computed(() => projects.value.length >= 4 && regions.value.length > 1)
const selectedRegion = computed({
  get() {
    const region = browseState.value[slug.value]?.region ?? ''
    return showRegionFilter.value && regions.value.includes(region) ? region : ''
  },
  set(region: string) {
    browseState.value[slug.value] = { region, limit: PAGE_SIZE }
  },
})
const filteredProjects = computed(() => selectedRegion.value
  ? projects.value.filter(project => projectRegionLabel(project.country) === selectedRegion.value)
  : projects.value)
const visibleCount = computed(() => browseState.value[slug.value]?.limit ?? PAGE_SIZE)
const visibleProjects = computed(() => filteredProjects.value.slice(0, visibleCount.value).map(toProjectCard))
const projectGrid = ref<HTMLElement | null>(null)

async function showMore() {
  const firstNewIndex = visibleProjects.value.length
  browseState.value[slug.value] = { region: selectedRegion.value, limit: visibleCount.value + PAGE_SIZE }
  await nextTick()
  projectGrid.value?.querySelectorAll<HTMLAnchorElement>(':scope > a')[firstNewIndex]?.focus()
}

const heroImages = computed(() => {
  const images = new Map<string, { src: string; variants?: Record<string, ImageVariantInfo> }>()
  const usedSources = new Set<string>()
  for (const project of projects.value) {
    const src = project.cover_image
    if (!src || usedSources.has(src)) continue
    const region = projectRegionLabel(project.country) || src
    if (images.has(region)) continue
    images.set(region, { src, variants: project.cover_image_variants })
    usedSources.add(src)
    if (images.size === 3) break
  }
  return [...images.values()]
})

const canonicalPath = computed(() => category.value ? `/projects/category/${category.value.slug}` : route.path)

const breadcrumbs = computed(() => category.value ? [
  HOME_BREADCRUMB,
  { label: '移民项目', href: '/projects' },
  { label: category.value.name, href: canonicalPath.value },
] : [])

useSeo({
  title: () => category.value?.seo_title || category.value?.name,
  description: () => category.value?.seo_description || categoryIntro.value,
  robots: () => projects.value.length ? 'index, follow' : 'noindex, follow',
  canonicalPath,
  breadcrumbs,
})
</script>

<style scoped>
.category-page {
  padding-bottom: clamp(48px, 7vw, 88px);
}

.category-content {
  padding-top: clamp(24px, 4vw, 40px);
}

.category-actions,
.category-list-tools,
.category-list-title,
.category-consultation {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
}

.category-actions {
  gap: 12px 20px;
}

.category-jump,
.category-hero-contact,
.category-consultation-link,
.category-directory-link {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-height: 44px;
  gap: 8px;
  font-size: var(--text-sm);
  font-weight: 600;
}

.category-jump,
.category-consultation-link {
  border-radius: var(--radius-sm);
  padding: 0 16px;
}

.category-jump {
  background: var(--gradient-gold);
  color: var(--color-primary-dark);
}

.category-jump:hover {
  background: var(--gradient-gold-hover);
}

.category-hero-contact {
  color: #fff;
}

.category-hero-contact:hover {
  text-decoration: underline;
  text-underline-offset: 3px;
}

.category-results {
  scroll-margin-top: calc(var(--header-height) + 16px);
}

.category-list-tools {
  justify-content: space-between;
  gap: 12px 20px;
}

.category-list-title {
  align-items: baseline;
  gap: 4px 12px;
}

.category-list-title h2 {
  font-size: clamp(22px, 2.8vw, 30px);
  line-height: 1.35;
}

.category-list-title span {
  color: var(--color-text-muted);
  font-size: var(--text-sm);
}

.category-region {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 8px;
  max-width: 100%;
  font-size: var(--text-sm);
}

.category-region select {
  max-width: 100%;
}

.category-projects {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: clamp(16px, 2vw, 24px);
  margin-top: 16px;
}

.category-more,
.category-state {
  margin-top: 24px;
  text-align: center;
}

.category-state {
  padding: 36px 24px;
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  color: var(--color-text-muted);
}

.category-state .category-directory-link {
  margin-top: 8px;
}

.category-directory-link {
  padding: 0 16px;
  border: 1px solid var(--color-border);
  border-radius: var(--radius-sm);
  background: var(--color-bg-surface);
  color: var(--color-text);
  font-family: inherit;
  cursor: pointer;
}

.category-directory-link:hover {
  border-color: var(--color-accent);
  background: var(--color-accent-soft);
}

.category-guide {
  display: grid;
  grid-template-columns: minmax(0, 2fr) minmax(0, 3fr);
  gap: 12px 32px;
  margin-top: clamp(32px, 5vw, 56px);
  padding-top: 24px;
  border-top: 1px solid var(--color-border);
}

.category-guide h2,
.category-consultation h2 {
  font-size: var(--text-lg);
}

.category-guide details {
  border-bottom: 1px solid var(--color-border);
}

.category-guide summary {
  min-height: 44px;
  padding: 10px 0;
  font-weight: 600;
  cursor: pointer;
}

.category-guide details p {
  padding-bottom: 12px;
  color: var(--color-text-secondary);
}

.category-consultation {
  justify-content: space-between;
  gap: 12px 24px;
  margin-top: 32px;
  padding: clamp(16px, 3vw, 28px);
  border-radius: var(--radius-md);
  background: var(--color-accent-soft);
}

.category-consultation p {
  color: var(--color-text-secondary);
}

.category-consultation-link {
  background: var(--color-primary);
  color: #fff;
}

.category-consultation-link:hover {
  background: var(--color-primary-hover);
}

@media (max-width: 1000px) {
  .category-projects {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
}

@media (max-width: 600px) {
  .category-projects,
  .category-guide {
    grid-template-columns: minmax(0, 1fr);
  }
}
</style>
