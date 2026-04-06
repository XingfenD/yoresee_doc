import { nextTick, onBeforeUnmount, onMounted, watch } from 'vue';
import { Editor } from '@tiptap/vue-3';
import StarterKit from '@tiptap/starter-kit';
import CodeBlockLowlight from '@tiptap/extension-code-block-lowlight';
import Placeholder from '@tiptap/extension-placeholder';
import Underline from '@tiptap/extension-underline';
import Link from '@tiptap/extension-link';
import { TextStyle } from '@tiptap/extension-text-style';
import Color from '@tiptap/extension-color';
import Highlight from '@tiptap/extension-highlight';
import { resolveRichTextComponentSystem } from '@/components/document/rich-text/components/registry';
import { CommentAnchorExtension } from '@/components/document/rich-text/extensions/commentAnchorExtension';

export function useRichTextEditorRuntime(options = {}) {
  const {
    editorRef,
    bodyScrollRef,
    modelValueRef,
    valueFormatIsJsonRef,
    placeholderTextRef,
    commentEnabledRef,
    enabledComponentsRef,
    externalExtensionsRef,
    reservedBridgeOptionsRef,
    componentToolbarItemsRef,
    lowlight,
    applyingModelValueRef,
    lastEmittedValueRef,
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
    selectionCommentVisibleRef,
    onUpdateModelValue,
    onCommit,
    onReady
  } = options;

  const buildReservedBridge = () => ({
    ...reservedBridgeOptionsRef.value,
    comments: null,
    collaboration: null
  });

  const handleSelectionOverlayScrollOrResize = () => {
    if (selectionCommentVisibleRef.value) {
      updateSelectionCommentTrigger();
    }
    updateCodeLanguageFloating();
  };

  const handleKeyDown = (view, event) => {
    if (event.key !== 'Tab' || event.metaKey || event.ctrlKey || event.altKey) {
      return false;
    }
    const instance = editorRef.value;
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
  };

  onMounted(() => {
    const initialValue = valueFormatIsJsonRef.value
      ? modelValueRef.value
      : String(modelValueRef.value || '');
    syncLastEmittedValue(initialValue);

    const componentSystem = resolveRichTextComponentSystem(enabledComponentsRef.value);
    componentToolbarItemsRef.value = componentSystem.toolbarItems;

    editorRef.value = new Editor({
      extensions: [
        StarterKit.configure({
          codeBlock: false,
          blockquote: true,
          heading: { levels: [1, 2, 3] }
        }),
        Underline,
        TextStyle,
        Color,
        Highlight.configure({
          multicolor: true
        }),
        Link.configure({
          autolink: true,
          openOnClick: false,
          linkOnPaste: true
        }),
        CodeBlockLowlight.configure({
          lowlight
        }),
        Placeholder.configure({
          placeholder: placeholderTextRef.value
        }),
        ...(commentEnabledRef.value ? [CommentAnchorExtension] : []),
        ...componentSystem.extensions,
        ...externalExtensionsRef.value
      ],
      content: resolveInitialEditorContent(modelValueRef.value),
      editorProps: {
        attributes: {
          class: 'yoresee-rich-text-content'
        },
        handleKeyDown
      },
      onUpdate: ({ editor: instance }) => {
        if (applyingModelValueRef.value) {
          return;
        }
        const nextValue = modelValueFromEditor(instance);
        const serialized = serializeModelValue(nextValue);
        if (serialized === lastEmittedValueRef.value) {
          return;
        }
        lastEmittedValueRef.value = serialized;
        onUpdateModelValue(nextValue);
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
        onCommit();
      }
    });

    onReady({
      editor: editorRef.value,
      reservedBridge: buildReservedBridge()
    });

    syncCodeLanguageDraft();
    updateCodeLanguageFloating();

    nextTick(() => {
      bodyScrollRef.value?.addEventListener('scroll', handleSelectionOverlayScrollOrResize, {
        passive: true
      });
      bodyScrollRef.value?.addEventListener('mouseup', requestSelectionCommentTriggerUpdate);
      bodyScrollRef.value?.addEventListener('keyup', requestSelectionCommentTriggerUpdate);
    });
    document.addEventListener('selectionchange', requestSelectionCommentTriggerUpdate);
    window.addEventListener('resize', handleSelectionOverlayScrollOrResize);
  });

  watch(modelValueRef, (nextValue) => {
    applyModelValueToEditor(nextValue);
    syncCodeLanguageDraft();
    nextTick(() => {
      updateCodeLanguageFloating();
    });
  });

  onBeforeUnmount(() => {
    clearSelectionRaf();
    bodyScrollRef.value?.removeEventListener('mouseup', requestSelectionCommentTriggerUpdate);
    bodyScrollRef.value?.removeEventListener('keyup', requestSelectionCommentTriggerUpdate);
    bodyScrollRef.value?.removeEventListener('scroll', handleSelectionOverlayScrollOrResize);
    document.removeEventListener('selectionchange', requestSelectionCommentTriggerUpdate);
    window.removeEventListener('resize', handleSelectionOverlayScrollOrResize);
    editorRef.value?.destroy();
    editorRef.value = null;
  });
}
