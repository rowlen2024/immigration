<template>
  <div>
    <div class="admin-page-header">
      <h2 class="admin-page-title">导航管理</h2>
      <div v-if="!isViewer" style="display:flex;align-items:center;gap:8px;">
        <el-button :icon="Refresh" circle @click="loadTree" :loading="loading" />
        <el-button type="primary" @click="openCreate()">新建导航</el-button>
      </div>
    </div>

    <div class="admin-table-wrap">
      <el-tree
        v-loading="loading"
        :data="treeData"
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
            <span class="row-meta">{{ data.link }}</span>
            <el-tag
              :type="linkTypeTag(data.link_type)"
              size="small"
              effect="plain"
            >
              {{ linkTypeLabel(data.link_type) }}
            </el-tag>
            <el-tag
              v-if="data.display_position && data.display_position !== 'header'"
              :type="data.display_position === 'both' ? 'success' : 'warning'"
              size="small"
              effect="plain"
            >
              {{ data.display_position === 'both' ? '头部+页脚' : '仅页脚' }}
            </el-tag>
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
              :label="`${p.title} (/${p.slug})`"
              :value="p.id"
            />
          </el-select>
          <div v-if="form.page_id" class="link-preview">
            链接预览：<code>/{{ selectedPageSlug }}</code>
          </div>
        </el-form-item>

        <el-form-item v-if="form.link_type === 'custom'" label="链接" prop="link">
          <el-input v-model="form.link" placeholder="/path" />
        </el-form-item>

        <el-form-item label="父级导航" prop="parent_id">
          <el-tree-select
            v-model="form.parent_id"
            :data="parentOptions"
            :props="{ children: 'children', label: 'label', value: 'id' }"
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
import { Refresh } from '@element-plus/icons-vue';
import { getIconSvg } from '~/composables/lucideIcons';

definePageMeta({ layout: 'admin', middleware: 'auth' });

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
  parent_id: null,
  sort_order: 0,
  status: true,
  display_position: 'header',
});

const form = reactive(defaultForm());

const linkTypeLabel = (t: string) => {
  const map: Record<string, string> = { project: '项目', page: '页面', custom: '自定义' };
  return map[t] || t;
};

const linkTypeTag = (t: string) => {
  const map: Record<string, string> = { project: 'success', page: 'warning', custom: 'info' };
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
        { required: true, message: '请输入链接', trigger: 'blur' },
        {
          validator: (_rule: any, value: string, callback: any) => {
            if (!value || !value.startsWith('/')) {
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

const parentOptions = computed(() => {
  const filterOut = (items: NavItem[]): NavItem[] =>
    items
      .filter((item) => item.id !== editingId.value)
      .map((item) => ({ ...item, children: filterOut(item.children || []) }));
  return filterOut(treeData.value);
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
    const [projData, pageData] = await Promise.all([
      api<any[]>('/admin/projects?all=true'),
      api<any[]>('/admin/pages?all=true'),
    ]);
    projects.value = (projData || []) as ProjectBrief[];
    pages.value = (pageData || []) as PageBrief[];
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
    parent_id: row.parent_id,
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
      parent_id: form.parent_id,
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
      body.link = form.link;
      body.project_id = null;
      body.page_id = null;
    }

    if (editingId.value) {
      await api(`/admin/navigation/${editingId.value}`, { method: 'PUT', body });
    } else {
      await api('/admin/navigation', { method: 'POST', body });
    }
    dialogVisible.value = false;
    loadTree();
  } catch (err: any) {
    ElMessage.error(err?.message || '操作失败');
  } finally {
    saving.value = false;
  }
};

const handleDelete = async (id: number) => {
  try {
    const api = useApi();
    await api(`/admin/navigation/${id}`, { method: 'DELETE' });
    loadTree();
  } catch (err: any) {
    ElMessage.error(err?.data?.message || err?.message || '删除失败');
  }
};

onMounted(() => {
  loadTree();
  loadOptions();
});
</script>

<style scoped>
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

.tree-node-actions {
  display: flex;
  gap: 2px;
  margin-left: auto;
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
