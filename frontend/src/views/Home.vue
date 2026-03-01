<template>
  <div class="home-container">
    <!-- 顶部导航栏 -->
    <header class="top-nav">
      <div class="nav-left">
        <h1 class="system-title">{{ systemName }}</h1>
      </div>
      <div class="nav-right">
        <!-- 语言切换 -->
        <el-dropdown trigger="click" @command="handleLanguageChange" class="nav-item">
          <span class="nav-link">
            <el-icon :size="18">
              <Flag v-if="currentLanguage === 'en'" />
              <ChatLineRound v-else />
            </el-icon>
          </span>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item command="en" :icon="'Flag'">
                {{ t('language.english') }}
              </el-dropdown-item>
              <el-dropdown-item command="zh" :icon="'ChatLineRound'">
                {{ t('language.chinese') }}
              </el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>

        <!-- 主题切换 -->
        <div class="nav-item theme-switch">
          <span class="nav-link" @click="toggleTheme">
            <el-icon :size="18">
              <Moon v-if="isDarkMode" />
              <Sunny v-else />
            </el-icon>
          </span>
        </div>

        <!-- 用户菜单 -->
        <el-dropdown trigger="click" class="nav-item">
          <span class="user-info">
            <el-avatar size="small" :src="userAvatar"></el-avatar>
            <span class="username">{{ userInfo?.username || '用户' }}</span>
            <el-icon class="el-icon--right">
              <ArrowDown />
            </el-icon>
          </span>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item @click="handleLogout">{{ t('button.logout') }}</el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </div>
    </header>

    <!-- 主内容区 -->
    <div class="main-content">
      <!-- 左侧导航 -->
      <SideNav :active-menu="activeMenu" @menu-select="handleMenuSelect" />

      <!-- 右侧内容 -->
      <div class="content-area">
        <!-- 操作栏 -->
        <div class="action-bar">
          <h2 class="page-title">{{ t('home.welcome') }}</h2>
        </div>

        <!-- 垂直布局：最近文档和最近知识库 -->
        <div class="home-vertical-layout">
          <!-- 最近文档部分 -->
          <RecentDocumentsSection
            :title="t('home.recentDocuments')"
            :items="recentDocuments"
            :empty-text="t('home.noRecentDocuments')"
            :show-view-all="true"
            @view-all="goToDocuments"
            @view-item="viewDocument"
            @edit-item="editDocument"
          />

          <!-- 最近知识库部分 -->
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
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from 'vue';
import { useRouter } from 'vue-router';
import { useUserStore } from '@/store/user';
import { useI18n } from 'vue-i18n';
import SideNav from '@/components/SideNav.vue';
import RecentDocumentsSection from '@/components/RecentDocumentsSection.vue';
import RecentKnowledgeBaseSection from '@/components/RecentKnowledgeBaseSection.vue';
import { ArrowDown, House, Plus, Upload, Search, User, Timer, View, Flag, ChatLineRound, Moon, Sunny, Collection } from '@element-plus/icons-vue';

const router = useRouter();
const userStore = useUserStore();
const { locale, t } = useI18n();

const systemName = ref('Yoresee');
const activeMenu = ref('home');
const searchKeyword = ref('');
const filterStatus = ref('all');
const isDarkMode = ref(false);

const userInfo = computed(() => userStore.userInfo);
const userAvatar = ref('');

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
    isPublic: true,
    documentsCount: 24
  },
  {
    externalId: 'kb2',
    name: '公司规章制度',
    description: '公司各项规章制度和政策',
    creatorName: '李四',
    updatedAt: '2024-01-14T15:30:00Z',
    isPublic: true,
    documentsCount: 15
  },
  {
    externalId: 'kb3',
    name: '技术分享',
    description: '团队技术分享资料',
    creatorName: '王五',
    updatedAt: '2024-01-13T11:20:00Z',
    isPublic: false,
    documentsCount: 8
  }
]);

// 页面跳转方法
const goToDocuments = () => {
  router.push('/');
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
.home-container {
  height: 100vh;
  display: flex;
  flex-direction: column;
  background-color: var(--bg-light);
  transition: all 0.3s ease;
}

/* 顶部导航栏 */
.top-nav {
  height: 60px;
  background-color: var(--bg-white);
  border-bottom: 1px solid var(--border-color);
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 var(--spacing-xl);
  box-shadow: var(--shadow-sm);
  transition: all 0.3s ease;
}

.system-title {
  font-size: 18px;
  font-weight: 600;
  color: var(--primary-color);
  margin: 0;
}

.nav-right {
  display: flex;
  align-items: center;
  gap: var(--spacing-sm);
}

.nav-item {
  display: flex;
  align-items: center;
  margin-left: var(--spacing-sm);
}

.nav-link {
  display: flex;
  align-items: center;
  padding: var(--spacing-xs) var(--spacing-sm);
  border-radius: var(--border-radius-md);
  color: var(--text-medium);
  transition: all 0.3s ease;
  cursor: pointer;

  &:hover {
    background-color: var(--bg-medium);
    color: var(--primary-color);
  }
}

.theme-switch {
  padding: var(--spacing-xs) var(--spacing-sm);
}

.user-info {
  display: flex;
  align-items: center;
  cursor: pointer;
  padding: var(--spacing-sm) var(--spacing-md);
  border-radius: var(--border-radius-md);
  transition: background-color 0.3s;
}

.user-info:hover {
  background-color: var(--bg-medium);
}

.username {
  margin-left: var(--spacing-sm);
  margin-right: 4px;
  color: var(--text-medium);
  font-size: 14px;
}

/* 主内容区 */
.main-content {
  flex: 1;
  display: flex;
  overflow: hidden;
}

/* 左侧导航 */
.side-nav {
  width: 240px;
  background-color: var(--bg-white);
  border-right: 1px solid var(--border-color);
  overflow-y: auto;
  transition: all 0.3s ease;
}

.side-menu {
  border-right: none;
}

.side-menu .el-menu-item {
  height: 48px;
  line-height: 48px;
  margin: 0;
  border-radius: 0;
  color: var(--text-medium);
  transition: all 0.3s ease;
}

.side-menu .el-menu-item:hover {
  background-color: var(--primary-light);
  color: var(--primary-color);
}

.side-menu .el-menu-item.is-active {
  background-color: var(--primary-light);
  color: var(--primary-color);
}

/* 右侧内容 */
.content-area {
  flex: 1;
  padding: var(--spacing-lg);
  overflow-y: auto;
}

/* 垂直布局 */
.home-vertical-layout {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-lg);
  height: auto;
}

/* 操作栏 */
.action-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: var(--spacing-lg);
  padding-bottom: var(--spacing-md);
  border-bottom: 1px solid var(--border-color);
}

.page-title {
  font-size: 20px;
  font-weight: 600;
  color: var(--text-dark);
  margin: 0;
}

.action-buttons {
  display: flex;
  gap: var(--spacing-md);
}

.primary-btn {
  border-radius: var(--border-radius-md);
  background-color: var(--primary-color);
  border-color: var(--primary-color);
}

.primary-btn:hover {
  background-color: var(--primary-dark);
  border-color: var(--primary-dark);
}

.secondary-btn {
  border-radius: var(--border-radius-md);
  background-color: var(--bg-white);
  border-color: var(--border-color);
  color: var(--text-medium);
}

.secondary-btn:hover {
  border-color: var(--primary-color);
  color: var(--primary-color);
}

/* 搜索和筛选 */
.search-filter {
  display: flex;
  gap: var(--spacing-md);
  margin-bottom: var(--spacing-lg);
}

.search-input {
  width: 300px;
  border-radius: var(--border-radius-md);
}

.filter-select {
  width: 120px;
  border-radius: var(--border-radius-md);
}

/* 文档列表 */
.document-list {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
  gap: var(--spacing-md);
}

.document-card {
  border-radius: var(--border-radius-md);
  transition: all 0.3s ease;
  background-color: var(--bg-white);
}

.document-card:hover {
  box-shadow: var(--shadow-lg);
  transform: translateY(-2px);
}

.document-card-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: var(--spacing-md);
}

.document-title {
  font-size: 16px;
  font-weight: 500;
  color: var(--text-dark);
  margin: 0;
  flex: 1;
}

.document-status {
  padding: 2px 8px;
  border-radius: 12px;
  font-size: 12px;
  font-weight: 500;
  margin-left: var(--spacing-md);
}

.status-draft {
  background-color: var(--bg-medium);
  color: var(--text-light);
}

.status-published {
  background-color: var(--primary-light);
  color: var(--primary-color);
}

.document-meta {
  display: flex;
  gap: var(--spacing-md);
  margin-bottom: var(--spacing-md);
  flex-wrap: wrap;
}

.meta-item {
  display: flex;
  align-items: center;
  font-size: 12px;
  color: var(--text-light);
  gap: 4px;
}

.document-actions {
  display: flex;
  gap: var(--spacing-md);
  justify-content: flex-end;
}

/* 空状态 */
.empty-state {
  margin-top: 60px;
  text-align: center;
}

/* 响应式设计 */
@media (max-width: 1200px) {
  .side-nav {
    width: 200px;
  }

  .document-list {
    grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  }
}

@media (max-width: 768px) {
  .top-nav {
    padding: 0 var(--spacing-md);
  }

  .system-title {
    font-size: 16px;
  }

  .side-nav {
    width: 60px;
  }

  .side-menu .el-menu-item span {
    display: none;
  }

  .content-area {
    padding: var(--spacing-md);
  }

  .action-bar {
    flex-direction: column;
    align-items: flex-start;
    gap: var(--spacing-md);
  }

  .search-filter {
    flex-direction: column;
  }

  .search-input,
  .filter-select {
    width: 100%;
  }

  .document-list {
    grid-template-columns: 1fr;
  }
}
</style>