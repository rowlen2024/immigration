<template>
  <div>
    <div class="admin-page-header">
      <h2 class="admin-page-title">媒体库</h2>
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

    <div class="admin-toolbar">
      <el-input
        v-model="searchQuery"
        placeholder="搜索文件名..."
        :prefix-icon="Search"
        clearable
        class="admin-search-input"
        @input="onSearch"
      />
      <el-button :icon="Refresh" circle @click="searchQuery='';loadList()" :loading="loading" />
    </div>

    <div class="admin-table-wrap">
      <el-table :data="list" v-loading="loading">
        <el-table-column label="缩略图" width="80">
          <template #default="{ row }">
            <img
              :src="row.url"
              :alt="row.filename"
              class="media-thumb"
              @click="openPreview(row.url)"
            />
          </template>
        </el-table-column>
        <el-table-column prop="original_name" label="文件名" min-width="180">
          <template #default="{ row }">
            <div class="row-title">{{ row.original_name || row.filename }}</div>
          </template>
        </el-table-column>
        <el-table-column label="URL" min-width="240">
          <template #default="{ row }">
            <div class="media-url-cell" :title="row.url">{{ row.url }}</div>
          </template>
        </el-table-column>
        <el-table-column label="大小" width="100">
          <template #default="{ row }">
            {{ formatSize(row.size_bytes) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="90">
          <template #default="{ row }">
            <el-popconfirm
              title="确定删除该文件？"
              confirm-button-text="删除"
              cancel-button-text="取消"
              @confirm="handleDelete(row.id)"
            >
              <template #reference>
                <el-button size="small" type="danger">删除</el-button>
              </template>
            </el-popconfirm>
          </template>
        </el-table-column>
      </el-table>
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
        <img :src="previewImage" alt="预览" />
        <button class="preview-close" @click="closePreview" v-html="getIconSvg('x', 24)"></button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ElMessage } from 'element-plus';
import { Refresh, Search } from '@element-plus/icons-vue';
import { useNotify } from '~/composables/useNotify';
import { getIconSvg } from '~/composables/lucideIcons';

definePageMeta({ layout: 'admin', middleware: 'auth' });

const notify = useNotify();

interface MediaItem {
  id: string;
  filename: string;
  original_name: string;
  url: string;
  size_bytes: number;
}

const list = ref<MediaItem[]>([]);
const loading = ref(false);
const searchQuery = ref('');
const page = ref(1);
const pageSize = ref(12);
const total = ref(0);
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

const handleUploadSuccess = () => {
  notify.success('上传成功');
  loadList();
};

const handleUploadError = () => {
  // error displayed by Element Plus upload component
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

onMounted(() => {
  loadList();
});
</script>

<style scoped>
.media-thumb {
  width: 56px;
  height: 56px;
  object-fit: cover;
  border-radius: 6px;
  cursor: pointer;
  border: 1px solid var(--border-color);
}

.media-url-cell {
  max-width: 300px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  font-size: 12px;
  color: var(--color-text-muted);
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
