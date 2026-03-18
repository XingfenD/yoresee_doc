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
    <div class="home-vertical-layout">
      <RecentDocumentsSection
        :title="t('home.recentDocuments')"
        :items="recentDocuments"
        :empty-text="t('home.noRecentDocuments')"
        :show-view-all="true"
        @view-all="goToDocuments"
        @view-item="viewDocument"
        @edit-item="editDocument"
      />

      <RecentKnowledgeBaseSection
        :title="t('home.recentKnowledgeBases')"
        :items="recentKnowledgeBases"
        :empty-text="t('home.noRecentKnowledgeBases')"
        :show-view-all="true"
        @view-all="goToKnowledgeBases"
        @view-item="viewKnowledgeBase"
        @access-item="accessKnowledgeBase"
      />
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
import RecentKnowledgeBaseSection from '@/components/RecentKnowledgeBaseSection.vue';

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
const recentKnowledgeBases = ref([
  {
    externalId: 'kb1',
    name: '项目知识库',
    description: '项目相关的技术文档和规范',
    creatorName: '张三',
    updatedAt: '2024-01-15T09:00:00Z',
    is_public: true,
    documentsCount: 24
  },
  {
    externalId: 'kb2',
    name: '公司规章制度',
    description: '公司各项规章制度和政策',
    creatorName: '李四',
    updatedAt: '2024-01-14T15:30:00Z',
    is_public: true,
    documentsCount: 15
  },
  {
    externalId: 'kb3',
    name: '技术分享',
    description: '团队技术分享资料',
    creatorName: '王五',
    updatedAt: '2024-01-13T11:20:00Z',
    is_public: false,
    documentsCount: 8
  }
]);

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

// 知识库相关方法
const viewKnowledgeBase = (kb) => {
  console.log('View knowledge base:', kb);
  // TODO: 实现查看知识库功能
};

const accessKnowledgeBase = (kb) => {
  console.log('Access knowledge base:', kb);
  // TODO: 实现访问知识库功能
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

onMounted(() => {
  fetchSystemInfo();
  initTheme();
  initLanguage();
});
</script>

<style scoped>
.home-vertical-layout {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-lg);
  height: auto;
}

</style>
