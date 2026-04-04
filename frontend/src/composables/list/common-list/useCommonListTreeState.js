import { computed, ref, watch } from 'vue';

export function useCommonListTreeState(props, emit, resolveRowKey) {
  const treeToggleButtonSize = 22;
  const toggleScrollLeft = ref(0);
  const internalExpanded = ref(new Set());

  const resolveTreeKey = (row, index) => {
    if (typeof props.treeKeyField === 'function') {
      return props.treeKeyField(row, index);
    }
    return row?.[props.treeKeyField] ?? resolveRowKey(row, index);
  };

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

  const isTreeColumn = (column, index, treeToggleColumnKey, columns = []) => {
    if (column.key === treeToggleColumnKey) {
      return false;
    }
    if (column.isIndexColumn) {
      return false;
    }
    if (props.treeColumnKey) {
      return column.key === props.treeColumnKey;
    }
    const firstContentColumnIndex = columns.findIndex(
      (item) => item?.key !== treeToggleColumnKey && !item?.isIndexColumn
    );
    if (firstContentColumnIndex < 0) {
      return index === 0;
    }
    return index === firstContentColumnIndex;
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

  const setToggleScrollLeft = (value) => {
    toggleScrollLeft.value = value || 0;
  };

  return {
    toggleScrollLeft,
    setToggleScrollLeft,
    toggleTreeNode,
    isTreeColumn,
    treeFlatRows,
    maxTreeIndentWidth
  };
}
