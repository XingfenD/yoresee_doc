<template>
  <PageLayout
    :system-name="systemName"
    :current-language="currentLanguage"
    :is-dark-mode="isDarkMode"
    :user-avatar="userAvatar"
    :username="userInfo?.username || t('common.unknown')"
    :active-menu="activeMenu"
    :title="t('templates.title')"
    layout="list"
    @change-language="handleLanguageChange"
    @toggle-theme="toggleTheme"
    @logout="handleLogout"
    @menu-select="handleMenuSelect"
  >
    <template #actions>
      <el-button class="page-action-btn" type="primary" size="small" @click="openCreateTemplateDialog">
        {{ t('templates.createNew') }}
      </el-button>
    </template>
    <div class="templates-layout">
      <el-tabs v-model="activeTab" class="templates-tabs">
        <el-tab-pane :label="t('templates.my')" name="my">
          <div v-loading="loadingMy">
            <TemplateListSection
              :title="t('templates.my')"
              :items="myTemplates"
              :empty-text="t('templates.noMy')"
              :fallback-description="t('templates.noDescription')"
              :tag-mapper="myTagMapper"
              :meta-mapper="templateMetaMapper"
              :action-label="t('common.open')"
              @open="openTemplate"
            />
          </div>
        </el-tab-pane>

        <el-tab-pane :label="t('templates.recent')" name="recent">
          <div v-loading="loadingRecent">
            <TemplateListSection
              :title="t('templates.recent')"
              :items="recentTemplates"
              :empty-text="t('templates.noRecent')"
              :fallback-description="t('templates.noDescription')"
              :tag-mapper="recentTagMapper"
              :meta-mapper="templateMetaMapper"
              :action-label="t('common.open')"
              @open="openTemplate"
            />
          </div>
        </el-tab-pane>

        <el-tab-pane :label="t('templates.public')" name="public">
          <div v-loading="loadingPublic">
            <TemplateListSection
              :title="t('templates.public')"
              :items="publicTemplates"
              :empty-text="t('templates.noPublic')"
              :fallback-description="t('templates.noDescription')"
              :tag-mapper="publicTagMapper"
              :meta-mapper="templateMetaMapper"
              :action-label="t('common.open')"
              @open="openTemplate"
            />
          </div>
        </el-tab-pane>
      </el-tabs>
    </div>
  </PageLayout>

  <TemplateCreateDialog
    v-model="showCreateDialog"
    :loading="creatingTemplate"
    :title="t('templates.createNewTitle')"
    :show-content="true"
    :show-kb-scope="false"
    :initial-name="templateDialogInit.name"
    :initial-description="templateDialogInit.description"
    :initial-scope="templateDialogInit.scope"
    :initial-tags="templateDialogInit.tags"
    :initial-content="templateDialogInit.content"
    @submit="submitCreateTemplate"
  />
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue';
import { useRouter } from 'vue-router';
import { useI18n } from 'vue-i18n';
import { ElMessage } from 'element-plus';
import { useUserStore } from '@/store/user';
import PageLayout from '@/components/PageLayout.vue';
import TemplateListSection from '@/components/TemplateListSection.vue';
import TemplateCreateDialog from '@/components/TemplateCreateDialog.vue';
import { listTemplates, listRecentTemplates, createTemplate as createTemplateApi } from '@/services/api';

const router = useRouter();
const userStore = useUserStore();
const { locale, t } = useI18n();

const systemName = ref(userStore.systemName || 'Yoresee');
const activeMenu = ref('templates');
const activeTab = ref('my');
const isDarkMode = computed(() => userStore.darkMode);
const showCreateDialog = ref(false);
const creatingTemplate = ref(false);
const templateDialogInit = ref({
  name: '',
  description: '',
  scope: 'own',
  tags: '',
  content: ''
});

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

const myTemplates = ref([]);
const recentTemplates = ref([]);
const publicTemplates = ref([]);
const loadingMy = ref(false);
const loadingRecent = ref(false);
const loadingPublic = ref(false);

const myTagMapper = () => ({ type: 'info', label: t('templates.private') });
const recentTagMapper = (tpl) =>
  tpl.scope === 'system'
    ? { type: 'success', label: t('templates.public') }
    : { type: 'info', label: t('templates.private') };
const publicTagMapper = () => ({ type: 'success', label: t('templates.public') });

const templateMetaMapper = (tpl) => [
  { label: t('templates.updatedAt'), value: formatDate(tpl.updated_at || tpl.updatedAt) }
];

const openTemplate = (tpl) => {
  if (!tpl?.id) return;
  router.push(`/template/${tpl.id}`);
};

const openCreateTemplateDialog = () => {
  templateDialogInit.value = {
    name: '',
    description: '',
    scope: 'own',
    tags: '',
    content: ''
  };
  showCreateDialog.value = true;
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

const formatDate = (value) => {
  if (!value) return t('common.unknown');
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) return value;
  return date.toLocaleDateString();
};

const fetchMyTemplates = async () => {
  if (loadingMy.value) return;
  loadingMy.value = true;
  try {
    const data = await listTemplates({
      only_mine: true,
      target_container: 'own',
      order_by: 'updated_at',
      order_desc: true,
      page: 1,
      page_size: 50
    });
    myTemplates.value = data.templates || [];
  } catch (err) {
    console.error('获取我的模板失败:', err);
  } finally {
    loadingMy.value = false;
  }
};

const fetchRecentTemplates = async () => {
  if (loadingRecent.value) return;
  loadingRecent.value = true;
  try {
    const data = await listRecentTemplates({
      page: 1,
      page_size: 50
    });
    recentTemplates.value = data.templates || [];
  } catch (err) {
    console.error('获取最近模板失败:', err);
  } finally {
    loadingRecent.value = false;
  }
};

const fetchPublicTemplates = async () => {
  if (loadingPublic.value) return;
  loadingPublic.value = true;
  try {
    const data = await listTemplates({
      target_container: 'public',
      order_by: 'updated_at',
      order_desc: true,
      page: 1,
      page_size: 50
    });
    publicTemplates.value = data.templates || [];
  } catch (err) {
    console.error('获取公开模板失败:', err);
  } finally {
    loadingPublic.value = false;
  }
};

const submitCreateTemplate = async (payload) => {
  if (creatingTemplate.value) return;
  try {
    creatingTemplate.value = true;
    const requestBody = {
      target_container: payload.scope,
      template_content: JSON.stringify({
        name: payload.name,
        description: payload.description,
        content: payload.content,
        tags: payload.tags || []
      })
    };
    await createTemplateApi(requestBody);
    showCreateDialog.value = false;
    ElMessage.success(t('templates.saveSuccess'));
    fetchTemplatesForTab(activeTab.value);
  } catch (err) {
    console.error('创建模板失败:', err);
    ElMessage.error(t('templates.saveFailed'));
  } finally {
    creatingTemplate.value = false;
  }
};

const fetchTemplatesForTab = (tab) => {
  if (tab === 'my') {
    fetchMyTemplates();
  } else if (tab === 'recent') {
    fetchRecentTemplates();
  } else if (tab === 'public') {
    fetchPublicTemplates();
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

onMounted(() => {
  fetchSystemInfo();
  initLanguage();
  fetchTemplatesForTab(activeTab.value);
});

watch(activeTab, (tab) => {
  fetchTemplatesForTab(tab);
});
</script>

<style scoped>
.templates-layout {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-lg);
  height: auto;
}

.templates-tabs :deep(.el-tabs__header) {
  margin: 0 0 var(--spacing-md);
  border-bottom: 1px solid var(--border-color);
}

.templates-tabs :deep(.el-tabs__nav-wrap) {
  padding: 0 var(--spacing-sm);
}

.templates-tabs :deep(.el-tabs__item) {
  color: var(--text-medium);
  font-weight: 500;
}

.templates-tabs :deep(.el-tabs__item.is-active) {
  color: var(--primary-color);
}

.templates-tabs :deep(.el-tabs__active-bar) {
  background-color: var(--primary-color);
}
</style>
