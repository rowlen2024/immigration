<template>
  <div class="quick-facts">
    <div v-for="(fact, index) in facts" :key="index" class="fact-card">
      <div class="fact-icon">
        <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round">
          <g v-if="index === 0">
            <polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2"/>
          </g>
          <g v-else-if="index === 1">
            <circle cx="12" cy="12" r="10"/>
            <polyline points="12 6 12 12 16 14"/>
          </g>
          <g v-else>
            <path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"/>
            <circle cx="9" cy="7" r="4"/>
            <path d="M23 21v-2a4 4 0 0 0-3-3.87"/>
            <path d="M16 3.13a4 4 0 0 1 0 7.75"/>
          </g>
        </svg>
      </div>
      <div class="fact-body">
        <span class="fact-value">{{ fact.value }}</span>
        <span class="fact-label">{{ fact.label }}</span>
      </div>
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
  grid-template-columns: repeat(3, 1fr);
  gap: 20px;
}

.fact-card {
  display: flex;
  align-items: center;
  gap: 16px;
  background: var(--bg-white);
  border: 1px solid var(--border-color);
  border-radius: var(--radius-lg);
  padding: 24px;
  transition: box-shadow 0.3s ease, border-color 0.3s ease;
}

.fact-card:hover {
  box-shadow: var(--shadow-md);
  border-color: rgba(200, 150, 62, 0.25);
}

.fact-card:active {
  transform: scale(0.985);
  background: var(--bg-light);
}

.fact-icon {
  flex-shrink: 0;
  width: 44px;
  height: 44px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: var(--radius-md);
  background: rgba(200, 150, 62, 0.1);
  color: var(--accent);
}

.fact-body {
  display: flex;
  flex-direction: column;
  gap: 4px;
  min-width: 0;
}

.fact-value {
  font-size: 20px;
  font-weight: 700;
  color: var(--text-primary);
  line-height: 1.2;
  word-break: break-word;
}

.fact-label {
  font-size: 13px;
  color: var(--text-light);
  line-height: 1.4;
}

@media (max-width: 767px) {
  .quick-facts {
    grid-template-columns: 1fr;
    gap: 12px;
  }

  .fact-card {
    padding: 18px 20px;
  }

  .fact-icon {
    width: 40px;
    height: 40px;
  }

  .fact-value {
    font-size: 18px;
  }
}
</style>
