import request from '../utils/request';

// 获取商品列表
export function getProductList(params) {
  return request({
    url: '/admin/products',
    method: 'get',
    params
  });
}

// 获取商品详情
export function getProductDetail(id) {
  return request({
    url: `/admin/products/${id}`,
    method: 'get'
  });
}

// 添加商品
export function addProduct(data) {
  return request({
    url: '/admin/products',
    method: 'post',
    data
  });
}

// 更新商品
export function updateProduct(id, data) {
  return request({
    url: `/admin/products/${id}`,
    method: 'put',
    data
  });
}

// 删除商品
export function deleteProduct(id) {
  return request({
    url: `/admin/products/${id}`,
    method: 'delete'
  });
}

// 上下架商品
export function updateProductStatus(id, status) {
  return request({
    url: `/admin/products/${id}/status`,
    method: 'put',
    data: { status }
  });
}

// 批量更新商品状态
export function batchUpdateProductStatus(ids, status) {
  return request({
    url: '/admin/products/batch/status',
    method: 'put',
    data: { ids, status }
  });
}

// 上传商品图片
export function uploadProductImage(file) {
  const formData = new FormData();
  formData.append('file', file);
  
  return request({
    url: '/admin/upload/image',
    method: 'post',
    data: formData,
    headers: {
      'Content-Type': 'multipart/form-data'
    }
  });
} 