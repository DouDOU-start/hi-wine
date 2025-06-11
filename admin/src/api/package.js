import { createApi } from '../utils/apiFactory';
import { get, post, put, del } from '../utils/request';

// 创建套餐基础API
const packageApi = createApi('/admin/packages');

// 创建用户套餐基础API
const userPackageApi = createApi('/admin/user-packages');

// 导出套餐基础API方法
export const {
  getList: getPackageList,
  create: createPackage,
  update: updatePackage,
  delete: deletePackage
} = packageApi;

// 导出用户套餐基础API方法
export const {
  getList: getUserPackages,
  getDetail: getUserPackageDetail
} = userPackageApi;

// 获取套餐详情（包含商品）
export function getPackageDetail(id) {
  return get(`/admin/packages/${id}/with-products`);
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

// 更新用户套餐状态
export function updateUserPackageStatus(id, data) {
  return put(`/admin/user-packages/${id}/status`, data);
}

// 获取套餐统计数据
export function getPackageStats(id) {
  return get(`/admin/packages/${id}/stats`);
} 