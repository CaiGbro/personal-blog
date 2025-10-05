// /frontend/src/store/ads.js
import { ref } from 'vue';
import api from '../services/api';

// --- 模块级状态 (单例) ---
const adsEnabled = ref(false);
const adList = ref([]);
const hasFetched = ref(false); // 防止重复获取

// --- 导出的组合式函数 ---
export function useAds() {

  const fetchAds = async () => {
    if (hasFetched.value) return; // 如果已经获取过，则不再执行
    try {
      const response = await api.get('/ads');
      adList.value = response.data || [];
      hasFetched.value = true;
    } catch (error) {
      console.error("获取广告列表失败:", error);
    }
  };

  const enableAds = () => {
    // 只有在广告列表为空时才去获取
    if (!hasFetched.value) {
      fetchAds();
    }
    adsEnabled.value = true;
  };

  const disableAds = () => {
    adsEnabled.value = false;
  };

  return {
    adsEnabled,
    adList,
    enableAds,
    disableAds,
  };
}