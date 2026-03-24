<template>
  <div class="common-list" :class="{ 'is-dark': isDark, 'is-tree': mode === 'tree' }">
    <div v-if="showTitleBar || showSearch" class="list-titlebar">
      <div class="list-title">
        <slot name="title">
          {{ title }}
        </slot>
      </div>
      <div class="list-title-actions">
        <slot name="toolbar-actions" />
        <el-input
          v-if="showSearch"
          v-model="searchValue"
          :placeholder="searchPlaceholder"
          clearable
          class="list-search"
          @clear="emitSearch"
          @input="emitSearch"
        />
        <slot name="toolbar-right" />
      </div>
    </div>
    <template v-if="rows.length === 0">
      <div class="list-empty">
        <el-empty :description="emptyText" />
      </div>
    </template>

    <template v-else-if="mode !== 'tree'">
      <div v-if="displayColumns.length > 0" class="list-head" :style="{ gridTemplateColumns }">
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

    <template v-else>
      <div class="tree-table" :style="{ '--tree-toggle-width': `${treeToggleWidth}px` }">
        <div class="tree-head">
          <div class="tree-toggle-head list-cell list-cell--head"></div>
          <div class="tree-data-head" :style="{ gridTemplateColumns: treeDataGridTemplate }">
            <div
              v-for="column in treeDataColumns"
              :key="`tree-head-${column.key}`"
              class="list-cell list-cell--head"
              :class="[column.className, alignClass(column.headerAlign || column.align)]"
            >
              <slot :name="`header-${column.key}`" :column="column">
                {{ column.label }}
              </slot>
            </div>
          </div>
        </div>
        <div class="tree-body" v-loading="treeLoading">
          <div
            v-for="(row, rowIndex) in treeFlatRows"
            :key="resolveRowKey(row.raw, rowIndex)"
            class="tree-row"
          >
            <div class="tree-toggle-cell">
              <div
                class="tree-toggle-inner"
                :style="{
                  width: `${maxTreeIndentWidth}px`,
                  transform: `translateX(${-toggleScrollLeft}px)`
                }"
              >
                <div class="tree-toggle-inner-content" :style="{ paddingLeft: `${row.level * treeIndent + treeBaseIndent}px` }">
                  <button
                    v-if="row.hasChildren"
                    class="tree-toggle"
                    type="button"
                    @click.stop="toggleTreeNode(row)"
                  >
                    <el-icon :size="14">
                      <Minus v-if="row.expanded" />
                      <Plus v-else />
                    </el-icon>
                  </button>
                  <span v-else class="tree-leaf-indicator" />
                </div>
              </div>
            </div>
            <div class="tree-data-row" :style="{ gridTemplateColumns: treeDataGridTemplate }">
              <div
                v-for="(column, columnIndex) in treeDataColumns"
                :key="`${resolveRowKey(row.raw, rowIndex)}-${column.key}`"
                class="list-cell list-cell--tree"
                :class="[column.className, alignClass(column.align)]"
              >
                <template v-if="isTreeColumn(column, columnIndex)">
                  <div class="tree-cell">
                    <slot name="tree-cell" :row="row.raw" :level="row.level">
                      <span class="tree-node-label">{{ row.raw?.[treeColumnResolvedKey] ?? '-' }}</span>
                    </slot>
                  </div>
                </template>
                <template v-else>
                  <slot
                    :name="`cell-${column.key}`"
                    :row="row.raw"
                    :row-index="rowIndex"
                    :column="column"
                    :value="row.raw?.[column.key]"
                  >
                    {{ row.raw?.[column.key] ?? '-' }}
                  </slot>
                </template>
              </div>
            </div>
          </div>
          <div class="tree-toggle-scrollbar" ref="toggleScrollbarRef">
            <div class="tree-toggle-scrollbar-inner" :style="{ width: `${maxTreeIndentWidth}px` }" />
          </div>
        </div>
      </div>
    </template>

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
import { computed, ref, watch, onMounted, onBeforeUnmount } from 'vue';
import { Plus, Minus } from '@element-plus/icons-vue';

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
  mode: {
    type: String,
    default: 'table'
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
  showTitleBar: {
    type: Boolean,
    default: false
  },
  title: {
    type: String,
    default: ''
  },
  searchQuery: {
    type: String,
    default: ''
  },
  searchPlaceholder: {
    type: String,
    default: ''
  },
  treeLoading: {
    type: Boolean,
    default: false
  },
  treeChildrenKey: {
    type: String,
    default: 'children'
  },
  treeColumnKey: {
    type: String,
    default: 'label'
  },
  treeColumnLabel: {
    type: String,
    default: ''
  },
  treeIndent: {
    type: Number,
    default: 16
  },
  treeBaseIndent: {
    type: Number,
    default: 6
  },
  treeDefaultExpandAll: {
    type: Boolean,
    default: false
  },
  treeExpandedKeys: {
    type: Array,
    default: null
  },
  treeKeyField: {
    type: [String, Function],
    default: 'id'
  }
});

const emit = defineEmits([
  'update:currentPage',
  'update:pageSize',
  'page-change',
  'size-change',
  'update:searchQuery',
  'search',
  'tree-node-click',
  'update:treeExpandedKeys',
  'tree-toggle'
]);

const normalizeWidth = (val) => {
  if (typeof val === 'number') {
    return `${val}px`;
  }
  return val;
};

const treeToggleColumnKey = '__tree_toggle__';
const treeToggleWidth = 81;
const treeToggleButtonSize = 22;

const displayColumns = computed(() => {
  if (props.mode === 'tree') {
    return [
      { key: treeToggleColumnKey, width: treeToggleWidth, align: 'center', className: 'tree-toggle-column' },
      ...props.columns
    ];
  }
  return props.columns;
});

const buildGridTemplate = (columns) => {
  if (!columns.length) {
    return '1fr';
  }
  return columns
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
};

const gridTemplateColumns = computed(() => buildGridTemplate(displayColumns.value));

const treeDataColumns = computed(() => props.columns || []);

const treeDataGridTemplate = computed(() => buildGridTemplate(treeDataColumns.value));

const resolveRowKey = (row, rowIndex) => {
  if (typeof props.rowKey === 'function') {
    return props.rowKey(row, rowIndex);
  }
  return row?.[props.rowKey] ?? rowIndex;
};

const treeColumnResolvedKey = computed(() => {
  if (props.treeColumnKey) {
    return props.treeColumnKey;
  }
  if (props.columns.length > 0) {
    return props.columns[0].key;
  }
  return 'label';
});

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

const resolveTreeKey = (row, index) => {
  if (typeof props.treeKeyField === 'function') {
    return props.treeKeyField(row, index);
  }
  return row?.[props.treeKeyField] ?? resolveRowKey(row, index);
};

const internalExpanded = ref(new Set());
const toggleScrollLeft = ref(0);
const toggleScrollbarRef = ref(null);

const initExpanded = () => {
  const next = new Set();
  if (Array.isArray(props.treeExpandedKeys)) {
    props.treeExpandedKeys.forEach((key) => next.add(String(key)));
  } else if (props.treeDefaultExpandAll) {
    const collect = (nodes) => {
      nodes.forEach((node, idx) => {
        const key = resolveTreeKey(node, idx);
        if (key !== undefined && key !== null) {
          next.add(String(key));
        }
        const children = node?.[props.treeChildrenKey];
        if (Array.isArray(children) && children.length) {
          collect(children);
        }
      });
    };
    collect(props.rows || []);
  }
  internalExpanded.value = next;
};

initExpanded();

watch(
  () => props.treeExpandedKeys,
  (value) => {
    if (Array.isArray(value)) {
      internalExpanded.value = new Set(value.map((key) => String(key)));
    }
  }
);

const handleToggleScroll = (event) => {
  toggleScrollLeft.value = event.target.scrollLeft || 0;
};

onMounted(() => {
  if (toggleScrollbarRef.value) {
    toggleScrollbarRef.value.addEventListener('scroll', handleToggleScroll, { passive: true });
  }
});

onBeforeUnmount(() => {
  if (toggleScrollbarRef.value) {
    toggleScrollbarRef.value.removeEventListener('scroll', handleToggleScroll);
  }
});

const isExpanded = (row, index) => {
  if (Array.isArray(props.treeExpandedKeys)) {
    const key = resolveTreeKey(row, index);
    return props.treeExpandedKeys.includes(key);
  }
  const key = resolveTreeKey(row, index);
  return internalExpanded.value.has(String(key));
};

const toggleTreeNode = (row) => {
  const key = resolveTreeKey(row.raw, row.index);
  if (key === undefined || key === null) {
    return;
  }
  const keyString = String(key);
  const next = new Set(internalExpanded.value);
  if (next.has(keyString)) {
    next.delete(keyString);
  } else {
    next.add(keyString);
  }
  if (!Array.isArray(props.treeExpandedKeys)) {
    internalExpanded.value = next;
  }
  emit('update:treeExpandedKeys', Array.from(next));
  emit('tree-toggle', { key, expanded: next.has(keyString), row: row.raw });
};

const isTreeColumn = (column, index) => {
  if (column.key === treeToggleColumnKey) {
    return false;
  }
  if (props.treeColumnKey) {
    return column.key === props.treeColumnKey;
  }
  return index === 0;
};

const treeFlatRows = computed(() => {
  if (props.mode !== 'tree') {
    return [];
  }
  const result = [];
  const walk = (nodes, level) => {
    nodes.forEach((node, index) => {
      const children = node?.[props.treeChildrenKey];
      const hasChildren = Array.isArray(children) && children.length > 0;
      const expanded = hasChildren ? isExpanded(node, index) : false;
      result.push({
        raw: node,
        level,
        hasChildren,
        expanded,
        index
      });
      if (hasChildren && expanded) {
        walk(children, level + 1);
      }
    });
  };
  walk(props.rows || [], 0);
  return result;
});

const maxTreeLevel = computed(() => {
  if (!treeFlatRows.value.length) {
    return 0;
  }
  return treeFlatRows.value.reduce((max, row) => Math.max(max, row.level), 0);
});

const maxTreeIndentWidth = computed(() => {
  return props.treeBaseIndent + maxTreeLevel.value * props.treeIndent + treeToggleButtonSize;
});

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

const handleTreeNodeClick = (data, node) => {
  emit('tree-node-click', data, node);
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
  width: 100%;
}

.common-list.is-tree {
  overflow: hidden;
}

.common-list.is-tree .list-head,
.common-list.is-tree .list-row {
  width: 100%;
}

.common-list.is-tree .tree-toggle-column {
  overflow-x: auto;
  scrollbar-width: thin;
}

.common-list.is-tree .tree-toggle-column::-webkit-scrollbar {
  height: 6px;
}

.common-list.is-tree .tree-toggle-column::-webkit-scrollbar-thumb {
  background: rgba(120, 132, 158, 0.5);
  border-radius: 999px;
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

.list-titlebar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  padding: 10px 12px;
  border-bottom: 1px solid var(--border-color);
  background: var(--bg-white);
}

.list-title {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 17px;
  font-weight: 600;
  color: var(--text-dark);
  line-height: 1.2;
}

.list-title-actions {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-left: auto;
}

.list-search {
  max-width: 320px;
  flex: 1;
  margin-top: -2px;
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

.tree-table {
  display: flex;
  flex-direction: column;
}

.tree-toggle-head {
  justify-content: center;
}

.tree-data-head {
  display: grid;
}

.tree-head {
  display: grid;
  grid-template-columns: var(--tree-toggle-width) 1fr;
}

.tree-body {
  display: flex;
  flex-direction: column;
}

.tree-row {
  display: grid;
  grid-template-columns: var(--tree-toggle-width) 1fr;
}

.tree-toggle-cell {
  border-right: 1px solid var(--border-color);
  overflow: hidden;
  display: flex;
  align-items: center;
}

.tree-toggle-inner {
  display: flex;
  align-items: center;
  height: 100%;
  will-change: transform;
}

.tree-toggle-inner-content {
  display: flex;
  align-items: center;
  height: 100%;
  padding: 12px 14px;
  box-sizing: border-box;
}

.tree-data-row {
  display: grid;
}

.list-cell--tree {
  border-bottom: 1px solid #d6dbe3;
}

.tree-row .list-cell--tree {
  border-bottom: 1px solid #d6dbe3;
}

.tree-toggle-scrollbar {
  width: var(--tree-toggle-width);
  overflow-x: auto;
  overflow-y: hidden;
  border-right: 1px solid var(--border-color);
}

.tree-toggle-scrollbar-inner {
  height: 1px;
}

.tree-cell {
  display: flex;
  align-items: center;
  gap: 8px;
  width: 100%;
  box-sizing: border-box;
}

.tree-toggle {
  width: 22px;
  height: 22px;
  border: 1px solid var(--border-color);
  border-radius: 4px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  color: var(--text-dark);
  background: var(--bg-white);
  cursor: pointer;
}

.tree-toggle-column {
  justify-content: center;
}

.tree-toggle:hover {
  color: var(--primary-color);
  border-color: var(--primary-color);
}

.tree-toggle-placeholder {
  width: 22px;
  height: 22px;
  display: inline-block;
}

.tree-leaf-indicator {
  width: 6px;
  height: 6px;
  border-radius: 999px;
  background: #9aa4b2;
  display: inline-block;
}

.tree-node-label {
  font-size: 14px;
  color: var(--text-medium);
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

.common-list.is-dark .list-titlebar {
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

.common-list.is-dark .list-cell--tree {
  border-bottom-color: #2a313a;
}

.common-list.is-dark .tree-toggle-cell {
  border-right-color: #2a313a;
}

.common-list.is-dark .tree-toggle-scrollbar {
  border-right-color: #2a313a;
}

.common-list.is-dark .tree-toggle {
  background: #1b2230;
  color: #e5e7eb;
  border-color: #2a313a;
}

.common-list.is-dark .tree-toggle:hover {
  color: #ffffff;
  border-color: #4a90ff;
}

.common-list.is-dark .tree-node-label {
  color: #e5e7eb;
}

.common-list.is-dark .tree-leaf-indicator {
  background: #6b7280;
}
</style>
