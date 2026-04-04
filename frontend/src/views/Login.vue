<template>
  <AuthLayout
    :current-language="currentLanguage"
    :is-dark-mode="isDarkMode"
    :system-name="systemName"
    :subtitle="t('login.title')"
    @change-language="handleLanguageChange"
    @toggle-theme="toggleTheme"
  >
    <el-form
      :model="loginForm"
      :rules="loginRules"
      ref="loginFormRef"
      class="login-form"
    >
      <el-form-item prop="email">
        <el-input
          v-model="loginForm.email"
          :placeholder="t('login.email')"
          prefix-icon="Message"
          :disabled="loading"
        />
      </el-form-item>
      <el-form-item prop="password">
        <el-input
          v-model="loginForm.password"
          type="password"
          :placeholder="t('login.password')"
          prefix-icon="Lock"
          :disabled="loading"
          show-password
          @keyup.enter="handleLogin"
        />
      </el-form-item>
      <el-form-item v-if="error" class="error-message">
        <el-alert
          :title="error"
          type="error"
          show-icon
          :closable="false"
          class="error-alert"
        />
      </el-form-item>
      <el-form-item>
        <el-button
          type="primary"
          class="login-button"
          :loading="loading"
          @click="handleLogin"
          :disabled="loading"
        >
          {{ t('login.signIn') }}
        </el-button>
      </el-form-item>
      <el-form-item class="register-link">
        <router-link to="/register">{{ t('login.noAccount') }} {{ t('login.signUp') }}</router-link>
      </el-form-item>
    </el-form>
  </AuthLayout>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import { useUserStore } from '../store/user';
import { useI18n } from 'vue-i18n';
import AuthLayout from '@/components/layout/AuthLayout.vue';
import { useAuthShell } from '@/composables/shell/useAuthShell';
import { usePageBoot } from '@/composables/shell/usePageBoot';

const router = useRouter();
const route = useRoute();
const userStore = useUserStore();
const { locale, t } = useI18n();
const loginFormRef = ref(null);
const loading = ref(false);
const error = ref('');
const {
  systemName,
  isDarkMode,
  currentLanguage,
  handleLanguageChange,
  toggleTheme,
  initLanguage,
  fetchSystemInfo
} = useAuthShell({ locale, userStore });
const { boot } = usePageBoot({ initLanguage, fetchSystemInfo });

const loginForm = reactive({
  email: route.query.email || '',
  password: ''
});

const loginRules = computed(() => ({
  email: [
    { required: true, message: t('login.validation.emailRequired'), trigger: 'blur' },
    { type: 'email', message: t('login.validation.emailFormat'), trigger: 'blur' }
  ],
  password: [
    { required: true, message: t('login.validation.passwordRequired'), trigger: 'blur' }
  ]
}));

const handleLogin = async () => {
  if (!loginFormRef.value) return;

  await loginFormRef.value.validate(async (valid) => {
    if (valid) {
      loading.value = true;
      error.value = '';

      try {
        const success = await userStore.login(loginForm.email, loginForm.password);
        if (success) {
          router.push('/');
        } else {
          error.value = t('login.failed');
        }
      } catch (err) {
        error.value = t('login.failed');
      } finally {
        loading.value = false;
      }
    }
  });
};

onMounted(() => {
  boot();
});
</script>

<style scoped>
.login-form {
  width: 100%;
}

.login-form .el-form-item {
  margin-bottom: var(--spacing-md);
}

.login-form .el-input {
  height: 40px;
  border-radius: var(--border-radius-md);
}

.login-form .el-input__wrapper {
  box-shadow: none;
}

.login-form .el-input__wrapper.is-focus {
  box-shadow: 0 0 0 2px var(--primary-light);
}

.login-button {
  width: 100%;
  height: 40px;
  border-radius: var(--border-radius-md);
  font-size: 14px;
  font-weight: 500;
  background-color: var(--primary-color);
  border-color: var(--primary-color);
}

.login-button:hover,
.login-button:active {
  background-color: var(--primary-dark);
  border-color: var(--primary-dark);
}

.error-message {
  margin-bottom: var(--spacing-md);
}

.error-alert {
  padding: 8px 12px;
  font-size: 12px;
}

.register-link {
  text-align: center;
  margin-top: var(--spacing-md);
}

.register-link a {
  color: var(--primary-color);
  font-size: 14px;
  text-decoration: none;
}

.register-link a:hover {
  text-decoration: underline;
}
</style>
