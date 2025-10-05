<!-- /frontend/src/components/DisclaimerMatrix.vue (双模式版) -->
<template>
  <!-- 1. 动态绑定 'cool-mode-active' 类 -->
  <div class="disclaimer-matrix" :class="{ 'cool-mode-active': adsEnabled }">
    
    <!-- 2. Canvas 只在广告模式开启时渲染 -->
    <canvas v-if="adsEnabled" ref="canvasRef" class="matrix-canvas"></canvas>
    
    <div class="disclaimer-content">
      <h3>--[ 免 责 声 明 ]--</h3>
      <ul>
        <!-- 3. data-text 也只在广告模式下应用，以激活 Glitch 特效 -->
        <li v-for="(item, index) in disclaimers" :key="index" :data-text="adsEnabled ? item : null">
          {{ item }}
        </li>
      </ul>
    </div>
  </div>
</template>

<script setup>
import { ref, watchEffect } from 'vue';
// 4. 引入 useAds 来感知全局状态
import { useAds } from '../store/ads';

const { adsEnabled } = useAds();
const canvasRef = ref(null);
const disclaimers = ref([
  "本站所有内容仅用于学习交流目的。",
  "请勿用于非法用途，本站不承担任何由此产生的法律责任。",
  "部分内容来源于网络，若有侵权，请联系删除。",
  "所有观点均为个人见解，不构成任何投资或专业建议。",
  "本站保留最终解释权，并随时准备跃迁至下一个维度。",
]);

// 5. 使用 watchEffect 来管理 Canvas 的创建和销毁
watchEffect((onCleanup) => {
  // 只有在广告模式开启且 canvas 元素已挂载时才执行
  if (adsEnabled.value && canvasRef.value) {
    const canvas = canvasRef.value;
    const ctx = canvas.getContext('2d');
    
    let width = canvas.width = canvas.offsetWidth;
    let height = canvas.height = canvas.offsetHeight;

    const characters = 'ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789';
    const charArray = characters.split('');
    const fontSize = 14;
    const columns = Math.floor(width / fontSize);
    const drops = Array(columns).fill(1);

    let animationFrameId;

    function draw() {
      ctx.fillStyle = 'rgba(0, 0, 0, 0.05)';
      ctx.fillRect(0, 0, width, height);
      ctx.fillStyle = '#0F0';
      ctx.font = `${fontSize}px monospace`;
      for (let i = 0; i < drops.length; i++) {
        const text = charArray[Math.floor(Math.random() * charArray.length)];
        ctx.fillText(text, i * fontSize, drops[i] * fontSize);
        if (drops[i] * fontSize > height && Math.random() > 0.975) drops[i] = 0;
        drops[i]++;
      }
    }

    const animate = () => {
      draw();
      animationFrameId = requestAnimationFrame(animate);
    };
    animate();

    const resizeHandler = () => {
      width = canvas.width = canvas.offsetWidth;
      height = canvas.height = canvas.offsetHeight;
      // 重置 drops 数组以适应新尺寸
      drops.fill(1);
    };
    window.addEventListener('resize', resizeHandler);

    // onCleanup 函数会在 watchEffect 重新运行前或组件卸载时执行
    onCleanup(() => {
      cancelAnimationFrame(animationFrameId);
      window.removeEventListener('resize', resizeHandler);
    });
  }
});
</script>

<style scoped>
/* ================================================== */
/* =============== 1. 正常模式样式 =============== */
/* ================================================== */
.disclaimer-matrix {
  position: relative;
  background-color: #f9f9f9; /* 浅灰色背景 */
  border: 1px solid #eee;    /* 细边框 */
  border-radius: 8px;
  padding: 2rem;
  overflow: hidden;
  transition: all 0.5s ease; /* 添加平滑过渡效果 */
}

.disclaimer-content {
  position: relative;
  z-index: 2;
  color: #555; /* 深灰色文字 */
  font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif;
}

.disclaimer-content h3 {
  text-align: center;
  margin-bottom: 1.5rem;
  color: #333;
  font-weight: 600;
}

.disclaimer-content ul {
  list-style: none;
  padding: 0;
  margin: 0;
}

.disclaimer-content li {
  position: relative;
  padding-left: 1.5rem;
  margin-bottom: 1rem;
  line-height: 1.6;
}

.disclaimer-content li::before {
  content: '•'; /* 使用简单的圆点作为列表符号 */
  position: absolute;
  left: 0;
  top: 0;
  color: #86b93a; /* 主题绿色 */
  font-weight: bold;
}


/* ================================================== */
/* =============== 2. 炫酷模式样式 =============== */
/* ================================================== */

/* 当 .cool-mode-active 类存在时，以下样式将覆盖默认样式 */

.cool-mode-active.disclaimer-matrix {
  background-color: #000;
  border: 2px solid #0F0;
  box-shadow: 0 0 15px #0F0, inset 0 0 10px #0F0;
  animation: pulse-glow 4s infinite ease-in-out;
}

.matrix-canvas {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  z-index: 1;
  opacity: 0.7;
}

.cool-mode-active .disclaimer-content {
  color: #0F0;
  font-family: 'Courier New', Courier, monospace;
  text-shadow: 0 0 5px #0F0;
  backdrop-filter: blur(1.5px);
}

.cool-mode-active .disclaimer-content h3 {
  color: #0F0;
  letter-spacing: 3px;
}

.cool-mode-active .disclaimer-content li {
  animation: glitch-load 2s steps(4, end) infinite;
}

/* 炫酷模式下的列表符号和 Glitch 特效 */
.cool-mode-active .disclaimer-content li::before {
  content: '>';
  animation: blink-caret 0.8s infinite;
}

.cool-mode-active .disclaimer-content li::after,
.cool-mode-active .disclaimer-content li::before {
  content: attr(data-text);
  position: absolute;
  top: 0;
  left: 1.5rem;
  width: 100%;
  background: #000;
  overflow: hidden;
}

.cool-mode-active .disclaimer-content li::after {
  left: 1.7rem;
  text-shadow: -1px 0 red;
  animation: glitch-anim-1 2.5s infinite linear alternate-reverse;
}

.cool-mode-active .disclaimer-content li::before {
  left: 1.3rem;
  text-shadow: 1px 0 blue;
  animation: glitch-anim-2 3s infinite linear alternate-reverse;
}


/* --- 动画定义 (仅用于炫酷模式) --- */
@keyframes pulse-glow {
  0%, 100% { box-shadow: 0 0 15px #0F0, inset 0 0 10px #0F0; }
  50% { box-shadow: 0 0 30px #0F0, inset 0 0 20px #0F0; }
}
@keyframes blink-caret {
  from, to { opacity: 0; }
  50% { opacity: 1; }
}
@keyframes glitch-load {
  0% { opacity: 0; }
  20% { opacity: 1; }
  100% { opacity: 1; }
}
@keyframes glitch-anim-1 {
  0% { clip-path: inset(3% 0 95% 0); } 20% { clip-path: inset(80% 0 5% 0); } 40% { clip-path: inset(45% 0 47% 0); } 60% { clip-path: inset(92% 0 2% 0); } 80% { clip-path: inset(25% 0 68% 0); } 100% { clip-path: inset(58% 0 33% 0); }
}
@keyframes glitch-anim-2 {
  0% { clip-path: inset(78% 0 12% 0); } 20% { clip-path: inset(15% 0 70% 0); } 40% { clip-path: inset(63% 0 23% 0); } 60% { clip-path: inset(2% 0 88% 0); } 80% { clip-path: inset(90% 0 5% 0); } 100% { clip-path: inset(40% 0 45% 0); }
}
</style>