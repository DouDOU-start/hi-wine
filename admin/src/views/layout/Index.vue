<template>
  <div class="app-wrapper">
    <!-- 侧边栏 -->
    <sidebar />
    
    <!-- 主内容区 -->
    <div class="main-container" :class="{ 'is-collapsed': isCollapsed }">
      <!-- 顶部导航栏 -->
      <navbar :admin-info="adminInfo" />
      
      <!-- 内容区 -->
      <div class="app-main">
        <router-view v-slot="{ Component }">
          <transition name="fade-transform" mode="out-in">
            <component :is="Component" />
          </transition>
        </router-view>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { getAdminInfo } from '@/utils/auth';
import Sidebar from '@/components/layout/Sidebar.vue';
import Navbar from '@/components/layout/Navbar.vue';
import { provideLayout } from '@/composables/useLayout';

// 提供布局状态
const { isCollapsed, initSidebarStatus } = provideLayout();

// 管理员信息
const adminInfo = ref({
  username: '管理员',
  nickname: '管理员',
  role: 'operator',
  avatar: ''
});

// 获取管理员信息
const loadAdminInfo = () => {
  const info = getAdminInfo();
  if (info) {
    adminInfo.value = {
      username: info.username || '管理员',
      nickname: info.nickname || info.username || '管理员',
      role: info.role || 'operator',
      avatar: info.avatarUrl || ''
    };
  }
};

// 初始化
onMounted(() => {
  // 初始化侧边栏状态
  initSidebarStatus();
  
  // 加载管理员信息
  loadAdminInfo();
});
</script>

<style scoped>
.app-wrapper {
  position: relative;
  height: 100vh;
  width: 100%;
}

.main-container {
  min-height: 100%;
  transition: margin-left 0.3s;
  margin-left: 210px;
  position: relative;
}

.main-container.is-collapsed {
  margin-left: 64px;
}

.app-main {
  padding: 20px;
  min-height: calc(100vh - 60px);
  position: relative;
  overflow: auto;
  box-sizing: border-box;
}

/* 过渡动画 */
.fade-transform-enter-active,
.fade-transform-leave-active {
  transition: all 0.3s;
}

.fade-transform-enter-from {
  opacity: 0;
  transform: translateX(30px);
}

.fade-transform-leave-to {
  opacity: 0;
  transform: translateX(-30px);
}
</style> 