<template>
  <div class="icon-picker">
    <div class="icon-picker__trigger" @click="dialogVisible = true">
      <div v-if="modelValue" class="icon-picker__preview">
        <span
          v-html="selectedIconSvg"
          class="icon-picker__svg"
        ></span>
      </div>
      <span v-else class="icon-picker__placeholder">选择图标</span>
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
const selectedIconSvg = computed(() => getIconSvg(props.modelValue, 20, '#c8963e'))

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
  width: 44px;
  height: 44px;
  border: 2px dashed #dcdfe6;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: border-color 0.2s;
}
.icon-picker__trigger:hover {
  border-color: #c8963e;
}
.icon-picker__preview {
  width: 36px;
  height: 36px;
  background: rgba(26, 58, 92, 0.08);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
}
.icon-picker__svg {
  width: 20px;
  height: 20px;
  color: #c8963e;
}
.icon-picker__placeholder {
  font-size: 11px;
  color: #c0c4cc;
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
  border-radius: 4px;
  cursor: pointer;
  color: #606266;
  background: #f5f7fa;
  transition: all 0.2s;
}
.icon-picker__tab.active {
  background: #1a3a5c;
  color: #fff;
}
.icon-picker__tab:hover:not(.active) {
  background: #e8eaed;
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
  border-radius: 6px;
  cursor: pointer;
  border: 2px solid transparent;
  transition: all 0.15s;
}
.icon-picker__item:hover {
  background: rgba(26, 58, 92, 0.06);
}
.icon-picker__item.selected {
  border-color: #c8963e;
  background: rgba(200, 150, 62, 0.1);
}
.icon-picker__item-svg {
  width: 22px;
  height: 22px;
  color: #606266;
}
.icon-picker__item.selected .icon-picker__item-svg {
  color: #c8963e;
}
.icon-picker__selected-name {
  margin-top: 12px;
  padding-top: 10px;
  border-top: 1px solid #ebeef5;
  font-size: 13px;
  color: #909399;
}
</style>
