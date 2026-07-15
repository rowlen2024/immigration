<template>
  <div class="cms-page">
    <div class="container">
      <div v-if="pending" class="page-skeleton-wrapper"><PageSkeleton variant="content" /></div>
      <div v-else-if="error" class="error-state">
        <template v-if="error === 'NOT_FOUND'">
          <h1 class="not-found-title">404</h1>
          <p class="not-found-desc">页面未找到</p>
          <NuxtLink to="/" class="btn-primary">返回首页</NuxtLink>
        </template>
        <template v-else>
          <p>{{ error }}</p>
        </template>
      </div>
      <template v-else-if="page">
        <!-- default layout -->
        <template v-if="template === 'default'">
          <template v-if="isNewsPage">
            <ProjectBreadcrumb :label="page.title" />
            <article class="news-article">
              <div class="news-layout">
                <div class="news-main">
                  <header class="news-header">
                    <div v-if="page.projects?.length" class="news-projects" aria-label="所属项目">
                      <NuxtLink
                        v-for="project in page.projects"
                        :key="project.id"
                        :to="`/projects/${project.slug}`"
                        class="news-project"
                      >
                        {{ project.name }}
                      </NuxtLink>
                </div>
                <h1 class="news-title">{{ page.title }}</h1>
                <ul v-if="page.tags?.length" class="news-tags" aria-label="文章标签">
                  <li v-for="tag in page.tags" :key="tag">{{ tag }}</li>
                </ul>
                <div class="news-meta" aria-label="文章信息">
                      <time v-if="formattedDate" :datetime="page.created_at">{{ formattedDate }}</time>
                      <span v-if="formattedDate" aria-hidden="true"></span>
                      <span>{{ readingMinutes }} 分钟阅读</span>
                    </div>
                    <ResponsiveImage
                      v-if="page.cover_image"
                      class="news-cover"
                      :src="page.cover_image"
                      variant="lg"
                      :variants="page.cover_image_variants"
                      sizes="(max-width: 767px) calc(100vw - 40px), 860px"
                      loading="eager"
                      fetchpriority="high"
                      :alt="page.title"
                    />
                  </header>
                  <div class="page-content news-content" v-html="page.content"></div>
                </div>
                <aside class="news-sidebar" aria-label="文章相关内容">
                  <section v-if="relatedPages.length" class="related-panel">
                    <div class="sidebar-heading">
                      <span>相关文章</span>
                      <small>RELATED</small>
                    </div>
                    <div class="related-list">
                      <NuxtLink
                        v-for="item in relatedPages"
                        :key="item.id"
                        :to="`/pages/${item.slug}`"
                        class="related-item"
                      >
                        <ResponsiveImage
                          v-if="item.cover_image"
                          class="related-cover"
                          :src="item.cover_image"
                          variant="thumb"
                          :variants="item.cover_image_variants"
                          sizes="84px"
                          :alt="item.title"
                        />
                        <span class="related-copy">
                          <strong>{{ item.title }}</strong>
                          <time
                            v-if="item.created_at && formatArticleDate(item.created_at)"
                            :datetime="item.created_at"
                          >
                            {{ formatArticleDate(item.created_at) }}
                          </time>
                        </span>
                      </NuxtLink>
                    </div>
                  </section>
                  <ConsultCTA
                    variant="sidebar"
                    title="想知道这条路径是否适合您？"
                    description="与顾问沟通您的家庭情况和投资目标，获取更有针对性的建议。"
                  />
                </aside>
              </div>
            </article>
          </template>
          <template v-else>
            <ProjectBreadcrumb :label="page.title" />
            <h1 class="page-title">{{ page.title }}</h1>
            <div class="page-content" v-html="page.content"></div>
          </template>
        </template>

        <!-- fullwidth layout -->
        <template v-else-if="template === 'fullwidth'">
          <ProjectBreadcrumb :label="page.title" />
          <div class="container">
            <h1 class="page-title">{{ page.title }}</h1>
          </div>
          <div class="page-content page-content--fullwidth" v-html="page.content"></div>
        </template>

        <!-- landing layout -->
        <template v-else-if="template === 'landing'">
          <div class="page-content page-content--landing" v-html="page.content"></div>
        </template>
      </template>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { ImageVariantInfo } from '~/utils/image';

const route = useRoute();

const slug = computed(() => {
  const s = route.params.slug;
  if (Array.isArray(s)) return s.join('/');
  return s || '';
});
const pageDataKey = computed(() => `public:page:${slug.value}`);

interface CmsPage {
  title: string;
  content: string;
  template?: string;
  cover_image?: string;
  cover_image_variants?: Record<string, ImageVariantInfo>;
  page_type?: string;
  created_at?: string;
  projects?: Array<{ id: number; name: string; slug: string }>;
  tags?: string[];
  meta_title?: string;
  meta_description?: string;
}

interface RelatedPage {
  id: number;
  title: string;
  slug: string;
  cover_image?: string;
  cover_image_variants?: Record<string, ImageVariantInfo>;
  created_at?: string;
}

interface RelatedPagesEnvelope {
  code: number;
  data: RelatedPage[];
}

usePublicDataFreshness(() => [pageDataKey.value]);

const { data, pending, error: fetchError, refresh } = await useFetch(
  () => `/api/v1/pages/${slug.value}`,
  {
    key: pageDataKey,
    transform: (response) => {
      const envelope = response as { code: number; data: CmsPage };
      return envelope?.data ?? null;
    },
  }
);

if (fetchError.value?.statusCode === 404) {
  throw createError({
    statusCode: 404,
    statusMessage: '页面不存在',
    fatal: true,
  });
}

const page = computed(() => data.value || null);

const template = computed(() => page.value?.template || 'default');
const isNewsPage = computed(() => page.value?.page_type === 'news');

const relatedDataKey = computed(() => `public:page:${slug.value}:related`);
const { data: relatedResponse } = await useFetch<RelatedPagesEnvelope>(
  () => `/api/v1/related-pages?slug=${encodeURIComponent(String(slug.value))}`,
  {
    key: relatedDataKey,
    immediate: isNewsPage.value && template.value === 'default',
    watch: false,
  },
);

const relatedPages = computed(() => relatedResponse.value?.data?.slice(0, 4) ?? []);

let relatedRequestToken = 0;
watch(
  [slug, isNewsPage, template],
  async ([currentSlug, currentIsNews, currentTemplate]) => {
    if (!import.meta.client) return;
    const requestToken = ++relatedRequestToken;
    relatedResponse.value = null;
    if (!currentIsNews || currentTemplate !== 'default') return;
    try {
      const response = await $fetch<RelatedPagesEnvelope>(
        `/api/v1/related-pages?slug=${encodeURIComponent(String(currentSlug))}`,
      );
      if (
        requestToken !== relatedRequestToken
        || slug.value !== currentSlug
        || !isNewsPage.value
        || template.value !== currentTemplate
      ) return;
      relatedResponse.value = response;
    } catch {
      // Related content is supplementary and must not block the article.
    }
  },
  { immediate: true },
);

const readingMinutes = computed(() => {
  const plainText = (page.value?.content || '')
    .replace(/<[^>]*>/g, ' ')
    .replace(/&[a-zA-Z0-9#]+;/g, ' ')
    .replace(/\s+/g, ' ')
    .trim();
  return Math.max(1, Math.ceil(plainText.length / 400));
});

function formatArticleDate(value?: string) {
  if (!value) return '';
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) return '';
  return new Intl.DateTimeFormat('zh-CN', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
  }).format(date);
}

const formattedDate = computed(() => page.value?.created_at ? formatArticleDate(page.value.created_at) : '');

const error = computed(() => {
  if (!fetchError.value) return null;

  const err = fetchError.value as { statusCode?: number };
  if (err.statusCode === 404) return 'NOT_FOUND';
  return '页面加载失败，请稍后重试';
});

useSeo({
  title: page.value?.meta_title || page.value?.title || '内容页面',
  description: page.value?.meta_description || '',
  breadcrumbLabel: page.value?.title,
});

</script>

<style scoped>
.news-article {
  padding-bottom: 64px;
}

.news-header {
  max-width: 860px;
  margin-bottom: 40px;
}

.news-projects {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-bottom: 18px;
}

.news-project {
  display: inline-flex;
  padding: 5px 12px;
  border-radius: var(--radius-full);
  background: var(--color-accent-soft);
  color: var(--color-accent-dark);
  font-size: 13px;
  font-weight: 600;
  transition: background-color var(--duration-fast) var(--ease-out), color var(--duration-fast) var(--ease-out);
}

.news-project:hover {
  background: var(--color-accent);
  color: #fff;
}

.news-title {
  max-width: 18em;
  margin-bottom: 18px;
  color: var(--text-primary);
  font-family: var(--font-serif);
  font-size: clamp(32px, 4vw, 48px);
  font-weight: 700;
  line-height: 1.3;
  letter-spacing: -0.02em;
}

.news-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-bottom: 18px;
}

.news-tags li {
  padding: 5px 11px;
  border: 1px solid rgba(200, 150, 62, 0.28);
  border-radius: var(--radius-full);
  background: var(--color-accent-soft);
  color: var(--color-primary-light);
  font-size: 12px;
  font-weight: 500;
  line-height: 1.4;
}

.news-meta {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 28px;
  color: var(--text-light);
  font-size: 14px;
}

.news-meta > span[aria-hidden='true'] {
  width: 3px;
  height: 3px;
  border-radius: 50%;
  background: var(--color-accent);
}

.news-cover {
  width: 100%;
  aspect-ratio: 16 / 9;
  object-fit: cover;
  border-radius: var(--radius-lg);
}

.news-layout {
  display: grid;
  grid-template-columns: minmax(0, 3fr) minmax(240px, 1fr);
  gap: 48px;
  align-items: start;
}

.news-main {
  min-width: 0;
}

.news-content {
  max-width: 760px;
  margin-bottom: 0;
  color: var(--color-text-secondary);
  font-size: 17px;
  line-height: 1.95;
}

.news-content :deep(h2) {
  padding-bottom: 12px;
  border-bottom: 1px solid var(--border-color);
  font-family: var(--font-serif);
}

.news-sidebar {
  position: sticky;
  top: calc(var(--header-scrolled-height) + 28px);
  display: grid;
  gap: 24px;
}

.related-panel {
  padding: 24px;
  border: 1px solid var(--border-color);
  border-radius: var(--radius-lg);
  background: var(--bg-white);
}

.sidebar-heading {
  display: flex;
  align-items: baseline;
  justify-content: space-between;
  padding-bottom: 14px;
  border-bottom: 1px solid var(--border-color);
  color: var(--text-primary);
  font-size: 18px;
  font-weight: 700;
}

.sidebar-heading small {
  color: var(--color-accent-dark);
  font-size: 10px;
  letter-spacing: 0.12em;
}

.related-list {
  display: grid;
}

.related-item {
  display: flex;
  gap: 12px;
  padding: 16px 0;
  border-bottom: 1px solid var(--color-border-light);
  transition: color var(--duration-fast) var(--ease-out);
}

.related-item:last-child {
  padding-bottom: 0;
  border-bottom: 0;
}

.related-item:hover {
  color: var(--color-accent-dark);
}

.related-cover {
  flex: 0 0 84px;
  width: 84px;
  height: 60px;
  object-fit: cover;
  border-radius: var(--radius-md);
}

.related-copy {
  display: grid;
  min-width: 0;
  gap: 7px;
}

.related-copy strong {
  display: -webkit-box;
  overflow: hidden;
  color: inherit;
  font-size: 14px;
  font-weight: 600;
  line-height: 1.55;
  -webkit-box-orient: vertical;
  -webkit-line-clamp: 2;
}

.related-copy time {
  color: var(--text-light);
  font-size: 12px;
}

.page-title {
  font-size: 36px;
  font-weight: 700;
  color: var(--text-primary);
  margin-bottom: 32px;
}

.page-content {
  font-size: 16px;
  color: var(--text-secondary);
  line-height: 1.9;
  margin-bottom: 48px;
}

.page-content :deep(h2) {
  font-size: 28px;
  font-weight: 700;
  color: var(--text-primary);
  margin: 40px 0 20px;
}

.page-content :deep(h3) {
  font-size: 22px;
  font-weight: 600;
  color: var(--text-primary);
  margin: 32px 0 16px;
}

.page-content :deep(p) {
  margin-bottom: 16px;
}

.page-content :deep(ul),
.page-content :deep(ol) {
  margin-bottom: 16px;
  padding-left: 24px;
}

.page-content :deep(ul) {
  list-style: disc;
}

.page-content :deep(ol) {
  list-style: decimal;
}

.page-content :deep(li) {
  margin-bottom: 8px;
}

.page-content :deep(img) {
  max-width: 100%;
  height: auto;
  border-radius: var(--radius-md);
  margin: 20px 0;
}

.page-content :deep(a) {
  color: var(--primary);
  text-decoration: underline;
  transition: color 0.2s ease;
}

.page-content :deep(a:hover) {
  color: var(--accent-dark);
}

.page-content :deep(table) {
  display: block;
  overflow-x: auto;
  -webkit-overflow-scrolling: touch;
  width: 100%;
  border-collapse: collapse;
  margin: 20px 0;
}

.page-content :deep(pre) {
  overflow-x: auto;
  -webkit-overflow-scrolling: touch;
}

.page-content :deep(th),
.page-content :deep(td) {
  padding: 12px 16px;
  border: 1px solid var(--border-color);
  text-align: left;
}

.page-content :deep(th) {
  background-color: var(--bg-light);
  font-weight: 600;
}

.page-content :deep(blockquote) {
  border-left: 4px solid var(--accent);
  padding: 12px 20px;
  margin: 20px 0;
  background-color: var(--bg-light);
  border-radius: 0 var(--radius-md) var(--radius-md) 0;
}

.loading-state {
  text-align: center;
  padding: 80px 20px;
  color: var(--text-light);
  font-size: 16px;
}

.error-state {
  text-align: center;
  padding: 80px 20px;
}

.not-found-title {
  font-size: 120px;
  font-weight: 800;
  color: var(--accent);
  line-height: 1;
  margin-bottom: 8px;
}

.not-found-desc {
  font-size: 20px;
  color: var(--text-secondary);
  margin-bottom: 32px;
}

@media (max-width: 890px) {
  .news-article {
    padding-bottom: calc(96px + env(safe-area-inset-bottom, 0px));
  }

  .news-header {
    margin-bottom: 32px;
  }

  .news-title {
    font-size: 30px;
  }

  .news-tags {
    gap: 6px;
  }

  .news-tags li {
    padding: 4px 10px;
  }

  .news-cover {
    border-radius: var(--radius-md);
  }

  .news-layout {
    grid-template-columns: minmax(0, 1fr);
    gap: 36px;
  }

  .news-sidebar {
    position: static;
  }

  .news-content {
    font-size: 16px;
    line-height: 1.85;
  }

  .related-panel {
    padding: 20px;
  }
}

@media (max-width: 767px) {
  .page-title {
    font-size: 28px;
    margin-bottom: 24px;
  }

  .not-found-title {
    font-size: 80px;
  }

  .page-content {
    font-size: 15px;
    line-height: 1.75;
  }

  .page-content :deep(h2) {
    font-size: 22px;
    margin: 32px 0 16px;
  }

  .page-content :deep(h3) {
    font-size: 18px;
    margin: 24px 0 12px;
  }

  .page-content :deep(ul),
  .page-content :deep(ol) {
    padding-left: 20px;
  }

  .page-content :deep(th),
  .page-content :deep(td) {
    padding: 10px 12px;
    font-size: 14px;
  }

  .page-content :deep(blockquote) {
    padding: 10px 16px;
    margin: 16px 0;
  }
}

@media (min-width: 891px) and (max-width: 1023px) {
  .news-layout {
    grid-template-columns: minmax(0, 2fr) minmax(220px, 1fr);
    gap: 32px;
  }
}

.page-content--fullwidth {
  padding: 0 16px;
}

.page-content--landing {
  /* No container — full bleed */
}

.page-content--landing :deep(img) {
  max-width: 100%;
  height: auto;
}

@media (min-width: 768px) {
  .page-content--fullwidth {
    padding: 0 24px;
  }
}
</style>
