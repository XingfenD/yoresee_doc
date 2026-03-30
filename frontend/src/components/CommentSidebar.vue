<template>
  <PanelSidebarShell
    :collapsed="collapsed"
    :resizing="resizing"
    resize-edge="left"
    :container-style="containerStyle"
    :panel-style="sidebarStyle"
    @resize-start="startResize"
  >
    <template #header>
      <div class="comment-header">
        <div class="comment-title">{{ displayTitle }}</div>
        <el-button text class="collapse-button" @click="$emit('toggle')" :title="collapseTitle">
          <el-icon>
            <ArrowRight />
          </el-icon>
        </el-button>
      </div>
    </template>
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
            :reply-label="t('common.reply')"
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
                  :placeholder="t('document.inlineCommentInputPlaceholder')"
                />
                <div class="inline-comment-editor-actions">
                  <el-button size="small" type="primary" :loading="item.saving" @click="saveEdit(item)">
                    {{ t('common.save') }}
                  </el-button>
                  <el-button size="small" text @click="cancelEdit(item)">
                    {{ t('common.cancel') }}
                  </el-button>
                </div>
              </div>
            </template>
          </CommentItem>
        </template>
      </CommentList>
    </div>
  </PanelSidebarShell>
</template>

<script setup>
import { computed } from 'vue';
import { useI18n } from 'vue-i18n';
import { ArrowRight } from '@element-plus/icons-vue';
import CommentList from '@/components/CommentList.vue';
import CommentItem from '@/components/CommentItem.vue';
import PanelSidebarShell from '@/components/PanelSidebarShell.vue';
import { usePanelSidebar } from '@/composables/usePanelSidebar';
import { useInlineComments } from '@/composables/useInlineComments';

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

const {
  width: commentWidth,
  resizing,
  startResize
} = usePanelSidebar({
  defaultWidth: 320,
  minWidth: 280,
  maxWidth: 560,
  resizeEdge: 'left',
  widthStorageKey: 'commentSidebarWidth',
  externalCollapsed: computed(() => props.collapsed),
  getMaxWidth: () => Math.min(560, Math.max(280, window.innerWidth - 360))
});

const containerStyle = computed(() => {
  if (props.collapsed) {
    return {};
  }
  return { width: `${commentWidth.value}px` };
});
const sidebarStyle = computed(() => ({ width: `${commentWidth.value}px` }));
const {
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
  reload
} = useInlineComments({
  t,
  getDocId: () => props.docId,
  getInlineEnabled: () => props.inlineEnabled,
  getUserInfo: () => props.userInfo,
  onAnchorClick: (id) => props.onAnchorClick?.(id),
  onAnchorHover: (id, hovering) => props.onAnchorHover?.(id, hovering),
  onAnchorRemove: (ids) => props.onAnchorRemove?.(ids),
  onCommentMutated: (payload) => props.onCommentMutated?.(payload)
});
const displayTitle = computed(() => {
  return titleWithCount.value ? `${props.title} (${titleWithCount.value})` : props.title;
});

defineExpose({
  handleInlineAnchorAdd,
  handleInlineAnchorRemove,
  reload
});
</script>

<style scoped>
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

.collapse-button {
  padding: 4px;
  color: var(--text-light);
}

.collapse-button:hover {
  color: var(--primary-color);
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

:global(.dark-mode) .comment-title {
  color: var(--text-dark);
}
</style>
