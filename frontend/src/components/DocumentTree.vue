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
        <DocumentTypeMenu
          v-if="showCreate"
          @open="closeContextMenu"
          @select="(type) => emitCreate(type, null)"
        >
          <el-button
            text
            class="tree-action-btn"
            :title="t('knowledgeBase.createDocument')"
            @click.stop
          >
            <el-icon :size="16">
              <Plus />
            </el-icon>
          </el-button>
        </DocumentTypeMenu>
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
              <DocumentTypeMenu
                class="node-action-dropdown"
                @open="closeContextMenu"
                @select="(type) => emitCreate(type, data)"
              >
                <button
                  type="button"
                  class="node-action-icon node-action-icon--hover"
                  :title="t('knowledgeBase.createDocument')"
                  @click.stop
                >
                  <el-icon :size="14">
                    <Plus />
                  </el-icon>
                </button>
              </DocumentTypeMenu>
              <button
                type="button"
                class="node-action-icon node-action-icon--hover"
                title="..."
                @click.stop="openContextMenuFromAction($event, data)"
              >
                <el-icon :size="14">
                  <MoreFilled />
                </el-icon>
              </button>
              <slot name="node-actions" :node="node" :data="data" />
            </div>
          </div>
        </template>
      </el-tree>
    </div>

    <AppMenu
      v-if="contextMenuEnabled"
      :visible="contextMenu.visible"
      :x="contextMenu.x"
      :y="contextMenu.y"
      @close="closeContextMenu"
    >
      <AppMenuItem
        v-if="showRename"
        @click="handleContextCommand('rename')"
      >
        <template #icon>
          <el-icon :size="14">
            <Edit />
          </el-icon>
        </template>
        {{ t('document.renameDocument') }}
      </AppMenuItem>
      <AppMenuItem
        v-if="showDelete"
        danger
        @click="handleContextCommand('delete')"
      >
        <template #icon>
          <el-icon :size="14">
            <Delete />
          </el-icon>
        </template>
        {{ t('document.deleteDocument') }}
      </AppMenuItem>
    </AppMenu>
  </div>
</template>

<script setup>
import { computed, ref } from 'vue';
import { useI18n } from 'vue-i18n';
import { Folder, FolderOpened, Document, Plus, Delete, Edit, MoreFilled } from '@element-plus/icons-vue';
import AppMenu from '@/components/AppMenu.vue';
import AppMenuItem from '@/components/AppMenuItem.vue';
import DocumentTypeMenu from '@/components/DocumentTypeMenu.vue';
import { useDocumentTreeContextMenu } from '@/composables/useDocumentTreeContextMenu';
import { useInlineRename } from '@/composables/useInlineRename';

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
const contextMenuEnabled = computed(() => props.contextMenuEnabled);
const {
  contextMenu,
  openContextMenu,
  closeContextMenu
} = useDocumentTreeContextMenu({
  enabled: contextMenuEnabled
});
const {
  setRenamingInputRef,
  startInlineRename,
  cancelInlineRename,
  confirmInlineRename
} = useInlineRename({
  revertRename: computed(() => props.revertRename),
  onConfirm: ({ node, title }) => {
    emit('rename', { data: node, title });
  }
});

const treeProps = {
  children: 'children',
  label: 'label',
  isLeaf: 'isLeaf'
};

const isSelected = (data) => String(data?.id) === String(props.currentId);

const handleNodeClick = (data) => {
  if (data?.isRenaming) {
    return;
  }
  emit('node-click', data);
};

const handleNodeContextMenu = (event, data) => {
  openContextMenu(event, data);
};

const emitCreate = (type, target) => {
  emit('create', {
    target: target || null,
    type
  });
  closeContextMenu();
};

const openContextMenuFromAction = (event, data) => {
  openContextMenu(event, data);
};

const handleContextCommand = (command) => {
  const target = contextMenu.value.data;
  if (!target) {
    return;
  }
  contextMenu.value.visible = false;
  if (command === 'delete') {
    emit('delete', target);
    return;
  }
  if (command === 'rename') {
    startInlineRename(target);
  }
};

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

.node-action-dropdown {
  opacity: 0;
  pointer-events: none;
  transition: opacity 0.18s ease;
}

.tree-node-content:hover .node-action-dropdown,
.tree-node-content.is-selected .node-action-dropdown {
  opacity: 1;
  pointer-events: auto;
}

.node-action-icon--hover {
  opacity: 0;
  pointer-events: none;
}

.tree-node-content:hover .node-action-icon--hover,
.tree-node-content.is-selected .node-action-icon--hover {
  opacity: 1;
  pointer-events: auto;
}

.node-action-icon {
  width: 22px;
  height: 22px;
  border: none;
  border-radius: 4px;
  background: transparent;
  color: var(--text-light);
  display: inline-flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.18s ease;
}

.node-action-icon:hover {
  background: color-mix(in srgb, var(--primary-color) 14%, transparent);
  color: var(--primary-color);
}

.dark-mode .node-action-icon {
  color: #9ca3af;
}

.dark-mode .node-action-icon:hover {
  color: #8fb2ff;
}

.inline-rename-input {
  max-width: 220px;
}
</style>
