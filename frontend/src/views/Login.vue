<template>
  <div class="login-container">
    <div class="login-form-wrapper">
      <div class="login-header">
        <h1>企业文档管理系统</h1>
        <p>请登录您的账号</p>
      </div>
      <el-form :model="loginForm" :rules="loginRules" ref="loginFormRef" label-position="top">
        <el-form-item prop="username">
          <el-input 
            v-model="loginForm.username" 
            placeholder="请输入用户名" 
            prefix-icon="User"
            :disabled="loading"
          />
        </el-form-item>
        <el-form-item prop="password">
          <el-input 
            v-model="loginForm.password" 
            type="password" 
            placeholder="请输入密码" 
            prefix-icon="Lock"
            :disabled="loading"
            show-password
          />
        </el-form-item>
        <el-form-item v-if="error" class="error-message">
          <el-alert
            :title="error"
            type="error"
            show-icon
            :closable="false"
            effect="dark"
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
      </el-form>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue';
import { useRouter } from 'vue-router';
import { useUserStore } from '@/store/user';

const router = useRouter();
const userStore = useUserStore();
const loginFormRef = ref(null);

const loginForm = reactive({
  username: '',
  password: ''
});

const loginRules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' }
  ]
};

const loading = ref(false);
const error = ref('');

const handleLogin = async () => {
  if (!loginFormRef.value) return;
  
  const validateResult = await loginFormRef.value.validate();
  if (!validateResult) return;
  
  loading.value = true;
  error.value = '';
  
  try {
    const success = await userStore.login(loginForm.username, loginForm.password);
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
};
</script>

<style scoped>
.login-container {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 20px;
}

.login-form-wrapper {
  width: 100%;
  max-width: 400px;
  background: white;
  border-radius: 12px;
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.1);
  padding: 40px;
  animation: fadeIn 0.5s ease-in-out;
}

.login-header {
  text-align: center;
  margin-bottom: 30px;
}

.login-header h1 {
  font-size: 24px;
  font-weight: 600;
  color: #333;
  margin-bottom: 10px;
}

.login-header p {
  font-size: 14px;
  color: #666;
}

.login-button {
  width: 100%;
  height: 48px;
  font-size: 16px;
  margin-top: 10px;
}

.error-message {
  margin-bottom: 15px;
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(-20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

@media (max-width: 768px) {
  .login-form-wrapper {
    padding: 30px 20px;
  }
}
</style>