<template>
  <PanelSidebarShell
    :collapsed="collapsed"
    :resizing="resizing"
    resize-edge="right"
    :panel-style="panelStyle"
    @resize-start="$emit('resize-start', $event)"
  >
    <div class="directory-sidebar">
      <div class="sidebar-header">
        <el-button text class="back-button" @click="$emit('back')">
          <el-icon>
            <ArrowLeft />
          </el-icon>
          {{ backLabel }}
        </el-button>
      </div>
      <div class="sidebar-title">
        {{ title }}
        <el-button text class="collapse-button" @click="$emit('toggle')" :title="collapseTitle">
          <el-icon>
            <ArrowLeft />
          </el-icon>
        </el-button>
      </div>
      <DocumentTree
        ref="treeComponentRef"
        class="directory-tree-panel"
        :nodes="nodes"
        :loading="loading"
        :current-id="currentId"
        :expand-all="expandAll"
        :disable-delete="disableDelete"
        @toggle-expand="$emit('toggle-expand')"
        @node-click="(data) => $emit('node-click', data)"
        @create="(target) => $emit('create', target)"
        @delete="(target) => $emit('delete', target)"
        @rename="(target) => $emit('rename', target)"
      />
    </div>
  </PanelSidebarShell>
</template>

<script setup>
import { ref } from 'vue';
import { ArrowLeft } from '@element-plus/icons-vue';
import DocumentTree from '@/components/document/DocumentTree.vue';
import PanelSidebarShell from '@/components/layout/PanelSidebarShell.vue';

defineProps({
  collapsed: { type: Boolean, default: false },
  resizing: { type: Boolean, default: false },
  title: { type: String, default: '' },
  collapseTitle: { type: String, default: '' },
  backLabel: { type: String, default: '' },
  nodes: { type: Array, default: () => [] },
  loading: { type: Boolean, default: false },
  currentId: { type: [String, Number], default: '' },
  expandAll: { type: Boolean, default: false },
  disableDelete: { type: Boolean, default: false }
});

defineEmits(['back', 'toggle', 'toggle-expand', 'node-click', 'create', 'delete', 'rename', 'resize-start']);

const treeComponentRef = ref(null);
const panelStyle = {
  width: 'var(--sidebar-width)',
  maxWidth: '520px'
};

defineExpose({
  getTreeRef: () => treeComponentRef.value.treeRef,
  closeContextMenu: () => treeComponentRef.value?.closeContextMenu?.()
});
</script>

<style scoped>
.directory-sidebar {
  display: flex;
  flex-direction: column;
  height: 100%;
  min-height: 0;
  overflow: hidden;
}

.directory-tree-panel {
  flex: 1;
  min-height: 0;
}

.sidebar-header {
  padding: var(--spacing-md);
  border-bottom: 1px solid var(--border-color);
}

:global(.dark-mode) .sidebar-header {
  border-color: var(--border-color);
}

.sidebar-title {
  padding: var(--spacing-md);
  font-size: 16px;
  font-weight: 600;
  color: var(--text-dark);
  border-bottom: 1px solid var(--border-color);
  display: flex;
  justify-content: space-between;
  align-items: center;
}

:global(.dark-mode) .sidebar-title {
  color: var(--text-dark);
  border-color: var(--border-color);
}

.collapse-button {
  padding: 4px;
  color: var(--text-light);
}

.collapse-button:hover {
  color: var(--primary-color);
}

</style>
