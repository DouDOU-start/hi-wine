// 布局组件
const Layout = () => import('@/views/layout/Index.vue');

export default {
  path: '/category',
  component: Layout,
  meta: { title: '分类管理', icon: 'Menu', requiresAuth: true },
  children: [
    {
      path: 'list',
      name: 'CategoryList',
      component: () => import('@/views/category/List.vue'),
      meta: { title: '分类列表', requiresAuth: true }
    }
  ]
}; 