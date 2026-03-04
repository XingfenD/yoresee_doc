<template>
  <div class="knowledge-base-container">
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
            <span class="username">{{ userInfo?.username || t('common.unknown') }}</span>
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
          <h2 class="page-title">{{ t('knowledgeBase.title') }}</h2>
        </div>

        <!-- 垂直布局 -->
        <div class="knowledge-base-vertical-layout">
          <!-- 第一部分：最近访问的知识库 -->
          <div class="vertical-section">
            <div class="section-header">
              <h3 class="section-title">{{ t('knowledgeBase.recent') }}</h3>
            </div>
            <div class="section-content">
              <el-empty :description="t('knowledgeBase.noRecent')" v-if="recentKnowledgeBases.length === 0" />
              <el-card v-for="kb in recentKnowledgeBases" :key="kb.externalId" class="knowledge-base-item">
                <template #header>
                  <div class="card-header">
                    <span class="kb-name">{{ kb.name }}</span>
                    <el-tag v-if="kb.isPublic" type="success" size="small">
                      {{ t('knowledgeBase.public') }}
                    </el-tag>
                  </div>
                </template>

                <p class="kb-description">{{ kb.description || t('knowledgeBase.noDescription') }}</p>

                <div class="kb-actions">
                  <el-button size="small" @click="accessKnowledgeBase(kb)">
                    {{ t('knowledgeBase.access') }}
                  </el-button>
                </div>
              </el-card>
            </div>
          </div>

          <!-- 第二部分：我的知识库 -->
          <div class="vertical-section">
            <div class="section-header">
              <h3 class="section-title">{{ t('knowledgeBase.my') }}</h3>
              <el-button type="primary" size="small" @click="createKnowledgeBase">
                {{ t('knowledgeBase.createNew') }}
              </el-button>
            </div>
            <div class="section-content">
              <el-card v-for="kb in myKnowledgeBases" :key="kb.externalId" class="knowledge-base-item">
                <template #header>
                  <div class="card-header">
                    <span class="kb-name">{{ kb.name }}</span>
                    <el-tag v-if="!kb.isPublic" type="info" size="small">
                      {{ t('knowledgeBase.private') }}
                    </el-tag>
                  </div>
                </template>

                <p class="kb-description">{{ kb.description || t('knowledgeBase.noDescription') }}</p>

                <div class="kb-details">
                  <div class="detail-item">
                    <span class="detail-label">{{ t('knowledgeBase.documentsCount') }}:</span>
                    <span class="detail-value">{{ kb.documentsCount || 0 }}</span>
                  </div>
                  <div class="detail-item">
                    <span class="detail-label">{{ t('knowledgeBase.updatedAt') }}:</span>
                    <span class="detail-value">{{ formatDate(kb.updatedAt) }}</span>
                  </div>
                </div>

                <div class="kb-actions">
                  <el-button size="small" @click="viewKnowledgeBase(kb)">
                    {{ t('common.view') }}
                  </el-button>
                  <el-button size="small" type="primary" @click="accessKnowledgeBase(kb)">
                    {{ t('knowledgeBase.access') }}
                  </el-button>
                </div>
              </el-card>

              <div class="load-more" v-if="myHasMore">
                <el-button @click="loadMoreMyKnowledgeBases" :loading="myLoading" plain>
                  {{ myLoading ? t('common.loading') : t('common.loadMore') }}
                </el-button>
              </div>
            </div>
          </div>

          <!-- 第三部分：公开知识库 -->
          <div class="vertical-section">
            <div class="section-header">
              <h3 class="section-title">{{ t('knowledgeBase.public') }}</h3>
            </div>
            <div class="section-content">
              <el-card v-for="kb in publicKnowledgeBases" :key="kb.externalId" class="knowledge-base-item">
                <template #header>
                  <div class="card-header">
                    <span class="kb-name">{{ kb.name }}</span>
                    <el-tag type="success" size="small">
                      {{ t('knowledgeBase.public') }}
                    </el-tag>
                  </div>
                </template>

                <p class="kb-description">{{ kb.description || t('knowledgeBase.noDescription') }}</p>

                <div class="kb-details">
                  <div class="detail-item">
                    <span class="detail-label">{{ t('knowledgeBase.documentsCount') }}:</span>
                    <span class="detail-value">{{ kb.documentsCount || 0 }}</span>
                  </div>
                  <div class="detail-item">
                    <span class="detail-label">{{ t('knowledgeBase.owner') }}:</span>
                    <span class="detail-value">{{ kb.creatorName || t('common.unknown') }}</span>
                  </div>
                </div>

                <div class="kb-actions">
                  <el-button size="small" @click="viewKnowledgeBase(kb)">
                    {{ t('common.view') }}
                  </el-button>
                  <el-button size="small" type="primary" @click="accessKnowledgeBase(kb)">
                    {{ t('knowledgeBase.access') }}
                  </el-button>
                </div>
              </el-card>

              <div class="load-more" v-if="publicHasMore">
                <el-button @click="loadMorePublicKnowledgeBases" :loading="publicLoading" plain>
                  {{ publicLoading ? t('common.loading') : t('common.loadMore') }}
                </el-button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/store/user'
import { ElMessage } from 'element-plus'
import { useI18n } from 'vue-i18n'
import SideNav from '@/components/SideNav.vue'
import * as api from '@/services/api'
import { ArrowDown, House, Collection, Flag, ChatLineRound, Moon, Sunny, Plus } from '@element-plus/icons-vue'

// 国际化
const { locale, t } = useI18n()
const router = useRouter()
const userStore = useUserStore()

// 系统信息
const systemName = ref('Yoresee')

// 导航相关
const activeMenu = ref('knowledge-base')
const isDarkMode = ref(false)
const currentLanguage = computed({
  get: () => locale.value,
  set: (value) => {
    locale.value = value
    localStorage.setItem('language', value)
  }
})

// 用户信息
const userInfo = computed(() => userStore.userInfo)
const userAvatar = ref('')

// 最近访问的知识库（API尚未实现，使用模拟数据）
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
])

// 我的知识库
const myKnowledgeBases = ref([])
const myPage = ref(1)
const myPageSize = ref(10)
const myTotal = ref(0)
const myLoading = ref(false)
const myHasMore = computed(() => myKnowledgeBases.value.length < myTotal.value)

// 公开知识库
const publicKnowledgeBases = ref([])
const publicPage = ref(1)
const publicPageSize = ref(10)
const publicTotal = ref(0)
const publicLoading = ref(false)
const publicHasMore = computed(() => publicKnowledgeBases.value.length < publicTotal.value)

// 获取我的知识库
const fetchMyKnowledgeBases = async (page = 1, pageSize = 10) => {
  if (myLoading.value) return

  myLoading.value = true

  try {
    const params = {
      page: page,
      page_size: pageSize,
      only_mine: true,
    }

    const data = await api.getKnowledgeBases(params)

    if (page === 1) {
      myKnowledgeBases.value = data.knowledge_bases || []
    } else {
      myKnowledgeBases.value.push(...(data.knowledge_bases || []))
    }

    myTotal.value = data.total || 0
  } catch (error) {
    console.error('Failed to fetch my knowledge bases:', error)
    ElMessage.error(t('knowledgeBase.fetchError'))
  } finally {
    myLoading.value = false
  }
}

// 获取公开知识库
const fetchPublicKnowledgeBases = async (page = 1, pageSize = 10) => {
  if (publicLoading.value) return

  publicLoading.value = true

  try {
    const params = {
      page: page,
      page_size: pageSize,
      is_public: true  // 获取公开知识库
    }

    const data = await api.getKnowledgeBases(params)

    if (page === 1) {
      publicKnowledgeBases.value = data.knowledge_bases || []
    } else {
      publicKnowledgeBases.value.push(...(data.knowledge_bases || []))
    }

    publicTotal.value = data.total || 0
  } catch (error) {
    console.error('Failed to fetch public knowledge bases:', error)
    ElMessage.error(t('knowledgeBase.fetchError'))
  } finally {
    publicLoading.value = false
  }
}

// 加载更多我的知识库
const loadMoreMyKnowledgeBases = async () => {
  if (!myHasMore.value || myLoading.value) return

  myPage.value++
  await fetchMyKnowledgeBases(myPage.value, myPageSize.value)
}

// 加载更多公开知识库
const loadMorePublicKnowledgeBases = async () => {
  if (!publicHasMore.value || publicLoading.value) return

  publicPage.value++
  await fetchPublicKnowledgeBases(publicPage.value, publicPageSize.value)
}

// 创建知识库
const createKnowledgeBase = () => {
  // TODO: 实现创建知识库功能
  ElMessage.info(t('knowledgeBase.createComingSoon'))
}

// 格式化日期
const formatDate = (dateString) => {
  if (!dateString) return t('common.unknown')

  const date = new Date(dateString)
  return date.toLocaleDateString()
}

// 查看知识库详情
const viewKnowledgeBase = (kb) => {
  // 跳转到知识库详情页面
  router.push(`/knowledge-base/${kb.externalId}`)
  ElMessage.info(t('knowledgeBase.viewInfo'))
}

// 访问知识库
const accessKnowledgeBase = (kb) => {
  // 跳转到知识库详情页面
  router.push(`/knowledge-base/${kb.externalId}`)
}

// 处理菜单选择
const handleMenuSelect = (key) => {
  activeMenu.value = key
}

// 处理语言切换
const handleLanguageChange = (command) => {
  currentLanguage.value = command
}

// 处理主题切换
const toggleTheme = () => {
  isDarkMode.value = !isDarkMode.value
  if (isDarkMode.value) {
    document.documentElement.classList.add('dark-mode')
    localStorage.setItem('darkMode', 'true')
  } else {
    document.documentElement.classList.remove('dark-mode')
    localStorage.setItem('darkMode', 'false')
  }
}

// 初始化主题
const initTheme = () => {
  const savedDarkMode = localStorage.getItem('darkMode')
  if (savedDarkMode === 'true') {
    isDarkMode.value = true
    document.documentElement.classList.add('dark-mode')
  }
}

// 初始化语言
const initLanguage = () => {
  const savedLanguage = localStorage.getItem('language')
  if (savedLanguage) {
    currentLanguage.value = savedLanguage
  }
}

// 获取系统信息
const fetchSystemInfo = async () => {
  try {
    const info = await userStore.fetchSystemInfo()
    systemName.value = info.system_name
  } catch (err) {
    console.error('获取系统信息失败:', err)
  }
}

// 登出处理
const handleLogout = () => {
  userStore.logout()
  router.push('/login')
}

onMounted(async () => {
  // 获取系统信息
  await fetchSystemInfo()

  // 初始化主题和语言
  initTheme()
  initLanguage()

  // 获取知识库数据
  await Promise.all([
    fetchMyKnowledgeBases(myPage.value, myPageSize.value),
    fetchPublicKnowledgeBases(publicPage.value, publicPageSize.value)
  ])
})
</script>

<style scoped>
.knowledge-base-container {
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
}

.nav-link:hover {
  background-color: var(--bg-medium);
  color: var(--primary-color);
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
  border-right: 3px solid var(--primary-color);
}

/* 右侧内容区域 */
.content-area {
  flex: 1;
  overflow-y: auto;
  padding: var(--spacing-xl);
  background-color: var(--bg-light);
}

.action-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: var(--spacing-lg);
  padding-bottom: var(--spacing-md);
  border-bottom: 1px solid var(--border-color);
}

.page-title {
  margin: 0;
  font-size: 20px;
  font-weight: 600;
  color: var(--text-dark);
}

/* 垂直布局 */
.knowledge-base-vertical-layout {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-lg);
  height: auto;
}

.vertical-section {
  display: flex;
  flex-direction: column;
  background-color: var(--bg-white);
  border-radius: var(--border-radius-md);
  box-shadow: var(--shadow-sm);
  overflow: hidden;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: var(--spacing-md);
  border-bottom: 1px solid var(--border-color);
  background-color: var(--bg-white);
}

.section-title {
  margin: 0;
  font-size: 18px;
  font-weight: 600;
  color: var(--text-primary);
}

.section-content {
  flex: 1;
  overflow-y: auto;
  padding: var(--spacing-md);
}

.knowledge-base-item {
  margin-bottom: var(--spacing-md);
  transition: box-shadow 0.3s ease;
}

.knowledge-base-item:hover {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-wrap: wrap;
  gap: 10px;
}

.kb-name {
  font-weight: bold;
  font-size: 16px;
  color: var(--el-text-color-primary);
  flex: 1;
}

.kb-description {
  color: var(--el-text-color-regular);
  margin: 10px 0;
  line-height: 1.5;
  word-break: break-word;
}

.kb-details {
  border-top: 1px solid var(--el-border-color-light);
  padding-top: 15px;
  margin-top: 10px;
}

.detail-item {
  display: flex;
  justify-content: space-between;
  margin-bottom: 8px;
}

.detail-label {
  color: var(--el-text-color-secondary);
  font-weight: 500;
}

.detail-value {
  color: var(--el-text-color-primary);
  font-weight: 400;
}

.kb-actions {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
  margin-top: 15px;
}

.load-more {
  text-align: center;
  margin-top: var(--spacing-md);
}

/* 响应式设计 */
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

  .section-header {
    flex-direction: column;
    align-items: stretch;
    gap: var(--spacing-md);
  }
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

/* 知识库卡片的深色模式支持 */
.dark-mode .kb-name {
  color: var(--text-dark); /* 使用更亮的文字颜色，确保在深色背景下清晰可见 */
}

.dark-mode .kb-description {
  color: var(--text-medium); /* 确保描述文字在深色背景下也清晰可见 */
}

.dark-mode .detail-label {
  color: var(--text-light); /* 确保详情标签在深色背景下也清晰可见 */
}

.dark-mode .detail-value {
  color: var(--text-medium); /* 确保详情值在深色背景下也清晰可见 */
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

.dark-mode .knowledge-base-item {
  background-color: var(--bg-white);
  border: 1px solid var(--border-color);
}
</style>