// 布局组件
const Layout = () => import('@/views/layout/Index.vue');

export default {
  path: '/package',
  component: Layout,
  meta: { title: '畅饮套餐管理', icon: 'Tickets', requiresAuth: true },
  children: [
    {
      path: 'list',
      name: 'PackageList',
      component: () => import('@/views/package/List.vue'),
      meta: { title: '套餐列表', requiresAuth: true }
    },
    {
      path: 'add',
      name: 'PackageAdd',
      component: () => import('@/views/package/Edit.vue'),
      meta: { title: '添加套餐', requiresAuth: true }
    },
    {
      path: 'edit/:id',
      name: 'PackageEdit',
      component: () => import('@/views/package/Edit.vue'),
      meta: { title: '编辑套餐', requiresAuth: true },
      props: true
    },
    {
      path: 'user-packages',
      name: 'UserPackages',
      component: () => import('@/views/package/UserPackages.vue'),
      meta: { title: '用户套餐', requiresAuth: true }
    }
  ]
}; 