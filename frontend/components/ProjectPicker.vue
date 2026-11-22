<template>
  <div class="project-picker">
    <div class="project-picker-bar">
      <template v-if="mode === 'filter'">
        <button
          type="button"
          class="filter-btn"
          :class="{ active: modelValue === ALL_PROJECT_SLUG }"
          :aria-pressed="modelValue === ALL_PROJECT_SLUG"
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
          :aria-pressed="modelValue === project.slug"
          @click="selectProject(project.slug)"
        >
          {{ project.name }}
        </button>
        <button
          v-if="selectedHiddenProject"
          type="button"
          class="filter-btn active selected-hidden-filter"
          aria-pressed="true"
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
        ref="trigger"
        type="button"
        class="filter-btn more-projects-btn"
        :class="{ open: showProjectPicker }"
        :aria-expanded="showProjectPicker"
        :aria-controls="panelId"
        aria-haspopup="dialog"
        @click="showProjectPicker = true"
      >
        更多项目
        <span class="more-projects-count">{{ hiddenProjects.length }}</span>
      </button>
    </div>

    <PublicPickerDialog
      :id="panelId"
      v-model:open="showProjectPicker"
      :anchor="trigger"
      :title="title"
      :subtitle="subtitle"
    >
      <template v-if="showProjectPicker">
        <div>
          <label :for="`${panelId}-search`" class="public-control-label">搜索项目</label>
          <input
            :id="`${panelId}-search`"
            v-model.trim="projectSearch"
            class="public-picker-search"
            type="search"
            placeholder="例如 EB-5 或香港"
            autocomplete="off"
            autofocus
          >
        </div>
        <div v-if="filteredProjects.length" class="public-picker-options">
          <template v-if="mode === 'filter'">
            <button
              v-if="!projectSearch"
              type="button"
              class="public-picker-option"
              :aria-pressed="modelValue === ALL_PROJECT_SLUG"
              @click="selectProject(ALL_PROJECT_SLUG)"
            >
              <span>全部项目</span>
              <span v-if="modelValue === ALL_PROJECT_SLUG" aria-hidden="true">✓</span>
            </button>
            <button
              v-for="project in filteredProjects"
              :key="project.slug"
              type="button"
              class="public-picker-option"
              :aria-pressed="modelValue === project.slug"
              @click="selectProject(project.slug)"
            >
              <span>{{ project.name }}</span>
              <span v-if="modelValue === project.slug" aria-hidden="true">✓</span>
            </button>
          </template>
          <template v-else>
            <NuxtLink
              v-for="project in filteredProjects"
              :key="project.slug"
              :to="`/projects/${project.slug}`"
              class="public-picker-option"
              @click="showProjectPicker = false"
            >
              {{ project.name }}
            </NuxtLink>
          </template>
        </div>
        <p v-else class="public-picker-empty" role="status">未找到相关项目</p>
      </template>
    </PublicPickerDialog>
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

const trigger = ref<HTMLButtonElement | null>(null)

/** 是否显示更多项目面板。 */
const showProjectPicker = ref(false)

/** 更多项目面板中的搜索关键词。 */
const projectSearch = ref('')

/** 当前视口下顶部直接展示的项目数量。 */
const visibleProjectLimit = ref(DESKTOP_VISIBLE_PROJECT_LIMIT)

/** 当前注册的响应式媒体查询对象。 */
let projectFilterMedia: MediaQueryList | null = null

/** 项目选择面板的稳定 DOM ID。 */
const panelId = useId()

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

/** 搜索当前页面关联的所有项目。 */
const filteredProjects = computed(() => {
  const keyword = projectSearch.value.trim().toLowerCase()
  if (!keyword) return props.projects
  return props.projects.filter((project) => {
    return project.name.toLowerCase().includes(keyword) || project.slug.toLowerCase().includes(keyword)
  })
})

/** 根据当前视口切换顶部直接展示的项目数量。 */
function updateVisibleProjectLimit(event: MediaQueryList | MediaQueryListEvent) {
  visibleProjectLimit.value = event.matches
    ? MOBILE_VISIBLE_PROJECT_LIMIT
    : DESKTOP_VISIBLE_PROJECT_LIMIT
}

/** 选择筛选项目并通知父页面刷新对应内容。 */
function selectProject(slug: string) {
  if (props.mode === 'filter') emit('update:modelValue', slug)
  showProjectPicker.value = false
}

watch(showProjectPicker, value => {
  if (!value) projectSearch.value = ''
})

onMounted(() => {
  projectFilterMedia = window.matchMedia(`(max-width: ${MOBILE_BREAKPOINT}px)`)
  updateVisibleProjectLimit(projectFilterMedia)
  projectFilterMedia.addEventListener('change', updateVisibleProjectLimit)
})

onBeforeUnmount(() => {
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
  color: var(--public-control-focus);
}

.filter-btn.active {
  border-color: var(--color-accent);
  background-color: var(--color-accent);
  color: var(--color-text);
}

.project-link {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  border-color: rgba(191, 154, 96, 0.28);
  background-color: var(--color-accent-soft);
  color: var(--public-control-focus);
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
  color: var(--public-control-focus);
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
  color: var(--public-control-focus);
  font-size: 12px;
  line-height: 1;
}

@media (max-width: 767px) {
  .project-picker-bar {
    gap: 8px;
  }

  .filter-btn,
  .project-link {
    max-width: 100%;
    padding: 8px 16px;
    font-size: 13px;
    white-space: normal;
    overflow-wrap: anywhere;
  }
}
</style>
