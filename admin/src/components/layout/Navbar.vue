<template>
  <div class="navbar">
    <div class="left-part">
      <el-icon class="toggle-button" @click="toggleSidebar">
        <component :is="isCollapsed ? 'Expand' : 'Fold'" />
      </el-icon>
      <breadcrumb />
    </div>
    <div class="right-part">
      <theme-switch />
      <el-dropdown trigger="click">
        <div class="avatar-container">
          <img :src="adminInfo.avatar || defaultAvatar" class="avatar-image">
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
            <el-dropdown-item divided @click="handleLogout">退出登录</el-dropdown-item>
          </el-dropdown-menu>
        </template>
      </el-dropdown>
    </div>
  </div>
</template>

<script setup>
import Breadcrumb from '@/views/layout/components/Breadcrumb.vue';
import { ElMessageBox } from 'element-plus';
import { useRouter } from 'vue-router';
import { clearAuth } from '@/utils/auth';
import { useLayout } from '@/composables/useLayout';
import ThemeSwitch from '@/components/common/ThemeSwitch.vue';

const props = defineProps({
  adminInfo: {
    type: Object,
    default: () => ({
      username: '管理员',
      nickname: '管理员',
      role: 'operator',
      avatar: ''
    })
  }
});

// 获取布局状态
const { isCollapsed, toggleSidebar } = useLayout();

const router = useRouter();

// 默认头像（内联）
const defaultAvatar = 'data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHZpZXdCb3g9IjAgMCAyMDAgMjAwIj48cmVjdCB3aWR0aD0iMjAwIiBoZWlnaHQ9IjIwMCIgZmlsbD0iI2U2ZTZlNiIvPjxjaXJjbGUgY3g9IjEwMCIgY3k9IjgwIiByPSI0MCIgZmlsbD0iI2JmYmZiZiIvPjxwYXRoIGQ9Ik0xNjAgMTgwYzAtMzMuMTM3LTI2Ljg2My02MC02MC02MHMtNjAgMjYuODYzLTYwIDYweiIgZmlsbD0iI2JmYmZiZiIvPjwvc3ZnPg==';

// 格式化角色名称
const formatRole = (role) => {
  const roleMap = {
    'super_admin': '超级管理员',
    'manager': '店长',
    'operator': '操作员'
  };
  return roleMap[role] || '操作员';
};

// 退出登录
const handleLogout = () => {
  ElMessageBox.confirm('确定要退出登录吗?', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(() => {
    clearAuth();
    router.push('/login');
  }).catch(() => {});
};
</script>

<style scoped>
.navbar {
  height: 60px;
  overflow: hidden;
  position: relative;
  background: #fff;
  box-shadow: 0 1px 4px rgba(0, 21, 41, 0.08);
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0 15px;
}

.left-part {
  display: flex;
  align-items: center;
}

.toggle-button {
  font-size: 20px;
  cursor: pointer;
  margin-right: 15px;
  color: #606266;
}

.right-part {
  display: flex;
  align-items: center;
}

.avatar-container {
  display: flex;
  align-items: center;
  cursor: pointer;
  padding: 0 8px;
}

.avatar-container:hover {
  background: rgba(0, 0, 0, 0.025);
}

.avatar-image {
  width: 30px;
  height: 30px;
  border-radius: 50%;
  margin-right: 8px;
}

.username {
  font-size: 14px;
  color: #606266;
}

.admin-info {
  padding: 5px 0;
  text-align: center;
}

.admin-role {
  font-size: 12px;
  color: #909399;
  margin-top: 5px;
}
</style> 