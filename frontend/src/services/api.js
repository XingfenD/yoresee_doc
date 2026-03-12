import axios from 'axios';
import { ElMessage } from 'element-plus';

// 创建axios实例
const api = axios.create({
  baseURL: '/',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json'
  }
});

// 请求拦截器
api.interceptors.request.use(
  config => {
    // 从localStorage获取token
    const token = localStorage.getItem('token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    // 从localStorage获取语言设置
    const language = localStorage.getItem('language') || 'en';
    config.headers['Accept-Language'] = language;
    return config;
  },
  error => {
    return Promise.reject(error);
  }
);

// 响应拦截器
api.interceptors.response.use(
  response => {
    return response;
  },
  error => {
    // 处理错误响应
    if (error.response) {
      // 服务器返回错误状态码
      switch (error.response.status) {
        case 401:
          // 未授权，清除token并跳转到登录页
          localStorage.removeItem('token');
          localStorage.removeItem('userInfo');
          window.location.href = '/login';
          break;
        case 403:
          // 禁止访问
          console.error('没有权限访问该资源');
          break;
        case 404:
          // 资源不存在
          console.error('请求的资源不存在');
          break;
        case 500:
          // 服务器错误
          console.error('服务器内部错误');
          break;
        default:
          console.error('请求失败:', error.response.data.message || '未知错误');
      }
    } else if (error.request) {
      // 请求已发送但没有收到响应
      console.error('网络错误，请检查网络连接');
    } else {
      // 请求配置错误
      console.error('请求配置错误:', error.message);
    }
    return Promise.reject(error);
  }
);

// 统一处理响应
function handleResponse(response) {
  if (response.data && response.data.code === 20000) {
    return response.data
  } else {
    const errorMessage = response.data?.message || '请求失败'
    ElMessage.error(errorMessage)
    throw new Error(errorMessage)
  }
}

// 获取知识库列表
export const getKnowledgeBases = (params) => {
  return api.get('/api/knowledge-bases', { params }).then(handleResponse)
};

// 获取知识库详情
export const getKnowledgeBaseDetail = (knowledgeBaseExternalID, params = {}) => {
  return api.get(`/api/knowledge-base/${knowledgeBaseExternalID}`, { params }).then(handleResponse)
};

// 获取知识库文档列表
export const getKnowledgeBaseDocuments = (knowledgeBaseExternalID, params = {}) => {
  return api.get(`/api/knowledge-base/${knowledgeBaseExternalID}`, { 
    params: {
      record_recent_log: false,
      page: 1,
      page_size: 1000,
      ...params
    }
  }).then(handleResponse)
};

// 创建文档
export const createDocument = (data) => {
  return api.post('/api/document/create', data).then(handleResponse)
};

export default api;