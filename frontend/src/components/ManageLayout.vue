<template>
  <div class="manage-layout">
    <slot />
  </div>
</template>

<style scoped>
.manage-layout {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-lg);
}
</style>
