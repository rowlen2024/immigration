<template>
  <div>
    <div class="admin-page-header">
      <h2 class="admin-page-title">咨询管理</h2>
    </div>

    <div class="admin-toolbar">
      <el-select v-model="statusFilter" placeholder="状态筛选" clearable @change="loadList">
        <el-option label="全部" value="" />
        <el-option label="新咨询" value="new" />
        <el-option label="已联系" value="contacted" />
        <el-option label="已认证" value="qualified" />
        <el-option label="已关闭" value="closed" />
      </el-select>
      <el-button :icon="Refresh" circle @click="statusFilter='';loadList()" :loading="loading" />
    </div>

    <div class="admin-table-wrap">
      <el-table :data="list" v-loading="loading" highlight-current-row>
        <el-table-column prop="name" label="姓名" width="110">
          <template #default="{ row }">
            <div class="row-title">{{ row.name }}</div>
          </template>
        </el-table-column>
        <el-table-column prop="phone" label="电话" width="140">
          <template #default="{ row }">{{ row.phone || '—' }}</template>
        </el-table-column>
        <el-table-column prop="email" label="邮箱" min-width="180">
          <template #default="{ row }">{{ row.email || '—' }}</template>
        </el-table-column>
        <el-table-column label="感兴趣项目" width="150">
          <template #default="{ row }">{{ row.project_name || row.interested_project || '—' }}</template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <span :class="['status-pill', statusPillClass(row.status)]">
              {{ statusLabel(row.status) }}
            </span>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="170">
          <template #default="{ row }">{{ row.created_at || '—' }}</template>
        </el-table-column>
        <el-table-column label="操作" width="90" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link size="small" @click.stop="openDetail(row)">详情</el-button>
          </template>
        </el-table-column>
      </el-table>
    </div>

    <div v-if="!loading && list.length === 0" class="admin-empty-state">
      <div class="empty-icon" v-html="getIconSvg('message-circle', 48)"></div>
      <div class="empty-title">暂无咨询</div>
    </div>

    <div class="admin-pagination-wrap" v-if="total > pageSize">
      <el-pagination v-model:current-page="page" :page-size="pageSize" :total="total" layout="total, prev, pager, next" @current-change="loadList" />
    </div>

    <el-drawer v-model="detailVisible" title="咨询详情" size="500px" destroy-on-close>
      <div v-if="currentLead" class="admin-detail">
        <div class="admin-detail-row">
          <span class="admin-detail-label">姓名</span>
          <span class="admin-detail-value">{{ currentLead.name }}</span>
        </div>
        <div class="admin-detail-row">
          <span class="admin-detail-label">电话</span>
          <span class="admin-detail-value">{{ currentLead.phone }}</span>
        </div>
        <div class="admin-detail-row">
          <span class="admin-detail-label">邮箱</span>
          <span class="admin-detail-value">{{ currentLead.email || '-' }}</span>
        </div>
        <div class="admin-detail-row">
          <span class="admin-detail-label">感兴趣项目</span>
          <span class="admin-detail-value">{{ currentLead.project_name || currentLead.interested_project || '-' }}</span>
        </div>
        <div class="admin-detail-row">
          <span class="admin-detail-label">创建时间</span>
          <span class="admin-detail-value">{{ currentLead.created_at }}</span>
        </div>

        <el-divider />

        <el-form label-position="top">
          <el-form-item label="状态">
            <el-select v-model="editStatus">
              <el-option label="新咨询" value="new" />
              <el-option label="已联系" value="contacted" />
              <el-option label="已认证" value="qualified" />
              <el-option label="已关闭" value="closed" />
            </el-select>
          </el-form-item>
          <el-form-item label="备注">
            <el-input v-model="editNotes" type="textarea" :rows="4" placeholder="添加备注..." />
          </el-form-item>
        </el-form>
      </div>
      <template #footer>
        <el-button @click="detailVisible = false">关闭</el-button>
        <el-button type="primary" :loading="updating" @click="handleUpdateStatus">保存</el-button>
      </template>
    </el-drawer>
  </div>
</template>

<script setup lang="ts">
import { ElMessage } from 'element-plus';
import { Refresh } from '@element-plus/icons-vue';
import { useNotify } from '~/composables/useNotify';
import { getIconSvg } from '~/composables/lucideIcons';

definePageMeta({ layout: 'admin', middleware: 'auth' });

const notify = useNotify();

interface Lead {
  id: string;
  name: string;
  phone: string;
  email: string;
  interested_project: string;
  project_name: string;
  status: string;
  notes: string;
  created_at: string;
}

const list = ref<Lead[]>([]);
const loading = ref(false);
const updating = ref(false);
const page = ref(1);
const pageSize = ref(10);
const total = ref(0);
const statusFilter = ref('');

const detailVisible = ref(false);
const currentLead = ref<Lead | null>(null);
const editStatus = ref('');
const editNotes = ref('');

const statusPillClass = (status: string) => {
  const map: Record<string, string> = {
    new: 'info',
    contacted: 'warning',
    qualified: 'published',
    closed: 'danger',
  };
  return map[status] || 'info';
};

const statusLabel = (status: string) => {
  const map: Record<string, string> = {
    new: '新咨询',
    contacted: '已联系',
    qualified: '已认证',
    closed: '已关闭',
  };
  return map[status] || status;
};

const loadList = async () => {
  loading.value = true;
  try {
    const api = useApi();
    let url = `/admin/leads?page=${page.value}&per_page=${pageSize.value}`;
    if (statusFilter.value) {
      url += `&status=${statusFilter.value}`;
    }
    const data = await api<{ items: Lead[]; total: number }>(url);
    list.value = data.items ?? [];
    total.value = data.total ?? 0;
  } catch {
    list.value = [];
    ElMessage.error('加载咨询列表失败');
  } finally {
    loading.value = false;
  }
};

const openDetail = (row: Lead) => {
  currentLead.value = row;
  editStatus.value = row.status;
  editNotes.value = row.notes || '';
  detailVisible.value = true;
};

const handleUpdateStatus = async () => {
  if (!currentLead.value) return;

  updating.value = true;
  try {
    const api = useApi();
    await api(`/admin/leads/${currentLead.value.id}`, {
      method: 'PUT',
      body: { status: editStatus.value, notes: editNotes.value },
    });
    notify.success('已更新');
    detailVisible.value = false;
    loadList();
  } catch (e) {
    notify.error(e, '操作失败');
  } finally {
    updating.value = false;
  }
};

onMounted(() => {
  loadList();
});
</script>

<style scoped>
.row-title {
  font-size: 14px;
  font-weight: 500;
  color: var(--color-text);
}
</style>
