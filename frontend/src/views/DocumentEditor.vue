<template>
  <PageLayout
    :system-name="systemName"
    :current-language="currentLanguage"
    :is-dark-mode="isDarkMode"
    :user-avatar="userAvatar"
    :username="userInfo?.username || '用户'"
    :active-menu="activeMenu"
    :title="''"
    content-padding="lg"
    @change-language="handleLanguageChange"
    @toggle-theme="toggleTheme"
    @logout="handleLogout"
    @menu-select="handleMenuSelect"
  >
    <div class="editor-layout">
      <DirectorySidebar
        ref="treeComponentRef"
        :collapsed="isSidebarCollapsed"
        :resizing="isSidebarResizing"
        :title="knowledgeBaseName"
        :collapse-title="t('common.collapse')"
        :back-label="t('common.back')"
        :nodes="directoryTree"
        :loading="treeLoading"
        :current-id="docId"
        :expand-all="isAllExpanded"
        :disable-delete="!docId"
        @back="goBack"
        @toggle="toggleSidebar"
        @resize-start="startResize"
        @toggle-expand="toggleExpandAll"
        @node-click="handleTreeNodeClick"
        @create="handleCreateFromTree"
        @delete="handleDeleteDocument"
        @rename="handleRenameFromTree"
      />

      <main class="editor-main">
        <DocumentEditorHeader
          :is-editing-title="isEditingTitle"
          :current-doc-title="currentDocTitle"
          :pending-title="pendingTitle"
          :is-sidebar-collapsed="isSidebarCollapsed"
          :title-placeholder="t('knowledgeBase.enterDocumentTitle')"
          :collapse-title="t('common.collapse')"
          :expand-title="t('common.expand')"
          :comments-title="t('document.comments')"
          :save-as-label="t('templates.saveAs')"
          @update:pending-title="pendingTitle = $event"
          @start-edit-title="startEditTitle"
          @commit-title="commitTitle"
          @cancel-edit-title="cancelEditTitle"
          @toggle-sidebar="toggleSidebar"
          @toggle-comment-sidebar="toggleCommentSidebar"
          @header-command="handleHeaderCommand"
        />
        <div class="editor-content">
          <div class="editor-wrapper">
            <div v-if="collabEnabled && !collabReady" class="editor-loading">
              {{ t('document.loading') }}
            </div>
            <MarkdownEditor
              ref="markdownEditorRef"
              v-model="editorContent"
              :placeholder="t('document.editorPlaceholder')"
              :collab-enabled="collabEnabled"
              :collab-room="collabRoom"
              :collab-url="collabUrl"
              :collab-token="collabToken"
              :comment-enabled="inlineCommentEnabled"
              @collab-sync="handleCollabSync"
              @comment-add="handleInlineCommentAdd"
              @comment-remove="handleInlineCommentRemove"
              @comment-changed="handleRemoteCommentChanged"
            />
          </div>
        </div>

      </main>
      <CommentSidebar
        ref="commentSidebarRef"
        :title="t('document.comments')"
        :collapse-title="t('common.collapse')"
        :collapsed="isCommentCollapsed"
        :doc-id="docId"
        :inline-enabled="inlineCommentEnabled"
        :user-info="userInfo"
        :on-anchor-click="scrollToInlineAnchor"
        :on-anchor-hover="handleAnchorHover"
        :on-anchor-remove="handleAnchorRemove"
        :on-comment-mutated="handleCommentMutated"
        @toggle="toggleCommentSidebar"
      />
    </div>
  </PageLayout>
  <DocumentCreateDialog v-model="showCreateDialog" :loading="creatingLoading"
    :parent-external-id="pendingParentId" :knowledge-base-id="kbId !== 'personal' ? kbId : ''"
    @submit="createDocument" @cancel="cancelCreateDocument" />
  <TemplateCreateDialog
    v-model="showTemplateDialog"
    :loading="savingTemplate"
    :title="t('templates.createDialogTitle')"
    :show-content="false"
    :show-kb-scope="kbId !== 'personal'"
    :initial-name="templateDialogInit.name"
    :initial-description="templateDialogInit.description"
    :initial-scope="templateDialogInit.scope"
    :initial-tags="templateDialogInit.tags"
    :initial-content="templateDialogInit.content"
    @submit="submitCreateTemplate"
  />
</template>

<script>
export default {
  inheritAttrs: false
};
</script>

<script setup>
import { ref, computed } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { useI18n } from 'vue-i18n';
import DirectorySidebar from '@/components/DirectorySidebar.vue';
import MarkdownEditor from '@/components/MarkdownEditor.vue';
import CommentSidebar from '@/components/CommentSidebar.vue';
import DocumentCreateDialog from '@/components/DocumentCreateDialog.vue';
import DocumentEditorHeader from '@/components/DocumentEditorHeader.vue';
import PageLayout from '@/components/PageLayout.vue';
import TemplateCreateDialog from '@/components/TemplateCreateDialog.vue';
import { usePanelSidebar } from '@/composables/usePanelSidebar';
import { useWorkspaceShell } from '@/composables/useWorkspaceShell';
import { useDocumentRouteContext } from '@/composables/useDocumentRouteContext';
import { useDirectoryTreeState } from '@/composables/useDirectoryTreeState';
import { useEditorCommentBridge } from '@/composables/useEditorCommentBridge';
import { useDocumentEditorActions } from '@/composables/useDocumentEditorActions';
import { useDocumentEditorLifecycle } from '@/composables/useDocumentEditorLifecycle';
import { useUserStore } from '@/store/user';
import {
  getKnowledgeBaseDocuments,
  getMyDocuments,
  recordRecentDocument
} from '@/services/api';

const props = defineProps({
  kbId: {
    type: String,
    default: ''
  },
  docId: {
    type: String,
    default: ''
  }
});

const { t, locale } = useI18n();
const route = useRoute();
const router = useRouter();
const userStore = useUserStore();
const {
  kbId,
  docId,
  resolveActiveMenu,
  collabEnabled,
  collabRoom,
  collabUrl,
  collabToken,
  collabReady,
  lastSyncedDocId
} = useDocumentRouteContext({ props, route });

const {
  systemName,
  activeMenu,
  isDarkMode,
  userInfo,
  userAvatar,
  currentLanguage,
  initLanguage,
  handleLanguageChange,
  toggleTheme,
  handleLogout,
  handleMenuSelect,
  fetchSystemInfo
} = useWorkspaceShell({
  locale,
  router,
  userStore,
  defaultActiveMenu: resolveActiveMenu(kbId.value)
});

const treeComponentRef = ref(null);
const {
  treeLoading,
  directoryTree,
  knowledgeBaseName,
  currentDocTitle,
  isAllExpanded,
  fetchDocuments,
  updateCurrentDocTitle,
  updateTreeNodeTitle,
  expandToCurrentDoc,
  goBack,
  handleTreeNodeClick,
  toggleExpandAll
} = useDirectoryTreeState({
  t,
  router,
  kbId,
  docId,
  treeComponentRef,
  getKnowledgeBaseDocuments,
  getMyDocuments
});

const editorContent = ref('');
const {
  collapsed: isSidebarCollapsed,
  resizing: isSidebarResizing,
  toggleCollapsed: toggleSidebar,
  startResize
} = usePanelSidebar({
  defaultWidth: 280,
  minWidth: 220,
  maxWidth: 520,
  resizeEdge: 'right',
  collapsedStorageKey: 'sidebarCollapsed',
  widthStorageKey: 'docSidebarWidth',
  getMaxWidth: () => {
    const layoutRect = document.querySelector('.editor-layout')?.getBoundingClientRect();
    if (!layoutRect) return 520;
    return Math.min(520, layoutRect.width - 320);
  },
  onWidthChange: (value) => {
    document.documentElement.style.setProperty('--sidebar-width', `${value}px`);
  }
});
const isCommentCollapsed = ref(false);
const markdownEditorRef = ref(null);
const commentSidebarRef = ref(null);
const inlineCommentEnabled = computed(() => !!docId.value && docId.value !== 'example');
const {
  isEditingTitle,
  pendingTitle,
  showCreateDialog,
  creatingLoading,
  pendingParentId,
  savingTemplate,
  showTemplateDialog,
  templateDialogInit,
  cancelCreateDocument,
  createDocument,
  handleCreateFromTree,
  handleDeleteDocument,
  handleRenameFromTree,
  startEditTitle,
  cancelEditTitle,
  commitTitle,
  handleHeaderCommand,
  submitCreateTemplate
} = useDocumentEditorActions({
  t,
  router,
  kbId,
  docId,
  currentDocTitle,
  editorContent,
  directoryTree,
  updateTreeNodeTitle,
  fetchDocuments
});
const {
  handleInlineCommentAdd,
  handleInlineCommentRemove,
  handleAnchorHover,
  handleAnchorRemove,
  handleCommentMutated,
  handleRemoteCommentChanged,
  scrollToInlineAnchor
} = useEditorCommentBridge({
  isCommentCollapsed,
  markdownEditorRef,
  commentSidebarRef
});
const {
  toggleCommentSidebar,
  handleCollabSync
} = useDocumentEditorLifecycle({
  props,
  route,
  initLanguage,
  fetchSystemInfo,
  kbId,
  docId,
  activeMenu,
  resolveActiveMenu,
  collabEnabled,
  collabReady,
  lastSyncedDocId,
  editorContent,
  currentDocTitle,
  knowledgeBaseName,
  fetchDocuments,
  updateCurrentDocTitle,
  expandToCurrentDoc,
  commentSidebarRef,
  isCommentCollapsed,
  cancelEditTitle,
  recordRecentDocument
});
</script>

<style scoped>
.editor-layout {
  display: flex;
  height: 100%;
  background-color: var(--bg-white);
  border-radius: var(--border-radius-md);
  box-shadow: var(--shadow-sm);
  overflow: hidden;
  transition: all 0.3s ease-in-out;
}

.editor-main {
  flex: 1;
  display: flex;
  flex-direction: column;
  min-height: 600px;
  background-color: var(--bg-white);
  transition: all 0.3s ease-in-out;
}

.dark-mode .editor-main {
  background-color: var(--bg-white);
}

.editor-content {
  flex: 1;
  display: flex;
  overflow: hidden;
  min-height: 0;
  transition: all 0.3s ease-in-out;
}

.editor-wrapper {
  flex: 1;
  display: flex;
  flex-direction: column;
  min-height: 500px;
  position: relative;
}
.editor-loading {
  position: absolute;
  inset: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(255, 255, 255, 0.85);
  color: var(--text-medium);
  font-size: 14px;
  z-index: 2;
}

.dark-mode .editor-loading {
  background: rgba(255, 255, 255, 0.9);
}

.editor-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: var(--spacing-sm) var(--spacing-lg);
  border-top: 1px solid var(--border-color);
  font-size: 12px;
  color: var(--text-light);
}

.dark-mode .editor-footer {
  border-color: var(--border-color);
  color: var(--text-light);
}

/* 深色模式对话框样式 */
.dark-mode .el-dialog {
  background-color: var(--bg-white);
  border: 1px solid var(--border-color);
  color: var(--text-dark);
}

.dark-mode .el-dialog__header {
  background-color: var(--bg-white);
  border-bottom: 1px solid var(--border-color);
  color: var(--text-dark);
}

.dark-mode .el-dialog__body {
  background-color: var(--bg-white);
  color: var(--text-dark);
}

.dark-mode .el-dialog__footer {
  background-color: var(--bg-white);
  border-top: 1px solid var(--border-color);
}

.dark-mode .el-form-item__label {
  color: var(--text-dark);
}

.dark-mode :deep(.el-input__wrapper) {
  background-color: var(--input-bg);
  border-color: var(--input-border);
  color: var(--input-text);
}

.dark-mode :deep(.el-input__inner) {
  background-color: var(--input-bg);
  border-color: var(--input-border);
  color: var(--input-text);
}

.dark-mode :deep(.el-select__wrapper) {
  background-color: var(--select-bg);
  border-color: var(--select-border);
  color: var(--select-text);
}

.dark-mode :deep(.el-select__input) {
  background-color: var(--select-bg);
  color: var(--select-text);
}

.dark-mode :deep(.el-select-dropdown__item) {
  background-color: var(--select-option-bg);
  color: var(--select-text);
}
</style>
