<template>
  <div class="contact-page">
    <div class="container">
      <ProjectBreadcrumb />

      <h1 class="page-title">联系北极星移民</h1>
      <p class="page-subtitle">通过北极星移民官方联系方式提交咨询，我们的专业顾问将在24小时内与您联系。</p>

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
                @blur="validateField('name')"
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
                @blur="validateField('phone')"
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
                @blur="validateField('email')"
              />
              <span v-if="errors.email" class="error-text">{{ errors.email }}</span>
            </div>

            <div class="form-group">
              <label class="form-label">意向项目</label>
              <div ref="projectSelectRef" class="custom-select" :class="{ open: projectDropdownOpen }">
                <div class="custom-select-trigger" @click="projectDropdownOpen = !projectDropdownOpen">
                  <span :class="{ placeholder: !selectedProjectName }">{{ selectedProjectName || '-- 请选择意向项目 --' }}</span>
                  <svg class="select-arrow" :class="{ rotated: projectDropdownOpen }" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><polyline points="6 9 12 15 18 9"/></svg>
                </div>
                <div v-show="projectDropdownOpen" class="custom-select-dropdown">
                  <div class="custom-select-option placeholder-option" @click="selectProject('')">-- 请选择意向项目 --</div>
                  <div v-for="p in projectOptions" :key="p.slug" class="custom-select-option" :class="{ selected: form.project === p.slug }" @click="selectProject(p.slug)">{{ p.name }}</div>
                  <div class="custom-select-option" :class="{ selected: form.project === 'other' }" @click="selectProject('other')">其他/尚不确定</div>
                </div>
              </div>
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

            <div v-if="submitError" class="submit-error">
              <span>{{ submitError }}</span>
              <button type="button" class="retry-btn" @click="handleSubmit">重试</button>
            </div>

            <button
              type="submit"
              class="btn-primary btn-submit"
              :disabled="submitting"
            >
              <span v-if="submitting" class="spinner"></span>
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
                <span class="info-label">客服电话</span>
                <a v-if="siteConfig?.contact_phone" :href="`tel:${siteConfig.contact_phone}`" class="info-value">{{ siteConfig.contact_phone }}</a>
                <span v-else class="info-value">400-xxx-xxxx</span>
              </li>
              <li v-if="siteConfig?.contact_phone_2" class="info-item">
                <span class="info-label">联系电话</span>
                <a :href="`tel:${siteConfig.contact_phone_2}`" class="info-value">{{ siteConfig.contact_phone_2 }}</a>
              </li>
              <!--
              <li class="info-item">
                <span class="info-label">邮箱</span>
                <a v-if="siteConfig?.contact_email" :href="`mailto:${siteConfig.contact_email}`" class="info-value">{{ siteConfig.contact_email }}</a>
                <span v-else class="info-value">info@mygo-immigration.com</span>
              </li>
              -->
              <li class="info-item">
                <span class="info-label">地址</span>
                <span class="info-value">{{ siteConfig?.contact_address || '上海市浦东新区陆家嘴金融中心' }}</span>
              </li>
              <li v-if="siteConfig?.contact_wechat" class="info-item">
                <span class="info-label">微信</span>
                <span class="info-value">
                  <img :src="siteConfig.contact_wechat" alt="微信" class="contact-qr-img" loading="lazy" />
                </span>
              </li>
              <li v-if="siteConfig?.contact_wechat_mp" class="info-item">
                <span class="info-label">微信公众号</span>
                <span class="info-value">
                  <img :src="siteConfig.contact_wechat_mp" alt="微信公众号" class="contact-qr-img" loading="lazy" />
                </span>
              </li>
              <li v-if="siteConfig?.contact_wechat_video" class="info-item">
                <span class="info-label">企业视频号</span>
                <span class="info-value">
                  <img :src="siteConfig.contact_wechat_video" alt="企业视频号" class="contact-qr-img" loading="lazy" />
                </span>
              </li>
            </ul>
          </div>
          <!--
          <div class="info-section">
            <h3 class="info-title">服务时间</h3>
            <div class="info-value">
              <p>周一至周五：9:00 - 18:00</p>
              <p>周六：10:00 - 16:00</p>
              <p>周日及法定节假日：休息</p>
            </div>
          </div>
        -->
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
useSeo({ title: '联系我们', description: '联系北极星移民，获取专业投资移民咨询。电话、微信、邮箱、地址等联系方式一应俱全。' });

const { siteConfig } = useMygoSiteConfig();

interface ProjectOption {
  slug: string;
  name: string;
}

usePublicDataFreshness([{ versionKey: 'public:projects:list', dataKey: 'public:projects:options:contact' }])

const projectDropdownOpen = ref(false);
const projectSelectRef = ref<HTMLElement | null>(null);

function handleClickOutside(e: MouseEvent) {
  if (projectSelectRef.value && !projectSelectRef.value.contains(e.target as Node)) {
    projectDropdownOpen.value = false;
  }
}

onMounted(() => {
  document.addEventListener('click', handleClickOutside);
})

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside);
})

const { data: projectListRaw } = await useFetch<{ data?: ProjectOption[] }>('/api/v1/projects/options', {
  key: 'public:projects:options:contact',
  query: { page: 1, per_page: 500 },
})

const projectOptions = computed<ProjectOption[]>(() => {
  const raw = projectListRaw.value as any
  return raw?.data ?? []
})

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

// Custom project select
const selectedProjectName = computed(() => {
  if (!form.project) return '';
  if (form.project === 'other') return '其他/尚不确定';
  const found = projectOptions.value.find(item => item.slug === form.project);
  return found?.name ?? '';
});
const selectProject = (slug: string) => {
  form.project = slug;
  projectDropdownOpen.value = false;
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

const validateField = (field: 'name' | 'phone' | 'email') => {
  const value = form[field].trim();
  if (field === 'name') {
    if (!value) {
      errors.name = '请输入姓名';
    } else if (value.length < 2) {
      errors.name = '姓名至少2个字符';
    } else {
      delete errors.name;
    }
  } else if (field === 'phone') {
    if (!value) {
      errors.phone = '请输入电话号码';
    } else if (!/^1\d{10}$/.test(value)) {
      errors.phone = '请输入有效的手机号码';
    } else {
      delete errors.phone;
    }
  } else if (field === 'email') {
    if (value && !/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(value)) {
      errors.email = '请输入有效的邮箱地址';
    } else {
      delete errors.email;
    }
  }
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
  color: var(--color-danger);
}

.form-input {
  width: 100%;
  padding: 12px 16px;
  font-size: 16px;
  font-family: var(--font-sans);
  border: 1.5px solid var(--color-border);
  border-radius: var(--radius-md);
  background-color: var(--bg-white);
  color: var(--color-text);
  transition: border-color var(--duration-fast) var(--ease-out),
              box-shadow var(--duration-fast) var(--ease-out);
  outline: none;
}

.form-input:focus {
  border-color: var(--color-accent);
  box-shadow: var(--shadow-focus);
}

.form-input.input-error {
  border-color: var(--color-danger);
}

.form-textarea {
  resize: vertical;
  min-height: 120px;
}

/* Custom Select */
.custom-select {
  position: relative;
  user-select: none;
}

.custom-select-trigger {
  width: 100%;
  padding: 12px 16px;
  font-size: 16px;
  font-family: var(--font-sans);
  border: 1.5px solid var(--color-border);
  border-radius: var(--radius-md);
  background-color: var(--bg-white);
  color: var(--color-text);
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
  min-height: 44px;
  transition: border-color var(--duration-fast) var(--ease-out),
              box-shadow var(--duration-fast) var(--ease-out);
}

.custom-select.open .custom-select-trigger {
  border-color: var(--color-accent);
  box-shadow: var(--shadow-focus);
}

.custom-select-trigger .placeholder {
  color: #999;
}

.select-arrow {
  flex-shrink: 0;
  color: var(--color-text-muted);
  transition: transform var(--duration-fast) var(--ease-out);
}

.select-arrow.rotated {
  transform: rotate(180deg);
}

.custom-select-dropdown {
  position: absolute;
  top: calc(100% + 4px);
  left: 0;
  right: 0;
  max-height: 240px;
  overflow-y: auto;
  background-color: var(--bg-white);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  box-shadow: var(--shadow-lg);
  z-index: 100;
}

.custom-select-option {
  padding: 10px 16px;
  font-size: 15px;
  color: var(--color-text);
  cursor: pointer;
  transition: background-color 0.15s;
}

.custom-select-option:hover {
  background-color: var(--color-accent-soft);
}

.custom-select-option.selected {
  color: var(--color-accent);
  font-weight: 600;
}

.custom-select-option.placeholder-option {
  color: #999;
}

.error-text {
  display: block;
  font-size: 13px;
  color: var(--color-danger);
  margin-top: 6px;
}

.submit-error {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  background-color: var(--color-danger-soft);
  color: var(--color-danger);
  padding: 12px 16px;
  border-radius: var(--radius-md);
  font-size: 14px;
  margin-bottom: 16px;
}

.retry-btn {
  flex-shrink: 0;
  padding: 4px 14px;
  background: none;
  border: 1.5px solid var(--color-danger);
  color: var(--color-danger);
  border-radius: var(--radius-sm);
  font-size: 13px;
  font-weight: 600;
  cursor: pointer;
  transition: all var(--duration-fast) var(--ease-out);
}

.retry-btn:hover {
  background: var(--color-danger);
  color: #fff;
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

.spinner {
  display: inline-block;
  width: 16px;
  height: 16px;
  border: 2px solid rgba(255, 255, 255, 0.3);
  border-top-color: #fff;
  border-radius: 50%;
  animation: spin 0.6s linear infinite;
  margin-right: 6px;
  vertical-align: middle;
}

@keyframes spin {
  to { transform: rotate(360deg); }
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
  background-color: var(--color-success-soft);
  color: var(--color-success);
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
  width: 80px;
  flex-shrink: 0;
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

.contact-qr-img {
  width: 120px;
  aspect-ratio: 1;
  object-fit: contain;
  border-radius: var(--radius-sm);
  border: 1px solid var(--border-color);
}

@media (max-width: 1023px) {
  .page-title {
    font-size: 32px;
  }

  .contact-grid {
    grid-template-columns: 1fr 1fr;
    gap: 32px;
  }

  .contact-qr-img {
    width: 100px;
  }
}

@media (max-width: 767px) {
  .page-title {
    font-size: 28px;
  }

  .page-subtitle {
    font-size: 15px;
    margin-bottom: 24px;
  }

  .contact-grid {
    grid-template-columns: 1fr;
    gap: 32px;
  }

  .contact-info {
    display: none;
  }

  .contact-form {
    padding: 24px;
  }

  .form-input {
    min-height: 44px;
  }

  .success-message {
    padding: 40px 24px;
  }

  .success-icon {
    width: 64px;
    height: 64px;
    font-size: 28px;
  }

  .contact-qr-img {
    width: 100px;
  }
}
</style>
