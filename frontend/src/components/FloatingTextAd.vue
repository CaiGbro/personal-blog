<!-- /frontend/src/components/FloatingTextAd.vue -->
<template>
  <div class="floating-text-ad" :style="adStyle">
    {{ text }}
  </div>
</template>

<script setup>
import { computed } from 'vue';

const props = defineProps({
  text: String,
  initialStyle: Object,
});

// 将从父组件接收的随机样式和动画绑定到元素上
const adStyle = computed(() => ({
  ...props.initialStyle,
  animation: `float-anim ${props.initialStyle.animationDuration} ease-in-out infinite`,
}));
</script>

<style scoped>
.floating-text-ad {
  position: absolute;
  font-weight: bold;
  text-shadow: 1px 1px 2px rgba(0, 0, 0, 0.5);
  white-space: nowrap;
  user-select: none;
  pointer-events: none; /* 确保这些背景广告不会捕获鼠标事件 */
}

/* 定义一个通用的浮动动画 */
@keyframes float-anim {
  0%, 100% {
    transform: translateY(0) rotate(-2deg);
  }
  50% {
    transform: translateY(-15px) rotate(2deg);
  }
}
</style>