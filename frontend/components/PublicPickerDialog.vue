<template>
  <dialog
    :id="id"
    ref="dialog"
    class="public-picker-dialog"
    :aria-labelledby="`${id}-title`"
    :aria-describedby="subtitle ? `${id}-description` : undefined"
    :style="position"
    @click.self="closeOnBackdrop"
    @close="open = false"
    @keydown="navigateOptions"
  >
    <div class="picker-header">
      <div>
        <h2 :id="`${id}-title`">{{ title }}</h2>
        <p v-if="subtitle" :id="`${id}-description`">{{ subtitle }}</p>
      </div>
      <button type="button" class="picker-close" aria-label="关闭项目选择" @click="open = false">×</button>
    </div>
    <slot />
  </dialog>
</template>

<script setup lang="ts">
const props = defineProps<{
  id: string
  title: string
  subtitle?: string
  anchor: HTMLElement | null
}>()

const open = defineModel<boolean>('open', { required: true })
const dialog = ref<HTMLDialogElement | null>(null)
const position = ref<Record<string, string>>({})

function placeDialog() {
  const anchor = props.anchor?.getBoundingClientRect()
  if (!anchor) return

  const gutter = 12
  const width = Math.min(Math.max(anchor.width, 332), window.innerWidth - gutter * 2)
  const left = Math.max(gutter, Math.min(anchor.left, window.innerWidth - width - gutter))
  const spaceBelow = window.innerHeight - anchor.bottom - gutter * 2
  const spaceAbove = anchor.top - gutter * 2
  const showBelow = spaceBelow >= 360 || spaceBelow >= spaceAbove

  position.value = {
    '--picker-width': `${width}px`,
    '--picker-left': `${left}px`,
    '--picker-top': showBelow ? `${anchor.bottom + gutter}px` : 'auto',
    '--picker-bottom': showBelow ? 'auto' : `${window.innerHeight - anchor.top + gutter}px`,
    '--picker-max-height': `${Math.min(440, showBelow ? spaceBelow : spaceAbove)}px`,
  }
}

watch(open, (value) => {
  if (!dialog.value) return
  if (value && !dialog.value.open) {
    placeDialog()
    dialog.value.showModal()
  } else if (!value && dialog.value.open) {
    dialog.value.close()
  }
}, { flush: 'post' })

function closeOnBackdrop(event: MouseEvent) {
  const bounds = dialog.value?.getBoundingClientRect()
  if (bounds && (event.clientX < bounds.left || event.clientX > bounds.right
    || event.clientY < bounds.top || event.clientY > bounds.bottom)) open.value = false
}

function navigateOptions(event: KeyboardEvent) {
  if (event.isComposing) return
  if (event.key === 'Escape') {
    event.preventDefault()
    event.stopPropagation()
    open.value = false
    return
  }
  if (event.key !== 'ArrowDown' && event.key !== 'ArrowUp') return
  const target = event.target as HTMLElement | null
  if (!target?.matches('.public-picker-search, .public-picker-option')) return

  const options = Array.from(dialog.value?.querySelectorAll<HTMLElement>('.public-picker-option') ?? [])
  if (!options.length) return
  const index = options.indexOf(target)
  const next = event.key === 'ArrowDown' ? index + 1 : index < 0 ? options.length - 1 : index - 1
  event.preventDefault()
  if (next < 0) dialog.value?.querySelector<HTMLElement>('.public-picker-search')?.focus()
  else options[Math.min(next, options.length - 1)]?.focus()
}

onMounted(() => window.addEventListener('resize', placeDialogWhenOpen))
onBeforeUnmount(() => window.removeEventListener('resize', placeDialogWhenOpen))

function placeDialogWhenOpen() {
  if (open.value) placeDialog()
}
</script>

<style scoped>
.public-picker-dialog {
  position: fixed;
  inset: var(--picker-top) auto var(--picker-bottom) var(--picker-left);
  width: var(--picker-width);
  max-width: none;
  max-height: var(--picker-max-height);
  margin: 0;
  padding: 14px;
  border: 1px solid var(--public-control-border);
  border-radius: var(--radius-lg);
  background: var(--bg-white);
  color: var(--color-text);
  box-shadow: var(--shadow-lg);
  overscroll-behavior: contain;
}

.public-picker-dialog[open] {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.public-picker-dialog::backdrop {
  background: transparent;
}

.picker-header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 12px;
}

.picker-header h2 {
  margin: 0;
  font-size: 17px;
}

.picker-header p {
  margin: 2px 0 0;
  color: var(--color-text-secondary);
  font-size: 13px;
}

.picker-close {
  flex: 0 0 44px;
  min-height: 44px;
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  background: var(--bg-white);
  color: var(--color-text);
  font: 22px/1 var(--font-sans);
  cursor: pointer;
}

@media (max-width: 767px) {
  .public-picker-dialog {
    inset: auto 0 0;
    width: 100%;
    max-height: 78vh;
    max-height: 78dvh;
    padding: 12px 16px calc(14px + env(safe-area-inset-bottom, 0px));
    border-radius: var(--radius-xl) var(--radius-xl) 0 0;
  }

  .public-picker-dialog::before {
    content: '';
    flex: 0 0 4px;
    width: 28px;
    margin: 0 auto;
    border-radius: var(--radius-full);
    background: var(--public-control-border);
  }

  .public-picker-dialog::backdrop {
    background: rgba(15, 27, 45, 0.55);
  }
}
</style>
