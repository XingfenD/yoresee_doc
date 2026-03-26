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
      <el-button class="page-action-btn" type="primary" size="small" @click="openCreateDocumentDialog">
        {{ t('document.createDocument') }}
      </el-button>
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
    :knowledge-base-id="''"
    @submit="createDocument"
    @cancel="cancelCreateDocument"
  />
</template>

<script setup>
import { ref, computed, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { useI18n } from 'vue-i18n';
import { useUserStore } from '@/store/user';
import PageLayout from '@/components/PageLayout.vue';
import DocumentListSection from '@/components/DocumentListSection.vue';
import DocumentCreateDialog from '@/components/DocumentCreateDialog.vue';
import { getMyDocuments, createDocument as createDocumentApi } from '@/services/api';

const router = useRouter();
const userStore = useUserStore();
const { locale, t } = useI18n();

const systemName = ref(userStore.systemName || 'Yoresee');
const activeMenu = ref('documents');
const isDarkMode = computed(() => userStore.darkMode);
const loading = ref(false);
const showCreateDialog = ref(false);
const creatingLoading = ref(false);

const page = ref(1);
const pageSize = ref(20);
const total = ref(0);

const userInfo = computed(() => userStore.userInfo);
const userAvatar = computed(
  () => userInfo.value?.avatar || 'https://cube.elemecdn.com/0/88/03b0d39583f48206768a7534e55bcpng.png'
);

const currentLanguage = computed({
  get: () => locale.value,
  set: (value) => {
    locale.value = value;
    localStorage.setItem('language', value);
  }
});

const documents = ref([]);

const mapMyDocuments = (docs) => {
  return (docs || []).map((doc) => ({
    id: doc.external_id || doc.id,
    title: doc.title || t('document.title'),
    author: userInfo.value?.username || t('common.unknown'),
    updatedAt: doc.updated_at || doc.updatedAt,
    views: doc.view_count || doc.views || 0,
    status: doc.status === 0 ? 'draft' : 'published'
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

const openCreateDocumentDialog = () => {
  showCreateDialog.value = true;
};

const cancelCreateDocument = () => {
  showCreateDialog.value = false;
};

const createDocument = async (payload) => {
  if (!payload?.title?.trim()) {
    return;
  }
  try {
    creatingLoading.value = true;
    const requestBody = {
      title: payload.title,
      type: payload.type || 'markdown',
      container_type: 'own'
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

const handleMenuSelect = (key) => {
  activeMenu.value = key;
};

const handleLanguageChange = (command) => {
  currentLanguage.value = command;
};

const toggleTheme = () => {
  userStore.toggleDarkMode();
};

const handleLogout = () => {
  userStore.logout();
  router.push('/login');
};

const initLanguage = () => {
  const savedLanguage = localStorage.getItem('language');
  if (savedLanguage) {
    currentLanguage.value = savedLanguage;
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

const viewDocument = (doc) => {
  router.push(`/mydocument/${doc.id}`);
};

onMounted(() => {
  fetchSystemInfo();
  initLanguage();
  fetchMyDocuments();
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
