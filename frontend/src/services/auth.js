import api from './api';

// 登录
export const login = async (email, password) => {
  const response = await api.post('/api/login', {
    email,
    password
  });
  return response.data;
};

// 注册
export const register = async (username, password, email, invitationCode) => {
  const response = await api.post('/api/register', {
    username,
    password,
    email,
    invitation_code: invitationCode
  });
  return response.data;
};

// 获取系统信息
export const getSystemInfo = async () => {
  const response = await api.get('/api/system-info');
  return response.data;
};