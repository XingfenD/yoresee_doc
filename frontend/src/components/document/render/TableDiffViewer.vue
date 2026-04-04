<template>
  <div class="table-diff-viewer">
    <div class="table-diff-head">
      <div class="head-cell">{{ leftTitle }}</div>
      <div class="head-cell">{{ rightTitle }}</div>
    </div>

    <div class="table-diff-body">
      <div class="table-pane">
        <table class="diff-table">
          <thead>
            <tr>
              <th class="index-head"></th>
              <th v-for="label in columnLabels" :key="`l-${label}`" class="col-head">{{ label }}</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="rowIndex in rowIndexes" :key="`left-row-${rowIndex}`">
              <th class="index-cell">{{ rowIndex + 1 }}</th>
              <td
                v-for="colIndex in colIndexes"
                :key="`left-cell-${rowIndex}-${colIndex}`"
                class="value-cell"
                :class="leftCellClass(rowIndex, colIndex)"
              >
                {{ resolveCell(leftRowsNormalized, rowIndex, colIndex) || '\u00A0' }}
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <div class="table-pane">
        <table class="diff-table">
          <thead>
            <tr>
              <th class="index-head"></th>
              <th v-for="label in columnLabels" :key="`r-${label}`" class="col-head">{{ label }}</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="rowIndex in rowIndexes" :key="`right-row-${rowIndex}`">
              <th class="index-cell">{{ rowIndex + 1 }}</th>
              <td
                v-for="colIndex in colIndexes"
                :key="`right-cell-${rowIndex}-${colIndex}`"
                class="value-cell"
                :class="rightCellClass(rowIndex, colIndex)"
              >
                {{ resolveCell(rightRowsNormalized, rowIndex, colIndex) || '\u00A0' }}
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <div class="table-diff-legend">
      <span class="legend-item is-added">Added</span>
      <span class="legend-item is-removed">Removed</span>
      <span class="legend-item is-changed">Changed</span>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue';

const props = defineProps({
  leftRows: {
    type: Array,
    default: () => []
  },
  rightRows: {
    type: Array,
    default: () => []
  },
  leftTitle: {
    type: String,
    default: ''
  },
  rightTitle: {
    type: String,
    default: ''
  }
});

const normalizeRows = (rows, width) =>
  Array.from({ length: Array.isArray(rows) ? rows.length : 0 }, (_, rowIndex) => {
    const row = Array.isArray(rows[rowIndex]) ? rows[rowIndex] : [];
    return Array.from({ length: width }, (_, colIndex) => {
      const value = row[colIndex];
      return value === null || value === undefined ? '' : String(value);
    });
  });

const columnLabelAt = (index) => {
  let value = Number(index) + 1;
  let label = '';
  while (value > 0) {
    const remainder = (value - 1) % 26;
    label = String.fromCharCode(65 + remainder) + label;
    value = Math.floor((value - 1) / 26);
  }
  return label || 'A';
};

const rowCount = computed(() => Math.max(
  Array.isArray(props.leftRows) ? props.leftRows.length : 0,
  Array.isArray(props.rightRows) ? props.rightRows.length : 0
));

const columnCount = computed(() => {
  const leftMax = Array.isArray(props.leftRows)
    ? Math.max(0, ...props.leftRows.map((row) => (Array.isArray(row) ? row.length : 0)))
    : 0;
  const rightMax = Array.isArray(props.rightRows)
    ? Math.max(0, ...props.rightRows.map((row) => (Array.isArray(row) ? row.length : 0)))
    : 0;
  return Math.max(1, leftMax, rightMax);
});

const leftRowsNormalized = computed(() => normalizeRows(props.leftRows, columnCount.value));
const rightRowsNormalized = computed(() => normalizeRows(props.rightRows, columnCount.value));
const rowIndexes = computed(() => Array.from({ length: rowCount.value }, (_, index) => index));
const colIndexes = computed(() => Array.from({ length: columnCount.value }, (_, index) => index));
const columnLabels = computed(() => colIndexes.value.map((index) => columnLabelAt(index)));

const resolveCell = (rows, rowIndex, colIndex) =>
  String(rows?.[rowIndex]?.[colIndex] ?? '');

const cellStatus = (rowIndex, colIndex) => {
  const leftValue = resolveCell(leftRowsNormalized.value, rowIndex, colIndex);
  const rightValue = resolveCell(rightRowsNormalized.value, rowIndex, colIndex);
  if (leftValue === rightValue) {
    return 'equal';
  }
  if (!leftValue && rightValue) {
    return 'added';
  }
  if (leftValue && !rightValue) {
    return 'removed';
  }
  return 'changed';
};

const leftCellClass = (rowIndex, colIndex) => {
  const status = cellStatus(rowIndex, colIndex);
  if (status === 'changed') {
    return 'is-changed';
  }
  if (status === 'removed') {
    return 'is-removed';
  }
  return '';
};

const rightCellClass = (rowIndex, colIndex) => {
  const status = cellStatus(rowIndex, colIndex);
  if (status === 'changed') {
    return 'is-changed';
  }
  if (status === 'added') {
    return 'is-added';
  }
  return '';
};
</script>

<style scoped>
.table-diff-viewer {
  border: 1px solid var(--border-color);
  border-radius: var(--border-radius-sm);
  background: var(--bg-white);
  overflow: hidden;
}

.table-diff-head {
  display: grid;
  grid-template-columns: 1fr 1fr;
  border-bottom: 1px solid var(--border-color);
  background: var(--bg-light);
}

.head-cell {
  padding: 8px 12px;
  font-size: 13px;
  color: var(--text-medium);
  border-right: 1px solid var(--border-color);
}

.head-cell:last-child {
  border-right: 0;
}

.table-diff-body {
  display: grid;
  grid-template-columns: 1fr 1fr;
  min-height: 320px;
}

.table-pane {
  overflow: auto;
  border-right: 1px solid var(--border-color);
}

.table-pane:last-child {
  border-right: 0;
}

.diff-table {
  width: max-content;
  min-width: 100%;
  border-collapse: collapse;
}

.diff-table th,
.diff-table td {
  border: 1px solid var(--border-color);
  padding: 6px 8px;
  font-size: 12px;
  min-width: 100px;
  max-width: 220px;
  white-space: pre-wrap;
  word-break: break-word;
}

.col-head {
  position: sticky;
  top: 0;
  z-index: 2;
  background: var(--bg-light);
  color: var(--text-medium);
  font-weight: 600;
}

.index-head,
.index-cell {
  min-width: 44px !important;
  max-width: 44px !important;
  text-align: center;
}

.index-head {
  position: sticky;
  left: 0;
  z-index: 3;
  background: var(--bg-light);
}

.index-cell {
  position: sticky;
  left: 0;
  z-index: 1;
  background: var(--bg-light);
  color: var(--text-light);
}

.value-cell {
  background: var(--bg-white);
  color: var(--text-dark);
}

.is-added {
  background: rgba(34, 197, 94, 0.18);
}

.is-removed {
  background: rgba(239, 68, 68, 0.18);
}

.is-changed {
  background: rgba(245, 158, 11, 0.18);
}

.table-diff-legend {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 12px;
  border-top: 1px solid var(--border-color);
}

.legend-item {
  padding: 2px 8px;
  border-radius: 999px;
  font-size: 11px;
  color: var(--text-medium);
}

.dark-mode .table-diff-head,
.dark-mode .index-head,
.dark-mode .index-cell,
.dark-mode .col-head {
  background: #1f2937;
}
</style>

