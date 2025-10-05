<!-- /frontend/src/components/VideoCommentSection.vue (最终修复版) -->
<template>
  <div class="video-comments-overlay" @click.self="$emit('close')">
    <div class="video-comments-section" @wheel.stop>
      <div class="comments-header">
        <h4>{{ totalCount }} 条评论</h4>
        <button @click="$emit('close')" class="close-btn">&times;</button>
      </div>
      <div class="comments-body">
        <VideoCommentItem
          v-for="comment in comments"
          :key="comment.id"
          :comment="comment"
          :video-filename="videoFilename"
          @comment-updated="fetchVideoComments"
        />
         <p v-if="comments.length === 0" class="no-comments">还没有评论，快来抢沙发吧！</p>
      </div>
      <div class="comments-footer">
        <VideoCommentForm :video-filename="videoFilename" @comment-posted="fetchVideoComments" />
      </div>
    </div>
  </div>
</template>

<script setup>
// --- Script 部分无改动 ---
import { ref, onMounted, watch } from 'vue';
import api from '../services/api';
import VideoCommentItem from './VideoCommentItem.vue';
import VideoCommentForm from './VideoCommentForm.vue';
const props = defineProps({ videoFilename: { type: String, required: true } });
const emit = defineEmits(['close', 'count-updated']);
const comments = ref([]);
const totalCount = ref(0);

const fetchVideoComments = async () => {
  if (!props.videoFilename) return;
  try {
    const response = await api.get('/video_comments', {
      params: { video_filename: props.videoFilename },
    });
    
    comments.value = response.data.comments || [];
    totalCount.value = response.data.totalCount || 0; // 直接从后端获取总数
    
    // 将真实总数发送给父组件
    emit('count-updated', totalCount.value);

  } catch (error) {
    console.error('获取视频评论失败:', error);
  }
};

watch(() => props.videoFilename, fetchVideoComments);
onMounted(fetchVideoComments);
</script>

<style scoped>
.video-comments-overlay { position: fixed; top: 0; left: 0; width: 100%; height: 100%; background: rgba(0, 0, 0, 0.5); z-index: 1000; display: flex; justify-content: flex-end; }
.video-comments-section { width: 400px; height: 100%; background: #121212; color: #fff; display: flex; flex-direction: column; animation: slideIn 0.3s ease-out; }
@keyframes slideIn { from { transform: translateX(100%); } to { transform: translateX(0); } }
.comments-header { padding: 1rem; display: flex; justify-content: space-between; align-items: center; border-bottom: 1px solid #333; }
.comments-header h4 { margin: 0; }
.close-btn { background: none; border: none; color: #fff; font-size: 1.5rem; cursor: pointer; }
.comments-body { flex-grow: 1; overflow-y: auto; padding: 1rem; }
.no-comments { text-align: center; color: #888; margin-top: 2rem; }
.comments-footer { padding: 0 1rem; }

/* ==================== 【最终修复方案】 ==================== */
/* 从父组件使用 :deep() 强制定义子组件在深色主题下的所有样式 */

/* 整体文字和边框颜色 */
:deep(.comment-wrapper) {
  color: #e0e0e0;
  border-bottom-color: #333;
}

/* 头像 */
:deep(.avatar-container) {
  background-color: #333;
  color: #fff;
}

/* 时间戳 */
:deep(.timestamp) {
  color: #888;
}

/* 连接线 */
:deep(.thread-line) {
  background-color: #444;
}

/* @用户 */
:deep(:deep(.mention)) { /* 需要双重:deep修复子组件的:deep */
  color: #4a90e2;
}

/* ==================== 核心修改：表情徽章样式 ==================== */
:deep(.reaction-badge) {
  background-color: #3a3a3a; /* 深色背景 */
  color: #e0e0e0;             /* 亮色文字，解决数字看不清的问题 */
  border: 1px solid #555;    /* 添加微妙的边框以增加立体感 */
  padding: 4px 8px;           /* 增加内边距，让内容更舒展 */
  display: inline-flex;       /* 使用flex布局 */
  align-items: center;        /* 确保表情和数字垂直居中 */
  gap: 4px;                   /* 在表情和数字之间增加一点间距 */
  font-size: 0.95rem;         /* 稍微增大字体，提升可读性 */
}
/* ============================================================ */

/* 表情选择器 */
:deep(.emoji-picker) {
  background-color: #222;
  border-color: #444;
}

/* --- 核心：确保所有按钮清晰可见 --- */
:deep(.add-reaction-btn),
:deep(.reply-btn),
:deep(.toggle-replies-btn),
:deep(.show-more-btn) {
  color: #a0a0a0; /* 一个清晰的灰色作为默认状态 */
  font-weight: 500;
  opacity: 1;
  transition: color 0.2s ease;
}

:deep(.add-reaction-btn:hover),
:deep(.reply-btn:hover),
:deep(.toggle-replies-btn:hover),
:deep(.show-more-btn:hover) {
  color: #ffffff; /* 鼠标悬停时变为纯白色，提供反馈 */
}
/* ====================================================== */
</style>