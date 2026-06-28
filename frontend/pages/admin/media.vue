<template>
  <div>
    <AdminPageHeader title="媒体库">
      <template #actions>
        <el-upload
          :action="uploadUrl"
          :headers="uploadHeaders"
          accept=".jpg,.jpeg,.png,.webp"
          :show-file-list="false"
          :disabled="uploading"
          :on-progress="handleUploadProgress"
          :on-success="handleUploadSuccess"
          :on-error="handleUploadError"
        >
          <el-button type="primary" :loading="uploading">
            <span class="admin-button-icon" v-html="getIconSvg('upload', 16)"></span>
            {{ uploading ? '上传中' : '上传图片' }}
          </el-button>
        </el-upload>
      </template>
    </AdminPageHeader>

    <AdminToolbar>
      <el-input
        v-model="searchQuery"
        placeholder="搜索文件名..."
        :prefix-icon="Search"
        clearable
        class="admin-search-input"
        @input="onSearch"
      />
      <el-button :icon="Refresh" circle @click="searchQuery='';loadList()" :loading="loading" />
      <el-button type="warning" @click="handleCleanupUnused" :loading="cleanupLoading">
        清理未使用
      </el-button>
    </AdminToolbar>

    <div class="admin-panel-shell">
      <AdminLoadingOverlay :show="loading" />
      <div v-if="list.length > 0" class="admin-media-grid media-grid">
        <article v-for="item in list" :key="item.id" class="media-card">
          <button class="media-preview-btn" type="button" @click="openPreview(item.url)">
            <ResponsiveImage
              :src="item.url"
              :alt="item.filename"
              variant="sm"
              class="media-card-thumb"
            />
          </button>
          <div class="media-card-body">
            <div class="media-card-title" :title="item.original_name || item.filename">
              {{ item.original_name || item.filename }}
            </div>
            <div class="media-card-meta">
              <span>{{ formatSize(item.size_bytes) }}</span>
              <span>{{ formatDateTime(item.created_at) }}</span>
            </div>
            <div class="media-card-url" :title="item.url">{{ item.url || '—' }}</div>
          </div>
          <div class="media-card-actions">
            <el-tooltip content="复制链接" placement="top">
              <button class="action-btn" type="button" title="复制链接" aria-label="复制链接" @click="copyLink(item.url)" v-html="getIconSvg('copy', 16)"></button>
            </el-tooltip>
            <el-tooltip content="预览" placement="top">
              <button class="action-btn" type="button" title="预览" aria-label="预览" @click="openPreview(item.url)" v-html="getIconSvg('eye', 16)"></button>
            </el-tooltip>
            <el-popconfirm
              title="确定删除该文件？"
              confirm-button-text="删除"
              cancel-button-text="取消"
              @confirm="handleDelete(item.id)"
            >
              <template #reference>
                <button class="action-btn danger" type="button" title="删除" aria-label="删除" v-html="getIconSvg('trash-2', 16)"></button>
              </template>
            </el-popconfirm>
          </div>
        </article>
      </div>
      <AdminEmptyState
        v-else-if="!loading"
        icon="file-text"
        title="暂无媒体文件"
        description="上传图片后可在这里预览、复制链接和清理未使用文件"
      />
    </div>

    <div class="admin-pagination-wrap" v-if="total > pageSize">
      <el-pagination
        v-model:current-page="page"
        :page-size="pageSize"
        :total="total"
        layout="total, prev, pager, next"
        @current-change="loadList"
      />
    </div>

    <div v-if="previewImage" class="preview-overlay" @click.self="closePreview">
      <div class="preview-container">
        <ResponsiveImage :src="previewImage" alt="预览" variant="md" />
        <button class="preview-close" @click="closePreview" v-html="getIconSvg('x', 24)"></button>
      </div>
    </div>

    <el-dialog
      v-model="showCleanupDialog"
      title="确认清理未使用媒体"
      width="600px"
      :close-on-click-modal="false"
    >
      <p>发现 {{ unusedMediaList.length }} 个未被引用的媒体文件，以下文件将被永久删除：</p>
      <div class="cleanup-file-list">
        <div
          v-for="item in unusedMediaList"
          :key="item.id"
          class="cleanup-file-item"
        >
          <ResponsiveImage :src="item.url" :alt="item.filename" variant="thumb" class="cleanup-file-thumb" />
          <div class="cleanup-file-info">
            <div class="cleanup-file-name">{{ item.original_name || item.filename }}</div>
            <div class="cleanup-file-url">{{ item.url }}</div>
          </div>
          <span class="cleanup-file-size">{{ formatSize(item.size_bytes) }}</span>
        </div>
      </div>
      <template #footer>
        <el-button @click="showCleanupDialog = false">取消</el-button>
        <el-button type="danger" @click="confirmCleanup" :loading="cleanupDeleting">
          确认删除 {{ unusedMediaList.length }} 个文件
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ElMessage } from 'element-plus';
import { Refresh, Search } from '@element-plus/icons-vue';
import { useNotify } from '~/composables/useNotify';
import { getIconSvg } from '~/composables/lucideIcons';
import { formatDateTime } from '~/utils/date';

definePageMeta({ layout: 'admin', middleware: 'auth' });

const notify = useNotify();

interface MediaItem {
  id: string;
  filename: string;
  original_name: string;
  url: string;
  size_bytes: number;
  created_at: string;
}

const list = ref<MediaItem[]>([]);
const loading = ref(false);
const searchQuery = ref('');
const page = ref(1);
const pageSize = ref(10);
const total = ref(0);
const uploading = ref(false);
let searchTimer: ReturnType<typeof setTimeout> | null = null;

const previewImage = ref<string | null>(null);

const uploadUrl = '/api/v1/admin/media/upload';

const uploadHeaders = computed(() => {
  const token = import.meta.client ? localStorage.getItem('token') : null;
  return token ? { Authorization: `Bearer ${token}` } : {};
});

const formatSize = (bytes: number) => {
  if (!bytes) return '未知大小';
  if (bytes < 1024) return `${bytes} B`;
  if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`;
  return `${(bytes / (1024 * 1024)).toFixed(1)} MB`;
};

const openPreview = (url: string) => {
  previewImage.value = url;
};

const closePreview = () => {
  previewImage.value = null;
};

const copyLink = async (url: string) => {
  if (!url) return;
  try {
    if (navigator.clipboard?.writeText) {
      await navigator.clipboard.writeText(url);
    } else {
      const textarea = document.createElement('textarea');
      textarea.value = url;
      textarea.style.position = 'fixed';
      textarea.style.opacity = '0';
      document.body.appendChild(textarea);
      textarea.select();
      document.execCommand('copy');
      document.body.removeChild(textarea);
    }
    notify.success('链接已复制');
  } catch (e) {
    notify.error(e, '复制失败');
  }
};

const onSearch = () => {
  if (searchTimer) clearTimeout(searchTimer);
  searchTimer = setTimeout(() => {
    page.value = 1;
    loadList();
  }, 300);
};

const loadList = async () => {
  loading.value = true;
  try {
    const api = useApi();
    const params = new URLSearchParams();
    params.set('page', String(page.value));
    params.set('per_page', String(pageSize.value));
    if (searchQuery.value) params.set('search', searchQuery.value);
    const data = await api<{ items: MediaItem[]; total: number }>(
      `/admin/media?${params.toString()}`
    );
    list.value = data.items ?? [];
    total.value = data.total ?? 0;
  } catch {
    list.value = [];
    ElMessage.error('加载媒体列表失败');
  } finally {
    loading.value = false;
  }
};

const handleUploadProgress = () => {
  uploading.value = true;
};

const handleUploadSuccess = () => {
  uploading.value = false;
  notify.success('上传成功');
  loadList();
};

const handleUploadError = () => {
  uploading.value = false;
  ElMessage.error('上传失败');
};

const handleDelete = async (id: string) => {
  try {
    const api = useApi();
    await api(`/admin/media/${id}`, { method: 'DELETE' });
    notify.success('已删除');
    loadList();
  } catch (e) {
    notify.error(e, '操作失败');
  }
};

const cleanupLoading = ref(false);
const unusedMediaList = ref<MediaItem[]>([]);
const showCleanupDialog = ref(false);
const cleanupDeleting = ref(false);

const handleCleanupUnused = async () => {
  cleanupLoading.value = true;
  try {
    const api = useApi();
    unusedMediaList.value = await api<MediaItem[]>('/admin/media/unused');
    if (unusedMediaList.value.length === 0) {
      ElMessage.success('没有未使用的媒体文件');
      return;
    }
    showCleanupDialog.value = true;
  } catch (e) {
    notify.error(e, '获取未使用媒体列表失败');
  } finally {
    cleanupLoading.value = false;
  }
};

const confirmCleanup = async () => {
  cleanupDeleting.value = true;
  try {
    const api = useApi();
    const ids = unusedMediaList.value.map((m) => Number(m.id));
    const result = await api<{ deleted: number; failed: string[] }>('/admin/media/cleanup', {
      method: 'POST',
      body: { ids },
    });
    notify.success(`已清理 ${result.deleted} 个文件`);
    if (result.failed && result.failed.length > 0) {
      ElMessage.warning(`部分文件清理失败: ${result.failed.join(', ')}`);
    }
    showCleanupDialog.value = false;
    loadList();
  } catch (e) {
    notify.error(e, '清理失败');
  } finally {
    cleanupDeleting.value = false;
  }
};

onMounted(() => {
  loadList();
});
</script>

<style scoped>
.media-grid {
  grid-template-columns: repeat(5, minmax(0, 1fr));
  min-height: 220px;
}

.media-card {
  display: flex;
  flex-direction: column;
  min-width: 0;
  background: var(--color-bg-surface);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  overflow: hidden;
  box-shadow: var(--shadow-xs);
  transition: border-color 0.15s ease, box-shadow 0.15s ease;
}

.media-card:hover {
  border-color: #bfdbfe;
  box-shadow: var(--shadow-md);
}

.media-preview-btn {
  display: block;
  width: 100%;
  aspect-ratio: 4 / 3;
  padding: 0;
  border: 0;
  background: var(--color-bg-app);
  cursor: pointer;
  overflow: hidden;
}

.media-card-thumb {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.media-card-body {
  min-width: 0;
  padding: 12px;
  border-top: 1px solid var(--color-border-light);
}

.media-card-title {
  min-height: 20px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  font-size: 14px;
  font-weight: 600;
  color: var(--color-text);
}

.media-card-meta {
  display: flex;
  justify-content: space-between;
  gap: 8px;
  margin-top: 6px;
  font-size: 12px;
  color: var(--color-text-muted);
}

.media-card-url {
  margin-top: 8px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  font-size: 12px;
  color: var(--color-text-muted);
}

.media-card-actions {
  display: flex;
  justify-content: flex-end;
  gap: 4px;
  padding: 8px 10px 10px;
}

.media-card-actions .action-btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-width: 32px;
  min-height: 32px;
  padding: 5px 8px;
  color: var(--color-text-secondary);
  background: transparent;
  border: 0;
  border-radius: var(--radius-sm);
  cursor: pointer;
  transition: color 0.15s ease, background-color 0.15s ease;
}

.media-card-actions .action-btn:hover {
  color: var(--color-primary);
  background: var(--color-info-soft);
}

.media-card-actions .action-btn.danger:hover {
  color: var(--color-danger);
  background: var(--color-danger-soft);
}

@media (max-width: 1280px) {
  .media-grid {
    grid-template-columns: repeat(4, minmax(0, 1fr));
  }
}

@media (max-width: 1024px) {
  .media-grid {
    grid-template-columns: repeat(3, minmax(0, 1fr));
  }
}

@media (max-width: 767px) {
  .media-grid {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
}

.preview-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.8);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 2000;
}

.preview-container {
  position: relative;
  max-width: 90vw;
  max-height: 90vh;
}

.preview-container img {
  max-width: 90vw;
  max-height: 90vh;
  object-fit: contain;
  border-radius: 4px;
}

.preview-close {
  position: absolute;
  top: -40px;
  right: 0;
  background: none;
  border: none;
  color: #fff;
  font-size: 28px;
  cursor: pointer;
  width: 36px;
  height: 36px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
}

.preview-close:hover {
  background: rgba(255, 255, 255, 0.15);
}

.cleanup-file-list {
  max-height: 300px;
  overflow-y: auto;
  margin: 16px 0;
}

.cleanup-file-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 8px 0;
  border-bottom: 1px solid var(--el-border-color-light);
}

.cleanup-file-thumb {
  width: 40px;
  height: 40px;
  object-fit: cover;
  border-radius: 4px;
  flex-shrink: 0;
}

.cleanup-file-info {
  flex: 1;
  min-width: 0;
}

.cleanup-file-name {
  font-size: 13px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.cleanup-file-url {
  font-size: 11px;
  color: var(--el-text-color-secondary);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.cleanup-file-size {
  font-size: 11px;
  color: var(--el-text-color-secondary);
  white-space: nowrap;
  flex-shrink: 0;
}
</style>
