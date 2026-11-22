<template>
  <div class="public-layout" :class="{ 'navigation-open': mobileMenuOpen }">
    <Header />
    <main class="main-content" :inert="mobileMenuOpen">
      <slot />
    </main>
    <Footer :inert="mobileMenuOpen" />
    <ContactSidebar />
  </div>
</template>

<script setup lang="ts">
import { buildOrganizationJsonLd, buildWebSiteJsonLd, toJsonLdConfig, toJsonLdScripts } from '~/utils/jsonld'

const { siteConfig } = useMygoSiteConfig()
const mobileMenuOpen = useState('public-mobile-menu-open', () => false)
const route = useRoute()

const isPublicPage = computed(() => !route.path.startsWith('/admin'))

useHead(() => {
  const sc = siteConfig.value
  const links: any[] = []
  const scripts: any[] = []

  if (sc?.site_favicon) {
    links.push({ rel: 'icon', href: sc.site_favicon })
  }

  if (sc?.ga_tracking_id) {
    scripts.push(
      { async: true, src: `https://www.googletagmanager.com/gtag/js?id=${sc.ga_tracking_id}` },
      { innerHTML: `window.dataLayer=window.dataLayer||[];function gtag(){dataLayer.push(arguments);}gtag('js',new Date());gtag('config','${sc.ga_tracking_id}');` },
    )
  }

  if (sc?.baidu_tongji_id) {
    scripts.push({
      innerHTML: `var _hmt=_hmt||[];(function(){var hm=document.createElement("script");hm.src="https://hm.baidu.com/hm.js?${sc.baidu_tongji_id}";var s=document.getElementsByTagName("script")[0];s.parentNode.insertBefore(hm,s);})();`,
    })
  }

  if (sc?.custom_head_code) {
    scripts.push({ innerHTML: sc.custom_head_code })
  }

  if (sc?.custom_body_code) {
    scripts.push({ innerHTML: sc.custom_body_code, tagPosition: 'bodyClose' })
  }

  if (isPublicPage.value) {
    const config = toJsonLdConfig(sc)
    scripts.push(...toJsonLdScripts(
      buildOrganizationJsonLd(config),
      buildWebSiteJsonLd(config),
    ))
  }

  return { htmlAttrs: { lang: 'zh-CN' }, link: links, script: scripts }
})
</script>

<style>
.public-layout {
  --public-contact-bar-height: 0px;
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  padding-bottom: var(--public-contact-bar-height);
}

.public-layout > .main-content {
  flex: 1;
  margin-top: var(--header-height);
}

html:has(.public-layout.navigation-open) {
  overflow: hidden;
}

@media (max-width: 890px) {
  .public-layout:has(.contact-sidebar) {
    --public-contact-bar-height: calc(49px + env(safe-area-inset-bottom, 0px));
  }

  .public-layout:is(.navigation-open, :has(.public-picker-dialog[open])) .contact-sidebar {
    visibility: hidden;
  }
}

@media (max-width: 767px) {
  .public-layout:has(.project-detail-actions, .compare-consult-bar) {
    --public-contact-bar-height: calc(66px + env(safe-area-inset-bottom, 0px));
  }

  .public-layout:has(.project-detail-actions, .compare-consult-bar) .contact-sidebar {
    display: none;
  }

  .public-layout:is(.navigation-open, :has(.public-picker-dialog[open])) :is(.project-detail-actions, .compare-consult-bar) {
    visibility: hidden;
  }
}
</style>
