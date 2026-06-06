<template>
  <img
    :src="resolvedSrc"
    :srcset="srcsetAttr"
    :sizes="sizesAttr"
    :loading="loading"
    :fetchpriority="fetchpriority"
    :decoding="loading === 'lazy' ? 'async' : 'sync'"
    v-bind="$attrs"
    @error="onError"
  />
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import {
  getVariantUrl,
  getSrcset,
  getVariantFromData,
  type ImageVariant,
  type ImageVariantInfo,
} from '~/utils/image'

const props = withDefaults(defineProps<{
  src: string | null | undefined
  variant: ImageVariant
  variants?: Record<string, ImageVariantInfo> | null
  sizes?: string
  loading?: 'lazy' | 'eager'
  fetchpriority?: 'high' | 'low' | 'auto'
}>(), {
  loading: 'lazy',
})

const failed = ref(false)

const resolvedSrc = computed(() => {
  if (failed.value || !props.src) return props.src ?? ''
  // 优先从 variants 数据精确取值，失败回退到盲拼
  if (props.variants) {
    return getVariantFromData(props.variants, props.variant, props.src) ?? getVariantUrl(props.src, props.variant)
  }
  return getVariantUrl(props.src, props.variant)
})

const srcsetAttr = computed(() => {
  if (failed.value) return undefined
  return getSrcset(props.variants)
})

const sizesAttr = computed(() => {
  return props.sizes ?? '(max-width: 640px) 100vw, (max-width: 1024px) 50vw, 33vw'
})

function onError() {
  if (!failed.value) {
    failed.value = true
  }
}
</script>
