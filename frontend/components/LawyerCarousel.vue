<template>
  <div v-if="lawyers.length > 0" class="lawyer-carousel">
    <div class="carousel-viewport" ref="viewportRef">
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
            <img v-if="lawyer.photo_url" :src="lawyer.photo_url" :alt="lawyer.name" loading="lazy" />
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
export interface LawyerItem {
  id: number;
  name: string;
  title: string;
  tags: string[];
  photo_url?: string;
}

const props = defineProps<{
  lawyers: LawyerItem[];
}>();

const viewportRef = ref<HTMLElement | null>(null);
const currentPage = ref(0);
const viewportWidth = ref(0);
let autoTimer: ReturnType<typeof setInterval> | null = null;

const cardsPerView = computed(() => {
  if (viewportWidth.value < 640) return 1;
  if (viewportWidth.value < 1024) return 2;
  return 3;
});

const gap = computed(() => {
  if (viewportWidth.value < 640) return 16;
  if (viewportWidth.value < 1024) return 24;
  return 28;
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
  margin: 0 -12px;
  padding: 12px;
}

.carousel-track {
  display: flex;
  transition: transform 0.5s ease;
}

/* ── Card ── */

.lawyer-card {
  display: flex;
  gap: 0;
  background: var(--bg-white);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-lg);
  overflow: hidden;
  box-shadow: var(--shadow-sm);
  transition: box-shadow 0.4s var(--ease-out),
              transform 0.4s var(--ease-spring),
              border-color 0.3s var(--ease-out);
}

.lawyer-card:hover {
  box-shadow: var(--shadow-lg);
  transform: translateY(-2px);
  border-color: rgba(200, 150, 62, 0.25);
}

/* ── Photo ── */

.lawyer-photo {
  flex: 0 0 40%;
  overflow: hidden;
  position: relative;
  background: linear-gradient(135deg, #0F1E3D, #1A3A5C);
}

.lawyer-photo img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  display: block;
  transition: transform 0.6s var(--ease-out);
}

.lawyer-card:hover .lawyer-photo img {
  transform: scale(1.04);
}

.lawyer-photo-placeholder {
  width: 100%;
  height: 100%;
  min-height: 150px;
  background: linear-gradient(135deg, #15294D, #1E3A6E);
}

/* ── Body ── */

.lawyer-body {
  display: flex;
  flex-direction: column;
  gap: 6px;
  flex: 1;
  padding: 20px 20px 20px 18px;
}

.lawyer-name {
  font-family: var(--font-serif);
  font-size: 19px;
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
  font-size: 13px;
  color: var(--color-accent);
  font-weight: 500;
  margin: 0;
  letter-spacing: 0.5px;
  transition: color 0.3s var(--ease-out);
}

.lawyer-divider {
  width: 40px;
  height: 2px;
  background: var(--gradient-gold);
  border-radius: 1px;
  margin: 6px 0 2px;
  transition: width 0.4s var(--ease-spring);
}

.lawyer-card:hover .lawyer-divider {
  width: 60px;
}

.lawyer-tags {
  list-style: none;
  margin: 0;
  padding: 0;
  display: flex;
  flex-direction: column;
  gap: 5px;
}

.lawyer-tags li {
  font-size: 12.5px;
  color: var(--color-text-secondary);
  display: flex;
  align-items: center;
  gap: 8px;
  line-height: 1.5;
  transition: color 0.3s var(--ease-out);
}

.lawyer-card:hover .lawyer-tags li {
  color: var(--color-text);
}

.lawyer-card:hover .lawyer-tags li:nth-child(1) { transition-delay: 0s; }
.lawyer-card:hover .lawyer-tags li:nth-child(2) { transition-delay: 0.05s; }
.lawyer-card:hover .lawyer-tags li:nth-child(3) { transition-delay: 0.10s; }
.lawyer-card:hover .lawyer-tags li:nth-child(4) { transition-delay: 0.15s; }

.lawyer-tags li::before {
  content: '';
  display: inline-block;
  width: 4px;
  height: 4px;
  border-radius: 50%;
  background: var(--color-accent);
  flex-shrink: 0;
  transition: transform 0.3s var(--ease-spring),
              box-shadow 0.3s var(--ease-out);
}

.lawyer-card:hover .lawyer-tags li::before {
  transform: scale(1.8);
  box-shadow: 0 0 5px rgba(200, 150, 62, 0.45);
}

.lawyer-card:hover .lawyer-tags li:nth-child(1)::before { transition-delay: 0s; }
.lawyer-card:hover .lawyer-tags li:nth-child(2)::before { transition-delay: 0.05s; }
.lawyer-card:hover .lawyer-tags li:nth-child(3)::before { transition-delay: 0.10s; }
.lawyer-card:hover .lawyer-tags li:nth-child(4)::before { transition-delay: 0.15s; }

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
  transition: all var(--duration-normal) var(--ease-out);
}

.carousel-dot:hover {
  background: var(--color-accent);
}

.carousel-dot.active {
  background: var(--color-accent);
  width: 24px;
}

@media (max-width: 767px) {
  .carousel-viewport {
    margin: 0 -8px;
    padding: 8px;
  }

  .lawyer-photo {
    flex: 0 0 36%;
    min-height: 140px;
  }

  .lawyer-body {
    padding: 16px 16px 16px 14px;
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
