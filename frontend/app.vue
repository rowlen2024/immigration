<template>
  <NuxtLayout>
    <NuxtPage />
  </NuxtLayout>
</template>

<script setup lang="ts">
import { buildOrganizationJsonLd, buildWebSiteJsonLd, toJsonLdConfig, toJsonLdScripts } from '~/utils/jsonld'

const { siteConfig, refreshSiteConfig } = useMygoSiteConfig()
const route = useRoute()

const isPublicPage = computed(() => !route.path.startsWith('/admin'))

useHead(() => {
  const sc = siteConfig.value
  const links: Record<string, unknown>[] = []
  const scripts: Record<string, unknown>[] = []

  // ── Favicon ──
  if (sc?.site_favicon) {
    links.push({ rel: 'icon', type: 'image/x-icon', href: sc.site_favicon })
  }

  // ── Google Analytics ──
  if (sc?.ga_tracking_id) {
    scripts.push(
      { async: true, src: `https://www.googletagmanager.com/gtag/js?id=${sc.ga_tracking_id}` },
      { innerHTML: `window.dataLayer=window.dataLayer||[];function gtag(){dataLayer.push(arguments);}gtag('js',new Date());gtag('config','${sc.ga_tracking_id}');` },
    )
  }

  // ── 百度统计 ──
  if (sc?.baidu_tongji_id) {
    scripts.push({
      innerHTML: `var _hmt=_hmt||[];(function(){var hm=document.createElement("script");hm.src="https://hm.baidu.com/hm.js?${sc.baidu_tongji_id}";var s=document.getElementsByTagName("script")[0];s.parentNode.insertBefore(hm,s);})();`,
    })
  }

  // ── 自定义 head 代码 ──
  if (sc?.custom_head_code) {
    scripts.push({ innerHTML: sc.custom_head_code })
  }

  // ── 自定义 body 尾代码 ──
  if (sc?.custom_body_code) {
    scripts.push({ innerHTML: sc.custom_body_code, tagPosition: 'bodyClose' })
  }

  // ── JSON-LD: Organization + WebSite（仅公开页） ──
  if (isPublicPage.value) {
    const config = toJsonLdConfig(sc)
    scripts.push(...toJsonLdScripts(
      buildOrganizationJsonLd(config),
      buildWebSiteJsonLd(config),
    ))
  }

  return { link: links, script: scripts }
})

onMounted(() => {
  refreshSiteConfig()
})
</script>
