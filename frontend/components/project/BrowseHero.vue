<template>
  <header class="project-browse-hero" :class="{ 'project-browse-hero--illustrated': images.length }">
    <div v-if="images.length" class="project-browse-hero__media" aria-hidden="true">
      <ResponsiveImage
        v-for="(image, index) in images"
        :key="image.src"
        :src="image.src"
        :variants="image.variants"
        :sizes="images.length > 1 ? '(max-width: 640px) 50vw, 22vw' : '(max-width: 640px) 85vw, 62vw'"
        variant="md"
        alt=""
        loading="eager"
        :fetchpriority="index === 0 ? 'high' : 'auto'"
      />
    </div>

    <div class="container project-browse-hero__content">
      <nav v-if="breadcrumbs.length" class="project-browse-hero__trail" aria-label="当前位置">
        <ol>
          <li v-for="(item, index) in breadcrumbs" :key="item.href">
            <span v-if="index === breadcrumbs.length - 1" aria-current="page">{{ item.label }}</span>
            <NuxtLink v-else :to="item.href">{{ item.label }}</NuxtLink>
          </li>
        </ol>
      </nav>
      <h1>{{ title }}</h1>
      <p v-if="description" class="project-browse-hero__description">{{ description }}</p>
      <div v-if="$slots.default" class="project-browse-hero__actions">
        <slot />
      </div>
    </div>
  </header>
</template>

<script setup lang="ts">
import type { BreadcrumbItem } from '~/utils/breadcrumb'
import type { ImageVariantInfo } from '~/utils/image'

defineProps<{
  title: string
  description: string
  breadcrumbs: readonly BreadcrumbItem[]
  images: readonly { src: string; variants?: Record<string, ImageVariantInfo> }[]
}>()
</script>

<style scoped>
.project-browse-hero {
  position: relative;
  isolation: isolate;
  background: var(--gradient-hero);
  color: var(--color-bg-surface);
}

.project-browse-hero__media {
  position: absolute;
  z-index: -2;
  inset: 0 0 0 auto;
  display: flex;
  width: 62%;
  overflow: hidden;
}

.project-browse-hero__media img {
  flex: 1;
  min-width: 0;
  height: 100%;
  object-fit: cover;
}

.project-browse-hero--illustrated::after {
  position: absolute;
  z-index: -1;
  inset: 0;
  background: linear-gradient(90deg, var(--color-primary) 31%, rgb(15 27 45 / 94%) 41%, rgb(15 27 45 / 46%) 70%, rgb(15 27 45 / 18%));
  content: '';
}

.project-browse-hero__content {
  padding-block: clamp(20px, 3vw, 34px) clamp(28px, 5vw, 52px);
}

.project-browse-hero__trail {
  margin-bottom: clamp(10px, 2vw, 18px);
}

.project-browse-hero__trail ol {
  display: flex;
  flex-wrap: wrap;
  gap: 4px 10px;
  font-size: var(--text-sm);
}

.project-browse-hero__trail li {
  display: flex;
  align-items: baseline;
  gap: 10px;
  min-width: 0;
  overflow-wrap: anywhere;
}

.project-browse-hero__trail li + li::before {
  color: rgb(255 255 255 / 70%);
  content: '/';
}

.project-browse-hero__trail a {
  color: var(--color-accent-soft);
}

.project-browse-hero__trail a:hover {
  text-decoration: underline;
}

.project-browse-hero h1 {
  max-width: 18ch;
  font: 700 clamp(30px, 4vw, var(--text-4xl)) / 1.25 var(--font-serif);
  overflow-wrap: anywhere;
}

.project-browse-hero__description {
  max-width: 34em;
  margin-top: 10px;
  color: #edf1f7;
  white-space: pre-line;
  overflow-wrap: anywhere;
}

.project-browse-hero__actions {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 8px clamp(12px, 2vw, 24px);
  margin-top: 18px;
}

.project-browse-hero :deep(:focus-visible) {
  outline: 2px solid var(--color-accent-hover);
  outline-offset: 4px;
}

@media (max-width: 767px) {
  .project-browse-hero__media {
    width: 90%;
  }

  .project-browse-hero--illustrated::after {
    background: linear-gradient(90deg, var(--color-primary) 8%, rgb(15 27 45 / 88%) 48%, rgb(15 27 45 / 63%));
  }
}

@media (max-width: 480px) {
  .project-browse-hero__media img:nth-child(n + 3) {
    display: none;
  }
}
</style>
