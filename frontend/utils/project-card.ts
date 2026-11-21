import type { ImageVariantInfo } from '~/utils/image'

export interface ProjectCardSource {
  slug: string
  name: string
  tagline?: string
  overview_text?: string
  cover_image?: string
  cover_image_variants?: Record<string, ImageVariantInfo>
  investment_amount?: string
  processing_period?: string
  target_crowd?: string
}

export interface ProjectCardItem {
  slug: string
  title: string
  description: string
  image: string
  imageVariants?: Record<string, ImageVariantInfo>
  features: string[]
  link: string
}

export function toProjectCard(project: ProjectCardSource): ProjectCardItem {
  const features: string[] = []
  if (project.investment_amount) features.push(`投资金额：${project.investment_amount}`)
  if (project.processing_period) features.push(`办理周期：${project.processing_period}`)
  if (project.target_crowd) features.push(`适合人群：${project.target_crowd}`)
  return {
    slug: project.slug,
    title: project.name,
    description: project.tagline || project.overview_text || '',
    image: project.cover_image || '',
    imageVariants: project.cover_image_variants,
    features,
    link: `/projects/${project.slug}`,
  }
}
