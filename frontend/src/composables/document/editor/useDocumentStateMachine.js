import { computed, ref } from 'vue';

export const DOC_STATE = Object.freeze({
  IDLE: 'idle',
  NAVIGATING: 'navigating',
  RESOLVING: 'resolving',
  READY: 'ready',
  SAVING: 'saving',
  ERROR: 'error'
});

export function useDocumentStateMachine() {
  const state = ref(DOC_STATE.IDLE);
  const targetDocId = ref('');
  const targetDocType = ref('1');
  const targetContent = ref('');
  const error = ref(null);
  const collabSynced = ref(true);

  const isIdle = computed(() => state.value === DOC_STATE.IDLE);
  const isNavigating = computed(() => state.value === DOC_STATE.NAVIGATING);
  const isResolving = computed(() => state.value === DOC_STATE.RESOLVING);
  const isReady = computed(() => state.value === DOC_STATE.READY);
  const isSaving = computed(() => state.value === DOC_STATE.SAVING);
  const isError = computed(() => state.value === DOC_STATE.ERROR);

  const startNavigation = (docId) => {
    targetDocId.value = docId || '';
    targetDocType.value = '1';
    targetContent.value = '';
    error.value = null;
    collabSynced.value = true;
    state.value = DOC_STATE.NAVIGATING;
  };

  const resolveMetadata = ({ docId, type }) => {
    if (docId !== targetDocId.value) {
      return;
    }
    targetDocType.value = type || '1';
    state.value = DOC_STATE.RESOLVING;
  };

  const resolveContent = ({ docId, content }) => {
    if (docId !== targetDocId.value) {
      return;
    }
    targetContent.value = content || '';
    if (state.value === DOC_STATE.RESOLVING || state.value === DOC_STATE.NAVIGATING) {
      state.value = DOC_STATE.READY;
    }
  };

  const markError = (err) => {
    error.value = err || new Error('document load failed');
    state.value = DOC_STATE.ERROR;
  };

  const startSave = () => {
    if (state.value === DOC_STATE.READY) {
      state.value = DOC_STATE.SAVING;
    }
  };

  const finishSave = () => {
    if (state.value === DOC_STATE.SAVING) {
      state.value = DOC_STATE.READY;
    }
  };

  const markCollabSynced = (isSynced, docId) => {
    if (docId !== targetDocId.value) {
      return;
    }
    collabSynced.value = isSynced;
  };

  return {
    state,
    targetDocId,
    targetDocType,
    targetContent,
    error,
    isIdle,
    isNavigating,
    isResolving,
    isReady,
    isSaving,
    isError,
    collabSynced,
    startNavigation,
    resolveMetadata,
    resolveContent,
    markError,
    startSave,
    finishSave,
    markCollabSynced
  };
}
