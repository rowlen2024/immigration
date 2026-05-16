<template>
  <div class="lawyers-admin">
    <div class="admin-page-header">
      <h2 class="admin-page-title">律师团队</h2>
      <el-button type="primary" @click="openAdd">新增律师</el-button>
    </div>

    <el-table :data="lawyers" stripe v-loading="loading" class="admin-table">
      <el-table-column label="照片" width="90">
        <template #default="{ row }">
          <img v-if="row.photo_url" :src="row.photo_url" class="lawyer-thumb" />
          <span v-else class="no-photo">—</span>
        </template>
      </el-table-column>
      <el-table-column prop="name" label="姓名" width="100" />
      <el-table-column prop="title" label="身份" min-width="140" />
      <el-table-column label="标签" min-width="200">
        <template #default="{ row }">
          <el-tag v-for="tag in row.tags" :key="tag" size="small" class="tag-chip">{{ tag }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="sort_order" label="排序" width="70" />
      <el-table-column label="操作" width="220">
        <template #default="{ row, $index }">
          <el-button size="small" @click="moveUp($index)" :disabled="$index === 0">↑</el-button>
          <el-button size="small" @click="moveDown($index)" :disabled="$index === lawyers.length - 1">↓</el-button>
          <el-button size="small" @click="openEdit(row)">编辑</el-button>
          <el-button size="small" type="danger" @click="removeLawyer(row.id)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-drawer
      v-model="dialogVisible"
      :title="editingId ? '编辑律师' : '新增律师'"
      size="500px"
      destroy-on-close
    >
      <el-form label-position="top">
        <el-form-item label="照片">
          <ImageInput v-model="form.photo_url" placeholder="照片地址" />
        </el-form-item>
        <el-form-item label="姓名" required>
          <el-input v-model="form.name" placeholder="律师姓名" />
        </el-form-item>
        <el-form-item label="身份">
          <el-input v-model="form.title" placeholder="如：首席移民律师" />
        </el-form-item>
        <el-form-item label="标签">
          <div class="tag-editor">
            <div v-for="(tag, i) in form.tags" :key="i" class="tag-editor-row">
              <el-input v-model="form.tags[i]" placeholder="标签内容" />
              <el-button type="danger" size="small" @click="form.tags.splice(i, 1)" :disabled="form.tags.length <= 1">×</el-button>
            </div>
            <el-button size="small" @click="form.tags.push('')">+ 添加标签</el-button>
          </div>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="saveLawyer">确定</el-button>
      </template>
    </el-drawer>
  </div>
</template>

<script setup lang="ts">
definePageMeta({ layout: 'admin', middleware: ['auth'] });
import ImageInput from '~/components/admin/ImageInput.vue';

interface Lawyer {
  id: number;
  photo_url: string;
  name: string;
  title: string;
  tags: string[];
  sort_order: number;
}

const lawyers = ref<Lawyer[]>([]);
const loading = ref(true);

const dialogVisible = ref(false);
const editingId = ref<number | null>(null);
const form = ref<{ photo_url: string; name: string; title: string; tags: string[] }>({
  photo_url: '',
  name: '',
  title: '',
  tags: [''],
});

const load = async () => {
  loading.value = true;
  try {
    const api = useApi();
    const data = await api<Lawyer[]>('/admin/lawyers');
    if (Array.isArray(data)) lawyers.value = data;
  } catch { /* ignore */ }
  finally { loading.value = false; }
};

function openAdd() {
  editingId.value = null;
  form.value = { photo_url: '', name: '', title: '', tags: [''] };
  dialogVisible.value = true;
}

function openEdit(row: Lawyer) {
  editingId.value = row.id;
  form.value = {
    photo_url: row.photo_url,
    name: row.name,
    title: row.title,
    tags: row.tags.length > 0 ? [...row.tags] : [''],
  };
  dialogVisible.value = true;
}

async function saveLawyer() {
  if (!form.value.name.trim()) { ElMessage.warning('请填写姓名'); return; }
  form.value.tags = form.value.tags.filter(t => t.trim() !== '');
  if (form.value.tags.length === 0) form.value.tags = [''];

  try {
    const api = useApi();
    if (editingId.value) {
      await api(`/admin/lawyers/${editingId.value}`, { method: 'PUT', body: form.value });
      ElMessage.success('已更新');
    } else {
      await api('/admin/lawyers', { method: 'POST', body: form.value });
      ElMessage.success('已添加');
    }
    dialogVisible.value = false;
    await load();
  } catch { ElMessage.error('操作失败'); }
}

async function removeLawyer(id: number) {
  try {
    await ElMessageBox.confirm('确定删除该律师？', '确认', { confirmButtonText: '删除', cancelButtonText: '取消', type: 'warning' });
    const api = useApi();
    await api(`/admin/lawyers/${id}`, { method: 'DELETE' });
    ElMessage.success('已删除');
    await load();
  } catch { /* cancelled */ }
}

async function moveUp(index: number) {
  if (index === 0) return;
  const items = [...lawyers.value];
  [items[index], items[index - 1]] = [items[index - 1], items[index]];
  lawyers.value = items;
  await syncSort();
}

async function moveDown(index: number) {
  if (index === lawyers.value.length - 1) return;
  const items = [...lawyers.value];
  [items[index], items[index + 1]] = [items[index + 1], items[index]];
  lawyers.value = items;
  await syncSort();
}

async function syncSort() {
  try {
    const api = useApi();
    for (let i = 0; i < lawyers.value.length; i++) {
      const item = lawyers.value[i];
      await api(`/admin/lawyers/${item.id}`, {
        method: 'PUT',
        body: {
          photo_url: item.photo_url,
          name: item.name,
          title: item.title,
          tags: item.tags,
          sort_order: i + 1,
        },
      });
    }
  } catch { ElMessage.error('排序更新失败'); }
}

onMounted(load);
</script>

<style scoped>
.lawyer-thumb { width: 56px; height: 64px; object-fit: cover; border-radius: 4px; }
.no-photo { color: #ccc; }
.tag-chip { margin-right: 4px; margin-bottom: 4px; }
.tag-editor { display: flex; flex-direction: column; gap: 8px; }
.tag-editor-row { display: flex; gap: 8px; align-items: center; }
</style>
