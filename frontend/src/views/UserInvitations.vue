<template>
  <PageLayout
    :system-name="systemName"
    :current-language="currentLanguage"
    :is-dark-mode="isDarkMode"
    :user-avatar="userAvatar"
    :username="userInfo?.username || '用户'"
    :active-menu="activeMenu"
    :side-menu-items="userMenuItems"
    :title="t('user.invite.title')"
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

    <InviteCreateDialog v-model="showCreateDialog" @submit="handleCreateInvite" />
  </PageLayout>
</template>

<script setup>
import { computed, ref, onMounted } from 'vue';
import { useI18n } from 'vue-i18n';
import { useRouter } from 'vue-router';
import { useUserStore } from '@/store/user';
import PageLayout from '@/components/PageLayout.vue';
import InviteCreateDialog from '@/components/InviteCreateDialog.vue';
import CommonList from '@/components/CommonList.vue';
import { House, User, Ticket, Setting } from '@element-plus/icons-vue';

const router = useRouter();
const userStore = useUserStore();
const { locale, t } = useI18n();

const systemName = ref('Yoresee');
const activeMenu = ref('user-invite');
const isDarkMode = computed(() => userStore.darkMode);

const userInfo = computed(() => userStore.userInfo);
const userAvatar = computed(() => userInfo.value?.avatar || 'https://cube.elemecdn.com/0/88/03b0d39583f48206768a7534e55bcpng.png');

const userMenuItems = [
  { key: 'home', labelKey: 'navigation.home', icon: House, route: '/' },
  { key: 'user-center', labelKey: 'user.menu.center', icon: User, route: '/user_info/example' },
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

const initLanguage = () => {
  const savedLanguage = localStorage.getItem('language');
  if (savedLanguage) {
    currentLanguage.value = savedLanguage;
  }
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
  { key: 'expires_at', label: t('user.invite.expiresAt'), minWidth: 160 },
  { key: 'actions', label: t('user.invite.actions'), minWidth: 160, align: 'right' }
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
    code: 'YORE-8K2P-9Q1M',
    created_at: '2026-03-18 10:12',
    created_by: 'admin',
    expires_at: '2026-04-18 23:59',
    status: 'active',
    max: 20,
    used: 6
  },
  {
    code: 'YORE-7D4X-1ZA3',
    created_at: '2026-03-10 09:30',
    created_by: 'admin',
    expires_at: '2026-03-31 23:59',
    status: 'active',
    max: 5,
    used: 5
  },
  {
    code: 'YORE-2M9K-4L8N',
    created_at: '2026-02-20 15:40',
    created_by: 'admin',
    expires_at: '2026-03-05 23:59',
    status: 'expired',
    max: 10,
    used: 7
  }
]);


const showCreateDialog = ref(false);

const openCreateDialog = () => {
  showCreateDialog.value = true;
};

const handleCreateInvite = (payload) => {
  // TODO: hook to backend
  console.log('create invite payload', payload);
};

const handlePauseInvite = (row) => {
  console.log('pause invite', row);
};

const handleDeleteInvite = (row) => {
  console.log('delete invite', row);
};

onMounted(() => {
  initLanguage();
});
</script>

<style scoped>
</style>
