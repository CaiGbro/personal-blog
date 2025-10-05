<!-- /frontend/src/components/VideoCommentItem.vue (已修复表情点击问题) -->
<template>
  <div class="comment-wrapper">
    <div class="comment-item">
      <div class="comment-gutter">
        <div class="avatar-container">
          <span>{{ comment.nickname.charAt(0).toUpperCase() }}</span>
        </div>
        <div v-if="hasReplies" class="thread-line"></div>
      </div>
      <div class="comment-main">
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
            <!-- 这个 emoji-picker 的样式被修改了 -->
            <div v-if="showEmojiPicker" class="emoji-picker">
              <span v-for="emoji in popularEmojis" :key="emoji" @click="addReaction(emoji)">{{ emoji }}</span>
            </div>
          </div>
          <button @click="showReplyForm = !showReplyForm" class="reply-btn">回复</button>
        </div>
        
        <VideoCommentForm 
          v-if="showReplyForm" 
          :parent-id="comment.id" 
          :initial-content="`@${comment.nickname} `" 
          :video-filename="videoFilename" 
          @comment-posted="onReplyPosted" 
        />
        
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
    <div v-if="hasReplies && isRepliesExpanded" class="replies-container">
      <VideoCommentItem
        v-for="reply in visibleReplies"
        :key="reply.id"
        :comment="reply"
        :video-filename="videoFilename"
        @comment-updated="$emit('comment-updated')"
      />
      <div v-if="hasMoreReplies" class="show-more-container">
        <button @click="loadMoreReplies" class="show-more-btn">
          查看下 {{ nextLoadCount }} 条回复
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
// --- Script 部分无改动 ---
import { ref, computed } from 'vue';
import api from '../services/api';
import VideoCommentForm from './VideoCommentForm.vue';
const props = defineProps({ comment: { type: Object, required: true }, videoFilename: { type: String, required: true } });
const emit = defineEmits(['comment-updated']);
const showReplyForm = ref(false);
const showEmojiPicker = ref(false);
const popularEmojis = ['👍', '❤️', '😂', '😮', '😢', '🔥'];
const hasReplies = computed(() => props.comment.replies && props.comment.replies.length > 0);
const REPLIES_PAGE_SIZE = 3;
const isRepliesExpanded = ref(false);
const visibleRepliesCount = ref(REPLIES_PAGE_SIZE);
const toggleReplies = () => { isRepliesExpanded.value = !isRepliesExpanded.value; if (!isRepliesExpanded.value) { visibleRepliesCount.value = REPLIES_PAGE_SIZE; } };
const loadMoreReplies = () => { visibleRepliesCount.value += REPLIES_PAGE_SIZE; };
const visibleReplies = computed(() => { if (!props.comment.replies) return []; return props.comment.replies.slice(0, visibleRepliesCount.value); });
const hasMoreReplies = computed(() => { return props.comment.replies && visibleRepliesCount.value < props.comment.replies.length; });
const remainingRepliesCount = computed(() => { if (!props.comment.replies) return 0; return props.comment.replies.length - visibleRepliesCount.value; });
const nextLoadCount = computed(() => { return Math.min(REPLIES_PAGE_SIZE, remainingRepliesCount.value); });
const countAllReplies = (repliesArray) => { if (!repliesArray || repliesArray.length === 0) return 0; let count = repliesArray.length; for (const reply of repliesArray) { count += countAllReplies(reply.replies); } return count; };
const totalRepliesCount = computed(() => countAllReplies(props.comment.replies));
const formatTimestamp = (isoString) => new Date(isoString).toLocaleString();
const formattedContent = computed(() => { return props.comment.content.replace( /^(@\S+\s)/, '<span class="mention">$1</span>' ); });
const addReaction = async (emoji) => { try { await api.post(`/video_comments/${props.comment.id}/react`, { emoji }); emit('comment-updated'); showEmojiPicker.value = false; } catch (error) { console.error('添加表情失败:', error); } };
const onReplyPosted = () => { showReplyForm.value = false; emit('comment-updated'); };
</script>

<style scoped>
/* 基础布局和默认样式 */
.comment-item { display: flex; gap: 12px; }
.comment-gutter { display: flex; flex-direction: column; align-items: center; flex-shrink: 0; }
.avatar-container { width: 40px; height: 40px; border-radius: 50%; background-color: #f0f2f5; display: flex; justify-content: center; align-items: center; font-weight: bold; color: #555; font-size: 1.1rem; }
.thread-line { width: 2px; background-color: #e1e8ed; flex-grow: 1; margin-top: 8px; }
.comment-main { flex-grow: 1; }
.comment-wrapper { padding-bottom: 1rem; border-bottom: 1px solid #eee; }
.replies-container > .comment-wrapper:last-child { padding-bottom: 0; border-bottom: none; }
.comment-header { margin-bottom: 0.5rem; }
.timestamp { font-size: 0.8rem; color: #999; margin-left: 0.5rem; }
.comment-content p { margin: 0; white-space: pre-wrap; word-break: break-word; }
.comment-attachment { max-width: 100%; max-height: 250px; border-radius: 8px; margin-top: 0.75rem; }
.comment-actions { display: flex; justify-content: space-between; align-items: center; margin-top: 1rem; gap: 1rem; }
.reactions { position: relative; display: flex; align-items: center; gap: 0.5rem; flex-wrap: wrap; }
.reaction-badge { background: #f0f0f0; padding: 0.25rem 0.5rem; border-radius: 1rem; cursor: pointer; font-size: 0.9rem; }
.add-reaction-btn, .reply-btn { background: none; border: none; cursor: pointer; color: #555; font-weight: 500; }
:deep(.mention) { color: #007bff; font-weight: 500; }
.replies-container { padding-left: 52px; margin-top: 1rem; }
.replies-container .replies-container { padding-left: 0; }
.replies-toggle-bar { margin-top: 1rem; }
.toggle-replies-btn, .show-more-btn { background: none; border: none; color: #007bff; cursor: pointer; font-weight: 500; font-size: 0.9rem; padding: 0.25rem 0; }

/* ==================== 【关键修复】 ==================== */
.emoji-picker {
  position: absolute;
  /* 将向上弹出 (bottom) 改为向下弹出 (top) */
  /* 'calc(100% + 8px)' 表示在父容器下方 8px 的位置弹出 */
  top: calc(100% + 8px); 
  left: 0;
  background: white;
  border: 1px solid #ccc;
  border-radius: 8px;
  padding: 0.5rem;
  display: flex;
  gap: 0.5rem;
  box-shadow: 0 4px 8px rgba(0,0,0,0.1);
  z-index: 20; /* 提升 z-index 以确保它显示在回复框等元素的上方 */
}
/* ====================================================== */

.emoji-picker span {
  cursor: pointer;
  font-size: 1.2rem;
}
</style>