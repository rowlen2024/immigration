<template>
  <div v-if="lawyers.length > 0" class="lawyer-carousel">
    <div class="carousel-viewport" ref="viewportRef" @touchstart="onTouchStart" @touchend="onTouchEnd">
      <div
        class="carousel-track"
        :style="{ transform: `translateX(${offset}px)`, gap: `${gap}px` }"
      >
        <div
          v-for="lawyer in lawyers"
          :key="lawyer.id"
          class="lawyer-card"
          :style="{ flex: `0 0 ${cardWidth}px`, width: `${cardWidth}px` }"
        >
          <div class="lawyer-photo">
            <ResponsiveImage
              v-if="lawyer.photo_url"
              class="lawyer-photo-img"
              :src="lawyer.photo_url"
              :alt="lawyer.name"
              variant="sm"
              :variants="lawyer.photo_variants"
              loading="lazy"
              sizes="(max-width: 640px) 120px, 180px"
            />
            <div v-else class="lawyer-photo-placeholder"></div>
          </div>
          <div class="lawyer-body">
            <h3 class="lawyer-name">{{ lawyer.name }}</h3>
            <p class="lawyer-title">{{ lawyer.title }}</p>
            <div class="lawyer-divider"></div>
            <ul class="lawyer-tags">
              <li v-for="tag in lawyer.tags" :key="tag">{{ tag }}</li>
            </ul>
          </div>
        </div>
      </div>
    </div>

    <div v-if="maxPage > 0" class="carousel-dots">
      <button
        v-for="p in maxPage + 1"
        :key="p"
        class="carousel-dot"
        :class="{ active: currentPage === p - 1 }"
        :aria-label="`第 ${p} 页`"
        @click="goToPage(p - 1)"
      ></button>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { ImageVariantInfo } from '~/utils/image'

export interface LawyerItem {
  id: number;
  name: string;
  title: string;
  tags: string[];
  photo_url?: string;
  photo_variants?: Record<string, ImageVariantInfo>;
}

const props = defineProps<{
  lawyers: LawyerItem[];
}>();

const viewportRef = ref<HTMLElement | null>(null);
const currentPage = ref(0);
const viewportWidth = ref(0);
let autoTimer: ReturnType<typeof setInterval> | null = null;
let touchStartX = 0;

function onTouchStart(e: TouchEvent) {
  touchStartX = e.touches[0].clientX;
}

function onTouchEnd(e: TouchEvent) {
  const diff = touchStartX - e.changedTouches[0].clientX;
  if (Math.abs(diff) > 50) {
    if (diff > 0) {
      goToPage(currentPage.value + 1);
    } else {
      goToPage(currentPage.value - 1);
    }
  }
}

const cardsPerView = computed(() => {
  if (viewportWidth.value < 640) return 1;
  return 2;
});

const gap = computed(() => {
  if (viewportWidth.value < 640) return 16;
  return 24;
});

const cardWidth = computed(() => {
  if (viewportWidth.value === 0) return 300;
  return (viewportWidth.value - gap.value * (cardsPerView.value - 1)) / cardsPerView.value;
});

const maxPage = computed(() => {
  return Math.max(0, props.lawyers.length - cardsPerView.value);
});

const offset = computed(() => {
  return -(currentPage.value * (cardWidth.value + gap.value));
});

function goToPage(page: number) {
  currentPage.value = Math.max(0, Math.min(page, maxPage.value));
  resetAutoPlay();
}

function updateViewport() {
  if (viewportRef.value) {
    viewportWidth.value = viewportRef.value.clientWidth;
  }
}

function resetAutoPlay() {
  stopAutoPlay();
  startAutoPlay();
}

function startAutoPlay() {
  if (props.lawyers.length <= cardsPerView.value) return;
  stopAutoPlay();
  autoTimer = setInterval(() => {
    if (currentPage.value >= maxPage.value) {
      currentPage.value = 0;
    } else {
      currentPage.value++;
    }
  }, 5000);
}

function stopAutoPlay() {
  if (autoTimer) {
    clearInterval(autoTimer);
    autoTimer = null;
  }
}

onMounted(() => {
  updateViewport();
  window.addEventListener('resize', updateViewport);
  nextTick(() => {
    updateViewport();
    startAutoPlay();
  });
});

onUnmounted(() => {
  stopAutoPlay();
  window.removeEventListener('resize', updateViewport);
});

defineExpose({ goToPage });
</script>

<style scoped>
.lawyer-carousel {
  position: relative;
}

.carousel-viewport {
  overflow: hidden;
  margin: 0 -16px;
  padding: 16px;
}

.carousel-track {
  display: flex;
  transition: transform 0.5s ease;
}

/* ── Card ── */

.lawyer-card {
  display: flex;
  align-items: flex-start;
  background: var(--bg-white);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-lg);
  overflow: hidden;
  box-shadow: 0 10px 28px rgba(15, 30, 61, 0.06);
  transition: box-shadow 0.4s var(--ease-out),
              transform 0.4s var(--ease-spring),
              border-color 0.3s var(--ease-out);
}

.lawyer-card:hover {
  box-shadow: 0 16px 36px rgba(15, 30, 61, 0.1);
  transform: translateY(-2px);
  border-color: rgba(200, 150, 62, 0.25);
}

/* ── Photo ── */

.lawyer-photo {
  flex: 0 0 178px;
  width: 178px;
  aspect-ratio: 3 / 4;
  margin: 22px 0 22px 22px;
  border-radius: 10px;
  overflow: hidden;
  position: relative;
  background: linear-gradient(135deg, #0F1E3D, #1A3A5C);
  box-shadow: 0 10px 24px rgba(15, 30, 61, 0.12);
}

.lawyer-photo-img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  object-position: center 16%;
  display: block;
  transition: transform 0.6s var(--ease-out);
}

.lawyer-card:hover .lawyer-photo-img {
  transform: scale(1.04);
}

.lawyer-photo-placeholder {
  width: 100%;
  height: 100%;
  background: linear-gradient(135deg, #15294D, #1E3A6E);
}

/* ── Body ── */

.lawyer-body {
  display: flex;
  flex-direction: column;
  gap: 8px;
  flex: 1;
  min-width: 0;
  padding: 24px 26px 26px 24px;
}

.lawyer-name {
  font-family: var(--font-serif);
  font-size: 20px;
  font-weight: 700;
  color: var(--color-text);
  margin: 0;
  line-height: 1.3;
  transition: color 0.3s var(--ease-out);
}

.lawyer-card:hover .lawyer-name {
  color: var(--color-primary-dark);
}

.lawyer-title {
  font-size: 14px;
  color: var(--color-accent);
  font-weight: 600;
  margin: 0;
  letter-spacing: 0.3px;
  line-height: 1.45;
  transition: color 0.3s var(--ease-out);
}

.lawyer-divider {
  width: 100%;
  height: 1px;
  background: linear-gradient(90deg, rgba(200, 150, 62, 0.45), rgba(200, 150, 62, 0.08), transparent);
  border-radius: 1px;
  margin: 6px 0 4px;
  transition: width 0.4s var(--ease-spring);
}

.lawyer-card:hover .lawyer-divider {
  width: 100%;
}

.lawyer-tags {
  list-style: none;
  margin: 4px 0 0;
  padding: 0;
  display: flex;
  flex-direction: column;
  gap: 0;
}

.lawyer-tags li {
  width: 100%;
  min-height: 34px;
  padding: 8px 0 8px 20px;
  border-top: 1px solid rgba(15, 30, 61, 0.08);
  font-size: 13px;
  color: var(--color-text-secondary);
  line-height: 1.45;
  overflow-wrap: anywhere;
  position: relative;
  transition: color 0.3s var(--ease-out),
              border-color 0.3s var(--ease-out);
}

.lawyer-tags li:first-child {
  border-top: none;
}

.lawyer-tags li::before {
  content: '';
  width: 5px;
  height: 5px;
  border-radius: 999px;
  background: var(--color-accent);
  position: absolute;
  left: 3px;
  top: 17px;
}

.lawyer-card:hover .lawyer-tags li {
  color: var(--color-text);
  border-color: rgba(200, 150, 62, 0.16);
}

.lawyer-card:hover .lawyer-tags li:nth-child(1) { transition-delay: 0s; }
.lawyer-card:hover .lawyer-tags li:nth-child(2) { transition-delay: 0.05s; }
.lawyer-card:hover .lawyer-tags li:nth-child(3) { transition-delay: 0.10s; }
.lawyer-card:hover .lawyer-tags li:nth-child(4) { transition-delay: 0.15s; }

/* ── Dots ── */

.carousel-dots {
  display: flex;
  justify-content: center;
  gap: 8px;
  margin-top: 24px;
}

.carousel-dot {
  width: 10px;
  height: 10px;
  border-radius: 5px;
  border: none;
  background: #d9d9d9;
  cursor: pointer;
  padding: 0;
  position: relative;
  transition: all var(--duration-normal) var(--ease-out);
}

.carousel-dot::after {
  content: '';
  position: absolute;
  top: -17px;
  left: -17px;
  right: -17px;
  bottom: -17px;
}

.carousel-dot:hover {
  background: var(--color-accent);
}

.carousel-dot.active {
  background: var(--color-accent);
  width: 24px;
}

.carousel-dot.active::after {
  left: -7px;
  right: -7px;
}

@media (max-width: 767px) {
  .carousel-viewport {
    margin: 0 -8px;
    padding: 8px;
  }

  .lawyer-card {
    display: grid;
    grid-template-columns: 132px minmax(0, 1fr);
    grid-template-areas:
      "photo name"
      "photo title"
      "divider divider"
      "tags tags";
    column-gap: 16px;
    align-items: center;
    padding: 18px;
  }

  .lawyer-photo {
    grid-area: photo;
    flex-basis: auto;
    width: 132px;
    aspect-ratio: 3 / 4;
    margin: 0;
    border-radius: 8px;
  }

  .lawyer-body {
    display: contents;
  }

  .lawyer-name {
    grid-area: name;
    align-self: end;
    font-size: 20px;
  }

  .lawyer-title {
    grid-area: title;
    align-self: start;
    font-size: 13.5px;
  }

  .lawyer-divider {
    grid-area: divider;
    margin-top: 16px;
  }

  .lawyer-tags {
    grid-area: tags;
    margin-top: 6px;
  }

  .lawyer-tags li {
    min-height: 34px;
    padding: 8px 0 8px 19px;
    font-size: 13px;
  }

  .lawyer-tags li::before {
    left: 2px;
    top: 16px;
  }

  .carousel-dot {
    width: 8px;
    height: 8px;
  }

  .carousel-dot.active {
    width: 20px;
  }
}
</style>
