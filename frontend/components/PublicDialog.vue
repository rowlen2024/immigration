<template>
  <dialog
    :id="id"
    ref="dialog"
    class="public-dialog"
    :class="{ 'public-dialog--anchored': anchor !== undefined }"
    :inert="!open"
    :aria-labelledby="`${id}-title`"
    :aria-describedby="subtitle ? `${id}-description` : undefined"
    :style="position"
    @cancel.prevent="open = false"
    @click.self="closeOnBackdrop"
    @close="syncClosed"
    @keydown="navigateOptions"
  >
    <div class="picker-header">
      <div>
        <h2 :id="`${id}-title`">{{ title }}</h2>
        <p v-if="subtitle" :id="`${id}-description`">{{ subtitle }}</p>
      </div>
      <button type="button" class="public-button dialog-close" :aria-label="closeLabel" @click="open = false">×</button>
    </div>
    <slot v-if="open || wasOpened" />
  </dialog>
</template>

<script setup lang="ts">
const props = withDefaults(defineProps<{
  id: string
  title: string
  subtitle?: string
  anchor?: HTMLElement | null
  closeLabel?: string
}>(), {
  closeLabel: '关闭',
})

const open = defineModel<boolean>('open', { required: true })
const dialog = ref<HTMLDialogElement | null>(null)
const wasOpened = ref(false)
const position = ref<Record<string, string>>({})

function placeDialog() {
  const anchor = props.anchor?.getBoundingClientRect()
  if (!anchor) return

  const gutter = 12
  const viewportWidth = document.body.clientWidth
  const width = Math.min(Math.max(anchor.width, 332), viewportWidth - gutter * 2)
  const left = Math.max(gutter, Math.min(anchor.left, viewportWidth - width - gutter))
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
    wasOpened.value = true
    props.anchor?.focus({ preventScroll: true })
    dialog.value.showModal()
  } else if (!value && dialog.value.open) {
    dialog.value.close()
  }
}, { flush: 'post' })

function syncClosed() {
  if (!dialog.value?.open) open.value = false
}

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
.public-dialog {
  --dialog-enter-duration: var(--public-motion-slow);
  --dialog-leave-duration: var(--public-motion-medium);
  --dialog-offset-y: 12px;
  position: fixed;
  inset: 0;
  width: min(420px, calc(100% - 32px));
  max-height: 85vh;
  max-height: 85dvh;
  margin: auto;
  padding: 14px;
  border: 1px solid var(--public-control-border);
  border-radius: var(--radius-lg);
  background: var(--bg-white);
  color: var(--color-text);
  box-shadow: var(--shadow-lg);
  overscroll-behavior: contain;
  flex-direction: column;
  gap: 10px;
  visibility: visible;
}

.public-dialog--anchored {
  --dialog-enter-duration: var(--public-motion-medium);
  --dialog-leave-duration: var(--public-motion-fast);
  --dialog-offset-y: -4px;
  inset: var(--picker-top) auto var(--picker-bottom) var(--picker-left);
  width: var(--picker-width);
  max-width: none;
  max-height: var(--picker-max-height);
  margin: 0;
}

.public-dialog[open] {
  display: flex;
}

.public-dialog::backdrop {
  background: rgba(15, 27, 45, 0.55);
}

.public-dialog--anchored::backdrop {
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

.dialog-close {
  flex: 0 0 44px;
  padding: 0;
  font-size: 22px;
}

@media (max-width: 767px) {
  .public-dialog--anchored {
    --dialog-enter-duration: var(--public-motion-slow);
    --dialog-leave-duration: var(--public-motion-medium);
    --dialog-offset-y: 12px;
    inset: auto 0 0;
    width: 100%;
    max-height: 78vh;
    max-height: 78dvh;
    padding: 12px 16px calc(14px + env(safe-area-inset-bottom, 0px));
    border-radius: var(--radius-xl) var(--radius-xl) 0 0;
  }

  .public-dialog--anchored::before {
    content: '';
    flex: 0 0 4px;
    width: 28px;
    margin: 0 auto;
    border-radius: var(--radius-full);
    background: var(--public-control-border);
  }

  .public-dialog--anchored::backdrop {
    background: rgba(15, 27, 45, 0.55);
  }
}

@supports (overlay: auto) and (transition-behavior: allow-discrete) {
  .public-dialog {
    opacity: 0;
    transform: translateY(var(--dialog-offset-y));
    transition: opacity var(--dialog-leave-duration) var(--ease-in),
                transform var(--dialog-leave-duration) var(--ease-in),
                display var(--dialog-leave-duration) allow-discrete,
                overlay var(--dialog-leave-duration) allow-discrete;
  }

  .public-dialog[open] {
    opacity: 1;
    transform: translateY(0);
    transition-timing-function: var(--ease-out);
    transition-duration: var(--dialog-enter-duration);
  }

  .public-dialog:not([open]),
  .public-dialog:not([open])::backdrop {
    pointer-events: none;
  }

  .public-dialog::backdrop {
    opacity: 0;
    transition: opacity var(--dialog-leave-duration) var(--ease-in),
                display var(--dialog-leave-duration) allow-discrete,
                overlay var(--dialog-leave-duration) allow-discrete;
  }

  .public-dialog[open]::backdrop {
    opacity: 1;
    transition-timing-function: var(--ease-out);
    transition-duration: var(--dialog-enter-duration);
  }

  @starting-style {
    .public-dialog[open] {
      opacity: 0;
      transform: translateY(var(--dialog-offset-y));
    }

    .public-dialog[open]::backdrop {
      opacity: 0;
    }
  }
}

@media (prefers-reduced-motion: reduce) {
  .public-dialog,
  .public-dialog::backdrop {
    transition: none !important;
  }
}
</style>
