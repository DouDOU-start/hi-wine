import request from '../utils/request';

// 获取套餐列表
export function getPackageList(params) {
  return request({
    url: '/admin/packages',
    method: 'get',
    params
  });
}

// 获取套餐详情
export function getPackageDetail(id) {
  return request({
    url: `/admin/packages/${id}/with-products`,
    method: 'get'
  });
}

// 创建套餐
export function createPackage(data) {
  return request({
    url: '/admin/packages',
    method: 'post',
    data
  });
}

// 更新套餐
export function updatePackage(id, data) {
  return request({
    url: `/admin/packages/${id}`,
    method: 'put',
    data
  });
}

// 删除套餐
export function deletePackage(id) {
  return request({
    url: `/admin/packages/${id}`,
    method: 'delete'
  });
}

// 关联套餐商品
export function associatePackageProducts(packageId, productIds) {
  return request({
    url: `/admin/packages/${packageId}/products`,
    method: 'post',
    data: { productIds }
  });
}

// 移除套餐中的单个商品
export function removeProductFromPackage(packageId, productId) {
  return request({
    url: `/admin/packages/${packageId}/products/${productId}`,
    method: 'delete'
  });
}

// 获取套餐关联的商品列表
export function getPackageProducts(packageId) {
  return request({
    url: `/admin/packages/${packageId}/products`,
    method: 'get'
  });
}

// 获取用户套餐列表
export function getUserPackageList(params) {
  return request({
    url: '/admin/user-packages',
    method: 'get',
    params
  });
}

// 获取用户套餐详情
export function getUserPackageDetail(id) {
  return request({
    url: `/admin/user-packages/${id}/full-detail`,
    method: 'get'
  });
}

// 更新用户套餐状态
export function updateUserPackageStatus(id, status) {
  return request({
    url: `/admin/user-packages/${id}/status`,
    method: 'put',
    data: { status }
  });
}

// 获取套餐统计数据
export function getPackageStats(id) {
  return request({
    url: `/admin/packages/${id}/stats`,
    method: 'get'
  });
} 