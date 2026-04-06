import { computed, ref } from 'vue';
import { createRichTextTurndown, htmlToMarkdown, markdownToHtml } from './value-bridge/markdownBridge';
import { cloneJsonDoc, normalizeJsonDoc, serializeJsonDoc } from './value-bridge/jsonDocBridge';

export function useRichTextValueBridge(options = {}) {
  const { editorRef, valueFormatRef } = options;
  const applyingModelValue = ref(false);
  const lastEmittedValue = ref('');
  const isJsonValueMode = computed(() => valueFormatRef.value === 'json');
  const turndown = createRichTextTurndown();

  const serializeModelValue = (value) => {
    if (isJsonValueMode.value) {
      return serializeJsonDoc(value);
    }
    return String(value || '');
  };
  const modelValueFromEditor = (instance) => {
    if (isJsonValueMode.value) {
      return cloneJsonDoc(instance.getJSON());
    }
    return htmlToMarkdown(instance.getHTML(), turndown);
  };
  const resolveInitialEditorContent = (modelValue) => {
    if (isJsonValueMode.value) {
      return normalizeJsonDoc(modelValue);
    }
    return markdownToHtml(String(modelValue || ''));
  };
  const applyModelValueToEditor = (modelValue) => {
    if (!editorRef.value) {
      return;
    }
    const serialized = serializeModelValue(modelValue);
    if (serialized === lastEmittedValue.value) {
      return;
    }
    applyingModelValue.value = true;
    if (isJsonValueMode.value) {
      editorRef.value.commands.setContent(normalizeJsonDoc(modelValue), false);
    } else {
      editorRef.value.commands.setContent(markdownToHtml(String(modelValue || '')), false, {
        preserveWhitespace: true
      });
    }
    lastEmittedValue.value = serialized;
    applyingModelValue.value = false;
  };
  const syncLastEmittedValue = (value) => {
    lastEmittedValue.value = serializeModelValue(value);
  };

  return {
    applyingModelValue,
    isJsonValueMode,
    lastEmittedValue,
    serializeModelValue,
    modelValueFromEditor,
    resolveInitialEditorContent,
    applyModelValueToEditor,
    syncLastEmittedValue
  };
}
