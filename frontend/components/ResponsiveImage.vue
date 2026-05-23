<template>
  <img :src="resolvedSrc" v-bind="$attrs" @error="onError" />
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import { getVariantUrl, type ImageVariant } from '~/utils/image'

const props = defineProps<{
  src: string | null | undefined
  variant: ImageVariant
}>()

const failed = ref(false)

const resolvedSrc = computed(() => {
  if (failed.value || !props.src) return props.src ?? ''
  return getVariantUrl(props.src, props.variant)
})

function onError() {
  if (!failed.value) {
    failed.value = true
  }
}
</script>
