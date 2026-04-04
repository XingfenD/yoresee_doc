<template>
  <div class="table-editor">
    <div ref="editorRef" class="sheet-container"></div>
  </div>
</template>

<script setup>
import { nextTick, onBeforeUnmount, onMounted, ref, watch } from 'vue';
import Spreadsheet from 'x-data-spreadsheet';
import 'x-data-spreadsheet/dist/xspreadsheet.css';

const DEFAULT_ROW_COUNT = 8;
const DEFAULT_COLUMN_COUNT = 4;
const MIN_COLUMN_LEN = 26;

const props = defineProps({
  modelValue: {
    type: String,
    default: ''
  }
});

const emit = defineEmits(['update:modelValue', 'commit']);

const editorRef = ref(null);
const sheetRef = ref(null);
const lastSerialized = ref('');
const applyingData = ref(false);
const wheelBlocker = ref(null);
const resizeHandler = ref(null);
const rafIdRef = ref(0);
const transitionContainerRef = ref(null);
const transitionEndHandlerRef = ref(null);
const themeObserverRef = ref(null);
const setBodyTableEditorOpen = (enabled) => {
  if (typeof document === 'undefined') {
    return;
  }
  document.body.classList.toggle('table-editor-open', Boolean(enabled));
};

const buildSheetStyle = () => ({
  bgcolor: '#ffffff',
  color: '#0a0a0a',
  align: 'left',
  valign: 'middle',
  textwrap: false,
  strike: false,
  underline: false,
  font: {
    name: 'Arial',
    size: 10,
    bold: false,
    italic: false
  },
  format: 'normal'
});

const createEmptyRows = (rowCount = DEFAULT_ROW_COUNT, colCount = DEFAULT_COLUMN_COUNT) =>
  Array.from({ length: rowCount }, () => Array.from({ length: colCount }, () => ''));

const normalizeRows = (sourceRows) => {
  if (!Array.isArray(sourceRows) || sourceRows.length === 0) {
    return createEmptyRows();
  }
  const width = Math.max(
    DEFAULT_COLUMN_COUNT,
    ...sourceRows.map((row) => (Array.isArray(row) ? row.length : 0))
  );
  return sourceRows.map((row) => {
    const values = Array.isArray(row) ? row : [];
    return Array.from({ length: width }, (_, index) => {
      const value = values[index];
      return value == null ? '' : String(value);
    });
  });
};

const parseRows = (rawContent) => {
  if (!rawContent || !String(rawContent).trim()) {
    return createEmptyRows();
  }
  try {
    const parsed = JSON.parse(rawContent);
    if (Array.isArray(parsed)) {
      return normalizeRows(parsed);
    }
    if (Array.isArray(parsed?.rows)) {
      return normalizeRows(parsed.rows);
    }
  } catch (error) {
    // ignore parse error and fallback
  }
  return createEmptyRows();
};

const serializeRows = (rows) =>
  JSON.stringify({
    type: 'table',
    version: 1,
    rows
  });

const rowsToSheetData = (rows) => {
  const rowCount = Math.max(DEFAULT_ROW_COUNT, rows.length);
  const colCount = Math.max(
    MIN_COLUMN_LEN,
    ...rows.map((row) => (Array.isArray(row) ? row.length : 0))
  );
  const sheetRows = { len: rowCount };
  rows.forEach((row, rowIndex) => {
    const cells = {};
    row.forEach((value, colIndex) => {
      cells[colIndex] = { text: value == null ? '' : String(value) };
    });
    sheetRows[rowIndex] = { cells };
  });
  return {
    name: 'Sheet1',
    freeze: 'A1',
    cols: { len: colCount },
    rows: sheetRows
  };
};

const firstSheet = (data) => {
  if (Array.isArray(data)) {
    return data[0] || {};
  }
  if (data && typeof data === 'object') {
    if (data.rows) {
      return data;
    }
    const candidates = Object.values(data);
    const found = candidates.find((item) => item && typeof item === 'object' && item.rows);
    return found || {};
  }
  return {};
};

const sheetDataToRows = (data) => {
  const sheet = firstSheet(data);
  const rowMap = sheet?.rows && typeof sheet.rows === 'object' ? sheet.rows : {};
  const rowIndexes = Object.keys(rowMap)
    .filter((key) => /^\d+$/.test(key))
    .map((key) => Number(key));

  let maxRow = DEFAULT_ROW_COUNT - 1;
  let maxCol = DEFAULT_COLUMN_COUNT - 1;

  rowIndexes.forEach((rowIndex) => {
    if (rowIndex > maxRow) {
      maxRow = rowIndex;
    }
    const cellMap = rowMap[rowIndex]?.cells || {};
    Object.keys(cellMap)
      .filter((key) => /^\d+$/.test(key))
      .forEach((key) => {
        const colIndex = Number(key);
        if (colIndex > maxCol) {
          maxCol = colIndex;
        }
      });
  });

  const rows = [];
  for (let r = 0; r <= maxRow; r += 1) {
    const current = [];
    const cellMap = rowMap[r]?.cells || {};
    for (let c = 0; c <= maxCol; c += 1) {
      const value = cellMap[c]?.text;
      current.push(value == null ? '' : String(value));
    }
    rows.push(current);
  }
  return rows;
};

const emitModelValueFromSheet = (sheetData) => {
  const rows = sheetDataToRows(sheetData);
  const serialized = serializeRows(rows);
  if (serialized === lastSerialized.value) {
    return;
  }
  lastSerialized.value = serialized;
  emit('update:modelValue', serialized);
};

const rerenderSheet = () => {
  const instance = sheetRef.value;
  if (!instance) {
    return;
  }
  instance.sheet?.reload?.();
  instance.reRender?.();
};

const scheduleRerender = () => {
  if (rafIdRef.value) {
    cancelAnimationFrame(rafIdRef.value);
  }
  rafIdRef.value = requestAnimationFrame(() => {
    rafIdRef.value = 0;
    rerenderSheet();
  });
};

const applyModelValue = async (value) => {
  const rows = parseRows(value);
  lastSerialized.value = serializeRows(rows);
  if (!sheetRef.value) {
    return;
  }
  applyingData.value = true;
  sheetRef.value.loadData(rowsToSheetData(rows));
  await nextTick();
  applyingData.value = false;
};

const applyThemeToSheet = () => {
  const instance = sheetRef.value;
  if (!instance?.data?.settings) {
    return;
  }
  const nextStyle = buildSheetStyle();
  const currentStyle = instance.data.settings.style || {};
  instance.data.settings.style = {
    ...currentStyle,
    ...nextStyle,
    font: {
      ...(currentStyle.font || {}),
      ...nextStyle.font
    }
  };
  scheduleRerender();
};

const initSpreadsheet = () => {
  if (!editorRef.value || sheetRef.value) {
    return;
  }
  const instance = new Spreadsheet(editorRef.value, {
    mode: 'edit',
    showToolbar: true,
    showGrid: true,
    showContextmenu: true,
    showBottomBar: false,
    row: {
      len: 100,
      height: 28
    },
    col: {
      len: 26,
      width: 120,
      indexWidth: 52,
      minWidth: 72
    },
    view: {
      height: () => editorRef.value?.clientHeight || 640,
      width: () => editorRef.value?.clientWidth || 960
    },
    style: buildSheetStyle()
  });
  sheetRef.value = instance;
  instance.change((data) => {
    if (applyingData.value) {
      return;
    }
    emitModelValueFromSheet(data);
  });
  instance.on('cell-edited', () => {
    emit('commit');
  });
};

defineExpose({
  reRender: () => scheduleRerender()
});

onMounted(async () => {
  setBodyTableEditorOpen(true);
  initSpreadsheet();
  const handler = (event) => {
    const container = editorRef.value;
    if (!container) {
      return;
    }
    if (!container.contains(event.target)) {
      return;
    }
    event.preventDefault();
  };
  window.addEventListener('wheel', handler, { passive: false, capture: true });
  window.addEventListener('mousewheel', handler, { passive: false, capture: true });
  window.addEventListener('DOMMouseScroll', handler, { passive: false, capture: true });
  wheelBlocker.value = handler;
  const onResize = () => rerenderSheet();
  window.addEventListener('resize', onResize);
  resizeHandler.value = onResize;
  const transitionContainer = editorRef.value?.closest('.editor-layout');
  if (transitionContainer) {
    const onTransitionEnd = () => {
      scheduleRerender();
    };
    transitionContainer.addEventListener('transitionend', onTransitionEnd);
    transitionContainerRef.value = transitionContainer;
    transitionEndHandlerRef.value = onTransitionEnd;
  }
  if (typeof MutationObserver !== 'undefined' && typeof document !== 'undefined') {
    const observer = new MutationObserver(() => {
      applyThemeToSheet();
    });
    observer.observe(document.body, { attributes: true, attributeFilter: ['class'] });
    themeObserverRef.value = observer;
  }
  applyThemeToSheet();
  await applyModelValue(props.modelValue);
  scheduleRerender();
});

onBeforeUnmount(() => {
  setBodyTableEditorOpen(false);
  if (themeObserverRef.value) {
    themeObserverRef.value.disconnect();
  }
  if (wheelBlocker.value) {
    window.removeEventListener('wheel', wheelBlocker.value, true);
    window.removeEventListener('mousewheel', wheelBlocker.value, true);
    window.removeEventListener('DOMMouseScroll', wheelBlocker.value, true);
  }
  if (sheetRef.value && typeof sheetRef.value.destroy === 'function') {
    sheetRef.value.destroy();
  }
  if (resizeHandler.value) {
    window.removeEventListener('resize', resizeHandler.value);
  }
  if (transitionContainerRef.value && transitionEndHandlerRef.value) {
    transitionContainerRef.value.removeEventListener('transitionend', transitionEndHandlerRef.value);
  }
  if (rafIdRef.value) {
    cancelAnimationFrame(rafIdRef.value);
  }
  rafIdRef.value = 0;
  themeObserverRef.value = null;
  transitionContainerRef.value = null;
  transitionEndHandlerRef.value = null;
  resizeHandler.value = null;
  wheelBlocker.value = null;
  sheetRef.value = null;
});

watch(
  () => props.modelValue,
  async (newValue) => {
    if (typeof newValue !== 'string') {
      return;
    }
    if (newValue === lastSerialized.value) {
      return;
    }
    await applyModelValue(newValue);
  }
);
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
