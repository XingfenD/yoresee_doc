<template>
  <PageLayout
    :system-name="systemName"
    :current-language="currentLanguage"
    :is-dark-mode="isDarkMode"
    :user-avatar="userAvatar"
    :username="userInfo?.username || '用户'"
    :active-menu="activeMenu"
    :side-menu-items="manageMenuItems"
    :title="t('system.user.title')"
    @change-language="handleLanguageChange"
    @toggle-theme="toggleTheme"
    @logout="handleLogout"
    @menu-select="handleMenuSelect"
  >
    <div class="manage-layout">
      <section class="manage-section">
        <div class="section-header">
          <h3 class="section-title">{{ t('system.user.placeholderTitle') }}</h3>
        </div>
        <div class="section-body">
          <CommonList
            :rows="userRows"
            :columns="userColumns"
            :is-dark="isDarkMode"
            row-key="email"
            :empty-text="t('message.empty')"
          />
        </div>
      </section>
    </div>
  </PageLayout>
</template>

<script setup>
import { computed, ref, onMounted } from 'vue';
import { useI18n } from 'vue-i18n';
import { useRouter } from 'vue-router';
import { useUserStore } from '@/store/user';
import PageLayout from '@/components/PageLayout.vue';
import CommonList from '@/components/CommonList.vue';
import { House, Setting, Ticket, User, UserFilled, OfficeBuilding } from '@element-plus/icons-vue';

const router = useRouter();
const userStore = useUserStore();
const { locale, t } = useI18n();

const systemName = ref('Yoresee');
const activeMenu = ref('manage-user');
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

const userRows = ref([
  { name: 'Alex Chen', email: 'alex.chen@yoresee.com', role: 'Owner', status: 'Active', created_at: '2026-03-01' },
  { name: 'Mia Liu', email: 'mia.liu@yoresee.com', role: 'Admin', status: 'Active', created_at: '2026-03-05' },
  { name: 'Sam Zhao', email: 'sam.zhao@yoresee.com', role: 'Editor', status: 'Active', created_at: '2026-03-10' },
  { name: 'Lina Wu', email: 'lina.wu@yoresee.com', role: 'Viewer', status: 'Invited', created_at: '2026-03-12' },
  { name: 'Eric Sun', email: 'eric.sun@yoresee.com', role: 'Editor', status: 'Suspended', created_at: '2026-03-15' }
]);

const userColumns = computed(() => [
  { key: 'name', label: t('user.name'), minWidth: 140 },
  { key: 'email', label: t('user.email'), minWidth: 220, flex: 1.4 },
  { key: 'role', label: t('user.role'), minWidth: 120 },
  { key: 'status', label: t('user.status'), minWidth: 120 },
  { key: 'created_at', label: t('common.createdAt'), minWidth: 140 }
]);

onMounted(() => {
  initLanguage();
});
</script>

<style scoped>
.manage-layout {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-lg);
}

.manage-section {
  background: var(--bg-white);
  border: 1px solid var(--border-color);
  border-radius: var(--border-radius-md);
  box-shadow: var(--shadow-sm);
  overflow: hidden;
}

.section-header {
  padding: var(--spacing-md);
  border-bottom: 1px solid var(--border-color);
  background: var(--bg-white);
}

.section-title {
  margin: 0;
  font-size: 16px;
  font-weight: 600;
  color: var(--text-dark);
}

.section-body {
  padding: var(--spacing-md);
}

.dark-mode .manage-section {
  background: #161b22;
  border-color: #2b2f36;
}

.dark-mode .section-header {
  background: #161b22;
  border-bottom-color: #2b2f36;
}
</style>
