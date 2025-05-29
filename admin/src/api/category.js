import request from '../utils/request';

// 获取分类列表
export function getCategoryList(params) {
  return request({
    url: '/api/admin/category/list',
    method: 'get',
    params
  });
}

// 添加分类
export function addCategory(data) {
  return request({
    url: '/api/admin/category/add',
    method: 'post',
    data
  });
}

// 更新分类
export function updateCategory(data) {
  return request({
    url: '/api/admin/category/update',
    method: 'post',
    data
  });
}

// 删除分类
export function deleteCategory(id) {
  return request({
    url: '/api/admin/category/delete',
    method: 'post',
    data: { id }
  });
}

// 获取分类详情
export function getCategoryDetail(id) {
  return request({
    url: '/api/admin/category/detail',
    method: 'get',
    params: { id }
  });
} 