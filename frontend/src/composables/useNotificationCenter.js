import { computed, ref, watch } from 'vue';
import { ElMessage } from 'element-plus';
import { listNotifications, markNotificationsRead, markAllNotificationsRead } from '@/services/api';

export function useNotificationCenter({ t }) {
  const activeTab = ref('all');
  const notifications = ref([]);
  const totalCount = ref(0);
  const currentPage = ref(1);
  const pageSize = ref(6);
  const loading = ref(false);
  const selectedIds = ref([]);

  const columns = computed(() => [
    { key: 'select', label: '', width: 36, align: 'center', className: 'select-column' },
    { key: 'title', label: t('user.notifications.columns.title'), minWidth: 320, flex: 1.6 },
    { key: 'type', label: t('user.notifications.columns.type'), minWidth: 120, align: 'center' },
    { key: 'created_at', label: t('user.notifications.columns.time'), minWidth: 180, align: 'center' },
    { key: 'actions', label: t('common.actions'), minWidth: 160, align: 'center' }
  ]);

  const displayNotifications = computed(() => notifications.value);

  const selectAll = computed(() => {
    if (!displayNotifications.value.length) return false;
    return displayNotifications.value.every((item) => selectedIds.value.includes(item.external_id));
  });

  const selectIndeterminate = computed(() => {
    const selectedCount = displayNotifications.value.filter((item) => selectedIds.value.includes(item.external_id)).length;
    return selectedCount > 0 && selectedCount < displayNotifications.value.length;
  });

  const toggleSelectAll = (value) => {
    if (value) {
      selectedIds.value = displayNotifications.value.map((item) => item.external_id);
      return;
    }
    selectedIds.value = selectedIds.value.filter(
      (id) => !displayNotifications.value.some((item) => item.external_id === id)
    );
  };

  const tagType = (value) => {
    if (value === 'mention') return 'warning';
    if (value === 'reply') return 'success';
    return 'info';
  };

  const tagLabel = (value) => {
    if (value === 'mention') return t('user.notifications.types.mention');
    if (value === 'reply') return t('user.notifications.types.reply');
    if (value === 'comment') return t('user.notifications.types.comment');
    if (value === 'system') return t('user.notifications.types.system');
    return value || t('user.notifications.types.system');
  };

  const buildTitle = (row) => row.title || row.content || '-';

  const formatDate = (value) => {
    if (!value) return '-';
    const date = new Date(value);
    if (Number.isNaN(date.getTime())) return value;
    return date.toLocaleString();
  };

  const refreshUnreadStatus = async () => {
    try {
      const resp = await listNotifications({ page: 1, page_size: 1, status: 'unread' });
      const hasUnread = Number(resp.total) > 0;
      window.dispatchEvent(new CustomEvent('notifications:unread', { detail: { hasUnread } }));
    } catch (err) {
      window.dispatchEvent(new CustomEvent('notifications:unread', { detail: { hasUnread: false } }));
    }
  };

  const loadNotifications = async () => {
    if (loading.value) return;
    loading.value = true;
    try {
      const resp = await listNotifications({
        page: currentPage.value,
        page_size: pageSize.value,
        status: activeTab.value === 'unread' ? 'unread' : undefined
      });
      notifications.value = resp.notifications || [];
      totalCount.value = Number(resp.total) || 0;
      selectedIds.value = [];
      if (activeTab.value === 'unread') {
        window.dispatchEvent(
          new CustomEvent('notifications:unread', { detail: { hasUnread: totalCount.value > 0 } })
        );
      }
    } catch (err) {
      notifications.value = [];
      totalCount.value = 0;
    } finally {
      loading.value = false;
    }
  };

  const markRead = async (row) => {
    if (!row?.external_id) return;
    try {
      await markNotificationsRead([row.external_id]);
      await loadNotifications();
      await refreshUnreadStatus();
    } catch (err) {
      ElMessage.error(t('common.requestFailed'));
    }
  };

  const handleMarkRead = async () => {
    try {
      if (selectedIds.value.length) {
        await markNotificationsRead(selectedIds.value);
      } else {
        await markAllNotificationsRead();
      }
      selectedIds.value = [];
      await loadNotifications();
      await refreshUnreadStatus();
    } catch (err) {
      ElMessage.error(t('common.requestFailed'));
    }
  };

  const handlePageChange = async () => {
    await loadNotifications();
  };

  const handleSizeChange = async () => {
    currentPage.value = 1;
    await loadNotifications();
  };

  watch(activeTab, async () => {
    currentPage.value = 1;
    await loadNotifications();
  });

  const init = async () => {
    await loadNotifications();
  };

  return {
    activeTab,
    totalCount,
    currentPage,
    pageSize,
    columns,
    displayNotifications,
    selectedIds,
    selectAll,
    selectIndeterminate,
    toggleSelectAll,
    tagType,
    tagLabel,
    buildTitle,
    formatDate,
    markRead,
    handleMarkRead,
    handlePageChange,
    handleSizeChange,
    init
  };
}
