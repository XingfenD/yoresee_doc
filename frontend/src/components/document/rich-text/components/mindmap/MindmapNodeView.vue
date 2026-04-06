<template>
  <NodeViewWrapper class="mindmap-node" contenteditable="false">
    <header class="mindmap-node-header" data-drag-handle>
      <span class="mindmap-badge">Mindmap</span>
      <div class="mindmap-actions">
        <button type="button" class="mindmap-btn" @click="toggleEditing">
          {{ editing ? 'Done' : 'Edit Source' }}
        </button>
        <button type="button" class="mindmap-btn mindmap-btn--danger" @click="deleteNode">
          Remove
        </button>
      </div>
    </header>

    <div class="mindmap-node-body" :class="{ 'is-editing': editing }">
      <textarea
        v-if="editing"
        v-model="draftSource"
        class="mindmap-source-input"
        @input="applySource"
      />
      <div ref="canvasWrapRef" class="mindmap-canvas-wrap">
        <svg ref="svgRef" class="mindmap-canvas"></svg>
      </div>
    </div>
  </NodeViewWrapper>
</template>

<script setup>
import { nextTick, onBeforeUnmount, onMounted, ref, watch } from 'vue';
import { NodeViewWrapper } from '@tiptap/vue-3';
import { Markmap } from 'markmap-view';
import { Transformer } from 'markmap-lib';
import { DEFAULT_MINDMAP_SOURCE } from './mindmapExtension';

const props = defineProps({
  node: {
    type: Object,
    required: true
  },
  updateAttributes: {
    type: Function,
    required: true
  },
  deleteNode: {
    type: Function,
    required: true
  }
});

const svgRef = ref(null);
const canvasWrapRef = ref(null);
const editing = ref(false);
const isDarkMode = ref(Boolean(document?.body?.classList?.contains('dark-mode')));
const draftSource = ref(props.node?.attrs?.source || DEFAULT_MINDMAP_SOURCE);
const markmapInstanceRef = ref(null);
const resizeObserverRef = ref(null);
const transformer = new Transformer();
const themeObserverRef = ref(null);

const resolveMindmapTextColor = () =>
  isDarkMode.value ? '#e5e7eb' : '#1f2937';

const fitMindmapToViewport = async () => {
  const instance = markmapInstanceRef.value;
  if (!instance || editing.value) {
    return;
  }
  await nextTick();
  await new Promise((resolve) => requestAnimationFrame(resolve));
  const fitPromise = instance.fit?.(Number.POSITIVE_INFINITY);
  if (fitPromise && typeof fitPromise.catch === 'function') {
    fitPromise.catch(() => {});
  }
};

const applyMindmapThemeVariables = () => {
  const svg = svgRef.value;
  if (!svg) {
    return;
  }
  const textColor = resolveMindmapTextColor();
  svg.style.setProperty('--markmap-text-color', textColor);
  svg.style.setProperty('--markmap-a-color', isDarkMode.value ? '#9ec5ff' : '#2563eb');
  svg.style.setProperty('--markmap-a-hover-color', isDarkMode.value ? '#c7dcff' : '#1d4ed8');
  svg.style.setProperty('--markmap-code-bg', isDarkMode.value ? '#111827' : '#f3f4f6');
  svg.style.setProperty('--markmap-code-color', isDarkMode.value ? '#e5e7eb' : '#374151');
  svg.style.setProperty('--markmap-circle-open-bg', isDarkMode.value ? '#1f2937' : '#ffffff');
};

const syncSvgViewport = () => {
  const svg = svgRef.value;
  const wrap = canvasWrapRef.value;
  if (!svg || !wrap) {
    return;
  }

  const rect = wrap.getBoundingClientRect();
  const width = Math.max(320, Math.round(rect.width || wrap.clientWidth || 0));
  const height = Math.max(220, Math.round(rect.height || wrap.clientHeight || 0));
  svg.setAttribute('width', String(width));
  svg.setAttribute('height', String(height));
  svg.setAttribute('viewBox', `0 0 ${width} ${height}`);
};

const renderMindmap = async () => {
  await nextTick();
  if (!svgRef.value) {
    return;
  }
  syncSvgViewport();
  applyMindmapThemeVariables();
  const source = String(props.node?.attrs?.source || DEFAULT_MINDMAP_SOURCE).trim() || DEFAULT_MINDMAP_SOURCE;
  const transformed = transformer.transform(source);
  const root = transformed?.root;
  if (!root) {
    return;
  }

  if (!markmapInstanceRef.value) {
    markmapInstanceRef.value = Markmap.create(
      svgRef.value,
      {
        autoFit: false,
        duration: 120,
        fitRatio: 0.98,
        color: () => resolveMindmapTextColor()
      },
      root
    );
    await fitMindmapToViewport();
    return;
  }

  await markmapInstanceRef.value.setData(root);
  await fitMindmapToViewport();
};

const applySource = () => {
  props.updateAttributes({
    source: draftSource.value || DEFAULT_MINDMAP_SOURCE
  });
};

const toggleEditing = () => {
  editing.value = !editing.value;
};

watch(
  () => props.node?.attrs?.source,
  (next) => {
    draftSource.value = next || DEFAULT_MINDMAP_SOURCE;
    renderMindmap();
  },
  { immediate: true }
);

watch(editing, (next) => {
  if (!next) {
    renderMindmap();
  }
});

onMounted(() => {
  if (typeof ResizeObserver !== 'undefined' && canvasWrapRef.value) {
    const observer = new ResizeObserver(() => {
      syncSvgViewport();
      void fitMindmapToViewport();
    });
    observer.observe(canvasWrapRef.value);
    resizeObserverRef.value = observer;
  }
  if (typeof MutationObserver !== 'undefined' && document?.body) {
    const observer = new MutationObserver(() => {
      isDarkMode.value = Boolean(document?.body?.classList?.contains('dark-mode'));
      renderMindmap();
    });
    observer.observe(document.body, { attributes: true, attributeFilter: ['class'] });
    themeObserverRef.value = observer;
  }
  renderMindmap();
});

onBeforeUnmount(() => {
  if (resizeObserverRef.value) {
    resizeObserverRef.value.disconnect();
    resizeObserverRef.value = null;
  }
  if (themeObserverRef.value) {
    themeObserverRef.value.disconnect();
    themeObserverRef.value = null;
  }
  markmapInstanceRef.value = null;
});
</script>

<style scoped>
.mindmap-node {
  border: 1px solid var(--border-color);
  border-radius: 10px;
  background: var(--bg-white);
  overflow: hidden;
}

.mindmap-node-header {
  height: 34px;
  padding: 0 10px;
  border-bottom: 1px solid var(--border-color);
  display: flex;
  align-items: center;
  justify-content: space-between;
  background: color-mix(in srgb, var(--bg-light) 70%, transparent);
}

.mindmap-badge {
  font-size: 12px;
  font-weight: 600;
  color: var(--text-medium);
}

.mindmap-actions {
  display: flex;
  align-items: center;
  gap: 6px;
}

.mindmap-btn {
  height: 24px;
  border: 1px solid var(--border-color);
  border-radius: 6px;
  background: var(--bg-white);
  color: var(--text-medium);
  padding: 0 8px;
  font-size: 12px;
  cursor: pointer;
}

.mindmap-btn--danger {
  color: #ef4444;
}

.mindmap-node-body {
  display: grid;
  grid-template-columns: minmax(180px, 0.45fr) minmax(0, 1fr);
  height: 320px;
  min-height: 280px;
}

.mindmap-node-body:not(.is-editing) {
  grid-template-columns: 1fr;
}

.mindmap-source-input {
  display: block;
  width: 100%;
  height: 100%;
  min-height: 0;
  border: none;
  border-right: 1px solid var(--border-color);
  outline: none;
  padding: 10px 12px;
  font-size: 12px;
  line-height: 1.55;
  resize: none;
  box-sizing: border-box;
  background: transparent;
  color: var(--text-primary);
}

.mindmap-canvas-wrap {
  min-width: 0;
  min-height: 280px;
  height: 100%;
  overflow: auto;
}

.mindmap-canvas {
  display: block;
  width: 100%;
  height: 100%;
  min-height: 280px;
}

:global(.dark-mode) .mindmap-node .mindmap-canvas :deep(.markmap-link) {
  stroke: #475569 !important;
}

@media (max-width: 640px) {
  .mindmap-node-body.is-editing {
    grid-template-columns: 1fr;
  }

  .mindmap-source-input {
    border-right: none;
    border-bottom: 1px solid var(--border-color);
    min-height: 120px;
  }
}
</style>
