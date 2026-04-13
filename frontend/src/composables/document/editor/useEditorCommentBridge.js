import { onBeforeUnmount, ref } from 'vue';

export function useEditorCommentBridge({
  isCommentCollapsed,
  markdownEditorRef,
  richTextEditorRef,
  isMarkdownDocument,
  isRichTextDocument,
  commentSidebarRef
}) {
  const remoteCommentReloadTimer = ref(null);

  const editorCandidates = [
    {
      active: isRichTextDocument,
      ref: richTextEditorRef
    },
    {
      active: isMarkdownDocument,
      ref: markdownEditorRef
    }
  ];

  const getActiveEditor = () => {
    const activeEditor = editorCandidates
      .find((candidate) => candidate.active?.value && candidate.ref?.value)
      ?.ref?.value;
    if (activeEditor) {
      return activeEditor;
    }
    return editorCandidates
      .map((candidate) => candidate.ref?.value)
      .find(Boolean) || null;
  };

  const callEditorMethod = (method, ...args) => {
    const editor = getActiveEditor();
    if (!editor || typeof editor[method] !== 'function') {
      return false;
    }
    editor[method](...args);
    return true;
  };

  const handleInlineCommentAdd = (payload) => {
    if (isCommentCollapsed.value) {
      isCommentCollapsed.value = false;
    }
    commentSidebarRef.value?.handleInlineAnchorAdd?.(payload);
  };

  const handleInlineCommentRemove = (ids) => {
    commentSidebarRef.value?.handleInlineAnchorRemove?.(ids);
  };

  const handleAnchorHover = (id, hovering) => {
    const targetMethod = hovering ? 'hlCommentIds' : 'unHlCommentIds';
    callEditorMethod(targetMethod, [id]);
  };

  const handleAnchorRemove = (ids) => {
    const targetIds = Array.isArray(ids) ? ids : [ids];
    callEditorMethod('removeCommentIds', targetIds);
  };

  const handleCommentMutated = () => {
    callEditorMethod('broadcastCommentChange');
  };

  const handleRemoteCommentChanged = () => {
    if (remoteCommentReloadTimer.value) {
      clearTimeout(remoteCommentReloadTimer.value);
    }
    remoteCommentReloadTimer.value = setTimeout(() => {
      commentSidebarRef.value?.reload?.();
      remoteCommentReloadTimer.value = null;
    }, 250);
  };

  const scrollToInlineAnchor = (id) => {
    callEditorMethod('scrollToCommentId', id);
  };

  onBeforeUnmount(() => {
    if (remoteCommentReloadTimer.value) {
      clearTimeout(remoteCommentReloadTimer.value);
      remoteCommentReloadTimer.value = null;
    }
  });

  return {
    handleInlineCommentAdd,
    handleInlineCommentRemove,
    handleAnchorHover,
    handleAnchorRemove,
    handleCommentMutated,
    handleRemoteCommentChanged,
    scrollToInlineAnchor
  };
}
