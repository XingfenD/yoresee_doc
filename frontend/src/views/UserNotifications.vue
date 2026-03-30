<template>
  <PageLayout
    :system-name="systemName"
    :current-language="currentLanguage"
    :is-dark-mode="isDarkMode"
    :user-avatar="userAvatar"
    :username="userInfo?.username || t('common.user')"
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
        <AppTag size="small" :type="tagType(value)">
          {{ tagLabel(value) }}
        </AppTag>
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
import { onMounted } from 'vue';
import { useI18n } from 'vue-i18n';
import { useRouter } from 'vue-router';
import { useUserStore } from '@/store/user';
import PageLayout from '@/components/PageLayout.vue';
import CommonList from '@/components/CommonList.vue';
import AppTag from '@/components/AppTag.vue';
import { useUserShell } from '@/composables/useUserShell';
import { useNotificationCenter } from '@/composables/useNotificationCenter';
import { usePageBoot } from '@/composables/usePageBoot';

const router = useRouter();
const userStore = useUserStore();
const { locale, t } = useI18n();

const {
  systemName,
  activeMenu,
  isDarkMode,
  userInfo,
  userAvatar,
  userMenuItems,
  currentLanguage,
  initLanguage,
  fetchSystemInfo,
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
const { boot } = usePageBoot({ initLanguage, fetchSystemInfo });

const {
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
} = useNotificationCenter({ t });

onMounted(() => {
  boot(init);
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
