import { onBeforeUnmount, ref } from 'vue';

export function useEditorCommentBridge({
  isCommentCollapsed,
  markdownEditorRef,
  commentSidebarRef
}) {
  const remoteCommentReloadTimer = ref(null);

  const getVditorInstance = () => markdownEditorRef.value?.getVditor?.();

  const handleInlineCommentAdd = (payload) => {
    if (isCommentCollapsed.value) {
      isCommentCollapsed.value = false;
    }
    commentSidebarRef.value?.handleInlineAnchorAdd?.(payload);
  };

  const handleInlineCommentRemove = (ids) => {
    commentSidebarRef.value?.handleInlineAnchorRemove?.(ids);
  };

  const highlightInlineComment = (id) => {
    const editor = getVditorInstance();
    if (editor && typeof editor.hlCommentIds === 'function') {
      editor.hlCommentIds([id]);
    }
  };

  const unhighlightInlineComment = (id) => {
    const editor = getVditorInstance();
    if (editor && typeof editor.unHlCommentIds === 'function') {
      editor.unHlCommentIds([id]);
    }
  };

  const handleAnchorHover = (id, hovering) => {
    if (hovering) {
      highlightInlineComment(id);
      return;
    }
    unhighlightInlineComment(id);
  };

  const handleAnchorRemove = (ids) => {
    const markdownEditor = markdownEditorRef.value;
    if (markdownEditor && typeof markdownEditor.removeCommentIds === 'function') {
      markdownEditor.removeCommentIds(ids);
      return;
    }
    const editor = getVditorInstance();
    if (!editor || typeof editor.removeCommentIds !== 'function') {
      return;
    }
    editor.removeCommentIds(Array.isArray(ids) ? ids : [ids]);
  };

  const handleCommentMutated = () => {
    markdownEditorRef.value?.broadcastCommentChange?.();
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
    const editor = getVditorInstance();
    if (!editor || typeof editor.getCommentIds !== 'function') {
      return;
    }
    const commentEntries = editor.getCommentIds();
    const target = Array.isArray(commentEntries)
      ? commentEntries.find((entry) => entry.id === id)
      : null;
    const container = editor?.vditor?.wysiwyg?.element || editor?.vditor?.ir?.element;
    if (!target || !container) {
      return;
    }
    const top = Math.max(target.top - 24, 0);
    if (typeof container.scrollTo === 'function') {
      container.scrollTo({ top, behavior: 'smooth' });
    } else {
      container.scrollTop = top;
    }
    if (typeof editor.hlCommentIds === 'function') {
      editor.hlCommentIds([id]);
    }
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
