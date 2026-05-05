<template>
  <div class="homepage-config">
    <div class="admin-page-header">
      <h2 class="admin-page-title">首页配置</h2>
    </div>

    <el-tabs v-model="activeConfigTab" type="border-card" class="homepage-tabs" v-loading="loading">
      <el-tab-pane label="轮播管理" name="slides">
      <!-- Hero Slides Card -->
      <el-card class="config-card">
        <template #header>
          <h3 class="admin-card-title">轮播管理</h3>
        </template>
        <div v-if="heroSlides.length === 0" class="admin-empty-hint">暂无轮播，点击"新增 Slide"添加。</div>
        <div v-else class="config-list">
          <div v-for="(slide, i) in heroSlides" :key="i" class="config-item">
            <img v-if="slide.image" :src="slide.image" class="slide-thumb" />
            <span v-else class="slide-label">(无图片)</span>
            <div class="config-item-actions">
              <button class="action-btn" :disabled="i === 0" @click="moveSlide(i, -1)">↑</button>
              <button class="action-btn" :disabled="i === heroSlides.length - 1" @click="moveSlide(i, 1)">↓</button>
              <button class="action-btn" @click="openEditSlide(i)">编辑</button>
              <button class="action-btn danger" @click="removeSlide(i)">删除</button>
            </div>
          </div>
        </div>
        <div class="config-list-actions">
          <el-button type="primary" size="small" @click="openAddSlide">新增 Slide</el-button>
        </div>
        <div class="card-footer">
          <el-button type="primary" :loading="slideSaving" @click="saveSlides">保存轮播</el-button>
        </div>
      </el-card>
      </el-tab-pane>

      <el-tab-pane label="项目展示区" name="showcase">
      <!-- Project Showcase Card -->
      <el-card class="config-card">
        <template #header><h3 class="admin-card-title">项目展示区</h3></template>
        <el-form label-width="100px">
          <el-form-item label="区域标题">
            <el-input v-model="projectShowcase.section_title" placeholder="精选移民项目" />
          </el-form-item>
          <el-form-item label="区域副标题">
            <el-input v-model="projectShowcase.section_subtitle" placeholder="为您量身定制的最佳移民方案" />
          </el-form-item>
          <el-form-item label="精选项目">
            <div class="featured-area">
              <div v-if="projectShowcase.featured_slugs.length === 0" class="admin-empty-hint">
                未选择精选项目，首页将展示全部项目。
              </div>
              <div v-else class="config-list">
                <div v-for="(slug, i) in projectShowcase.featured_slugs" :key="slug" class="config-item">
                  <span class="config-item-name">{{ getProjectTitle(slug) }}</span>
                  <div class="config-item-actions">
                    <button class="action-btn" :disabled="i === 0" @click="moveFeatured(i, -1)">↑</button>
                    <button class="action-btn" :disabled="i === projectShowcase.featured_slugs.length - 1" @click="moveFeatured(i, 1)">↓</button>
                    <button class="action-btn danger" @click="removeFeatured(i)">移除</button>
                  </div>
                </div>
              </div>
              <el-select
                v-if="availableProjects.length > 0"
                value=""
                placeholder="添加项目..."
                clearable
                @change="(val: string) => { if (val) addFeatured(val) }"
                class="add-project-select"
              >
                <el-option
                  v-for="p in availableProjects"
                  :key="p.slug"
                  :label="p.name"
                  :value="p.slug"
                />
              </el-select>
            </div>
          </el-form-item>
        </el-form>
        <div class="card-footer">
          <el-button type="primary" :loading="showcaseSaving" @click="saveShowcase">保存</el-button>
        </div>
      </el-card>
      </el-tab-pane>

      <el-tab-pane label="优势管理" name="advantages">
      <!-- Advantage Items Card -->
      <el-card class="config-card">
        <template #header>
          <h3 class="admin-card-title">优势管理</h3>
        </template>
        <el-form label-width="100px" class="section-form">
          <el-form-item label="区域标题">
            <el-input v-model="advantageSection.section_title" placeholder="为什么选择 北极星移民？" />
          </el-form-item>
          <el-form-item label="区域副标题">
            <el-input v-model="advantageSection.section_subtitle" placeholder="专业服务，值得信赖" />
          </el-form-item>
          <el-form-item label="区域图片">
            <ImageInput v-model="advantageSection.image" placeholder="图片地址（选填）" />
          </el-form-item>
        </el-form>
        <div v-if="advantageItems.length === 0" class="admin-empty-hint">暂无优势项，点击"新增优势项"添加。</div>
        <div v-else class="config-list">
          <div v-for="(item, i) in advantageItems" :key="i" class="config-item">
            <div class="adv-icon-preview">
              <span
                v-if="getIconByName(item.icon)"
                v-html="getIconSvg(item.icon, 18, '#c8963e')"
                class="adv-icon-svg"
              ></span>
              <span v-else class="adv-icon-emoji">{{ item.icon }}</span>
            </div>
            <div class="config-item-info">
              <strong>{{ item.title }}</strong>
              <span class="config-item-desc">{{ item.description }}</span>
            </div>
            <div class="config-item-actions">
              <button class="action-btn" :disabled="i === 0" @click="moveAdv(i, -1)">↑</button>
              <button class="action-btn" :disabled="i === advantageItems.length - 1" @click="moveAdv(i, 1)">↓</button>
              <button class="action-btn" @click="openEditAdv(i)">编辑</button>
              <button class="action-btn danger" @click="removeAdv(i)">删除</button>
            </div>
          </div>
        </div>
        <div class="config-list-actions">
          <el-button type="primary" size="small" @click="openAddAdv">新增优势项</el-button>
        </div>
        <div class="card-footer">
          <el-button type="primary" :loading="advSaving" @click="saveAdvantages">保存优势设置</el-button>
        </div>
      </el-card>
      </el-tab-pane>
    </el-tabs>

    <!-- Slide Edit Drawer -->
    <el-drawer
      v-model="slideDialogVisible"
      :title="slideEditIndex >= 0 ? '编辑 Slide' : '新增 Slide'"
      size="560px"
      destroy-on-close
    >
      <el-form label-position="top">
        <el-form-item label="背景图" required>
          <ImageInput v-model="slideForm.image" placeholder="图片地址" />
        </el-form-item>
        <el-form-item label="标题">
          <el-input v-model="slideForm.title" placeholder="主标题(可选)" />
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="slideForm.desc" placeholder="描述文案(可选)" />
        </el-form-item>
        <el-form-item label="关联项目">
          <el-select v-model="slideForm.project_slug" placeholder="(可选)" clearable>
            <el-option v-for="p in allProjects" :key="p.slug" :label="p.title" :value="p.slug" />
          </el-select>
        </el-form-item>
        <el-form-item label="背景渐变色">
          <el-input v-model="slideForm.gradient" placeholder="linear-gradient(135deg, #1a3a5c, #2d5a8e)" />
        </el-form-item>
        <el-form-item label="跳转链接">
          <el-input v-model="slideForm.link" placeholder="点击跳转链接(可选)" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="slideDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="saveSlide">确定</el-button>
      </template>
    </el-drawer>

    <!-- Advantage Edit Drawer -->
    <el-drawer
      v-model="advDialogVisible"
      :title="advEditIndex >= 0 ? '编辑优势项' : '新增优势项'"
      size="500px"
      destroy-on-close
    >
      <el-form label-position="top">
        <el-form-item label="图标">
          <IconPicker v-model="advForm.icon" />
        </el-form-item>
        <el-form-item label="标题" required>
          <el-input v-model="advForm.title" placeholder="标题" />
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="advForm.description" type="textarea" :rows="3" placeholder="描述文案" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="advDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="saveAdv">确定</el-button>
      </template>
    </el-drawer>
  </div>
</template>

<script setup lang="ts">
definePageMeta({ layout: 'admin', middleware: ['auth'] });
import ImageInput from '~/components/admin/ImageInput.vue';
import IconPicker from '~/components/admin/IconPicker.vue';
import { getIconByName, getIconSvg } from '~/composables/lucideIcons';

interface HeroSlide {
  title: string;
  desc: string;
  project_slug: string;
  gradient: string;
  image: string;
  link: string;
}

interface AdvantageItem {
  icon: string;
  icon_type: string;
  title: string;
  description: string;
}

interface ProjectShowcase {
  section_title: string;
  section_subtitle: string;
  featured_slugs: string[];
}

interface ProjectOption {
  slug: string;
  name: string;
}

const activeConfigTab = ref('slides');
const heroSlides = ref<HeroSlide[]>([]);
const advantageItems = ref<AdvantageItem[]>([]);
const advantageSection = ref<{ section_title: string; section_subtitle: string; image: string }>({
  section_title: '',
  section_subtitle: '',
  image: '',
});
const projectShowcase = ref<ProjectShowcase>({
  section_title: '',
  section_subtitle: '',
  featured_slugs: [],
});
const allProjects = ref<ProjectOption[]>([]);
const loading = ref(true);

const load = async () => {
  loading.value = true;
  try {
    const api = useApi();
    const [config, projects] = await Promise.all([
      api<{
        hero_slides: HeroSlide[];
        advantage_items: AdvantageItem[];
        advantage_section: { section_title: string; section_subtitle: string; image: string } | null;
        project_showcase: ProjectShowcase | null;
      }>('/admin/home-config'),
      api<{ items: ProjectOption[] }>('/projects'),
    ]);

    if (config) {
      heroSlides.value = config.hero_slides || [];
      advantageItems.value = config.advantage_items || [];
      if (config.advantage_section) {
        advantageSection.value = config.advantage_section;
      }
      if (config.project_showcase) {
        projectShowcase.value = config.project_showcase;
      }
    }

    if (projects?.items) {
      const seen = new Set<string>();
      allProjects.value = projects.items.filter((p) => {
        if (seen.has(p.slug)) return false;
        seen.add(p.slug);
        return true;
      });
    }
  } finally {
    loading.value = false;
  }
};

function getProjectTitle(slug: string): string {
  return allProjects.value.find((p) => p.slug === slug)?.name || slug;
}

// --- Hero Slides ---
const slideDialogVisible = ref(false);
const slideForm = ref<HeroSlide>(blankSlide());
const slideEditIndex = ref(-1);
const slideSaving = ref(false);

function blankSlide(): HeroSlide {
  return { title: '', desc: '', project_slug: '', gradient: '', image: '', link: '' };
}

function openAddSlide() {
  slideEditIndex.value = -1;
  slideForm.value = blankSlide();
  slideDialogVisible.value = true;
}

function openEditSlide(index: number) {
  slideEditIndex.value = index;
  slideForm.value = { ...heroSlides.value[index] };
  slideDialogVisible.value = true;
}

function removeSlide(index: number) {
  heroSlides.value.splice(index, 1);
}

function moveSlide(index: number, direction: -1 | 1) {
  const target = index + direction;
  if (target < 0 || target >= heroSlides.value.length) return;
  const items = [...heroSlides.value];
  [items[index], items[target]] = [items[target], items[index]];
  heroSlides.value = items;
}

async function saveSlide() {
  if (!slideForm.value.image.trim()) {
    ElMessage.warning('请上传背景图');
    return;
  }
  if (slideEditIndex.value >= 0) {
    heroSlides.value[slideEditIndex.value] = { ...slideForm.value };
  } else {
    heroSlides.value.push({ ...slideForm.value });
  }
  slideDialogVisible.value = false;
  await saveSlides();
}

async function saveSlides() {
  slideSaving.value = true;
  try {
    const api = useApi();
    await api('/admin/home-config', {
      method: 'PUT',
      body: { hero_slides: heroSlides.value },
    });
    ElMessage.success('轮播已保存');
  } catch {
    ElMessage.error('保存失败');
  } finally {
    slideSaving.value = false;
  }
}

// --- Project Showcase ---
const showcaseSaving = ref(false);

const availableProjects = computed(() => {
  const featured = new Set(projectShowcase.value.featured_slugs);
  return allProjects.value.filter((p) => !featured.has(p.slug));
});

function moveFeatured(index: number, direction: -1 | 1) {
  const target = index + direction;
  if (target < 0 || target >= projectShowcase.value.featured_slugs.length) return;
  const slugs = [...projectShowcase.value.featured_slugs];
  [slugs[index], slugs[target]] = [slugs[target], slugs[index]];
  projectShowcase.value.featured_slugs = slugs;
}

function removeFeatured(index: number) {
  projectShowcase.value.featured_slugs.splice(index, 1);
}

function addFeatured(slug: string) {
  if (!projectShowcase.value.featured_slugs.includes(slug)) {
    projectShowcase.value.featured_slugs.push(slug);
  }
}

async function saveShowcase() {
  showcaseSaving.value = true;
  try {
    const api = useApi();
    await api('/admin/home-config', {
      method: 'PUT',
      body: { project_showcase: projectShowcase.value },
    });
    ElMessage.success('项目展示区已保存');
  } catch {
    ElMessage.error('保存失败');
  } finally {
    showcaseSaving.value = false;
  }
}

// --- Advantage Items ---
const advDialogVisible = ref(false);
const advForm = ref<AdvantageItem>({ icon: '', icon_type: 'lucide', title: '', description: '' });
const advEditIndex = ref(-1);
const advSaving = ref(false);

function openAddAdv() {
  advEditIndex.value = -1;
  advForm.value = { icon: '', icon_type: 'lucide', title: '', description: '' };
  advDialogVisible.value = true;
}

function openEditAdv(index: number) {
  advEditIndex.value = index;
  advForm.value = { ...advantageItems.value[index] };
  advDialogVisible.value = true;
}

function removeAdv(index: number) {
  advantageItems.value.splice(index, 1);
}

function moveAdv(index: number, direction: -1 | 1) {
  const target = index + direction;
  if (target < 0 || target >= advantageItems.value.length) return;
  const items = [...advantageItems.value];
  [items[index], items[target]] = [items[target], items[index]];
  advantageItems.value = items;
}

async function saveAdv() {
  if (!advForm.value.title.trim()) {
    ElMessage.warning('请填写标题');
    return;
  }
  if (advEditIndex.value >= 0) {
    advantageItems.value[advEditIndex.value] = { ...advForm.value };
  } else {
    advantageItems.value.push({ ...advForm.value });
  }
  advDialogVisible.value = false;
  await saveAdvantages();
}

async function saveAdvantages() {
  advSaving.value = true;
  try {
    const api = useApi();
    await api('/admin/home-config', {
      method: 'PUT',
      body: {
        advantage_items: advantageItems.value,
        advantage_section: advantageSection.value,
      },
    });
    ElMessage.success('优势设置已保存');
  } catch {
    ElMessage.error('保存失败');
  } finally {
    advSaving.value = false;
  }
}

onMounted(load);
</script>

<style scoped>
/* Homepage tabs wrapper */
.homepage-tabs {
  background: var(--color-bg-surface);
  border-radius: var(--radius-md);
}

/* Empty hint override for homepage */
.admin-empty-hint {
  color: var(--color-text-muted);
  font-size: 14px;
  padding: 16px 0;
  text-align: center;
}
</style>
