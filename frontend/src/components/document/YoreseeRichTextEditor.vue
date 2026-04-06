<template>
  <div class="rich-text-editor">
    <header class="rich-text-toolbar">
      <button type="button" class="toolbar-btn" :class="{ 'is-active': isActive('bold') }" @click="toggleBold">
        B
      </button>
      <button type="button" class="toolbar-btn" :class="{ 'is-active': isActive('italic') }" @click="toggleItalic">
        I
      </button>
      <button type="button" class="toolbar-btn" :class="{ 'is-active': isActive('strike') }" @click="toggleStrike">
        S
      </button>
      <span class="toolbar-divider"></span>
      <button
        type="button"
        class="toolbar-btn"
        :class="{ 'is-active': isHeadingActive(1) }"
        @click="toggleHeading(1)"
      >
        H1
      </button>
      <button
        type="button"
        class="toolbar-btn"
        :class="{ 'is-active': isHeadingActive(2) }"
        @click="toggleHeading(2)"
      >
        H2
      </button>
      <button
        type="button"
        class="toolbar-btn"
        :class="{ 'is-active': isActive('bulletList') }"
        @click="toggleBulletList"
      >
        UL
      </button>
      <button
        type="button"
        class="toolbar-btn"
        :class="{ 'is-active': isActive('orderedList') }"
        @click="toggleOrderedList"
      >
        OL
      </button>
      <button
        type="button"
        class="toolbar-btn"
        :class="{ 'is-active': isActive('blockquote') }"
        @click="toggleBlockquote"
      >
        Quote
      </button>
      <button
        type="button"
        class="toolbar-btn"
        :class="{ 'is-active': isActive('codeBlock') }"
        @click="toggleCodeBlock"
      >
        Code
      </button>
      <button
        v-if="commentEnabled"
        type="button"
        class="toolbar-btn"
        @mousedown.prevent
        @click="addInlineComment"
      >
        Comment
      </button>
      <button
        v-for="item in componentToolbarItems"
        :key="item.key"
        type="button"
        class="toolbar-btn"
        @click="runComponentCommand(item)"
      >
        {{ item.label }}
      </button>
      <span class="toolbar-spacer"></span>
      <button type="button" class="toolbar-btn" :disabled="!canUndo" @click="undo">
        Undo
      </button>
      <button type="button" class="toolbar-btn" :disabled="!canRedo" @click="redo">
        Redo
      </button>
    </header>

    <div v-if="editor" ref="bodyScrollRef" class="rich-text-body">
      <div
        v-show="codeLanguageFloatingVisible"
        class="code-language-floating"
        :style="codeLanguageFloatingStyle"
        @mouseenter="codeLanguageInteracting = true"
        @mouseleave="codeLanguageInteracting = false"
      >
        <el-autocomplete
          v-model="codeLanguageDraft"
          class="code-language-input"
          :fetch-suggestions="queryCodeLanguageSuggestions"
          :trigger-on-focus="true"
          :debounce="0"
          value-key="label"
          placeholder="语言"
          @focus="codeLanguageInteracting = true"
          @blur="handleCodeLanguageInputBlur"
          @change="applyCodeLanguageFromDraft"
          @select="handleCodeLanguageSelect"
        >
          <template #default="{ item }">
            <div class="code-language-option">
              <span class="code-language-option-label">{{ item.label }}</span>
              <span class="code-language-option-value">{{ item.value }}</span>
            </div>
          </template>
        </el-autocomplete>
      </div>
      <button
        v-show="selectionCommentVisible"
        type="button"
        class="selection-comment-trigger"
        :style="selectionCommentStyle"
        :title="t('document.comments')"
        @mousedown.prevent
        @mouseenter="selectionCommentHovering = true"
        @mouseleave="selectionCommentHovering = false"
        @click="handleSelectionCommentClick"
      >
        <el-icon>
          <ChatDotRound />
        </el-icon>
      </button>
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
import { computed, nextTick, onBeforeUnmount, onMounted, ref, watch } from 'vue';
import { Editor, EditorContent } from '@tiptap/vue-3';
import StarterKit from '@tiptap/starter-kit';
import CodeBlockLowlight from '@tiptap/extension-code-block-lowlight';
import Placeholder from '@tiptap/extension-placeholder';
import { ChatDotRound } from '@element-plus/icons-vue';
import { useI18n } from 'vue-i18n';
import RichTextParagraphHandle from '@/components/document/rich-text/RichTextParagraphHandle.vue';
import { resolveRichTextComponentSystem } from '@/components/document/rich-text/components/registry';
import { CommentAnchorExtension } from '@/components/document/rich-text/extensions/commentAnchorExtension';
import { useRichTextParagraphHandle } from '@/composables/document/editor/rich-text-editor/useRichTextParagraphHandle';
import { useRichTextValueBridge } from '@/composables/document/editor/rich-text-editor/useRichTextValueBridge';
import { useRichTextCodeLanguageOverlay } from '@/composables/document/editor/rich-text-editor/useRichTextCodeLanguageOverlay';
import { useRichTextCommentBridge } from '@/composables/document/editor/rich-text-editor/useRichTextCommentBridge';

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
const canUndo = computed(() => editor.value?.can().chain().focus().undo().run() ?? false);
const canRedo = computed(() => editor.value?.can().chain().focus().redo().run() ?? false);

const runCommand = (command) => {
  if (!editor.value) {
    return;
  }
  command(editor.value.chain().focus()).run();
};

const isActive = (name) => editor.value?.isActive(name) ?? false;
const isHeadingActive = (level) => editor.value?.isActive('heading', { level }) ?? false;

const toggleBold = () => runCommand((chain) => chain.toggleBold());
const toggleItalic = () => runCommand((chain) => chain.toggleItalic());
const toggleStrike = () => runCommand((chain) => chain.toggleStrike());
const toggleHeading = (level) => runCommand((chain) => chain.toggleHeading({ level }));
const toggleBulletList = () => runCommand((chain) => chain.toggleBulletList());
const toggleOrderedList = () => runCommand((chain) => chain.toggleOrderedList());
const toggleBlockquote = () => runCommand((chain) => chain.toggleBlockquote());
const toggleCodeBlock = () => {
  if (!editor.value) {
    return;
  }
  const chain = editor.value.chain().focus();
  if (editor.value.isActive('codeBlock')) {
    chain.toggleCodeBlock().run();
    hideCodeLanguageFloating();
    return;
  }
  const nextLanguage = resolveCodeBlockLanguage(codeLanguageDraft.value);
  codeLanguageDraft.value = nextLanguage;
  chain.setCodeBlock({ language: nextLanguage }).run();
  updateCodeLanguageFloating();
};
const undo = () => runCommand((chain) => chain.undo());
const redo = () => runCommand((chain) => chain.redo());

const runComponentCommand = (item) => {
  if (!editor.value || !item?.command) {
    return;
  }
  const chain = editor.value.chain().focus();
  const commandMethod = chain[item.command];
  if (typeof commandMethod !== 'function') {
    return;
  }
  commandMethod.call(chain, item.commandArgs || {});
  chain.run();
};

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

const handleSelectionCommentClick = () => {
  addInlineComment();
};

const handleSelectionOverlayScrollOrResize = () => {
  if (selectionCommentVisible.value) {
    updateSelectionCommentTrigger();
  }
  updateCodeLanguageFloating();
};

const buildReservedBridge = () => ({
  ...props.reservedBridgeOptions,
  comments: null,
  collaboration: null
});

onMounted(() => {
  const initialValue = isJsonValueMode.value ? props.modelValue : String(props.modelValue || '');
  syncLastEmittedValue(initialValue);
  const componentSystem = resolveRichTextComponentSystem(props.enabledComponents);
  componentToolbarItems.value = componentSystem.toolbarItems;

  editor.value = new Editor({
    extensions: [
      StarterKit.configure({
        codeBlock: false,
        blockquote: true,
        heading: { levels: [1, 2, 3] }
      }),
      CodeBlockLowlight.configure({
        lowlight
      }),
      Placeholder.configure({
        placeholder: placeholderText.value
      }),
      ...(props.commentEnabled ? [CommentAnchorExtension] : []),
      ...componentSystem.extensions,
      ...props.externalExtensions
    ],
    content: resolveInitialEditorContent(props.modelValue),
    editorProps: {
      attributes: {
        class: 'yoresee-rich-text-content'
      },
      handleKeyDown: (view, event) => {
        if (event.key !== 'Tab' || event.metaKey || event.ctrlKey || event.altKey) {
          return false;
        }
        const instance = editor.value;
        if (instance?.isActive('listItem')) {
          event.preventDefault();
          if (event.shiftKey) {
            instance.chain().focus().liftListItem('listItem').run();
          } else {
            instance.chain().focus().sinkListItem('listItem').run();
          }
          return true;
        }
        if (event.shiftKey) {
          return false;
        }
        event.preventDefault();
        view.dispatch(view.state.tr.insertText('    '));
        return true;
      }
    },
    onUpdate: ({ editor: instance }) => {
      if (applyingModelValue.value) {
        return;
      }
      const nextValue = modelValueFromEditor(instance);
      const serialized = serializeModelValue(nextValue);
      if (serialized === lastEmittedValue.value) {
        return;
      }
      lastEmittedValue.value = serialized;
      emit('update:modelValue', nextValue);
      syncCodeLanguageDraft();
      updateCodeLanguageFloating();
      updateSelectionCommentTrigger();
    },
    onSelectionUpdate: () => {
      syncCodeLanguageDraft();
      updateCodeLanguageFloating();
      updateSelectionCommentTrigger();
    },
    onBlur: () => {
      handleSelectionCommentEditorBlur();
      handleCodeLanguageEditorBlur();
      emit('commit');
    }
  });

  emit('ready', {
    editor: editor.value,
    reservedBridge: buildReservedBridge()
  });

  syncCodeLanguageDraft();
  updateCodeLanguageFloating();

  nextTick(() => {
    bodyScrollRef.value?.addEventListener('scroll', handleSelectionOverlayScrollOrResize, { passive: true });
    bodyScrollRef.value?.addEventListener('mouseup', requestSelectionCommentTriggerUpdate);
    bodyScrollRef.value?.addEventListener('keyup', requestSelectionCommentTriggerUpdate);
  });
  document.addEventListener('selectionchange', requestSelectionCommentTriggerUpdate);
  window.addEventListener('resize', handleSelectionOverlayScrollOrResize);
});

watch(
  () => props.modelValue,
  (nextValue) => {
    applyModelValueToEditor(nextValue);
    syncCodeLanguageDraft();
    nextTick(() => {
      updateCodeLanguageFloating();
    });
  }
);

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

onBeforeUnmount(() => {
  clearSelectionRaf();
  bodyScrollRef.value?.removeEventListener('mouseup', requestSelectionCommentTriggerUpdate);
  bodyScrollRef.value?.removeEventListener('keyup', requestSelectionCommentTriggerUpdate);
  bodyScrollRef.value?.removeEventListener('scroll', handleSelectionOverlayScrollOrResize);
  document.removeEventListener('selectionchange', requestSelectionCommentTriggerUpdate);
  window.removeEventListener('resize', handleSelectionOverlayScrollOrResize);
  editor.value?.destroy();
  editor.value = null;
});
</script>

<style scoped src="./rich-text/yoreseeRichTextEditor.css"></style>
