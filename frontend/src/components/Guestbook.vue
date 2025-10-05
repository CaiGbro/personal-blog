<!-- /frontend/src/components/Guestbook.vue -->
<template>
  <!-- template 部分无需改动 -->
  <div class="guestbook-section">
    <h3>留下新评论</h3>
    <CommentForm @comment-posted="() => fetchComments(1)" />

    <div class="comment-list">
      <CommentItem
        v-for="comment in comments"
        :key="comment.id"
        :comment="comment"
        @comment-updated="() => fetchComments(currentPage)"
      />
    </div>

    <div v-if="totalPages > 1" class="pagination">
      <button @click="changePage(page)" v-for="page in totalPages" :key="page" :class="{ active: page === currentPage }">
        {{ page }}
      </button>
    </div>
  </div>
</template>

<script setup>
// script 部分无需改动
import { ref, onMounted, computed } from 'vue';
import api from '../services/api';
import CommentForm from './CommentForm.vue';
import CommentItem from './CommentItem.vue';

const comments = ref([]);
const currentPage = ref(1);
const pageSize = ref(5);
const totalCount = ref(0);

const totalPages = computed(() => Math.ceil(totalCount.value / pageSize.value));

const fetchComments = async (page) => {
  try {
    const response = await api.get('/comments', {
      params: { page, pageSize: pageSize.value }
    });
    comments.value = response.data.comments;
    totalCount.value = response.data.totalCount;
    currentPage.value = page;
  } catch (error) {
    console.error('获取留言失败:', error);
  }
};

const changePage = (page) => {
  if (page > 0 && page <= totalPages.value) {
    fetchComments(page);
  }
};

onMounted(() => {
  fetchComments(currentPage.value);
});
</script>

<style scoped>
/* 样式已简化 */
.guestbook-section { margin-top: 3rem; }
.comment-list {
  /* 移除边框，因为边框现在由 CommentItem 内部处理 */
  margin-top: 2rem;
}
.pagination { margin-top: 1.5rem; display: flex; gap: 0.5rem; }
.pagination button { padding: 0.5rem 0.75rem; border: 1px solid #ccc; cursor: pointer; }
.pagination button.active { background-color: #86b93a; color: white; border-color: #86b93a; }
</style>