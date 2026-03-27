<template>
  <PageLayout
    :system-name="systemName"
    :current-language="currentLanguage"
    :is-dark-mode="isDarkMode"
    :user-avatar="userAvatar"
    :username="userInfo?.username || t('common.unknown')"
    :active-menu="activeMenu"
    :title="t('knowledgeBase.title')"
    layout="list"
    @change-language="handleLanguageChange"
    @toggle-theme="toggleTheme"
    @logout="handleLogout"
    @menu-select="handleMenuSelect"
  >
    <template #actions>
      <el-button class="page-action-btn" type="primary" size="small" @click="createKnowledgeBase">
        {{ t("knowledgeBase.createNew") }}
      </el-button>
    </template>

    <div class="knowledge-base-vertical-layout">
      <el-tabs v-model="activeTab" class="common-tabs">
        <el-tab-pane :label="t('knowledgeBase.my')" name="my">
          <KnowledgeBaseListSection
            :title="t('knowledgeBase.my')"
            :items="myKnowledgeBases"
            :empty-text="t('knowledgeBase.noRecent')"
            :tag-mapper="myTagMapper"
            :fallback-description="t('knowledgeBase.noDescription')"
            :meta-mapper="myMetaMapper"
            :show-load-more="myHasMore"
            :loading="myLoading"
            :load-more-label="t('common.loadMore')"
            :loading-label="t('common.loading')"
            :action-label="t('common.open')"
            @open="viewKnowledgeBase"
            @load-more="loadMoreMyKnowledgeBases"
          />
        </el-tab-pane>

        <el-tab-pane :label="t('knowledgeBase.recent')" name="recent">
          <KnowledgeBaseListSection
            :title="t('knowledgeBase.recent')"
            :items="recentKnowledgeBases"
            :empty-text="t('knowledgeBase.noRecent')"
            :tag-mapper="recentTagMapper"
            :fallback-description="t('knowledgeBase.noDescription')"
            :meta-mapper="null"
            :show-load-more="false"
            :action-label="t('common.open')"
            @open="accessKnowledgeBase"
          />
        </el-tab-pane>

        <el-tab-pane :label="t('knowledgeBase.publicList')" name="public">
          <KnowledgeBaseListSection
            :title="t('knowledgeBase.publicList')"
            :items="publicKnowledgeBases"
            :empty-text="t('knowledgeBase.noRecent')"
            :tag-type="'success'"
            :tag-label="t('knowledgeBase.public')"
            :fallback-description="t('knowledgeBase.noDescription')"
            :meta-mapper="publicMetaMapper"
            :show-load-more="publicHasMore"
            :loading="publicLoading"
            :load-more-label="t('common.loadMore')"
            :loading-label="t('common.loading')"
            :action-label="t('common.open')"
            @open="viewKnowledgeBase"
            @load-more="loadMorePublicKnowledgeBases"
          />
        </el-tab-pane>
      </el-tabs>
    </div>
  </PageLayout>
</template>

<script setup>
import { onMounted } from "vue";
import { useRouter } from "vue-router";
import { useUserStore } from "@/store/user";
import { useI18n } from "vue-i18n";
import PageLayout from "@/components/PageLayout.vue";
import KnowledgeBaseListSection from "@/components/KnowledgeBaseListSection.vue";
import { useWorkspaceShell } from "@/composables/useWorkspaceShell";
import { useKnowledgeBaseListPage } from "@/composables/useKnowledgeBaseListPage";
import { usePageBoot } from "@/composables/usePageBoot";

const { locale, t } = useI18n();
const router = useRouter();
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
  defaultActiveMenu: "knowledge-base"
});
const { boot } = usePageBoot({ initLanguage, fetchSystemInfo });

const {
  activeTab,
  recentKnowledgeBases,
  myKnowledgeBases,
  publicKnowledgeBases,
  myLoading,
  publicLoading,
  myHasMore,
  publicHasMore,
  myTagMapper,
  recentTagMapper,
  myMetaMapper,
  publicMetaMapper,
  createKnowledgeBase,
  loadMoreMyKnowledgeBases,
  loadMorePublicKnowledgeBases,
  viewKnowledgeBase,
  accessKnowledgeBase,
  init
} = useKnowledgeBaseListPage({
  t,
  router
});

onMounted(() => {
  boot(init);
});
</script>

<style scoped>
/* 垂直布局 */
.knowledge-base-vertical-layout {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-lg);
  height: auto;
}

/* 深色模式支持 */
.dark-mode .search-input :deep(.el-input__wrapper) {
  background-color: var(--input-bg);
  border-color: var(--input-border);
  color: var(--input-text);
}

.dark-mode .search-input :deep(.el-input__inner) {
  background-color: var(--input-bg);
  border-color: var(--input-border);
  color: var(--input-text);
}

.dark-mode .filter-select :deep(.el-input__wrapper) {
  background-color: var(--select-bg);
  border-color: var(--select-border);
  color: var(--select-text);
}

.dark-mode .filter-select :deep(.el-input__inner) {
  background-color: var(--select-bg);
  border-color: var(--select-border);
  color: var(--select-text);
}

.dark-mode .filter-select :deep(.el-select__dropdown) {
  background-color: var(--select-bg);
  border-color: var(--select-border);
}

.dark-mode .filter-select :deep(.el-popper) {
  background-color: var(--select-bg);
  border-color: var(--select-border);
  color: var(--select-text);
}

.dark-mode .filter-select :deep(.el-select-dropdown__item) {
  background-color: var(--select-option-bg);
  color: var(--select-text);
}

.dark-mode .filter-select :deep(.el-select-dropdown__item:hover) {
  background-color: var(--select-option-hover);
}

.dark-mode .el-pagination :deep(.el-pager li) {
  background-color: var(--bg-white);
  color: var(--text-primary);
  border-color: var(--border-color);
}

.dark-mode .el-pagination :deep(button) {
  background-color: var(--bg-white);
  color: var(--text-primary);
  border-color: var(--border-color);
}

.dark-mode .el-pagination.is-background :deep(.btn-next),
.dark-mode .el-pagination.is-background :deep(.btn-prev),
.dark-mode .el-pagination.is-background :deep(.el-pager li) {
  background-color: var(--bg-medium);
  color: var(--text-primary);
}

.dark-mode .el-pagination.is-background :deep(.el-pager li:not(.is-disabled):hover) {
  color: var(--primary-color);
}

.dark-mode .el-pagination.is-background :deep(.el-pager li.is-active) {
  background-color: var(--primary-color);
  color: white;
}

.dark-mode .el-button {
  border-color: var(--border-color);
  color: var(--text-primary);
}

.dark-mode .el-button--primary {
  background-color: var(--primary-color);
  border-color: var(--primary-color);
  color: white;
}

.dark-mode .el-button--primary:hover {
  background-color: var(--primary-light);
  border-color: var(--primary-light);
  color: white;
}

.dark-mode .column {
  background-color: var(--bg-medium);
  border: 1px solid var(--border-color);
}

.dark-mode .column-header {
  background-color: var(--bg-medium);
  border-bottom: 1px solid var(--border-color);
}
</style>
