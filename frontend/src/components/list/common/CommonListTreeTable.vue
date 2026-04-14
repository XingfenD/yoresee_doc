<template>
  <div class="tree-table-scroll">
    <div
      class="tree-table"
      :style="{
        '--tree-toggle-width': `${resolvedToggleWidth}px`,
        '--tree-data-template': resizableDataTemplate
      }"
    >
      <!-- Header -->
      <div class="tree-head">
        <div class="tree-toggle-head list-cell list-cell--head" />
        <div class="tree-data-head">
          <div
            v-for="column in effectiveDataColumns"
            :key="`tree-head-${column.key}`"
            class="list-cell list-cell--head"
            :class="[column.className, alignClass(column.headerAlign || column.align)]"
          >
            <slot :name="`header-${column.key}`" :column="column">
              {{ column.label }}
            </slot>
            <span
              class="col-resize-handle"
              @mousedown="startResize($event, column)"
            />
          </div>
        </div>
      </div>

      <!-- Body -->
      <div class="tree-body" v-loading="treeLoading">
        <div
          v-for="(row, rowIndex) in treeFlatRows"
          :key="resolveRowKey(row.raw, rowIndex)"
          class="tree-row"
        >
          <!-- Sticky indent/toggle cell -->
          <div class="tree-toggle-cell">
            <div
              class="tree-toggle-inner-content"
              :style="{ paddingLeft: `${row.level * treeIndent + treeBaseIndent}px` }"
            >
              <button
                v-if="row.hasChildren"
                class="tree-toggle"
                type="button"
                @click.stop="toggleTreeNode(row)"
              >
                <el-icon :size="14">
                  <Minus v-if="row.expanded" />
                  <Plus v-else />
                </el-icon>
              </button>
              <span v-else class="tree-leaf-indicator" />
            </div>
          </div>

          <!-- Data cells -->
          <div class="tree-data-row">
            <div
              v-for="(column, columnIndex) in effectiveDataColumns"
              :key="`${resolveRowKey(row.raw, rowIndex)}-${column.key}`"
              class="list-cell list-cell--tree"
              :class="[column.className, alignClass(column.align)]"
            >
              <template v-if="isTreeColumn(column, columnIndex)">
                <div class="tree-cell">
                  <slot name="tree-cell" :row="row.raw" :level="row.level">
                    <span class="tree-node-label">{{ row.raw?.[treeColumnResolvedKey] ?? '-' }}</span>
                  </slot>
                </div>
              </template>
              <template v-else-if="column.isIndexColumn">
                <slot
                  :name="`cell-${column.key}`"
                  :row="row.raw"
                  :row-index="rowIndex"
                  :column="column"
                  :value="row.raw?.[column.key]"
                >
                  {{ resolveSerialNumber(rowIndex, currentPage, pageSize) }}
                </slot>
              </template>
              <template v-else>
                <slot
                  :name="`cell-${column.key}`"
                  :row="row.raw"
                  :row-index="rowIndex"
                  :column="column"
                  :value="row.raw?.[column.key]"
                >
                  {{ row.raw?.[column.key] ?? '-' }}
                </slot>
              </template>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue';
import { Plus, Minus } from '@element-plus/icons-vue';
import { useColumnResize } from '@/composables/list/useColumnResize.js';

const props = defineProps({
  treeToggleWidth: { type: Number, default: 81 },
  treeDataColumns: { type: Array, default: () => [] },
  treeDataGridTemplate: { type: String, default: '1fr' },
  currentPage: { type: Number, default: 1 },
  pageSize: { type: Number, default: 10 },
  treeLoading: { type: Boolean, default: false },
  treeFlatRows: { type: Array, default: () => [] },
  maxTreeIndentWidth: { type: Number, default: 0 },
  toggleScrollLeft: { type: Number, default: 0 },
  treeIndent: { type: Number, default: 16 },
  treeBaseIndent: { type: Number, default: 6 },
  treeColumnResolvedKey: { type: String, default: 'label' },
  resolveRowKey: { type: Function, required: true },
  alignClass: { type: Function, required: true },
  toggleTreeNode: { type: Function, required: true },
  isTreeColumn: { type: Function, required: true },
  buildGridTemplate: { type: Function, default: null },
});

defineEmits(['toggle-scroll']);

const resolvedToggleWidth = computed(() =>
  Math.max(props.treeToggleWidth, props.maxTreeIndentWidth)
);

const treeDataColumnsRef = computed(() => props.treeDataColumns);

const buildFn = computed(() =>
  props.buildGridTemplate ??
  ((cols) =>
    cols
      .map((c) => (c.width ? `${typeof c.width === 'number' ? c.width + 'px' : c.width}` : `${c.flex || 1}fr`))
      .join(' '))
);

const { effectiveColumns: effectiveDataColumns, gridTemplateColumns: resizableDataTemplate, startResize } =
  useColumnResize(treeDataColumnsRef, (cols) => buildFn.value(cols));

const resolveSerialNumber = (rowIndex, currentPage, pageSize) => {
  const page = Number.isFinite(Number(currentPage)) ? Number(currentPage) : 1;
  const size = Number.isFinite(Number(pageSize)) ? Number(pageSize) : 0;
  if (size <= 0) return rowIndex + 1;
  return (Math.max(page, 1) - 1) * size + rowIndex + 1;
};
</script>

<style scoped>
.tree-table-scroll {
  overflow-x: auto;
}

.tree-table {
  display: flex;
  flex-direction: column;
  min-width: max-content;
  width: 100%;
}

.tree-head {
  display: grid;
  grid-template-columns: var(--tree-toggle-width) 1fr;
}

.tree-toggle-head {
  justify-content: center;
  position: sticky;
  left: 0;
  z-index: 2;
}

.tree-data-head {
  display: grid;
  grid-template-columns: var(--tree-data-template);
}

.tree-body {
  display: flex;
  flex-direction: column;
}

.tree-row {
  display: grid;
  grid-template-columns: var(--tree-toggle-width) 1fr;
}

.tree-toggle-cell {
  position: sticky;
  left: 0;
  z-index: 1;
  background: var(--list-cell-bg, var(--bg-white));
  border-right: 1px solid var(--list-cell-border, var(--border-color));
  border-bottom: 1px solid var(--list-cell-border, #d6dbe3);
  display: flex;
  align-items: center;
}

.tree-toggle-inner-content {
  display: flex;
  align-items: center;
  height: 100%;
  padding: 12px 14px;
  box-sizing: border-box;
  white-space: nowrap;
}

.tree-data-row {
  display: grid;
  grid-template-columns: var(--tree-data-template);
}

.list-cell {
  min-width: 0;
  padding: 12px 14px;
  font-size: 13px;
  color: var(--list-cell-text, var(--text-dark));
  border-bottom: 1px solid var(--list-cell-border, #d6dbe3);
  background: var(--list-cell-bg, var(--bg-white));
  display: flex;
  align-items: center;
  gap: 8px;
  position: relative;
}

.list-cell--head {
  font-size: 12px;
  font-weight: 700;
  color: var(--list-head-text, #1f2937);
  background: var(--list-head-bg, #e5ebf2);
  border-bottom-color: var(--list-head-border, #aeb8c6);
  overflow: visible;
}

.list-cell--tree {
  border-bottom: 1px solid var(--list-cell-border, #d6dbe3);
}

.tree-row:last-child .tree-toggle-cell,
.tree-row:last-child .list-cell--tree {
  border-bottom: none;
}

.is-left   { justify-content: flex-start; text-align: left; }
.is-center { justify-content: center;     text-align: center; }
.is-right  { justify-content: flex-end;   text-align: right; }

.tree-cell {
  display: flex;
  align-items: center;
  gap: 8px;
  width: 100%;
  box-sizing: border-box;
}

.tree-toggle {
  width: 22px;
  height: 22px;
  flex-shrink: 0;
  border: 1px solid var(--border-color);
  border-radius: 4px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  color: var(--list-cell-text, var(--text-dark));
  background: var(--list-cell-bg, var(--bg-white));
  cursor: pointer;
}

.tree-toggle:hover {
  color: var(--primary-color);
  border-color: var(--primary-color);
}

.tree-leaf-indicator {
  width: 6px;
  height: 6px;
  border-radius: 999px;
  background: #9aa4b2;
  display: inline-block;
  flex-shrink: 0;
}

.tree-node-label {
  font-size: 14px;
  color: var(--text-medium);
}

/* ── Resize handle ── */
.col-resize-handle {
  position: absolute;
  right: 0;
  top: 0;
  bottom: 0;
  width: 6px;
  cursor: col-resize;
  z-index: 1;
}

.col-resize-handle::after {
  content: '';
  position: absolute;
  right: 2px;
  top: 20%;
  bottom: 20%;
  width: 2px;
  border-radius: 1px;
  background: transparent;
  transition: background 0.15s;
}

.list-cell--head:hover .col-resize-handle::after,
.col-resize-handle:hover::after {
  background: var(--primary-color, #165dff);
}
</style>
