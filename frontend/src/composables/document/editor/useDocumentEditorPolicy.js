import { computed } from 'vue';
import { normalizeDocumentType } from '@/utils/documentType';

export function useDocumentEditorPolicy(options = {}) {
  const {
    kbId,
    docId,
    currentDocType
  } = options;

  const isPersonalDocument = computed(() => kbId.value === 'personal');
  const hasDocument = computed(() => Boolean(docId.value));
  const normalizedDocType = computed(() => normalizeDocumentType(currentDocType.value, '1'));
  const isMarkdownDocument = computed(() => normalizedDocType.value === '1');
  const isTableDocument = computed(() => normalizedDocType.value === '2');
  const isSlideDocument = computed(() => normalizedDocType.value === '3');
  const isRichTextDocument = computed(() => normalizedDocType.value === '4');

  const canManageAttachments = computed(() => hasDocument.value);
  const canManageSettings = computed(() => hasDocument.value);
  const collabEnabled = computed(
    () => hasDocument.value && (isMarkdownDocument.value || isRichTextDocument.value)
  );
  const inlineCommentEnabled = computed(
    () => hasDocument.value && (isMarkdownDocument.value || isRichTextDocument.value)
  );
  const createDialogKnowledgeBaseId = computed(() => (isPersonalDocument.value ? '' : kbId.value));
  const showTemplateDialogKbScope = computed(() => !isPersonalDocument.value);

  return {
    isPersonalDocument,
    hasDocument,
    normalizedDocType,
    isMarkdownDocument,
    isTableDocument,
    isSlideDocument,
    isRichTextDocument,
    canManageAttachments,
    canManageSettings,
    collabEnabled,
    inlineCommentEnabled,
    createDialogKnowledgeBaseId,
    showTemplateDialogKbScope
  };
}
