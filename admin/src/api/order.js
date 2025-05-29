import request from '../utils/request';

// 获取订单列表
export function getOrderList(params) {
  return request({
    url: '/api/order/list',
    method: 'get',
    params
  });
}

// 获取订单详情
export function getOrderDetail(id) {
  return request({
    url: '/api/order/detail',
    method: 'get',
    params: { id }
  });
}

// 更新订单状态
export function updateOrderStatus(id, status) {
  return request({
    url: '/api/order/updateStatus',
    method: 'post',
    data: { id, status }
  });
}

// 获取订单统计数据
export function getOrderStats() {
  return request({
    url: '/api/order/stats',
    method: 'get'
  });
}

// 导出订单数据
export function exportOrders(params) {
  return request({
    url: '/api/order/export',
    method: 'get',
    params,
    responseType: 'blob'
  });
} 