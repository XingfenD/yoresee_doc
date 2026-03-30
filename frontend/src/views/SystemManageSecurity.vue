<template>
  <PageLayout
    :system-name="systemName"
    :current-language="currentLanguage"
    :is-dark-mode="isDarkMode"
    :user-avatar="userAvatar"
    :username="userInfo?.username || t('common.user')"
    :active-menu="activeMenu"
    :side-menu-items="manageMenuItems"
    sidebar-scene="manage"
    :title="t('system.security.title')"
    @change-language="handleLanguageChange"
    @toggle-theme="toggleTheme"
    @logout="handleLogout"
    @menu-select="handleMenuSelect"
  >
    <template #actions>
      <el-button class="page-action-btn" type="primary" size="small" :loading="isSaving" @click="handleSave">
        {{ t('common.save') }}
      </el-button>
    </template>
    <ManageLayout>
      <ManageSection
        v-for="group in settingGroups"
        :key="group.key"
        :title="resolveText(group.title_key, group.title)"
        body-padding="md"
      >
        <div v-for="item in group.items" :key="item.key" class="setting-row setting-row--stacked">
          <div class="setting-label">
            {{ resolveText(item.label_key, item.label) }}
            <span v-if="item.required" class="required-mark">*</span>
          </div>
          <div v-if="item.description || item.description_key" class="setting-desc">
            {{ resolveText(item.description_key, item.description) }}
          </div>
          <div class="setting-control">
            <el-radio-group
              v-if="item.ui?.component === 'radio'"
              v-model="settingValues[item.key]"
              :disabled="item.readonly"
            >
              <el-radio
                v-for="opt in item.ui.options"
                :key="opt.value"
                :value="opt.value"
              >
                {{ resolveText(opt.label_key, opt.label) }}
              </el-radio>
            </el-radio-group>
            <el-select
              v-else-if="item.ui?.component === 'select'"
              v-model="settingValues[item.key]"
              :placeholder="resolveText(item.ui?.placeholder_key, item.ui?.placeholder)"
              :disabled="item.readonly"
            >
              <el-option
                v-for="opt in item.ui.options"
                :key="opt.value"
                :label="resolveText(opt.label_key, opt.label)"
                :value="opt.value"
              />
            </el-select>
            <el-switch
              v-else-if="item.ui?.component === 'switch'"
              v-model="settingValues[item.key]"
              :disabled="item.readonly"
            />
            <el-input
              v-else
              v-model="settingValues[item.key]"
              :placeholder="resolveText(item.ui?.placeholder_key, item.ui?.placeholder)"
              :disabled="item.readonly"
            />
          </div>
        </div>
      </ManageSection>
    </ManageLayout>
  </PageLayout>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { useI18n } from 'vue-i18n';
import { useRouter } from 'vue-router';
import { useUserStore } from '@/store/user';
import PageLayout from '@/components/PageLayout.vue';
import ManageLayout from '@/components/ManageLayout.vue';
import ManageSection from '@/components/ManageSection.vue';
import { useManageShell } from '@/composables/useManageShell';
import { usePageBoot } from '@/composables/usePageBoot';
import { useApiAction } from '@/composables/useApiAction';
import { getSettings, updateSettings } from '@/services/api';

const router = useRouter();
const userStore = useUserStore();
const { locale, t } = useI18n();

const {
  systemName,
  activeMenu,
  isDarkMode,
  userInfo,
  userAvatar,
  manageMenuItems,
  currentLanguage,
  initLanguage,
  fetchSystemInfo,
  handleLanguageChange,
  toggleTheme,
  handleLogout,
  handleMenuSelect
} = useManageShell({
  locale,
  router,
  userStore,
  defaultActiveMenu: 'manage-security'
});
const { boot } = usePageBoot({ initLanguage, fetchSystemInfo });
const { runSilent, runWithLoading } = useApiAction({ t });
const isSaving = ref(false);
const settingGroups = ref([]);
const settingValues = ref({});

const resolveText = (key, fallback) => {
  if (key) {
    return t(key);
  }
  return fallback || '';
};

const loadSettings = async () => {
  await runSilent(
    () => getSettings('system'),
    {
      context: 'loadSystemSettings',
      onSuccess: (resp) => {
        settingGroups.value = resp.groups || [];
        const nextValues = {};
        settingGroups.value.forEach((group) => {
          (group.items || []).forEach((item) => {
            if (item.type === 'bool') {
              nextValues[item.key] = item.value === 'true';
            } else {
              nextValues[item.key] = item.value ?? item.default_value ?? '';
            }
          });
        });
        settingValues.value = nextValues;
      },
      onError: () => {
        settingGroups.value = [];
        settingValues.value = {};
      }
    }
  );
};

const handleSave = async () => {
  await runWithLoading(
    isSaving,
    async () => {
      const updates = [];
      settingGroups.value.forEach((group) => {
        (group.items || []).forEach((item) => {
          let value = settingValues.value[item.key];
          if (item.type === 'bool') {
            value = value ? 'true' : 'false';
          } else if (value == null) {
            value = '';
          } else {
            value = String(value);
          }
          updates.push({ key: item.key, value });
        });
      });
      await updateSettings(updates);
    },
    {
      context: 'saveSystemSettings',
      successMessage: t('message.saveSuccessGeneric'),
      errorMessage: t('message.saveFailedGeneric')
    }
  );
};

onMounted(() => {
  boot(loadSettings);
});
</script>

<style scoped>
.setting-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: var(--spacing-md);
}

.setting-row--stacked {
  align-items: flex-start;
  flex-direction: column;
}

.setting-label {
  color: var(--text-medium);
  font-size: 14px;
}

.setting-control {
  min-width: 200px;
}

.setting-desc {
  color: var(--text-light);
  font-size: 12px;
}

.required-mark {
  color: #f56c6c;
  margin-left: 4px;
}
</style>
