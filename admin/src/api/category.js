import request from '../utils/request';

// 获取分类列表
export function getCategoryList(params) {
  return request({
    url: '/admin/categories',
    method: 'get',
    params
  });
}

// 获取分类详情
export function getCategoryDetail(id) {
  return request({
    url: `/admin/categories/${id}`,
    method: 'get'
  });
}

// 添加分类
export function addCategory(data) {
  return request({
    url: '/admin/categories',
    method: 'post',
    data
  });
}

// 更新分类
export function updateCategory(id, data) {
  return request({
    url: `/admin/categories/${id}`,
    method: 'put',
    data
  });
}

// 删除分类
export function deleteCategory(id) {
  return request({
    url: `/admin/categories/${id}`,
    method: 'delete'
  });
}

// 更新分类状态
export function updateCategoryStatus(id, status) {
  return request({
    url: `/admin/categories/${id}`,
    method: 'put',
    data: { isActive: status }
  });
} 