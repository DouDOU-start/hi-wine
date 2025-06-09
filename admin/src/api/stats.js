import request from '../utils/request';

// 获取仪表盘概览数据
export function getDashboardStats() {
  return request({
    url: '/admin/statistics/dashboard',
    method: 'get'
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

// 获取商品销量排行
export function getProductRanking(params) {
  return request({
    url: '/admin/statistics/products/ranking',
    method: 'get',
    params
  });
} 