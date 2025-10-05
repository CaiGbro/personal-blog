<!-- /frontend/src/components/WritingCard.vue (已清理日志) -->
<template>
  <router-link :to="`/writings/${encodeURIComponent(writing.filename)}`" class="writing-card">
    <div class="thumbnail-container">
      <img v-if="isThumbnailAnImage" :src="writing.thumbnailUrl" alt="Thumbnail" class="thumbnail"/>
      <div v-else class="thumbnail-placeholder">
        <span>{{ fileExtension }}</span>
      </div>
    </div>
    <div class="card-info">
      <h3 class="title">{{ formattedTitle }}</h3>
      <p class="filename">{{ writing.filename }}</p>
    </div>
  </router-link>
</template>

<script setup>
import { computed } from 'vue';

const props = defineProps({
  writing: { type: Object, required: true }
});

const fileExtension = computed(() => {
  const parts = props.writing.filename.split('.');
  return parts.length > 1 ? parts.pop().toUpperCase() : 'FILE';
});

const formattedTitle = computed(() => {
  if (!props.writing.filename) return '';
  let title = props.writing.filename.split('.').slice(0, -1).join('.');
  title = title.replace(/[-_]/g, ' ');
  return title;
});

const isThumbnailAnImage = computed(() => {
  const url = props.writing.thumbnailUrl;
  
  if (!url || typeof url !== 'string') {
    return false;
  }
  
  // 正则表达式判断 URL 是否以常见的图片格式结尾
  return /\.(png|jpg|jpeg|gif|webp|svg)$/i.test(url);
});

</script>

<style scoped>
/* 样式保持不变 */
.writing-card {
  display: flex;
  flex-direction: column;
  break-inside: avoid;
  margin-bottom: 1.5rem;
  border-radius: 8px;
  overflow: hidden;
  box-shadow: 0 4px 12px rgba(0,0,0,0.08);
  transition: transform 0.2s ease, box-shadow 0.2s ease;
  background-color: #fff;
}
.thumbnail-container {
  width: 100%;
  aspect-ratio: 16 / 10;
  background-color: #f0f2f5;
  border-bottom: 1px solid #eee;
}
.thumbnail {
  width: 100%;
  height: 100%;
  object-fit: cover;
}
.thumbnail-placeholder {
  width: 100%;
  height: 100%;
  display: flex;
  justify-content: center;
  align-items: center;
  font-size: 1.5rem;
  font-weight: bold;
  color: #a0a0a0;
}
.card-info {
  padding: 1rem 1.25rem;
}
.title {
  font-size: 1.1rem;
  font-weight: 600;
  color: #333;
  margin: 0 0 0.5rem 0;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
.filename {
  margin: 0;
  font-weight: 400;
  font-size: 0.85rem;
  color: #999;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
</style>