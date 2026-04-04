import { computed } from 'vue';
import { normalizeDocumentType } from '@/utils/documentType';

export function useDocumentEditorPolicy(options = {}) {
  const {
    kbId,
    docId,
    currentDocType
  } = options;

  const isPersonalDocument = computed(() => kbId.value === 'personal');
  const isExampleDocument = computed(() => docId.value === 'example');
  const hasDocument = computed(() => Boolean(docId.value));
  const isMarkdownDocument = computed(() => normalizeDocumentType(currentDocType.value, '1') === '1');

  const canManageAttachments = computed(() => hasDocument.value && !isExampleDocument.value);
  const canManageSettings = computed(() => hasDocument.value && !isExampleDocument.value);
  const collabEnabled = computed(() => hasDocument.value && !isExampleDocument.value && isMarkdownDocument.value);
  const inlineCommentEnabled = computed(() => collabEnabled.value);
  const createDialogKnowledgeBaseId = computed(() => (isPersonalDocument.value ? '' : kbId.value));
  const showTemplateDialogKbScope = computed(() => !isPersonalDocument.value);

  return {
    isPersonalDocument,
    isExampleDocument,
    hasDocument,
    isMarkdownDocument,
    canManageAttachments,
    canManageSettings,
    collabEnabled,
    inlineCommentEnabled,
    createDialogKnowledgeBaseId,
    showTemplateDialogKbScope
  };
}
