import { computed, nextTick, ref, watch } from 'vue';
import { ElMessage } from 'element-plus';
import Vditor from 'vditor';
import { createDocument as createDocumentApi, getTemplate } from '@/services/api';

export function useTemplatePreviewPage({ t, route, router, isDarkMode }) {
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
      theme: isDarkMode.value ? 'dark' : 'classic',
      hljs: { style: isDarkMode.value ? 'monokai' : 'github' }
    });
  };

  const fetchTemplate = async () => {
    if (!templateId.value) return;
    loading.value = true;
    try {
      const data = await getTemplate(templateId.value, { record_recent_log: true });
      template.value = data.template;
    } catch (error) {
      console.error('获取模板失败:', error);
      ElMessage.error(t('common.requestFailed'));
    } finally {
      loading.value = false;
    }
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
    if (!payload?.title?.trim()) {
      ElMessage.error(t('knowledgeBase.titleRequired'));
      return;
    }
    try {
      creatingLoading.value = true;
      const kbExternalId = template.value?.knowledge_base_external_id || '';
      const isKnowledgeBase = template.value?.scope === 'knowledge_base' && !!kbExternalId;
      const requestBody = {
        title: payload.title,
        type: payload.type || 'markdown',
        container_type: isKnowledgeBase ? 'knowledge_base' : 'own'
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
      showCreateDialog.value = false;
      if (response?.external_id) {
        if (isKnowledgeBase) {
          router.push(`/knowledge-base/${kbExternalId}/document/${response.external_id}`);
        } else {
          router.push(`/mydocument/${response.external_id}`);
        }
      }
    } catch (error) {
      console.error('创建文档失败:', error);
      ElMessage.error(t('knowledgeBase.createDocumentError'));
    } finally {
      creatingLoading.value = false;
    }
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
