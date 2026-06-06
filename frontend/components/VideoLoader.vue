<template>
  <div ref="containerRef" class="lazy-video" :style="containerStyle">
    <!-- Poster / placeholder -->
    <div v-if="!activated" class="lazy-video-poster" @click="activate">
      <img
        v-if="poster"
        :src="poster"
        :alt="alt"
        class="lazy-video-poster-img"
        loading="lazy"
      />
      <div class="lazy-video-play-btn" aria-label="播放视频">
        <svg width="64" height="64" viewBox="0 0 24 24" fill="none">
          <circle cx="12" cy="12" r="11" stroke="white" stroke-width="1.5" fill="rgba(0,0,0,0.4)"/>
          <polygon points="9,7 9,17 18,12" fill="white"/>
        </svg>
      </div>
    </div>

    <!-- Native video (loaded on demand) -->
    <video
      v-if="activated && !youtubeId"
      ref="videoRef"
      :src="src"
      :poster="poster"
      :controls="controls"
      :autoplay="autoplay"
      :loop="loop"
      :muted="muted"
      :playsinline="playsinline"
      preload="none"
      class="lazy-video-el"
      v-bind="$attrs"
      @canplaythrough="$emit('ready')"
      @error="$emit('error', $event)"
    />

    <!-- YouTube iframe (loaded on demand) -->
    <iframe
      v-if="activated && youtubeId"
      :src="youtubeSrc"
      :title="alt"
      class="lazy-video-iframe"
      frameborder="0"
      allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture"
      allowfullscreen
      loading="lazy"
    />
  </div>
</template>

<script setup lang="ts">
const props = withDefaults(defineProps<{
  src?: string
  poster?: string
  youtubeId?: string
  alt?: string
  aspectRatio?: string
  controls?: boolean
  autoplay?: boolean
  loop?: boolean
  muted?: boolean
  playsinline?: boolean
}>(), {
  aspectRatio: '16 / 9',
  controls: true,
  playsinline: true,
  muted: false,
})

defineEmits<{
  ready: []
  error: [event: Event]
}>()

const containerRef = ref<HTMLElement | null>(null)
const videoRef = ref<HTMLVideoElement | null>(null)
const activated = ref(false)

const containerStyle = computed(() => ({
  aspectRatio: props.aspectRatio,
}))

const youtubeSrc = computed(() =>
  props.youtubeId
    ? `https://www.youtube-nocookie.com/embed/${props.youtubeId}?rel=0&modestbranding=1`
    : ''
)

function activate() {
  activated.value = true
  nextTick(() => {
    if (videoRef.value && props.autoplay) {
      videoRef.value.play().catch(() => {
        // 浏览器可能阻止自动播放，用户需要手动点击
      })
    }
  })
}

// YouTube: 滚动到附近时预激活 iframe（不自动播放，用户仍需点击）
let observer: IntersectionObserver | null = null
onMounted(() => {
  if (props.youtubeId && containerRef.value) {
    observer = new IntersectionObserver(
      (entries) => {
        if (entries[0]?.isIntersecting) {
          activated.value = true
          observer?.disconnect()
        }
      },
      { rootMargin: '300px 0px' }
    )
    observer.observe(containerRef.value)
  }
})

onUnmounted(() => {
  observer?.disconnect()
})
</script>

<style scoped>
.lazy-video {
  position: relative;
  width: 100%;
  background: #000;
  border-radius: var(--radius-md, 8px);
  overflow: hidden;
}

.lazy-video-poster {
  position: absolute;
  inset: 0;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
}

.lazy-video-poster-img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.lazy-video-play-btn {
  position: absolute;
  z-index: 1;
  transition: transform 0.2s ease, opacity 0.2s ease;
}

.lazy-video-poster:hover .lazy-video-play-btn {
  transform: scale(1.1);
  opacity: 0.9;
}

.lazy-video-el,
.lazy-video-iframe {
  width: 100%;
  height: 100%;
  display: block;
}
</style>
