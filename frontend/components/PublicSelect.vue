<template>
  <select v-model="model" class="public-control public-select">
    <option v-for="option in options" :key="option.value" :value="option.value">
      {{ option.label }}
    </option>
  </select>
</template>

<script setup lang="ts" generic="T extends string | number">
defineProps<{
  options: readonly { value: T; label: string }[]
}>()

const model = defineModel<T>({ required: true })
</script>
