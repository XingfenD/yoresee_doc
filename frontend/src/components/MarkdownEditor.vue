<template>
  <div class="markdown-editor">
    <div ref="editorRef" class="vditor-container"></div>
  </div>
</template>

<script setup>
import { ref, onMounted, onBeforeUnmount, watch } from 'vue';
import 'vditor/dist/index.css';
import { useVditorCore } from '@/composables/markdown-editor/useVditorCore';
import { useYjsCollaboration } from '@/composables/markdown-editor/useYjsCollaboration';
import { useCommentBridge } from '@/composables/markdown-editor/useCommentBridge';

const props = defineProps({
  modelValue: {
    type: String,
    default: ''
  },
  placeholder: {
    type: String,
    default: ''
  },
  height: {
    type: [String, Number],
    default: '100%'
  },
  collabEnabled: {
    type: Boolean,
    default: false
  },
  collabRoom: {
    type: String,
    default: ''
  },
  collabUrl: {
    type: String,
    default: '/collab'
  },
  collabToken: {
    type: String,
    default: ''
  },
  commentEnabled: {
    type: Boolean,
    default: false
  }
});

const emit = defineEmits([
  'update:modelValue',
  'collab-sync',
  'ready',
  'comment-add',
  'comment-remove',
  'comment-scroll',
  'comment-adjust',
  'comment-changed'
]);

const editorRef = ref(null);
const vditorRef = ref(null);
const isVditorReady = ref(false);
const suppressInput = ref(false);

let getEditableElement = () => editorRef.value;
let getValueSafely = (fallback = '') => fallback;
let setValueSafely = () => {};

const collaboration = useYjsCollaboration({
  props,
  emit,
  editorRef,
  vditorRef,
  isVditorReady,
  suppressInput,
  getEditableElement: (...args) => getEditableElement(...args),
  getValueSafely: (...args) => getValueSafely(...args),
  setValueSafely: (...args) => setValueSafely(...args)
});

const syncEditorToYjs = (fallbackValue = '') => {
  const nextValue = getValueSafely(fallbackValue);
  collaboration.syncContentToYjs(nextValue);
  emit('update:modelValue', nextValue);
};

const { getCommentOptions, removeCommentIds, broadcastCommentChange } = useCommentBridge({
  props,
  emit,
  vditorRef,
  ydocRef: collaboration.ydocRef,
  syncEditorToYjs
});

const vditorCore = useVditorCore({
  props,
  emit,
  editorRef,
  vditorRef,
  isVditorReady,
  suppressInput,
  getCommentOptions,
  onEditorInput: (value) => {
    queueMicrotask(() => {
      syncEditorToYjs(value);
    });
  }
});

getEditableElement = vditorCore.getEditableElement;
getValueSafely = vditorCore.getValueSafely;
setValueSafely = vditorCore.setValueSafely;

defineExpose({
  getVditor: () => vditorRef.value,
  removeCommentIds,
  broadcastCommentChange
});

onMounted(() => {
  vditorCore.initVditor();
  collaboration.setupCollaboration();
});

onBeforeUnmount(() => {
  vditorCore.destroyVditor();
  collaboration.teardownCollaboration();
});

watch(() => props.modelValue, (newValue) => {
  collaboration.handleModelValueChange(newValue);
});

watch(() => props.collabRoom, () => {
  collaboration.handleRoomChange();
});
</script>

<style scoped>
.markdown-editor {
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: column;
}

.vditor-container {
  flex: 1;
  min-height: 500px;
}

.markdown-editor :deep(.vditor) {
  border: none;
  background-color: var(--bg-white);
}

.dark-mode .markdown-editor :deep(.vditor) {
  background-color: var(--bg-white);
}

.markdown-editor :deep(.vditor-toolbar) {
  background-color: var(--bg-white);
  border-bottom: 1px solid var(--border-color);
}

.dark-mode .markdown-editor :deep(.vditor-toolbar) {
  background-color: var(--bg-medium);
  border-color: var(--border-color);
}

.markdown-editor :deep(.vditor-toolbar__item) {
  color: var(--text-medium);
}

.dark-mode .markdown-editor :deep(.vditor-toolbar__item) {
  color: var(--text-medium);
}

.markdown-editor :deep(.vditor-toolbar__item:hover) {
  color: var(--primary-color);
}

.markdown-editor :deep(.vditor-toolbar__item--current) {
  color: var(--primary-color);
}

.markdown-editor :deep(.vditor-content) {
  background-color: var(--bg-white);
  color: var(--text-dark);
}

.dark-mode .markdown-editor :deep(.vditor-content) {
  background-color: var(--bg-medium);
  color: var(--text-dark);
}

.markdown-editor :deep(.vditor-ir) {
  background-color: var(--bg-white);
  color: var(--text-dark);
}

.dark-mode .markdown-editor :deep(.vditor-ir) {
  background-color: var(--bg-medium);
  color: var(--text-dark);
}

.markdown-editor :deep(.vditor-ir__node) {
  color: var(--text-dark);
}

.dark-mode .markdown-editor :deep(.vditor-ir__node) {
  color: var(--text-dark);
}

.markdown-editor :deep(.vditor-ir__link) {
  color: var(--primary-color);
}

.markdown-editor :deep(.vditor-ir__link:hover) {
  color: var(--primary-color);
}

.markdown-editor :deep(.vditor-ir__marker) {
  color: var(--text-light);
}

.markdown-editor :deep(.vditor-ir__marker--heading) {
  color: var(--text-medium);
}

.dark-mode .markdown-editor :deep(.vditor-ir__marker) {
  color: var(--text-light);
}
</style>
