// ═══════════════════════════════════════════════
// JSON-LD 结构化数据纯函数工厂
// 无 Vue 依赖，可独立单测、跨项目复用
// ═══════════════════════════════════════════════

import { stripHtml } from './html'

// ── 共享类型 ──────────────────────────────────

/** SiteConfig 中 JSON-LD 生成所需字段（子集） */
export interface JsonLdSiteConfig {
  siteName?: string | null
  canonicalBase?: string | null
  organizationName?: string | null
  organizationUrl?: string | null
  organizationDescription?: string | null
  organizationLogo?: string | null
  sameAs?: (string | null)[] | null
  contactPhone?: string | null
  contactEmail?: string | null
  contactAddress?: string | null
}

export interface ServiceJsonLdInput {
  name: string
  description: string
  category?: string
  url?: string
  image?: string
  investmentAmount?: string
  avgRating?: number | null
  reviewCount?: number
  providerName?: string
}

export interface ArticleJsonLdInput {
  headline: string
  description: string
  url?: string
  image?: string
  datePublished?: string
  authorName?: string
  publisherName?: string
}

export interface FaqItem {
  question: string
  answer: string
}

export interface BreadcrumbItem {
  label: string
  link?: string
}

// ── 工具函数 ──────────────────────────────────

const context = (type: string, body: Record<string, unknown>) => ({
  '@context': 'https://schema.org',
  '@type': type,
  ...body,
})

// ── 构建函数 ──────────────────────────────────

/** WebSite + SearchAction */
export function buildWebSiteJsonLd(config: JsonLdSiteConfig): Record<string, unknown> | null {
  const name = config.siteName
  if (!name) return null

  const siteUrl = config.canonicalBase || config.organizationUrl || ''

  const result: Record<string, unknown> = {
    '@context': 'https://schema.org',
    '@type': 'WebSite',
    name,
    url: siteUrl,
  }

  if (siteUrl) {
    result.potentialAction = {
      '@type': 'SearchAction',
      target: {
        '@type': 'EntryPoint',
        urlTemplate: `${siteUrl}/search?q={search_term_string}`,
      },
      'query-input': 'required name=search_term_string',
    }
  }

  return result
}

/** Organization + ContactPoint + PostalAddress */
export function buildOrganizationJsonLd(config: JsonLdSiteConfig): Record<string, unknown> | null {
  if (!config.organizationName) return null

  const org: Record<string, unknown> = context('Organization', {
    name: config.organizationName,
    url: config.organizationUrl || '',
  })

  if (config.organizationDescription) org.description = config.organizationDescription
  if (config.organizationLogo) org.logo = config.organizationLogo

  if (config.sameAs?.length) {
    const filtered = config.sameAs.filter((s): s is string => typeof s === 'string' && s.trim() !== '')
    if (filtered.length) org.sameAs = filtered
  }

  if (config.contactPhone || config.contactEmail) {
    org.contactPoint = {
      '@type': 'ContactPoint',
      telephone: config.contactPhone || '',
      email: config.contactEmail || '',
      contactType: 'customer service',
    }
  }

  if (config.contactAddress) {
    org.address = {
      '@type': 'PostalAddress',
      streetAddress: config.contactAddress,
    }
  }

  return org
}

/** Service（移民项目） */
export function buildServiceJsonLd(
  input: ServiceJsonLdInput,
  config: JsonLdSiteConfig,
): Record<string, unknown> | null {
  if (!input.name) return null

  const service: Record<string, unknown> = context('Service', {
    name: input.name,
    description: stripHtml(input.description, 300),
    category: input.category || '移民服务',
    provider: {
      '@type': 'Organization',
      name: input.providerName || config.organizationName || '北极星移民',
    },
  })

  if (input.image) service.image = input.image
  if (input.url) service.url = input.url

  if (input.investmentAmount) {
    service.offers = {
      '@type': 'Offer',
      description: `投资金额: ${input.investmentAmount}`,
      availability: 'https://schema.org/InStock',
    }
  }

  if (input.avgRating != null && input.reviewCount != null && input.reviewCount > 0) {
    service.aggregateRating = {
      '@type': 'AggregateRating',
      ratingValue: input.avgRating,
      reviewCount: input.reviewCount,
      bestRating: '5',
      worstRating: '1',
    }
  }

  return service
}

/** Article（成功案例详情） */
export function buildArticleJsonLd(
  input: ArticleJsonLdInput,
  config: JsonLdSiteConfig,
): Record<string, unknown> | null {
  if (!input.headline) return null

  const orgName = config.organizationName || '北极星移民'
  const article: Record<string, unknown> = context('Article', {
    headline: input.headline,
    description: input.description,
    author: { '@type': 'Organization', name: input.authorName || orgName },
    publisher: { '@type': 'Organization', name: input.publisherName || orgName },
  })

  if (input.url) article.url = input.url
  if (input.image) article.image = input.image
  if (input.datePublished) article.datePublished = input.datePublished

  return article
}

/** FAQPage */
export function buildFAQPageJsonLd(items: FaqItem[]): Record<string, unknown> | null {
  if (!items.length) return null

  return context('FAQPage', {
    mainEntity: items.map((faq) => ({
      '@type': 'Question',
      name: faq.question,
      acceptedAnswer: { '@type': 'Answer', text: stripHtml(faq.answer, 500) },
    })),
  })
}

/** BreadcrumbList */
export function buildBreadcrumbListJsonLd(
  items: BreadcrumbItem[],
  baseUrl: string,
): Record<string, unknown> | null {
  if (!items.length) return null

  return context('BreadcrumbList', {
    itemListElement: items.map((item, index) => ({
      '@type': 'ListItem',
      position: index + 1,
      name: item.label,
      item: item.link ? baseUrl + item.link : undefined,
    })),
  })
}

/** AboutPage */
export function buildAboutPageJsonLd(
  name: string,
  description: string,
  aboutOrg?: Record<string, unknown>,
): Record<string, unknown> {
  const page: Record<string, unknown> = context('AboutPage', { name, description })
  if (aboutOrg) page.about = aboutOrg
  return page
}

// ── 辅助：将 JSON-LD 对象转为 Nuxt useHead script 条目 ──

/** 从站点配置原始对象提取 JSON-LD 所需字段 */
export function toJsonLdConfig(raw: unknown): JsonLdSiteConfig {
  const r = raw as Record<string, unknown> | null | undefined
  if (!r) return {}
  return {
    siteName: (r.site_name as string) ?? null,
    canonicalBase: (r.canonical_base as string) ?? null,
    organizationName: (r.organization_name as string) ?? null,
    organizationUrl: (r.organization_url as string) ?? null,
    organizationDescription: (r.organization_description as string) ?? null,
    organizationLogo: (r.organization_logo as string) ?? null,
    sameAs: (r.same_as as (string | null)[]) ?? null,
    contactPhone: (r.contact_phone as string) ?? null,
    contactEmail: (r.contact_email as string) ?? null,
    contactAddress: (r.contact_address as string) ?? null,
  }
}

export function toJsonLdScript(ld: Record<string, unknown> | null): { type: string; innerHTML: string } | null {
  if (!ld) return null
  return { type: 'application/ld+json', innerHTML: JSON.stringify(ld) }
}

/** 批量转换，过滤 null */
export function toJsonLdScripts(...items: (Record<string, unknown> | null)[]): { type: string; innerHTML: string }[] {
  return items.map(toJsonLdScript).filter((s): s is { type: string; innerHTML: string } => s !== null)
}
