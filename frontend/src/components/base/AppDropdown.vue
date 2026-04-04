<template>
  <el-dropdown
    :trigger="trigger"
    :placement="resolvedPlacement"
    :disabled="disabled"
    :hide-on-click="hideOnClick"
    :teleported="teleported"
    :popper-class="resolvedPopperClass"
    @command="(command) => emit('command', command)"
    @visible-change="handleVisibleChange"
  >
    <slot />
    <template #dropdown>
      <slot name="dropdown" />
    </template>
  </el-dropdown>
</template>

<script setup>
import { computed } from 'vue';

const props = defineProps({
  trigger: {
    type: String,
    default: 'click'
  },
  placement: {
    type: String,
    default: undefined
  },
  disabled: {
    type: Boolean,
    default: false
  },
  hideOnClick: {
    type: Boolean,
    default: true
  },
  teleported: {
    type: Boolean,
    default: true
  },
  popperClass: {
    type: String,
    default: ''
  }
});

const emit = defineEmits(['command', 'visible-change']);

const validPlacements = new Set([
  'top',
  'top-start',
  'top-end',
  'bottom',
  'bottom-start',
  'bottom-end',
  'left',
  'left-start',
  'left-end',
  'right',
  'right-start',
  'right-end'
]);

const resolvedPlacement = computed(() => {
  const placement = (props.placement || '').trim();
  if (!placement) {
    return undefined;
  }
  return validPlacements.has(placement) ? placement : undefined;
});

const resolvedPopperClass = computed(() =>
  ['app-dropdown-popper', props.popperClass].filter(Boolean).join(' ')
);

const handleVisibleChange = (visible) => {
  emit('visible-change', visible);
  if (visible) {
    return;
  }
  requestAnimationFrame(() => {
    const active = document.activeElement;
    if (!active || !(active instanceof HTMLElement)) {
      return;
    }
    if (!active.classList.contains('el-dropdown-menu__item')) {
      return;
    }
    active.blur();
  });
};
</script>

<style scoped>
:global(.dark-mode .app-dropdown-popper.el-popper) {
  background-color: var(--bg-white) !important;
  border-color: var(--border-color) !important;
}

:global(.dark-mode .app-dropdown-popper .el-popper__arrow::before) {
  background-color: var(--bg-white) !important;
  border-color: var(--border-color) !important;
}

:global(.dark-mode .app-dropdown-popper .el-dropdown-menu) {
  background-color: var(--bg-white) !important;
}

:global(.dark-mode .app-dropdown-popper .el-dropdown-menu__item) {
  color: var(--text-dark) !important;
}

:global(.dark-mode .app-dropdown-popper .el-dropdown-menu__item:not(.is-disabled):hover),
:global(.dark-mode .app-dropdown-popper .el-dropdown-menu__item:not(.is-disabled):focus),
:global(.dark-mode .app-dropdown-popper .el-dropdown-menu__item.is-hovering),
:global(.dark-mode .app-dropdown-popper .el-dropdown-menu__item.hover) {
  background-color: var(--select-option-hover) !important;
  color: var(--text-dark) !important;
}

:global(.dark-mode .app-dropdown-popper .el-dropdown-menu__item.hover:not(:hover):not(:focus):not(.is-hovering)) {
  background-color: transparent !important;
}
</style>
