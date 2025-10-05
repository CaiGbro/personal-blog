<!-- /frontend/src/components/WritingCommentItem.vue -->
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
        <div class="comment-header"><strong>{{ comment.nickname }}</strong><span class="timestamp">{{ formatTimestamp(comment.createdAt) }}</span></div>
        <div class="comment-content">
          <p v-html="formattedContent"></p>
          <img v-if="comment.attachmentUrl" :src="comment.attachmentUrl" class="comment-attachment" />
        </div>
        <div class="comment-actions">
          <div class="reactions">
            <span v-for="(count, emoji) in comment.reactions" :key="emoji" class="reaction-badge" @click="addReaction(emoji)">{{ emoji }} {{ count }}</span>
            <button class="add-reaction-btn" @click.stop="showEmojiPicker = !showEmojiPicker">☺</button>
            <div v-if="showEmojiPicker" class="emoji-picker">
              <span v-for="emoji in popularEmojis" :key="emoji" @click="addReaction(emoji)">{{ emoji }}</span>
            </div>
          </div>
          <button @click="showReplyForm = !showReplyForm" class="reply-btn">回复</button>
        </div>
        <WritingCommentForm v-if="showReplyForm" :parent-id="comment.id" :initial-content="`@${comment.nickname} `" :writing-filename="writingFilename" @comment-posted="onReplyPosted" />
        <div v-if="hasReplies" class="replies-toggle-bar">
          <button @click="toggleReplies" class="toggle-replies-btn">
            <span v-if="isRepliesExpanded">收起回复</span><span v-else>展开 {{ totalRepliesCount }} 条回复</span>
          </button>
        </div>
      </div>
    </div>
    <div v-if="hasReplies && isRepliesExpanded" class="replies-container">
      <WritingCommentItem v-for="reply in visibleReplies" :key="reply.id" :comment="reply" :writing-filename="writingFilename" @comment-updated="$emit('comment-updated')" />
      <div v-if="hasMoreReplies" class="show-more-container"><button @click="loadMoreReplies" class="show-more-btn">查看下 {{ nextLoadCount }} 条回复</button></div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue';
import api from '../services/api';
import WritingCommentForm from './WritingCommentForm.vue';

const props = defineProps({ comment: Object, writingFilename: String });
const emit = defineEmits(['comment-updated']);
const showReplyForm = ref(false);
const showEmojiPicker = ref(false);
const popularEmojis = ['👍', '❤️', '😂', '😮', '😢', '🔥'];
const hasReplies = computed(() => props.comment.replies && props.comment.replies.length > 0);
const REPLIES_PAGE_SIZE = 3;
const isRepliesExpanded = ref(false);
const visibleRepliesCount = ref(REPLIES_PAGE_SIZE);
const toggleReplies = () => { isRepliesExpanded.value = !isRepliesExpanded.value; if (!isRepliesExpanded.value) visibleRepliesCount.value = REPLIES_PAGE_SIZE; };
const loadMoreReplies = () => visibleRepliesCount.value += REPLIES_PAGE_SIZE;
const visibleReplies = computed(() => props.comment.replies ? props.comment.replies.slice(0, visibleRepliesCount.value) : []);
const hasMoreReplies = computed(() => props.comment.replies && visibleRepliesCount.value < props.comment.replies.length);
const remainingRepliesCount = computed(() => props.comment.replies ? props.comment.replies.length - visibleRepliesCount.value : 0);
const nextLoadCount = computed(() => Math.min(REPLIES_PAGE_SIZE, remainingRepliesCount.value));
const countAllReplies = (r) => !r || r.length === 0 ? 0 : r.length + r.reduce((sum, reply) => sum + countAllReplies(reply.replies), 0);
const totalRepliesCount = computed(() => countAllReplies(props.comment.replies));
const formatTimestamp = (iso) => new Date(iso).toLocaleString();
const formattedContent = computed(() => props.comment.content.replace(/^(@\S+\s)/, '<span class="mention">$1</span>'));
const addReaction = async (emoji) => {
  try {
    await api.post(`/writing_comments/${props.comment.id}/react`, { emoji });
    emit('comment-updated');
    showEmojiPicker.value = false;
  } catch (error) { console.error('添加表情失败:', error); }
};
const onReplyPosted = () => { showReplyForm.value = false; emit('comment-updated'); };
</script>

<style scoped>
/* 样式与 CommentItem.vue 相同，这里省略以保持简洁 */
.comment-item { display: flex; gap: 12px; }
.comment-gutter { display: flex; flex-direction: column; align-items: center; flex-shrink: 0; }
.avatar-container { width: 40px; height: 40px; border-radius: 50%; background-color: #f0f2f5; display: flex; justify-content: center; align-items: center; font-weight: bold; }
.thread-line { width: 2px; background-color: #e1e8ed; flex-grow: 1; margin-top: 8px; }
.comment-main { flex-grow: 1; }
.comment-wrapper { padding-bottom: 1rem; border-bottom: 1px solid #eee; }
.replies-container > .comment-wrapper:last-child { padding-bottom: 0; border-bottom: none; }
.comment-header { margin-bottom: 0.5rem; }
.timestamp { font-size: 0.8rem; color: #999; margin-left: 0.5rem; }
.comment-content p { margin: 0; white-space: pre-wrap; word-break: break-word; }
.comment-attachment { max-width: 100%; max-height: 250px; border-radius: 8px; margin-top: 0.75rem; }
.comment-actions { display: flex; justify-content: space-between; align-items: center; margin-top: 1rem; }
.reactions { position: relative; display: flex; align-items: center; gap: 0.5rem; }
.reaction-badge { background: #f0f0f0; padding: 0.25rem 0.5rem; border-radius: 1rem; cursor: pointer; }
.add-reaction-btn, .reply-btn { background: none; border: none; cursor: pointer; font-weight: 500; }
.emoji-picker { position: absolute; top: calc(100% + 8px); left: 0; background: white; border: 1px solid #ccc; border-radius: 8px; padding: 0.5rem; display: flex; gap: 0.5rem; box-shadow: 0 4px 8px rgba(0,0,0,0.1); z-index: 10; }
.emoji-picker span { cursor: pointer; font-size: 1.2rem; }
:deep(.mention) { color: #007bff; font-weight: 500; }
.replies-container { padding-left: 52px; margin-top: 1rem; }
/* ==================== 新增的修复规则 ==================== */
.replies-container .replies-container {
  padding-left: 0;
}
/* ====================================================== */
.toggle-replies-btn, .show-more-btn { background: none; border: none; color: #007bff; cursor: pointer; font-weight: 500; }
</style>