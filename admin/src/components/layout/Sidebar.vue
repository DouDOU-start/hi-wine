<template>
  <div class="sidebar-container" :class="{ 'is-collapsed': isCollapsed }">
    <div class="logo-container">
      <img :src="logoImage" alt="Logo" class="logo-image" v-if="!isCollapsed">
      <img :src="logoSmall" alt="Logo" class="logo-small" v-else>
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
      </el-sub-menu>
      
      <el-menu-item index="/user/list">
        <el-icon><component :is="'User'" /></el-icon>
        <template #title>用户管理</template>
      </el-menu-item>
    </el-menu>
  </div>
</template>

<script setup>
import { computed } from 'vue';
import { useRoute } from 'vue-router';
import { useLayout } from '@/composables/useLayout';

// 获取布局状态
const { isCollapsed } = useLayout();

const route = useRoute();

// 内联logo图片
const logoImage = 'data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHZpZXdCb3g9IjAgMCAyMDAgNTAiPjxyZWN0IHdpZHRoPSIyMDAiIGhlaWdodD0iNTAiIGZpbGw9Im5vbmUiLz48dGV4dCB4PSIxMCIgeT0iMzUiIGZvbnQtZmFtaWx5PSJBcmlhbCIgZm9udC1zaXplPSIyNCIgZmlsbD0iI2JmY2JkOSI+SGktV2luZTwvdGV4dD48L3N2Zz4=';
const logoSmall = 'data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHZpZXdCb3g9IjAgMCA1MCA1MCI+PGNpcmNsZSBjeD0iMjUiIGN5PSIyNSIgcj0iMjAiIGZpbGw9Im5vbmUiIHN0cm9rZT0iI2JmY2JkOSIgc3Ryb2tlLXdpZHRoPSIyIi8+PHRleHQgeD0iMTIiIHk9IjM1IiBmb250LWZhbWlseT0iQXJpYWwiIGZvbnQtc2l6ZT0iMjQiIGZpbGw9IiNiZmNiZDkiPkg8L3RleHQ+PC9zdmc+';

// 当前激活的菜单
const activeMenu = computed(() => {
  const { path } = route;
  return path;
});
</script>

<style scoped>
.sidebar-container {
  height: 100%;
  background-color: #304156;
  transition: width 0.3s;
  width: 210px;
  position: fixed;
  left: 0;
  top: 0;
  bottom: 0;
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
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: #263445;
}

.logo-image {
  height: 32px;
  max-width: 100%;
}

.logo-small {
  height: 32px;
  width: 32px;
}

.sidebar-menu {
  border-right: none;
}

.sidebar-menu:not(.el-menu--collapse) {
  width: 210px;
}
</style> 