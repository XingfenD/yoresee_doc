<template>
  <div class="comment-container" :class="{ collapsed: collapsed, resizing: resizing }" :style="containerStyle">
    <div
      class="comment-resizer"
      role="separator"
      aria-orientation="vertical"
      @mousedown="startResize"
    ></div>
    <aside class="comment-sidebar" :style="sidebarStyle">
      <div class="comment-header">
        <div class="comment-title">{{ titleWithCount }}</div>
        <el-button text class="collapse-button" @click="$emit('toggle')" :title="collapseTitle">
          <el-icon>
            <ArrowRight />
          </el-icon>
        </el-button>
      </div>
      <div class="comment-body">
        <CommentList
          :show-title="false"
          :items="displayComments"
          :empty-text="inlineEmptyText"
          key-field="local_id"
        >
          <template #item="{ item }">
            <CommentItem
              :class="{ 'comment-item--reply': item.level > 0 }"
              :style="getIndentStyle(item)"
              :avatar="item.creator_avatar"
              :author="item.creator_name"
              :time="formatCommentTime(item.created_at)"
              :content="item.content"
              :reply-text="item.parent_external_id ? replyLabel(item) : ''"
              :actions="getActions(item)"
              :editing="item.editing"
              :content-clickable="true"
              :replyable="true"
              reply-label="回复"
              @action="(action) => handleAction(item, action)"
              @reply="() => handleReply(item)"
              @content-click="() => handleContentClick(item)"
              @mouseenter="() => handleHover(item, true)"
              @mouseleave="() => handleHover(item, false)"
            >
              <template #editor>
                <div class="inline-comment-editor">
                  <el-input
                    v-model="item.draft"
                    type="textarea"
                    :autosize="{ minRows: 2, maxRows: 4 }"
                    placeholder="输入评论内容..."
                  />
                  <div class="inline-comment-editor-actions">
                    <el-button size="small" type="primary" :loading="item.saving" @click="saveEdit(item)">保存</el-button>
                    <el-button size="small" text @click="cancelEdit(item)">取消</el-button>
                  </div>
                </div>
              </template>
            </CommentItem>
          </template>
        </CommentList>
      </div>
    </aside>
  </div>
</template>

<script setup>
import { computed, onBeforeUnmount, ref, watch } from 'vue';
import { ElMessage, ElMessageBox } from 'element-plus';
import { useI18n } from 'vue-i18n';
import { ArrowRight } from '@element-plus/icons-vue';
import CommentList from '@/components/CommentList.vue';
import CommentItem from '@/components/CommentItem.vue';
import {
  listDocumentComments,
  createDocumentComment,
  deleteDocumentComment,
  updateDocumentComment
} from '@/services/api';

const props = defineProps({
  title: { type: String, default: '' },
  collapseTitle: { type: String, default: '' },
  collapsed: { type: Boolean, default: false },
  docId: { type: String, default: '' },
  inlineEnabled: { type: Boolean, default: false },
  userInfo: { type: Object, default: null },
  onAnchorClick: { type: Function, default: null },
  onAnchorHover: { type: Function, default: null },
  onAnchorRemove: { type: Function, default: null },
  onCommentMutated: { type: Function, default: null }
});

defineEmits(['toggle']);

const { t } = useI18n();

const commentList = ref([]);
const loading = ref(false);
const loadVersion = ref(0);
const inlineEmptyText = '暂无行内评论';
const COMMENT_WIDTH_KEY = 'commentSidebarWidth';
const MIN_COMMENT_WIDTH = 280;
const DEFAULT_COMMENT_WIDTH = 320;
const MAX_COMMENT_WIDTH = 560;
const commentWidth = ref(DEFAULT_COMMENT_WIDTH);
const resizing = ref(false);
const resizeStartX = ref(0);
const resizeStartWidth = ref(DEFAULT_COMMENT_WIDTH);

const userDisplayName = computed(() => props.userInfo?.nickname || props.userInfo?.username || '我');
const containerStyle = computed(() => {
  if (props.collapsed) {
    return {};
  }
  return { width: `${commentWidth.value}px` };
});
const sidebarStyle = computed(() => ({ width: `${commentWidth.value}px` }));

const titleWithCount = computed(() => {
  const count = commentList.value.length;
  return count > 0 ? `${props.title} (${count})` : props.title;
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
    // Temporary editor cards have empty external_id. Do not resolve children by empty key,
    // otherwise root nodes recurse into the whole root list.
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
  const currentExternalId = props.userInfo?.external_id;
  if (!currentExternalId) return false;
  if (item.creator_user_external_id === currentExternalId) {
    return true;
  }
  return props.userInfo?.username === 'admin';
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

const createTempComment = ({ parent, anchorId }) => ({
  local_id: buildLocalId(),
  external_id: '',
  parent_external_id: parent?.external_id || '',
  content: '',
  created_at: '',
  creator_name: userDisplayName.value,
  creator_user_external_id: props.userInfo?.external_id || '',
  creator_avatar: props.userInfo?.avatar || '',
  anchor_id: `${anchorId || parent?.anchor_id || ''}`.trim(),
  editing: true,
  draft: '',
  saving: false
});

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
  if (item?.anchor_id && typeof props.onAnchorClick === 'function') {
    props.onAnchorClick(item.anchor_id);
  }
};

const handleHover = (item, hovering) => {
  if (item?.anchor_id && typeof props.onAnchorHover === 'function') {
    props.onAnchorHover(item.anchor_id, hovering);
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
    ElMessage.error('请输入评论内容');
    return;
  }
  if (!item.anchor_id) {
    ElMessage.error('未找到行内锚点');
    return;
  }
  if (item.saving) return;
  if (!props.docId || props.docId === 'example') return;

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
      if (typeof props.onCommentMutated === 'function') {
        props.onCommentMutated({ type: 'update', comment_id: item.external_id, anchor_id: item.anchor_id });
      }
    } else {
      const resp = await createDocumentComment({
        document_external_id: props.docId,
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
      if (typeof props.onCommentMutated === 'function') {
        props.onCommentMutated({ type: 'create', comment_id: item.external_id, anchor_id: item.anchor_id });
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
      if (typeof props.onCommentMutated === 'function') {
        props.onCommentMutated({ type: 'delete', comment_id: item.external_id, anchor_id: deletedAnchorId });
      }
    } catch (error) {
      return;
    }
  }
  const nextList = commentList.value.filter((entry) => entry.local_id !== item.local_id);
  commentList.value = nextList;
  maybeRemoveAnchor(item.anchor_id, nextList);
};

const maybeRemoveAnchor = (anchorId, list = commentList.value) => {
  const id = `${anchorId || ''}`.trim();
  if (!id) {
    return;
  }
  const hasAnyComment = (list || []).some((entry) => {
    return `${entry?.anchor_id || ''}`.trim() === id;
  });
  if (hasAnyComment) {
    return;
  }
  if (typeof props.onAnchorRemove === 'function') {
    props.onAnchorRemove([id]);
  }
};

const getMaxAllowedWidth = () => {
  const viewportMax = Math.max(MIN_COMMENT_WIDTH, window.innerWidth - 360);
  return Math.min(MAX_COMMENT_WIDTH, viewportMax);
};

const clampCommentWidth = (width) => {
  return Math.min(Math.max(width, MIN_COMMENT_WIDTH), getMaxAllowedWidth());
};

const onResizeMove = (event) => {
  if (!resizing.value) return;
  const delta = resizeStartX.value - event.clientX;
  commentWidth.value = clampCommentWidth(resizeStartWidth.value + delta);
};

const stopResize = () => {
  if (!resizing.value) return;
  resizing.value = false;
  document.body.style.cursor = '';
  document.body.style.userSelect = '';
  window.removeEventListener('mousemove', onResizeMove);
  window.removeEventListener('mouseup', stopResize);
};

const startResize = (event) => {
  if (props.collapsed) return;
  event.preventDefault();
  resizing.value = true;
  resizeStartX.value = event.clientX;
  resizeStartWidth.value = commentWidth.value;
  document.body.style.cursor = 'col-resize';
  document.body.style.userSelect = 'none';
  window.addEventListener('mousemove', onResizeMove);
  window.addEventListener('mouseup', stopResize);
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

const loadComments = async () => {
  const requestVersion = ++loadVersion.value;
  const requestDocId = props.docId;
  if (!props.docId || props.docId === 'example' || !props.inlineEnabled) {
    commentList.value = [];
    if (requestVersion === loadVersion.value) {
      loading.value = false;
    }
    return;
  }

  loading.value = true;
  try {
    const resp = await listDocumentComments({
      document_external_id: props.docId,
      page: 1,
      page_size: 100
    });
    if (requestVersion !== loadVersion.value || requestDocId !== props.docId) {
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
    if (requestVersion !== loadVersion.value || requestDocId !== props.docId) {
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
  () => props.docId,
  async () => {
    commentList.value = [];
    loading.value = false;
    await loadComments();
  },
  { immediate: true }
);

watch(
  () => props.inlineEnabled,
  async (enabled) => {
    if (!enabled) {
      commentList.value = [];
      return;
    }
    await loadComments();
  }
);

watch(commentWidth, (value) => {
  localStorage.setItem(COMMENT_WIDTH_KEY, `${value}`);
});

watch(
  () => props.collapsed,
  () => {
    stopResize();
  }
);

onBeforeUnmount(() => {
  stopResize();
});

const savedCommentWidth = Number(localStorage.getItem(COMMENT_WIDTH_KEY) || DEFAULT_COMMENT_WIDTH);
if (Number.isFinite(savedCommentWidth)) {
  commentWidth.value = clampCommentWidth(savedCommentWidth);
}

defineExpose({
  handleInlineAnchorAdd,
  handleInlineAnchorRemove,
  reload: async () => {
    await loadComments();
  }
});
</script>

<style scoped>
.comment-container {
  position: relative;
  width: 320px;
  overflow: hidden;
  flex-shrink: 0;
  border-left: 1px solid var(--border-color);
  background-color: var(--bg-white);
  transition: all 0.3s ease-in-out;
  display: flex;
}

.comment-container.resizing,
.comment-container.resizing .comment-sidebar,
.comment-container.resizing .comment-resizer {
  transition: none !important;
}

.comment-container.collapsed {
  width: 0;
  opacity: 0;
  visibility: hidden;
  pointer-events: none;
  border-left: none;
}

.comment-sidebar {
  display: flex;
  flex-direction: column;
  width: 320px;
  transition: transform 0.3s ease-in-out, opacity 0.3s ease-in-out;
  background-color: var(--bg-white);
}

.comment-resizer {
  width: 6px;
  flex-shrink: 0;
  cursor: col-resize;
  background-color: var(--bg-light);
  border-right: 1px solid var(--border-color);
  transition: background-color 0.2s ease;
}

.comment-resizer:hover {
  background-color: var(--bg-medium);
}

.comment-container.collapsed .comment-resizer {
  display: none;
}

.comment-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: var(--spacing-md);
  border-bottom: 1px solid var(--border-color);
}

.comment-title {
  font-size: 16px;
  font-weight: 600;
  color: var(--text-dark);
}

.comment-body {
  flex: 1;
  overflow-y: auto;
  padding: var(--spacing-md);
}

.inline-comment-editor {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.inline-comment-editor-actions {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
}

:global(.dark-mode) .comment-container,
:global(.dark-mode) .comment-sidebar {
  background-color: var(--bg-white);
  border-color: var(--border-color);
}

:global(.dark-mode) .comment-resizer {
  background-color: var(--bg-medium);
  border-color: var(--border-color);
}

:global(.dark-mode) .comment-resizer:hover {
  background-color: var(--bg-white);
}

:global(.dark-mode) .comment-title {
  color: var(--text-dark);
}
</style>
