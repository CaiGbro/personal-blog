<!-- /frontend/src/views/VideoGalleryPage.vue (已适配移动端单列) -->
<template>
  <div class="gallery-container">
    <div class="page-header">
      <h1>美与欣赏</h1>
      <p>光影流转，瞬间永恒。</p>
    </div>
    <div v-if="videos.length > 0" class="gallery-grid">
      <VideoCard v-for="videoSrc in videos" :key="videoSrc" :src="videoSrc" />
    </div>
    <div v-else class="status-message">
      <p>正在加载视频...</p>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import api from '../services/api';
import VideoCard from '../components/VideoCard.vue';

const videos = ref([]);

const fetchVideos = async () => {
  try {
    const response = await api.get('/videos');
    if (response.data && response.data.length > 0) {
      videos.value = response.data;
    }
  } catch (error) {
    console.error('获取视频列表失败:', error);
  }
};

onMounted(fetchVideos);
</script>

<script>
export default {
  name: 'VideoGalleryPage'
}
</script>

<style scoped>
.gallery-container {
  max-width: 960px;
  margin: 2rem auto;
  padding: 2rem;
}
.page-header {
  text-align: center;
  margin-bottom: 3rem;
}
.gallery-grid {
  /* 默认在桌面端显示为两列 */
  column-count: 2;
  column-gap: 1.5rem;
}
.status-message {
  text-align: center;
  padding: 4rem 0;
  color: #888;
}

/* ==================== 新增：移动端适配 ==================== */
/* 当屏幕宽度小于或等于 768px 时 (通常是平板和手机) */
@media (max-width: 768px) {
  .gallery-grid {
    /* 在小屏幕上，将列数改为 1 */
    column-count: 1;
  }

  .gallery-container {
    /* 同时减小页面边距，让内容区域更大 */
    padding: 1rem;
  }
}
/* ====================================================== */
</style>