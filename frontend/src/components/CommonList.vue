<template>
  <div class="common-list" :class="{ 'is-dark': isDark }">
    <div class="list-head" :style="{ gridTemplateColumns }">
      <div
        v-for="column in columns"
        :key="`head-${column.key}`"
        class="list-cell list-cell--head"
        :class="[column.className, alignClass(column.headerAlign || column.align)]"
      >
        <slot :name="`header-${column.key}`" :column="column">
          {{ column.label }}
        </slot>
      </div>
    </div>

    <div v-if="rows.length === 0" class="list-empty">
      <el-empty :description="emptyText" />
    </div>

    <div v-else class="list-body">
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
  </div>
</template>

<script setup>
import { computed } from 'vue';

const props = defineProps({
  rows: {
    type: Array,
    default: () => []
  },
  columns: {
    type: Array,
    default: () => []
  },
  rowKey: {
    type: [String, Function],
    default: 'id'
  },
  emptyText: {
    type: String,
    default: ''
  },
  isDark: {
    type: Boolean,
    default: false
  }
});

const normalizeWidth = (val) => {
  if (typeof val === 'number') {
    return `${val}px`;
  }
  return val;
};

const gridTemplateColumns = computed(() => {
  if (!props.columns.length) {
    return '1fr';
  }
  return props.columns
    .map((column) => {
      if (column.width) {
        return normalizeWidth(column.width);
      }
      if (column.minWidth) {
        return `minmax(${normalizeWidth(column.minWidth)}, ${column.flex || 1}fr)`;
      }
      return `${column.flex || 1}fr`;
    })
    .join(' ');
});

const resolveRowKey = (row, rowIndex) => {
  if (typeof props.rowKey === 'function') {
    return props.rowKey(row, rowIndex);
  }
  return row?.[props.rowKey] ?? rowIndex;
};

const alignClass = (align) => {
  if (align === 'center') return 'is-center';
  if (align === 'right') return 'is-right';
  return 'is-left';
};
</script>

<style scoped>
.common-list {
  background: var(--bg-white);
  border: 1px solid var(--border-color);
  border-radius: var(--border-radius-md);
  box-shadow: var(--shadow-sm);
  overflow: hidden;
}

.list-head,
.list-row {
  display: grid;
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

.list-empty {
  padding: 24px 0;
}

.common-list.is-dark {
  background: #161b22;
  border-color: #161b22;
}

.common-list.is-dark .list-cell {
  color: #e5e7eb;
  border-bottom-color: #2a313a;
  background: #161b22;
}

.common-list.is-dark .list-cell--head {
  color: #e5edf8;
  background: #202734;
  border-bottom-color: #4a5668;
}
</style>
