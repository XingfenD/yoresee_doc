import { computed, ref, watch } from 'vue';
import { flattenCommentTree, getIndentStyle as buildIndentStyle, buildReplyLabel } from './inline-comments/tree';
import { useInlineCommentAnchor } from './inline-comments/anchor';
import { useInlineCommentCrud } from './inline-comments/crud';

export function useInlineComments(options = {}) {
  const {
    t,
    getDocId,
    getInlineEnabled,
    getUserInfo,
    onAnchorClick,
    onAnchorHover,
    onAnchorRemove,
    onCommentMutated
  } = options;

  const commentList = ref([]);
  const loading = ref(false);
  const loadVersion = ref(0);

  const inlineEmptyText = computed(() => t('document.inlineCommentEmpty'));
  const userDisplayName = computed(() => {
    const userInfo = getUserInfo?.();
    return userInfo?.nickname || userInfo?.username || '我';
  });
  const titleWithCount = computed(() => {
    const count = commentList.value.length;
    return count > 0 ? `${count}` : '';
  });
  const displayComments = computed(() => flattenCommentTree(commentList.value));
  const getIndentStyle = (item) => buildIndentStyle(item);
  const replyLabel = (item) => buildReplyLabel(item, commentList.value, t);

  let maybeRemoveAnchor = () => {};

  const commentCrud = useInlineCommentCrud({
    t,
    getDocId,
    getInlineEnabled,
    getUserInfo,
    commentList,
    userDisplayName,
    maybeRemoveAnchor: (...args) => maybeRemoveAnchor(...args),
    onCommentMutated
  });

  const anchorBridge = useInlineCommentAnchor({
    commentList,
    onAnchorClick,
    onAnchorHover,
    onAnchorRemove,
    createTempComment: commentCrud.createTempComment,
    deleteRemoteComment: commentCrud.deleteRemoteComment
  });
  maybeRemoveAnchor = anchorBridge.maybeRemoveAnchor;

  const loadComments = async () => {
    const requestVersion = loadVersion.value + 1;
    loading.value = true;
    try {
      await commentCrud.loadComments(loadVersion);
    } finally {
      if (requestVersion === loadVersion.value) {
        loading.value = false;
      }
    }
  };

  watch(
    () => getDocId?.(),
    async () => {
      commentList.value = [];
      loading.value = false;
      await loadComments();
    },
    { immediate: true }
  );

  watch(
    () => !!getInlineEnabled?.(),
    async (enabled) => {
      if (!enabled) {
        commentList.value = [];
        return;
      }
      await loadComments();
    }
  );

  return {
    inlineEmptyText,
    titleWithCount,
    displayComments,
    getIndentStyle,
    getActions: commentCrud.getActions,
    replyLabel,
    handleAction: commentCrud.handleAction,
    handleReply: commentCrud.handleReply,
    handleContentClick: anchorBridge.handleContentClick,
    handleHover: anchorBridge.handleHover,
    saveEdit: commentCrud.saveEdit,
    cancelEdit: commentCrud.cancelEdit,
    formatCommentTime: commentCrud.formatCommentTime,
    handleInlineAnchorAdd: anchorBridge.handleInlineAnchorAdd,
    handleInlineAnchorRemove: anchorBridge.handleInlineAnchorRemove,
    reload: loadComments
  };
}
