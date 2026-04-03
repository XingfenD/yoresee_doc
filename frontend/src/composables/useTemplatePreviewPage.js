import { computed, nextTick, ref, watch } from 'vue';
import Vditor from 'vditor';
import { createDocument as createDocumentApi, getTemplate } from '@/services/api';
import { useApiAction } from '@/composables/useApiAction';

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
    await runWithLoading(
      creatingLoading,
      async () => {
        const kbExternalId = template.value?.knowledge_base_external_id || '';
        const isKnowledgeBase = template.value?.scope === 'knowledge_base' && !!kbExternalId;
        const requestBody = {
          title,
          type: payload.type || 'markdown',
          container_type: isKnowledgeBase ? 'knowledge_base' : 'own',
          is_public: false
        };
        if (isKnowledgeBase) {
          requestBody.knowledge_base_external_id = kbExternalId;
        }
        if (payload?.parent_external_id) {
          requestBody.parent_external_id = payload.parent_external_id;
        }
        if (payload?.template) {
          requestBody.template_id = payload.template;
        }
        const response = await createDocumentApi(requestBody);
        return { response, isKnowledgeBase, kbExternalId };
      },
      {
        context: 'createDocumentFromTemplate',
        errorMessage: t('knowledgeBase.createDocumentError'),
        onSuccess: ({ response, isKnowledgeBase, kbExternalId }) => {
          showCreateDialog.value = false;
          if (response?.external_id) {
            if (isKnowledgeBase) {
              router.push(`/knowledge-base/${kbExternalId}/document/${response.external_id}`);
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
