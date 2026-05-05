<template>
  <footer class="site-footer">
    <div class="footer-container">
      <div class="footer-column footer-brand">
        <NuxtLink to="/" class="footer-logo">
          <span v-if="!siteConfig?.site_logo" class="footer-logo-mark">M</span>
          <img v-if="siteConfig?.site_logo" :src="siteConfig.site_logo" :alt="siteConfig?.site_name || '北极星移民'" class="footer-logo-img" />
          <span class="footer-logo-text">{{ siteConfig?.site_name || '北极星移民' }}</span>
        </NuxtLink>
        <p class="footer-tagline">{{ siteConfig?.footer_tagline || '为高净值人群提供美国EB-5、香港投资、巴拿马购房三大精品移民路径的深度决策支持。' }}</p>
      </div>

      <div v-for="col in footerNav" :key="col.id" class="footer-column">
        <h3 class="footer-heading">{{ col.label }}</h3>
        <ul class="footer-links">
          <li v-for="child in col.children" :key="child.id">
            <NuxtLink :to="child.link" class="footer-link">{{ child.label }}</NuxtLink>
          </li>
        </ul>
      </div>

      <div class="footer-column">
        <h3 class="footer-heading">联系方式</h3>
        <ul class="footer-contact">
          <li class="contact-item">
            <span class="contact-label">电话：</span>
            <a v-if="siteConfig?.contact_phone" :href="`tel:${siteConfig.contact_phone}`" class="footer-link">{{ siteConfig.contact_phone }}</a>
            <span v-else class="footer-link">400-xxx-xxxx</span>
          </li>
          <li class="contact-item">
            <span class="contact-label">邮箱：</span>
            <a v-if="siteConfig?.contact_email" :href="`mailto:${siteConfig.contact_email}`" class="footer-link">{{ siteConfig.contact_email }}</a>
            <span v-else class="footer-link">info@mygo-immigration.com</span>
          </li>
          <li class="contact-item">
            <span class="contact-label">地址：</span>
            <span>{{ siteConfig?.contact_address || '上海市浦东新区陆家嘴金融中心' }}</span>
          </li>
          <li class="contact-item">
            <span class="contact-label">微信：</span>
            <span>{{ siteConfig?.contact_wechat || 'MyGo_Immigration' }}</span>
          </li>
        </ul>
      </div>
    </div>

    <div class="footer-bottom">
      <div class="footer-container">
        <p class="copyright">{{ copyrightText }}</p>
        <!--
        <div class="footer-legal">
          <span>隐私政策</span>
          <span>免责声明</span>
          <span>服务条款</span>
        </div>
        -->
      </div>
    </div>
  </footer>
</template>

<script setup lang="ts">
const { siteConfig, fetch: fetchSiteConfig } = useSiteConfig();

interface NavItem {
  id: number;
  label: string;
  link: string;
  children: NavItem[];
  status: boolean;
}

const FALLBACK_FOOTER: NavItem[] = [
  {
    id: 1, label: '移民项目', link: '', status: true,
    children: [
      { id: 11, label: '美国EB-5投资移民', link: '/projects/eb5', children: [], status: true },
      { id: 12, label: '香港投资移民', link: '/projects/cies', children: [], status: true },
      { id: 13, label: '巴拿马购房移民', link: '/projects/panama', children: [], status: true },
      { id: 14, label: '项目对比', link: '/compare', children: [], status: true },
    ],
  },
  {
    id: 2, label: '关于我们', link: '', status: true,
    children: [
      { id: 21, label: '公司简介', link: '/about', children: [], status: true },
      { id: 22, label: '成功案例', link: '/cases', children: [], status: true },
      { id: 23, label: '常见问题', link: '/faq', children: [], status: true },
      { id: 24, label: '联系我们', link: '/contact', children: [], status: true },
    ],
  },
];

const footerNav = ref<NavItem[]>([]);

const fetchFooterNav = async () => {
  try {
    const api = useApi();
    const data = await api<NavItem[]>('/navigation?position=footer');
    if (data && (data as NavItem[]).length > 0) {
      footerNav.value = data as NavItem[];
    } else {
      footerNav.value = FALLBACK_FOOTER;
    }
  } catch {
    footerNav.value = FALLBACK_FOOTER;
  }
};

const copyrightText = computed(() => {
  const template = siteConfig.value?.copyright_text || '© {year} {site_name}. All rights reserved.';
  return template
    .replace('{year}', String(new Date().getFullYear()))
    .replace('{site_name}', siteConfig.value?.site_name || '北极星移民');
});

onMounted(() => {
  fetchFooterNav();
  fetchSiteConfig();
});
</script>

<style scoped>
.site-footer {
  background: linear-gradient(180deg, var(--primary-dark), #0F1E3D);
  color: var(--bg-white);
  padding-top: 56px;
}

.footer-container {
  max-width: var(--max-width);
  margin: 0 auto;
  padding: 0 24px;
  display: grid;
  grid-template-columns: 1.5fr 1fr 1fr 1fr;
  gap: 40px;
}

.footer-brand {
  display: flex;
  flex-direction: column;
  gap: 14px;
}

.footer-logo {
  display: flex;
  align-items: center;
  gap: 10px;
}

.footer-logo-mark {
  width: 30px;
  height: 30px;
  background: linear-gradient(135deg, var(--accent), var(--accent-dark));
  border-radius: 7px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--bg-white);
  font-weight: 800;
  font-size: 12px;
}

.footer-logo-img {
  height: 30px;
  width: auto;
  filter: brightness(10) saturate(0);
}

.footer-logo-text {
  font-size: 18px;
  font-weight: 700;
  color: var(--bg-white);
}

.footer-tagline {
  font-size: 13px;
  color: rgba(255, 255, 255, 0.5);
  line-height: 1.8;
}

.footer-heading {
  font-size: 14px;
  font-weight: 600;
  margin-bottom: 16px;
  color: var(--accent);
}

.footer-links {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.footer-link {
  color: rgba(255, 255, 255, 0.55);
  font-size: 13px;
  transition: color 0.2s ease;
}

.footer-link:hover {
  color: var(--accent-light);
}

.footer-contact {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.contact-item {
  font-size: 13px;
  color: rgba(255, 255, 255, 0.55);
  line-height: 1.8;
}

.contact-label {
  color: rgba(255, 255, 255, 0.4);
}

.footer-bottom {
  margin-top: 44px;
  padding: 20px 0;
  border-top: 1px solid rgba(200, 150, 62, 0.12);
}

.footer-bottom .footer-container {
  display: flex;
  justify-content: space-between;
  align-items: center;
  grid-template-columns: none;
}

.copyright {
  font-size: 12px;
  color: rgba(255, 255, 255, 0.4);
}

.footer-legal {
  display: flex;
  gap: 16px;
  font-size: 12px;
  color: rgba(255, 255, 255, 0.4);
}

@media (max-width: 1023px) {
  .footer-container {
    grid-template-columns: 1fr 1fr;
  }
}

@media (max-width: 767px) {
  .footer-container {
    grid-template-columns: 1fr;
    gap: 32px;
  }

  .footer-bottom .footer-container {
    flex-direction: column;
    gap: 8px;
    text-align: center;
  }
}
</style>
