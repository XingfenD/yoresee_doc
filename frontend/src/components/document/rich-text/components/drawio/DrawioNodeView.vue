<template>
  <NodeViewWrapper class="drawio-node" contenteditable="false">
    <header class="drawio-header" data-drag-handle>
      <span class="drawio-badge">Draw.io</span>
      <div class="drawio-actions">
        <button type="button" class="drawio-btn" @click="openEditor">
          {{ hasDiagram ? 'Edit Diagram' : 'Create Diagram' }}
        </button>
        <button
          v-if="hasDiagram"
          type="button"
          class="drawio-btn"
          @click="clearDiagram"
        >
          Clear
        </button>
      </div>
    </header>

    <div class="drawio-body">
      <div v-if="hasDiagram" class="drawio-preview-wrap">
        <iframe
          ref="previewFrameRef"
          class="drawio-preview-frame"
          :src="DRAWIO_PREVIEW_URL"
          tabindex="-1"
        ></iframe>
      </div>
    </div>

    <Teleport to="body">
      <div v-if="editorVisible" class="drawio-overlay" @click.self="closeEditor">
        <div class="drawio-modal">
          <header class="drawio-modal-header">
            <span>Draw.io Editor</span>
            <button type="button" class="drawio-btn" @click="closeEditor">Close</button>
          </header>
          <iframe ref="iframeRef" class="drawio-frame" :src="DRAWIO_EMBED_URL"></iframe>
        </div>
      </div>
    </Teleport>
  </NodeViewWrapper>
</template>

<script setup>
import { computed, onBeforeUnmount, ref, watch } from 'vue';
import { NodeViewWrapper } from '@tiptap/vue-3';
import { DEFAULT_DRAWIO_XML } from './drawioExtension';

const DRAWIO_ORIGIN = 'https://embed.diagrams.net';
const DRAWIO_EMBED_URL = `${DRAWIO_ORIGIN}/?embed=1&ui=min&spin=1&proto=json&saveAndExit=1&noSaveBtn=1`;
const DRAWIO_PREVIEW_URL = `${DRAWIO_ORIGIN}/?embed=1&ui=min&spin=1&proto=json&chrome=0&toolbar=0&pages=0&layers=0&status=0&noSaveBtn=1`;

const props = defineProps({
  node: {
    type: Object,
    required: true
  },
  updateAttributes: {
    type: Function,
    required: true
  }
});

const iframeRef = ref(null);
const previewFrameRef = ref(null);
const editorVisible = ref(false);
const previewReady = ref(false);

const currentDiagram = computed(() => String(props.node?.attrs?.diagram || ''));
const hasDiagram = computed(() => currentDiagram.value.trim().length > 0);

const postToFrame = (frameRef, message) => {
  const targetWindow = frameRef.value?.contentWindow;
  if (!targetWindow) {
    return;
  }
  targetWindow.postMessage(JSON.stringify(message), DRAWIO_ORIGIN);
};

const postToEditorFrame = (message) => {
  postToFrame(iframeRef, message);
};

const loadPreviewDiagram = () => {
  if (!previewReady.value || !hasDiagram.value) {
    return;
  }
  postToFrame(previewFrameRef, {
    action: 'load',
    autosave: 0,
    xml: currentDiagram.value
  });
};

const closeEditor = () => {
  editorVisible.value = false;
};

const openEditor = () => {
  editorVisible.value = true;
  postToEditorFrame({
    action: 'load',
    autosave: 1,
    xml: currentDiagram.value || DEFAULT_DRAWIO_XML
  });
};

const clearDiagram = () => {
  props.updateAttributes({ diagram: '' });
};

const handleMessage = (event) => {
  if (event.origin !== DRAWIO_ORIGIN) {
    return;
  }
  const isEditorFrame = event.source === iframeRef.value?.contentWindow;
  const isPreviewFrame = event.source === previewFrameRef.value?.contentWindow;
  if (!isEditorFrame && !isPreviewFrame) {
    return;
  }
  let payload = event.data;
  if (typeof payload === 'string') {
    try {
      payload = JSON.parse(payload);
    } catch (_) {
      return;
    }
  }
  if (!payload || typeof payload !== 'object') {
    return;
  }

  const eventName = String(payload.event || '').trim();
  if (eventName === 'init') {
    if (isPreviewFrame) {
      previewReady.value = true;
      loadPreviewDiagram();
      return;
    }
    if (!editorVisible.value) {
      return;
    }
    postToEditorFrame({
      action: 'load',
      autosave: 1,
      xml: currentDiagram.value || DEFAULT_DRAWIO_XML
    });
    return;
  }

  if (!isEditorFrame || !editorVisible.value) {
    return;
  }

  if (eventName === 'autosave' && typeof payload.xml === 'string') {
    props.updateAttributes({ diagram: payload.xml });
    return;
  }

  if (eventName === 'save') {
    if (typeof payload.xml === 'string') {
      props.updateAttributes({ diagram: payload.xml });
    }
    postToEditorFrame({ action: 'exit' });
    return;
  }

  if (eventName === 'exit') {
    closeEditor();
  }
};

if (typeof window !== 'undefined') {
  window.addEventListener('message', handleMessage);
}

watch(
  currentDiagram,
  () => {
    loadPreviewDiagram();
  }
);

watch(
  hasDiagram,
  (next) => {
    if (!next) {
      previewReady.value = false;
    }
  }
);

onBeforeUnmount(() => {
  if (typeof window !== 'undefined') {
    window.removeEventListener('message', handleMessage);
  }
});
</script>

<style scoped>
.drawio-node {
  border: 1px solid var(--border-color);
  border-radius: 10px;
  overflow: hidden;
  background: var(--bg-white);
}

.drawio-header {
  height: 36px;
  padding: 0 10px;
  border-bottom: 1px solid var(--border-color);
  display: flex;
  align-items: center;
  justify-content: space-between;
  background: color-mix(in srgb, var(--bg-light) 70%, transparent);
}

.drawio-badge {
  font-size: 12px;
  font-weight: 600;
  color: var(--text-medium);
}

.drawio-actions {
  display: flex;
  align-items: center;
  gap: 6px;
}

.drawio-btn {
  height: 24px;
  border-radius: 6px;
  border: 1px solid var(--border-color);
  background: var(--bg-white);
  color: var(--text-medium);
  font-size: 12px;
  padding: 0 8px;
  cursor: pointer;
}

.drawio-body {
  min-height: 120px;
  display: flex;
  flex-direction: column;
  align-items: stretch;
  padding: 12px;
}

.drawio-preview-wrap {
  border: 1px solid var(--border-color);
  border-radius: 8px;
  overflow: hidden;
  width: 100%;
  min-height: 220px;
  height: 280px;
  background: #ffffff;
}

.drawio-preview-frame {
  width: 100%;
  height: 100%;
  border: none;
  pointer-events: none;
}

.drawio-overlay {
  position: fixed;
  inset: 0;
  background: color-mix(in srgb, #000 45%, transparent);
  z-index: 4000;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 20px;
  box-sizing: border-box;
}

.drawio-modal {
  width: min(1200px, calc(100vw - 40px));
  height: min(760px, calc(100vh - 40px));
  border-radius: 12px;
  border: 1px solid var(--border-color);
  background: var(--bg-white);
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

.drawio-modal-header {
  height: 46px;
  border-bottom: 1px solid var(--border-color);
  padding: 0 14px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  color: var(--text-primary);
  font-size: 14px;
  font-weight: 600;
}

.drawio-frame {
  width: 100%;
  flex: 1;
  border: none;
  background: #ffffff;
}
</style>
