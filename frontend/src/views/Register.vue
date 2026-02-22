<template>
  <div class="login-container">
    <!-- 顶部导航栏 -->
    <header class="register-nav">
      <div class="nav-right">
        <!-- 语言切换 -->
        <el-dropdown trigger="click" @command="handleLanguageChange" class="nav-item">
          <span class="nav-link">
            <el-icon :size="18"><Flag v-if="currentLanguage === 'en'" /><ChatLineRound v-else /></el-icon>
          </span>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item command="en" :icon="'Flag'">
                {{ t('language.english') }}
              </el-dropdown-item>
              <el-dropdown-item command="zh" :icon="'ChatLineRound'">
                {{ t('language.chinese') }}
              </el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
        
        <!-- 主题切换 -->
        <div class="nav-item theme-switch">
          <span class="nav-link" @click="toggleTheme">
            <el-icon :size="18"><Moon v-if="isDarkMode" /><Sunny v-else /></el-icon>
          </span>
        </div>
      </div>
    </header>
    
    <div class="login-form-wrapper">
      <div class="login-header">
        <h2>{{ systemName }}</h2>
        <p>{{ t('register.title') }}</p>
      </div>
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
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { useUserStore } from '../store/user';
import { useI18n } from 'vue-i18n';
import { Flag, ChatLineRound, Moon, Sunny } from '@element-plus/icons-vue';

const router = useRouter();
const userStore = useUserStore();
const { locale, t } = useI18n();
const registerFormRef = ref(null);
const loading = ref(false);
const error = ref('');
const systemName = ref('Yoresee');
const systemRegisterMode = ref('invite');
const isDarkMode = ref(false);

// 计算当前语言
const currentLanguage = ref(localStorage.getItem('language') || 'en');

// 处理语言切换
const handleLanguageChange = (command) => {
  currentLanguage.value = command;
  locale.value = command;
  localStorage.setItem('language', command);
};

// 处理主题切换
const toggleTheme = () => {
  isDarkMode.value = !isDarkMode.value;
  if (isDarkMode.value) {
    document.documentElement.classList.add('dark-mode');
    localStorage.setItem('darkMode', 'true');
  } else {
    document.documentElement.classList.remove('dark-mode');
    localStorage.setItem('darkMode', 'false');
  }
};

// 初始化主题
const initTheme = () => {
  const savedDarkMode = localStorage.getItem('darkMode');
  if (savedDarkMode === 'true') {
    isDarkMode.value = true;
    document.documentElement.classList.add('dark-mode');
  }
};

// 初始化语言
const initLanguage = () => {
  const savedLanguage = localStorage.getItem('language');
  if (savedLanguage) {
    currentLanguage.value = savedLanguage;
    locale.value = savedLanguage;
  }
};

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
          error.value = userStore.error;
        }
      } catch (err) {
        error.value = '注册失败，请稍后重试';
      } finally {
        loading.value = false;
      }
    }
  });
};

const fetchSystemInfo = async () => {
  try {
    const info = await userStore.fetchSystemInfo();
    systemName.value = info.system_name;
    systemRegisterMode.value = info.system_register_mode || 'open';
  } catch (err) {
    console.error('获取系统信息失败:', err);
  }
};

onMounted(() => {
  fetchSystemInfo();
  initTheme();
  initLanguage();
});
</script>

<style scoped>
.login-container {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  background-color: var(--bg-light);
  padding: var(--spacing-md);
  transition: all 0.3s ease;
  position: relative;
}

/* 注册页面导航栏 */
.register-nav {
  position: absolute;
  top: 0;
  right: 0;
  padding: var(--spacing-md);
}

.nav-right {
  display: flex;
  align-items: center;
  gap: var(--spacing-sm);
}

.nav-item {
  display: flex;
  align-items: center;
  margin-left: var(--spacing-sm);
}

.nav-link {
  display: flex;
  align-items: center;
  padding: var(--spacing-xs) var(--spacing-sm);
  border-radius: var(--border-radius-md);
  color: var(--text-medium);
  transition: all 0.3s ease;
  cursor: pointer;
  
  &:hover {
    background-color: var(--bg-medium);
    color: var(--primary-color);
  }
}

.theme-switch {
  padding: var(--spacing-xs) var(--spacing-sm);
}

.login-form-wrapper {
  margin: auto;
  width: 100%;
  max-width: 400px;
  background-color: var(--bg-white);
  border-radius: var(--border-radius-lg);
  box-shadow: var(--shadow-md);
  padding: var(--spacing-xl);
  transition: all 0.3s ease;
}

.login-form-wrapper:hover {
  box-shadow: var(--shadow-lg);
}

.login-header {
  text-align: center;
  margin-bottom: var(--spacing-lg);
}

.login-header h2 {
  font-size: 24px;
  font-weight: 600;
  color: var(--text-dark);
  margin-bottom: var(--spacing-sm);
}

.login-header p {
  font-size: 14px;
  color: var(--text-light);
  margin: 0;
}

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

.login-button:hover {
  background-color: var(--primary-dark);
  border-color: var(--primary-dark);
}

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

@media (max-width: 768px) {
  .login-form-wrapper {
    padding: var(--spacing-lg) var(--spacing-md);
  }
  
  .login-header h2 {
    font-size: 20px;
  }
}
</style>