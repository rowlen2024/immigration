<template>
  <footer class="site-footer">
    <div class="footer-container" :style="footerGridStyle">
      <div class="footer-column footer-brand">
        <NuxtLink to="/" class="footer-logo">
          <span v-if="!siteConfig?.site_logo" class="footer-logo-mark">M</span>
          <img v-if="siteConfig?.site_logo" :src="siteConfig.site_logo" :alt="siteConfig?.site_name || '北极星移民'" class="footer-logo-img" />
          <span class="footer-logo-text">{{ siteConfig?.site_name || '北极星移民' }}</span>
        </NuxtLink>
        <p class="footer-tagline">{{ siteConfig?.footer_tagline || '为高净值人群提供美国EB-5、香港投资、巴拿马购房三大精品移民路径的深度决策支持。' }}</p>

        <ul class="footer-brand-contact">
          <li v-if="siteConfig?.contact_phone" class="brand-contact-item">
            <svg class="bc-icon" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M3 18v-6a9 9 0 0 1 18 0v6"/><path d="M21 19a2 2 0 0 1-2 2h-1a2 2 0 0 1-2-2v-3a2 2 0 0 1 2-2h3zM3 19a2 2 0 0 0 2 2h1a2 2 0 0 0 2-2v-3a2 2 0 0 0-2-2H3z"/></svg>
            <span class="brand-contact-label">客服</span>
            <a :href="`tel:${siteConfig.contact_phone}`" class="footer-link">{{ siteConfig.contact_phone }}</a>
          </li>
          <li v-if="siteConfig?.contact_phone_2" class="brand-contact-item">
            <svg class="bc-icon" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M22 16.92v3a2 2 0 0 1-2.18 2 19.79 19.79 0 0 1-8.63-3.07 19.5 19.5 0 0 1-6-6 19.79 19.79 0 0 1-3.07-8.67A2 2 0 0 1 4.11 2h3a2 2 0 0 1 2 1.72 12.84 12.84 0 0 0 .7 2.81 2 2 0 0 1-.45 2.11L8.09 9.91a16 16 0 0 0 6 6l1.27-1.27a2 2 0 0 1 2.11-.45 12.84 12.84 0 0 0 2.81.7A2 2 0 0 1 22 16.92z"/></svg>
            <span class="brand-contact-label">电话</span>
            <a :href="`tel:${siteConfig.contact_phone_2}`" class="footer-link">{{ siteConfig.contact_phone_2 }}</a>
          </li>
          <li v-if="siteConfig?.contact_email" class="brand-contact-item">
            <svg class="bc-icon" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M4 4h16c1.1 0 2 .9 2 2v12c0 1.1-.9 2-2 2H4c-1.1 0-2-.9-2-2V6c0-1.1.9-2 2-2z"/><polyline points="22,6 12,13 2,6"/></svg>
            <a :href="`mailto:${siteConfig.contact_email}`" class="footer-link">{{ siteConfig.contact_email }}</a>
          </li>
          <li v-if="siteConfig?.contact_address" class="brand-contact-item">
            <svg class="bc-icon" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M21 10c0 7-9 13-9 13s-9-6-9-13a9 9 0 0 1 18 0z"/><circle cx="12" cy="10" r="3"/></svg>
            <span>{{ siteConfig.contact_address }}</span>
          </li>
        </ul>

        <div v-if="hasQRCodes" class="footer-qr-row">
          <div v-if="siteConfig?.contact_wechat" class="footer-qr-item">
            <img :src="siteConfig.contact_wechat" alt="微信咨询" class="footer-qr-img" />
            <span class="footer-qr-label">微信咨询</span>
          </div>
          <div v-if="siteConfig?.contact_wechat_mp" class="footer-qr-item">
            <img :src="siteConfig.contact_wechat_mp" alt="公众号" class="footer-qr-img" />
            <span class="footer-qr-label">公众号</span>
          </div>
          <div v-if="siteConfig?.contact_wechat_video" class="footer-qr-item">
            <img :src="siteConfig.contact_wechat_video" alt="视频号" class="footer-qr-img" />
            <span class="footer-qr-label">视频号</span>
          </div>
        </div>
      </div>

      <div v-for="col in footerNav" :key="col.id" class="footer-column">
        <h3 class="footer-heading">{{ col.label }}</h3>
        <ul class="footer-links">
          <li v-for="child in col.children" :key="child.id">
            <NuxtLink :to="child.link" class="footer-link">{{ child.label }}</NuxtLink>
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
const { siteConfig } = useSiteConfig();
const { navItems: footerNav } = useNavigation('footer');

const footerGridStyle = computed(() => {
  const count = footerNav.value.length || 2;
  return { '--footer-cols': `2fr ${Array(count).fill('1fr').join(' ')}` };
});

const hasQRCodes = computed(() => {
  const c = siteConfig.value;
  return !!(c?.contact_wechat || c?.contact_wechat_mp || c?.contact_wechat_video);
});

const copyrightText = computed(() => {
  const template = siteConfig.value?.copyright_text || '© {year} {site_name}. All rights reserved.';
  return template
    .replace('{year}', String(new Date().getFullYear()))
    .replace('{site_name}', siteConfig.value?.site_name || '北极星移民');
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
  grid-template-columns: var(--footer-cols, 1.5fr 1fr 1fr);
  gap: 16px;
}

.footer-brand {
  display: flex;
  flex-direction: column;
  gap: 14px;
  padding-right: 48px;
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

.footer-qr-row {
  display: flex;
  gap: 16px;
  margin-top: 6px;
  padding-top: 10px;
  border-top: 1px solid rgba(200, 150, 62, 0.2);
}

.footer-qr-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 6px;
}

.footer-qr-img {
  width: 80px;
  height: 80px;
  border-radius: 8px;
  border: 1px solid rgba(200, 150, 62, 0.25);
  background: #fff;
  object-fit: contain;
  padding: 4px;
}

.footer-qr-label {
  font-size: 11px;
  color: rgba(255, 255, 255, 0.45);
  white-space: nowrap;
}

.footer-brand-contact {
  display: flex;
  flex-direction: column;
  gap: 6px;
  margin-top: 2px;
  padding-top: 10px;
  border-top: 1px solid rgba(200, 150, 62, 0.2);
}

.brand-contact-item {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 13px;
  color: rgba(255, 255, 255, 0.55);
}

.brand-contact-label {
  font-size: 11px;
  color: rgba(255, 255, 255, 0.35);
  min-width: 28px;
}

.bc-icon {
  color: rgba(200, 150, 62, 0.6);
  flex-shrink: 0;
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
    grid-template-columns: 1fr 1fr !important;
  }
}

@media (max-width: 767px) {
  .footer-container {
    grid-template-columns: 1fr !important;
    gap: 32px;
  }

  .footer-bottom .footer-container {
    flex-direction: column;
    gap: 8px;
    text-align: center;
  }
}
</style>
