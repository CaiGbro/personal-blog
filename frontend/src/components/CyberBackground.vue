<!-- /frontend/src/components/CyberBackground.vue (新增文件) -->
<template>
  <transition name="fade">
    <div v-if="adsEnabled" class="cyber-background"></div>
  </transition>
</template>

<script setup>
import { useAds } from '../store/ads';
const { adsEnabled } = useAds();
</script>

<style scoped>
.cyber-background {
  position: fixed;
  top: 0;
  left: 0;
  width: 100vw;
  height: 100vh;
  /* z-index: -1 确保它在所有内容的后面 */
  z-index: -1; 
  overflow: hidden;

  /* 1. 深邃的星空/宇宙背景 */
  background-color: #0a0a14;
  background-image: radial-gradient(
    circle at 20% 20%,
    rgba(52, 152, 219, 0.2) 0%,
    transparent 30%
  ),
  radial-gradient(
    circle at 80% 70%,
    rgba(155, 89, 182, 0.2) 0%,
    transparent 30%
  );

  /* 2. 动态移动的赛博网格 */
  background-image: 
    repeating-linear-gradient(
      0deg,
      rgba(0, 255, 255, 0.1),
      rgba(0, 255, 255, 0.1) 1px,
      transparent 1px,
      transparent 50px
    ),
    repeating-linear-gradient(
      90deg,
      rgba(0, 255, 255, 0.1),
      rgba(0, 255, 255, 0.1) 1px,
      transparent 1px,
      transparent 50px
    );
  animation: move-grid 20s linear infinite;
}

/* 3. 叠加一层从上到下缓慢移动的扫描线，增加复古感 */
.cyber-background::after {
  content: '';
  position: absolute;
  top: -10%;
  left: 0;
  width: 100%;
  height: 4px;
  background: cyan;
  box-shadow: 0 0 15px cyan;
  opacity: 0.5;
  animation: scan-line 8s linear infinite;
}

/* --- 动画定义区 --- */

/* 网格移动动画 */
@keyframes move-grid {
  from {
    background-position: 0 0;
  }
  to {
    background-position: 50px 50px;
  }
}

/* 扫描线动画 */
@keyframes scan-line {
  from {
    top: -10%;
  }
  to {
    top: 110%;
  }
}

/* 背景板淡入淡出动画 */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.5s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>