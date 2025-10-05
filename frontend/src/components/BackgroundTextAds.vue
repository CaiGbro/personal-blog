<!-- /frontend/src/components/BackgroundTextAds.vue (最终版) -->
<template>
  <div class="background-ads-container" v-if="adsEnabled">
    <!-- 1. 左侧列 -->
    <div class="ad-column left-column">
      <!-- 文字广告卡片循环 -->
      <SloganCard 
        v-for="(card, i) in leftAdCardList" 
        :key="'left-' + i" 
        v-bind="card.props"
      />
      <!-- 2. 在卡片循环结束后，添加新的浮动图片广告 -->
      <FloatingImageAd src="/static/ads/corner/左.jpeg" />
    </div>

    <!-- 3. 右侧列 -->
    <div class="ad-column right-column">
      <!-- 文字广告卡片循环 -->
      <SloganCard 
        v-for="(card, i) in rightAdCardList" 
        :key="'right-' + i" 
        v-bind="card.props"
      />
      <!-- 4. 添加右侧的浮动图片广告 -->
      <FloatingImageAd src="/static/ads/corner/右.jpeg" />
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import SloganCard from './SloganCard.vue';
// 5. 导入新的图片广告组件
import FloatingImageAd from './FloatingImageAd.vue'; 
import { useAds } from '../store/ads';

const { adsEnabled } = useAds();
const leftAdCardList = ref([]);
const rightAdCardList = ref([]);

const slogansPool = [
  "探索未知", "代码艺术", "创新无界", "技术驱动", "构建未来",
  "灵感迸发", "极客精神", "持续学习", "拥抱变化", "追求卓越",
  "化繁为简", "数据智能", "云端漫步", "开源力量", "安全第一",
  "用户至上", "设计美学", "高效协作", "敏捷开发", "深度思考",
  "算法之舞", "架构之美", "智能纪元", "万物互联", "数字孪生"
];

const randomBetween = (min, max) => Math.random() * (max - min) + min;

const generateRandomStyle = () => {
  const animations = ['float-anim', 'pulse-bg', 'gentle-sway'];
  const backgrounds = [
    'linear-gradient(45deg, #1a1a1a, #333333)',
    'linear-gradient(45deg, #2c3e50, #34495e)',
    'linear-gradient(45deg, #232526, #414345)'
  ];
  const textColors = ['#e0e0e0', '#f1c40f', '#3498db', '#e74c3c', '#9b59b6'];
  const glowColors = {
    '#f1c40f': 'rgba(241, 196, 15, 0.6)',
    '#3498db': 'rgba(52, 152, 219, 0.6)',
    '#e74c3c': 'rgba(231, 76, 60, 0.6)',
    '#9b59b6': 'rgba(155, 89, 182, 0.6)',
    '#e0e0e0': 'rgba(224, 224, 224, 0.5)',
  };
  const selectedTextColor = textColors[Math.floor(Math.random() * textColors.length)];
  const selectedGlow = glowColors[selectedTextColor];

  return {
    cardStyle: {
      background: backgrounds[Math.floor(Math.random() * backgrounds.length)],
      border: `1px solid ${selectedGlow}`,
      boxShadow: `0 0 15px ${selectedGlow}`,
    },
    textStyle: { color: selectedTextColor, textShadow: `0 0 10px ${selectedGlow}` },
    animationName: animations[Math.floor(Math.random() * animations.length)],
    animationDuration: `${randomBetween(6, 12).toFixed(1)}s`,
    sloganChangeInterval: randomBetween(3000, 5000),
  };
};

const generateRandomAdCards = (count) => {
  const cards = [];
  const usedSlogans = new Set();

  for (let i = 0; i < count; i++) {
    const cardSlogans = new Set();
    // 6. 将每个卡片的标语数量增加到 5，以匀入多出来的内容
    while (cardSlogans.size < 5) { 
      const randomSlogan = slogansPool[Math.floor(Math.random() * slogansPool.length)];
      if (!usedSlogans.has(randomSlogan)) {
        cardSlogans.add(randomSlogan);
        usedSlogans.add(randomSlogan);
      }
    }
    cards.push({
      props: {
        slogans: Array.from(cardSlogans),
        ...generateRandomStyle(),
      }
    });
  }
  return cards;
};

onMounted(() => {
  // 7. 将卡片数量固定为 3
  const count = 3; 
  leftAdCardList.value = generateRandomAdCards(count);
  rightAdCardList.value = generateRandomAdCards(count);
});
</script>

<style scoped>
/* 样式保持不变 */
.background-ads-container {
  position: fixed; top: 0; left: 0; width: 100vw; height: 100vh;
  z-index: 9990; pointer-events: none; overflow: hidden;
}
.ad-column {
  position: absolute; top: 80px; bottom: 20px;
  width: calc((100vw - 1024px) / 2 - 4rem);
  display: flex; flex-direction: column; align-items: center;
  gap: 1.5rem; padding: 1rem 0; overflow-y: auto;
}
.ad-column::-webkit-scrollbar { display: none; }
.left-column { left: 2rem; }
.right-column { right: 2rem; }
@media (max-width: 1600px) {
  .background-ads-container { display: none; }
}
</style>