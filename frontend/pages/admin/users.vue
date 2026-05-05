<template>
  <div>
    <div class="admin-page-header">
      <h2 class="admin-page-title">用户管理</h2>
      <div style="display:flex;align-items:center;gap:8px;">
        <el-button :icon="Refresh" circle @click="loadList" :loading="loading" />
        <el-button v-if="isAdmin" type="primary" @click="openCreate">新建用户</el-button>
      </div>
    </div>

    <div class="admin-table-wrap">
      <el-table :data="list" v-loading="loading">
        <el-table-column prop="username" label="用户名" min-width="140">
          <template #default="{ row }">
            <div class="row-title">{{ row.username }}</div>
          </template>
        </el-table-column>
        <el-table-column prop="display_name" label="显示名称" min-width="140" />
        <el-table-column prop="role" label="角色" width="110">
          <template #default="{ row }">
            <span :class="['status-pill', row.role === 'admin' ? 'danger' : row.role === 'editor' ? 'warning' : 'info']">
              {{ row.role === 'admin' ? '管理员' : row.role === 'editor' ? '编辑者' : '只读用户' }}
            </span>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="90">
          <template #default="{ row }">
            <span :class="['status-pill', row.status === 1 ? 'published' : 'warning']">
              {{ row.status === 1 ? '正常' : '禁用' }}
            </span>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <div class="table-actions" v-if="isAdmin">
              <button class="action-btn" @click="openEditRole(row)">编辑角色</button>
              <button class="action-btn" :class="{ danger: row.status === 1 }" @click="handleToggleStatus(row)">
                {{ row.status === 1 ? '禁用' : '启用' }}
              </button>
            </div>
          </template>
        </el-table-column>
      </el-table>
    </div>

    <div v-if="!loading && list.length === 0" class="admin-empty-state">
      <div class="empty-icon" v-html="getIconSvg('shield', 48)"></div>
      <div class="empty-title">暂无用户</div>
    </div>

    <div class="admin-pagination-wrap" v-if="total > pageSize">
      <el-pagination v-model:current-page="page" :page-size="pageSize" :total="total" layout="total, prev, pager, next" @current-change="loadList" />
    </div>

    <el-drawer v-model="drawerVisible" title="新建用户" size="500px" destroy-on-close>
      <el-form ref="formRef" :model="form" :rules="rules" label-position="top">
        <el-form-item label="用户名" prop="username">
          <el-input v-model="form.username" />
        </el-form-item>
        <el-form-item label="密码" prop="password">
          <el-input v-model="form.password" type="password" show-password />
        </el-form-item>
        <el-form-item label="显示名称" prop="display_name">
          <el-input v-model="form.display_name" />
        </el-form-item>
        <el-form-item label="角色" prop="role">
          <el-select v-model="form.role">
            <el-option label="管理员" value="admin" />
            <el-option label="编辑者" value="editor" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="drawerVisible = false">取消</el-button>
        <el-button type="primary" :loading="saving" @click="handleCreateUser">创建</el-button>
      </template>
    </el-drawer>

    <el-dialog v-model="roleDialogVisible" title="编辑角色" width="400px" destroy-on-close>
      <el-form label-width="80px" label-position="top">
        <el-form-item label="角色">
          <el-select v-model="editRoleValue">
            <el-option label="管理员" value="admin" />
            <el-option label="编辑者" value="editor" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="roleDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="saving" @click="handleSaveRole">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ElMessage, type FormInstance, type FormRules } from 'element-plus';
import { Refresh } from '@element-plus/icons-vue';
import { getIconSvg } from '~/composables/lucideIcons';

definePageMeta({ layout: 'admin', middleware: 'auth' });

interface User {
  id: string;
  username: string;
  display_name: string;
  role: string;
  status: number;
}

const list = ref<User[]>([]);
const loading = ref(false);
const saving = ref(false);
const page = ref(1);
const pageSize = ref(10);
const total = ref(0);

const isAdmin = ref(false);

const drawerVisible = ref(false);
const formRef = ref<FormInstance>();

const form = reactive({
  username: '',
  password: '',
  display_name: '',
  role: 'editor',
});

const rules: FormRules = {
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  password: [{ required: true, message: '请输入密码', trigger: 'blur', min: 6 }],
  display_name: [{ required: true, message: '请输入显示名称', trigger: 'blur' }],
  role: [{ required: true, message: '请选择角色', trigger: 'change' }],
};

const roleDialogVisible = ref(false);
const editingUserId = ref<string | null>(null);
const editRoleValue = ref('editor');

const checkAdmin = () => {
  const { user } = useAuth();
  isAdmin.value = (user.value as any)?.role === 'admin';
};

const loadList = async () => {
  loading.value = true;
  try {
    const api = useApi();
    const data = await api<{ items: User[]; total: number }>(
      `/admin/users?page=${page.value}&per_page=${pageSize.value}`
    );
    list.value = data.items ?? [];
    total.value = data.total ?? 0;
  } catch {
    list.value = [];
    ElMessage.error('加载用户列表失败');
  } finally {
    loading.value = false;
  }
};

const openCreate = () => {
  form.username = '';
  form.password = '';
  form.display_name = '';
  form.role = 'editor';
  drawerVisible.value = true;
};

const handleCreateUser = async () => {
  const valid = await formRef.value?.validate().catch(() => false);
  if (!valid) return;

  saving.value = true;
  try {
    const api = useApi();
    await api('/admin/users', { method: 'POST', body: form });
    drawerVisible.value = false;
    loadList();
  } catch {
    ElMessage.error('创建用户失败');
  } finally {
    saving.value = false;
  }
};

const openEditRole = (row: User) => {
  editingUserId.value = row.id;
  editRoleValue.value = row.role;
  roleDialogVisible.value = true;
};

const handleSaveRole = async () => {
  if (!editingUserId.value) return;

  saving.value = true;
  try {
    const api = useApi();
    await api(`/admin/users/${editingUserId.value}`, {
      method: 'PUT',
      body: { role: editRoleValue.value },
    });
    roleDialogVisible.value = false;
    loadList();
  } catch {
    ElMessage.error('更新角色失败');
  } finally {
    saving.value = false;
  }
};

const handleToggleStatus = async (row: User) => {
  const newStatus = row.status === 1 ? 0 : 1;
  try {
    const api = useApi();
    await api(`/admin/users/${row.id}`, { method: 'PUT', body: { status: newStatus } });
    loadList();
  } catch {
    ElMessage.error('更新用户状态失败');
  }
};

onMounted(() => {
  checkAdmin();
  loadList();
});
</script>

<style scoped>
.row-title {
  font-size: 14px;
  font-weight: 500;
  color: var(--color-text);
}
</style>
