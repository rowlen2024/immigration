<template>
  <div>
    <AdminPageHeader
      title="控制台"
      description="查看关键指标、最近咨询和常用后台入口"
    />

    <div class="admin-dashboard-body">
      <AdminLoadingOverlay :show="loading" />

      <!-- 5 stat cards -->
      <div class="admin-stat-grid">
        <div class="admin-stat-card" v-for="card in statCards" :key="card.label">
          <div class="admin-stat-top">
            <span class="admin-stat-label">{{ card.label }}</span>
            <span :class="['admin-stat-icon', `admin-stat-icon--${card.tone}`]" v-html="card.icon"></span>
          </div>
          <div class="admin-stat-value">{{ card.value }}</div>
          <div class="admin-stat-trend" :class="card.trendClass">{{ card.trend }}</div>
        </div>
      </div>

      <!-- 2-column: recent leads + quick actions -->
      <div class="admin-dashboard-grid">
        <div class="admin-dashboard-left">
          <div class="admin-card">
            <h3 class="admin-section-title">最近咨询</h3>
            <div v-if="recentLeads.length === 0" class="admin-empty-hint">暂无咨询记录</div>
            <div v-else class="admin-recent-list">
              <div v-for="lead in recentLeads" :key="lead.id" class="admin-recent-item">
                <span class="admin-recent-name">{{ lead.name }}</span>
                <span class="admin-recent-project">{{ lead.project_name || lead.interested_project || '—' }}</span>
                <span class="admin-recent-time">{{ formatRelativeTime(lead.created_at) }}</span>
              </div>
            </div>
          </div>
        </div>
        <div class="admin-dashboard-right">
          <div class="admin-card">
            <h3 class="admin-section-title">快捷操作</h3>
            <div class="admin-links-grid">
              <NuxtLink to="/admin/projects" class="admin-quick-link">
                <span v-html="getIconSvg('folder', 18)"></span>
                <span>新建项目</span>
              </NuxtLink>
              <NuxtLink to="/admin/pages" class="admin-quick-link">
                <span v-html="getIconSvg('file-text', 18)"></span>
                <span>新建页面</span>
              </NuxtLink>
              <NuxtLink to="/admin/faqs" class="admin-quick-link">
                <span v-html="getIconSvg('message-circle', 18)"></span>
                <span>新建 FAQ</span>
              </NuxtLink>
              <NuxtLink to="/admin/leads" class="admin-quick-link">
                <span v-html="getIconSvg('messages-square', 18)"></span>
                <span>查看咨询</span>
              </NuxtLink>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { getIconSvg } from '~/composables/lucideIcons';

definePageMeta({ layout: 'admin', middleware: 'auth' });

interface Trend {
  key: string;
  direction: 'up' | 'down' | 'neutral';
  percent: number;
  label: string;
}

interface DashboardStats {
  totalProjects: number;
  totalPages: number;
  totalLeads: number;
  totalCases: number;
  unreadLeads: number;
  trends: Trend[];
}

interface RecentLead {
  id: string;
  name: string;
  interested_project: string;
  project_name: string;
  created_at: string;
}

const stats = reactive<DashboardStats>({
  totalProjects: 0,
  totalPages: 0,
  totalLeads: 0,
  totalCases: 0,
  unreadLeads: 0,
  trends: [],
});

const recentLeads = ref<RecentLead[]>([]);
const loading = ref(true);

const loadStats = async () => {
  try {
    const api = useApi();
    const data = await api<DashboardStats>('/admin/dashboard/stats');
    Object.assign(stats, data);
  } catch {
    // silently fail, stats stay at 0
  }
};

const loadRecentLeads = async () => {
  try {
    const api = useApi();
    const data = await api<{ items: RecentLead[] }>('/admin/leads?page=1&per_page=5');
    recentLeads.value = data.items ?? [];
  } catch {
    recentLeads.value = [];
  }
};

function getTrend(key: string): { text: string; cls: string } {
  const t = stats.trends.find((t) => t.key === key);
  if (!t) return { text: '--', cls: 'neutral' };
  const cls = t.direction === 'up' ? 'up' : t.direction === 'down' ? 'down' : 'neutral';
  return { text: t.label, cls };
}

function formatRelativeTime(iso: string): string {
  if (!iso) return '';
  const now = Date.now();
  const then = new Date(iso).getTime();
  const diffMs = now - then;
  const diffSec = Math.floor(diffMs / 1000);
  const diffMin = Math.floor(diffSec / 60);
  const diffHour = Math.floor(diffMin / 60);
  const diffDay = Math.floor(diffHour / 24);

  if (diffMin < 1) return '刚刚';
  if (diffMin < 60) return `${diffMin}分钟前`;
  if (diffHour < 24) return `${diffHour}小时前`;
  if (diffDay === 1) return '昨天';
  if (diffDay < 7) return `${diffDay}天前`;

  const d = new Date(iso);
  return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')}`;
}

const statCards = computed(() => [
  {
    label: '项目总数',
    value: stats.totalProjects ?? 0,
    icon: getIconSvg('folder', 20),
    tone: 'primary',
    trend: getTrend('projects').text,
    trendClass: getTrend('projects').cls,
  },
  {
    label: '咨询总数',
    value: stats.totalLeads ?? 0,
    icon: getIconSvg('messages-square', 20),
    tone: 'warning',
    trend: getTrend('leads').text,
    trendClass: getTrend('leads').cls,
  },
  {
    label: '页面总数',
    value: stats.totalPages ?? 0,
    icon: getIconSvg('file-text', 20),
    tone: 'success',
    trend: getTrend('pages').text,
    trendClass: getTrend('pages').cls,
  },
  {
    label: '案例总数',
    value: stats.totalCases ?? 0,
    icon: getIconSvg('shield', 20),
    tone: 'primary',
    trend: getTrend('cases').text,
    trendClass: getTrend('cases').cls,
  },
  {
    label: '未读咨询',
    value: stats.unreadLeads ?? 0,
    icon: getIconSvg('bell', 20),
    tone: 'danger',
    trend: '待处理',
    trendClass: 'neutral',
  },
]);

onMounted(async () => {
  loading.value = true;
  try {
    await Promise.all([loadStats(), loadRecentLeads()]);
  } finally {
    loading.value = false;
  }
});
</script>

<style scoped>
.admin-dashboard-body {
  position: relative;
}

/* Stat card top row: label + icon side by side */
.admin-stat-top {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  margin-bottom: 12px;
}

/* Colored rounded square icon container */
.admin-stat-icon {
  width: 36px;
  height: 36px;
  border-radius: var(--radius-sm);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.admin-stat-icon :deep(svg) {
  display: block;
}

.admin-stat-icon--primary {
  color: var(--color-primary);
  background: var(--color-info-soft);
}

.admin-stat-icon--warning {
  color: var(--color-warning);
  background: var(--color-warning-soft);
}

.admin-stat-icon--success {
  color: var(--color-success);
  background: var(--color-success-soft);
}

.admin-stat-icon--danger {
  color: var(--color-danger);
  background: var(--color-danger-soft);
}

/* 5-column stat grid */
.admin-stat-grid {
  display: grid;
  grid-template-columns: repeat(5, 1fr);
  gap: 16px;
  margin-bottom: 16px;
}

@media (max-width: 1200px) {
  .admin-stat-grid {
    grid-template-columns: repeat(3, 1fr);
  }
}

@media (max-width: 767px) {
  .admin-stat-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}

/* 2-column dashboard layout */
.admin-dashboard-grid {
  display: grid;
  grid-template-columns: 2fr 1fr;
  gap: 16px;
}

@media (max-width: 767px) {
  .admin-dashboard-grid {
    grid-template-columns: 1fr;
  }
}

/* Recent leads list */
.admin-recent-list {
  display: flex;
  flex-direction: column;
}

.admin-recent-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 0;
  border-bottom: 1px solid var(--color-border-light);
}

.admin-recent-item:last-child {
  border-bottom: none;
}

.admin-recent-name {
  font-size: 14px;
  font-weight: 500;
  color: var(--color-text);
  flex-shrink: 0;
}

.admin-recent-project {
  font-size: 13px;
  color: var(--color-text-secondary);
  flex: 1;
  min-width: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.admin-recent-time {
  font-size: 12px;
  color: var(--color-text-muted);
  flex-shrink: 0;
}

.admin-empty-hint {
  text-align: center;
  padding: 32px 0;
  font-size: 13px;
  color: var(--color-text-muted);
}

/* Quick link icon sizing */
.admin-quick-link :deep(svg) {
  display: block;
  flex-shrink: 0;
}
</style>
