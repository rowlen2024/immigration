<template>
  <div class="lawyers-admin">
    <div class="admin-page-header">
      <h2 class="admin-page-title">律师团队</h2>
      <el-button type="primary" @click="openAdd">新增律师</el-button>
    </div>

    <div class="admin-toolbar">
      <el-input
        v-model="searchQuery"
        placeholder="搜索姓名..."
        :prefix-icon="Search"
        clearable
        class="admin-search-input"
        @input="onSearch"
      />
      <el-button :icon="Refresh" circle @click="searchQuery='';loadList()" :loading="loading" />
    </div>

    <div class="admin-table-wrap">
      <div v-if="loading" class="lawyers-loading-mask">
        <span class="lawyers-loading-spinner"></span>
      </div>
      <el-table :data="lawyers" stripe class="admin-table">
        <el-table-column label="照片" width="90">
          <template #default="{ row }">
            <ResponsiveImage v-if="row.photo_url" :src="row.photo_url" variant="thumb" class="lawyer-thumb" />
            <span v-else class="no-photo">—</span>
          </template>
        </el-table-column>
        <el-table-column prop="name" label="姓名" width="100">
          <template #default="{ row }">{{ row.name || '—' }}</template>
        </el-table-column>
        <el-table-column prop="title" label="身份" min-width="140">
          <template #default="{ row }">{{ row.title || '—' }}</template>
        </el-table-column>
        <el-table-column label="标签" min-width="200">
          <template #default="{ row }">
            <template v-if="row.tags?.length">
              <el-tag v-for="tag in row.tags" :key="tag" size="small" class="tag-chip">{{ tag }}</el-tag>
            </template>
            <span v-else class="text-muted">—</span>
          </template>
        </el-table-column>
        <el-table-column prop="sort_order" label="排序" width="70">
          <template #default="{ row }">{{ row.sort_order ?? '—' }}</template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="160">
          <template #default="{ row }">{{ formatDateTime(row.created_at) }}</template>
        </el-table-column>
        <el-table-column prop="updated_at" label="修改时间" width="160">
          <template #default="{ row }">{{ formatDateTime(row.updated_at) }}</template>
        </el-table-column>
        <el-table-column label="操作" width="140" fixed="right">
          <template #default="{ row }">
            <div class="table-actions">
              <button class="action-btn" @click="openEdit(row)">编辑</button>
              <el-popconfirm title="确定删除该律师？" confirm-button-text="删除" cancel-button-text="取消" @confirm="removeLawyer(row.id)">
                <template #reference>
                  <button class="action-btn danger">删除</button>
                </template>
              </el-popconfirm>
            </div>
          </template>
        </el-table-column>
      </el-table>
    </div>

    <div v-if="!loading && lawyers.length === 0" class="admin-empty-state">
      <div class="empty-icon" v-html="getIconSvg('users', 48)"></div>
      <div class="empty-title">暂无律师</div>
      <div class="empty-desc">点击上方按钮添加第一位律师</div>
      <el-button type="primary" @click="openAdd">新增律师</el-button>
    </div>

    <div class="admin-pagination-wrap" v-if="total > pageSize">
      <el-pagination v-model:current-page="page" :page-size="pageSize" :total="total" layout="total, prev, pager, next" @current-change="loadList" />
    </div>

    <el-drawer
      v-model="dialogVisible"
      :title="editingId ? '编辑律师' : '新增律师'"
      size="500px"
      destroy-on-close
    >
      <el-form label-position="top">
        <el-form-item label="照片">
          <ImageInput v-model="form.photo_url" placeholder="照片地址" size-hint="推荐 600×800px (3:4 竖版)" preview-ratio="3 / 4" context="lawyer" />
        </el-form-item>
        <el-form-item label="姓名" required>
          <el-input v-model="form.name" placeholder="律师姓名" />
        </el-form-item>
        <el-form-item label="身份">
          <el-input v-model="form.title" placeholder="如：首席移民律师" />
        </el-form-item>
        <el-form-item label="标签">
          <div class="tag-editor">
            <div v-for="(tag, i) in form.tags" :key="i" class="tag-editor-row">
              <el-input v-model="form.tags[i]" placeholder="标签内容" />
              <el-button type="danger" size="small" @click="form.tags.splice(i, 1)" :disabled="form.tags.length <= 1">×</el-button>
            </div>
            <el-button size="small" @click="form.tags.push('')">+ 添加标签</el-button>
          </div>
        </el-form-item>
        <el-form-item label="排序">
          <el-input-number v-model="form.sort_order" :min="0" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="saveLawyer">确定</el-button>
      </template>
    </el-drawer>
  </div>
</template>

<script setup lang="ts">
definePageMeta({ layout: 'admin', middleware: ['auth'] });
import { Search, Refresh } from '@element-plus/icons-vue';
import { ElMessage } from 'element-plus';
import { getIconSvg } from '~/composables/lucideIcons';
import { useNotify } from '~/composables/useNotify';
import { formatDateTime } from '~/utils/date';
import ImageInput from '~/components/admin/ImageInput.vue';

interface Lawyer {
  id: number;
  photo_url: string;
  name: string;
  title: string;
  tags: string[];
  sort_order: number;
}

const notify = useNotify();

const lawyers = ref<Lawyer[]>([]);
const loading = ref(true);
const page = ref(1);
const pageSize = ref(10);
const total = ref(0);

const dialogVisible = ref(false);
const editingId = ref<number | null>(null);
const form = ref<{ photo_url: string; name: string; title: string; tags: string[]; sort_order: number }>({
  photo_url: '',
  name: '',
  title: '',
  tags: [''],
  sort_order: 0,
});

const searchQuery = ref('');

let searchTimer: ReturnType<typeof setTimeout>;
const onSearch = () => {
  clearTimeout(searchTimer);
  searchTimer = setTimeout(() => {
    page.value = 1;
    loadList();
  }, 300);
};

const loadList = async () => {
  loading.value = true;
  try {
    const api = useApi();
    let url = `/admin/lawyers?page=${page.value}&per_page=${pageSize.value}`;
    if (searchQuery.value) url += `&name=${encodeURIComponent(searchQuery.value)}`;
    const data = await api<{ items: Lawyer[]; total: number }>(url);
    lawyers.value = data.items ?? [];
    total.value = data.total ?? 0;
  } catch {
    lawyers.value = [];
  } finally {
    loading.value = false;
  }
};

function openAdd() {
  editingId.value = null;
  form.value = { photo_url: '', name: '', title: '', tags: [''], sort_order: 0 };
  dialogVisible.value = true;
}

function openEdit(row: Lawyer) {
  editingId.value = row.id;
  form.value = {
    photo_url: row.photo_url,
    name: row.name,
    title: row.title,
    tags: row.tags.length > 0 ? [...row.tags] : [''],
    sort_order: row.sort_order,
  };
  dialogVisible.value = true;
}

async function saveLawyer() {
  if (!form.value.name.trim()) { ElMessage.warning('请填写姓名'); return; }
  form.value.tags = form.value.tags.filter(t => t.trim() !== '');
  if (form.value.tags.length === 0) form.value.tags = [''];

  try {
    const api = useApi();
    if (editingId.value) {
      await api(`/admin/lawyers/${editingId.value}`, { method: 'PUT', body: form.value });
      notify.success('已更新');
    } else {
      await api('/admin/lawyers', { method: 'POST', body: form.value });
      notify.success('已添加');
    }
    dialogVisible.value = false;
    await loadList();
  } catch (e) { notify.error(e, '操作失败'); }
}

async function removeLawyer(id: number) {
  try {
    const api = useApi();
    await api(`/admin/lawyers/${id}`, { method: 'DELETE' });
    notify.success('已删除');
    await loadList();
  } catch (e) { notify.error(e, '操作失败'); }
}

onMounted(() => {
  loadList();
});
</script>

<style scoped>
.admin-table-wrap {
  position: relative;
}

.lawyers-loading-mask {
  position: absolute;
  inset: 0;
  z-index: 10;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(255, 255, 255, 0.85);
}

.lawyers-loading-spinner {
  width: 32px;
  height: 32px;
  border: 3px solid var(--el-border-color-light);
  border-top-color: var(--el-color-primary);
  border-radius: 50%;
  animation: lawyers-loading-spin 0.8s linear infinite;
}

@keyframes lawyers-loading-spin {
  to {
    transform: rotate(360deg);
  }
}

.lawyer-thumb { width: 56px; height: 64px; object-fit: cover; border-radius: 4px; }
.no-photo { color: #ccc; }
.tag-chip { margin-right: 4px; margin-bottom: 4px; }
.tag-editor { display: flex; flex-direction: column; gap: 8px; }
.tag-editor-row { display: flex; gap: 8px; align-items: center; }
</style>
