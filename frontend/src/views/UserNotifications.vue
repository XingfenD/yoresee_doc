<template>
  <PageLayout
    :system-name="systemName"
    :current-language="currentLanguage"
    :is-dark-mode="isDarkMode"
    :user-avatar="userAvatar"
    :username="userInfo?.username || '用户'"
    :active-menu="activeMenu"
    :side-menu-items="userMenuItems"
    sidebar-scene="user_info"
    :title="t('user.notifications.title')"
    layout="list"
    @change-language="handleLanguageChange"
    @toggle-theme="toggleTheme"
    @logout="handleLogout"
    @menu-select="handleMenuSelect"
  >
    <el-tabs v-model="activeTab" class="common-tabs">
      <el-tab-pane :label="t('user.notifications.tabs.all')" name="all" />
      <el-tab-pane :label="t('user.notifications.tabs.unread')" name="unread" />
    </el-tabs>

    <CommonList
      :rows="displayNotifications"
      :columns="columns"
      :is-dark="isDarkMode"
      row-key="external_id"
      :empty-text="t('user.notifications.empty')"
      :show-title-bar="true"
      :title="t('user.notifications.title')"
      :show-pagination="true"
      :total="totalCount"
      v-model:current-page="currentPage"
      v-model:page-size="pageSize"
      :page-sizes="[6, 10, 20]"
      @page-change="handlePageChange"
      @size-change="handleSizeChange"
    >
      <template #toolbar-actions>
        <el-button size="small" type="primary" @click="handleMarkRead">
          {{ selectedIds.length ? t('user.notifications.markSelected') : t('user.notifications.markAll') }}
        </el-button>
      </template>
      <template #header-select>
        <el-checkbox
          :indeterminate="selectIndeterminate"
          :model-value="selectAll"
          @change="toggleSelectAll"
          class="checkbox-only"
        />
      </template>
      <template #cell-select="{ row }">
        <el-checkbox v-model="selectedIds" :label="row.external_id" class="checkbox-only" />
      </template>
      <template #cell-title="{ row }">
        <div class="notice-title" :class="{ 'is-unread': !row.read }">
          <span v-if="!row.read" class="notice-dot" />
          <div class="notice-text">
            <div class="notice-main">{{ buildTitle(row) }}</div>
            <div class="notice-sub">{{ row.content || '-' }}</div>
          </div>
        </div>
      </template>
      <template #cell-type="{ value }">
        <el-tag size="small" :type="tagType(value)">
          {{ tagLabel(value) }}
        </el-tag>
      </template>
      <template #cell-created_at="{ value }">
        {{ formatDate(value) }}
      </template>
      <template #cell-actions="{ row }">
        <el-button size="small" text type="primary" @click="markRead(row)">
          {{ t('user.notifications.markRead') }}
        </el-button>
        <el-button size="small" text type="primary" @click="openDetail(row)">
          {{ t('user.notifications.view') }}
        </el-button>
      </template>
    </CommonList>

  </PageLayout>
</template>

<script setup>
import { computed, ref, watch, onMounted } from 'vue';
import { useI18n } from 'vue-i18n';
import { useRouter } from 'vue-router';
import { useUserStore } from '@/store/user';
import PageLayout from '@/components/PageLayout.vue';
import CommonList from '@/components/CommonList.vue';
import { ElMessage } from 'element-plus';
import { listNotifications, markNotificationsRead, markAllNotificationsRead } from '@/services/api';
import { useUserShell } from '@/composables/useUserShell';

const router = useRouter();
const userStore = useUserStore();
const { locale, t } = useI18n();

const activeTab = ref('all');
const {
  systemName,
  activeMenu,
  isDarkMode,
  userInfo,
  userAvatar,
  userMenuItems,
  currentLanguage,
  initLanguage,
  handleLanguageChange,
  toggleTheme,
  handleLogout,
  handleMenuSelect
} = useUserShell({
  locale,
  router,
  userStore,
  defaultActiveMenu: 'user-notifications'
});

const notifications = ref([]);
const totalCount = ref(0);
const currentPage = ref(1);
const pageSize = ref(6);
const loading = ref(false);
const sending = ref(false);

const columns = computed(() => [
  { key: 'select', label: '', width: 36, align: 'center', className: 'select-column' },
  { key: 'title', label: t('user.notifications.columns.title'), minWidth: 320, flex: 1.6 },
  { key: 'type', label: t('user.notifications.columns.type'), minWidth: 120, align: 'center' },
  { key: 'created_at', label: t('user.notifications.columns.time'), minWidth: 180, align: 'center' },
  { key: 'actions', label: t('common.actions'), minWidth: 160, align: 'center' }
]);


const displayNotifications = computed(() => notifications.value);

const selectedIds = ref([]);

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
  } else {
    selectedIds.value = selectedIds.value.filter(
      (id) => !displayNotifications.value.some((item) => item.external_id === id)
    );
  }
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

const refreshUnreadStatus = async () => {
  try {
    const resp = await listNotifications({ page: 1, page_size: 1, status: 'unread' });
    const hasUnread = Number(resp.total) > 0;
    window.dispatchEvent(new CustomEvent('notifications:unread', { detail: { hasUnread } }));
  } catch (err) {
    window.dispatchEvent(new CustomEvent('notifications:unread', { detail: { hasUnread: false } }));
  }
};


watch(activeTab, async () => {
  currentPage.value = 1;
  await loadNotifications();
});

onMounted(() => {
  initLanguage();
  loadNotifications();
});
const openDetail = () => {
  // Sample only; would route to doc + comment anchor in real integration.
};
</script>

<style scoped>
.notice-title {
  display: flex;
  align-items: flex-start;
  gap: 10px;
}

.notice-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  margin-top: 6px;
  background: var(--primary-color);
  flex-shrink: 0;
}

.notice-text {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.notice-main {
  font-weight: 600;
  color: var(--text-dark);
}

.notice-title.is-unread .notice-main {
  color: var(--primary-color);
}

.notice-sub {
  font-size: 12px;
  color: var(--text-light);
}

.checkbox-only :deep(.el-checkbox__label) {
  display: none;
}

:deep(.select-column) {
  justify-content: center;
  padding-left: 6px;
  padding-right: 6px;
}

.dark-mode .notice-main {
  color: #e5e7eb;
}

.dark-mode .notice-title.is-unread .notice-main {
  color: #8ab4ff;
}

.dark-mode .notice-sub {
  color: #9aa4b2;
}
</style>
