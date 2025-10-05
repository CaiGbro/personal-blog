<!-- /frontend/src/components/WalkingAd.vue (已适配移动端尺寸) -->
<template>
  <div
    class="walking-ad"
    ref="elementRef"
    :style="style"
    @mouseenter="pause"
    @mouseleave="resume"
  >
    <img :src="adUrl" alt="Walking Ad" />
  </div>
</template>

<script setup>
import { ref } from 'vue';
import { useWalkingAnimation } from '../composables/useWalkingAnimation';

defineProps({
  adUrl: { type: String, required: true },
});

const elementRef = ref(null);
const { style, pause, resume } = useWalkingAnimation(elementRef);
</script>

<style scoped>
.walking-ad {
  position: fixed;
  top: 0;
  left: 0;
  /* 桌面端的默认宽度 */
  width: 280px;
  border: 3px solid #ffcc00;
  box-shadow: 0 0 20px rgba(255, 204, 0, 0.8);
  border-radius: 8px;
  z-index: 9998;
  overflow: hidden;
  will-change: transform;
  cursor: pointer;
}
.walking-ad img {
  width: 100%;
  height: auto;
  display: block;
}

/* ==================== 新增：移动端适配 ==================== */
/* 当屏幕宽度小于或等于 768px 时 (通常是手机) */
@media (max-width: 768px) {
  .walking-ad {
    /* 1. 将宽度改为基于视口宽度的相对单位，例如 45vw (视口宽度的45%) */
    width: 35vw;
    /* 2. 同时设置一个最大宽度，防止在某些宽屏手机上过大 */
    max-width: 180px;
    
    /* 3. (可选) 相应地减小边框和阴影，使其看起来更协调 */
    border-width: 2px;
    box-shadow: 0 0 10px rgba(255, 204, 0, 0.7);
  }
}
/* ====================================================== */
</style>