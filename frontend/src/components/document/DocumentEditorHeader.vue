<template>
  <div class="editor-header">
    <div class="editor-title">
      <div
        v-if="!isEditingTitle"
        class="doc-title"
        :title="titlePlaceholder"
        @click="$emit('start-edit-title')"
      >
        {{ currentDocTitle || titlePlaceholder }}
      </div>
      <el-input
        v-else
        ref="titleInputRef"
        :model-value="pendingTitle"
        class="doc-title-input"
        maxlength="200"
        @update:model-value="$emit('update:pending-title', $event)"
        @blur="$emit('commit-title')"
        @keyup.enter="$emit('commit-title')"
        @keyup.esc="$emit('cancel-edit-title')"
      />
    </div>
    <div class="editor-actions">
      <slot name="title-action" />
      <el-button
        class="editor-action-button"
        text
        :title="isSidebarCollapsed ? expandTitle : collapseTitle"
        @click="$emit('toggle-sidebar')"
      >
        <el-icon>
          <component :is="isSidebarCollapsed ? Expand : Fold" />
        </el-icon>
      </el-button>
      <el-button
        class="editor-action-button"
        text
        :title="commentsTitle"
        @click="$emit('toggle-comment-sidebar')"
      >
        <el-icon><ChatLineRound /></el-icon>
      </el-button>
      <AppDropdown
        trigger="click"
        :teleported="!isFullscreen"
        @command="$emit('header-command', $event)"
      >
        <el-button class="editor-action-button" text>
          <el-icon><MoreFilled /></el-icon>
        </el-button>
        <template #dropdown>
          <el-dropdown-menu>
            <el-dropdown-item command="show_history" :disabled="!canManageAttachments">
              {{ historyLabel }}
            </el-dropdown-item>
            <el-dropdown-item command="manage_attachments" :disabled="!canManageAttachments">
              {{ attachmentsLabel }}
            </el-dropdown-item>
            <el-dropdown-item command="document_settings" :disabled="!canManageSettings">
              {{ settingsLabel }}
            </el-dropdown-item>
            <el-dropdown-item command="create_template">
              {{ saveAsLabel }}
            </el-dropdown-item>
          </el-dropdown-menu>
        </template>
      </AppDropdown>
    </div>
  </div>
</template>

<script setup>
import { nextTick, ref, watch } from 'vue';
import { Expand, Fold, MoreFilled, ChatLineRound } from '@element-plus/icons-vue';
import AppDropdown from '@/components/base/AppDropdown.vue';

const props = defineProps({
  isEditingTitle: { type: Boolean, default: false },
  currentDocTitle: { type: String, default: '' },
  pendingTitle: { type: String, default: '' },
  isSidebarCollapsed: { type: Boolean, default: false },
  titlePlaceholder: { type: String, default: '' },
  collapseTitle: { type: String, default: '' },
  expandTitle: { type: String, default: '' },
  commentsTitle: { type: String, default: '' },
  historyLabel: { type: String, default: '' },
  saveAsLabel: { type: String, default: '' },
  attachmentsLabel: { type: String, default: '' },
  settingsLabel: { type: String, default: '' },
  isFullscreen: { type: Boolean, default: false },
  canManageAttachments: { type: Boolean, default: false },
  canManageSettings: { type: Boolean, default: false }
});

defineEmits([
  'update:pending-title',
  'start-edit-title',
  'commit-title',
  'cancel-edit-title',
  'toggle-sidebar',
  'toggle-comment-sidebar',
  'header-command'
]);

const titleInputRef = ref(null);

watch(
  () => props.isEditingTitle,
  async (editing) => {
    if (!editing) return;
    await nextTick();
    titleInputRef.value?.focus?.();
  }
);
</script>

<style scoped>
.editor-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  min-width: 0;
  padding: var(--spacing-md) var(--spacing-lg);
  border-bottom: 1px solid var(--border-color);
}

.editor-title {
  display: flex;
  align-items: center;
  flex: 1;
  min-width: 0;
  overflow: hidden;
  min-height: 28px;
}

.editor-actions {
  display: flex;
  align-items: center;
  flex-shrink: 0;
  gap: var(--spacing-sm);
}

.editor-action-button {
  color: var(--text-medium);
}

.editor-action-button:hover {
  color: var(--primary-color);
}

.doc-title {
  font-size: 18px;
  font-weight: 600;
  color: var(--text-dark);
  cursor: text;
  max-width: 100%;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.doc-title-input :deep(.el-input__wrapper) {
  box-shadow: none;
  border-radius: 0;
  background-color: transparent;
  padding: 0;
}

.doc-title-input {
  width: min(100%, 560px);
}

.doc-title-input :deep(.el-input__inner) {
  font-size: 18px;
  font-weight: 600;
  color: var(--text-dark);
  padding: 0;
  height: 28px;
  line-height: 28px;
}

.dark-mode .editor-header {
  border-color: var(--border-color);
}

.dark-mode .doc-title {
  color: var(--text-dark);
}

.dark-mode .doc-title-input :deep(.el-input__inner) {
  color: var(--text-dark);
}
</style>
