<template>
  <div>
    <div class="admin-page-header">
      <h2 class="admin-page-title">FAQ 管理</h2>
      <el-button type="primary" @click="openCreate">新建 FAQ</el-button>
    </div>

    <div class="admin-toolbar">
      <div class="admin-toolbar-row">
        <el-input
          v-model="searchQuery"
          placeholder="搜索问题..."
          :prefix-icon="Search"
          clearable
          class="admin-search-input"
          @input="onSearch"
        />
        <el-select
          v-model="projectFilter"
          placeholder="按项目筛选"
          clearable
          class="admin-project-filter"
          @change="onFilterChange"
        >
          <el-option
            v-for="p in projects"
            :key="p.id"
            :label="p.name"
            :value="String(p.id)"
          />
        </el-select>
      <el-button :icon="Refresh" circle @click="searchQuery='';projectFilter=null;loadList()" :loading="loading" />
      </div>
    </div>

    <div class="admin-table-wrap">
      <el-table :data="list" v-loading="loading">
        <el-table-column prop="question" label="问题" min-width="220">
          <template #default="{ row }">
            <div class="row-title">{{ row.question }}</div>
          </template>
        </el-table-column>
        <el-table-column prop="project_name" label="所属项目" width="160" />
        <el-table-column prop="is_global" label="全局" width="80">
          <template #default="{ row }">
            <span :class="['status-pill', row.is_global ? 'published' : 'draft']">
              {{ row.is_global ? '是' : '否' }}
            </span>
          </template>
        </el-table-column>
        <el-table-column prop="sort_order" label="排序" width="70" />
        <el-table-column label="操作" width="140" fixed="right">
          <template #default="{ row }">
            <div class="table-actions">
              <button class="action-btn" @click="openEdit(row)">编辑</button>
              <el-popconfirm title="确定删除该 FAQ？" confirm-button-text="删除" cancel-button-text="取消" @confirm="handleDelete(row.id)">
                <template #reference>
                  <button class="action-btn danger">删除</button>
                </template>
              </el-popconfirm>
            </div>
          </template>
        </el-table-column>
      </el-table>
    </div>

    <div v-if="!loading && list.length === 0" class="admin-empty-state">
      <div class="empty-icon" v-html="getIconSvg('help-circle', 48)"></div>
      <div class="empty-title">暂无 FAQ</div>
      <div class="empty-desc">点击上方按钮创建第一个 FAQ</div>
      <el-button type="primary" @click="openCreate">新建 FAQ</el-button>
    </div>

    <div class="admin-pagination-wrap" v-if="total > pageSize">
      <el-pagination v-model:current-page="page" :page-size="pageSize" :total="total" layout="total, prev, pager, next" @current-change="loadList" />
    </div>

    <el-drawer v-model="drawerVisible" :title="editingId ? '编辑 FAQ' : '新建 FAQ'" size="560px" destroy-on-close>
      <el-form ref="formRef" :model="form" :rules="rules" label-position="top">
        <el-form-item label="问题" prop="question">
          <el-input v-model="form.question" />
        </el-form-item>
        <el-form-item label="回答" prop="answer">
          <el-input v-model="form.answer" type="textarea" :rows="4" />
        </el-form-item>
        <el-form-item label="所属项目" prop="project_id">
          <el-select v-model="form.project_id" placeholder="选择项目（留空为全局）" clearable filterable>
            <el-option v-for="p in projects" :key="p.id" :label="p.name" :value="p.id" />
          </el-select>
        </el-form-item>
        <el-row :gutter="12">
          <el-col :span="12">
            <el-form-item label="全局" prop="is_global">
              <el-switch v-model="form.is_global" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="排序" prop="sort_order">
              <el-input-number v-model="form.sort_order" :min="0" />
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>
      <template #footer>
        <el-button @click="drawerVisible = false">取消</el-button>
        <el-button type="primary" :loading="saving" @click="handleSave">保存</el-button>
      </template>
    </el-drawer>
  </div>
</template>

<script setup lang="ts">
import { Search, Refresh } from '@element-plus/icons-vue';
import { ElMessage, type FormInstance, type FormRules } from 'element-plus';
import { getIconSvg } from '~/composables/lucideIcons';

definePageMeta({ layout: 'admin', middleware: 'auth' });

interface Faq {
  id: string;
  question: string;
  answer: string;
  project_id: string | null;
  project_name: string;
  is_global: boolean;
  sort_order: number;
}

const list = ref<Faq[]>([]);
const loading = ref(false);
const saving = ref(false);
const page = ref(1);
const pageSize = ref(10);
const total = ref(0);
const projects = ref<{ id: string; name: string }[]>([]);

const drawerVisible = ref(false);
const editingId = ref<string | null>(null);
const formRef = ref<FormInstance>();

const searchQuery = ref('');
const projectFilter = ref<string | null>(null);

let searchTimer: ReturnType<typeof setTimeout>;
const onSearch = () => {
  clearTimeout(searchTimer);
  searchTimer = setTimeout(() => {
    page.value = 1;
    loadList();
  }, 300);
};

const onFilterChange = () => {
  page.value = 1;
  loadList();
};

const defaultForm = () => ({
  question: '',
  answer: '',
  project_id: null as string | null,
  is_global: false,
  sort_order: 0,
});

const form = reactive(defaultForm());

const rules: FormRules = {
  question: [{ required: true, message: '请输入问题', trigger: 'blur' }],
  answer: [{ required: true, message: '请输入回答', trigger: 'blur' }],
};

const loadProjects = async () => {
  try {
    const api = useApi();
    const data = await api<{ items: { id: string; name: string }[] }>('/admin/projects?all=true');
    projects.value = data?.items ?? [];
  } catch {
    projects.value = [];
  }
};

const loadList = async () => {
  loading.value = true;
  try {
    const api = useApi();
    let url = `/admin/faqs?page=${page.value}&per_page=${pageSize.value}`;
    if (searchQuery.value) url += `&search=${encodeURIComponent(searchQuery.value)}`;
    if (projectFilter.value) url += `&project_id=${projectFilter.value}`;
    const data = await api<{ items: Faq[]; total: number }>(url);
    list.value = data.items ?? [];
    total.value = data.total ?? 0;
  } catch {
    list.value = [];
    ElMessage.error('加载FAQ列表失败');
  } finally {
    loading.value = false;
  }
};

const openCreate = () => {
  editingId.value = null;
  Object.assign(form, defaultForm());
  drawerVisible.value = true;
};

const openEdit = (row: Faq) => {
  editingId.value = row.id;
  Object.assign(form, {
    question: row.question,
    answer: row.answer,
    project_id: row.project_id,
    is_global: row.is_global,
    sort_order: row.sort_order,
  });
  drawerVisible.value = true;
};

const handleSave = async () => {
  const valid = await formRef.value?.validate().catch(() => false);
  if (!valid) return;

  saving.value = true;
  try {
    const api = useApi();
    if (editingId.value) {
      await api(`/admin/faqs/${editingId.value}`, { method: 'PUT', body: form });
    } else {
      await api('/admin/faqs', { method: 'POST', body: form });
    }
    drawerVisible.value = false;
    loadList();
  } catch {
    ElMessage.error('操作失败');
  } finally {
    saving.value = false;
  }
};

const handleDelete = async (id: string) => {
  try {
    const api = useApi();
    await api(`/admin/faqs/${id}`, { method: 'DELETE' });
    loadList();
  } catch {
    ElMessage.error('操作失败');
  }
};

onMounted(() => {
  loadProjects();
  loadList();
});
</script>

<style scoped>
.admin-toolbar-row {
  display: flex;
  gap: 12px;
  align-items: center;
}
.admin-project-filter {
  width: 200px;
  flex-shrink: 0;
}
.row-title {
  font-size: 14px;
  font-weight: 500;
  color: var(--color-text);
}
</style>
