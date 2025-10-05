<!-- /frontend/src/views/WritingsPage.vue (最终修复版) -->
<template>
  <div class="writings-container">
    <div class="page-header">
      <h1>我的写作</h1>
      <p>一些思考与记录的沉淀。</p>
      <div class="search-container">
        <input 
          type="text" 
          v-model="searchQuery" 
          placeholder="按文件名搜索..."
          class="search-box"
        />
      </div>
    </div>
    
    <!-- === 模板修改：从单个 v-for 变为两个并排的列 === -->
    <div class="writings-grid">
      <!-- 第一列 -->
      <div class="writing-column">
        <WritingCard v-for="writing in column1Writings" :key="writing.filename" :writing="writing" />
      </div>
      <!-- 第二列 -->
      <div class="writing-column">
        <WritingCard v-for="writing in column2Writings" :key="writing.filename" :writing="writing" />
      </div>
    </div>
    <!-- ====================================================== -->
    
    <div ref="loadMoreTrigger" class="load-more-trigger">
      <p v-if="isLoading">正在加载...</p>
      <p v-if="!hasMore && !isLoading && writings.length > 0">已加载全部内容</p>
      <p v-if="!isLoading && writings.length === 0 && hasLoadedOnce">没有找到匹配的文件</p>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, watch } from 'vue';
import api from '../services/api';
import WritingCard from '../components/WritingCard.vue';

// === Script 修改：使用两个数组来管理列数据 ===
const writings = ref([]); // 保留一个原始列表，用于判断是否加载完成
const column1Writings = ref([]);
const column2Writings = ref([]);
// ===============================================

const currentPage = ref(1);
const isLoading = ref(false);
const hasMore = ref(true);
const loadMoreTrigger = ref(null);
const searchQuery = ref('');
let debounceTimer = null;
const hasLoadedOnce = ref(false);

const fetchWritings = async () => {
  if (isLoading.value || !hasMore.value) return;
  isLoading.value = true;
  try {
    const response = await api.get('/writings', {
      params: { 
        page: currentPage.value, 
        pageSize: 10, // 可以适当增加每次获取的数量，以更好地分配
        search: searchQuery.value,
      }
    });

    if (response.data.writings && response.data.writings.length > 0) {
      // === Script 修改：将获取到的数据交替分配到两个列数组中 ===
      response.data.writings.forEach((writing, index) => {
        // 通过判断当前总条目数的奇偶性来决定放入哪一列
        if ((writings.value.length + index) % 2 === 0) {
          column1Writings.value.push(writing);
        } else {
          column2Writings.value.push(writing);
        }
      });
      // ========================================================
      
      writings.value.push(...response.data.writings);
      currentPage.value++;
    }
    hasMore.value = response.data.hasMore;
  } catch (error) {
    console.error('获取文章列表失败:', error);
  } finally {
    isLoading.value = false;
    hasLoadedOnce.value = true;
  }
};

watch(searchQuery, () => {
  clearTimeout(debounceTimer);
  debounceTimer = setTimeout(() => {
    // 重置所有状态以进行新的搜索
    currentPage.value = 1;
    writings.value = [];
    // === Script 修改：清空列数组 ===
    column1Writings.value = [];
    column2Writings.value = [];
    // ===============================
    hasMore.value = true;
    hasLoadedOnce.value = false;
    fetchWritings();
  }, 500);
});

let observer;
onMounted(() => {
  fetchWritings();
  
  observer = new IntersectionObserver((entries) => {
    if (entries[0].isIntersecting && hasMore.value) {
      fetchWritings();
    }
  }, { threshold: 1.0 });

  if (loadMoreTrigger.value) {
    observer.observe(loadMoreTrigger.value);
  }
});

onUnmounted(() => {
  if (observer && loadMoreTrigger.value) {
    observer.unobserve(loadMoreTrigger.value);
  }
  clearTimeout(debounceTimer);
});
</script>

<script>
export default {
  name: 'WritingsPage'
}
</script>

<style scoped>
.writings-container {
  max-width: 960px;
  margin: 2rem auto;
  padding: 2rem;
}
.page-header {
  text-align: center;
  margin-bottom: 3rem;
}
.search-container {
  margin-top: 1.5rem;
  display: flex;
  justify-content: center;
}
.search-box {
  width: 100%;
  max-width: 400px;
  padding: 0.75rem 1rem;
  font-size: 1rem;
  border: 1px solid #ddd;
  border-radius: 2rem;
  outline: none;
  transition: border-color 0.2s, box-shadow 0.2s;
}
.search-box:focus {
  border-color: #86b93a;
  box-shadow: 0 0 0 3px rgba(134, 185, 58, 0.2);
}

/* === 样式修改：将 Grid 布局改为 Flexbox 布局 === */
.writings-grid {
  display: flex;
  gap: 1.5rem; /* 控制两列之间的间距 */
  align-items: flex-start; /* 关键：让列从顶部对齐 */
}

.writing-column {
  flex: 1; /* 让两列平分宽度 */
  min-width: 0; /* 防止内容溢出 */
  display: flex;
  flex-direction: column;
  gap: 1.5rem; /* 控制列内部卡片的垂直间距 */
}
/* ====================================================== */

.load-more-trigger {
  height: 50px;
  text-align: center;
  color: #888;
}

@media (max-width: 768px) {
  .writings-grid {
    /* 将 flex 布局改为垂直方向，自然形成单列 */
    flex-direction: column;
    gap: 0; /* 移除列间距 */
  }

  /* 因为现在是单列，所以不需要 writing-column 这个容器了，
     但为了兼容现有结构，我们让它占满宽度 */
  .writing-column {
    width: 100%;
  }
}
</style>