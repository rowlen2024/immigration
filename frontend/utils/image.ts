export type ImageVariant = 'thumb' | 'sm' | 'md' | 'lg'

const variantExtRegex = /\.(jpe?g|png|webp|gif)$/i

export function getVariantUrl(url: string | null | undefined, variant: ImageVariant): string {
  if (!url) return ''
  return url.replace(variantExtRegex, `_${variant}.jpg`)
}
