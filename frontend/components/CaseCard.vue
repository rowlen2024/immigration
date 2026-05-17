<template>
  <NuxtLink :to="`/case/${slug}`" class="case-card">
    <div class="case-image">
      <div class="case-image-overlay"></div>
      <div class="case-image-goldline"></div>
      <img
        :src="image || ''"
        :alt="name"
        loading="lazy"
      />
    </div>
    <div class="case-body">
      <div class="case-title-row">
        <span class="case-title-bar"></span>
        <h3 class="case-name">{{ name }}</h3>
      </div>
      <div v-if="metaText" class="case-meta-text">{{ metaText }}</div>
      <div v-else class="case-meta">
        <span class="case-country">{{ country }}</span>
        <span v-if="project" class="case-project">{{ project }}</span>
      </div>
      <p v-if="summary" class="case-desc">{{ summary }}</p>
      <div v-if="showResult" class="case-result">
        <span class="result-badge">{{ resultText }}</span>
      </div>
    </div>
  </NuxtLink>
</template>

<script setup lang="ts">
defineProps<{
  slug: string;
  name: string;
  country: string;
  image?: string;
  project?: string;
  summary?: string;
  showResult?: boolean;
  resultText?: string;
  metaText?: string;
}>();
</script>

<style scoped>
.case-card {
  display: block;
  text-decoration: none;
  color: inherit;
  background-color: var(--bg-white);
  border-radius: var(--radius-lg);
  overflow: hidden;
  border: 1px solid var(--color-border);
  box-shadow: var(--shadow-sm);
  transition: box-shadow 0.35s var(--ease-out),
              transform 0.35s var(--ease-spring),
              border-color 0.3s var(--ease-out);
}

.case-card:hover {
  box-shadow: var(--shadow-lg), 0 0 0 1px rgba(200, 150, 62, 0.15);
  transform: translateY(-4px);
  border-color: rgba(200, 150, 62, 0.3);
}

/* ── Image area ── */

.case-image {
  height: 200px;
  overflow: hidden;
  background: linear-gradient(135deg, #0F1E3D, #1A3A5C);
  position: relative;
}

.case-image-overlay {
  position: absolute;
  inset: 0;
  background: linear-gradient(180deg, rgba(15, 30, 61, 0.5) 0%, transparent 60%);
  opacity: 0;
  transition: opacity 0.35s var(--ease-out);
  z-index: 1;
}

.case-card:hover .case-image-overlay {
  opacity: 1;
}

.case-image-goldline {
  position: absolute;
  top: 0;
  left: 20px;
  right: 20px;
  height: 2px;
  background: var(--gradient-gold);
  opacity: 0;
  transform: scaleX(0);
  transform-origin: left;
  transition: opacity 0.2s var(--ease-out),
              transform 0.45s var(--ease-spring);
  z-index: 2;
}

.case-card:hover .case-image-goldline {
  opacity: 1;
  transform: scaleX(1);
}

.case-image img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: transform 0.45s var(--ease-out);
}

.case-card:hover .case-image img {
  transform: scale(1.06);
}

/* ── Body ── */

.case-body {
  padding: 20px;
  position: relative;
}

.case-title-row {
  display: flex;
  align-items: flex-start;
  gap: 10px;
  margin-bottom: 10px;
}

.case-title-bar {
  display: inline-block;
  width: 3px;
  height: 22px;
  background: var(--gradient-gold);
  border-radius: 2px;
  flex-shrink: 0;
  margin-top: 1px;
  opacity: 0;
  transform: scaleY(0);
  transition: opacity 0.25s var(--ease-out),
              transform 0.35s var(--ease-spring);
}

.case-card:hover .case-title-bar {
  opacity: 1;
  transform: scaleY(1);
}

.case-meta-text {
  font-size: 13px;
  color: var(--color-text-muted);
  margin-bottom: 8px;
}

.case-meta {
  display: flex;
  gap: 8px;
  margin-bottom: 12px;
}

.case-country,
.case-project {
  font-size: 12px;
  font-weight: 600;
  padding: 2px 10px;
  border-radius: var(--radius-full);
  transition: background var(--duration-normal) var(--ease-out),
              color var(--duration-normal) var(--ease-out);
}

.case-country {
  background-color: rgba(15, 27, 45, 0.08);
  color: var(--color-primary);
}

.case-card:hover .case-country {
  background-color: rgba(15, 27, 45, 0.14);
  color: var(--color-primary-dark);
}

.case-project {
  background-color: var(--color-accent-soft);
  color: var(--color-accent-dark);
}

.case-card:hover .case-project {
  background-color: rgba(200, 150, 62, 0.18);
  color: var(--color-accent-dark);
}

.case-name {
  font-size: 18px;
  font-weight: 700;
  color: var(--color-text);
  margin-bottom: 0;
  transition: color 0.3s var(--ease-out);
}

.case-card:hover .case-name {
  color: var(--color-primary-dark);
}

.case-desc {
  font-size: 14px;
  color: var(--color-text-secondary);
  line-height: 1.7;
  margin-bottom: 14px;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  transition: transform 0.3s var(--ease-out);
}

.case-card:hover .case-desc {
  transform: translateY(-2px);
}

.case-result {
  display: flex;
}

.result-badge {
  font-size: 13px;
  font-weight: 600;
  color: var(--color-success);
  background-color: var(--color-success-soft);
  padding: 4px 14px;
  border-radius: var(--radius-full);
  transition: background var(--duration-normal) var(--ease-out);
}

.case-card:hover .result-badge {
  background-color: #c6f6d5;
}

@media (max-width: 767px) {
  .case-image {
    height: 180px;
  }
}
</style>
