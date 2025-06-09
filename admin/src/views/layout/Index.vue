<template>
  <div class="app-wrapper">
    <!-- 侧边栏 -->
    <div class="sidebar-container" :class="{ 'is-collapsed': isCollapsed }">
      <div class="logo-container">
        <img src="../../assets/logo.png" alt="Logo" class="logo-image" v-if="!isCollapsed">
        <img src="../../assets/logo-small.png" alt="Logo" class="logo-small" v-else>
      </div>
      
      <el-menu
        :default-active="activeMenu"
        class="sidebar-menu"
        :collapse="isCollapsed"
        background-color="#304156"
        text-color="#bfcbd9"
        active-text-color="#409EFF"
        :collapse-transition="false"
        router
      >
        <el-menu-item index="/dashboard">
          <el-icon><component :is="'Odometer'" /></el-icon>
          <template #title>仪表盘</template>
        </el-menu-item>
        
        <el-sub-menu index="/product">
          <template #title>
            <el-icon><component :is="'Goods'" /></el-icon>
            <span>商品管理</span>
          </template>
          <el-menu-item index="/product/list">商品列表</el-menu-item>
          <el-menu-item index="/product/add">添加商品</el-menu-item>
        </el-sub-menu>
        
        <el-menu-item index="/category/list">
          <el-icon><component :is="'Menu'" /></el-icon>
          <template #title>分类管理</template>
        </el-menu-item>
        
        <el-menu-item index="/order/list">
          <el-icon><component :is="'List'" /></el-icon>
          <template #title>订单管理</template>
        </el-menu-item>
        
        <el-sub-menu index="/package">
          <template #title>
            <el-icon><component :is="'Tickets'" /></el-icon>
            <span>畅饮套餐管理</span>
          </template>
          <el-menu-item index="/package/list">套餐列表</el-menu-item>
          <el-menu-item index="/package/add">添加套餐</el-menu-item>
          <el-menu-item index="/package/user-packages">用户套餐</el-menu-item>
        </el-sub-menu>
        
        <el-sub-menu index="/table">
          <template #title>
            <el-icon><component :is="'Grid'" /></el-icon>
            <span>桌号管理</span>
          </template>
          <el-menu-item index="/table/list">桌号列表</el-menu-item>
          <el-menu-item index="/table/add">添加桌号</el-menu-item>
          <el-menu-item index="/table/edit">编辑桌号</el-menu-item>
        </el-sub-menu>
        
        <el-menu-item index="/user/list">
          <el-icon><component :is="'User'" /></el-icon>
          <template #title>用户管理</template>
        </el-menu-item>
      </el-menu>
    </div>
    
    <!-- 主内容区 -->
    <div class="main-container">
      <!-- 顶部导航栏 -->
      <div class="navbar">
        <div class="left-part">
          <el-icon class="toggle-button" @click="toggleSidebar">
            <component :is="isCollapsed ? 'Expand' : 'Fold'" />
          </el-icon>
          <breadcrumb />
        </div>
        <div class="right-part">
          <el-dropdown trigger="click">
            <div class="avatar-container">
              <img src="../../assets/avatar.png" class="avatar-image">
              <span class="username">{{ adminInfo.nickname }}</span>
            </div>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item disabled>
                  <div class="admin-info">
                    <div>{{ adminInfo.nickname }}</div>
                    <div class="admin-role">{{ formatRole(adminInfo.role) }}</div>
                  </div>
                </el-dropdown-item>
                <el-dropdown-item divided @click="logout">退出登录</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </div>
      
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
import { ref, computed, onMounted } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import { ElMessageBox } from 'element-plus';
import Breadcrumb from './components/Breadcrumb.vue';
import { getAdminInfo, clearAuth } from '../../utils/auth';

const router = useRouter();
const route = useRoute();

// 侧边栏折叠状态
const isCollapsed = ref(false);

// 管理员信息
const adminInfo = ref({
  username: '管理员',
  nickname: '管理员',
  role: 'operator',
  avatar: '../../assets/avatar.png'
});

// 格式化角色名称
const formatRole = (role) => {
  const roleMap = {
    'super_admin': '超级管理员',
    'manager': '店长',
    'operator': '操作员'
  };
  return roleMap[role] || '操作员';
};

// 切换侧边栏折叠状态
const toggleSidebar = () => {
  isCollapsed.value = !isCollapsed.value;
  localStorage.setItem('sidebarStatus', isCollapsed.value ? '1' : '0');
};

// 当前激活的菜单
const activeMenu = computed(() => {
  const { path } = route;
  return path;
});

// 退出登录
const logout = () => {
  ElMessageBox.confirm('确定要退出登录吗?', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(() => {
    clearAuth();
    router.push('/login');
  }).catch(() => {});
};

// 获取管理员信息
const loadAdminInfo = () => {
  const info = getAdminInfo();
  if (info) {
    adminInfo.value = {
      username: info.username || '管理员',
      nickname: info.nickname || info.username || '管理员',
      role: info.role || 'operator',
      avatar: info.avatarUrl || '../../assets/avatar.png'
    };
  }
};

// 初始化
onMounted(() => {
  // 从本地存储读取侧边栏状态
  const sidebarStatus = localStorage.getItem('sidebarStatus');
  isCollapsed.value = sidebarStatus === '1';
  
  // 加载管理员信息
  loadAdminInfo();
});
</script>

<style scoped>
.app-wrapper {
  position: relative;
  height: 100vh;
  width: 100%;
  display: flex;
}

.sidebar-container {
  width: 210px;
  height: 100%;
  position: fixed;
  left: 0;
  top: 0;
  bottom: 0;
  background-color: #304156;
  transition: width 0.28s;
  z-index: 1001;
  overflow-y: auto;
  overflow-x: hidden;
}

.sidebar-container.is-collapsed {
  width: 64px;
}

.logo-container {
  height: 60px;
  padding: 10px;
  text-align: center;
  overflow: hidden;
  display: flex;
  align-items: center;
  justify-content: center;
}

.logo-image {
  height: 40px;
  max-width: 100%;
}

.logo-small {
  height: 30px;
  max-width: 100%;
}

.sidebar-menu {
  border-right: none;
}

.main-container {
  flex: 1;
  margin-left: 210px;
  position: relative;
  display: flex;
  flex-direction: column;
  transition: margin-left 0.28s;
  min-height: 100%;
  background-color: #f0f2f5;
}

.is-collapsed + .main-container {
  margin-left: 64px;
}

.navbar {
  height: 60px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 15px;
  background-color: #fff;
  box-shadow: 0 1px 4px rgba(0, 21, 41, 0.08);
  position: relative;
}

.left-part, .right-part {
  display: flex;
  align-items: center;
}

.toggle-button {
  font-size: 20px;
  cursor: pointer;
  margin-right: 15px;
  color: #606266;
}

.avatar-container {
  display: flex;
  align-items: center;
  cursor: pointer;
}

.avatar-image {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  margin-right: 8px;
}

.username {
  font-size: 14px;
  color: #606266;
}

.admin-info {
  padding: 6px 0;
}

.admin-role {
  font-size: 12px;
  color: #909399;
  margin-top: 4px;
}

.app-main {
  padding: 15px;
  flex: 1;
  overflow-y: auto;
}

/* 过渡动画 */
.fade-transform-enter-active,
.fade-transform-leave-active {
  transition: all 0.3s;
}

.fade-transform-enter-from {
  opacity: 0;
  transform: translateX(20px);
}

.fade-transform-leave-to {
  opacity: 0;
  transform: translateX(-20px);
}
</style> 