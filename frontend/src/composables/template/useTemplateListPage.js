import { computed, reactive, ref, watch } from 'vue';
import { ElMessage } from 'element-plus';
import {
  createTemplate as createTemplateApi,
  updateTemplateSettings as updateTemplateSettingsApi
} from '@/services/api';
import { useLazyTabLoader } from '@/composables/list/useLazyTabLoader';
import { useTemplateCatalog } from '@/composables/template/useTemplateCatalog';
import { useApiAction } from '@/composables/actions/useApiAction';

export function useTemplateListPage({ t, router }) {
  const { runWithLoading, createApiErrorHandler } = useApiAction({ t });
  const handleTemplateLoadError = createApiErrorHandler({
    context: 'loadTemplates',
    errorMessage: t('common.requestFailed')
  });

  const showCreateDialog = ref(false);
  const creatingTemplate = ref(false);
  const keyword = ref('');
  const currentPage = ref(1);
  const pageSize = ref(9);
  const previewTemplate = ref(null);
  const showPreviewDialog = ref(false);
  const showSettingsDialog = ref(false);
  const savingTemplateSettings = ref(false);
  const templateSettingsForm = reactive({
    id: '',
    name: '',
    description: '',
    is_public: false,
    scope: 'private'
  });
  const templateDialogInit = ref({
    name: '',
    description: '',
    scope: 'own',
    tags: '',
    content: ''
  });

  const {
    myTemplates,
    recentTemplates,
    publicTemplates,
    loadingMy,
    loadingRecent,
    loadingPublic,
    ensureLoaded,
    invalidateScope,
    isLoaded
  } = useTemplateCatalog({
    includeKnowledgeBase: false,
    onError: handleTemplateLoadError
  });

  const formatDate = (value) => {
    if (!value) return t('common.unknown');
    const date = new Date(value);
    if (Number.isNaN(date.getTime())) return value;
    return date.toLocaleDateString();
  };

  const openTemplate = (tpl) => {
    if (!tpl?.id) return;
    router.push(`/template/${tpl.id}`);
  };

  const openCreateTemplateDialog = () => {
    templateDialogInit.value = {
      name: '',
      description: '',
      scope: 'own',
      tags: '',
      content: ''
    };
    showCreateDialog.value = true;
  };

  const { activeTab, ensureTabLoaded } = useLazyTabLoader({
    initialTab: 'my',
    tabs: {
      my: {
        loaded: () => isLoaded('my'),
        load: () => ensureLoaded('my')
      },
      recent: {
        loaded: () => isLoaded('recent'),
        load: () => ensureLoaded('recent')
      },
      public: {
        loaded: () => isLoaded('public'),
        load: () => ensureLoaded('public')
      }
    }
  });

  const refreshCurrentTab = async () => {
    invalidateScope(activeTab.value);
    await ensureTabLoaded(activeTab.value);
  };

  watch(activeTab, () => {
    keyword.value = '';
    currentPage.value = 1;
  });

  const currentTemplates = computed(() => {
    if (activeTab.value === 'recent') {
      return recentTemplates.value;
    }
    if (activeTab.value === 'public') {
      return publicTemplates.value;
    }
    return myTemplates.value;
  });

  const currentLoading = computed(() => {
    if (activeTab.value === 'recent') {
      return loadingRecent.value;
    }
    if (activeTab.value === 'public') {
      return loadingPublic.value;
    }
    return loadingMy.value;
  });

  const filteredTemplates = computed(() => {
    const query = keyword.value.trim().toLowerCase();
    if (!query) {
      return currentTemplates.value;
    }
    return currentTemplates.value.filter((tpl) => {
      const name = String(tpl?.name || '').toLowerCase();
      const description = String(tpl?.description || '').toLowerCase();
      return name.includes(query) || description.includes(query);
    });
  });

  const paginationTotal = computed(() => filteredTemplates.value.length);

  const pagedTemplates = computed(() => {
    const start = (currentPage.value - 1) * pageSize.value;
    const end = start + pageSize.value;
    return filteredTemplates.value.slice(start, end);
  });

  watch(filteredTemplates, (list) => {
    const maxPage = Math.max(1, Math.ceil(list.length / pageSize.value));
    if (currentPage.value > maxPage) {
      currentPage.value = maxPage;
    }
  });

  const handlePageChange = (nextPage) => {
    currentPage.value = nextPage;
  };

  const handleSearch = () => {
    currentPage.value = 1;
  };

  const parseTemplateContent = (raw) => {
    if (!raw) return '';
    try {
      const parsed = JSON.parse(raw);
      if (typeof parsed?.content === 'string') {
        return parsed.content;
      }
    } catch (error) {
      // not json payload
    }
    return raw;
  };

  const previewTitle = computed(
    () => previewTemplate.value?.name || t('templates.untitled')
  );
  const previewContent = computed(() => parseTemplateContent(previewTemplate.value?.content || ''));

  const openPreviewDialog = (tpl) => {
    if (!tpl?.id) {
      return;
    }
    previewTemplate.value = tpl;
    showPreviewDialog.value = true;
  };

  const closePreviewDialog = () => {
    showPreviewDialog.value = false;
    previewTemplate.value = null;
  };

  const openTemplateSettings = (tpl) => {
    if (!tpl?.id) {
      return;
    }
    templateSettingsForm.id = tpl.id;
    templateSettingsForm.name = tpl.name || '';
    templateSettingsForm.description = tpl.description || '';
    templateSettingsForm.scope = tpl.scope || 'private';
    templateSettingsForm.is_public = tpl.scope === 'system';
    showSettingsDialog.value = true;
  };

  const submitTemplateSettings = async () => {
    const nextName = templateSettingsForm.name.trim();
    if (!nextName) {
      ElMessage.warning(t('templates.nameRequired'));
      return;
    }
    await runWithLoading(
      savingTemplateSettings,
      async () => {
        await updateTemplateSettingsApi({
          template_id: templateSettingsForm.id,
          name: nextName,
          description: templateSettingsForm.description.trim(),
          is_public: templateSettingsForm.is_public
        });
        invalidateScope('my');
        invalidateScope('recent');
        invalidateScope('public');
        await ensureLoaded(activeTab.value, { force: true });
      },
      {
        context: 'updateTemplateSettings',
        successMessage: t('message.saveSuccessGeneric'),
        errorMessage: t('message.saveFailedGeneric'),
        onSuccess: () => {
          showSettingsDialog.value = false;
        }
      }
    );
  };

  const submitCreateTemplate = async (payload) => {
    await runWithLoading(
      creatingTemplate,
      async () => {
        const requestBody = {
          target_container: payload.scope,
          template_content: JSON.stringify({
            name: payload.name,
            description: payload.description,
            content: payload.content,
            tags: payload.tags || []
          })
        };
        await createTemplateApi(requestBody);
      },
      {
        context: 'createTemplate',
        successMessage: t('templates.saveSuccess'),
        errorMessage: t('templates.saveFailed'),
        onSuccess: async () => {
          showCreateDialog.value = false;
          await refreshCurrentTab();
        }
      }
    );
  };

  const init = async () => {
    await ensureTabLoaded('my');
  };

  return {
    activeTab,
    showCreateDialog,
    creatingTemplate,
    keyword,
    currentPage,
    pageSize,
    paginationTotal,
    currentLoading,
    filteredTemplates,
    pagedTemplates,
    showPreviewDialog,
    previewTitle,
    previewContent,
    showSettingsDialog,
    savingTemplateSettings,
    templateSettingsForm,
    templateDialogInit,
    myTemplates,
    recentTemplates,
    publicTemplates,
    loadingMy,
    loadingRecent,
    loadingPublic,
    formatDate,
    openTemplate,
    openPreviewDialog,
    closePreviewDialog,
    openTemplateSettings,
    submitTemplateSettings,
    handlePageChange,
    handleSearch,
    openCreateTemplateDialog,
    submitCreateTemplate,
    init
  };
}
