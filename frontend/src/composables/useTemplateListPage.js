import { computed, ref } from 'vue';
import { ElMessage } from 'element-plus';
import { createTemplate as createTemplateApi, listRecentTemplates, listTemplates } from '@/services/api';
import { useLazyTabLoader } from '@/composables/useLazyTabLoader';

export function useTemplateListPage({ t, router }) {
  const showCreateDialog = ref(false);
  const creatingTemplate = ref(false);
  const templateDialogInit = ref({
    name: '',
    description: '',
    scope: 'own',
    tags: '',
    content: ''
  });

  const myTemplates = ref([]);
  const recentTemplates = ref([]);
  const publicTemplates = ref([]);
  const loadingMy = ref(false);
  const loadingRecent = ref(false);
  const loadingPublic = ref(false);
  const myLoaded = ref(false);
  const recentLoaded = ref(false);
  const publicLoaded = ref(false);

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

  const fetchMyTemplates = async () => {
    if (loadingMy.value) return;
    loadingMy.value = true;
    try {
      const data = await listTemplates({
        only_mine: true,
        target_container: 'own',
        order_by: 'updated_at',
        order_desc: true,
        page: 1,
        page_size: 50
      });
      myTemplates.value = data.templates || [];
      myLoaded.value = true;
    } catch (err) {
      console.error('获取我的模板失败:', err);
    } finally {
      loadingMy.value = false;
    }
  };

  const fetchRecentTemplates = async () => {
    if (loadingRecent.value) return;
    loadingRecent.value = true;
    try {
      const data = await listRecentTemplates({
        page: 1,
        page_size: 50
      });
      recentTemplates.value = data.templates || [];
      recentLoaded.value = true;
    } catch (err) {
      console.error('获取最近模板失败:', err);
    } finally {
      loadingRecent.value = false;
    }
  };

  const fetchPublicTemplates = async () => {
    if (loadingPublic.value) return;
    loadingPublic.value = true;
    try {
      const data = await listTemplates({
        target_container: 'public',
        order_by: 'updated_at',
        order_desc: true,
        page: 1,
        page_size: 50
      });
      publicTemplates.value = data.templates || [];
      publicLoaded.value = true;
    } catch (err) {
      console.error('获取公开模板失败:', err);
    } finally {
      loadingPublic.value = false;
    }
  };

  const { activeTab, ensureTabLoaded } = useLazyTabLoader({
    initialTab: 'my',
    tabs: {
      my: {
        loaded: () => myLoaded.value,
        load: fetchMyTemplates
      },
      recent: {
        loaded: () => recentLoaded.value,
        load: fetchRecentTemplates
      },
      public: {
        loaded: () => publicLoaded.value,
        load: fetchPublicTemplates
      }
    }
  });

  const refreshCurrentTab = async () => {
    if (activeTab.value === 'my') {
      myLoaded.value = false;
    } else if (activeTab.value === 'recent') {
      recentLoaded.value = false;
    } else if (activeTab.value === 'public') {
      publicLoaded.value = false;
    }
    await ensureTabLoaded(activeTab.value);
  };

  const submitCreateTemplate = async (payload) => {
    if (creatingTemplate.value) return;
    try {
      creatingTemplate.value = true;
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
      showCreateDialog.value = false;
      ElMessage.success(t('templates.saveSuccess'));
      await refreshCurrentTab();
    } catch (err) {
      console.error('创建模板失败:', err);
      ElMessage.error(t('templates.saveFailed'));
    } finally {
      creatingTemplate.value = false;
    }
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
