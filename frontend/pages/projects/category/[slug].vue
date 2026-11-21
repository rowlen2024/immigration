<template>
  <div class="container category-page">
    <ProjectBreadcrumb :items="breadcrumbs" />

    <PageSkeleton v-if="pending && !category" variant="cards" :count="6" />

    <template v-else-if="category">
      <header class="category-heading">
        <h1>{{ category.name }}</h1>
        <p v-if="category.intro">{{ category.intro }}</p>
      </header>

      <div v-if="projects.length" class="category-projects">
        <ProjectCard
          v-for="(project, index) in projects"
          :key="project.slug"
          v-bind="project"
          :image-variant="index"
        />
      </div>
      <div v-else class="category-state">
        <p>该类目暂无已发布项目。</p>
        <NuxtLink to="/projects" class="category-directory-link">查看全部移民项目</NuxtLink>
      </div>
    </template>

    <div v-else class="category-state" role="status">
      <p>类目信息加载失败。</p>
      <button type="button" class="category-directory-link" @click="refresh()">重试</button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { HOME_BREADCRUMB } from '~/utils/breadcrumb'
import { toProjectCard, type ProjectCardSource } from '~/utils/project-card'

interface ProjectCategory {
  name: string
  slug: string
  intro: string
  seo_title: string
  seo_description: string
  projects: ProjectCardSource[]
}

const route = useRoute()
const slug = computed(() => String(route.params.slug))
const categoryKey = computed(() => `public:project-categories:detail:${slug.value}`)
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

const category = computed(() => data.value?.data)
const projects = computed(() => category.value?.projects.map(toProjectCard) ?? [])
const canonicalPath = computed(() => category.value ? `/projects/category/${category.value.slug}` : route.path)

const breadcrumbs = computed(() => category.value ? [
  HOME_BREADCRUMB,
  { label: '移民项目', href: '/projects' },
  { label: category.value.name, href: canonicalPath.value },
] : [])

useSeo({
  title: () => category.value?.seo_title || category.value?.name,
  description: () => category.value?.seo_description || category.value?.intro,
  robots: () => projects.value.length ? 'index, follow' : 'noindex, follow',
  canonicalPath,
  breadcrumbs,
})
</script>

<style scoped>
.category-page {
  padding-bottom: clamp(48px, 7vw, 88px);
}

.category-heading {
  margin: clamp(12px, 3vw, 28px) 0 32px;
}

.category-heading h1 {
  font-size: var(--text-3xl);
  color: var(--color-text);
  overflow-wrap: anywhere;
}

.category-heading p {
  max-width: 68ch;
  margin-top: 12px;
  color: var(--color-text-muted);
  white-space: pre-line;
  overflow-wrap: anywhere;
}

.category-projects {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 24px;
}

.category-state {
  margin: 24px 0;
  padding: 36px 24px;
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  color: var(--color-text-muted);
  text-align: center;
}

.category-directory-link {
  display: inline-flex;
  align-items: center;
  min-height: 44px;
  margin-top: 8px;
  padding: 0 16px;
  border: 1px solid var(--color-border);
  border-radius: var(--radius-sm);
  background: var(--color-bg-surface);
  color: var(--color-text);
  font: inherit;
  cursor: pointer;
}

@media (max-width: 1000px) {
  .category-projects {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
}

@media (max-width: 600px) {
  .category-projects {
    grid-template-columns: minmax(0, 1fr);
  }
}
</style>
