// 布局组件
const Layout = () => import('@/views/layout/Index.vue');

export default {
  path: '/product',
  component: Layout,
  meta: { title: '商品管理', icon: 'Goods', requiresAuth: true },
  children: [
    {
      path: 'list',
      name: 'ProductList',
      component: () => import('@/views/product/List.vue'),
      meta: { title: '商品列表', requiresAuth: true }
    },
    {
      path: 'add',
      name: 'ProductAdd',
      component: () => import('@/views/product/Edit.vue'),
      meta: { title: '添加商品', requiresAuth: true }
    },
    {
      path: 'edit/:id',
      name: 'ProductEdit',
      component: () => import('@/views/product/Edit.vue'),
      meta: { title: '编辑商品', requiresAuth: true },
      props: true
    }
  ]
}; 