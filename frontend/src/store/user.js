import { defineStore } from 'pinia';
import { login, getSystemInfo } from '../services/auth';

export const useUserStore = defineStore('user', {
  state: () => ({
    token: localStorage.getItem('token') || '',
    userInfo: JSON.parse(localStorage.getItem('userInfo')) || null,
    systemName: '',
    loading: false,
    error: ''
  }),
  
  actions: {
    async login(username, password) {
      this.loading = true;
      this.error = '';
      try {
        const response = await login(username, password);
        this.token = response.token;
        this.userInfo = response.user;
        localStorage.setItem('token', response.token);
        localStorage.setItem('userInfo', JSON.stringify(response.user));
        return true;
      } catch (error) {
        this.error = error.response?.data?.message || '登录失败';
        return false;
      } finally {
        this.loading = false;
      }
    },
    
    logout() {
      this.token = '';
      this.userInfo = null;
      localStorage.removeItem('token');
      localStorage.removeItem('userInfo');
    },
    
    async fetchSystemInfo() {
      try {
        const response = await getSystemInfo();
        this.systemName = response.system_name;
        return response.system_name;
      } catch (error) {
        console.error('获取系统信息失败:', error);
        return '文档管理系统';
      }
    }
  }
});