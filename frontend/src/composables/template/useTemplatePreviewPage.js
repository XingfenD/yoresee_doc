import { computed, nextTick, ref, watch } from 'vue';
import { ElMessage } from 'element-plus';
import Vditor from 'vditor';
import { createDocument as createDocumentApi, getTemplate } from '@/services/api';
import { useApiAction } from '@/composables/actions/useApiAction';
import { DEFAULT_DOCUMENT_TYPE, normalizeDocumentType } from '@/utils/documentType';

export function useTemplatePreviewPage({ t, route, router, isDarkMode }) {
  const { runWithLoading } = useApiAction({ t });

  const loading = ref(false);
  const template = ref(null);
  const previewRef = ref(null);
  const showCreateDialog = ref(false);
  const creatingLoading = ref(false);

  const templateId = computed(() => route.params.templateId || route.params.id);

  const formatDate = (dateString) => {
    if (!dateString) return t('common.unknown');
    const date = new Date(dateString);
    return date.toLocaleDateString();
  };

  const scopeLabel = computed(() => {
    const scope = template.value?.scope;
    if (scope === 'system') return t('templates.public');
    if (scope === 'knowledge_base') return t('templates.scopeKb');
    if (scope === 'private') return t('templates.private');
    return '';
  });

  const scopeTagType = computed(() => {
    const scope = template.value?.scope;
    if (scope === 'system') return 'success';
    if (scope === 'knowledge_base') return 'warning';
    return 'info';
  });

  const parseTemplateContent = (raw) => {
    if (!raw) return '';
    try {
      const parsed = JSON.parse(raw);
      if (parsed && typeof parsed.content === 'string') {
        return parsed.content;
      }
    } catch (error) {
      // not json
    }
    return raw;
  };

  const previewContent = computed(() => parseTemplateContent(template.value?.content || ''));

  const renderPreview = async () => {
    if (!previewRef.value) return;
    await nextTick();
    const content = previewContent.value || '';
    if (!content) {
      previewRef.value.innerHTML = '';
      return;
    }
    await Vditor.preview(previewRef.value, content, {
      mode: isDarkMode.value ? 'dark' : 'light',
      theme: {
        current: isDarkMode.value ? 'dark' : 'light'
      },
      hljs: { style: isDarkMode.value ? 'monokai' : 'github' }
    });
  };

  const fetchTemplate = async () => {
    if (!templateId.value) return;
    await runWithLoading(
      loading,
      () => getTemplate(templateId.value, { record_recent_log: true }),
      {
        context: 'fetchTemplate',
        errorMessage: t('common.requestFailed'),
        onSuccess: (data) => {
          template.value = data.template;
        }
      }
    );
  };

  const goBack = () => {
    router.push('/templates');
  };

  const openCreateDocumentDialog = () => {
    showCreateDialog.value = true;
  };

  const cancelCreateDocument = () => {
    showCreateDialog.value = false;
  };

  const createDocument = async (payload) => {
    const title = payload?.title?.trim() || t('document.untitledDefaultTitle');
    const containerType = payload?.container_type === 'knowledge_base' ? 'knowledge_base' : 'own';
    const selectedKBExternalID =
      payload?.knowledge_base_external_id || template.value?.knowledge_base_external_id || '';
    if (containerType === 'knowledge_base' && !selectedKBExternalID) {
      ElMessage.error(t('knowledgeBase.selectKnowledgeBase'));
      return;
    }

    await runWithLoading(
      creatingLoading,
      async () => {
        const requestBody = {
          title,
          type: normalizeDocumentType(payload?.type || template.value?.type || DEFAULT_DOCUMENT_TYPE),
          container_type: containerType,
          is_public: typeof payload?.is_public === 'boolean' ? payload.is_public : false
        };
        if (containerType === 'knowledge_base') {
          requestBody.knowledge_base_external_id = selectedKBExternalID;
        }
        if (payload?.parent_external_id) {
          requestBody.parent_external_id = payload.parent_external_id;
        }
        if (payload?.template) {
          requestBody.template_id = payload.template;
        }
        const response = await createDocumentApi(requestBody);
        return { response, containerType, selectedKBExternalID };
      },
      {
        context: 'createDocumentFromTemplate',
        errorMessage: t('knowledgeBase.createDocumentError'),
        onSuccess: ({ response, containerType, selectedKBExternalID }) => {
          showCreateDialog.value = false;
          if (response?.external_id) {
            if (containerType === 'knowledge_base') {
              router.push(`/knowledge-base/${selectedKBExternalID}/document/${response.external_id}`);
            } else {
              router.push(`/mydocument/${response.external_id}`);
            }
          }
        }
      }
    );
  };

  watch(previewContent, () => {
    renderPreview();
  });

  watch(isDarkMode, () => {
    renderPreview();
  });

  const init = async () => {
    await fetchTemplate();
    await renderPreview();
  };

  return {
    loading,
    template,
    previewRef,
    showCreateDialog,
    creatingLoading,
    previewContent,
    scopeLabel,
    scopeTagType,
    formatDate,
    goBack,
    openCreateDocumentDialog,
    cancelCreateDocument,
    createDocument,
    init
  };
}
