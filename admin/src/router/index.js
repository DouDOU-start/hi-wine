import { createRouter, createWebHistory } from 'vue-router';
import { getToken, getAdminInfo, setAdminInfo } from '../utils/auth';
import { getAdminInfo as fetchAdminInfo } from '../api/user';
import { ElMessage } from 'element-plus';

// 导入路由模块
import dashboardRoutes from './modules/dashboard';
import userRoutes from './modules/user';
import categoryRoutes from './modules/category';
import productRoutes from './modules/product';
import orderRoutes from './modules/order';
import packageRoutes from './modules/package';
import tableRoutes from './modules/table';
import authRoutes from './modules/auth';
import errorRoutes from './modules/error';

// 路由配置
const routes = [
  ...authRoutes,
  dashboardRoutes,
  userRoutes,
  categoryRoutes,
  productRoutes,
  orderRoutes,
  packageRoutes,
  tableRoutes,
  ...errorRoutes
];

// 创建路由实例
const router = createRouter({
  history: createWebHistory(),
  routes
});

// 路由守卫
router.beforeEach(async (to, from, next) => {
  // 设置页面标题
  document.title = to.meta.title ? `${to.meta.title} - Hi-Wine酒水销售系统` : 'Hi-Wine酒水销售系统';
  
  // 检查是否需要登录认证
  if (to.meta.requiresAuth !== false) {
    const token = getToken();
    if (!token) {
      // 没有token，跳转到登录页
      next({ path: '/login', query: { redirect: to.fullPath } });
      return;
    }
    
    // 检查是否有管理员信息
    const adminInfo = getAdminInfo();
    if (!adminInfo) {
      try {
        // 尝试获取管理员信息
        const response = await fetchAdminInfo();
        if (response.data) {
          // 存储管理员信息
          setAdminInfo(response.data);
          next();
        } else {
          // 没有管理员信息，可能token已失效
          ElMessage.error('登录已过期，请重新登录');
          next({ path: '/login', query: { redirect: to.fullPath } });
        }
      } catch (error) {
        console.error('获取管理员信息失败:', error);
        // 获取管理员信息失败，可能是token无效
        ElMessage.error('登录已过期，请重新登录');
        next({ path: '/login', query: { redirect: to.fullPath } });
      }
    } else {
      // 有管理员信息，直接通过
      next();
    }
  } else {
    // 不需要认证的路由，直接通过
    next();
  }
});

export default router; 