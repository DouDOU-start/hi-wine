import { defineStore } from 'pinia';
import { ref } from 'vue';
import { getAdminInfo as fetchAdminInfo } from '@/api/user';
import { getToken, setAdminInfo, getAdminInfo as getStoredAdminInfo, clearAuth } from '@/utils/auth';
import { ElMessage } from 'element-plus';
import router from '@/router';

export const useUserStore = defineStore('user', () => {
  // 用户信息
  const userInfo = ref({
    username: '',
    nickname: '',
    role: '',
    avatar: '',
    permissions: []
  });

  // 登录状态
  const isLoggedIn = ref(false);

  // 加载状态
  const loading = ref(false);

  // 初始化用户信息
  const initUserInfo = async () => {
    const token = getToken();
    if (!token) {
      return false;
    }

    // 尝试从本地存储获取
    const storedInfo = getStoredAdminInfo();
    if (storedInfo) {
      userInfo.value = { ...storedInfo };
      isLoggedIn.value = true;
      return true;
    }

    // 从服务器获取
    loading.value = true;
    try {
      const response = await fetchAdminInfo();
      if (response && response.code === 200 && response.data) {
        userInfo.value = { ...response.data };
        setAdminInfo(response.data);
        isLoggedIn.value = true;
        return true;
      }
      return false;
    } catch (error) {
      console.error('获取用户信息失败:', error);
      return false;
    } finally {
      loading.value = false;
    }
  };

  // 退出登录
  const logout = () => {
    clearAuth();
    userInfo.value = {
      username: '',
      nickname: '',
      role: '',
      avatar: '',
      permissions: []
    };
    isLoggedIn.value = false;
    router.push('/login');
  };

  // 检查权限
  const hasPermission = (permission) => {
    if (userInfo.value.role === 'super_admin') {
      return true;
    }
    
    if (!permission) {
      return true;
    }
    
    return userInfo.value.permissions && userInfo.value.permissions.includes(permission);
  };

  return {
    userInfo,
    isLoggedIn,
    loading,
    initUserInfo,
    logout,
    hasPermission
  };
}); 