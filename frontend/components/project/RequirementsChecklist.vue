<template>
  <ul class="requirements-checklist">
    <li
      v-for="(item, index) in items"
      :key="index"
      class="requirement-item"
      :class="{ met: item.met, 'not-met': !item.met }"
    >
      <span class="requirement-icon">
        <span v-if="item.met" class="icon-check">&#10003;</span>
        <span v-else class="icon-cross">&#10007;</span>
      </span>
      <span class="requirement-text">{{ item.text }}</span>
    </li>
  </ul>
</template>

<script setup lang="ts">
interface RequirementItem {
  text: string;
  met: boolean;
}

defineProps<{
  items: RequirementItem[];
}>();
</script>

<style scoped>
.requirements-checklist {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.requirement-item {
  display: flex;
  align-items: flex-start;
  gap: 12px;
  padding: 12px 16px;
  border-radius: var(--radius-md);
  background-color: var(--bg-light);
}

.requirement-icon {
  flex-shrink: 0;
  width: 24px;
  height: 24px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  font-size: 12px;
  font-weight: 700;
}

.requirement-item.met .requirement-icon {
  background-color: #e6f4ea;
  color: #1e7e34;
}

.requirement-item.not-met .requirement-icon {
  background-color: #fce8e8;
  color: #c62828;
}

.requirement-text {
  font-size: 15px;
  color: var(--text-primary);
  line-height: 1.6;
  padding-top: 2px;
}
</style>
