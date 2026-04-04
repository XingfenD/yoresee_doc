import { computed, onBeforeUnmount, ref, watch } from 'vue';
import { normalizeDocumentType } from '@/utils/documentType';
import { useApiAction } from '@/composables/actions/useApiAction';

export function useTypedDocumentPersistence(options = {}) {
  const {
    type,
    docId,
    currentDocType,
    editorContent,
    t,
    getDocumentContent,
    updateDocument,
    saveContext = 'saveTypedDocument',
    loadContext = 'loadTypedDocument',
    rerender = null
  } = options;

  const normalizedType = normalizeDocumentType(type, '');
  const { runSilent } = useApiAction({ t });

  const isCurrentType = computed(() => normalizeDocumentType(currentDocType.value, '1') === normalizedType);
  const dirty = ref(false);
  const saveTimer = ref(null);
  const saveInFlight = ref(false);
  const saveSeq = ref(0);
  const lastSaved = ref('');
  const skipWatcher = ref(false);
  const loadSeq = ref(0);

  const clearSaveTimer = () => {
    if (!saveTimer.value) {
      return;
    }
    clearTimeout(saveTimer.value);
    saveTimer.value = null;
  };

  const persistContent = async (content, requestSeq) => {
    if (!docId.value || docId.value === 'example' || !isCurrentType.value) {
      return;
    }
    if (saveInFlight.value) {
      scheduleSave();
      return;
    }
    if (content === lastSaved.value) {
      dirty.value = false;
      return;
    }

    saveInFlight.value = true;
    await runSilent(
      () => updateDocument(docId.value, { content }),
      {
        context: saveContext,
        onSuccess: () => {
          if (requestSeq !== saveSeq.value) {
            return;
          }
          lastSaved.value = content;
          dirty.value = false;
        }
      }
    );
    saveInFlight.value = false;
  };

  const scheduleSave = () => {
    if (!isCurrentType.value || !dirty.value) {
      return;
    }
    clearSaveTimer();
    saveTimer.value = setTimeout(() => {
      saveSeq.value += 1;
      const seq = saveSeq.value;
      persistContent(editorContent.value, seq);
    }, 900);
  };

  const flushSave = async () => {
    if (!isCurrentType.value || !dirty.value) {
      return;
    }
    clearSaveTimer();
    saveSeq.value += 1;
    const seq = saveSeq.value;
    await persistContent(editorContent.value, seq);
  };

  const loadContent = async () => {
    if (!docId.value || docId.value === 'example' || !isCurrentType.value) {
      return;
    }
    const nextLoadSeq = ++loadSeq.value;
    const response = await runSilent(
      () => getDocumentContent(docId.value),
      { context: loadContext }
    );
    if (!response || nextLoadSeq !== loadSeq.value) {
      return;
    }
    skipWatcher.value = true;
    editorContent.value = response.content || '';
    lastSaved.value = editorContent.value;
    dirty.value = false;
    skipWatcher.value = false;
  };

  const rerenderEditor = () => {
    if (!isCurrentType.value || typeof rerender !== 'function') {
      return;
    }
    requestAnimationFrame(() => {
      rerender();
    });
  };

  watch(
    () => [docId.value, currentDocType.value],
    async () => {
      clearSaveTimer();
      dirty.value = false;
      lastSaved.value = '';
      if (!isCurrentType.value) {
        return;
      }
      await loadContent();
    },
    { immediate: true }
  );

  watch(
    () => editorContent.value,
    () => {
      if (!isCurrentType.value || skipWatcher.value) {
        return;
      }
      dirty.value = true;
      scheduleSave();
    }
  );

  watch(
    () => isCurrentType.value,
    (enabled) => {
      if (!enabled) {
        return;
      }
      rerenderEditor();
    }
  );

  onBeforeUnmount(() => {
    clearSaveTimer();
  });

  return {
    isCurrentType,
    flushSave,
    rerenderEditor
  };
}
