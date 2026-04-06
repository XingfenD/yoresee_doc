import { computed, onBeforeUnmount, onMounted, ref, watch } from 'vue';

const DEFAULT_TEXT_COLOR = '#1f2937';
const DEFAULT_BACKGROUND_COLOR = '';

const textColorItems = [
  '#1f2937',
  '#9ca3af',
  '#ef4444',
  '#f97316',
  '#eab308',
  '#22c55e',
  '#2563eb',
  '#7c3aed'
];

const backgroundColorItems = [
  '',
  '#d1d5db',
  '#f0b4b4',
  '#efcfaa',
  '#efe687',
  '#b7dfb3',
  '#b3c2e3',
  '#c7b5e8',
  '#d1d5db',
  '#aeb3bc',
  '#f26767',
  '#f8a231',
  '#f7df1e',
  '#59c24a',
  '#8ea7e1',
  '#ad90e0'
];

const normalizeHexColor = (value) => {
  const normalized = String(value || '').trim().toLowerCase();
  return /^#([0-9a-f]{3}|[0-9a-f]{6})$/.test(normalized) ? normalized : '';
};

export function useRichTextSelectionColorPanel({
  visibleRef,
  selectionTextColorRef,
  selectionBackgroundColorRef,
  onTextColorChange,
  onBackgroundColorChange
}) {
  const colorPanelVisible = ref(false);
  const colorPanelRef = ref(null);
  const colorTriggerRef = ref(null);

  const normalizedTextColor = computed(
    () => normalizeHexColor(selectionTextColorRef.value) || DEFAULT_TEXT_COLOR
  );
  const normalizedBackgroundColor = computed(
    () => normalizeHexColor(selectionBackgroundColorRef.value) || DEFAULT_BACKGROUND_COLOR
  );

  const colorTriggerAStyle = computed(() => ({
    color: normalizedTextColor.value,
    backgroundColor: normalizedBackgroundColor.value || 'transparent'
  }));

  const pickTextColor = (color) => {
    onTextColorChange(color);
  };

  const pickBackgroundColor = (color) => {
    onBackgroundColorChange(color || '');
  };

  const isTextColorActive = (color) => normalizedTextColor.value === normalizeHexColor(color);
  const isBackgroundColorActive = (color) => normalizedBackgroundColor.value === normalizeHexColor(color);

  const resetColorToDefault = () => {
    onTextColorChange(DEFAULT_TEXT_COLOR);
    onBackgroundColorChange(DEFAULT_BACKGROUND_COLOR);
  };

  const toggleColorPanel = () => {
    colorPanelVisible.value = !colorPanelVisible.value;
  };

  const closeColorPanel = () => {
    colorPanelVisible.value = false;
  };

  const handleDocumentPointerDown = (event) => {
    if (!colorPanelVisible.value) {
      return;
    }
    const target = event.target;
    if (colorPanelRef.value?.contains(target) || colorTriggerRef.value?.contains(target)) {
      return;
    }
    closeColorPanel();
  };

  onMounted(() => {
    window.addEventListener('mousedown', handleDocumentPointerDown, true);
  });

  onBeforeUnmount(() => {
    window.removeEventListener('mousedown', handleDocumentPointerDown, true);
  });

  watch(visibleRef, (next) => {
    if (!next) {
      closeColorPanel();
    }
  });

  return {
    textColorItems,
    backgroundColorItems,
    colorPanelVisible,
    colorPanelRef,
    colorTriggerRef,
    colorTriggerAStyle,
    pickTextColor,
    pickBackgroundColor,
    isTextColorActive,
    isBackgroundColorActive,
    resetColorToDefault,
    toggleColorPanel
  };
}
