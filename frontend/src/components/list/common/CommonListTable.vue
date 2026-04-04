<template>
  <template v-if="displayColumns.length > 0">
    <div class="list-head" :style="{ gridTemplateColumns }">
      <div
        v-for="column in displayColumns"
        :key="`head-${column.key}`"
        class="list-cell list-cell--head"
        :class="[column.className, alignClass(column.headerAlign || column.align)]"
      >
        <slot :name="`header-${column.key}`" :column="column">
          {{ column.key === treeToggleColumnKey ? '' : column.label }}
        </slot>
      </div>
    </div>
    <div class="list-body">
      <div
        v-for="(row, rowIndex) in rows"
        :key="resolveRowKey(row, rowIndex)"
        class="list-row"
        :style="{ gridTemplateColumns }"
      >
        <div
          v-for="column in displayColumns"
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
  </template>
</template>

<script setup>
defineProps({
  rows: { type: Array, default: () => [] },
  columns: { type: Array, default: () => [] },
  displayColumns: { type: Array, default: () => [] },
  gridTemplateColumns: { type: String, default: '1fr' },
  currentPage: { type: Number, default: 1 },
  pageSize: { type: Number, default: 10 },
  treeToggleColumnKey: { type: String, default: '__tree_toggle__' },
  resolveRowKey: { type: Function, required: true },
  alignClass: { type: Function, required: true }
});

const resolveSerialNumber = (rowIndex, currentPage, pageSize) => {
  const page = Number.isFinite(Number(currentPage)) ? Number(currentPage) : 1;
  const size = Number.isFinite(Number(pageSize)) ? Number(pageSize) : 0;
  if (size <= 0) {
    return rowIndex + 1;
  }
  return (Math.max(page, 1) - 1) * size + rowIndex + 1;
};
</script>

<style scoped>
.list-head,
.list-row {
  display: grid;
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
}

.list-cell--head {
  font-size: 12px;
  font-weight: 700;
  color: var(--list-head-text, #1f2937);
  background: var(--list-head-bg, #e5ebf2);
  border-bottom-color: var(--list-head-border, #aeb8c6);
}

.list-row:last-child .list-cell {
  border-bottom: none;
}

.is-left {
  justify-content: flex-start;
  text-align: left;
}

.is-center {
  justify-content: center;
  text-align: center;
}

.is-right {
  justify-content: flex-end;
  text-align: right;
}

</style>
