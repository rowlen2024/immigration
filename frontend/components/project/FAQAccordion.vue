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
        <span class="faq-icon">
          <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <line x1="12" y1="5" x2="12" y2="19" />
            <line x1="5" y1="12" x2="19" y2="12" />
          </svg>
        </span>
      </button>
      <div class="faq-answer-wrapper" :ref="(el) => setAnswerRef(el as HTMLElement, index)">
        <div class="faq-answer">
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
const answerRefs = ref<Record<number, HTMLElement | null>>({});

const setAnswerRef = (el: HTMLElement | null, index: number) => {
  if (el) answerRefs.value[index] = el;
};

const toggle = (index: number) => {
  if (activeIndex.value === index) {
    // 收起当前
    const wrapper = answerRefs.value[index];
    if (wrapper) {
      wrapper.style.maxHeight = wrapper.scrollHeight + 'px';
      wrapper.offsetHeight;
      wrapper.style.maxHeight = '0px';
    }
    activeIndex.value = null;
  } else {
    // 收起上一个
    if (activeIndex.value !== null) {
      const prev = answerRefs.value[activeIndex.value];
      if (prev) {
        prev.style.maxHeight = prev.scrollHeight + 'px';
        prev.offsetHeight;
        prev.style.maxHeight = '0px';
      }
    }
    // 展开当前
    activeIndex.value = index;
    const wrapper = answerRefs.value[index];
    if (wrapper) {
      wrapper.style.maxHeight = wrapper.scrollHeight + 'px';
    }
  }
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
  transition: color var(--duration-fast) var(--ease-out);
}

.faq-question:hover {
  color: var(--color-primary);
}

.faq-icon {
  flex-shrink: 0;
  display: flex;
  align-items: center;
  color: var(--color-accent);
  transition: transform var(--duration-normal) var(--ease-out);
}

.faq-item.active .faq-icon {
  transform: rotate(45deg);
  color: var(--color-primary);
}

.faq-answer-wrapper {
  max-height: 0;
  overflow: hidden;
  transition: max-height 400ms var(--ease-out);
}

.faq-answer {
  padding-bottom: 20px;
  font-size: 15px;
  color: var(--color-text-secondary);
  line-height: 1.8;
}
</style>
