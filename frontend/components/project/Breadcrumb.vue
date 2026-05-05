<template>
  <nav class="breadcrumb" aria-label="Breadcrumb">
    <ol>
      <li><NuxtLink to="/">首页</NuxtLink></li>
      <li v-for="(item, i) in crumbs" :key="i">
        <NuxtLink v-if="item.link && i < crumbs.length - 1" :to="item.link">{{ item.label }}</NuxtLink>
        <span v-else>{{ item.label }}</span>
      </li>
    </ol>
  </nav>
</template>

<script setup lang="ts">
const props = defineProps<{
  label?: string;
}>();

const route = useRoute();
const { getBreadcrumb } = useNavigation();

const crumbs = computed(() => getBreadcrumb(route.path, props.label));
</script>

<style scoped>
.breadcrumb {
  padding: 16px 0;
}

.breadcrumb ol {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
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
</style>
