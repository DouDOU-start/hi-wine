import { createApi } from '../utils/apiFactory';
import { get, post, download } from '../utils/request';

// 创建基础API
const orderApi = createApi('/admin/orders', {
  statusUpdateUrl: '/status',
  batchUrl: '/batch'
});

// 导出基础API方法
export const {
  getList: getOrderList,
  getDetail: getOrderDetail,
  updateStatus: updateOrderStatus,
  batchUpdate: batchUpdateOrders
} = orderApi;

// 获取订单统计数据
export function getOrderStats() {
  return get('/admin/statistics/orders');
}

// 批量更新订单状态
export function batchUpdateOrderStatus(ids, status) {
  return post('/admin/orders/batch-status', { ids, status });
}

// 导出订单数据
export function exportOrders(params) {
  return download('/admin/orders/export', params, 'orders.xlsx');
}

// 获取订单趋势数据
export function getOrderTrends(params) {
  return get('/admin/statistics/order-trends', params);
} 