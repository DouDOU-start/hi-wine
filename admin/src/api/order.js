import request from '../utils/request';

// 获取订单列表
export function getOrderList(params) {
  return request({
    url: '/admin/orders',
    method: 'get',
    params
  });
}

// 获取订单详情
export function getOrderDetail(id) {
  return request({
    url: `/admin/orders/${id}`,
    method: 'get'
  });
}

// 更新订单状态
export function updateOrderStatus(id, status) {
  return request({
    url: `/admin/orders/${id}/status`,
    method: 'put',
    data: { status }
  });
}

// 获取订单统计数据
export function getOrderStats() {
  return request({
    url: '/admin/statistics/orders',
    method: 'get'
  });
}

// 导出订单数据
export function exportOrders(params) {
  return request({
    url: '/admin/orders/export',
    method: 'get',
    params,
    responseType: 'blob'
  });
} 