<template>
  <div class="faq-accordion">
    <div
      v-for="(faq, index) in items"
      :key="index"
      class="faq-item"
      :class="{ active: activeIndex === index }"
    >
      <button
        class="faq-question"
        @click="toggle(index)"
        :aria-expanded="activeIndex === index"
      >
        <span>{{ faq.question }}</span>
        <span class="faq-icon" aria-hidden="true">
          <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <line x1="12" y1="5" x2="12" y2="19" />
            <line x1="5" y1="12" x2="19" y2="12" />
          </svg>
        </span>
      </button>
      <div
        class="public-disclosure"
        :class="{ 'is-expanded': activeIndex === index }"
        :aria-hidden="activeIndex !== index"
      >
        <div class="public-disclosure__body faq-answer">
          <p>{{ faq.answer }}</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
interface FaqItem {
  question: string;
  answer: string;
}

defineProps<{
  items: FaqItem[];
}>();

const activeIndex = ref<number | null>(null);

const toggle = (index: number) => {
  activeIndex.value = activeIndex.value === index ? null : index;
};
</script>

<style scoped>
.faq-accordion {
  border-top: 1px solid var(--color-border);
}

.faq-item {
  border-bottom: 1px solid var(--color-border);
}

.faq-question {
  display: flex;
  align-items: center;
  justify-content: space-between;
  width: 100%;
  padding: 20px 0;
  background: none;
  border: none;
  cursor: pointer;
  font-family: var(--font-sans);
  font-size: var(--text-base);
  font-weight: 600;
  color: var(--color-text);
  text-align: left;
  gap: 16px;
  transition: color var(--public-motion-fast) var(--ease-out);
}

.faq-question:hover {
  color: var(--color-primary);
}

.faq-icon {
  flex-shrink: 0;
  display: flex;
  align-items: center;
  color: var(--color-accent);
  transition: transform var(--public-motion-medium) var(--ease-in),
              color var(--public-motion-medium) var(--ease-in);
}

.faq-item.active .faq-icon {
  transform: rotate(45deg);
  color: var(--color-primary);
  transition-timing-function: var(--ease-out);
}

.faq-answer {
  font-size: 15px;
  color: var(--color-text-secondary);
  line-height: 1.8;
}

.faq-answer p {
  padding-bottom: 20px;
}
</style>
