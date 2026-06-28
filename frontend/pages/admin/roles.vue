<template>
  <div>
    <AdminPageHeader title="角色权限">
      <template #actions>
        <el-button :icon="Refresh" circle @click="loadRoles" :loading="loading" />
        <el-button v-if="canWriteRoles" type="primary" @click="openCreate">新建角色</el-button>
      </template>
    </AdminPageHeader>

    <AdminTableShell v-if="loading || roles.length > 0" :loading="loading">
      <el-table :data="roles">
        <el-table-column prop="name" label="角色名称" min-width="140" />
        <el-table-column prop="code" label="角色编码" min-width="140" />
        <el-table-column prop="description" label="描述" min-width="180">
          <template #default="{ row }">{{ row.description || '-' }}</template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="90">
          <template #default="{ row }">
            <span :class="['status-pill', row.status === 1 ? 'published' : 'warning']">
              {{ row.status === 1 ? '启用' : '停用' }}
            </span>
          </template>
        </el-table-column>
        <el-table-column prop="is_system" label="系统角色" width="100">
          <template #default="{ row }">{{ row.is_system ? '是' : '否' }}</template>
        </el-table-column>
        <el-table-column label="操作" width="180" fixed="right">
          <template #default="{ row }">
            <AdminRowActions>
              <button v-if="canWriteRoles" class="action-btn" type="button" title="编辑" aria-label="编辑" @click="openEdit(row)" v-html="getIconSvg('pencil', 16)"></button>
              <el-popconfirm
                v-if="canWriteRoles && !row.is_system"
                title="确定删除该角色？"
                confirm-button-text="删除"
                cancel-button-text="取消"
                @confirm="handleDelete(row.id)"
              >
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
      v-else
      icon="shield-check"
      title="暂无角色"
      description="点击上方按钮创建第一个角色"
      :action-label="canWriteRoles ? '新建角色' : undefined"
      @action="openCreate"
    />

    <el-drawer v-model="drawerVisible" :title="editingId ? '编辑角色' : '新建角色'" size="760px" destroy-on-close>
      <el-form ref="formRef" :model="form" :rules="rules" label-position="top">
        <el-form-item label="角色编码" prop="code">
          <el-input v-model="form.code" :disabled="!!editingId" placeholder="例如 content_editor" />
        </el-form-item>
        <el-form-item label="角色名称" prop="name">
          <el-input v-model="form.name" />
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="form.description" type="textarea" :rows="2" />
        </el-form-item>
        <el-form-item label="状态">
          <el-switch v-model="form.enabled" active-text="启用" inactive-text="停用" />
        </el-form-item>

        <el-divider content-position="left">权限配置</el-divider>
        <div class="permission-groups">
          <section v-for="group in permissionGroups" :key="group.module" class="permission-group">
            <h3>{{ moduleLabel(group.module) }}</h3>
            <el-checkbox-group v-model="selectedPermissions" :disabled="isEditingAdminRole">
              <el-checkbox v-for="permission in group.items" :key="permission.code" :label="permission.code">
                {{ permission.name }}
                <span class="permission-code">{{ permission.code }}</span>
              </el-checkbox>
            </el-checkbox-group>
          </section>
        </div>
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
import { ElMessage, type FormInstance, type FormRules } from 'element-plus';
import { Refresh } from '@element-plus/icons-vue';
import { getIconSvg } from '~/composables/lucideIcons';

definePageMeta({ layout: 'admin', middleware: 'auth' });

interface Role {
  id: number;
  code: string;
  name: string;
  description: string;
  status: number;
  is_system: boolean;
  permission_codes?: string[];
}

interface Permission {
  code: string;
  name: string;
  module: string;
}

const { hasPermission } = usePermissions();
const canWriteRoles = computed(() => hasPermission('roles:write'));

const roles = ref<Role[]>([]);
const permissions = ref<Permission[]>([]);
const selectedPermissions = ref<string[]>([]);
const loading = ref(false);
const saving = ref(false);
const drawerVisible = ref(false);
const editingId = ref<number | null>(null);
const editingRoleCode = ref('');
const formRef = ref<FormInstance>();

const form = reactive({
  code: '',
  name: '',
  description: '',
  enabled: true,
});

const rules: FormRules = {
  code: [{ required: true, message: '请输入角色编码', trigger: 'blur' }],
  name: [{ required: true, message: '请输入角色名称', trigger: 'blur' }],
};

const permissionGroups = computed(() => {
  const groups = new Map<string, Permission[]>();
  for (const permission of permissions.value) {
    if (!groups.has(permission.module)) groups.set(permission.module, []);
    groups.get(permission.module)!.push(permission);
  }
  return Array.from(groups.entries()).map(([module, items]) => ({ module, items }));
});

const isEditingAdminRole = computed(() => editingRoleCode.value === 'admin');

const moduleLabel = (module: string) => {
  const labels: Record<string, string> = {
    dashboard: '控制台',
    projects: '项目',
    homepage: '首页配置',
    navigation: '导航',
    pages: '页面',
    media: '媒体库',
    faqs: 'FAQ',
    cases: '案例',
    lawyers: '律师团队',
    testimonials: '客户评价',
    leads: '咨询',
    settings: '网站设置',
    users: '用户',
    roles: '角色权限',
  };
  return labels[module] || module;
};

const loadPermissions = async () => {
  const api = useApi();
  permissions.value = await api<Permission[]>('/admin/permissions') ?? [];
};

const loadRoles = async () => {
  loading.value = true;
  try {
    const api = useApi();
    roles.value = await api<Role[]>('/admin/roles') ?? [];
  } catch {
    roles.value = [];
    ElMessage.error('加载角色失败');
  } finally {
    loading.value = false;
  }
};

const openCreate = () => {
  editingId.value = null;
  editingRoleCode.value = '';
  Object.assign(form, { code: '', name: '', description: '', enabled: true });
  selectedPermissions.value = [];
  drawerVisible.value = true;
};

const openEdit = async (row: Role) => {
  editingId.value = row.id;
  editingRoleCode.value = row.code;
  Object.assign(form, {
    code: row.code,
    name: row.name,
    description: row.description,
    enabled: row.status === 1,
  });
  try {
    const api = useApi();
    const detail = await api<Role>(`/admin/roles/${row.id}`);
    selectedPermissions.value = detail?.permission_codes ?? [];
  } catch {
    selectedPermissions.value = [];
    ElMessage.error('加载角色权限失败');
  }
  drawerVisible.value = true;
};

const handleSave = async () => {
  const valid = await formRef.value?.validate().catch(() => false);
  if (!valid) return;
  saving.value = true;
  try {
    const api = useApi();
    const body: {
      code: string;
      name: string;
      description: string;
      status: number;
      permission_codes?: string[];
    } = {
      code: form.code,
      name: form.name,
      description: form.description,
      status: form.enabled ? 1 : 0,
    };
    if (!isEditingAdminRole.value) {
      body.permission_codes = selectedPermissions.value;
    }
    if (editingId.value) {
      await api(`/admin/roles/${editingId.value}`, { method: 'PUT', body });
      ElMessage.success('角色已更新');
    } else {
      await api('/admin/roles', { method: 'POST', body });
      ElMessage.success('角色已创建');
    }
    drawerVisible.value = false;
    loadRoles();
  } catch (e: any) {
    ElMessage.error(e?.data?.message || '保存角色失败');
  } finally {
    saving.value = false;
  }
};

const handleDelete = async (id: number) => {
  try {
    const api = useApi();
    await api(`/admin/roles/${id}`, { method: 'DELETE' });
    ElMessage.success('角色已删除');
    loadRoles();
  } catch (e: any) {
    ElMessage.error(e?.data?.message || '删除角色失败');
  }
};

onMounted(async () => {
  await loadPermissions();
  await loadRoles();
});
</script>

<style scoped>
.permission-groups {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(260px, 1fr));
  gap: 16px;
}

.permission-group {
  border: 1px solid var(--color-border);
  border-radius: var(--radius-sm);
  padding: 12px;
}

.permission-group h3 {
  margin: 0 0 10px;
  font-size: 14px;
}

.permission-group :deep(.el-checkbox-group) {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.permission-code {
  margin-left: 6px;
  color: var(--color-text-muted);
  font-size: 12px;
}
</style>
