<template>
  <nav v-if="items.length > 0" class="breadcrumb" aria-label="Breadcrumb">
    <!-- 桌面端：完整路径 -->
    <ol class="hidden-mobile">
      <li v-for="(item, i) in items" :key="item.href">
        <NuxtLink v-if="i < items.length - 1" :to="item.href">{{ item.label }}</NuxtLink>
        <span v-else>{{ item.label }}</span>
      </li>
    </ol>
    <!-- 移动端：仅 首页 / 目标label -->
    <ol class="hidden-desktop breadcrumb-mobile">
      <li v-if="items.length > 1"><NuxtLink :to="items[0].href">{{ items[0].label }}</NuxtLink></li>
      <li class="breadcrumb-mobile-leaf"><span>{{ leafLabel }}</span></li>
    </ol>
  </nav>
</template>

<script setup lang="ts">
import type { BreadcrumbItem } from '~/utils/breadcrumb'

const props = defineProps<{
  /** 页面传入的完整语义路径。 */
  items: readonly BreadcrumbItem[];
}>();

/** 返回当前路径的末级名称。 */
const leafLabel = computed(() => {
  const lastItem = props.items[props.items.length - 1];
  return lastItem?.label || '';
});
</script>

<style scoped>
/* 桌面/移动端 显示切换 */
@media (min-width: 768px) {
  .breadcrumb .hidden-desktop {
    display: none;
  }
}

@media (max-width: 767px) {
  .breadcrumb .hidden-mobile {
    display: none;
  }
}

.breadcrumb {
  padding: 16px 0;
}

.breadcrumb ol {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
  row-gap: 4px;
}

.breadcrumb li {
  display: flex;
  align-items: center;
  font-size: 14px;
  color: var(--text-light);
}

.breadcrumb li:not(:last-child)::after {
  content: '/';
  margin-left: 8px;
  color: var(--border-color);
}

.breadcrumb a {
  color: var(--text-secondary);
  transition: color 0.2s ease;
}

.breadcrumb a:hover {
  color: var(--primary);
}

.breadcrumb span {
  color: var(--text-primary);
}

@media (max-width: 767px) {
  .breadcrumb {
    padding: 12px 0;
  }

  .breadcrumb li {
    font-size: 13px;
  }

  .breadcrumb ol {
    gap: 4px;
  }

  .breadcrumb li:not(:last-child)::after {
    margin-left: 4px;
  }

  .breadcrumb-mobile {
    gap: 4px;
    flex-wrap: nowrap;
  }

  .breadcrumb-mobile-leaf {
    max-width: 180px;
    overflow: hidden;
    white-space: nowrap;
  }

  .breadcrumb-mobile-leaf span {
    text-overflow: ellipsis;
    overflow: hidden;
  }
}
</style>
