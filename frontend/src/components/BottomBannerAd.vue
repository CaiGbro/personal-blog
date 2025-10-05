<!-- /frontend/src/components/BottomBannerAd.vue (高度削减版) -->
<template>
  <transition name="slide-up">
    <div v-if="adsEnabled" class="bottom-banner-ad">
      <div class="marquee-track">
        <p class="marquee-content">
          <span>{{ adText }}</span>
          <span>{{ adText }}</span>
        </p>
      </div>
    </div>
  </transition>
</template>

<script setup>
import { ref } from 'vue';
import { useAds } from '../store/ads';

const { adsEnabled } = useAds();

const adText = ref(
  "【警告：未来已来】 >> 您的数字资产正在贬值！立即拥抱 AGI 浪潮，链接全球顶级算力网络，释放无限潜能。 >> 立即咨询，抢占下一个万亿级风口！ >> "
);
</script>

<style scoped>
.bottom-banner-ad {
  position: fixed;
  bottom: 0;
  left: 0;
  width: 100%;
  /* 1. 高度从 100px 削减约 1/3，调整为 70px */
  height: 70px; 
  z-index: 9995;
  display: flex;
  align-items: center;
  overflow: hidden;
  color: #fff;
  
  background: linear-gradient(125deg, #0f0c29, #302b63, #24243e, #e74c3c, #3498db);
  background-size: 400% 400%;
  animation: 
    aurora-bg 15s ease infinite,
    pulse-glow-border 4s infinite;

  box-shadow: 0 0 15px rgba(52, 152, 219, 0.5), inset 0 0 10px rgba(0, 0, 0, 0.5);
  border-top: 2px solid rgba(255, 255, 255, 0.3);
}

.bottom-banner-ad::before {
  content: '';
  position: absolute;
  top: 0;
  left: -150%;
  width: 60%;
  height: 100%;
  background: linear-gradient(
    to right,
    rgba(255, 255, 255, 0) 0%,
    rgba(0, 255, 255, 0.4) 50%,
    rgba(255, 255, 255, 0) 100%
  );
  transform: skewX(-25deg);
  animation: shimmer-sweep-fast 3s infinite linear;
  z-index: 1;
}

.bottom-banner-ad::after {
  content: '';
  position: absolute;
  top: 0;
  left: -150%;
  width: 40%;
  height: 100%;
  background: linear-gradient(
    to right,
    rgba(255, 255, 255, 0) 0%,
    rgba(231, 76, 60, 0.5) 50%,
    rgba(255, 255, 255, 0) 100%
  );
  transform: skewX(-25deg);
  animation: shimmer-sweep-slow 5s infinite linear;
  animation-delay: 1.5s;
  z-index: 1;
}

.marquee-track {
  width: 100%;
  white-space: nowrap;
  will-change: transform;
  position: relative;
  z-index: 2;
}

.marquee-content {
  display: inline-block;
  margin: 0;
  animation: marquee 30s linear infinite;
}

.marquee-content span {
  /* 2. 字体大小从 2rem 调整为 1.6rem，以适应新高度 */
  font-size: 1.6rem; 
  font-weight: 700;
  padding-right: 5rem;
  background: linear-gradient(90deg, #f1c40f, #e74c3c, #9b59b6, #3498db, #f1c40f);
  background-size: 200% 200%;
  background-clip: text;
  -webkit-background-clip: text;
  color: transparent;
  animation: rainbow-text 5s ease infinite;
  text-shadow: 1px 1px 2px rgba(0, 0, 0, 0.5);
}

/* --- 动画定义区 (保持不变) --- */
@keyframes aurora-bg { 0%{background-position:0% 50%} 50%{background-position:100% 50%} 100%{background-position:0% 50%} }
@keyframes pulse-glow-border { 0%, 100% { box-shadow: 0 0 15px rgba(52, 152, 219, 0.5), inset 0 0 10px rgba(0, 0, 0, 0.5); } 50% { box-shadow: 0 0 30px rgba(231, 76, 60, 0.8), inset 0 0 14px rgba(0, 0, 0, 0.5); } }
@keyframes shimmer-sweep-fast { 100% { left: 150%; } }
@keyframes shimmer-sweep-slow { 100% { left: 150%; } }
@keyframes rainbow-text { 0%{background-position:0% 50%} 50%{background-position:100% 50%} 100%{background-position:0% 50%} }
@keyframes marquee { from { transform: translateX(0%); } to { transform: translateX(-50%); } }
.slide-up-enter-active, .slide-up-leave-active { transition: transform 0.4s ease-out; }
.slide-up-enter-from, .slide-up-leave-to { transform: translateY(100%); }
</style>