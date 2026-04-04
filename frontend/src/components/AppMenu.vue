<template>
  <Teleport to="body">
    <div
      v-if="visible"
      ref="menuRef"
      class="app-menu"
      :class="menuClass"
      :style="menuStyle"
      @click.stop
    >
      <slot />
    </div>
  </Teleport>
</template>

<script setup>
import { computed, onBeforeUnmount, onMounted, ref } from 'vue';

const props = defineProps({
  visible: {
    type: Boolean,
    default: false
  },
  x: {
    type: Number,
    default: 0
  },
  y: {
    type: Number,
    default: 0
  },
  minWidth: {
    type: [Number, String],
    default: 150
  },
  zIndex: {
    type: Number,
    default: 3000
  },
  menuClass: {
    type: String,
    default: ''
  },
  closeOnOutside: {
    type: Boolean,
    default: true
  },
  closeOnScroll: {
    type: Boolean,
    default: true
  },
  closeOnResize: {
    type: Boolean,
    default: true
  },
  closeOnEscape: {
    type: Boolean,
    default: true
  },
  closeOnContextOutside: {
    type: Boolean,
    default: true
  },
  ignoreElements: {
    type: Array,
    default: () => []
  }
});

const emit = defineEmits(['close']);

const menuRef = ref(null);
const getMenuElement = () => menuRef.value;

const resolveElement = (entry) => {
  if (!entry) {
    return null;
  }
  if (entry instanceof HTMLElement) {
    return entry;
  }
  if (entry?.value instanceof HTMLElement) {
    return entry.value;
  }
  if (entry?.$el instanceof HTMLElement) {
    return entry.$el;
  }
  return null;
};

const isInsideIgnoredTargets = (target) => {
  if (!target || !props.ignoreElements.length) {
    return false;
  }
  for (const entry of props.ignoreElements) {
    const el = resolveElement(entry);
    if (el?.contains(target)) {
      return true;
    }
  }
  return false;
};

const requestClose = () => {
  if (!props.visible) {
    return;
  }
  emit('close');
};

const handlePointerDown = (event) => {
  if (!props.visible || !props.closeOnOutside) {
    return;
  }
  const target = event.target;
  if (menuRef.value?.contains(target) || isInsideIgnoredTargets(target)) {
    return;
  }
  requestClose();
};

const handleContextMenu = (event) => {
  if (!props.visible || !props.closeOnContextOutside) {
    return;
  }
  const target = event.target;
  if (menuRef.value?.contains(target) || isInsideIgnoredTargets(target)) {
    return;
  }
  requestClose();
};

const handleWindowScroll = () => {
  if (props.closeOnScroll) {
    requestClose();
  }
};

const handleWindowResize = () => {
  if (props.closeOnResize) {
    requestClose();
  }
};

const handleWindowKeydown = (event) => {
  if (props.closeOnEscape && event.key === 'Escape') {
    requestClose();
  }
};

const menuStyle = computed(() => {
  const minWidth = typeof props.minWidth === 'number' ? `${props.minWidth}px` : props.minWidth;
  return {
    left: `${props.x}px`,
    top: `${props.y}px`,
    minWidth,
    zIndex: props.zIndex
  };
});

onMounted(() => {
  window.addEventListener('mousedown', handlePointerDown, true);
  window.addEventListener('contextmenu', handleContextMenu, true);
  window.addEventListener('scroll', handleWindowScroll, true);
  window.addEventListener('resize', handleWindowResize);
  window.addEventListener('keydown', handleWindowKeydown);
});

onBeforeUnmount(() => {
  window.removeEventListener('mousedown', handlePointerDown, true);
  window.removeEventListener('contextmenu', handleContextMenu, true);
  window.removeEventListener('scroll', handleWindowScroll, true);
  window.removeEventListener('resize', handleWindowResize);
  window.removeEventListener('keydown', handleWindowKeydown);
});

defineExpose({
  getMenuElement
});
</script>

<style scoped>
.app-menu {
  position: fixed;
  background-color: var(--bg-white);
  border: 1px solid var(--border-color);
  border-radius: var(--border-radius-sm);
  box-shadow: var(--shadow-md);
  padding: var(--spacing-xs) 0;
}

.dark-mode .app-menu {
  background-color: var(--bg-white);
  border-color: var(--border-color);
}
</style>
