<template>
  <div class="admin-table-wrap">
    <AdminLoadingOverlay :show="loading" />
    <slot />
  </div>
</template>

<script setup lang="ts">
withDefaults(defineProps<{
  loading?: boolean;
}>(), {
  loading: false,
});
</script>
