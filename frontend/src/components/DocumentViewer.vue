<!-- /frontend/src/components/DocumentViewer.vue (已修复移动端显示问题) -->
<template>
  <div class="viewer-container">
    <div v-if="isLoading" class="status-message">正在加载文档...</div>
    <div v-else-if="error" class="status-message error">{{ error }}</div>
    <template v-else>
      <div v-if="fileType === 'pdf' && pdfSource" class="pdf-viewer">
        <vue-pdf-embed :source="pdfSource" />
      </div>
      <div v-else-if="fileType === 'docx'" ref="docxContainer" class="docx-viewer"></div>
      <div v-else-if="fileType === 'md'" v-html="markdownContent" class="markdown-viewer"></div>
      <div v-else class="status-message">不支持预览此文件类型。</div>
    </template>
  </div>
</template>

<script setup>
import { ref, computed, watch, nextTick } from 'vue';
import VuePdfEmbed from 'vue-pdf-embed';
import { marked } from 'marked';
import { renderAsync } from 'docx-preview';
import axios from 'axios';

const props = defineProps({
  url: { type: String, required: true }
});

const isLoading = ref(true);
const error = ref(null);
const docxContainer = ref(null);
const markdownContent = ref('');
const pdfSource = ref(null);

const fileType = computed(() => {
  if (!props.url) return 'unsupported';
  const parts = props.url.split('.');
  const ext = parts.pop().toLowerCase();
  if (ext === 'md' || ext === 'markdown') return 'md';
  if (ext === 'docx') return 'docx';
  if (ext === 'pdf') return 'pdf';
  return 'unsupported';
});

const loadDocument = async () => {
  isLoading.value = true;
  error.value = null;
  markdownContent.value = '';
  pdfSource.value = null;

  try {
    if (fileType.value === 'pdf') {
      pdfSource.value = props.url;
      isLoading.value = false;
    } else if (fileType.value === 'md') {
      const response = await axios.get(props.url);
      markdownContent.value = marked(response.data);
      isLoading.value = false;
    } else if (fileType.value === 'docx') {
      const response = await axios.get(props.url, { responseType: 'blob' });
      const docxBlob = response.data;
      isLoading.value = false;
      await nextTick();
      if (docxContainer.value) {
        await renderAsync(docxBlob, docxContainer.value);
      } else {
        throw new Error("DOCX container is still not available after DOM update.");
      }
    }
  } catch (e) {
    console.error("加载或渲染文档失败:", e);
    error.value = "无法加载或渲染文档，请检查文件是否存在或稍后重试。";
    isLoading.value = false;
  }
};

watch(() => props.url, loadDocument, { immediate: true });
</script>

<style>
/* --- 核心修改 --- */
.viewer-container {
  width: 100%;
  height: 100%;
  background-color: #f8f9fa;
  overflow-y: auto;
  /* 1. 允许横向滚动以查看溢出内容 */
  overflow-x: auto; 
  padding: 2rem;
  box-sizing: border-box;
}

/* 2. 针对移动端优化内边距 */
@media (max-width: 768px) {
  .viewer-container {
    padding: 1rem 0.5rem; /* 减小左右内边距 */
  }
  .docx-viewer .docx-wrapper {
    padding: 1rem !important; /* 同样减小 DOCX 预览的内边距 */
  }
}
/* --- 修改结束 --- */

.status-message { text-align: center; padding-top: 4rem; color: #666; }
.error { color: red; }
.pdf-viewer .vue-pdf-embed > div {
  margin-bottom: 1rem;
  box-shadow: 0 4px 12px rgba(0,0,0,0.1);
  display: flex;
  justify-content: center;
}
.pdf-viewer .vue-pdf-embed canvas {
  transform-origin: top center;
}
.markdown-viewer { line-height: 1.7; }
.markdown-viewer h1, .markdown-viewer h2, .markdown-viewer h3 { border-bottom: 1px solid #eee; padding-bottom: 0.5rem; margin-top: 2rem; }
.markdown-viewer pre { background-color: #2d2d2d; color: #f8f8f2; padding: 1rem; border-radius: 8px; overflow-x: auto; }
.markdown-viewer code { font-family: 'Courier New', Courier, monospace; }
.markdown-viewer blockquote { border-left: 4px solid #ccc; padding-left: 1rem; color: #666; }
.docx-viewer .docx-wrapper { background-color: #fff !important; padding: 2rem !important; }
.docx-viewer section.docx { box-shadow: 0 4px 12px rgba(0,0,0,0.1) !important; }
</style>