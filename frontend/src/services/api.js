// /frontend/src/services/api.js (已修复循环依赖)

import axios from 'axios';
// ==================== 【核心修改】 ====================
// 不再导入 useAuth，而是直接导入 token 和 clearAuthData
import { token, clearAuthData } from './auth';
// ======================================================

const api = axios.create({
  baseURL: '/api',
});

api.interceptors.request.use(
  (config) => {
    // 直接使用导入的 token ref
    if (token.value) {
      config.headers['Authorization'] = `Bearer ${token.value}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

api.interceptors.response.use(
  (response) => response,
  (error) => {
    const status = error.response?.status;

    if (status === 401 || status === 403) {
      if (window.location.pathname !== '/login') {
        clearAuthData(); 
        window.location.href = '/login?mode=maintenance';
      }
    }
    
    return Promise.reject(error);
  }
);

export default api;