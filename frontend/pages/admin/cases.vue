<template>
  <div>
    <AdminPageHeader title="案例管理">
      <template #actions>
        <el-button type="primary" @click="openCreate">新建案例</el-button>
      </template>
    </AdminPageHeader>

    <AdminToolbar>
      <el-input
        v-model="searchQuery"
        placeholder="搜索姓名..."
        :prefix-icon="Search"
        clearable
        class="admin-search-input"
        @input="onSearch"
      />
      <el-button :icon="Refresh" circle @click="searchQuery='';loadList()" :loading="loading" />
    </AdminToolbar>

    <AdminTableShell :loading="loading">
      <el-table :data="list">
        <el-table-column prop="name" label="姓名" min-width="140">
          <template #default="{ row }">
            <div class="admin-row-title">{{ row.name }}</div>
          </template>
        </el-table-column>
        <el-table-column prop="country_from" label="来源国家" min-width="120">
          <template #default="{ row }">{{ row.country_from || '—' }}</template>
        </el-table-column>
        <el-table-column label="封面图" width="80">
          <template #default="{ row }">
            <ResponsiveImage v-if="row.photo_url" :src="row.photo_url" variant="thumb" class="admin-thumb" />
            <span v-else class="admin-no-thumb">—</span>
          </template>
        </el-table-column>
        <el-table-column prop="project" label="项目" width="160" >
          <template #default="{ row }">
            <div class="admin-row-title">{{ row.project?.name || '—' }}</div>
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
            <AdminRowActions>
              <button class="action-btn" type="button" title="编辑" aria-label="编辑" @click="openEdit(row)" v-html="getIconSvg('pencil', 16)"></button>
              <el-popconfirm title="确定删除该案例？" confirm-button-text="删除" cancel-button-text="取消" @confirm="handleDelete(row.id)">
                <template #reference>
                  <button class="action-btn danger" type="button" title="删除" aria-label="删除" v-html="getIconSvg('trash-2', 16)"></button>
                </template>
              </el-popconfirm>
            </AdminRowActions>
          </template>
        </el-table-column>
      </el-table>
    </AdminTableShell>

    <AdminEmptyState
      v-if="!loading && list.length === 0"
      icon="users"
      title="暂无案例"
      description="点击上方按钮创建第一个案例"
      action-label="新建案例"
      @action="openCreate"
    />

    <div class="admin-pagination-wrap" v-if="total > pageSize">
      <el-pagination v-model:current-page="page" :page-size="pageSize" :total="total" layout="total, prev, pager, next" @current-change="loadList" />
    </div>

    <el-drawer v-model="drawerVisible" :title="editingId ? '编辑案例' : '新建案例'" size="900px" destroy-on-close>
      <el-form ref="formRef" :model="form" :rules="rules" label-position="top">
        <el-form-item label="姓名" prop="name">
          <el-input v-model="form.name" />
        </el-form-item>
        <el-form-item label="来源国家" prop="country_from">
          <el-input v-model="form.country_from" />
        </el-form-item>
        <el-form-item label="所属项目" prop="project_id">
          <el-select v-model="form.project_id" placeholder="选择项目" filterable>
            <el-option v-for="p in projects" :key="p.id" :label="p.name" :value="p.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="封面图片">
          <ImageInput v-model="form.photo_url" placeholder="图片 URL 或上传" size-hint="推荐 800×450px (16:9 横向)" context="case" />
        </el-form-item>
        <el-form-item label="投资金额" prop="investment_amount">
          <el-input v-model="form.investment_amount" placeholder="如：80万美元" />
        </el-form-item>
        <el-form-item label="投资数额" prop="investment_value">
          <el-input-number v-model="form.investment_value" :min="0" :precision="2" class="admin-full-width" />
        </el-form-item>
        <el-form-item label="办理周期" prop="processing_period">
          <el-input v-model="form.processing_period" placeholder="如：28个月" />
        </el-form-item>
        <el-form-item label="内容" prop="content">
          <RichEditor v-model="form.content" />
        </el-form-item>
        <el-form-item label="排序" prop="sort_order">
          <el-input-number v-model="form.sort_order" :min="0" />
        </el-form-item>
      </el-form>
      <template #footer>
        <AdminDrawerFooter
          :loading="saving"
          @cancel="drawerVisible = false"
          @confirm="handleSave"
        />
      </template>
    </el-drawer>
  </div>
</template>

<script setup lang="ts">
import { Search, Refresh } from '@element-plus/icons-vue';
import { ElMessage, type FormInstance, type FormRules } from 'element-plus';
import { useNotify } from '~/composables/useNotify';
import { formatDateTime } from '~/utils/date';
import ImageInput from '~/components/admin/ImageInput.vue';
import RichEditor from '~/components/RichEditor.vue';
import { getIconSvg } from '~/composables/lucideIcons';

definePageMeta({ layout: 'admin', middleware: 'auth' });

const notify = useNotify();

interface CaseItem {
  id: string;
  name: string;
  country_from: string;
  project: string;
  project_id: string;
  photo_url: string;
  content: string;
  investment_amount: string;
  investment_value: number;
  processing_period: string;
  sort_order: number;
}

const list = ref<CaseItem[]>([]);
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

let searchTimer: ReturnType<typeof setTimeout>;
const onSearch = () => {
  clearTimeout(searchTimer);
  searchTimer = setTimeout(() => {
    page.value = 1;
    loadList();
  }, 300);
};

const defaultForm = () => ({
  name: '',
  country_from: '',
  project_id: '',
  photo_url: '',
  content: '',
  investment_amount: '',
  investment_value: 0,
  processing_period: '',
  sort_order: 0,
});

const form = reactive(defaultForm());

const rules: FormRules = {
  name: [{ required: true, message: '请输入姓名', trigger: 'blur' }],
  country_from: [{ required: true, message: '请输入来源国家', trigger: 'blur' }],
  project_id: [{ required: true, message: '请选择项目', trigger: 'blur' }],
};

const loadProjects = async () => {
  try {
    const api = useApi();
    const data = await api<{ items: { id: string; name: string }[] }>('/admin/projects/options?page=1&per_page=500');
    projects.value = data.items ?? [];
  } catch {
    projects.value = [];
  }
};

const loadList = async () => {
  loading.value = true;
  try {
    const api = useApi();
    let url = `/admin/cases?page=${page.value}&per_page=${pageSize.value}`;
    if (searchQuery.value) url += `&name=${encodeURIComponent(searchQuery.value)}`;
    const data = await api<{ items: CaseItem[]; total: number }>(url);
    list.value = data.items ?? [];
    total.value = data.total ?? 0;
  } catch {
    list.value = [];
    ElMessage.error('加载案例列表失败');
  } finally {
    loading.value = false;
  }
};

const openCreate = () => {
  editingId.value = null;
  Object.assign(form, defaultForm());
  drawerVisible.value = true;
};

const openEdit = (row: CaseItem) => {
  editingId.value = row.id;
  Object.assign(form, {
    name: row.name,
    country_from: row.country_from,
    project_id: row.project_id,
    photo_url: row.photo_url,
    content: row.content,
    investment_amount: row.investment_amount,
    investment_value: row.investment_value,
    processing_period: row.processing_period,
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
      await api(`/admin/cases/${editingId.value}`, { method: 'PUT', body: form });
      notify.success('更新成功');
    } else {
      await api('/admin/cases', { method: 'POST', body: form });
      notify.success('添加成功');
    }
    drawerVisible.value = false;
    loadList();
  } catch (e) {
    notify.error(e, '操作失败');
  } finally {
    saving.value = false;
  }
};

const handleDelete = async (id: string) => {
  try {
    const api = useApi();
    await api(`/admin/cases/${id}`, { method: 'DELETE' });
    notify.success('已删除');
    loadList();
  } catch (e) {
    notify.error(e, '操作失败');
  }
};

onMounted(() => {
  loadProjects();
  loadList();
});
</script>
