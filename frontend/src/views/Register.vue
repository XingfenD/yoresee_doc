<template>
  <div class="login-container">
    <div class="login-form-wrapper">
      <div class="login-header">
        <h2>{{ systemName }}</h2>
        <p>创建新账户</p>
      </div>
      <div class="theme-toggle-container">
        <el-button
          type="text"
          class="theme-toggle"
          @click="toggleDarkMode"
          :icon="darkMode ? 'Sunny' : 'Moon'"
        />
      </div>
      <el-form :model="registerForm" :rules="registerRules" ref="registerFormRef" class="login-form">
        <el-form-item prop="username">
          <el-input
            v-model="registerForm.username"
            placeholder="用户名"
            prefix-icon="User"
            :disabled="loading"
          />
        </el-form-item>
        <el-form-item prop="email">
          <el-input
            v-model="registerForm.email"
            type="email"
            placeholder="邮箱"
            prefix-icon="Message"
            :disabled="loading"
          />
        </el-form-item>
        <el-form-item prop="password">
          <el-input
            v-model="registerForm.password"
            type="password"
            placeholder="密码"
            prefix-icon="Lock"
            :disabled="loading"
            show-password
          />
        </el-form-item>
        <el-form-item prop="confirmPassword">
          <el-input
            v-model="registerForm.confirmPassword"
            type="password"
            placeholder="确认密码"
            prefix-icon="Check"
            :disabled="loading"
            show-password
            @keyup.enter="handleRegister"
          />
        </el-form-item>
        <el-form-item v-if="systemRegisterMode === 'invite'" prop="invitationCode">
          <el-input
            v-model="registerForm.invitationCode"
            placeholder="邀请码"
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
            注册
          </el-button>
        </el-form-item>
        <el-form-item class="register-link">
          <router-link to="/login">已有账户？立即登录</router-link>
        </el-form-item>
      </el-form>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from 'vue';
import { useRouter } from 'vue-router';
import { useUserStore } from '../store/user';
import { Moon, Sunny } from '@element-plus/icons-vue';

const router = useRouter();
const userStore = useUserStore();
const registerFormRef = ref(null);
const loading = ref(false);
const error = ref('');
const systemName = ref('文档管理系统');
const systemRegisterMode = ref('invite');
const darkMode = computed(() => userStore.darkMode);

const toggleDarkMode = () => {
  userStore.toggleDarkMode();
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
    { required: true, message: '请输入用户名', trigger: 'blur' }
  ],
  email: [
    { required: true, message: '请输入邮箱', trigger: 'blur' },
    { type: 'email', message: '请输入正确的邮箱格式', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' }
  ],
  confirmPassword: [
    { required: true, message: '请确认密码', trigger: 'blur' },
    {
      validator: (rule, value, callback) => {
        if (value !== registerForm.password) {
          callback(new Error('两次输入的密码不一致'));
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
          callback(new Error('请输入邀请码'));
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
});
</script>

<style scoped>
.login-container {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: var(--bg-light);
  padding: var(--spacing-md);
  transition: all 0.3s ease;
}

.login-form-wrapper {
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

.theme-toggle-container {
  display: flex;
  justify-content: flex-end;
  margin-bottom: var(--spacing-md);
}

.theme-toggle {
  font-size: 18px;
  color: var(--text-medium);
  transition: all 0.3s ease;
}

.theme-toggle:hover {
  color: var(--primary-color);
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