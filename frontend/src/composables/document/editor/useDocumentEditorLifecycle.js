import { onMounted, watch } from 'vue';

export function useDocumentEditorLifecycle({
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
}) {
  const toggleCommentSidebar = () => {
    isCommentCollapsed.value = !isCommentCollapsed.value;
  };

  const handleCollabSync = (isSynced) => {
    if (!collabEnabled.value) {
      collabReady.value = true;
      return;
    }
    collabReady.value = isSynced;
    if (isSynced) {
      lastSyncedDocId.value = docId.value || '';
    }
  };

  onMounted(async () => {
    initLanguage();
    activeMenu.value = resolveActiveMenu(kbId.value);

    if (kbId.value === 'example' && docId.value === 'example') {
      knowledgeBaseName.value = '示例知识库';
      currentDocTitle.value = '示例文档';
    } else {
      await fetchDocuments();
      if (docId.value) {
        recordRecentDocument(docId.value).catch(() => {});
      }
      if (lastSyncedDocId.value !== docId.value) {
        collabReady.value = !collabEnabled.value;
      }
    }

    await fetchSystemInfo();
  });

  watch(
    () => props.docId || route.params.docId,
    async (newDocId) => {
      docId.value = newDocId;
      editorContent.value = '';
      currentDocTitle.value = '';
      cancelEditTitle();
      await commentSidebarRef.value?.reload?.();
      if (docId.value && docId.value !== 'example') {
        recordRecentDocument(docId.value).catch(() => {});
      }
      await expandToCurrentDoc();
      updateCurrentDocTitle();
      if (lastSyncedDocId.value !== docId.value) {
        collabReady.value = !collabEnabled.value;
      }
    }
  );

  watch(
    () => props.kbId || route.params.kbId,
    async (newKbId) => {
      if (!newKbId) {
        return;
      }
      kbId.value = newKbId;
      activeMenu.value = resolveActiveMenu(kbId.value);
      cancelEditTitle();
      await fetchDocuments();
      if (lastSyncedDocId.value !== docId.value) {
        collabReady.value = !collabEnabled.value;
      }
      updateCurrentDocTitle();
    }
  );

  return {
    toggleCommentSidebar,
    handleCollabSync
  };
}
