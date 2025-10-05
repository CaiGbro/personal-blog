<!-- /frontend/src/components/WritingCommentForm.vue -->
<template>
  <form @submit.prevent="submit" class="comment-form">
    <div class="form-row">
      <input v-model="comment.nickname" type="text" placeholder="你的昵称" required />
      <label class="file-upload-btn">
        📷
        <input type="file" @change="handleFileUpload" accept="image/*" />
      </label>
      <button type="submit" :disabled="isSubmitting">
        {{ isSubmitting ? '...' : '提交' }}
      </button>
    </div>
    <textarea ref="textareaRef" v-model="comment.content" placeholder="留下你的评论..." required></textarea>
    <div v-if="uploadPreview" class="upload-preview">
      <img :src="uploadPreview" alt="Upload preview" />
      <button @click="removeAttachment" type="button" class="remove-btn">&times;</button>
    </div>
    <div v-if="uploadError" class="error-message">{{ uploadError }}</div>
  </form>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import api from '../services/api';

const props = defineProps({
  parentId: { type: Number, default: null },
  initialContent: { type: String, default: '' },
  writingFilename: { type: String, required: true },
});
const emit = defineEmits(['comment-posted']);
const isSubmitting = ref(false);
const comment = ref({
  nickname: '',
  content: props.initialContent,
  parentId: props.parentId,
  attachmentUrl: null,
  writingFilename: props.writingFilename,
});
const uploadPreview = ref(null);
const uploadError = ref(null);
const textareaRef = ref(null);

onMounted(() => { if(textareaRef.value) { textareaRef.value.focus(); } });

const handleFileUpload = async (event) => {
  const file = event.target.files[0];
  if (!file) return;
  if (file.size > 5 * 1024 * 1024) { uploadError.value = '文件大小不能超过 5MB。'; return; }
  const formData = new FormData();
  formData.append('file', file);
  isSubmitting.value = true;
  uploadError.value = null;
  try {
    const response = await api.post('/upload/writing_comment', formData);
    comment.value.attachmentUrl = response.data.url;
    uploadPreview.value = response.data.url;
  } catch (error) {
    uploadError.value = '图片上传失败，请重试。';
  } finally {
    isSubmitting.value = false;
    event.target.value = '';
  }
};

const removeAttachment = () => {
    comment.value.attachmentUrl = null;
    uploadPreview.value = null;
}

const submit = async () => {
  isSubmitting.value = true;
  try {
    await api.post('/writing_comments', comment.value);
    emit('comment-posted');
    comment.value.content = '';
    comment.value.attachmentUrl = null;
    uploadPreview.value = null;
  } catch (error) {
    console.error('提交评论失败:', error);
  } finally {
    isSubmitting.value = false;
  }
};
</script>

<style scoped>
/* 样式与 CommentForm.vue 相同，这里省略以保持简洁 */
.comment-form { display: flex; flex-direction: column; gap: 0.75rem; padding: 1rem 0; }
.form-row { display: flex; gap: 0.5rem; align-items: center; }
input, textarea { background-color: #fff; color: #333; border: 1px solid #ccc; border-radius: 4px; padding: 0.75rem; width: 100%; box-sizing: border-box; }
input { flex-grow: 1; }
textarea { min-height: 80px; resize: vertical; }
button { padding: 0.75rem 1rem; background-color: #86b93a; color: white; border: none; border-radius: 4px; cursor: pointer; white-space: nowrap; }
button:disabled { background-color: #ccc; cursor: not-allowed; }
.file-upload-btn { position: relative; cursor: pointer; font-size: 1.5rem; padding: 0 0.5rem; }
.file-upload-btn input[type="file"] { display: none; }
.upload-preview { position: relative; max-width: 100px; margin-top: 0.5rem; }
.upload-preview img { width: 100%; border-radius: 4px; }
.remove-btn { position: absolute; top: -8px; right: -8px; background: rgba(0,0,0,0.7); color: white; border-radius: 50%; width: 22px; height: 22px; border: none; font-weight: bold; cursor: pointer; }
.error-message { color: #ff4d4d; font-size: 0.9rem; }
</style>