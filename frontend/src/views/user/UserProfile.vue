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
    :title="t('user.profileTitle')"
    @change-language="handleLanguageChange"
    @toggle-theme="toggleTheme"
    @logout="handleLogout"
    @menu-select="handleMenuSelect"
  >
    <div class="user-profile-page">
      <div class="user-card">
        <div class="user-card-header">
          <el-avatar :size="64" :src="userAvatar" class="profile-avatar" />
          <div class="user-card-title">
            <div class="user-name">{{ userInfo?.username || t('common.user') }}</div>
            <div class="user-subtitle">{{ t('user.profile') }}</div>
          </div>
          <el-button class="setting-btn" type="primary" plain @click="goToSetting">
            {{ t('user.editAccount') }}
          </el-button>
        </div>
        <div class="user-card-body">
          <div class="info-row">
            <span class="info-label">{{ t('user.name') }}</span>
            <span class="info-value">{{ userInfo?.username || '-' }}</span>
          </div>
          <div class="info-row">
            <span class="info-label">{{ t('user.email') }}</span>
            <span class="info-value">{{ userInfo?.email || '-' }}</span>
          </div>
          <div class="info-row">
            <span class="info-label">{{ t('user.nickname') }}</span>
            <span class="info-value">{{ userInfo?.nickname || '-' }}</span>
          </div>
          <div class="info-row">
            <span class="info-label">{{ t('user.createdAt') }}</span>
            <span class="info-value">{{ formatDate(userInfo?.created_at) }}</span>
          </div>
        </div>
      </div>
    </div>
  </PageLayout>
</template>

<script setup>
import { onMounted } from 'vue';
import { useI18n } from 'vue-i18n';
import { useRouter } from 'vue-router';
import { useUserStore } from '@/store/user';
import PageLayout from '@/components/layout/PageLayout.vue';
import { useUserShell } from '@/composables/shell/useUserShell';
import { usePageBoot } from '@/composables/shell/usePageBoot';

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
  defaultActiveMenu: 'user-center'
});
const { boot } = usePageBoot({ initLanguage, fetchSystemInfo });

const goToSetting = () => {
  router.push('/user_info/setting');
};

const formatDate = (value) => {
  if (!value) return '-';
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) return value;
  return date.toLocaleString();
};

onMounted(() => {
  boot();
});
</script>

<style scoped>
.user-profile-page {
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

.profile-avatar {
  flex-shrink: 0;
}

.user-card-title {
  display: flex;
  flex-direction: column;
  gap: 4px;
  min-width: 0;
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

.setting-btn {
  margin-left: auto;
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
  max-width: 60%;
  text-align: right;
  word-break: break-all;
}
</style>
