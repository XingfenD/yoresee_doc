import { ElMessage, ElMessageBox } from 'element-plus';
import {
  listDocumentComments,
  createDocumentComment,
  deleteDocumentComment,
  updateDocumentComment
} from '@/services/api';

export function useInlineCommentCrud({
  t,
  getDocId,
  getInlineEnabled,
  getUserInfo,
  commentList,
  userDisplayName,
  maybeRemoveAnchor,
  onCommentMutated
}) {
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

  const insertAfter = (list, predicate, item) => {
    const index = list.findIndex(predicate);
    if (index === -1) {
      list.unshift(item);
      return;
    }
    list.splice(index + 1, 0, item);
  };

  const handleReply = (item) => {
    const temp = createTempComment({ parent: item });
    insertAfter(commentList.value, (entry) => entry.external_id === item.external_id, temp);
  };

  const startEdit = (item) => {
    if (!item || item.editing) return;
    item.draft = item.content || '';
    item.editing = true;
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
      ElMessage.error(t('document.inlineCommentContentRequired'));
      return;
    }
    if (!item.anchor_id) {
      ElMessage.error(t('document.inlineCommentAnchorMissing'));
      return;
    }
    if (item.saving) return;

    const docId = getDocId?.();
    if (!docId || docId === 'example') {
      return;
    }

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

  const formatCommentTime = (value) => {
    if (!value) return '';
    const date = new Date(value);
    if (Number.isNaN(date.getTime())) return value;
    return date.toLocaleString();
  };

  const loadComments = async (loadVersionRef) => {
    const requestVersion = ++loadVersionRef.value;
    const requestDocId = getDocId?.();
    const inlineEnabled = !!getInlineEnabled?.();

    if (!requestDocId || requestDocId === 'example' || !inlineEnabled) {
      commentList.value = [];
      if (requestVersion === loadVersionRef.value) {
        return false;
      }
      return true;
    }

    try {
      const resp = await listDocumentComments({
        document_external_id: requestDocId,
        page: 1,
        page_size: 100
      });
      if (requestVersion !== loadVersionRef.value || requestDocId !== getDocId?.()) {
        return false;
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
      return true;
    } catch (error) {
      if (requestVersion !== loadVersionRef.value || requestDocId !== getDocId?.()) {
        return false;
      }
      commentList.value = [];
      return true;
    }
  };

  return {
    createTempComment,
    getActions,
    handleReply,
    handleAction,
    saveEdit,
    cancelEdit,
    formatCommentTime,
    loadComments,
    deleteRemoteComment: deleteDocumentComment
  };
}
