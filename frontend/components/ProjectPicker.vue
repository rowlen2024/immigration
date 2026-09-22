<template>
  <div ref="pickerRef" class="project-picker">
    <div class="project-picker-bar">
      <template v-if="mode === 'filter'">
        <button
          type="button"
          class="filter-btn"
          :class="{ active: modelValue === ALL_PROJECT_SLUG }"
          @click="selectProject(ALL_PROJECT_SLUG)"
        >
          {{ allLabel }}
        </button>
        <button
          v-for="project in visibleProjects"
          :key="project.slug"
          type="button"
          class="filter-btn"
          :class="{ active: modelValue === project.slug }"
          @click="selectProject(project.slug)"
        >
          {{ project.name }}
        </button>
        <button
          v-if="selectedHiddenProject"
          type="button"
          class="filter-btn active selected-hidden-filter"
          @click="selectProject(selectedHiddenProject.slug)"
        >
          {{ selectedHiddenProject.name }}
        </button>
      </template>

      <template v-else>
        <NuxtLink
          v-for="project in visibleProjects"
          :key="project.slug"
          :to="`/projects/${project.slug}`"
          class="project-link"
        >
          {{ project.name }}
        </NuxtLink>
      </template>

      <button
        v-if="hiddenProjects.length > 0"
        type="button"
        class="filter-btn more-projects-btn"
        :class="{ open: showProjectPicker }"
        :aria-expanded="showProjectPicker"
        :aria-controls="panelId"
        @click="toggleProjectPicker"
      >
        更多项目
        <span class="more-projects-count">{{ hiddenProjects.length }}</span>
      </button>
    </div>

    <div
      v-if="showProjectPicker"
      :id="panelId"
      class="project-picker-panel"
    >
      <div class="project-picker-handle" aria-hidden="true"></div>
      <div class="project-picker-header">
        <div>
          <h2 class="project-picker-title">{{ title }}</h2>
          <p class="project-picker-subtitle">{{ subtitle }}</p>
        </div>
        <button
          type="button"
          class="project-picker-close"
          aria-label="关闭项目选择"
          @click="closeProjectPicker"
        >
          ×
        </button>
      </div>

      <input
        v-model.trim="projectSearch"
        class="project-picker-search"
        type="search"
        placeholder="搜索项目名称"
      >

      <div v-if="filteredHiddenProjects.length > 0" class="project-picker-grid">
        <template v-if="mode === 'filter'">
          <button
            v-for="project in filteredHiddenProjects"
            :key="project.slug"
            type="button"
            class="project-picker-option"
            :class="{ active: modelValue === project.slug }"
            @click="selectProject(project.slug)"
          >
            <span>{{ project.name }}</span>
            <span v-if="modelValue === project.slug" class="project-picker-check">✓</span>
          </button>
        </template>
        <template v-else>
          <NuxtLink
            v-for="project in filteredHiddenProjects"
            :key="project.slug"
            :to="`/projects/${project.slug}`"
            class="project-picker-option project-picker-option-link"
            @click="closeProjectPicker"
          >
            {{ project.name }}
          </NuxtLink>
        </template>
      </div>
      <div v-else class="project-picker-empty">未找到相关项目</div>
    </div>
  </div>
</template>

<script setup lang="ts">
/** 项目选择器支持的项目最小字段集合。 */
interface ProjectPickerItem {
  /** 项目的数据库标识。 */
  id: number | string
  /** 项目展示名称。 */
  name: string
  /** 项目路由或筛选使用的唯一标识。 */
  slug: string
}

/** 项目选择器的交互模式。 */
type ProjectPickerMode = 'filter' | 'links'

/** “全部项目”筛选项使用的稳定标识。 */
const ALL_PROJECT_SLUG = 'all'

/** 桌面端直接展示的项目数量。 */
const DESKTOP_VISIBLE_PROJECT_LIMIT = 5

/** 移动端直接展示的项目数量。 */
const MOBILE_VISIBLE_PROJECT_LIMIT = 1

/** 与项目页面移动端断点保持一致。 */
const MOBILE_BREAKPOINT = 767

/** 项目选择器接收的配置属性。 */
const props = withDefaults(defineProps<{
  /** 当前页面可选择或跳转的项目列表。 */
  projects: ProjectPickerItem[]
  /** 项目选择器的交互模式。 */
  mode?: ProjectPickerMode
  /** 当前筛选模式下选中的项目标识。 */
  modelValue?: string
  /** 筛选模式下“全部”按钮的显示文本。 */
  allLabel?: string
  /** 更多项目面板的标题。 */
  title?: string
  /** 更多项目面板的说明文本。 */
  subtitle?: string
}>(), {
  mode: 'filter',
  modelValue: ALL_PROJECT_SLUG,
  allLabel: '全部',
  title: '选择项目',
  subtitle: '按项目查看对应内容',
})

/** 项目选择器向父级发送筛选项目变化事件。 */
const emit = defineEmits<{
  (event: 'update:modelValue', value: string): void
}>()

/** 项目选择器根节点，用于判断面板外部点击。 */
const pickerRef = ref<HTMLElement | null>(null)

/** 是否显示更多项目面板。 */
const showProjectPicker = ref(false)

/** 更多项目面板中的搜索关键词。 */
const projectSearch = ref('')

/** 当前视口下顶部直接展示的项目数量。 */
const visibleProjectLimit = ref(DESKTOP_VISIBLE_PROJECT_LIMIT)

/** 当前注册的响应式媒体查询对象。 */
let projectFilterMedia: MediaQueryList | null = null

/** 项目选择面板的稳定 DOM ID。 */
const panelId = 'project-picker-panel'

/** 当前应直接展示的项目列表。 */
const visibleProjects = computed(() => props.projects.slice(0, visibleProjectLimit.value))

/** 需要收纳到“更多项目”面板中的项目列表。 */
const hiddenProjects = computed(() => props.projects.slice(visibleProjectLimit.value))

/** 当前选中的隐藏项目，确保筛选结果不会从顶部操作栏消失。 */
const selectedHiddenProject = computed(() => {
  if (props.mode !== 'filter' || props.modelValue === ALL_PROJECT_SLUG) return null
  const visible = visibleProjects.value.some((project) => project.slug === props.modelValue)
  if (visible) return null
  return hiddenProjects.value.find((project) => project.slug === props.modelValue) ?? null
})

/** 根据搜索关键词筛选隐藏项目。 */
const filteredHiddenProjects = computed(() => {
  const keyword = projectSearch.value.trim().toLowerCase()
  if (!keyword) return hiddenProjects.value
  return hiddenProjects.value.filter((project) => {
    return project.name.toLowerCase().includes(keyword) || project.slug.toLowerCase().includes(keyword)
  })
})

/** 根据当前视口切换顶部直接展示的项目数量。 */
function updateVisibleProjectLimit(event: MediaQueryList | MediaQueryListEvent) {
  visibleProjectLimit.value = event.matches
    ? MOBILE_VISIBLE_PROJECT_LIMIT
    : DESKTOP_VISIBLE_PROJECT_LIMIT
}

/** 打开或关闭更多项目面板。 */
function toggleProjectPicker() {
  showProjectPicker.value = !showProjectPicker.value
}

/** 关闭项目面板并清空临时搜索词。 */
function closeProjectPicker() {
  showProjectPicker.value = false
  projectSearch.value = ''
}

/** 选择筛选项目并通知父页面刷新对应内容。 */
function selectProject(slug: string) {
  if (props.mode === 'filter') emit('update:modelValue', slug)
  closeProjectPicker()
}

/** 点击面板外部时关闭项目面板。 */
function handleProjectPickerClickOutside(event: MouseEvent) {
  if (!showProjectPicker.value) return
  const target = event.target as Node
  if (pickerRef.value && !pickerRef.value.contains(target)) closeProjectPicker()
}

/** 按 Escape 键关闭项目面板。 */
function handleProjectPickerKeydown(event: KeyboardEvent) {
  if (event.key === 'Escape') closeProjectPicker()
}

onMounted(() => {
  document.addEventListener('click', handleProjectPickerClickOutside)
  document.addEventListener('keydown', handleProjectPickerKeydown)
  projectFilterMedia = window.matchMedia(`(max-width: ${MOBILE_BREAKPOINT}px)`)
  updateVisibleProjectLimit(projectFilterMedia)
  projectFilterMedia.addEventListener('change', updateVisibleProjectLimit)
})

onBeforeUnmount(() => {
  document.removeEventListener('click', handleProjectPickerClickOutside)
  document.removeEventListener('keydown', handleProjectPickerKeydown)
  projectFilterMedia?.removeEventListener('change', updateVisibleProjectLimit)
})
</script>

<style scoped>
.project-picker {
  position: relative;
}

.project-picker-bar {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 12px;
}

.filter-btn,
.project-link {
  min-height: 44px;
  padding: 10px 20px;
  border: 1.5px solid var(--color-border);
  border-radius: var(--radius-full);
  background-color: var(--bg-white);
  color: var(--color-text-secondary);
  cursor: pointer;
  font-family: var(--font-sans);
  font-size: 13px;
  font-weight: 500;
  line-height: 1.4;
  text-decoration: none;
  transition: all var(--duration-fast) var(--ease-out);
  white-space: nowrap;
  flex-shrink: 0;
}

.filter-btn:hover,
.project-link:hover {
  border-color: var(--color-accent);
  color: var(--color-accent);
}

.filter-btn.active {
  border-color: var(--color-accent);
  background-color: var(--color-accent);
  color: #fff;
}

.project-link {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  border-color: rgba(191, 154, 96, 0.28);
  background-color: var(--color-accent-soft);
  color: var(--color-accent-dark);
}

.selected-hidden-filter {
  box-shadow: 0 4px 14px rgba(191, 154, 96, 0.22);
}

.more-projects-btn {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  background-color: var(--bg-light);
  color: var(--color-text);
}

.more-projects-btn.open,
.more-projects-btn:hover {
  border-color: var(--color-accent);
  color: var(--color-accent);
}

.more-projects-count {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-width: 22px;
  height: 22px;
  padding: 0 7px;
  border-radius: var(--radius-full);
  background-color: rgba(191, 154, 96, 0.12);
  color: var(--color-accent);
  font-size: 12px;
  line-height: 1;
}

.project-picker-panel {
  position: absolute;
  top: calc(100% + 12px);
  left: 0;
  z-index: 20;
  width: min(720px, 100%);
  padding: 18px;
  border: 1px solid var(--color-border);
  border-radius: var(--radius-lg);
  background-color: var(--bg-white);
  box-shadow: 0 24px 60px rgba(32, 32, 32, 0.14);
}

.project-picker-handle {
  display: none;
}

.project-picker-header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 16px;
  margin-bottom: 14px;
}

.project-picker-title {
  margin: 0 0 4px;
  color: var(--color-text);
  font-size: 16px;
  font-weight: 700;
}

.project-picker-subtitle {
  margin: 0;
  color: var(--color-text-secondary);
  font-size: 13px;
}

.project-picker-close {
  width: 32px;
  height: 32px;
  border: 0;
  border-radius: var(--radius-full);
  background-color: var(--bg-light);
  color: var(--color-text-secondary);
  cursor: pointer;
  font-size: 20px;
  line-height: 1;
}

.project-picker-close:hover {
  color: var(--color-accent);
}

.project-picker-search {
  width: 100%;
  height: 42px;
  margin-bottom: 14px;
  padding: 0 14px;
  border: 1.5px solid var(--color-border);
  border-radius: var(--radius-md);
  outline: none;
  background-color: var(--bg-white);
  color: var(--color-text);
  font-family: var(--font-sans);
  font-size: 14px;
}

.project-picker-search:focus {
  border-color: var(--color-accent);
  box-shadow: 0 0 0 3px rgba(191, 154, 96, 0.12);
}

.project-picker-grid {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 8px;
  max-height: 280px;
  overflow-y: auto;
  padding-right: 2px;
}

.project-picker-option {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 10px;
  min-height: 44px;
  padding: 10px 12px;
  border: 1.5px solid transparent;
  border-radius: var(--radius-md);
  background-color: var(--bg-white);
  color: var(--color-text);
  cursor: pointer;
  font-family: var(--font-sans);
  font-size: 13px;
  text-align: left;
  text-decoration: none;
  transition: all var(--duration-fast) var(--ease-out);
}

.project-picker-option:hover {
  border-color: var(--color-border);
  background-color: var(--bg-light);
}

.project-picker-option.active {
  border-color: rgba(191, 154, 96, 0.45);
  background-color: rgba(191, 154, 96, 0.12);
  color: var(--color-accent);
  font-weight: 600;
}

.project-picker-option-link {
  display: flex;
}

.project-picker-check {
  flex-shrink: 0;
}

.project-picker-empty {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 88px;
  border: 1px dashed var(--color-border);
  border-radius: var(--radius-md);
  color: var(--color-text-secondary);
  font-size: 14px;
}

@media (max-width: 767px) {
  .project-picker-bar {
    gap: 8px;
  }

  .filter-btn,
  .project-link {
    padding: 8px 16px;
    min-height: 44px;
    font-size: 12px;
  }

  .project-picker-panel {
    position: fixed;
    top: auto;
    right: 0;
    bottom: 0;
    left: 0;
    width: 100%;
    max-height: 78vh;
    padding: 12px 16px 18px;
    border-right: 0;
    border-bottom: 0;
    border-left: 0;
    border-radius: 24px 24px 0 0;
    box-shadow: 0 -18px 45px rgba(32, 32, 32, 0.18);
  }

  .project-picker-handle {
    display: block;
    width: 44px;
    height: 4px;
    margin: 0 auto 14px;
    border-radius: var(--radius-full);
    background-color: var(--color-border);
  }

  .project-picker-grid {
    grid-template-columns: 1fr;
    max-height: calc(78vh - 164px);
  }

  .project-picker-title {
    font-size: 15px;
  }

  .project-picker-subtitle {
    font-size: 12px;
  }
}
</style>
