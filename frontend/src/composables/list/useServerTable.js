import { onScopeDispose, ref } from 'vue';

const normalizeKeyword = (value) => {
  const keyword = `${value ?? ''}`.trim();
  return keyword || undefined;
};

const toSafeTotal = (value, fallback = 0) => {
  const total = Number(value);
  return Number.isFinite(total) ? total : fallback;
};

export function useServerTable(options = {}) {
  const {
    fetcher,
    mapRows = (response) => response?.rows || [],
    mapTotal = (response, rows) => response?.total ?? rows.length,
    initialPage = 1,
    initialPageSize = 10,
    initialKeyword = '',
    searchDelay = 300,
    onError
  } = options;

  if (typeof fetcher !== 'function') {
    throw new Error('[useServerTable] fetcher is required');
  }

  const rows = ref([]);
  const total = ref(0);
  const loading = ref(false);
  const page = ref(initialPage);
  const pageSize = ref(initialPageSize);
  const keyword = ref(initialKeyword);

  let searchTimer = null;
  let latestRequestId = 0;

  const clearSearchTimer = () => {
    if (!searchTimer) {
      return;
    }
    clearTimeout(searchTimer);
    searchTimer = null;
  };

  const buildQuery = (overrides = {}) => ({
    page: overrides.page ?? page.value,
    page_size: overrides.page_size ?? overrides.pageSize ?? pageSize.value,
    keyword:
      overrides.keyword !== undefined
        ? normalizeKeyword(overrides.keyword)
        : normalizeKeyword(keyword.value)
  });

  const load = async (overrides = {}) => {
    const requestId = ++latestRequestId;
    loading.value = true;
    try {
      const response = await fetcher(buildQuery(overrides));
      const nextRows = mapRows(response) || [];
      if (requestId === latestRequestId) {
        rows.value = nextRows;
        total.value = toSafeTotal(mapTotal(response, nextRows), nextRows.length);
      }
      return response;
    } catch (error) {
      if (requestId === latestRequestId) {
        rows.value = [];
        total.value = 0;
        if (typeof onError === 'function') {
          onError(error);
        } else {
          console.error('[useServerTable] load failed', error);
        }
      }
      return null;
    } finally {
      if (requestId === latestRequestId) {
        loading.value = false;
      }
    }
  };

  const reload = async () => load();

  const handlePageChange = async (nextPage) => {
    clearSearchTimer();
    page.value = nextPage;
    await load();
  };

  const handleSizeChange = async (nextPageSize) => {
    clearSearchTimer();
    pageSize.value = nextPageSize;
    page.value = 1;
    await load();
  };

  const handleSearch = () => {
    clearSearchTimer();
    searchTimer = setTimeout(async () => {
      page.value = 1;
      await load();
    }, searchDelay);
  };

  const reset = ({ keepKeyword = false } = {}) => {
    page.value = initialPage;
    pageSize.value = initialPageSize;
    if (!keepKeyword) {
      keyword.value = '';
    }
    rows.value = [];
    total.value = 0;
  };

  onScopeDispose(() => {
    clearSearchTimer();
  });

  return {
    rows,
    total,
    loading,
    page,
    pageSize,
    keyword,
    load,
    reload,
    reset,
    handlePageChange,
    handleSizeChange,
    handleSearch,
    clearSearchTimer
  };
}
