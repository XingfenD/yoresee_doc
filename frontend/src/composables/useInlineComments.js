import { computed, ref, watch } from 'vue';
import { ElMessage, ElMessageBox } from 'element-plus';
import {
  listDocumentComments,
  createDocumentComment,
  deleteDocumentComment,
  updateDocumentComment
} from '@/services/api';

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
  const inlineEmptyText = '暂无行内评论';

  const userDisplayName = computed(() => {
    const userInfo = getUserInfo?.();
    return userInfo?.nickname || userInfo?.username || '我';
  });

  const titleWithCount = computed(() => {
    const count = commentList.value.length;
    return count > 0 ? `${count}` : '';
  });

  const flattenComments = (items) => {
    if (!Array.isArray(items) || items.length === 0) return [];
    const childrenMap = new Map();
    const idSet = new Set();
    items.forEach((item) => {
      if (item?.external_id) {
        idSet.add(item.external_id);
      }
    });
    items.forEach((item) => {
      const parentId = item.parent_external_id || '';
      if (!childrenMap.has(parentId)) {
        childrenMap.set(parentId, []);
      }
      childrenMap.get(parentId).push(item);
    });

    const result = [];
    const walk = (node, level) => {
      node.level = level;
      result.push(node);
      if (!node.external_id) {
        return;
      }
      const children = childrenMap.get(node.external_id) || [];
      children.forEach((child) => walk(child, level + 1));
    };

    const roots = [];
    items.forEach((item) => {
      const parentId = item.parent_external_id || '';
      if (!parentId || !idSet.has(parentId)) {
        roots.push(item);
      }
    });
    roots.forEach((root) => walk(root, 0));
    return result;
  };

  const displayComments = computed(() => flattenComments(commentList.value));

  const getIndentStyle = (item) => ({
    paddingLeft: `${Math.min(item.level || 0, 3) * 16}px`
  });

  const canModify = (item) => {
    if (!item) return false;
    if (!item.external_id) return true;
    const userInfo = getUserInfo?.();
    const currentExternalId = userInfo?.external_id;
    if (!currentExternalId) return false;
    if (item.creator_user_external_id === currentExternalId) {
      return true;
    }
    return userInfo?.username === 'admin';
  };

  const getActions = (item) => {
    const editable = canModify(item);
    return [
      { key: 'copy', label: t('common.copy') },
      { key: 'edit', label: t('common.edit'), disabled: !editable },
      { key: 'delete', label: t('document.commentDelete'), danger: true, disabled: !editable }
    ];
  };

  const buildLocalId = () => `temp_${Date.now()}_${Math.random().toString(36).slice(2)}`;

  const createTempComment = ({ parent, anchorId }) => {
    const userInfo = getUserInfo?.();
    return {
      local_id: buildLocalId(),
      external_id: '',
      parent_external_id: parent?.external_id || '',
      content: '',
      created_at: '',
      creator_name: userDisplayName.value,
      creator_user_external_id: userInfo?.external_id || '',
      creator_avatar: userInfo?.avatar || '',
      anchor_id: `${anchorId || parent?.anchor_id || ''}`.trim(),
      editing: true,
      draft: '',
      saving: false
    };
  };

  const insertAfter = (list, predicate, item) => {
    const index = list.findIndex(predicate);
    if (index === -1) {
      list.unshift(item);
      return;
    }
    list.splice(index + 1, 0, item);
  };

  const replyLabel = (item) => {
    if (!item?.parent_external_id) return '';
    const parent = commentList.value.find((entry) => entry.external_id === item.parent_external_id);
    return t('document.commentReplyTo', { name: parent?.creator_name || t('document.commentUnknown') });
  };

  const handleReply = (item) => {
    const temp = createTempComment({ parent: item });
    insertAfter(commentList.value, (entry) => entry.external_id === item.external_id, temp);
  };

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

  const startEdit = (item) => {
    if (!item || item.editing) return;
    item.draft = item.content || '';
    item.editing = true;
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

  const cancelEdit = (item) => {
    if (!item) return;
    if (!item.external_id) {
      const nextList = commentList.value.filter((entry) => entry.local_id !== item.local_id);
      commentList.value = nextList;
      maybeRemoveAnchor(item.anchor_id, nextList);
      return;
    }
    item.editing = false;
    item.draft = '';
  };

  const saveEdit = async (item) => {
    if (!item) return;
    const content = (item.draft || '').trim();
    if (!content) {
      ElMessage.error('请输入评论内容');
      return;
    }
    if (!item.anchor_id) {
      ElMessage.error('未找到行内锚点');
      return;
    }
    if (item.saving) return;
    const docId = getDocId?.();
    if (!docId || docId === 'example') return;

    item.saving = true;
    try {
      if (item.external_id) {
        const resp = await updateDocumentComment({
          external_id: item.external_id,
          content
        });
        const saved = resp.comment;
        item.content = saved?.content || content;
        item.created_at = saved?.created_at || item.created_at;
        item.creator_name = saved?.creator_name || item.creator_name;
        item.creator_avatar = saved?.creator_avatar || item.creator_avatar;
        item.creator_user_external_id = saved?.creator_user_external_id || item.creator_user_external_id;
        item.editing = false;
        item.draft = '';
        if (typeof onCommentMutated === 'function') {
          onCommentMutated({ type: 'update', comment_id: item.external_id, anchor_id: item.anchor_id });
        }
      } else {
        const resp = await createDocumentComment({
          document_external_id: docId,
          content,
          parent_external_id: item.parent_external_id || undefined,
          anchor_id: item.anchor_id
        });
        const saved = resp.comment;
        item.content = saved?.content || content;
        item.external_id = saved?.external_id || item.external_id;
        item.local_id = item.external_id || item.local_id;
        item.created_at = saved?.created_at || new Date().toISOString();
        item.creator_name = saved?.creator_name || item.creator_name;
        item.creator_avatar = saved?.creator_avatar || item.creator_avatar;
        item.creator_user_external_id = saved?.creator_user_external_id || item.creator_user_external_id;
        item.anchor_id = saved?.anchor_id || item.anchor_id;
        item.editing = false;
        item.draft = '';
        if (typeof onCommentMutated === 'function') {
          onCommentMutated({ type: 'create', comment_id: item.external_id, anchor_id: item.anchor_id });
        }
      }
    } catch (error) {
      ElMessage.error(t('common.requestFailed'));
    } finally {
      item.saving = false;
    }
  };

  const deleteCommentItem = async (item) => {
    if (!item) return;
    const deletedAnchorId = item.anchor_id;
    if (item.external_id) {
      try {
        await ElMessageBox.confirm(
          t('document.commentDeleteConfirm'),
          t('document.commentDelete'),
          {
            confirmButtonText: t('button.confirm'),
            cancelButtonText: t('button.cancel'),
            type: 'warning'
          }
        );
        await deleteDocumentComment(item.external_id);
        if (typeof onCommentMutated === 'function') {
          onCommentMutated({ type: 'delete', comment_id: item.external_id, anchor_id: deletedAnchorId });
        }
      } catch (error) {
        return;
      }
    }
    const nextList = commentList.value.filter((entry) => entry.local_id !== item.local_id);
    commentList.value = nextList;
    maybeRemoveAnchor(item.anchor_id, nextList);
  };

  const copyComment = async (item) => {
    const text = item?.content || '';
    if (!text) return;
    try {
      await navigator.clipboard.writeText(text);
      ElMessage.success(t('common.copySuccess'));
    } catch (error) {
      ElMessage.error(t('common.copyFailed'));
    }
  };

  const handleAction = (item, action) => {
    if (action === 'copy') {
      copyComment(item);
      return;
    }
    if (action === 'edit') {
      if (!canModify(item)) return;
      startEdit(item);
      return;
    }
    if (action === 'delete') {
      if (!canModify(item)) return;
      deleteCommentItem(item);
    }
  };

  const loadComments = async () => {
    const requestVersion = ++loadVersion.value;
    const requestDocId = getDocId?.();
    const inlineEnabled = !!getInlineEnabled?.();
    if (!requestDocId || requestDocId === 'example' || !inlineEnabled) {
      commentList.value = [];
      if (requestVersion === loadVersion.value) {
        loading.value = false;
      }
      return;
    }

    loading.value = true;
    try {
      const resp = await listDocumentComments({
        document_external_id: requestDocId,
        page: 1,
        page_size: 100
      });
      if (requestVersion !== loadVersion.value || requestDocId !== getDocId?.()) {
        return;
      }
      commentList.value = (resp.comments || [])
        .filter((item) => item.anchor_id)
        .map((item) => ({
          local_id: item.external_id || item.anchor_id,
          anchor_id: item.anchor_id,
          external_id: item.external_id,
          parent_external_id: item.parent_external_id || '',
          content: item.content,
          created_at: item.created_at,
          creator_name: item.creator_name,
          creator_user_external_id: item.creator_user_external_id,
          creator_avatar: item.creator_avatar,
          editing: false,
          draft: '',
          saving: false
        }));
    } catch (error) {
      if (requestVersion !== loadVersion.value || requestDocId !== getDocId?.()) {
        return;
      }
      commentList.value = [];
    } finally {
      if (requestVersion === loadVersion.value) {
        loading.value = false;
      }
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
    if (!Array.isArray(ids) || ids.length === 0) return;
    const idSet = new Set(ids);
    const targets = commentList.value.filter((item) => idSet.has(item.anchor_id));
    commentList.value = commentList.value.filter((item) => !idSet.has(item.anchor_id));
    await Promise.all(
      targets
        .filter((item) => item.external_id)
        .map((item) => deleteDocumentComment(item.external_id).catch(() => {}))
    );
  };

  const formatCommentTime = (value) => {
    if (!value) return '';
    const date = new Date(value);
    if (Number.isNaN(date.getTime())) return value;
    return date.toLocaleString();
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
    getActions,
    replyLabel,
    handleAction,
    handleReply,
    handleContentClick,
    handleHover,
    saveEdit,
    cancelEdit,
    formatCommentTime,
    handleInlineAnchorAdd,
    handleInlineAnchorRemove,
    reload: loadComments
  };
}
