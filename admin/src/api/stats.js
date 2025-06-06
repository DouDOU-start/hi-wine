import request from '../utils/request';

// 获取仪表盘概览数据
export function getDashboardStats() {
  return request({
    url: '/admin/statistics/dashboard',
    method: 'get'
  });
}

// 获取销售统计数据
export function getSalesStats(params) {
  return request({
    url: '/admin/statistics/sales',
    method: 'get',
    params
  });
}

// 获取销售趋势数据
export function getSalesTrend(params) {
  return request({
    url: '/admin/statistics/sales/trend',
    method: 'get',
    params
  });
}

// 获取热门商品数据
export function getHotProducts(params) {
  return request({
    url: '/admin/statistics/products/hot',
    method: 'get',
    params
  });
}

// 获取套餐销售统计
export function getPackageSalesStats(params) {
  return request({
    url: '/admin/statistics/packages/sales',
    method: 'get',
    params
  });
}

// 获取套餐使用统计
export function getPackageUsageStats(params) {
  return request({
    url: '/admin/statistics/packages/usage',
    method: 'get',
    params
  });
}

// 获取商品销量排行
export function getProductRanking(params) {
  return request({
    url: '/admin/statistics/products/ranking',
    method: 'get',
    params
  });
}

// 获取用户消费排行
export function getUserRanking(params) {
  return request({
    url: '/admin/statistics/users/ranking',
    method: 'get',
    params
  });
}

// 获取分类销售统计
export function getCategorySales() {
  return request({
    url: '/admin/statistics/categories/sales',
    method: 'get'
  });
} 