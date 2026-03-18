<template>
  <PageLayout
    :system-name="systemName"
    :current-language="currentLanguage"
    :is-dark-mode="isDarkMode"
    :user-avatar="userAvatar"
    :username="userInfo?.username || '用户'"
    :active-menu="activeMenu"
    :side-menu-items="manageMenuItems"
    :title="t('system.security.title')"
    @change-language="handleLanguageChange"
    @toggle-theme="toggleTheme"
    @logout="handleLogout"
    @menu-select="handleMenuSelect"
  >
    <template #actions>
      <el-button class="page-action-btn" type="primary" size="small" :loading="isSaving" @click="handleSave">
        {{ t('common.save') }}
      </el-button>
    </template>
    <div class="manage-layout">
      <section class="manage-section">
        <div class="section-header">
          <h3 class="section-title">{{ t('system.security.registration') }}</h3>
        </div>
        <div class="section-body">
          <div class="setting-row setting-row--stacked">
            <div class="setting-label">{{ t('system.security.registrationMode') }}</div>
            <el-radio-group v-model="registrationMode">
              <el-radio value="open">{{ t('system.security.freeRegister') }}</el-radio>
              <el-radio value="invite">{{ t('system.security.inviteOnly') }}</el-radio>
            </el-radio-group>
          </div>
        </div>
      </section>

      <section class="manage-section">
        <div class="section-header">
          <h3 class="section-title">{{ t('system.security.placeholderTitle') }}</h3>
        </div>
        <div class="section-body">
          <el-alert
            type="info"
            :closable="false"
            :title="t('system.managementPlaceholder')"
            show-icon
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
import { House, Setting, Ticket } from '@element-plus/icons-vue';

const router = useRouter();
const userStore = useUserStore();
const { locale, t } = useI18n();

const systemName = ref('Yoresee');
const activeMenu = ref('manage-security');
const isDarkMode = ref(false);
const registrationMode = ref('open');
const isSaving = ref(false);

const userInfo = computed(() => userStore.userInfo);
const userAvatar = computed(() => userInfo.value?.avatar || 'https://cube.elemecdn.com/0/88/03b0d39583f48206768a7534e55bcpng.png');

const manageMenuItems = [
  { key: 'home', labelKey: 'navigation.home', icon: House, route: '/' },
  { key: 'manage-security', labelKey: 'system.menu.security', icon: Setting, route: '/manage/security' },
  { key: 'manage-invite', labelKey: 'system.menu.invite', icon: Ticket, route: '/manage/invitations' }
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

const initTheme = () => {
  const savedDarkMode = localStorage.getItem('darkMode');
  if (savedDarkMode === 'true') {
    isDarkMode.value = true;
    document.documentElement.classList.add('dark-mode');
  }
};

const initLanguage = () => {
  const savedLanguage = localStorage.getItem('language');
  if (savedLanguage) {
    currentLanguage.value = savedLanguage;
  }
};

const toggleTheme = () => {
  isDarkMode.value = !isDarkMode.value;
  if (isDarkMode.value) {
    document.documentElement.classList.add('dark-mode');
    localStorage.setItem('darkMode', 'true');
  } else {
    document.documentElement.classList.remove('dark-mode');
    localStorage.setItem('darkMode', 'false');
  }
};

const handleLogout = () => {
  userStore.logout();
  router.push('/login');
};

const handleMenuSelect = (key) => {
  activeMenu.value = key;
};

const handleSave = () => {
  if (isSaving.value) {
    return;
  }
  isSaving.value = true;
  setTimeout(() => {
    isSaving.value = false;
  }, 500);
};

onMounted(() => {
  initTheme();
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

.setting-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: var(--spacing-md);
}

.setting-row--stacked {
  align-items: flex-start;
  flex-direction: column;
}

.setting-label {
  color: var(--text-medium);
  font-size: 14px;
}

.setting-control {
  min-width: 200px;
}
</style>
