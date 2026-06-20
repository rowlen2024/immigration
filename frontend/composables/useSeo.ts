import { buildBreadcrumbListJsonLd, toJsonLdScript } from '~/utils/jsonld'
import type { Link, Meta, Script } from '@unhead/vue'

interface SeoOptions {
  title?: string
  description?: string
  /** @deprecated 请使用 utils/jsonld.ts 的纯函数 + useHead 替代 */
  jsonLd?: Record<string, unknown>
  breadcrumbLabel?: string
  robots?: string
}

export const useSeo = (options: SeoOptions) => {
  const { siteConfig } = useMygoSiteConfig()
  const { getBreadcrumb } = useNavigation()
  const route = useRoute()

  const siteName = computed(() => siteConfig.value?.site_name || '北极星移民')

  const fullTitle = computed(() => {
    if (options.title) {
      return `${options.title} | ${siteName.value}`
    }
    const seoTitle = siteConfig.value?.seo_title || '{site_name} | 专业投资移民服务'
    return seoTitle.replace('{site_name}', siteName.value)
  })

  const head = computed(() => {
    const metas: Meta[] = []
    const links: Link[] = []
    const scripts: Script[] = []

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
    if (siteConfig.value?.canonical_base) {
      metas.push({ property: 'og:url', content: siteConfig.value.canonical_base + route.path })
    }
    if (siteConfig.value?.og_image) {
      metas.push({ property: 'og:image', content: siteConfig.value.og_image })
    }

    // ── Twitter Card ──
    metas.push({ name: 'twitter:card', content: 'summary_large_image' })
    if (siteConfig.value?.og_image) {
      metas.push({ name: 'twitter:image', content: siteConfig.value.og_image })
    }

    // ── Canonical ──
    if (siteConfig.value?.canonical_base) {
      links.push({ rel: 'canonical', href: siteConfig.value.canonical_base + route.path })
    }

    // ── Robots ──
    metas.push({ name: 'robots', content: options.robots || 'index, follow' })

    // ── JSON-LD: custom (兼容旧用法) ──
    if (options.jsonLd) {
      scripts.push(toJsonLdScript(options.jsonLd)!)
    }

    // ── JSON-LD: BreadcrumbList ──
    const crumbs = getBreadcrumb(route.path, options.breadcrumbLabel)
    if (crumbs.length > 0) {
      const base = siteConfig.value?.canonical_base || ''
      // SSR 场景下导航数据可能未加载，此时面包屑标签是原始路径段
      const hasCJK = (s: string) => /[一-鿿]/.test(s)
      const fixedCrumbs = crumbs.map((item, index) => {
        const isLast = index === crumbs.length - 1
        const label = (!hasCJK(item.label) && isLast && options.breadcrumbLabel)
          ? options.breadcrumbLabel
          : item.label
        return { label, link: item.link }
      })
      const breadcrumbLd = buildBreadcrumbListJsonLd(fixedCrumbs, base)
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
