import { get, post, put, del } from '../utils/request';

// 获取套餐列表
export function getPackageList(params) {
  return get('/admin/packages', params);
}

// 获取套餐详情
export function getPackageDetail(id) {
  return get(`/admin/packages/${id}/with-products`);
}

// 创建套餐
export function createPackage(data) {
  return post('/admin/packages', data);
}

// 更新套餐
export function updatePackage(id, data) {
  return put(`/admin/packages/${id}`, data);
}

// 删除套餐
export function deletePackage(id) {
  return del(`/admin/packages/${id}`);
}

// 关联套餐商品
export function associatePackageProducts(packageId, productIds) {
  return post(`/admin/packages/${packageId}/products`, { productIds });
}

// 移除套餐中的单个商品
export function removeProductFromPackage(packageId, productId) {
  return del(`/admin/packages/${packageId}/products/${productId}`);
}

// 获取套餐关联的商品列表
export function getPackageProducts(packageId) {
  return get(`/admin/packages/${packageId}/products`);
}

// 获取用户套餐列表
export function getUserPackages(params) {
  return get('/admin/user-packages', params);
}

// 获取用户套餐详情
export function getUserPackageDetail(id) {
  return get(`/admin/user-packages/${id}/full-detail`);
}

// 更新用户套餐状态
export function updateUserPackageStatus(id, data) {
  return put(`/admin/user-packages/${id}/status`, data);
}

// 获取套餐统计数据
export function getPackageStats(id) {
  return get(`/admin/packages/${id}/stats`);
} 