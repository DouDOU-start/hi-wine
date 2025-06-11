// 布局组件
const Layout = () => import('@/views/layout/Index.vue');

export default {
  path: '/',
  component: Layout,
  redirect: '/dashboard',
  children: [
    {
      path: 'dashboard',
      name: 'Dashboard',
      component: () => import('@/views/dashboard/Index.vue'),
      meta: { title: '仪表盘', icon: 'Dashboard', requiresAuth: true }
    }
  ]
}; 