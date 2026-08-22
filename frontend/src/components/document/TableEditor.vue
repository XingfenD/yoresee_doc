<template>
  <div class="table-editor">
    <div ref="editorRef" class="sheet-container"></div>
  </div>
</template>

<script setup>
import { ref, toRef } from 'vue';
import 'x-data-spreadsheet/dist/xspreadsheet.css';
import { useTableEditor } from '@/composables/document/editor/table-editor/useTableEditor';

const props = defineProps({
  modelValue: {
    type: String,
    default: ''
  }
});

const emit = defineEmits(['update:modelValue', 'commit']);
const editorRef = ref(null);

const { reRender } = useTableEditor(editorRef, {
  modelValue: toRef(props, 'modelValue'),
  onModelValue: (value) => emit('update:modelValue', value),
  onCommit: () => emit('commit')
});

defineExpose({
  reRender
});
</script>

<style scoped>
.table-editor {
  width: 100%;
  height: 100%;
  min-width: 0;
  min-height: 0;
  background: var(--bg-white);
  display: flex;
  flex-direction: column;
}

.sheet-container {
  flex: 1;
  min-width: 0;
  min-height: 0;
  overscroll-behavior: contain;
}

.table-editor :deep(.x-spreadsheet) {
  width: 100% !important;
  height: 100% !important;
  border: none;
  font-family: inherit;
  box-shadow: none;
  background: var(--bg-white);
}

.table-editor :deep(.x-spreadsheet-toolbar),
.table-editor :deep(.x-spreadsheet-bottombar) {
  width: 100%;
  box-sizing: border-box;
  background: var(--bg-white);
}
</style>

<style>
body.table-editor-open > .x-spreadsheet-dimmer {
  display: none !important;
  opacity: 0 !important;
  pointer-events: none !important;
  background: transparent !important;
}

body.dark-mode .table-editor {
  background: #0f172a;
}

body.dark-mode .table-editor .x-spreadsheet {
  background: #0f172a;
}

body.dark-mode .table-editor .x-spreadsheet-toolbar,
body.dark-mode .table-editor .x-spreadsheet-bottombar {
  background: #111827;
  border-color: #374151;
}

body.dark-mode .table-editor .x-spreadsheet-toolbar .x-spreadsheet-toolbar-btn:hover,
body.dark-mode .table-editor .x-spreadsheet-toolbar .x-spreadsheet-toolbar-btn.active {
  background: rgba(148, 163, 184, 0.18);
}

body.dark-mode .table-editor .x-spreadsheet-toolbar .x-spreadsheet-icon .x-spreadsheet-icon-img,
body.dark-mode .table-editor .x-spreadsheet-dropdown-header .x-spreadsheet-icon .x-spreadsheet-icon-img {
  filter: invert(0.9);
  opacity: 0.92;
}

body.dark-mode .table-editor .x-spreadsheet-toolbar-divider {
  border-right-color: #374151;
}

body.dark-mode .table-editor .x-spreadsheet-menu > li,
body.dark-mode .table-editor .x-spreadsheet-item {
  color: #cbd5e1;
}

body.dark-mode .table-editor .x-spreadsheet-menu > li.active,
body.dark-mode .table-editor .x-spreadsheet-item:hover,
body.dark-mode .table-editor .x-spreadsheet-item.active {
  background: rgba(148, 163, 184, 0.18);
  color: #f8fafc;
}

body.dark-mode .table-editor .x-spreadsheet-sheet,
body.dark-mode .table-editor .x-spreadsheet-table {
  background: #0b1220;
}

body.dark-mode .table-editor .x-spreadsheet-overlayer,
body.dark-mode .table-editor .x-spreadsheet-overlayer-content {
  background: transparent;
}

body.dark-mode .table-editor .x-spreadsheet-table {
  filter: invert(1) hue-rotate(180deg);
}

body.dark-mode .table-editor .x-spreadsheet-editor .x-spreadsheet-editor-area,
body.dark-mode .table-editor .x-spreadsheet-selector .x-spreadsheet-selector-area {
  background: rgba(37, 99, 235, 0.12);
}

body.dark-mode .table-editor .x-spreadsheet-editor .x-spreadsheet-editor-area textarea {
  background: #111827;
  color: #f8fafc;
}

body.dark-mode .table-editor .x-spreadsheet-scrollbar {
  background-color: #111827;
}

body.dark-mode .table-editor .x-spreadsheet-scrollbar.horizontal > div,
body.dark-mode .table-editor .x-spreadsheet-scrollbar.vertical > div {
  background: #4b5563;
}

body.dark-mode .table-editor .x-spreadsheet-dropdown .x-spreadsheet-dropdown-content,
body.dark-mode .table-editor .x-spreadsheet-contextmenu,
body.dark-mode .table-editor .x-spreadsheet-suggest,
body.dark-mode .table-editor .x-spreadsheet-sort-filter,
body.dark-mode .table-editor .x-spreadsheet-calendar,
body.dark-mode .table-editor .x-spreadsheet-modal,
body.dark-mode .table-editor .x-spreadsheet-toast {
  border-color: #374151;
  background: #111827;
  color: #e5e7eb;
}
</style>
