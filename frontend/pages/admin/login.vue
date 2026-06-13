<template>
  <div class="login-page">
    <div class="login-card">
      <div class="login-logo">
        <svg width="48" height="48" viewBox="0 0 48 48" fill="none" xmlns="http://www.w3.org/2000/svg">
          <rect width="48" height="48" rx="12" fill="#0f172a"/>
          <path d="M14 24L21 17L27 23L34 16" stroke="#e2a83e" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"/>
          <path d="M14 32L21 25L27 31L34 24" stroke="#e2a83e" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"/>
          <circle cx="14" cy="24" r="2" fill="#e2a83e"/>
          <circle cx="34" cy="16" r="2" fill="#e2a83e"/>
          <circle cx="14" cy="32" r="2" fill="#e2a83e"/>
          <circle cx="34" cy="24" r="2" fill="#e2a83e"/>
        </svg>
      </div>
      <h1 class="login-title">北极星移民管理后台</h1>
      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        label-position="top"
        @submit.prevent="handleLogin"
      >
        <el-form-item label="用户名" prop="username">
          <el-input
            v-model="form.username"
            placeholder="请输入用户名"
            size="large"
          />
        </el-form-item>
        <el-form-item label="密码" prop="password">
          <el-input
            v-model="form.password"
            type="password"
            placeholder="请输入密码"
            size="large"
            show-password
          />
        </el-form-item>
        <el-form-item>
          <el-button
            type="primary"
            size="large"
            class="login-btn"
            :loading="loading"
            native-type="submit"
          >
            登 录
          </el-button>
        </el-form-item>
      </el-form>
      <p v-if="errorMsg" class="error-msg">{{ errorMsg }}</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { FormInstance, FormRules } from 'element-plus';

definePageMeta({ layout: false });

const formRef = ref<FormInstance>();
const loading = ref(false);
const errorMsg = ref('');

const form = reactive({
  username: '',
  password: '',
});

const rules: FormRules = {
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  password: [{ required: true, message: '请输入密码', trigger: 'blur' }],
};

const handleLogin = async () => {
  const valid = await formRef.value?.validate().catch(() => false);
  if (!valid) return;

  loading.value = true;
  errorMsg.value = '';

  try {
    const { login } = useAuth();
    await login({ username: form.username, password: form.password });
    const router = useRouter();
    router.push('/admin');
  } catch (err: any) {
    errorMsg.value = err?.data?.message || err?.message || '登录失败，请检查用户名和密码';
  } finally {
    loading.value = false;
  }
};
</script>

<style scoped>
.login-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #0f172a 0%, #1e293b 100%);
}

.login-card {
  width: 400px;
  padding: 48px 40px 40px;
  background: var(--color-bg-surface);
  border-radius: var(--radius-md);
  box-shadow: var(--shadow-lg);
}

.login-logo {
  text-align: center;
  margin-bottom: 16px;
}

.login-title {
  font-size: 22px;
  font-weight: 600;
  text-align: center;
  margin-bottom: 32px;
  color: var(--color-primary);
}

.login-btn {
  width: 100%;
  height: 42px;
}

.error-msg {
  color: var(--color-danger);
  text-align: center;
  font-size: 14px;
  margin-top: 8px;
}

/* Override Element Plus input border-radius */
:deep(.el-input__wrapper) {
  border-radius: var(--radius-md);
}

:deep(.el-button--primary) {
  --el-button-bg-color: var(--color-primary);
  --el-button-border-color: var(--color-primary);
  --el-button-hover-bg-color: #1e293b;
  --el-button-hover-border-color: #1e293b;
}
</style>
