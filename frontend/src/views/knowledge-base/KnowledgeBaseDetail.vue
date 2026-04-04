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
    <TitleBar :show-back="true" :compact="true" :back-text="t('common.back')" @back="goBackToKnowledgeBase">
      <template #actions>
        <DocumentTypeMenu @select="openCreateDocumentDialog">
          <el-button type="primary">
            {{ t('knowledgeBase.createDocument') }}
          </el-button>
        </DocumentTypeMenu>
      </template>
    </TitleBar>

    <div class="detail-content" v-loading="loading">
      <InfoStatsCard
        v-if="knowledgeBaseData"
        :title="knowledgeBaseName"
        :description="knowledgeBaseDescription"
        :stats="knowledgeBaseStats"
      />
      <div v-else-if="!loading" class="empty-state">
        <el-empty :description="t('message.empty')" />
      </div>

      <div class="detail-columns">
        <KnowledgeBaseDocumentTreePanel
          class="detail-tree-panel"
          v-model:search-keyword="searchKeyword"
          v-model:sort-by="sortBy"
          :title="t('knowledgeBase.documentStructure')"
          :search-placeholder="t('knowledgeBase.searchDocuments')"
          :sort-placeholder="t('knowledgeBase.sortBy')"
          :sort-options="sortOptions"
          :nodes="directoryTreeData"
          :loading="loading"
          :empty-text="t('knowledgeBase.noDocuments')"
          :total="totalDocumentsCount"
          :current-page="currentPage"
          :page-size="pageSize"
          :open-label="t('common.open')"
          :rename-label="t('common.rename')"
          :share-label="t('document.share')"
          :delete-label="t('common.delete')"
          @node-click="handleTreeNodeClick"
          @open="openDocument"
          @node-action="handleNodeAction"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />

        <KnowledgeBaseTemplatesPanel
          class="detail-templates-panel"
          :loading="kbTemplatesLoading"
          :title="t('knowledgeBase.templates')"
          :items="kbTemplates"
          :empty-text="t('templates.noMy')"
          :fallback-description="t('templates.noDescription')"
          :tag-mapper="templateTagMapper"
          :meta-mapper="templateMetaMapper"
          :action-label="t('common.open')"
          @open="openTemplate"
        />
      </div>
    </div>
  </PageLayout>

  <DocumentCreateDialog
    v-model="showCreateDialog"
    :loading="creatingLoading"
    :initial-document-type="selectedDocumentType"
    :knowledge-base-id="route.params.id || ''"
    @submit="createDocument"
    @cancel="cancelCreateDocument"
  />
</template>

<script setup>
import { onMounted } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import { useUserStore } from '@/store/user';
import { useI18n } from 'vue-i18n';
import PageLayout from '@/components/layout/PageLayout.vue';
import TitleBar from '@/components/layout/TitleBar.vue';
import DocumentCreateDialog from '@/components/document/DocumentCreateDialog.vue';
import DocumentTypeMenu from '@/components/document/DocumentTypeMenu.vue';
import InfoStatsCard from '@/components/shared/InfoStatsCard.vue';
import KnowledgeBaseDocumentTreePanel from '@/components/knowledge-base/KnowledgeBaseDocumentTreePanel.vue';
import KnowledgeBaseTemplatesPanel from '@/components/knowledge-base/KnowledgeBaseTemplatesPanel.vue';
import { useKnowledgeBaseDetailPage } from '@/composables/knowledge-base/useKnowledgeBaseDetailPage';

const { locale, t } = useI18n();
const router = useRouter();
const route = useRoute();
const userStore = useUserStore();

const {
  systemName,
  activeMenu,
  isDarkMode,
  currentLanguage,
  userInfo,
  userAvatar,
  knowledgeBaseName,
  knowledgeBaseDescription,
  knowledgeBaseData,
  loading,
  kbTemplates,
  kbTemplatesLoading,
  knowledgeBaseStats,
  searchKeyword,
  sortBy,
  sortOptions,
  currentPage,
  pageSize,
  totalDocumentsCount,
  showCreateDialog,
  creatingLoading,
  selectedDocumentType,
  createDocument,
  cancelCreateDocument,
  openCreateDocumentDialog,
  handleTreeNodeClick,
  openDocument,
  handleNodeAction,
  handleSizeChange,
  handleCurrentChange,
  handleMenuSelect,
  handleLanguageChange,
  toggleTheme,
  handleLogout,
  goBackToKnowledgeBase,
  templateTagMapper,
  templateMetaMapper,
  openTemplate,
  directoryTreeData,
  init
} = useKnowledgeBaseDetailPage({
  t,
  router,
  route,
  userStore,
  locale
});

onMounted(async () => {
  await init();
});
</script>

<style scoped>
.detail-content {
  display: flex;
  flex-direction: column;
  flex: 1;
  min-height: 0;
  gap: var(--spacing-lg);
}

.detail-columns {
  display: flex;
  flex: 1;
  min-height: 0;
  gap: var(--spacing-lg);
  align-items: stretch;
}

.detail-tree-panel {
  flex: 1;
  min-width: 0;
  min-height: 0;
}

.detail-templates-panel {
  min-height: 0;
}

.empty-state {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 200px;
  background-color: var(--bg-white);
  border-radius: var(--border-radius-md);
  box-shadow: var(--shadow-sm);
  border: 1px solid var(--border-color);
}

@media (max-width: 1200px) {
  .detail-columns {
    flex-direction: column;
  }
}
</style>
