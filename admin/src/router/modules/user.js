// 布局组件
const Layout = () => import('@/views/layout/Index.vue');

export default {
  path: '/user',
  component: Layout,
  meta: { title: '用户管理', icon: 'User', requiresAuth: true },
  children: [
    {
      path: 'list',
      name: 'UserList',
      component: () => import('@/views/user/List.vue'),
      meta: { title: '用户列表', requiresAuth: true }
    }
  ]
}; 