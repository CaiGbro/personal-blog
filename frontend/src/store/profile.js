// /frontend/src/store/profile.js
import { ref } from 'vue';
// 确保你的头像图片放在 assets 文件夹下
import profileAvatar from '../assets/profile-photo.jpg';

// 将 profile 数据移到这里，以便多个组件可以共享
export const profile = ref({
  name: "CaiGbro",
  title: "AI应用与全栈工程师",
  location: "上海",
  email: "CaiGbro@163.com",
  github: "CaiGbro",
  website: "CaiGbro.top",
  bio: "CaiGbro， 互联网牛马。  本科机械、研究生计算机，从放牛到看店、从罐头厂到电子厂、从服务员到家教，从 C 、 C++ 到 JAVA 实习、全栈实习、自媒体，从 Agent 到大模型 、自动化...， 从体育到历史、哲学到科学...，哥们东西都学杂了。",
  avatarUrl: profileAvatar,
});