<template>
  <NodeViewWrapper class="rich-table-node" contenteditable="false">
    <header class="rich-table-header" data-drag-handle>
      <span class="rich-table-badge">Table</span>
      <div class="rich-table-actions">
        <button type="button" class="rich-table-btn" @click="appendRow">
          + Row
        </button>
        <button type="button" class="rich-table-btn" @click="appendColumn">
          + Col
        </button>
      </div>
    </header>

    <div class="rich-table-body">
      <table class="rich-table-grid">
        <tbody>
          <tr v-for="(row, rowIndex) in tableModel.rows" :key="`row-${rowIndex}`">
            <td v-for="(cell, colIndex) in row" :key="`cell-${rowIndex}-${colIndex}`">
              <input
                :value="cell"
                class="rich-table-cell-input"
                @mousedown.stop
                @click.stop
                @keydown.stop
                @input="updateCell(rowIndex, colIndex, $event.target.value)"
              />
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </NodeViewWrapper>
</template>

<script setup>
import { ref, watch } from 'vue';
import { NodeViewWrapper } from '@tiptap/vue-3';
import {
  cloneRichTableModel,
  createRichTableModel,
  normalizeRichTableModel
} from './richTableModel';

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

const tableModel = ref(createRichTableModel());

const serializeRows = (rows) => JSON.stringify(rows || []);

const applyModel = (next) => {
  const normalized = normalizeRichTableModel(next);
  tableModel.value = normalized;
  props.updateAttributes({
    table: cloneRichTableModel(normalized)
  });
};

const updateCell = (rowIndex, colIndex, value) => {
  const next = cloneRichTableModel(tableModel.value);
  if (!Array.isArray(next.rows[rowIndex])) {
    return;
  }
  next.rows[rowIndex][colIndex] = String(value ?? '');
  applyModel(next);
};

const appendRow = () => {
  const next = cloneRichTableModel(tableModel.value);
  const columnCount = Math.max(1, Number(next.rows?.[0]?.length || 1));
  next.rows.push(Array.from({ length: columnCount }, () => ''));
  applyModel(next);
};

const appendColumn = () => {
  const next = cloneRichTableModel(tableModel.value);
  next.rows = next.rows.map((row) => [...row, '']);
  applyModel(next);
};

watch(
  () => props.node?.attrs?.table,
  (nextValue) => {
    const normalized = normalizeRichTableModel(nextValue);
    if (serializeRows(normalized.rows) === serializeRows(tableModel.value.rows)) {
      return;
    }
    tableModel.value = normalized;
  },
  { immediate: true }
);
</script>

<style scoped>
.rich-table-node {
  border: 1px solid var(--border-color);
  border-radius: 10px;
  background: var(--bg-white);
  overflow: hidden;
}

.rich-table-header {
  height: 34px;
  padding: 0 10px;
  border-bottom: 1px solid var(--border-color);
  display: flex;
  align-items: center;
  justify-content: space-between;
  background: color-mix(in srgb, var(--bg-light) 70%, transparent);
}

.rich-table-badge {
  font-size: 12px;
  font-weight: 600;
  color: var(--text-medium);
}

.rich-table-actions {
  display: inline-flex;
  align-items: center;
  gap: 6px;
}

.rich-table-btn {
  height: 24px;
  border: 1px solid var(--border-color);
  border-radius: 6px;
  background: var(--bg-white);
  color: var(--text-medium);
  padding: 0 8px;
  font-size: 12px;
  cursor: pointer;
}

.rich-table-body {
  overflow: auto;
  max-height: 380px;
}

.rich-table-grid {
  width: 100%;
  border-collapse: collapse;
  table-layout: fixed;
}

.rich-table-grid td {
  border: 1px solid color-mix(in srgb, var(--border-color) 90%, transparent);
  min-width: 120px;
  padding: 0;
}

.rich-table-cell-input {
  width: 100%;
  height: 36px;
  border: none;
  outline: none;
  padding: 7px 10px;
  box-sizing: border-box;
  background: transparent;
  color: var(--text-primary);
  font-size: 13px;
}

.rich-table-cell-input:focus {
  box-shadow: inset 0 0 0 2px color-mix(in srgb, var(--primary-color) 60%, transparent);
}
</style>
