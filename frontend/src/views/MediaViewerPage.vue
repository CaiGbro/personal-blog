<!-- /frontend/src/views/MediaViewerPage.vue (已修复广告Bug) -->
<template>
  <div>
    <VideoPlayer 
      v-if="videoList.length > 0 && startWithSrc" 
      :videos="videoList"
      :start-with-src="startWithSrc"
    />
    <div v-else class="status-message">
      <p>{{ message }}</p>
    </div>
  </div>
</template>

<script setup>
// 1. 从 'vue' 中导入 onUnmounted
import { ref, onMounted, onUnmounted } from 'vue';
import { useRoute } from 'vue-router';
import api from '../services/api';
import VideoPlayer from '../components/VideoPlayer.vue';
// 2. 导入广告控制逻辑
import { useAds } from '../store/ads';

const videoList = ref([]);
const message = ref('正在加载视频...');
const route = useRoute();
const startWithSrc = ref(null);

// 3. 获取广告状态和控制函数
const { adsEnabled, enableAds, disableAds } = useAds();
// 4. 创建一个本地变量，用于记住进入此页面前的广告状态
const adsWereEnabledBeforeEntering = ref(false);

const fetchVideos = async () => {
  try {
    const response = await api.get('/videos');
    if (response.data && response.data.length > 0) {
      videoList.value = response.data;
      startWithSrc.value = route.query.start_with || response.data[0];
    } else {
      message.value = '没有找到可播放的视频。';
    }
  } catch (error) {
    console.error('获取视频列表失败:', error);
    message.value = '加载视频失败，请稍后再试。';
  }
};

// 5. 使用 onMounted 生命周期钩子
onMounted(() => {
  // 在组件挂载时（即进入页面时）
  // a. 记录下当前的广告状态
  adsWereEnabledBeforeEntering.value = adsEnabled.value;
  
  // b. 如果广告是开启的，则立即关闭它们
  if (adsEnabled.value) {
    disableAds();
  }

  // c. 继续执行原有的视频加载逻辑
  fetchVideos();
});

// 6. 使用 onUnmounted 生命周期钩子
onUnmounted(() => {
  // 在组件卸载时（即离开页面时）
  // a. 检查我们之前记录的状态
  // b. 如果广告在进入前是开启的，就重新开启它们
  if (adsWereEnabledBeforeEntering.value) {
    enableAds();
  }
});
</script>

<style scoped>
.status-message {
  display: flex;
  justify-content: center;
  align-items: center;
  height: calc(100vh - 60px);
  color: #fff;
  background-color: #000;
  font-size: 1.2rem;
}
</style>