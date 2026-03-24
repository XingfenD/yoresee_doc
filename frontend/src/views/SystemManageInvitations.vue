<template>
  <PageLayout
    :system-name="systemName"
    :current-language="currentLanguage"
    :is-dark-mode="isDarkMode"
    :user-avatar="userAvatar"
    :username="userInfo?.username || '用户'"
    :active-menu="activeMenu"
    :side-menu-items="manageMenuItems"
    sidebar-scene="manage"
    :title="t('system.invite.title')"
    @change-language="handleLanguageChange"
    @toggle-theme="toggleTheme"
    @logout="handleLogout"
    @menu-select="handleMenuSelect"
  >
    <template #actions>
      <el-button class="page-action-btn" type="primary" size="small" @click="openCreateDialog">
        {{ t('user.invite.create') }}
      </el-button>
    </template>

    <el-tabs v-model="activeTab" class="manage-tabs">
      <el-tab-pane :label="t('system.invite.tabs.list')" name="list">
        <CommonList
          :rows="inviteList"
          :columns="inviteColumns"
          :is-dark="isDarkMode"
          row-key="code"
          :empty-text="t('message.empty')"
        >
          <template #cell-status="{ value }">
            <el-tag :type="inviteStatusType(value)" size="small">
              {{ inviteStatusLabel(value) }}
            </el-tag>
          </template>
          <template #cell-usage="{ row }">
            {{ row.used }}/{{ row.max }}
          </template>
          <template #cell-actions="{ row }">
            <el-button size="small" text type="primary" @click="handlePauseInvite(row)">
              {{ t('user.invite.pause') }}
            </el-button>
            <el-button size="small" text type="danger" @click="handleDeleteInvite(row)">
              {{ t('user.invite.delete') }}
            </el-button>
          </template>
        </CommonList>
      </el-tab-pane>
      <el-tab-pane :label="t('system.invite.tabs.records')" name="records">
        <CommonList
          :rows="inviteRecords"
          :columns="recordColumns"
          :is-dark="isDarkMode"
          row-key="id"
          :empty-text="t('system.invite.records.empty')"
        >
          <template #cell-status="{ value }">
            <el-tag :type="value === 'success' ? 'success' : 'warning'" size="small">
              {{ value === 'success' ? t('system.invite.records.success') : t('system.invite.records.failed') }}
            </el-tag>
          </template>
        </CommonList>
      </el-tab-pane>
    </el-tabs>

    <InviteCreateDialog v-model="showCreateDialog" @submit="handleCreateInvite" />
  </PageLayout>
</template>

<script setup>
import { computed, ref, onMounted } from 'vue';
import { useI18n } from 'vue-i18n';
import { useRouter } from 'vue-router';
import { useUserStore } from '@/store/user';
import PageLayout from '@/components/PageLayout.vue';
import CommonList from '@/components/CommonList.vue';
import InviteCreateDialog from '@/components/InviteCreateDialog.vue';
import { House, Setting, Ticket, User, UserFilled, OfficeBuilding } from '@element-plus/icons-vue';

const router = useRouter();
const userStore = useUserStore();
const { locale, t } = useI18n();

const systemName = ref('Yoresee');
const activeMenu = ref('manage-invite');
const activeTab = ref('list');
const isDarkMode = computed(() => userStore.darkMode);

const userInfo = computed(() => userStore.userInfo);
const userAvatar = computed(() => userInfo.value?.avatar || 'https://cube.elemecdn.com/0/88/03b0d39583f48206768a7534e55bcpng.png');

const manageMenuItems = [
  { key: 'home', labelKey: 'navigation.home', icon: House, route: '/' },
  { key: 'manage-user', labelKey: 'system.menu.user', icon: User, route: '/manage/user' },
  { key: 'manage-user-group', labelKey: 'system.menu.userGroup', icon: UserFilled, route: '/manage/user_group' },
  { key: 'manage-organization', labelKey: 'system.menu.organization', icon: OfficeBuilding, route: '/manage/organization' },
  { key: 'manage-invite', labelKey: 'system.menu.invite', icon: Ticket, route: '/manage/invitations' },
  { key: 'manage-security', labelKey: 'system.menu.security', icon: Setting, route: '/manage/security' }
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

const inviteColumns = computed(() => [
  { key: 'code', label: t('user.invite.code'), minWidth: 180 },
  { key: 'status', label: t('user.invite.status'), minWidth: 120, align: 'center' },
  { key: 'usage', label: t('user.invite.usage'), minWidth: 110, align: 'center' },
  { key: 'created_at', label: t('user.invite.createdAt'), minWidth: 160 },
  { key: 'created_by', label: t('user.invite.createdBy'), minWidth: 140 },
  { key: 'expires_at', label: t('user.invite.expiresAt'), minWidth: 160 },
  { key: 'actions', label: t('user.invite.actions'), minWidth: 160, align: 'right' }
]);

const recordColumns = computed(() => [
  { key: 'code', label: t('system.invite.records.code'), minWidth: 180 },
  { key: 'used_by', label: t('system.invite.records.usedBy'), minWidth: 160 },
  { key: 'used_at', label: t('system.invite.records.usedAt'), minWidth: 180 },
  { key: 'status', label: t('system.invite.records.result'), minWidth: 120, align: 'center' }
]);

const inviteStatusType = (status) => {
  if (status === 'active') return 'success';
  if (status === 'expired') return 'info';
  return 'warning';
};

const inviteStatusLabel = (status) => {
  if (status === 'active') return t('user.invite.active');
  if (status === 'expired') return t('user.invite.expired');
  return t('user.invite.disabled');
};

const inviteList = ref([
  {
    code: 'SYS-9KD2-8LMQ',
    created_at: '2026-03-19 09:12',
    created_by: 'admin',
    expires_at: '2026-05-01 23:59',
    status: 'active',
    max: 50,
    used: 12
  },
  {
    code: 'SYS-3PA9-2XKD',
    created_at: '2026-03-01 11:20',
    created_by: 'admin',
    expires_at: '2026-03-25 23:59',
    status: 'expired',
    max: 5,
    used: 3
  }
]);

const inviteRecords = ref([
  {
    id: 1,
    code: 'SYS-9KD2-8LMQ',
    used_by: 'user_lee',
    used_at: '2026-03-19 10:01',
    status: 'success'
  },
  {
    id: 2,
    code: 'SYS-9KD2-8LMQ',
    used_by: 'user_wang',
    used_at: '2026-03-19 10:05',
    status: 'failed'
  }
]);

const showCreateDialog = ref(false);

const openCreateDialog = () => {
  showCreateDialog.value = true;
};

const handleCreateInvite = (payload) => {
  // TODO: hook to backend
  console.log('create system invite payload', payload);
};

const handlePauseInvite = (row) => {
  console.log('pause invite', row);
};

const handleDeleteInvite = (row) => {
  console.log('delete invite', row);
};

const initLanguage = () => {
  const savedLanguage = localStorage.getItem('language');
  if (savedLanguage) {
    currentLanguage.value = savedLanguage;
  }
};

onMounted(() => {
  initLanguage();
});
</script>

<style scoped>
.manage-tabs :deep(.el-tabs__header) {
  margin-bottom: var(--spacing-lg);
}
</style>
