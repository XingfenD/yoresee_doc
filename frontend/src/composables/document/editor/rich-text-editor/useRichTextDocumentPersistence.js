import { useTypedDocumentPersistence } from '@/composables/document/editor/shared/useTypedDocumentPersistence';

export function useRichTextDocumentPersistence(options = {}) {
  const {
    docId,
    currentDocType,
    editorContent,
    richTextEditorRef,
    t,
    getDocumentContent,
    updateDocument
  } = options;

  const {
    isCurrentType,
    flushSave,
    rerenderEditor
  } = useTypedDocumentPersistence({
    type: '4',
    docId,
    currentDocType,
    editorContent,
    t,
    getDocumentContent,
    updateDocument,
    saveContext: 'saveRichTextDocument',
    loadContext: 'loadRichTextDocument',
    rerender: () => richTextEditorRef.value?.reRender?.()
  });

  return {
    isRichTextDocument: isCurrentType,
    flushRichTextSave: flushSave,
    rerenderRichTextEditor: rerenderEditor
  };
}
