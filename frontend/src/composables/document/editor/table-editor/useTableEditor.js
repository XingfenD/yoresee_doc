import { nextTick, onBeforeUnmount, onMounted, ref, watch } from 'vue';
import Spreadsheet from 'x-data-spreadsheet';
import {
  buildSheetStyle,
  parseRows,
  serializeRows,
  rowsToSheetData,
  sheetDataToRows
} from './tableData';

export function useTableEditor(editorRef, { modelValue, onModelValue, onCommit }) {
  const sheetRef = ref(null);
  const lastSerialized = ref('');
  const applyingData = ref(false);
  const wheelBlocker = ref(null);
  const resizeHandler = ref(null);
  const rafIdRef = ref(0);
  const transitionContainerRef = ref(null);
  const transitionEndHandlerRef = ref(null);
  const themeObserverRef = ref(null);

  const setBodyTableEditorOpen = (enabled) => {
    if (typeof document === 'undefined') {
      return;
    }
    document.body.classList.toggle('table-editor-open', Boolean(enabled));
  };

  const emitModelValueFromSheet = (sheetData) => {
    const rows = sheetDataToRows(sheetData);
    const serialized = serializeRows(rows);
    if (serialized === lastSerialized.value) {
      return;
    }
    lastSerialized.value = serialized;
    onModelValue(serialized);
  };

  const rerenderSheet = () => {
    const instance = sheetRef.value;
    if (!instance) {
      return;
    }
    instance.sheet?.reload?.();
    instance.reRender?.();
  };

  const scheduleRerender = () => {
    if (rafIdRef.value) {
      cancelAnimationFrame(rafIdRef.value);
    }
    rafIdRef.value = requestAnimationFrame(() => {
      rafIdRef.value = 0;
      rerenderSheet();
    });
  };

  const applyModelValue = async (value) => {
    const rows = parseRows(value);
    lastSerialized.value = serializeRows(rows);
    if (!sheetRef.value) {
      return;
    }
    applyingData.value = true;
    sheetRef.value.loadData(rowsToSheetData(rows));
    await nextTick();
    applyingData.value = false;
  };

  const applyThemeToSheet = () => {
    const instance = sheetRef.value;
    if (!instance?.data?.settings) {
      return;
    }
    const nextStyle = buildSheetStyle();
    const currentStyle = instance.data.settings.style || {};
    instance.data.settings.style = {
      ...currentStyle,
      ...nextStyle,
      font: {
        ...(currentStyle.font || {}),
        ...nextStyle.font
      }
    };
    scheduleRerender();
  };

  const initSpreadsheet = () => {
    if (!editorRef.value || sheetRef.value) {
      return;
    }
    const instance = new Spreadsheet(editorRef.value, {
      mode: 'edit',
      showToolbar: true,
      showGrid: true,
      showContextmenu: true,
      showBottomBar: false,
      row: {
        len: 100,
        height: 28
      },
      col: {
        len: 26,
        width: 120,
        indexWidth: 52,
        minWidth: 72
      },
      view: {
        height: () => editorRef.value?.clientHeight || 640,
        width: () => editorRef.value?.clientWidth || 960
      },
      style: buildSheetStyle()
    });
    sheetRef.value = instance;
    instance.change((data) => {
      if (applyingData.value) {
        return;
      }
      emitModelValueFromSheet(data);
    });
    instance.on('cell-edited', () => {
      onCommit();
    });
  };

  onMounted(async () => {
    setBodyTableEditorOpen(true);
    initSpreadsheet();
    const handler = (event) => {
      const container = editorRef.value;
      if (!container) {
        return;
      }
      if (!container.contains(event.target)) {
        return;
      }
      event.preventDefault();
    };
    window.addEventListener('wheel', handler, { passive: false, capture: true });
    window.addEventListener('mousewheel', handler, { passive: false, capture: true });
    window.addEventListener('DOMMouseScroll', handler, { passive: false, capture: true });
    wheelBlocker.value = handler;
    const onResize = () => rerenderSheet();
    window.addEventListener('resize', onResize);
    resizeHandler.value = onResize;
    const transitionContainer = editorRef.value?.closest('.editor-layout');
    if (transitionContainer) {
      const onTransitionEnd = () => {
        scheduleRerender();
      };
      transitionContainer.addEventListener('transitionend', onTransitionEnd);
      transitionContainerRef.value = transitionContainer;
      transitionEndHandlerRef.value = onTransitionEnd;
    }
    if (typeof MutationObserver !== 'undefined' && typeof document !== 'undefined') {
      const observer = new MutationObserver(() => {
        applyThemeToSheet();
      });
      observer.observe(document.body, { attributes: true, attributeFilter: ['class'] });
      themeObserverRef.value = observer;
    }
    applyThemeToSheet();
    await applyModelValue(modelValue.value);
    scheduleRerender();
  });

  onBeforeUnmount(() => {
    setBodyTableEditorOpen(false);
    if (themeObserverRef.value) {
      themeObserverRef.value.disconnect();
    }
    if (wheelBlocker.value) {
      window.removeEventListener('wheel', wheelBlocker.value, true);
      window.removeEventListener('mousewheel', wheelBlocker.value, true);
      window.removeEventListener('DOMMouseScroll', wheelBlocker.value, true);
    }
    if (sheetRef.value && typeof sheetRef.value.destroy === 'function') {
      sheetRef.value.destroy();
    }
    if (resizeHandler.value) {
      window.removeEventListener('resize', resizeHandler.value);
    }
    if (transitionContainerRef.value && transitionEndHandlerRef.value) {
      transitionContainerRef.value.removeEventListener('transitionend', transitionEndHandlerRef.value);
    }
    if (rafIdRef.value) {
      cancelAnimationFrame(rafIdRef.value);
    }
    rafIdRef.value = 0;
    themeObserverRef.value = null;
    transitionContainerRef.value = null;
    transitionEndHandlerRef.value = null;
    resizeHandler.value = null;
    wheelBlocker.value = null;
    sheetRef.value = null;
  });

  watch(
    () => modelValue.value,
    async (newValue) => {
      if (typeof newValue !== 'string') {
        return;
      }
      if (newValue === lastSerialized.value) {
        return;
      }
      await applyModelValue(newValue);
    }
  );

  return {
    applyModelValue,
    reRender: scheduleRerender
  };
}
