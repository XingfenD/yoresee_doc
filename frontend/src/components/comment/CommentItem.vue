<template>
  <div class="comment-item" :class="rootClass">
    <el-avatar v-if="showAvatar" :size="28" :src="avatar" />
    <div class="comment-item-main">
      <div class="comment-meta">
        <span class="comment-author">{{ author }}</span>
        <div class="comment-meta-actions">
          <span class="comment-time">{{ time }}</span>
          <AppDropdown v-if="hasActions" trigger="click" @command="emitAction">
            <el-button text class="comment-more">
              <el-icon><MoreFilled /></el-icon>
            </el-button>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item
                  v-for="action in actions"
                  :key="action.key"
                  :command="action.key"
                  :divided="action.divided"
                  :disabled="action.disabled"
                  :class="{ 'comment-action-danger': action.danger }"
                >
                  {{ action.label }}
                </el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </AppDropdown>
        </div>
      </div>
      <div v-if="replyText" class="comment-reply">{{ replyText }}</div>
      <slot v-if="editing" name="editor" />
      <div
        v-else
        class="comment-text"
        :class="{ 'comment-text--clickable': contentClickable }"
        @click="handleContentClick"
      >
        {{ content }}
      </div>
      <button
        v-if="!editing && replyable"
        type="button"
        class="comment-reply-action"
        @click="$emit('reply')"
      >
        {{ resolvedReplyLabel }}
      </button>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue';
import { useI18n } from 'vue-i18n';
import { MoreFilled } from '@element-plus/icons-vue';
import AppDropdown from '@/components/base/AppDropdown.vue';

const props = defineProps({
  avatar: {
    type: String,
    default: ''
  },
  author: {
    type: String,
    default: ''
  },
  time: {
    type: String,
    default: ''
  },
  content: {
    type: String,
    default: ''
  },
  replyText: {
    type: String,
    default: ''
  },
  editing: {
    type: Boolean,
    default: false
  },
  variant: {
    type: String,
    default: 'normal'
  },
  showAvatar: {
    type: Boolean,
    default: true
  },
  actions: {
    type: Array,
    default: () => []
  },
  contentClickable: {
    type: Boolean,
    default: false
  },
  replyable: {
    type: Boolean,
    default: false
  },
  replyLabel: {
    type: String,
    default: ''
  }
});

const emit = defineEmits(['content-click', 'action', 'reply']);
const { t } = useI18n();

const rootClass = computed(() => '');
const hasActions = computed(() => Array.isArray(props.actions) && props.actions.length > 0);
const resolvedReplyLabel = computed(() => props.replyLabel || t('common.reply'));

const emitAction = (key) => {
  emit('action', key);
};

const handleContentClick = () => {
  if (!props.contentClickable) {
    return;
  }
  emit('content-click');
};
</script>

<style scoped>
.comment-item {
  display: flex;
  align-items: flex-start;
  gap: 10px;
  position: relative;
}

.comment-item--reply {
  padding-top: 8px;
  padding-bottom: 8px;
  background: var(--bg-light);
  border-radius: var(--border-radius-sm);
}

.comment-item--reply::before {
  content: '';
  position: absolute;
  left: 8px;
  top: 12px;
  bottom: 12px;
  width: 2px;
  background-color: var(--border-color);
  opacity: 0.6;
  border-radius: 2px;
}

.comment-item--reply .comment-item-main {
  padding-left: 6px;
}

.comment-item-main {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.comment-meta {
  display: flex;
  align-items: center;
  justify-content: space-between;
  font-size: 12px;
  color: var(--text-light);
}

.comment-meta-actions {
  display: flex;
  align-items: center;
  gap: 6px;
}

.comment-more {
  padding: 2px 4px;
  color: var(--text-light);
}

.comment-more:hover {
  color: var(--primary-color);
}

.comment-action-danger {
  color: #e53935;
}

.comment-author {
  font-weight: 600;
  color: var(--text-dark);
}

.comment-time {
  font-size: 12px;
  color: var(--text-light);
}

.comment-text {
  font-size: 13px;
  color: var(--text-medium);
  line-height: 1.5;
  word-break: break-word;
}

.comment-text--clickable {
  cursor: pointer;
}

.comment-reply-action {
  align-self: flex-start;
  border: none;
  background: transparent;
  padding: 2px 6px;
  font-size: 11px;
  color: var(--text-light);
  cursor: pointer;
  border-radius: 10px;
  transition: color 0.2s ease, background-color 0.2s ease;
  opacity: 0;
  visibility: hidden;
  pointer-events: none;
}

.comment-item:hover .comment-reply-action {
  opacity: 1;
  visibility: visible;
  pointer-events: auto;
}

.comment-reply-action:hover {
  color: var(--text-medium);
  background-color: rgba(15, 23, 42, 0.06);
}

.comment-reply {
  font-size: 12px;
  color: var(--text-light);
  background: rgba(59, 130, 246, 0.08);
  padding: 2px 6px;
  border-radius: 10px;
  width: fit-content;
}

:global(.dark-mode) .comment-author {
  color: var(--text-dark);
}

:global(.dark-mode) .comment-text {
  color: var(--text-dark);
}

:global(.dark-mode) .comment-meta {
  color: var(--text-light);
}
</style>
