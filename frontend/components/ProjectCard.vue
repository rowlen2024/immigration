<template>
  <NuxtLink :to="link" class="project-card">
    <div class="card-image" :class="`card-image--${variantIdx}`">
      <div class="card-image-glow"></div>
      <div class="card-image-overlay"></div>
      <ResponsiveImage v-if="image" :src="image" :alt="title" variant="sm" :variants="imageVariants" loading="lazy" />
    </div>
    <div class="card-body">
      <h3 class="card-title">{{ title }}</h3>
      <p class="card-desc">{{ description }}</p>
      <dl v-if="keyFacts.length" class="card-facts">
        <div v-for="fact in keyFacts" :key="fact.label" class="card-fact">
          <dt>{{ fact.label }}</dt>
          <dd :title="fact.value">{{ fact.value }}</dd>
        </div>
      </dl>
      <div v-if="audience" class="card-audience">
        <span class="audience-label">适合人群</span>
        <p :title="audience">{{ audience }}</p>
      </div>
      <span class="card-link">
        了解详情
        <span class="link-arrow" v-html="getIconSvg('chevron-right', 14, 'currentColor')"></span>
      </span>
    </div>
    <div class="card-bottom-line"></div>
  </NuxtLink>
</template>

<script setup lang="ts">
import { getIconSvg } from '~/composables/lucideIcons'
import type { ImageVariantInfo } from '~/utils/image'

const props = defineProps<{
  slug: string
  title: string
  description: string
  image: string
  features: string[]
  link: string
  imageVariant?: number
  imageVariants?: Record<string, ImageVariantInfo> | null
}>()

const variantIdx = computed(() => (props.imageVariant ?? 0) % 3)

function splitFeature(feature: string) {
  const separatorIndex = feature.search(/[：:]/)
  if (separatorIndex === -1) return { label: '', value: feature.trim() }

  return {
    label: feature.slice(0, separatorIndex).trim(),
    value: feature.slice(separatorIndex + 1).trim(),
  }
}

const parsedFeatures = computed(() => props.features.filter(Boolean).map(splitFeature))
const keyFacts = computed(() => parsedFeatures.value.filter(feature =>
  feature.label.includes('投资金额') || feature.label.includes('办理周期'),
))
const audience = computed(() => parsedFeatures.value.find(feature => feature.label.includes('适合人群'))?.value || '')
</script>

<style scoped>
.project-card {
  display: flex;
  flex-direction: column;
  height: 100%;
  color: inherit;
  text-decoration: none;
  position: relative;
  background-color: var(--bg-white);
  border-radius: var(--radius-lg);
  overflow: hidden;
  border: 1px solid var(--color-border);
  box-shadow: var(--shadow-sm);
  transition: box-shadow var(--duration-slow) var(--ease-out),
              transform 0.35s var(--ease-spring),
              border-color var(--duration-normal) var(--ease-out);
}

.project-card:hover {
  box-shadow: var(--shadow-xl), 0 0 0 1px rgba(200, 150, 62, 0.25);
  transform: translateY(-6px);
  border-color: rgba(200, 150, 62, 0.4);
}

.project-card:focus-visible {
  outline: 2px solid var(--color-accent);
  outline-offset: 3px;
}

.card-image {
  height: 200px;
  overflow: hidden;
  position: relative;
}

.card-image-overlay {
  position: absolute;
  inset: 0;
  z-index: 1;
  background: linear-gradient(180deg, transparent 50%, rgba(200, 150, 62, 0.25) 100%);
  opacity: 0;
  transition: opacity var(--duration-slow) var(--ease-out);
}

.card-image::after {
  content: '';
  position: absolute;
  bottom: 0;
  left: 0;
  right: 0;
  height: 60px;
  background: linear-gradient(to top, rgba(200, 150, 62, 0.3), transparent);
  z-index: 2;
  transition: height var(--duration-slow) var(--ease-out),
              opacity var(--duration-slow) var(--ease-out);
}

.project-card:hover .card-image-overlay {
  opacity: 1;
}

.project-card:hover .card-image::after {
  height: 100px;
  background: linear-gradient(to top, rgba(200, 150, 62, 0.4), transparent 80%);
}

.card-image--0 {
  background: linear-gradient(135deg, #0F1E3D, #15294D);
}

.card-image--1 {
  background: linear-gradient(135deg, #15294D, #1A3A5C);
}

.card-image--2 {
  background: linear-gradient(135deg, #1A3A5C, #1E3A6E);
}

.card-image-glow {
  position: absolute;
  top: -30px;
  right: -20px;
  width: 140px;
  height: 140px;
  background: radial-gradient(circle, rgba(200, 150, 62, 0.12), transparent 70%);
  border-radius: 50%;
  z-index: 1;
  transition: all var(--duration-slow) var(--ease-out);
}

.project-card:hover .card-image-glow {
  width: 220px;
  height: 220px;
  top: -60px;
  right: -60px;
  background: radial-gradient(circle, rgba(200, 150, 62, 0.22), transparent 65%);
}

.card-image :deep(img) {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: transform 0.45s var(--ease-out);
}

.project-card:hover .card-image :deep(img) {
  transform: scale(1.08);
}

.card-body {
  display: flex;
  flex: 1;
  flex-direction: column;
  padding: 22px 24px;
  position: relative;
  z-index: 2;
}

.card-title {
  display: -webkit-box;
  overflow: hidden;
  min-height: 50px;
  font-size: 18px;
  font-weight: 700;
  color: var(--color-text);
  margin-bottom: 8px;
  transition: color var(--duration-normal) var(--ease-out);
  -webkit-box-orient: vertical;
  -webkit-line-clamp: 2;
}

.project-card:hover .card-title {
  color: var(--color-accent-dark);
}

.card-desc {
  font-size: 14px;
  color: var(--color-text-secondary);
  line-height: 1.7;
  margin-bottom: 16px;
  min-height: 48px;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.card-facts {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  margin: 0;
  padding: 14px 0;
  border-top: 1px solid var(--color-border);
  border-bottom: 1px solid var(--color-border);
}

.card-fact {
  min-width: 0;
  padding-right: 14px;
}

.card-fact:only-child {
  grid-column: 1 / -1;
  padding-right: 0;
}

.card-fact + .card-fact {
  padding-right: 0;
  padding-left: 16px;
  border-left: 1px solid var(--color-border);
}

.card-fact dt,
.audience-label {
  display: block;
  margin-bottom: 5px;
  color: var(--color-text-muted);
  font-size: 12px;
  font-weight: 500;
  line-height: 1.4;
}

.card-fact dd {
  overflow: hidden;
  margin: 0;
  color: var(--color-accent-dark);
  font-size: 15px;
  font-weight: 600;
  line-height: 1.5;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.card-audience {
  min-height: 84px;
  padding: 14px 0 12px;
}

.card-audience p {
  display: -webkit-box;
  overflow: hidden;
  margin: 0;
  color: var(--color-text-secondary);
  font-size: 13px;
  line-height: 1.6;
  -webkit-box-orient: vertical;
  -webkit-line-clamp: 2;
}

.card-link {
  display: inline-flex;
  align-items: center;
  gap: 5px;
  font-size: 14px;
  font-weight: 600;
  color: var(--primary);
  margin-top: auto;
  text-decoration: none;
  transition: gap var(--duration-normal) var(--ease-out),
              color var(--duration-normal) var(--ease-out);
}

.card-link:hover {
  gap: 8px;
  color: var(--accent-dark);
}

.link-arrow {
  display: inline-flex;
  align-items: center;
  transition: transform var(--duration-normal) var(--ease-out);
}

.card-link:hover .link-arrow {
  transform: translateX(2px);
}

.card-bottom-line {
  position: absolute;
  bottom: 0;
  left: 0;
  right: 0;
  height: 2px;
  background: var(--gradient-gold);
  transform: scaleX(0);
  transform-origin: left;
  transition: transform 0.45s var(--ease-spring);
  z-index: 3;
}

.project-card:hover .card-bottom-line {
  transform: scaleX(1);
}

@media (max-width: 767px) {
  .card-image {
    height: 180px;
  }

  .card-body {
    padding: 18px 20px;
  }

  .card-title {
    min-height: auto;
  }

  .card-desc {
    min-height: auto;
  }

  .card-audience {
    min-height: 80px;
  }
}

@media (prefers-reduced-motion: reduce) {
  .project-card,
  .card-image-glow,
  .card-image-overlay,
  .card-image::after,
  .card-image :deep(img),
  .card-title,
  .card-link,
  .link-arrow,
  .card-bottom-line {
    transition: none;
  }

  .project-card:hover {
    transform: none;
  }

  .project-card:hover .card-image :deep(img) {
    transform: none;
  }
}
</style>
