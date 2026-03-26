<template>
  <div class="tree-table" :style="{ '--tree-toggle-width': `${treeToggleWidth}px` }">
    <div class="tree-head">
      <div class="tree-toggle-head list-cell list-cell--head"></div>
      <div class="tree-data-head" :style="{ gridTemplateColumns: treeDataGridTemplate }">
        <div
          v-for="column in treeDataColumns"
          :key="`tree-head-${column.key}`"
          class="list-cell list-cell--head"
          :class="[column.className, alignClass(column.headerAlign || column.align)]"
        >
          <slot :name="`header-${column.key}`" :column="column">
            {{ column.label }}
          </slot>
        </div>
      </div>
    </div>
    <div class="tree-body" v-loading="treeLoading">
      <div
        v-for="(row, rowIndex) in treeFlatRows"
        :key="resolveRowKey(row.raw, rowIndex)"
        class="tree-row"
      >
        <div class="tree-toggle-cell">
          <div
            class="tree-toggle-inner"
            :style="{
              width: `${maxTreeIndentWidth}px`,
              transform: `translateX(${-toggleScrollLeft}px)`
            }"
          >
            <div class="tree-toggle-inner-content" :style="{ paddingLeft: `${row.level * treeIndent + treeBaseIndent}px` }">
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
        </div>
        <div class="tree-data-row" :style="{ gridTemplateColumns: treeDataGridTemplate }">
          <div
            v-for="(column, columnIndex) in treeDataColumns"
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
      <div class="tree-toggle-scrollbar" @scroll="$emit('toggle-scroll', $event.target.scrollLeft || 0)">
        <div class="tree-toggle-scrollbar-inner" :style="{ width: `${maxTreeIndentWidth}px` }" />
      </div>
    </div>
  </div>
</template>

<script setup>
import { Plus, Minus } from '@element-plus/icons-vue';

defineProps({
  treeToggleWidth: { type: Number, default: 81 },
  treeDataColumns: { type: Array, default: () => [] },
  treeDataGridTemplate: { type: String, default: '1fr' },
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
  isTreeColumn: { type: Function, required: true }
});

defineEmits(['toggle-scroll']);
</script>

<style scoped>
.tree-table {
  display: flex;
  flex-direction: column;
}

.tree-toggle-head {
  justify-content: center;
}

.tree-data-head {
  display: grid;
}

.tree-head {
  display: grid;
  grid-template-columns: var(--tree-toggle-width) 1fr;
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
  border-right: 1px solid var(--list-cell-border, var(--border-color));
  overflow: hidden;
  display: flex;
  align-items: center;
}

.tree-toggle-inner {
  display: flex;
  align-items: center;
  height: 100%;
  will-change: transform;
}

.tree-toggle-inner-content {
  display: flex;
  align-items: center;
  height: 100%;
  padding: 12px 14px;
  box-sizing: border-box;
}

.tree-data-row {
  display: grid;
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
}

.list-cell--head {
  font-size: 12px;
  font-weight: 700;
  color: var(--list-head-text, #1f2937);
  background: var(--list-head-bg, #e5ebf2);
  border-bottom-color: var(--list-head-border, #aeb8c6);
}

.list-cell--tree {
  border-bottom: 1px solid var(--list-cell-border, #d6dbe3);
}

.is-left {
  justify-content: flex-start;
  text-align: left;
}

.is-center {
  justify-content: center;
  text-align: center;
}

.is-right {
  justify-content: flex-end;
  text-align: right;
}

.tree-toggle-scrollbar {
  width: var(--tree-toggle-width);
  overflow-x: auto;
  overflow-y: hidden;
  border-right: 1px solid var(--list-cell-border, var(--border-color));
  background: var(--list-cell-bg, var(--bg-white));
}

.tree-toggle-scrollbar-inner {
  height: 1px;
}

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
}

.tree-node-label {
  font-size: 14px;
  color: var(--text-medium);
}

</style>
