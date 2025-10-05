<template>
   <div 
     class="video-player-wrapper" 
     @wheel.prevent="handleWheel" 
     @touchstart="handleTouchStart"
     @touchmove="handleTouchMove"
     @touchend="handleTouchEnd"
     ref="playerWrapper" 
     tabindex="0"
   >
    <div class="video-container" v-if="videos.length > 0">
      <!-- 使用两个 video 元素实现无缝切换 -->
      <video
        ref="videoRefA"
        :class="['video-player', { active: activePlayer === 'A' }]"
        @play="onPlay"
        @pause="onPause"
        muted
        loop
        playsinline
      ></video>
      <video
        ref="videoRefB"
        :class="['video-player', { active: activePlayer === 'B' }]"
        @play="onPlay"
        @pause="onPause"
        muted
        loop
        playsinline
      ></video>
      <!-- click-overlay 的点击事件现在统一处理播放/暂停和导航栏切换 -->
      <div class="click-overlay" @click.stop="togglePlay">
        <div class="play-pause-icon-container" :class="{ 'is-visible': !isPlaying && !userWantsToPlay && !isSwitching }">
          <button class="play-pause-btn">▶</button>
        </div>
      </div>
      
      <div class="controls-overlay"></div>

      <div class="progress-indicator">
        {{ currentIndex + 1 }} / {{ videos.length }}
      </div>

      <div class="volume-control">
        <button @click.stop="toggleMute" class="volume-btn">
          {{ isMuted || volumeLevel === 0 ? '🔇' : '🔊' }}
        </button>
        <input
          type="range"
          min="0"
          max="1"
          step="0.01"
          :value="isMuted ? 0 : volumeLevel"
          class="volume-slider"
          @input="handleVolumeChange"
        />
      </div>

      <div class="side-controls">
        <button class="control-btn" @click.stop="addReaction('❤️')">
          <span class="icon">❤️</span>
          <span class="count">{{ reactionCount }}</span>
        </button>
        <button class="control-btn" @click.stop="showComments = true">
          <span class="icon">💬</span>
          <span class="count">{{ commentCount }}</span>
        </button>
        <button class="control-btn" @click.stop="downloadCurrentVideo">
          <span class="icon">📥</span>
          <span class="count">下载</span>
        </button>
      </div>

    </div>
    <div v-else class="loading">
      <p>加载视频中...</p>
    </div>

    <VideoCommentSection
      v-if="showComments && currentVideoSrc"
      :video-filename="currentVideoSrc"
      @close="showComments = false"
      @count-updated="updateCommentCount" 
    />
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted, watch } from 'vue';
import api from '../services/api';
import VideoCommentSection from './VideoCommentSection.vue';
import { useLayout } from '../store/layout'; 

// 1. 添加一个简单的函数来判断是否为移动设备 (通过检测触摸事件)
const isMobileDevice = () => 'ontouchstart' in window || navigator.maxTouchPoints > 0;

// 2. 根据设备类型来初始化 isMuted 的值
//    - 如果是移动设备，isMuted 默认为 false (不静音)
//    - 如果是桌面设备，isMuted 默认为 true (静音)
const isMuted = ref(isMobileDevice() ? false : true);

const props = defineProps({
  videos: { type: Array, required: true },
  startWithSrc: { type: String, required: true },
});

const playerWrapper = ref(null);
const videoRefA = ref(null);
const videoRefB = ref(null);
const activePlayer = ref('A');
const currentVideoElement = computed(() => activePlayer.value === 'A' ? videoRefA.value : videoRefB.value);
const currentIndex = ref(0);
const isPlaying = ref(false);
const userWantsToPlay = ref(true);
const isSwitching = ref(false);

const volumeLevel = ref(0.2);
const currentVideoSrc = computed(() => props.videos[currentIndex.value]);
const showComments = ref(false);
const commentCount = ref(0);

// --- 新增：用于处理触摸滑动的状态 ---
const touchStartY = ref(0);
const touchEndY = ref(0);
const touchMoveThreshold = 50; // 至少滑动50像素才算有效切换

const updateCommentCount = (newCount) => { 
  commentCount.value = newCount; 
};

const { hideNavBar, showNavBar, toggleNavBar } = useLayout();

const fetchCommentCount = async (filename) => {
  if (!filename) return;
  try {
    const response = await api.get('/video_comments', { params: { video_filename: filename } });
    commentCount.value = response.data.totalCount || 0;
  } catch (error) {
    console.error("获取评论数失败:", error);
    commentCount.value = 0;
  }
};

const reactions = ref({});
const reactionCount = computed(() => reactions.value['❤️'] || 0);

const fetchReactions = async (filename) => {
  if (!filename) return;
  try {
    const response = await api.get(`/videos/reactions${filename}`);
    reactions.value = response.data || {};
  } catch (error) {
    console.error("获取视频回应失败:", error);
    reactions.value = {};
  }
};

const addReaction = async (emoji) => {
  if (!currentVideoSrc.value) return;
  try {
    const response = await api.post(`/videos/react${currentVideoSrc.value}`, { emoji });
    reactions.value = response.data;
  } catch (error) {
    console.error("添加回应失败:", error);
  }
};

const downloadCurrentVideo = () => {
  if (!currentVideoSrc.value) return;
  const link = document.createElement('a');
  link.href = currentVideoSrc.value;
  link.download = currentVideoSrc.value.split('/').pop() || 'video.mp4';
  document.body.appendChild(link);
  link.click();
  document.body.removeChild(link);
};

const playVideo = async () => {
  const video = currentVideoElement.value;
  if (!video) return;
  try {
    await video.play();
  } catch (error) {
    console.error("Autoplay failed:", error);
    isPlaying.value = false;
  }
};

const pauseVideo = () => {
  const video = currentVideoElement.value;
  if (video) video.pause();
};

// ==================== 2. 修改 togglePlay 函数 ====================
const togglePlay = () => {
  // 每次点击视频区域时，都切换导航栏的显示状态
  toggleNavBar();

  const video = currentVideoElement.value;
  if (!video) return;

  // 切换视频的播放/暂停状态
  if (video.paused) {
    userWantsToPlay.value = true;
    playVideo();
  } else {
    userWantsToPlay.value = false;
    pauseVideo();
  }
};
// ================================================================

const onPlay = () => { 
  isPlaying.value = true; 
};
const onPause = () => { 
  if (!isSwitching.value) isPlaying.value = false; 
};

const changeVideo = (newIndex) => {
  if (newIndex < 0 || newIndex >= props.videos.length || newIndex === currentIndex.value || isSwitching.value) {
    return;
  }
  isSwitching.value = true;
  userWantsToPlay.value = true; // 假设用户希望新视频也播放

  const inactivePlayerRef = activePlayer.value === 'A' ? videoRefB : videoRefA;
  const activePlayerRef = activePlayer.value === 'A' ? videoRefA : videoRefB;

  const inactivePlayerEl = inactivePlayerRef.value;
  const activePlayerEl = activePlayerRef.value;

  if (!inactivePlayerEl || !activePlayerEl) {
    isSwitching.value = false;
    return;
  }

  inactivePlayerEl.src = props.videos[newIndex];
  inactivePlayerEl.load(); // 加载新视频

  const onCanPlay = () => {
    inactivePlayerEl.play().then(() => {
      activePlayer.value = (activePlayer.value === 'A' ? 'B' : 'A'); // 切换活跃播放器
      currentIndex.value = newIndex;
      activePlayerEl.pause(); // 暂停旧的活跃播放器
      // 延迟一小段时间，确保切换动画完成
      setTimeout(() => { isSwitching.value = false; }, 300); 
    }).catch(error => {
      console.error("New video play failed:", error);
      isSwitching.value = false;
    });
    inactivePlayerEl.removeEventListener('canplay', onCanPlay); // 移除监听器
  };

  inactivePlayerEl.addEventListener('canplay', onCanPlay, { once: true });
};

const nextVideo = () => {
  changeVideo(currentIndex.value + 1);
};

const prevVideo = () => {
  changeVideo(currentIndex.value - 1);
};

const toggleMute = () => {
  const video = currentVideoElement.value;
  if (!video) return;

  isMuted.value = !isMuted.value;
  // 同时更新两个视频元素的 muted 状态
  if(videoRefA.value) videoRefA.value.muted = isMuted.value;
  if(videoRefB.value) videoRefB.value.muted = isMuted.value;
};

const handleVolumeChange = (event) => {
  const newVolume = parseFloat(event.target.value);
  volumeLevel.value = newVolume;
  isMuted.value = newVolume === 0; // 如果音量为0，则视为静音

  // 同时更新两个视频元素的音量和静音状态
  if (videoRefA.value) {
    videoRefA.value.volume = newVolume;
    videoRefA.value.muted = isMuted.value;
  }
  if (videoRefB.value) {
    videoRefB.value.volume = newVolume;
    videoRefB.value.muted = isMuted.value;
  }
};

let lastScrollTime = 0;
const handleWheel = (event) => {
  const now = Date.now();
  // 简单的防抖，防止滚轮过快触发多次切换
  if (now - lastScrollTime < 500) return; 

  if (event.deltaY > 0) {
    nextVideo();
  } else if (event.deltaY < 0) {
    prevVideo();
  }
  lastScrollTime = now;
};

const handleKeyDown = (event) => {
  // 如果评论区打开，则不响应键盘事件
  if (showComments.value) return;

  if (event.key === 'ArrowDown') {
    event.preventDefault(); // 阻止页面滚动
    nextVideo();
  } else if (event.key === 'ArrowUp') {
    event.preventDefault(); // 阻止页面滚动
    prevVideo();
  } else if (event.key === ' ') {
    event.preventDefault(); // 阻止空格键的默认行为（通常是滚动页面）
    togglePlay();
  }
};

// --- 新增：触摸事件处理函数 ---
const handleTouchStart = (event) => {
  touchStartY.value = event.touches[0].clientY;
};

const handleTouchMove = (event) => {
  // 阻止默认的滚动行为，以实现自定义的视频切换滑动
  event.preventDefault(); 
  touchEndY.value = event.touches[0].clientY;
};

const handleTouchEnd = () => {
  if (touchStartY.value === 0 || touchEndY.value === 0) return; // 没有有效触摸

  const deltaY = touchEndY.value - touchStartY.value;

  if (deltaY < -touchMoveThreshold) {
    // 向上滑动，切换到下一个视频
    nextVideo();
  } else if (deltaY > touchMoveThreshold) {
    // 向下滑动，切换到上一个视频
    prevVideo();
  }

  // 重置触摸坐标
  touchStartY.value = 0;
  touchEndY.value = 0;
};

watch(currentVideoSrc, (newSrc) => {
  if (newSrc) {
    fetchCommentCount(newSrc);
    fetchReactions(newSrc);
  }
}, { immediate: true }); // 立即执行一次，以加载初始视频的评论和回应

onMounted(() => {
  hideNavBar(); // 进入视频播放页时隐藏导航栏
  window.addEventListener('keydown', handleKeyDown);
  playerWrapper.value?.focus(); // 确保 wrapper 元素获得焦点，以便响应键盘事件

  // 根据 startWithSrc 确定初始播放的视频
  const startIndex = props.videos.findIndex(v => v === props.startWithSrc);
  if (startIndex !== -1) {
    currentIndex.value = startIndex;
  }

  if (props.videos.length > 0) {
    const videoA = videoRefA.value;
    const videoB = videoRefB.value;

     if (videoA && videoB) {
      // --- 核心修改 2：在挂载时应用正确的静音状态 ---
      // 这里的 isMuted.value 已经根据设备类型被正确设置了
      videoA.volume = volumeLevel.value;
      videoA.muted = isMuted.value;
      videoB.volume = volumeLevel.value;
      videoB.muted = isMuted.value;
      // --- 修改结束 ---

      // 设置初始视频源并播放
      videoA.src = props.videos[currentIndex.value];
      videoA.load();
      playVideo();

      // 预加载下一个视频，实现无缝切换
      const nextIndex = (currentIndex.value + 1) % props.videos.length;
      if (nextIndex !== currentIndex.value) { // 确保不是只有一个视频的情况
        videoB.src = props.videos[nextIndex];
        videoB.load();
      }
    }
  }
});

onUnmounted(() => {
  showNavBar(); // 离开视频播放页时显示导航栏
  window.removeEventListener('keydown', handleKeyDown);
});
</script>

<style scoped>
.video-player-wrapper {
  position: fixed;
  top: 0;
  left: 0;
  width: 100vw;
  /* ==================== 核心修复 1：使用 dvh 替代 vh ==================== */
  /* dvh (dynamic viewport height) 会根据浏览器UI（如地址栏）的出现和消失动态调整高度 */
  height: 100dvh;
  /* ====================================================================== */
  background-color: #000;
  display: flex;
  justify-content: center;
  align-items: center;
  outline: none;
  /* 防止在iOS上拖动页面产生回弹效果 */
  overflow: hidden;
}
.video-container { 
  width: 100%; 
  height: 100%; 
  position: relative; 
  overflow: hidden; 
}
.video-player { 
  position: absolute; 
  top: 0; 
  left: 0; 
  width: 100%; 
  height: 100%; 
  object-fit: contain; 
  opacity: 0; 
  transition: opacity 0.3s ease-in-out; 
}
.video-player.active { 
  opacity: 1; 
  z-index: 2; 
}
.click-overlay { 
  position: absolute; 
  top: 0; 
  left: 0; 
  width: 100%; 
  height: 100%; 
  display: flex; 
  justify-content: center; 
  align-items: center; 
  z-index: 5; 
}
.play-pause-icon-container { 
  background-color: rgba(0, 0, 0, 0.3); 
  width: 100%; 
  height: 100%; 
  display: flex; 
  justify-content: center; 
  align-items: center; 
  opacity: 0; 
  pointer-events: none; 
  transition: opacity 0.3s ease; 
}
.play-pause-icon-container.is-visible { 
  opacity: 1; 
}
.play-pause-btn { 
  background: rgba(255, 255, 255, 0.2); 
  border: 2px solid white; 
  color: white; 
  width: 80px; 
  height: 80px; 
  border-radius: 50%; 
  font-size: 2.5rem; 
  line-height: 80px; 
  padding-left: 10px; 
  cursor: pointer; 
}
.controls-overlay { 
  position: absolute; 
  top: 0; 
  left: 0; 
  width: 100%; 
  height: 100%; 
  opacity: 0; 
  transition: opacity 0.3s ease; 
  pointer-events: none; 
}
.video-container:hover .controls-overlay { 
  opacity: 1; 
}
.progress-indicator { 
  position: absolute; 
  bottom: 20px; 
  right: 20px; 
  background-color: rgba(0,0,0,0.6); 
  color: white; 
  padding: 5px 10px; 
  border-radius: 5px; 
  font-size: 0.9rem; 
  z-index: 10; 
}
.loading { 
  color: #fff; 
}
.volume-control { 
  position: absolute; 
  bottom: 20px; 
  left: 20px; 
  z-index: 10; 
  display: flex; 
  align-items: center; 
}
.volume-btn { 
  background: rgba(0, 0, 0, 0.5); 
  border: none; 
  color: white; 
  width: 44px; 
  height: 44px; 
  border-radius: 50%; 
  font-size: 1.2rem; 
  cursor: pointer; 
  display: flex; 
  justify-content: center; 
  align-items: center; 
  transition: background-color 0.2s; 
}
.volume-btn:hover { 
  background-color: rgba(0, 0, 0, 0.8); 
}
.volume-slider { 
  margin-left: 8px; 
  width: 80px; 
  -webkit-appearance: none; 
  appearance: none; 
  height: 5px; 
  background: rgba(255, 255, 255, 0.4); 
  border-radius: 5px; 
  outline: none; 
  opacity: 0.7; 
  transition: opacity .2s; 
  cursor: pointer; 
}
.volume-slider:hover { 
  opacity: 1; 
}
.volume-slider::-webkit-slider-thumb { 
  -webkit-appearance: none; 
  appearance: none; 
  width: 15px; 
  height: 15px; 
  background: white; 
  border-radius: 50%; 
  cursor: pointer; 
}
.volume-slider::-moz-range-thumb { 
  width: 15px; 
  height: 15px; 
  background: white; 
  border-radius: 50%; 
  cursor: pointer; 
}
.side-controls { 
  position: absolute; 
  bottom: 90px; 
  right: 15px; 
  z-index: 10; 
  display: flex; 
  flex-direction: column; 
  gap: 1.5rem; 
}
.control-btn { 
  background: none; 
  border: none; 
  color: white; 
  cursor: pointer; 
  display: flex; 
  flex-direction: column; 
  align-items: center; 
  gap: 0.3rem; 
  font-size: 0.9rem; 
  padding: 0; 
}
.control-btn .icon { 
  font-size: 2.2rem; 
  text-shadow: 0 2px 4px rgba(0,0,0,0.5); 
  transition: transform 0.2s; 
}
.control-btn:hover .icon { 
  transform: scale(1.1); 
}
.control-btn .count { 
  font-weight: 500; 
  text-shadow: 0 1px 3px rgba(0,0,0,0.5); 
}

/* ==================== 核心修复 2：为移动端安全区域添加边距 ==================== */
/* env(safe-area-inset-bottom) 是一个CSS变量，代表iOS等设备底部安全区域的高度 */
@media (max-width: 768px) {
  .side-controls {
    bottom: calc(90px + env(safe-area-inset-bottom, 0px));
  }
  .volume-control, .progress-indicator {
    bottom: calc(20px + env(safe-area-inset-bottom, 0px));
  }
}

/* 当屏幕宽度小于等于 768px 时 (通常是手机) */
@media (max-width: 768px) {
  /* 隐藏音量控制条，因为用户会使用物理按键 */
  .volume-control {
    display: none;
  }

  /* (可选) 调整侧边按钮的位置，给底部安全区留出更多空间 */
  .side-controls {
    bottom: calc(100px + env(safe-area-inset-bottom, 0px));
  }
  .progress-indicator {
    bottom: calc(30px + env(safe-area-inset-bottom, 0px));
  }
}
/* ====================================================================== */
</style>