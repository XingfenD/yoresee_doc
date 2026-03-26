<template>
  <AuthLayout
    :current-language="currentLanguage"
    :is-dark-mode="isDarkMode"
    :system-name="systemName"
    :subtitle="t('register.title')"
    @change-language="handleLanguageChange"
    @toggle-theme="toggleTheme"
  >
    <el-form :model="registerForm" :rules="registerRules" ref="registerFormRef" class="register-form">
      <el-form-item prop="username">
        <el-input
          v-model="registerForm.username"
          :placeholder="t('register.username')"
          prefix-icon="User"
          :disabled="loading"
        />
      </el-form-item>
      <el-form-item prop="email">
        <el-input
          v-model="registerForm.email"
          type="email"
          :placeholder="t('register.email')"
          prefix-icon="Message"
          :disabled="loading"
        />
      </el-form-item>
      <el-form-item prop="password">
        <el-input
          v-model="registerForm.password"
          type="password"
          :placeholder="t('register.password')"
          prefix-icon="Lock"
          :disabled="loading"
          show-password
        />
      </el-form-item>
      <el-form-item prop="confirmPassword">
        <el-input
          v-model="registerForm.confirmPassword"
          type="password"
          :placeholder="t('register.confirmPassword')"
          prefix-icon="Check"
          :disabled="loading"
          show-password
          @keyup.enter="handleRegister"
        />
      </el-form-item>
      <el-form-item v-if="systemRegisterMode === 'invite'" prop="invitationCode">
        <el-input
          v-model="registerForm.invitationCode"
          :placeholder="t('register.invitationCode')"
          prefix-icon="Ticket"
          :disabled="loading"
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
          @click="handleRegister"
          :disabled="loading"
        >
          {{ t('register.signUp') }}
        </el-button>
      </el-form-item>
      <el-form-item class="register-link">
        <router-link to="/login">{{ t('register.haveAccount') }} {{ t('register.signIn') }}</router-link>
      </el-form-item>
    </el-form>
  </AuthLayout>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { useUserStore } from '../store/user';
import { useI18n } from 'vue-i18n';
import AuthLayout from '@/components/AuthLayout.vue';
import { useAuthShell } from '@/composables/useAuthShell';

const router = useRouter();
const userStore = useUserStore();
const { locale, t } = useI18n();
const registerFormRef = ref(null);
const loading = ref(false);
const error = ref('');
const systemRegisterMode = ref('invite');
const {
  systemName,
  isDarkMode,
  currentLanguage,
  handleLanguageChange,
  toggleTheme,
  initLanguage,
  fetchSystemInfo
} = useAuthShell({ locale, userStore });

const registerForm = reactive({
  username: '',
  email: '',
  password: '',
  confirmPassword: '',
  invitationCode: ''
});

const registerRules = {
  username: [
    { required: true, message: t('register.usernameRequired'), trigger: 'blur' }
  ],
  email: [
    { required: true, message: t('register.emailRequired'), trigger: 'blur' },
    { type: 'email', message: t('register.emailFormat'), trigger: 'blur' }
  ],
  password: [
    { required: true, message: t('register.passwordRequired'), trigger: 'blur' }
  ],
  confirmPassword: [
    { required: true, message: t('register.confirmPasswordRequired'), trigger: 'blur' },
    {
      validator: (rule, value, callback) => {
        if (value !== registerForm.password) {
          callback(new Error(t('register.passwordMismatch')));
        } else {
          callback();
        }
      },
      trigger: 'blur'
    }
  ],
  invitationCode: [
    {
      validator: (rule, value, callback) => {
        if (systemRegisterMode.value === 'invite' && !value) {
          callback(new Error(t('register.invitationCodeRequired')));
        } else {
          callback();
        }
      },
      trigger: 'blur'
    }
  ]
};

const handleRegister = async () => {
  if (!registerFormRef.value) return;

  await registerFormRef.value.validate(async (valid) => {
    if (valid) {
      loading.value = true;
      error.value = '';

      try {
        const success = await userStore.register(
          registerForm.username,
          registerForm.password,
          registerForm.email,
          registerForm.invitationCode || null
        );
        if (success) {
          router.push({ path: '/login', query: { email: registerForm.email } });
        } else {
          error.value = t('common.requestFailed');
        }
      } catch (err) {
        error.value = t('common.requestFailed');
      } finally {
        loading.value = false;
      }
    }
  });
};

onMounted(() => {
  fetchSystemInfo((info) => {
    systemRegisterMode.value = info.system_register_mode || 'open';
  });
  initLanguage();
});
</script>

<style scoped>
.register-form {
  width: 100%;
}

.register-form .el-form-item {
  margin-bottom: var(--spacing-md);
}

.register-form .el-input {
  height: 40px;
  border-radius: var(--border-radius-md);
}

.register-form .el-input__wrapper {
  box-shadow: none;
}

.register-form .el-input__wrapper.is-focus {
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
