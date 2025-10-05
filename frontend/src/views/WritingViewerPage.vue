<!-- /frontend/src/views/WritingViewerPage.vue (支持多格式下载) -->
<template>
  <div class="viewer-page-container">
    <main class="document-pane">
      <DocumentViewer v-if="documentUrlToDisplay" :url="documentUrlToDisplay" />
      
      <!-- 控制按钮 -->
      <div class="side-controls">
        <button class="control-btn" @click="addReaction('❤️')">
          <span class="icon">❤️</span>
          <span class="count">{{ reactionCount }}</span>
        </button>
        <button class="control-btn" @click="showComments = !showComments">
          <span class="icon">💬</span>
          <span class="count">{{ commentCount }}</span>
        </button>

        <!-- ==================== 下载按钮区域 ==================== -->
        <div class="download-control">
          <!-- Case 1: 如果存在多种格式，点击按钮会展开选项 -->
          <button v-if="pdfFileExists && originalFileExists" class="control-btn" @click="showDownloadOptions = !showDownloadOptions">
            <span class="icon">📥</span>
            <span class="count">下载</span>
          </button>
          <!-- Case 2: 如果只存在一种格式，直接下载 -->
          <button v-else class="control-btn" @click="downloadSingleFile">
            <span class="icon">📥</span>
            <span class="count">下载</span>
          </button>

          <!-- 下载选项浮窗 -->
          <transition name="fade">
            <div v-if="showDownloadOptions" class="download-options">
              <a :href="pdfFileUrl" :download="pdfFilename" class="option-item">下载 PDF</a>
              <a :href="originalFileUrl" :download="originalFilename" class="option-item">下载 {{ originalFileExtension.toUpperCase() }}</a>
            </div>
          </transition>
        </div>
        <!-- ======================================================= -->

      </div>
    </main>

    <Transition name="slide-fade">
      <aside v-if="showComments" class="comment-pane">
        <WritingCommentSection 
          :writing-filename="filename" 
          @count-updated="updateCommentCount"
        />
      </aside>
    </Transition>
  </div>
</template>

<script setup>
import { ref, computed, watch } from 'vue';
import axios from 'axios';
import api from '../services/api';
import DocumentViewer from '../components/DocumentViewer.vue';
import WritingCommentSection from '../components/WritingCommentSection.vue';

const props = defineProps({
  filename: { type: String, required: true }
});

const documentUrlToDisplay = ref('');
const showComments = ref(false);
const commentCount = ref(0);
const reactions = ref({}); 
const reactionCount = computed(() => reactions.value['❤️'] || 0);

// --- 下载相关状态 ---
const showDownloadOptions = ref(false);
const originalFileExists = ref(false);
const pdfFileExists = ref(false);

// --- 计算各种文件格式的名称和URL ---
const baseFilename = computed(() => props.filename.replace(/\.(docx|md|pdf)$/i, ''));
const originalFileExtension = computed(() => {
    if (!props.filename.toLowerCase().endsWith('.pdf')) {
        return props.filename.split('.').pop() || 'docx';
    }
    // 默认猜测原始格式为 docx，您可以根据您的主要文件类型调整
    return 'docx'; 
});
const originalFilename = computed(() => `${baseFilename.value}.${originalFileExtension.value}`);
const pdfFilename = computed(() => `${baseFilename.value}.pdf`);

const originalFileUrl = computed(() => `/static/writings/${originalFilename.value}`);
const pdfFileUrl = computed(() => `/static/writings/${pdfFilename.value}`);


/** 检查文件是否存在 */
const checkFileExistence = async (url) => {
  try {
    await axios.head(url);
    return true;
  } catch (error) {
    return false;
  }
};

/** 检查所有可用的文件版本 */
const checkAllFileVersions = async () => {
    // 并行检查两个文件是否存在，速度更快
    const [pdfExists, originalExists] = await Promise.all([
        checkFileExistence(pdfFileUrl.value),
        checkFileExistence(originalFileUrl.value)
    ]);
    pdfFileExists.value = pdfExists;
    originalFileExists.value = originalExists;

    // 根据检查结果，决定显示哪个文件（优先PDF）
    if (pdfExists) {
        documentUrlToDisplay.value = pdfFileUrl.value;
    } else if (originalExists) {
        documentUrlToDisplay.value = originalFileUrl.value;
    } else {
        // 如果两个都不存在（异常情况），显示原始URL让其404
        documentUrlToDisplay.value = originalFileUrl.value;
    }
};

/** 下载唯一存在的文件 */
const downloadSingleFile = () => {
    const url = pdfFileExists.value ? pdfFileUrl.value : originalFileUrl.value;
    const filename = pdfFileExists.value ? pdfFilename.value : originalFilename.value;
    
    const link = document.createElement('a');
    link.href = url;
    link.download = filename;
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
};

const updateCommentCount = (newCount) => {
  commentCount.value = newCount;
};

const fetchReactions = async (filename) => {
  if (!filename) return;
  const apiUrl = `/writings/reactions/static/writings/${filename}`; 
  try {
    const response = await api.get(apiUrl);
    reactions.value = response.data || {};
  } catch (error) {
    console.error("获取文章回应失败:", error);
    reactions.value = {};
  }
};

const addReaction = async (emoji) => {
  if (!props.filename) return;
  const apiUrl = `/writings/react/static/writings/${props.filename}`;
  try {
    const response = await api.post(apiUrl, { emoji });
    reactions.value = response.data;
  } catch (error) {
    console.error("添加文章回应失败:", error);
  }
};

const fetchCommentCount = async (filename) => {
  if (!filename) return;
  try {
    const response = await api.get('/writing_comments/count', {
      params: { writing_filename: filename }
    });
    commentCount.value = response.data.totalCount || 0;
  } catch (error) {
    console.error("获取初始评论数失败:", error);
    commentCount.value = 0;
  }
};

// 监听文件名变化，重新执行所有数据获取和URL检查
watch(() => props.filename, (newFilename) => {
  if(newFilename) {
    documentUrlToDisplay.value = '';
    showDownloadOptions.value = false; // 切换文件时关闭下载选项
    checkAllFileVersions(); // 检查所有文件版本
    fetchReactions(newFilename);
    fetchCommentCount(newFilename);
  }
}, { immediate: true });

</script>

<style scoped>
.viewer-page-container {
  display: flex;
  height: calc(100vh - 60px);
  width: 100vw;
  position: fixed;
  top: 60px;
  left: 0;
  overflow: hidden;
}

.document-pane {
  flex-grow: 1;
  height: 100%;
  position: relative;
}

.comment-pane {
  flex-shrink: 0;
  width: 400px;
  max-width: 450px;
  height: 100%;
  border-left: 1px solid #e0e0e0;
  box-shadow: -4px 0 12px rgba(0,0,0,0.05);
}

.side-controls {
  position: absolute;
  bottom: 30px;
  right: 30px;
  z-index: 10;
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.control-btn {
  background: rgba(255, 255, 255, 0.8);
  backdrop-filter: blur(5px);
  border: 1px solid #eee;
  color: #333;
  cursor: pointer;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 0.3rem;
  font-size: 0.9rem;
  padding: 0.75rem;
  border-radius: 50%;
  width: 60px;
  height: 60px;
  justify-content: center;
  box-shadow: 0 4px 12px rgba(0,0,0,0.1);
  transition: all 0.2s;
}

.control-btn:hover {
  transform: scale(1.1);
  box-shadow: 0 6px 16px rgba(0,0,0,0.12);
}

.control-btn .icon {
  font-size: 1.5rem;
}

.control-btn .count {
  font-weight: 500;
  font-size: 0.8rem;
}

.slide-fade-enter-active, .slide-fade-leave-active {
  transition: all 0.3s ease-out;
}

.slide-fade-enter-from, .slide-fade-leave-to {
  transform: translateX(100%);
  opacity: 0;
}

.download-control {
  position: relative;
}

.download-options {
  position: absolute;
  bottom: 75px;
  right: 0;
  background-color: white;
  border-radius: 8px;
  box-shadow: 0 6px 20px rgba(0,0,0,0.15);
  padding: 0.5rem;
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
  width: 150px;
  z-index: 20;
}

.option-item {
  padding: 0.75rem 1rem;
  text-align: center;
  text-decoration: none;
  color: #333;
  border-radius: 6px;
  transition: background-color 0.2s;
}

.option-item:hover {
  background-color: #f0f0f0;
}

.fade-enter-active, .fade-leave-active {
  transition: opacity 0.2s ease, transform 0.2s ease;
}

.fade-enter-from, .fade-leave-to {
  opacity: 0;
  transform: translateY(10px);
}

@media (max-width: 768px) {
  .viewer-page-container {
    flex-direction: column;
    height: auto;
    position: static;
    overflow-y: auto;
  }

  .document-pane {
    height: auto;
  }

  .comment-pane {
    width: 100%;
    max-width: 100%;
    height: auto;
    border-left: none;
    border-top: 1px solid #e0e0e0;
  }

  .side-controls {
    bottom: 15px;
    right: 15px;
  }
  
  .download-options {
    bottom: 70px;
    right: 0;
  }
}
</style>