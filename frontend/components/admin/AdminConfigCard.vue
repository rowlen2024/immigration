<template>
  <el-card class="config-card">
    <template #header>
      <h3 class="admin-card-title">{{ title }}</h3>
    </template>
    <slot />
    <div v-if="$slots.actions" class="config-list-actions">
      <slot name="actions" />
    </div>
    <div v-if="$slots.footer" class="card-footer">
      <slot name="footer" />
    </div>
  </el-card>
</template>

<script setup lang="ts">
defineProps<{
  title: string;
}>();
</script>
