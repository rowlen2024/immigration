<template>
  <nav v-if="totalPages > 1" class="pg-root">
    <button
      class="pg-btn"
      :disabled="page <= 1"
      @click="go(page - 1)"
    >
      <span class="pg-arrow">&laquo;</span>
      <span class="pg-label">上一页</span>
    </button>

    <div class="pg-pages">
      <button
        v-for="p in visiblePages"
        :key="p"
        class="pg-num"
        :class="{ active: p === page }"
        @click="go(p)"
      >
        {{ p }}
      </button>
    </div>

    <button
      class="pg-btn"
      :disabled="page >= totalPages"
      @click="go(page + 1)"
    >
      <span class="pg-label">下一页</span>
      <span class="pg-arrow">&raquo;</span>
    </button>
  </nav>
</template>

<script setup lang="ts">
const props = withDefaults(defineProps<{
  page: number;
  perPage?: number;
  total: number;
}>(), {
  perPage: 10,
});

const emit = defineEmits<{
  (e: 'change', page: number): void;
}>();

const totalPages = computed(() => Math.max(1, Math.ceil(props.total / props.perPage)));

const visiblePages = computed(() => {
  const total = totalPages.value;
  const current = props.page;
  const pages: number[] = [];

  if (total <= 7) {
    for (let i = 1; i <= total; i++) pages.push(i);
    return pages;
  }

  pages.push(1);
  if (current > 3) pages.push(-1); // ellipsis
  const start = Math.max(2, current - 1);
  const end = Math.min(total - 1, current + 1);
  for (let i = start; i <= end; i++) pages.push(i);
  if (current < total - 2) pages.push(-1);
  pages.push(total);
  return pages;
});

const go = (p: number) => {
  if (p < 1 || p > totalPages.value || p === props.page) return;
  emit('change', p);
};
</script>

<style scoped>
.pg-root {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  gap: 4px;
  padding: 32px 0 48px;
  user-select: none;
}

.pg-btn {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 8px 14px;
  font-size: 14px;
  font-family: var(--font-sans);
  font-weight: 500;
  color: var(--text-secondary);
  background: transparent;
  border: 1px solid var(--border-color);
  border-radius: var(--radius-md);
  cursor: pointer;
  transition: all 0.2s ease;
}

.pg-btn:hover:not(:disabled) {
  color: var(--accent-dark);
  border-color: var(--accent);
}

.pg-btn:disabled {
  opacity: 0.35;
  cursor: not-allowed;
}

.pg-arrow {
  font-size: 16px;
  line-height: 1;
}

.pg-pages {
  display: flex;
  align-items: center;
  gap: 2px;
  margin: 0 4px;
}

.pg-num {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-width: 36px;
  height: 36px;
  padding: 0 6px;
  font-size: 14px;
  font-family: var(--font-sans);
  font-weight: 500;
  color: var(--text-secondary);
  background: transparent;
  border: none;
  border-radius: var(--radius-md);
  cursor: pointer;
  transition: all 0.2s ease;
}

.pg-num:hover {
  color: var(--accent-dark);
  background: var(--bg-light);
}

.pg-num.active {
  color: var(--bg-white);
  background: var(--primary);
}

@media (max-width: 767px) {
  .pg-root {
    gap: 2px;
  }

  .pg-label {
    display: none;
  }

  .pg-btn {
    padding: 8px 10px;
  }

  .pg-num {
    min-width: 32px;
    height: 32px;
    font-size: 13px;
  }
}
</style>
