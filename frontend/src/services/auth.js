import api from './api';

// 登录
export const login = async (username, password) => {
  const response = await api.post('/api/login', {
    username,
    password
  });
  return response.data;
};

// 获取系统信息
export const getSystemInfo = async () => {
  const response = await api.get('/api/system-info');
  return response.data;
};