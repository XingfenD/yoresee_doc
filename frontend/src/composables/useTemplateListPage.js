import { ref } from 'vue';
import { createTemplate as createTemplateApi } from '@/services/api';
import { useLazyTabLoader } from '@/composables/useLazyTabLoader';
import { useTemplateCatalog } from '@/composables/useTemplateCatalog';
import { useApiAction } from '@/composables/useApiAction';

export function useTemplateListPage({ t, router }) {
  const { runWithLoading, createApiErrorHandler } = useApiAction({ t });
  const handleTemplateLoadError = createApiErrorHandler({
    context: 'loadTemplates',
    errorMessage: t('common.requestFailed')
  });

  const showCreateDialog = ref(false);
  const creatingTemplate = ref(false);
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

  const myTagMapper = () => ({ type: 'info', label: t('templates.private') });
  const recentTagMapper = (tpl) =>
    tpl.scope === 'system'
      ? { type: 'success', label: t('templates.public') }
      : { type: 'info', label: t('templates.private') };
  const publicTagMapper = () => ({ type: 'success', label: t('templates.public') });

  const formatDate = (value) => {
    if (!value) return t('common.unknown');
    const date = new Date(value);
    if (Number.isNaN(date.getTime())) return value;
    return date.toLocaleDateString();
  };

  const templateMetaMapper = (tpl) => [
    { label: t('templates.updatedAt'), value: formatDate(tpl.updated_at || tpl.updatedAt) }
  ];

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
    templateDialogInit,
    myTemplates,
    recentTemplates,
    publicTemplates,
    loadingMy,
    loadingRecent,
    loadingPublic,
    myTagMapper,
    recentTagMapper,
    publicTagMapper,
    templateMetaMapper,
    openTemplate,
    openCreateTemplateDialog,
    submitCreateTemplate,
    init
  };
}
