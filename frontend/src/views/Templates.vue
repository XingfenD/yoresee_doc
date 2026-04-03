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
      <el-tabs v-model="activeTab" class="common-tabs templates-tabs">
        <el-tab-pane :label="t('templates.my')" name="my" />
        <el-tab-pane :label="t('templates.recent')" name="recent" />
        <el-tab-pane :label="t('templates.public')" name="public" />
      </el-tabs>

      <CommonList
        v-loading="currentLoading"
        :rows="pagedTemplates"
        :columns="templateColumns"
        row-key="id"
        :is-dark="isDarkMode"
        :empty-text="currentEmptyText"
        :show-pagination="true"
        :total="paginationTotal"
        v-model:current-page="currentPage"
        :page-size="pageSize"
        :page-sizes="[9]"
        @page-change="handlePageChange"
        :show-search="true"
        v-model:search-query="keyword"
        :search-placeholder="t('templates.searchPlaceholder')"
        @search="handleSearch"
        :show-title-bar="true"
        :title="currentTitle"
      >
        <template #cell-name="{ row }">
          <span class="template-name">{{ row?.name || t('templates.untitled') }}</span>
        </template>
        <template #cell-description="{ row }">
          <span class="template-description">
            {{ row?.description || t('templates.noDescription') }}
          </span>
        </template>
        <template #cell-updated_at="{ row }">
          {{ formatDate(row?.updated_at || row?.updatedAt) }}
        </template>
        <template #cell-actions="{ row }">
          <el-button type="primary" text size="small" @click="openPreviewDialog(row)">
            {{ t('common.preview') }}
          </el-button>
          <el-button type="primary" text size="small" @click="openTemplate(row)">
            {{ t('common.open') }}
          </el-button>
          <el-button type="primary" text size="small" @click="openTemplateSettings(row)">
            {{ t('common.settings') }}
          </el-button>
        </template>
      </CommonList>
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

  <TemplatePreviewDialog
    v-model="showPreviewDialog"
    :title="`${t('templates.previewTitle')} · ${previewTitle}`"
    :content="previewContent"
    :is-dark-mode="isDarkMode"
    @closed="closePreviewDialog"
  />

  <TemplateSettingsDialog
    v-model="showSettingsDialog"
    :title="t('common.settings')"
    :loading="savingTemplateSettings"
    :form-state="templateSettingsForm"
    @submit="submitTemplateSettings"
  />
</template>

<script setup>
import { computed, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { useI18n } from 'vue-i18n';
import { useUserStore } from '@/store/user';
import PageLayout from '@/components/PageLayout.vue';
import CommonList from '@/components/CommonList.vue';
import TemplateCreateDialog from '@/components/TemplateCreateDialog.vue';
import TemplatePreviewDialog from '@/components/TemplatePreviewDialog.vue';
import TemplateSettingsDialog from '@/components/TemplateSettingsDialog.vue';
import { useWorkspaceShell } from '@/composables/useWorkspaceShell';
import { usePageBoot } from '@/composables/usePageBoot';
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
const { boot } = usePageBoot({ initLanguage, fetchSystemInfo });

const {
  activeTab,
  keyword,
  currentPage,
  pageSize,
  paginationTotal,
  currentLoading,
  pagedTemplates,
  showPreviewDialog,
  previewTitle,
  previewContent,
  showSettingsDialog,
  savingTemplateSettings,
  templateSettingsForm,
  showCreateDialog,
  creatingTemplate,
  templateDialogInit,
  formatDate,
  openTemplate,
  openPreviewDialog,
  closePreviewDialog,
  openTemplateSettings,
  submitTemplateSettings,
  handlePageChange,
  handleSearch,
  openCreateTemplateDialog,
  submitCreateTemplate,
  init
} = useTemplateListPage({
  t,
  router
});

const templateColumns = computed(() => [
  { key: 'name', label: t('templates.nameLabel'), minWidth: 180, flex: 1.2 },
  { key: 'description', label: t('templates.descLabel'), minWidth: 240, flex: 1.8 },
  { key: 'updated_at', label: t('templates.updatedAt'), minWidth: 160 },
  { key: 'actions', label: t('common.actions'), minWidth: 220, align: 'center' }
]);

const currentTitle = computed(() => {
  if (activeTab.value === 'recent') return t('templates.recent');
  if (activeTab.value === 'public') return t('templates.public');
  return t('templates.my');
});

const currentEmptyText = computed(() => {
  if (activeTab.value === 'recent') return t('templates.noRecent');
  if (activeTab.value === 'public') return t('templates.noPublic');
  return t('templates.noMy');
});

onMounted(async () => {
  await boot(init);
});
</script>

<style scoped>
.templates-layout {
  display: flex;
  flex-direction: column;
  gap: 12px;
  height: auto;
}

.templates-tabs {
  margin-bottom: 2px;
}

.template-name {
  font-weight: 600;
}

.template-description {
  color: var(--text-medium);
}
</style>
