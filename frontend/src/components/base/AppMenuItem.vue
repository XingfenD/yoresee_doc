<template>
  <button
    type="button"
    class="app-menu-item"
    :class="{
      'is-danger': danger,
      'is-disabled': disabled
    }"
    :disabled="disabled"
    @click="handleClick"
  >
    <slot name="icon" />
    <slot />
  </button>
</template>

<script setup>
const props = defineProps({
  danger: {
    type: Boolean,
    default: false
  },
  disabled: {
    type: Boolean,
    default: false
  }
});

const emit = defineEmits(['click']);

const handleClick = () => {
  if (props.disabled) {
    return;
  }
  emit('click');
};
</script>

<style scoped>
.app-menu-item {
  width: 100%;
  padding: var(--spacing-xs) var(--spacing-md);
  background: transparent;
  border: none;
  text-align: left;
  display: flex;
  align-items: center;
  gap: var(--spacing-xs);
  color: var(--text-medium);
  cursor: pointer;
}

.app-menu-item:hover {
  background-color: var(--bg-light);
  color: var(--primary-color);
}

.dark-mode .app-menu-item:hover {
  background-color: rgba(255, 255, 255, 0.08);
}

.app-menu-item.is-danger:hover {
  color: #f56c6c;
}

.app-menu-item.is-disabled {
  color: var(--text-light);
  cursor: not-allowed;
}

.app-menu-item.is-disabled:hover {
  color: var(--text-light);
  background: transparent;
}

.app-menu-item :deep(.el-icon) {
  color: var(--text-light);
}

.app-menu-item:hover :deep(.el-icon) {
  color: var(--primary-color);
}

.app-menu-item.is-danger:hover :deep(.el-icon) {
  color: #f56c6c;
}
</style>
