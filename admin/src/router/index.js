import { createRouter, createWebHistory } from 'vue-router';

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
router.beforeEach((to, from, next) => {
  // 设置页面标题
  document.title = to.meta.title ? `${to.meta.title} - 酒馆后台管理系统` : '酒馆后台管理系统';
  
  // 检查是否需要登录认证
  if (to.meta.requiresAuth !== false) {
    const token = localStorage.getItem('token');
    if (!token) {
      next({ path: '/login', query: { redirect: to.fullPath } });
      return;
    }
  }
  
  next();
});

export default router; 