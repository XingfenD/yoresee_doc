import { defineStore } from 'pinia';
import { login, register, getSystemInfo } from '../services/auth';

export const useUserStore = defineStore('user', {
  state: () => ({
    token: localStorage.getItem('token') || '',
    userInfo: (() => {
      try {
        const userInfo = localStorage.getItem('userInfo');
        return userInfo ? JSON.parse(userInfo) : null;
      } catch {
        return null;
      }
    })(),
    systemName: '',
    systemRegisterMode: 'open',
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
    
    async register(username, password, email, invitationCode) {
      this.loading = true;
      this.error = '';
      try {
        const response = await register(username, password, email, invitationCode);
        this.token = response.token;
        this.userInfo = response.user;
        localStorage.setItem('token', response.token);
        localStorage.setItem('userInfo', JSON.stringify(response.user));
        return true;
      } catch (error) {
        this.error = error.response?.data?.message || '注册失败';
        return false;
      } finally {
        this.loading = false;
      }
    },
    
    async fetchSystemInfo() {
      try {
        const response = await getSystemInfo();
        this.systemName = response.system_name;
        this.systemRegisterMode = response.system_register_mode || 'invite';
        return response;
      } catch (error) {
        console.error('获取系统信息失败:', error);
        return { system_name: 'Yoresee', system_register_mode: 'invite' };
      }
    }
  }
});