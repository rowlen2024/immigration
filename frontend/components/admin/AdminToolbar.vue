<template>
  <div class="admin-toolbar">
    <div class="admin-toolbar-main">
      <slot />
    </div>
    <div v-if="$slots.actions" class="admin-toolbar-actions">
      <slot name="actions" />
    </div>
  </div>
</template>
