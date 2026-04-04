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
        {{ t('document.settings.title') }}
      </template>
      <template #actions>
        <el-button type="primary" :loading="saving" @click="saveSettings">
          {{ t('common.save') }}
        </el-button>
      </template>
    </TitleBar>

    <section class="settings-panel" v-loading="loading">
      <div class="settings-header">
        <p>{{ t('document.settings.publicSectionHint') }}</p>
      </div>

      <el-form label-position="top" class="settings-form">
        <el-form-item :label="t('document.settings.docIdLabel')">
          <el-input :model-value="docId" disabled />
        </el-form-item>

        <el-form-item :label="t('document.settings.publicLabel')">
          <div class="settings-switch-row">
            <el-switch v-model="settings.isPublic" />
            <p class="settings-help">
              {{
                settings.isPublic
                  ? t('document.settings.publicEnabledDesc')
                  : t('document.settings.publicDisabledDesc')
              }}
            </p>
          </div>
        </el-form-item>
      </el-form>
    </section>
  </PageLayout>
</template>

<script setup>
import { computed, onMounted, reactive, ref } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { useI18n } from 'vue-i18n';
import PageLayout from '@/components/layout/PageLayout.vue';
import TitleBar from '@/components/layout/TitleBar.vue';
import { useWorkspaceShell } from '@/composables/shell/useWorkspaceShell';
import { useApiAction } from '@/composables/actions/useApiAction';
import { useUserStore } from '@/store/user';
import { getDocumentSettings, updateDocumentSettings } from '@/services/api';

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
const { runWithLoading } = useApiAction({ t });
const loading = ref(false);
const saving = ref(false);

const settings = reactive({
  isPublic: false
});

const goBackToDocument = () => {
  if (kbId.value === 'personal') {
    router.push(`/mydocument/${docId.value}`);
    return;
  }
  router.push(`/knowledge-base/${kbId.value}/document/${docId.value}`);
};

const loadSettings = async () => {
  if (!docId.value) {
    settings.isPublic = false;
    return;
  }
  await runWithLoading(
    loading,
    () => getDocumentSettings(docId.value),
    {
      context: 'loadDocumentSettings',
      errorMessage: t('common.requestFailed'),
      onSuccess: (resp) => {
        if (typeof resp?.is_public === 'boolean') {
          settings.isPublic = resp.is_public;
        }
      }
    }
  );
};

const saveSettings = async () => {
  if (!docId.value) {
    return;
  }
  await runWithLoading(
    saving,
    () => updateDocumentSettings(docId.value, { is_public: Boolean(settings.isPublic) }),
    {
      context: 'saveDocumentSettings',
      successMessage: t('document.settings.saved'),
      errorMessage: t('common.requestFailed')
    }
  );
};

onMounted(async () => {
  await initLanguage();
  await fetchSystemInfo();
  await loadSettings();
});
</script>

<style scoped>
.settings-panel {
  background: var(--bg-white);
  border: 1px solid var(--border-color);
  border-radius: var(--border-radius-md);
  padding: var(--spacing-lg);
}

.settings-header {
  margin-bottom: var(--spacing-lg);
}

.settings-header p {
  margin: 0;
  color: var(--text-light);
  font-size: 13px;
}

.settings-form {
  max-width: 640px;
}

.settings-switch-row {
  display: flex;
  align-items: center;
  gap: 10px;
}

.settings-help {
  margin: 0;
  font-size: 12px;
  color: var(--text-light);
  line-height: 1;
}
</style>
