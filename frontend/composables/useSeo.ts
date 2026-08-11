import { buildBreadcrumbListJsonLd, toJsonLdScript } from '~/utils/jsonld'
import type { Link, Meta, Script } from '@unhead/vue'
import type { MaybeRefOrGetter } from 'vue'
import type { BreadcrumbItem } from '~/utils/breadcrumb'

interface SeoOptions {
  title?: string
  description?: string
  /** @deprecated 请使用 utils/jsonld.ts 的纯函数 + useHead 替代 */
  jsonLd?: Record<string, unknown>
  /** 页面提供的稳定语义面包屑。 */
  breadcrumbs?: MaybeRefOrGetter<readonly BreadcrumbItem[] | null | undefined>
  robots?: string
}

export const useSeo = (options: SeoOptions) => {
  const { siteConfig } = useMygoSiteConfig()
  const route = useRoute()
  const runtimeConfig = useRuntimeConfig()

  const siteName = computed(() => siteConfig.value?.site_name || '北极星移民')
  const canonicalBase = computed(() => {
    const configuredBase = siteConfig.value?.canonical_base?.trim()
    const fallbackBase = String(runtimeConfig.public.siteUrl || '').trim()
    return (configuredBase || fallbackBase).replace(/\/+$/, '')
  })
  const canonicalUrl = computed(() => `${canonicalBase.value}${route.path}`)

  const fullTitle = computed(() => {
    const pageTitle = options.title?.trim()
    if (pageTitle) {
      return pageTitle.includes(siteName.value)
        ? pageTitle
        : `${pageTitle} - ${siteName.value}`
    }
    if (route.path !== '/') return siteName.value

    const homeTitle = siteConfig.value?.seo_title?.trim() || '{site_name} - 专业投资移民服务'
    return homeTitle.replace('{site_name}', siteName.value)
  })

  const head = computed(() => {
    const metas: Meta[] = []
    const links: Link[] = []
    const scripts: Script[] = []

    // 百度统计
    metas.push({name: "baidu-site-verification", content: "codeva-spiwx4inGs"})

    // 必应
    metas.push({name: "msvalidate.01", content: "EFC906F85E9EE0CA6B69D87204515B66"})

    // ── Description / OG / Twitter ──
    const desc = options.description || siteConfig.value?.seo_description || ''
    if (desc) {
      metas.push(
        { name: 'description', content: desc },
        { property: 'og:title', content: fullTitle.value },
        { property: 'og:description', content: desc },
        { name: 'twitter:title', content: fullTitle.value },
        { name: 'twitter:description', content: desc },
      )
    }

    // ── Keywords ──
    if (siteConfig.value?.seo_keywords) {
      metas.push({ name: 'keywords', content: siteConfig.value.seo_keywords })
    }

    // ── Open Graph ──
    metas.push(
      { property: 'og:type', content: 'website' },
      { property: 'og:site_name', content: siteName.value },
    )
    metas.push({ property: 'og:url', content: canonicalUrl.value })
    if (siteConfig.value?.og_image) {
      metas.push({ property: 'og:image', content: siteConfig.value.og_image })
    }

    // ── Twitter Card ──
    metas.push({ name: 'twitter:card', content: 'summary_large_image' })
    if (siteConfig.value?.og_image) {
      metas.push({ name: 'twitter:image', content: siteConfig.value.og_image })
    }

    // ── Canonical ──
    links.push({ rel: 'canonical', href: canonicalUrl.value })

    // ── Robots ──
    metas.push({ name: 'robots', content: options.robots || 'index, follow' })

    // ── JSON-LD: custom (兼容旧用法) ──
    if (options.jsonLd) {
      scripts.push(toJsonLdScript(options.jsonLd)!)
    }

    // ── JSON-LD: BreadcrumbList ──
    const breadcrumbs = toValue(options.breadcrumbs) ?? []
    if (breadcrumbs.length > 1) {
      const breadcrumbLd = buildBreadcrumbListJsonLd(breadcrumbs, canonicalBase.value)
      const script = toJsonLdScript(breadcrumbLd)
      if (script) scripts.push(script)
    }

    return {
      title: fullTitle.value,
      meta: metas,
      link: links,
      script: scripts,
    }
  })

  useHead(head)
}
