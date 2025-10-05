<!-- /frontend/src/components/ProfileCard.vue (已加入炫酷特效) -->
<template>
  <div class="profile-card">
    <img v-if="profile.avatarUrl" :src="profile.avatarUrl" alt="Profile photo" class="avatar" />
    <div class="spotify-status">
      <span>🎵</span> Not Playing — Spotify
    </div>
    <div class="card-content">
      <h2>{{ profile.name }}</h2>
      <p class="subtitle">学者 | 实干</p>
      <ul class="info-list">
        <li><span class="icon">💼</span> {{ profile.title }}</li>
        <li><span class="icon">📍</span> {{ profile.location }}</li>
        <li><span class="icon">📧</span> {{ profile.email }}</li>
        <li>
          <span class="icon">
            <svg class="github-svg" viewBox="0 0 16 16" version="1.1" width="16" height="16" aria-hidden="true"><path fill-rule="evenodd" d="M8 0C3.58 0 0 3.58 0 8c0 3.54 2.29 6.53 5.47 7.59.4.07.55-.17.55-.38 0-.19-.01-.82-.01-1.49-2.01.37-2.53-.49-2.69-.94-.09-.23-.48-.94-.82-1.13-.28-.15-.68-.52-.01-.53.63-.01 1.08.58 1.23.82.72 1.21 1.87.87 2.33.66.07-.52.28-.87.51-1.07-1.78-.2-3.64-.89-3.64-3.95 0-.87.31-1.59.82-2.15-.08-.2-.36-1.02.08-2.12 0 0 .67-.21 2.2.82.64-.18 1.32-.27 2-.27.68 0 1.36.09 2 .27 1.53-1.04 2.2-.82 2.2-.82.44 1.1.16 1.92.08 2.12.51.56.82 1.27.82 2.15 0 3.07-1.87 3.75-3.65 3.95.29.25.54.73.54 1.48 0 1.07-.01 1.93-.01 2.2 0 .21.15.46.55.38A8.013 8.013 0 0016 8c0-4.42-3.58-8-8-8z"></path></svg>
          </span>
          <a :href="'https://github.com/' + profile.github" target="_blank" rel="noopener noreferrer">
            /{{ profile.github }}
          </a>
          <span class="separator">|</span>
          <span class="icon">🌐</span>
          <a :href="'https://' + profile.website" target="_blank" rel="noopener noreferrer">
            /{{ profile.website }}
          </a>
        </li>
      </ul>
    </div>
  </div>
</template>

<script setup>
defineProps({
  profile: Object
});
</script>

<style scoped>
/* --- 默认模式下的样式 (保持不变) --- */
.profile-card {
  background-color: white;
  border-radius: 16px;
  overflow: hidden;
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.05);
  border: 1px solid #f0f0f0;
  transition: all 0.4s ease; /* 延长过渡时间，让切换更平滑 */
}
.avatar {
  width: 100%;
  height: auto;
  display: block;
  transition: all 0.4s ease;
}
.spotify-status {
  padding: 8px 16px;
  font-size: 0.8rem;
  background-color: #f5f5f5;
  color: #555;
  border-bottom: 1px solid #f0f0f0;
  transition: all 0.4s ease;
}
.card-content { padding: 1.5rem; }
h2 { margin: 0 0 0.25rem; }
.subtitle { font-size: 0.9rem; color: #777; margin: 0 0 1.5rem; }
.info-list { list-style: none; padding: 0; margin: 0; font-size: 0.9rem; }
.info-list li { margin-bottom: 0.75rem; display: flex; align-items: center; }
.icon { margin-right: 0.75rem; display: inline-flex; align-items: center; }
.info-list a { text-decoration: none; color: inherit; transition: color 0.2s; }
.info-list a:hover { color: #86b93a; }
.separator { margin: 0 0.5rem; }
.github-svg { width: 18px; height: 18px; fill: currentColor; }

/* ================================================================== */
/* --- 新增：广告激活时的“全息显示屏”样式 --- */
/* ================================================================== */

/* 1. 卡片主体：变为半透明的“毛玻璃”材质，并添加辉光边框 */
:is(.ads-active) .profile-card {
  background: rgba(15, 23, 42, 0.6); /* 深蓝半透明背景 */
  backdrop-filter: blur(12px); /* 毛玻璃效果 */
  border: 1px solid rgba(0, 255, 255, 0.3); /* 青色边框 */
  box-shadow: 0 0 25px rgba(0, 255, 255, 0.3); /* 青色辉光 */
}

/* 2. 头像：添加赛博朋克风格的扫描线和脉冲动画 */
:is(.ads-active) .avatar {
  opacity: 0.8;
  filter: contrast(1.1) saturate(1.2); /* 增强对比度和饱和度 */
  animation: avatar-pulse 6s infinite ease-in-out;
  position: relative;
}
:is(.ads-active) .avatar::after {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: repeating-linear-gradient(
    0deg,
    rgba(0, 255, 255, 0.15),
    rgba(0, 255, 255, 0.15) 1px,
    transparent 1px,
    transparent 3px
  );
}

/* 3. Spotify 栏：与卡片主体风格统一 */
:is(.ads-active) .spotify-status {
  background-color: rgba(0, 0, 0, 0.4);
  border-bottom: 1px solid rgba(0, 255, 255, 0.2);
  color: #c7d2fe; /* 淡紫色文字 */
  text-shadow: 0 0 5px rgba(0, 255, 255, 0.5);
}

/* 4. 内容区文字和图标：添加发光效果 */
:is(.ads-active) .card-content h2,
:is(.ads-active) .card-content .subtitle {
  color: #fff;
  text-shadow: 0 0 8px rgba(52, 152, 219, 0.8);
}
:is(.ads-active) .info-list {
  color: #c7d2fe;
  text-shadow: 0 0 5px rgba(52, 152, 219, 0.5);
}
:is(.ads-active) .icon {
  color: #00ffff; /* 图标变为亮青色 */
  filter: drop-shadow(0 0 3px #00ffff); /* 给图标添加辉光 */
}
:is(.ads-active) .github-svg {
  fill: #00ffff;
}

/* 5. 链接：悬停时变为更亮的颜色 */
:is(.ads-active) .info-list a:hover {
  color: #f1c40f; /* 悬停时变为金色 */
}

/* --- 动画定义 --- */
@keyframes avatar-pulse {
  0%, 100% { transform: scale(1); }
  50% { transform: scale(1.02); }
}
</style>