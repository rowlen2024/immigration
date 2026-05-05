<template>
  <div class="timeline">
    <div
      v-for="(phase, index) in phases"
      :key="index"
      class="timeline-item"
      :class="{ 'is-last': index === phases.length - 1 }"
    >
      <div class="timeline-marker">
        <div class="timeline-dot"></div>
        <div v-if="index < phases.length - 1" class="timeline-line"></div>
      </div>
      <div class="timeline-content">
        <div class="timeline-phase">{{ phase.phase }}</div>
        <h4 class="timeline-title">{{ phase.title }}</h4>
        <p class="timeline-desc">{{ phase.description }}</p>
        <div v-if="phase.period" class="timeline-period">{{ phase.period }}</div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
interface TimelinePhase {
  phase: string;
  title: string;
  description: string;
  period?: string;
}

defineProps<{
  phases: TimelinePhase[];
}>();
</script>

<style scoped>
.timeline {
  position: relative;
  padding-left: 40px;
}

.timeline-item {
  position: relative;
  padding-bottom: 32px;
  display: flex;
  gap: 24px;
}

.timeline-item.is-last {
  padding-bottom: 0;
}

.timeline-marker {
  position: absolute;
  left: -40px;
  top: 0;
  bottom: 0;
  display: flex;
  flex-direction: column;
  align-items: center;
}

.timeline-dot {
  width: 16px;
  height: 16px;
  border-radius: 50%;
  background-color: var(--accent);
  border: 3px solid var(--bg-white);
  box-shadow: 0 0 0 3px var(--accent);
  flex-shrink: 0;
  z-index: 1;
}

.timeline-line {
  width: 2px;
  flex: 1;
  background-color: var(--accent);
  margin-top: 4px;
}

.timeline-content {
  flex: 1;
}

.timeline-phase {
  display: inline-block;
  font-size: 12px;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 1px;
  color: var(--accent-dark);
  background-color: rgba(200, 150, 62, 0.1);
  padding: 2px 10px;
  border-radius: var(--radius-sm);
  margin-bottom: 8px;
}

.timeline-title {
  font-size: 18px;
  font-weight: 700;
  color: var(--text-primary);
  margin-bottom: 6px;
}

.timeline-desc {
  font-size: 14px;
  color: var(--text-secondary);
  line-height: 1.7;
  margin-bottom: 4px;
}

.timeline-period {
  font-size: 13px;
  color: var(--text-light);
  font-style: italic;
}
</style>
