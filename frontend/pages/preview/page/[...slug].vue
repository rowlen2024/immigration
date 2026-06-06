<template>
  <div class="preview-page">
    <div v-if="pending" class="preview-loading">加载中...</div>
    <div v-else-if="error" class="preview-error">
      <template v-if="error === 'NOT_FOUND'">
        <h1 class="not-found-title">404</h1>
        <p class="not-found-desc">页面未找到</p>
        <NuxtLink to="/" class="btn-primary">返回首页</NuxtLink>
      </template>
      <template v-else>
        <p>{{ error }}</p>
        <el-button @click="load">重试</el-button>
      </template>
    </div>
    <template v-else-if="page">
      <ProjectBreadcrumb :label="page.title" />
      <h1 class="page-title">{{ page.title }}</h1>
      <div class="page-content" v-html="page.content"></div>
    </template>
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
  status: string;
  template?: string;
}

const page = ref<CmsPage | null>(null);
const pending = ref(true);
const error = ref<string | null>(null);

const load = async () => {
  pending.value = true;
  error.value = null;
  try {
    const api = useApi();
    const data = await api<CmsPage>(`/admin/pages/preview?slug=${encodeURIComponent(slug.value)}`);
    if (!data) {
      error.value = 'NOT_FOUND';
      return;
    }
    page.value = data as CmsPage;
  } catch (e: any) {
    if (e?.statusCode === 404 || e?.response?.status === 404) {
      error.value = 'NOT_FOUND';
    } else {
      error.value = e?.data?.message || e?.message || '加载失败';
    }
  } finally {
    pending.value = false;
  }
};

useSeo({
  title: page.value?.title || '预览页面',
  description: '',
  robots: 'noindex, nofollow',
});

onMounted(load);
</script>

<style scoped>
.preview-page {
  max-width: var(--max-width);
  margin: 0 auto;
  padding: 40px 24px;
}

.preview-loading {
  text-align: center;
  padding: 80px 20px;
  color: var(--color-text-muted);
}

.preview-error {
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
  color: var(--color-text-secondary);
  margin-bottom: 32px;
}

.btn-primary {
  display: inline-block;
  padding: 12px 32px;
  background: var(--accent);
  color: #fff;
  border-radius: 8px;
  font-size: 16px;
  font-weight: 600;
  text-decoration: none;
}

.page-title {
  font-size: 36px;
  font-weight: 700;
  color: var(--color-text);
  margin-bottom: 32px;
}

.page-content {
  font-size: 16px;
  color: var(--color-text-secondary);
  line-height: 1.9;
  margin-bottom: 48px;
}

.page-content :deep(h2) { font-size: 28px; font-weight: 700; margin: 40px 0 20px; color: var(--color-text); }
.page-content :deep(h3) { font-size: 22px; font-weight: 600; margin: 32px 0 16px; color: var(--color-text); }
.page-content :deep(p) { margin-bottom: 16px; }
.page-content :deep(ul), .page-content :deep(ol) { margin-bottom: 16px; padding-left: 24px; }
.page-content :deep(li) { margin-bottom: 8px; }
.page-content :deep(img) { max-width: 100%; border-radius: 4px; margin: 20px 0; }
.page-content :deep(a) { color: var(--accent); text-decoration: underline; }
.page-content :deep(table) { width: 100%; border-collapse: collapse; margin: 20px 0; }
.page-content :deep(th), .page-content :deep(td) { padding: 12px 16px; border: 1px solid var(--border-color); text-align: left; }
.page-content :deep(th) { background: var(--bg-light); font-weight: 600; }
.page-content :deep(blockquote) { border-left: 4px solid var(--accent); padding: 12px 20px; margin: 20px 0; background: var(--bg-light); border-radius: 0 4px 4px 0; }

@media (max-width: 767px) {
  .page-title { font-size: 28px; }
  .not-found-title { font-size: 80px; }
  .preview-page { padding: 24px 16px; }
}
</style>
