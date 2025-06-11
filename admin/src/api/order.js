import { createApi } from '../utils/apiFactory';
import { get, download } from '../utils/request';

// 创建基础API
const orderApi = createApi('/admin/orders', {
  statusUpdateUrl: '/status'
});

// 导出基础API方法
export const {
  getList: getOrderList,
  getDetail: getOrderDetail,
  updateStatus: updateOrderStatus
} = orderApi;

// 获取订单统计数据
export function getOrderStats() {
  return get('/admin/statistics/orders');
}

// 导出订单数据
export function exportOrders(params) {
  return download('/admin/orders/export', params, 'orders.xlsx');
} 