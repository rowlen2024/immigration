<template>
  <div class="case-detail-page">
    <div class="container">
      <ProjectBreadcrumb :label="item?.name" parentLabel="成功案例" parentLink="/cases" />

      <div v-if="pending" class="page-skeleton-wrapper"><PageSkeleton variant="detail" /></div>
      <div v-else-if="error" class="error-state">{{ error }}</div>
      <div v-else-if="item" class="case-detail">
        <h1 class="case-title">{{ item.name }}</h1>

        <div class="case-meta">
          <div class="meta-item" v-if="item.country_from">
            <span class="meta-label">来源国家</span>
            <span class="meta-value">{{ item.country_from }}</span>
          </div>
          <div class="meta-item" v-if="item.project?.name">
            <span class="meta-label">所属项目</span>
            <NuxtLink :to="'/projects/' + item.project_id" class="meta-value link">
              {{ item.project.name }}
            </NuxtLink>
          </div>
          <div class="meta-item" v-if="item.investment_amount">
            <span class="meta-label">投资金额</span>
            <span class="meta-value">{{ item.investment_amount }}</span>
          </div>
          <div class="meta-item" v-if="item.processing_period">
            <span class="meta-label">办理周期</span>
            <span class="meta-value">{{ item.processing_period }}</span>
          </div>
        </div>

        <div class="case-content" v-html="item.content" />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
const route = useRoute();
const slug = route.params.slug as string;

const { data, pending, error, refresh } = await useFetch<{ data: any }>(`/api/v1/cases/${slug}`);

const item = computed(() => data.value?.data ?? null);

useSeo({
  title: item.value?.name ?? '案例详情',
});

onMounted(() => {
  $fetch(`/api/v1/cases/${slug}`).then(v => { data.value = v }).catch(() => {})
})
</script>

<style scoped>
.case-detail-page {
  padding: 60px 0;
}

.case-title {
  font-size: 36px;
  font-weight: 700;
  color: var(--text-primary);
  margin-bottom: 24px;
}

.case-meta {
  display: flex;
  flex-wrap: wrap;
  gap: 32px;
  margin-bottom: 40px;
  padding: 24px;
  background: var(--bg-white);
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-sm);
}

.meta-item {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.meta-label {
  font-size: 13px;
  color: var(--text-light);
}

.meta-value {
  font-size: 16px;
  font-weight: 600;
  color: var(--text-primary);
}

.meta-value.link {
  color: var(--primary);
  text-decoration: none;
}

.meta-value.link:hover {
  text-decoration: underline;
}

.case-content {
  font-size: 16px;
  line-height: 1.85;
  color: var(--text-primary);
  background: var(--bg-white);
  border-radius: var(--radius-lg);
  padding: 32px;
  box-shadow: var(--shadow-sm);
}

.case-content :deep(img) {
  max-width: 100%;
  height: auto;
  border-radius: 4px;
}

.case-content :deep(table) {
  border-collapse: collapse;
  width: 100%;
  margin: 16px 0;
}

.case-content :deep(th),
.case-content :deep(td) {
  border: 1px solid var(--el-border-color);
  padding: 8px 12px;
  text-align: left;
}

.case-content :deep(th) {
  background: var(--el-fill-color-light);
  font-weight: 600;
}

.case-content :deep(blockquote) {
  border-left: 3px solid var(--el-border-color-dark);
  padding-left: 16px;
  color: var(--el-text-color-secondary);
  margin: 16px 0;
}

.case-content :deep(pre) {
  background: #1e1e1e;
  color: #d4d4d4;
  border-radius: 6px;
  padding: 16px;
  overflow-x: auto;
}

.case-content :deep(pre code) {
  background: none;
  color: inherit;
  font-size: 13px;
}

.loading-state,
.error-state {
  text-align: center;
  padding: 80px 20px;
  color: var(--text-light);
  font-size: 16px;
}

.error-state {
  color: #c62828;
}

@media (max-width: 767px) {
  .case-title {
    font-size: 28px;
  }

  .case-meta {
    gap: 16px;
    padding: 16px;
  }

  .case-content {
    padding: 20px;
  }
}
</style>
