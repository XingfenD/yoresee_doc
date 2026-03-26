<template>
  <div v-if="showPagination" class="pagination-container">
    <el-pagination
      v-model:current-page="innerCurrentPage"
      v-model:page-size="innerPageSize"
      :page-sizes="pageSizes"
      :total="total"
      :layout="paginationLayout"
      popper-class="common-list-pagination-popper"
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
  align-items: center;
  justify-content: flex-end;
  padding: 10px 12px;
  border-top: 1px solid var(--list-pagination-border, var(--border-color));
  background: var(--list-pagination-bg, var(--bg-white));
}

.pagination-container :deep(.el-pagination) {
  --el-pagination-bg-color: transparent;
  --el-pagination-text-color: var(--list-pagination-text, var(--text-medium));
  --el-pagination-button-bg-color: var(--list-pagination-item-bg, #ffffff);
  --el-pagination-button-color: var(--list-pagination-strong-text, var(--text-dark));
  --el-pagination-button-disabled-bg-color: var(--list-pagination-disabled-bg, #f4f6fa);
  --el-pagination-button-disabled-color: var(--list-pagination-disabled-text, #b6bfcc);
  --el-pagination-hover-color: var(--primary-color);
  display: flex;
  align-items: center;
  gap: 6px;
  color: var(--list-pagination-text, var(--text-medium));
  font-weight: 500;
}

.pagination-container :deep(.btn-prev),
.pagination-container :deep(.btn-next),
.pagination-container :deep(.el-pager li) {
  min-width: 32px;
  height: 32px;
  line-height: 30px;
  border-radius: 8px;
  background: var(--list-pagination-item-bg, #ffffff);
  color: var(--list-pagination-strong-text, var(--text-dark));
  border: 1px solid var(--list-pagination-item-border, var(--border-color));
  transition: all 0.16s ease;
}

.pagination-container :deep(.btn-prev:hover:not(:disabled)),
.pagination-container :deep(.btn-next:hover:not(:disabled)),
.pagination-container :deep(.el-pager li:not(.is-active):hover) {
  background: var(--list-pagination-item-hover-bg, #f3f7ff);
  border-color: var(--list-pagination-item-hover-border, #b8cbff);
  color: var(--primary-color);
}

.pagination-container :deep(.btn-prev:disabled),
.pagination-container :deep(.btn-next:disabled),
.pagination-container :deep(.el-pager li.is-disabled) {
  background: var(--list-pagination-disabled-bg, #f4f6fa);
  color: var(--list-pagination-disabled-text, #b6bfcc);
  border-color: var(--list-pagination-item-border, var(--border-color));
}

.pagination-container :deep(.el-pager li.is-active) {
  background: var(--list-pagination-item-active-bg, var(--primary-color));
  color: var(--list-pagination-item-active-text, #ffffff);
  border-color: var(--list-pagination-item-active-bg, var(--primary-color));
}

.pagination-container :deep(.btn-quicknext),
.pagination-container :deep(.btn-quickprev) {
  color: var(--list-pagination-text, var(--text-medium));
}

.pagination-container :deep(.el-pagination__total),
.pagination-container :deep(.el-pagination__jump) {
  color: var(--list-pagination-text, var(--text-medium));
}

.pagination-container :deep(.el-pagination__jump .el-input__wrapper),
.pagination-container :deep(.el-pagination__sizes .el-input__wrapper),
.pagination-container :deep(.el-pagination__sizes .el-select__wrapper) {
  background: var(--list-pagination-input-bg, #ffffff);
  border: 1px solid var(--list-pagination-input-border, var(--border-color));
  box-shadow: none;
  border-radius: 8px;
  transition: all 0.16s ease;
}

.pagination-container :deep(.el-pagination__jump .el-input__wrapper:hover),
.pagination-container :deep(.el-pagination__sizes .el-input__wrapper:hover),
.pagination-container :deep(.el-pagination__sizes .el-select__wrapper:hover) {
  border-color: var(--list-pagination-item-hover-border, #b8cbff);
}

.pagination-container :deep(.el-pagination__jump .el-input.is-focus .el-input__wrapper),
.pagination-container :deep(.el-pagination__sizes .el-select .el-input.is-focus .el-input__wrapper),
.pagination-container :deep(.el-pagination__sizes .el-select.is-focused .el-select__wrapper) {
  border-color: var(--primary-color);
  box-shadow: 0 0 0 3px var(--list-pagination-focus-ring, rgba(22, 93, 255, 0.18));
}

.pagination-container :deep(.el-pagination__jump .el-input__inner),
.pagination-container :deep(.el-pagination__sizes .el-input__inner),
.pagination-container :deep(.el-pagination__sizes .el-select__selected-item),
.pagination-container :deep(.el-pagination__sizes .el-select__placeholder) {
  color: var(--list-pagination-input-text, var(--text-dark));
}

.pagination-container :deep(.el-pagination__sizes .el-input) {
  width: 96px;
}

:global(.common-list-pagination-popper .el-select-dropdown__item) {
  font-size: 13px;
}

:global(.dark-mode .common-list-pagination-popper.el-popper) {
  background: #1a2230;
  border-color: #2f3b4d;
}

:global(.dark-mode .common-list-pagination-popper .el-select-dropdown__item) {
  color: #e5e7eb;
}

:global(.dark-mode .common-list-pagination-popper .el-select-dropdown__item.is-hovering),
:global(.dark-mode .common-list-pagination-popper .el-select-dropdown__item:hover) {
  background: #233147;
}

:global(.dark-mode .common-list-pagination-popper .el-select-dropdown__item.is-selected) {
  color: #ffffff;
  background: #3370ff;
}
</style>
