<template>
  <PageLayout
    :system-name="systemName"
    :current-language="currentLanguage"
    :is-dark-mode="isDarkMode"
    :user-avatar="userAvatar"
    :username="userInfo?.username || t('common.unknown')"
    :active-menu="activeMenu"
    :title="t('home.myDocuments')"
    layout="list"
    @change-language="handleLanguageChange"
    @toggle-theme="toggleTheme"
    @logout="handleLogout"
    @menu-select="handleMenuSelect"
  >
    <template #actions>
      <DocumentTypeMenu @select="openCreateDocumentDialog">
        <el-button class="page-action-btn" type="primary" size="small">
          {{ t('document.createDocument') }}
        </el-button>
      </DocumentTypeMenu>
    </template>

    <div class="documents-content" v-loading="loading">
      <DocumentListSection
        :title="t('home.myDocuments')"
        :items="documents"
        :empty-text="t('home.noMyDocuments')"
        :show-view-all="false"
        :single-action="true"
        :primary-action-label="t('common.open')"
        @view-item="viewDocument"
      />

      <div class="pagination-container" v-if="total > pageSize">
        <el-pagination
          v-model:current-page="page"
          v-model:page-size="pageSize"
          :page-sizes="[20, 50, 100]"
          :total="total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </div>
  </PageLayout>

  <DocumentCreateDialog
    v-model="showCreateDialog"
    :loading="creatingLoading"
    :initial-document-type="selectedDocumentType"
    :knowledge-base-id="''"
    @submit="createDocument"
    @cancel="cancelCreateDocument"
  />
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { useI18n } from 'vue-i18n';
import { useUserStore } from '@/store/user';
import PageLayout from '@/components/layout/PageLayout.vue';
import DocumentListSection from '@/components/list/DocumentListSection.vue';
import DocumentCreateDialog from '@/components/document/DocumentCreateDialog.vue';
import DocumentTypeMenu from '@/components/document/DocumentTypeMenu.vue';
import { getMyDocuments, createDocument as createDocumentApi } from '@/services/api';
import { useWorkspaceShell } from '@/composables/shell/useWorkspaceShell';
import { usePageBoot } from '@/composables/shell/usePageBoot';
import { DEFAULT_DOCUMENT_TYPE, normalizeDocumentType } from '@/utils/documentType';

const router = useRouter();
const userStore = useUserStore();
const { locale, t } = useI18n();

const loading = ref(false);
const showCreateDialog = ref(false);
const creatingLoading = ref(false);
const selectedDocumentType = ref(DEFAULT_DOCUMENT_TYPE);

const page = ref(1);
const pageSize = ref(20);
const total = ref(0);

const {
  systemName,
  activeMenu,
  isDarkMode,
  userInfo,
  userAvatar,
  currentLanguage,
  initLanguage,
  handleLanguageChange,
  toggleTheme,
  handleLogout,
  handleMenuSelect,
  fetchSystemInfo
} = useWorkspaceShell({
  locale,
  router,
  userStore,
  defaultActiveMenu: 'documents'
});
const { boot } = usePageBoot({ initLanguage, fetchSystemInfo });

const documents = ref([]);

const mapMyDocuments = (docs) => {
  return (docs || []).map((doc) => ({
    id: doc.external_id || doc.id,
    title: doc.title || t('document.title'),
    author: userInfo.value?.username || t('common.unknown'),
    updatedAt: doc.updated_at || doc.updatedAt
  }));
};

const fetchMyDocuments = async () => {
  loading.value = true;
  try {
    const response = await getMyDocuments({
      page: page.value,
      page_size: pageSize.value
    });
    documents.value = mapMyDocuments(response.documents);
    total.value =
      response.total_count ||
      response.total ||
      response.totalCount ||
      response.count ||
      documents.value.length;
  } catch (error) {
    console.error('获取个人文档失败:', error);
  } finally {
    loading.value = false;
  }
};

const openCreateDocumentDialog = (documentType = DEFAULT_DOCUMENT_TYPE) => {
  selectedDocumentType.value = normalizeDocumentType(documentType);
  showCreateDialog.value = true;
};

const cancelCreateDocument = () => {
  showCreateDialog.value = false;
};

const createDocument = async (payload) => {
  const title = payload?.title?.trim() || t('document.untitledDefaultTitle');
  try {
    creatingLoading.value = true;
    const requestBody = {
      title,
      type: normalizeDocumentType(payload?.type || selectedDocumentType.value),
      container_type: 'own',
      is_public: false
    };
    if (payload?.parent_external_id) {
      requestBody.parent_external_id = payload.parent_external_id;
    }
    if (payload?.template) {
      requestBody.template_id = payload.template;
    }
    const response = await createDocumentApi(requestBody);
    showCreateDialog.value = false;
    await fetchMyDocuments();
    if (response?.external_id) {
      router.push(`/mydocument/${response.external_id}`);
    }
  } catch (error) {
    console.error('创建文档失败:', error);
  } finally {
    creatingLoading.value = false;
  }
};

const handleSizeChange = (size) => {
  pageSize.value = size;
  page.value = 1;
  fetchMyDocuments();
};

const handleCurrentChange = (current) => {
  page.value = current;
  fetchMyDocuments();
};

const viewDocument = (doc) => {
  router.push(`/mydocument/${doc.id}`);
};

onMounted(() => {
  boot(fetchMyDocuments);
});
</script>

<style scoped>
.documents-content {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-lg);
}

.pagination-container {
  display: flex;
  justify-content: flex-end;
  padding: var(--spacing-md) 0;
}
</style>
