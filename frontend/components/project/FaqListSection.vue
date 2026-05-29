<template>
  <div class="fls-root">
    <!-- Loading -->
    <div v-if="loading" class="fls-skeleton">
      <PageSkeleton variant="list" :count="6" />
    </div>

    <!-- Error -->
    <div v-else-if="error" class="fls-error">{{ error }}</div>

    <!-- Content -->
    <template v-else>
      <div class="fls-list">
        <ProjectFAQAccordion :items="items" />
      </div>

      <div v-if="items.length === 0" class="fls-empty">{{ emptyText }}</div>

      <Pagination
        :page="page"
        :per-page="perPage"
        :total="total"
        @change="(p: number) => emit('pageChange', p)"
      />
    </template>
  </div>
</template>

<script setup lang="ts">
interface FaqItem {
  question: string
  answer: string
}

const props = withDefaults(defineProps<{
  items: FaqItem[]
  loading?: boolean
  error?: string | null
  emptyText?: string
  page: number
  perPage?: number
  total: number
}>(), {
  loading: false,
  error: null,
  emptyText: '暂无常见问题',
  perPage: 10,
})

const emit = defineEmits<{
  (e: 'pageChange', page: number): void
}>()
</script>

<style scoped>
.fls-skeleton {
  padding: 40px 0;
}

.fls-error {
  text-align: center;
  padding: 60px 20px;
  color: #c62828;
  font-size: 16px;
}

.fls-list {
  margin-bottom: 48px;
}

.fls-empty {
  text-align: center;
  padding: 60px 20px;
  color: var(--color-text-muted);
  font-size: 16px;
}

@media (max-width: 767px) {
  .fls-skeleton {
    padding: 24px 0;
  }

  .fls-error {
    padding: 40px 16px;
    font-size: 14px;
  }

  .fls-list {
    margin-bottom: 32px;
  }

  .fls-empty {
    padding: 40px 16px;
    font-size: 14px;
  }
}
</style>
