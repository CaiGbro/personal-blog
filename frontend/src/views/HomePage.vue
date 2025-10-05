<!-- /frontend/src/views/HomePage.vue (已修复按钮文本Bug) -->
<template>
  <div class="home-container">
    <div class="home-layout-grid">
      <main class="main-content">
        <section class="intro-section">
          <h1 class="greeting"> 欢迎来访！ <span class="wave">👋</span></h1>
          <p class="bio">{{ profile.bio }}</p>
          
          <nav class="quick-links-grid">
            <router-link to="/works" class="link-card">
              <span class="icon">🧐</span>
              <span class="title">发现更多</span>
            </router-link>
            <router-link to="/writings" class="link-card">
              <span class="icon">📝</span>
              <span class="title">我的写作</span>
            </router-link>
            <router-link to="/guestbook" class="link-card">
              <span class="icon">📊</span>
              <span class="title">访客评论</span>
            </router-link>
            <a href="#" class="link-card disabled">
              <span class="icon">😎</span>
              <span class="title">关于我</span>
            </a>
            <a href="#" @click.prevent="toggleAds" class="link-card ad-trigger">
              <span class="icon">💰</span>
              <!-- ==================== 核心修复：修正这里的文本逻辑 ==================== -->
              <span class="title">{{ adsEnabled ? '关闭广告' : '广告招商' }}</span>
              <!-- ==================================================================== -->
            </a>
          </nav>

          <div class="closing-remark" :class="{ 'cool-mode': adsEnabled }">
            <span class="remark-text">{{ displayedRemark }}</span>
            <span class="typing-cursor"></span>
          </div>

        </section>
      </main>
      <aside class="sidebar">
        <ProfileCard :profile="profile" />
        <ShimmeringAd v-if="adsEnabled" />
      </aside>
    </div>

    <DisclaimerMatrix />

  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue';
import api from '../services/api';
import ProfileCard from '../components/ProfileCard.vue';
import ShimmeringAd from '../components/ShimmeringAd.vue';
import DisclaimerMatrix from '../components/DisclaimerMatrix.vue';
import { profile } from '../store/profile';
import { useAds } from '../store/ads';

const { adsEnabled, enableAds, disableAds } = useAds();
const toggleAds = () => { adsEnabled.value ? disableAds() : enableAds(); };

const fullRemark = '祝您浏览愉快 📖';
const displayedRemark = ref('');
let typingInterval = null;

const startTypingAnimation = () => {
  let currentIndex = 0;
  displayedRemark.value = '';
  if (typingInterval) clearInterval(typingInterval);
  typingInterval = setInterval(() => {
    if (currentIndex < fullRemark.length) {
      displayedRemark.value += fullRemark.charAt(currentIndex);
      currentIndex++;
    } else {
      clearInterval(typingInterval);
    }
  }, 150);
};

const prefetchFirstVideos = async () => { /* ... 此处代码不变 ... */ };

onMounted(() => {
  prefetchFirstVideos();
  startTypingAnimation();
});

onUnmounted(() => {
  if (typingInterval) {
    clearInterval(typingInterval);
  }
});
</script>

<style scoped>
/* --- 样式部分无需任何改动 --- */
.home-container { display: flex; flex-direction: column; gap: 3rem; max-width: 1024px; margin: 2rem auto; padding: 2rem; }
.home-layout-grid { display: flex; flex-wrap: wrap; gap: 2rem; }
.main-content { flex: 2; min-width: 300px; }
.sidebar { flex: 1; min-width: 280px; }
.greeting { font-size: 3.5rem; font-weight: 700; color: #86b93a; margin-bottom: 1.5rem; }
.wave { display: inline-block; animation: wave-animation 2s infinite; }
@keyframes wave-animation { 0% { transform: rotate( 0.0deg) } 10% { transform: rotate(14.0deg) } 20% { transform: rotate(-8.0deg) } 30% { transform: rotate(14.0deg) } 40% { transform: rotate(-4.0deg) } 50% { transform: rotate(10.0deg) } 60% { transform: rotate( 0.0deg) } 100% { transform: rotate( 0.0deg) } }
.bio { line-height: 1.8; font-size: 1.1rem; color: var(--text-secondary); }
.quick-links-grid { display: grid; grid-template-columns: repeat(2, 1fr); gap: 1rem; margin-top: 2.5rem; }
.link-card { display: flex; flex-direction: column; align-items: center; justify-content: center; padding: 1.5rem; background-color: var(--bg-secondary); border: 1px solid var(--border-color); border-radius: 12px; text-decoration: none; color: var(--text-primary); transition: transform 0.2s ease, box-shadow 0.2s ease; }
.link-card:hover { transform: translateY(-5px); box-shadow: 0 8px 24px rgba(0, 0, 0, 0.08); }
.link-card .icon { font-size: 2rem; margin-bottom: 0.75rem; }
.link-card .title { font-weight: 500; font-size: 1rem; }
.link-card.disabled { opacity: 0.6; cursor: not-allowed; }
.closing-remark { margin-top: 3rem; text-align: center; font-size: 1.5rem; font-weight: 500; position: relative; overflow: hidden; padding: 0.5rem 0; color: #333; font-family: 'KaiTi', 'STKaiti', serif; transition: all 0.5s ease; }
.typing-cursor { display: inline-block; background-color: #333; margin-left: 8px; animation: blink 0.7s infinite; transition: all 0.5s ease; width: 3px; height: 1.5rem; }
.closing-remark.cool-mode { color: #00ffff; font-family: 'Courier New', Courier, monospace; text-shadow: 0 0 5px #00ffff, 0 0 10px #00ffff; }
.closing-remark.cool-mode .remark-text::after { content: ''; position: absolute; top: -50%; left: -10%; width: 120%; height: 4px; background: cyan; box-shadow: 0 0 15px cyan; opacity: 0.6; transform: rotate(-10deg); animation: scan-line 4s linear infinite; }
.closing-remark.cool-mode .typing-cursor { background-color: #00ffff; box-shadow: 0 0 5px #00ffff; width: 1rem; height: 1.6rem; }
@keyframes blink { 0%, 100% { opacity: 1; } 50% { opacity: 0; } }
@keyframes scan-line { 0% { transform: translate(-10%, -100%) rotate(-10deg); } 100% { transform: translate(10%, 500%) rotate(-10deg); } }
@media (max-width: 768px) { .home-layout-grid { flex-direction: column; } .greeting { font-size: 2.5rem; } .quick-links-grid { gap: 0.75rem; } .closing-remark { font-size: 1.2rem; } .typing-cursor { height: 1.2rem; } .closing-remark.cool-mode .typing-cursor { height: 1.3rem; width: 0.8rem; } }
</style>