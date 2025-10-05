<!-- /frontend/src/App.vue (已添加移动端广告隐藏规则) -->
<template>
  <div id="app">
    <CyberBackground />
    
    <AdManager />
    <BackgroundTextAds /> 
    <BottomBannerAd />
    <NavBar />
    <router-view v-slot="{ Component }">
      <keep-alive include="WritingsPage,VideoGalleryPage">
        <component :is="Component" />
      </keep-alive>
    </router-view>
  </div>
</template>

<script setup>
import { onMounted, watch } from 'vue';
import NavBar from './components/NavBar.vue';
import AdManager from './components/AdManager.vue';
import BackgroundTextAds from './components/BackgroundTextAds.vue';
import BottomBannerAd from './components/BottomBannerAd.vue';
import CyberBackground from './components/CyberBackground.vue';
import { useAuth } from './services/auth';
import { useAds } from './store/ads';

const { fetchSystemStatus } = useAuth();
const { adsEnabled } = useAds();

onMounted(() => {
  fetchSystemStatus();
});

watch(adsEnabled, (isNowEnabled) => {
  const appElement = document.getElementById('app');
  if (appElement) {
    appElement.classList.toggle('ads-active', isNowEnabled);
  }
}, { immediate: true });
</script>

<style>
/* ==================== 核心修改：添加/强化这里的媒体查询规则 ==================== */

/* 
  当屏幕宽度小于 1600px 时 (通常是笔记本电脑)，
  隐藏屏幕两侧的文字广告，因为空间不足。
*/
@media (max-width: 1600px) {
  .background-ads-container {
    display: none !important;
  }
}

/* 
  当屏幕宽度小于 768px 时 (通常是手机)，
  隐藏底部的横幅广告。
  注意：两侧的文字广告因为上面的规则已经隐藏了，所以这里无需重复。
*/
@media (max-width: 768px) {
  .bottom-banner-ad {
    display: none !important;
  }
}
/* ============================================================================ */
</style>