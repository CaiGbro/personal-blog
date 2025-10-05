<!-- /frontend/src/components/FloatingImageAd.vue (究极·花里胡哨版) -->
<template>
  <div class="floating-image-ad">
    <img :src="src" alt="Cool Ad Image" class="ad-image" />
  </div>
</template>

<script setup>
defineProps({
  src: {
    type: String,
    required: true,
  },
});
</script>

<style scoped>
.floating-image-ad {
  width: 90%;
  margin-top: 1.5rem;
  border-radius: 12px;
  overflow: hidden;
  position: relative;
  background-color: #1a1a1a;
  padding: 8px;
  
  /* 1. 基础3D浮动效果 */
  transform: perspective(1000px) rotateY(-5deg);
  transition: transform 0.4s cubic-bezier(0.25, 0.46, 0.45, 0.94);
  will-change: transform, box-shadow;
  
  /* 2. 默认辉光边框 */
  box-shadow: 0 0 15px rgba(0, 255, 255, 0.4), 0 0 20px rgba(155, 89, 182, 0.3);
}

/* 3. 鼠标悬停时，效果增强且动画暂停，提供清晰的交互反馈 */
.floating-image-ad:hover {
  transform: perspective(1000px) rotateY(0deg) scale(1.08);
  box-shadow: 0 0 35px rgba(0, 255, 255, 0.8), 0 0 50px rgba(155, 89, 182, 0.6);
}
.floating-image-ad:hover .ad-image,
.floating-image-ad:hover::before,
.floating-image-ad:hover::after {
  animation-play-state: paused; /* 暂停所有子动画 */
}

.ad-image {
  width: 100%;
  height: auto;
  display: block;
  border-radius: 6px;
  opacity: 0.9;
  
  /* 4. 新增：色彩循环动画 (Hue-Rotate) */
  animation: hue-cycle 10s infinite linear;
}

/* 5. 新增：赛博朋克故障艺术 (Glitch) 效果 - 第一层 */
.floating-image-ad::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: inherit; /* 继承父元素的背景色 */
  mix-blend-mode: screen; /* 混合模式，产生颜色叠加效果 */
  opacity: 0.8;
  animation: glitch-effect-1 3s infinite linear alternate-reverse;
}

/* 6. 新增：赛博朋克故障艺术 (Glitch) 效果 - 第二层 */
.floating-image-ad::after {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: #e74c3c; /* 故障颜色 */
  mix-blend-mode: screen;
  opacity: 0.8;
  animation: glitch-effect-2 2s infinite linear alternate-reverse;
}

/* --- 动画定义区 --- */

/* 色彩循环动画 */
@keyframes hue-cycle {
  from { filter: hue-rotate(0deg); }
  to { filter: hue-rotate(360deg); }
}

/* 故障动画 1 (错位和切割) */
@keyframes glitch-effect-1 {
  0%, 100% { transform: translate(0); clip-path: inset(0); }
  10% { transform: translate(-2px, 2px); clip-path: inset(20% 0 60% 0); }
  20% { transform: translate(2px, -2px); clip-path: inset(80% 0 5% 0); }
  30% { transform: translate(-2px, -2px); clip-path: inset(40% 0 40% 0); }
  40% { transform: translate(2px, 2px); clip-path: inset(90% 0 2% 0); }
  50% { transform: translate(0); clip-path: inset(50% 0 30% 0); }
  60% { transform: translate(2px, -2px); clip-path: inset(10% 0 70% 0); }
  70% { transform: translate(-2px, 2px); clip-path: inset(75% 0 10% 0); }
  80% { transform: translate(0); clip-path: inset(0); }
}

/* 故障动画 2 (快速闪烁) */
@keyframes glitch-effect-2 {
  0%, 2%, 4%, 6%, 100% { opacity: 0; transform: scaleX(0); }
  1%, 3%, 5% { opacity: 0.6; transform: scaleX(1); }
}
</style>