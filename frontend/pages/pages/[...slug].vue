<template>
  <div class="cms-page">
    <div class="container">
      <div v-if="pending" class="loading-state">加载中...</div>
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
          <div class="container">
            <ProjectBreadcrumb :label="page.title" />
            <h1 class="page-title">{{ page.title }}</h1>
            <div class="page-content" v-html="page.content"></div>
          </div>
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
const route = useRoute();

const slug = computed(() => {
  const s = route.params.slug;
  if (Array.isArray(s)) return s.join('/');
  return s || '';
});

interface CmsPage {
  title: string;
  content: string;
  template?: string;
  cover_image?: string;
  meta_title?: string;
  meta_description?: string;
}

const { data, pending, error: fetchError } = await useFetch(
  () => `/api/v1/pages/${slug.value}`,
  {
    transform: (response) => {
      const envelope = response as { code: number; data: CmsPage };
      return envelope?.data ?? null;
    },
  }
);

const page = computed(() => data.value || null);

const template = computed(() => page.value?.template || 'default');

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
  width: 100%;
  border-collapse: collapse;
  margin: 20px 0;
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

@media (max-width: 767px) {
  .page-title {
    font-size: 28px;
  }

  .not-found-title {
    font-size: 80px;
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
