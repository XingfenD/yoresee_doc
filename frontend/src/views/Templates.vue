<template>
  <PageLayout
    :system-name="systemName"
    :current-language="currentLanguage"
    :is-dark-mode="isDarkMode"
    :user-avatar="userAvatar"
    :username="userInfo?.username || t('common.unknown')"
    :active-menu="activeMenu"
    :title="t('templates.title')"
    layout="list"
    @change-language="handleLanguageChange"
    @toggle-theme="toggleTheme"
    @logout="handleLogout"
    @menu-select="handleMenuSelect"
  >
    <template #actions>
      <el-button class="page-action-btn" type="primary" size="small" @click="openCreateTemplateDialog">
        {{ t('templates.createNew') }}
      </el-button>
    </template>
    <div class="templates-layout">
      <el-tabs v-model="activeTab" class="common-tabs">
        <el-tab-pane :label="t('templates.my')" name="my">
          <div v-loading="loadingMy">
            <TemplateListSection
              :title="t('templates.my')"
              :items="myTemplates"
              :empty-text="t('templates.noMy')"
              :fallback-description="t('templates.noDescription')"
              :tag-mapper="myTagMapper"
              :meta-mapper="templateMetaMapper"
              :action-label="t('common.open')"
              @open="openTemplate"
            />
          </div>
        </el-tab-pane>

        <el-tab-pane :label="t('templates.recent')" name="recent">
          <div v-loading="loadingRecent">
            <TemplateListSection
              :title="t('templates.recent')"
              :items="recentTemplates"
              :empty-text="t('templates.noRecent')"
              :fallback-description="t('templates.noDescription')"
              :tag-mapper="recentTagMapper"
              :meta-mapper="templateMetaMapper"
              :action-label="t('common.open')"
              @open="openTemplate"
            />
          </div>
        </el-tab-pane>

        <el-tab-pane :label="t('templates.public')" name="public">
          <div v-loading="loadingPublic">
            <TemplateListSection
              :title="t('templates.public')"
              :items="publicTemplates"
              :empty-text="t('templates.noPublic')"
              :fallback-description="t('templates.noDescription')"
              :tag-mapper="publicTagMapper"
              :meta-mapper="templateMetaMapper"
              :action-label="t('common.open')"
              @open="openTemplate"
            />
          </div>
        </el-tab-pane>
      </el-tabs>
    </div>
  </PageLayout>

  <TemplateCreateDialog
    v-model="showCreateDialog"
    :loading="creatingTemplate"
    :title="t('templates.createNewTitle')"
    :show-content="true"
    :show-kb-scope="false"
    :initial-name="templateDialogInit.name"
    :initial-description="templateDialogInit.description"
    :initial-scope="templateDialogInit.scope"
    :initial-tags="templateDialogInit.tags"
    :initial-content="templateDialogInit.content"
    @submit="submitCreateTemplate"
  />
</template>

<script setup>
import { onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { useI18n } from 'vue-i18n';
import { useUserStore } from '@/store/user';
import PageLayout from '@/components/PageLayout.vue';
import TemplateListSection from '@/components/TemplateListSection.vue';
import TemplateCreateDialog from '@/components/TemplateCreateDialog.vue';
import { useWorkspaceShell } from '@/composables/useWorkspaceShell';
import { useTemplateListPage } from '@/composables/useTemplateListPage';

const router = useRouter();
const userStore = useUserStore();
const { locale, t } = useI18n();

const {
  systemName,
  activeMenu,
  isDarkMode,
  userInfo,
  userAvatar,
  currentLanguage,
  initLanguage,
  handleLanguageChange,
  toggleTheme,
  handleLogout,
  handleMenuSelect,
  fetchSystemInfo
} = useWorkspaceShell({
  locale,
  router,
  userStore,
  defaultActiveMenu: 'templates'
});

const {
  activeTab,
  showCreateDialog,
  creatingTemplate,
  templateDialogInit,
  myTemplates,
  recentTemplates,
  publicTemplates,
  loadingMy,
  loadingRecent,
  loadingPublic,
  myTagMapper,
  recentTagMapper,
  publicTagMapper,
  templateMetaMapper,
  openTemplate,
  openCreateTemplateDialog,
  submitCreateTemplate,
  init
} = useTemplateListPage({
  t,
  router
});

onMounted(async () => {
  await fetchSystemInfo();
  initLanguage();
  await init();
});
</script>

<style scoped>
.templates-layout {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-lg);
  height: auto;
}

</style>
