<template>
  <div class="public-layout">
    <Header />
    <main class="main-content">
      <slot />
    </main>
    <Footer />
    <ContactSidebar />
  </div>
</template>

<script setup lang="ts">
import { buildOrganizationJsonLd, buildWebSiteJsonLd, toJsonLdConfig, toJsonLdScripts } from '~/utils/jsonld'

const { siteConfig } = useMygoSiteConfig()
const route = useRoute()

const isPublicPage = computed(() => !route.path.startsWith('/admin'))

useHead(() => {
  const sc = siteConfig.value
  const links: any[] = []
  const scripts: any[] = []

  if (sc?.site_favicon) {
    links.push({ rel: 'icon', type: 'image/x-icon', href: sc.site_favicon })
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

  return { link: links, script: scripts }
})
</script>

<style scoped>
.public-layout {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

.main-content {
  flex: 1;
  margin-top: var(--header-height);
}
</style>
