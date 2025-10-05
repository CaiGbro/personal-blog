<!-- /frontend/src/components/SloganCard.vue (已重构为深色流光主题) -->
<template>
  <div class="slogan-card" :style="[cardStyle, animationStyle]">
    <transition name="fade" mode="out-in">
      <span :key="currentSlogan" class="slogan-text" :style="textStyle">{{ currentSlogan }}</span>
    </transition>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, computed } from 'vue';

const props = defineProps({
  slogans: { type: Array, required: true, default: () => ['默认标语'] },
  cardStyle: { type: Object, default: () => ({}) },
  textStyle: { type: Object, default: () => ({}) },
  animationName: { type: String, default: 'float-anim' },
  animationDuration: { type: String, default: '8s' },
  sloganChangeInterval: { type: Number, default: 3500 },
});

const currentSlogan = ref(props.slogans[0]);
let sloganIndex = 0;
let intervalId = null;

const animationStyle = computed(() => ({
  animation: `${props.animationName} ${props.animationDuration} ease-in-out infinite`,
}));

onMounted(() => {
  intervalId = setInterval(() => {
    sloganIndex = (sloganIndex + 1) % props.slogans.length;
    currentSlogan.value = props.slogans[sloganIndex];
  }, props.sloganChangeInterval);
});

onUnmounted(() => {
  if (intervalId) clearInterval(intervalId);
});
</script>

<style scoped>
.slogan-card {
  width: 90%;
  padding: 1.5rem 1rem;
  border-radius: 12px;
  text-align: center;
  overflow: hidden; /* 关键：隐藏流光溢出的部分 */
  position: relative; /* 关键：为流光伪元素提供定位上下文 */

  /* --- 1. 默认样式改为深色主题 --- */
  background: linear-gradient(45deg, #1a1a1a, #333333);
  border: 1px solid #444;
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.3);
}

/* --- 2. 新增：流光效果伪元素 --- */
.slogan-card::before {
  content: '';
  position: absolute;
  top: 0;
  left: -150%; /* 从左侧屏幕外开始 */
  width: 80%;
  height: 100%;
  background: linear-gradient(
    to right,
    rgba(255, 255, 255, 0) 0%,
    rgba(255, 255, 255, 0.2) 50%, /* 流光亮度 */
    rgba(255, 255, 255, 0) 100%
  );
  transform: skewX(-25deg); /* 倾斜效果 */
  animation: shimmer-sweep 4s infinite linear; /* 应用动画 */
  z-index: 1;
}

.slogan-text {
  display: inline-block;
  position: relative; /* 确保文字在流光之上 */
  z-index: 2;

  /* --- 3. 默认文字样式改为亮色 --- */
  font-size: 1.1rem;
  font-weight: 500;
  color: #e0e0e0;
  text-shadow: 0 0 5px rgba(255, 255, 255, 0.5);
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.5s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

/* --- 4. 新增：流光动画 --- */
@keyframes shimmer-sweep {
  100% {
    left: 150%; /* 移动到右侧屏幕外结束 */
  }
}

/* 现有的其他动画效果保持不变 */
@keyframes float-anim {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-8px); }
}
@keyframes pulse-bg {
  0%, 100% { box-shadow: 0 4px 15px rgba(0, 0, 0, 0.3); }
  50% { box-shadow: 0 8px 30px rgba(52, 152, 219, 0.4); } /* 脉冲光效颜色 */
}
@keyframes gentle-sway {
  0%, 100% { transform: rotate(-1deg); }
  50% { transform: rotate(1deg); }
}
</style>