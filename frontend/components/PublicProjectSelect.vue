<template>
  <div class="public-project-select">
    <label :for="controlId" class="public-control-label">{{ label }}</label>
    <button
      :id="controlId"
      ref="trigger"
      class="public-select"
      type="button"
      :disabled="disabled || loading || !!error"
      :aria-describedby="error ? `${controlId}-value ${controlId}-error` : `${controlId}-value`"
      :aria-expanded="open"
      :aria-controls="panelId"
      aria-haspopup="dialog"
      @click="open = true"
    >
      <span :id="`${controlId}-value`" :class="{ 'public-select-placeholder': !selectedLabel }">
        {{ loading ? '正在加载项目…' : selectedLabel || placeholder }}
      </span>
    </button>
    <div v-if="error" :id="`${controlId}-error`" class="public-picker-error">
      <span role="status">{{ error }}</span>
      <button type="button" @click="$emit('retry')">重试</button>
    </div>

    <PublicPickerDialog :id="panelId" v-model:open="open" :anchor="trigger" :title="label">
      <template v-if="open">
        <div>
          <label :for="`${controlId}-search`" class="public-control-label">搜索项目</label>
          <input
            :id="`${controlId}-search`"
            v-model="keyword"
            class="public-picker-search"
            type="search"
            placeholder="例如 EB-5 或香港"
            autocomplete="off"
            autofocus
            @keydown.enter="!$event.isComposing && $event.preventDefault()"
          >
        </div>
        <div v-if="filteredOptions.length" class="public-picker-options">
          <button
            v-for="option in filteredOptions"
            :key="option.value"
            class="public-picker-option"
            type="button"
            :aria-pressed="model === option.value"
            @click="select(option.value)"
          >
            <span>{{ option.label }}</span>
            <span v-if="model === option.value" aria-hidden="true">✓</span>
          </button>
        </div>
        <p v-else class="public-picker-empty" role="status">
          {{ options.length ? '未找到相关项目' : '暂无可选项目' }}
        </p>
      </template>
    </PublicPickerDialog>
  </div>
</template>

<script setup lang="ts">
const props = withDefaults(defineProps<{
  id?: string
  label: string
  options: { value: string; label: string }[]
  placeholder?: string
  disabled?: boolean
  loading?: boolean
  error?: string | null
}>(), {
  placeholder: '请选择项目',
})

defineEmits<{ retry: [] }>()

const model = defineModel<string>({ required: true })
const controlId = props.id || useId()
const panelId = `${controlId}-options`
const trigger = ref<HTMLButtonElement | null>(null)
const open = ref(false)
const keyword = ref('')

const selectedLabel = computed(() => model.value
  ? props.options.find(option => option.value === model.value)?.label
  : undefined,
)
const filteredOptions = computed(() => {
  const search = keyword.value.trim().toLowerCase()
  if (!search) return props.options
  return props.options.filter(option =>
    option.label.toLowerCase().includes(search) || option.value.toLowerCase().includes(search),
  )
})

watch(open, value => {
  if (!value) keyword.value = ''
})

function select(value: string) {
  model.value = value
  open.value = false
}
</script>
