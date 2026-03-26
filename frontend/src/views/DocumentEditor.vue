<template>
  <div class="document-editor-container">
    <TopNav
      :system-name="systemName"
      :current-language="currentLanguage"
      :is-dark-mode="isDarkMode"
      :user-avatar="userAvatar"
      :username="userInfo?.username || '用户'"
      @change-language="handleLanguageChange"
      @toggle-theme="toggleTheme"
      @logout="handleLogout"
    />

    <!-- 主内容区 -->
    <div class="main-content">
      <!-- 左侧导航 -->
      <SideNav :active-menu="activeMenu" scene="home" @menu-select="handleMenuSelect" />

      <!-- 右侧内容 -->
      <div class="content-area">
        <div class="editor-layout">
          <div class="sidebar-container" :class="{ 'collapsed': isSidebarCollapsed }">
            <aside class="sidebar">
              <div class="sidebar-header">
                <el-button text class="back-button" @click="goBack">
                  <el-icon>
                    <ArrowLeft />
                  </el-icon>
                  {{ t('common.back') }}
                </el-button>
              </div>
              <div class="sidebar-title">
                {{ knowledgeBaseName }}
              </div>
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
          </div>

          <main class="editor-main">
            <div class="editor-header">
              <div class="editor-title">
                <div
                  class="doc-title"
                  v-if="!isEditingTitle"
                  @click="startEditTitle"
                  :title="t('knowledgeBase.enterDocumentTitle')"
                >
                  {{ currentDocTitle || t('knowledgeBase.enterDocumentTitle') }}
                </div>
                <el-input
                  v-else
                  ref="titleInputRef"
                  v-model="pendingTitle"
                  class="doc-title-input"
                  maxlength="200"
                  @blur="commitTitle"
                  @keyup.enter="commitTitle"
                  @keyup.esc="cancelEditTitle"
                />
              </div>
              <div class="editor-actions">
                <el-button
                  class="editor-action-button"
                  text
                  @click="toggleSidebar"
                  :title="isSidebarCollapsed ? t('common.expand') : t('common.collapse')"
                >
                  <el-icon>
                    <component :is="isSidebarCollapsed ? Expand : Fold" />
                  </el-icon>
                </el-button>
                <el-button
                  class="editor-action-button"
                  text
                  @click="toggleCommentSidebar"
                  :title="t('document.comments')"
                >
                  <el-icon><ChatLineRound /></el-icon>
                </el-button>
                <el-dropdown trigger="click" @command="handleHeaderCommand">
                  <el-button class="editor-action-button" text>
                    <el-icon><MoreFilled /></el-icon>
                  </el-button>
                  <template #dropdown>
                    <el-dropdown-menu>
                      <el-dropdown-item command="create_template">
                        {{ t('templates.saveAs') }}
                      </el-dropdown-item>
                    </el-dropdown-menu>
                  </template>
                </el-dropdown>
              </div>
            </div>
            <div class="editor-content">
              <div class="editor-wrapper">
                <div v-if="collabEnabled && !collabReady" class="editor-loading">
                  {{ t('document.loading') }}
                </div>
                <MarkdownEditor
                  ref="markdownEditorRef"
                  v-model="editorContent"
                  :placeholder="t('document.editorPlaceholder')"
                  :collab-enabled="collabEnabled"
                  :collab-room="collabRoom"
                  :collab-url="collabUrl"
                  :collab-token="collabToken"
                  :comment-enabled="inlineCommentEnabled"
                  @collab-sync="handleCollabSync"
                  @comment-add="handleInlineCommentAdd"
                  @comment-remove="handleInlineCommentRemove"
                />
              </div>
            </div>

          </main>
          <CommentSidebar
            ref="commentSidebarRef"
            :title="t('document.comments')"
            :collapse-title="t('common.collapse')"
            :collapsed="isCommentCollapsed"
            :doc-id="docId"
            :inline-enabled="inlineCommentEnabled"
            :user-info="userInfo"
            :on-anchor-click="scrollToInlineAnchor"
            :on-anchor-hover="handleAnchorHover"
            :on-anchor-remove="handleAnchorRemove"
            @toggle="toggleCommentSidebar"
          />
        </div>
      </div>
    </div>
    <DocumentCreateDialog v-model="showCreateDialog" :loading="creatingLoading"
      :parent-external-id="pendingParentId" :knowledge-base-id="kbId !== 'personal' ? kbId : ''"
      @submit="createDocument" @cancel="cancelCreateDocument" />
    <TemplateCreateDialog
      v-model="showTemplateDialog"
      :loading="savingTemplate"
      :title="t('templates.createDialogTitle')"
      :show-content="false"
      :show-kb-scope="kbId !== 'personal'"
      :initial-name="templateDialogInit.name"
      :initial-description="templateDialogInit.description"
      :initial-scope="templateDialogInit.scope"
      :initial-tags="templateDialogInit.tags"
      :initial-content="templateDialogInit.content"
      @submit="submitCreateTemplate"
    />
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
import { ArrowLeft, Expand, Fold, MoreFilled, ChatLineRound } from '@element-plus/icons-vue';
import MarkdownEditor from '@/components/MarkdownEditor.vue';
import CommentSidebar from '@/components/CommentSidebar.vue';
import DocumentCreateDialog from '@/components/DocumentCreateDialog.vue';
import DocumentTree from '@/components/DocumentTree.vue';
import SideNav from '@/components/SideNav.vue';
import TopNav from '@/components/TopNav.vue';
import TemplateCreateDialog from '@/components/TemplateCreateDialog.vue';
import { useUserStore } from '@/store/user';
import {
  getKnowledgeBaseDocuments,
  createDocument as createDocumentApi,
  createTemplate as createTemplateApi,
  getMyDocuments,
  updateDocumentMeta,
  recordRecentDocument
} from '@/services/api';

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
const isDarkMode = computed(() => userStore.darkMode);
const resolveActiveMenu = (currentKbId) => {
  if (currentKbId === 'personal') return 'documents';
  if (currentKbId) return 'knowledge-base';
  return 'home';
};

const activeMenu = ref(resolveActiveMenu(kbId.value));
const userInfo = computed(() => userStore.userInfo);

const knowledgeBaseName = ref('示例知识库');
const currentDocTitle = ref('示例文档');
const isEditingTitle = ref(false);
const pendingTitle = ref('');
const titleInputRef = ref(null);
const savingTitle = ref(false);
const editorContent = ref('');
const collabEnabled = computed(() => !!docId.value && docId.value !== 'example');
const collabRoom = computed(() => (docId.value ? `${docId.value}` : ''));
const collabUrl = computed(() => '/ws/doc');
const collabToken = computed(() => localStorage.getItem('token') || '');
const collabReady = ref(false);
const lastSyncedDocId = ref('');

const treeLoading = ref(false);
const sidebarWidth = ref(280);
const isResizingSidebar = ref(false);
const showCreateDialog = ref(false);
const creatingLoading = ref(false);
const pendingParentId = ref(null);
const treeComponentRef = ref(null);

const treeRef = computed(() => treeComponentRef.value?.treeRef);
const isAllExpanded = ref(true);
const savedState = localStorage.getItem('sidebarCollapsed');
const isSidebarCollapsed = ref(savedState ? JSON.parse(savedState) : false);
const isCommentCollapsed = ref(false);
const markdownEditorRef = ref(null);
const commentSidebarRef = ref(null);
const inlineCommentEnabled = computed(() => !!docId.value && docId.value !== 'example');
const savingTemplate = ref(false);
const showTemplateDialog = ref(false);
const templateDialogInit = ref({
  name: '',
  description: '',
  scope: 'own',
  tags: '',
  content: ''
});

const getVditorInstance = () => markdownEditorRef.value?.getVditor?.();

const handleInlineCommentAdd = (payload) => {
  if (isCommentCollapsed.value) {
    isCommentCollapsed.value = false;
  }
  commentSidebarRef.value?.handleInlineAnchorAdd?.(payload);
};

const handleInlineCommentRemove = (ids) => {
  commentSidebarRef.value?.handleInlineAnchorRemove?.(ids);
};

const handleAnchorHover = (id, hovering) => {
  if (hovering) {
    highlightInlineComment(id);
    return;
  }
  unhighlightInlineComment(id);
};

const handleAnchorRemove = (ids) => {
  const markdownEditor = markdownEditorRef.value;
  if (markdownEditor && typeof markdownEditor.removeCommentIds === 'function') {
    markdownEditor.removeCommentIds(ids);
    return;
  }
  const editor = getVditorInstance();
  if (!editor || typeof editor.removeCommentIds !== 'function') {
    return;
  }
  editor.removeCommentIds(Array.isArray(ids) ? ids : [ids]);
};

const highlightInlineComment = (id) => {
  const editor = getVditorInstance();
  if (editor && typeof editor.hlCommentIds === 'function') {
    editor.hlCommentIds([id]);
  }
};

const unhighlightInlineComment = (id) => {
  const editor = getVditorInstance();
  if (editor && typeof editor.unHlCommentIds === 'function') {
    editor.unHlCommentIds([id]);
  }
};

const scrollToInlineAnchor = (id) => {
  const editor = getVditorInstance();
  if (!editor || typeof editor.getCommentIds !== 'function') {
    return;
  }
  const commentEntries = editor.getCommentIds();
  const target = Array.isArray(commentEntries)
    ? commentEntries.find((entry) => entry.id === id)
    : null;
  const container = editor?.vditor?.wysiwyg?.element || editor?.vditor?.ir?.element;
  if (!target || !container) {
    return;
  }
  const top = Math.max(target.top - 24, 0);
  if (typeof container.scrollTo === 'function') {
    container.scrollTo({ top, behavior: 'smooth' });
  } else {
    container.scrollTop = top;
  }
  if (typeof editor.hlCommentIds === 'function') {
    editor.hlCommentIds([id]);
  }
};

// 更新CSS变量以支持宽度调节
const updateSidebarWidth = () => {
  document.documentElement.style.setProperty('--sidebar-width', `${sidebarWidth.value}px`);
};

// 初始化和监听宽度变化
onMounted(() => {
  updateSidebarWidth();
});

watch(sidebarWidth, () => {
  updateSidebarWidth();
});

const directoryTree = ref([]);

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
    const isPersonal = kbId.value === 'personal';
    const requestBody = {
      title: payload.title,
      type: payload.type || 'markdown',
      container_type: isPersonal ? 'own' : 'knowledge_base'
    };
    if (!isPersonal) {
      requestBody.knowledge_base_external_id = kbId.value;
    }
    if (payload?.parent_external_id) {
      requestBody.parent_external_id = payload.parent_external_id;
    } else if (pendingParentId.value) {
      requestBody.parent_external_id = pendingParentId.value;
    }
    if (payload?.template) {
      requestBody.template_id = payload.template;
    }
    const response = await createDocumentApi(requestBody);

    showCreateDialog.value = false;
    pendingParentId.value = null;
    await fetchDocuments();
    if (response?.external_id) {
      if (isPersonal) {
        router.push(`/mydocument/${response.external_id}`);
      } else {
        router.push(`/knowledge-base/${kbId.value}/document/${response.external_id}`);
      }
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

const handleHeaderCommand = (command) => {
  if (command === 'create_template') {
    openCreateTemplateDialog();
  }
};

const openCreateTemplateDialog = () => {
  const defaultScope = kbId.value && kbId.value !== 'personal' ? 'knowledge_base' : 'own';
  templateDialogInit.value = {
    name: currentDocTitle.value || t('templates.untitled'),
    description: '',
    scope: defaultScope,
    tags: '',
    content: editorContent.value || ''
  };
  showTemplateDialog.value = true;
};

const submitCreateTemplate = async (payload) => {
  if (savingTemplate.value) {
    return;
  }
  if (!editorContent.value || !editorContent.value.trim()) {
    ElMessage.error(t('templates.emptyContent'));
    return;
  }

  try {
    savingTemplate.value = true;
    const requestBody = {
      target_container: payload.scope,
      template_content: JSON.stringify({
        name: payload.name,
        description: payload.description,
        content: editorContent.value,
        tags: payload.tags || []
      })
    };
    if (payload.scope === 'knowledge_base' && kbId.value && kbId.value !== 'personal') {
      requestBody.knowledge_base_id = kbId.value;
    }
    await createTemplateApi(requestBody);
    showTemplateDialog.value = false;
    ElMessage.success(t('templates.saveSuccess'));
  } catch (error) {
    console.error('创建模板失败:', error);
    ElMessage.error(t('templates.saveFailed'));
  } finally {
    savingTemplate.value = false;
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
  // 重新启用过渡动画
  const sidebarContainer = document.querySelector('.sidebar-container');
  const sidebar = document.querySelector('.sidebar');
  if (sidebarContainer) {
    sidebarContainer.style.transition = 'all 0.3s ease-in-out';
  }
  if (sidebar) {
    sidebar.style.transition = 'transform 0.3s ease-in-out, opacity 0.3s ease-in-out, width 0.3s ease-in-out';
  }
  window.removeEventListener('mousemove', onResizeMove);
  window.removeEventListener('mouseup', stopResize);
};

const startResize = (event) => {
  event.preventDefault();
  isResizingSidebar.value = true;
  document.body.style.cursor = 'col-resize';
  document.body.style.userSelect = 'none';
  // 禁用过渡动画，使调整更跟手
  const sidebarContainer = document.querySelector('.sidebar-container');
  const sidebar = document.querySelector('.sidebar');
  if (sidebarContainer) {
    sidebarContainer.style.transition = 'none';
  }
  if (sidebar) {
    sidebar.style.transition = 'none';
  }
  window.addEventListener('mousemove', onResizeMove);
  window.addEventListener('mouseup', stopResize);
};


const transformDocumentsToTree = (documents, parentId = null) => {
  const tree = [];

  documents.forEach(doc => {
    const treeNode = {
      id: doc.external_id,
      label: doc.title,
      isFolder: !!doc.has_children,
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
  if (kbId.value === 'example') {
    return;
  }

  treeLoading.value = true;
  try {
    if (kbId.value === 'personal') {
      const response = await getMyDocuments({ page: 1, page_size: 1000, directory_only: true });
      knowledgeBaseName.value = t('home.myDocuments');
      directoryTree.value = transformDocumentsToTree(response.documents || []);
    } else {
      const response = await getKnowledgeBaseDocuments(kbId.value, { directory_only: true });
      knowledgeBaseName.value = response.knowledge_base.name;
      directoryTree.value = transformDocumentsToTree(response.documents);
    }

    await expandToCurrentDoc();
    updateCurrentDocTitle();
  } catch (error) {
    console.error('获取文档列表失败:', error);
    ElMessage.error(t('knowledgeBase.fetchError'));
  } finally {
    treeLoading.value = false;
  }
};

const findNodeById = (nodes, targetId) => {
  for (const node of nodes) {
    if (String(node.id) === String(targetId)) {
      return node;
    }
    if (node.children && node.children.length > 0) {
      const found = findNodeById(node.children, targetId);
      if (found) {
        return found;
      }
    }
  }
  return null;
};

const updateCurrentDocTitle = () => {
  if (!docId.value || docId.value === 'example') {
    return;
  }
  const node = findNodeById(directoryTree.value, docId.value);
  if (node) {
    currentDocTitle.value = node.label;
  }
};

const updateTreeNodeTitle = (nodes, targetId, title) => {
  for (const node of nodes) {
    if (String(node.id) === String(targetId)) {
      node.label = title;
      return true;
    }
    if (node.children && node.children.length > 0) {
      if (updateTreeNodeTitle(node.children, targetId, title)) {
        return true;
      }
    }
  }
  return false;
};

const startEditTitle = async () => {
  if (!docId.value || docId.value === 'example') {
    return;
  }
  isEditingTitle.value = true;
  pendingTitle.value = currentDocTitle.value || '';
  await nextTick();
  titleInputRef.value?.focus?.();
};

const cancelEditTitle = () => {
  isEditingTitle.value = false;
  pendingTitle.value = '';
};

const commitTitle = async () => {
  if (!isEditingTitle.value) {
    return;
  }
  const nextTitle = pendingTitle.value.trim();
  if (!nextTitle) {
    ElMessage.error(t('knowledgeBase.titleRequired'));
    return;
  }
  if (nextTitle === currentDocTitle.value) {
    cancelEditTitle();
    return;
  }
  if (!docId.value) {
    cancelEditTitle();
    return;
  }
  if (savingTitle.value) {
    return;
  }
  savingTitle.value = true;
  try {
    await updateDocumentMeta(docId.value, { title: nextTitle });
    currentDocTitle.value = nextTitle;
    updateTreeNodeTitle(directoryTree.value, docId.value, nextTitle);
    cancelEditTitle();
  } catch (error) {
    console.error('更新文档标题失败:', error);
    ElMessage.error(t('common.requestFailed'));
  } finally {
    savingTitle.value = false;
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
    router.push('/mydocuments');
  } else {
    router.push(`/knowledge-base/${kbId.value}`);
  }
};

const handleTreeNodeClick = (data) => {
  if (data?.isCreating) {
    return;
  }
  if (!data?.id) {
    return;
  }
  if (kbId.value === 'personal') {
    router.push(`/mydocument/${data.id}`);
  } else {
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

const toggleSidebar = () => {
  isSidebarCollapsed.value = !isSidebarCollapsed.value;
  localStorage.setItem('sidebarCollapsed', JSON.stringify(isSidebarCollapsed.value));
};

const toggleCommentSidebar = () => {
  isCommentCollapsed.value = !isCommentCollapsed.value;
};


const handleCollabSync = (isSynced) => {
  if (!collabEnabled.value) {
    collabReady.value = true;
    return;
  }
  collabReady.value = isSynced;
  if (isSynced) {
    lastSyncedDocId.value = docId.value || '';
  }
};


const handleLanguageChange = (command) => {
  locale.value = command;
  localStorage.setItem('language', command);
};

const toggleTheme = () => {
  userStore.toggleDarkMode();
};

const handleMenuSelect = (menu) => {
  activeMenu.value = menu;
};

const handleLogout = () => {
  userStore.logout();
  router.push('/login');
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
  initLanguage();

  activeMenu.value = resolveActiveMenu(kbId.value);

  if (kbId.value === 'example' && docId.value === 'example') {
    knowledgeBaseName.value = '示例知识库';
    currentDocTitle.value = '示例文档';
  } else {
    await fetchDocuments();
    if (docId.value) {
      recordRecentDocument(docId.value).catch(() => {});
    }
    if (lastSyncedDocId.value !== docId.value) {
      collabReady.value = !collabEnabled.value;
    }
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
    editorContent.value = '';
    currentDocTitle.value = '';
    cancelEditTitle();
    await commentSidebarRef.value?.reload?.();
    if (docId.value && docId.value !== 'example') {
      recordRecentDocument(docId.value).catch(() => {});
    }
    if (lastSyncedDocId.value !== docId.value) {
      collabReady.value = !collabEnabled.value;
    }
    await expandToCurrentDoc();
    updateCurrentDocTitle();
  }
);

watch(
  () => props.kbId || route.params.kbId,
  async (newKbId) => {
    if (!newKbId) {
      return;
    }
    kbId.value = newKbId;
    cancelEditTitle();
    await fetchDocuments();
    if (lastSyncedDocId.value !== docId.value) {
      collabReady.value = !collabEnabled.value;
    }
    updateCurrentDocTitle();
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
  transition: all 0.3s ease-in-out;
}

.sidebar-container {
  display: flex;
  align-items: stretch;
  position: relative;
  width: calc(var(--sidebar-width) + 6px);
  overflow: hidden;
  flex-shrink: 0;
  transition: all 0.3s ease-in-out;
}

.sidebar-container.collapsed {
  width: 0;
}

.sidebar-container.collapsed .sidebar {
  transform: translateX(-100%);
  opacity: 0;
  pointer-events: none;
}

.sidebar-container.collapsed .sidebar-resizer {
  display: none;
}

.sidebar {
  background-color: var(--bg-white);
  border-right: 1px solid var(--border-color);
  display: flex;
  flex-direction: column;
  width: var(--sidebar-width);
  max-width: 520px;
  transition: transform 0.3s ease-in-out, opacity 0.3s ease-in-out, width 0.3s ease-in-out;
  overflow: hidden;
  flex-shrink: 0;
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
  display: flex;
  justify-content: space-between;
  align-items: center;
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
  transition: all 0.3s ease-in-out;
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

.editor-title {
  display: flex;
  align-items: center;
  min-height: 28px;
}

.editor-actions {
  display: flex;
  align-items: center;
  gap: var(--spacing-sm);
}

.editor-action-button {
  color: var(--text-medium);
}

.editor-action-button:hover {
  color: var(--primary-color);
}

.dark-mode .editor-header {
  border-color: var(--border-color);
}

.doc-title {
  font-size: 18px;
  font-weight: 600;
  color: var(--text-dark);
  cursor: text;
}

.dark-mode .doc-title {
  color: var(--text-dark);
}

.doc-title-input :deep(.el-input__wrapper) {
  box-shadow: none;
  border-radius: 0;
  background-color: transparent;
  padding: 0;
}

.doc-title-input :deep(.el-input__inner) {
  font-size: 18px;
  font-weight: 600;
  color: var(--text-dark);
  padding: 0;
  height: 28px;
  line-height: 28px;
}

.dark-mode .doc-title-input :deep(.el-input__inner) {
  color: var(--text-dark);
}

.editor-content {
  flex: 1;
  display: flex;
  overflow: hidden;
  min-height: 0;
  transition: all 0.3s ease-in-out;
}

.editor-wrapper {
  flex: 1;
  display: flex;
  flex-direction: column;
  min-height: 500px;
  position: relative;
}
.editor-loading {
  position: absolute;
  inset: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(255, 255, 255, 0.85);
  color: var(--text-medium);
  font-size: 14px;
  z-index: 2;
}

.dark-mode .editor-loading {
  background: rgba(255, 255, 255, 0.9);
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
