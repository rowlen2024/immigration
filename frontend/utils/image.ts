export type ImageVariant = 'thumb' | 'sm' | 'md' | 'lg'

// 后端返回的变体信息
export interface ImageVariantInfo {
  url: string
  width: number
}

const variantExtRegex = /\.(jpe?g|png|webp|gif)$/i

/** 兼容旧用法：盲拼变体 URL（当没有 variants 数据时作为 fallback） */
export function getVariantUrl(url: string | null | undefined, variant: ImageVariant): string {
  if (!url) return ''
  return url.replace(variantExtRegex, `_${variant}.jpg`)
}

/** 从后端返回的 variants 数据生成 width-descriptor srcset */
export function getSrcset(variants: Record<string, ImageVariantInfo> | null | undefined): string {
  if (!variants) return ''
  return Object.values(variants)
    .sort((a, b) => a.width - b.width)
    .map(v => `${v.url} ${v.width}w`)
    .join(', ')
}

/** 从 variants 数据中精确取某个变体的 URL（先精确取值，没有则 fallback 到盲拼） */
export function getVariantFromData(
  variants: Record<string, ImageVariantInfo> | null | undefined,
  variant: ImageVariant,
  fallbackUrl?: string | null,
): string | undefined {
  if (variants?.[variant]) return variants[variant].url
  if (fallbackUrl) return getVariantUrl(fallbackUrl, variant)
  return undefined
}
