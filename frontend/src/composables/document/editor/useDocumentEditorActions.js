import { ref } from 'vue';
import { ElMessage, ElMessageBox } from 'element-plus';
import {
  createDocument as createDocumentApi,
  createTemplate as createTemplateApi,
  deleteDocument as deleteDocumentApi,
  updateDocumentMeta
} from '@/services/api';
import { useApiAction } from '@/composables/actions/useApiAction';
import { DEFAULT_DOCUMENT_TYPE, normalizeDocumentType } from '@/utils/documentType';

export function useDocumentEditorActions({
  t,
  router,
  kbId,
  docId,
  currentDocType,
  currentDocTitle,
  markdownContent,
  tableContent,
  slideContent,
  richTextContent,
  docMachine,
  directoryTree,
  updateTreeNodeTitle,
  fetchDocuments
}) {
  const { runWithLoading } = useApiAction({ t });

  const resolveActiveDocType = () => {
    if (docMachine?.targetDocType.value) {
      return normalizeDocumentType(docMachine.targetDocType.value, '1');
    }
    return normalizeDocumentType(currentDocType?.value, '1');
  };

  const resolveActiveContent = () => {
    const type = resolveActiveDocType();
    if (type === '2') return tableContent?.value || '';
    if (type === '3') return slideContent?.value || '';
    if (type === '4') return richTextContent?.value || '';
    return markdownContent?.value || '';
  };

  const isEditingTitle = ref(false);
  const pendingTitle = ref('');
  const savingTitle = ref(false);
  const renamingNode = ref(false);

  const showCreateDialog = ref(false);
  const creatingLoading = ref(false);
  const pendingParentId = ref(null);
  const selectedDocumentType = ref(DEFAULT_DOCUMENT_TYPE);

  const savingTemplate = ref(false);
  const showTemplateDialog = ref(false);
  const templateDialogInit = ref({
    name: '',
    description: '',
    scope: 'own',
    tags: '',
    content: ''
  });

  const openCreateDocumentDialog = (parentId = null, documentType = DEFAULT_DOCUMENT_TYPE) => {
    pendingParentId.value = parentId;
    selectedDocumentType.value = normalizeDocumentType(documentType);
    showCreateDialog.value = true;
  };

  const cancelCreateDocument = () => {
    showCreateDialog.value = false;
  };

  const navigateToDocument = (externalId) => {
    if (kbId.value === 'personal') {
      router.push(`/mydocument/${externalId}`);
    } else {
      router.push(`/knowledge-base/${kbId.value}/document/${externalId}`);
    }
  };

  const navigateToList = () => {
    if (kbId.value === 'personal') {
      router.push('/mydocuments');
    } else {
      router.push(`/knowledge-base/${kbId.value}`);
    }
  };

  const createDocument = async (payload) => {
    const title = payload?.title?.trim() || t('document.untitledDefaultTitle');
    await runWithLoading(
      creatingLoading,
      async () => {
        const isPersonal = kbId.value === 'personal';
        const requestBody = {
          title,
          type: normalizeDocumentType(payload?.type || selectedDocumentType.value),
          container_type: isPersonal ? 'own' : 'knowledge_base',
          is_public: false
        };
        if (!isPersonal) {
          requestBody.knowledge_base_external_id = kbId.value;
        }
        if (payload?.parent_external_id) {
          requestBody.parent_external_id = payload.parent_external_id;
        } else if (pendingParentId.value) {
          requestBody.parent_external_id = pendingParentId.value;
        }
        if (payload?.template) {
          requestBody.template_id = payload.template;
        }
        return createDocumentApi(requestBody);
      },
      {
        context: 'createDocument',
        errorMessage: t('knowledgeBase.createDocumentError'),
        onSuccess: async (response) => {
          showCreateDialog.value = false;
          pendingParentId.value = null;
          await fetchDocuments();
          if (response?.external_id) {
            navigateToDocument(response.external_id);
          }
        }
      }
    );
  };

  const handleCreateFromTree = (payload) => {
    if (!payload) {
      openCreateDocumentDialog(null, DEFAULT_DOCUMENT_TYPE);
      return;
    }
    if (typeof payload === 'object' && ('target' in payload || 'type' in payload)) {
      openCreateDocumentDialog(
        payload?.target?.id || null,
        payload?.type || DEFAULT_DOCUMENT_TYPE
      );
      return;
    }
    openCreateDocumentDialog(payload?.id || null, DEFAULT_DOCUMENT_TYPE);
  };

  const deletingDocument = ref(false);

  const handleDeleteDocument = async (target) => {
    const targetId = target?.id || docId.value;
    if (!targetId) {
      return;
    }
    try {
      await ElMessageBox.confirm(t('document.deleteDocumentConfirm'), t('document.deleteDocument'), {
        confirmButtonText: t('button.confirm'),
        cancelButtonText: t('button.cancel'),
        type: 'warning'
      });
    } catch {
      return;
    }

    await runWithLoading(
      deletingDocument,
      () => deleteDocumentApi(targetId),
      {
        context: 'deleteDocument',
        successMessage: t('document.deleteSuccess'),
        errorMessage: t('common.requestFailed'),
        onSuccess: async () => {
          await fetchDocuments();
          if (String(targetId) === String(docId.value)) {
            navigateToList();
          }
        }
      }
    );
  };

  const handleRenameFromTree = async (payload) => {
    const targetNode = payload?.data;
    const nextTitle = payload?.title?.trim?.() || '';
    if (!targetNode?.id || !nextTitle) {
      ElMessage.error(t('knowledgeBase.titleRequired'));
      return;
    }
    const targetId = String(targetNode.id);
    const currentTitle = String(targetNode.label || '').trim();
    if (nextTitle === currentTitle) {
      return;
    }

    await runWithLoading(
      renamingNode,
      () => updateDocumentMeta(targetId, { title: nextTitle }),
      {
        context: 'renameDocumentFromTree',
        errorMessage: t('common.requestFailed'),
        onSuccess: () => {
          updateTreeNodeTitle(directoryTree.value, targetId, nextTitle);
          if (String(docId.value || '') === targetId) {
            currentDocTitle.value = nextTitle;
          }
        }
      }
    );
  };

  const startEditTitle = async () => {
    if (!docId.value) {
      return;
    }
    isEditingTitle.value = true;
    pendingTitle.value = currentDocTitle.value || '';
  };

  const cancelEditTitle = () => {
    isEditingTitle.value = false;
    pendingTitle.value = '';
  };

  const commitTitle = async () => {
    if (!isEditingTitle.value) {
      return;
    }
    const nextTitle = pendingTitle.value.trim();
    if (!nextTitle) {
      ElMessage.error(t('knowledgeBase.titleRequired'));
      return;
    }
    if (nextTitle === currentDocTitle.value) {
      cancelEditTitle();
      return;
    }
    if (!docId.value) {
      cancelEditTitle();
      return;
    }
    await runWithLoading(
      savingTitle,
      () => updateDocumentMeta(docId.value, { title: nextTitle }),
      {
        context: 'updateDocumentMeta',
        errorMessage: t('common.requestFailed'),
        onSuccess: () => {
          currentDocTitle.value = nextTitle;
          updateTreeNodeTitle(directoryTree.value, docId.value, nextTitle);
          cancelEditTitle();
        }
      }
    );
  };

  const openCreateTemplateDialog = () => {
    const defaultScope = kbId.value && kbId.value !== 'personal' ? 'knowledge_base' : 'own';
    templateDialogInit.value = {
      name: currentDocTitle.value || t('templates.untitled'),
      description: '',
      scope: defaultScope,
      tags: '',
      content: resolveActiveContent()
    };
    showTemplateDialog.value = true;
  };

  const handleHeaderCommand = (command) => {
    if (command === 'create_template') {
      openCreateTemplateDialog();
      return true;
    }
    return false;
  };

  const submitCreateTemplate = async (payload) => {
    const activeContent = resolveActiveContent();
    if (!activeContent || !activeContent.trim()) {
      ElMessage.error(t('templates.emptyContent'));
      return;
    }

    await runWithLoading(
      savingTemplate,
      async () => {
        const requestBody = {
          target_container: payload.scope,
          type: normalizeDocumentType(resolveActiveDocType() || DEFAULT_DOCUMENT_TYPE),
          template_content: JSON.stringify({
            name: payload.name,
            description: payload.description,
            content: activeContent,
            tags: payload.tags || []
          })
        };
        if (payload.scope === 'knowledge_base' && kbId.value && kbId.value !== 'personal') {
          requestBody.knowledge_base_id = kbId.value;
        }
        await createTemplateApi(requestBody);
      },
      {
        context: 'createTemplate',
        successMessage: t('templates.saveSuccess'),
        errorMessage: t('templates.saveFailed'),
        onSuccess: () => {
          showTemplateDialog.value = false;
        }
      }
    );
  };

  return {
    isEditingTitle,
    pendingTitle,
    showCreateDialog,
    creatingLoading,
    pendingParentId,
    selectedDocumentType,
    savingTemplate,
    showTemplateDialog,
    templateDialogInit,
    openCreateDocumentDialog,
    cancelCreateDocument,
    createDocument,
    handleCreateFromTree,
    handleDeleteDocument,
    handleRenameFromTree,
    startEditTitle,
    cancelEditTitle,
    commitTitle,
    handleHeaderCommand,
    submitCreateTemplate
  };
}
