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
        <span class="faq-icon">{{ activeIndex === index ? '−' : '+' }}</span>
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
    const wrapper = answerRefs.value[index];
    if (wrapper) {
      wrapper.style.maxHeight = wrapper.scrollHeight + 'px';
      requestAnimationFrame(() => {
        wrapper.style.maxHeight = '0px';
      });
    }
    activeIndex.value = null;
  } else {
    const prev = activeIndex.value !== null ? answerRefs.value[activeIndex.value] : null;
    if (prev) {
      prev.style.maxHeight = '0px';
    }
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
  border-top: 1px solid var(--border-color);
}

.faq-item {
  border-bottom: 1px solid var(--border-color);
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
  font-size: 16px;
  font-weight: 600;
  color: var(--text-primary);
  text-align: left;
  gap: 16px;
}

.faq-question:hover {
  color: var(--primary);
}

.faq-icon {
  flex-shrink: 0;
  font-size: 20px;
  color: var(--accent);
  transition: transform 0.3s ease;
}

.faq-item.active .faq-icon {
  color: var(--primary);
}

.faq-answer-wrapper {
  max-height: 0;
  overflow: hidden;
  transition: max-height 0.3s ease;
}

.faq-answer {
  padding-bottom: 20px;
  font-size: 15px;
  color: var(--text-secondary);
  line-height: 1.8;
}
</style>
