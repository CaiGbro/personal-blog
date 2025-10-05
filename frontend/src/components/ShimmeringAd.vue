<!-- /frontend/src/components/ShimmeringAd.vue (已修复) -->
<template>
  <!-- 1. 将 <a> 标签改为 <div>，移除链接功能，防止页面跳转 -->
  <div class="shimmering-ad">
    <span class="ad-text">{{ currentAdText }}</span>
  </div>
</template>

<script setup>
// script 部分无需改动
import { ref, onMounted, onUnmounted } from 'vue';

const adSlogans = [
  "前沿科技，触手可及",
  "代码改变世界",
  "加入我们，共创未来",
  "您的专属技术伙伴",
  "在这里，遇见思想的火花",
];

const currentAdText = ref(adSlogans[0]);
let intervalId = null;

onMounted(() => {
  let currentIndex = 0;
  intervalId = setInterval(() => {
    currentIndex = (currentIndex + 1) % adSlogans.length;
    currentAdText.value = adSlogans[currentIndex];
  }, 4000);
});

onUnmounted(() => {
  if (intervalId) {
    clearInterval(intervalId);
  }
});
</script>

<style scoped>
/* 2. 移除 text-decoration，并为 div 添加一个默认光标样式 */
.shimmering-ad {
  display: flex;
  justify-content: center;
  align-items: center;
  position: relative;
  border-radius: 12px;
  overflow: hidden;
  margin-top: 2rem;
  height: 120px;
  text-decoration: none; /* 确保没有下划线 */
  cursor: default; /* 使用默认光标，表明它不可点击 */
  background: linear-gradient(45deg, #1a1a1a, #333333);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.2);
  border: 1px solid #444;
}

/* 其他样式保持不变 */
.ad-text {
  font-size: 2rem;
  font-weight: bold;
  color: #fff;
  text-shadow: 0 0 5px #fff, 0 0 10px #fff, 0 0 15px #86b93a;
  position: relative;
  z-index: 2;
}

.shimmering-ad::before {
  content: '';
  position: absolute;
  top: 0;
  left: -150%;
  width: 80%;
  height: 100%;
  background: linear-gradient(
    to right,
    rgba(255, 255, 255, 0) 0%,
    rgba(255, 255, 255, 0.3) 50%,
    rgba(255, 255, 255, 0) 100%
  );
  transform: skewX(-25deg);
  animation: shimmer-sweep 3s infinite linear;
  z-index: 1;
}

@keyframes shimmer-sweep {
  100% {
    left: 150%;
  }
}

.shimmering-ad::after {
  content: '';
  position: absolute;
  top: 50%;
  left: 50%;
  width: 200px;
  height: 200px;
  background: radial-gradient(circle, rgba(134, 185, 58, 0.4) 0%, rgba(134, 185, 58, 0) 70%);
  transform: translate(-50%, -50%);
  animation: pulse-glow 5s infinite ease-in-out;
  z-index: 0;
}

@keyframes pulse-glow {
  0%, 100% {
    opacity: 0.5;
    transform: translate(-50%, -50%) scale(0.8);
  }
  50% {
    opacity: 1;
    transform: translate(-50%, -50%) scale(1.2);
  }
}
</style>