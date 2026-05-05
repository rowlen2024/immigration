<template>
  <div class="contact-page">
    <div class="container">
      <ProjectBreadcrumb />

      <h1 class="page-title">联系我们</h1>
      <p class="page-subtitle">填写下方表单，我们的专业顾问将在24小时内与您联系</p>

      <div class="contact-grid">
        <!-- Contact Form -->
        <div class="contact-form-wrapper">
          <div v-if="submitSuccess" class="success-message">
            <div class="success-icon">&#10003;</div>
            <h3>提交成功</h3>
            <p>提交成功，我们会尽快联系您</p>
            <button class="btn-primary" @click="resetForm">再次提交</button>
          </div>

          <form v-else class="contact-form" @submit.prevent="handleSubmit">
            <div class="form-group">
              <label class="form-label" for="name">
                姓名 <span class="required">*</span>
              </label>
              <input
                id="name"
                v-model="form.name"
                type="text"
                class="form-input"
                :class="{ 'input-error': errors.name }"
                placeholder="请输入您的姓名"
              />
              <span v-if="errors.name" class="error-text">{{ errors.name }}</span>
            </div>

            <div class="form-group">
              <label class="form-label" for="phone">
                电话 <span class="required">*</span>
              </label>
              <input
                id="phone"
                v-model="form.phone"
                type="tel"
                class="form-input"
                :class="{ 'input-error': errors.phone }"
                placeholder="请输入您的手机号码"
              />
              <span v-if="errors.phone" class="error-text">{{ errors.phone }}</span>
            </div>

            <div class="form-group">
              <label class="form-label" for="email">
                邮箱
              </label>
              <input
                id="email"
                v-model="form.email"
                type="email"
                class="form-input"
                :class="{ 'input-error': errors.email }"
                placeholder="请输入您的邮箱地址"
              />
              <span v-if="errors.email" class="error-text">{{ errors.email }}</span>
            </div>

            <div class="form-group">
              <label class="form-label" for="project">意向项目</label>
              <select id="project" v-model="form.project" class="form-input">
                <option value="">-- 请选择意向项目 --</option>
                <option v-for="p in projectOptions" :key="p.slug" :value="p.slug">{{ p.name }}</option>
                <option value="other">其他/尚不确定</option>
              </select>
            </div>

            <div class="form-group">
              <label class="form-label" for="message">留言</label>
              <textarea
                id="message"
                v-model="form.message"
                class="form-input form-textarea"
                :class="{ 'input-error': errors.message }"
                rows="5"
                placeholder="请描述您的情况和需求（可选）"
              ></textarea>
            </div>

            <div v-if="submitError" class="submit-error">{{ submitError }}</div>

            <button
              type="submit"
              class="btn-primary btn-submit"
              :disabled="submitting"
            >
              {{ submitting ? '提交中...' : '提交咨询' }}
            </button>
          </form>
        </div>

        <!-- Contact Info Sidebar -->
        <div class="contact-info">
          <div class="info-section">
            <h3 class="info-title">联系方式</h3>
            <ul class="info-list">
              <li class="info-item">
                <span class="info-label">电话</span>
                <a v-if="siteConfig?.contact_phone" :href="`tel:${siteConfig.contact_phone}`" class="info-value">{{ siteConfig.contact_phone }}</a>
                <span v-else class="info-value">400-xxx-xxxx</span>
              </li>
              <li class="info-item">
                <span class="info-label">邮箱</span>
                <a v-if="siteConfig?.contact_email" :href="`mailto:${siteConfig.contact_email}`" class="info-value">{{ siteConfig.contact_email }}</a>
                <span v-else class="info-value">info@mygo-immigration.com</span>
              </li>
              <li class="info-item">
                <span class="info-label">地址</span>
                <span class="info-value">{{ siteConfig?.contact_address || '上海市浦东新区陆家嘴金融中心' }}</span>
              </li>
              <li class="info-item">
                <span class="info-label">微信</span>
                <span class="info-value">{{ siteConfig?.contact_wechat || 'MyGo_Immigration' }}</span>
              </li>
            </ul>
          </div>

          <div class="info-section">
            <h3 class="info-title">服务时间</h3>
            <div class="info-value">
              <p>周一至周五：9:00 - 18:00</p>
              <p>周六：10:00 - 16:00</p>
              <p>周日及法定节假日：休息</p>
            </div>
          </div>

          <div class="info-section">
            <h3 class="info-title">常见需求</h3>
            <ul class="info-quick-links">
              <li v-for="link in quickLinks" :key="link.label">
                <NuxtLink :to="link.link">{{ link.label }}</NuxtLink>
              </li>
            </ul>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
useSeo({ title: '联系我们' });

const { siteConfig, fetch: fetchSiteConfig } = useSiteConfig();

interface ProjectOption {
  slug: string;
  name: string;
}

interface QuickLink {
  label: string;
  link: string;
}

const projectOptions = ref<ProjectOption[]>([]);
const quickLinks = ref<QuickLink[]>([]);

interface ContactForm {
  name: string;
  phone: string;
  email: string;
  project: string;
  message: string;
}

const form = reactive<ContactForm>({
  name: '',
  phone: '',
  email: '',
  project: '',
  message: '',
});

const errors = reactive<Record<string, string>>({});
const submitting = ref(false);
const submitSuccess = ref(false);
const submitError = ref<string | null>(null);

const fetchProjects = async () => {
  try {
    const api = useApi();
    const data = await api<{ items: ProjectOption[] }>('/projects?per_page=100');
    if (data?.items?.length) {
      projectOptions.value = data.items;
    }
  } catch { /* keep empty, dropdown falls back to base options */ }
};

const fetchQuickLinks = async () => {
  try {
    const api = useApi();
    const navItems = await api<QuickLink[]>('/navigation?position=header');
    if (navItems && Array.isArray(navItems)) {
      // Flatten: take top-level items that have children, grab their children
      const links: QuickLink[] = [];
      for (const item of navItems) {
        const children = (item as any).children as QuickLink[] | undefined;
        if (children?.length) {
          links.push(...children);
        } else if ((item as any).link) {
          links.push({ label: (item as any).label, link: (item as any).link });
        }
      }
      quickLinks.value = links;
    }
  } catch { /* keep empty, sidebar hides empty list naturally */ }
};

const validate = (): boolean => {
  const newErrors: Record<string, string> = {};

  if (!form.name.trim()) {
    newErrors.name = '请输入姓名';
  } else if (form.name.trim().length < 2) {
    newErrors.name = '姓名至少2个字符';
  }

  if (!form.phone.trim()) {
    newErrors.phone = '请输入电话号码';
  } else if (!/^1\d{10}$/.test(form.phone.trim())) {
    newErrors.phone = '请输入有效的手机号码';
  }

  if (!form.email.trim()) {
    //newErrors.email = '请输入邮箱地址';
  } else if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(form.email.trim())) {
    newErrors.email = '请输入有效的邮箱地址';
  }

  Object.keys(errors).forEach((key) => delete errors[key]);
  Object.assign(errors, newErrors);
  return Object.keys(newErrors).length === 0;
};

const handleSubmit = async () => {
  if (!validate()) return;

  submitting.value = true;
  submitError.value = null;

  try {
    await $fetch('/api/v1/leads', {
      method: 'POST',
      body: {
        name: form.name.trim(),
        phone: form.phone.trim(),
        email: form.email.trim(),
        interested_project: form.project,
        message: form.message.trim(),
        source: 'website_contact',
      },
    });
    submitSuccess.value = true;
  } catch (err) {
    submitError.value = err instanceof Error ? err.message : '提交失败，请稍后重试';
  } finally {
    submitting.value = false;
  }
};

const resetForm = () => {
  form.name = '';
  form.phone = '';
  form.email = '';
  form.project = '';
  form.message = '';
  submitSuccess.value = false;
  submitError.value = null;
  Object.keys(errors).forEach((key) => delete errors[key]);
};

onMounted(() => {
  fetchSiteConfig();
  fetchProjects();
  fetchQuickLinks();
});
</script>

<style scoped>
.page-title {
  font-size: 36px;
  font-weight: 700;
  color: var(--text-primary);
  margin-bottom: 12px;
}

.page-subtitle {
  font-size: 16px;
  color: var(--text-light);
  margin-bottom: 40px;
}

.contact-grid {
  display: grid;
  grid-template-columns: 1.2fr 0.8fr;
  gap: 48px;
  margin-bottom: 48px;
}

/* Form */
.contact-form {
  background-color: var(--bg-white);
  padding: 40px;
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-sm);
}

.form-group {
  margin-bottom: 24px;
}

.form-label {
  display: block;
  font-size: 14px;
  font-weight: 600;
  color: var(--text-primary);
  margin-bottom: 8px;
}

.required {
  color: #c62828;
}

.form-input {
  width: 100%;
  padding: 12px 16px;
  font-size: 15px;
  font-family: var(--font-sans);
  border: 2px solid var(--border-color);
  border-radius: var(--radius-md);
  background-color: var(--bg-white);
  color: var(--text-primary);
  transition: border-color 0.3s ease;
  outline: none;
}

.form-input:focus {
  border-color: var(--accent);
}

.form-input.input-error {
  border-color: #c62828;
}

.form-textarea {
  resize: vertical;
  min-height: 120px;
}

.error-text {
  display: block;
  font-size: 13px;
  color: #c62828;
  margin-top: 6px;
}

.submit-error {
  background-color: #fce8e8;
  color: #c62828;
  padding: 12px 16px;
  border-radius: var(--radius-md);
  font-size: 14px;
  margin-bottom: 16px;
}

.btn-submit {
  width: 100%;
  font-size: 16px;
  padding: 14px;
}

.btn-submit:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

/* Success Message */
.success-message {
  text-align: center;
  padding: 60px 40px;
  background-color: var(--bg-white);
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-sm);
}

.success-icon {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 80px;
  height: 80px;
  border-radius: 50%;
  background-color: #e6f4ea;
  color: #1e7e34;
  font-size: 36px;
  margin-bottom: 20px;
}

.success-message h3 {
  font-size: 24px;
  font-weight: 700;
  color: var(--text-primary);
  margin-bottom: 8px;
}

.success-message p {
  font-size: 16px;
  color: var(--text-secondary);
  margin-bottom: 28px;
}

/* Contact Info */
.contact-info {
  display: flex;
  flex-direction: column;
  gap: 32px;
}

.info-section {
  background-color: var(--bg-light);
  padding: 28px;
  border-radius: var(--radius-lg);
}

.info-title {
  font-size: 18px;
  font-weight: 700;
  color: var(--text-primary);
  margin-bottom: 16px;
  padding-bottom: 12px;
  border-bottom: 2px solid var(--accent);
}

.info-list {
  display: flex;
  flex-direction: column;
  gap: 14px;
}

.info-item {
  display: flex;
  gap: 12px;
  font-size: 14px;
}

.info-label {
  font-weight: 600;
  color: var(--text-primary);
  white-space: nowrap;
  min-width: 48px;
}

.info-value {
  color: var(--text-secondary);
  line-height: 1.8;
}

.info-value a {
  color: var(--primary);
  transition: color 0.2s ease;
}

.info-value a:hover {
  color: var(--accent-dark);
}

.info-quick-links {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.info-quick-links a {
  font-size: 14px;
  color: var(--primary);
  transition: color 0.2s ease;
}

.info-quick-links a:hover {
  color: var(--accent-dark);
}

@media (max-width: 1023px) {
  .contact-grid {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 767px) {
  .page-title {
    font-size: 28px;
  }

  .contact-form {
    padding: 24px;
  }
}
</style>
