import { computed, onBeforeUnmount, ref, watch } from 'vue';
import { normalizeDocumentType } from '@/utils/documentType';
import { useApiAction } from '@/composables/actions/useApiAction';

export function useTableDocumentPersistence(options = {}) {
  const {
    docId,
    currentDocType,
    editorContent,
    tableEditorRef,
    t,
    getDocumentContent,
    updateDocument
  } = options;
  const { runSilent } = useApiAction({ t });

  const isTableDocument = computed(() => normalizeDocumentType(currentDocType.value, '1') === '2');
  const tableDirty = ref(false);
  const tableSaveTimer = ref(null);
  const tableSaveInFlight = ref(false);
  const tableSaveSeq = ref(0);
  const tableLastSaved = ref('');
  const tableSkipWatcher = ref(false);
  const tableLoadSeq = ref(0);

  const clearTableSaveTimer = () => {
    if (tableSaveTimer.value) {
      clearTimeout(tableSaveTimer.value);
      tableSaveTimer.value = null;
    }
  };

  const persistTableContent = async (content, requestSeq) => {
    if (!docId.value || docId.value === 'example' || !isTableDocument.value) {
      return;
    }
    if (tableSaveInFlight.value) {
      scheduleTableSave();
      return;
    }
    if (content === tableLastSaved.value) {
      tableDirty.value = false;
      return;
    }
    tableSaveInFlight.value = true;
    await runSilent(
      () => updateDocument(docId.value, { content }),
      {
        context: 'saveTableDocument',
        onSuccess: () => {
          if (requestSeq !== tableSaveSeq.value) {
            return;
          }
          tableLastSaved.value = content;
          tableDirty.value = false;
        }
      }
    );
    tableSaveInFlight.value = false;
  };

  const scheduleTableSave = () => {
    if (!isTableDocument.value || !tableDirty.value) {
      return;
    }
    clearTableSaveTimer();
    tableSaveTimer.value = setTimeout(() => {
      tableSaveSeq.value += 1;
      const seq = tableSaveSeq.value;
      persistTableContent(editorContent.value, seq);
    }, 900);
  };

  const flushTableSave = async () => {
    if (!isTableDocument.value || !tableDirty.value) {
      return;
    }
    clearTableSaveTimer();
    tableSaveSeq.value += 1;
    const seq = tableSaveSeq.value;
    await persistTableContent(editorContent.value, seq);
  };

  const loadTableContent = async () => {
    if (!docId.value || docId.value === 'example' || !isTableDocument.value) {
      return;
    }
    const loadSeq = ++tableLoadSeq.value;
    const response = await runSilent(
      () => getDocumentContent(docId.value),
      { context: 'loadTableDocument' }
    );
    if (!response || loadSeq !== tableLoadSeq.value) {
      return;
    }
    tableSkipWatcher.value = true;
    editorContent.value = response.content || '';
    tableLastSaved.value = editorContent.value;
    tableDirty.value = false;
    tableSkipWatcher.value = false;
  };

  const rerenderTableEditor = () => {
    if (!isTableDocument.value) {
      return;
    }
    requestAnimationFrame(() => {
      tableEditorRef.value?.reRender?.();
    });
  };

  watch(
    () => [docId.value, currentDocType.value],
    async () => {
      clearTableSaveTimer();
      tableDirty.value = false;
      tableLastSaved.value = '';
      if (!isTableDocument.value) {
        return;
      }
      await loadTableContent();
    },
    { immediate: true }
  );

  watch(
    () => editorContent.value,
    () => {
      if (!isTableDocument.value || tableSkipWatcher.value) {
        return;
      }
      tableDirty.value = true;
      scheduleTableSave();
    }
  );

  watch(
    () => isTableDocument.value,
    (enabled) => {
      if (!enabled) {
        return;
      }
      rerenderTableEditor();
    }
  );

  onBeforeUnmount(() => {
    clearTableSaveTimer();
  });

  return {
    isTableDocument,
    flushTableSave,
    rerenderTableEditor
  };
}
