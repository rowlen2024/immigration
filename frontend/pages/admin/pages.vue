<template>
  <div>
    <div class="admin-page-header">
      <h2 class="admin-page-title">页面管理</h2>
      <el-button type="primary" @click="openCreate">新建页面</el-button>
    </div>

    <div class="admin-toolbar">
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
    </div>

    <div class="admin-table-wrap">
      <el-table :data="list" v-loading="loading">
        <el-table-column prop="title" label="标题" min-width="180">
          <template #default="{ row }">
            <div>
              <div class="row-title">{{ row.title }}</div>
              <div class="row-meta">/pages/{{ row.slug }}</div>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="template" label="模板" width="100" />
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
        <el-table-column label="操作" width="230" fixed="right">
          <template #default="{ row }">
            <div class="table-actions">
              <button class="action-btn" @click="handlePreview(row)">预览</button>
              <button
                class="action-btn"
                :class="row.status === 'published' ? 'warning' : 'success'"
                @click="handleToggleStatus(row)"
              >{{ row.status === 'published' ? '下架' : '发布' }}</button>
              <button class="action-btn" @click="openEdit(row)">编辑</button>
              <el-popconfirm title="确定删除该页面？" confirm-button-text="删除" cancel-button-text="取消" @confirm="handleDelete(row.id)">
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
      <div class="empty-icon" v-html="getIconSvg('file-text', 48)"></div>
      <div class="empty-title">暂无页面</div>
      <div class="empty-desc">点击上方按钮创建第一个页面</div>
      <el-button type="primary" @click="openCreate">新建页面</el-button>
    </div>

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
          <ImageInput v-model="form.cover_image" placeholder="封面图片URL" />
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
        <el-button @click="drawerVisible = false">取消</el-button>
        <el-button type="primary" :loading="saving" @click="handleSave">保存</el-button>
      </template>
    </el-drawer>
  </div>
</template>

<script setup lang="ts">
import { Search, Refresh } from '@element-plus/icons-vue';
import { ElMessage, type FormInstance, type FormRules } from 'element-plus';
import { useNotify } from '~/composables/useNotify';
import { getIconSvg } from '~/composables/lucideIcons';
import { pinyin } from 'pinyin-pro';
import ImageInput from '~/components/admin/ImageInput.vue';

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

const defaultForm = (): Partial<Page> => ({
  id: undefined,
  title: '',
  slug: '',
  content: '',
  cover_image: '',
  meta_title: '',
  meta_description: '',
  template: 'default',
  page_type: 'default',
  status: 'draft',
});

const form = reactive<Partial<Page>>(defaultForm());

const rules: FormRules = {
  title: [{ required: true, message: '请输入标题', trigger: 'blur' }],
  slug: [{ required: true, message: '请输入标识', trigger: 'blur' }],
};

const loadList = async () => {
  loading.value = true;
  try {
    const api = useApi();
    let url = `/admin/pages?page=${page.value}&per_page=${pageSize.value}`;
    if (searchQuery.value) url += `&search=${encodeURIComponent(searchQuery.value)}`;
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
  const arr = pinyin(form.title, { toneType: 'none', type: 'array' });
  const nonEmpty = arr.filter((s: string) => s.trim() !== '');
  if (nonEmpty.length === 0) {
    ElMessage.warning('未识别到可生成拼音的文字');
    return;
  }
  form.slug = nonEmpty.join('-').toLowerCase();
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

.row-meta {
  font-size: 12px;
  color: var(--color-text-muted);
  margin-top: 2px;
}
</style>
