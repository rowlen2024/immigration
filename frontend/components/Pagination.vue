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
      <template v-for="(p, index) in visiblePages" :key="typeof p === 'number' ? p : `${p}-${index}`">
        <span v-if="p === PAGINATION_ELLIPSIS" class="pg-ellipsis" aria-hidden="true">…</span>
        <button
          v-else
          class="pg-num"
          :class="{ active: p === page }"
          @click="go(p)"
        >
          {{ p }}
        </button>
      </template>
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
/** 分页组件接收的页码、每页数量和总记录数。 */
const props = withDefaults(defineProps<{
  page: number;
  perPage?: number;
  total: number;
}>(), {
  perPage: 10,
});

/** 分页组件向父级发送页码变化事件。 */
const emit = defineEmits<{
  (e: 'change', page: number): void;
}>();

/** 页码列表中的省略号标识，避免把哨兵值直接渲染成负数。 */
const PAGINATION_ELLIPSIS = 'ellipsis' as const;

/** 当前数据对应的最大页码。 */
const totalPages = computed(() => Math.max(1, Math.ceil(props.total / props.perPage)));

/** 当前分页控件需要展示的页码和省略号。 */
const visiblePages = computed(() => {
  const total = totalPages.value;
  const current = Math.min(Math.max(props.page, 1), total);
  const pages: Array<number | typeof PAGINATION_ELLIPSIS> = [];

  if (total <= 7) {
    for (let i = 1; i <= total; i++) pages.push(i);
    return pages;
  }

  pages.push(1);
  if (current > 3) pages.push(PAGINATION_ELLIPSIS);
  const start = Math.max(2, current - 1);
  const end = Math.min(total - 1, current + 1);
  for (let i = start; i <= end; i++) pages.push(i);
  if (current < total - 2) pages.push(PAGINATION_ELLIPSIS);
  pages.push(total);
  return pages;
});

/** 校验并发送用户选择的目标页码。 */
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

.pg-ellipsis {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-width: 36px;
  height: 36px;
  color: var(--text-light);
  font-size: 14px;
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
    justify-content: center;
  }

  .pg-label {
    display: none;
  }

  .pg-btn {
    padding: 8px 10px;
  }

  .pg-num {
    min-width: 44px;
    height: 44px;
    font-size: 13px;
  }

  .pg-ellipsis {
    min-width: 24px;
    height: 44px;
  }
}
</style>
