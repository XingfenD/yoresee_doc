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
    layout="list"
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

    <ManageLayout>
      <ManageSection>
        <InvitationCenter ref="invitationCenterRef" mode="system" :is-dark-mode="isDarkMode" />
      </ManageSection>
    </ManageLayout>
  </PageLayout>
</template>

<script setup>
import { onMounted, ref } from 'vue';
import { useI18n } from 'vue-i18n';
import { useRouter } from 'vue-router';
import { useUserStore } from '@/store/user';
import PageLayout from '@/components/PageLayout.vue';
import ManageLayout from '@/components/ManageLayout.vue';
import ManageSection from '@/components/ManageSection.vue';
import InvitationCenter from '@/components/InvitationCenter.vue';
import { useManageShell } from '@/composables/useManageShell';
import { usePageBoot } from '@/composables/usePageBoot';

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
  manageMenuItems,
  currentLanguage,
  initLanguage,
  fetchSystemInfo,
  handleLanguageChange,
  toggleTheme,
  handleLogout,
  handleMenuSelect
} = useManageShell({
  locale,
  router,
  userStore,
  defaultActiveMenu: 'manage-invite'
});
const { boot } = usePageBoot({ initLanguage, fetchSystemInfo });

const handleCreateClick = () => {
  invitationCenterRef.value?.openCreateDialog();
};

onMounted(() => {
  boot();
});
</script>
