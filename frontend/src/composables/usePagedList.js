import { computed, ref } from 'vue';

export function usePagedList(options = {}) {
  const {
    initialPage = 1,
    initialPageSize = 10,
    fetcher
  } = options;

  const items = ref([]);
  const page = ref(initialPage);
  const pageSize = ref(initialPageSize);
  const total = ref(0);
  const loading = ref(false);
  const loaded = ref(false);
  const hasMore = computed(() => items.value.length < total.value);

  const fetchPage = async (targetPage = page.value, targetPageSize = pageSize.value) => {
    if (loading.value || typeof fetcher !== 'function') return null;
    loading.value = true;
    try {
      const result = await fetcher({ page: targetPage, pageSize: targetPageSize });
      const nextItems = Array.isArray(result?.items) ? result.items : [];
      if (targetPage === 1) {
        items.value = nextItems;
      } else {
        items.value.push(...nextItems);
      }
      total.value = Number(result?.total) || items.value.length;
      page.value = targetPage;
      pageSize.value = targetPageSize;
      loaded.value = true;
      return result;
    } finally {
      loading.value = false;
    }
  };

  const loadMore = async () => {
    if (!hasMore.value || loading.value) return;
    await fetchPage(page.value + 1, pageSize.value);
  };

  const reset = () => {
    page.value = initialPage;
    pageSize.value = initialPageSize;
    total.value = 0;
    items.value = [];
    loaded.value = false;
  };

  return {
    items,
    page,
    pageSize,
    total,
    loading,
    loaded,
    hasMore,
    fetchPage,
    loadMore,
    reset
  };
}
