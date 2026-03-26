<template>
  <PageLayout
    :system-name="systemName"
    :current-language="currentLanguage"
    :is-dark-mode="isDarkMode"
    :user-avatar="userAvatar"
    :username="userInfo?.username || '用户'"
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
    <div class="manage-layout">
      <section v-for="group in settingGroups" :key="group.key" class="manage-section">
        <div class="section-header">
          <h3 class="section-title">{{ resolveText(group.title_key, group.title) }}</h3>
        </div>
        <div class="section-body">
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
        </div>
      </section>
    </div>
  </PageLayout>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { useI18n } from 'vue-i18n';
import { useRouter } from 'vue-router';
import { useUserStore } from '@/store/user';
import PageLayout from '@/components/PageLayout.vue';
import { useManageShell } from '@/composables/useManageShell';
import { getSettings, updateSettings } from '@/services/api';
import { ElMessage } from 'element-plus';

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
  try {
    const resp = await getSettings('system');
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
  } catch (error) {
    settingGroups.value = [];
  }
};

const handleSave = async () => {
  if (isSaving.value) {
    return;
  }
  isSaving.value = true;
  try {
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
    ElMessage.success(t('message.saveSuccessGeneric'));
  } catch (error) {
    ElMessage.error(t('message.saveFailedGeneric'));
  } finally {
    isSaving.value = false;
  }
};

onMounted(() => {
  initLanguage();
  loadSettings();
});
</script>

<style scoped>
.manage-layout {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-lg);
}

.manage-section {
  background: var(--bg-white);
  border: 1px solid var(--border-color);
  border-radius: var(--border-radius-md);
  box-shadow: var(--shadow-sm);
  overflow: hidden;
}

.section-header {
  padding: var(--spacing-md);
  border-bottom: 1px solid var(--border-color);
  background: var(--bg-white);
}

.section-title {
  margin: 0;
  font-size: 16px;
  font-weight: 600;
  color: var(--text-dark);
}

.section-body {
  padding: var(--spacing-md);
}

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
