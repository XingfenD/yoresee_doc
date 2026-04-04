import { ElMessageBox } from 'element-plus';
import * as api from '@/services/api';
import { usePagedList } from '@/composables/list/usePagedList';
import { useLazyTabLoader } from '@/composables/list/useLazyTabLoader';
import { isActionCancelled, useApiAction } from '@/composables/actions/useApiAction';

export function useKnowledgeBaseListPage({ t, router }) {
  const { runApi } = useApiAction({ t });
  const runFetch = (context, action) =>
    runApi(action, {
      context,
      errorMessage: t('knowledgeBase.fetchError')
    });

  const myList = usePagedList({
    initialPage: 1,
    initialPageSize: 10,
    fetcher: async ({ page, pageSize }) => {
      const data = await api.getKnowledgeBases({ page, page_size: pageSize, only_mine: true });
      return { items: data.knowledge_bases || [], total: data.total || 0 };
    }
  });

  const recentList = usePagedList({
    initialPage: 1,
    initialPageSize: 10,
    fetcher: async ({ page, pageSize }) => {
      const data = await api.getRecentKnowledgeBases({ page, page_size: pageSize });
      return { items: data.knowledge_bases || [], total: data.total || 0 };
    }
  });

  const publicList = usePagedList({
    initialPage: 1,
    initialPageSize: 10,
    fetcher: async ({ page, pageSize }) => {
      const data = await api.getKnowledgeBases({ page, page_size: pageSize, is_public: true });
      return { items: data.knowledge_bases || [], total: data.total || 0 };
    }
  });

  const formatDate = (dateString) => {
    if (!dateString) return t('common.unknown');
    const date = new Date(dateString);
    return date.toLocaleDateString();
  };

  const myTagMapper = (kb) => (kb.is_public ? null : { type: 'info', label: t('knowledgeBase.private') });
  const recentTagMapper = (kb) => (kb.is_public ? { type: 'success', label: t('knowledgeBase.public') } : null);

  const myMetaMapper = (kb) => [
    { label: t('knowledgeBase.documentsCount'), value: kb.documents_count || 0 },
    { label: t('knowledgeBase.updatedAt'), value: formatDate(kb.updated_at) }
  ];

  const publicMetaMapper = (kb) => [
    { label: t('knowledgeBase.documentsCount'), value: kb.documents_count || 0 },
    { label: t('knowledgeBase.owner'), value: kb.creator_name || t('common.unknown') }
  ];

  const loadMoreMyKnowledgeBases = async () => {
    await runFetch('loadMoreMyKnowledgeBases', () => myList.loadMore());
  };

  const loadMorePublicKnowledgeBases = async () => {
    await runFetch('loadMorePublicKnowledgeBases', () => publicList.loadMore());
  };

  const viewKnowledgeBase = (kb) => {
    if (!kb?.external_id) return;
    router.push(`/knowledge-base/${kb.external_id}`);
  };

  const accessKnowledgeBase = (kb) => {
    if (!kb?.external_id) return;
    router.push(`/knowledge-base/${kb.external_id}`);
  };

  const createKnowledgeBase = async () => {
    await runApi(
      async () => {
        const { value } = await ElMessageBox.prompt(
          t('knowledgeBase.createPrompt'),
          t('knowledgeBase.createNew'),
          {
            confirmButtonText: t('button.confirm'),
            cancelButtonText: t('button.cancel'),
            inputPlaceholder: t('knowledgeBase.createPlaceholder'),
            inputValidator: (val) => {
              if (!val || !val.trim()) {
                return t('knowledgeBase.titleRequired');
              }
              return true;
            }
          }
        );

        const name = value.trim();
        const resp = await api.createKnowledgeBase({ name, is_public: false });

        activeTab.value = 'my';
        await myList.fetchPage(1, myList.pageSize.value);

        if (resp?.external_id) {
          router.push(`/knowledge-base/${resp.external_id}`);
        }
      },
      {
        context: 'createKnowledgeBase',
        successMessage: t('knowledgeBase.createSuccess'),
        errorMessage: t('common.requestFailed'),
        ignoreError: isActionCancelled
      }
    );
  };

  const { activeTab, ensureTabLoaded } = useLazyTabLoader({
    initialTab: 'my',
    tabs: {
      my: {
        loaded: () => myList.loaded.value,
        load: () => runFetch('fetchMyKnowledgeBases', () => myList.fetchPage(myList.page.value, myList.pageSize.value))
      },
      recent: {
        loaded: () => recentList.loaded.value,
        load: () => runFetch('fetchRecentKnowledgeBases', () => recentList.fetchPage(recentList.page.value, recentList.pageSize.value))
      },
      public: {
        loaded: () => publicList.loaded.value,
        load: () => runFetch('fetchPublicKnowledgeBases', () => publicList.fetchPage(publicList.page.value, publicList.pageSize.value))
      }
    }
  });

  const init = async () => {
    await ensureTabLoaded('my');
  };

  return {
    activeTab,
    recentKnowledgeBases: recentList.items,
    myKnowledgeBases: myList.items,
    publicKnowledgeBases: publicList.items,
    myLoading: myList.loading,
    publicLoading: publicList.loading,
    myHasMore: myList.hasMore,
    publicHasMore: publicList.hasMore,
    myTagMapper,
    recentTagMapper,
    myMetaMapper,
    publicMetaMapper,
    createKnowledgeBase,
    loadMoreMyKnowledgeBases,
    loadMorePublicKnowledgeBases,
    viewKnowledgeBase,
    accessKnowledgeBase,
    init
  };
}
