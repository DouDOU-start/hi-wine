// 布局组件
const Layout = () => import('@/views/layout/Index.vue');

export default {
  path: '/table',
  component: Layout,
  meta: { title: '桌号管理', icon: 'Grid', requiresAuth: true },
  children: [
    {
      path: 'list',
      name: 'TableList',
      component: () => import('@/views/table/List.vue'),
      meta: { title: '桌号列表', requiresAuth: true }
    },
    {
      path: 'add',
      name: 'TableAdd',
      component: () => import('@/views/table/Edit.vue'),
      meta: { title: '添加桌号', requiresAuth: true }
    },
    {
      path: 'edit/:id',
      name: 'TableEdit',
      component: () => import('@/views/table/Edit.vue'),
      meta: { title: '编辑桌号', requiresAuth: true },
      props: true
    }
  ]
}; 