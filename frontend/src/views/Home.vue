<template>
  <div class="home-container">
    <el-container>
      <el-header>
        <div class="header-left">
          <h1>企业文档管理系统</h1>
        </div>
        <div class="header-right">
          <el-dropdown>
            <span class="user-info">
              <el-avatar :size="32" :src="userAvatar"></el-avatar>
              <span class="username">{{ userInfo?.username }}</span>
              <el-icon class="el-icon--right"><ArrowDown /></el-icon>
            </span>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item>个人中心</el-dropdown-item>
                <el-dropdown-item @click="handleLogout">退出登录</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </el-header>
      <el-main>
        <div class="welcome-section">
          <h2>欢迎回来，{{ userInfo?.username }}！</h2>
          <p>您现在可以访问和管理企业文档。</p>
        </div>
        <el-card class="stats-card">
          <template #header>
            <div class="card-header">
              <span>文档统计</span>
            </div>
          </template>
          <div class="stats-grid">
            <div class="stat-item">
              <el-icon class="stat-icon"><Document /></el-icon>
              <div class="stat-content">
                <div class="stat-value">0</div>
                <div class="stat-label">我的文档</div>
              </div>
            </div>
            <div class="stat-item">
              <el-icon class="stat-icon"><Folder /></el-icon>
              <div class="stat-content">
                <div class="stat-value">0</div>
                <div class="stat-label">共享文档</div>
              </div>
            </div>
            <div class="stat-item">
              <el-icon class="stat-icon"><Timer /></el-icon>
              <div class="stat-content">
                <div class="stat-value">0</div>
                <div class="stat-label">最近编辑</div>
              </div>
            </div>
            <div class="stat-item">
              <el-icon class="stat-icon"><Star /></el-icon>
              <div class="stat-content">
                <div class="stat-value">0</div>
                <div class="stat-label">收藏文档</div>
              </div>
            </div>
          </div>
        </el-card>
      </el-main>
    </el-container>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { useUserStore } from '@/store/user';
import { Document, Folder, Timer, Star, ArrowDown } from '@element-plus/icons-vue';

const router = useRouter();
const userStore = useUserStore();

const userInfo = computed(() => userStore.userInfo);
const userAvatar = ref('https://a0ai.marscode.cn/api/ide/v1/text_to_image?prompt=professional%20user%20avatar%20icon&image_size=square');

const handleLogout = () => {
  userStore.logout();
  router.push('/login');
};

onMounted(() => {
  if (!userStore.token) {
    router.push('/login');
  }
});
</script>

<style scoped>
.home-container {
  min-height: 100vh;
  background-color: #f5f7fa;
}

.el-header {
  background-color: white;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0 30px;
}

.header-left h1 {
  font-size: 20px;
  font-weight: 600;
  color: #303133;
  margin: 0;
}

.header-right {
  display: flex;
  align-items: center;
}

.user-info {
  display: flex;
  align-items: center;
  cursor: pointer;
  padding: 8px 12px;
  border-radius: 20px;
  transition: all 0.3s;
}

.user-info:hover {
  background-color: #f5f7fa;
}

.username {
  margin-left: 10px;
  font-size: 14px;
  color: #303133;
}

.el-main {
  padding: 30px;
}

.welcome-section {
  margin-bottom: 30px;
}

.welcome-section h2 {
  font-size: 24px;
  font-weight: 600;
  color: #303133;
  margin-bottom: 10px;
}

.welcome-section p {
  font-size: 14px;
  color: #606266;
  margin: 0;
}

.stats-card {
  margin-bottom: 30px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 20px;
  margin-top: 20px;
}

.stat-item {
  display: flex;
  align-items: center;
  padding: 20px;
  background-color: #f5f7fa;
  border-radius: 8px;
  transition: all 0.3s;
}

.stat-item:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.stat-icon {
  font-size: 24px;
  color: #409eff;
  margin-right: 16px;
}

.stat-content {
  flex: 1;
}

.stat-value {
  font-size: 24px;
  font-weight: 600;
  color: #303133;
  margin-bottom: 4px;
}

.stat-label {
  font-size: 14px;
  color: #606266;
}

@media (max-width: 768px) {
  .el-header {
    padding: 0 20px;
  }
  
  .el-main {
    padding: 20px;
  }
  
  .stats-grid {
    grid-template-columns: 1fr;
  }
}
</style>