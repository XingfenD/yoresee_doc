<template>
  <PageLayout
    :system-name="systemName"
    :current-language="currentLanguage"
    :is-dark-mode="isDarkMode"
    :user-avatar="userAvatar"
    :username="userInfo?.username || t('common.unknown')"
    :active-menu="activeMenu"
    :title="''"
    content-padding="xl"
    @change-language="handleLanguageChange"
    @toggle-theme="toggleTheme"
    @logout="handleLogout"
    @menu-select="handleMenuSelect"
  >
    <TitleBar :show-back="true" :back-text="t('common.back')" @back="goBackToKnowledgeBase">
      <template #actions>
        <el-button type="primary" @click="openCreateDocumentDialog">
          {{ t("knowledgeBase.createDocument") }}
        </el-button>
      </template>
    </TitleBar>

    <div class="detail-content" v-loading="loading">
      <InfoStatsCard
        v-if="knowledgeBaseData"
        :title="knowledgeBaseName"
        :description="knowledgeBaseDescription"
        :stats="knowledgeBaseStats"
      />
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
  </PageLayout>

  <DocumentCreateDialog v-model="showCreateDialog" :loading="creatingLoading"
    :knowledge-base-id="route.params.id || ''" @submit="createDocument"
    @cancel="cancelCreateDocument" />
</template>

<script setup>
import { onMounted } from "vue";
import { useRouter, useRoute } from "vue-router";
import { useUserStore } from "@/store/user";
import { useI18n } from "vue-i18n";
import PageLayout from "@/components/PageLayout.vue";
import TitleBar from "@/components/TitleBar.vue";
import DocumentTree from "@/components/DocumentTree.vue";
import DocumentCreateDialog from "@/components/DocumentCreateDialog.vue";
import TemplateListSection from "@/components/TemplateListSection.vue";
import InfoStatsCard from "@/components/InfoStatsCard.vue";
import { MoreFilled } from "@element-plus/icons-vue";
import { useKnowledgeBaseDetailPage } from "@/composables/useKnowledgeBaseDetailPage";

const { locale, t } = useI18n();
const router = useRouter();
const route = useRoute();
const userStore = useUserStore();

const {
  systemName,
  activeMenu,
  isDarkMode,
  currentLanguage,
  userInfo,
  userAvatar,
  knowledgeBaseName,
  knowledgeBaseDescription,
  knowledgeBaseData,
  loading,
  kbTemplates,
  kbTemplatesLoading,
  knowledgeBaseStats,
  searchKeyword,
  sortBy,
  sortOptions,
  currentPage,
  pageSize,
  totalDocumentsCount,
  showCreateDialog,
  creatingLoading,
  createDocument,
  cancelCreateDocument,
  openCreateDocumentDialog,
  handleTreeNodeClick,
  openDocument,
  handleNodeAction,
  handleSizeChange,
  handleCurrentChange,
  handleMenuSelect,
  handleLanguageChange,
  toggleTheme,
  handleLogout,
  goBackToKnowledgeBase,
  templateTagMapper,
  templateMetaMapper,
  openTemplate,
  directoryTreeData,
  init
} = useKnowledgeBaseDetailPage({
  t,
  router,
  route,
  userStore,
  locale
});

onMounted(async () => {
  await init();
});
</script>

<style scoped>
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
