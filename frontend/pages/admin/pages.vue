<template>
  <div>
    <AdminPageHeader title="页面管理">
      <template #actions>
        <el-button type="primary" @click="openCreate">新建页面</el-button>
      </template>
    </AdminPageHeader>

    <AdminToolbar>
      <el-input
        v-model="searchQuery"
        placeholder="搜索页面标题..."
        :prefix-icon="Search"
        clearable
        class="admin-search-input"
        @input="onSearch"
      />
      <el-select v-model="statusFilter" placeholder="状态筛选" clearable class="admin-filter-select" @change="loadList">
        <el-option label="全部" value="" />
        <el-option label="已发布" value="published" />
        <el-option label="草稿" value="draft" />
      </el-select>
      <el-select v-model="pageTypeFilter" placeholder="页面类型" clearable class="admin-filter-select" @change="loadList">
        <el-option label="全部" value="" />
        <el-option label="默认" value="default" />
        <el-option label="新闻" value="news" />
      </el-select>
      <el-button :icon="Refresh" circle @click="searchQuery='';statusFilter='';pageTypeFilter='';loadList()" :loading="loading" />
    </AdminToolbar>

    <AdminTableShell :loading="loading">
      <el-table :data="list">
        <el-table-column prop="title" label="标题" min-width="180">
          <template #default="{ row }">
            <div>
              <div class="admin-row-title">{{ row.title }}</div>
              <div class="admin-row-meta">/pages/{{ row.slug }}</div>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="template" label="模板" width="100">
          <template #default="{ row }">{{ row.template || '—' }}</template>
        </el-table-column>
        <el-table-column prop="page_type" label="类型" width="80">
          <template #default="{ row }">
            <span>{{ row.page_type === 'news' ? '新闻' : '默认' }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <span :class="['status-pill', row.status === 'published' ? 'published' : 'draft']">
              {{ row.status === 'published' ? '已发布' : '草稿' }}
            </span>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="160">
          <template #default="{ row }">{{ formatDateTime(row.created_at) }}</template>
        </el-table-column>
        <el-table-column prop="updated_at" label="修改时间" width="160">
          <template #default="{ row }">{{ formatDateTime(row.updated_at) }}</template>
        </el-table-column>
        <el-table-column label="操作" width="230" fixed="right">
          <template #default="{ row }">
            <AdminRowActions>
              <button class="action-btn" type="button" title="预览" aria-label="预览" @click="handlePreview(row)" v-html="getIconSvg('external-link', 16)"></button>
              <button
                class="action-btn"
                :class="row.status === 'published' ? 'warning' : 'success'"
                type="button"
                :title="row.status === 'published' ? '下架' : '发布'"
                :aria-label="row.status === 'published' ? '下架' : '发布'"
                @click="handleToggleStatus(row)"
                v-html="getIconSvg(row.status === 'published' ? 'archive' : 'send', 16)"
              ></button>
              <button class="action-btn" type="button" title="编辑" aria-label="编辑" @click="openEdit(row)" v-html="getIconSvg('pencil', 16)"></button>
              <el-popconfirm title="确定删除该页面？" confirm-button-text="删除" cancel-button-text="取消" @confirm="handleDelete(row.id)">
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
      icon="file-text"
      title="暂无页面"
      description="点击上方按钮创建第一个页面"
      action-label="新建页面"
      @action="openCreate"
    />

    <div class="admin-pagination-wrap" v-if="total > pageSize">
      <el-pagination v-model:current-page="page" :page-size="pageSize" :total="total" layout="total, prev, pager, next" @current-change="loadList" />
    </div>

    <el-drawer v-model="drawerVisible" :title="editingId ? '编辑页面' : '新建页面'" size="900px" destroy-on-close>
      <el-form ref="formRef" :model="form" :rules="rules" label-position="top">
        <el-form-item label="标题" prop="title">
          <el-input v-model="form.title" />
        </el-form-item>
        <el-form-item label="标识(slug)" prop="slug">
          <el-input v-model="form.slug">
            <template #suffix>
              <el-button
                link
                type="primary"
                size="small"
                :disabled="!form.title || !form.title.trim()"
                @click="generateSlug"
              >自动生成</el-button>
            </template>
          </el-input>
        </el-form-item>
        <el-form-item label="内容" prop="content">
          <RichEditor v-model="form.content" />
        </el-form-item>
        <el-form-item label="封面图片" prop="cover_image">
          <ImageInput v-model="form.cover_image" placeholder="封面图片URL" size-hint="推荐 360×240px (3:2 横向)" context="page-cover" />
        </el-form-item>
        <el-row :gutter="12">
          <el-col :span="12">
            <el-form-item label="SEO 标题" prop="meta_title">
              <el-input v-model="form.meta_title" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="SEO 描述" prop="meta_description">
              <el-input v-model="form.meta_description" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="12">
          <el-col :span="8">
            <el-form-item label="模板" prop="template">
              <el-select v-model="form.template">
                <el-option label="默认" value="default" />
                <el-option label="全宽" value="fullwidth" />
                <el-option label="落地页" value="landing" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="页面类型" prop="page_type">
              <el-select v-model="form.page_type">
                <el-option label="默认" value="default" />
                <el-option label="新闻" value="news" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="状态" prop="status">
              <el-select v-model="form.status">
                <el-option label="草稿" value="draft" />
                <el-option label="已发布" value="published" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
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
import { generateSlugFromText } from '~/utils/slug';
import ImageInput from '~/components/admin/ImageInput.vue';
import { getIconSvg } from '~/composables/lucideIcons';

definePageMeta({ layout: 'admin', middleware: 'auth' });

const notify = useNotify();

  interface Page {
    id: string;
    title: string;
    slug: string;
    content: string;
    cover_image: string;
    meta_title: string;
    meta_description: string;
    template: string;
    page_type: string;
    status: string;
    created_at: string;
    updated_at: string;
    deleted_at?: string;
  }

const list = ref<Page[]>([]);
const loading = ref(false);
const saving = ref(false);
const page = ref(1);
const pageSize = ref(10);
const total = ref(0);

const drawerVisible = ref(false);
const editingId = ref<string | null>(null);
const formRef = ref<FormInstance>();

const searchQuery = ref('');
const statusFilter = ref('');
const pageTypeFilter = ref('');

let searchTimer: ReturnType<typeof setTimeout>;
const onSearch = () => {
  clearTimeout(searchTimer);
  searchTimer = setTimeout(() => {
    page.value = 1;
    loadList();
  }, 300);
};

const defaultForm = () => ({
  id: undefined as string | undefined,
  title: '',
  slug: '',
  content: '',
  cover_image: '',
  meta_title: '',
  meta_description: '',
  template: 'default',
  page_type: 'default',
  status: 'draft',
} as Page);

const form = reactive(defaultForm());

const rules: FormRules = {
  title: [{ required: true, message: '请输入标题', trigger: 'blur' }],
  slug: [{ required: true, message: '请输入标识', trigger: 'blur' }],
};

const loadList = async () => {
  loading.value = true;
  try {
    const api = useApi();
    let url = `/admin/pages?page=${page.value}&per_page=${pageSize.value}`;
    if (searchQuery.value) url += `&title=${encodeURIComponent(searchQuery.value)}`;
    if (statusFilter.value) url += `&status=${statusFilter.value}`;
    if (pageTypeFilter.value) url += `&page_type=${encodeURIComponent(pageTypeFilter.value)}`;
    const data = await api<{ items: Page[]; total: number }>(url);
    list.value = data.items ?? [];
    total.value = data.total ?? 0;
  } catch {
    list.value = [];
    ElMessage.error('加载页面列表失败');
  } finally {
    loading.value = false;
  }
};

const openCreate = () => {
  editingId.value = null;
  Object.assign(form, defaultForm());
  drawerVisible.value = true;
};

const openEdit = (row: Page) => {
  editingId.value = row.id;
  const { id, created_at, updated_at, deleted_at, ...cleanRow } = row;
  Object.assign(form, cleanRow);
  drawerVisible.value = true;
};

const handleSave = async () => {
  const valid = await formRef.value?.validate().catch(() => false);
  if (!valid) return;

  saving.value = true;
  try {
    const api = useApi();
    if (editingId.value) {
      await api(`/admin/pages/${editingId.value}`, {
        method: 'PUT',
        body: form,
      });
      notify.success('更新成功');
    } else {
      await api('/admin/pages', { method: 'POST', body: form });
      notify.success('创建成功');
    }
    drawerVisible.value = false;
    loadList();
  } catch (e) {
    notify.error(e, editingId.value ? '更新页面失败' : '创建页面失败');
  } finally {
    saving.value = false;
  }
};

const handlePreview = (row: Page) => {
  window.open(`/preview/page/${row.slug}`, '_blank');
};

const handleToggleStatus = async (row: Page) => {
  const newStatus = row.status === 'published' ? 'draft' : 'published';
  try {
    const api = useApi();
    await api(`/admin/pages/${row.id}`, {
      method: 'PUT',
      body: { ...row, status: newStatus },
    });
    row.status = newStatus;
    notify.success(newStatus === 'published' ? '已发布' : '已下架');
  } catch (e) {
    notify.error(e, '操作失败');
  }
};

const handleDelete = async (id: string) => {
  try {
    const api = useApi();
    await api(`/admin/pages/${id}`, { method: 'DELETE' });
    notify.success('已删除');
    loadList();
  } catch (e) {
    notify.error(e, '删除页面失败');
  }
};

const generateSlug = () => {
  if (!form.title || !form.title.trim()) {
    ElMessage.warning('请先输入标题');
    return;
  }
  const slug = generateSlugFromText(form.title);
  if (!slug) {
    ElMessage.warning('未识别到可生成 slug 的有效内容');
    return;
  }
  form.slug = slug;
};

onMounted(() => {
  loadList();
});
</script>
