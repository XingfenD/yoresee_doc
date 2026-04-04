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
const resizeObserverRef = ref(null);
const rafIdRef = ref(0);

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
  const sheetRows = {};
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
  sheetRef.value?.reRender?.();
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
    }
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
  if (typeof ResizeObserver !== 'undefined' && editorRef.value) {
    const observer = new ResizeObserver(() => {
      scheduleRerender();
    });
    observer.observe(editorRef.value);
    resizeObserverRef.value = observer;
  }
  await applyModelValue(props.modelValue);
  scheduleRerender();
});

onBeforeUnmount(() => {
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
  if (resizeObserverRef.value) {
    resizeObserverRef.value.disconnect();
  }
  if (rafIdRef.value) {
    cancelAnimationFrame(rafIdRef.value);
  }
  resizeObserverRef.value = null;
  rafIdRef.value = 0;
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
}

.table-editor :deep(.x-spreadsheet-toolbar),
.table-editor :deep(.x-spreadsheet-bottombar) {
  background: var(--bg-white);
}

.dark-mode .table-editor :deep(.x-spreadsheet-toolbar),
.dark-mode .table-editor :deep(.x-spreadsheet-bottombar),
.dark-mode .table-editor :deep(.x-spreadsheet-sheet) {
  filter: invert(0.92) hue-rotate(180deg);
}

</style>
