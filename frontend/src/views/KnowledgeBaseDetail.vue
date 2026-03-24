<template>
  <div class="knowledge-base-detail-container">
    <TopNav
      :system-name="systemName"
      :current-language="currentLanguage"
      :is-dark-mode="isDarkMode"
      :user-avatar="userAvatar"
      :username="userInfo?.username || t('common.unknown')"
      @change-language="handleLanguageChange"
      @toggle-theme="toggleTheme"
      @logout="handleLogout"
    />

    <!-- 主内容区 -->
    <div class="main-content">
      <!-- 左侧导航 -->
      <SideNav :active-menu="activeMenu" @menu-select="handleMenuSelect" />

      <!-- 右侧内容 -->
      <div class="content-area">
        <!-- 操作栏 -->
        <TitleBar :show-back="true" :back-text="t('common.back')" @back="goBackToKnowledgeBase">
          <template #actions>
            <el-button type="primary" @click="openCreateDocumentDialog">
              {{ t("knowledgeBase.createDocument") }}
            </el-button>
          </template>
        </TitleBar>

        <!-- 知识库详情内容 -->
        <div class="detail-content" v-loading="loading">
          <div class="kb-info-card" v-if="knowledgeBaseData">
            <h2 class="kb-title">{{ knowledgeBaseName }}</h2>
            <p class="kb-description">{{ knowledgeBaseDescription }}</p>

            <div class="kb-stats">
              <div class="stat-item">
                <el-icon>
                  <Document />
                </el-icon>
                <span>{{ t("knowledgeBase.documentsCount") }}: {{ totalDocuments }}</span>
              </div>
              <div class="stat-item">
                <el-icon>
                  <Clock />
                </el-icon>
                <span>{{ t("knowledgeBase.lastUpdated") }}:
                  {{ formatDate(lastUpdated) }}</span>
              </div>
              <div class="stat-item">
                <el-icon>
                  <User />
                </el-icon>
                <span>{{ t("knowledgeBase.owner") }}: {{ ownerName }}</span>
              </div>
            </div>
          </div>
          <div v-else-if="!loading" class="empty-state">
            <el-empty :description="t('message.empty')" />
          </div>

          <div class="detail-columns">
            <!-- 文档树形结构 -->
            <div class="document-tree-section">
            <div class="section-header">
              <h3 class="section-title">{{ t("knowledgeBase.documentStructure") }}</h3>

              <div class="tree-controls">
                <el-input v-model="searchKeyword" :placeholder="t('knowledgeBase.searchDocuments')" prefix-icon="Search"
                  clearable class="search-input" />

                <el-select v-model="sortBy" :placeholder="t('knowledgeBase.sortBy')" class="sort-select">
                  <el-option v-for="option in sortOptions" :key="option.value" :label="option.label"
                    :value="option.value" />
                </el-select>
              </div>
            </div>

            <div class="tree-content" v-loading="loading">
              <DocumentTree
                v-if="directoryTreeData.length > 0"
                :nodes="directoryTreeData"
                :loading="loading"
                :show-toolbar="false"
                :show-create="false"
                :show-delete="false"
                :context-menu-enabled="false"
                @node-click="handleTreeNodeClick"
              >
                <template #node-extra="{ data }">
                  <el-tag v-if="data.tags && data.tags.length > 0" size="small" type="info" class="node-tag">
                    {{ data.tags[0] }}
                  </el-tag>
                </template>
                <template #node-actions="{ data }">
                  <el-button size="small" type="primary" text @click.stop="openDocument(data)">
                    {{ t("common.open") }}
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
                          {{ t("common.rename") }}
                        </el-dropdown-item>
                        <el-dropdown-item command="share" divided>
                          {{ t("document.share") }}
                        </el-dropdown-item>
                        <el-dropdown-item command="delete" divided>
                          {{ t("common.delete") }}
                        </el-dropdown-item>
                      </el-dropdown-menu>
                    </template>
                  </el-dropdown>
                </template>
              </DocumentTree>
              <div v-else-if="!loading" class="empty-tree-state">
                <el-empty :description="t('knowledgeBase.noDocuments')" :image-size="64" />
              </div>
            </div>

            <!-- 分页控件 -->
            <div class="pagination-container" v-if="totalDocumentsCount > pageSize">
              <el-pagination v-model:current-page="currentPage" v-model:page-size="pageSize" :page-sizes="[20, 50, 100]"
                :total="totalDocumentsCount" layout="total, sizes, prev, pager, next, jumper"
                @size-change="handleSizeChange" @current-change="handleCurrentChange" />
            </div>
          </div>
            <div class="kb-templates-section">
              <div v-loading="kbTemplatesLoading">
                <TemplateListSection
                  :title="t('knowledgeBase.templates')"
                  :items="kbTemplates"
                  :empty-text="t('templates.noMy')"
                  :fallback-description="t('templates.noDescription')"
                  :tag-mapper="templateTagMapper"
                  :meta-mapper="templateMetaMapper"
                  :action-label="t('common.open')"
                  @open="openTemplate"
                />
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>

  <DocumentCreateDialog v-model="showCreateDialog" :loading="creatingLoading"
    :knowledge-base-id="route.params.id || ''" @submit="createDocument"
    @cancel="cancelCreateDocument" />
</template>

<script setup>
import { ref, onMounted, computed } from "vue";
import { useRouter, useRoute } from "vue-router";
import { useUserStore } from "@/store/user";
import { ElMessage } from "element-plus";
import { useI18n } from "vue-i18n";
import SideNav from "@/components/SideNav.vue";
import TopNav from "@/components/TopNav.vue";
import TitleBar from "@/components/TitleBar.vue";
import DocumentTree from "@/components/DocumentTree.vue";
import DocumentCreateDialog from "@/components/DocumentCreateDialog.vue";
import TemplateListSection from "@/components/TemplateListSection.vue";
import { getKnowledgeBaseDetail, createDocument as createDocumentApi, listTemplates } from "@/services/api.js";
import {
  Document,
  Clock,
  User,
  Folder,
  FolderOpened,
  MoreFilled,
  Search,
} from "@element-plus/icons-vue";

// 国际化
const { locale, t } = useI18n();
const router = useRouter();
const route = useRoute();
const userStore = useUserStore();

// 系统信息
const systemName = ref("Yoresee");

// 导航相关
const activeMenu = ref("knowledge-base");
const isDarkMode = computed(() => userStore.darkMode);
const currentLanguage = computed({
  get: () => locale.value,
  set: (value) => {
    locale.value = value;
    localStorage.setItem("language", value);
  },
});

// 用户信息
const userInfo = computed(() => userStore.userInfo);
const userAvatar = computed(() => userInfo.value?.avatar || 'https://cube.elemecdn.com/0/88/03b0d39583f48206768a7534e55bcpng.png');

// 知识库信息
const knowledgeBaseName = ref("");
const knowledgeBaseDescription = ref("");
const totalDocuments = ref(0);
const lastUpdated = ref("");
const ownerName = ref("");
const knowledgeBaseData = ref(null);
const loading = ref(false);
const kbTemplates = ref([]);
const kbTemplatesLoading = ref(false);

// 文档树相关
const searchKeyword = ref("");
const sortBy = ref("name");
const currentPage = ref(1);
const pageSize = ref(50);
const totalDocumentsCount = ref(0);

// 创建文档对话框相关
const showCreateDialog = ref(false);
const creatingLoading = ref(false);
const sortOptions = ref([
  { value: "name", label: t("knowledgeBase.sortByName") },
  { value: "date", label: t("knowledgeBase.sortByDate") },
  { value: "type", label: t("knowledgeBase.sortByType") },
]);

// 文档树数据 - 从API获取
const directoryTreeData = computed(() => {
  if (!knowledgeBaseData.value || !knowledgeBaseData.value.documents) {
    return [];
  }

  const mapDoc = (doc, parentId = null) => ({
    id: doc.external_id,
    label: doc.title,
    type: doc.type,
    isFolder: doc.hasChildren || (doc.children && doc.children.length > 0),
    isLeaf: !(doc.hasChildren || (doc.children && doc.children.length > 0)),
    tags: doc.tags || [],
    parentId,
    children: doc.children ? doc.children.map((child) => mapDoc(child, doc.external_id)) : []
  });

  return knowledgeBaseData.value.documents.map((doc) => mapDoc(doc));
});


// 加载知识库详情数据
const loadKnowledgeBaseDetail = async () => {
  const knowledgeBaseExternalID = route.params.id;
  if (!knowledgeBaseExternalID) {
    ElMessage.error(t("message.knowledgeBaseNotFound"));
    return;
  }

  loading.value = true;
  try {
    // 调用API获取知识库详情，同时记录最近访问
    const data = await getKnowledgeBaseDetail(knowledgeBaseExternalID, {
      record_recent_log: true,
      page: currentPage.value,
      page_size: pageSize.value
    });

    knowledgeBaseData.value = data;

    // 更新知识库信息
    if (data.knowledge_base) {
      knowledgeBaseName.value = data.knowledge_base.name;
      knowledgeBaseDescription.value = data.knowledge_base.description;
      lastUpdated.value = data.knowledge_base.updated_at;
      totalDocuments.value = data.knowledge_base.documents_count || 0;
      totalDocumentsCount.value = data.total_count || 0;
      ownerName.value = data.knowledge_base.creator_name || "未知用户";
    }
  } catch (error) {
    console.error("加载知识库详情失败:", error);
    ElMessage.error(t("message.loadKnowledgeBaseError"));
  } finally {
    loading.value = false;
  }
};

// 打开创建文档对话框
const openCreateDocumentDialog = () => {
  showCreateDialog.value = true;
};

// 创建文档
const createDocument = async (payload) => {
  if (!payload?.title?.trim()) {
    ElMessage.error(t("knowledgeBase.titleRequired"));
    return;
  }

  try {
    const knowledgeBaseExternalID = route.params.id;
    if (!knowledgeBaseExternalID) {
      ElMessage.error(t("message.knowledgeBaseNotFound"));
      return;
    }

    creatingLoading.value = true;

    // 调用API创建文档
    const response = await createDocumentApi({
      title: payload.title,
      type: payload.type || 'markdown',
      container_type: "knowledge_base",
      knowledge_base_external_id: knowledgeBaseExternalID,
      parent_external_id: payload.parent_external_id || undefined,
      template_id: payload.template || undefined
    });

    // 关闭对话框
    showCreateDialog.value = false;

    // 重新加载知识库详情以显示新创建的文档
    await loadKnowledgeBaseDetail();
  } catch (error) {
    console.error("创建文档失败:", error);
    ElMessage.error(t("knowledgeBase.createDocumentError"));
  } finally {
    creatingLoading.value = false;
  }
};

// 取消创建文档
const cancelCreateDocument = () => {
  showCreateDialog.value = false;
};

// 格式化日期
const formatDate = (dateString) => {
  if (!dateString) return t("common.unknown");

  const date = new Date(dateString);
  return date.toLocaleDateString();
};

// 处理树节点点击
const handleTreeNodeClick = (data) => {
  console.log("Node clicked:", data);
  openDocument(data);
};

// 打开文档
const openDocument = (data) => {
  router.push(`/knowledge-base/${route.params.id}/document/${data.id}`);
};

// 处理节点操作
const handleNodeAction = (command, data) => {
  switch (command) {
    case "rename":
      ElMessage.info(`${t("message.renameDocument")}: ${data.name}`);
      break;
    case "share":
      ElMessage.info(`${t("message.shareDocument")}: ${data.name}`);
      break;
    case "delete":
      ElMessage.confirm(t("message.confirmDelete"), t("common.warning"), {
        confirmButtonText: t("button.confirm"),
        cancelButtonText: t("button.cancel"),
        type: "warning",
      }).then(() => {
        ElMessage.success(t("message.deleteSuccess"));
      });
      break;
  }
};

// 处理菜单选择
const handleMenuSelect = (key) => {
  activeMenu.value = key;
};

// 处理语言切换
const handleLanguageChange = (command) => {
  currentLanguage.value = command;
};

// 处理主题切换
const toggleTheme = () => {
  userStore.toggleDarkMode();
};

// 初始化语言
const initLanguage = () => {
  const savedLanguage = localStorage.getItem("language");
  if (savedLanguage) {
    currentLanguage.value = savedLanguage;
  }
};

// 获取系统信息
const fetchSystemInfo = async () => {
  try {
    const info = await userStore.fetchSystemInfo();
    systemName.value = info.system_name;
  } catch (err) {
    console.error("获取系统信息失败:", err);
  }
};

// 分页大小改变
const handleSizeChange = (val) => {
  pageSize.value = val;
  currentPage.value = 1; // 重置到第一页
  loadKnowledgeBaseDetail();
};

// 当前页改变
const handleCurrentChange = (val) => {
  currentPage.value = val;
  loadKnowledgeBaseDetail();
};

// 登出处理
const handleLogout = () => {
  userStore.logout();
  router.push("/login");
};

const goBackToKnowledgeBase = () => {
  router.push("/knowledge-base");
};

const templateTagMapper = () => ({ type: "info", label: t("templates.private") });
const templateMetaMapper = (tpl) => [
  { label: t("templates.updatedAt"), value: formatDate(tpl.updated_at || tpl.updatedAt) }
];

const openTemplate = (tpl) => {
  if (!tpl?.id) return;
  router.push(`/template/${tpl.id}`);
};

const fetchKnowledgeBaseTemplates = async () => {
  const knowledgeBaseExternalID = route.params.id;
  if (!knowledgeBaseExternalID || kbTemplatesLoading.value) {
    return;
  }
  kbTemplatesLoading.value = true;
  try {
    const data = await listTemplates({
      target_container: "knowledge_base",
      knowledge_base_id: knowledgeBaseExternalID,
      order_by: "updated_at",
      order_desc: true,
      page: 1,
      page_size: 50
    });
    kbTemplates.value = data.templates || [];
  } catch (error) {
    console.error("获取知识库模板失败:", error);
  } finally {
    kbTemplatesLoading.value = false;
  }
};

onMounted(async () => {
  // 获取系统信息
  await fetchSystemInfo();

  // 初始化主题和语言
  initLanguage();

  // 加载知识库详情数据
  await loadKnowledgeBaseDetail();
  await fetchKnowledgeBaseTemplates();
});
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

/* 知识库详情内容 */
.detail-content {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-lg);
}

.detail-columns {
  display: flex;
  gap: var(--spacing-lg);
  align-items: flex-start;
}

.kb-templates-section {
  width: 320px;
  flex-shrink: 0;
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

@media (max-width: 1200px) {
  .detail-columns {
    flex-direction: column;
  }

  .kb-templates-section {
    width: 100%;
  }
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

/* 分页容器样式 */
.pagination-container {
  padding: var(--spacing-md);
  border-top: 1px solid var(--border-color);
  background-color: var(--bg-white);
  display: flex;
  justify-content: center;
  align-items: center;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .content-area {
    padding: var(--spacing-md);
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

/* 空状态样式 */
.empty-state {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 200px;
  background-color: var(--bg-white);
  border-radius: var(--border-radius-md);
  box-shadow: var(--shadow-sm);
  border: 1px solid var(--border-color);
}

.empty-tree-state {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 300px;
  background-color: var(--bg-white);
  border-radius: var(--border-radius-md);
  box-shadow: var(--shadow-sm);
  border: 1px solid var(--border-color);
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

.dark-mode :deep(.el-select-dropdown__item.hover) {
  background-color: var(--select-option-hover);
}

.dark-mode :deep(.el-select-dropdown__item:hover) {
  background-color: var(--select-option-hover);
}

.dark-mode :deep(.el-button) {
  --el-button-bg-color: var(--bg-light);
  --el-button-text-color: var(--text-dark);
  --el-button-hover-bg-color: var(--bg-medium);
  --el-button-hover-text-color: var(--text-dark);
  --el-button-border-color: var(--border-color);
}

.dark-mode :deep(.el-button--primary) {
  --el-button-bg-color: var(--primary-color);
  --el-button-text-color: var(--text-light);
  --el-button-hover-bg-color: var(--primary-color);
  --el-button-hover-text-color: var(--text-light);
}
</style>
