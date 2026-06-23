<template>
  <div class="homepage-config">
    <div class="admin-page-header">
      <h2 class="admin-page-title">首页配置</h2>
    </div>

    <div class="homepage-tabs-loading">
      <div v-if="loading" class="homepage-loading-mask">
        <span class="homepage-loading-spinner"></span>
      </div>
      <el-tabs v-model="activeConfigTab" type="border-card" class="homepage-tabs">
      <el-tab-pane label="轮播管理" name="slides">
      <!-- Hero Slides Card -->
      <el-card class="config-card">
        <template #header>
          <h3 class="admin-card-title">轮播管理</h3>
        </template>
        <div v-if="heroSlides.length === 0" class="admin-empty-hint">暂无轮播，点击"新增 Slide"添加。</div>
        <div v-else class="config-list">
          <div v-for="(slide, i) in heroSlides" :key="i" class="config-item">
            <ResponsiveImage v-if="slide.image" :src="slide.image" variant="thumb" class="slide-thumb" />
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
            <ImageInput v-model="advantageSection.image" placeholder="图片地址（选填）" size-hint="推荐 1920×480px (约4:1 横向)" context="general" />
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

      <el-tab-pane label="案例展示区" name="cases">
        <el-card class="config-card">
          <template #header><h3 class="admin-card-title">案例展示区</h3></template>
          <el-form label-width="100px">
            <el-form-item label="区域标题">
              <el-input v-model="caseShowcase.section_title" placeholder="成功案例" />
            </el-form-item>
            <el-form-item label="区域副标题">
              <el-input v-model="caseShowcase.section_subtitle" placeholder="数百家庭的信赖之选" />
            </el-form-item>
            <el-form-item label="精选案例">
              <div class="featured-area">
                <div v-if="caseShowcase.featured_case_ids.length === 0" class="admin-empty-hint">
                  未选择精选案例，首页将展示全部案例。
                </div>
                <div v-else class="config-list">
                  <div v-for="(id, i) in caseShowcase.featured_case_ids" :key="id" class="config-item">
                    <span class="config-item-name">{{ getCaseTitle(id) }}</span>
                    <div class="config-item-actions">
                      <button class="action-btn" :disabled="i === 0" @click="moveCaseFeatured(i, -1)">↑</button>
                      <button class="action-btn" :disabled="i === caseShowcase.featured_case_ids.length - 1" @click="moveCaseFeatured(i, 1)">↓</button>
                      <button class="action-btn danger" @click="removeCaseFeatured(i)">移除</button>
                    </div>
                  </div>
                </div>
                <el-select
                  v-if="availableCases.length > 0"
                  value=""
                  placeholder="添加案例..."
                  clearable
                  @change="(val: number | string) => { if (val) addCaseFeatured(val as number) }"
                  class="add-project-select"
                >
                  <el-option
                    v-for="c in availableCases"
                    :key="c.id"
                    :label="c.name"
                    :value="c.id"
                  />
                </el-select>
              </div>
            </el-form-item>
          </el-form>
          <div class="card-footer">
            <el-button type="primary" :loading="caseSaving" @click="saveCaseShowcase">保存</el-button>
          </div>
        </el-card>
      </el-tab-pane>

      <el-tab-pane label="评价展示区" name="testimonials">
        <el-card class="config-card">
          <template #header><h3 class="admin-card-title">评价展示区</h3></template>
          <el-form label-width="100px">
            <el-form-item label="区域标题">
              <el-input v-model="testimonialShowcase.section_title" placeholder="客户评价" />
            </el-form-item>
            <el-form-item label="区域副标题">
              <el-input v-model="testimonialShowcase.section_subtitle" placeholder="来自真实客户的评价" />
            </el-form-item>
            <el-form-item label="精选评价">
              <div class="featured-area">
                <div v-if="testimonialShowcase.featured_testimonial_ids.length === 0" class="admin-empty-hint">
                  未选择精选评价，首页将展示全部评价。
                </div>
                <div v-else class="config-list">
                  <div v-for="(id, i) in testimonialShowcase.featured_testimonial_ids" :key="id" class="config-item">
                    <span class="config-item-name">{{ getTestimonialTitle(id) }}</span>
                    <div class="config-item-actions">
                      <button class="action-btn" :disabled="i === 0" @click="moveTestimonialFeatured(i, -1)">↑</button>
                      <button class="action-btn" :disabled="i === testimonialShowcase.featured_testimonial_ids.length - 1" @click="moveTestimonialFeatured(i, 1)">↓</button>
                      <button class="action-btn danger" @click="removeTestimonialFeatured(i)">移除</button>
                    </div>
                  </div>
                </div>
                <el-select
                  v-if="availableTestimonials.length > 0"
                  value=""
                  placeholder="添加评价..."
                  clearable
                  @change="(val: number | string) => { if (val) addTestimonialFeatured(val as number) }"
                  class="add-project-select"
                >
                  <el-option
                    v-for="t in availableTestimonials"
                    :key="t.id"
                    :label="t.nickname"
                    :value="t.id"
                  />
                </el-select>
              </div>
            </el-form-item>
          </el-form>
          <div class="card-footer">
            <el-button type="primary" :loading="testimonialSaving" @click="saveTestimonialShowcase">保存</el-button>
          </div>
        </el-card>
      </el-tab-pane>

      <el-tab-pane label="信任数据" name="trust">
        <el-card class="config-card">
          <template #header>
            <h3 class="admin-card-title">信任数据</h3>
          </template>
          <div v-if="trustItems.length === 0" class="admin-empty-hint">暂无信任数据，点击"新增条目"添加。</div>
          <div v-else class="config-list">
            <div v-for="(item, i) in trustItems" :key="i" class="config-item">
              <div class="config-item-info">
                <strong>{{ item.number }}</strong>
                <span class="config-item-desc">{{ item.label }}</span>
              </div>
              <div class="config-item-actions">
                <button class="action-btn" :disabled="i === 0" @click="moveTrust(i, -1)">↑</button>
                <button class="action-btn" :disabled="i === trustItems.length - 1" @click="moveTrust(i, 1)">↓</button>
                <button class="action-btn" @click="openEditTrust(i)">编辑</button>
                <button class="action-btn danger" @click="removeTrust(i)">删除</button>
              </div>
            </div>
          </div>
          <div class="config-list-actions">
            <el-button type="primary" size="small" @click="openAddTrust">新增条目</el-button>
          </div>
          <div class="card-footer">
            <el-button type="primary" :loading="trustSaving" @click="saveTrust">保存</el-button>
          </div>
        </el-card>
      </el-tab-pane>
      </el-tabs>
    </div>

    <!-- Slide Edit Drawer -->
    <el-drawer
      v-model="slideDialogVisible"
      :title="slideEditIndex >= 0 ? '编辑 Slide' : '新增 Slide'"
      size="560px"
      destroy-on-close
    >
      <el-form label-position="top">
        <el-form-item label="背景图" required>
          <ImageInput v-model="slideForm.image" placeholder="图片地址" size-hint="推荐 1920×800px (约2.4:1 横向)" context="homepage-slide" />
        </el-form-item>
        <el-form-item label="标题">
          <el-input v-model="slideForm.title" placeholder="主标题(可选)" />
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="slideForm.desc" placeholder="描述文案(可选)" />
        </el-form-item>
        <el-form-item label="关联项目">
          <el-select v-model="slideForm.project_slug" placeholder="(可选)" clearable>
            <el-option v-for="p in allProjects" :key="p.slug" :label="p.name" :value="p.slug" />
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

    <!-- Trust Edit Drawer -->
    <el-drawer
      v-model="trustDialogVisible"
      :title="trustEditIndex >= 0 ? '编辑信任数据' : '新增信任数据'"
      size="500px"
      destroy-on-close
    >
      <el-form label-position="top">
        <el-form-item label="数值" required>
          <el-input v-model="trustForm.number" placeholder="如：3,000+" />
        </el-form-item>
        <el-form-item label="标签" required>
          <el-input v-model="trustForm.label" placeholder="如：服务家庭" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="trustDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="saveTrustItem">确定</el-button>
      </template>
    </el-drawer>
  </div>
</template>

<script setup lang="ts">
definePageMeta({ layout: 'admin', middleware: ['auth'] });

const notify = useNotify();

import ImageInput from '~/components/admin/ImageInput.vue';
import IconPicker from '~/components/admin/IconPicker.vue';
import { getIconByName, getIconSvg } from '~/composables/lucideIcons';
import { useNotify } from '~/composables/useNotify';

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

interface FeaturedCaseData {
  id: number;
  name: string;
}

interface CaseShowcase {
  section_title: string;
  section_subtitle: string;
  featured_case_ids: number[];
  featured_cases?: FeaturedCaseData[];
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

const caseShowcase = ref<CaseShowcase>({
  section_title: '',
  section_subtitle: '',
  featured_case_ids: [],
});

interface CaseOption {
  id: number;
  name: string;
}

interface FeaturedTestimonialData {
  id: number;
  nickname: string;
}

interface TestimonialShowcase {
  section_title: string;
  section_subtitle: string;
  featured_testimonial_ids: number[];
  featured_testimonials?: FeaturedTestimonialData[];
}

interface TestimonialOption {
  id: number;
  nickname: string;
}

interface TrustItem {
  number: string;
  label: string;
}

const trustItems = ref<TrustItem[]>([]);
const trustDialogVisible = ref(false);
const trustForm = ref<TrustItem>({ number: '', label: '' });
const trustEditIndex = ref(-1);
const trustSaving = ref(false);

const allCases = ref<CaseOption[]>([]);
const caseSaving = ref(false);

const testimonialShowcase = ref<TestimonialShowcase>({
  section_title: '',
  section_subtitle: '',
  featured_testimonial_ids: [],
});
const allTestimonials = ref<TestimonialOption[]>([]);
const testimonialSaving = ref(false);

const loading = ref(true);

const load = async () => {
  loading.value = true;
  try {
    const api = useApi();
    const [config, projectsData, casesData, testimonialsData] = await Promise.all([
      api<{
        hero_slides: HeroSlide[];
        advantage_items: AdvantageItem[];
        advantage_section: { section_title: string; section_subtitle: string; image: string } | null;
        project_showcase: ProjectShowcase | null;
        case_showcase: CaseShowcase | null;
        testimonial_showcase: TestimonialShowcase | null;
        hero_trust: TrustItem[] | null;
      }>('/admin/home-config'),
      api<{ items: ProjectOption[] }>('/admin/projects/options?page=1&per_page=500'),
      api<{ items: CaseOption[] }>('/admin/cases/options?page=1&per_page=500'),
      api<{ items: TestimonialOption[] }>('/admin/testimonials/options?page=1&per_page=500'),
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
      if (config.case_showcase) {
        caseShowcase.value = config.case_showcase;
      }
      if (config.testimonial_showcase) {
        testimonialShowcase.value = config.testimonial_showcase;
      }
      trustItems.value = config.hero_trust || [];
    }

    if (projectsData?.items) {
      const seen = new Set<string>();
      allProjects.value = projectsData.items.filter((p) => {
        if (seen.has(p.slug)) return false;
        seen.add(p.slug);
        return true;
      });
    }

    if (casesData?.items) {
      const seen = new Set<number>();
      allCases.value = casesData.items.filter((c) => {
        if (seen.has(c.id)) return false;
        seen.add(c.id);
        return true;
      });
    }

    if (testimonialsData?.items) {
      allTestimonials.value = testimonialsData.items;
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
    notify.success('轮播已保存');
  } catch (e) {
    notify.error(e, '保存失败');
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

// --- Case Showcase ---
const availableCases = computed(() => {
  const featured = new Set(caseShowcase.value.featured_case_ids);
  return allCases.value.filter((c) => !featured.has(c.id));
});

function getCaseTitle(id: number): string {
  const featured = caseShowcase.value.featured_cases?.find((c) => c.id === id);
  if (featured?.name) return featured.name;
  return allCases.value.find((c) => c.id === id)?.name || String(id);
}

function moveCaseFeatured(index: number, direction: -1 | 1) {
  const target = index + direction;
  if (target < 0 || target >= caseShowcase.value.featured_case_ids.length) return;
  const ids = [...caseShowcase.value.featured_case_ids];
  [ids[index], ids[target]] = [ids[target], ids[index]];
  caseShowcase.value.featured_case_ids = ids;
}

function removeCaseFeatured(index: number) {
  caseShowcase.value.featured_case_ids.splice(index, 1);
}

function addCaseFeatured(id: number) {
  if (!caseShowcase.value.featured_case_ids.includes(id)) {
    caseShowcase.value.featured_case_ids.push(id);
  }
}

async function saveCaseShowcase() {
  caseSaving.value = true;
  try {
    const api = useApi();
    await api('/admin/home-config', {
      method: 'PUT',
      body: { case_showcase: caseShowcase.value },
    });
    notify.success('案例展示区已保存');
  } catch (e) {
    notify.error(e, '保存失败');
  } finally {
    caseSaving.value = false;
  }
}

// --- Testimonial Showcase ---
const availableTestimonials = computed(() => {
  const featured = new Set(testimonialShowcase.value.featured_testimonial_ids);
  return allTestimonials.value.filter((t) => !featured.has(t.id));
});

function getTestimonialTitle(id: number): string {
  const featured = testimonialShowcase.value.featured_testimonials?.find((t) => t.id === id);
  if (featured?.nickname) return featured.nickname;
  return allTestimonials.value.find((t) => t.id === id)?.nickname || String(id);
}

function moveTestimonialFeatured(index: number, direction: -1 | 1) {
  const target = index + direction;
  if (target < 0 || target >= testimonialShowcase.value.featured_testimonial_ids.length) return;
  const ids = [...testimonialShowcase.value.featured_testimonial_ids];
  [ids[index], ids[target]] = [ids[target], ids[index]];
  testimonialShowcase.value.featured_testimonial_ids = ids;
}

function removeTestimonialFeatured(index: number) {
  testimonialShowcase.value.featured_testimonial_ids.splice(index, 1);
}

function addTestimonialFeatured(id: number) {
  if (!testimonialShowcase.value.featured_testimonial_ids.includes(id)) {
    testimonialShowcase.value.featured_testimonial_ids.push(id);
  }
}

async function saveTestimonialShowcase() {
  testimonialSaving.value = true;
  try {
    const api = useApi();
    await api('/admin/home-config', {
      method: 'PUT',
      body: { testimonial_showcase: testimonialShowcase.value },
    });
    notify.success('评价展示区已保存');
  } catch (e) {
    notify.error(e, '保存失败');
  } finally {
    testimonialSaving.value = false;
  }
}

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
    notify.success('项目展示区已保存');
  } catch (e) {
    notify.error(e, '保存失败');
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
    notify.success('优势设置已保存');
  } catch (e) {
    notify.error(e, '保存失败');
  } finally {
    advSaving.value = false;
  }
}

// --- Trust Items ---
function openAddTrust() {
  trustEditIndex.value = -1;
  trustForm.value = { number: '', label: '' };
  trustDialogVisible.value = true;
}

function openEditTrust(index: number) {
  trustEditIndex.value = index;
  trustForm.value = { ...trustItems.value[index] };
  trustDialogVisible.value = true;
}

function removeTrust(index: number) {
  trustItems.value.splice(index, 1);
}

function moveTrust(index: number, direction: -1 | 1) {
  const target = index + direction;
  if (target < 0 || target >= trustItems.value.length) return;
  const items = [...trustItems.value];
  [items[index], items[target]] = [items[target], items[index]];
  trustItems.value = items;
}

async function saveTrustItem() {
  if (!trustForm.value.number.trim() || !trustForm.value.label.trim()) {
    ElMessage.warning('请填写数值和标签');
    return;
  }
  if (trustEditIndex.value >= 0) {
    trustItems.value[trustEditIndex.value] = { ...trustForm.value };
  } else {
    trustItems.value.push({ ...trustForm.value });
  }
  trustDialogVisible.value = false;
  await saveTrust();
}

async function saveTrust() {
  trustSaving.value = true;
  try {
    const api = useApi();
    await api('/admin/home-config', {
      method: 'PUT',
      body: { hero_trust: trustItems.value },
    });
    notify.success('信任数据已保存');
  } catch (e) {
    notify.error(e, '保存失败');
  } finally {
    trustSaving.value = false;
  }
}

onMounted(load);
</script>

<style scoped>
/* Homepage tabs wrapper */
.homepage-tabs-loading {
  position: relative;
  border-radius: var(--radius-md);
}

.homepage-loading-mask {
  position: absolute;
  inset: 0;
  z-index: 10;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: var(--radius-md);
  background: rgba(255, 255, 255, 0.85);
}

.homepage-loading-spinner {
  width: 32px;
  height: 32px;
  border: 3px solid var(--el-border-color-light);
  border-top-color: var(--el-color-primary);
  border-radius: 50%;
  animation: homepage-loading-spin 0.8s linear infinite;
}

@keyframes homepage-loading-spin {
  to {
    transform: rotate(360deg);
  }
}

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
