<template>
  <PageLayout
    :system-name="systemName"
    :current-language="currentLanguage"
    :is-dark-mode="isDarkMode"
    :user-avatar="userAvatar"
    :username="userInfo?.username || t('common.user')"
    :active-menu="activeMenu"
    :title="t('home.welcome')"
    @change-language="handleLanguageChange"
    @toggle-theme="toggleTheme"
    @logout="handleLogout"
    @menu-select="handleMenuSelect"
  >
    <div class="home-horizontal-layout">
      <div class="home-column">
        <DocumentListSection
          :title="t('home.recentDocuments')"
          :items="recentDocuments"
          :empty-text="t('home.noRecentDocuments')"
          :show-view-all="true"
          :single-action="true"
          :primary-action-label="t('common.open')"
          @view-all="goToDocuments"
          @view-item="viewDocument"
        />
      </div>

      <div class="home-column">
        <KnowledgeBaseListSection
          :title="t('home.recentKnowledgeBases')"
          :items="recentKnowledgeBases"
          :empty-text="t('home.noRecentKnowledgeBases')"
          :tag-mapper="recentTagMapper"
          :fallback-description="t('knowledgeBase.noDescription')"
          :action-label="t('common.open')"
          @open="accessKnowledgeBase"
        />
      </div>
    </div>
  </PageLayout>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { useUserStore } from '@/store/user';
import { useI18n } from 'vue-i18n';
import PageLayout from '@/components/layout/PageLayout.vue';
import DocumentListSection from '@/components/list/DocumentListSection.vue';
import KnowledgeBaseListSection from '@/components/knowledge-base/KnowledgeBaseListSection.vue';
import { getRecentKnowledgeBases, getRecentDocuments } from '@/services/api';
import { useWorkspaceShell } from '@/composables/shell/useWorkspaceShell';
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
  defaultActiveMenu: 'home'
});
const { boot } = usePageBoot({ initLanguage, fetchSystemInfo });

// 最近文档数据
const recentDocuments = ref([]);

// 最近知识库数据
const recentKnowledgeBases = ref([]);
const recentTagMapper = (kb) =>
  kb.is_public ? { type: 'success', label: t('knowledgeBase.public') } : null;

// 页面跳转方法
const goToDocuments = () => {
  router.push('/mydocuments');
};

// 文档相关方法
const viewDocument = (doc) => {
  router.push(`/mydocument/${doc.id}`);
};

const accessKnowledgeBase = (kb) => {
  const kbId = kb?.external_id || kb?.externalId;
  if (!kbId) {
    return;
  }
  router.push(`/knowledge-base/${kbId}`);
};

const fetchRecentKnowledgeBases = async () => {
  try {
    const data = await getRecentKnowledgeBases({ page: 1, page_size: 6 });
    recentKnowledgeBases.value = data.knowledge_bases || [];
  } catch (err) {
    console.error('获取最近知识库失败:', err);
    recentKnowledgeBases.value = [];
  }
};

const mapRecentDocuments = (docs) => {
  return (docs || []).map((doc) => ({
    id: doc.external_id || doc.id,
    title: doc.title || t('document.title'),
    author: doc.creator_name || doc.creatorName || userInfo.value?.username || t('common.unknown'),
    updatedAt: doc.updated_at || doc.updatedAt
  }));
};

const fetchRecentDocuments = async () => {
  try {
    const data = await getRecentDocuments({ page: 1, page_size: 6 });
    recentDocuments.value = mapRecentDocuments(data.documents);
  } catch (err) {
    console.error('获取最近文档失败:', err);
    recentDocuments.value = [];
  }
};

onMounted(() => {
  boot(fetchRecentKnowledgeBases, fetchRecentDocuments);
});
</script>

<style scoped>
.home-horizontal-layout {
  display: flex;
  gap: var(--spacing-lg);
  flex: 1;
  min-height: 0;
  align-items: stretch;
}

.home-column {
  flex: 1;
  min-width: 0;
  display: flex;
}

.home-column :deep(.section-content) {
  flex: 1;
  min-height: 0;
  overflow-y: auto;
}

.home-column :deep(.card-list-section),
.home-column :deep(.vertical-section),
.home-column :deep(.document-list-section) {
  width: 100%;
  height: 100%;
}

@media (max-width: 1024px) {
  .home-horizontal-layout {
    flex-direction: column;
    height: auto;
  }
}

</style>
