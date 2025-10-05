<!-- /frontend/src/views/GuestbookPage.vue -->
<template>
  <div class="guestbook-container">
    <div class="page-header">
      <h1>留言板</h1>
      <p>在这里留下你的足迹吧！</p>
    </div>
    <!-- 将原有的 Guestbook 组件内容直接放在这里 -->
    <Guestbook />
  </div>
</template>

<script setup>
import Guestbook from '../components/Guestbook.vue';
</script>

<style scoped>
.guestbook-container {
  max-width: 768px;
  margin: 2rem auto;
  padding: 2rem;
  background-color: #fff;
  border-radius: 8px;
  box-shadow: 0 4px 12px rgba(0,0,0,0.05);
}
.page-header {
  text-align: center;
  margin-bottom: 2.5rem;
  border-bottom: 1px solid #eee;
  padding-bottom: 1.5rem;
}
.page-header h1 {
  margin: 0;
  font-size: 2.5rem;
}
.page-header p {
  font-size: 1.1rem;
  color: #777;
  margin-top: 0.5rem;
}
</style>