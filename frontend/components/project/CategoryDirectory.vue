<template>
  <div v-if="$slots.leading || canBrowse" class="category-directory" @keydown.esc="closeDirectory">
    <div class="category-directory__bar">
      <div v-if="$slots.leading" class="category-directory__leading">
        <slot name="leading" />
      </div>
      <button
        v-if="canBrowse"
        ref="trigger"
        class="category-directory__trigger"
        type="button"
        :aria-expanded="open"
        :aria-controls="panelId"
        @click="open = !open"
      >
        {{ label }}
        <span class="category-directory__chevron" aria-hidden="true"></span>
      </button>
    </div>

    <section v-if="canBrowse" v-show="open" :id="panelId" class="category-directory__panel" :aria-labelledby="headingId">
      <div class="category-directory__heading">
        <h2 :id="headingId">项目类型</h2>
        <button type="button" class="category-directory__text-button" @click="closeDirectory()">收起</button>
      </div>
      <div v-if="error && !available.length" class="category-directory__message" role="status">
        项目类型加载失败。
        <button type="button" class="category-directory__text-button" @click="refresh()">重试</button>
      </div>
      <p v-else-if="pending && !available.length" class="category-directory__message" role="status">正在加载项目类型…</p>
      <template v-else>
        <ul v-if="regularCategories.length" class="category-directory__list">
          <li v-for="category in regularCategories" :key="category.slug">
            <NuxtLink :to="`/projects/category/${category.slug}`" class="category-directory__category">
              <span>
                <strong>{{ category.name }}</strong>
                <small>{{ descriptions.get(category.slug) ?? `${category.project_count} 个项目` }}</small>
              </span>
              <span aria-hidden="true">→</span>
            </NuxtLink>
          </li>
        </ul>
        <NuxtLink
          v-if="featured"
          :to="`/projects/category/${featured.slug}`"
          class="category-directory__featured"
        >
          查看{{ featured.name }} <span aria-hidden="true">→</span>
        </NuxtLink>
      </template>
    </section>
  </div>
</template>

<script setup lang="ts">
interface PublishedCategory {
  slug: string
  name: string
  project_count: number
}

const props = defineProps<{
  label: string
  currentSlug?: string
}>()

const descriptions = new Map([
  ['investment', '按目的地查看不同投资方式'],
  ['property', '了解海外置业与身份申请'],
  ['entrepreneur', '按目的地查看创业和经营方案'],
  ['citizenship-by-donation', '按目的地查看入籍方案'],
  ['talent-and-study', '从职业背景或升学计划出发'],
  ['long-stay-visas', '查看长期工作和居留选择'],
  ['cross-border-services', '查看开户、税号等服务'],
])

const key = 'public:project-categories:list'
usePublicDataFreshness([key])
const { data, pending, error, refresh } = await useFetch<{ data: PublishedCategory[] }>('/api/v1/project-categories', { key })

const available = computed(() => (data.value?.data ?? []).filter(category => category.slug !== props.currentSlug))
const featured = computed(() => available.value.find(category => category.slug === 'hot'))
const regularCategories = computed(() => available.value.filter(category => category !== featured.value))
const canBrowse = computed(() => available.value.length > 0 || pending.value || Boolean(error.value))

const open = ref(false)
const trigger = ref<HTMLButtonElement | null>(null)
const panelId = useId()
const headingId = `${panelId}-heading`

watch(() => props.currentSlug, () => {
  open.value = false
})

function closeDirectory(event?: KeyboardEvent) {
  if (!open.value) return
  event?.preventDefault()
  open.value = false
  trigger.value?.focus()
}
</script>

<style scoped>
.category-directory__bar {
  display: flex;
  flex-wrap: wrap;
  align-items: flex-end;
  gap: 12px 16px;
}

.category-directory__leading {
  flex: 1 1 32rem;
  min-width: 0;
}

.category-directory__trigger,
.category-directory__text-button,
.category-directory__featured {
  display: inline-flex;
  align-items: center;
  justify-content: space-between;
  gap: 10px;
  min-height: 44px;
  font: inherit;
  font-weight: 600;
}

.category-directory__trigger {
  padding: 7px 14px;
  border: 1px solid var(--color-border);
  border-radius: var(--radius-sm);
  background: var(--color-bg-surface);
  color: var(--color-text);
  cursor: pointer;
}

.category-directory__trigger:hover,
.category-directory__trigger[aria-expanded='true'] {
  border-color: var(--color-accent-dark);
  background: var(--color-accent-soft);
}

.category-directory__chevron {
  width: 7px;
  height: 7px;
  margin: -3px 2px 0;
  border-right: 1.5px solid currentColor;
  border-bottom: 1.5px solid currentColor;
  transform: rotate(45deg);
}

.category-directory__trigger[aria-expanded='true'] .category-directory__chevron {
  margin-top: 3px;
  transform: rotate(225deg);
}

.category-directory__panel {
  margin-top: 12px;
  padding: 12px clamp(12px, 3vw, 24px);
  border: 1px solid var(--color-accent-dark);
  border-radius: var(--radius-md);
  background: var(--color-bg-surface);
}

.category-directory__heading {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
}

.category-directory__heading h2 {
  font-size: var(--text-base);
}

.category-directory__list {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  column-gap: 24px;
}

.category-directory__category {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 12px;
  min-height: 44px;
  padding: 8px 4px;
  border-bottom: 1px solid var(--color-border);
  color: var(--color-text);
  overflow-wrap: anywhere;
}

.category-directory__category small {
  display: block;
  color: var(--color-text-secondary);
  font-size: var(--text-sm);
}

.category-directory__category:hover {
  background: var(--color-accent-soft);
}

.category-directory__featured,
.category-directory__text-button {
  color: var(--color-text);
}

.category-directory__featured:hover,
.category-directory__text-button:hover {
  text-decoration: underline;
}

.category-directory__text-button {
  padding-inline: 6px;
  border: 0;
  background: transparent;
  cursor: pointer;
}

.category-directory__message {
  color: var(--color-text-secondary);
}

.category-directory :focus-visible {
  outline: 2px solid var(--color-accent-dark);
  outline-offset: 3px;
}

@media (max-width: 600px) {
  .category-directory__trigger {
    width: 100%;
  }

  .category-directory__list {
    grid-template-columns: minmax(0, 1fr);
  }
}
</style>
