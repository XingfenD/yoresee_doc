import { useTypedDocumentPersistence } from '@/composables/document/editor/shared/useTypedDocumentPersistence';

export function useTableDocumentPersistence(options = {}) {
  const {
    docId,
    currentDocType,
    editorContent,
    docMachine,
    tableEditorRef,
    t,
    getDocumentContent,
    updateDocument
  } = options;

  const {
    isCurrentType,
    flushSave,
    rerenderEditor
  } = useTypedDocumentPersistence({
    type: '2',
    docId,
    currentDocType,
    editorContent,
    docMachine,
    t,
    getDocumentContent,
    updateDocument,
    saveContext: 'saveTableDocument',
    loadContext: 'loadTableDocument',
    rerender: () => tableEditorRef.value?.reRender?.()
  });

  return {
    isTableDocument: isCurrentType,
    flushTableSave: flushSave,
    rerenderTableEditor: rerenderEditor
  };
}
