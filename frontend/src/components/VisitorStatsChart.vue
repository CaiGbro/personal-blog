<!-- /frontend/src/components/VisitorStatsChart.vue (新增文件) -->
<template>
  <div class="chart-container">
    <div v-if="isLoading" class="status-message">正在加载统计数据...</div>
    <div v-else-if="error" class="status-message error">{{ error }}</div>
    <Line v-else :data="chartData" :options="chartOptions" />
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { Line } from 'vue-chartjs';
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend,
  Filler, // 导入 Filler 以支持区域填充
} from 'chart.js';
import api from '../services/api';

// 注册 Chart.js 所需的组件
ChartJS.register(
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend,
  Filler
);

const isLoading = ref(true);
const error = ref(null);
const chartData = ref({
  labels: [],
  datasets: [],
});

const chartOptions = ref({
  responsive: true,
  maintainAspectRatio: false,
  plugins: {
    legend: {
      display: false, // 不显示图例
    },
    tooltip: {
      mode: 'index',
      intersect: false,
    },
  },
  scales: {
    x: {
      grid: {
        display: false, // 隐藏 x 轴网格线
      },
    },
    y: {
      beginAtZero: true,
      grid: {
        color: '#f0f0f0', // y 轴网格线颜色
      },
    },
  },
  interaction: {
    mode: 'nearest',
    axis: 'x',
    intersect: false,
  },
  elements: {
    line: {
      tension: 0.4, // 使线条更平滑
    },
  },
});

onMounted(async () => {
  try {
    const response = await api.get('/admin/stats/visitors');
    const stats = response.data;

    if (!stats || stats.length === 0) {
      error.value = '暂无访客数据。';
      return;
    }

    chartData.value = {
      labels: stats.map(s => s.date),
      datasets: [
        {
          label: '独立访客数',
          backgroundColor: (context) => {
            const ctx = context.chart.ctx;
            const gradient = ctx.createLinearGradient(0, 0, 0, 200);
            gradient.addColorStop(0, 'rgba(134, 185, 58, 0.5)');
            gradient.addColorStop(1, 'rgba(134, 185, 58, 0)');
            return gradient;
          },
          borderColor: '#86b93a',
          pointBackgroundColor: '#86b93a',
          pointBorderColor: '#fff',
          pointHoverBackgroundColor: '#fff',
          pointHoverBorderColor: '#86b93a',
          data: stats.map(s => s.count),
          fill: true, // 开启区域填充
        },
      ],
    };
  } catch (e) {
    console.error('获取访客数据失败:', e);
    error.value = '无法加载统计数据，请稍后重试。';
  } finally {
    isLoading.value = false;
  }
});
</script>

<style scoped>
.chart-container {
  position: relative;
  height: 300px; /* 给图表一个固定的高度 */
  padding: 1rem;
  background-color: #fafafa;
  border-radius: 8px;
}
.status-message {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100%;
  color: #777;
}
.error {
  color: #d9534f;
}
</style>