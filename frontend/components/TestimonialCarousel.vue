<template>
  <div v-if="testimonials.length > 0" class="testimonial-carousel">
    <div class="carousel-viewport" ref="viewportRef" @touchstart="onTouchStart" @touchend="onTouchEnd">
      <div
        class="carousel-track"
        :style="{ transform: `translateX(${offset}px)`, gap: `${gap}px` }"
      >
        <div
          v-for="item in testimonials"
          :key="item.id"
          class="testimonial-card"
          :style="{ flex: `0 0 ${cardWidth}px`, width: `${cardWidth}px` }"
        >
          <!-- 顶部金色点缀 -->
          <div class="tm-accent-top"></div>

          <!-- 星级评分 -->
          <div class="tm-stars">
            <el-rate v-model="item.rating" disabled size="small" />
          </div>

          <!-- 评价正文 -->
          <blockquote class="tm-desc">{{ item.content }}</blockquote>

          <!-- 用户信息 -->
          <div class="tm-footer">
            <div class="tm-avatar">
              <ResponsiveImage v-if="item.avatar_url" :src="item.avatar_url" :alt="item.nickname" variant="thumb" :variants="item.avatar_variants" loading="lazy" />
              <svg v-else width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="8" r="4"/><path d="M4 22c0-4.4 3.6-8 8-8s8 3.6 8 8"/></svg>
            </div>
            <div class="tm-meta">
              <span class="tm-name">{{ item.nickname }}</span>
              <span class="tm-verified">
                <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><polyline points="20 6 9 17 4 12"/></svg>
                已签约客户
              </span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 小方块指示器 — 点击主动轮播 -->
    <div v-if="maxPage > 0" class="carousel-dots">
      <button
        v-for="p in maxPage + 1"
        :key="p"
        class="carousel-dot"
        :class="{ active: currentPage === p - 1 }"
        :aria-label="`第 ${p} 条评价`"
        @click="goToPage(p - 1)"
      ></button>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { ImageVariantInfo } from '~/utils/image'

export interface TestimonialItem {
  id: number;
  nickname: string;
  content: string;
  rating: number;
  avatar_url?: string;
  avatar_variants?: Record<string, ImageVariantInfo>;
}

const props = defineProps<{
  testimonials: TestimonialItem[];
}>();

const viewportRef = ref<HTMLElement | null>(null);
const currentPage = ref(0);
const viewportWidth = ref(0);
let autoTimer: ReturnType<typeof setInterval> | null = null;
let touchStartX = 0;

// ---- Touch swipe ----

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

// ---- 响应式计算 ----

const cardsPerView = computed(() => {
  if (viewportWidth.value < 640) return 1;
  if (viewportWidth.value < 1024) return 2;
  return 3;
});

const gap = computed(() => {
  if (viewportWidth.value < 640) return 16;
  if (viewportWidth.value < 1024) return 32;
  return 48;
});

const cardWidth = computed(() => {
  if (viewportWidth.value === 0) return 300;
  return (viewportWidth.value - gap.value * (cardsPerView.value - 1)) / cardsPerView.value;
});

const maxPage = computed(() => {
  return Math.max(0, props.testimonials.length - cardsPerView.value);
});

const offset = computed(() => {
  return -(currentPage.value * (cardWidth.value + gap.value));
});

// ---- 公开方法 ----

function goToPage(page: number) {
  currentPage.value = Math.max(0, Math.min(page, maxPage.value));
  resetAutoPlay();
}

function updateViewport() {
  if (viewportRef.value) {
    viewportWidth.value = viewportRef.value.clientWidth;
  }
}

// ---- 自动播放 ----

function resetAutoPlay() {
  stopAutoPlay();
  startAutoPlay();
}

function startAutoPlay() {
  if (props.testimonials.length <= cardsPerView.value) return;
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

// ---- 生命周期 ----

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

// 暴露方法供父组件手动控制
defineExpose({ goToPage });
</script>

<style scoped>
.testimonial-carousel {
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

.testimonial-card {
  position: relative;
  background: linear-gradient(135deg, #FEFEFE 0%, #FDF9F2 100%);
  border-radius: var(--radius-lg);
  padding: 32px 28px 24px;
  border-left: 3px solid var(--color-accent);
  box-shadow: var(--shadow-sm);
  transition: box-shadow var(--duration-normal) var(--ease-out),
              transform var(--duration-normal) var(--ease-out);
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.testimonial-card:hover {
  box-shadow: var(--shadow-md);
  transform: translateY(-2px);
}

/* 顶部金色点缀条 */
.tm-accent-top {
  position: absolute;
  top: 0;
  left: 20px;
  right: 20px;
  height: 2px;
  background: linear-gradient(90deg, transparent, var(--color-accent), transparent);
  border-radius: 0 0 1px 1px;
}

/* 星级 — 居中 */
.tm-stars {
  display: flex;
  justify-content: center;
  margin-bottom: 16px;
  position: relative;
  z-index: 1;
}

/* 评价正文 */
.tm-desc {
  font-size: 15px;
  line-height: 1.85;
  color: var(--color-text-secondary);
  margin: 0 0 20px;
  position: relative;
  z-index: 1;
  flex: 1;
}

/* 底部分隔 */
.tm-footer {
  display: flex;
  align-items: center;
  gap: 12px;
  padding-top: 16px;
  border-top: 1px solid rgba(200, 150, 62, 0.15);
  position: relative;
  z-index: 1;
}

.tm-avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  overflow: hidden;
  flex-shrink: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #e8ecf4, #dfe6f0);
  box-shadow: 0 0 0 2px var(--color-accent-soft);
}

.tm-avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.tm-avatar svg {
  color: #94a3b8;
}

.tm-meta {
  display: flex;
  flex-direction: column;
  gap: 2px;
  min-width: 0;
}

.tm-name {
  font-size: 14px;
  font-weight: 600;
  color: var(--color-text);
}

.tm-verified {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  font-size: 11px;
  color: var(--color-accent-dark);
  font-weight: 500;
}

.tm-verified svg {
  flex-shrink: 0;
  color: var(--color-accent);
}

/* 小方块指示器 */
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

  .carousel-dot {
    width: 8px;
    height: 8px;
  }

  .carousel-dot.active {
    width: 20px;
  }
}
</style>
