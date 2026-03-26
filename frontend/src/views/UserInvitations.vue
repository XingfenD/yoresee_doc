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
    :title="t('user.invite.title')"
    @change-language="handleLanguageChange"
    @toggle-theme="toggleTheme"
    @logout="handleLogout"
    @menu-select="handleMenuSelect"
  >
    <template #actions>
      <el-button class="page-action-btn" type="primary" size="small" @click="handleCreateClick">
        {{ t('user.invite.create') }}
      </el-button>
    </template>

    <InvitationCenter ref="invitationCenterRef" mode="user" :is-dark-mode="isDarkMode" />
  </PageLayout>
</template>

<script setup>
import { onMounted, ref } from 'vue';
import { useI18n } from 'vue-i18n';
import { useRouter } from 'vue-router';
import { useUserStore } from '@/store/user';
import PageLayout from '@/components/PageLayout.vue';
import InvitationCenter from '@/components/InvitationCenter.vue';
import { useUserShell } from '@/composables/useUserShell';

const router = useRouter();
const userStore = useUserStore();
const { locale, t } = useI18n();
const invitationCenterRef = ref(null);

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
  defaultActiveMenu: 'user-invite'
});

const handleCreateClick = () => {
  invitationCenterRef.value?.openCreateDialog();
};

onMounted(() => {
  initLanguage();
});
</script>
