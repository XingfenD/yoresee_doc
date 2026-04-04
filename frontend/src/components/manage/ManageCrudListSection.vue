<template>
  <ManageLayout>
    <ManageSection :title="sectionTitle" :plain="sectionPlain" :body-padding="bodyPadding">
      <CommonList
        :rows="rows"
        :columns="columns"
        :is-dark="isDark"
        :row-key="rowKey"
        :mode="mode"
        :empty-text="emptyText"
        :show-pagination="showPagination"
        :total="total"
        :current-page="currentPage"
        :page-size="pageSize"
        :page-sizes="pageSizes"
        :pagination-layout="paginationLayout"
        :show-search="showSearch"
        :search-query="searchQuery"
        :search-placeholder="searchPlaceholder"
        :show-index-column="showIndexColumn"
        :index-column-label="indexColumnLabel"
        :index-column-width="indexColumnWidth"
        :index-column-align="indexColumnAlign"
        :show-title-bar="showTitleBar"
        :title="title"
        :tree-loading="treeLoading"
        :tree-column-key="treeColumnKey"
        :tree-key-field="treeKeyField"
        @update:current-page="emit('update:currentPage', $event)"
        @update:page-size="emit('update:pageSize', $event)"
        @page-change="emit('page-change', $event)"
        @size-change="emit('size-change', $event)"
        @update:search-query="emit('update:searchQuery', $event)"
        @search="emit('search', $event)"
      >
        <template v-for="name in forwardedSlotNames" :key="name" #[name]="slotProps">
          <slot :name="name" v-bind="slotProps || {}" />
        </template>
      </CommonList>
    </ManageSection>
  </ManageLayout>
</template>

<script setup>
import { useForwardedSlotNames } from '@/composables/list/useForwardedSlotNames';
import ManageLayout from '@/components/manage/ManageLayout.vue';
import ManageSection from '@/components/manage/ManageSection.vue';
import CommonList from '@/components/list/CommonList.vue';

defineProps({
  sectionTitle: {
    type: String,
    default: ''
  },
  sectionPlain: {
    type: Boolean,
    default: false
  },
  bodyPadding: {
    type: String,
    default: 'none'
  },
  rows: {
    type: Array,
    default: () => []
  },
  columns: {
    type: Array,
    default: () => []
  },
  isDark: {
    type: Boolean,
    default: false
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
  showPagination: {
    type: Boolean,
    default: true
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
  showTitleBar: {
    type: Boolean,
    default: false
  },
  title: {
    type: String,
    default: ''
  },
  treeLoading: {
    type: Boolean,
    default: false
  },
  treeColumnKey: {
    type: String,
    default: 'name'
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
  'search'
]);

const forwardedSlotNames = useForwardedSlotNames(['default']);
</script>
