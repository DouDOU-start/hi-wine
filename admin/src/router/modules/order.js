// 布局组件
const Layout = () => import('@/views/layout/Index.vue');

export default {
  path: '/order',
  component: Layout,
  meta: { title: '订单管理', icon: 'List', requiresAuth: true },
  children: [
    {
      path: 'list',
      name: 'OrderList',
      component: () => import('@/views/order/List.vue'),
      meta: { title: '订单列表', requiresAuth: true }
    },
    {
      path: 'detail/:id',
      name: 'OrderDetail',
      component: () => import('@/views/order/Detail.vue'),
      meta: { title: '订单详情', requiresAuth: true },
      props: true
    }
  ]
}; 