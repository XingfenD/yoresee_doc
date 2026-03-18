<template>
  <PageLayout
    :system-name="systemName"
    :current-language="currentLanguage"
    :is-dark-mode="isDarkMode"
    :user-avatar="userAvatar"
    :username="userInfo?.username || '用户'"
    :active-menu="activeMenu"
    :title="t('home.welcome')"
    @change-language="handleLanguageChange"
    @toggle-theme="toggleTheme"
    @logout="handleLogout"
    @menu-select="handleMenuSelect"
  >
    <div class="home-horizontal-layout">
      <div class="home-column">
        <RecentDocumentsSection
          :title="t('home.recentDocuments')"
          :items="recentDocuments"
          :empty-text="t('home.noRecentDocuments')"
          :show-view-all="true"
          @view-all="goToDocuments"
          @view-item="viewDocument"
          @edit-item="editDocument"
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
import { ref, onMounted, computed } from 'vue';
import { useRouter } from 'vue-router';
import { useUserStore } from '@/store/user';
import { useI18n } from 'vue-i18n';
import PageLayout from '@/components/PageLayout.vue';
import RecentDocumentsSection from '@/components/RecentDocumentsSection.vue';
import KnowledgeBaseListSection from '@/components/KnowledgeBaseListSection.vue';
import { getRecentKnowledgeBases } from '@/services/api';

const router = useRouter();
const userStore = useUserStore();
const { locale, t } = useI18n();

const systemName = ref('Yoresee');
const activeMenu = ref('home');
const isDarkMode = ref(false);

const userInfo = computed(() => userStore.userInfo);
const userAvatar = computed(() => userInfo.value?.avatar || 'https://cube.elemecdn.com/0/88/03b0d39583f48206768a7534e55bcpng.png');

// 计算当前语言
const currentLanguage = computed({
  get: () => locale.value,
  set: (value) => {
    locale.value = value;
    localStorage.setItem('language', value);
  }
});

// 处理语言切换
const handleLanguageChange = (command) => {
  currentLanguage.value = command;
};

// 处理主题切换
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

// 初始化主题
const initTheme = () => {
  const savedDarkMode = localStorage.getItem('darkMode');
  if (savedDarkMode === 'true') {
    isDarkMode.value = true;
    document.documentElement.classList.add('dark-mode');
  }
};

// 初始化语言
const initLanguage = () => {
  const savedLanguage = localStorage.getItem('language');
  if (savedLanguage) {
    currentLanguage.value = savedLanguage;
  }
};

// 最近文档数据
const recentDocuments = ref([
  {
    id: 1,
    title: '产品需求文档',
    author: '张三',
    updatedAt: '2024-01-15T10:30:00Z',
    views: 120,
    status: 'published'
  },
  {
    id: 2,
    title: '技术架构设计',
    author: '李四',
    updatedAt: '2024-01-14T16:45:00Z',
    views: 85,
    status: 'published'
  },
  {
    id: 3,
    title: '会议纪要',
    author: '王五',
    updatedAt: '2024-01-13T09:15:00Z',
    views: 45,
    status: 'draft'
  },
  {
    id: 4,
    title: '用户手册',
    author: '赵六',
    updatedAt: '2024-01-12T14:20:00Z',
    views: 200,
    status: 'published'
  }
]);

// 最近知识库数据
const recentKnowledgeBases = ref([]);
const recentTagMapper = (kb) =>
  kb.is_public ? { type: 'success', label: t('knowledgeBase.public') } : null;

// 页面跳转方法
const goToDocuments = () => {
  router.push('/mydocuments');
};

const goToKnowledgeBases = () => {
  router.push('/knowledge-base');
};

// 文档相关方法
const viewDocument = (doc) => {
  console.log('View document:', doc);
  // TODO: 实现查看文档功能
};

const editDocument = (doc) => {
  console.log('Edit document:', doc);
  // TODO: 实现编辑文档功能
};

const accessKnowledgeBase = (kb) => {
  const kbId = kb?.external_id || kb?.externalId;
  if (!kbId) {
    return;
  }
  router.push(`/knowledge-base/${kbId}`);
};

const handleMenuSelect = (key) => {
  activeMenu.value = key;
};

const handleLogout = () => {
  userStore.logout();
  router.push('/login');
};

const formatDate = (dateString) => {
  const date = new Date(dateString);
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  });
};

const fetchSystemInfo = async () => {
  try {
    const info = await userStore.fetchSystemInfo();
    systemName.value = info.system_name;
  } catch (err) {
    console.error('获取系统信息失败:', err);
  }
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

onMounted(() => {
  fetchSystemInfo();
  fetchRecentKnowledgeBases();
  initTheme();
  initLanguage();
});
</script>

<style scoped>
.home-horizontal-layout {
  display: flex;
  gap: var(--spacing-lg);
  height: calc(100vh - 60px - 120px);
  align-items: stretch;
}

.home-column {
  flex: 1;
  min-width: 0;
  display: flex;
}

.home-column :deep(.section-content) {
  max-height: 100%;
  overflow-y: auto;
}

.home-column :deep(.vertical-section),
.home-column :deep(.recent-documents-section) {
  width: 100%;
}

@media (max-width: 1024px) {
  .home-horizontal-layout {
    flex-direction: column;
    height: auto;
  }
}

</style>
