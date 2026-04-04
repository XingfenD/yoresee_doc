export function useInlineCommentAnchor({
  commentList,
  onAnchorClick,
  onAnchorHover,
  onAnchorRemove,
  createTempComment,
  deleteRemoteComment
}) {
  const handleContentClick = (item) => {
    if (item?.anchor_id && typeof onAnchorClick === 'function') {
      onAnchorClick(item.anchor_id);
    }
  };

  const handleHover = (item, hovering) => {
    if (item?.anchor_id && typeof onAnchorHover === 'function') {
      onAnchorHover(item.anchor_id, hovering);
    }
  };

  const maybeRemoveAnchor = (anchorId, list = commentList.value) => {
    const id = `${anchorId || ''}`.trim();
    if (!id) {
      return;
    }
    const hasAnyComment = (list || []).some((entry) => `${entry?.anchor_id || ''}`.trim() === id);
    if (hasAnyComment) {
      return;
    }
    if (typeof onAnchorRemove === 'function') {
      onAnchorRemove([id]);
    }
  };

  const handleInlineAnchorAdd = (payload) => {
    const rawId = typeof payload === 'string' ? payload : payload?.id;
    const id = `${rawId || ''}`.trim();
    if (!id) return;
    if (commentList.value.some((item) => item.anchor_id === id && item.editing)) return;
    commentList.value.unshift(createTempComment({ parent: null, anchorId: id }));
  };

  const handleInlineAnchorRemove = async (ids) => {
    if (!Array.isArray(ids) || ids.length === 0) {
      return;
    }
    const idSet = new Set(ids);
    const targets = commentList.value.filter((item) => idSet.has(item.anchor_id));
    commentList.value = commentList.value.filter((item) => !idSet.has(item.anchor_id));
    await Promise.all(
      targets
        .filter((item) => item.external_id)
        .map((item) => deleteRemoteComment(item.external_id).catch(() => {}))
    );
  };

  return {
    handleContentClick,
    handleHover,
    maybeRemoveAnchor,
    handleInlineAnchorAdd,
    handleInlineAnchorRemove
  };
}
