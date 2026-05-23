<template>
  <el-dialog
    :model-value="modelValue"
    @update:model-value="$emit('update:modelValue', $event)"
    title="选择图片"
    width="780px"
    destroy-on-close
  >
    <div class="picker-toolbar">
      <el-input
        v-model="searchText"
        placeholder="搜索文件名..."
        clearable
        style="width: 260px"
        @input="onSearch"
      />
      <el-upload
        :action="uploadUrl"
        :headers="uploadHeaders"
        accept=".jpg,.jpeg,.png,.webp"
        :show-file-list="false"
        :on-success="onUploadSuccess"
      >
        <el-button type="primary">上传</el-button>
      </el-upload>
    </div>

    <div class="picker-body" v-loading="loading">
      <div class="picker-grid">
        <div
          v-for="item in list"
          :key="item.id"
          class="picker-item"
          :class="{ selected: selectedId === item.id }"
          @click="selectedId = item.id"
        >
          <ResponsiveImage :src="item.url" :alt="item.original_name" variant="sm" />
        </div>
        <div v-if="!loading && list.length === 0" class="empty-hint">
          暂无图片
        </div>
      </div>

      <div class="picker-detail" v-if="selectedItem">
        <div class="detail-preview">
          <ResponsiveImage :src="selectedItem.url" :alt="selectedItem.original_name" variant="md" />
        </div>
        <div class="detail-info">
          <p class="detail-name">{{ selectedItem.original_name || selectedItem.filename }}</p>
          <p class="detail-meta">{{ selectedItem.url }}</p>
          <p class="detail-meta">{{ formatSize(selectedItem.size_bytes) }}</p>
        </div>
      </div>
    </div>

    <div class="picker-footer">
      <el-pagination
        v-if="total > pageSize"
        v-model:current-page="page"
        :page-size="pageSize"
        :total="total"
        layout="total, prev, pager, next"
        small
        @current-change="loadList"
      />
      <div class="footer-actions">
        <el-button @click="$emit('update:modelValue', false)">取消</el-button>
        <el-button type="primary" :disabled="!selectedId" @click="onConfirm">
          确认选择
        </el-button>
      </div>
    </div>
  </el-dialog>
</template>

<script setup lang="ts">
import { ElMessage } from 'element-plus';

interface MediaItem {
  id: number;
  filename: string;
  original_name: string;
  url: string;
  mime_type: string;
  size_bytes: number;
}

const props = defineProps<{
  modelValue: boolean;
}>();

const emit = defineEmits<{
  'update:modelValue': [value: boolean];
  select: [url: string];
}>();

const list = ref<MediaItem[]>([]);
const loading = ref(false);
const page = ref(1);
const pageSize = ref(12);
const total = ref(0);
const searchText = ref('');
const selectedId = ref<number | null>(null);

const uploadUrl = '/api/v1/admin/media/upload';

const uploadHeaders = computed(() => {
  const token = import.meta.client ? localStorage.getItem('token') : null;
  return token ? { Authorization: `Bearer ${token}` } : {};
});

const selectedItem = computed(() =>
  list.value.find((item) => item.id === selectedId.value) || null
);

let searchTimer: ReturnType<typeof setTimeout> | null = null;

const formatSize = (bytes: number) => {
  if (!bytes) return '未知大小';
  if (bytes < 1024) return `${bytes} B`;
  if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`;
  return `${(bytes / (1024 * 1024)).toFixed(1)} MB`;
};

const loadList = async () => {
  loading.value = true;
  try {
    const api = useApi();
    const data = await api<{ items: MediaItem[]; total: number }>(
      `/admin/media?page=${page.value}&per_page=${pageSize.value}&search=${encodeURIComponent(searchText.value)}`
    );
    list.value = data.items ?? [];
    total.value = data.total ?? 0;
  } catch {
    ElMessage.error('加载媒体列表失败');
  } finally {
    loading.value = false;
  }
};

const onSearch = () => {
  if (searchTimer) clearTimeout(searchTimer);
  searchTimer = setTimeout(() => {
    page.value = 1;
    loadList();
  }, 300);
};

const onUploadSuccess = () => {
  page.value = 1;
  loadList();
};

const onConfirm = () => {
  if (selectedItem.value) {
    emit('select', selectedItem.value.url);
    emit('update:modelValue', false);
  }
};

watch(
  () => props.modelValue,
  (val) => {
    if (val) {
      selectedId.value = null;
      searchText.value = '';
      page.value = 1;
      loadList();
    }
  }
);
</script>

<style scoped>
.picker-toolbar {
  display: flex;
  gap: 12px;
  align-items: center;
  margin-bottom: 16px;
}

.picker-body {
  display: flex;
  gap: 16px;
  min-height: 320px;
}

.picker-grid {
  flex: 1;
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 8px;
  align-content: start;
}

.picker-item {
  aspect-ratio: 1;
  border: 2px solid transparent;
  border-radius: 6px;
  overflow: hidden;
  cursor: pointer;
  background: #f5f7fa;
  display: flex;
  align-items: center;
  justify-content: center;
}

.picker-item img {
  max-width: 100%;
  max-height: 100%;
  object-fit: contain;
}

.picker-item.selected {
  border-color: #409eff;
  box-shadow: 0 0 0 2px rgba(64, 158, 255, 0.2);
}

.picker-detail {
  width: 200px;
  flex-shrink: 0;
  border-left: 1px solid #ebeef5;
  padding-left: 16px;
}

.detail-preview {
  aspect-ratio: 1;
  background: #f5f7fa;
  border-radius: 6px;
  display: flex;
  align-items: center;
  justify-content: center;
  overflow: hidden;
  margin-bottom: 12px;
}

.detail-preview img {
  max-width: 100%;
  max-height: 100%;
  object-fit: contain;
}

.detail-name {
  font-size: 13px;
  color: #303133;
  word-break: break-all;
  margin: 0 0 4px;
}

.detail-meta {
  font-size: 12px;
  color: #909399;
  word-break: break-all;
  margin: 0 0 2px;
}

.empty-hint {
  grid-column: 1 / -1;
  text-align: center;
  padding: 60px 0;
  color: #909399;
  font-size: 14px;
}

.picker-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-top: 16px;
}

.footer-actions {
  display: flex;
  gap: 8px;
  margin-left: auto;
}
</style>
