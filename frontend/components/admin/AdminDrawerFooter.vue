<template>
  <div class="admin-drawer-footer">
    <el-button @click="$emit('cancel')">{{ cancelText }}</el-button>
    <el-button type="primary" :loading="loading" @click="$emit('confirm')">
      {{ confirmText }}
    </el-button>
  </div>
</template>

<script setup lang="ts">
withDefaults(defineProps<{
  loading?: boolean;
  cancelText?: string;
  confirmText?: string;
}>(), {
  loading: false,
  cancelText: '取消',
  confirmText: '保存',
});

defineEmits<{
  cancel: [];
  confirm: [];
}>();
</script>
