import { computed, ref } from 'vue';
import { ElMessage, ElMessageBox } from 'element-plus';
import { Document, Clock, User } from '@element-plus/icons-vue';
import {
  getKnowledgeBaseDetail,
  createDocument as createDocumentApi
} from '@/services/api.js';
import { useWorkspaceShell } from '@/composables/useWorkspaceShell';
import { useTemplateCatalog } from '@/composables/useTemplateCatalog';
import { usePageBoot } from '@/composables/usePageBoot';

export function useKnowledgeBaseDetailPage({ t, router, route, userStore, locale }) {
  const {
    systemName,
    activeMenu,
    isDarkMode,
    currentLanguage,
    userInfo,
    userAvatar,
    initLanguage,
    handleLanguageChange,
    toggleTheme,
    handleLogout,
    handleMenuSelect,
    fetchSystemInfo
  } = useWorkspaceShell({
    locale,
    router,
    userStore,
    defaultActiveMenu: 'knowledge-base'
  });
  const { boot } = usePageBoot({ initLanguage, fetchSystemInfo });

  const knowledgeBaseName = ref('');
  const knowledgeBaseDescription = ref('');
  const totalDocuments = ref(0);
  const lastUpdated = ref('');
  const ownerName = ref('');
  const knowledgeBaseData = ref(null);
  const loading = ref(false);
  const {
    kbTemplates,
    loadingKb: kbTemplatesLoading,
    ensureLoaded: ensureTemplateCatalogLoaded
  } = useTemplateCatalog({
    includeKnowledgeBase: true,
    knowledgeBaseId: computed(() => route.params.id || ''),
    onError: (error) => {
      console.error('获取知识库模板失败:', error);
    }
  });

  const searchKeyword = ref('');
  const sortBy = ref('name');
  const currentPage = ref(1);
  const pageSize = ref(50);
  const totalDocumentsCount = ref(0);

  const showCreateDialog = ref(false);
  const creatingLoading = ref(false);

  const sortOptions = computed(() => [
    { value: 'name', label: t('knowledgeBase.sortByName') },
    { value: 'date', label: t('knowledgeBase.sortByDate') },
    { value: 'type', label: t('knowledgeBase.sortByType') }
  ]);

  const formatDate = (dateString) => {
    if (!dateString) return t('common.unknown');
    const date = new Date(dateString);
    return date.toLocaleDateString();
  };

  const knowledgeBaseStats = computed(() => [
    { key: 'documents', icon: Document, label: t('knowledgeBase.documentsCount'), value: totalDocuments.value },
    { key: 'updated', icon: Clock, label: t('knowledgeBase.lastUpdated'), value: formatDate(lastUpdated.value) },
    { key: 'owner', icon: User, label: t('knowledgeBase.owner'), value: ownerName.value }
  ]);

  const directoryTreeData = computed(() => {
    if (!knowledgeBaseData.value?.documents) return [];

    const mapDoc = (doc, parentId = null) => ({
      id: doc.external_id,
      label: doc.title,
      type: doc.type,
      isFolder: doc.hasChildren || (doc.children && doc.children.length > 0),
      isLeaf: !(doc.hasChildren || (doc.children && doc.children.length > 0)),
      tags: doc.tags || [],
      parentId,
      children: doc.children ? doc.children.map((child) => mapDoc(child, doc.external_id)) : []
    });

    return knowledgeBaseData.value.documents.map((doc) => mapDoc(doc));
  });

  const loadKnowledgeBaseDetail = async () => {
    const knowledgeBaseExternalID = route.params.id;
    if (!knowledgeBaseExternalID) {
      ElMessage.error(t('message.knowledgeBaseNotFound'));
      return;
    }

    loading.value = true;
    try {
      const data = await getKnowledgeBaseDetail(knowledgeBaseExternalID, {
        record_recent_log: true,
        page: currentPage.value,
        page_size: pageSize.value
      });

      knowledgeBaseData.value = data;

      if (data.knowledge_base) {
        knowledgeBaseName.value = data.knowledge_base.name;
        knowledgeBaseDescription.value = data.knowledge_base.description;
        lastUpdated.value = data.knowledge_base.updated_at;
        totalDocuments.value = data.knowledge_base.documents_count || 0;
        totalDocumentsCount.value = data.total_count || 0;
        ownerName.value = data.knowledge_base.creator_name || t('common.unknown');
      }
    } catch (error) {
      console.error('加载知识库详情失败:', error);
      ElMessage.error(t('message.loadKnowledgeBaseError'));
    } finally {
      loading.value = false;
    }
  };

  const fetchKnowledgeBaseTemplates = async () => {
    await ensureTemplateCatalogLoaded('knowledge_base');
  };

  const openCreateDocumentDialog = () => {
    showCreateDialog.value = true;
  };

  const createDocument = async (payload) => {
    if (!payload?.title?.trim()) {
      ElMessage.error(t('knowledgeBase.titleRequired'));
      return;
    }

    try {
      const knowledgeBaseExternalID = route.params.id;
      if (!knowledgeBaseExternalID) {
        ElMessage.error(t('message.knowledgeBaseNotFound'));
        return;
      }

      creatingLoading.value = true;

      await createDocumentApi({
        title: payload.title,
        type: payload.type || 'markdown',
        container_type: 'knowledge_base',
        knowledge_base_external_id: knowledgeBaseExternalID,
        parent_external_id: payload.parent_external_id || undefined,
        template_id: payload.template || undefined
      });

      showCreateDialog.value = false;
      await loadKnowledgeBaseDetail();
    } catch (error) {
      console.error('创建文档失败:', error);
      ElMessage.error(t('knowledgeBase.createDocumentError'));
    } finally {
      creatingLoading.value = false;
    }
  };

  const cancelCreateDocument = () => {
    showCreateDialog.value = false;
  };

  const openDocument = (data) => {
    router.push(`/knowledge-base/${route.params.id}/document/${data.id}`);
  };

  const handleTreeNodeClick = (data) => {
    openDocument(data);
  };

  const handleNodeAction = (command, data) => {
    switch (command) {
      case 'rename':
        ElMessage.info(`${t('message.renameDocument')}: ${data.label}`);
        break;
      case 'share':
        ElMessage.info(`${t('message.shareDocument')}: ${data.label}`);
        break;
      case 'delete':
        ElMessageBox.confirm(t('message.confirmDelete'), t('common.warning'), {
          confirmButtonText: t('button.confirm'),
          cancelButtonText: t('button.cancel'),
          type: 'warning'
        }).then(() => {
          ElMessage.success(t('message.deleteSuccess'));
        });
        break;
      default:
        break;
    }
  };

  const handleSizeChange = async (val) => {
    pageSize.value = val;
    currentPage.value = 1;
    await loadKnowledgeBaseDetail();
  };

  const handleCurrentChange = async (val) => {
    currentPage.value = val;
    await loadKnowledgeBaseDetail();
  };

  const goBackToKnowledgeBase = () => {
    router.push('/knowledge-base');
  };

  const templateTagMapper = () => ({ type: 'info', label: t('templates.private') });
  const templateMetaMapper = (tpl) => [
    { label: t('templates.updatedAt'), value: formatDate(tpl.updated_at || tpl.updatedAt) }
  ];

  const openTemplate = (tpl) => {
    if (!tpl?.id) return;
    router.push(`/template/${tpl.id}`);
  };

  const init = async () => {
    await boot(loadKnowledgeBaseDetail, fetchKnowledgeBaseTemplates);
  };

  return {
    systemName,
    activeMenu,
    isDarkMode,
    currentLanguage,
    userInfo,
    userAvatar,
    knowledgeBaseName,
    knowledgeBaseDescription,
    knowledgeBaseData,
    loading,
    kbTemplates,
    kbTemplatesLoading,
    knowledgeBaseStats,
    searchKeyword,
    sortBy,
    sortOptions,
    currentPage,
    pageSize,
    totalDocumentsCount,
    showCreateDialog,
    creatingLoading,
    createDocument,
    cancelCreateDocument,
    openCreateDocumentDialog,
    handleTreeNodeClick,
    openDocument,
    handleNodeAction,
    handleSizeChange,
    handleCurrentChange,
    handleMenuSelect,
    handleLanguageChange,
    toggleTheme,
    handleLogout,
    goBackToKnowledgeBase,
    templateTagMapper,
    templateMetaMapper,
    openTemplate,
    directoryTreeData,
    init
  };
}
