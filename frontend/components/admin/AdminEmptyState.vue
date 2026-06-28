<template>
  <div class="admin-empty-state">
    <div v-if="icon" class="empty-icon" v-html="getIconSvg(icon, 48)"></div>
    <div class="empty-title">{{ title }}</div>
    <div v-if="description" class="empty-desc">{{ description }}</div>
    <el-button v-if="actionLabel" type="primary" @click="$emit('action')">
      {{ actionLabel }}
    </el-button>
  </div>
</template>

<script setup lang="ts">
import { getIconSvg } from '~/composables/lucideIcons';

defineProps<{
  icon?: string;
  title: string;
  description?: string;
  actionLabel?: string;
}>();

defineEmits<{
  action: [];
}>();
</script>
