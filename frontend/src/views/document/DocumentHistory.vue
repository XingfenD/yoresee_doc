<template>
  <PageLayout
    :system-name="systemName"
    :current-language="currentLanguage"
    :is-dark-mode="isDarkMode"
    :user-avatar="userAvatar"
    :username="userInfo?.username || t('common.user')"
    :active-menu="activeMenu"
    :title="''"
    content-padding="xl"
    @change-language="handleLanguageChange"
    @toggle-theme="toggleTheme"
    @logout="handleLogout"
    @menu-select="handleMenuSelect"
  >
    <TitleBar :show-back="true" :compact="true" :back-text="t('common.back')" @back="goBackToDocument">
      <template #title>
        {{ t('history.title') }}
      </template>
      <template #actions>
        <el-button type="primary" @click="openCompareDialog">
          {{ t('history.compare') }}
        </el-button>
      </template>
    </TitleBar>

    <div class="history-content">
      <CommonList
        :rows="versions"
        :columns="columns"
        :row-key="'version'"
        :is-dark="isDarkMode"
        :show-title-bar="true"
        :title="t('history.versionList')"
        :empty-text="t('history.empty')"
        :show-pagination="true"
        :total="total"
        v-model:current-page="page"
        v-model:page-size="pageSize"
        :page-sizes="[20, 50, 100]"
        @current-change="handleCurrentChange"
        @size-change="handleSizeChange"
      >
        <template #cell-actions="{ row }">
          <el-button link type="primary" @click="rollbackVersion(row.version)">
            {{ t('history.rollback') }}
          </el-button>
        </template>
      </CommonList>
    </div>

    <el-dialog
      v-model="compareDialogVisible"
      :title="t('history.compareDialogTitle')"
      width="80%"
      top="6vh"
      destroy-on-close
    >
      <div class="compare-toolbar">
        <el-select v-model="leftVersion" style="width: 220px" @change="loadLeftContent">
          <el-option
            v-for="item in versions"
            :key="`left-${item.version}`"
            :value="item.version"
            :label="`${t('history.versionPrefix')} ${item.version}`"
          />
        </el-select>
        <span class="compare-vs">{{ t('history.vs') }}</span>
        <el-select v-model="rightVersion" style="width: 220px" @change="loadRightContent">
          <el-option
            v-for="item in versions"
            :key="`right-${item.version}`"
            :value="item.version"
            :label="`${t('history.versionPrefix')} ${item.version}`"
          />
        </el-select>
      </div>

      <TextDiffViewer
        :left-text="leftContent ?? t('history.selectHint')"
        :right-text="rightContent ?? t('history.selectHint')"
        :left-title="`${t('history.leftVersion')} ${leftVersion || '-'}`"
        :right-title="`${t('history.rightVersion')} ${rightVersion || '-'}`"
      />
    </el-dialog>
  </PageLayout>
</template>

<script setup>
import { computed, onMounted, ref } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { useI18n } from 'vue-i18n';
import { ElMessage } from 'element-plus';
import PageLayout from '@/components/layout/PageLayout.vue';
import TitleBar from '@/components/layout/TitleBar.vue';
import CommonList from '@/components/list/CommonList.vue';
import TextDiffViewer from '@/components/document/TextDiffViewer.vue';
import { useWorkspaceShell } from '@/composables/shell/useWorkspaceShell';
import { useUserStore } from '@/store/user';
import { listDocumentVersions, getDocumentVersionContent } from '@/services/api';

const route = useRoute();
const router = useRouter();
const userStore = useUserStore();
const { locale, t } = useI18n();

const kbId = computed(() => route.params.kbId || 'personal');
const docId = computed(() => route.params.docId || '');

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
  defaultActiveMenu: kbId.value === 'personal' ? 'documents' : 'knowledge-base'
});

const loading = ref(false);
const versions = ref([]);
const total = ref(0);
const page = ref(1);
const pageSize = ref(20);
const leftVersion = ref(null);
const rightVersion = ref(null);
const leftContent = ref(null);
const rightContent = ref(null);
const contentCache = ref({});
const compareDialogVisible = ref(false);

const columns = computed(() => ([
  { key: 'version', label: t('history.version'), width: '110px' },
  { key: 'title', label: t('history.titleColumn'), minWidth: '220px' },
  { key: 'change_summary', label: t('history.changeSummary'), minWidth: '220px' },
  { key: 'created_at', label: t('history.createdAt'), minWidth: '180px' },
  { key: 'actions', label: t('history.actions'), width: '120px', align: 'center' }
]));

const fetchVersions = async () => {
  if (!docId.value) return;
  loading.value = true;
  try {
    const resp = await listDocumentVersions({
      document_external_id: docId.value,
      page: page.value,
      page_size: pageSize.value
    });
    versions.value = resp.versions || [];
    total.value = Number(resp.total || 0);
  } catch (error) {
    versions.value = [];
    total.value = 0;
    ElMessage.error(t('history.loadFailed'));
  } finally {
    loading.value = false;
  }
};

const loadVersionContent = async (version) => {
  if (!version) return '';
  const cached = contentCache.value[version];
  if (typeof cached === 'string') {
    return cached;
  }
  const resp = await getDocumentVersionContent(docId.value, Number(version));
  const content = resp.version?.content || '';
  contentCache.value = { ...contentCache.value, [version]: content };
  return content;
};

const loadLeftContent = async () => {
  if (!leftVersion.value) {
    leftContent.value = null;
    return;
  }
  try {
    leftContent.value = await loadVersionContent(leftVersion.value);
  } catch (_) {
    leftContent.value = null;
    ElMessage.error(t('history.loadFailed'));
  }
};

const loadRightContent = async () => {
  if (!rightVersion.value) {
    rightContent.value = null;
    return;
  }
  try {
    rightContent.value = await loadVersionContent(rightVersion.value);
  } catch (_) {
    rightContent.value = null;
    ElMessage.error(t('history.loadFailed'));
  }
};

const openCompareDialog = async () => {
  compareDialogVisible.value = true;
  if (versions.value.length === 0) {
    return;
  }
  const latest = versions.value[0]?.version ?? null;
  const previous = versions.value[1]?.version ?? latest;

  leftVersion.value = previous;
  rightVersion.value = latest;
  await Promise.all([loadLeftContent(), loadRightContent()]);
};

const rollbackVersion = (version) => {
  ElMessage.info(`${t('history.rollbackPending')} (${t('history.versionPrefix')} ${version})`);
};

const handleCurrentChange = async (current) => {
  page.value = current;
  await fetchVersions();
};

const handleSizeChange = async (size) => {
  pageSize.value = size;
  page.value = 1;
  await fetchVersions();
};

const goBackToDocument = () => {
  if (kbId.value === 'personal') {
    router.push(`/mydocument/${docId.value}`);
    return;
  }
  router.push(`/knowledge-base/${kbId.value}/document/${docId.value}`);
};

onMounted(async () => {
  await initLanguage();
  await fetchSystemInfo();
  await fetchVersions();
});
</script>

<style scoped>
.history-content {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-md);
}

.compare-toolbar {
  display: flex;
  align-items: center;
  gap: var(--spacing-sm);
  margin-bottom: var(--spacing-sm);
  flex-wrap: wrap;
}

.compare-vs {
  color: var(--text-light);
  font-size: 13px;
}

</style>
