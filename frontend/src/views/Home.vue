<template>
  <div class="home-container">
    <!-- 顶部导航栏 -->
    <header class="top-nav">
      <div class="nav-left">
        <h1 class="system-title">{{ systemName }}</h1>
      </div>
      <div class="nav-right">
        <el-dropdown trigger="click">
          <span class="user-info">
            <el-avatar size="small" :src="userAvatar"></el-avatar>
            <span class="username">{{ userInfo?.username || '用户' }}</span>
            <el-icon class="el-icon--right"><ArrowDown /></el-icon>
          </span>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item @click="handleLogout">退出登录</el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </div>
    </header>

    <!-- 主内容区 -->
    <div class="main-content">
      <!-- 左侧导航 -->
      <aside class="side-nav">
        <el-menu
          :default-active="activeMenu"
          class="side-menu"
          @select="handleMenuSelect"
        >
          <el-menu-item index="documents">
            <el-icon><Document /></el-icon>
            <span>文档管理</span>
          </el-menu-item>
          <el-menu-item index="folders">
            <el-icon><Folder /></el-icon>
            <span>文件夹</span>
          </el-menu-item>
          <el-menu-item index="trash">
            <el-icon><Delete /></el-icon>
            <span>回收站</span>
          </el-menu-item>
        </el-menu>
      </aside>

      <!-- 右侧内容 -->
      <div class="content-area">
        <!-- 操作栏 -->
        <div class="action-bar">
          <h2 class="page-title">文档管理</h2>
          <div class="action-buttons">
            <el-button type="primary" class="primary-btn">
              <el-icon><Plus /></el-icon>
              新建文档
            </el-button>
            <el-button class="secondary-btn">
              <el-icon><Upload /></el-icon>
              上传文件
            </el-button>
          </div>
        </div>

        <!-- 搜索和筛选 -->
        <div class="search-filter">
          <el-input
            v-model="searchKeyword"
            placeholder="搜索文档..."
            prefix-icon="Search"
            class="search-input"
          />
          <el-select v-model="filterStatus" placeholder="状态" class="filter-select">
            <el-option label="全部" value="all"></el-option>
            <el-option label="草稿" value="draft"></el-option>
            <el-option label="已发布" value="published"></el-option>
          </el-select>
        </div>

        <!-- 文档列表 -->
        <div class="document-list">
          <el-card
            v-for="doc in documents"
            :key="doc.id"
            class="document-card"
            hover
          >
            <div class="document-card-header">
              <h3 class="document-title">{{ doc.title }}</h3>
              <span class="document-status" :class="`status-${doc.status}`">
                {{ doc.status === 'draft' ? '草稿' : '已发布' }}
              </span>
            </div>
            <div class="document-meta">
              <span class="meta-item">
                <el-icon><User /></el-icon>
                {{ doc.author }}
              </span>
              <span class="meta-item">
                <el-icon><Timer /></el-icon>
                {{ formatDate(doc.updatedAt) }}
              </span>
              <span class="meta-item">
                <el-icon><View /></el-icon>
                {{ doc.views }} 次查看
              </span>
            </div>
            <div class="document-actions">
              <el-button size="small" text>查看</el-button>
              <el-button size="small" text>编辑</el-button>
              <el-button size="small" text>分享</el-button>
            </div>
          </el-card>
        </div>

        <!-- 空状态 -->
        <div v-if="documents.length === 0" class="empty-state">
          <el-empty description="暂无文档" />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from 'vue';
import { useRouter } from 'vue-router';
import { useUserStore } from '../store/user';
import { ArrowDown, Document, Folder, Delete, Plus, Upload, Search, User, Timer, View } from '@element-plus/icons-vue';

const router = useRouter();
const userStore = useUserStore();

const systemName = ref('文档管理系统');
const activeMenu = ref('documents');
const searchKeyword = ref('');
const filterStatus = ref('all');

const userInfo = computed(() => userStore.userInfo);
const userAvatar = ref('');

// 模拟文档数据
const documents = ref([
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
    const name = await userStore.fetchSystemInfo();
    systemName.value = name;
  } catch (err) {
    console.error('获取系统信息失败:', err);
  }
};

onMounted(() => {
  fetchSystemInfo();
});
</script>

<style scoped>
.home-container {
  height: 100vh;
  display: flex;
  flex-direction: column;
  background-color: #F7F8FA;
}

/* 顶部导航栏 */
.top-nav {
  height: 60px;
  background-color: #FFFFFF;
  border-bottom: 1px solid #C9CDD4;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 32px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
}

.system-title {
  font-size: 18px;
  font-weight: 600;
  color: #165DFF;
  margin: 0;
}

.user-info {
  display: flex;
  align-items: center;
  cursor: pointer;
  padding: 8px 12px;
  border-radius: 8px;
  transition: background-color 0.3s;
}

.user-info:hover {
  background-color: #F2F3F5;
}

.username {
  margin-left: 8px;
  margin-right: 4px;
  color: #4E5969;
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
  background-color: #FFFFFF;
  border-right: 1px solid #C9CDD4;
  overflow-y: auto;
}

.side-menu {
  border-right: none;
}

.side-menu .el-menu-item {
  height: 48px;
  line-height: 48px;
  margin: 0;
  border-radius: 0;
  color: #4E5969;
}

.side-menu .el-menu-item:hover {
  background-color: #E8F0FF;
  color: #165DFF;
}

.side-menu .el-menu-item.is-active {
  background-color: #E8F0FF;
  color: #165DFF;
}

/* 右侧内容 */
.content-area {
  flex: 1;
  padding: 24px;
  overflow-y: auto;
}

/* 操作栏 */
.action-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
}

.page-title {
  font-size: 20px;
  font-weight: 600;
  color: #1D2129;
  margin: 0;
}

.action-buttons {
  display: flex;
  gap: 12px;
}

.primary-btn {
  border-radius: 8px;
  background-color: #165DFF;
  border-color: #165DFF;
}

.primary-btn:hover {
  background-color: #4080FF;
  border-color: #4080FF;
}

.secondary-btn {
  border-radius: 8px;
  background-color: #FFFFFF;
  border-color: #C9CDD4;
  color: #4E5969;
}

.secondary-btn:hover {
  border-color: #165DFF;
  color: #165DFF;
}

/* 搜索和筛选 */
.search-filter {
  display: flex;
  gap: 12px;
  margin-bottom: 24px;
}

.search-input {
  width: 300px;
  border-radius: 8px;
}

.filter-select {
  width: 120px;
  border-radius: 8px;
}

/* 文档列表 */
.document-list {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
  gap: 16px;
}

.document-card {
  border-radius: 8px;
  transition: all 0.3s ease;
}

.document-card:hover {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.12);
  transform: translateY(-2px);
}

.document-card-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 12px;
}

.document-title {
  font-size: 16px;
  font-weight: 500;
  color: #1D2129;
  margin: 0;
  flex: 1;
}

.document-status {
  padding: 2px 8px;
  border-radius: 12px;
  font-size: 12px;
  font-weight: 500;
  margin-left: 12px;
}

.status-draft {
  background-color: #F2F3F5;
  color: #86909C;
}

.status-published {
  background-color: #E8F0FF;
  color: #165DFF;
}

.document-meta {
  display: flex;
  gap: 16px;
  margin-bottom: 16px;
  flex-wrap: wrap;
}

.meta-item {
  display: flex;
  align-items: center;
  font-size: 12px;
  color: #86909C;
  gap: 4px;
}

.document-actions {
  display: flex;
  gap: 12px;
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
    padding: 0 16px;
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
    padding: 16px;
  }
  
  .action-bar {
    flex-direction: column;
    align-items: flex-start;
    gap: 12px;
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