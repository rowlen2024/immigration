<template>
  <div>
    <AdminPageHeader
      title="网站设置"
      description="统一管理站点基础信息、SEO、联系方式和第三方代码"
    />

    <div class="settings-body admin-panel-shell">
      <AdminLoadingOverlay :show="loading" />

      <el-card v-for="group in groups" :key="group.key" class="admin-settings-card">
        <template #header>
          <h3 class="admin-card-title">{{ group.label }}</h3>
        </template>

        <el-form label-width="160px" label-position="right">
          <el-form-item
            v-for="field in group.fields"
            :key="field.key"
            :label="field.label"
          >
            <div class="admin-field-wrap">
              <div v-if="field.key === 'same_as'" class="array-input">
                <div
                  v-for="(item, idx) in form.same_as"
                  :key="idx"
                  class="array-row"
                >
                  <el-input v-model="form.same_as[idx]" placeholder="https://" />
                  <el-button
                    type="danger"
                    :icon="Delete"
                    circle
                    size="small"
                    @click="removeSameAs(idx)"
                  />
                </div>
                <el-button type="primary" text @click="addSameAs">
                  + 添加链接
                </el-button>
              </div>

              <el-input
                v-else-if="field.textarea"
                v-model="form[field.key]"
                type="textarea"
                :rows="field.rows || 6"
                class="monospace-input"
              />

              <ImageInput
                v-else-if="field.image"
                v-model="form[field.key]"
                :placeholder="field.placeholder"
                :size-hint="field.sizeHint"
                :preview-ratio="field.previewRatio"
                :context="field.context"
              />

              <el-input v-else v-model="form[field.key]" :placeholder="field.placeholder" />

              <el-tooltip
                :content="field.tip"
                placement="right"
                effect="dark"
                raw-content
              >
                <span class="admin-tip-icon" v-html="getIconSvg('help-circle', 14)"></span>
              </el-tooltip>
            </div>
          </el-form-item>
        </el-form>
      </el-card>

      <div class="admin-save-bar">
        <el-button type="primary" size="large" :loading="saving" @click="save">
          保存设置
        </el-button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { Delete } from '@element-plus/icons-vue';
import ImageInput from '~/components/admin/ImageInput.vue';
import { getIconSvg } from '~/composables/lucideIcons';
import { useNotify } from '~/composables/useNotify';

definePageMeta({ layout: 'admin', middleware: ['auth'] });

const notify = useNotify();

interface SiteConfig {
  [key: string]: any;
  site_name: string;
  site_logo: string;
  site_favicon: string;
  seo_title: string;
  seo_description: string;
  seo_keywords: string;
  og_image: string;
  canonical_base: string;
  organization_name: string;
  organization_description: string;
  organization_logo: string;
  organization_url: string;
  same_as: string[];
  contact_phone: string;
  contact_phone_2: string;
  contact_email: string;
  contact_address: string;
  contact_wechat: string;
  contact_wechat_mp: string;
  contact_wechat_video: string;
  ga_tracking_id: string;
  baidu_tongji_id: string;
  custom_head_code: string;
  custom_body_code: string;
  copyright_text: string;
  icp_number: string;
  footer_tagline: string;
}

const defaultForm = (): SiteConfig => ({
  site_name: '北极星移民',
  site_logo: '',
  site_favicon: '',
  seo_title: '{site_name} | 专业投资移民服务',
  seo_description: '',
  seo_keywords: '',
  og_image: '',
  canonical_base: '',
  organization_name: '',
  organization_description: '',
  organization_logo: '',
  organization_url: '',
  same_as: [],
  contact_phone: '',
  contact_phone_2: '',
  contact_email: '',
  contact_address: '',
  contact_wechat: '',
  contact_wechat_mp: '',
  contact_wechat_video: '',
  ga_tracking_id: '',
  baidu_tongji_id: '',
  custom_head_code: '',
  custom_body_code: '',
  copyright_text: '© {year} {site_name}. All rights reserved.',
  icp_number: '',
  footer_tagline: '',
});

const tips: Record<string, string> = {
  site_name: '网站在导航栏、浏览器标题栏等位置显示的名称',
  site_logo: '网站主 Logo 图片地址，显示在页头导航栏',
  site_favicon: '浏览器标签页和书签栏显示的小图标，建议 32×32px',
  seo_title: '搜索引擎结果页显示的标题模板。<br/><code>{site_name}</code> 会被替换为网站名称',
  seo_description: '搜索引擎结果页显示的页面描述，建议 120-160 字，各内页无独立描述时使用此值',
  seo_keywords: '页面关键词，用逗号分隔。现代搜索引擎权重已降低，但仍建议填写',
  og_image: '在微信、Facebook 等社交平台分享链接时显示的默认预览图',
  canonical_base: '网站首选访问域名，用于规范搜索引擎索引（如 https://www.example.com）',
  organization_name: '向 Google 知识图谱和 AI 搜索声明的企业法律实体全称',
  organization_description: '用于生成 AI 搜索摘要和知识面板的企业介绍，建议包含成立时间、核心业务和规模',
  organization_logo: 'Google 知识面板展示的 Logo，建议 112×112px 以上高清 PNG',
  organization_url: '声明机构的官方网址，用于多平台交叉验证网站真实性',
  same_as: '与官网关联的社交媒体链接（LinkedIn、公众号、小红书等）。AI 搜索引擎通过双向链接验证实体可信度',
  contact_phone: '网站底部和结构化数据中展示的客服电话',
  contact_phone_2: '网站底部和结构化数据中展示的联系电话',
  contact_email: '网站底部和结构化数据中展示的客服邮箱',
  contact_address: '公司办公地址，显示在网站底部',
  contact_wechat: '微信二维码图片，显示在联系我们页面',
  contact_wechat_mp: '微信公众号二维码图片，显示在联系我们页面',
  contact_wechat_video: '企业视频号二维码图片，显示在联系我们页面',
  ga_tracking_id: 'Google Analytics 4 衡量 ID（格式：G-XXXXXXXX），用于网站流量分析',
  baidu_tongji_id: '百度统计站点 ID，用于中国国内流量分析',
  custom_head_code: '插入到每个页面 &lt;head&gt; 标签内的自定义代码（meta 验证标签、第三方脚本等）',
  custom_body_code: '插入到每个页面 &lt;/body&gt; 闭合标签前的自定义代码',
  copyright_text: '页脚版权声明。<code>{year}</code> 动态替换为当前年份，<code>{site_name}</code> 替换为网站名称',
  icp_number: 'ICP 备案号（中国大陆运营网站必需），如 沪ICP备XXXXXXXX号',
  footer_tagline: '网站 Footer 区域显示的标语/简介，用于传达品牌核心价值主张',
};

const form = ref<SiteConfig>(defaultForm());
const loading = ref(true);
const saving = ref(false);

interface FieldDef {
  key: string;
  label: string;
  placeholder?: string;
  textarea?: boolean;
  rows?: number;
  image?: boolean;
  sizeHint?: string;
  previewRatio?: string;
  context?: string;
  tip: string;
}

interface GroupDef {
  key: string;
  label: string;
  fields: FieldDef[];
}

const groups: GroupDef[] = [
  {
    key: 'basic', label: '基础信息',
    fields: [
      { key: 'site_name', label: '网站名称', placeholder: '北极星移民', tip: tips.site_name },
      { key: 'site_logo', label: '网站 Logo', image: true, placeholder: '/images/logo.png', sizeHint: '推荐 SVG 或 200×60px PNG（透明背景）', previewRatio: '16 / 5', context: 'general', tip: tips.site_logo },
      { key: 'site_favicon', label: 'Favicon', image: true, placeholder: '/favicon.ico', sizeHint: '推荐 32×32px PNG 或 ICO', previewRatio: '1 / 1', context: 'favicon', tip: tips.site_favicon },
    ],
  },
  {
    key: 'seo', label: 'SEO 设置',
    fields: [
      { key: 'seo_title', label: '标题模板', tip: tips.seo_title },
      { key: 'seo_description', label: 'Meta 描述', tip: tips.seo_description },
      { key: 'seo_keywords', label: '关键词', tip: tips.seo_keywords },
      { key: 'og_image', label: 'OG 分享图', image: true, sizeHint: '推荐 1200×630px (1.91:1)', previewRatio: '1.91 / 1', context: 'og-image', tip: tips.og_image },
      { key: 'canonical_base', label: '首选域名', placeholder: 'https://www.example.com', tip: tips.canonical_base },
    ],
  },
  {
    key: 'geo', label: 'GEO / 结构化数据',
    fields: [
      { key: 'organization_name', label: '机构名称', tip: tips.organization_name },
      { key: 'organization_description', label: '机构描述', tip: tips.organization_description },
      { key: 'organization_logo', label: '机构 Logo', image: true, sizeHint: '推荐 SVG 或 200×60px PNG（透明背景）', previewRatio: '16 / 5', context: 'general', tip: tips.organization_logo },
      { key: 'organization_url', label: '官网 URL', placeholder: 'https://www.example.com', tip: tips.organization_url },
      { key: 'same_as', label: '社交媒体链接', tip: tips.same_as },
    ],
  },
  {
    key: 'contact', label: '联系方式',
    fields: [
      { key: 'contact_phone', label: '客服电话', placeholder: '400-xxx-xxxx', tip: tips.contact_phone },
      { key: 'contact_phone_2', label: '联系电话', placeholder: '400-xxx-xxxx', tip: tips.contact_phone_2 },
      { key: 'contact_email', label: '客服邮箱', placeholder: 'info@example.com', tip: tips.contact_email },
      { key: 'contact_address', label: '公司地址', tip: tips.contact_address },
      { key: 'contact_wechat', label: '微信', image: true, sizeHint: '推荐 500×500px (1:1 正方形)', previewRatio: '1 / 1', context: 'qr-code', tip: tips.contact_wechat },
      { key: 'contact_wechat_mp', label: '微信公众号', image: true, sizeHint: '推荐 500×500px (1:1 正方形)', previewRatio: '1 / 1', context: 'qr-code', tip: tips.contact_wechat_mp },
      { key: 'contact_wechat_video', label: '企业视频号', image: true, sizeHint: '推荐 500×500px (1:1 正方形)', previewRatio: '1 / 1', context: 'qr-code', tip: tips.contact_wechat_video },
    ],
  },
  {
    key: 'third_party', label: '第三方代码',
    fields: [
      { key: 'ga_tracking_id', label: 'Google Analytics ID', placeholder: 'G-XXXXXXXX', tip: tips.ga_tracking_id },
      { key: 'baidu_tongji_id', label: '百度统计 ID', tip: tips.baidu_tongji_id },
      { key: 'custom_head_code', label: '自定义 Head', textarea: true, rows: 6, tip: tips.custom_head_code },
      { key: 'custom_body_code', label: '自定义 Body', textarea: true, rows: 6, tip: tips.custom_body_code },
    ],
  },
  {
    key: 'footer', label: '页脚设置',
    fields: [
      { key: 'copyright_text', label: '版权声明', tip: tips.copyright_text },
      { key: 'icp_number', label: 'ICP 备案号', placeholder: '沪ICP备XXXXXXXX号', tip: tips.icp_number },
      { key: 'footer_tagline', label: '页脚标语', tip: tips.footer_tagline },
    ],
  },
];

const load = async () => {
  loading.value = true;
  try {
    const api = useApi();
    const data = await api<SiteConfig>('/admin/site-config');
    if (data) {
      form.value = { ...defaultForm(), ...data };
    }
  } finally {
    loading.value = false;
  }
};

const save = async () => {
  saving.value = true;
  try {
    const api = useApi();
    await api('/admin/site-config', {
      method: 'PUT',
      body: JSON.parse(JSON.stringify(form.value)),
    });
    notify.success('设置已保存');
  } catch (e) {
    notify.error(e, '保存失败');
  } finally {
    saving.value = false;
  }
};

const addSameAs = () => {
  if (!Array.isArray(form.value.same_as)) form.value.same_as = [];
  form.value.same_as.push('');
};
const removeSameAs = (idx: number) => form.value.same_as.splice(idx, 1);

onMounted(load);
</script>

<style scoped>
.settings-body {
  position: relative;
}

.monospace-input :deep(textarea) {
  font-family: 'Courier New', Courier, monospace;
  font-size: 13px;
}

.array-input {
  flex: 1;
}

.array-row {
  display: flex;
  gap: 8px;
  margin-bottom: 8px;
  align-items: center;
}

.admin-tip-icon {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 20px;
  height: 20px;
  border-radius: 50%;
  background: var(--color-border-light);
  color: var(--color-text-muted);
  cursor: help;
  flex-shrink: 0;
  margin-left: 6px;
  transition: background 0.15s, color 0.15s;
}

.admin-tip-icon:hover {
  background: var(--color-text-muted);
  color: #fff;
}
</style>
