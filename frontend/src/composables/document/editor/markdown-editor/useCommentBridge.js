import { YJS_COMMENT_META_FIELD } from '@/composables/document/editor/collab/useYjsContentBridge';

export function useCommentBridge({ props, emit, vditorRef, ydocRef, syncEditorToYjs }) {
  const getCommentOptions = () => {
    if (!props.commentEnabled) {
      return { enable: false };
    }
    return {
      enable: true,
      add: (id, text, commentsData) => {
        emit('comment-add', { id, text, commentsData });
        queueMicrotask(() => {
          syncEditorToYjs();
        });
      },
      remove: (ids) => {
        emit('comment-remove', ids);
        queueMicrotask(() => {
          syncEditorToYjs();
        });
      },
      scroll: (top) => {
        emit('comment-scroll', top);
      },
      adjustTop: (commentsData) => {
        emit('comment-adjust', commentsData);
      }
    };
  };

  const removeCommentIds = (ids = []) => {
    const vditor = vditorRef.value;
    if (!vditor || typeof vditor.removeCommentIds !== 'function') {
      return;
    }
    if (!Array.isArray(ids) || ids.length === 0) {
      return;
    }
    vditor.removeCommentIds(ids);
    queueMicrotask(() => {
      syncEditorToYjs();
    });
  };

  const broadcastCommentChange = () => {
    if (!props.collabEnabled || !ydocRef.value) {
      return;
    }
    const commentMeta = ydocRef.value.getMap(YJS_COMMENT_META_FIELD);
    commentMeta.set('tick', `${Date.now()}_${Math.random().toString(36).slice(2)}`);
  };

  return {
    getCommentOptions,
    removeCommentIds,
    broadcastCommentChange
  };
}
