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
          v-for="column in columns"
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
            {{ row?.[column.key] ?? '-' }}
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
  treeToggleColumnKey: { type: String, default: '__tree_toggle__' },
  resolveRowKey: { type: Function, required: true },
  alignClass: { type: Function, required: true }
});
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
  color: var(--text-dark);
  border-bottom: 1px solid #d6dbe3;
  display: flex;
  align-items: center;
  gap: 8px;
}

.list-cell--head {
  font-size: 12px;
  font-weight: 700;
  color: #1f2937;
  background: #e5ebf2;
  border-bottom-color: #aeb8c6;
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

:global(.common-list.is-dark) .list-cell {
  color: #e5e7eb;
  border-bottom-color: #2a313a;
  background: #161b22;
}

:global(.common-list.is-dark) .list-cell--head {
  color: #e5edf8;
  background: #202734;
  border-bottom-color: #4a5668;
}
</style>
