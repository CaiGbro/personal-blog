<!-- /frontend/src/components/VideoCommentForm.vue -->
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
    <textarea ref="textareaRef" v-model="comment.content" placeholder="留下你的足迹..." required></textarea>
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
  videoFilename: { type: String, required: true },
});

const emit = defineEmits(['comment-posted']);

const isSubmitting = ref(false);
const comment = ref({
  nickname: '',
  content: props.initialContent,
  parentId: props.parentId,
  attachmentUrl: null,
  videoFilename: props.videoFilename, // 关键：将 prop 存入表单数据
});
const uploadPreview = ref(null);
const uploadError = ref(null);
const textareaRef = ref(null); // 用于自动聚焦

// 当组件挂载时（尤其是作为回复表单出现时），自动聚焦输入框
onMounted(() => {
    if(textareaRef.value) {
        textareaRef.value.focus();
    }
});

const handleFileUpload = async (event) => {
  const file = event.target.files[0];
  if (!file) return;

  if (file.size > 5 * 1024 * 1024) { // 限制文件大小为 5MB
    uploadError.value = '文件大小不能超过 5MB。';
    return;
  }

  const formData = new FormData();
  formData.append('file', file);
  isSubmitting.value = true;
  uploadError.value = null;

  try {
    // 调用为视频评论专设的上传接口
    const response = await api.post('/upload/video_comment', formData, {
      headers: { 'Content-Type': 'multipart/form-data' },
    });
    comment.value.attachmentUrl = response.data.url;
    uploadPreview.value = response.data.url;
  } catch (error) {
    console.error('上传失败:', error);
    uploadError.value = '图片上传失败，请重试。';
  } finally {
    isSubmitting.value = false;
    // 清空文件输入，以便用户可以重新上传同一个文件
    event.target.value = '';
  }
};

const removeAttachment = () => {
    comment.value.attachmentUrl = null;
    uploadPreview.value = null;
}

const submit = async () => {
  if (isSubmitting.value || !comment.value.nickname.trim() || !comment.value.content.trim()) {
      return;
  }
  isSubmitting.value = true;
  try {
    // 调用视频评论的提交接口
    await api.post('/video_comments', comment.value);
    emit('comment-posted');
    
    // 提交成功后重置内容和附件，但保留昵称
    comment.value.content = '';
    comment.value.attachmentUrl = null;
    uploadPreview.value = null;
    uploadError.value = null;

  } catch (error) {
    console.error('提交视频评论失败:', error);
    // 可以在这里设置一个提交失败的提示
  } finally {
    isSubmitting.value = false;
  }
};
</script>

<style scoped>
/* 这些样式是为视频播放器的深色主题设计的 */
.comment-form { 
  display: flex; 
  flex-direction: column; 
  gap: 0.75rem; 
  padding: 1rem 0; 
  border-top: 1px solid #333; 
}
.form-row { 
  display: flex; 
  gap: 0.5rem; 
  align-items: center; 
}
input, textarea { 
  background-color: #222; 
  color: #fff; 
  border: 1px solid #444; 
  border-radius: 4px; 
  padding: 0.75rem;
  width: 100%;
  box-sizing: border-box; /* 确保 padding 不会撑大宽度 */
}
input {
  flex-grow: 1; /* 让昵称输入框填满剩余空间 */
}
textarea { 
  min-height: 80px; 
  resize: vertical; /* 允许用户垂直调整大小 */
}
button { 
  padding: 0.75rem 1rem; 
  background-color: #86b93a; 
  color: white; 
  border: none; 
  border-radius: 4px; 
  cursor: pointer;
  white-space: nowrap; /* 防止按钮文字换行 */
  transition: background-color 0.2s;
}
button:hover:not(:disabled) {
    background-color: #76a333;
}
button:disabled { 
  background-color: #555;
  cursor: not-allowed;
}
.file-upload-btn { 
  position: relative; 
  cursor: pointer; 
  font-size: 1.5rem; 
  padding: 0 0.5rem;
}
.file-upload-btn input[type="file"] { 
  display: none; 
}
.upload-preview { 
  position: relative; 
  max-width: 100px;
  margin-top: 0.5rem;
}
.upload-preview img { 
  width: 100%; 
  border-radius: 4px; 
}
.remove-btn { 
  position: absolute; 
  top: -8px; 
  right: -8px; 
  background: rgba(0,0,0,0.7); 
  color: white; 
  border-radius: 50%; 
  width: 22px; 
  height: 22px; 
  border: none; 
  font-weight: bold; 
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 0;
  line-height: 1;
}
.error-message { 
  color: #ff4d4d; 
  font-size: 0.9rem; 
}
</style>