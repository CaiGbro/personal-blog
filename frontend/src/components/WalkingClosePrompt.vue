<!-- /frontend/src/components/WalkingClosePrompt.vue (已适配移动端尺寸) -->
<template>
  <div
    class="walking-prompt"
    ref="elementRef"
    :style="style"
    @mouseenter="pause"
    @mouseleave="resume"
  >
    <h3>温馨提示</h3>
    <p>“牛皮癣”广告已开启！</p>
    <button @click="disableAds" class="close-all-btn">
      一键关闭
    </button>
  </div>
</template>

<script setup>
import { ref } from 'vue';
import { useWalkingAnimation } from '../composables/useWalkingAnimation';
import { useAds } from '../store/ads';

const { disableAds } = useAds();
const elementRef = ref(null);
const { style, pause, resume } = useWalkingAnimation(elementRef);
</script>

<style scoped>
.walking-prompt {
  position: fixed;
  top: 0;
  left: 0;
  background-color: white;
  /* 桌面端的默认内边距 */
  padding: 1.5rem 2rem;
  border-radius: 12px;
  text-align: center;
  box-shadow: 0 10px 40px rgba(255, 77, 79, 0.5);
  z-index: 9999; /* 最高层级 */
  will-change: transform;
  cursor: pointer;
  border: 3px solid #ff4d4f;
}
.walking-prompt h3 {
  margin-top: 0;
  /* 桌面端的默认字体大小 */
  font-size: 1.5rem;
  color: #333;
}
.walking-prompt p {
  color: #666;
  margin-bottom: 1.5rem;
}
.close-all-btn {
  background-color: #ff4d4f;
  color: white;
  border: none;
  /* 桌面端的默认内边距和字体大小 */
  padding: 0.8rem 2rem;
  font-size: 1rem;
  font-weight: bold;
  border-radius: 8px;
  cursor: pointer;
  transition: background-color 0.2s, transform 0.2s;
}
.close-all-btn:hover {
  background-color: #d9363e;
  transform: scale(1.05);
}

/* ==================== 新增：移动端适配 ==================== */
/* 当屏幕宽度小于或等于 768px 时 (通常是手机) */
@media (max-width: 768px) {
  .walking-prompt {
    /* 减小内边距，让弹窗整体变小 */
    padding: 1rem 1.2rem;
    border-width: 2px;
    box-shadow: 0 5px 20px rgba(255, 77, 79, 0.5);
  }
  .walking-prompt h3 {
    /* 减小标题字体大小 */
    font-size: 1.2rem;
  }
  .walking-prompt p {
    /* 减小段落文字大小和下边距 */
    font-size: 0.9rem;
    margin-bottom: 1rem;
  }
  .close-all-btn {
    /* 减小按钮的内边距和字体大小 */
    padding: 0.6rem 1.5rem;
    font-size: 0.9rem;
  }
}
/* ====================================================== */
</style>