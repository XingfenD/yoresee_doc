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
      :rows="pagedNotifications"
      :columns="columns"
      :is-dark="isDarkMode"
      row-key="id"
      :empty-text="t('user.notifications.empty')"
      :show-title-bar="true"
      :title="t('user.notifications.title')"
      :show-pagination="true"
      :total="totalCount"
      v-model:current-page="currentPage"
      v-model:page-size="pageSize"
      :page-sizes="[5, 10, 20]"
      @page-change="handlePageChange"
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
        <el-checkbox v-model="selectedIds" :label="row.id" class="checkbox-only" />
      </template>
      <template #cell-title="{ row }">
        <div class="notice-title" :class="{ 'is-unread': !row.read }">
          <span v-if="!row.read" class="notice-dot" />
          <div class="notice-text">
            <div class="notice-main">{{ buildTitle(row) }}</div>
            <div class="notice-sub">{{ row.snippet || '-' }}</div>
          </div>
        </div>
      </template>
      <template #cell-type="{ value }">
        <el-tag size="small" :type="tagType(value)">
          {{ tagLabel(value) }}
        </el-tag>
      </template>
      <template #cell-time="{ value }">
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
import { computed, ref, watch } from 'vue';
import { useI18n } from 'vue-i18n';
import { useRouter } from 'vue-router';
import { useUserStore } from '@/store/user';
import PageLayout from '@/components/PageLayout.vue';
import CommonList from '@/components/CommonList.vue';
import { House, User, Ticket, Setting, Bell } from '@element-plus/icons-vue';

const router = useRouter();
const userStore = useUserStore();
const { locale, t } = useI18n();

const systemName = ref('Yoresee');
const activeMenu = ref('user-notifications');
const activeTab = ref('all');
const isDarkMode = computed(() => userStore.darkMode);

const userInfo = computed(() => userStore.userInfo);
const userAvatar = computed(() => userInfo.value?.avatar || 'https://cube.elemecdn.com/0/88/03b0d39583f48206768a7534e55bcpng.png');

const userMenuItems = [
  { key: 'home', labelKey: 'navigation.home', icon: House, route: '/' },
  { key: 'user-center', labelKey: 'user.menu.center', icon: User, route: '/user_info/example' },
  { key: 'user-notifications', labelKey: 'user.menu.notifications', icon: Bell, route: '/user_info/notifications' },
  { key: 'user-invite', labelKey: 'user.menu.invite', icon: Ticket, route: '/user_info/invatations' },
  { key: 'user-security', labelKey: 'user.menu.security', icon: Setting, route: '/user_info/example' }
];

const currentLanguage = computed({
  get: () => locale.value,
  set: (value) => {
    locale.value = value;
    localStorage.setItem('language', value);
  }
});

const handleLanguageChange = (command) => {
  currentLanguage.value = command;
};

const toggleTheme = () => {
  userStore.toggleDarkMode();
};

const handleLogout = () => {
  userStore.logout();
  router.push('/login');
};

const handleMenuSelect = (key) => {
  activeMenu.value = key;
};

const notifications = ref([
  {
    id: 'n1',
    type: 'mention',
    actor: '张三',
    docTitle: '产品需求文档',
    snippet: '请看第 3 节的评审结论',
    time: '2026-03-25T09:40:00Z',
    read: false
  },
  {
    id: 'n2',
    type: 'comment',
    actor: '李四',
    docTitle: '技术架构设计',
    snippet: '这部分是否可以拆分为两个服务？',
    time: '2026-03-24T18:22:00Z',
    read: false
  },
  {
    id: 'n3',
    type: 'reply',
    actor: '王五',
    docTitle: '会议纪要',
    snippet: '已补充行动项。',
    time: '2026-03-24T12:05:00Z',
    read: true
  }
]);

const columns = computed(() => [
  { key: 'select', label: '', width: 36, align: 'center', className: 'select-column' },
  { key: 'title', label: t('user.notifications.columns.title'), minWidth: 320, flex: 1.6 },
  { key: 'type', label: t('user.notifications.columns.type'), minWidth: 120, align: 'center' },
  { key: 'time', label: t('user.notifications.columns.time'), minWidth: 180, align: 'center' },
  { key: 'actions', label: t('common.actions'), minWidth: 160, align: 'center' }
]);

const displayNotifications = computed(() => {
  if (activeTab.value === 'unread') {
    return notifications.value.filter((item) => !item.read);
  }
  return notifications.value;
});

const selectedIds = ref([]);
const currentPage = ref(1);
const pageSize = ref(10);

const totalCount = computed(() => displayNotifications.value.length);

const pagedNotifications = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value;
  return displayNotifications.value.slice(start, start + pageSize.value);
});

const selectAll = computed(() => {
  if (!pagedNotifications.value.length) return false;
  return pagedNotifications.value.every((item) => selectedIds.value.includes(item.id));
});

const selectIndeterminate = computed(() => {
  const selectedCount = pagedNotifications.value.filter((item) => selectedIds.value.includes(item.id)).length;
  return selectedCount > 0 && selectedCount < pagedNotifications.value.length;
});

const toggleSelectAll = (value) => {
  if (value) {
    selectedIds.value = Array.from(
      new Set([...selectedIds.value, ...pagedNotifications.value.map((item) => item.id)])
    );
  } else {
    selectedIds.value = selectedIds.value.filter(
      (id) => !pagedNotifications.value.some((item) => item.id === id)
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
  return t('user.notifications.types.comment');
};

const buildTitle = (row) => `${row.actor} ${t('user.notifications.inDoc')}「${row.docTitle}」`;

const formatDate = (value) => {
  if (!value) return '-';
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) return value;
  return date.toLocaleString();
};

const markRead = (row) => {
  const target = notifications.value.find((item) => item.id === row.id);
  if (target) {
    target.read = true;
  }
};

const markSelectedRead = () => {
  notifications.value = notifications.value.map((item) =>
    selectedIds.value.includes(item.id) ? { ...item, read: true } : item
  );
  selectedIds.value = [];
};

const markAllRead = () => {
  notifications.value = notifications.value.map((item) => ({ ...item, read: true }));
  selectedIds.value = [];
};

const handleMarkRead = () => {
  if (selectedIds.value.length) {
    markSelectedRead();
  } else {
    markAllRead();
  }
};

const handlePageChange = () => {
  // client-side pagination
};

watch(activeTab, () => {
  currentPage.value = 1;
  selectedIds.value = [];
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
