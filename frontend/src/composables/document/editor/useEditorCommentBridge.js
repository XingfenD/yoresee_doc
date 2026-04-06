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

  const getActiveEditorRef = () => {
    if (isRichTextDocument?.value) {
      return richTextEditorRef.value;
    }
    if (isMarkdownDocument?.value) {
      return markdownEditorRef.value;
    }
    return markdownEditorRef.value || richTextEditorRef.value;
  };

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
    const editorRef = getActiveEditorRef();
    if (editorRef && typeof editorRef.hlCommentIds === 'function') {
      editorRef.hlCommentIds([id]);
      return;
    }
    const vditor = getVditorInstance();
    if (vditor && typeof vditor.hlCommentIds === 'function') {
      vditor.hlCommentIds([id]);
    }
  };

  const unhighlightInlineComment = (id) => {
    const editorRef = getActiveEditorRef();
    if (editorRef && typeof editorRef.unHlCommentIds === 'function') {
      editorRef.unHlCommentIds([id]);
      return;
    }
    const vditor = getVditorInstance();
    if (vditor && typeof vditor.unHlCommentIds === 'function') {
      vditor.unHlCommentIds([id]);
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
    const editorRef = getActiveEditorRef();
    if (editorRef && typeof editorRef.removeCommentIds === 'function') {
      editorRef.removeCommentIds(ids);
      return;
    }
    const vditor = getVditorInstance();
    if (!vditor || typeof vditor.removeCommentIds !== 'function') {
      return;
    }
    vditor.removeCommentIds(Array.isArray(ids) ? ids : [ids]);
  };

  const handleCommentMutated = () => {
    getActiveEditorRef()?.broadcastCommentChange?.();
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
    const editorRef = getActiveEditorRef();
    if (editorRef && typeof editorRef.scrollToCommentId === 'function') {
      editorRef.scrollToCommentId(id);
      return;
    }
    const vditor = getVditorInstance();
    if (!vditor || typeof vditor.getCommentIds !== 'function') {
      return;
    }
    const commentEntries = vditor.getCommentIds();
    const target = Array.isArray(commentEntries) ? commentEntries.find((entry) => entry.id === id) : null;
    const container = vditor?.vditor?.wysiwyg?.element || vditor?.vditor?.ir?.element;
    if (!target || !container) {
      return;
    }
    const top = Math.max(target.top - 24, 0);
    if (typeof container.scrollTo === 'function') {
      container.scrollTo({ top, behavior: 'smooth' });
    } else {
      container.scrollTop = top;
    }
    if (typeof vditor.hlCommentIds === 'function') {
      vditor.hlCommentIds([id]);
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
