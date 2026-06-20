<template>
  <div class="compare-detail-page">
    <div class="container">
      <ProjectBreadcrumb :label="`${slugA} vs ${slugB}`" />

      <div v-if="pending" class="page-skeleton-wrapper"><PageSkeleton variant="detail" /></div>
      <div v-else-if="error" class="error-state">{{ error }}</div>
      <template v-else-if="comparison">
        <h1 class="page-title">{{ comparison.projectA }} vs {{ comparison.projectB }}</h1>
        <p class="page-subtitle">详细项目对比分析</p>

        <!-- Cost Comparison -->
        <section class="detail-section">
          <h2 class="section-title">投资金额对比</h2>
          <div class="dual-cards">
            <div class="detail-card card-a">
              <div class="card-project-name">{{ comparison.projectA }}</div>
              <div class="card-value">{{ comparison.costA }}</div>
              <p class="card-desc">{{ comparison.costDescA }}</p>
            </div>
            <div class="detail-card card-b">
              <div class="card-project-name">{{ comparison.projectB }}</div>
              <div class="card-value">{{ comparison.costB }}</div>
              <p class="card-desc">{{ comparison.costDescB }}</p>
            </div>
          </div>
        </section>

        <!-- Timeline Comparison -->
        <section class="detail-section">
          <h2 class="section-title">办理周期对比</h2>
          <div class="dual-cards">
            <div class="detail-card card-a">
              <div class="card-project-name">{{ comparison.projectA }}</div>
              <div class="card-value">{{ comparison.timeA }}</div>
            </div>
            <div class="detail-card card-b">
              <div class="card-project-name">{{ comparison.projectB }}</div>
              <div class="card-value">{{ comparison.timeB }}</div>
            </div>
          </div>
        </section>

        <!-- Requirements Comparison -->
        <section class="detail-section">
          <h2 class="section-title">申请条件对比</h2>
          <table class="comparison-table">
            <thead>
              <tr>
                <th>条件项</th>
                <th>{{ comparison.projectA }}</th>
                <th>{{ comparison.projectB }}</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="(row, index) in comparison.requirementRows" :key="index">
                <td class="row-label">{{ row.label }}</td>
                <td>
                  <span v-if="row.a" class="icon-check">&#10003;</span>
                  <span v-else class="icon-cross">&#10007;</span>
                </td>
                <td>
                  <span v-if="row.b" class="icon-check">&#10003;</span>
                  <span v-else class="icon-cross">&#10007;</span>
                </td>
              </tr>
            </tbody>
          </table>
        </section>

        <!-- Summary -->
        <section class="detail-section">
          <h2 class="section-title">综合分析</h2>
          <div class="summary-content">
            <p>{{ comparison.summary }}</p>
          </div>
        </section>

        <!-- CTA -->
        <ConsultCTA
          title="需要专业建议？"
          description="我们的专业顾问可以帮您分析最适合您的移民方案"
        />
      </template>
    </div>
  </div>
</template>

<script setup lang="ts">
const route = useRoute();
const slugA = computed(() => (route.params.a as string) || '');
const slugB = computed(() => (route.params.b as string) || '');

useSeo({ title: '项目对比详情', breadcrumbLabel: `${slugA.value} vs ${slugB.value}` });

interface DetailedComparison {
  projectA: string;
  projectB: string;
  costA: string;
  costB: string;
  costDescA: string;
  costDescB: string;
  timeA: string;
  timeB: string;
  requirementRows: Array<{ label: string; a: boolean; b: boolean }>;
  summary: string;
}

const { data, pending, error, refresh } = await useFetch<DetailedComparison>(
  () => `/api/v1/projects/compare?slugs=${slugA.value},${slugB.value}`,
  { key: computed(() => `public:compare:${slugA.value}:${slugB.value}`) },
);
usePublicDataFreshness(() => [`public:compare:${slugA.value}:${slugB.value}`]);

const comparison = computed(() => {
  if (data.value) return data.value;

  const projectMap: Record<string, Record<string, string | boolean>> = {
    eb5: {
      name: '美国EB-5投资移民',
      cost: '80万美元起',
      costDesc: '投资到TEA目标就业区项目，5年后可尝试返还。额外的律师费和管理费约7-13万美元。',
      time: '24-36个月',
      age18: true,
      investment: true,
      noCriminal: true,
      language: true,
      edu: true,
      residency: true,
    },
    cies: {
      name: '香港资本投资者入境计划',
      cost: '3000万港元',
      costDesc: '投资于获许金融资产（股票、债券、基金等），可自由调配。律师费和审计费约12-20万港元。',
      time: '6-12个月',
      age18: true,
      investment: true,
      noCriminal: true,
      language: true,
      edu: true,
      residency: false,
    },
    panama: {
      name: '巴拿马购房移民',
      cost: '30万美元起',
      costDesc: '购买政府批准的巴拿马房产，律师费和文件费约6000-12000美元。无移民监要求。',
      time: '3-6个月',
      age18: true,
      investment: true,
      noCriminal: true,
      language: true,
      edu: true,
      residency: true,
      noResidencyReq: true,
    },
  };

  const a = projectMap[slugA.value] || {};
  const b = projectMap[slugB.value] || {};

  const getBool = (val: unknown) => Boolean(val);

  return {
    projectA: (a.name as string) || slugA.value,
    projectB: (b.name as string) || slugB.value,
    costA: (a.cost as string) || '',
    costB: (b.cost as string) || '',
    costDescA: (a.costDesc as string) || '',
    costDescB: (b.costDesc as string) || '',
    timeA: (a.time as string) || '',
    timeB: (b.time as string) || '',
    requirementRows: [
      { label: '年满18周岁', a: getBool(a.age18), b: getBool(b.age18) },
      { label: '无犯罪记录', a: getBool(a.noCriminal), b: getBool(b.noCriminal) },
      { label: '无语言要求', a: getBool(a.language), b: getBool(b.language) },
      { label: '无学历要求', a: getBool(a.edu), b: getBool(b.edu) },
      { label: '无居住要求', a: getBool(a.residency || a.noResidencyReq), b: getBool(b.residency || b.noResidencyReq) },
    ],
    summary:
      '以上三个移民项目各有特点：EB-5适合希望获得美国身份的高净值家庭；香港CIES适合希望在亚洲金融中心定居的投资者；巴拿马购房移民门槛最低、速度最快，适合寻求快速获得海外身份的投资者。建议根据自身资产规模、移民目的和时间规划综合考虑。',
  };
});

onMounted(() => {
  $fetch<any>(`/api/v1/projects/compare?slugs=${slugA.value},${slugB.value}`).then(v => { data.value = v }).catch(() => {})
})
</script>

<style scoped>
.page-title {
  font-size: 36px;
  font-weight: 700;
  color: var(--text-primary);
  margin-bottom: 8px;
}

.page-subtitle {
  font-size: 16px;
  color: var(--text-light);
  margin-bottom: 40px;
}

.detail-section {
  padding: 48px 0;
  border-bottom: 1px solid var(--border-color);
}

.detail-section:last-of-type {
  border-bottom: none;
}

.section-title {
  font-size: 28px;
  font-weight: 700;
  color: var(--text-primary);
  margin-bottom: 28px;
  position: relative;
  padding-bottom: 12px;
}

.section-title::after {
  content: '';
  position: absolute;
  bottom: 0;
  left: 0;
  width: 60px;
  height: 3px;
  background-color: var(--accent);
  border-radius: 2px;
}

.dual-cards {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 24px;
}

.detail-card {
  padding: 32px;
  border-radius: var(--radius-lg);
  border: 2px solid var(--border-color);
}

.card-a {
  border-color: var(--primary);
}

.card-b {
  border-color: var(--accent);
}

.card-project-name {
  font-size: 14px;
  font-weight: 600;
  color: var(--text-light);
  margin-bottom: 12px;
}

.card-value {
  font-size: 28px;
  font-weight: 800;
  margin-bottom: 12px;
}

.card-a .card-value {
  color: var(--primary);
}

.card-b .card-value {
  color: var(--accent-dark);
}

.card-desc {
  font-size: 14px;
  color: var(--text-secondary);
  line-height: 1.7;
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
  padding: 14px 16px;
  font-weight: 600;
  text-align: center;
}

.comparison-table td {
  padding: 12px 16px;
  text-align: center;
  border-bottom: 1px solid var(--border-color);
}

.comparison-table tbody tr:nth-child(even) {
  background-color: var(--bg-light);
}

.row-label {
  text-align: left;
  font-weight: 600;
  color: var(--text-primary);
}

.icon-check {
  color: #1e7e34;
  font-weight: 700;
  font-size: 16px;
}

.icon-cross {
  color: #c62828;
  font-weight: 700;
  font-size: 16px;
}

.summary-content {
  font-size: 16px;
  color: var(--text-secondary);
  line-height: 1.9;
}

/* detail-cta styles moved to ConsultCTA component */

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
  .dual-cards {
    grid-template-columns: 1fr;
  }

  .page-title {
    font-size: 26px;
  }
}
</style>
