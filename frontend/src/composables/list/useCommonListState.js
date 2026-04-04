import { useCommonListSharedState } from '@/composables/list/common-list/useCommonListSharedState';
import { useCommonListTreeState } from '@/composables/list/common-list/useCommonListTreeState';

export function useCommonListState(props, emit) {
  const shared = useCommonListSharedState(props, emit);
  const tree = useCommonListTreeState(props, emit, shared.resolveRowKey);

  return {
    ...shared,
    toggleScrollLeft: tree.toggleScrollLeft,
    setToggleScrollLeft: tree.setToggleScrollLeft,
    toggleTreeNode: tree.toggleTreeNode,
    isTreeColumn: (column, index) =>
      tree.isTreeColumn(column, index, shared.treeToggleColumnKey, shared.treeDataColumns.value),
    treeFlatRows: tree.treeFlatRows,
    maxTreeIndentWidth: tree.maxTreeIndentWidth
  };
}
