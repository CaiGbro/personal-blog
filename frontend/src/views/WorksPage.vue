<!-- /frontend/src/views/WorksPage.vue (已修复样式Bug) -->
<template>
  <div class="works-container">
    <div class="page-header">
      <h1>探索发现</h1>
      <p>一些技术与创意的探索。</p>
    </div>
    <div class="projects-grid">
      <!-- 项目一：美与欣赏 -->
      <router-link to="/works/video-gallery" class="project-card">
        <div class="card-content">
          <div class="card-icon">🎬</div>
          <h3>美与欣赏</h3>
          <p>沉浸式短视频浏览器，探索光影之美。</p>
        </div>
      </router-link>

      <!-- 项目二：新闻热点 -->
      <div class="project-card placeholder">
        <div class="card-content">
          <div class="card-icon">📰</div>
          <h3>智能情报</h3>
          <p>AI赋能，实时热点，深度发掘、长效记忆、自我进化。
            沉淀深刻道理、预测发展趋势。</p>
          <span>(开发中)</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
</script>

<style scoped>
.works-container {
  max-width: 960px;
  margin: 2rem auto;
  padding: 2rem;
}
.page-header {
  text-align: center;
  margin-bottom: 3rem;
}
.projects-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
  gap: 2rem;
}
.project-card {
  background-color: #fff;
  border-radius: 12px;
  box-shadow: 0 4px 12px rgba(0,0,0,0.05);
  padding: 2rem;
  text-align: center;
  transition: transform 0.3s ease, box-shadow 0.3s ease;
  border: 1px solid #f0f0f0;

  /* ==================== 【核心修复】 ==================== */
  /* 强制重置链接的默认样式，确保视觉统一 */
  color: inherit; 
  text-decoration: none;
  /* ====================================================== */
}
.project-card:not(.placeholder):hover {
  transform: translateY(-5px);
  box-shadow: 0 8px 24px rgba(0,0,0,0.08);
}
.card-content h3 {
  margin: 1rem 0 0.5rem;
  font-size: 1.5rem;
  /* 明确指定标题颜色，防止被 a 标签影响 */
  color: #333; 
}
.card-content p {
  color: #666;
  line-height: 1.6;
}
.card-icon {
  font-size: 3rem;
}
.placeholder {
  background-color: #f9f9f9;
  cursor: not-allowed;
  opacity: 0.7;
}
.placeholder span {
  display: inline-block;
  margin-top: 1rem;
  font-size: 0.9rem;
  color: #999;
  font-weight: bold;
}
</style>