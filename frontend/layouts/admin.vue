<template>
  <div class="admin-layout">
    <div
      class="sidebar-overlay"
      :class="{ 'overlay-visible': mobileOpen }"
      @click="mobileOpen = false"
    ></div>

    <aside
      class="admin-sidebar"
      :class="{ collapsed: sidebarCollapsed, 'mobile-open': mobileOpen }"
    >
      <!-- Logo -->
      <div class="sidebar-header">
        <NuxtLink to="/admin" class="sidebar-logo" @click="closeMobile">
          <div class="logo-icon">
            <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <path d="m3 9 9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z"/>
              <polyline points="9 22 9 12 15 12 15 22"/>
            </svg>
          </div>
          <span class="logo-text">北极星移民管理后台</span>
        </NuxtLink>
      </div>

      <!-- Navigation -->
      <nav class="sidebar-nav">
        <template v-for="group in navGroups" :key="group.label">
          <!-- Standalone item -->
          <NuxtLink
            v-if="!group.children"
            :to="group.to!"
            class="nav-item"
            :class="{ active: isActive(group.to!) }"
            @click="closeMobile"
          >
            <span class="nav-icon" v-html="group.icon"></span>
            <span class="nav-label">{{ group.label }}</span>
          </NuxtLink>

          <!-- Group with children -->
          <div v-else class="nav-group">
            <button
              class="nav-item nav-group-title"
              :class="{ active: isGroupActive(group) }"
              @click="toggleGroup(group.label)"
            >
              <span class="nav-icon" v-html="group.icon"></span>
              <span class="nav-label">{{ group.label }}</span>
              <span class="nav-chevron" :class="{ open: expandedGroups.has(group.label) }">
                <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <polyline points="6 9 12 15 18 9"/>
                </svg>
              </span>
            </button>
            <div v-show="expandedGroups.has(group.label)" class="nav-sub-items">
              <NuxtLink
                v-for="child in group.children"
                :key="child.to"
                :to="child.to"
                class="nav-item nav-sub-item"
                :class="{ active: isActive(child.to) }"
                @click="closeMobile"
              >
                <span class="nav-label">{{ child.label }}</span>
              </NuxtLink>
            </div>
          </div>
        </template>
      </nav>

      <!-- User footer -->
      <div class="sidebar-footer">
        <div class="user-info">
          <div class="user-avatar">{{ userInitial }}</div>
          <div class="user-meta">
            <div class="user-name">{{ userName }}</div>
            <div class="user-role">{{ userRoleLabel }}</div>
          </div>
        </div>
        <button class="sidebar-collapse-btn" @click="toggleSidebar" :title="sidebarCollapsed ? '展开侧边栏' : '折叠侧边栏'">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <polyline v-if="sidebarCollapsed" points="9 18 15 12 9 6"/>
            <polyline v-else points="15 18 9 12 15 6"/>
          </svg>
        </button>
      </div>
    </aside>

    <!-- Main area -->
    <div class="admin-main">
      <header class="admin-topbar">
        <button class="menu-toggle" @click="toggleSidebar">
          <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <line x1="4" x2="20" y1="6" y2="6"/>
            <line x1="4" x2="20" y1="12" y2="12"/>
            <line x1="4" x2="20" y1="18" y2="18"/>
          </svg>
        </button>

        <div class="topbar-right">
          <NuxtLink to="/" target="_blank" class="topbar-link">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M18 13v6a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V8a2 2 0 0 1 2-2h6"/><polyline points="15 3 21 3 21 9"/><line x1="10" x2="21" y1="14" y2="3"/></svg>
            <span>访问网站</span>
          </NuxtLink>
          <button class="topbar-link" @click="openPasswordDrawer">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect width="18" height="11" x="3" y="11" rx="2" ry="2"/><path d="M7 11V7a5 5 0 0 1 10 0v4"/></svg>
            <span>修改密码</span>
          </button>
          <button class="topbar-link logout-btn" @click="handleLogout">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"/><polyline points="16 17 21 12 16 7"/><line x1="21" x2="9" y1="12" y2="12"/></svg>
            <span>退出登录</span>
          </button>
        </div>
      </header>

      <main class="admin-content">
        <slot />
      </main>
    </div>

    <el-drawer
      v-model="passwordDrawerVisible"
      title="修改密码"
      size="min(420px, 92vw)"
      destroy-on-close
      @closed="resetPasswordForm"
    >
      <el-form
        ref="passwordFormRef"
        :model="passwordForm"
        :rules="passwordRules"
        label-position="top"
      >
        <el-form-item label="当前密码" prop="old_password">
          <el-input v-model="passwordForm.old_password" type="password" show-password autocomplete="current-password" />
        </el-form-item>
        <el-form-item label="新密码" prop="new_password">
          <el-input v-model="passwordForm.new_password" type="password" show-password autocomplete="new-password" />
        </el-form-item>
        <el-form-item label="确认新密码" prop="confirm_password">
          <el-input v-model="passwordForm.confirm_password" type="password" show-password autocomplete="new-password" />
        </el-form-item>
      </el-form>
      <template #footer>
        <AdminDrawerFooter
          :loading="passwordSaving"
          @cancel="passwordDrawerVisible = false"
          @confirm="handleChangePassword"
        />
      </template>
    </el-drawer>
  </div>
</template>

<script setup lang="ts">
import { ElMessage, type FormInstance, type FormRules } from 'element-plus';
import { getIconSvg } from '~/composables/lucideIcons';

const route = useRoute();

const sidebarCollapsed = ref(false);
const mobileOpen = ref(false);
const passwordDrawerVisible = ref(false);
const passwordSaving = ref(false);
const passwordFormRef = ref<FormInstance>();

const passwordForm = reactive({
  old_password: '',
  new_password: '',
  confirm_password: '',
});

const passwordRules = computed<FormRules>(() => ({
  old_password: [{ required: true, message: '请输入当前密码', trigger: 'blur' }],
  new_password: [{ required: true, min: 6, message: '新密码至少 6 位', trigger: 'blur' }],
  confirm_password: [
    { required: true, message: '请再次输入新密码', trigger: 'blur' },
    {
      validator: (_rule: unknown, value: string, callback: (error?: Error) => void) => {
        if (value !== passwordForm.new_password) {
          callback(new Error('两次输入的新密码不一致'));
          return;
        }
        callback();
      },
      trigger: 'blur',
    },
  ],
}));

interface NavItem {
  to: string;
  label: string;
  permissions?: string[];
}

interface NavGroup {
  label: string;
  icon: string;
  to?: string;
  children?: NavItem[];
  permissions?: string[];
}

const baseNavGroups: NavGroup[] = [
  { to: '/admin', label: '控制台', icon: getIconSvg('bar-chart', 20) },
  {
    label: '内容管理',
    icon: getIconSvg('file-text', 20),
    children: [
      { to: '/admin/homepage', label: '首页配置' },
      { to: '/admin/navigation', label: '导航管理' },
      { to: '/admin/pages', label: '页面管理' },
      { to: '/admin/media', label: '媒体库' },
    ],
  },
  { to: '/admin/projects', label: '项目管理', icon: getIconSvg('folder', 20) },
  { to: '/admin/faqs', label: 'FAQ 管理', icon: getIconSvg('help-circle', 20) },
  { to: '/admin/cases', label: '案例管理', icon: getIconSvg('users', 20) },
  { to: '/admin/lawyers', label: '律师团队', icon: getIconSvg('briefcase', 20) },
  { to: '/admin/leads', label: '咨询管理', icon: getIconSvg('message-circle', 20) },
  {
    label: '系统设置',
    icon: getIconSvg('settings', 20),
    children: [
      { to: '/admin/users', label: '用户管理' },
      { to: '/admin/settings', label: '网站设置' },
    ],
  },
];

const routePermissions: Record<string, string[]> = {
  '/admin': ['dashboard:read'],
  '/admin/homepage': ['homepage:read'],
  '/admin/navigation': ['navigation:read'],
  '/admin/pages': ['pages:read'],
  '/admin/media': ['media:read'],
  '/admin/projects': ['projects:read'],
  '/admin/faqs': ['faqs:read'],
  '/admin/cases': ['cases:read'],
  '/admin/lawyers': ['lawyers:read'],
  '/admin/leads': ['leads:read'],
  '/admin/users': ['users:read'],
  '/admin/roles': ['roles:read'],
  '/admin/settings': ['settings:read'],
};

const { loadPermissions, hasAnyPermission } = usePermissions();

const withRoutePermissions = <T extends NavItem | NavGroup>(item: T): T => ({
  ...item,
  permissions: item.permissions ?? (item.to ? routePermissions[item.to] : undefined),
});

const navGroups = computed<NavGroup[]>(() => {
  return baseNavGroups.reduce((groups: NavGroup[], baseGroup) => {
    const group = withRoutePermissions(baseGroup);
    if (group.children) {
      const children = group.children.map(withRoutePermissions);
      if (children.some((child) => child.to === '/admin/users') && !children.some((child) => child.to === '/admin/roles')) {
        children.splice(children.findIndex((child) => child.to === '/admin/users') + 1, 0, {
          to: '/admin/roles',
          label: '角色权限',
          permissions: routePermissions['/admin/roles'],
        });
      }
      const visibleChildren = children.filter((child) => hasAnyPermission(child.permissions));
      if (visibleChildren.length > 0) {
        groups.push({ ...group, children: visibleChildren });
      }
      return groups;
    }
    if (hasAnyPermission(group.permissions)) {
      groups.push(group);
    }
    return groups;
  }, []);
});

const expandedGroups = ref(new Set<string>());

const isActive = (to: string) => {
  if (to === '/admin') return route.path === '/admin';
  return route.path.startsWith(to);
};

const isGroupActive = (group: NavGroup) => {
  return group.children?.some((child) => isActive(child.to)) ?? false;
};

const toggleGroup = (label: string) => {
  if (expandedGroups.value.has(label)) {
    expandedGroups.value.delete(label);
  } else {
    expandedGroups.value.add(label);
  }
  expandedGroups.value = new Set(expandedGroups.value);
};

// Auto-expand group when a child route is active
watch(
  () => route.path,
  () => {
    for (const group of navGroups.value) {
      if (group.children && isGroupActive(group)) {
        if (!expandedGroups.value.has(group.label)) {
          expandedGroups.value.add(group.label);
          expandedGroups.value = new Set(expandedGroups.value);
        }
      }
    }
  },
  { immediate: true }
);

function toggleSidebar() {
  if (window.innerWidth < 768) {
    mobileOpen.value = !mobileOpen.value;
  } else {
    sidebarCollapsed.value = !sidebarCollapsed.value;
  }
}

function closeMobile() {
  mobileOpen.value = false;
}

const resetPasswordForm = () => {
  passwordFormRef.value?.resetFields();
  Object.assign(passwordForm, {
    old_password: '',
    new_password: '',
    confirm_password: '',
  });
};

const openPasswordDrawer = () => {
  resetPasswordForm();
  passwordDrawerVisible.value = true;
};

const handleLogout = () => {
  logout();
};

const handleChangePassword = async () => {
  const valid = await passwordFormRef.value?.validate().catch(() => false);
  if (!valid) return;

  passwordSaving.value = true;
  try {
    const api = useApi();
    await api('/admin/me/password', {
      method: 'PUT',
      body: {
        old_password: passwordForm.old_password,
        new_password: passwordForm.new_password,
      },
    });
    ElMessage.success('密码已修改，请重新登录');
    passwordDrawerVisible.value = false;
    logout();
  } catch (e: any) {
    ElMessage.error(e?.data?.message || '修改密码失败');
  } finally {
    passwordSaving.value = false;
  }
};

const { user, logout } = useAuth();

onMounted(() => {
  loadPermissions();
});

const userName = computed(() => (user.value as any)?.display_name || (user.value as any)?.username || '管理员');
const userRoleLabel = computed(() => {
  const role = (user.value as any)?.role;
  if (role === 'admin') return '管理员';
  if (role === 'editor') return '编辑者';
  return '只读用户';
});
const userInitial = computed(() => userName.value.charAt(0).toUpperCase());
</script>

<style scoped>
.admin-layout {
  display: flex;
  min-height: 100vh;
  background: var(--color-bg-app);
  color: var(--color-text);
}

/* Sidebar overlay (mobile) */
.sidebar-overlay {
  display: none;
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.4);
  z-index: 150;
}

/* Sidebar */
.admin-sidebar {
  width: var(--sidebar-width);
  background: var(--color-bg-sidebar);
  color: var(--color-text-sidebar);
  display: flex;
  flex-direction: column;
  flex-shrink: 0;
  transition: width 0.25s ease;
  overflow: hidden;
  position: sticky;
  top: 0;
  height: 100vh;
  border-right: 1px solid var(--color-border-sidebar);
  box-shadow: 1px 0 0 rgba(15, 23, 42, 0.02);
}

.sidebar-header {
  padding: 18px 16px;
  border-bottom: 1px solid var(--color-border-sidebar);
}

.sidebar-logo {
  display: flex;
  align-items: center;
  gap: 10px;
  color: var(--color-text);
  text-decoration: none;
}

.logo-icon {
  width: 32px;
  height: 32px;
  background: var(--color-primary);
  border-radius: var(--radius-sm);
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  flex-shrink: 0;
}

.logo-text {
  font-size: 15px;
  font-weight: 650;
  white-space: nowrap;
  transition: opacity 0.2s ease;
}

/* Nav */
.sidebar-nav {
  flex: 1;
  padding: 12px 8px;
  overflow-y: auto;
  overflow-x: hidden;
}

.nav-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px 12px;
  margin-bottom: 2px;
  color: var(--color-text-sidebar);
  border-radius: var(--radius-sm);
  font-size: 14px;
  font-weight: 500;
  white-space: nowrap;
  border: 1px solid transparent;
  transition: color 0.15s ease, background-color 0.15s ease, border-color 0.15s ease;
}

.nav-item:hover {
  color: var(--color-text);
  background: var(--color-bg-sidebar-hover);
}

.nav-item.active {
  color: var(--color-text-sidebar-active);
  background: var(--color-bg-sidebar-active);
  border-color: #bfdbfe;
}

.nav-icon {
  width: 20px;
  height: 20px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.nav-icon :deep(svg) {
  width: 20px;
  height: 20px;
}

.nav-label {
  transition: opacity 0.2s ease;
}

/* Group title */
.nav-group-title {
  width: 100%;
  background: none;
  border: none;
  cursor: pointer;
  font: inherit;
  color: inherit;
}

.nav-chevron {
  margin-left: auto;
  display: flex;
  align-items: center;
  transition: transform 0.2s ease;
}

.nav-chevron.open {
  transform: rotate(180deg);
}

/* Sub-items */
.nav-sub-items {
  overflow: hidden;
}

.nav-sub-item {
  padding-left: 44px;
  font-size: 13px;
}

/* Sidebar footer */
.sidebar-footer {
  padding: 12px;
  border-top: 1px solid var(--color-border-sidebar);
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 10px;
  min-width: 0;
}

.user-avatar {
  width: 32px;
  height: 32px;
  border-radius: var(--radius-sm);
  background: var(--color-info-soft);
  color: var(--color-primary);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 14px;
  font-weight: 600;
  flex-shrink: 0;
}

.user-meta {
  min-width: 0;
}

.user-name {
  font-size: 13px;
  font-weight: 500;
  color: var(--color-text);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.user-role {
  font-size: 11px;
  color: var(--color-text-sidebar);
}

.sidebar-collapse-btn {
  width: 28px;
  height: 28px;
  border: none;
  background: none;
  color: var(--color-text-sidebar);
  cursor: pointer;
  border-radius: var(--radius-sm);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  transition: all 0.15s ease;
}

.sidebar-collapse-btn:hover {
  color: var(--color-primary);
  background: var(--color-bg-sidebar-hover);
}

/* Collapsed sidebar */
.admin-sidebar.collapsed {
  width: var(--sidebar-collapsed-width);
}

.admin-sidebar.collapsed .logo-text,
.admin-sidebar.collapsed .nav-label,
.admin-sidebar.collapsed .user-meta,
.admin-sidebar.collapsed .user-avatar {
  opacity: 0;
  width: 0;
  overflow: hidden;
  padding: 0;
  margin: 0;
}

.admin-sidebar.collapsed .nav-chevron {
  display: none;
}

.admin-sidebar.collapsed .nav-sub-items {
  display: none;
}

.admin-sidebar.collapsed .sidebar-logo {
  justify-content: center;
}

.admin-sidebar.collapsed .nav-item {
  justify-content: center;
  padding: 10px;
  border-left: none;
}

.admin-sidebar.collapsed .nav-item.active {
  border-left: none;
  border-radius: var(--radius-sm);
}

.admin-sidebar.collapsed .sidebar-footer {
  justify-content: center;
}

.admin-sidebar.collapsed .sidebar-header {
  padding: 16px 8px;
  text-align: center;
}

/* Main area */
.admin-main {
  flex: 1;
  display: flex;
  flex-direction: column;
  min-width: 0;
}

/* Topbar */
.admin-topbar {
  height: var(--topbar-height);
  background: var(--color-bg-surface);
  border-bottom: 1px solid var(--color-border);
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 20px;
  flex-shrink: 0;
  position: sticky;
  top: 0;
  z-index: 50;
  box-shadow: var(--shadow-xs);
}

.menu-toggle {
  background: none;
  border: none;
  cursor: pointer;
  padding: 6px;
  color: var(--color-text-secondary);
  border-radius: var(--radius-sm);
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.15s ease;
}

.menu-toggle:hover {
  color: var(--color-text);
  background: var(--color-bg-app);
}

.topbar-right {
  display: flex;
  align-items: center;
  gap: 4px;
}

.topbar-link {
  font-size: 13px;
  color: var(--color-text-secondary);
  background: none;
  border: none;
  cursor: pointer;
  padding: 8px 12px;
  border-radius: var(--radius-sm);
  display: flex;
  align-items: center;
  gap: 6px;
  transition: all 0.15s ease;
}

.topbar-link:hover {
  color: var(--color-primary);
  background: var(--color-info-soft);
}

.logout-btn:hover {
  color: var(--color-danger);
}

/* Content */
.admin-content {
  flex: 1;
  padding: 24px;
  overflow-y: auto;
}

@media (max-width: 767px) {
  .admin-content {
    padding: 16px;
  }
}

/* Mobile */
@media (max-width: 767px) {
  .sidebar-overlay {
    display: block;
    opacity: 0;
    pointer-events: none;
    transition: opacity 0.25s ease;
  }

  .sidebar-overlay.overlay-visible {
    opacity: 1;
    pointer-events: auto;
  }

  .admin-sidebar {
    position: fixed;
    top: 0;
    left: 0;
    height: 100vh;
    z-index: 200;
    transform: translateX(-100%);
    transition: transform 0.25s ease;
  }

  .admin-sidebar.mobile-open {
    transform: translateX(0);
  }

  .admin-sidebar.collapsed {
    width: var(--sidebar-width);
  }
}
</style>
