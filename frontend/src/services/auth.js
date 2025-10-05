// /frontend/src/services/auth.js (已修复循环依赖)

import { ref, computed } from 'vue';
import { useRouter } from 'vue-router';
import api from './api';

// --- 模块级状态 (单例) ---
const token = ref(localStorage.getItem('token') || null);
const user = ref(JSON.parse(localStorage.getItem('user')) || null);
const isVisitorModeEnabled = ref(true);
const hasCheckedStatus = ref(false);

// ==================== 【核心修改】 ====================
// 单独导出 token，以便其他模块（如 api.js）可以直接访问而无需调用整个 useAuth()
export { token };
// ======================================================

// --- 纯粹的、不依赖组件上下文的函数 ---
const setAuthData = (newToken, newUser) => {
  token.value = newToken;
  user.value = newUser;
  localStorage.setItem('token', newToken);
  localStorage.setItem('user', JSON.stringify(newUser));
};

export const clearAuthData = () => {
  token.value = null;
  user.value = null;
  localStorage.removeItem('token');
  localStorage.removeItem('user');
};

// --- Vue 组合式函数 (只能在组件 setup 中使用) ---
export function useAuth() {
  const router = useRouter(); // 在这里调用是安全的

  const isLoggedIn = computed(() => !!token.value);

  const login = async (email, code) => {
    const response = await api.post('/auth/login', { email, code });
    setAuthData(response.data.token, response.data.user);
  };

  const logout = () => {
    clearAuthData();
    router.push('/login');
  };
  
  const fetchSystemStatus = async () => {
    try {
        const response = await api.get('/settings');
        isVisitorModeEnabled.value = response.data.visitor_mode_enabled;
    } catch (error) {
        console.error("无法获取系统设置:", error);
        isVisitorModeEnabled.value = false; 
    } finally {
        hasCheckedStatus.value = true;
    }
  };

  return {
    // token 不再需要从这里导出给 api.js，但组件可能仍需要它
    token, 
    user,
    isLoggedIn,
    isVisitorModeEnabled,
    hasCheckedStatus,
    login,
    logout,
    fetchSystemStatus,
  };
}