// /frontend/src/router/index.js
import { createRouter, createWebHistory } from 'vue-router';
import { useAuth } from '../services/auth'; // <-- 引入 auth
import HomePage from '../views/HomePage.vue';
import GuestbookPage from '../views/GuestbookPage.vue';
import WorksPage from '../views/WorksPage.vue';
import MediaViewerPage from '../views/MediaViewerPage.vue';
import WritingsPage from '../views/WritingsPage.vue';
import WritingViewerPage from '../views/WritingViewerPage.vue';
// --- 新增引入 ---
import VideoGalleryPage from '../views/VideoGalleryPage.vue';
import LoginPage from '../views/LoginPage.vue';
import AdminPage from '../views/AdminPage.vue';


const routes = [
  { path: '/', name: 'Home', component: HomePage },
  { path: '/guestbook', name: 'Guestbook', component: GuestbookPage },
  { path: '/works', name: 'Works', component: WorksPage },
  // --- 新增路由 ---
  { path: '/works/video-gallery', name: 'VideoGallery', component: VideoGalleryPage },
  // MediaViewer 路由保持不变，它现在是播放器页面
  { path: '/works/media-viewer', name: 'MediaViewer', component: MediaViewerPage },
  { path: '/writings', name: 'Writings', component: WritingsPage },
  { path: '/writings/:filename', name: 'WritingViewer', component: WritingViewerPage, props: true },

  { path: '/login', name: 'Login', component: LoginPage },
  { path: '/admin', name: 'Admin', component: AdminPage, meta: { requiresAuth: true } }, // 添加 meta 标记
];

const router = createRouter({
  history: createWebHistory(),
  routes,
  
  // +++ 新增下面的 scrollBehavior 函数 +++
  scrollBehavior(to, from, savedPosition) {
    // 如果 savedPosition 存在，说明是用户通过浏览器的“后退/前进”按钮导航
    if (savedPosition) {
      // 直接返回保存的位置
      return savedPosition;
    } else {
      // 否则，是新的导航，滚动到页面顶部
      return { top: 0 };
    }
  },
  // +++ 新增结束 +++
});

// --- 新增：全局路由守卫 ---
router.beforeEach(async (to, from, next) => {
  const { 
    isLoggedIn, 
    isVisitorModeEnabled, 
    fetchSystemStatus, 
    hasCheckedStatus 
  } = useAuth();

  // 如果 App 还未从后端获取过系统状态，则先获取
  if (!hasCheckedStatus.value) {
    await fetchSystemStatus();
  }

  // 规则 1: 如果游客模式开启，则放行所有
  if (isVisitorModeEnabled.value) {
    return next();
  }

  // 规则 2: 如果游客模式关闭
  const requiresAuth = to.matched.some(record => record.meta.requiresAuth);
  
  if (to.name !== 'Login' && !isLoggedIn.value) {
    // 如果用户未登录且想访问的不是登录页，则重定向到登录页，并带上“维护模式”的提示
    return next({ name: 'Login', query: { mode: 'maintenance' } });
  }
  
  if (requiresAuth && !isLoggedIn.value) {
     // 如果页面需要认证但用户未登录 (这部分逻辑与上面重叠，但作为双重保险)
     return next({ name: 'Login', query: { mode: 'maintenance' } });
  }

  // 其他情况（已登录或访问登录页）
  next();
});

export default router;