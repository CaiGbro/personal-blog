<!-- /frontend/src/components/CommentItem.vue (默认折叠 + 增量加载) -->
<template>
  <div class="comment-wrapper">
    <!-- 评论主体 -->
    <div class="comment-item">
      <!-- ... (左侧头像和连接线部分，无改动) ... -->
      <div class="comment-gutter">
        <div class="avatar-container">
          <span>{{ comment.nickname.charAt(0).toUpperCase() }}</span>
        </div>
        <div v-if="hasReplies" class="thread-line"></div>
      </div>

      <!-- 右侧：主要内容 -->
      <div class="comment-main">
        <!-- ... (评论头部、内容、操作区，无改动) ... -->
        <div class="comment-header">
          <strong>{{ comment.nickname }}</strong>
          <span class="timestamp">{{ formatTimestamp(comment.createdAt) }}</span>
        </div>
        <div class="comment-content">
          <p v-html="formattedContent"></p>
          <img v-if="comment.attachmentUrl" :src="comment.attachmentUrl" class="comment-attachment" />
        </div>
        <div class="comment-actions">
          <div class="reactions">
            <span v-for="(count, emoji) in comment.reactions" :key="emoji" class="reaction-badge" @click="addReaction(emoji)">
              {{ emoji }} {{ count }}
            </span>
            <button class="add-reaction-btn" @click.stop="showEmojiPicker = !showEmojiPicker">☺</button>
            <div v-if="showEmojiPicker" class="emoji-picker">
              <span v-for="emoji in popularEmojis" :key="emoji" @click="addReaction(emoji)">{{ emoji }}</span>
            </div>
          </div>
          <button @click="showReplyForm = !showReplyForm" class="reply-btn">回复</button>
        </div>
        <CommentForm v-if="showReplyForm" :parent-id="comment.id" :initial-content="`@${comment.nickname} `" @comment-posted="onReplyPosted" />
        
        <div v-if="hasReplies" class="replies-toggle-bar">
          <button @click="toggleReplies" class="toggle-replies-btn">
            <span v-if="isRepliesExpanded">收起回复</span>
            <span v-else>
              展开 {{ totalRepliesCount }} 条回复
            </span>
          </button>
        </div>
      </div>
    </div>

    <!-- 回复列表容器 -->
    <div v-if="hasReplies && isRepliesExpanded" class="replies-container">
      <CommentItem
        v-for="reply in visibleReplies"
        :key="reply.id"
        :comment="reply"
        @comment-updated="$emit('comment-updated')"
      />
      <div v-if="hasMoreReplies" class="show-more-container">
        <button @click="loadMoreReplies" class="show-more-btn">
          <!-- ==================== 按钮文本优化 ==================== -->
          查看下 {{ nextLoadCount }} 条回复
          <!-- ====================================================== -->
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue';
import api from '../services/api';
import CommentForm from './CommentForm.vue';

const props = defineProps({
  comment: { type: Object, required: true },
});

const emit = defineEmits(['comment-updated']);

const showReplyForm = ref(false);
const showEmojiPicker = ref(false);
const popularEmojis = ['👍', '❤️', '😂', '😮', '😢', '🔥'];
const hasReplies = computed(() => props.comment.replies && props.comment.replies.length > 0);

// ==================== 核心改动 ====================

// === 核心改动 1: 调整分页大小和默认状态 ===
const REPLIES_PAGE_SIZE = 3; // 每次显示的回复数量改为 3
const isRepliesExpanded = ref(false); // 回复默认是折叠的
const visibleRepliesCount = ref(REPLIES_PAGE_SIZE); // 初始可见数量仍为分页大小

// === 核心改动 2: 修改加载和切换逻辑 ===

// 方法：切换回复的折叠/展开状态
const toggleReplies = () => {
  isRepliesExpanded.value = !isRepliesExpanded.value;
  
  // 优化体验：当用户收起回复时，重置可见回复的数量。
  // 这样下次展开时，总是从第一页（3条）开始看。
  if (!isRepliesExpanded.value) {
    visibleRepliesCount.value = REPLIES_PAGE_SIZE;
  }
};

// 方法：加载更多回复（增量加载）
const loadMoreReplies = () => {
  visibleRepliesCount.value += REPLIES_PAGE_SIZE;
};

// ========================================================


// --- 计算属性部分 ---

const visibleReplies = computed(() => {
  if (!props.comment.replies) return [];
  return props.comment.replies.slice(0, visibleRepliesCount.value);
});

const hasMoreReplies = computed(() => {
  return props.comment.replies && visibleRepliesCount.value < props.comment.replies.length;
});

const remainingRepliesCount = computed(() => {
  if (!props.comment.replies) return 0;
  return props.comment.replies.length - visibleRepliesCount.value;
});

// 新增计算属性：用于优化 "查看更多" 按钮的文本
const nextLoadCount = computed(() => {
  return Math.min(REPLIES_PAGE_SIZE, remainingRepliesCount.value);
});


// --- 递归计数部分 (无改动) ---

const countAllReplies = (repliesArray) => {
  if (!repliesArray || repliesArray.length === 0) {
    return 0;
  }
  let count = repliesArray.length;
  for (const reply of repliesArray) {
    count += countAllReplies(reply.replies);
  }
  return count;
};

const totalRepliesCount = computed(() => {
  return countAllReplies(props.comment.replies);
});


// --- 其他方法 (无改动) ---

const formatTimestamp = (isoString) => new Date(isoString).toLocaleString();

const formattedContent = computed(() => {
  return props.comment.content.replace(
    /^(@\S+\s)/,
    '<span class="mention">$1</span>'
  );
});

const addReaction = async (emoji) => {
  try {
    await api.post(`/comments/${props.comment.id}/react`, { emoji });
    emit('comment-updated');
    showEmojiPicker.value = false;
  } catch (error) {
    console.error('添加表情失败:', error);
  }
};

const onReplyPosted = () => {
  showReplyForm.value = false;
  emit('comment-updated');
};
</script>

<style scoped>
/* ... (所有样式保持不变) ... */
.comment-item {
  display: flex;
  gap: 12px;
}
.comment-gutter {
  display: flex;
  flex-direction: column;
  align-items: center;
  flex-shrink: 0;
}
.avatar-container {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background-color: #f0f2f5;
  display: flex;
  justify-content: center;
  align-items: center;
  font-weight: bold;
  color: #555;
  font-size: 1.1rem;
}
.thread-line {
  width: 2px;
  background-color: #e1e8ed;
  flex-grow: 1;
  margin-top: 8px;
}
.comment-main {
  flex-grow: 1;
}
.comment-wrapper {
  padding-bottom: 1rem;
  border-bottom: 1px solid #eee;
}
.guestbook-section .comment-list > .comment-wrapper:last-child {
  border-bottom: none;
  padding-bottom: 0;
}
.replies-container > .comment-wrapper:last-child {
    padding-bottom: 0;
}
.comment-header { margin-bottom: 0.5rem; }
.timestamp { font-size: 0.8rem; color: #999; margin-left: 0.5rem; }
.comment-content p { margin: 0; white-space: pre-wrap; word-break: break-word; }
.comment-attachment { max-width: 100%; max-height: 250px; border-radius: 8px; margin-top: 0.75rem; }
.comment-actions { display: flex; justify-content: space-between; align-items: center; margin-top: 1rem; gap: 1rem; }
.reactions { position: relative; display: flex; align-items: center; gap: 0.5rem; flex-wrap: wrap; }
.reaction-badge { background: #f0f0f0; padding: 0.25rem 0.5rem; border-radius: 1rem; cursor: pointer; font-size: 0.9rem; }
.add-reaction-btn, .reply-btn { background: none; border: none; cursor: pointer; color: #555; font-weight: 500; }
.emoji-picker { position: absolute; bottom: 120%; left: 0; background: white; border: 1px solid #ccc; border-radius: 8px; padding: 0.5rem; display: flex; gap: 0.5rem; box-shadow: 0 4px 8px rgba(0,0,0,0.1); z-index: 10; }
.emoji-picker span { cursor: pointer; font-size: 1.2rem; }

:deep(.mention) {
  color: #007bff;
  font-weight: 500;
}

.replies-container {
  padding-left: 52px;
  margin-top: 1rem;
}
.replies-container .replies-container {
  padding-left: 0;
}
.replies-container .comment-main {
    border-bottom: none;
}
.replies-toggle-bar {
  margin-top: 1rem;
}
.toggle-replies-btn {
  background: none;
  border: none;
  color: #007bff;
  cursor: pointer;
  font-weight: 500;
  font-size: 0.9rem;
  padding: 0.25rem 0;
}
.show-more-container {
  display: flex;
  margin-top: 1rem;
}
.show-more-btn {
  background: none;
  border: none;
  color: #007bff;
  cursor: pointer;
  font-weight: 500;
  font-size: 0.9rem;
  padding: 0.25rem 0;
}
</style>