<template>
  <nav class="breadcrumb" aria-label="Breadcrumb">
    <!-- 桌面端：完整路径 -->
    <ol class="hidden-mobile">
      <li><NuxtLink to="/">首页</NuxtLink></li>
      <li v-for="(item, i) in crumbs" :key="i">
        <NuxtLink v-if="item.link && i < crumbs.length - 1" :to="item.link">{{ item.label }}</NuxtLink>
        <span v-else>{{ item.label }}</span>
      </li>
    </ol>
    <!-- 移动端：仅 首页 / 目标label -->
    <ol class="hidden-desktop breadcrumb-mobile">
      <li><NuxtLink to="/">首页</NuxtLink></li>
      <li class="breadcrumb-mobile-leaf"><span>{{ leafLabel }}</span></li>
    </ol>
  </nav>
</template>

<script setup lang="ts">
const props = defineProps<{
  label?: string;
  parentLabel?: string;
  parentLink?: string;
}>();

const route = useRoute();
const { getBreadcrumb } = useNavigation();

const parentCrumb = computed(() =>
  props.parentLabel
    ? { label: props.parentLabel, link: props.parentLink || undefined }
    : undefined,
);

const crumbs = computed(() => getBreadcrumb(route.path, props.label, parentCrumb.value));

const leafLabel = computed(() => {
  const c = crumbs.value;
  if (c.length > 0) return c[c.length - 1].label;
  return props.label || '';
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
