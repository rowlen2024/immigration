<template>
  <div class="cases-page">
    <div class="container">
      <ProjectBreadcrumb />

      <h1 class="page-title">成功案例</h1>
      <p class="page-subtitle">每一位客户的成功获批，都是我们最大的骄傲</p>

      <div v-if="pending" class="loading-state">加载中...</div>
      <div v-else-if="error" class="error-state">{{ error }}</div>
      <div v-else class="cases-grid">
        <div v-for="item in cases" :key="item.id" class="case-card">
          <div class="case-image">
            <img
              :src="item.image || ''"
              :alt="item.name"
              loading="lazy"
            />
          </div>
          <div class="case-body">
            <div class="case-meta">
              <span class="case-country">{{ item.country }}</span>
              <span class="case-project">{{ item.project }}</span>
            </div>
            <h3 class="case-name">{{ item.name }}</h3>
            <p class="case-desc">{{ item.description }}</p>
            <div class="case-result">
              <span class="result-badge">{{ item.result }}</span>
            </div>
          </div>
        </div>
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
  name: string;
  country_from: string;
  photo_url: string;
  description: string;
  project?: { name: string };
}

const { data, pending, error } = await useFetch<{ data: ApiCaseItem[] }>('/api/v1/cases');

const cases = computed(() => {
  const apiData = data.value as { data?: ApiCaseItem[] } | null;
  const items = apiData?.data ?? [];
  return items.map((c) => ({
    id: String(c.id),
    name: c.name,
    country: c.country_from,
    project: c.project?.name ?? '',
    description: c.description,
    result: '成功获批',
    image: c.photo_url,
  }));
});
</script>

<style scoped>
.page-title {
  font-size: 36px;
  font-weight: 700;
  color: var(--text-primary);
  margin-bottom: 12px;
}

.page-subtitle {
  font-size: 16px;
  color: var(--text-light);
  margin-bottom: 40px;
}

.cases-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 32px;
  margin-bottom: 48px;
}

.case-card {
  background-color: var(--bg-white);
  border-radius: var(--radius-lg);
  overflow: hidden;
  box-shadow: var(--shadow-sm);
  transition: box-shadow 0.3s ease, transform 0.3s ease;
}

.case-card:hover {
  box-shadow: var(--shadow-md);
  transform: translateY(-4px);
}

.case-image {
  height: 200px;
  overflow: hidden;
}

.case-image img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: transform 0.3s ease;
}

.case-card:hover .case-image img {
  transform: scale(1.05);
}

.case-body {
  padding: 24px;
}

.case-meta {
  display: flex;
  gap: 8px;
  margin-bottom: 12px;
}

.case-country,
.case-project {
  font-size: 12px;
  font-weight: 600;
  padding: 2px 10px;
  border-radius: var(--radius-sm);
}

.case-country {
  background-color: rgba(26, 58, 92, 0.1);
  color: var(--primary);
}

.case-project {
  background-color: rgba(200, 150, 62, 0.1);
  color: var(--accent-dark);
}

.case-name {
  font-size: 18px;
  font-weight: 700;
  color: var(--text-primary);
  margin-bottom: 8px;
}

.case-desc {
  font-size: 14px;
  color: var(--text-secondary);
  line-height: 1.7;
  margin-bottom: 16px;
}

.case-result {
  display: flex;
}

.result-badge {
  font-size: 13px;
  font-weight: 600;
  color: #1e7e34;
  background-color: #e6f4ea;
  padding: 4px 14px;
  border-radius: var(--radius-sm);
}

.loading-state,
.error-state,
.empty-state {
  text-align: center;
  padding: 80px 20px;
  color: var(--text-light);
  font-size: 16px;
}

.error-state {
  color: #c62828;
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
