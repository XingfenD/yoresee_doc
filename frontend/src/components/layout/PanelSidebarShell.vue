<template>
  <div
    class="panel-sidebar-container"
    :class="{ collapsed, resizing, 'edge-left': resizeEdge === 'left', 'edge-right': resizeEdge === 'right' }"
    :style="containerStyle"
  >
    <div
      v-if="showResizer && resizeEdge === 'left'"
      class="panel-sidebar-resizer"
      role="separator"
      aria-orientation="vertical"
      @mousedown="$emit('resize-start', $event)"
    ></div>
    <aside class="panel-sidebar" :style="panelStyle">
      <slot name="header"></slot>
      <slot></slot>
    </aside>
    <div
      v-if="showResizer && resizeEdge === 'right'"
      class="panel-sidebar-resizer"
      role="separator"
      aria-orientation="vertical"
      @mousedown="$emit('resize-start', $event)"
    ></div>
  </div>
</template>

<script setup>
defineProps({
  collapsed: { type: Boolean, default: false },
  resizing: { type: Boolean, default: false },
  resizeEdge: { type: String, default: 'right' },
  showResizer: { type: Boolean, default: true },
  containerStyle: { type: [Object, Array], default: () => ({}) },
  panelStyle: { type: [Object, Array], default: () => ({}) }
});

defineEmits(['resize-start']);
</script>

<style scoped>
.panel-sidebar-container {
  display: flex;
  align-items: stretch;
  position: relative;
  overflow: hidden;
  flex-shrink: 0;
  transition: all 0.3s ease-in-out;
  background-color: var(--bg-white);
}

.panel-sidebar-container.edge-right {
  width: calc(var(--sidebar-width) + 6px);
}

.panel-sidebar-container.collapsed {
  width: 0;
  opacity: 0;
  visibility: hidden;
  pointer-events: none;
  border: none;
}

.panel-sidebar-container.resizing,
.panel-sidebar-container.resizing .panel-sidebar,
.panel-sidebar-container.resizing .panel-sidebar-resizer {
  transition: none !important;
}

.panel-sidebar {
  display: flex;
  flex-direction: column;
  overflow: hidden;
  flex-shrink: 0;
  background-color: var(--bg-white);
}

.panel-sidebar-container.edge-right .panel-sidebar {
  border-right: 1px solid var(--border-color);
}

.panel-sidebar-container.edge-left {
  border-left: 1px solid var(--border-color);
}

.panel-sidebar-resizer {
  width: 6px;
  flex-shrink: 0;
  cursor: col-resize;
  background-color: var(--bg-light);
  border-right: 1px solid var(--border-color);
  transition: background-color 0.2s ease;
}

.panel-sidebar-resizer:hover {
  background-color: var(--bg-medium);
}

:global(.dark-mode) .panel-sidebar-container,
:global(.dark-mode) .panel-sidebar {
  background-color: var(--bg-white);
  border-color: var(--border-color);
}

:global(.dark-mode) .panel-sidebar-resizer {
  background-color: var(--bg-medium);
  border-color: var(--border-color);
}

:global(.dark-mode) .panel-sidebar-resizer:hover {
  background-color: var(--bg-white);
}
</style>
