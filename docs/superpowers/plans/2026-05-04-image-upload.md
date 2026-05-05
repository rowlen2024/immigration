# 图片上传功能统一改造 — 实现计划

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 将所有管理页面的图片 URL 输入框统一改造为「上传 + 媒体库浏览」组件，修复媒体库 NaN 和 URL 显示问题。

**Architecture:** 新增两个可复用组件 — `ImageInput.vue`（替代 `el-input` 的图片字段组件）和 `MediaPicker.vue`（媒体库选择对话框），后端补充 `search` 查询参数支持。ImageInput 通过 v-model 与父组件通信，内部调用 MediaPicker 对话框。

**Tech Stack:** Vue 3 + Element Plus 2.8 + Nuxt 3 (SPA), Go + Gin + GORM

---

### Task 1: 后端 — 媒体列表支持 search 参数

**Files:**
- Modify: `backend/internal/repository/media_repo.go:13-22`
- Modify: `backend/internal/service/media_svc.go:17-40`
- Modify: `backend/internal/handler/media_handler.go:71-81`

- [ ] **Step 1: 给 FindAll 添加 search 参数**

修改 `backend/internal/repository/media_repo.go` 的 `FindAll` 方法：

```go
func (r *MediaRepo) FindAll(search string) ([]model.Media, error) {
	var media []model.Media
	q := r.db.Order("created_at desc")
	if search != "" {
		like := "%" + search + "%"
		q = q.Where("filename LIKE ? OR original_name LIKE ?", like, like)
	}
	err := q.Find(&media).Error
	if err != nil {
		return nil, err
	}
	return media, nil
}
```

- [ ] **Step 2: 给 Service List 添加 search 参数**

修改 `backend/internal/service/media_svc.go` 的 `List` 方法签名和调用：

```go
func (s *MediaService) List(page, perPage int, search string) ([]model.Media, int64, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 10
	}

	media, err := s.repo.FindAll(search)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list media: %w", err)
	}

	total := int64(len(media))
	start := (page - 1) * perPage
	if start >= len(media) {
		return []model.Media{}, total, nil
	}
	end := start + perPage
	if end > len(media) {
		end = len(media)
	}
	return media[start:end], total, nil
}
```

- [ ] **Step 3: Handler 读取 search 查询参数**

修改 `backend/internal/handler/media_handler.go` 的 `ListMedia` 方法：

```go
func (h *Handler) ListMedia(c *gin.Context) {
	page, perPage := parsePagination(c)
	search := c.Query("search")

	mediaList, total, err := h.svc.Media.List(page, perPage, search)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Error(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessPaginated(mediaList, page, perPage, total))
}
```

- [ ] **Step 4: 验证后端编译通过**

Run: `cd backend && go build ./...`
Expected: 编译成功，无报错

- [ ] **Step 5: 运行后端测试**

Run: `cd backend && go test ./... -v`
Expected: 全部测试通过

- [ ] **Step 6: Commit**

```bash
git add backend/internal/repository/media_repo.go backend/internal/service/media_svc.go backend/internal/handler/media_handler.go
git commit -m "feat: add search query param to media list API"
```

---

### Task 2: 创建 MediaPicker.vue 组件

**Files:**
- Create: `frontend/components/admin/MediaPicker.vue`

- [ ] **Step 1: 创建 MediaPicker.vue**

写入 `frontend/components/admin/MediaPicker.vue`：

```vue
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
          <img :src="item.url" :alt="item.original_name" />
        </div>
        <div v-if="!loading && list.length === 0" class="empty-hint">
          暂无图片
        </div>
      </div>

      <div class="picker-detail" v-if="selectedItem">
        <div class="detail-preview">
          <img :src="selectedItem.url" :alt="selectedItem.original_name" />
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
```

- [ ] **Step 2: Commit**

```bash
git add frontend/components/admin/MediaPicker.vue
git commit -m "feat: add MediaPicker component with search, upload, and detail panel"
```

---

### Task 3: 创建 ImageInput.vue 组件

**Files:**
- Create: `frontend/components/admin/ImageInput.vue`

- [ ] **Step 1: 创建 ImageInput.vue**

写入 `frontend/components/admin/ImageInput.vue`：

```vue
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
    <div v-if="urlValue" class="preview">
      <img :src="urlValue" alt="预览" @error="previewError = true" v-show="!previewError" />
    </div>

    <MediaPicker
      v-model="pickerVisible"
      @select="onPickerSelect"
    />
  </div>
</template>

<script setup lang="ts">
import MediaPicker from './MediaPicker.vue';

const props = defineProps<{
  modelValue: string;
  placeholder?: string;
}>();

const emit = defineEmits<{
  'update:modelValue': [value: string];
}>();

const urlValue = ref(props.modelValue);
const pickerVisible = ref(false);
const previewError = ref(false);

const uploadUrl = '/api/v1/admin/media/upload';

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
  height: 68px;
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
</style>
```

- [ ] **Step 2: Commit**

```bash
git add frontend/components/admin/ImageInput.vue
git commit -m "feat: add ImageInput component with upload, browse, and preview"
```

---

### Task 4: 修复 media.vue 的 NaN 和缺少路径

**Files:**
- Modify: `frontend/pages/admin/media.vue:27-32,64-69`

- [ ] **Step 1: 修复字段名和添加 URL 显示**

修改 `frontend/pages/admin/media.vue`：

**接口定义 (line 64-69)** — `size` 改为 `size_bytes`：
```typescript
interface MediaItem {
  id: string;
  filename: string;
  original_name: string;
  url: string;
  size_bytes: number;
}
```

**模板 (lines 27-32)** — 将 `media-info` div 改为：
```html
<div class="media-info">
  <div class="media-filename" :title="item.filename">{{ item.original_name || item.filename }}</div>
  <div class="media-url" :title="item.url">{{ item.url }}</div>
  <div class="media-size">{{ formatSize(item.size_bytes) }}</div>
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
```

**formatSize (line 84-88)** — 添加空值保护：
```typescript
const formatSize = (bytes: number) => {
  if (!bytes) return '未知大小';
  if (bytes < 1024) return `${bytes} B`;
  if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`;
  return `${(bytes / (1024 * 1024)).toFixed(1)} MB`;
};
```

**新增样式** — 在 `<style scoped>` 中添加 `.media-url`：
```css
.media-url {
  font-size: 11px;
  color: #c0c4cc;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  margin-bottom: 4px;
}
```

- [ ] **Step 2: Commit**

```bash
git add frontend/pages/admin/media.vue
git commit -m "fix: fix NaN file size and add URL display in media library"
```

---

### Task 5: 改造 homepage.vue 轮播背景图

**Files:**
- Modify: `frontend/pages/admin/homepage.vue:125-127`

- [ ] **Step 1: 替换 el-input 为 ImageInput**

修改 `frontend/pages/admin/homepage.vue`：

**模板 (line 125-127)** — 将：
```html
<el-form-item label="背景图 URL">
  <el-input v-model="slideForm.image" placeholder="图片地址" />
</el-form-item>
```
替换为：
```html
<el-form-item label="背景图">
  <ImageInput v-model="slideForm.image" placeholder="图片地址" />
</el-form-item>
```

**script setup** — 添加 import（在 `<script setup>` 顶部已有其他 import 之后）：
```typescript
import ImageInput from '~/components/admin/ImageInput.vue';
```

- [ ] **Step 2: Commit**

```bash
git add frontend/pages/admin/homepage.vue
git commit -m "feat: replace carousel image URL input with ImageInput component"
```

---

### Task 6: 改造 projects.vue 封面图片

**Files:**
- Modify: `frontend/pages/admin/projects.vue:131-143,575-589`

- [ ] **Step 1: 替换内联上传为 ImageInput**

修改 `frontend/pages/admin/projects.vue`：

**模板 (lines 131-143)** — 将：
```html
<el-form-item label="封面图片">
  <div style="display: flex; gap: 8px; width: 100%">
    <el-input v-model="form.cover_image" placeholder="图片 URL 或上传" style="flex: 1" />
    <el-upload
      :action="uploadUrl"
      :headers="uploadHeaders"
      accept=".jpg,.jpeg,.png,.webp"
      :show-file-list="false"
      :on-success="handleCoverUploadSuccess"
    >
      <el-button>上传</el-button>
    </el-upload>
  </div>
</el-form-item>
```
替换为：
```html
<el-form-item label="封面图片">
  <ImageInput v-model="form.cover_image" placeholder="图片 URL 或上传" />
</el-form-item>
```

**script setup** — 添加 import：
```typescript
import ImageInput from '~/components/admin/ImageInput.vue';
```

**script setup** — 删除 `handleCoverUploadSuccess` (line 583-589)：
```typescript
// 删除这个函数，ImageInput 内部已处理
```

如果 `uploadUrl` 和 `uploadHeaders` 仅用于封面图片（未被其他字段使用），也可以一并删除。

- [ ] **Step 2: Commit**

```bash
git add frontend/pages/admin/projects.vue
git commit -m "feat: replace inline upload with ImageInput in project form"
```

---

### Task 7: 改造 settings.vue 图片字段

**Files:**
- Modify: `frontend/pages/admin/settings.vue:48,160-173,175-229`

- [ ] **Step 1: 添加 image 类型支持**

修改 `frontend/pages/admin/settings.vue`：

**FieldDef 接口 (line 160-167)** — 添加 `image` 属性：
```typescript
interface FieldDef {
  key: string;
  label: string;
  placeholder?: string;
  textarea?: boolean;
  rows?: number;
  image?: boolean;
  tip: string;
}
```

**模板 (line 48)** — 在 `<el-input v-else>` 之前插入：
```html
<ImageInput
  v-else-if="field.image"
  v-model="form[field.key]"
  :placeholder="field.placeholder"
/>
```

**groups 数组 (lines 178,181,190,199)** — 给四个图片字段添加 `image: true`：

```typescript
{ key: 'site_logo', label: '网站 Logo', image: true, placeholder: '/images/logo.png', tip: tips.site_logo },
{ key: 'site_favicon', label: 'Favicon', image: true, placeholder: '/favicon.ico', tip: tips.site_favicon },
```
```typescript
{ key: 'og_image', label: 'OG 分享图', image: true, tip: tips.og_image },
```
```typescript
{ key: 'organization_logo', label: '机构 Logo', image: true, tip: tips.organization_logo },
```

**script setup** — 添加 import：
```typescript
import ImageInput from '~/components/admin/ImageInput.vue';
```

- [ ] **Step 2: Commit**

```bash
git add frontend/pages/admin/settings.vue
git commit -m "feat: replace image URL inputs with ImageInput in settings page"
```

---

### Task 8: 端到端验证

- [ ] **Step 1: 检查前端类型**

Run: `cd frontend && npx nuxi typecheck`
Expected: 类型检查通过，无新增错误

- [ ] **Step 2: 检查前端构建**

Run: `cd frontend && npx nuxi build`
Expected: 构建成功

- [ ] **Step 3: 手动验证清单**

启动开发环境后，逐项验证：

| 页面 | 验证项 |
|------|--------|
| 媒体库 `/admin/media` | 文件大小正常显示，URL 可见 |
| 媒体库 `/admin/media` | 上传、搜索、删除功能正常 |
| 首页配置 `/admin/homepage` | 轮播背景图：手动输入 / 上传 / 浏览选择 均正常 |
| 项目编辑 `/admin/projects` | 封面图片：三种方式均正常 |
| 网站设置 `/admin/settings` | Logo/Favicon/OG图/机构Logo：三种方式均正常 |
