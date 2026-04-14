import { ref, computed } from 'vue';

/**
 * Composable for drag-to-resize column widths.
 *
 * @param {import('vue').ComputedRef} displayColumns - reactive column definitions
 * @param {(columns: any[]) => string} buildGridTemplate - fn to build grid template string
 * @param {number} [minColumnWidth=60] - minimum column width in px
 */
export function useColumnResize(displayColumns, buildGridTemplate, minColumnWidth = 60) {
  // Map of column key -> overridden width in px (null = use column def default)
  const columnWidthOverrides = ref({});

  // Reset overrides when columns change identity (e.g. navigating to a different list)
  // Not done automatically — caller can call resetOverrides() if needed.

  const effectiveColumns = computed(() =>
    displayColumns.value.map((col) => {
      const override = columnWidthOverrides.value[col.key];
      if (override != null) {
        return { ...col, width: override, minWidth: undefined };
      }
      return col;
    })
  );

  const gridTemplateColumns = computed(() => buildGridTemplate(effectiveColumns.value));

  let dragState = null;

  function startResize(event, column) {
    event.preventDefault();
    event.stopPropagation();

    // Resolve starting width: override > column.width > measure from DOM
    let startWidth = columnWidthOverrides.value[column.key];
    if (startWidth == null) {
      if (column.width) {
        startWidth = typeof column.width === 'number' ? column.width : parseInt(column.width);
      } else {
        // Measure the header cell from the DOM
        const cell = event.currentTarget?.closest?.('.list-cell');
        startWidth = cell ? cell.getBoundingClientRect().width : 120;
      }
    }

    dragState = {
      column,
      startX: event.clientX,
      startWidth,
    };

    document.addEventListener('mousemove', onMouseMove);
    document.addEventListener('mouseup', onMouseUp);
    document.body.style.cursor = 'col-resize';
    document.body.style.userSelect = 'none';
  }

  function onMouseMove(event) {
    if (!dragState) return;
    const delta = event.clientX - dragState.startX;
    const newWidth = Math.max(minColumnWidth, dragState.startWidth + delta);
    columnWidthOverrides.value = {
      ...columnWidthOverrides.value,
      [dragState.column.key]: newWidth,
    };
  }

  function onMouseUp() {
    dragState = null;
    document.removeEventListener('mousemove', onMouseMove);
    document.removeEventListener('mouseup', onMouseUp);
    document.body.style.cursor = '';
    document.body.style.userSelect = '';
  }

  function resetOverrides() {
    columnWidthOverrides.value = {};
  }

  return {
    effectiveColumns,
    gridTemplateColumns,
    startResize,
    resetOverrides,
  };
}
