import { createRouter, createWebHistory } from 'vue-router';
import { getToken, getAdminInfo, setAdminInfo } from '../utils/auth';
import { getAdminInfo as fetchAdminInfo } from '../api/user';
import { ElMessage } from 'element-plus';

// 布局组件
const Layout = () => import('../views/layout/Index.vue');

// 路由配置
const routes = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('../views/login/Index.vue'),
    meta: { title: '登录', requiresAuth: false }
  },
  {
    path: '/',
    component: Layout,
    redirect: '/dashboard',
    children: [
      {
        path: 'dashboard',
        name: 'Dashboard',
        component: () => import('../views/dashboard/Index.vue'),
        meta: { title: '仪表盘', icon: 'Dashboard', requiresAuth: true }
      }
    ]
  },
  {
    path: '/user',
    component: Layout,
    meta: { title: '用户管理', icon: 'User', requiresAuth: true },
    children: [
      {
        path: 'list',
        name: 'UserList',
        component: () => import('../views/user/List.vue'),
        meta: { title: '用户列表', requiresAuth: true }
      }
    ]
  },
  {
    path: '/category',
    component: Layout,
    meta: { title: '分类管理', icon: 'Menu', requiresAuth: true },
    children: [
      {
        path: 'list',
        name: 'CategoryList',
        component: () => import('../views/category/List.vue'),
        meta: { title: '分类列表', requiresAuth: true }
      }
    ]
  },
  {
    path: '/product',
    component: Layout,
    meta: { title: '商品管理', icon: 'Goods', requiresAuth: true },
    children: [
      {
        path: 'list',
        name: 'ProductList',
        component: () => import('../views/product/List.vue'),
        meta: { title: '商品列表', requiresAuth: true }
      },
      {
        path: 'add',
        name: 'ProductAdd',
        component: () => import('../views/product/Edit.vue'),
        meta: { title: '添加商品', requiresAuth: true }
      },
      {
        path: 'edit/:id',
        name: 'ProductEdit',
        component: () => import('../views/product/Edit.vue'),
        meta: { title: '编辑商品', requiresAuth: true },
        props: true
      }
    ]
  },
  {
    path: '/order',
    component: Layout,
    meta: { title: '订单管理', icon: 'List', requiresAuth: true },
    children: [
      {
        path: 'list',
        name: 'OrderList',
        component: () => import('../views/order/List.vue'),
        meta: { title: '订单列表', requiresAuth: true }
      },
      {
        path: 'detail/:id',
        name: 'OrderDetail',
        component: () => import('../views/order/Detail.vue'),
        meta: { title: '订单详情', requiresAuth: true },
        props: true
      }
    ]
  },
  {
    path: '/package',
    component: Layout,
    meta: { title: '畅饮套餐管理', icon: 'Tickets', requiresAuth: true },
    children: [
      {
        path: 'list',
        name: 'PackageList',
        component: () => import('../views/package/List.vue'),
        meta: { title: '套餐列表', requiresAuth: true }
      },
      {
        path: 'add',
        name: 'PackageAdd',
        component: () => import('../views/package/Edit.vue'),
        meta: { title: '添加套餐', requiresAuth: true }
      },
      {
        path: 'edit/:id',
        name: 'PackageEdit',
        component: () => import('../views/package/Edit.vue'),
        meta: { title: '编辑套餐', requiresAuth: true },
        props: true
      },
      {
        path: 'user-packages',
        name: 'UserPackages',
        component: () => import('../views/package/UserPackages.vue'),
        meta: { title: '用户套餐', requiresAuth: true }
      }
    ]
  },
  {
    path: '/table',
    component: Layout,
    meta: { title: '桌号管理', icon: 'Grid', requiresAuth: true },
    children: [
      {
        path: 'list',
        name: 'TableList',
        component: () => import('../views/table/List.vue'),
        meta: { title: '桌号列表', requiresAuth: true }
      },
      {
        path: 'add',
        name: 'TableAdd',
        component: () => import('../views/table/Edit.vue'),
        meta: { title: '添加桌号', requiresAuth: true }
      }
    ]
  },
  {
    path: '/statistics',
    component: Layout,
    meta: { title: '数据统计', icon: 'DataAnalysis', requiresAuth: true },
    children: [
      {
        path: 'sales',
        name: 'SalesStats',
        component: () => import('../views/statistics/Sales.vue'),
        meta: { title: '销售统计', requiresAuth: true }
      },
      {
        path: 'packages',
        name: 'PackageStats',
        component: () => import('../views/statistics/Packages.vue'),
        meta: { title: '套餐统计', requiresAuth: true }
      }
    ]
  },
  {
    path: '/404',
    name: 'NotFound',
    component: () => import('../views/error/404.vue'),
    meta: { title: '404', requiresAuth: false }
  },
  {
    path: '/:pathMatch(.*)*',
    redirect: '/404'
  }
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