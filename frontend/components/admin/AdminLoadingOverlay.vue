<template>
  <div v-if="show" class="admin-loading-overlay">
    <span class="admin-loading-spinner"></span>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  show: boolean;
}>();
</script>

<style scoped>
.admin-loading-overlay {
  position: absolute;
  inset: 0;
  z-index: 10;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: inherit;
  background: rgba(255, 255, 255, 0.85);
}

.admin-loading-spinner {
  width: 32px;
  height: 32px;
  border: 3px solid var(--el-border-color-light);
  border-top-color: var(--el-color-primary);
  border-radius: 50%;
  animation: admin-loading-spin 0.8s linear infinite;
}

@keyframes admin-loading-spin {
  to {
    transform: rotate(360deg);
  }
}
</style>
