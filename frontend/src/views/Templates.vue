<template>
  <PageLayout
    :system-name="systemName"
    :current-language="currentLanguage"
    :is-dark-mode="isDarkMode"
    :user-avatar="userAvatar"
    :username="userInfo?.username || t('common.unknown')"
    :active-menu="activeMenu"
    :title="t('templates.title')"
    @change-language="handleLanguageChange"
    @toggle-theme="toggleTheme"
    @logout="handleLogout"
    @menu-select="handleMenuSelect"
  >
    <div class="templates-layout">
      <el-tabs v-model="activeTab" class="templates-tabs">
        <el-tab-pane :label="t('templates.my')" name="my">
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
        </el-tab-pane>

        <el-tab-pane :label="t('templates.recent')" name="recent">
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
        </el-tab-pane>

        <el-tab-pane :label="t('templates.public')" name="public">
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
        </el-tab-pane>
      </el-tabs>
    </div>
  </PageLayout>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue';
import { useRouter } from 'vue-router';
import { useI18n } from 'vue-i18n';
import { useUserStore } from '@/store/user';
import PageLayout from '@/components/PageLayout.vue';
import TemplateListSection from '@/components/TemplateListSection.vue';

const router = useRouter();
const userStore = useUserStore();
const { locale, t } = useI18n();

const systemName = ref(userStore.systemName || 'Yoresee');
const activeMenu = ref('templates');
const activeTab = ref('my');
const isDarkMode = computed(() => userStore.darkMode);

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

const myTemplates = ref([
  { id: 'tpl-my-1', name: '产品需求模板', description: 'PRD 标准模板', updatedAt: '2026-03-10', isPublic: false },
  { id: 'tpl-my-2', name: '技术方案模板', description: '技术方案统一结构', updatedAt: '2026-03-08', isPublic: false }
]);

const recentTemplates = ref([
  { id: 'tpl-recent-1', name: '周报模板', description: '周报结构化模板', updatedAt: '2026-03-11', isPublic: true },
  { id: 'tpl-recent-2', name: '会议纪要模板', description: '会议纪要整理模板', updatedAt: '2026-03-09', isPublic: true }
]);

const publicTemplates = ref([
  { id: 'tpl-public-1', name: '项目复盘模板', description: '复盘要点与改进项', updatedAt: '2026-03-07', isPublic: true },
  { id: 'tpl-public-2', name: '用户访谈模板', description: '访谈问题与记录', updatedAt: '2026-03-05', isPublic: true }
]);

const myTagMapper = () => ({ type: 'info', label: t('templates.private') });
const recentTagMapper = () => ({ type: 'success', label: t('templates.public') });
const publicTagMapper = () => ({ type: 'success', label: t('templates.public') });

const templateMetaMapper = (tpl) => [
  { label: t('templates.updatedAt'), value: tpl.updatedAt }
];

const openTemplate = (tpl) => {
  console.log('open template', tpl);
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

onMounted(() => {
  fetchSystemInfo();
  initLanguage();
});

watch(activeTab, () => {
  // fake data for now
});
</script>

<style scoped>
.templates-layout {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-lg);
}

.templates-tabs :deep(.el-tabs__header) {
  margin: 0 0 var(--spacing-md);
  border-bottom: 1px solid var(--border-color);
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
