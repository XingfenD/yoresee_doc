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
        :current-page="paginationPage"
        :page-size="paginationPageSize"
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
        :current-page="paginationPage"
        :page-size="paginationPageSize"
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
import { useCommonListState } from '@/composables/list/useCommonListState';
import { useForwardedSlotNames } from '@/composables/list/useForwardedSlotNames';
import CommonListToolbar from '@/components/list/common/CommonListToolbar.vue';
import CommonListPagination from '@/components/list/common/CommonListPagination.vue';
import CommonListTable from '@/components/list/common/CommonListTable.vue';
import CommonListTreeTable from '@/components/list/common/CommonListTreeTable.vue';

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
  showIndexColumn: {
    type: Boolean,
    default: false
  },
  indexColumnLabel: {
    type: String,
    default: '序号'
  },
  indexColumnWidth: {
    type: [String, Number],
    default: 72
  },
  indexColumnAlign: {
    type: String,
    default: 'center'
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
  --list-pagination-bg: var(--bg-white);
  --list-pagination-border: #d6dbe3;
  --list-pagination-text: var(--text-medium);
  --list-pagination-strong-text: var(--text-dark);
  --list-pagination-item-bg: #ffffff;
  --list-pagination-item-border: #d6dbe3;
  --list-pagination-item-hover-bg: #f3f7ff;
  --list-pagination-item-hover-border: #b8cbff;
  --list-pagination-item-active-bg: var(--primary-color);
  --list-pagination-item-active-text: #ffffff;
  --list-pagination-disabled-bg: #f4f6fa;
  --list-pagination-disabled-text: #b6bfcc;
  --list-pagination-input-bg: #ffffff;
  --list-pagination-input-border: #d6dbe3;
  --list-pagination-input-text: var(--text-dark);
  --list-pagination-focus-ring: rgba(22, 93, 255, 0.18);
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
  --list-pagination-bg: #141a23;
  --list-pagination-border: #2a313a;
  --list-pagination-text: #96a1b4;
  --list-pagination-strong-text: #e5e7eb;
  --list-pagination-item-bg: #1a2230;
  --list-pagination-item-border: #2f3b4d;
  --list-pagination-item-hover-bg: #233147;
  --list-pagination-item-hover-border: #45679a;
  --list-pagination-item-active-bg: #3370ff;
  --list-pagination-item-active-text: #ffffff;
  --list-pagination-disabled-bg: #131a24;
  --list-pagination-disabled-text: #5f6a7d;
  --list-pagination-input-bg: #1a2230;
  --list-pagination-input-border: #2f3b4d;
  --list-pagination-input-text: #e5e7eb;
  --list-pagination-focus-ring: rgba(64, 128, 255, 0.28);
  background: #161b22;
  border-color: #161b22;
}
</style>
