<!-- /frontend/src/views/AdminPage.vue (修改此文件) -->
<template>
  <div class="admin-container">
    <h1>后台管理</h1>

    <!-- 新增：访客统计卡片 -->
    <div class="settings-card">
      <h3>访客统计 (近30天)</h3>
      <VisitorStatsChart />
    </div>

    <!-- 现有：系统设置卡片 -->
    <div class="settings-card">
      <h3>系统设置</h3>
      <div class="setting-item">
        <div class="setting-label">
          <strong>游客模式</strong>
          <p>开启后，所有人都可以访问网站内容。关闭后，只有管理员可以登录和访问。</p>
        </div>
        <div class="setting-control">
          <label class="switch">
            <input type="checkbox" v-model="visitorModeEnabled" @change="updateSettings" />
            <span class="slider round"></span>
          </label>
        </div>
      </div>
      <p v-if="message" class="message">{{ message }}</p>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import api from '../services/api';
// --- 新增：导入图表组件 ---
import VisitorStatsChart from '../components/VisitorStatsChart.vue';

const visitorModeEnabled = ref(true);
const message = ref('');

const fetchSettings = async () => {
  try {
    const response = await api.get('/admin/settings');
    visitorModeEnabled.value = response.data.visitor_mode_enabled;
  } catch (error) {
    console.error('获取设置失败:', error);
    message.value = '无法加载当前设置。';
  }
};

const updateSettings = async () => {
  try {
    message.value = '正在更新...';
    const response = await api.post('/admin/settings', {
      visitor_mode_enabled: visitorModeEnabled.value,
    });
    message.value = response.data.message;
    setTimeout(() => message.value = '', 3000);
  } catch (error) {
    console.error('更新设置失败:', error);
    message.value = '更新失败，请重试。';
    visitorModeEnabled.value = !visitorModeEnabled.value;
  }
};

onMounted(fetchSettings);
</script>

<style scoped>
.admin-container {
  max-width: 768px;
  margin: 2rem auto;
  padding: 2rem;
  display: flex; /* 使用 flex 布局 */
  flex-direction: column; /* 垂直排列 */
  gap: 2rem; /* 设置卡片之间的间距 */
}
h1 { text-align: center; margin-bottom: 2rem; }
.settings-card { background: white; border-radius: 8px; padding: 2rem; box-shadow: 0 4px 12px rgba(0,0,0,0.05); }
.setting-item { display: flex; justify-content: space-between; align-items: center; padding: 1.5rem 0; border-bottom: 1px solid #f0f0f0; }
.setting-item:last-child { border-bottom: none; }
.setting-label p { color: #777; font-size: 0.9rem; margin-top: 0.25rem; }
.message { margin-top: 1rem; color: #86b93a; text-align: center; }

/* 开关样式 */
.switch { position: relative; display: inline-block; width: 60px; height: 34px; }
.switch input { opacity: 0; width: 0; height: 0; }
.slider { position: absolute; cursor: pointer; top: 0; left: 0; right: 0; bottom: 0; background-color: #ccc; transition: .4s; }
.slider:before { position: absolute; content: ""; height: 26px; width: 26px; left: 4px; bottom: 4px; background-color: white; transition: .4s; }
input:checked + .slider { background-color: #86b93a; }
input:focus + .slider { box-shadow: 0 0 1px #86b93a; }
input:checked + .slider:before { transform: translateX(26px); }
.slider.round { border-radius: 34px; }
.slider.round:before { border-radius: 50%; }
</style>