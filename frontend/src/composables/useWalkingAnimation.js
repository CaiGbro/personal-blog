// /frontend/src/composables/useWalkingAnimation.js
import { ref, onMounted, onUnmounted, computed } from 'vue';

export function useWalkingAnimation(elementRef) {
  const position = ref({ x: 0, y: 0 });
  const velocity = ref({ vx: 0, vy: 0 });
  let animationFrameId = null;

  // 将位置转换为 CSS transform 样式
  const style = computed(() => ({
    transform: `translate(${position.value.x}px, ${position.value.y}px)`,
  }));

  const getRandomVelocity = () => {
    const speed = 1.5; // 基础速度
    const angle = Math.random() * 2 * Math.PI;
    return {
      vx: Math.cos(angle) * speed,
      vy: Math.sin(angle) * speed,
    };
  };

  const animate = () => {
    if (!elementRef.value) return;

    const el = elementRef.value;
    let { x, y } = position.value;
    let { vx, vy } = velocity.value;

    x += vx;
    y += vy;

    if (x <= 0 || x + el.offsetWidth >= window.innerWidth) vx = -vx;
    if (y <= 0 || y + el.offsetHeight >= window.innerHeight) vy = -vy;

    x = Math.max(0, Math.min(x, window.innerWidth - el.offsetWidth));
    y = Math.max(0, Math.min(y, window.innerHeight - el.offsetHeight));

    position.value = { x, y };
    velocity.value = { vx, vy };

    animationFrameId = requestAnimationFrame(animate);
  };

  const pause = () => {
    if (animationFrameId) {
      cancelAnimationFrame(animationFrameId);
      animationFrameId = null;
    }
  };

  const resume = () => {
    if (!animationFrameId) {
      animationFrameId = requestAnimationFrame(animate);
    }
  };

  onMounted(() => {
    // 延迟一帧以确保元素尺寸已计算
    requestAnimationFrame(() => {
        if (elementRef.value) {
            position.value = {
                x: Math.random() * (window.innerWidth - elementRef.value.offsetWidth),
                y: Math.random() * (window.innerHeight - elementRef.value.offsetHeight),
            };
            velocity.value = getRandomVelocity();
            resume(); // 启动动画
        }
    });
  });

  onUnmounted(() => {
    pause(); // 组件销毁时停止动画
  });

  return { style, pause, resume };
}