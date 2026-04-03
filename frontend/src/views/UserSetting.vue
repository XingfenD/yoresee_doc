<template>
  <PageLayout
    :system-name="systemName"
    :current-language="currentLanguage"
    :is-dark-mode="isDarkMode"
    :user-avatar="userAvatar"
    :username="userInfo?.username || t('common.user')"
    :active-menu="activeMenu"
    :side-menu-items="userMenuItems"
    sidebar-scene="user_info"
    :title="t('user.settingTitle')"
    @change-language="handleLanguageChange"
    @toggle-theme="toggleTheme"
    @logout="handleLogout"
    @menu-select="handleMenuSelect"
  >
    <template #actions>
      <el-button class="page-action-btn" size="small" @click="resetForm">
        {{ t('button.cancel') }}
      </el-button>
      <el-button class="page-action-btn" size="small" :loading="saving" type="primary" @click="handleSave">
        {{ t('common.save') }}
      </el-button>
    </template>

    <div class="setting-page">
      <div class="setting-card">
        <div class="setting-row setting-row--avatar">
          <div class="setting-label-group">
            <div class="setting-label">{{ t('user.avatar') }}</div>
            <div class="setting-desc">{{ t('user.avatarHint') }}</div>
          </div>
          <div class="avatar-editor">
            <el-avatar :size="64" :src="previewAvatar" />
            <el-upload
              ref="avatarUploadRef"
              :auto-upload="false"
              :show-file-list="false"
              accept="image/jpeg,image/png,image/webp"
              :before-upload="() => false"
              :on-change="handleAvatarFileChange"
            >
              <el-button class="avatar-upload-btn" size="small" :loading="avatarProcessing">
                {{ t('user.uploadAvatar') }}
              </el-button>
            </el-upload>
          </div>
        </div>

        <div class="setting-row">
          <div class="setting-label-group">
            <div class="setting-label">{{ t('user.name') }}</div>
          </div>
          <el-input :model-value="userInfo?.username || ''" disabled />
        </div>

        <div class="setting-row">
          <div class="setting-label-group">
            <div class="setting-label">{{ t('user.nickname') }}</div>
          </div>
          <el-input v-model="form.nickname" :placeholder="t('user.nicknamePlaceholder')" />
        </div>

        <div class="setting-row">
          <div class="setting-label-group">
            <div class="setting-label">{{ t('user.email') }}</div>
          </div>
          <el-input v-model="form.email" :placeholder="t('login.email')" />
        </div>

        <div class="setting-row">
          <div class="setting-label-group">
            <div class="setting-label">{{ t('user.newPassword') }}</div>
            <div class="setting-desc">{{ t('user.passwordHint') }}</div>
          </div>
          <el-input
            v-model="form.password"
            show-password
            clearable
            :placeholder="t('user.newPasswordPlaceholder')"
          />
        </div>

        <div class="setting-row">
          <div class="setting-label-group">
            <div class="setting-label">{{ t('user.confirmPassword') }}</div>
          </div>
          <el-input
            v-model="form.confirmPassword"
            show-password
            clearable
            :placeholder="t('user.confirmPassword')"
          />
        </div>

      </div>
    </div>
  </PageLayout>
</template>

<script setup>
import { computed, onBeforeUnmount, onMounted, reactive, ref, watch } from 'vue';
import { useI18n } from 'vue-i18n';
import { useRouter } from 'vue-router';
import { ElMessage } from 'element-plus';
import { useUserStore } from '@/store/user';
import PageLayout from '@/components/PageLayout.vue';
import { useUserShell } from '@/composables/useUserShell';
import { usePageBoot } from '@/composables/usePageBoot';

const router = useRouter();
const userStore = useUserStore();
const { locale, t } = useI18n();

const {
  systemName,
  activeMenu,
  isDarkMode,
  userInfo,
  userAvatar,
  userMenuItems,
  currentLanguage,
  initLanguage,
  fetchSystemInfo,
  handleLanguageChange,
  toggleTheme,
  handleLogout,
  handleMenuSelect
} = useUserShell({
  locale,
  router,
  userStore,
  defaultActiveMenu: 'user-security'
});
const { boot } = usePageBoot({ initLanguage, fetchSystemInfo });

const form = reactive({
  email: '',
  nickname: '',
  password: '',
  confirmPassword: ''
});
const avatarUploadRef = ref();
const avatarProcessing = ref(false);
const saving = ref(false);
const selectedAvatar = ref(null);
const maxAvatarSizeBytes = 5 * 1024 * 1024;
const avatarContentTypeByExt = {
  '.jpg': 'image/jpeg',
  '.jpeg': 'image/jpeg',
  '.png': 'image/png',
  '.webp': 'image/webp'
};

const previewAvatar = computed(() => {
  if (selectedAvatar.value?.previewUrl) {
    return selectedAvatar.value.previewUrl;
  }
  return userAvatar.value;
});

const revokeSelectedPreview = () => {
  const url = selectedAvatar.value?.previewUrl;
  if (url) {
    URL.revokeObjectURL(url);
  }
};

const resetForm = () => {
  revokeSelectedPreview();
  form.email = userInfo.value?.email || '';
  form.nickname = userInfo.value?.nickname || '';
  form.password = '';
  form.confirmPassword = '';
  selectedAvatar.value = null;
  avatarUploadRef.value?.clearFiles?.();
};

watch(
  () => userInfo.value,
  () => {
    resetForm();
  },
  { immediate: true }
);

const resolveAvatarContentType = (file) => {
  let rawType = String(file?.type || '').trim().toLowerCase();
  if (rawType === 'image/jpg') {
    rawType = 'image/jpeg';
  }
  if (rawType === 'image/jpeg' || rawType === 'image/png' || rawType === 'image/webp') {
    return rawType;
  }
  const name = String(file?.name || '').trim().toLowerCase();
  const ext = name.lastIndexOf('.') > -1 ? name.slice(name.lastIndexOf('.')) : '';
  return avatarContentTypeByExt[ext] || '';
};

const handleAvatarFileChange = async (uploadFile) => {
  const rawFile = uploadFile?.raw;
  if (!rawFile || avatarProcessing.value) return;

  if (rawFile.size > maxAvatarSizeBytes) {
    ElMessage.error(t('user.avatarSizeError'));
    avatarUploadRef.value?.clearFiles?.();
    return;
  }
  const contentType = resolveAvatarContentType(rawFile);
  if (!contentType) {
    ElMessage.error(t('user.avatarTypeError'));
    avatarUploadRef.value?.clearFiles?.();
    return;
  }

  avatarProcessing.value = true;
  try {
    revokeSelectedPreview();
    selectedAvatar.value = {
      file: rawFile,
      contentType,
      previewUrl: URL.createObjectURL(rawFile)
    };
  } finally {
    avatarProcessing.value = false;
  }
};

const buildPayload = async () => {
  const payload = {};
  const originalEmail = (userInfo.value?.email || '').trim();
  const originalNickname = userInfo.value?.nickname ?? '';

  const nextEmail = form.email.trim();
  if (nextEmail && nextEmail !== originalEmail) {
    payload.email = nextEmail;
  }

  if (form.nickname !== originalNickname) {
    payload.nickname = form.nickname;
  }

  const password = form.password.trim();
  if (password) {
    payload.password = password;
  }

  if (selectedAvatar.value?.file) {
    payload.avatar_file = new Uint8Array(await selectedAvatar.value.file.arrayBuffer());
    payload.avatar_filename = selectedAvatar.value.file.name || 'avatar';
    payload.avatar_content_type = selectedAvatar.value.contentType;
  }

  return payload;
};

const handleSave = async () => {
  if (form.password !== form.confirmPassword) {
    ElMessage.error(t('user.passwordMismatch'));
    return;
  }
  saving.value = true;
  try {
    const payload = await buildPayload();
    if (!Object.keys(payload).length) {
      ElMessage.info(t('user.noChanges'));
      return;
    }
    await userStore.updateProfile(payload);
    ElMessage.success(t('user.saveProfileSuccess'));
    resetForm();
  } catch (error) {
    ElMessage.error(t('user.saveProfileFailed'));
  } finally {
    saving.value = false;
  }
};

onMounted(() => {
  boot();
});

onBeforeUnmount(() => {
  revokeSelectedPreview();
});
</script>

<style scoped>
.setting-page {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-lg);
}

.setting-card {
  background: var(--bg-white);
  border: 1px solid var(--border-color);
  border-radius: var(--border-radius-md);
  box-shadow: var(--shadow-sm);
  padding: var(--spacing-lg);
  display: flex;
  flex-direction: column;
  gap: var(--spacing-md);
}

.setting-row {
  display: grid;
  grid-template-columns: 220px minmax(0, 1fr);
  gap: var(--spacing-md);
  align-items: center;
}

.setting-row--avatar {
  align-items: flex-start;
}

.setting-label-group {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.setting-label {
  color: var(--text-dark);
  font-size: 14px;
  font-weight: 500;
}

.setting-desc {
  color: var(--text-light);
  font-size: 12px;
}

.avatar-editor {
  display: flex;
  align-items: center;
  gap: 12px;
}

.avatar-upload-btn {
  border-color: var(--border-color);
  background: var(--bg-white);
  color: var(--text-medium);
}

.avatar-upload-btn:hover,
.avatar-upload-btn:focus {
  border-color: var(--primary-color);
  color: var(--primary-color);
  background: var(--primary-light);
}

@media (max-width: 768px) {
  .setting-row {
    grid-template-columns: 1fr;
    align-items: stretch;
  }
}
</style>
