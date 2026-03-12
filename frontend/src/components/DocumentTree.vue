<template>
  <div class="document-tree">
    <div v-if="showToolbar" class="tree-toolbar">
      <div class="tree-toolbar-left">
        <el-button text @click="emit('toggle-expand')"
          :title="expandAll ? t('common.collapseAll') : t('common.expandAll')" class="tree-expand-btn">
          <el-icon :size="14">
            <FolderOpened v-if="expandAll" />
            <Folder v-else />
          </el-icon>
        </el-button>
      </div>
      <div class="tree-toolbar-actions">
        <el-button v-if="showCreate" text class="tree-action-btn" :title="t('knowledgeBase.createDocument')"
          @click="emit('create', null)">
          <el-icon :size="16">
            <Plus />
          </el-icon>
        </el-button>
        <el-button v-if="showDelete" text class="tree-action-btn tree-action-btn--danger" :disabled="disableDelete"
          :title="t('document.deleteDocument')" @click="emit('delete', null)">
          <el-icon :size="16">
            <Delete />
          </el-icon>
        </el-button>
      </div>
    </div>

    <div class="directory-tree" v-loading="loading" @click="closeContextMenu">
      <el-tree ref="treeRef" :data="nodes" :props="treeProps" node-key="id" :default-expand-all="false"
        :expand-on-click-node="false"
        @node-click="handleNodeClick" class="editor-tree">
        <template #default="{ node, data }">
          <div class="tree-node-content" :class="{ 'is-selected': isSelected(data) }" @click="closeContextMenu"
            @contextmenu.prevent="(event) => handleNodeContextMenu(event, data)">
            <div class="node-icon">
              <el-icon v-if="data.isFolder">
                <FolderOpened v-if="node.expanded" />
                <Folder v-else />
              </el-icon>
              <el-icon v-else>
                <Document />
              </el-icon>
            </div>
            <div class="node-info">
              <el-input
                v-if="data.isRenaming"
                :ref="(el) => setRenamingInputRef(el, data)"
                v-model="data.renameValue"
                size="small"
                class="inline-rename-input"
                @keyup.enter="confirmInlineRename(data)"
                @keydown.esc.prevent="cancelInlineRename(data)"
                @blur="confirmInlineRename(data)"
              />
              <span v-else class="node-label">{{ node.label }}</span>
              <slot name="node-extra" :node="node" :data="data" />
            </div>
            <div class="node-actions">
              <slot name="node-actions" :node="node" :data="data" />
            </div>
          </div>
        </template>
      </el-tree>
    </div>

    <div v-if="contextMenuEnabled" v-show="contextMenu.visible" class="tree-context-menu" :style="contextMenuStyle">
      <button class="context-item" type="button" @click="handleContextCommand('create')">
        <el-icon :size="14" class="context-icon">
          <Plus />
        </el-icon>
        {{ t('document.createDocument') }}
      </button>
      <button v-if="showRename" class="context-item" type="button" @click="handleContextCommand('rename')">
        <el-icon :size="14" class="context-icon">
          <Edit />
        </el-icon>
        {{ t('document.renameDocument') }}
      </button>
      <button v-if="showDelete" class="context-item is-danger" type="button" @click="handleContextCommand('delete')">
        <el-icon :size="14" class="context-icon">
          <Delete />
        </el-icon>
        {{ t('document.deleteDocument') }}
      </button>
    </div>
  </div>
</template>

<script setup>
import { computed, ref, nextTick, onMounted, onBeforeUnmount } from 'vue';
import { useI18n } from 'vue-i18n';
import { Folder, FolderOpened, Document, Plus, Delete, Edit } from '@element-plus/icons-vue';

const props = defineProps({
  nodes: {
    type: Array,
    default: () => []
  },
  loading: {
    type: Boolean,
    default: false
  },
  currentId: {
    type: [String, Number],
    default: ''
  },
  expandAll: {
    type: Boolean,
    default: false
  },
  showToolbar: {
    type: Boolean,
    default: true
  },
  showCreate: {
    type: Boolean,
    default: true
  },
  showDelete: {
    type: Boolean,
    default: true
  },
  disableDelete: {
    type: Boolean,
    default: false
  },
  contextMenuEnabled: {
    type: Boolean,
    default: true
  },
  showRename: {
    type: Boolean,
    default: true
  },
  revertRename: {
    type: Boolean,
    default: true
  }
});

const emit = defineEmits(['node-click', 'toggle-expand', 'create', 'delete', 'rename']);
const { t } = useI18n();

const treeRef = ref(null);
const contextMenu = ref({
  visible: false,
  x: 0,
  y: 0,
  data: null
});
const renamingInputRef = ref(null);

const treeProps = {
  children: 'children',
  label: 'label',
  isLeaf: 'isLeaf'
};

const isSelected = (data) => String(data?.id) === String(props.currentId);

const setRenamingInputRef = (el, data) => {
  if (el && data?.isRenaming) {
    renamingInputRef.value = el;
  }
};

const handleNodeClick = (data) => {
  if (data?.isRenaming) {
    return;
  }
  emit('node-click', data);
};

const handleNodeContextMenu = (event, data) => {
  if (!props.contextMenuEnabled) {
    return;
  }
  event.preventDefault();
  event.stopPropagation();
  const menuWidth = 150;
  const menuHeight = 120;
  let x = event.clientX;
  let y = event.clientY;
  if (x + menuWidth > window.innerWidth) {
    x = window.innerWidth - menuWidth - 8;
  }
  if (y + menuHeight > window.innerHeight) {
    y = window.innerHeight - menuHeight - 8;
  }
  contextMenu.value = {
    visible: true,
    x,
    y,
    data
  };
};

const closeContextMenu = () => {
  if (contextMenu.value.visible) {
    contextMenu.value.visible = false;
  }
};

const handleContextCommand = (command) => {
  const target = contextMenu.value.data;
  if (!target) {
    return;
  }
  contextMenu.value.visible = false;
  if (command === 'create') {
    emit('create', target);
    return;
  }
  if (command === 'delete') {
    emit('delete', target);
    return;
  }
  if (command === 'rename') {
    startInlineRename(target);
  }
};

const startInlineRename = (node) => {
  if (!node || node.isRenaming) {
    return;
  }
  node.isRenaming = true;
  node.originalLabel = node.label;
  node.renameValue = node.label;
  nextTick(() => {
    const inputEl = renamingInputRef.value?.input || renamingInputRef.value;
    if (inputEl && typeof inputEl.focus === 'function') {
      inputEl.focus();
      if (typeof inputEl.select === 'function') {
        inputEl.select();
      }
    }
  });
};

const cancelInlineRename = (node) => {
  if (!node?.isRenaming) {
    return;
  }
  node.label = node.originalLabel || node.label;
  node.renameValue = '';
  node.isRenaming = false;
};

const confirmInlineRename = (node) => {
  if (!node?.isRenaming) {
    return;
  }
  const nextName = node.renameValue?.trim();
  if (!nextName) {
    cancelInlineRename(node);
    return;
  }
  node.isRenaming = false;
  node.renameValue = '';
  emit('rename', { data: node, title: nextName });
  if (props.revertRename) {
    node.label = node.originalLabel || node.label;
  } else {
    node.label = nextName;
  }
};

const contextMenuStyle = computed(() => ({
  left: `${contextMenu.value.x}px`,
  top: `${contextMenu.value.y}px`
}));

onMounted(() => {
  window.addEventListener('click', closeContextMenu);
  window.addEventListener('scroll', closeContextMenu, true);
});

onBeforeUnmount(() => {
  window.removeEventListener('click', closeContextMenu);
  window.removeEventListener('scroll', closeContextMenu, true);
});

defineExpose({ treeRef, closeContextMenu });
</script>

<style scoped>
.document-tree {
  display: flex;
  flex-direction: column;
  height: 100%;
}

.tree-toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: var(--spacing-sm);
  padding: var(--spacing-xs) var(--spacing-md);
  border-bottom: 1px solid var(--border-color);
}

.dark-mode .tree-toolbar {
  border-color: var(--border-color);
}

.tree-toolbar-actions {
  display: flex;
  align-items: center;
  gap: var(--spacing-xs);
}

.tree-action-btn {
  padding: 2px 4px;
  color: var(--text-light);
}

.tree-action-btn:hover {
  color: var(--primary-color);
}

.tree-action-btn--danger:hover {
  color: #f56c6c;
}

.directory-tree {
  flex: 1;
  overflow-y: auto;
  overflow-x: auto;
  padding: var(--spacing-sm);
}

.editor-tree {
  background: transparent;
}

.editor-tree :deep(.el-tree-node__label) {
  white-space: nowrap;
}

.editor-tree :deep(.el-tree-node__content) {
  min-width: max-content;
}

.editor-tree :deep(.el-tree-node__children) {
  min-width: max-content;
}

.editor-tree :deep(.el-tree-node__content:hover) {
  background-color: var(--bg-light);
}

.dark-mode .editor-tree :deep(.el-tree-node__content:hover) {
  background-color: rgba(255, 255, 255, 0.08);
}

.editor-tree :deep(.el-tree-node.is-current > .el-tree-node__content),
.editor-tree :deep(.el-tree-node.is-focus > .el-tree-node__content),
.editor-tree :deep(.el-tree-node:focus-within > .el-tree-node__content),
.editor-tree :deep(.el-tree-node__content:focus),
.editor-tree :deep(.el-tree-node__content:focus-visible) {
  background-color: transparent;
  outline: none;
}

.dark-mode .editor-tree :deep(.el-tree-node.is-current > .el-tree-node__content),
.dark-mode .editor-tree :deep(.el-tree-node.is-focus > .el-tree-node__content),
.dark-mode .editor-tree :deep(.el-tree-node:focus-within > .el-tree-node__content),
.dark-mode .editor-tree :deep(.el-tree-node__content:focus),
.dark-mode .editor-tree :deep(.el-tree-node__content:focus-visible) {
  background-color: transparent;
  outline: none;
}

.tree-node-content.is-selected {
  background-color: rgba(22, 93, 255, 0.1);
}

.dark-mode .tree-node-content.is-selected {
  background-color: rgba(64, 128, 255, 0.35);
}

.tree-node-content {
  display: flex;
  align-items: center;
  gap: var(--spacing-xs);
  padding: var(--spacing-xs) var(--spacing-sm);
  border-radius: var(--border-radius-sm);
  width: 100%;
  box-sizing: border-box;
}

.node-icon {
  display: flex;
  align-items: center;
  color: var(--text-light);
}

.node-info {
  display: flex;
  align-items: center;
  gap: var(--spacing-xs);
  flex: 1;
}

.node-label {
  font-size: 14px;
  color: var(--text-medium);
  cursor: pointer;
}

.node-label:hover {
  color: var(--primary-color);
}

.node-actions {
  display: flex;
  align-items: center;
  gap: var(--spacing-xs);
}

.tree-context-menu {
  position: fixed;
  z-index: 3000;
  min-width: 150px;
  background-color: var(--bg-white);
  border: 1px solid var(--border-color);
  border-radius: var(--border-radius-sm);
  box-shadow: var(--shadow-md);
  padding: var(--spacing-xs) 0;
}

.dark-mode .tree-context-menu {
  background-color: var(--bg-white);
  border-color: var(--border-color);
}

.context-item {
  width: 100%;
  padding: var(--spacing-xs) var(--spacing-md);
  background: transparent;
  border: none;
  text-align: left;
  display: flex;
  align-items: center;
  gap: var(--spacing-xs);
  color: var(--text-medium);
  cursor: pointer;
}

.context-item:hover {
  background-color: var(--bg-light);
  color: var(--primary-color);
}

.dark-mode .context-item:hover {
  background-color: rgba(255, 255, 255, 0.08);
}

.context-item.is-danger:hover {
  color: #f56c6c;
}

.context-icon {
  color: var(--text-light);
}

.context-item:hover .context-icon {
  color: var(--primary-color);
}

.context-item.is-danger:hover .context-icon {
  color: #f56c6c;
}

.inline-rename-input {
  max-width: 220px;
}
</style>
