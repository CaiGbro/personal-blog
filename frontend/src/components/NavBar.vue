<!-- /frontend/src/components/NavBar.vue (最终修复版) -->
<template>
  <!-- 1. 确保 :class 中绑定了 'menu-open': isMenuOpen -->
  <nav class="navbar" :class="{ 'scrolled': isScrolled, 'menu-open': isMenuOpen }" v-if="isNavBarVisible">
    <div class="nav-container">
      <router-link to="/" class="nav-brand">
        <img :src="profile.avatarUrl" alt="Avatar" class="nav-avatar" />
      </router-link>

      <!-- 2. 确保汉堡按钮的 @click 事件正确设置 -->
      <button class="hamburger" @click="isMenuOpen = !isMenuOpen">
        <span></span>
        <span></span>
        <span></span>
      </button>

      <!-- 导航菜单 (这部分不需要修改) -->
      <ul class="nav-menu">
        <li><router-link to="/">Blog</router-link></li>
        <li><router-link to="/writings">Snippets</router-link></li>
        <li><router-link to="/works">Projects</router-link></li>
        <li><router-link to="/">About</router-link></li>
        <template v-if="isLoggedIn">
          <li><router-link to="/admin">后台管理</router-link></li>
          <li><a href="#" @click.prevent="logout">退出登录</a></li>
        </template>
        <template v-else>
          <li><router-link to="/login">管理员登录</router-link></li>
        </template>
      </ul>
    </div>
  </nav>
</template>

<script setup>
// 确保从 'vue' 导入了 ref
import { ref, onMounted, onUnmounted } from 'vue';
import { profile } from '../store/profile';
import { useAuth } from '../services/auth';
import { useLayout } from '../store/layout';

const { isLoggedIn, logout } = useAuth();
const { isNavBarVisible } = useLayout();

// --- 关键：确保这一行存在 ---
const isMenuOpen = ref(false); // 新增状态，控制菜单是否展开

const isScrolled = ref(false);

const handleScroll = () => {
  isScrolled.value = window.scrollY > 10;
};

onMounted(() => {
  window.addEventListener('scroll', handleScroll);
});

onUnmounted(() => {
  window.removeEventListener('scroll', handleScroll);
});
</script>

<style scoped>
.navbar {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  /* 7. 初始背景色可以更透明一些，因为下方就是白色 */
  background-color: rgba(255, 255, 255, 0.7);
  backdrop-filter: blur(10px);
  z-index: 1000;
  height: 60px;
  display: flex;
  align-items: center;
  
  /* 8. 默认状态下，完全没有阴影 */
  box-shadow: none;

  /* 9. 为阴影和背景色的变化添加平滑的过渡效果 */
  transition: box-shadow 0.3s ease, background-color 0.3s ease;
}

/* 10. 当 .navbar 元素拥有 .scrolled 类时，应用这些样式 */
.navbar.scrolled {
  /* 恢复背景色的不透明度，使其在复杂内容上更清晰 */
  background-color: rgba(255, 255, 255, 0.8);
  /* 添加柔和的阴影 */
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
}

.nav-container {
  max-width: 1024px;
  width: 100%;
  margin: 0 auto;
  padding: 0 2rem;
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.nav-brand {
  display: flex;
  align-items: center;
}
.nav-avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  object-fit: cover;
}
.nav-menu {
  list-style: none;
  display: flex;
  gap: 1.5rem;
  margin: 0;
  padding: 0;
}
.nav-menu a {
  color: #333;
  font-weight: 500;
  transition: color 0.2s;
}
.nav-menu a:hover {
  color: #86b93a;
}

/* --- 移动端导航栏样式 --- */
.hamburger {
  display: none; /* 默认隐藏 */
  flex-direction: column;
  justify-content: space-around;
  width: 2rem;
  height: 2rem;
  background: transparent;
  border: none;
  cursor: pointer;
  padding: 0;
  z-index: 10;
}
.hamburger span {
  width: 2rem;
  height: 0.25rem;
  background: #333;
  border-radius: 10px;
  transition: all 0.3s linear;
  position: relative;
  transform-origin: 1px;
}

@media (max-width: 768px) {
  .nav-menu {
    position: fixed;
    top: 60px; /* 从导航栏下方开始 */
    left: 0;
    width: 100%;
    height: calc(100vh - 60px);
    background-color: rgba(255, 255, 255, 0.95);
    backdrop-filter: blur(10px);
    flex-direction: column;
    align-items: center;
    padding-top: 2rem;
    
    /* 关键：根据 isMenuOpen 状态来显示或隐藏 */
    transform: translateX(100%);
    transition: transform 0.3s ease-in-out;
  }

  .navbar.menu-open .nav-menu {
    transform: translateX(0);
  }

  .hamburger {
    display: flex; /* 在小屏幕上显示汉堡按钮 */
  }
}

</style>