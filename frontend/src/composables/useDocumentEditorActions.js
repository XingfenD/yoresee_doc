import { ref } from 'vue';
import { ElMessage, ElMessageBox } from 'element-plus';
import {
  createDocument as createDocumentApi,
  createTemplate as createTemplateApi,
  updateDocumentMeta
} from '@/services/api';
import { useApiAction } from '@/composables/useApiAction';

export function useDocumentEditorActions({
  t,
  router,
  kbId,
  docId,
  currentDocTitle,
  editorContent,
  directoryTree,
  updateTreeNodeTitle,
  fetchDocuments
}) {
  const { runWithLoading } = useApiAction({ t });

  const isEditingTitle = ref(false);
  const pendingTitle = ref('');
  const savingTitle = ref(false);

  const showCreateDialog = ref(false);
  const creatingLoading = ref(false);
  const pendingParentId = ref(null);

  const savingTemplate = ref(false);
  const showTemplateDialog = ref(false);
  const templateDialogInit = ref({
    name: '',
    description: '',
    scope: 'own',
    tags: '',
    content: ''
  });

  const openCreateDocumentDialog = (parentId = null) => {
    pendingParentId.value = parentId;
    showCreateDialog.value = true;
  };

  const cancelCreateDocument = () => {
    showCreateDialog.value = false;
  };

  const createDocument = async (payload) => {
    const title = payload?.title?.trim() || t('document.untitledDefaultTitle');
    await runWithLoading(
      creatingLoading,
      async () => {
        const isPersonal = kbId.value === 'personal';
        const requestBody = {
          title,
          type: payload.type || 'markdown',
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
          const isPersonal = kbId.value === 'personal';
          showCreateDialog.value = false;
          pendingParentId.value = null;
          await fetchDocuments();
          if (response?.external_id) {
            if (isPersonal) {
              router.push(`/mydocument/${response.external_id}`);
            } else {
              router.push(`/knowledge-base/${kbId.value}/document/${response.external_id}`);
            }
          }
        }
      }
    );
  };

  const handleCreateFromTree = (target) => {
    openCreateDocumentDialog(target?.id || null);
  };

  const handleDeleteDocument = async () => {
    if (!docId.value) {
      return;
    }
    try {
      await ElMessageBox.confirm(t('document.deleteDocumentConfirm'), t('document.deleteDocument'), {
        confirmButtonText: t('button.confirm'),
        cancelButtonText: t('button.cancel'),
        type: 'warning'
      });
      ElMessage.warning(t('document.deleteNotSupported'));
    } catch (error) {
      // cancel
    }
  };

  const handleRenameFromTree = () => {
    ElMessage.warning(t('document.renameNotSupported'));
  };

  const startEditTitle = async () => {
    if (!docId.value || docId.value === 'example') {
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
      content: editorContent.value || ''
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
    if (!editorContent.value || !editorContent.value.trim()) {
      ElMessage.error(t('templates.emptyContent'));
      return;
    }

    await runWithLoading(
      savingTemplate,
      async () => {
        const requestBody = {
          target_container: payload.scope,
          template_content: JSON.stringify({
            name: payload.name,
            description: payload.description,
            content: editorContent.value,
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
