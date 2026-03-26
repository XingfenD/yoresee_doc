import { useCommonListSharedState } from '@/composables/common-list/useCommonListSharedState';
import { useCommonListTreeState } from '@/composables/common-list/useCommonListTreeState';

export function useCommonListState(props, emit) {
  const shared = useCommonListSharedState(props, emit);
  const tree = useCommonListTreeState(props, emit, shared.resolveRowKey);

  return {
    ...shared,
    toggleScrollLeft: tree.toggleScrollLeft,
    setToggleScrollLeft: tree.setToggleScrollLeft,
    toggleTreeNode: tree.toggleTreeNode,
    isTreeColumn: (column, index) => tree.isTreeColumn(column, index, shared.treeToggleColumnKey),
    treeFlatRows: tree.treeFlatRows,
    maxTreeIndentWidth: tree.maxTreeIndentWidth
  };
}
