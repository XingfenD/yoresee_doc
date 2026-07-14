<template>
  <div v-if="displayColumns.length > 0" class="list-table-scroll">
    <div class="list-head" :style="{ gridTemplateColumns: resizableGridTemplate }">
      <div
        v-for="column in effectiveColumns"
        :key="`head-${column.key}`"
        class="list-cell list-cell--head"
        :class="[column.className, alignClass(column.headerAlign || column.align)]"
      >
        <slot :name="`header-${column.key}`" :column="column">
          {{ column.key === treeToggleColumnKey ? '' : column.label }}
        </slot>
        <span
          v-if="column.key !== treeToggleColumnKey"
          class="col-resize-handle"
          @mousedown="startResize($event, column)"
        />
      </div>
    </div>
    <div class="list-body">
      <div
        v-for="(row, rowIndex) in rows"
        :key="resolveRowKey(row, rowIndex)"
        class="list-row"
        :style="{ gridTemplateColumns: resizableGridTemplate }"
      >
        <div
          v-for="column in effectiveColumns"
          :key="`${resolveRowKey(row, rowIndex)}-${column.key}`"
          class="list-cell"
          :class="[column.className, alignClass(column.align)]"
        >
          <slot
            :name="`cell-${column.key}`"
            :row="row"
            :row-index="rowIndex"
            :column="column"
            :value="row?.[column.key]"
          >
            <template v-if="column.isIndexColumn">
              {{ resolveSerialNumber(rowIndex, currentPage, pageSize) }}
            </template>
            <template v-else>
              {{ row?.[column.key] ?? '-' }}
            </template>
          </slot>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue';
import { useColumnResize } from '@/composables/list/useColumnResize.js';

const props = defineProps({
  rows: { type: Array, default: () => [] },
  columns: { type: Array, default: () => [] },
  displayColumns: { type: Array, default: () => [] },
  gridTemplateColumns: { type: String, default: '1fr' },
  currentPage: { type: Number, default: 1 },
  pageSize: { type: Number, default: 10 },
  treeToggleColumnKey: { type: String, default: '__tree_toggle__' },
  resolveRowKey: { type: Function, required: true },
  alignClass: { type: Function, required: true },
  buildGridTemplate: { type: Function, required: true },
  resolveSerialNumber: { type: Function, required: true },
});

const displayColumnsRef = computed(() => props.displayColumns);

const { effectiveColumns, gridTemplateColumns: resizableGridTemplate, startResize } =
  useColumnResize(displayColumnsRef, (cols) => props.buildGridTemplate(cols));
</script>

<style scoped>
.list-table-scroll {
  overflow-x: auto;
}

.list-head,
.list-row {
  display: grid;
  min-width: max-content;
  width: 100%;
}

.list-cell {
  min-width: 0;
  padding: 12px 14px;
  font-size: 13px;
  color: var(--list-cell-text, var(--text-dark));
  border-bottom: 1px solid var(--list-cell-border, #d6dbe3);
  background: var(--list-cell-bg, var(--bg-white));
  display: flex;
  align-items: center;
  gap: 8px;
  position: relative;
}

.list-cell--head {
  font-size: 12px;
  font-weight: 700;
  color: var(--list-head-text, #1f2937);
  background: var(--list-head-bg, #e5ebf2);
  border-bottom-color: var(--list-head-border, #aeb8c6);
  overflow: visible;
}

.list-row:last-child .list-cell {
  border-bottom: none;
}

.is-left   { justify-content: flex-start; text-align: left; }
.is-center { justify-content: center;     text-align: center; }
.is-right  { justify-content: flex-end;   text-align: right; }
</style>

<style>
@import '@/styles/column-resize.css';
</style>
