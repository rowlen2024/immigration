<template>
  <div>
    <div class="admin-page-header">
      <h2 class="admin-page-title">导航管理</h2>
      <el-button v-if="!isViewer" type="primary" @click="openCreate()">新建导航</el-button>
    </div>

    <div class="admin-toolbar">
      <el-input
        v-model="searchQuery"
        placeholder="搜索名称..."
        :prefix-icon="Search"
        clearable
        class="admin-search-input"
      />
      <el-select v-model="typeFilter" placeholder="类型筛选" clearable class="admin-filter-select">
        <el-option label="全部" value="" />
        <el-option label="项目" value="project" />
        <el-option label="页面" value="page" />
        <el-option label="自定义" value="custom" />
      </el-select>
      <el-select v-model="positionFilter" placeholder="显示位置" clearable class="admin-filter-select">
        <el-option label="全部" value="" />
        <el-option label="头部" value="header" />
        <el-option label="底部" value="footer" />
      </el-select>
      <el-button :icon="Refresh" circle @click="searchQuery='';typeFilter='';positionFilter='';loadTree()" :loading="loading" />
    </div>

    <div class="admin-table-wrap">
      <div class="nav-tree-header">
        <span class="nav-th nav-th-name">名称</span>
        <span class="nav-th nav-th-sort">排序</span>
        <span class="nav-th nav-th-link">链接</span>
        <span class="nav-th nav-th-type">类型</span>
        <span class="nav-th nav-th-pos">显示位置</span>
        <span class="nav-th nav-th-time">创建时间</span>
        <span v-if="!isViewer" class="nav-th nav-th-actions">操作</span>
      </div>
      <el-tree
        v-loading="loading"
        :data="filteredTreeData"
        node-key="id"
        default-expand-all
        :props="{ children: 'children', label: 'label' }"
        class="nav-tree"
      >
        <template #default="{ node, data }">
          <div class="tree-node">
            <div class="tree-node-info">
              <span class="row-title">{{ data.label }}</span>
              <span v-if="!data.status" class="status-pill draft">已隐藏</span>
            </div>
            <span class="row-sort">{{ data.sort_order }}</span>
            <span class="row-meta">{{ data.link }}</span>
            <span class="row-type">
              <el-tag
                :type="linkTypeTag(data.link_type)"
                size="small"
                effect="plain"
              >
                {{ linkTypeLabel(data.link_type) }}
              </el-tag>
            </span>
            <span class="row-pos">
              <el-tag
                v-if="data.display_position === 'both'"
                type="success"
                size="small"
                effect="plain"
              >头部+底部</el-tag>
              <el-tag
                v-else-if="data.display_position === 'header'"
                size="small"
                effect="plain"
              >头部</el-tag>
              <el-tag
                v-else-if="data.display_position === 'footer'"
                type="warning"
                size="small"
                effect="plain"
              >底部</el-tag>
            </span>
            <span class="row-time">{{ formatDateTime(data.created_at) }}</span>
            <span v-if="!isViewer" class="tree-node-actions">
              <el-tooltip content="添加子级" placement="top">
                <button class="action-btn" @click.stop="openCreate(data.id)" v-html="getIconSvg('plus', 16)"></button>
              </el-tooltip>
              <el-tooltip content="编辑" placement="top">
                <button class="action-btn" @click.stop="openEdit(data)" v-html="getIconSvg('pencil', 16)"></button>
              </el-tooltip>
              <el-popconfirm
                title="确定删除该导航项？"
                confirm-button-text="删除"
                cancel-button-text="取消"
                @confirm="handleDelete(data.id)"
              >
                <template #reference>
                  <button class="action-btn danger" v-html="getIconSvg('trash-2', 16)"></button>
                </template>
              </el-popconfirm>
            </span>
          </div>
        </template>
      </el-tree>
    </div>

    <el-drawer
      v-model="dialogVisible"
      :title="editingId ? '编辑导航' : '新建导航'"
      size="560px"
      destroy-on-close
    >
      <el-form ref="formRef" :model="form" :rules="rules" label-position="top">
        <el-form-item label="名称" prop="label">
          <el-input v-model="form.label" maxlength="255" />
        </el-form-item>

        <el-form-item label="链接类型" prop="link_type">
          <el-radio-group v-model="form.link_type" @change="onLinkTypeChange">
            <el-radio value="project">项目链接</el-radio>
            <el-radio value="page">页面链接</el-radio>
            <el-radio value="custom">自定义链接</el-radio>
          </el-radio-group>
        </el-form-item>

        <el-form-item v-if="form.link_type === 'project'" label="选择项目" prop="project_id">
          <el-select
            v-model="form.project_id"
            placeholder="搜索并选择项目"
            filterable
            clearable
            style="width: 100%"
            @change="onProjectChange"
          >
            <el-option
              v-for="p in projects"
              :key="p.id"
              :label="`${p.name} (/${p.slug})`"
              :value="p.id"
            />
          </el-select>
          <div v-if="form.project_id" class="link-preview">
            链接预览：<code>/projects/{{ selectedProjectSlug }}</code>
          </div>
        </el-form-item>

        <el-form-item v-if="form.link_type === 'page'" label="选择页面" prop="page_id">
          <el-select
            v-model="form.page_id"
            placeholder="搜索并选择页面"
            filterable
            clearable
            style="width: 100%"
            @change="onPageChange"
          >
            <el-option
              v-for="p in pages"
              :key="p.id"
              :label="`${p.title} (/pages/${p.slug})`"
              :value="p.id"
            />
          </el-select>
          <div v-if="form.page_id" class="link-preview">
            链接预览：<code>/pages/{{ selectedPageSlug }}</code>
          </div>
        </el-form-item>

        <el-form-item v-if="form.link_type === 'custom'" label="链接" prop="link">
          <el-input v-model="form.link" placeholder="/path" />
        </el-form-item>

        <el-form-item label="父级导航" prop="parent_id">
          <el-tree-select
            v-model="form.parent_id"
            :data="parentOptions"
            :props="{ children: 'children', label: 'label' }"
            placeholder="选择父级（留空为顶级）"
            clearable
            check-strictly
          />
        </el-form-item>
        <el-row :gutter="16">
          <el-col :span="12">
            <el-form-item label="排序" prop="sort_order">
              <el-input-number v-model="form.sort_order" :min="0" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="状态" prop="status">
              <el-switch v-model="form.status" active-text="显示" inactive-text="隐藏" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="显示位置" prop="display_position">
          <el-radio-group v-model="form.display_position">
            <el-radio value="header">仅头部</el-radio>
            <el-radio value="footer">仅页脚</el-radio>
            <el-radio value="both">头部+页脚</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="saving" @click="handleSave">保存</el-button>
      </template>
    </el-drawer>
  </div>
</template>

<script setup lang="ts">
import { ElMessage, type FormInstance, type FormRules } from 'element-plus';
import { Refresh, Search } from '@element-plus/icons-vue';
import { getIconSvg } from '~/composables/lucideIcons';
import { useNotify } from '~/composables/useNotify';
import { formatDateTime } from '~/utils/date';

definePageMeta({ layout: 'admin', middleware: 'auth' });

const notify = useNotify();

interface NavItem {
  id: number;
  label: string;
  link: string | null;
  link_type: string;
  project_id: number | null;
  page_id: number | null;
  parent_id: number | null;
  sort_order: number;
  status: boolean;
  display_position?: string;
  children: NavItem[];
}

interface ProjectBrief {
  id: number;
  name: string;
  slug: string;
}

interface PageBrief {
  id: number;
  title: string;
  slug: string;
}

const { user } = useAuth();
const isViewer = computed(() => user.value?.role === 'viewer');

const treeData = ref<NavItem[]>([]);
const projects = ref<ProjectBrief[]>([]);
const pages = ref<PageBrief[]>([]);
const loading = ref(false);
const saving = ref(false);
const searchQuery = ref('');
const typeFilter = ref('');
const positionFilter = ref('');

const dialogVisible = ref(false);
const editingId = ref<number | null>(null);
const formRef = ref<FormInstance>();

const defaultForm = (): {
  label: string;
  link: string;
  link_type: string;
  project_id: number | null;
  page_id: number | null;
  parent_id: number | null;
  sort_order: number;
  status: boolean;
  display_position: string;
} => ({
  label: '',
  link: '',
  link_type: 'custom',
  project_id: null,
  page_id: null,
  parent_id: ROOT_PARENT_ID,
  sort_order: 0,
  status: true,
  display_position: 'header',
});

const form = reactive(defaultForm());

const linkTypeLabel = (t: string) => {
  const map: Record<string, string> = { project: '项目', page: '页面', custom: '自定义' };
  return map[t] || t;
};

const linkTypeTag = (t: string): 'info' | 'primary' | 'success' | 'warning' | 'danger' => {
  const map: Record<string, 'success' | 'warning' | 'info'> = { project: 'success', page: 'warning', custom: 'info' };
  return map[t] || 'info';
};

const selectedProjectSlug = computed(() => {
  if (!form.project_id) return '';
  const p = projects.value.find((p) => p.id === form.project_id);
  return p ? p.slug : '';
});

const selectedPageSlug = computed(() => {
  if (!form.page_id) return '';
  const p = pages.value.find((p) => p.id === form.page_id);
  return p ? p.slug : '';
});

const rules = computed<FormRules>(() => ({
  label: [{ required: true, message: '请输入名称', trigger: 'blur' }],
  project_id: form.link_type === 'project'
    ? [{ required: true, message: '请选择项目', trigger: 'change' }]
    : [],
  page_id: form.link_type === 'page'
    ? [{ required: true, message: '请选择页面', trigger: 'change' }]
    : [],
  link: form.link_type === 'custom'
    ? [
        {
          validator: (_rule: any, value: string, callback: any) => {
            if (!value) {
              // Link is optional for custom items
              callback();
            } else if (!value.startsWith('/')) {
              callback(new Error('链接必须以 / 开头'));
            } else if (value.includes('://')) {
              callback(new Error('仅支持内部链接'));
            } else {
              callback();
            }
          },
          trigger: 'blur',
        },
      ]
    : [],
}));

const ROOT_PARENT_ID = 0;

const parentOptions = computed(() => {
  const filterOut = (items: NavItem[]): NavItem[] =>
    items
      .filter((item) => item.id !== editingId.value)
      .map((item) => ({ ...item, children: filterOut(item.children || []) }));
  return [
    { id: ROOT_PARENT_ID, label: '根节点（顶级）', link: null, link_type: 'custom', project_id: null, page_id: null, parent_id: null, sort_order: 0, status: true, children: [] as NavItem[] },
    ...filterOut(treeData.value),
  ];
});

const filterTree = (items: NavItem[]): NavItem[] => {
  if (!items) return [];
  return items.reduce((acc: NavItem[], item) => {
    const matchName = !searchQuery.value || item.label.toLowerCase().includes(searchQuery.value.toLowerCase());
    const matchType = !typeFilter.value || item.link_type === typeFilter.value;
    const matchPos = !positionFilter.value || item.display_position === positionFilter.value || item.display_position === 'both';
    const matches = matchName && matchType && matchPos;
    const filteredChildren = filterTree(item.children || []);
    if (matches || filteredChildren.length > 0) {
      acc.push({ ...item, children: filteredChildren });
    }
    return acc;
  }, []);
};

const filteredTreeData = computed(() => {
  if (!searchQuery.value && !typeFilter.value && !positionFilter.value) {
    return treeData.value;
  }
  return filterTree(treeData.value);
});

const loadTree = async () => {
  loading.value = true;
  try {
    const api = useApi();
    const data = await api<NavItem[]>('/admin/navigation');
    treeData.value = (data || []) as NavItem[];
  } catch {
    treeData.value = [];
    ElMessage.error('加载导航数据失败');
  } finally {
    loading.value = false;
  }
};

const loadOptions = async () => {
  try {
    const api = useApi();
    const [projRes, pageRes] = await Promise.all([
      api<ProjectBrief[]>('/admin/projects'),
      api<PageBrief[]>('/admin/pages'),
    ]);
    projects.value = projRes || [];
    pages.value = pageRes || [];
  } catch {
    // non-critical; dropdowns will be empty
  }
};

const onLinkTypeChange = () => {
  form.project_id = null;
  form.page_id = null;
  form.link = '';
};

const onProjectChange = () => {
  // auto-fill link preview is reactive via selectedProjectSlug
};

const onPageChange = () => {
  // auto-fill link preview is reactive via selectedPageSlug
};

const openCreate = (parentId?: number) => {
  editingId.value = null;
  Object.assign(form, defaultForm());
  if (parentId !== undefined) {
    form.parent_id = parentId;
  }
  dialogVisible.value = true;
};

const openEdit = (row: NavItem) => {
  editingId.value = row.id;
  Object.assign(form, {
    label: row.label,
    link: row.link || '',
    link_type: row.link_type || 'custom',
    project_id: row.project_id,
    page_id: row.page_id,
    parent_id: row.parent_id ?? ROOT_PARENT_ID,
    sort_order: row.sort_order,
    status: row.status,
    display_position: row.display_position || 'header',
  });
  dialogVisible.value = true;
};

const handleSave = async () => {
  const valid = await formRef.value?.validate().catch(() => false);
  if (!valid) return;

  saving.value = true;
  try {
    const api = useApi();
    const body: Record<string, any> = {
      label: form.label,
      link_type: form.link_type,
      parent_id: form.parent_id === ROOT_PARENT_ID ? null : form.parent_id,
      sort_order: form.sort_order,
      status: form.status,
      display_position: form.display_position,
    };

    if (form.link_type === 'project') {
      body.project_id = form.project_id;
      body.link = null;
    } else if (form.link_type === 'page') {
      body.page_id = form.page_id;
      body.link = null;
    } else {
      body.link = form.link || null;
      body.project_id = null;
      body.page_id = null;
    }

    if (editingId.value) {
      await api(`/admin/navigation/${editingId.value}`, { method: 'PUT', body });
      notify.success('更新成功');
    } else {
      await api('/admin/navigation', { method: 'POST', body });
      notify.success('添加成功');
    }
    dialogVisible.value = false;
    loadTree();
  } catch (err: any) {
    notify.error(err, '操作失败');
  } finally {
    saving.value = false;
  }
};

const handleDelete = async (id: number) => {
  try {
    const api = useApi();
    await api(`/admin/navigation/${id}`, { method: 'DELETE' });
    notify.success('已删除');
    loadTree();
  } catch (err: any) {
    notify.error(err, '删除失败');
  }
};

onMounted(() => {
  loadTree();
  loadOptions();
});
</script>

<style scoped>
.nav-tree-header {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 10px 16px 10px 40px;
  border-bottom: 2px solid var(--border-color);
  font-size: 12px;
  font-weight: 600;
  color: var(--color-text-muted);
  background: var(--bg-gray-50, #fafafa);
}

.nav-th-name { min-width: 160px; }
.nav-th-sort { width: 100px; text-align: center; }
.nav-th-link { flex: 1; }
.nav-th-type { width: 70px; text-align: center; }
.nav-th-pos { width: 90px; text-align: center; }
.nav-th-time { width: 130px; text-align: center; }
.nav-th-actions { width: 110px; text-align: right; }

.nav-tree {
  background: transparent;
}

.nav-tree :deep(.el-tree-node__content) {
  height: auto;
  padding: 0 16px;
}

.tree-node {
  display: flex;
  align-items: center;
  gap: 12px;
  flex: 1;
  padding: 8px 0;
  width: 100%;
}

.tree-node-info {
  display: flex;
  align-items: center;
  gap: 8px;
  min-width: 160px;
}

.row-title {
  font-size: 14px;
  font-weight: 500;
  color: var(--color-text);
}

.row-meta {
  font-size: 12px;
  color: var(--color-text-muted);
  flex: 1;
}

.row-sort {
  width: 100px;
  font-size: 13px;
  color: var(--color-text-muted);
  text-align: center;
}

.row-type {
  width: 70px;
  display: flex;
  justify-content: center;
}

.row-pos {
  width: 90px;
  display: flex;
  justify-content: center;
}

.row-time {
  width: 130px;
  font-size: 12px;
  color: var(--color-text-muted);
  text-align: center;
}

.tree-node-actions {
  display: flex;
  gap: 2px;
  width: 110px;
  justify-content: flex-end;
  align-items: center;
}

.link-preview {
  margin-top: 6px;
  font-size: 13px;
  color: #606266;
}

.link-preview code {
  background: #f5f7fa;
  padding: 2px 6px;
  border-radius: 3px;
  font-family: monospace;
}
</style>
