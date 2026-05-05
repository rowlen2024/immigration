<template>
  <div class="quick-facts">
    <div v-for="(fact, index) in facts" :key="index" class="fact-card">
      <div class="fact-value">{{ fact.value }}</div>
      <div class="fact-label">{{ fact.label }}</div>
    </div>
  </div>
</template>

<script setup lang="ts">
interface QuickFact {
  label: string;
  value: string;
}

defineProps<{
  facts: QuickFact[];
}>();
</script>

<style scoped>
.quick-facts {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 20px;
}

.fact-card {
  background-color: var(--bg-light);
  border: 1px solid var(--border-color);
  border-radius: var(--radius-lg);
  padding: 28px 24px;
  text-align: center;
  transition: box-shadow 0.3s ease;
}

.fact-card:hover {
  box-shadow: var(--shadow-md);
}

.fact-value {
  font-size: 24px;
  font-weight: 700;
  color: var(--accent-dark);
  margin-bottom: 8px;
}

.fact-label {
  font-size: 14px;
  color: var(--text-secondary);
}
</style>
