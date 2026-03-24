<template>
  <div class="common-list" :class="{ 'is-dark': isDark }">
    <div v-if="showSearch" class="list-toolbar">
      <slot name="toolbar-actions" />
      <el-input
        v-model="searchValue"
        :placeholder="searchPlaceholder"
        clearable
        class="list-search"
        @clear="emitSearch"
        @input="emitSearch"
      />
    </div>
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

    <div v-if="showPagination" class="pagination-container">
      <el-pagination
        v-model:current-page="paginationPage"
        v-model:page-size="paginationPageSize"
        :page-sizes="pageSizes"
        :total="paginationTotal"
        :layout="paginationLayout"
        :hide-on-single-page="false"
        @current-change="handlePageChange"
        @size-change="handleSizeChange"
      />
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
  },
  showPagination: {
    type: Boolean,
    default: false
  },
  total: {
    type: Number,
    default: 0
  },
  currentPage: {
    type: Number,
    default: 1
  },
  pageSize: {
    type: Number,
    default: 10
  },
  pageSizes: {
    type: Array,
    default: () => [10, 20, 50, 100]
  },
  paginationLayout: {
    type: String,
    default: 'total, prev, pager, next, jumper'
  },
  showSearch: {
    type: Boolean,
    default: false
  },
  searchQuery: {
    type: String,
    default: ''
  },
  searchPlaceholder: {
    type: String,
    default: ''
  }
});

const emit = defineEmits([
  'update:currentPage',
  'update:pageSize',
  'page-change',
  'size-change',
  'update:searchQuery',
  'search'
]);

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

const paginationPage = computed({
  get: () => props.currentPage,
  set: (value) => emit('update:currentPage', value)
});

const paginationPageSize = computed({
  get: () => props.pageSize,
  set: (value) => emit('update:pageSize', value)
});

const paginationTotal = computed(() => {
  if (typeof props.total === 'number' && props.total > 0) {
    return props.total;
  }
  if (typeof props.total === 'string' && props.total.trim() !== '') {
    const parsed = Number(props.total);
    if (Number.isFinite(parsed) && parsed > 0) {
      return parsed;
    }
  }
  return Array.isArray(props.rows) ? props.rows.length : 0;
});

const searchValue = computed({
  get: () => props.searchQuery,
  set: (value) => emit('update:searchQuery', value)
});

const emitSearch = () => {
  emit('search', searchValue.value);
};

const handlePageChange = (page) => {
  emit('page-change', page);
};

const handleSizeChange = (size) => {
  emit('size-change', size);
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

.list-toolbar {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  gap: 12px;
  padding: var(--spacing-md);
  border-bottom: 1px solid var(--border-color);
  background: var(--bg-white);
}

.list-search {
  max-width: 320px;
  flex: 1;
  margin-left: auto;
}

.pagination-container {
  display: flex;
  justify-content: flex-end;
  padding: var(--spacing-md);
  border-top: 1px solid var(--border-color);
  background: var(--bg-white);
}

.pagination-container :deep(.el-pagination) {
  color: var(--text-dark);
}

.pagination-container :deep(.btn-prev),
.pagination-container :deep(.btn-next),
.pagination-container :deep(.el-pager li) {
  background: #ffffff;
  color: var(--text-dark);
  border: 1px solid var(--border-color);
}

.pagination-container :deep(.btn-quicknext),
.pagination-container :deep(.btn-quickprev) {
  pointer-events: none;
  cursor: default;
  background: #ffffff;
  color: var(--text-dark);
  border: 1px solid var(--border-color);
}

.pagination-container :deep(.el-pager li.is-active) {
  background: var(--primary-color);
  color: #ffffff;
  border-color: var(--primary-color);
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

.common-list.is-dark .list-toolbar {
  background: #161b22;
  border-bottom-color: #2a313a;
}

.common-list.is-dark .pagination-container {
  background: #161b22;
  border-top-color: #2a313a;
}

.common-list.is-dark .pagination-container :deep(.el-pagination) {
  color: #e5e7eb;
}

.common-list.is-dark .pagination-container :deep(.btn-prev),
.common-list.is-dark .pagination-container :deep(.btn-next),
.common-list.is-dark .pagination-container :deep(.el-pager li) {
  background: #1b2230;
  color: #e5e7eb;
  border-color: #2a313a;
}

.common-list.is-dark .pagination-container :deep(.btn-quicknext),
.common-list.is-dark .pagination-container :deep(.btn-quickprev) {
  pointer-events: none;
  cursor: default;
  background: #1b2230;
  color: #e5e7eb;
  border-color: #2a313a;
}

.common-list.is-dark .pagination-container :deep(.el-pager li.is-active) {
  background: var(--primary-color);
  color: #ffffff;
  border-color: var(--primary-color);
}
</style>
