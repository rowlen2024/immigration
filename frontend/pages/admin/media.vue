<template>
  <div>
    <div class="admin-page-header">
      <h2 class="admin-page-title">媒体库</h2>
      <div style="display:flex;align-items:center;gap:8px;">
        <el-button :icon="Refresh" circle @click="loadList" :loading="loading" />
        <el-upload
          :action="uploadUrl"
          :headers="uploadHeaders"
          accept=".jpg,.jpeg,.png,.webp"
          :show-file-list="false"
          :on-success="handleUploadSuccess"
          :on-error="handleUploadError"
        >
          <el-button type="primary">
            <span v-html="getIconSvg('upload', 16)" style="margin-right:6px;vertical-align:middle"></span>
            上传图片
          </el-button>
        </el-upload>
      </div>
    </div>

    <div v-if="loading" class="admin-loading-state">
      <p>加载中...</p>
    </div>

    <div v-else-if="list.length === 0" class="admin-empty-state">
      暂无媒体文件
    </div>

    <div v-else class="admin-media-grid">
      <div v-for="item in list" :key="item.id" class="admin-media-card">
        <div class="admin-media-preview" @click="openPreview(item.url)">
          <img :src="item.url" :alt="item.filename" />
        </div>
        <div class="admin-media-info">
          <div class="admin-media-filename" :title="item.filename">{{ item.original_name || item.filename }}</div>
          <div class="admin-media-url" :title="item.url">{{ item.url }}</div>
          <div class="admin-media-size">{{ formatSize(item.size_bytes) }}</div>
          <el-popconfirm
            title="确定删除该文件？"
            confirm-button-text="删除"
            cancel-button-text="取消"
            @confirm="handleDelete(item.id)"
          >
            <template #reference>
              <el-button size="small" type="danger" class="delete-btn">删除</el-button>
            </template>
          </el-popconfirm>
        </div>
      </div>
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

    <!-- Image preview overlay -->
    <div v-if="previewImage" class="preview-overlay" @click.self="closePreview">
      <div class="preview-container">
        <img :src="previewImage" alt="预览" />
        <button class="preview-close" @click="closePreview" v-html="getIconSvg('x', 24)"></button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ElMessage } from 'element-plus';
import { Refresh } from '@element-plus/icons-vue';
import { getIconSvg } from '~/composables/lucideIcons';

definePageMeta({ layout: 'admin', middleware: 'auth' });

interface MediaItem {
  id: string;
  filename: string;
  original_name: string;
  url: string;
  size_bytes: number;
}

const list = ref<MediaItem[]>([]);
const loading = ref(false);
const page = ref(1);
const pageSize = ref(12);
const total = ref(0);

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

const loadList = async () => {
  loading.value = true;
  try {
    const api = useApi();
    const data = await api<{ items: MediaItem[]; total: number }>(
      `/admin/media?page=${page.value}&per_page=${pageSize.value}`
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

const handleUploadSuccess = () => {
  loadList();
};

const handleUploadError = () => {
  // error displayed by Element Plus upload component
};

const handleDelete = async (id: string) => {
  try {
    const api = useApi();
    await api(`/admin/media/${id}`, { method: 'DELETE' });
    loadList();
  } catch {
    ElMessage.error('操作失败');
  }
};

onMounted(() => {
  loadList();
});
</script>

<style scoped>
.delete-btn {
  width: 100%;
}

.admin-media-preview {
  cursor: pointer;
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
</style>
