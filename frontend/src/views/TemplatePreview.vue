<template>
  <PageLayout
    :system-name="systemName"
    :current-language="currentLanguage"
    :is-dark-mode="isDarkMode"
    :user-avatar="userAvatar"
    :username="userInfo?.username || t('common.unknown')"
    :active-menu="activeMenu"
    :title="t('templates.previewTitle')"
    @change-language="handleLanguageChange"
    @toggle-theme="toggleTheme"
    @logout="handleLogout"
    @menu-select="handleMenuSelect"
  >
    <template #actions>
      <el-button size="small" @click="goBack">{{ t('common.back') }}</el-button>
    </template>

    <div class="template-preview" v-loading="loading">
      <div class="template-preview-header">
        <div class="template-preview-title">
          {{ template?.name || t('templates.untitled') }}
        </div>
        <div class="template-preview-meta">
          <el-tag v-if="scopeLabel" size="small" :type="scopeTagType">
            {{ scopeLabel }}
          </el-tag>
          <span class="template-preview-date" v-if="template?.updated_at || template?.updatedAt">
            {{ t('templates.updatedAt') }}: {{ formatDate(template?.updated_at || template?.updatedAt) }}
          </span>
        </div>
        <div v-if="template?.description" class="template-preview-desc">
          {{ template.description }}
        </div>
      </div>

      <div class="template-preview-body">
        <div v-if="!previewContent" class="template-preview-empty">
          <el-empty :description="t('templates.contentEmpty')" />
        </div>
        <div v-else ref="previewRef" class="template-preview-content"></div>
      </div>
    </div>
  </PageLayout>
</template>

<script setup>
import { computed, nextTick, onMounted, ref, watch } from 'vue';
import { useI18n } from 'vue-i18n';
import { useRoute, useRouter } from 'vue-router';
import { ElMessage } from 'element-plus';
import Vditor from 'vditor';
import 'vditor/dist/index.css';
import { useUserStore } from '@/store/user';
import PageLayout from '@/components/PageLayout.vue';
import { getTemplate } from '@/services/api';

const route = useRoute();
const router = useRouter();
const { locale, t } = useI18n();
const userStore = useUserStore();

const systemName = ref(userStore.systemName || 'Yoresee');
const activeMenu = ref('templates');
const isDarkMode = computed(() => userStore.darkMode);
const currentLanguage = computed({
  get: () => locale.value,
  set: (val) => (locale.value = val)
});
const userInfo = computed(() => userStore.userInfo);
const userAvatar = computed(
  () => userInfo.value?.avatar || 'https://cube.elemecdn.com/0/88/03b0d39583f48206768a7534e55bcpng.png'
);

const loading = ref(false);
const template = ref(null);
const previewRef = ref(null);

const templateId = computed(() => route.params.templateId || route.params.id);

const formatDate = (dateString) => {
  if (!dateString) return t('common.unknown');
  const date = new Date(dateString);
  return date.toLocaleDateString();
};

const scopeLabel = computed(() => {
  const scope = template.value?.scope;
  if (scope === 'system') return t('templates.public');
  if (scope === 'knowledge_base') return t('templates.scopeKb');
  if (scope === 'private') return t('templates.private');
  return '';
});

const scopeTagType = computed(() => {
  const scope = template.value?.scope;
  if (scope === 'system') return 'success';
  if (scope === 'knowledge_base') return 'warning';
  return 'info';
});

const parseTemplateContent = (raw) => {
  if (!raw) return '';
  try {
    const parsed = JSON.parse(raw);
    if (parsed && typeof parsed.content === 'string') {
      return parsed.content;
    }
  } catch (error) {
    // not json
  }
  return raw;
};

const previewContent = computed(() => parseTemplateContent(template.value?.content || ''));

const renderPreview = async () => {
  if (!previewRef.value) return;
  await nextTick();
  const content = previewContent.value || '';
  if (!content) {
    previewRef.value.innerHTML = '';
    return;
  }
  await Vditor.preview(previewRef.value, content, {
    theme: isDarkMode.value ? 'dark' : 'classic',
    hljs: { style: isDarkMode.value ? 'monokai' : 'github' }
  });
};

const fetchTemplate = async () => {
  if (!templateId.value) return;
  loading.value = true;
  try {
    const data = await getTemplate(templateId.value);
    template.value = data.template;
  } catch (error) {
    console.error('获取模板失败:', error);
    ElMessage.error(t('common.requestFailed'));
  } finally {
    loading.value = false;
  }
};

const goBack = () => {
  router.push('/templates');
};

const handleMenuSelect = (key) => {
  activeMenu.value = key;
};

const handleLanguageChange = (command) => {
  currentLanguage.value = command;
  localStorage.setItem('language', command);
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

watch(previewContent, () => {
  renderPreview();
});

watch(isDarkMode, () => {
  renderPreview();
});

onMounted(async () => {
  initLanguage();
  await fetchTemplate();
  renderPreview();
});
</script>

<style scoped>
.template-preview {
  background: var(--bg-white);
  border-radius: var(--border-radius-md);
  box-shadow: var(--shadow-sm);
  padding: var(--spacing-lg);
}

.template-preview-header {
  padding-bottom: var(--spacing-md);
  border-bottom: 1px solid var(--border-color);
  margin-bottom: var(--spacing-md);
}

.template-preview-title {
  font-size: 22px;
  font-weight: 600;
  color: var(--text-primary);
  margin-bottom: 6px;
}

.template-preview-meta {
  display: flex;
  align-items: center;
  gap: 12px;
  color: var(--text-medium);
  font-size: 12px;
}

.template-preview-date {
  color: var(--text-light);
}

.template-preview-desc {
  margin-top: 8px;
  color: var(--text-medium);
}

.template-preview-body {
  min-height: 220px;
}

.template-preview-content {
  padding: 8px 0;
}

.template-preview-empty {
  padding: var(--spacing-lg) 0;
}

.dark-mode .template-preview {
  background: var(--bg-white);
}
</style>
