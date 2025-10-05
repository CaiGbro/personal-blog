// /frontend/src/store/layout.js

import { ref } from 'vue';

// 创建一个响应式的变量来控制导航栏的可见性，默认为 true (可见)
const isNavBarVisible = ref(true);

// 导出一个组合式函数，以便任何组件都可以访问和修改这个状态
export function useLayout() {

  const showNavBar = () => {
    isNavBarVisible.value = true;
  };

  const hideNavBar = () => {
    isNavBarVisible.value = false;
  };

  const toggleNavBar = () => {
    isNavBarVisible.value = !isNavBarVisible.value;
  };

  return {
    isNavBarVisible,
    showNavBar,
    hideNavBar,
    toggleNavBar,
  };
}