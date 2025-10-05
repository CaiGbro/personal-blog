<!-- /frontend/src/components/VideoCard.vue -->
<template>
  <router-link :to="playerUrl" class="video-card" @mouseover="playPreview" @mouseleave="pausePreview">
    <video
      ref="videoPreview"
      :src="src"
      :poster="posterUrl" 
      class="video-preview"
      muted
      loop
      playsinline
      preload="metadata"
    ></video>
    <div class="overlay">
      <span class="play-icon">▶</span>
    </div>
  </router-link>
</template>

<script setup>
import { ref, computed } from 'vue';

const props = defineProps({
  src: { type: String, required: true },
});

const videoPreview = ref(null);

// 新增：计算封面图的 URL
const posterUrl = computed(() => {
  // 假设视频 src 是 /static/video/my-video.mp4
  // 我们把它替换成 /static/video/my-video.jpg
  return props.src.replace(/\.mp4$/, '.jpg');
});

const playerUrl = computed(() => {
  return `/works/media-viewer?start_with=${encodeURIComponent(props.src)}`;
});

const playPreview = () => {
  // 在桌面端，我们仍然保留悬停预览功能
  if (videoPreview.value) {
    videoPreview.value.play().catch(e => console.log("预览播放失败:", e));
  }
};

const pausePreview = () => {
  if (videoPreview.value) {
    videoPreview.value.pause();
    videoPreview.value.currentTime = 0;
  }
};
</script>

<style scoped>
.video-card {
  position: relative;
  display: block;
  break-inside: avoid;
  margin-bottom: 1.5rem;
  border-radius: 8px;
  overflow: hidden;
  box-shadow: 0 4px 12px rgba(0,0,0,0.08);
  cursor: pointer;
}
.video-preview {
  width: 100%;
  height: 100%;
  object-fit: cover;
  display: block;
  transition: transform 0.3s ease;
}
.overlay {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: rgba(0, 0, 0, 0.2);
  display: flex;
  justify-content: center;
  align-items: center;
  opacity: 1;
  transition: opacity 0.3s ease;
}
.play-icon {
  color: white;
  font-size: 2rem;
  text-shadow: 0 2px 4px rgba(0,0,0,0.5);
}
.video-card:hover .overlay {
  opacity: 0;
}
.video-card:hover .video-preview {
  transform: scale(1.05);
}

@media (max-width: 768px) {
  /* 在手机上，默认隐藏播放图标遮罩层，让封面图更清晰 */
  .overlay {
    opacity: 0;
  }

  /* 移除手机上的悬停放大效果，因为它可能导致布局问题 */
  .video-card:hover .video-preview {
    transform: none;
  }
}

</style>