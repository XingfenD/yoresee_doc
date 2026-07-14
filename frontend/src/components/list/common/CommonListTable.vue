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

.list-row:last-child .list-cell {
  border-bottom: none;
}
</style>

<style>
@import '@/styles/list-cell.css';
@import '@/styles/column-resize.css';
</style>
