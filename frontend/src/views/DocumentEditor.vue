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
            <el-avatar v-if="userAvatar" size="small" :src="userAvatar"></el-avatar>
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
        <div class="editor-layout">
          <aside class="sidebar" :style="{ width: `${sidebarWidth}px` }">
            <div class="sidebar-header">
              <el-button text @click="goBack">
                <el-icon>
                  <ArrowLeft />
                </el-icon>
                {{ t('common.back') }}
              </el-button>
            </div>
            <div class="sidebar-title">{{ knowledgeBaseName }}</div>
            <DocumentTree
              ref="treeComponentRef"
              :nodes="directoryTree"
              :loading="treeLoading"
              :current-id="docId"
              :expand-all="isAllExpanded"
              :disable-delete="!docId"
              @toggle-expand="toggleExpandAll"
              @node-click="handleTreeNodeClick"
              @create="handleCreateFromTree"
              @delete="handleDeleteDocument"
              @rename="handleRenameFromTree"
            />
          </aside>
          <div class="sidebar-resizer" role="separator" aria-orientation="vertical" @mousedown="startResize"></div>

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
    <DocumentCreateDialog v-model="showCreateDialog" :loading="creatingLoading"
      :parent-external-id="pendingParentId" @submit="createDocument"
      @cancel="cancelCreateDocument" />
  </div>
</template>

<script>
export default {
  inheritAttrs: false
};
</script>

<script setup>
import { ref, onMounted, onBeforeUnmount, computed, nextTick, watch } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { useI18n } from 'vue-i18n';
import { ElMessage, ElMessageBox } from 'element-plus';
import { ArrowLeft, Check, Flag, ChatLineRound, Moon, Sunny, ArrowDown } from '@element-plus/icons-vue';
import MarkdownEditor from '@/components/MarkdownEditor.vue';
import DocumentCreateDialog from '@/components/DocumentCreateDialog.vue';
import DocumentTree from '@/components/DocumentTree.vue';
import SideNav from '@/components/SideNav.vue';
import { useUserStore } from '@/store/user';
import { getKnowledgeBaseDocuments, getDocumentContent, createDocument as createDocumentApi } from '@/services/api';

const props = defineProps({
  kbId: {
    type: String,
    default: ''
  },
  docId: {
    type: String,
    default: ''
  }
});

const { t, locale } = useI18n();
const route = useRoute();
const router = useRouter();
const userStore = useUserStore();

const kbId = ref(props.kbId || route.params.kbId);
const docId = ref(props.docId || route.params.docId);

const systemName = ref(userStore.systemName || 'Yoresee');
const userAvatar = computed(() => userInfo.value?.avatar || 'https://cube.elemecdn.com/0/88/03b0d39583f48206768a7534e55bcpng.png');
const currentLanguage = computed(() => locale.value);
const isDarkMode = ref(document.documentElement.classList.contains('dark-mode'));
const activeMenu = ref('knowledge-base');
const userInfo = computed(() => userStore.userInfo);

const knowledgeBaseName = ref('示例知识库');
const currentDocTitle = ref('示例文档');
const editorContent = ref('# 欢迎使用文档编辑器\n\n这是一个示例文档内容。\n\n## 功能说明\n\n- 左侧显示知识库目录\n- 右侧为 Markdown 编辑器\n- 支持实时预览\n- 支持语法高亮\n\n```javascript\n// 代码块示例\nconst hello = "Hello World";\nconsole.log(hello);\n```\n\n> 引用块示例\n\n**粗体** 和 *斜体* 文本\n\n| 表格 | 示例 | |\n|------|------|---|\n| 列1  | 列2  | 列3 |\n');
const lastSavedTime = ref('--');
const treeLoading = ref(false);
const sidebarWidth = ref(280);
const isResizingSidebar = ref(false);
const showCreateDialog = ref(false);
const creatingLoading = ref(false);
const pendingParentId = ref(null);
const treeComponentRef = ref(null);

const treeRef = computed(() => treeComponentRef.value?.treeRef);
const isAllExpanded = ref(true);

const directoryTree = ref([]);

const isCurrentDoc = (data) => String(data.id) === String(docId.value);

const openCreateDocumentDialog = (parentId = null) => {
  pendingParentId.value = parentId;
  showCreateDialog.value = true;
};

const cancelCreateDocument = () => {
  showCreateDialog.value = false;
};

const createDocument = async (payload) => {
  if (!payload?.title?.trim()) {
    ElMessage.error(t('knowledgeBase.titleRequired'));
    return;
  }

  try {
    creatingLoading.value = true;
    const requestBody = {
      title: payload.title,
      type: payload.type || 'markdown',
      container_type: 'knowledge_base',
      knowledge_base_external_id: kbId.value
    };
    if (payload?.parent_external_id) {
      requestBody.parent_external_id = payload.parent_external_id;
    } else if (pendingParentId.value) {
      requestBody.parent_external_id = pendingParentId.value;
    }
    const response = await createDocumentApi(requestBody);

    showCreateDialog.value = false;
    pendingParentId.value = null;
    await fetchDocuments();
    if (response?.external_id) {
      router.push(`/knowledge-base/${kbId.value}/document/${response.external_id}`);
    }
  } catch (error) {
    console.error('创建文档失败:', error);
    ElMessage.error(t('knowledgeBase.createDocumentError'));
  } finally {
    creatingLoading.value = false;
  }
};

const handleDeleteDocument = async () => {
  if (!docId.value) {
    return;
  }
  try {
    await ElMessageBox.confirm(
      t('document.deleteDocumentConfirm'),
      t('document.deleteDocument'),
      {
        confirmButtonText: t('button.confirm'),
        cancelButtonText: t('button.cancel'),
        type: 'warning'
      }
    );
    ElMessage.warning(t('document.deleteNotSupported'));
  } catch (error) {
    // cancel
  }
};

const onResizeMove = (event) => {
  if (!isResizingSidebar.value) {
    return;
  }
  const layoutRect = document.querySelector('.editor-layout')?.getBoundingClientRect();
  if (!layoutRect) {
    return;
  }
  const minWidth = 220;
  const maxWidth = Math.min(520, layoutRect.width - 320);
  const nextWidth = Math.min(Math.max(event.clientX - layoutRect.left, minWidth), maxWidth);
  sidebarWidth.value = nextWidth;
};

const stopResize = () => {
  if (!isResizingSidebar.value) {
    return;
  }
  isResizingSidebar.value = false;
  document.body.style.cursor = '';
  document.body.style.userSelect = '';
  window.removeEventListener('mousemove', onResizeMove);
  window.removeEventListener('mouseup', stopResize);
};

const startResize = (event) => {
  event.preventDefault();
  isResizingSidebar.value = true;
  document.body.style.cursor = 'col-resize';
  document.body.style.userSelect = 'none';
  window.addEventListener('mousemove', onResizeMove);
  window.addEventListener('mouseup', stopResize);
};


const transformDocumentsToTree = (documents, parentId = null) => {
  const tree = [];

  documents.forEach(doc => {
    const treeNode = {
      id: doc.external_id,
      label: doc.title,
      isFolder: doc.has_children,
      isLeaf: !doc.has_children,
      type: doc.type,
      parentId,
      children: []
    };

    if (doc.children && doc.children.length > 0) {
      treeNode.children = transformDocumentsToTree(doc.children, treeNode.id);
    }

    tree.push(treeNode);
  });

  return tree;
};

const handleCreateFromTree = (target) => {
  openCreateDocumentDialog(target?.id || null);
};

const handleRenameFromTree = () => {
  ElMessage.warning(t('document.renameNotSupported'));
};

const closeContextMenu = () => {
  if (treeComponentRef.value) {
    treeComponentRef.value.closeContextMenu?.();
  }
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

const fetchDocumentContent = async () => {
  if (kbId.value === 'example' || docId.value === 'example') {
    return;
  }

  try {
    const response = await getDocumentContent(docId.value);
    if (response.content !== undefined) {
      editorContent.value = response.content;
    }
    if (response.document) {
      currentDocTitle.value = response.document.title;
    }
  } catch (error) {
    console.error('获取文档内容失败:', error);
    ElMessage.error(t('document.fetchContentError'));
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
  if (data?.isCreating) {
    return;
  }
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
  const savedDarkMode = localStorage.getItem('darkMode');
  if (savedTheme === 'dark' || savedDarkMode === 'true') {
    isDarkMode.value = true;
    document.documentElement.classList.add('dark-mode');
  } else if (savedTheme === 'light' || savedDarkMode === 'false') {
    isDarkMode.value = false;
    document.documentElement.classList.remove('dark-mode');
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
  initTheme();
  initLanguage();

  if (kbId.value === 'example' && docId.value === 'example') {
    knowledgeBaseName.value = '示例知识库';
    currentDocTitle.value = '示例文档';
  } else {
    await fetchDocuments();
    await fetchDocumentContent();
  }

  await fetchSystemInfo();
});

onBeforeUnmount(() => {
  stopResize();
  window.removeEventListener('click', closeContextMenu);
  window.removeEventListener('scroll', closeContextMenu, true);
});

onMounted(() => {
  window.addEventListener('click', closeContextMenu);
  window.addEventListener('scroll', closeContextMenu, true);
});

watch(
  () => props.docId || route.params.docId,
  async (newDocId) => {
    docId.value = newDocId;
    await expandToCurrentDoc();
    await fetchDocumentContent();
  }
);

watch(
  () => props.kbId || route.params.kbId,
  async (newKbId) => {
    if (!newKbId) {
      return;
    }
    kbId.value = newKbId;
    await fetchDocuments();
    await fetchDocumentContent();
  }
);
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
  background-color: var(--bg-white);
  border-right: 1px solid var(--border-color);
  display: flex;
  flex-direction: column;
  min-width: 220px;
  max-width: 520px;
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

.sidebar-resizer {
  width: 6px;
  cursor: col-resize;
  background-color: var(--bg-light);
  border-right: 1px solid var(--border-color);
  transition: background-color 0.2s ease;
}

.sidebar-resizer:hover {
  background-color: var(--bg-medium);
}

.dark-mode .sidebar-resizer {
  background-color: var(--bg-medium);
  border-color: var(--border-color);
}

.dark-mode .sidebar-resizer:hover {
  background-color: var(--bg-white);
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

/* 深色模式对话框样式 */
.dark-mode .el-dialog {
  background-color: var(--bg-white);
  border: 1px solid var(--border-color);
  color: var(--text-dark);
}

.dark-mode .el-dialog__header {
  background-color: var(--bg-white);
  border-bottom: 1px solid var(--border-color);
  color: var(--text-dark);
}

.dark-mode .el-dialog__body {
  background-color: var(--bg-white);
  color: var(--text-dark);
}

.dark-mode .el-dialog__footer {
  background-color: var(--bg-white);
  border-top: 1px solid var(--border-color);
}

.dark-mode .el-form-item__label {
  color: var(--text-dark);
}

.dark-mode :deep(.el-input__wrapper) {
  background-color: var(--input-bg);
  border-color: var(--input-border);
  color: var(--input-text);
}

.dark-mode :deep(.el-input__inner) {
  background-color: var(--input-bg);
  border-color: var(--input-border);
  color: var(--input-text);
}

.dark-mode :deep(.el-select__wrapper) {
  background-color: var(--select-bg);
  border-color: var(--select-border);
  color: var(--select-text);
}

.dark-mode :deep(.el-select__input) {
  background-color: var(--select-bg);
  color: var(--select-text);
}

.dark-mode :deep(.el-select-dropdown__item) {
  background-color: var(--select-option-bg);
  color: var(--select-text);
}
</style>
