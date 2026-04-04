<template>
  <PageLayout
    :system-name="systemName"
    :current-language="currentLanguage"
    :is-dark-mode="isDarkMode"
    :user-avatar="userAvatar"
    :username="userInfo?.username || t('common.unknown')"
    :active-menu="activeMenu"
    :title="t('search.title')"
    layout="list"
    @change-language="handleLanguageChange"
    @toggle-theme="toggleTheme"
    @logout="handleLogout"
    @menu-select="handleMenuSelect"
  >
    <div class="search-page" v-loading="loading">
      <div class="search-toolbar">
        <el-input
          v-model="keyword"
          :placeholder="t('search.placeholder')"
          clearable
          class="search-input"
          @keyup.enter="submitSearch"
        >
          <template #prefix>
            <el-icon><Search /></el-icon>
          </template>
        </el-input>
        <el-button type="primary" @click="submitSearch">
          {{ t('common.search') }}
        </el-button>
      </div>

      <div class="search-result-meta" v-if="searched">
        {{ t('search.resultCount', { count: total }) }}
      </div>

      <DocumentListSection
        :title="t('search.title')"
        :items="documents"
        :empty-text="t('search.empty')"
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
</template>

<script setup>
import { computed, onMounted, ref, watch } from 'vue';
import { useI18n } from 'vue-i18n';
import { useRoute, useRouter } from 'vue-router';
import { Search } from '@element-plus/icons-vue';
import PageLayout from '@/components/layout/PageLayout.vue';
import DocumentListSection from '@/components/list/DocumentListSection.vue';
import { listDocuments } from '@/services/api';
import { useUserStore } from '@/store/user';
import { useWorkspaceShell } from '@/composables/shell/useWorkspaceShell';
import { usePageBoot } from '@/composables/shell/usePageBoot';

const route = useRoute();
const router = useRouter();
const userStore = useUserStore();
const { locale, t } = useI18n();

const loading = ref(false);
const searched = ref(false);
const keyword = ref('');
const page = ref(1);
const pageSize = ref(20);
const total = ref(0);
const documents = ref([]);

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
  defaultActiveMenu: 'search'
});
const { boot } = usePageBoot({ initLanguage, fetchSystemInfo });

const normalizedRouteKeyword = computed(() => `${route.query.q || ''}`.trim());

const mapSearchDocuments = (docs) =>
  (docs || []).map((doc) => ({
    id: doc.external_id || '',
    title: doc.title || t('document.title'),
    author: userInfo.value?.username || t('common.unknown'),
    updatedAt: doc.updated_at || doc.updatedAt,
    views: doc.view_count || doc.views || 0
  }));

const fetchDocuments = async () => {
  const q = normalizedRouteKeyword.value;
  if (!q) {
    documents.value = [];
    total.value = 0;
    searched.value = false;
    return;
  }

  loading.value = true;
  try {
    const response = await listDocuments({
      title_keyword: q,
      page: page.value,
      page_size: pageSize.value,
      options: {
        include_children: false,
        recursive: false
      }
    });
    documents.value = mapSearchDocuments(response.documents);
    total.value = Number(response.total_count || 0);
    searched.value = true;
  } catch (error) {
    console.error('search documents failed', error);
    documents.value = [];
    total.value = 0;
    searched.value = true;
  } finally {
    loading.value = false;
  }
};

const submitSearch = () => {
  const q = `${keyword.value || ''}`.trim();
  page.value = 1;
  router.push({
    path: '/search',
    query: q ? { q } : {}
  });
};

const handleSizeChange = async (size) => {
  pageSize.value = size;
  page.value = 1;
  await fetchDocuments();
};

const handleCurrentChange = async (current) => {
  page.value = current;
  await fetchDocuments();
};

const viewDocument = (doc) => {
  if (!doc?.id) return;
  router.push(`/mydocument/${doc.id}`);
};

watch(
  () => route.query.q,
  async () => {
    keyword.value = normalizedRouteKeyword.value;
    page.value = 1;
    await fetchDocuments();
  }
);

onMounted(async () => {
  await boot();
  keyword.value = normalizedRouteKeyword.value;
  await fetchDocuments();
});
</script>

<style scoped>
.search-page {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-md);
}

.search-toolbar {
  display: flex;
  gap: var(--spacing-sm);
  align-items: center;
}

.search-input {
  max-width: 560px;
}

.search-result-meta {
  font-size: 13px;
  color: var(--text-light);
}

.pagination-container {
  display: flex;
  justify-content: flex-end;
  padding: var(--spacing-md) 0;
}

@media (max-width: 768px) {
  .search-toolbar {
    flex-direction: column;
    align-items: stretch;
  }

  .search-input {
    max-width: none;
    width: 100%;
  }
}
</style>
