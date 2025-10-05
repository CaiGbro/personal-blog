<!-- /frontend/src/components/AdManager.vue (恢复默认行为) -->
<template>
  <div v-if="adsEnabled" class="ad-overlay-container">
    <!-- 
      核心修改：
      移除了 v-if="!isMobile"，让 WalkingClosePrompt 在所有设备上都显示。
    -->
    <WalkingClosePrompt />

    <!-- WalkingAd (GIF广告) 始终渲染 -->
    <WalkingAd
      v-for="(ad, index) in adList"
      :key="index"
      :ad-url="ad"
    />
  </div>
</template>

<script setup>
// --- 移除了 isMobile 相关的逻辑，恢复到最简状态 ---
import { useAds } from '../store/ads';
import WalkingAd from './WalkingAd.vue';
import WalkingClosePrompt from './WalkingClosePrompt.vue';

const { adsEnabled, adList } = useAds();
</script>