<template>
  <div class="cases-page">
    <div class="container">
      <ProjectBreadcrumb />

      <h1 class="page-title">成功案例</h1>
      <p class="page-subtitle">每一位客户的成功获批，都是我们最大的骄傲</p>

      <div v-if="pending" class="loading-state">加载中...</div>
      <div v-else-if="error" class="error-state">{{ error }}</div>
      <div v-else class="cases-grid">
        <CaseCard
          v-for="item in cases"
          :key="item.id"
          :slug="item.slug"
          :name="item.name"
          :country="item.country"
          :project="item.project || undefined"
          :summary="item.summary"
          :image="item.image"
          :show-result="true"
          result-text="成功获批"
        />
      </div>

      <div v-if="!pending && cases.length === 0" class="empty-state">
        暂无成功案例展示
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
useSeo({ title: '成功案例' });

interface ApiCaseItem {
  id: number;
  slug: string;
  name: string;
  country_from: string;
  photo_url: string;
  content: string;
  project?: { name: string };
}

function stripHtml(html: string): string {
  if (!html) return '';
  return html.replace(/<[^>]+>/g, '').replace(/&nbsp;/g, ' ').slice(0, 80);
}

const { data, pending, error } = await useFetch<{ data: ApiCaseItem[] }>('/api/v1/cases');

const cases = computed(() => {
  const apiData = data.value as { data?: ApiCaseItem[] } | null;
  const items = apiData?.data ?? [];
  return items.map((c) => ({
    id: String(c.id),
    slug: c.slug,
    name: c.name,
    country: c.country_from,
    project: c.project?.name ?? '',
    summary: stripHtml(c.content),
    image: c.photo_url,
  }));
});
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
  margin-bottom: 48px;
}

.loading-state,
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
