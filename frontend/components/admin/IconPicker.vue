<template>
  <div class="icon-picker">
    <div class="icon-picker__trigger" @click="dialogVisible = true">
      <span v-if="modelValue" class="icon-picker__selected">
        <span v-html="selectedIconSvg" class="icon-picker__svg"></span>
        <span class="icon-picker__name">{{ modelValue }}</span>
      </span>
      <span v-else class="icon-picker__placeholder">
        <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="icon-picker__search-icon"><circle cx="11" cy="11" r="8"/><path d="m21 21-4.35-4.35"/></svg>
        点击选择图标
      </span>
      <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="icon-picker__chevron"><polyline points="6 9 12 15 18 9"/></svg>
    </div>

    <el-dialog
      v-model="dialogVisible"
      title="选择图标"
      width="620px"
      destroy-on-close
    >
      <div class="icon-picker__dialog">
        <div class="icon-picker__tabs">
          <span
            v-for="cat in iconCategories"
            :key="cat"
            class="icon-picker__tab"
            :class="{ active: activeCategory === cat }"
            @click="activeCategory = cat"
          >{{ cat }}</span>
        </div>
        <el-input
          v-model="searchQuery"
          placeholder="搜索图标..."
          clearable
          class="icon-picker__search"
        />
        <div class="icon-picker__grid">
          <div
            v-for="icon in filteredIcons"
            :key="icon.name"
            class="icon-picker__item"
            :class="{ selected: modelValue === icon.name }"
            :title="icon.name"
            @click="select(icon.name)"
          >
            <span
              v-html="getIconSvg(icon.name, 22)"
              class="icon-picker__item-svg"
            ></span>
          </div>
        </div>
        <div v-if="modelValue" class="icon-picker__selected-name">
          已选: <strong>{{ modelValue }}</strong>
        </div>
      </div>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { lucideIcons, iconCategories, searchIcons, getIconsByCategory, getIconByName, getIconSvg } from '~/composables/lucideIcons'

const props = defineProps<{ modelValue: string }>()
const emit = defineEmits<{ 'update:modelValue': [value: string] }>()

const dialogVisible = ref(false)
const activeCategory = ref('全部')
const searchQuery = ref('')

const selectedIcon = computed(() => getIconByName(props.modelValue))
const selectedIconSvg = computed(() => getIconSvg(props.modelValue, 20))

const filteredIcons = computed(() => {
  if (searchQuery.value) return searchIcons(searchQuery.value)
  return getIconsByCategory(activeCategory.value)
})

function select(name: string) {
  emit('update:modelValue', name)
  dialogVisible.value = false
}
</script>

<style scoped>
.icon-picker__trigger {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  width: 100%;
  height: 32px;
  padding: 0 12px;
  border: 1px solid var(--color-border);
  border-radius: var(--radius-sm);
  background: var(--color-bg-surface);
  cursor: pointer;
  transition: border-color 0.2s;
  box-sizing: border-box;
}
.icon-picker__trigger:hover {
  border-color: var(--color-primary);
}
.icon-picker__selected {
  display: flex;
  align-items: center;
  gap: 8px;
  flex: 1;
}
.icon-picker__svg {
  width: 18px;
  height: 18px;
  color: var(--color-primary);
  flex-shrink: 0;
}
.icon-picker__name {
  font-size: 13px;
  color: var(--color-text);
}
.icon-picker__placeholder {
  font-size: 13px;
  color: var(--color-text-muted);
  display: flex;
  align-items: center;
  gap: 6px;
  flex: 1;
}
.icon-picker__search-icon {
  color: var(--color-text-muted);
  flex-shrink: 0;
}
.icon-picker__chevron {
  color: var(--color-text-muted);
  flex-shrink: 0;
  transition: transform 0.2s;
}
.icon-picker__trigger:hover .icon-picker__chevron {
  color: var(--color-primary);
}
.icon-picker__dialog {
  max-height: 460px;
  overflow-y: auto;
}
.icon-picker__tabs {
  display: flex;
  gap: 4px;
  margin-bottom: 12px;
  flex-wrap: wrap;
}
.icon-picker__tab {
  padding: 4px 12px;
  font-size: 12px;
  border-radius: var(--radius-sm);
  cursor: pointer;
  color: var(--color-text-secondary);
  background: var(--color-bg-app);
  transition: all 0.2s;
}
.icon-picker__tab.active {
  background: var(--color-primary);
  color: #fff;
}
.icon-picker__tab:hover:not(.active) {
  background: var(--color-border-light);
}
.icon-picker__search {
  margin-bottom: 12px;
}
.icon-picker__grid {
  display: grid;
  grid-template-columns: repeat(8, 1fr);
  gap: 6px;
}
.icon-picker__item {
  aspect-ratio: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: var(--radius-sm);
  cursor: pointer;
  border: 2px solid transparent;
  transition: all 0.15s;
}
.icon-picker__item:hover {
  background: var(--color-info-soft);
}
.icon-picker__item.selected {
  border-color: var(--color-primary);
  background: var(--color-info-soft);
}
.icon-picker__item-svg {
  width: 22px;
  height: 22px;
  color: var(--color-text-secondary);
}
.icon-picker__item.selected .icon-picker__item-svg {
  color: var(--color-primary);
}
.icon-picker__selected-name {
  margin-top: 12px;
  padding-top: 10px;
  border-top: 1px solid var(--color-border-light);
  font-size: 13px;
  color: var(--color-text-muted);
}
</style>
