import { computed } from 'vue';

const normalizeColorForPicker = (value, fallback) => {
  const normalized = String(value || '').trim().toLowerCase();
  return /^#([0-9a-f]{3}|[0-9a-f]{6})$/.test(normalized) ? normalized : fallback;
};

export function useRichTextSelectionColors({ editorRef }) {
  const selectionTextColor = computed(() => {
    const value = editorRef.value?.getAttributes('textStyle')?.color;
    return normalizeColorForPicker(value, '#1f2937');
  });

  const selectionBackgroundColor = computed(() => {
    const value = editorRef.value?.getAttributes('highlight')?.color;
    return normalizeColorForPicker(value, '');
  });

  return {
    selectionTextColor,
    selectionBackgroundColor
  };
}
