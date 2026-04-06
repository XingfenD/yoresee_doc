<template>
  <PageLayout
    :system-name="systemName"
    :current-language="currentLanguage"
    :is-dark-mode="isDarkMode"
    :user-avatar="userAvatar"
    :username="userInfo?.username || t('common.user')"
    :active-menu="activeMenu"
    :title="''"
    content-padding="lg"
    @change-language="handleLanguageChange"
    @toggle-theme="toggleTheme"
    @logout="handleLogout"
    @menu-select="handleMenuSelect"
  >
    <div ref="editorLayoutRef" class="editor-layout" :class="{ 'is-fullscreen': isEditorFullscreen }">
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
          :history-label="t('document.history')"
          :save-as-label="t('templates.saveAs')"
          :attachments-label="t('document.attachments.title')"
          :settings-label="t('document.settings.title')"
          :is-fullscreen="isEditorFullscreen"
          :can-manage-attachments="canManageAttachments"
          :can-manage-settings="canManageSettings"
          @update:pending-title="pendingTitle = $event"
          @start-edit-title="startEditTitle"
          @commit-title="commitTitle"
          @cancel-edit-title="cancelEditTitle"
          @toggle-sidebar="toggleSidebar"
          @toggle-comment-sidebar="toggleCommentSidebar"
          @header-command="onHeaderCommand"
        >
          <template #title-action>
            <el-button
              v-if="canToggleEditorFullscreen"
              class="editor-fullscreen-header-btn"
              text
              :title="isEditorFullscreen ? t('common.exitFullscreen') : t('common.fullscreen')"
              @click="toggleEditorFullscreen"
            >
              <el-icon>
                <ScaleToOriginal v-if="isEditorFullscreen" />
                <FullScreen v-else />
              </el-icon>
            </el-button>
          </template>
        </DocumentEditorHeader>
        <div class="editor-content">
          <div class="editor-wrapper">
            <div v-if="collabEnabled && !collabReady" class="editor-loading">
              {{ t('document.loading') }}
            </div>
            <MarkdownEditor
              v-if="isMarkdownDocument"
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
            <TableEditor
              v-else-if="isTableDocument"
              ref="tableEditorRef"
              v-model="editorContent"
              @commit="flushTableSave"
            />
            <SlideEditor
              v-else-if="isSlideDocument"
              ref="slideEditorRef"
              v-model="editorContent"
              @commit="flushSlideSave"
            />
            <YoreseeRichTextEditor
              v-else-if="isRichTextDocument"
              ref="richTextEditorRef"
              v-model="editorContent"
              :placeholder="t('document.editorPlaceholder')"
              :comment-enabled="inlineCommentEnabled"
              @commit="flushRichTextSave"
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
        @width-change="handleCommentWidthChange"
        @toggle="toggleCommentSidebar"
      />
    </div>
  </PageLayout>
  <DocumentCreateDialog v-model="showCreateDialog" :loading="creatingLoading"
    :parent-external-id="pendingParentId" :initial-document-type="selectedDocumentType"
    :knowledge-base-id="createDialogKnowledgeBaseId"
    @submit="createDocument" @cancel="cancelCreateDocument" />
  <TemplateCreateDialog
    v-model="showTemplateDialog"
    :loading="savingTemplate"
    :title="t('templates.createDialogTitle')"
    :show-content="false"
    :show-kb-scope="showTemplateDialogKbScope"
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
import { defineAsyncComponent, ref } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { useI18n } from 'vue-i18n';
import { FullScreen, ScaleToOriginal } from '@element-plus/icons-vue';
import DirectorySidebar from '@/components/document/DirectorySidebar.vue';
import MarkdownEditor from '@/components/document/MarkdownEditor.vue';
import TableEditor from '@/components/document/TableEditor.vue';
import YoreseeRichTextEditor from '@/components/document/YoreseeRichTextEditor.vue';
import CommentSidebar from '@/components/comment/CommentSidebar.vue';
import DocumentCreateDialog from '@/components/document/DocumentCreateDialog.vue';
import DocumentEditorHeader from '@/components/document/DocumentEditorHeader.vue';
import PageLayout from '@/components/layout/PageLayout.vue';
import TemplateCreateDialog from '@/components/template/TemplateCreateDialog.vue';
import { useWorkspaceShell } from '@/composables/shell/useWorkspaceShell';
import { useDocumentRouteContext } from '@/composables/document/editor/useDocumentRouteContext';
import { useDirectoryTreeState } from '@/composables/document/tree/useDirectoryTreeState';
import { useEditorCommentBridge } from '@/composables/document/editor/useEditorCommentBridge';
import { useDocumentEditorActions } from '@/composables/document/editor/useDocumentEditorActions';
import { useDocumentEditorLifecycle } from '@/composables/document/editor/useDocumentEditorLifecycle';
import { useEditorFullscreen } from '@/composables/document/editor/useEditorFullscreen';
import { useEditorPanelConstraints } from '@/composables/document/editor/useEditorPanelConstraints';
import { useDocumentEditorPolicy } from '@/composables/document/editor/useDocumentEditorPolicy';
import { useDocumentHeaderRouting } from '@/composables/document/editor/useDocumentHeaderRouting';
import { useTableDocumentPersistence } from '@/composables/document/editor/table-editor/useTableDocumentPersistence';
import { useSlideDocumentPersistence } from '@/composables/document/editor/slide-editor/useSlideDocumentPersistence';
import { useRichTextDocumentPersistence } from '@/composables/document/editor/rich-text-editor/useRichTextDocumentPersistence';
import { useUserStore } from '@/store/user';
import {
  getKnowledgeBaseDocuments,
  getMyDocuments,
  getDocumentContent,
  updateDocument,
  recordRecentDocument
} from '@/services/api';

const SlideEditor = defineAsyncComponent(() => import('@/components/document/SlideEditor.vue'));

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
  currentDocType,
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
const editorLayoutRef = ref(null);
const isCommentCollapsed = ref(false);
const markdownEditorRef = ref(null);
const tableEditorRef = ref(null);
const slideEditorRef = ref(null);
const richTextEditorRef = ref(null);
const commentSidebarRef = ref(null);
const {
  flushTableSave,
  rerenderTableEditor
} = useTableDocumentPersistence({
  docId,
  currentDocType,
  editorContent,
  tableEditorRef,
  t,
  getDocumentContent,
  updateDocument
});
const {
  flushSlideSave,
  rerenderSlideEditor
} = useSlideDocumentPersistence({
  docId,
  currentDocType,
  editorContent,
  slideEditorRef,
  t,
  getDocumentContent,
  updateDocument
});
const {
  flushRichTextSave,
  rerenderRichTextEditor
} = useRichTextDocumentPersistence({
  docId,
  currentDocType,
  editorContent,
  richTextEditorRef,
  t,
  getDocumentContent,
  updateDocument
});
const {
  isSidebarCollapsed,
  isSidebarResizing,
  toggleSidebar,
  startResize,
  handleCommentWidthChange,
  clampSidebarWidth
} = useEditorPanelConstraints({
  editorLayoutRef,
  commentSidebarRef,
  isCommentCollapsed,
  onLayoutChange: () => {
    rerenderTableEditor();
    rerenderSlideEditor();
    rerenderRichTextEditor();
  }
});
const {
  isMarkdownDocument,
  isTableDocument,
  isSlideDocument,
  isRichTextDocument,
  canManageAttachments,
  canManageSettings,
  collabEnabled,
  inlineCommentEnabled,
  createDialogKnowledgeBaseId,
  showTemplateDialogKbScope
} = useDocumentEditorPolicy({
  kbId,
  docId,
  currentDocType
});
const {
  isEditingTitle,
  pendingTitle,
  showCreateDialog,
  creatingLoading,
  pendingParentId,
  selectedDocumentType,
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
  currentDocType,
  currentDocTitle,
  editorContent,
  directoryTree,
  updateTreeNodeTitle,
  fetchDocuments
});

const { handleHeaderCommand: onHeaderCommand } = useDocumentHeaderRouting({
  router,
  kbId,
  docId,
  onCommand: handleHeaderCommand
});

const {
  isEditorFullscreen,
  canToggleEditorFullscreen,
  toggleEditorFullscreen
} = useEditorFullscreen({
  editorLayoutRef,
  docId,
  onChange: () => {
    rerenderTableEditor();
    rerenderSlideEditor();
    rerenderRichTextEditor();
    clampSidebarWidth();
  }
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
  richTextEditorRef,
  isMarkdownDocument,
  isRichTextDocument,
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
  width: 100%;
  max-width: 100%;
  min-width: 0;
  height: 100%;
  min-height: 0;
  background-color: var(--bg-white);
  border-radius: var(--border-radius-md);
  box-shadow: var(--shadow-sm);
  overflow: hidden;
  transition: all 0.3s ease-in-out;
}

.editor-layout.is-fullscreen {
  width: 100vw;
  height: 100vh;
  border-radius: 0;
  box-shadow: none;
}

.editor-main {
  flex: 1;
  display: flex;
  flex-direction: column;
  min-width: 0;
  min-height: 0;
  overflow: hidden;
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
  min-width: 0;
  min-height: 0;
  transition: all 0.3s ease-in-out;
}

.editor-wrapper {
  flex: 1;
  display: flex;
  flex-direction: column;
  min-width: 0;
  min-height: 0;
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
  background: rgba(11, 17, 26, 0.72);
  color: var(--text-dark);
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

.editor-fullscreen-header-btn {
  color: var(--text-medium);
}

.editor-fullscreen-header-btn:hover {
  color: var(--primary-color);
}
</style>
