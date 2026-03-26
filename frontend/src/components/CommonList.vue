<template>
  <div class="common-list" :class="{ 'is-dark': isDark, 'is-tree': mode === 'tree' }">
    <CommonListToolbar
      :show-title-bar="showTitleBar"
      :show-search="showSearch"
      :title="title"
      :search-placeholder="searchPlaceholder"
      :search-query="searchValue"
      @update:search-query="searchValue = $event"
      @search="emitSearch"
    >
      <template #title>
        <slot name="title">
          {{ title }}
        </slot>
      </template>
      <template #toolbar-actions>
        <slot name="toolbar-actions" />
      </template>
      <template #toolbar-right>
        <slot name="toolbar-right" />
      </template>
    </CommonListToolbar>
    <template v-if="rows.length === 0">
      <div class="list-empty">
        <el-empty :description="emptyText" />
      </div>
    </template>

    <template v-else>
      <CommonListTable
        v-if="mode !== 'tree'"
        :rows="rows"
        :columns="columns"
        :display-columns="displayColumns"
        :grid-template-columns="gridTemplateColumns"
        :tree-toggle-column-key="treeToggleColumnKey"
        :resolve-row-key="resolveRowKey"
        :align-class="alignClass"
      >
        <template v-for="name in forwardedSlotNames" :key="`table-${name}`" #[name]="slotProps">
          <slot :name="name" v-bind="slotProps || {}" />
        </template>
      </CommonListTable>

      <CommonListTreeTable
        v-else
        :tree-toggle-width="treeToggleWidth"
        :tree-data-columns="treeDataColumns"
        :tree-data-grid-template="treeDataGridTemplate"
        :tree-loading="treeLoading"
        :tree-flat-rows="treeFlatRows"
        :max-tree-indent-width="maxTreeIndentWidth"
        :toggle-scroll-left="toggleScrollLeft"
        :tree-indent="treeIndent"
        :tree-base-indent="treeBaseIndent"
        :tree-column-resolved-key="treeColumnResolvedKey"
        :resolve-row-key="resolveRowKey"
        :align-class="alignClass"
        :toggle-tree-node="toggleTreeNode"
        :is-tree-column="isTreeColumn"
        @toggle-scroll="setToggleScrollLeft"
      >
        <template v-for="name in forwardedSlotNames" :key="`tree-${name}`" #[name]="slotProps">
          <slot :name="name" v-bind="slotProps || {}" />
        </template>
      </CommonListTreeTable>
    </template>

    <CommonListPagination
      :show-pagination="showPagination"
      :current-page="paginationPage"
      :page-size="paginationPageSize"
      :page-sizes="pageSizes"
      :total="paginationTotal"
      :pagination-layout="paginationLayout"
      @update:current-page="paginationPage = $event"
      @update:page-size="paginationPageSize = $event"
      @page-change="handlePageChange"
      @size-change="handleSizeChange"
    />
  </div>
</template>

<script setup>
import { useCommonListState } from '@/composables/useCommonListState';
import { useForwardedSlotNames } from '@/composables/useForwardedSlotNames';
import CommonListToolbar from '@/components/common-list/CommonListToolbar.vue';
import CommonListPagination from '@/components/common-list/CommonListPagination.vue';
import CommonListTable from '@/components/common-list/CommonListTable.vue';
import CommonListTreeTable from '@/components/common-list/CommonListTreeTable.vue';

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
  'update:treeExpandedKeys',
  'tree-toggle'
]);

const forwardedSlotNames = useForwardedSlotNames(['title', 'toolbar-actions', 'toolbar-right']);

const {
  treeToggleColumnKey,
  treeToggleWidth,
  displayColumns,
  gridTemplateColumns,
  treeDataColumns,
  treeDataGridTemplate,
  resolveRowKey,
  treeColumnResolvedKey,
  paginationPage,
  paginationPageSize,
  paginationTotal,
  searchValue,
  emitSearch,
  toggleScrollLeft,
  setToggleScrollLeft,
  toggleTreeNode,
  isTreeColumn,
  treeFlatRows,
  maxTreeIndentWidth,
  handlePageChange,
  handleSizeChange,
  alignClass
} = useCommonListState(props, emit);
</script>

<style scoped>
.common-list {
  --list-cell-bg: var(--bg-white);
  --list-cell-text: var(--text-dark);
  --list-cell-border: #d6dbe3;
  --list-head-bg: #e5ebf2;
  --list-head-text: #1f2937;
  --list-head-border: #aeb8c6;
  background: var(--bg-white);
  border: 1px solid var(--border-color);
  border-radius: var(--border-radius-md);
  box-shadow: var(--shadow-sm);
  overflow: hidden;
}

.common-list.is-tree {
  overflow: hidden;
}

.list-empty {
  padding: 24px 0;
}

.common-list.is-dark {
  --list-cell-bg: #161b22;
  --list-cell-text: #e5e7eb;
  --list-cell-border: #2a313a;
  --list-head-bg: #202734;
  --list-head-text: #e5edf8;
  --list-head-border: #4a5668;
  background: #161b22;
  border-color: #161b22;
}
</style>
