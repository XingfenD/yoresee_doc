<template>
  <div class="document-editor-container">
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
        <div class="editor-layout">
          <aside class="sidebar">
            <div class="sidebar-header">
              <el-button text @click="goBack">
                <el-icon>
                  <ArrowLeft />
                </el-icon>
                {{ t('common.back') }}
              </el-button>
            </div>
            <div class="sidebar-title">{{ knowledgeBaseName }}</div>
            <div class="tree-toolbar">
              <el-button text @click="toggleExpandAll"
                :title="isAllExpanded ? t('common.collapseAll') : t('common.expandAll')" class="tree-expand-btn">
                <el-icon :size="14">
                  <FolderOpened v-if="isAllExpanded" />
                  <Folder v-else />
                </el-icon>
              </el-button>
            </div>
            <div class="directory-tree" v-loading="treeLoading">
              <el-tree ref="treeRef" :data="directoryTree" :props="treeProps" node-key="id" :default-expand-all="false"
                :expand-on-click-node="false" @node-click="handleTreeNodeClick" class="editor-tree">
                <template #default="{ node, data }">
                  <div class="tree-node-content">
                    <div class="node-icon">
                      <el-icon v-if="data.isFolder">
                        <FolderOpened v-if="node.expanded" />
                        <Folder v-else />
                      </el-icon>
                      <el-icon v-else>
                        <Document />
                      </el-icon>
                    </div>
                    <span class="node-label">
                      {{ node.label }}
                    </span>
                  </div>
                </template>
              </el-tree>
            </div>
          </aside>

          <main class="editor-main">
            <div class="editor-header">
              <div class="doc-title">{{ currentDocTitle }}</div>
              <div class="editor-actions">
                <el-button type="primary" @click="saveDocument">
                  <el-icon>
                    <Check />
                  </el-icon>
                  {{ t('common.save') }}
                </el-button>
              </div>
            </div>
            <div class="editor-content">
              <div class="editor-wrapper">
                <MarkdownEditor v-model="editorContent" :placeholder="t('document.editorPlaceholder')" />
              </div>
            </div>
            <div class="editor-footer">
              <div class="editor-info">
                <span>{{ t('document.lastSaved') }}: {{ lastSavedTime }}</span>
              </div>
            </div>
          </main>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onBeforeUnmount, computed, nextTick } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { useI18n } from 'vue-i18n';
import { ElMessage } from 'element-plus';
import { ArrowLeft, Folder, FolderOpened, Document, Check, Flag, ChatLineRound, Moon, Sunny, ArrowDown } from '@element-plus/icons-vue';
import MarkdownEditor from '@/components/MarkdownEditor.vue';
import SideNav from '@/components/SideNav.vue';
import { useUserStore } from '@/store/user';
import { getKnowledgeBaseDocuments } from '@/services/api';

const { t, locale } = useI18n();
const route = useRoute();
const router = useRouter();
const userStore = useUserStore();

const kbId = ref(route.params.kbId);
const docId = ref(route.params.docId);

const systemName = ref('Yoresee 知识库');
const userAvatar = ref('https://cube.elemecdn.com/0/88/03b0d39583f48206768a7534e55bcpng.png');
const currentLanguage = computed(() => locale.value);
const isDarkMode = ref(false);
const activeMenu = ref('knowledge-base');

const knowledgeBaseName = ref('示例知识库');
const currentDocTitle = ref('示例文档');
const editorContent = ref('# 欢迎使用文档编辑器\n\n这是一个示例文档内容。\n\n## 功能说明\n\n- 左侧显示知识库目录\n- 右侧为 Markdown 编辑器\n- 支持实时预览\n- 支持语法高亮\n\n```javascript\n// 代码块示例\nconst hello = "Hello World";\nconsole.log(hello);\n```\n\n> 引用块示例\n\n**粗体** 和 *斜体* 文本\n\n| 表格 | 示例 | |\n|------|------|---|\n| 列1  | 列2  | 列3 |\n');
const lastSavedTime = ref('--');
const treeLoading = ref(false);

const treeRef = ref(null);
const isAllExpanded = ref(true);

const treeProps = {
  children: 'children',
  label: 'label'
};

const directoryTree = ref([]);

const transformDocumentsToTree = (documents) => {
  const tree = [];

  documents.forEach(doc => {
    const treeNode = {
      id: doc.external_id,
      label: doc.title,
      isFolder: doc.has_children,
      children: []
    };

    if (doc.children && doc.children.length > 0) {
      treeNode.children = transformDocumentsToTree(doc.children);
    }

    tree.push(treeNode);
  });

  return tree;
};

const fetchDocuments = async () => {
  if (kbId.value === 'example' || kbId.value === 'personal') {
    return;
  }

  treeLoading.value = true;
  try {
    const response = await getKnowledgeBaseDocuments(kbId.value);
    knowledgeBaseName.value = response.knowledge_base.name;
    directoryTree.value = transformDocumentsToTree(response.documents);

    await expandToCurrentDoc();
  } catch (error) {
    console.error('获取文档列表失败:', error);
    ElMessage.error(t('knowledgeBase.fetchError'));
  } finally {
    treeLoading.value = false;
  }
};

const findPathToNode = (tree, targetId, path = []) => {
  for (const node of tree) {
    if (node.id === targetId) {
      return [...path, node];
    }
    if (node.children && node.children.length > 0) {
      const result = findPathToNode(node.children, targetId, [...path, node]);
      if (result) {
        return result;
      }
    }
  }
  return null;
};

const expandToCurrentDoc = async () => {
  if (!treeRef.value || !docId.value || directoryTree.value.length === 0) {
    return;
  }

  const path = findPathToNode(directoryTree.value, docId.value);
  if (path && path.length > 0) {
    await nextTick();
    for (let i = 0; i < path.length - 1; i++) {
      const node = treeRef.value.getNode(path[i].id);
      if (node) {
        node.expanded = true;
      }
    }
  }
};

const goBack = () => {
  if (kbId.value === 'personal' || kbId.value === 'example') {
    router.push('/');
  } else {
    router.push(`/knowledge-base/${kbId.value}`);
  }
};

const handleTreeNodeClick = (data) => {
  if (!data.isFolder) {
    router.push(`/knowledge-base/${kbId.value}/document/${data.id}`);
  }
};

const toggleExpandAll = () => {
  isAllExpanded.value = !isAllExpanded.value;
  if (treeRef.value) {
    const nodes = treeRef.value.store?.nodesMap;
    if (nodes) {
      Object.values(nodes).forEach(node => {
        node.expanded = isAllExpanded.value;
      });
    }
  }
};

const saveDocument = () => {
  const now = new Date();
  lastSavedTime.value = now.toLocaleTimeString();
  ElMessage.success(t('message.saveSuccess'));
};

const handleLanguageChange = (command) => {
  locale.value = command;
  localStorage.setItem('language', command);
};

const toggleTheme = () => {
  isDarkMode.value = !isDarkMode.value;
  document.documentElement.classList.toggle('dark-mode', isDarkMode.value);
  localStorage.setItem('theme', isDarkMode.value ? 'dark' : 'light');
};

const handleMenuSelect = (menu) => {
  activeMenu.value = menu;
  if (menu === 'home') {
    router.push('/');
  } else if (menu === 'knowledge-base') {
    router.push('/knowledge-base');
  } else if (menu === 'documents') {
    router.push('/documents');
  } else if (menu === 'folders') {
    router.push('/folders');
  } else if (menu === 'trash') {
    router.push('/trash');
  } else if (menu === 'templates') {
    router.push('/templates');
  } else if (menu === 'settings') {
    router.push('/settings');
  }
};

const handleLogout = () => {
  userStore.logout();
  router.push('/login');
};

const initTheme = () => {
  const savedTheme = localStorage.getItem('theme');
  if (savedTheme === 'dark') {
    isDarkMode.value = true;
    document.documentElement.classList.add('dark-mode');
  }
};

const initLanguage = () => {
  const savedLanguage = localStorage.getItem('language');
  if (savedLanguage) {
    locale.value = savedLanguage;
  }
};

const fetchSystemInfo = async () => {
  try {
    const info = await userStore.fetchSystemInfo();
    systemName.value = info.system_name;
  } catch (err) {
    console.error('获取系统信息失败:', err);
  }
};

onMounted(async () => {
  if (kbId.value === 'example' && docId.value === 'example') {
    knowledgeBaseName.value = '示例知识库';
    currentDocTitle.value = '示例文档';
  } else {
    await fetchDocuments();
  }

  await fetchSystemInfo();
  initTheme();
  initLanguage();
});
</script>

<style scoped>
.document-editor-container {
  height: 100vh;
  display: flex;
  flex-direction: column;
  background-color: var(--bg-light);
}

.dark-mode .document-editor-container {
  background-color: var(--bg-light);
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

/* 主内容区 */
.main-content {
  flex: 1;
  display: flex;
  overflow: hidden;
}

.content-area {
  flex: 1;
  overflow-y: auto;
  padding: var(--spacing-lg);
  background-color: var(--bg-light);
}

.editor-layout {
  display: flex;
  height: 100%;
  background-color: var(--bg-white);
  border-radius: var(--border-radius-md);
  box-shadow: var(--shadow-sm);
  overflow: hidden;
}

.sidebar {
  width: 280px;
  background-color: var(--bg-white);
  border-right: 1px solid var(--border-color);
  display: flex;
  flex-direction: column;
}

.dark-mode .sidebar {
  background-color: var(--bg-white);
  border-color: var(--border-color);
}

.sidebar-header {
  padding: var(--spacing-md);
  border-bottom: 1px solid var(--border-color);
}

.dark-mode .sidebar-header {
  border-color: var(--border-color);
}

.sidebar-title {
  padding: var(--spacing-md);
  font-size: 16px;
  font-weight: 600;
  color: var(--text-dark);
  border-bottom: 1px solid var(--border-color);
}

.dark-mode .sidebar-title {
  color: var(--text-dark);
  border-color: var(--border-color);
}

.tree-toolbar {
  display: flex;
  justify-content: flex-start;
  padding: var(--spacing-xs) var(--spacing-md);
  border-bottom: 1px solid var(--border-color);
}

.dark-mode .tree-toolbar {
  border-color: var(--border-color);
}

.tree-expand-btn {
  padding: 2px 4px;
}

.tree-expand-btn .el-icon {
  color: var(--text-light);
}

.tree-expand-btn:hover .el-icon {
  color: var(--primary-color);
}

.directory-tree {
  flex: 1;
  overflow-y: auto;
  padding: var(--spacing-sm);
}

.editor-tree {
  background: transparent;
}

.tree-node-content {
  display: flex;
  align-items: center;
  gap: var(--spacing-xs);
  padding: var(--spacing-xs) 0;
  border-radius: var(--border-radius-sm);
}

.tree-node-content.is-active {
  background-color: rgba(22, 93, 255, 0.1);
}

.dark-mode .tree-node-content.is-active {
  background-color: rgba(64, 128, 255, 0.2);
}

.node-icon {
  display: flex;
  align-items: center;
  color: var(--text-light);
}

.dark-mode .node-icon {
  color: var(--text-light);
}

.node-label {
  font-size: 14px;
  color: var(--text-medium);
  cursor: pointer;
}

.dark-mode .node-label {
  color: var(--text-medium);
}

.node-label:hover {
  color: var(--primary-color);
}

.node-label.is-active {
  color: var(--primary-color);
  font-weight: 500;
}

.editor-main {
  flex: 1;
  display: flex;
  flex-direction: column;
  min-height: 600px;
  background-color: var(--bg-white);
}

.dark-mode .editor-main {
  background-color: var(--bg-white);
}

.editor-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: var(--spacing-md) var(--spacing-lg);
  border-bottom: 1px solid var(--border-color);
}

.dark-mode .editor-header {
  border-color: var(--border-color);
}

.doc-title {
  font-size: 18px;
  font-weight: 600;
  color: var(--text-dark);
}

.dark-mode .doc-title {
  color: var(--text-dark);
}

.editor-content {
  flex: 1;
  display: flex;
  overflow: hidden;
  min-height: 0;
}

.editor-wrapper {
  flex: 1;
  display: flex;
  flex-direction: column;
  min-height: 500px;
}

.editor-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: var(--spacing-sm) var(--spacing-lg);
  border-top: 1px solid var(--border-color);
  font-size: 12px;
  color: var(--text-light);
}

.dark-mode .editor-footer {
  border-color: var(--border-color);
  color: var(--text-light);
}
</style>