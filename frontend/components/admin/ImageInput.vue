<template>
  <div class="image-input">
    <div class="input-row">
      <el-input v-model="urlValue" :placeholder="placeholder" @input="onInput" />
      <el-upload
        :action="uploadUrl"
        :headers="uploadHeaders"
        accept=".jpg,.jpeg,.png,.webp"
        :show-file-list="false"
        :on-success="onUploadSuccess"
      >
        <el-button>上传</el-button>
      </el-upload>
      <el-button @click="pickerVisible = true">浏览</el-button>
    </div>
    <div v-if="urlValue" class="preview" :style="{ aspectRatio: previewRatio }">
      <img :src="previewSrc" alt="预览" @error="previewError = true" v-show="!previewError" />
    </div>
    <p v-if="sizeHint" class="size-hint">{{ sizeHint }}</p>

    <MediaPicker
      v-model="pickerVisible"
      @select="onPickerSelect"
    />
  </div>
</template>

<script setup lang="ts">
import { getVariantUrl } from '~/utils/image'
import MediaPicker from './MediaPicker.vue';

const props = defineProps<{
  modelValue: string;
  placeholder?: string;
  sizeHint?: string;
  previewRatio?: string;
  context?: string;
}>();

const emit = defineEmits<{
  'update:modelValue': [value: string];
}>();

const urlValue = ref(props.modelValue);
const pickerVisible = ref(false);
const previewError = ref(false);

const previewSrc = computed(() => getVariantUrl(urlValue.value, 'sm'))

const uploadUrl = computed(() => {
  const ctx = props.context || 'general';
  return `/api/v1/admin/media/upload?context=${ctx}`;
});

const uploadHeaders = computed(() => {
  const token = import.meta.client ? localStorage.getItem('token') : null;
  return token ? { Authorization: `Bearer ${token}` } : {};
});

const onInput = (val: string) => {
  previewError.value = false;
  emit('update:modelValue', val);
};

const onUploadSuccess = (res: any) => {
  const url = res?.data?.url || res?.url || '';
  if (url) {
    urlValue.value = url;
    previewError.value = false;
    emit('update:modelValue', url);
  }
};

const onPickerSelect = (url: string) => {
  urlValue.value = url;
  previewError.value = false;
  emit('update:modelValue', url);
};

watch(
  () => props.modelValue,
  (val) => {
    urlValue.value = val;
    previewError.value = false;
  }
);
</script>

<style scoped>
.image-input {
  flex: 1;
  min-width: 0;
}

.input-row {
  display: flex;
  gap: 8px;
}

.input-row > .el-input {
  flex: 1;
}

.input-row > .el-button {
  flex-shrink: 0;
}

.preview {
  margin-top: 8px;
  width: 120px;
  aspect-ratio: 16 / 9;
  border-radius: 4px;
  overflow: hidden;
  background: #f5f7fa;
  border: 1px solid #ebeef5;
}

.preview img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.size-hint {
  margin-top: 6px;
  font-size: 12px;
  color: var(--color-text-muted);
  line-height: 1.4;
}
</style>
