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
        <el-button type="primary" @click="saveSettingsSample">
          {{ t('common.save') }}
        </el-button>
      </template>
    </TitleBar>

    <section class="settings-panel">
      <div class="settings-header">
        <h3>{{ t('document.settings.sampleSectionTitle') }}</h3>
        <p>{{ t('document.settings.sampleSectionHint') }}</p>
      </div>

      <el-form label-position="top" class="settings-form">
        <el-form-item :label="t('document.settings.docIdLabel')">
          <el-input :model-value="docId" disabled />
        </el-form-item>

        <el-form-item :label="t('document.settings.defaultPermission')">
          <el-radio-group v-model="settings.defaultPermission">
            <el-radio value="private">{{ t('document.settings.permissionPrivate') }}</el-radio>
            <el-radio value="team">{{ t('document.settings.permissionTeam') }}</el-radio>
          </el-radio-group>
        </el-form-item>

        <el-form-item :label="t('document.settings.autoSaveInterval')">
          <el-select v-model="settings.autoSaveInterval" style="width: 100%">
            <el-option
              v-for="option in autoSaveOptions"
              :key="option.value"
              :label="option.label"
              :value="option.value"
            />
          </el-select>
        </el-form-item>

        <el-form-item :label="t('document.settings.enableAttachments')">
          <el-switch v-model="settings.enableAttachments" />
        </el-form-item>

        <el-form-item :label="t('document.settings.enableInlineComments')">
          <el-switch v-model="settings.enableInlineComments" />
        </el-form-item>
      </el-form>
    </section>
  </PageLayout>
</template>

<script setup>
import { computed, onMounted, reactive } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { useI18n } from 'vue-i18n';
import { ElMessage } from 'element-plus';
import PageLayout from '@/components/PageLayout.vue';
import TitleBar from '@/components/TitleBar.vue';
import { useWorkspaceShell } from '@/composables/useWorkspaceShell';
import { useUserStore } from '@/store/user';

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

const settings = reactive({
  defaultPermission: 'private',
  autoSaveInterval: '30s',
  enableAttachments: true,
  enableInlineComments: true
});

const autoSaveOptions = computed(() => [
  { value: '10s', label: t('document.settings.autoSave10s') },
  { value: '30s', label: t('document.settings.autoSave30s') },
  { value: '60s', label: t('document.settings.autoSave60s') }
]);

const goBackToDocument = () => {
  if (kbId.value === 'personal') {
    router.push(`/mydocument/${docId.value}`);
    return;
  }
  router.push(`/knowledge-base/${kbId.value}/document/${docId.value}`);
};

const saveSettingsSample = () => {
  ElMessage.success(t('document.settings.sampleSaved'));
};

onMounted(async () => {
  await initLanguage();
  await fetchSystemInfo();
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

.settings-header h3 {
  margin: 0 0 6px;
  font-size: 18px;
  color: var(--text-dark);
}

.settings-header p {
  margin: 0;
  color: var(--text-light);
  font-size: 13px;
}

.settings-form {
  max-width: 640px;
}
</style>
