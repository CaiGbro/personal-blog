<!-- /frontend/src/views/LoginPage.vue (已更新) -->
<template>
  <div class="login-container">
    <div class="login-box">
      <!-- 1. 维护模式下的头部信息 -->
      <div v-if="isMaintenanceMode" class="maintenance-header">
        <h2>功能维护中</h2>
        <!-- ==================== 文本已修改 ==================== -->
        <p>游客功能暂时关闭，当前仅开发者访问，敬请谅解。</p>
        <!-- ====================================================== -->
      </div>

      <!-- 2. 只有在非维护模式，或用户点击了“开发者访问”后，才显示“管理员登录”标题 -->
      <h2 v-if="!isMaintenanceMode || showLoginForm">管理员登录</h2>
      
      <!-- 3. 新增的“开发者访问”按钮，只在维护模式且表单未显示时出现 -->
      <div v-if="isMaintenanceMode && !showLoginForm" class="developer-access">
        <button @click="showLoginForm = true" class="developer-btn">
          开发者访问
        </button>
      </div>

      <!-- 4. 登录表单，现在由 showLoginForm 控制 -->
      <form v-if="showLoginForm" @submit.prevent="handleLogin">
        <div class="form-group">
          <label for="email">邮箱</label>
          <!-- ==================== 移除了 placeholder 默认值 ==================== -->
          <input type="email" v-model="email" id="email" required placeholder="请输入管理员邮箱" />
          <!-- ==================================================================== -->
        </div>
        <div class="form-group">
          <label for="code">验证码</label>
          <div class="code-input-group">
            <input type="text" v-model="code" id="code" required />
            <button type="button" @click="sendCode" :disabled="isSendingCode || countdown > 0">
              {{ countdown > 0 ? `${countdown}s` : '发送验证码' }}
            </button>
          </div>
        </div>
        <p v-if="message" class="message" :class="{ 'error': isError }">{{ message }}</p>
        <button type="submit" class="login-btn" :disabled="isLoggingIn">
          {{ isLoggingIn ? '登录中...' : '登录' }}
        </button>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import api from '../services/api';
import { useAuth } from '../services/auth';

const router = useRouter();
const route = useRoute();
const { login } = useAuth();

// ==================== 移除了 email 的默认值 ====================
const email = ref('');
const code = ref('');
// ===============================================================

const isSendingCode = ref(false);
const isLoggingIn = ref(false);
const countdown = ref(0);
const message = ref('');
const isError = ref(false);
const isMaintenanceMode = ref(false);

// ==================== 新增状态，控制表单显示 ====================
const showLoginForm = ref(false);
// ===============================================================

onMounted(() => {
    if (route.query.mode === 'maintenance') {
        isMaintenanceMode.value = true;
    } else {
        // 如果不是维护模式，则直接显示登录框
        isMaintenanceMode.value = false;
        showLoginForm.value = true;
    }
});

const sendCode = async () => {
  if (!email.value) {
    message.value = '请输入管理员邮箱。';
    isError.value = true;
    return;
  }
  isSendingCode.value = true;
  message.value = '正在发送...';
  isError.value = false;
  try {
    const response = await api.post('/auth/send-code', { email: email.value });
    message.value = response.data.message;
    countdown.value = 60;
    const timer = setInterval(() => {
      countdown.value--;
      if (countdown.value <= 0) {
        clearInterval(timer);
      }
    }, 1000);
  } catch (error) {
    message.value = error.response?.data?.error || '发送失败，请重试。';
    isError.value = true;
  } finally {
    isSendingCode.value = false;
  }
};

const handleLogin = async () => {
  isLoggingIn.value = true;
  message.value = '';
  isError.value = false;
  try {
    await login(email.value, code.value);
    router.push('/');
  } catch (error) {
    message.value = error.response?.data?.error || '登录失败，请检查验证码。';
    isError.value = true;
  } finally {
    isLoggingIn.value = false;
  }
};
</script>

<style scoped>
.login-container { display: flex; justify-content: center; align-items: center; min-height: calc(100vh - 60px); background-color: #f4f4f4; }
.login-box { background: white; padding: 2rem 2.5rem; border-radius: 8px; box-shadow: 0 4px 20px rgba(0,0,0,0.1); width: 100%; max-width: 400px; text-align: center; }
h2 { text-align: center; margin-bottom: 1.5rem; }
.maintenance-header { text-align: center; padding-bottom: 1rem; margin-bottom: 1rem; }
.maintenance-header h2 { color: #d9534f; }
.maintenance-header p { color: #777; font-size: 0.9rem; }
.form-group { margin-bottom: 1.5rem; text-align: left; } /* 让 label 左对齐 */
label { display: block; margin-bottom: 0.5rem; font-weight: 500; }
input { width: 100%; padding: 0.75rem; border: 1px solid #ccc; border-radius: 4px; box-sizing: border-box; }
.code-input-group { display: flex; gap: 0.5rem; }
.code-input-group button { white-space: nowrap; padding: 0.75rem; border: 1px solid #ccc; background: #f0f0f0; cursor: pointer; }
.login-btn { width: 100%; padding: 0.8rem; background-color: #86b93a; color: white; border: none; border-radius: 4px; font-size: 1rem; cursor: pointer; }
.login-btn:disabled { background-color: #ccc; }
.message { text-align: center; margin-bottom: 1rem; font-size: 0.9rem; }
.message.error { color: red; }

/* ==================== 新增样式 ==================== */
.developer-access {
  margin-top: 1.5rem;
  padding-top: 1.5rem;
  border-top: 1px solid #f0f0f0;
}
.developer-btn {
  width: 100%;
  padding: 0.8rem;
  background-color: #5bc0de;
  color: white;
  border: none;
  border-radius: 4px;
  font-size: 1rem;
  cursor: pointer;
  transition: background-color 0.2s;
}
.developer-btn:hover {
  background-color: #31b0d5;
}
/* ====================================================== */
</style>