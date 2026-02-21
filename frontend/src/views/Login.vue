<template>
  <div class="login-container">
    <div class="login-form-wrapper">
      <div class="login-header">
        <h2>{{ systemName }}</h2>
        <p>请登录您的账户</p>
      </div>
      <el-form :model="loginForm" :rules="loginRules" ref="loginFormRef" class="login-form">
        <el-form-item prop="email">
          <el-input
            v-model="loginForm.email"
            placeholder="邮箱"
            prefix-icon="Message"
            :disabled="loading"
          />
        </el-form-item>
        <el-form-item prop="password">
          <el-input
            v-model="loginForm.password"
            type="password"
            placeholder="密码"
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
            登录
          </el-button>
        </el-form-item>
        <el-form-item class="register-link">
          <router-link to="/register">没有账户？立即注册</router-link>
        </el-form-item>
      </el-form>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import { useUserStore } from '../store/user';

const router = useRouter();
const route = useRoute();
const userStore = useUserStore();
const loginFormRef = ref(null);
const loading = ref(false);
const error = ref('');
const systemName = ref('Yoresee');

const loginForm = reactive({
  email: route.query.email || '',
  password: ''
});

const loginRules = {
  email: [
    { required: true, message: '请输入邮箱', trigger: 'blur' },
    { type: 'email', message: '请输入正确的邮箱格式', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' }
  ]
};

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
          error.value = userStore.error;
        }
      } catch (err) {
        error.value = '登录失败，请稍后重试';
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
  background-color: #F7F8FA;
  padding: 20px;
}

.login-form-wrapper {
  width: 100%;
  max-width: 400px;
  background-color: #FFFFFF;
  border-radius: 12px;
  box-shadow: 0px 4px 16px rgba(0, 0, 0, 0.1);
  padding: 40px;
  transition: all 0.3s ease;
}

.login-form-wrapper:hover {
  box-shadow: 0px 6px 20px rgba(0, 0, 0, 0.12);
}

.login-header {
  text-align: center;
  margin-bottom: 32px;
}

.login-header h2 {
  font-size: 24px;
  font-weight: 600;
  color: #1D2129;
  margin-bottom: 8px;
}

.login-header p {
  font-size: 14px;
  color: #86909C;
  margin: 0;
}

.login-form {
  width: 100%;
}

.login-form .el-form-item {
  margin-bottom: 20px;
}

.login-form .el-input {
  height: 40px;
  border-radius: 8px;
}

.login-form .el-input__wrapper {
  box-shadow: none;
}

.login-form .el-input__wrapper.is-focus {
  box-shadow: 0 0 0 2px rgba(22, 93, 255, 0.2);
}

.login-button {
  width: 100%;
  height: 40px;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 500;
  background-color: #165DFF;
  border-color: #165DFF;
}

.login-button:hover {
  background-color: #4080FF;
  border-color: #4080FF;
}

.login-button:active {
  background-color: #0E42D2;
  border-color: #0E42D2;
}

.error-message {
  margin-bottom: 16px;
}

.error-alert {
  padding: 8px 12px;
  font-size: 12px;
}

.register-link {
  text-align: center;
  margin-top: 16px;
}

.register-link a {
  color: #165DFF;
  font-size: 14px;
  text-decoration: none;
}

.register-link a:hover {
  text-decoration: underline;
}

@media (max-width: 768px) {
  .login-form-wrapper {
    padding: 32px 24px;
  }
  
  .login-header h2 {
    font-size: 20px;
  }
}
</style>