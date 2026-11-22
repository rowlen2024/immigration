<template>
  <div>
    <NuxtPage v-if="route.params.a && route.params.b" />
    <ProjectComparison v-else />
  </div>
</template>

<script setup lang="ts">
const route = useRoute()
</script>
