<template>
  <div class="contact-page">
    <div class="container">
      <ProjectBreadcrumb :items="breadcrumbs" />

      <h1 class="page-title">联系北极星移民</h1>
      <p class="page-subtitle">留下联系方式，与顾问讨论您的需求。</p>

      <div class="contact-grid">
        <!-- Contact Form -->
        <div class="contact-form-wrapper">
          <div v-if="submitSuccess" ref="successRef" class="success-message" tabindex="-1" role="status">
            <div class="success-icon" aria-hidden="true">&#10003;</div>
            <h3>提交成功</h3>
            <p>顾问会尽快联系您。</p>
            <button type="button" class="btn-primary" @click="resetForm">再次提交</button>
          </div>

          <form v-else ref="formRef" class="contact-form" novalidate @submit.prevent="handleSubmit">
            <div class="form-group">
              <label class="form-label" for="name">
                姓名 <span class="required">*</span>
              </label>
              <input
                id="name"
                v-model="form.name"
                type="text"
                autocomplete="name"
                required
                class="form-input"
                :class="{ 'input-error': errors.name }"
                :aria-invalid="!!errors.name"
                :aria-describedby="errors.name ? 'contact-name-error' : undefined"
                placeholder="请输入姓名"
                @input="errors.name && validateField('name')"
                @blur="validateField('name')"
              />
              <div class="field-feedback">
                <span v-if="errors.name" id="contact-name-error" class="error-text">{{ errors.name }}</span>
              </div>
            </div>

            <div class="form-group">
              <label class="form-label" for="phone">
                联系电话 <span class="required">*</span>
              </label>
              <input
                id="phone"
                v-model="form.phone"
                type="tel"
                autocomplete="tel"
                required
                class="form-input"
                :class="{ 'input-error': errors.phone }"
                :aria-invalid="!!errors.phone"
                :aria-describedby="errors.phone ? 'contact-phone-error' : 'contact-phone-hint'"
                placeholder="请输入联系电话"
                @input="errors.phone && validateField('phone')"
                @blur="validateField('phone')"
              />
              <div class="field-feedback">
                <span v-if="errors.phone" id="contact-phone-error" class="error-text">{{ errors.phone }}</span>
                <span v-else id="contact-phone-hint">境外号码请加区号，例如 +852。</span>
              </div>
            </div>

            <div class="form-group">
              <label class="form-label" for="email">
                邮箱（选填）
              </label>
              <input
                id="email"
                v-model="form.email"
                type="email"
                autocomplete="email"
                class="form-input"
                :class="{ 'input-error': errors.email }"
                :aria-invalid="!!errors.email"
                :aria-describedby="errors.email ? 'contact-email-error' : undefined"
                placeholder="请输入邮箱"
                @input="errors.email && validateField('email')"
                @blur="validateField('email')"
              />
              <div class="field-feedback">
                <span v-if="errors.email" id="contact-email-error" class="error-text">{{ errors.email }}</span>
              </div>
            </div>

            <div class="form-group">
              <PublicProjectSelect
                v-model="selectedProject"
                id="contact-project"
                label="意向项目（选填）"
                placeholder="请选择意向项目"
                :options="projectOptions"
                :loading="projectListPending"
                :error="projectListError ? '项目列表加载失败，请重试' : undefined"
                @retry="refreshProjectList()"
              />
              <p v-if="showProjectPrefillHint" class="field-hint">已根据您浏览的项目预选，可更改。</p>
            </div>

            <div class="form-group">
              <label class="form-label" for="message">留言（选填）</label>
              <textarea
                id="message"
                v-model="message"
                class="form-input form-textarea"
                rows="3"
                placeholder="告诉我们您想了解的问题"
              ></textarea>
            </div>

            <div v-if="submitError" class="submit-error">
              <span role="alert">{{ submitError }}</span>
              <button type="button" class="retry-btn" :disabled="submitting" @click="handleSubmit">重试</button>
            </div>

            <button
              type="submit"
              class="btn-primary btn-submit"
              :disabled="submitting"
            >
              <span v-if="submitting" class="spinner" aria-hidden="true"></span>
              {{ submitting ? '提交中...' : '提交咨询' }}
            </button>
            <p class="submission-hint">顾问将使用您填写的联系方式回复本次咨询。</p>
          </form>
          <div v-if="siteConfig?.contact_phone || siteConfig?.contact_wechat" class="contact-quick-links">
            <a v-if="siteConfig?.contact_phone" :href="`tel:${siteConfig.contact_phone}`">拨打咨询热线</a>
            <details v-if="siteConfig?.contact_wechat">
              <summary>微信咨询</summary>
              <img class="contact-qr-img" :src="siteConfig.contact_wechat" alt="微信咨询二维码" loading="lazy" />
            </details>
          </div>
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
import { HOME_BREADCRUMB } from '~/utils/breadcrumb'

/** 联系页面的稳定面包屑路径。 */
const breadcrumbs = [
  HOME_BREADCRUMB,
  { label: '联系我们', href: '/contact' },
]

useSeo({
  title: '联系我们',
  description: '联系北极星移民，获取专业投资移民咨询。电话、微信、邮箱、地址等联系方式一应俱全。',
  breadcrumbs,
});

const { siteConfig } = useMygoSiteConfig();
const route = useRoute();

interface ProjectOption {
  slug: string;
  name: string;
}

usePublicDataFreshness([{ versionKey: 'public:projects:list', dataKey: 'public:projects:options:contact' }])

const {
  data: projectListRaw,
  pending: projectListPending,
  error: projectListError,
  refresh: refreshProjectList,
} = await useFetch<{ data?: ProjectOption[] }>('/api/v1/projects/options', {
  key: 'public:projects:options:contact',
  query: { page: 1, per_page: 500 },
})

const publishedProjects = computed(() => projectListRaw.value?.data ?? []);
const projectOptions = computed(() => [
  { value: '', label: '暂不选择' },
  ...publishedProjects.value.map(({ slug, name }) => ({ value: slug, label: name })),
  { value: 'other', label: '其他或尚不确定' },
]);

const form = reactive({
  name: '',
  phone: '',
  email: '',
});

const sourceProject = computed(() => publishedProjects.value.find(({ slug }) => slug === route.query.project));
const sourceComparison = computed(() => publishedProjects.value.find(({ slug }) => slug === route.query.compare));
const projectOverride = ref<string | null>(null);
const messageOverride = ref<string | null>(null);
const selectedProject = computed({
  get: () => projectOverride.value ?? sourceProject.value?.slug ?? '',
  set: (value) => { projectOverride.value = value; },
});
const showProjectPrefillHint = computed(() => projectOverride.value === null && !!sourceProject.value);
const comparisonMessage = computed(() => {
  const primary = sourceProject.value;
  const secondary = sourceComparison.value;
  if (!primary || !secondary || primary.slug === secondary.slug || selectedProject.value !== primary.slug) return '';
  return `我想咨询「${primary.name}」与「${secondary.name}」的对比。`;
});
const message = computed({
  get: () => messageOverride.value ?? comparisonMessage.value,
  set: (value) => { messageOverride.value = value; },
});

const fields = ['name', 'phone', 'email'] as const;
type ContactField = typeof fields[number];
const errors = reactive<Partial<Record<ContactField, string>>>({});
const formRef = ref<HTMLFormElement | null>(null);
const successRef = ref<HTMLElement | null>(null);
const submitting = ref(false);
const submitSuccess = ref(false);
const submitError = ref<string | null>(null);

function phoneNumber(value: string) {
  return value.replace(/[\s()-]/g, '');
}

function fieldError(field: ContactField) {
  const value = form[field].trim();
  switch (field) {
    case 'name':
      if (!value) return '请输入姓名';
      if (value.length < 2) return '姓名至少 2 个字符';
      break;
    case 'phone':
      if (!value) return '请输入电话号码';
      if (!/^(1\d{10}|\+[1-9]\d{6,14})$/.test(phoneNumber(value))) {
        return '请输入有效的联系电话，境外号码请加区号';
      }
      break;
    case 'email':
      if (value && !/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(value)) return '请输入有效的邮箱地址';
      break;
  }
}

function validateField(field: ContactField) {
  const error = fieldError(field);
  if (error) errors[field] = error;
  else delete errors[field];
}

async function handleSubmit() {
  if (submitting.value) return;
  fields.forEach(validateField);
  if (fields.some((field) => errors[field])) {
    await nextTick();
    formRef.value?.querySelector<HTMLElement>('[aria-invalid="true"]')?.focus();
    return;
  }

  submitting.value = true;
  submitError.value = null;

  try {
    await $fetch('/api/v1/leads', {
      method: 'POST',
      body: {
        name: form.name.trim(),
        phone: phoneNumber(form.phone),
        email: form.email.trim(),
        interested_project: selectedProject.value,
        message: message.value.trim(),
        source: 'website_contact',
      },
    });
    submitSuccess.value = true;
    await nextTick();
    successRef.value?.focus();
  } catch {
    submitError.value = '提交失败，请稍后重试';
  } finally {
    submitting.value = false;
  }
}

async function resetForm() {
  form.name = '';
  form.phone = '';
  form.email = '';
  projectOverride.value = null;
  messageOverride.value = null;
  submitSuccess.value = false;
  submitError.value = null;
  fields.forEach((field) => delete errors[field]);
  await nextTick();
  formRef.value?.querySelector<HTMLInputElement>('#name')?.focus();
}

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
  color: var(--text-secondary);
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
  padding: 32px;
  border: 1px solid var(--color-border);
  border-radius: var(--radius-lg);
}

.form-group {
  margin-bottom: 16px;
}

.form-group:has(.field-feedback) {
  margin-bottom: 8px;
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
  min-height: 44px;
  padding: 10px 12px;
  font-size: 16px;
  font-family: var(--font-sans);
  border: 1px solid var(--public-control-border);
  border-radius: var(--radius-md);
  background-color: var(--bg-white);
  color: var(--color-text);
  transition: border-color var(--duration-fast) var(--ease-out),
              box-shadow var(--duration-fast) var(--ease-out);
}

.form-input:focus-visible {
  outline: 2px solid var(--public-control-focus);
  outline-offset: 2px;
  border-radius: var(--radius-md);
}

.form-input::placeholder {
  color: var(--color-text-secondary);
}

.form-input.input-error {
  border-color: var(--color-danger);
}

.form-textarea {
  resize: vertical;
  min-height: 88px;
}

.field-hint,
.field-feedback,
.submission-hint {
  margin-top: 6px;
  color: var(--color-text-secondary);
  font-size: 13px;
}

.field-feedback {
  min-height: 1.6em;
}

.error-text {
  color: var(--color-danger);
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
  font-weight: 700;
  min-height: 44px;
  color: var(--color-primary);
  text-shadow: none;
  box-shadow: none;
}

.btn-submit:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.btn-submit:hover {
  background: var(--gradient-gold-hover);
  transform: none;
  box-shadow: none;
}

.spinner {
  display: inline-block;
  width: 16px;
  height: 16px;
  border: 2px solid transparent;
  border-top-color: currentColor;
  border-right-color: currentColor;
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

.contact-quick-links {
  display: none;
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
    gap: 16px;
  }

  .contact-info {
    display: none;
  }

  .contact-form {
    padding: 0;
    border: 0;
  }

  .contact-quick-links {
    display: flex;
    flex-wrap: wrap;
    align-items: flex-start;
    gap: 8px 16px;
    margin-top: 14px;
    padding-top: 4px;
    border-top: 1px solid var(--color-border);
    color: var(--color-primary);
    font-size: 14px;
    font-weight: 600;
  }

  .contact-quick-links a,
  .contact-quick-links summary {
    min-height: 44px;
    padding-block: 11px;
    cursor: pointer;
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
