<template>
    <div class="knowledge-base-detail-container">
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
                    <div class="breadcrumb">
                        <el-breadcrumb separator="/">
                            <el-breadcrumb-item>
                                <div class="breadcrumb-home">
                                    <el-icon class="home-icon">
                                        <House />
                                    </el-icon>
                                    <span class="home-text">{{ t('navigation.home') }}</span>
                                </div>
                            </el-breadcrumb-item>
                            <el-breadcrumb-item>
                                <span>{{ t('navigation.knowledgeBase') }}</span>
                            </el-breadcrumb-item>
                            <el-breadcrumb-item>
                                <span>{{ knowledgeBaseName }}</span>
                            </el-breadcrumb-item>
                        </el-breadcrumb>
                    </div>

                    <div class="actions">
                        <el-button type="primary" @click="createDocument">
                            {{ t('knowledgeBase.createDocument') }}
                        </el-button>
                        <el-button @click="refreshTree" :icon="Refresh" />
                    </div>
                </div>

                <!-- 知识库详情内容 -->
                <div class="detail-content">
                    <div class="kb-info-card">
                        <h2 class="kb-title">{{ knowledgeBaseName }}</h2>
                        <p class="kb-description">{{ knowledgeBaseDescription }}</p>

                        <div class="kb-stats">
                            <div class="stat-item">
                                <el-icon>
                                    <Document />
                                </el-icon>
                                <span>{{ t('knowledgeBase.documentsCount') }}: {{ totalDocuments }}</span>
                            </div>
                            <div class="stat-item">
                                <el-icon>
                                    <Clock />
                                </el-icon>
                                <span>{{ t('knowledgeBase.lastUpdated') }}: {{ formatDate(lastUpdated) }}</span>
                            </div>
                            <div class="stat-item">
                                <el-icon>
                                    <User />
                                </el-icon>
                                <span>{{ t('knowledgeBase.owner') }}: {{ ownerName }}</span>
                            </div>
                        </div>
                    </div>

                    <!-- 文档树形结构 -->
                    <div class="document-tree-section">
                        <div class="section-header">
                            <h3 class="section-title">{{ t('knowledgeBase.documentStructure') }}</h3>

                            <div class="tree-controls">
                                <el-input v-model="searchKeyword" :placeholder="t('knowledgeBase.searchDocuments')"
                                    prefix-icon="Search" clearable class="search-input" />

                                <el-select v-model="sortBy" :placeholder="t('knowledgeBase.sortBy')"
                                    class="sort-select">
                                    <el-option v-for="option in sortOptions" :key="option.value" :label="option.label"
                                        :value="option.value" />
                                </el-select>
                            </div>
                        </div>

                        <div class="tree-content">
                            <el-tree :data="documentTreeData" :props="treeProps" node-key="id"
                                :default-expand-all="false" :expand-on-click-node="false" :accordion="true"
                                @node-click="handleTreeNodeClick" class="custom-tree">
                                <template #default="{ node, data }">
                                    <div class="tree-node-content">
                                        <div class="node-icon">
                                            <el-icon v-if="data.isParent">
                                                <FolderOpened v-if="node.expanded" />
                                                <Folder v-else />
                                            </el-icon>
                                            <el-icon v-else>
                                                <Document />
                                            </el-icon>
                                        </div>

                                        <div class="node-info">
                                            <span class="node-label">{{ node.label }}</span>
                                            <el-tag v-if="data.tags && data.tags.length > 0" size="small" type="info"
                                                class="node-tag">
                                                {{ data.tags[0] }}
                                            </el-tag>
                                        </div>

                                        <div class="node-actions">
                                            <el-button size="small" type="primary" text
                                                @click.stop="openDocument(data)">
                                                {{ t('common.open') }}
                                            </el-button>

                                            <el-dropdown trigger="click" @command="handleNodeAction($event, data)">
                                                <el-button size="small" text @click.stop>
                                                    <el-icon>
                                                        <MoreFilled />
                                                    </el-icon>
                                                </el-button>
                                                <template #dropdown>
                                                    <el-dropdown-menu>
                                                        <el-dropdown-item command="rename">
                                                            {{ t('common.rename') }}
                                                        </el-dropdown-item>
                                                        <el-dropdown-item command="share" divided>
                                                            {{ t('document.share') }}
                                                        </el-dropdown-item>
                                                        <el-dropdown-item command="delete" divided>
                                                            {{ t('common.delete') }}
                                                        </el-dropdown-item>
                                                    </el-dropdown-menu>
                                                </template>
                                            </el-dropdown>
                                        </div>
                                    </div>
                                </template>
                            </el-tree>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useUserStore } from '@/store/user'
import { ElMessage } from 'element-plus'
import { useI18n } from 'vue-i18n'
import SideNav from '@/components/SideNav.vue'
import {
    ArrowDown,
    House,
    Flag,
    ChatLineRound,
    Moon,
    Sunny,
    Document,
    Clock,
    User,
    Folder,
    FolderOpened,
    MoreFilled,
    Refresh,
    Search
} from '@element-plus/icons-vue'

// 国际化
const { locale, t } = useI18n()
const router = useRouter()
const route = useRoute()
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

// 知识库信息
const knowledgeBaseName = ref('项目知识库')
const knowledgeBaseDescription = ref('项目相关的技术文档和规范')
const totalDocuments = ref(24)
const lastUpdated = ref('2024-01-15T09:00:00Z')
const ownerName = ref('张三')

// 文档树相关
const searchKeyword = ref('')
const sortBy = ref('name')
const sortOptions = ref([
    { value: 'name', label: t('knowledgeBase.sortByName') },
    { value: 'date', label: t('knowledgeBase.sortByDate') },
    { value: 'type', label: t('knowledgeBase.sortByType') }
])



// 模拟文档树数据 - 体现父文档和子文档的关系

const treeProps = {
    children: 'children',
    label: 'name'
}

const documentTreeData = ref([
    {
        id: 1,
        name: '项目启动文档',
        type: 'document',
        isParent: true, // 表明这是一个父文档
        children: [
            {
                id: 11,
                name: '项目章程',
                type: 'document',
                isParent: false,
                tags: ['重要'],
                size: '2.4 MB',
                modified: '2024-01-15T09:00:00Z'
            },
            {
                id: 12,
                name: '项目范围说明书',
                type: 'document',
                isParent: false,
                tags: ['重要'],
                size: '1.8 MB',
                modified: '2024-01-14T14:30:00Z'
            }
        ],
        tags: ['重要'],
        size: '4.2 MB',
        modified: '2024-01-15T09:00:00Z'
    },
    {
        id: 2,
        name: '需求分析',
        type: 'document',
        isParent: true,
        children: [
            {
                id: 21,
                name: '功能需求',
                type: 'document',
                isParent: false,
                tags: ['需求'],
                size: '1.2 MB',
                modified: '2024-01-13T11:20:00Z'
            },
            {
                id: 22,
                name: '非功能需求',
                type: 'document',
                isParent: false,
                tags: ['需求'],
                size: '0.6 MB',
                modified: '2024-01-12T16:45:00Z'
            }
        ],
        tags: ['需求'],
        size: '1.8 MB',
        modified: '2024-01-13T11:20:00Z'
    },
    {
        id: 3,
        name: '产品路线图',
        type: 'document',
        isParent: false,
        tags: ['计划'],
        size: '0.5 MB',
        modified: '2024-01-11T10:15:00Z'
    },
    {
        id: 4,
        name: '系统架构',
        type: 'document',
        isParent: true,
        children: [
            {
                id: 41,
                name: '系统架构图',
                type: 'document',
                isParent: false,
                tags: ['设计'],
                size: '3.2 MB',
                modified: '2024-01-10T16:45:00Z'
            },
            {
                id: 42,
                name: '数据库设计规范',
                type: 'document',
                isParent: false,
                tags: ['规范'],
                size: '0.3 MB',
                modified: '2024-01-09T10:15:00Z'
            },
            {
                id: 43,
                name: '接口设计规范',
                type: 'document',
                isParent: false,
                tags: ['规范'],
                size: '0.4 MB',
                modified: '2024-01-08T14:20:00Z'
            }
        ],
        tags: ['设计'],
        size: '3.9 MB',
        modified: '2024-01-10T16:45:00Z'
    },
    {
        id: 5,
        name: 'API文档',
        type: 'document',
        isParent: true,
        children: [
            {
                id: 51,
                name: '用户服务API',
                type: 'document',
                isParent: false,
                tags: ['API'],
                size: '1.1 MB',
                modified: '2024-01-07T09:30:00Z'
            },
            {
                id: 52,
                name: '订单服务API',
                type: 'document',
                isParent: false,
                tags: ['API'],
                size: '0.9 MB',
                modified: '2024-01-06T15:20:00Z'
            }
        ],
        tags: ['API'],
        size: '2.0 MB',
        modified: '2024-01-07T09:30:00Z'
    },
    {
        id: 6,
        name: '开发规范',
        type: 'document',
        isParent: true,
        children: [
            {
                id: 61,
                name: '前端开发规范',
                type: 'document',
                isParent: false,
                tags: ['规范'],
                size: '0.2 MB',
                modified: '2024-01-05T14:10:00Z'
            },
            {
                id: 62,
                name: '后端开发规范',
                type: 'document',
                isParent: false,
                tags: ['规范'],
                size: '0.3 MB',
                modified: '2024-01-05T13:05:00Z'
            }
        ],
        tags: ['规范'],
        size: '0.5 MB',
        modified: '2024-01-05T14:10:00Z'
    },
    {
        id: 7,
        name: '测试文档',
        type: 'document',
        isParent: true,
        children: [
            {
                id: 71,
                name: '测试计划',
                type: 'document',
                isParent: false,
                tags: ['测试'],
                size: '0.8 MB',
                modified: '2024-01-04T10:45:00Z'
            },
            {
                id: 72,
                name: '测试用例',
                type: 'document',
                isParent: false,
                tags: ['测试'],
                size: '1.5 MB',
                modified: '2024-01-03T16:30:00Z'
            },
            {
                id: 73,
                name: '测试报告',
                type: 'document',
                isParent: false,
                tags: ['测试'],
                size: '2.1 MB',
                modified: '2024-01-02T12:20:00Z'
            }
        ],
        tags: ['测试'],
        size: '4.4 MB',
        modified: '2024-01-04T10:45:00Z'
    },
    {
        id: 8,
        name: '会议纪要',
        type: 'document',
        isParent: true,
        children: [
            {
                id: 81,
                name: '第1周会议纪要',
                type: 'document',
                isParent: false,
                tags: ['会议'],
                size: '0.2 MB',
                modified: '2024-01-01T11:00:00Z'
            },
            {
                id: 82,
                name: '第2周会议纪要',
                type: 'document',
                isParent: false,
                tags: ['会议'],
                size: '0.2 MB',
                modified: '2023-12-31T11:00:00Z'
            }
        ],
        tags: ['会议'],
        size: '0.4 MB',
        modified: '2024-01-01T11:00:00Z'
    }
])

// 创建文档
const createDocument = () => {
    ElMessage.success(t('knowledgeBase.createDocumentSuccess'))
}

// 刷新树
const refreshTree = () => {
    ElMessage.info(t('knowledgeBase.treeRefreshed'))
}

// 格式化日期
const formatDate = (dateString) => {
    if (!dateString) return t('common.unknown')

    const date = new Date(dateString)
    return date.toLocaleDateString()
}

// 处理树节点点击
const handleTreeNodeClick = (data) => {
    console.log('Node clicked:', data)
    openDocument(data)
}

// 打开文档
const openDocument = (data) => {
    ElMessage.success(`${t('common.open')} ${data.name}`)
    console.log('Opening document:', data)
}

// 处理节点操作
const handleNodeAction = (command, data) => {
    switch (command) {
        case 'rename':
            ElMessage.info(`${t('common.rename')} ${data.name}`)
            break
        case 'share':
            ElMessage.info(`${t('document.share')} ${data.name}`)
            break
        case 'delete':
            ElMessage.confirm(
                t('message.confirmDelete'),
                t('common.warning'),
                {
                    confirmButtonText: t('button.confirm'),
                    cancelButtonText: t('button.cancel'),
                    type: 'warning'
                }
            )
                .then(() => {
                    ElMessage.success(t('message.deleteSuccess'))
                })
            break
    }
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

    // 设置当前知识库ID
    const kbId = route.params.id
    console.log('Current knowledge base ID:', kbId)
})
</script>

<style scoped>
.knowledge-base-detail-container {
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

.breadcrumb {
    display: flex;
    align-items: center;
}

.breadcrumb-home {
    display: flex;
    align-items: center;
    gap: var(--spacing-xs);
}

.home-icon {
    vertical-align: middle;
}

.home-text {
    vertical-align: middle;
}

.actions {
    display: flex;
    gap: var(--spacing-sm);
}

.page-title {
    margin: 0;
    font-size: 20px;
    font-weight: 600;
    color: var(--text-dark);
}

/* 知识库详情内容 */
.detail-content {
    display: flex;
    flex-direction: column;
    gap: var(--spacing-lg);
}

.kb-info-card {
    background-color: var(--bg-white);
    border-radius: var(--border-radius-md);
    box-shadow: var(--shadow-sm);
    padding: var(--spacing-lg);
    border: 1px solid var(--border-color);
}

.kb-title {
    margin: 0 0 var(--spacing-md) 0;
    font-size: 24px;
    font-weight: 600;
    color: var(--text-dark);
}

.kb-description {
    margin: 0 0 var(--spacing-lg) 0;
    font-size: 16px;
    color: var(--text-medium);
    line-height: 1.6;
}

.kb-stats {
    display: flex;
    gap: var(--spacing-lg);
    flex-wrap: wrap;
}

.stat-item {
    display: flex;
    align-items: center;
    gap: var(--spacing-xs);
    color: var(--text-medium);
    font-size: 14px;
}

.stat-item .el-icon {
    color: var(--primary-color);
}

/* 文档树形结构区域 */
.document-tree-section {
    background-color: var(--bg-white);
    border-radius: var(--border-radius-md);
    box-shadow: var(--shadow-sm);
    overflow: hidden;
    flex: 1;
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

.tree-controls {
    display: flex;
    gap: var(--spacing-md);
    align-items: center;
}

.search-input {
    width: 200px;
}

.sort-select {
    width: 150px;
}

.tree-content {
    padding: var(--spacing-md);
    min-height: 400px;
    max-height: 60vh;
    overflow-y: auto;
}

.custom-tree {
    width: 100%;
}

.tree-node-content {
    display: flex;
    align-items: center;
    width: 100%;
    padding: var(--spacing-xs) 0;
}

.node-icon {
    width: 24px;
    margin-right: var(--spacing-sm);
    color: var(--primary-color);
}

.node-info {
    flex: 1;
    display: flex;
    align-items: center;
    gap: var(--spacing-xs);
}

.node-label {
    color: var(--text-medium);
    font-size: 14px;
}

.node-tag {
    height: 22px;
    padding: 0 var(--spacing-xs);
    font-size: 12px;
    line-height: 20px;
}

.node-actions {
    display: flex;
    align-items: center;
    gap: var(--spacing-xs);
    opacity: 0;
    transition: opacity 0.2s ease;
}

.custom-tree :deep(.el-tree-node__content):hover .node-actions {
    opacity: 1;
}

/* 响应式设计 */
@media (max-width: 768px) {
    .top-nav {
        padding: 0 var(--spacing-md);
    }

    .system-title {
        font-size: 16px;
    }

    .content-area {
        padding: var(--spacing-md);
    }

    .action-bar {
        flex-direction: column;
        align-items: stretch;
        gap: var(--spacing-md);
    }

    .tree-controls {
        flex-direction: column;
        align-items: stretch;
    }

    .search-input {
        width: 100%;
    }

    .sort-select {
        width: 100%;
    }

    .kb-stats {
        flex-direction: column;
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

.dark-mode .sort-select :deep(.el-input__wrapper) {
    background-color: var(--select-bg);
    border-color: var(--select-border);
    color: var(--select-text);
}

.dark-mode .sort-select :deep(.el-input__inner) {
    background-color: var(--select-bg);
    border-color: var(--select-border);
    color: var(--select-text);
}

.dark-mode .sort-select :deep(.el-select__dropdown) {
    background-color: var(--select-bg);
    border-color: var(--select-border);
}

.dark-mode .sort-select :deep(.el-popper) {
    background-color: var(--select-bg);
    border-color: var(--select-border);
    color: var(--select-text);
}

.dark-mode .sort-select :deep(.el-select-dropdown__item) {
    background-color: var(--select-option-bg);
    color: var(--select-text);
}

.dark-mode .sort-select :deep(.el-select-dropdown__item:hover) {
    background-color: var(--select-option-hover);
}

/* 文档树的深色模式支持 */
.dark-mode .document-tree-section {
    background-color: var(--bg-medium);
    border-color: var(--border-color);
}

.dark-mode .section-header {
    background-color: var(--bg-medium);
    border-color: var(--border-color);
}

.dark-mode .section-title {
    color: var(--text-dark);
}

.dark-mode .custom-tree :deep(.el-tree-node__content) {
    background-color: var(--bg-medium);
    color: var(--text-medium);
}

.dark-mode .custom-tree :deep(.el-tree-node__content:hover) {
    background-color: var(--bg-white);
    color: var(--text-dark);
}

.dark-mode .node-label {
    color: var(--text-medium);
}

.dark-mode .node-icon {
    color: var(--primary-color);
}

.dark-mode .node-actions {
    color: var(--text-medium);
}

.dark-mode .kb-info-card {
    background-color: var(--bg-medium);
    border-color: var(--border-color);
}

.dark-mode .kb-title {
    color: var(--text-dark);
}

.dark-mode .kb-description {
    color: var(--text-medium);
}

.dark-mode .stat-item {
    color: var(--text-medium);
}

.dark-mode .stat-item .el-icon {
    color: var(--primary-color);
}

/* 夜间模式下的标签样式 */
.dark-mode .node-tag {
    background-color: rgba(64, 128, 255, 0.1);
    /* 更暗的背景色 */
    border-color: rgba(64, 128, 255, 0.2);
    /* 更暗的边框色 */
    color: var(--text-light);
    /* 调整文字颜色为更浅的灰色 */
}

/* 夜间模式下的Element Plus标签样式 */
.dark-mode :deep(.el-tag--info) {
    background-color: rgba(64, 128, 255, 0.1);
    border-color: rgba(64, 128, 255, 0.2);
    color: var(--text-light);
}

.dark-mode :deep(.el-tag--success) {
    background-color: rgba(51, 209, 122, 0.1);
    border-color: rgba(51, 209, 122, 0.2);
    color: var(--text-light);
}

.dark-mode :deep(.el-tag--warning) {
    background-color: rgba(255, 152, 0, 0.1);
    border-color: rgba(255, 152, 0, 0.2);
    color: var(--text-light);
}

.dark-mode :deep(.el-tag--danger) {
    background-color: rgba(255, 82, 82, 0.1);
    border-color: rgba(255, 82, 82, 0.2);
    color: var(--text-light);
}
</style>