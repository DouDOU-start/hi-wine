import { get } from '../utils/request';

// 获取仪表盘概览数据
export function getDashboardStats() {
  return get('/admin/statistics/dashboard');
}

// 获取热门商品数据
export function getHotProducts(params) {
  return get('/admin/statistics/products/hot', params);
}

// 获取商品销量排行
export function getProductRanking(params) {
  return get('/admin/statistics/products/ranking', params);
} 