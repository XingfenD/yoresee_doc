<template>
  <div v-if="showPagination" class="pagination-container">
    <el-pagination
      v-model:current-page="innerCurrentPage"
      v-model:page-size="innerPageSize"
      :page-sizes="pageSizes"
      :total="total"
      :layout="paginationLayout"
      :hide-on-single-page="false"
      @current-change="(page) => $emit('page-change', page)"
      @size-change="(size) => $emit('size-change', size)"
    />
  </div>
</template>

<script setup>
import { computed } from 'vue';

const props = defineProps({
  showPagination: { type: Boolean, default: false },
  currentPage: { type: Number, default: 1 },
  pageSize: { type: Number, default: 10 },
  pageSizes: { type: Array, default: () => [10, 20, 50, 100] },
  total: { type: Number, default: 0 },
  paginationLayout: { type: String, default: 'total, prev, pager, next, jumper' }
});

const emit = defineEmits(['update:currentPage', 'update:pageSize', 'page-change', 'size-change']);

const innerCurrentPage = computed({
  get: () => props.currentPage,
  set: (value) => emit('update:currentPage', value)
});

const innerPageSize = computed({
  get: () => props.pageSize,
  set: (value) => emit('update:pageSize', value)
});
</script>

<style scoped>
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

:global(.common-list.is-dark) .pagination-container {
  background: #161b22;
  border-top-color: #2a313a;
}

:global(.common-list.is-dark) .pagination-container :deep(.el-pagination) {
  color: #e5e7eb;
}

:global(.common-list.is-dark) .pagination-container :deep(.btn-prev),
:global(.common-list.is-dark) .pagination-container :deep(.btn-next),
:global(.common-list.is-dark) .pagination-container :deep(.el-pager li) {
  background: #1b2230;
  color: #e5e7eb;
  border-color: #2a313a;
}

:global(.common-list.is-dark) .pagination-container :deep(.btn-quicknext),
:global(.common-list.is-dark) .pagination-container :deep(.btn-quickprev) {
  background: #1b2230;
  color: #e5e7eb;
  border-color: #2a313a;
}
</style>
