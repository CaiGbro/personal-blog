<!-- /frontend/src/components/CommentForm.vue (完全替换) -->
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
    <!-- 引用 ref="textareaRef" 以便聚焦 -->
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
  initialContent: { type: String, default: '' }, // 新增：接收初始内容
});

const emit = defineEmits(['comment-posted']);

const isSubmitting = ref(false);
const comment = ref({
  nickname: '',
  content: props.initialContent, // 修改：用 prop 初始化 content
  parentId: props.parentId,
  attachmentUrl: null,
});
const uploadPreview = ref(null);
const uploadError = ref(null);
const textareaRef = ref(null); // 新增：模板引用

// 新增：当组件挂载（即回复表单显示时），自动聚焦到输入框
onMounted(() => {
    if(textareaRef.value) {
        textareaRef.value.focus();
    }
});

const handleFileUpload = async (event) => {
  const file = event.target.files[0];
  if (!file) return;

  if (file.size > 5 * 1024 * 1024) {
    uploadError.value = '文件大小不能超过 5MB。';
    return;
  }

  const formData = new FormData();
  formData.append('file', file);
  isSubmitting.value = true;
  uploadError.value = null;

  try {
    const response = await api.post('/upload', formData, {
      headers: { 'Content-Type': 'multipart/form-data' },
    });
    comment.value.attachmentUrl = response.data.url;
    uploadPreview.value = response.data.url;
  } catch (error) {
    console.error('上传失败:', error);
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
    await api.post('/comments', comment.value);
    emit('comment-posted');
    // 重置时不再清空 nickname
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
/* 样式无需改动 */
.comment-form { display: flex; flex-direction: column; gap: 0.75rem; margin-bottom: 2rem; margin-top: 1rem; }
.form-row { display: flex; gap: 0.5rem; align-items: center; }
.form-row input { flex-grow: 1; padding: 0.5rem; border: 1px solid #ccc; border-radius: 4px; }
textarea { padding: 0.5rem; border: 1px solid #ccc; border-radius: 4px; min-height: 60px; }
button { padding: 0.5rem 1rem; background-color: #86b93a; color: white; border: none; border-radius: 4px; cursor: pointer; }
button:disabled { background-color: #ccc; }
.file-upload-btn { position: relative; cursor: pointer; font-size: 1.5rem; }
.file-upload-btn input[type="file"] { display: none; }
.upload-preview { position: relative; max-width: 150px; }
.upload-preview img { width: 100%; border-radius: 4px; }
.remove-btn { position: absolute; top: -5px; right: -5px; background: #333; color: white; border-radius: 50%; width: 20px; height: 20px; border: none; font-weight: bold; cursor: pointer; }
.error-message { color: red; font-size: 0.9rem; }
</style>