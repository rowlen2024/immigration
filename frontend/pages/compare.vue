<template>
  <div class="compare-page">
    <div class="container">
      <ProjectBreadcrumb />

      <h1 class="page-title">项目对比</h1>
      <p class="page-subtitle">选择两个移民项目进行详细对比，帮助您做出最佳选择</p>

      <!-- Project Selectors -->
      <div class="compare-selectors">
        <div class="selector-group">
          <label class="selector-label">项目A</label>
          <select v-model="selectedA" class="selector-dropdown" @change="onSelect" :disabled="projectListPending">
            <option value="">{{ projectListPending ? '加载中...' : '-- 请选择项目 --' }}</option>
            <option v-for="proj in projectOptions" :key="proj.slug" :value="proj.slug">
              {{ proj.title }}
            </option>
          </select>
        </div>

        <div class="selector-vs">VS</div>

        <div class="selector-group">
          <label class="selector-label">项目B</label>
          <select v-model="selectedB" class="selector-dropdown" @change="onSelect" :disabled="projectListPending">
            <option value="">{{ projectListPending ? '加载中...' : '-- 请选择项目 --' }}</option>
            <option v-for="proj in projectOptions" :key="proj.slug" :value="proj.slug">
              {{ proj.title }}
            </option>
          </select>
        </div>
      </div>

      <!-- Same project warning -->
      <div v-if="sameProjectWarning" class="same-project-warning">
        请选择两个不同的项目进行对比
      </div>

      <!-- Comparison Table -->
      <div v-if="selectedA && selectedB && !sameProjectWarning" class="comparison-result">
        <div v-if="comparePending" class="loading-state">加载对比数据...</div>
        <div v-else-if="compareError" class="error-state">{{ compareError }}</div>
        <div v-else-if="comparison" class="comparison-table-wrapper">
          <table class="comparison-table">
            <thead>
              <tr>
                <th class="row-label-col">对比项目</th>
                <th class="proj-col-a">{{ comparison.projects[0]?.title || '项目A' }}</th>
                <th class="proj-col-b">{{ comparison.projects[1]?.title || '项目B' }}</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="row in comparison.rows" :key="row.label">
                <td class="row-label">{{ row.label }}</td>
                <td :class="getColClass(row.a, row.b, 'a')">{{ row.a }}</td>
                <td :class="getColClass(row.a, row.b, 'b')">{{ row.b }}</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <div v-if="!(selectedA && selectedB) && !sameProjectWarning" class="compare-placeholder">
        <p>请从上方选择两个项目进行对比</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
useSeo({ title: '项目对比' });

interface ProjectOption {
  slug: string;
  title: string;
}

const { data: projectListRaw, pending: projectListPending } = await useFetch<{
  data?: Array<{ slug: string; name: string }>;
}>('/api/v1/projects', {
  query: { per_page: 100 },
  onResponseError() {
    // dropdown will be empty if API fails
  },
});

const projectOptions = computed<ProjectOption[]>(() => {
  const raw = projectListRaw.value as any;
  const items = raw?.data as Array<{ slug: string; name: string }> | undefined;
  if (items && items.length > 0) {
    return items.map((p) => ({ slug: p.slug, title: p.name }));
  }
  return [];
});

const selectedA = ref('');
const selectedB = ref('');

const sameProjectWarning = computed(() =>
  selectedA.value && selectedB.value && selectedA.value === selectedB.value
);

interface ComparisonData {
  projects: Array<{ title: string; slug: string }>;
  rows: Array<{ label: string; a: string; b: string }>;
}

const {
  data: comparisonRaw,
  pending: comparePending,
  error: compareErrorRaw,
  refresh: refreshComparison,
} = useFetch<ComparisonData>('/api/v1/projects/compare', {
  query: computed(() => {
    if (selectedA.value && selectedB.value && selectedA.value !== selectedB.value) {
      return { slugs: `${selectedA.value},${selectedB.value}` };
    }
    return {};
  }),
  immediate: false,
});

const comparison = computed(() => {
  const raw = comparisonRaw.value as any;
  if (raw?.rows) return raw;
  if (raw?.data?.rows) return raw.data;
  return null;
});
const compareError = computed(() =>
  compareErrorRaw.value ? '加载对比数据失败，请重试' : null
);

const onSelect = () => {
  if (selectedA.value && selectedB.value && selectedA.value !== selectedB.value) {
    refreshComparison();
  }
};

const getColClass = (_valueA: string, _valueB: string, col: 'a' | 'b') => {
  if (col === 'a') return 'col-a';
  return 'col-b';
};

// Trigger initial fetch if both selected from query params
const route = useRoute();
onMounted(() => {
  const queryA = route.query.a as string | undefined;
  const queryB = route.query.b as string | undefined;
  if (queryA && queryB && queryA !== queryB) {
    selectedA.value = queryA;
    selectedB.value = queryB;
    nextTick(() => {
      if (selectedA.value && selectedB.value) {
        refreshComparison();
      }
    });
  }
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

.compare-selectors {
  display: flex;
  align-items: flex-end;
  gap: 24px;
  margin-bottom: 48px;
  background-color: var(--bg-light);
  padding: 32px;
  border-radius: var(--radius-lg);
}

.selector-group {
  flex: 1;
}

.selector-label {
  display: block;
  font-size: 14px;
  font-weight: 600;
  color: var(--text-secondary);
  margin-bottom: 8px;
}

.selector-dropdown {
  width: 100%;
  padding: 12px 16px;
  font-size: 15px;
  font-family: var(--font-sans);
  border: 2px solid var(--border-color);
  border-radius: var(--radius-md);
  background-color: var(--bg-white);
  color: var(--text-primary);
  cursor: pointer;
  transition: border-color 0.3s ease;
}

.selector-dropdown:focus {
  outline: none;
  border-color: var(--accent);
}

.selector-vs {
  font-size: 24px;
  font-weight: 800;
  color: var(--accent);
  padding-bottom: 12px;
  flex-shrink: 0;
}

.comparison-table-wrapper {
  overflow-x: auto;
}

.comparison-table {
  width: 100%;
  border-collapse: collapse;
  font-size: 14px;
}

.comparison-table thead {
  background-color: var(--primary);
  color: var(--bg-white);
}

.comparison-table th {
  padding: 16px;
  font-weight: 600;
  text-align: left;
}

.comparison-table th:first-child {
  border-radius: var(--radius-md) 0 0 0;
}

.comparison-table th:last-child {
  border-radius: 0 var(--radius-md) 0 0;
}

.row-label-col {
  width: 150px;
}

.comparison-table td {
  padding: 14px 16px;
  border-bottom: 1px solid var(--border-color);
  color: var(--text-secondary);
  line-height: 1.6;
}

.comparison-table tbody tr:nth-child(even) {
  background-color: var(--bg-light);
}

.row-label {
  font-weight: 600;
  color: var(--text-primary);
  white-space: nowrap;
}

.col-a,
.col-b {
  min-width: 200px;
}

.same-project-warning {
  text-align: center;
  padding: 16px 20px;
  margin-bottom: 24px;
  color: #e65100;
  font-size: 15px;
  font-weight: 500;
  background-color: #fff3e0;
  border: 1px solid #ffcc80;
  border-radius: var(--radius-md);
}

.compare-placeholder {
  text-align: center;
  padding: 80px 20px;
  color: var(--text-light);
  font-size: 16px;
  background-color: var(--bg-light);
  border-radius: var(--radius-lg);
}

.comparison-result {
  margin-bottom: 60px;
}

.loading-state,
.error-state {
  text-align: center;
  padding: 40px;
  color: var(--text-light);
  font-size: 16px;
}

.error-state {
  color: #c62828;
}

@media (max-width: 767px) {
  .compare-selectors {
    flex-direction: column;
    align-items: stretch;
    gap: 16px;
  }

  .selector-vs {
    text-align: center;
    padding-bottom: 0;
  }

  .page-title {
    font-size: 28px;
  }
}
</style>
