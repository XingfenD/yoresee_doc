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

    <InviteList :items="inviteList" :is-dark="isDarkMode" />

    <InviteCreateDialog v-model="showCreateDialog" @submit="handleCreateInvite" />
  </PageLayout>
</template>

<script setup>
import { computed, ref, onMounted, onBeforeUnmount } from 'vue';
import { useI18n } from 'vue-i18n';
import { useRouter } from 'vue-router';
import { useUserStore } from '@/store/user';
import PageLayout from '@/components/PageLayout.vue';
import InviteCreateDialog from '@/components/InviteCreateDialog.vue';
import InviteList from '@/components/InviteList.vue';
import { House, User, Ticket, Setting } from '@element-plus/icons-vue';

const router = useRouter();
const userStore = useUserStore();
const { locale, t } = useI18n();

const systemName = ref('Yoresee');
const activeMenu = ref('user-invite');
const isDarkMode = ref(document.documentElement.classList.contains('dark-mode'));
let classObserver = null;

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
  const next = !document.documentElement.classList.contains('dark-mode');
  document.documentElement.classList.toggle('dark-mode', next);
  localStorage.setItem('darkMode', next ? 'true' : 'false');
  isDarkMode.value = next;
};

const handleLogout = () => {
  userStore.logout();
  router.push('/login');
};

const handleMenuSelect = (key) => {
  activeMenu.value = key;
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

const sectionStyle = computed(() => ({}));

onMounted(() => {
  initLanguage();
  classObserver = new MutationObserver(() => {
    isDarkMode.value = document.documentElement.classList.contains('dark-mode');
  });
  classObserver.observe(document.documentElement, { attributes: true, attributeFilter: ['class'] });
});

onBeforeUnmount(() => {
  classObserver?.disconnect();
});
</script>

<style scoped>
</style>
