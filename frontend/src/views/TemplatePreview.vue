<template>
  <PageLayout
    :system-name="systemName"
    :current-language="currentLanguage"
    :is-dark-mode="isDarkMode"
    :user-avatar="userAvatar"
    :username="userInfo?.username || t('common.unknown')"
    :active-menu="activeMenu"
    :title="''"
    content-padding="xl"
    @change-language="handleLanguageChange"
    @toggle-theme="toggleTheme"
    @logout="handleLogout"
    @menu-select="handleMenuSelect"
  >
    <TitleBar :show-back="true" :back-text="t('common.back')" @back="goBack">
      <template #actions>
        <el-button type="primary" @click="openCreateDocumentDialog">
          {{ t('document.createDocument') }}
        </el-button>
      </template>
    </TitleBar>

    <div class="template-preview" v-loading="loading">
      <div class="template-preview-header">
        <div class="template-preview-title">
          {{ template?.name || t('templates.untitled') }}
        </div>
        <div class="template-preview-meta">
          <el-tag v-if="scopeLabel" size="small" :type="scopeTagType">
            {{ scopeLabel }}
          </el-tag>
          <span class="template-preview-date" v-if="template?.updated_at || template?.updatedAt">
            {{ t('templates.updatedAt') }}: {{ formatDate(template?.updated_at || template?.updatedAt) }}
          </span>
        </div>
        <div v-if="template?.description" class="template-preview-desc">
          {{ template.description }}
        </div>
      </div>

      <div class="template-preview-body">
        <div v-if="!previewContent" class="template-preview-empty">
          <el-empty :description="t('templates.contentEmpty')" />
        </div>
        <div v-else ref="previewRef" class="template-preview-content"></div>
      </div>
    </div>
  </PageLayout>

  <DocumentCreateDialog
    v-model="showCreateDialog"
    :loading="creatingLoading"
    :initial-title="template?.name || ''"
    :initial-template-id="template?.id"
    :initial-template-meta="template || null"
    :knowledge-base-id="template?.knowledge_base_external_id || ''"
    :show-template-selector="false"
    @submit="createDocument"
    @cancel="cancelCreateDocument"
  />
</template>

<script setup>
import { onMounted } from 'vue';
import { useI18n } from 'vue-i18n';
import { useRoute, useRouter } from 'vue-router';
import 'vditor/dist/index.css';
import { useUserStore } from '@/store/user';
import PageLayout from '@/components/PageLayout.vue';
import TitleBar from '@/components/TitleBar.vue';
import DocumentCreateDialog from '@/components/DocumentCreateDialog.vue';
import { useWorkspaceShell } from '@/composables/useWorkspaceShell';
import { useTemplatePreviewPage } from '@/composables/useTemplatePreviewPage';

const route = useRoute();
const router = useRouter();
const { locale, t } = useI18n();
const userStore = useUserStore();

const {
  systemName,
  activeMenu,
  isDarkMode,
  currentLanguage,
  userInfo,
  userAvatar,
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
  loading,
  template,
  previewRef,
  showCreateDialog,
  creatingLoading,
  previewContent,
  scopeLabel,
  scopeTagType,
  formatDate,
  goBack,
  openCreateDocumentDialog,
  cancelCreateDocument,
  createDocument,
  init
} = useTemplatePreviewPage({
  t,
  route,
  router,
  isDarkMode
});

onMounted(async () => {
  await fetchSystemInfo();
  initLanguage();
  await init();
});
</script>

<style scoped>
:deep(.page-body) {
  gap: 0;
}

.template-preview {
  background: var(--bg-white);
  border-radius: var(--border-radius-md);
  box-shadow: var(--shadow-sm);
  border: 1px solid var(--border-color);
  padding: var(--spacing-lg);
}

.template-preview-header {
  padding-bottom: var(--spacing-md);
  border-bottom: 1px solid var(--border-color);
  margin-bottom: var(--spacing-md);
}

.template-preview-title {
  font-size: 22px;
  font-weight: 600;
  color: var(--text-primary);
  margin-bottom: 6px;
}

.template-preview-meta {
  display: flex;
  align-items: center;
  gap: 12px;
  color: var(--text-medium);
  font-size: 12px;
}

.template-preview-date {
  color: var(--text-light);
}

.template-preview-desc {
  margin-top: 8px;
  color: var(--text-medium);
}

.template-preview-body {
  min-height: 220px;
}

.template-preview-content {
  padding: 8px 0;
}

.template-preview-empty {
  padding: var(--spacing-lg) 0;
}

.dark-mode .template-preview {
  background: var(--bg-white);
}
</style>
