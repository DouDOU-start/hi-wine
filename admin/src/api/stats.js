import request from '../utils/request';

// 获取仪表盘概览数据
export function getDashboardStats() {
  return request({
    url: '/api/stats/dashboard',
    method: 'get'
  });
}

// 获取销售统计数据
export function getSalesStats(params) {
  return request({
    url: '/api/stats/sales',
    method: 'get',
    params
  });
}

// 获取商品销量排行
export function getProductRanking(params) {
  return request({
    url: '/api/stats/products/ranking',
    method: 'get',
    params
  });
}

// 获取用户消费排行
export function getUserRanking(params) {
  return request({
    url: '/api/stats/users/ranking',
    method: 'get',
    params
  });
}

// 获取分类销售统计
export function getCategorySales() {
  return request({
    url: '/api/stats/category/sales',
    method: 'get'
  });
} 