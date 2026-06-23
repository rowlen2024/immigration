<template>
  <div>
    <div class="admin-page-header">
      <h2 class="admin-page-title">用户管理</h2>
      <div class="header-actions">
        <el-button :icon="Refresh" circle @click="loadList" :loading="loading" />
        <el-button v-if="canWriteUsers" type="primary" @click="openCreate">新建用户</el-button>
      </div>
    </div>

    <div class="admin-table-wrap">
      <el-table :data="list" v-loading="loading">
        <el-table-column prop="username" label="用户名" min-width="140" />
        <el-table-column prop="display_name" label="显示名称" min-width="140">
          <template #default="{ row }">{{ row.display_name || '-' }}</template>
        </el-table-column>
        <el-table-column prop="role" label="角色" width="140">
          <template #default="{ row }">{{ roleLabel(row.role) }}</template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="90">
          <template #default="{ row }">
            <span :class="['status-pill', row.status === 1 ? 'published' : 'warning']">
              {{ row.status === 1 ? '正常' : '禁用' }}
            </span>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="160">
          <template #default="{ row }">{{ formatDateTime(row.created_at) }}</template>
        </el-table-column>
        <el-table-column label="操作" width="220" fixed="right">
          <template #default="{ row }">
            <div class="table-actions" v-if="canWriteUsers">
              <button class="action-btn" @click="openEdit(row)">编辑</button>
              <button class="action-btn" :class="{ danger: row.status === 1 }" @click="handleToggleStatus(row)">
                {{ row.status === 1 ? '禁用' : '启用' }}
              </button>
            </div>
          </template>
        </el-table-column>
      </el-table>
    </div>

    <div class="admin-pagination-wrap" v-if="total > pageSize">
      <el-pagination v-model:current-page="page" :page-size="pageSize" :total="total" layout="total, prev, pager, next" @current-change="loadList" />
    </div>

    <el-drawer v-model="drawerVisible" :title="editingUserId ? '编辑用户' : '新建用户'" size="760px" destroy-on-close>
      <el-form ref="formRef" :model="form" :rules="rules" label-position="top">
        <el-form-item label="用户名" prop="username">
          <el-input v-model="form.username" :disabled="!!editingUserId" />
        </el-form-item>
        <el-form-item :label="editingUserId ? '新密码' : '密码'" prop="password">
          <el-input v-model="form.password" type="password" show-password placeholder="编辑时留空表示不修改密码" />
        </el-form-item>
        <el-form-item label="显示名称" prop="display_name">
          <el-input v-model="form.display_name" />
        </el-form-item>
        <el-form-item label="角色" prop="role">
          <el-select v-model="form.role" style="width: 100%">
            <el-option v-for="role in roles" :key="role.code" :label="role.name" :value="role.code" />
          </el-select>
        </el-form-item>

        <el-divider content-position="left">个人权限微调</el-divider>
        <div class="permission-groups">
          <section v-for="group in permissionGroups" :key="group.module" class="permission-group">
            <h3>{{ moduleLabel(group.module) }}</h3>
            <div class="permission-options">
              <div v-for="permission in group.items" :key="permission.code" class="permission-option">
                <div class="permission-info">
                  <span class="permission-name">{{ permission.name }}</span>
                  <span class="permission-code">{{ permission.code }}</span>
                </div>
                <el-radio-group
                  :model-value="overrideMap[permission.code] || ''"
                  size="small"
                  @change="(value) => setOverride(permission.code, value)"
                >
                  <el-radio-button value="">继承角色</el-radio-button>
                  <el-radio-button value="allow">额外允许</el-radio-button>
                  <el-radio-button value="deny">明确拒绝</el-radio-button>
                </el-radio-group>
              </div>
            </div>
          </section>
        </div>
      </el-form>
      <template #footer>
        <el-button @click="drawerVisible = false">取消</el-button>
        <el-button type="primary" :loading="saving" @click="handleSave">保存</el-button>
      </template>
    </el-drawer>
  </div>
</template>

<script setup lang="ts">
import { ElMessage, type FormInstance, type FormRules } from 'element-plus';
import { Refresh } from '@element-plus/icons-vue';
import { formatDateTime } from '~/utils/date';

definePageMeta({ layout: 'admin', middleware: 'auth' });

interface User {
  id: string;
  username: string;
  display_name: string;
  role: string;
  status: number;
  created_at: string;
}

interface Role {
  id: number;
  code: string;
  name: string;
}

interface Permission {
  code: string;
  name: string;
  module: string;
}

const { hasPermission } = usePermissions();
const canWriteUsers = computed(() => hasPermission('users:write'));

const list = ref<User[]>([]);
const roles = ref<Role[]>([]);
const permissions = ref<Permission[]>([]);
const overrideMap = reactive<Record<string, string>>({});
const loading = ref(false);
const saving = ref(false);
const page = ref(1);
const pageSize = ref(10);
const total = ref(0);
const drawerVisible = ref(false);
const editingUserId = ref<string | null>(null);
const formRef = ref<FormInstance>();

const form = reactive({
  username: '',
  password: '',
  display_name: '',
  role: 'viewer',
});

const rules = computed<FormRules>(() => ({
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  password: editingUserId.value ? [] : [{ required: true, message: '请输入密码', trigger: 'blur', min: 6 }],
  display_name: [{ required: true, message: '请输入显示名称', trigger: 'blur' }],
  role: [{ required: true, message: '请选择角色', trigger: 'change' }],
}));

const roleLabel = (code: string) => roles.value.find((role) => role.code === code)?.name || code;

const permissionGroups = computed(() => {
  const groups = new Map<string, Permission[]>();
  for (const permission of permissions.value) {
    if (!groups.has(permission.module)) groups.set(permission.module, []);
    groups.get(permission.module)!.push(permission);
  }
  return Array.from(groups.entries()).map(([module, items]) => ({ module, items }));
});

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

const resetOverrides = () => {
  for (const key of Object.keys(overrideMap)) {
    delete overrideMap[key];
  }
};

const setOverride = (code: string, value: string | number | boolean | undefined) => {
  if (value === 'allow' || value === 'deny') {
    overrideMap[code] = value;
    return;
  }
  delete overrideMap[code];
};

const loadMeta = async () => {
  const api = useApi();
  const [roleData, permissionData] = await Promise.all([
    api<Role[]>('/admin/roles'),
    api<Permission[]>('/admin/permissions'),
  ]);
  roles.value = roleData ?? [];
  permissions.value = permissionData ?? [];
};

const loadList = async () => {
  loading.value = true;
  try {
    const api = useApi();
    const data = await api<{ items: User[]; total: number }>(`/admin/users?page=${page.value}&per_page=${pageSize.value}`);
    list.value = data.items ?? [];
    total.value = data.total ?? 0;
  } catch {
    list.value = [];
    ElMessage.error('加载用户列表失败');
  } finally {
    loading.value = false;
  }
};

const openCreate = async () => {
  editingUserId.value = null;
  Object.assign(form, { username: '', password: '', display_name: '', role: roles.value[0]?.code || 'viewer' });
  resetOverrides();
  drawerVisible.value = true;
};

const openEdit = async (row: User) => {
  editingUserId.value = row.id;
  Object.assign(form, { username: row.username, password: '', display_name: row.display_name, role: row.role });
  resetOverrides();
  try {
    const api = useApi();
    const data = await api<any>(`/admin/users/${row.id}`);
    const overrides = data?.permission_overrides ?? [];
    for (const item of overrides) {
      const code = item.permission?.code;
      if (code) overrideMap[code] = item.effect;
    }
  } catch {
    ElMessage.error('加载用户权限失败');
  }
  drawerVisible.value = true;
};

const buildOverrides = () => Object.entries(overrideMap)
  .filter(([, effect]) => effect === 'allow' || effect === 'deny')
  .map(([permission_code, effect]) => ({ permission_code, effect }));

const handleSave = async () => {
  const valid = await formRef.value?.validate().catch(() => false);
  if (!valid) return;

  saving.value = true;
  try {
    const api = useApi();
    const body: Record<string, any> = {
      display_name: form.display_name,
      role: form.role,
      permission_overrides: buildOverrides(),
    };
    if (form.password) body.password = form.password;
    if (editingUserId.value) {
      await api(`/admin/users/${editingUserId.value}`, { method: 'PUT', body });
      ElMessage.success('用户已更新');
    } else {
      await api('/admin/users', {
        method: 'POST',
        body: { ...body, username: form.username, password: form.password },
      });
      ElMessage.success('用户已创建');
    }
    drawerVisible.value = false;
    loadList();
  } catch (e: any) {
    ElMessage.error(e?.data?.message || '保存用户失败');
  } finally {
    saving.value = false;
  }
};

const handleToggleStatus = async (row: User) => {
  try {
    const api = useApi();
    await api(`/admin/users/${row.id}`, { method: 'PUT', body: { status: row.status === 1 ? 0 : 1 } });
    ElMessage.success('状态已更新');
    loadList();
  } catch (e: any) {
    ElMessage.error(e?.data?.message || '更新用户状态失败');
  }
};

onMounted(async () => {
  await loadMeta();
  await loadList();
});
</script>

<style scoped>
.header-actions {
  display: flex;
  align-items: center;
  gap: 8px;
}

.permission-groups {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
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

.permission-options {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.permission-option {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}

.permission-info {
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.permission-name {
  font-size: 14px;
  color: var(--color-text);
}

.permission-code {
  color: var(--color-text-muted);
  font-size: 12px;
}

.permission-option :deep(.el-radio-group) {
  flex-shrink: 0;
}

@media (max-width: 767px) {
  .permission-option {
    align-items: flex-start;
    flex-direction: column;
  }
}
</style>
