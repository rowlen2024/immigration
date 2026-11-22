<template>
  <div class="public-form-field" :class="{ 'public-form-field--reserve-feedback': reserveFeedback }">
    <label :for="id" class="public-control-label">
      {{ label }}<span v-if="required" class="public-field-required" aria-hidden="true"> *</span>
    </label>
    <slot :control="control" />
    <div v-if="feedback || reserveFeedback" class="public-field-feedback">
      <span v-if="feedback" :id="feedbackId" :class="{ 'public-field-error': error }">{{ feedback }}</span>
    </div>
  </div>
</template>

<script setup lang="ts">
const props = defineProps<{
  id: string
  label: string
  required?: boolean
  hint?: string
  error?: string
  reserveFeedback?: boolean
}>()

const feedback = computed(() => props.error || props.hint)
const feedbackId = computed(() => feedback.value
  ? `${props.id}-${props.error ? 'error' : 'hint'}`
  : undefined)
const control = computed(() => ({
  id: props.id,
  class: 'public-control',
  required: props.required,
  'aria-invalid': !!props.error,
  'aria-describedby': feedbackId.value,
}))
</script>

<style scoped>
.public-field-required,
.public-field-error {
  color: var(--color-danger);
}

.public-field-feedback {
  margin-top: 6px;
  color: var(--color-text-secondary);
  font-size: 13px;
}

.public-form-field--reserve-feedback .public-field-feedback {
  min-height: 1.6em;
}
</style>
