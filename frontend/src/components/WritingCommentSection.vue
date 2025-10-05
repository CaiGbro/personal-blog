<!-- /frontend/src/components/WritingCommentSection.vue -->
<template>
  <!-- template 部分无需改动 -->
  <div class="comments-section">
    <div class="comments-header">
      <h4>{{ totalCount }} 条评论</h4>
    </div>
    <div class="comments-body">
      <WritingCommentItem
        v-for="comment in comments"
        :key="comment.id"
        :comment="comment"
        :writing-filename="writingFilename"
        @comment-updated="fetchComments"
      />
      <p v-if="comments.length === 0" class="no-comments">还没有评论，快来抢沙发吧！</p>
    </div>
    <div class="comments-footer">
      <WritingCommentForm :writing-filename="writingFilename" @comment-posted="fetchComments" />
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue';
import api from '../services/api';
import WritingCommentItem from './WritingCommentItem.vue';
import WritingCommentForm from './WritingCommentForm.vue';

const props = defineProps({ writingFilename: { type: String, required: true } });
// ==================== 修改代码 ====================
const emit = defineEmits(['count-updated']);
// ================================================

const comments = ref([]);
const totalCount = ref(0);

const fetchComments = async () => {
  if (!props.writingFilename) return;
  try {
    const response = await api.get('/writing_comments', {
      params: { writing_filename: props.writingFilename },
    });
    comments.value = response.data.comments || [];
    totalCount.value = response.data.totalCount || 0;
    
    // ==================== 修改代码 ====================
    // 将总数发送给父组件
    emit('count-updated', totalCount.value);
    // ================================================

  } catch (error) {
    console.error('获取文章评论失败:', error);
  }
};
watch(() => props.writingFilename, fetchComments);
onMounted(fetchComments);
</script>

<style scoped>
/* style 部分无需改动 */
.comments-section {
  width: 100%;
  height: 100%;
  background-color: #fff;
  display: flex;
  flex-direction: column;
  overflow-y: hidden;
}
.comments-header { padding: 1rem; border-bottom: 1px solid #eee; }
.comments-header h4 { margin: 0; }
.comments-body { flex-grow: 1; overflow-y: auto; padding: 1rem; }
.no-comments { text-align: center; color: #888; margin-top: 2rem; }
.comments-footer { padding: 0 1rem; border-top: 1px solid #eee; }
</style>