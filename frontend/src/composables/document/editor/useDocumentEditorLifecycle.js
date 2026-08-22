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
  markdownContent,
  tableContent,
  slideContent,
  richTextContent,
  currentDocTitle,
  knowledgeBaseName,
  fetchDocuments,
  updateCurrentDocTitle,
  expandToCurrentDoc,
  commentSidebarRef,
  isCommentCollapsed,
  cancelEditTitle,
  recordRecentDocument,
  docMachine
}) {
  const clearEditorContents = () => {
    if (markdownContent) markdownContent.value = '';
    if (tableContent) tableContent.value = '';
    if (slideContent) slideContent.value = '';
    if (richTextContent) richTextContent.value = '';
  };

  const toggleCommentSidebar = () => {
    isCommentCollapsed.value = !isCommentCollapsed.value;
  };

  const handleCollabSync = (isSynced) => {
    docMachine.markCollabSynced(collabEnabled.value ? isSynced : true, docId.value);
  };

  const syncCollabReadyFlag = () => {
    docMachine.markCollabSynced(!collabEnabled.value, docId.value);
  };

  const resolveTargetDocument = async (targetDocId) => {
    if (!targetDocId) {
      docMachine.startNavigation('');
      docId.value = '';
      clearEditorContents();
      currentDocTitle.value = '';
      cancelEditTitle();
      docMachine.resolveContent({ docId: '', content: '' });
      return;
    }

    docMachine.startNavigation(targetDocId);
    docId.value = targetDocId;
    clearEditorContents();
    currentDocTitle.value = '';
    cancelEditTitle();

    await commentSidebarRef.value?.reload?.();
    if (docId.value) {
      recordRecentDocument(docId.value).catch(() => {});
    }

    await expandToCurrentDoc();
    const { type } = updateCurrentDocTitle();

    docMachine.resolveMetadata({ docId: targetDocId, type });

    // For table/slide the typed persistence will load content once the machine
    // reaches READY. For markdown/rich-text the collaboration layer is the
    // source of truth, so we mount with empty local content.
    docMachine.resolveContent({ docId: targetDocId, content: '' });

    syncCollabReadyFlag();
  };

  onMounted(async () => {
    initLanguage();
    activeMenu.value = resolveActiveMenu(kbId.value);

    await fetchDocuments();
    if (docId.value) {
      await resolveTargetDocument(docId.value);
    } else {
      docMachine.startNavigation('');
      docMachine.resolveContent({ docId: '', content: '' });
    }

    await fetchSystemInfo();
  });

  watch(
    () => props.docId || route.params.docId,
    resolveTargetDocument
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
      updateCurrentDocTitle();
    }
  );

  return {
    toggleCommentSidebar,
    handleCollabSync
  };
}
