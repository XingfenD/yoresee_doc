import { computed } from 'vue';

export function useCommonListSharedState(props, emit) {
  const normalizeWidth = (val) => {
    if (typeof val === 'number') {
      return `${val}px`;
    }
    return val;
  };

  const treeToggleColumnKey = '__tree_toggle__';
  const treeToggleWidth = 81;
  const indexColumnKey = '__list_index__';
  const indexColumn = computed(() => {
    if (!props.showIndexColumn) {
      return null;
    }
    return {
      key: indexColumnKey,
      label: props.indexColumnLabel,
      width: props.indexColumnWidth,
      align: props.indexColumnAlign,
      headerAlign: props.indexColumnAlign,
      className: 'list-index-column',
      isIndexColumn: true
    };
  });

  const displayColumns = computed(() => {
    if (props.mode === 'tree') {
      return [
        { key: treeToggleColumnKey, width: treeToggleWidth, align: 'center', className: 'tree-toggle-column' },
        ...(indexColumn.value ? [indexColumn.value] : []),
        ...props.columns
      ];
    }
    return [...(indexColumn.value ? [indexColumn.value] : []), ...props.columns];
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
  const treeDataColumns = computed(() => {
    if (props.mode !== 'tree') {
      return props.columns || [];
    }
    return [...(indexColumn.value ? [indexColumn.value] : []), ...(props.columns || [])];
  });
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

  return {
    treeToggleColumnKey,
    treeToggleWidth,
    indexColumnKey,
    displayColumns,
    gridTemplateColumns,
    treeDataColumns,
    treeDataGridTemplate,
    buildGridTemplate,
    resolveRowKey,
    treeColumnResolvedKey,
    paginationPage,
    paginationPageSize,
    paginationTotal,
    searchValue,
    emitSearch,
    handlePageChange,
    handleSizeChange,
    alignClass
  };
}
