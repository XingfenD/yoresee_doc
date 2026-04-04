<template>
  <div class="table-preview-viewer">
    <el-empty v-if="!hasRows" />
    <div v-else class="table-preview-scroll">
      <table class="table-preview-grid">
        <thead>
          <tr>
            <th class="row-index-head"></th>
            <th v-for="label in columnLabels" :key="label" class="col-head">
              {{ label }}
            </th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="(row, rowIndex) in normalizedRows" :key="`row-${rowIndex}`">
            <th class="row-index-cell">{{ rowIndex + 1 }}</th>
            <td v-for="(cell, colIndex) in row" :key="`cell-${rowIndex}-${colIndex}`" class="value-cell">
              {{ cell || '\u00A0' }}
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue';

const props = defineProps({
  rows: {
    type: Array,
    default: () => []
  }
});

const normalizeRows = (rows) => {
  if (!Array.isArray(rows) || rows.length === 0) {
    return [];
  }
  const width = Math.max(1, ...rows.map((row) => (Array.isArray(row) ? row.length : 0)));
  return rows.map((row) => {
    const values = Array.isArray(row) ? row : [];
    return Array.from({ length: width }, (_, index) => {
      const value = values[index];
      return value === null || value === undefined ? '' : String(value);
    });
  });
};

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

const normalizedRows = computed(() => normalizeRows(props.rows));
const hasRows = computed(() => normalizedRows.value.length > 0);
const columnLabels = computed(() => {
  const width = normalizedRows.value[0]?.length || 0;
  return Array.from({ length: width }, (_, index) => columnLabelAt(index));
});
</script>

<style scoped>
.table-preview-viewer {
  width: 100%;
  height: 100%;
  min-height: 0;
  border: 1px solid var(--border-color);
  border-radius: var(--border-radius-sm);
  background: var(--bg-white);
  overflow: hidden;
}

.table-preview-scroll {
  width: 100%;
  height: 100%;
  overflow: auto;
}

.table-preview-grid {
  width: max-content;
  min-width: 100%;
  border-collapse: collapse;
  table-layout: fixed;
}

.table-preview-grid th,
.table-preview-grid td {
  border: 1px solid var(--border-color);
  padding: 8px 10px;
  font-size: 13px;
  min-width: 120px;
  max-width: 260px;
  white-space: pre-wrap;
  word-break: break-word;
}

.table-preview-grid thead th {
  position: sticky;
  top: 0;
  z-index: 2;
  background: var(--bg-light);
  color: var(--text-medium);
  font-weight: 600;
}

.row-index-head,
.row-index-cell {
  min-width: 52px !important;
  max-width: 52px !important;
  text-align: center;
}

.row-index-cell {
  position: sticky;
  left: 0;
  z-index: 1;
  background: var(--bg-light);
  color: var(--text-light);
  font-weight: 500;
}

.row-index-head {
  position: sticky;
  left: 0;
  z-index: 3;
}

.value-cell {
  color: var(--text-dark);
  background: var(--bg-white);
}
</style>

