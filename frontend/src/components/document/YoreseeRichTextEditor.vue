<template>
  <div class="rich-text-editor">
    <div v-if="editor" ref="bodyScrollRef" class="rich-text-body">
      <RichTextCodeLanguageFloating
        v-model="codeLanguageDraft"
        :visible="codeLanguageFloatingVisible"
        :style-object="codeLanguageFloatingStyle"
        :query-suggestions="queryCodeLanguageSuggestions"
        @enter="codeLanguageInteracting = true"
        @leave="codeLanguageInteracting = false"
        @focus="codeLanguageInteracting = true"
        @blur="handleCodeLanguageInputBlur"
        @change="applyCodeLanguageFromDraft"
        @select="handleCodeLanguageSelect"
      />
      <RichTextSelectionActionBubble
        :visible="selectionCommentVisible"
        :style-object="selectionCommentStyle"
        :comment-enabled="commentEnabled"
        :comment-title="t('document.comments')"
        :selection-bold-active="isActive('bold')"
        :selection-underline-active="isActive('underline')"
        :selection-italic-active="isActive('italic')"
        :selection-strike-active="isActive('strike')"
        :selection-inline-code-active="isActive('code')"
        :selection-link-active="isActive('link')"
        :selection-highlight-active="isActive('highlight')"
        :selection-text-color="selectionTextColor"
        :selection-background-color="selectionBackgroundColor"
        @enter="selectionCommentHovering = true"
        @leave="selectionCommentHovering = false"
        @toggle-bold="toggleBold"
        @toggle-underline="toggleUnderline"
        @toggle-italic="toggleItalic"
        @toggle-strike="toggleStrike"
        @toggle-inline-code="toggleInlineCode"
        @link-click="openSelectionLinkDialog"
        @text-color-change="setTextColor"
        @background-color-change="setBackgroundColor"
        @comment-click="handleSelectionCommentClick"
      />
      <RichTextLinkHoverCard
        :visible="linkHoverVisible"
        :style-object="linkHoverStyle"
        :href="linkHoverHref"
        @enter="handleLinkCardEnter"
        @leave="handleLinkCardLeave"
        @edit="openHoverLinkDialog"
        @remove="removeHoverLink"
      />
      <RichTextLinkDialog
        v-model="linkDialogValue"
        :visible="linkDialogVisible"
        @cancel="closeLinkDialog"
        @confirm="applyLinkDialog"
      />
      <RichTextParagraphHandle
        :visible="paragraphHandleVisible"
        :style="paragraphHandleStyle"
        :mode="paragraphHandleMode"
        :block-type="paragraphType"
        :actions="paragraphActions"
        plus-title="新增段落"
        more-title="更多操作"
        :type-label="paragraphTypeLabel"
        @action="runParagraphAction"
        @mouseenter="handleHandleEnter"
        @mouseleave="handleHandleLeave"
      />
      <EditorContent :editor="editor" class="rich-text-content-host" />
    </div>
    <div v-else class="rich-text-placeholder">{{ placeholderText }}</div>

    <footer class="rich-text-hint">
      {{ t('document.richTextShortcutHint') }}
    </footer>
  </div>
</template>

<script setup>
import { computed, ref } from 'vue';
import { EditorContent } from '@tiptap/vue-3';
import { useI18n } from 'vue-i18n';
import RichTextParagraphHandle from '@/components/document/rich-text/RichTextParagraphHandle.vue';
import RichTextCodeLanguageFloating from '@/components/document/rich-text/overlay/RichTextCodeLanguageFloating.vue';
import RichTextSelectionActionBubble from '@/components/document/rich-text/overlay/RichTextSelectionActionBubble.vue';
import RichTextLinkHoverCard from '@/components/document/rich-text/overlay/RichTextLinkHoverCard.vue';
import RichTextLinkDialog from '@/components/document/rich-text/overlay/RichTextLinkDialog.vue';
import { useRichTextParagraphHandle } from '@/composables/document/editor/rich-text-editor/useRichTextParagraphHandle';
import { useRichTextValueBridge } from '@/composables/document/editor/rich-text-editor/useRichTextValueBridge';
import { useRichTextCodeLanguageOverlay } from '@/composables/document/editor/rich-text-editor/useRichTextCodeLanguageOverlay';
import { useRichTextCommentBridge } from '@/composables/document/editor/rich-text-editor/useRichTextCommentBridge';
import { useRichTextLinkOverlay } from '@/composables/document/editor/rich-text-editor/useRichTextLinkOverlay';
import { useRichTextToolbarActions } from '@/composables/document/editor/rich-text-editor/useRichTextToolbarActions';
import { useRichTextEditorRuntime } from '@/composables/document/editor/rich-text-editor/useRichTextEditorRuntime';
import { useRichTextSelectionColors } from '@/composables/document/editor/rich-text-editor/useRichTextSelectionColors';
import { useRichTextParagraphActions } from '@/composables/document/editor/rich-text-editor/useRichTextParagraphActions';

const props = defineProps({
  modelValue: {
    type: [String, Object],
    default: ''
  },
  valueFormat: {
    type: String,
    default: 'markdown',
    validator: (value) => ['markdown', 'json'].includes(value)
  },
  placeholder: {
    type: String,
    default: ''
  },
  externalExtensions: {
    type: Array,
    default: () => []
  },
  enabledComponents: {
    type: Array,
    default: () => ['mindmap', 'drawio']
  },
  commentEnabled: {
    type: Boolean,
    default: false
  },
  reservedBridgeOptions: {
    type: Object,
    default: () => ({})
  }
});

const emit = defineEmits([
  'update:modelValue',
  'commit',
  'ready',
  'comment-add',
  'comment-remove',
  'comment-changed'
]);

const { t } = useI18n();
const editor = ref(null);
const bodyScrollRef = ref(null);
const componentToolbarItems = ref([]);

const valueFormatRef = computed(() => props.valueFormat);
const modelValueRef = computed(() => props.modelValue);
const enabledComponentsRef = computed(() => props.enabledComponents);
const externalExtensionsRef = computed(() => props.externalExtensions);
const reservedBridgeOptionsRef = computed(() => props.reservedBridgeOptions);
const {
  applyingModelValue,
  isJsonValueMode,
  lastEmittedValue,
  serializeModelValue,
  modelValueFromEditor,
  resolveInitialEditorContent,
  applyModelValueToEditor,
  syncLastEmittedValue
} = useRichTextValueBridge({
  editorRef: editor,
  valueFormatRef
});

const {
  lowlight,
  codeLanguageDraft,
  codeLanguageFloatingVisible,
  codeLanguageFloatingStyle,
  codeLanguageInteracting,
  queryCodeLanguageSuggestions,
  applyCodeLanguageFromDraft,
  handleCodeLanguageSelect,
  handleCodeLanguageInputBlur,
  syncCodeLanguageDraft,
  updateCodeLanguageFloating,
  hideCodeLanguageFloating,
  handleEditorBlur: handleCodeLanguageEditorBlur,
  resolveCodeBlockLanguage
} = useRichTextCodeLanguageOverlay({
  editorRef: editor,
  scrollContainerRef: bodyScrollRef
});

const commentEnabledRef = computed(() => props.commentEnabled);
const {
  selectionCommentVisible,
  selectionCommentStyle,
  selectionCommentHovering,
  addInlineComment,
  updateSelectionCommentTrigger,
  requestSelectionCommentTriggerUpdate,
  handleEditorBlur: handleSelectionCommentEditorBlur,
  clearSelectionRaf,
  getCommentIds,
  hlCommentIds,
  unHlCommentIds,
  removeCommentIds,
  scrollToCommentId
} = useRichTextCommentBridge({
  editorRef: editor,
  scrollContainerRef: bodyScrollRef,
  commentEnabledRef,
  onCommentAdd: (payload) => emit('comment-add', payload),
  onCommentRemove: (ids) => emit('comment-remove', ids),
  onCommentChanged: () => emit('comment-changed')
});

const placeholderText = computed(() => props.placeholder || t('document.editorPlaceholder'));

const {
  isActive,
  toggleBold,
  toggleUnderline,
  toggleItalic,
  toggleStrike,
  toggleInlineCode,
  toggleBulletList,
  toggleOrderedList,
  toggleBlockquote,
  toggleCodeBlock,
  setLink,
  unsetLink,
  setTextColor,
  setBackgroundColor,
  undo,
  redo,
  runComponentCommand
} = useRichTextToolbarActions({
  editorRef: editor,
  codeLanguageDraftRef: codeLanguageDraft,
  resolveCodeBlockLanguage,
  hideCodeLanguageFloating,
  updateCodeLanguageFloating
});

const {
  linkDialogVisible,
  linkDialogValue,
  linkHoverVisible,
  linkHoverStyle,
  linkHoverHref,
  openSelectionLinkDialog,
  openHoverLinkDialog,
  closeLinkDialog,
  applyLinkDialog,
  removeHoverLink,
  handleLinkCardEnter,
  handleLinkCardLeave
} = useRichTextLinkOverlay({
  editorRef: editor,
  scrollContainerRef: bodyScrollRef,
  setLink,
  unsetLink
});

const { emptyParagraphActions } = useRichTextParagraphActions({
  editorRef: editor,
  componentToolbarItemsRef: componentToolbarItems,
  runComponentCommand,
  toggleBulletList,
  toggleOrderedList,
  toggleBlockquote,
  toggleCodeBlock,
  undo,
  redo
});

const {
  paragraphHandleVisible,
  paragraphHandleStyle,
  paragraphHandleMode,
  paragraphType,
  paragraphActions,
  runParagraphAction,
  handleHandleEnter,
  handleHandleLeave
} = useRichTextParagraphHandle({
  editorRef: editor,
  scrollContainerRef: bodyScrollRef,
  labels: {
    addAbove: '上方插入空段落',
    addBelow: '下方插入空段落',
    delete: t('document.delete')
  },
  resolveActions: ({ isEmpty, defaults }) => {
    if (!isEmpty) {
      return defaults;
    }
    return [...emptyParagraphActions.value, ...defaults];
  },
  onMutated: () => {
    emit('comment-changed');
  }
});

const paragraphTypeLabelMap = {
  paragraph: '段落',
  heading: '标题',
  list: '列表',
  quote: '引用',
  code: '代码块',
  table: '表格',
  divider: '分割线',
  mindmap: 'Mindmap',
  drawio: 'Draw.io'
};
const paragraphTypeLabel = computed(() => paragraphTypeLabelMap[paragraphType.value] || '段落');
const { selectionTextColor, selectionBackgroundColor } = useRichTextSelectionColors({
  editorRef: editor
});

const handleSelectionCommentClick = () => {
  addInlineComment();
};

useRichTextEditorRuntime({
  editorRef: editor,
  bodyScrollRef,
  modelValueRef,
  valueFormatIsJsonRef: isJsonValueMode,
  placeholderTextRef: placeholderText,
  commentEnabledRef,
  enabledComponentsRef,
  externalExtensionsRef,
  reservedBridgeOptionsRef,
  componentToolbarItemsRef: componentToolbarItems,
  lowlight,
  applyingModelValueRef: applyingModelValue,
  lastEmittedValueRef: lastEmittedValue,
  serializeModelValue,
  modelValueFromEditor,
  resolveInitialEditorContent,
  applyModelValueToEditor,
  syncLastEmittedValue,
  syncCodeLanguageDraft,
  updateCodeLanguageFloating,
  updateSelectionCommentTrigger,
  handleSelectionCommentEditorBlur,
  handleCodeLanguageEditorBlur,
  requestSelectionCommentTriggerUpdate,
  clearSelectionRaf,
  selectionCommentVisibleRef: selectionCommentVisible,
  onUpdateModelValue: (value) => emit('update:modelValue', value),
  onCommit: () => emit('commit'),
  onReady: (payload) => emit('ready', payload)
});

defineExpose({
  reRender: () => {
    const instance = editor.value;
    if (!instance?.view) {
      return;
    }
    instance.view.updateState(instance.state);
  },
  getEditor: () => editor.value,
  getCommentIds,
  hlCommentIds,
  unHlCommentIds,
  removeCommentIds,
  scrollToCommentId,
  broadcastCommentChange: () => {
    emit('comment-changed');
  }
});
</script>

<style scoped src="./rich-text/yoreseeRichTextEditor.css"></style>
