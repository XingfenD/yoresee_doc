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
    :title="t('user.center')"
    @change-language="handleLanguageChange"
    @toggle-theme="toggleTheme"
    @logout="handleLogout"
    @menu-select="handleMenuSelect"
  >
    <div class="user-center">
      <div v-if="activeMenu === 'user-center'" class="user-card">
        <div class="user-card-header">
          <el-avatar :size="56" :src="userAvatar" />
          <div class="user-card-title">
            <div class="user-name">{{ userInfo?.username || '用户' }}</div>
            <div class="user-subtitle">{{ t('user.profile') }}</div>
          </div>
        </div>
        <div class="user-card-body">
          <div class="info-row">
            <span class="info-label">{{ t('user.basicInfo') }}</span>
            <span class="info-value">{{ userInfo?.email || '-' }}</span>
          </div>
          <div class="info-row">
            <span class="info-label">{{ t('user.account') }}</span>
            <span class="info-value">{{ userInfo?.username || '-' }}</span>
          </div>
          <div class="info-row">
            <span class="info-label">{{ t('user.security') }}</span>
            <span class="info-value">-</span>
          </div>
        </div>
      </div>

      <div v-else class="user-placeholder">
        <el-alert
          type="info"
          :closable="false"
          :title="t('user.placeholder')"
          show-icon
        />
      </div>
    </div>
  </PageLayout>
</template>

<script setup>
import { computed, ref, onMounted } from 'vue';
import { useI18n } from 'vue-i18n';
import { useRouter } from 'vue-router';
import { useUserStore } from '@/store/user';
import PageLayout from '@/components/PageLayout.vue';
import { House, User, Ticket, Setting } from '@element-plus/icons-vue';

const router = useRouter();
const userStore = useUserStore();
const { locale, t } = useI18n();

const systemName = ref('Yoresee');
const activeMenu = ref('user-center');
const isDarkMode = computed(() => userStore.darkMode);

const userInfo = computed(() => userStore.userInfo);
const userAvatar = computed(() => userInfo.value?.avatar || 'https://cube.elemecdn.com/0/88/03b0d39583f48206768a7534e55bcpng.png');

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

const userMenuItems = [
  { key: 'home', labelKey: 'navigation.home', icon: House, route: '/' },
  { key: 'user-center', labelKey: 'user.menu.center', icon: User, route: '/user_info/example' },
  { key: 'user-invite', labelKey: 'user.menu.invite', icon: Ticket, route: '/user_info/invatations' },
  { key: 'user-security', labelKey: 'user.menu.security', icon: Setting, route: '/user_info/example' }
];

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


onMounted(() => {
  initLanguage();
});
</script>

<style scoped>
.user-center {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-lg);
}

.user-card {
  background: var(--bg-white);
  border: 1px solid var(--border-color);
  border-radius: var(--border-radius-md);
  padding: var(--spacing-lg);
  box-shadow: var(--shadow-sm);
}

.user-card-header {
  display: flex;
  align-items: center;
  gap: var(--spacing-md);
  margin-bottom: var(--spacing-md);
}

.user-card-title {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.user-name {
  font-size: 18px;
  font-weight: 600;
  color: var(--text-dark);
}

.user-subtitle {
  font-size: 13px;
  color: var(--text-light);
}

.user-card-body {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-sm);
}

.info-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: var(--spacing-xs) 0;
  border-bottom: 1px dashed var(--border-color);
}

.info-row:last-child {
  border-bottom: none;
}

.info-label {
  color: var(--text-light);
  font-size: 13px;
}

.info-value {
  color: var(--text-dark);
  font-size: 14px;
}

.user-placeholder {
  display: flex;
  flex-direction: column;
}
</style>
