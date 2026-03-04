import api from './api';

// 统一处理响应
function handleResponse(response) {
  if (response.data && response.data.code === 20000) {
    return response.data
  } else {
    const errorMessage = response.data?.message || '请求失败'
    throw new Error(errorMessage)
  }
}

// 登录
export const login = async (email, password) => {
  const response = await api.post('/api/login', {
    email,
    password
  });
  return handleResponse(response);
};

// 注册
export const register = async (username, password, email, invitationCode) => {
  const response = await api.post('/api/register', {
    username,
    password,
    email,
    invitation_code: invitationCode
  });
  return handleResponse(response);
};

// 获取系统信息
export const getSystemInfo = async () => {
  const response = await api.get('/api/system-info');
  return handleResponse(response);
};