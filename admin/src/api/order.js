import { get, post, put, download } from '../utils/request';

// 获取订单列表
export function getOrderList(params) {
  return get('/admin/orders', params);
}

// 获取订单详情
export function getOrderDetail(id) {
  return get(`/admin/orders/${id}`);
}

// 更新订单状态
export function updateOrderStatus(id, orderStatus) {
  return put(`/admin/orders/${id}/status`, { order_status: orderStatus });
}

// 获取订单统计数据
export function getOrderStats() {
  return get('/admin/statistics/orders');
}

// 导出订单数据
export function exportOrders(params) {
  return download('/admin/orders/export', params, 'orders.xlsx');
} 