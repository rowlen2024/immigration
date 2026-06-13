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
            <NuxtLink :to="'/projects/' + item.project.slug" class="meta-value link">
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

        <div class="case-layout">
          <div class="case-main">
            <div class="case-content" v-html="item.content" />
          </div>

          <aside class="case-sidebar">
            <!-- 关键信息 -->
            <div class="sb-card">
              <h4 class="sb-title">关键信息</h4>
              <dl class="sb-info-list">
                <div v-if="item.country_from" class="sb-info-item">
                  <dt>来源国家</dt>
                  <dd>{{ item.country_from }}</dd>
                </div>
                <div v-if="item.investment_amount" class="sb-info-item">
                  <dt>投资金额</dt>
                  <dd>{{ item.investment_amount }}</dd>
                </div>
                <div v-if="item.processing_period" class="sb-info-item">
                  <dt>办理周期</dt>
                  <dd>{{ item.processing_period }}</dd>
                </div>
                <div v-if="item.project?.name" class="sb-info-item">
                  <dt>所属项目</dt>
                  <dd>
                    <NuxtLink :to="'/projects/' + item.project.slug">{{ item.project.name }}</NuxtLink>
                  </dd>
                </div>
              </dl>
            </div>

            <!-- 相关案例 -->
            <div v-if="relatedCases.length > 0" class="sb-card">
              <h4 class="sb-title">相关案例</h4>
              <ul class="sb-related">
                <li v-for="rc in relatedCases" :key="rc.slug">
                  <NuxtLink :to="'/case/' + rc.slug">
                    <span class="sb-related-name">{{ rc.name }}</span>
                    <span class="sb-related-meta">{{ rc.country_from }}</span>
                  </NuxtLink>
                </li>
              </ul>
            </div>

            <!-- CTA -->
            <div class="sb-cta">
              <p class="sb-cta-text">正在规划类似方案？</p>
              <NuxtLink to="/contact" class="btn-primary sb-cta-btn">咨询同类方案</NuxtLink>
            </div>
          </aside>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { getIconSvg } from '~/composables/lucideIcons'

const route = useRoute();
const slug = route.params.slug as string;

const { data, pending, error, refresh } = await useFetch<{ data: any }>(`/api/v1/cases/${slug}`);

const item = computed(() => data.value?.data ?? null);

import { stripHtml } from '~/utils/html'
import { buildArticleJsonLd, toJsonLdConfig, toJsonLdScripts } from '~/utils/jsonld'

useSeo({
  title: item.value?.name ?? '案例详情',
  description: (() => {
    const c = item.value;
    if (!c) return '';
    const text = c.content ? stripHtml(c.content, 160) : '';
    return text || `${c.name} — ${c.country_from}${c.project?.name ? ' · ' + c.project.name : ''}`;
  })(),
});

const { siteConfig: csConfig } = useMygoSiteConfig();

// Article structured data for rich results
useHead(() => {
  const c = item.value;
  if (!c?.name) return {};
  const base = csConfig.value?.canonical_base || '';
  const pageUrl = base ? base + route.path : undefined;

  return {
    script: toJsonLdScripts(
      buildArticleJsonLd({
        headline: c.name,
        description: c.content ? stripHtml(c.content, 300) : `${c.name} — ${c.country_from}${c.project?.name ? ' · ' + c.project.name : ''}`,
        url: pageUrl,
        image: c.photo_url,
        datePublished: c.created_at,
      }, toJsonLdConfig(csConfig.value)),
    ),
  };
});

// 相关案例：同项目或同国家，排除当前，最多 4 条
const relatedCases = ref<any[]>([])

async function fetchRelatedCases() {
  try {
    const cur = item.value
    if (!cur) return
    let cases: any[] = []

    // 先按项目筛选
    if (cur.project?.id) {
      const res: any = await $fetch(`/api/v1/cases?project_id=${cur.project.id}&page=1&per_page=4`)
      cases = Array.isArray(res?.data)? res.data: []
    }

    // 项目结果不足 4 条，再按国家补充
    if (cases.length < 4 && cur.country_from) {
      const remain = 4 - cases.length
      const existingIds = new Set(cases.map((c: any) => c.id))
      const res2: any = await $fetch(`/api/v1/cases?country_from=${encodeURIComponent(cur.country_from)}&page=1&per_page=${remain}`)
      const byCountry = Array.isArray(res2?.data)? res2.data: []
      const filtered = byCountry.filter((c: any) => !existingIds.has(c.id))
      cases = [...cases, ...filtered].slice(0, 4)
    }

    // 前端过滤掉当前案例
    relatedCases.value = cases.filter((c: any) => c.slug !== slug)
  } catch { relatedCases.value = [] }
}

onMounted(() => {
  $fetch<{ data: any }>(`/api/v1/cases/${slug}`).then(v => { data.value = v }).catch(() => {})
  if (item.value) fetchRelatedCases()
})

watch(item, (v) => { if (v) fetchRelatedCases() })
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
  min-width: 0;
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
  display: inline-block;
  color: var(--primary);
  text-decoration: none;
}

.meta-value.link:hover {
  text-decoration: underline;
}

.meta-value.link:active {
  opacity: 0.7;
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

/* ══════ Two-column layout ══════ */
.case-layout {
  display: grid;
  grid-template-columns: 1fr 320px;
  gap: 40px;
  align-items: start;
}

/* ══════ Sidebar ══════ */
.case-sidebar {
  position: sticky;
  top: calc(var(--header-height) + 24px);
}

.sb-card {
  background: var(--bg-white);
  border-radius: var(--radius-lg);
  padding: 24px;
  box-shadow: var(--shadow-sm);
  margin-bottom: 20px;
}

.sb-title {
  font-family: var(--font-serif);
  font-size: 17px;
  font-weight: 700;
  color: var(--text-primary);
  margin-bottom: 16px;
  padding-bottom: 10px;
  border-bottom: 2px solid var(--accent-soft);
}

.sb-info-list {
  display: flex;
  flex-direction: column;
  gap: 14px;
}

.sb-info-item {
  display: flex;
  justify-content: space-between;
  align-items: baseline;
  gap: 8px;
}

.sb-info-item dt {
  font-size: 13px;
  color: var(--text-light);
  flex-shrink: 0;
}

.sb-info-item dd {
  font-size: 14px;
  font-weight: 600;
  color: var(--text-primary);
  text-align: right;
}

.sb-info-item dd a {
  color: var(--primary);
  text-decoration: none;
}

.sb-info-item dd a:hover {
  text-decoration: underline;
}

.sb-related {
  list-style: none;
  padding: 0;
  margin: 0;
  display: flex;
  flex-direction: column;
  gap: 0;
}

.sb-related li {
  border-bottom: 1px solid var(--border-color);
}

.sb-related li:last-child {
  border-bottom: none;
}

.sb-related a {
  display: flex;
  flex-direction: column;
  gap: 2px;
  padding: 12px 0;
  text-decoration: none;
  transition: color .2s;
}

.sb-related a:hover .sb-related-name {
  color: var(--accent-dark);
}

.sb-related-name {
  font-size: 14px;
  font-weight: 500;
  color: var(--text-primary);
  line-height: 1.4;
  transition: color .2s;
}

.sb-related-meta {
  font-size: 12px;
  color: var(--text-light);
}

.sb-cta {
  background: var(--gradient-hero);
  border-radius: var(--radius-lg);
  padding: 28px 24px;
  text-align: center;
}

.sb-cta-text {
  font-size: 15px;
  color: rgba(255,255,255,.78);
  margin-bottom: 16px;
}

.sb-cta-btn {
  width: 100%;
  display: block;
  text-align: center;
  padding: 12px 0;
  font-size: 15px;
}

.case-content :deep(img) {
  max-width: 100%;
  height: auto;
  border-radius: 4px;
}

.case-content :deep(table) {
  display: block;
  overflow-x: auto;
  -webkit-overflow-scrolling: touch;
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

@media (max-width: 1023px) {
  .case-layout {
    grid-template-columns: 1fr;
    gap: 32px;
  }

  .case-sidebar {
    position: static;
  }
}

@media (max-width: 767px) {
  .case-detail-page {
    padding: 40px 0;
  }

  .case-title {
    font-size: 28px;
  }

  .case-meta {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 16px;
    padding: 16px;
  }

  .meta-label {
    font-size: 14px;
  }

  .meta-value {
    word-break: break-word;
  }

  .meta-value.link {
    min-height: 44px;
    line-height: 44px;
  }

  .case-content {
    padding: 20px;
  }

  .sb-card {
    padding: 20px;
  }

  .sb-cta {
    padding: 24px 20px;
  }
}

@media (max-width: 374px) {
  .case-detail-page {
    padding: 32px 0;
  }

  .case-title {
    font-size: 24px;
  }
}
</style>
