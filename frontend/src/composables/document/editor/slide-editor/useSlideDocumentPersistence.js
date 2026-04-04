import { useTypedDocumentPersistence } from '@/composables/document/editor/shared/useTypedDocumentPersistence';

export function useSlideDocumentPersistence(options = {}) {
  const {
    docId,
    currentDocType,
    editorContent,
    slideEditorRef,
    t,
    getDocumentContent,
    updateDocument
  } = options;

  const {
    isCurrentType,
    flushSave,
    rerenderEditor
  } = useTypedDocumentPersistence({
    type: '3',
    docId,
    currentDocType,
    editorContent,
    t,
    getDocumentContent,
    updateDocument,
    saveContext: 'saveSlideDocument',
    loadContext: 'loadSlideDocument',
    rerender: () => slideEditorRef.value?.reRender?.()
  });

  return {
    isSlideDocument: isCurrentType,
    flushSlideSave: flushSave,
    rerenderSlideEditor: rerenderEditor
  };
}
