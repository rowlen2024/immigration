<template>
  <div>
    <AdminPageHeader title="项目类目">
      <template #actions>
        <el-button v-if="canWrite" type="primary" @click="openCreate">新建类目</el-button>
      </template>
    </AdminPageHeader>

    <AdminToolbar>
      <el-input
        v-model="searchName"
        aria-label="类目名称"
        placeholder="搜索类目名称..."
        :prefix-icon="Search"
        clearable
        class="admin-search-input"
        @keyup.enter="applyFilters"
        @clear="applyFilters"
      />
      <el-select
        v-model="statusFilter"
        aria-label="类目状态"
        placeholder="状态筛选"
        clearable
        class="admin-filter-select"
        @change="applyFilters"
      >
        <el-option label="全部状态" value="" />
        <el-option label="已发布" value="published" />
        <el-option label="草稿" value="draft" />
      </el-select>
      <el-button @click="applyFilters">搜索</el-button>
      <el-button :icon="Refresh" circle aria-label="刷新类目" :loading="loading" @click="loadList" />
    </AdminToolbar>

    <AdminTableShell :loading="loading">
      <el-table :data="categories" empty-text="暂无匹配的项目类目">
        <el-table-column label="类目名称" min-width="200">
          <template #default="{ row }">
            <div class="admin-row-title">{{ row.name }}</div>
            <div class="admin-row-meta">/projects/category/{{ row.slug }}</div>
          </template>
        </el-table-column>
        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <span :class="['status-pill', row.status]">{{ row.status === 'published' ? '已发布' : '草稿' }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="project_count" label="关联项目" width="110" />
        <el-table-column label="更新时间" width="165">
          <template #default="{ row }">{{ formatDateTime(row.updated_at) }}</template>
        </el-table-column>
        <el-table-column label="操作" width="165" fixed="right">
          <template #default="{ row }">
            <AdminRowActions>
              <el-button link type="primary" @click="openMembers(row)">查看项目</el-button>
              <button
                v-if="canWrite"
                class="action-btn"
                type="button"
                title="编辑类目"
                :aria-label="`编辑 ${row.name}`"
                @click="openEdit(row)"
                v-html="getIconSvg('pencil', 16)"
              ></button>
              <el-popconfirm
                v-if="canWrite"
                title="确定删除该类目？"
                confirm-button-text="删除"
                cancel-button-text="取消"
                @confirm="deleteCategory(row.id)"
              >
                <template #reference>
                  <button
                    class="action-btn danger"
                    type="button"
                    title="删除类目"
                    :aria-label="`删除 ${row.name}`"
                    :disabled="deleting"
                    v-html="getIconSvg('trash-2', 16)"
                  ></button>
                </template>
              </el-popconfirm>
            </AdminRowActions>
          </template>
        </el-table-column>
      </el-table>
    </AdminTableShell>

    <div v-if="total > pageSize" class="admin-pagination-wrap">
      <el-pagination
        :current-page="page"
        :page-size="pageSize"
        :total="total"
        layout="total, prev, pager, next"
        @update:current-page="changePage"
      />
    </div>

    <el-drawer
      v-model="editorVisible"
      :title="editingId === null ? '新建类目' : '编辑类目'"
      size="580px"
      :before-close="closeEditor"
      destroy-on-close
    >
      <el-form ref="formRef" :model="form" :rules="rules" label-position="top">
        <el-form-item label="类目名称" prop="name">
          <el-input v-model="form.name" maxlength="128" />
        </el-form-item>
        <el-form-item label="URL 标识" prop="slug">
          <el-input v-model="form.slug" maxlength="128" :disabled="slugLocked" placeholder="例如 real-estate">
            <template v-if="!slugLocked" #suffix>
              <el-button link type="primary" :disabled="!form.name.trim()" @click="form.slug = generateSlugFromText(form.name)">
                自动生成
              </el-button>
            </template>
          </el-input>
          <div v-if="slugLocked" class="admin-row-meta">类目发布后 URL 标识固定。</div>
        </el-form-item>
        <el-form-item label="类目简介">
          <el-input v-model="form.intro" type="textarea" :rows="4" />
        </el-form-item>
        <el-form-item label="SEO 标题">
          <el-input v-model="form.seo_title" maxlength="128" placeholder="留空时使用类目名称" />
        </el-form-item>
        <el-form-item label="SEO 描述">
          <el-input v-model="form.seo_description" type="textarea" :rows="3" maxlength="512" />
        </el-form-item>
        <el-form-item label="发布状态">
          <el-radio-group v-model="form.status">
            <el-radio value="draft">草稿</el-radio>
            <el-radio value="published">已发布</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>
      <template #footer>
        <AdminDrawerFooter :loading="saving" @cancel="cancelEditor" @confirm="saveCategory" />
      </template>
    </el-drawer>

    <el-drawer v-model="membersVisible" :title="`${activeCategory?.name || ''} · 关联项目`" size="680px" destroy-on-close>
      <AdminTableShell :loading="membersLoading">
        <el-table :data="members" empty-text="暂无关联项目">
          <el-table-column label="项目名称" min-width="175">
            <template #default="{ row }">
              <div class="admin-row-title">{{ row.name }}</div>
              <div class="admin-row-meta">/projects/{{ row.slug }}</div>
            </template>
          </el-table-column>
          <el-table-column label="状态" width="95">
            <template #default="{ row }">
              <span :class="['status-pill', row.status === 1 ? 'published' : 'draft']">
                {{ row.status === 1 ? '已上线' : '未上线' }}
              </span>
            </template>
          </el-table-column>
          <el-table-column label="排序" width="128">
            <template #default="{ row }">
              <el-input-number
                v-if="editingMemberId === row.id"
                v-model="memberSortOrder"
                :min="-2147483648"
                :max="2147483647"
                :precision="0"
                :controls="false"
                :disabled="savingMember"
                :aria-label="`${row.name}的排序`"
                class="admin-full-width"
              />
              <span v-else>{{ row.sort_order }}</span>
            </template>
          </el-table-column>
          <el-table-column v-if="canWrite" label="操作" width="122">
            <template #default="{ row }">
              <template v-if="editingMemberId === row.id">
                <el-button link type="primary" :disabled="!Number.isInteger(memberSortOrder) || savingMember" @click="saveMemberOrder(row.id)">保存</el-button>
                <el-button link :disabled="savingMember" @click="editingMemberId = null">取消</el-button>
              </template>
              <el-button v-else link type="primary" :disabled="savingMember" @click="editMemberOrder(row)">调整排序</el-button>
            </template>
          </el-table-column>
        </el-table>
      </AdminTableShell>
      <div v-if="memberTotal > pageSize" class="admin-pagination-wrap">
        <el-pagination
          :current-page="memberPage"
          :page-size="pageSize"
          :total="memberTotal"
          :disabled="savingMember"
          layout="total, prev, pager, next"
          @update:current-page="changeMemberPage"
        />
      </div>
    </el-drawer>
  </div>
</template>

<script setup lang="ts">
import { Refresh, Search } from '@element-plus/icons-vue';
import type { FormInstance, FormRules } from 'element-plus';
import { getIconSvg } from '~/composables/lucideIcons';
import { useNotify } from '~/composables/useNotify';
import type { AdminProjectCategoryOption, ProjectCategoryStatus } from '~/types/admin-project-category';
import { formatDateTime } from '~/utils/date';
import { generateSlugFromText } from '~/utils/slug';

definePageMeta({ layout: 'admin', middleware: 'auth' });

interface Category extends AdminProjectCategoryOption {
  intro: string;
  seo_title: string;
  seo_description: string;
  project_count: number;
  updated_at: string;
  published_at?: string | null;
}

interface CategoryProject {
  id: number;
  name: string;
  slug: string;
  status: number;
  sort_order: number;
}

const notify = useNotify();
const { hasPermission } = usePermissions();
const canWrite = computed(() => hasPermission('projects:write'));
const pageSize = 20;
const categories = ref<Category[]>([]);
const page = ref(1);
const total = ref(0);
const searchName = ref('');
const statusFilter = ref<ProjectCategoryStatus | ''>('');
const filters = reactive({ name: '', status: '' as ProjectCategoryStatus | '' });
const loading = ref(false);
const deleting = ref(false);
let listRequest = 0;

const loadList = async () => {
  const request = ++listRequest;
  loading.value = true;
  try {
    const result = await useApi()<{ items: Category[]; total: number }>('/admin/project-categories', {
      query: { page: page.value, per_page: pageSize, ...filters },
    });
    if (request !== listRequest) return;
    categories.value = result.items;
    total.value = result.total;
  } catch (err) {
    if (request === listRequest) notify.error(err, '加载类目失败');
  } finally {
    if (request === listRequest) loading.value = false;
  }
};

const applyFilters = () => {
  filters.name = searchName.value;
  filters.status = statusFilter.value || '';
  page.value = 1;
  loadList();
};

const changePage = (value: number) => {
  if (value === page.value) return;
  page.value = value;
  loadList();
};

const defaultForm = () => ({ name: '', slug: '', intro: '', seo_title: '', seo_description: '', status: 'draft' as ProjectCategoryStatus });
const form = reactive(defaultForm());
const formRef = ref<FormInstance>();
const editorVisible = ref(false);
const editingId = ref<number | null>(null);
const slugLocked = ref(false);
const saving = ref(false);
const rules: FormRules = {
  name: [{ required: true, whitespace: true, message: '请输入类目名称', trigger: 'blur' }],
  slug: [
    { required: true, message: '请输入 URL 标识', trigger: 'blur' },
    { max: 128, message: 'URL 标识最多 128 个字符', trigger: 'blur' },
    { pattern: /^[a-z0-9]+(?:-[a-z0-9]+)*$/, message: '使用小写英文字母、数字和连字符', trigger: 'blur' },
  ],
};

const openCreate = () => {
  editingId.value = null;
  slugLocked.value = false;
  Object.assign(form, defaultForm());
  editorVisible.value = true;
};

const openEdit = (row: Category) => {
  const { id, name, slug, intro, seo_title, seo_description, status, published_at } = row;
  editingId.value = id;
  slugLocked.value = status === 'published' || Boolean(published_at);
  Object.assign(form, { name, slug, intro, seo_title, seo_description, status });
  editorVisible.value = true;
};

const closeEditor = (done: () => void) => {
  if (!saving.value) done();
};

const cancelEditor = () => closeEditor(() => { editorVisible.value = false; });

const saveCategory = async () => {
  if (saving.value) return;
  const id = editingId.value;
  saving.value = true;
  try {
    if (!await formRef.value?.validate().catch(() => false)) return;
    const api = useApi();
    const body = { ...form };
    if (id === null) await api('/admin/project-categories', { method: 'POST', body });
    else await api(`/admin/project-categories/${id}`, { method: 'PUT', body });
    notify.success('类目已保存');
    editorVisible.value = false;
    await loadList();
  } catch (err) {
    notify.error(err, '保存类目失败');
  } finally {
    saving.value = false;
  }
};

const deleteCategory = async (id: number) => {
  if (deleting.value) return;
  deleting.value = true;
  try {
    await useApi()(`/admin/project-categories/${id}`, { method: 'DELETE' });
    notify.success('类目已删除');
    if (page.value > 1 && categories.value.length === 1 && categories.value[0]?.id === id) page.value--;
    await loadList();
  } catch (err) {
    notify.error(err, '删除类目失败');
  } finally {
    deleting.value = false;
  }
};

const activeCategory = ref<Category | null>(null);
const membersVisible = ref(false);
const members = ref<CategoryProject[]>([]);
const memberPage = ref(1);
const memberTotal = ref(0);
const membersLoading = ref(false);
let membersRequest = 0;
const editingMemberId = ref<number | null>(null);
const memberSortOrder = ref<number>();
const savingMember = ref(false);

const loadMembers = async () => {
  const id = activeCategory.value?.id;
  if (!id || !membersVisible.value) return;
  const request = ++membersRequest;
  membersLoading.value = true;
  try {
    const result = await useApi()<{ items: CategoryProject[]; total: number }>(`/admin/project-categories/${id}/projects`, {
      query: { page: memberPage.value, per_page: pageSize },
    });
    if (request !== membersRequest || !membersVisible.value) return;
    members.value = result.items;
    memberTotal.value = result.total;
  } catch (err) {
    if (request === membersRequest && membersVisible.value) notify.error(err, '加载关联项目失败');
  } finally {
    if (request === membersRequest) membersLoading.value = false;
  }
};

const openMembers = (row: Category) => {
  activeCategory.value = row;
  memberPage.value = 1;
  memberTotal.value = 0;
  members.value = [];
  editingMemberId.value = null;
  membersVisible.value = true;
  loadMembers();
};

const changeMemberPage = (value: number) => {
  if (value === memberPage.value) return;
  memberPage.value = value;
  editingMemberId.value = null;
  loadMembers();
};

const editMemberOrder = (row: CategoryProject) => {
  editingMemberId.value = row.id;
  memberSortOrder.value = row.sort_order;
};

const saveMemberOrder = async (projectId: number) => {
  const categoryId = activeCategory.value?.id;
  const order = memberSortOrder.value;
  if (!categoryId || typeof order !== 'number' || !Number.isInteger(order) || savingMember.value) return;
  savingMember.value = true;
  try {
    await useApi()(`/admin/project-categories/${categoryId}/projects/${projectId}`, {
      method: 'PUT',
      body: { sort_order: order },
    });
    notify.success('项目排序已保存');
    if (activeCategory.value?.id === categoryId && membersVisible.value) {
      editingMemberId.value = null;
      await loadMembers();
    }
  } catch (err) {
    notify.error(err, '保存项目排序失败');
  } finally {
    savingMember.value = false;
  }
};

onMounted(loadList);
</script>
